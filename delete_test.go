package main

import (
	"regexp"
	"testing"
)

func TestParseDeletePatterns(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{
			name:  "valid patterns",
			input: "- \"^/myapp/legacy/\"\n- \"^/myapp/deprecated-.*\"\n",
			want:  2,
		},
		{
			name:  "empty list",
			input: "[]\n",
			want:  0,
		},
		{
			name:    "invalid regex",
			input:   "- \"[invalid\"\n",
			wantErr: true,
		},
		{
			name:    "non-YAML",
			input:   "not: a: list: {{{\n",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDeletePatterns([]byte(tt.input))
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != tt.want {
				t.Errorf("got %d patterns, want %d", len(got), tt.want)
			}
		})
	}
}

func TestPlanDeletes(t *testing.T) {
	compile := func(patterns ...string) []*regexp.Regexp {
		var res []*regexp.Regexp
		for _, p := range patterns {
			res = append(res, regexp.MustCompile(p))
		}
		return res
	}

	tests := []struct {
		name          string
		existing      map[string]string
		yamlKeys      map[string]bool
		patterns      []*regexp.Regexp
		wantDeletes   int
		wantConflicts int
		wantUnmanaged int
	}{
		{
			name:          "key matches pattern and not in YAML → delete",
			existing:      map[string]string{"/myapp/legacy/old": "v"},
			yamlKeys:      map[string]bool{},
			patterns:      compile("^/myapp/legacy/"),
			wantDeletes:   1,
			wantConflicts: 0,
			wantUnmanaged: 0,
		},
		{
			name:          "key matches pattern and in YAML → conflict",
			existing:      map[string]string{"/myapp/legacy/active": "v"},
			yamlKeys:      map[string]bool{"/myapp/legacy/active": true},
			patterns:      compile("^/myapp/legacy/"),
			wantDeletes:   0,
			wantConflicts: 1,
			wantUnmanaged: 0,
		},
		{
			name:          "key matches no pattern and not in YAML → unmanaged",
			existing:      map[string]string{"/other/team/key": "v"},
			yamlKeys:      map[string]bool{},
			patterns:      compile("^/myapp/"),
			wantDeletes:   0,
			wantConflicts: 0,
			wantUnmanaged: 1,
		},
		{
			name:          "key in YAML and no pattern match → not affected",
			existing:      map[string]string{"/myapp/prod/key": "v"},
			yamlKeys:      map[string]bool{"/myapp/prod/key": true},
			patterns:      compile("^/myapp/legacy/"),
			wantDeletes:   0,
			wantConflicts: 0,
			wantUnmanaged: 0,
		},
		{
			name:          "multiple patterns match same key → single delete",
			existing:      map[string]string{"/myapp/legacy/old": "v"},
			yamlKeys:      map[string]bool{},
			patterns:      compile("^/myapp/", "legacy"),
			wantDeletes:   1,
			wantConflicts: 0,
			wantUnmanaged: 0,
		},
		{
			name:          "no AWS keys match any pattern → empty result",
			existing:      map[string]string{"/other/key": "v"},
			yamlKeys:      map[string]bool{},
			patterns:      compile("^/myapp/legacy/"),
			wantDeletes:   0,
			wantConflicts: 0,
			wantUnmanaged: 1,
		},
		{
			name: "mixed scenario",
			existing: map[string]string{
				"/myapp/legacy/old":  "v1",
				"/myapp/prod/key":    "v2",
				"/myapp/legacy/keep": "v3",
				"/other/team/key":    "v4",
			},
			yamlKeys: map[string]bool{
				"/myapp/prod/key":    true,
				"/myapp/legacy/keep": true,
			},
			patterns:      compile("^/myapp/legacy/"),
			wantDeletes:   1, // /myapp/legacy/old
			wantConflicts: 1, // /myapp/legacy/keep (in YAML + matches pattern)
			wantUnmanaged: 1, // /other/team/key
		},
		{
			name:          "mass deletion — all keys match",
			existing:      map[string]string{"/a": "1", "/b": "2", "/c": "3"},
			yamlKeys:      map[string]bool{},
			patterns:      compile(".*"),
			wantDeletes:   3,
			wantConflicts: 0,
			wantUnmanaged: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deletes, conflicts, unmanaged := planDeletes(tt.existing, tt.yamlKeys, tt.patterns)
			if len(deletes) != tt.wantDeletes {
				t.Errorf("deletes = %d, want %d", len(deletes), tt.wantDeletes)
			}
			if len(conflicts) != tt.wantConflicts {
				t.Errorf("conflicts = %d, want %d", len(conflicts), tt.wantConflicts)
			}
			if len(unmanaged) != tt.wantUnmanaged {
				t.Errorf("unmanaged = %d, want %d", len(unmanaged), tt.wantUnmanaged)
			}
		})
	}
}
