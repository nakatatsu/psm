package main

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

const ssmTestPrefix = "/psm-test/"

func skipIfNoAWS(t *testing.T) {
	t.Helper()
	if os.Getenv("PSM_INTEGRATION_TEST") == "" {
		t.Skip("skipping integration test: set PSM_INTEGRATION_TEST=1")
	}
}

func loadAWSConfig(t *testing.T) aws.Config {
	t.Helper()
	_ = os.Unsetenv("AWS_PROFILE")
	profile := os.Getenv("PSM_TEST_PROFILE")
	var opts []func(*awsconfig.LoadOptions) error
	if profile != "" {
		opts = append(opts, awsconfig.WithSharedConfigProfile(profile))
	}
	cfg, err := awsconfig.LoadDefaultConfig(context.Background(), opts...)
	if err != nil {
		t.Fatalf("failed to load AWS config: %v", err)
	}
	return cfg
}

func setupSSMTestData(t *testing.T, client *ssm.Client, data map[string]string) {
	t.Helper()
	ctx := context.Background()
	for k, v := range data {
		_, err := client.PutParameter(ctx, &ssm.PutParameterInput{
			Name:      aws.String(k),
			Value:     aws.String(v),
			Type:      ssmtypes.ParameterTypeSecureString,
			Overwrite: aws.Bool(true),
		})
		if err != nil {
			t.Fatalf("setup: failed to put %s: %v", k, err)
		}
	}
}

func cleanupSSMTestData(t *testing.T, client *ssm.Client, keys []string) {
	t.Helper()
	if len(keys) == 0 {
		return
	}
	ctx := context.Background()
	// Delete in batches of 10
	for i := 0; i < len(keys); i += 10 {
		end := i + 10
		if end > len(keys) {
			end = len(keys)
		}
		_, err := client.DeleteParameters(ctx, &ssm.DeleteParametersInput{
			Names: keys[i:end],
		})
		if err != nil {
			t.Logf("cleanup warning: %v", err)
		}
	}
}

// cleanAllSSMTestParams removes all parameters under /psm-test/ to prevent cross-test pollution.
func cleanAllSSMTestParams(t *testing.T, client *ssm.Client) {
	t.Helper()
	ctx := context.Background()
	var names []string
	var nextToken *string
	for {
		out, err := client.GetParametersByPath(ctx, &ssm.GetParametersByPathInput{
			Path:      aws.String(ssmTestPrefix),
			Recursive: aws.Bool(true),
			NextToken: nextToken,
		})
		if err != nil {
			t.Logf("cleanAll warning: %v", err)
			return
		}
		for _, p := range out.Parameters {
			names = append(names, *p.Name)
		}
		if out.NextToken == nil {
			break
		}
		nextToken = out.NextToken
	}
	cleanupSSMTestData(t, client, names)
}

// getSSMParam retrieves a single parameter value directly (strongly consistent).
func getSSMParam(t *testing.T, client *ssm.Client, key string) (string, bool) {
	t.Helper()
	out, err := client.GetParameter(context.Background(), &ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", false
	}
	return *out.Parameter.Value, true
}

func TestSSMStoreGetAll(t *testing.T) {
	skipIfNoAWS(t)
	cfg := loadAWSConfig(t)
	client := ssm.NewFromConfig(cfg)
	store := NewSSMStore(cfg)
	ctx := context.Background()

	cleanAllSSMTestParams(t, client)

	keys := []string{
		ssmTestPrefix + "getall-1",
		ssmTestPrefix + "getall-2",
		ssmTestPrefix + "getall-3",
	}
	data := map[string]string{
		keys[0]: "value1",
		keys[1]: "value2",
		keys[2]: "value3",
	}
	setupSSMTestData(t, client, data)
	defer cleanupSSMTestData(t, client, keys)

	result, err := store.GetAll(ctx)
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}

	for k, v := range data {
		got, ok := result[k]
		if !ok {
			t.Errorf("key %q not found in GetAll result", k)
			continue
		}
		if got != v {
			t.Errorf("key %q: got %q, want %q", k, got, v)
		}
	}
}

func TestSSMStorePutAndDelete(t *testing.T) {
	skipIfNoAWS(t)
	cfg := loadAWSConfig(t)
	store := NewSSMStore(cfg)
	ctx := context.Background()

	key := ssmTestPrefix + "put-delete-test"
	defer cleanupSSMTestData(t, ssm.NewFromConfig(cfg), []string{key})

	// Put
	if err := store.Put(ctx, key, "testvalue"); err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	// Verify via GetParameter (strongly consistent, unlike GetParametersByPath)
	got, ok := getSSMParam(t, ssm.NewFromConfig(cfg), key)
	if !ok {
		t.Fatal("after Put: key not found via GetParameter")
	}
	if got != "testvalue" {
		t.Errorf("after Put: got %q, want %q", got, "testvalue")
	}

	// Delete
	if err := store.Delete(ctx, []string{key}); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify deleted
	if _, exists := getSSMParam(t, ssm.NewFromConfig(cfg), key); exists {
		t.Error("key still exists after Delete")
	}
}

func TestSSMPrune(t *testing.T) {
	skipIfNoAWS(t)
	cfg := loadAWSConfig(t)
	client := ssm.NewFromConfig(cfg)
	store := NewSSMStore(cfg)
	ctx := context.Background()

	cleanAllSSMTestParams(t, client)

	keys := []string{
		ssmTestPrefix + "prune-k1",
		ssmTestPrefix + "prune-k2",
		ssmTestPrefix + "prune-k3",
	}
	setupSSMTestData(t, client, map[string]string{
		keys[0]: "v1",
		keys[1]: "v2",
		keys[2]: "v3",
	})
	defer cleanupSSMTestData(t, client, keys)

	// YAML has only k1, k2
	entries := []Entry{
		{Key: keys[0], Value: "v1"},
		{Key: keys[1], Value: "v2"},
	}

	// Build existing map via GetParameter (strongly consistent, unlike GetParametersByPath)
	existing := make(map[string]string)
	for _, k := range keys {
		if v, ok := getSSMParam(t, client, k); ok {
			existing[k] = v
		}
	}

	// With prune
	actions := plan(entries, existing, true)
	var stdout, stderr bytes.Buffer
	summary := execute(ctx, actions, store, false, &stdout, &stderr)

	if summary.Deleted != 1 {
		t.Errorf("deleted = %d, want 1", summary.Deleted)
	}
	out := stdout.String()
	if !strings.Contains(out, "delete: "+keys[2]) {
		t.Errorf("stdout missing delete line: %s", out)
	}

	// Verify k3 is actually deleted
	after, _ := store.GetAll(ctx)
	if _, exists := after[keys[2]]; exists {
		t.Error("k3 should be deleted after prune")
	}
}

func TestSSMNoPrune(t *testing.T) {
	skipIfNoAWS(t)
	cfg := loadAWSConfig(t)
	client := ssm.NewFromConfig(cfg)
	store := NewSSMStore(cfg)
	ctx := context.Background()

	cleanAllSSMTestParams(t, client)

	keys := []string{
		ssmTestPrefix + "noprune-k1",
		ssmTestPrefix + "noprune-k2",
	}
	setupSSMTestData(t, client, map[string]string{
		keys[0]: "v1",
		keys[1]: "v2",
	})
	defer cleanupSSMTestData(t, client, keys)

	entries := []Entry{{Key: keys[0], Value: "v1"}}
	// Build existing via GetParameter (strongly consistent)
	existing := make(map[string]string)
	for _, k := range keys {
		if v, ok := getSSMParam(t, client, k); ok {
			existing[k] = v
		}
	}
	actions := plan(entries, existing, false)
	var stdout, stderr bytes.Buffer
	summary := execute(ctx, actions, store, false, &stdout, &stderr)

	if summary.Deleted != 0 {
		t.Errorf("deleted = %d, want 0", summary.Deleted)
	}
	// Verify k2 still exists via GetParameter
	if _, exists := getSSMParam(t, client, keys[1]); !exists {
		t.Error("k2 should still exist without prune")
	}
}

func TestSSMDryRun(t *testing.T) {
	skipIfNoAWS(t)
	cfg := loadAWSConfig(t)
	client := ssm.NewFromConfig(cfg)
	store := NewSSMStore(cfg)
	ctx := context.Background()

	cleanAllSSMTestParams(t, client)

	key := ssmTestPrefix + "dryrun-existing"
	setupSSMTestData(t, client, map[string]string{key: "original"})
	defer cleanupSSMTestData(t, client, []string{key, ssmTestPrefix + "dryrun-new"})

	entries := []Entry{
		{Key: ssmTestPrefix + "dryrun-new", Value: "newval"},
		{Key: key, Value: "changed"},
	}
	existing, _ := store.GetAll(ctx)
	actions := plan(entries, existing, false)
	var stdout, stderr bytes.Buffer
	summary := execute(ctx, actions, store, true, &stdout, &stderr)

	if summary.Failed != 0 {
		t.Errorf("failed = %d, want 0", summary.Failed)
	}

	// Verify AWS not changed (use GetParameter for strong consistency)
	got, ok := getSSMParam(t, client, key)
	if !ok {
		t.Fatal("dry-run: original key not found")
	}
	if got != "original" {
		t.Errorf("dry-run changed AWS value: got %q, want %q", got, "original")
	}
	if _, exists := getSSMParam(t, client, ssmTestPrefix+"dryrun-new"); exists {
		t.Error("dry-run created a new parameter")
	}
}

func TestSSMSyncExecute(t *testing.T) {
	skipIfNoAWS(t)
	cfg := loadAWSConfig(t)
	client := ssm.NewFromConfig(cfg)
	store := NewSSMStore(cfg)
	ctx := context.Background()

	cleanAllSSMTestParams(t, client)

	keys := []string{
		ssmTestPrefix + "exec-existing",
		ssmTestPrefix + "exec-unchanged",
	}
	setupSSMTestData(t, client, map[string]string{
		keys[0]: "oldvalue",
		keys[1]: "sameval",
	})
	allKeys := append(keys, ssmTestPrefix+"exec-new")
	defer cleanupSSMTestData(t, client, allKeys)

	entries := []Entry{
		{Key: ssmTestPrefix + "exec-new", Value: "newval"},
		{Key: ssmTestPrefix + "exec-existing", Value: "newvalue2"},
		{Key: ssmTestPrefix + "exec-unchanged", Value: "sameval"},
	}

	// Build existing map via GetParameter (strongly consistent, unlike GetParametersByPath)
	existing := make(map[string]string)
	for _, k := range keys {
		if v, ok := getSSMParam(t, client, k); ok {
			existing[k] = v
		}
	}
	actions := plan(entries, existing, false)

	var stdout, stderr bytes.Buffer
	summary := execute(ctx, actions, store, false, &stdout, &stderr)

	// Check summary
	if summary.Created != 1 {
		t.Errorf("created = %d, want 1", summary.Created)
	}
	if summary.Updated != 1 {
		t.Errorf("updated = %d, want 1", summary.Updated)
	}
	if summary.Unchanged != 1 {
		t.Errorf("unchanged = %d, want 1", summary.Unchanged)
	}
	if summary.Failed != 0 {
		t.Errorf("failed = %d, want 0", summary.Failed)
	}

	// Check stdout format
	out := stdout.String()
	if !strings.Contains(out, "create: "+ssmTestPrefix+"exec-new") {
		t.Errorf("stdout missing create line: %s", out)
	}
	if !strings.Contains(out, "update: "+ssmTestPrefix+"exec-existing") {
		t.Errorf("stdout missing update line: %s", out)
	}
	if strings.Contains(out, ssmTestPrefix+"exec-unchanged") {
		t.Errorf("stdout should not contain unchanged key: %s", out)
	}
	if !strings.Contains(out, "1 created, 1 updated, 0 deleted, 1 unchanged, 0 failed") {
		t.Errorf("stdout missing summary line: %s", out)
	}

	// Values must never appear in output
	if strings.Contains(out, "newval") || strings.Contains(out, "newvalue2") || strings.Contains(out, "sameval") {
		t.Error("stdout contains values - must never output values")
	}
	if strings.Contains(stderr.String(), "newval") || strings.Contains(stderr.String(), "newvalue2") {
		t.Error("stderr contains values - must never output values")
	}
}
