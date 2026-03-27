# Feature Specification: Integration Test Script for example/

**Feature Branch**: `021-verify-quickstart-scenarios`
**Created**: 2026-03-25
**Status**: Draft
**Input**: GitHub Issue #21 — Integration testing using an AWS development account

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Run All Integration Tests in One Command (Priority: P1)

A developer who has completed AWS SSO login inside the DevContainer wants to verify all psm behaviors against a real AWS sandbox by running a single script, instead of executing each README scenario by hand.

**Why this priority**: This is the core purpose of the feature. Without a runnable script, nothing else matters.

**Independent Test**: Run `./test.sh` in the DevContainer after AWS login. The script executes all scenarios and reports pass/fail for each.

**Acceptance Scenarios**:

1. **Given** the DevContainer is running and AWS SSO session is active, **When** the user runs `./test.sh`, **Then** all test scenarios execute sequentially with clear pass/fail output for each
2. **Given** one scenario fails, **When** the script continues, **Then** remaining scenarios still execute and the final summary shows total passed/failed counts
3. **Given** AWS SSO session has expired, **When** the user runs `./test.sh`, **Then** the first AWS-dependent scenario fails with a clear error and the script exits early

---

### User Story 2 - Verify Each psm Behavior Individually (Priority: P1)

The script must cover the following 7 behavioral scenarios, each independently verifiable:

1. **Dry-run**: `--dry-run` shows plan with `(dry-run)` suffix, no AWS changes
2. **Sync with `--skip-approve`**: piped input + `--skip-approve` applies changes to AWS
3. **Delete with `--delete`**: delete pattern file removes matching keys from AWS
4. **Conflict detection**: sync YAML and delete pattern overlap causes abort (exit code 1)
5. **Debug logging**: `--debug` produces `level=DEBUG` on stderr
6. **Non-terminal auto-decline**: piped input without `--skip-approve` makes no changes
7. **No changes**: syncing identical data shows summary only, no prompt

**Why this priority**: Each scenario validates a distinct behavioral requirement from `specs/requirements.md`. Missing any one leaves a gap.

**Independent Test**: Each scenario can be run and verified in isolation by checking exit codes, stdout content, stderr content, and AWS state.

**Acceptance Scenarios**:

1. **Given** test data is seeded in AWS, **When** dry-run scenario executes, **Then** stdout contains `(dry-run)` and AWS state is unchanged
2. **Given** test data differs from AWS, **When** skip-approve scenario executes, **Then** AWS parameters match the test data
3. **Given** a delete pattern matches an existing AWS key not in sync YAML, **When** delete scenario executes, **Then** the key is removed from AWS
4. **Given** a delete pattern matches a key that is also in sync YAML, **When** conflict scenario executes, **Then** exit code is 1 and no changes are made
5. **Given** `--debug` is passed, **When** debug scenario executes, **Then** stderr contains `level=DEBUG`
6. **Given** input is piped without `--skip-approve`, **When** auto-decline scenario executes, **Then** AWS state is unchanged
7. **Given** AWS state already matches sync YAML, **When** no-changes scenario executes, **Then** stdout shows `0 created, 0 updated, 0 deleted` with no prompt displayed

---

### User Story 3 - Clean Test Environment (Priority: P2)

The script must set up and tear down its own test data so that repeated runs do not pollute the AWS sandbox or interfere with each other.

**Why this priority**: Without cleanup, test runs leave orphaned parameters that cause false positives/negatives in subsequent runs.

**Independent Test**: Run the script twice. The second run produces identical results.

**Acceptance Scenarios**:

1. **Given** the script has run before, **When** the script runs again, **Then** all scenarios pass with the same results
2. **Given** the script is interrupted midway, **When** the user runs it again, **Then** the setup phase cleans residual data before proceeding

---

### Edge Cases

- What happens when the AWS profile is not set or invalid? The script should fail early with a clear message.
- What happens when psm binary is not found on PATH? The script should fail at the prerequisite check.
- What happens when a test scenario leaves unexpected parameters? The cleanup phase at the start of the script removes all test-prefixed parameters.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The script MUST be a single executable shell script at `example/test.sh`
- **FR-002**: The script MUST verify prerequisites (psm on PATH, AWS credentials valid) before running any test scenario
- **FR-003**: The script MUST use a dedicated key prefix (e.g., `/psm-test/`) to isolate test data from real parameters
- **FR-004**: The script MUST clean up all test-prefixed parameters at the start of each run (setup phase)
- **FR-005**: The script MUST clean up all test-prefixed parameters at the end of each run (teardown phase)
- **FR-006**: The script MUST report pass/fail for each scenario with a final summary count
- **FR-007**: The script MUST exit with code 0 when all scenarios pass, and non-zero when any scenario fails
- **FR-008**: The script MUST document prerequisites (AWS login, DevContainer) in a header comment
- **FR-009**: Each scenario MUST verify its expected outcome (exit code, stdout content, stderr content, or AWS state) rather than just running the command
- **FR-010**: The script MUST use temporary files for test YAML data, created and cleaned up automatically

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All 7 behavioral scenarios are covered and independently verifiable
- **SC-002**: A developer can run the full test suite with a single command (`./test.sh`) in under 2 minutes
- **SC-003**: Consecutive runs produce identical results (no state leakage between runs)
- **SC-004**: The script exit code reliably indicates overall pass/fail status
