package main

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

// fakeStore is a minimal Store for unit testing execute without AWS.
type fakeStore struct {
	puts    map[string]string
	deleted []string
	putErr  map[string]error
}

func newFakeStore() *fakeStore {
	return &fakeStore{puts: make(map[string]string), putErr: make(map[string]error)}
}

func (f *fakeStore) GetAll(_ context.Context) (map[string]string, error) { return nil, nil }
func (f *fakeStore) Put(_ context.Context, key, value string) error {
	if err, ok := f.putErr[key]; ok {
		return err
	}
	f.puts[key] = value
	return nil
}

func (f *fakeStore) Delete(_ context.Context, keys []string) error {
	f.deleted = append(f.deleted, keys...)
	return nil
}

func TestExecuteDryRun(t *testing.T) {
	fs := newFakeStore()
	actions := []Action{
		{Key: "k1", Type: ActionCreate, Value: "v1"},
		{Key: "k2", Type: ActionUpdate, Value: "v2"},
		{Key: "k3", Type: ActionSkip},
	}
	var stdout, stderr bytes.Buffer
	summary := execute(context.Background(), actions, fs, true, &stdout, &stderr)

	if summary.Created != 1 || summary.Updated != 1 || summary.Unchanged != 1 || summary.Failed != 0 {
		t.Errorf("unexpected summary: %+v", summary)
	}
	if len(fs.puts) != 0 {
		t.Error("dry-run should not call Put")
	}

	out := stdout.String()
	if !strings.Contains(out, "create: k1") || !strings.Contains(out, "update: k2") {
		t.Errorf("unexpected stdout: %s", out)
	}
	if strings.Contains(out, "k3") {
		t.Error("skip should not appear in output")
	}
}

func TestExecutePartialFailure(t *testing.T) {
	fs := newFakeStore()
	fs.putErr["k2"] = context.DeadlineExceeded
	actions := []Action{
		{Key: "k1", Type: ActionCreate, Value: "v1"},
		{Key: "k2", Type: ActionCreate, Value: "v2"},
	}
	var stdout, stderr bytes.Buffer
	summary := execute(context.Background(), actions, fs, false, &stdout, &stderr)

	if summary.Created != 1 {
		t.Errorf("created = %d, want 1", summary.Created)
	}
	if summary.Failed != 1 {
		t.Errorf("failed = %d, want 1", summary.Failed)
	}
	if !strings.Contains(stderr.String(), "error: k2") {
		t.Errorf("stderr missing error for k2: %s", stderr.String())
	}
	// Values must never appear
	if strings.Contains(stdout.String(), "v1") || strings.Contains(stdout.String(), "v2") {
		t.Error("values must not appear in stdout")
	}
}

func TestPlan(t *testing.T) {
	tests := []struct {
		name     string
		entries  []Entry
		existing map[string]string
		prune    bool
		want     []Action
	}{
		{
			name:     "new key creates",
			entries:  []Entry{{Key: "/app/key1", Value: "val1"}},
			existing: map[string]string{},
			want:     []Action{{Key: "/app/key1", Type: ActionCreate, Value: "val1"}},
		},
		{
			name:     "same value skips",
			entries:  []Entry{{Key: "/app/key1", Value: "val1"}},
			existing: map[string]string{"/app/key1": "val1"},
			want:     []Action{{Key: "/app/key1", Type: ActionSkip}},
		},
		{
			name:     "changed value updates",
			entries:  []Entry{{Key: "/app/key1", Value: "newval"}},
			existing: map[string]string{"/app/key1": "oldval"},
			want:     []Action{{Key: "/app/key1", Type: ActionUpdate, Value: "newval"}},
		},
		{
			name:     "mixed create update skip",
			entries:  []Entry{{Key: "k1", Value: "v1"}, {Key: "k2", Value: "changed"}, {Key: "k3", Value: "same"}},
			existing: map[string]string{"k2": "original", "k3": "same"},
			want: []Action{
				{Key: "k1", Type: ActionCreate, Value: "v1"},
				{Key: "k2", Type: ActionUpdate, Value: "changed"},
				{Key: "k3", Type: ActionSkip},
			},
		},
		{
			name:     "prune deletes missing keys",
			entries:  []Entry{{Key: "k1", Value: "v1"}},
			existing: map[string]string{"k1": "v1", "k2": "v2", "k3": "v3"},
			prune:    true,
			want: []Action{
				{Key: "k1", Type: ActionSkip},
				{Key: "k2", Type: ActionDelete},
				{Key: "k3", Type: ActionDelete},
			},
		},
		{
			name:     "no prune keeps missing keys",
			entries:  []Entry{{Key: "k1", Value: "v1"}},
			existing: map[string]string{"k1": "v1", "k2": "v2"},
			prune:    false,
			want: []Action{
				{Key: "k1", Type: ActionSkip},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := plan(tt.entries, tt.existing, tt.prune)
			if len(got) != len(tt.want) {
				t.Fatalf("got %d actions, want %d", len(got), len(tt.want))
			}
			// Build map for lookup since delete order is not guaranteed
			wantMap := make(map[string]Action)
			for _, a := range tt.want {
				wantMap[a.Key] = a
			}
			for _, g := range got {
				w, ok := wantMap[g.Key]
				if !ok {
					t.Errorf("unexpected action for key %q", g.Key)
					continue
				}
				if g.Type != w.Type {
					t.Errorf("key %q: type = %v, want %v", g.Key, g.Type, w.Type)
				}
				if w.Type == ActionCreate || w.Type == ActionUpdate {
					if g.Value != w.Value {
						t.Errorf("key %q: value = %q, want %q", g.Key, g.Value, w.Value)
					}
				}
			}
		})
	}
}
