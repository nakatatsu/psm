# Tasks: Integration Test Script

**Input**: Design documents from `/specs/021-verify-quickstart-scenarios/`
**Prerequisites**: plan.md, spec.md, research.md, quickstart.md

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup

**Purpose**: Script skeleton and helper functions

- [x] T001 Create `example/test.sh` with shebang, header comment (prerequisites documentation), and `set -uo pipefail`
- [x] T002 Implement helper functions in `example/test.sh`: `pass()`, `fail()`, `cleanup_aws()` for test reporting and AWS parameter cleanup using `aws ssm get-parameters-by-path` + `aws ssm delete-parameters`
- [x] T003 Implement prerequisite checks in `example/test.sh`: verify `psm` on PATH, verify AWS credentials with `aws sts get-caller-identity`

---

## Phase 2: Foundational (Test Data Setup)

**Purpose**: SOPS/age key generation and test data preparation that all scenarios depend on

- [x] T004 Implement test data setup in `example/test.sh`: generate ephemeral age key, create `.sops.yaml`, create test YAML with `/psm-test/` prefixed keys, encrypt with sops
- [x] T005 Implement initial cleanup in `example/test.sh`: call `cleanup_aws()` to remove all `/psm-test/` parameters before scenarios start
- [x] T006 Seed baseline data in `example/test.sh`: each scenario seeds its own data as needed via `put_param()`

**Checkpoint**: Test environment is ready — scenario implementation can begin

---

## Phase 3: User Story 1 — Run All Scenarios (Priority: P1) 🎯 MVP

**Goal**: All 8 behavioral scenarios execute and verify outcomes

**Independent Test**: Run `./test.sh` — each scenario reports PASS/FAIL

### Scenario implementations in `example/test.sh`

- [x] T007 [US1] Scenario 1: Dry-run — run `sops -d | psm sync --store ssm --dry-run /dev/stdin`, assert stdout contains `(dry-run)`, assert AWS parameters unchanged via `aws ssm get-parameter`
- [x] T008 [US1] Scenario 2: Sync with `--skip-approve` — run `sops -d | psm sync --store ssm --skip-approve /dev/stdin`, assert exit code 0, assert AWS parameters match test YAML via `aws ssm get-parameter`
- [x] T009 [US1] Scenario 3: Delete with `--delete` — create delete pattern file, seed a deletable parameter, run `psm sync --store ssm --skip-approve --delete <pattern-file> <yaml>`, assert parameter deleted via `aws ssm get-parameter`
- [x] T010 [US1] Scenario 4: Conflict detection — create delete pattern that conflicts with sync YAML, run `psm sync --store ssm --skip-approve --delete <pattern-file> <yaml>`, assert exit code 1, assert no AWS changes
- [x] T011 [US1] Scenario 5: Debug logging — run `psm sync --store ssm --debug --dry-run <yaml>`, assert stderr contains `level=DEBUG`
- [x] T012 [US1] Scenario 6: `--prune` error — run `psm sync --store ssm --prune <yaml>`, assert exit code 1, assert stderr contains `--delete`
- [x] T013 [US1] Scenario 7: Non-terminal auto-decline — pipe input without `--skip-approve`, assert AWS parameters unchanged
- [x] T014 [US1] Scenario 8: No changes — sync YAML identical to AWS state, assert stdout contains `0 created, 0 updated, 0 deleted`

**Checkpoint**: All 8 scenarios pass independently

---

## Phase 4: User Story 3 — Clean Test Environment (Priority: P2)

**Goal**: Teardown and idempotent re-runs

**Independent Test**: Run `./test.sh` twice — second run produces identical results

- [x] T015 [US3] Implement final teardown in `example/test.sh`: call `cleanup_aws()` after all scenarios, remove temporary files via trap
- [x] T016 [US3] Implement summary output in `example/test.sh`: print `=== Results: N passed, M failed ===`, exit with code 0 if all passed, 1 otherwise

---

## Phase 5: Polish

**Purpose**: Final verification

- [x] T017 Make `example/test.sh` executable (`chmod +x`)
- [x] T018 Update `example/README.md` File Structure section to reflect test.sh description
- [ ] T019 Run quickstart.md validation: execute `example/test.sh` end-to-end in DevContainer with active AWS session

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: No dependencies — start immediately
- **Phase 2 (Foundational)**: Depends on Phase 1 completion
- **Phase 3 (Scenarios)**: Depends on Phase 2 completion — all scenario tasks are sequential (shared AWS state)
- **Phase 4 (Cleanup)**: Depends on Phase 3 completion
- **Phase 5 (Polish)**: Depends on all phases complete

### Within Phase 3

Scenarios T007–T014 are sequential because they share AWS state and some scenarios depend on the state left by previous scenarios (e.g., T008 creates data that T009/T014 verify against).

---

## Implementation Strategy

### MVP First

1. Complete Phases 1–3: script with all 8 scenarios
2. **STOP and VALIDATE**: Run in DevContainer with active AWS session
3. Complete Phase 4: cleanup and summary
4. Complete Phase 5: final polish
