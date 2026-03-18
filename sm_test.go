package main

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

const smTestPrefix = "psm-test/"

func loadAWSConfigSM(t *testing.T) aws.Config {
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

func setupSMTestData(t *testing.T, client *secretsmanager.Client, data map[string]string) {
	t.Helper()
	ctx := context.Background()
	for k, v := range data {
		_, err := client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(k),
			SecretString: aws.String(v),
		})
		if err != nil {
			// Secret may be marked for deletion — restore it first
			_, _ = client.RestoreSecret(ctx, &secretsmanager.RestoreSecretInput{
				SecretId: aws.String(k),
			})
			_, err2 := client.PutSecretValue(ctx, &secretsmanager.PutSecretValueInput{
				SecretId:     aws.String(k),
				SecretString: aws.String(v),
			})
			if err2 != nil {
				t.Fatalf("setup: failed to put %s: %v (create: %v)", k, err2, err)
			}
		}
	}
}

func cleanupSMTestData(t *testing.T, client *secretsmanager.Client, keys []string) {
	t.Helper()
	ctx := context.Background()
	for _, k := range keys {
		_, err := client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
			SecretId:                   aws.String(k),
			ForceDeleteWithoutRecovery: aws.Bool(true),
		})
		if err != nil {
			t.Logf("cleanup warning: %v", err)
		}
	}
}

// getSMSecret retrieves a single secret value directly (strongly consistent).
func getSMSecret(t *testing.T, client *secretsmanager.Client, key string) (string, bool) {
	t.Helper()
	out, err := client.GetSecretValue(context.Background(), &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(key),
	})
	if err != nil {
		return "", false
	}
	if out.SecretString == nil {
		return "", false
	}
	return *out.SecretString, true
}

// buildSMExisting builds an existing map by directly getting each known key (avoids ListSecrets eventual consistency).
func buildSMExisting(t *testing.T, client *secretsmanager.Client, keys []string) map[string]string {
	t.Helper()
	result := make(map[string]string)
	for _, k := range keys {
		v, ok := getSMSecret(t, client, k)
		if ok {
			result[k] = v
		}
	}
	return result
}

func TestSMStoreGetAll(t *testing.T) {
	skipIfNoAWS(t)
	cfg := loadAWSConfigSM(t)
	client := secretsmanager.NewFromConfig(cfg)
	_ = NewSMStore(cfg) // verify construction
	ctx := context.Background()

	keys := []string{
		smTestPrefix + "getall-1",
		smTestPrefix + "getall-2",
	}
	data := map[string]string{
		keys[0]: "smval1",
		keys[1]: "smval2",
	}
	setupSMTestData(t, client, data)
	defer cleanupSMTestData(t, client, keys)

	// Use direct GetSecretValue to verify (ListSecrets is eventually consistent)
	for k, v := range data {
		got, ok := getSMSecret(t, client, k)
		if !ok {
			t.Errorf("key %q not found via GetSecretValue", k)
			continue
		}
		if got != v {
			t.Errorf("key %q: got %q, want %q", k, got, v)
		}
	}

	// Also verify ListSecrets-based GetAll works (may see stale data, so just check no error)
	_, err := NewSMStore(cfg).GetAll(ctx)
	if err != nil {
		t.Errorf("GetAll returned error: %v", err)
	}
}

func TestSMStorePutAndDelete(t *testing.T) {
	skipIfNoAWS(t)
	cfg := loadAWSConfigSM(t)
	store := NewSMStore(cfg)
	ctx := context.Background()

	key := smTestPrefix + "put-delete-test"
	defer cleanupSMTestData(t, secretsmanager.NewFromConfig(cfg), []string{key})

	if err := store.Put(ctx, key, "testvalue"); err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	// Verify via direct GetSecretValue (strongly consistent)
	got, ok := getSMSecret(t, secretsmanager.NewFromConfig(cfg), key)
	if !ok {
		t.Fatal("after Put: key not found via GetSecretValue")
	}
	if got != "testvalue" {
		t.Errorf("after Put: got %q, want %q", got, "testvalue")
	}

	if err := store.Delete(ctx, []string{key}); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestSMSyncExecute(t *testing.T) {
	skipIfNoAWS(t)
	cfg := loadAWSConfigSM(t)
	client := secretsmanager.NewFromConfig(cfg)
	store := NewSMStore(cfg)
	ctx := context.Background()

	keys := []string{
		smTestPrefix + "exec-existing",
		smTestPrefix + "exec-unchanged",
	}
	setupSMTestData(t, client, map[string]string{
		keys[0]: "oldvalue",
		keys[1]: "sameval",
	})
	allKeys := append(keys, smTestPrefix+"exec-new")
	defer cleanupSMTestData(t, client, allKeys)

	entries := []Entry{
		{Key: smTestPrefix + "exec-new", Value: "newval"},
		{Key: smTestPrefix + "exec-existing", Value: "updated"},
		{Key: smTestPrefix + "exec-unchanged", Value: "sameval"},
	}

	// Build existing map via direct GetSecretValue (ListSecrets is eventually consistent)
	existing := buildSMExisting(t, client, keys)
	actions := plan(entries, existing)

	var stdout, stderr bytes.Buffer
	summary := execute(ctx, actions, store, &stdout, &stderr)

	if summary.Created != 1 {
		t.Errorf("created = %d, want 1", summary.Created)
	}
	if summary.Updated != 1 {
		t.Errorf("updated = %d, want 1", summary.Updated)
	}
	if summary.Unchanged != 1 {
		t.Errorf("unchanged = %d, want 1", summary.Unchanged)
	}

	out := stdout.String()
	if !strings.Contains(out, "create: "+smTestPrefix+"exec-new") {
		t.Errorf("stdout missing create line: %s", out)
	}
	if !strings.Contains(out, "update: "+smTestPrefix+"exec-existing") {
		t.Errorf("stdout missing update line: %s", out)
	}
}
