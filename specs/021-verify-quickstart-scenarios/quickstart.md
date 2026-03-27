# Quickstart: Integration Test Script

## Prerequisites

1. DevContainer is running (`example/.devcontainer/`)
2. AWS SSO login is complete:
   ```bash
   aws sso login --sso-session psm-sandbox --use-device-code
   ```

## Run Tests

```bash
cd /workspace/example
./test.sh
```

## Expected Output

```
=== Prerequisites ===
psm ... ok
aws credentials ... ok

=== Setup: cleaning /psm-test/ parameters ===
cleanup ... ok

=== Scenario 1/7: Dry-run ===
PASS

=== Scenario 2/7: Sync with --skip-approve ===
PASS

=== Scenario 3/7: Delete with --delete ===
PASS

=== Scenario 4/7: Conflict detection ===
PASS

=== Scenario 5/7: Debug logging ===
PASS

=== Scenario 6/7: Non-terminal auto-decline ===
PASS

=== Scenario 7/7: No changes ===
PASS

=== Results: 7 passed, 0 failed ===
```

## Troubleshooting

| Symptom | Cause | Fix |
|---------|-------|-----|
| `aws credentials ... FAIL` | SSO session expired | Re-run `aws sso login --sso-session psm-sandbox --use-device-code` |
| `psm ... FAIL` | psm not on PATH | Check DevContainer build or install psm |
| Scenario fails with AccessDenied | IAM permissions insufficient | Verify sandbox role has SSM full access |
