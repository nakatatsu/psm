package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
)

var version = "dev"

func main() {
	cfg, err := parseArgs(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if cfg.ShowVersion {
		fmt.Printf("psm version %s\n", version)
		os.Exit(0)
	}

	level := slog.LevelInfo
	if cfg.Debug {
		level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	})))

	exitCode, err := run(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(exitCode)
}

func run(cfg Config) (int, error) {
	ctx := context.Background()

	// FR-002: always ignore AWS_PROFILE env var
	if err := os.Unsetenv("AWS_PROFILE"); err != nil {
		return 1, fmt.Errorf("failed to unset AWS_PROFILE: %w", err)
	}

	var awsOpts []func(*awsconfig.LoadOptions) error
	if cfg.Profile != "" {
		awsOpts = append(awsOpts, awsconfig.WithSharedConfigProfile(cfg.Profile))
	}
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, awsOpts...)
	if err != nil {
		return 1, fmt.Errorf("failed to load AWS config: %w", err)
	}

	var store Store
	switch cfg.Store {
	case "ssm":
		store = NewSSMStore(awsCfg)
	}

	switch cfg.Subcommand {
	case "sync":
		fi, _ := os.Stdin.Stat()
		isTerminal := fi != nil && fi.Mode()&os.ModeCharDevice != 0
		streams := &IOStreams{
			Stdin:      os.Stdin,
			Stdout:     os.Stdout,
			Stderr:     os.Stderr,
			IsTerminal: func() bool { return isTerminal },
		}
		return runSync(ctx, cfg, store, streams)
	case "export":
		return runExport(ctx, cfg, store)
	}

	return 1, fmt.Errorf("unknown subcommand: %s", cfg.Subcommand)
}

func runSync(ctx context.Context, cfg Config, store Store, io *IOStreams) (int, error) {
	data, err := os.ReadFile(cfg.File)
	if err != nil {
		return 1, fmt.Errorf("failed to read file %s: %w", cfg.File, err)
	}

	entries, err := parseYAML(data)
	if err != nil {
		return 1, err
	}

	existing, err := store.GetAll(ctx)
	if err != nil {
		return 1, fmt.Errorf("failed to get existing parameters: %w", err)
	}

	actions := plan(entries, existing)

	// Delete flow
	if cfg.DeleteFile != "" {
		deleteData, err := os.ReadFile(cfg.DeleteFile)
		if err != nil {
			return 1, fmt.Errorf("failed to read delete file %s: %w", cfg.DeleteFile, err)
		}

		patterns, err := parseDeletePatterns(deleteData)
		if err != nil {
			return 1, err
		}

		yamlKeys := make(map[string]bool)
		for _, e := range entries {
			yamlKeys[e.Key] = true
		}

		deletes, conflicts, unmanaged := planDeletes(existing, yamlKeys, patterns)

		// Conflict detection — abort before any execution
		if len(conflicts) > 0 {
			for _, k := range conflicts {
				slog.Error("conflict: deletion candidate exists in sync YAML", "key", k)
			}
			return 1, fmt.Errorf("aborting — %d conflict(s) detected, no changes made", len(conflicts))
		}

		// Warn about unmanaged keys
		for _, k := range unmanaged {
			slog.Warn("unmanaged key detected", "key", k)
		}

		actions = append(actions, deletes...)
	}

	// Check if there are any changes
	hasChanges := false
	for _, a := range actions {
		if a.Type != ActionSkip {
			hasChanges = true
			break
		}
	}

	if cfg.DryRun {
		displayPlan(actions, io.Stdout)
		printSummary(actions, true, io.Stdout)
		return 0, nil
	}

	if !hasChanges {
		printSummary(actions, false, io.Stdout)
		return 0, nil
	}

	displayPlan(actions, io.Stdout)

	// Approval flow
	if !cfg.SkipApprove {
		if !io.IsTerminal() {
			return 0, nil
		}
		if !promptApprove(io.Stdin, io.Stderr) {
			return 0, nil
		}
	}

	summary := execute(ctx, actions, store, io.Stdout, io.Stderr)

	if summary.Failed > 0 {
		return 1, nil
	}
	return 0, nil
}

func parseArgs(args []string) (Config, error) {
	if len(args) < 2 {
		return Config{}, fmt.Errorf("usage: psm <sync|export> [flags] <file>")
	}

	if args[1] == "--version" {
		return Config{ShowVersion: true}, nil
	}

	sub := args[1]
	if sub != "sync" && sub != "export" {
		return Config{}, fmt.Errorf("unknown subcommand %q. usage: psm <sync|export> [flags] <file>", sub)
	}

	fs := flag.NewFlagSet(sub, flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	store := fs.String("store", "", "store type: ssm")
	profile := fs.String("profile", "", "AWS profile name")

	var dryRun, skipApprove, debug bool
	var deleteFile string
	if sub == "sync" {
		fs.BoolVar(&dryRun, "dry-run", false, "show plan without executing")
		fs.BoolVar(&skipApprove, "skip-approve", false, "skip approval prompt")
		fs.StringVar(&deleteFile, "delete", "", "YAML file with delete regex patterns")
	}
	fs.BoolVar(&debug, "debug", false, "enable debug logging")

	// Detect removed --prune flag before parsing
	for _, arg := range args[2:] {
		if arg == "--prune" {
			return Config{}, fmt.Errorf("--prune has been removed. Use --delete <file> with regex patterns instead")
		}
	}

	if err := fs.Parse(args[2:]); err != nil {
		return Config{}, fmt.Errorf("usage: psm %s [flags] <file>", sub)
	}

	if *store == "" {
		return Config{}, fmt.Errorf("--store is required. usage: psm %s --store ssm [flags] <file>", sub)
	}
	if *store != "ssm" {
		return Config{}, fmt.Errorf("invalid --store value %q: must be ssm", *store)
	}

	remaining := fs.Args()
	if len(remaining) != 1 {
		return Config{}, fmt.Errorf("exactly one file argument required. usage: psm %s --store ssm [flags] <file>", sub)
	}

	return Config{
		Subcommand:  sub,
		Store:       *store,
		Profile:     *profile,
		DryRun:      dryRun,
		SkipApprove: skipApprove,
		Debug:       debug,
		DeleteFile:  deleteFile,
		File:        remaining[0],
	}, nil
}
