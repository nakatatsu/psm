# Tasks: Add Approve Flow and Replace --prune with --delete

**Input**: Design documents from `/specs/018-approve-delete-flow/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/cli-schema.md, quickstart.md

**Tests**: Required — Constitution Principle III (Test-First) is NON-NEGOTIABLE.

**Organization**: Tasks are grouped by user story. Test-first: write test → confirm fail → implement → confirm pass.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup

**Purpose**: slog integration and Config struct update — shared infrastructure for all stories

- [x] T001 Initialize slog with TextHandler to stderr in main() of main.go — default level Info. Flag registration is handled in T009, slog level is set based on Config.Debug after parseArgs returns (research R7)
- [x] T002 Update Config struct in store.go — remove Prune field, add DeleteFile string, SkipApprove bool, Debug bool (research R8, data-model.md)

---

## Phase 2: Foundational — Display Plan + Approve Prompt

**Purpose**: Core plan display and approval mechanism that ALL user stories depend on

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T003 [P] Write tests for displayPlan() in sync_test.go — table-driven: actions with creates/updates/deletes produce expected stdout output, skip actions do not appear, empty actions produce no output (research R6)
- [x] T004 [P] Write tests for approve prompt in sync_test.go — test cases: input "y" returns true, input "Y" returns true, input "N" returns false, empty input (Enter) returns false, arbitrary input returns false. Use io.Reader injection for testability
- [x] T005 Implement displayPlan(actions []Action, stdout io.Writer) in sync.go — render action lines to stdout without executing (research R6)
- [x] T006 Implement promptApprove(reader io.Reader, writer io.Writer) bool in sync.go — write "Proceed? [y/N] " to writer, read one line from reader, return true only for "y" or "Y" (research R1)

**Checkpoint**: displayPlan and promptApprove are functional and tested

---

## Phase 3: User Story 1 — Operator Reviews Changes Before Execution (Priority: P1) 🎯 MVP

**Goal**: All sync operations require approval before execution

**Independent Test**: Run psm sync with changes pending. Verify plan is displayed, prompt appears, and execution depends on user input.

### Tests for User Story 1

- [x] T007 [P] [US1] Write tests for runSync approve flow in main_test.go or sync_test.go — test cases: (1) changes exist + user approves → execute called, (2) changes exist + user declines → no execution + exit 0, (3) --skip-approve → execute called without prompt, (4) --dry-run → displayPlan called + no prompt + no execute, (5) no changes (all skip) → summary only + no prompt (FR-015), (6) non-terminal stdin without --skip-approve → exit 0 no changes (research R2)
- [x] T008 [P] [US1] Write tests for parseArgs new flags in main_test.go — test cases: --skip-approve sets SkipApprove=true, --debug sets Debug=true, --dry-run still works, all flags combined, unknown flag → error

### Implementation for User Story 1

- [x] T009 [US1] Update parseArgs in main.go — register --skip-approve and --debug flags for sync subcommand, --debug for export subcommand, map to Config fields (research R8)
- [x] T010 [US1] Update runSync in main.go — insert approval flow between plan() and execute(): call displayPlan(), check if changes exist (any non-skip action), if --dry-run just show summary, if no changes show summary only, check terminal (os.Stdin.Stat), prompt if interactive and not --skip-approve, execute only if approved (research R1, R2, R6)
- [x] T011 [US1] Update execute() in sync.go — remove all display logic from execute(). Display is fully handled by displayPlan() before approval. execute() receives an io.Writer for progress/result output during execution (e.g., errors). Dry-run flow uses displayPlan() + summary only, never calls execute()

**Checkpoint**: psm sync requires approval. --skip-approve, --dry-run, non-terminal all work correctly.

---

## Phase 4: User Story 4 — Remove --prune Flag (Priority: P1)

**Goal**: --prune returns an error with migration message

### Tests for User Story 4

- [x] T012 [P] [US4] Write tests for --prune removal in main_test.go — test case: passing --prune returns error containing "removed" and "--delete" (FR-005)

### Implementation for User Story 4

- [x] T013 [US4] Remove --prune flag registration from parseArgs in main.go — detect if "--prune" is in args manually, return error "‑‑prune has been removed. Use --delete <file> with regex patterns instead." (FR-005)
- [x] T014 [US4] Remove prune parameter from plan() in sync.go — remove the prune bool parameter and the prune block (lines 29-40). Update all callers
- [x] T015 [US4] Update existing prune tests in sync_test.go — remove "prune deletes missing keys" and "no prune keeps missing keys" test cases, update plan() call signatures in remaining tests
- [x] T016 [US4] Update integration tests in ssm_test.go and sm_test.go — remove TestSSMPrune, TestSSMNoPrune and equivalent SM tests, update remaining tests that call plan() with prune parameter

**Checkpoint**: --prune returns error. plan() no longer has prune parameter. All existing tests pass.

---

## Phase 5: User Story 2 — Operator Deletes Obsolete Keys Using Pattern File (Priority: P1)

**Goal**: --delete <file> with regex patterns, conflict detection, unmanaged key warnings

### Tests for User Story 2

- [x] T017 [P] [US2] Write tests for parseDeletePatterns() in delete_test.go — test cases: valid patterns compile, invalid regex returns error with pattern name, empty list returns empty slice, non-YAML file returns error (research R3)
- [x] T018 [P] [US2] Write tests for planDeletes() in delete_test.go — test cases: (1) key matches pattern + not in YAML → delete action, (2) key matches pattern + in YAML → conflict, (3) key matches no pattern + not in YAML → unmanaged, (4) key in YAML + no pattern match → not affected, (5) multiple patterns match same key → single delete, (6) no AWS keys match any pattern → empty result (research R5)
- [x] T019 [P] [US2] Write tests for conflict detection abort in delete_test.go or sync_test.go — test case: if planDeletes returns conflicts, entire operation (including sync creates/updates) must not execute
- [x] T020 [P] [US2] Write tests for parseArgs --delete flag in main_test.go — test cases: --delete sets DeleteFile, --delete without sync file → error (FR-008), --delete with sync file → both parsed

### Implementation for User Story 2

- [x] T021 [US2] Implement parseDeletePatterns(data []byte) ([]*regexp.Regexp, error) in delete.go — parse YAML list of strings, compile each as regexp, return error on invalid pattern with pattern text in message (research R3, FR-010)
- [x] T022 [US2] Implement planDeletes(existing map[string]string, yamlKeys map[string]bool, patterns []*regexp.Regexp) (deletes []Action, conflicts []string, unmanaged []string) in delete.go — iterate existing keys, match against patterns, classify as delete/conflict/unmanaged (research R4, R5)
- [x] T023 [US2] Add --delete flag to parseArgs in main.go — string flag for sync subcommand, store in Config.DeleteFile (research R8)
- [x] T024 [US2] Integrate delete flow into runSync in main.go — if DeleteFile set: read file, parseDeletePatterns, call planDeletes, check conflicts (abort if any), log unmanaged keys as slog.Warn, merge delete actions into action list, display via displayPlan, prompt, execute (FR-006 through FR-012)

**Checkpoint**: psm sync --delete works. Conflicts abort. Unmanaged keys warn. All integrated with approval flow.

---

## Phase 6: User Story 3 — Debug Logging (Priority: P3)

**Goal**: --debug enables Debug-level slog output

### Tests for User Story 3

- [x] T025 [P] [US3] Write tests for debug logging in main_test.go or sync_test.go — test cases: (1) with --debug, Debug-level messages are emitted, (2) without --debug, Debug-level messages are hidden (FR-013, FR-014)

### Implementation for User Story 3

- [x] T026 [US3] Add slog.Debug calls at key points — in planDeletes for regex match results, in execute for API calls, in parseDeletePatterns for pattern compilation (FR-013)

**Checkpoint**: --debug shows diagnostic output. Without --debug, only Error/Warn/Info visible.

---

## Phase 7: Polish & Cross-Cutting Concerns

- [x] T027 Update README.md — document new flags (--delete, --skip-approve, --debug), remove --prune documentation, add delete pattern file format, add migration guide from --prune
- [x] T028 Run full test suite and verify all tests pass — `go test ./...`, `go vet ./...`
- [ ] T029 Run quickstart.md verification checklist — execute each scenario and confirm expected behavior

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Phase 1 (Config struct, slog) — BLOCKS all user stories
- **US1 Approve (Phase 3)**: Depends on Phase 2 (displayPlan, promptApprove)
- **US4 Remove --prune (Phase 4)**: Depends on Phase 3 (approve flow must be in place before removing prune)
- **US2 Delete Patterns (Phase 5)**: Depends on Phase 4 (prune removed, plan() signature clean)
- **US3 Debug (Phase 6)**: Depends on Phase 5 (needs delete flow code to add debug logging to)
- **Polish (Phase 7)**: Depends on all user stories complete

### Within Each User Story

- Tests MUST be written and FAIL before implementation (Constitution III)
- Test tasks marked [P] within a story can run in parallel
- Implementation tasks are sequential within a story

### Parallel Opportunities

- T003 + T004 (Phase 2 tests)
- T007 + T008 (US1 tests)
- T017 + T018 + T019 + T020 (US2 tests)
- T025 (US3 tests — independent)

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Phase 1: Setup (slog, Config)
2. Phase 2: Foundational (displayPlan, promptApprove)
3. Phase 3: US1 Approve flow
4. **STOP and VALIDATE**: psm sync now requires approval

### Incremental Delivery

1. Setup + Foundational → approval infrastructure ready
2. US1 → approval works → deployable MVP
3. US4 → --prune removed → safety improved
4. US2 → --delete added → full replacement complete
5. US3 → --debug added → troubleshooting capability
6. Polish → README, full validation

---

## Notes

- Test-first is NON-NEGOTIABLE per Constitution Principle III
- All new functions accept io.Reader/io.Writer for testability (no direct os.Stdin/os.Stdout in logic)
- slog messages go to stderr; program output (diff lines, summaries) to stdout per logging policy
- No new third-party dependencies — stdlib only (regexp, bufio, os, log/slog)
