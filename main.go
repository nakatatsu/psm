package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
)

func main() {
	cfg, err := parseArgs(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

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
	case "sm":
		store = NewSMStore(awsCfg)
	}

	switch cfg.Subcommand {
	case "sync":
		return runSync(ctx, cfg, store)
	case "export":
		return runExport(ctx, cfg, store)
	}

	return 1, fmt.Errorf("unknown subcommand: %s", cfg.Subcommand)
}

func runSync(ctx context.Context, cfg Config, store Store) (int, error) {
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

	actions := plan(entries, existing, cfg.Prune)
	summary := execute(ctx, actions, store, cfg.DryRun, os.Stdout, os.Stderr)

	if summary.Failed > 0 {
		return 1, nil
	}
	return 0, nil
}

func parseArgs(args []string) (Config, error) {
	if len(args) < 2 {
		return Config{}, fmt.Errorf("usage: psm <sync|export> [flags] <file>")
	}

	sub := args[1]
	if sub != "sync" && sub != "export" {
		return Config{}, fmt.Errorf("unknown subcommand %q. usage: psm <sync|export> [flags] <file>", sub)
	}

	fs := flag.NewFlagSet(sub, flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	store := fs.String("store", "", "store type: ssm or sm")
	profile := fs.String("profile", "", "AWS profile name")

	var prune, dryRun bool
	if sub == "sync" {
		fs.BoolVar(&prune, "prune", false, "delete keys not in YAML")
		fs.BoolVar(&dryRun, "dry-run", false, "show plan without executing")
	}

	if err := fs.Parse(args[2:]); err != nil {
		return Config{}, fmt.Errorf("usage: psm %s [flags] <file>", sub)
	}

	if *store == "" {
		return Config{}, fmt.Errorf("--store is required. usage: psm %s --store <ssm|sm> [flags] <file>", sub)
	}
	if *store != "ssm" && *store != "sm" {
		return Config{}, fmt.Errorf("invalid --store value %q: must be ssm or sm", *store)
	}

	remaining := fs.Args()
	if len(remaining) != 1 {
		return Config{}, fmt.Errorf("exactly one file argument required. usage: psm %s --store <ssm|sm> [flags] <file>", sub)
	}

	return Config{
		Subcommand: sub,
		Store:      *store,
		Profile:    *profile,
		Prune:      prune,
		DryRun:     dryRun,
		File:       remaining[0],
	}, nil
}
