package main

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

func TestExportFileExists(t *testing.T) {
	// Create a temp file
	f, err := os.CreateTemp("", "psm-export-test-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(f.Name()) }()

	cfg := Config{File: f.Name()}
	code, err := runExport(context.Background(), cfg, newFakeStore())
	if code != 1 || err == nil {
		t.Errorf("expected error for existing file, got code=%d err=%v", code, err)
	}
}

func TestExportEmptyStore(t *testing.T) {
	tmpFile := "/tmp/psm-export-empty-test.yaml"
	defer func() { _ = os.Remove(tmpFile) }()

	cfg := Config{File: tmpFile}
	code, err := runExport(context.Background(), cfg, &emptyStore{})
	if code != 1 || err == nil {
		t.Errorf("expected error for empty store, got code=%d err=%v", code, err)
	}
}

type emptyStore struct{}

func (e *emptyStore) GetAll(_ context.Context) (map[string]string, error) {
	return map[string]string{}, nil
}
func (e *emptyStore) Put(_ context.Context, _, _ string) error   { return nil }
func (e *emptyStore) Delete(_ context.Context, _ []string) error { return nil }

func TestExportRoundTrip(t *testing.T) {
	skipIfNoAWS(t)
	cfg := loadAWSConfig(t)
	client := ssm.NewFromConfig(cfg)
	store := NewSSMStore(cfg)
	ctx := context.Background()

	keys := []string{
		ssmTestPrefix + "export-k1",
		ssmTestPrefix + "export-k2",
	}
	data := map[string]string{
		keys[0]: "val1",
		keys[1]: "val2",
	}

	for k, v := range data {
		_, err := client.PutParameter(ctx, &ssm.PutParameterInput{
			Name:      aws.String(k),
			Value:     aws.String(v),
			Type:      ssmtypes.ParameterTypeSecureString,
			Overwrite: aws.Bool(true),
		})
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}
	}
	defer cleanupSSMTestData(t, client, keys)

	// Export
	tmpFile := "/tmp/psm-export-roundtrip.yaml"
	defer func() { _ = os.Remove(tmpFile) }()

	exportCfg := Config{File: tmpFile}
	code, err := runExport(ctx, exportCfg, store)
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}
	if code != 0 {
		t.Fatalf("export exit code = %d, want 0", code)
	}

	// Read back and sync — should be 0 changes (SC-006)
	yamlData, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatal(err)
	}
	entries, err := parseYAML(yamlData)
	if err != nil {
		t.Fatalf("parseYAML failed: %v", err)
	}

	existing, _ := store.GetAll(ctx)
	actions := plan(entries, existing)

	for _, a := range actions {
		if a.Type != ActionSkip {
			t.Errorf("round-trip: key %q has action %v, want skip", a.Key, a.Type)
		}
	}
}
