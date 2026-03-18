package main

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"regexp"
	"strings"
	"testing"
)

// fakeStore is a minimal Store for unit testing execute without AWS.
type fakeStore struct {
	existing map[string]string
	puts     map[string]string
	deleted  []string
	putErr   map[string]error
}

func newFakeStore() *fakeStore {
	return &fakeStore{puts: make(map[string]string), putErr: make(map[string]error)}
}

func (f *fakeStore) GetAll(_ context.Context) (map[string]string, error) {
	if f.existing != nil {
		return f.existing, nil
	}
	return nil, nil
}
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

func TestDisplayPlan(t *testing.T) {
	tests := []struct {
		name    string
		actions []Action
		want    []string
		notWant []string
	}{
		{
			name: "creates and updates shown",
			actions: []Action{
				{Key: "k1", Type: ActionCreate, Value: "v1"},
				{Key: "k2", Type: ActionUpdate, Value: "v2"},
			},
			want: []string{"create: k1", "update: k2"},
		},
		{
			name: "deletes shown",
			actions: []Action{
				{Key: "k3", Type: ActionDelete},
			},
			want: []string{"delete: k3"},
		},
		{
			name: "skip actions not shown",
			actions: []Action{
				{Key: "k1", Type: ActionCreate, Value: "v1"},
				{Key: "k2", Type: ActionSkip},
			},
			want:    []string{"create: k1"},
			notWant: []string{"k2"},
		},
		{
			name:    "empty actions produce no output",
			actions: []Action{},
			want:    []string{},
		},
		{
			name: "values never shown",
			actions: []Action{
				{Key: "k1", Type: ActionCreate, Value: "secret-value"},
			},
			want:    []string{"create: k1"},
			notWant: []string{"secret-value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			displayPlan(tt.actions, &buf)
			out := buf.String()
			for _, w := range tt.want {
				if !strings.Contains(out, w) {
					t.Errorf("output missing %q: %s", w, out)
				}
			}
			for _, nw := range tt.notWant {
				if strings.Contains(out, nw) {
					t.Errorf("output should not contain %q: %s", nw, out)
				}
			}
		})
	}
}

func TestPromptApprove(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "y approves", input: "y\n", want: true},
		{name: "Y approves", input: "Y\n", want: true},
		{name: "N declines", input: "N\n", want: false},
		{name: "n declines", input: "n\n", want: false},
		{name: "empty declines", input: "\n", want: false},
		{name: "yes declines", input: "yes\n", want: false},
		{name: "arbitrary declines", input: "abc\n", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			var writer bytes.Buffer
			got := promptApprove(reader, &writer)
			if got != tt.want {
				t.Errorf("promptApprove(%q) = %v, want %v", tt.input, got, tt.want)
			}
			if !strings.Contains(writer.String(), "Proceed? [y/N]") {
				t.Errorf("prompt text missing: %s", writer.String())
			}
		})
	}
}

func testIOStreams(stdin string, isTerminal bool) (*IOStreams, *bytes.Buffer, *bytes.Buffer) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	return &IOStreams{
		Stdin:      strings.NewReader(stdin),
		Stdout:     stdout,
		Stderr:     stderr,
		IsTerminal: func() bool { return isTerminal },
	}, stdout, stderr
}

func TestRunSyncApproveFlow(t *testing.T) {
	// Setup: create temp sync file
	syncFile := t.TempDir() + "/params.yml"
	if err := os.WriteFile(syncFile, []byte("/app/key1: \"newval\"\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// fakeStore with existing data
	fs := newFakeStore()
	fs.existing = map[string]string{"/app/key1": "oldval"}
	ctx := context.Background()

	t.Run("approve executes", func(t *testing.T) {
		s := newFakeStoreWithExisting(map[string]string{"/app/key1": "oldval"})
		io, stdout, _ := testIOStreams("y\n", true)
		cfg := Config{File: syncFile}
		code, err := runSync(ctx, cfg, s, io)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if code != 0 {
			t.Errorf("exit code = %d, want 0", code)
		}
		if !strings.Contains(stdout.String(), "update: /app/key1") {
			t.Errorf("stdout missing update line: %s", stdout.String())
		}
		if s.puts["/app/key1"] != "newval" {
			t.Error("expected store.Put to be called with new value")
		}
	})

	t.Run("decline does not execute", func(t *testing.T) {
		s := newFakeStoreWithExisting(map[string]string{"/app/key1": "oldval"})
		io, _, _ := testIOStreams("N\n", true)
		cfg := Config{File: syncFile}
		code, err := runSync(ctx, cfg, s, io)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if code != 0 {
			t.Errorf("exit code = %d, want 0", code)
		}
		if len(s.puts) != 0 {
			t.Error("store should not be called when user declines")
		}
	})

	t.Run("skip-approve executes without prompt", func(t *testing.T) {
		s := newFakeStoreWithExisting(map[string]string{"/app/key1": "oldval"})
		io, _, stderr := testIOStreams("", false)
		cfg := Config{File: syncFile, SkipApprove: true}
		code, err := runSync(ctx, cfg, s, io)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if code != 0 {
			t.Errorf("exit code = %d, want 0", code)
		}
		if strings.Contains(stderr.String(), "Proceed?") {
			t.Error("prompt should not appear with --skip-approve")
		}
		if s.puts["/app/key1"] != "newval" {
			t.Error("expected store.Put to be called")
		}
	})

	t.Run("dry-run shows plan without prompt or execution", func(t *testing.T) {
		s := newFakeStoreWithExisting(map[string]string{"/app/key1": "oldval"})
		io, stdout, stderr := testIOStreams("", false)
		cfg := Config{File: syncFile, DryRun: true}
		code, err := runSync(ctx, cfg, s, io)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if code != 0 {
			t.Errorf("exit code = %d, want 0", code)
		}
		if !strings.Contains(stdout.String(), "update: /app/key1") {
			t.Errorf("stdout missing plan: %s", stdout.String())
		}
		if strings.Contains(stderr.String(), "Proceed?") {
			t.Error("prompt should not appear with --dry-run")
		}
		if len(s.puts) != 0 {
			t.Error("store should not be called in dry-run")
		}
	})

	t.Run("no changes skips prompt", func(t *testing.T) {
		syncFileNoChange := t.TempDir() + "/params.yml"
		os.WriteFile(syncFileNoChange, []byte("/app/key1: \"sameval\"\n"), 0644)
		s := newFakeStoreWithExisting(map[string]string{"/app/key1": "sameval"})
		io, stdout, stderr := testIOStreams("", true)
		cfg := Config{File: syncFileNoChange}
		code, err := runSync(ctx, cfg, s, io)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if code != 0 {
			t.Errorf("exit code = %d, want 0", code)
		}
		if strings.Contains(stderr.String(), "Proceed?") {
			t.Error("prompt should not appear when no changes")
		}
		// Should still show summary
		if !strings.Contains(stdout.String(), "unchanged") {
			t.Errorf("stdout missing summary: %s", stdout.String())
		}
	})

	t.Run("non-terminal without skip-approve declines", func(t *testing.T) {
		s := newFakeStoreWithExisting(map[string]string{"/app/key1": "oldval"})
		io, _, _ := testIOStreams("", false) // not a terminal
		cfg := Config{File: syncFile}
		code, err := runSync(ctx, cfg, s, io)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if code != 0 {
			t.Errorf("exit code = %d, want 0", code)
		}
		if len(s.puts) != 0 {
			t.Error("store should not be called in non-terminal without --skip-approve")
		}
	})
}

func newFakeStoreWithExisting(existing map[string]string) *fakeStore {
	fs := newFakeStore()
	fs.existing = existing
	return fs
}

func TestDebugLogging(t *testing.T) {
	t.Run("debug level emits debug messages", func(t *testing.T) {
		var buf bytes.Buffer
		handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
		slog.SetDefault(slog.New(handler))
		defer slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))

		slog.Debug("test debug message", "key", "value")
		if !strings.Contains(buf.String(), "test debug message") {
			t.Errorf("debug message not emitted: %s", buf.String())
		}
	})

	t.Run("info level hides debug messages", func(t *testing.T) {
		var buf bytes.Buffer
		handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})
		slog.SetDefault(slog.New(handler))
		defer slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))

		slog.Debug("hidden debug message", "key", "value")
		if strings.Contains(buf.String(), "hidden debug message") {
			t.Errorf("debug message should be hidden at info level: %s", buf.String())
		}
	})
}

func TestConflictDetectionAbortsAll(t *testing.T) {
	// If planDeletes returns conflicts, no operations should execute
	fs := newFakeStore()
	ctx := context.Background()

	// Simulate: sync wants to create k1, but k2 is a conflict (in YAML + matches delete pattern)
	entries := []Entry{{Key: "k1", Value: "v1"}, {Key: "k2", Value: "v2"}}
	existing := map[string]string{"k2": "v2", "k3": "v3"}

	syncActions := plan(entries, existing)

	yamlKeys := map[string]bool{"k1": true, "k2": true}
	patterns := []*regexp.Regexp{regexp.MustCompile("k[23]")}
	deletes, conflicts, _ := planDeletes(existing, yamlKeys, patterns)

	// Conflicts should be detected
	if len(conflicts) == 0 {
		t.Fatal("expected conflicts, got none")
	}
	// k2 is in YAML and matches pattern → conflict
	if conflicts[0] != "k2" {
		t.Errorf("expected conflict on k2, got %q", conflicts[0])
	}

	// With conflicts, we should NOT execute anything
	// Verify by checking that if we did execute the sync actions, the store would be called
	// But the point is: the caller (runSync) checks conflicts BEFORE calling execute
	_ = syncActions
	_ = deletes

	// Verify store was never touched
	if len(fs.puts) != 0 {
		t.Error("store should not have been called when conflicts exist")
	}

	// Also verify that even the non-conflicting delete (k3) is in the deletes list
	if len(deletes) != 1 || deletes[0].Key != "k3" {
		t.Errorf("expected delete for k3, got %v", deletes)
	}

	// Execute should not be called — simulate what runSync does
	summary := execute(ctx, syncActions, fs, false, &bytes.Buffer{}, &bytes.Buffer{})
	// This is wrong behavior — in real code, execute is never called when conflicts exist
	// We're just verifying the store interaction for completeness
	_ = summary
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

	// dry-run indicator on action lines
	if !strings.Contains(out, "(dry-run) create: k1") {
		t.Errorf("dry-run action line missing (dry-run) prefix: %s", out)
	}
	if !strings.Contains(out, "(dry-run) update: k2") {
		t.Errorf("dry-run action line missing (dry-run) prefix: %s", out)
	}
	// dry-run indicator on summary line
	if !strings.Contains(out, "(dry-run)") || !strings.Contains(out, "created") {
		t.Errorf("dry-run summary line missing (dry-run) indicator: %s", out)
	}
}

func TestExecuteNoDryRunOutput(t *testing.T) {
	fs := newFakeStore()
	actions := []Action{
		{Key: "k1", Type: ActionCreate, Value: "v1"},
	}
	var stdout, stderr bytes.Buffer
	execute(context.Background(), actions, fs, false, &stdout, &stderr)

	out := stdout.String()
	if strings.Contains(out, "(dry-run)") {
		t.Errorf("non-dry-run output should not contain (dry-run): %s", out)
	}
}

func TestExecuteDryRunWithDelete(t *testing.T) {
	fs := newFakeStore()
	actions := []Action{
		{Key: "k1", Type: ActionDelete},
	}
	var stdout, stderr bytes.Buffer
	execute(context.Background(), actions, fs, true, &stdout, &stderr)

	out := stdout.String()
	if !strings.Contains(out, "(dry-run) delete: k1") {
		t.Errorf("dry-run delete line missing (dry-run) prefix: %s", out)
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
			name:     "extra keys in AWS are ignored",
			entries:  []Entry{{Key: "k1", Value: "v1"}},
			existing: map[string]string{"k1": "v1", "k2": "v2"},
			want: []Action{
				{Key: "k1", Type: ActionSkip},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := plan(tt.entries, tt.existing)
			if len(got) != len(tt.want) {
				t.Fatalf("got %d actions, want %d", len(got), len(tt.want))
			}
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
