package main

import (
	"testing"
)

func TestParseYAML(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []Entry
		wantErr string
	}{
		{
			name:  "valid key-values",
			input: "/db/url: postgres://localhost\n/db/port: \"5432\"\n",
			want: []Entry{
				{Key: "/db/url", Value: "postgres://localhost"},
				{Key: "/db/port", Value: "5432"},
			},
		},
		{
			name:  "sops key filtered",
			input: "/db/url: myvalue\nsops:\n  version: \"3.7\"\n",
			want: []Entry{
				{Key: "/db/url", Value: "myvalue"},
			},
		},
		{
			name:  "integer value converted to string",
			input: "port: 5432\n",
			want:  []Entry{{Key: "port", Value: "5432"}},
		},
		{
			name:  "boolean value converted to string",
			input: "debug: true\n",
			want:  []Entry{{Key: "debug", Value: "true"}},
		},
		{
			name:  "float value converted to string",
			input: "ratio: 3.14\n",
			want:  []Entry{{Key: "ratio", Value: "3.14"}},
		},
		{
			name:  "empty string value is valid",
			input: "key: \"\"\n",
			want:  []Entry{{Key: "key", Value: ""}},
		},
		{
			name:    "duplicate key error",
			input:   "key1: val1\nkey1: val2\n",
			wantErr: "key1",
		},
		{
			name:    "null value error",
			input:   "key1: null\n",
			wantErr: "key1",
		},
		{
			name:    "map value error",
			input:   "key1:\n  nested: value\n",
			wantErr: "key1",
		},
		{
			name:    "array value error",
			input:   "key1:\n  - item1\n  - item2\n",
			wantErr: "key1",
		},
		{
			name:    "empty key error",
			input:   "\"\": value\n",
			wantErr: "empty",
		},
		{
			name:    "zero keys after sops filter",
			input:   "sops:\n  version: \"3.7\"\n",
			wantErr: "no keys",
		},
		{
			name:    "completely empty yaml",
			input:   "",
			wantErr: "no keys",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseYAML([]byte(tt.input))
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.wantErr)
				}
				if !containsStr(err.Error(), tt.wantErr) {
					t.Fatalf("expected error containing %q, got %q", tt.wantErr, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("got %d entries, want %d", len(got), len(tt.want))
			}
			for i, e := range got {
				if e.Key != tt.want[i].Key || e.Value != tt.want[i].Value {
					t.Errorf("entry[%d] = {%q, %q}, want {%q, %q}", i, e.Key, e.Value, tt.want[i].Key, tt.want[i].Value)
				}
			}
		})
	}
}

func TestWriteYAML(t *testing.T) {
	entries := []Entry{
		{Key: "/app/db_url", Value: "postgres://localhost"},
		{Key: "/app/api_key", Value: "sk-123"},
	}
	data, err := writeYAML(entries)
	if err != nil {
		t.Fatalf("writeYAML failed: %v", err)
	}

	// Parse it back to verify round-trip
	parsed, err := parseYAML(data)
	if err != nil {
		t.Fatalf("parseYAML of written data failed: %v", err)
	}
	if len(parsed) != len(entries) {
		t.Fatalf("got %d entries, want %d", len(parsed), len(entries))
	}
	for i, e := range parsed {
		if e.Key != entries[i].Key || e.Value != entries[i].Value {
			t.Errorf("entry[%d] = {%q, %q}, want {%q, %q}", i, e.Key, e.Value, entries[i].Key, entries[i].Value)
		}
	}
}

func TestWriteYAMLEmpty(t *testing.T) {
	data, err := writeYAML([]Entry{})
	if err != nil {
		t.Fatalf("writeYAML failed: %v", err)
	}
	// Should produce valid YAML (empty mapping)
	if len(data) == 0 {
		t.Error("expected non-empty output")
	}
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && searchStr(s, substr)
}

func searchStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
