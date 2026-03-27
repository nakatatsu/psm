# Research: Integration Test Script

## R1: Test Data Isolation Strategy

**Decision**: Use `/psm-test/` prefix for all test parameters, matching the existing convention in `ssm_test.go`.

**Rationale**: The Go integration tests already use `ssmTestPrefix = "/psm-test/"` with `cleanAllSSMTestParams()` for cleanup. Reusing the same prefix and pattern ensures consistency. The cleanup uses `aws ssm get-parameters-by-path` + batch delete.

**Alternatives considered**:
- Random prefix per run — adds complexity, harder to debug orphaned data
- Timestamp-based prefix — same issues as random

## R2: AWS State Verification Method

**Decision**: Use `aws ssm get-parameter` (strongly consistent) to verify individual parameter values after operations.

**Rationale**: `GetParametersByPath` has eventual consistency. The Go tests (`ssm_test.go:104-115`) already use `GetParameter` for this reason. The shell script follows the same approach.

**Alternatives considered**:
- `GetParametersByPath` — eventually consistent, may cause flaky tests
- Sleep-and-retry — fragile, slow

## R3: SOPS/age Integration in Tests

**Decision**: The test script generates its own age key and SOPS config, then encrypts test YAML. This tests the full SOPS pipeline as documented in `example/README.md`.

**Rationale**: Some scenarios (dry-run, skip-approve) use piped `sops -d` input per the README workflow. Generating ephemeral keys keeps tests self-contained.

**Alternatives considered**:
- Use plaintext YAML only — skips the SOPS pipe path, which is a key user workflow
- Use committed encrypted file — requires a committed key, security concern

## R4: Script Error Handling

**Decision**: Each scenario runs independently. A failure increments the fail counter but does not abort subsequent scenarios. The script exits non-zero if any scenario failed.

**Rationale**: Maximum diagnostic value — seeing all failures at once is more useful than stopping at the first one.

**Alternatives considered**:
- `set -e` (fail-fast) — hides subsequent failures, less useful for debugging
