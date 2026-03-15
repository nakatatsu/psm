package main

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// parseYAML parses YAML bytes into entries, filtering sops key and validating per FR-020.
func parseYAML(data []byte) ([]Entry, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("yaml parse error: %w", err)
	}

	// Empty document
	if doc.Kind == 0 || len(doc.Content) == 0 {
		return nil, fmt.Errorf("no keys found after filtering")
	}

	root := doc.Content[0]
	if root.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("expected YAML mapping at top level")
	}

	// First pass: filter sops key (FR-006), collect key-value nodes
	type kvPair struct {
		keyNode *yaml.Node
		valNode *yaml.Node
	}
	var pairs []kvPair
	for i := 0; i+1 < len(root.Content); i += 2 {
		keyNode := root.Content[i]
		valNode := root.Content[i+1]
		if keyNode.Value == "sops" {
			continue
		}
		pairs = append(pairs, kvPair{keyNode, valNode})
	}

	// Validate (FR-020)
	if len(pairs) == 0 {
		return nil, fmt.Errorf("no keys found after filtering")
	}

	var errors []string
	seen := make(map[string]bool)
	for _, p := range pairs {
		key := p.keyNode.Value

		if key == "" {
			errors = append(errors, "empty key found")
			continue
		}

		if seen[key] {
			errors = append(errors, fmt.Sprintf("duplicate key: %s", key))
			continue
		}
		seen[key] = true

		switch p.valNode.Kind {
		case yaml.MappingNode:
			errors = append(errors, fmt.Sprintf("value for key %q is a map (must be scalar)", key))
		case yaml.SequenceNode:
			errors = append(errors, fmt.Sprintf("value for key %q is an array (must be scalar)", key))
		case yaml.ScalarNode:
			if p.valNode.Tag == "!!null" {
				errors = append(errors, fmt.Sprintf("value for key %q is null", key))
			}
		}
	}

	if len(errors) > 0 {
		return nil, fmt.Errorf("validation errors:\n  %s", strings.Join(errors, "\n  "))
	}

	// Build entries
	entries := make([]Entry, 0, len(pairs))
	for _, p := range pairs {
		entries = append(entries, Entry{
			Key:   p.keyNode.Value,
			Value: p.valNode.Value,
		})
	}

	return entries, nil
}

// writeYAML converts entries to YAML bytes.
func writeYAML(entries []Entry) ([]byte, error) {
	root := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	for _, e := range entries {
		root.Content = append(root.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: e.Key},
			&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: e.Value, Style: yaml.DoubleQuotedStyle},
		)
	}
	doc := &yaml.Node{Kind: yaml.DocumentNode, Content: []*yaml.Node{root}}
	out, err := yaml.Marshal(doc)
	if err != nil {
		return nil, fmt.Errorf("yaml marshal error: %w", err)
	}
	return out, nil
}
