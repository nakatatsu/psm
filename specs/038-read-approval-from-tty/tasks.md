# Tasks: Read Approval Prompt from /dev/tty

**Input**: Design documents from `/specs/038-read-approval-from-tty/`
**Prerequisites**: plan.md, spec.md, research.md, quickstart.md

**Tests**: Included — constitution mandates test-first (Red-Green cycle).

**Organization**: Tasks grouped by user story for independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Add TtyOpener field to IOStreams for dependency injection

- [x] T001 Add `TtyOpener func() (io.ReadCloser, error)` field to `IOStreams` struct in store.go
- [x] T002 Wire production `TtyOpener` (opens `/dev/tty`) in `runSync` IOStreams initialization in main.go
- [x] T003 Update `testIOStreams` helper to accept and set `TtyOpener` in sync_test.go

**Checkpoint**: IOStreams has TtyOpener field, production wiring compiles, test helper updated

---

## Phase 2: User Story 1 - Pipe Input with Interactive Approval (Priority: P1) MVP

**Goal**: When stdin is piped, read approval from `/dev/tty` so `sops -d ... | psm sync ...` works interactively.

**Independent Test**: Run with piped stdin and injected TtyOpener; verify prompt reads from tty reader.

### Tests for User Story 1

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T004 [US1] Write test: piped stdin + TtyOpener returns "y" reader → approval succeeds, sync executes in sync_test.go
- [x] T005 [US1] Write test: piped stdin + TtyOpener returns "N" reader → approval declined, no changes in sync_test.go

### Implementation for User Story 1

- [x] T006 [US1] Update approval flow in `runSync` in main.go: when `!IsTerminal()` and `!SkipApprove`, call `TtyOpener()` to get reader, pass to `promptApprove`
- [x] T007 [US1] Ensure `/dev/tty` file descriptor is closed after prompt in main.go

**Checkpoint**: Piped stdin with available tty prompts user and respects y/N response

---

## Phase 3: User Story 2 - Non-Interactive Failsafe (Priority: P2)

**Goal**: Auto-cancel with descriptive message when no terminal is available and `--skip-approve` not set.

**Independent Test**: Inject TtyOpener that returns error; verify auto-cancel with message to stderr.

### Tests for User Story 2

- [x] T008 [US2] Write test: TtyOpener returns error → auto-cancel with message to stderr, exit 0, no changes in sync_test.go
- [x] T009 [US2] Write test: TtyOpener returns error + `--skip-approve` → sync proceeds normally in sync_test.go

### Implementation for User Story 2

- [x] T010 [US2] Handle `TtyOpener()` error in `runSync` in main.go: print "No terminal available for approval prompt. Use --skip-approve for non-interactive usage." to stderr, return 0

**Checkpoint**: No-tty scenario auto-cancels with actionable message; --skip-approve bypasses as before

---

## Phase 4: User Story 3 - Direct Terminal Usage Unchanged (Priority: P3)

**Goal**: Existing non-piped terminal usage continues to work identically.

**Independent Test**: Existing tests pass; stdin-is-terminal path unchanged.

### Tests for User Story 3

- [x] T011 [US3] Verify existing approval tests still pass (terminal stdin reads from stdin, not tty) in sync_test.go

### Implementation for User Story 3

- [x] T012 [US3] Ensure `runSync` approval flow only calls `TtyOpener` when `!IsTerminal()` — terminal stdin path remains unchanged in main.go

**Checkpoint**: All existing tests pass, backward compatibility confirmed

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Documentation and final validation

- [x] T013 Update requirements-v2.md section 3.1 to reflect /dev/tty behavior
- [x] T014 Update README.md if approval prompt behavior is documented there
- [x] T015 Run `go test ./...` and `go vet ./...` for final validation

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — T001 → T002, T003 can be parallel after T001
- **US1 (Phase 2)**: Depends on Phase 1 completion
- **US2 (Phase 3)**: Depends on Phase 1 completion, can run parallel with US1
- **US3 (Phase 4)**: Depends on Phase 1 completion, can run parallel with US1/US2
- **Polish (Phase 5)**: Depends on all user stories complete

### Within Each User Story

- Tests MUST be written and FAIL before implementation (constitution mandate)
- Implementation follows Red-Green cycle

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001-T003)
2. Complete Phase 2: US1 tests then implementation (T004-T007)
3. **STOP and VALIDATE**: `go test ./...` passes
4. Continue with US2, US3, Polish

### Sequential Execution (Single Developer)

1. T001 → T002 + T003 → T004 → T005 → T006 → T007 → T008 → T009 → T010 → T011 → T012 → T013 → T014 → T015

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story
- Constitution mandates test-first: write test → confirm fail → implement → confirm pass
- Commit after each phase checkpoint
