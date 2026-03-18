package main

import (
	"fmt"
	"log/slog"
	"regexp"
	"sort"

	"gopkg.in/yaml.v3"
)

// parseDeletePatterns parses a YAML list of regex pattern strings and compiles them.
func parseDeletePatterns(data []byte) ([]*regexp.Regexp, error) {
	var patterns []string
	if err := yaml.Unmarshal(data, &patterns); err != nil {
		return nil, fmt.Errorf("delete file parse error: %w", err)
	}

	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("invalid regex pattern %q: %w", p, err)
		}
		slog.Debug("compiled delete pattern", "pattern", p)
		compiled = append(compiled, re)
	}
	return compiled, nil
}

// planDeletes classifies existing AWS keys into delete candidates, conflicts, and unmanaged keys.
// - deletes: keys matching any pattern AND not in yamlKeys
// - conflicts: keys matching any pattern AND in yamlKeys (should abort)
// - unmanaged: keys not in yamlKeys AND not matching any pattern (warning)
func planDeletes(existing map[string]string, yamlKeys map[string]bool, patterns []*regexp.Regexp) (deletes []Action, conflicts []string, unmanaged []string) {
	for k := range existing {
		if yamlKeys[k] {
			// Check if this YAML key matches any delete pattern (conflict)
			for _, re := range patterns {
				if re.MatchString(k) {
					conflicts = append(conflicts, k)
					break
				}
			}
			continue
		}

		// Key not in YAML — check patterns
		matched := false
		for _, re := range patterns {
			if re.MatchString(k) {
				slog.Debug("regex match", "pattern", re.String(), "key", k, "matched", true)
				matched = true
				break
			}
		}

		if matched {
			deletes = append(deletes, Action{Key: k, Type: ActionDelete})
		} else {
			unmanaged = append(unmanaged, k)
		}
	}

	sort.Slice(deletes, func(i, j int) bool { return deletes[i].Key < deletes[j].Key })
	sort.Strings(conflicts)
	sort.Strings(unmanaged)

	return deletes, conflicts, unmanaged
}
