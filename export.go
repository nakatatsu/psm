package main

import (
	"context"
	"fmt"
	"os"
	"sort"
)

func runExport(ctx context.Context, cfg Config, store Store) (int, error) {
	// FR-022: check output file does not exist
	if _, err := os.Stat(cfg.File); err == nil {
		return 1, fmt.Errorf("output file %s already exists", cfg.File)
	}

	// Get all parameters/secrets
	all, err := store.GetAll(ctx)
	if err != nil {
		return 1, fmt.Errorf("failed to get parameters: %w", err)
	}

	// FR-023: error if zero items
	if len(all) == 0 {
		return 1, fmt.Errorf("no parameters/secrets found in store")
	}

	// Convert to sorted entries
	entries := make([]Entry, 0, len(all))
	for k, v := range all {
		entries = append(entries, Entry{Key: k, Value: v})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	// Write YAML
	data, err := writeYAML(entries)
	if err != nil {
		return 1, err
	}

	if err := os.WriteFile(cfg.File, data, 0o600); err != nil {
		return 1, fmt.Errorf("failed to write file %s: %w", cfg.File, err)
	}

	return 0, nil
}
