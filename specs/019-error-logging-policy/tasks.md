# Tasks: Define Error Handling and Logging Policy

**Input**: Design documents from `/specs/019-error-logging-policy/`
**Prerequisites**: plan.md, spec.md, research.md, quickstart.md

**Tests**: Not applicable — this is a documentation-only feature (FR-011).

**Organization**: Tasks are grouped by user story. All tasks operate on documentation files only — no code changes.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup

**Purpose**: Prepare the document structure

- [x] T001 Delete redundant documents/backend/coding-standard.md (research R6 — content is subset of deliverable)

---

## Phase 2: User Story 1 — Error Handling Conventions (Priority: P1)

**Goal**: A developer can read the policy and know exactly how to handle errors in psm code

**Independent Test**: Review the Error Handling section and confirm it covers return conventions, panic prohibition, wrapping guidance, and exit codes with MUST/MUST NOT language and examples

- [x] T002 [US1] Write the Error Handling section in documents/backend/error-handling-and-logging.md — include Relationship to Constitution subsection referencing constitution line 54 as authority (research R5), Return Conventions subsection (FR-006: functions that can fail return `(..., error)`), Panic Prohibition subsection (FR-007: no panic anywhere — extends constitution's "library code" scope), Error Wrapping subsection (FR-008: `fmt.Errorf` with `%w`, context at each call site, concrete examples), Exit Codes subsection (FR-009: 0 = success, 1 = failure)

**Checkpoint**: Error handling policy is complete and independently reviewable

---

## Phase 3: User Story 2 — Logging Conventions (Priority: P1)

**Goal**: A developer can read the policy and know which logger to use, what log level to assign, and where output goes

**Independent Test**: Review the Logging section and confirm it covers logger choice, log levels with examples, default level, output routing, log format, and sensitive data rules with MUST/MUST NOT language

- [x] T003 [US2] Write the Logging section in documents/backend/error-handling-and-logging.md — include Logger subsection (FR-001: slog only, default logger via `slog.SetDefault`, no constructor injection per research R1), Log Levels subsection (FR-002: Error/Warn/Info/Debug table with usage guidelines and concrete psm examples per research R3), Default Level subsection (FR-003: Info), Log Format subsection (text format via `slog.TextHandler` per research R2), Sensitive Data subsection (FR-005: never log secret values/tokens/passwords, key paths permitted)

**Checkpoint**: Logging policy is complete and independently reviewable

---

## Phase 4: User Story 3 — Output Routing (Priority: P2)

**Goal**: A developer can distinguish program output (stdout) from log messages and prompts (stderr)

**Independent Test**: Review the Output Routing section and confirm it clearly defines what goes to stdout vs stderr with examples

- [x] T004 [US3] Write the Output Routing subsection in the Logging section of documents/backend/error-handling-and-logging.md — define stdout rules (program output only: diff lines, summaries, action plans per FR-004 and research R4), define stderr rules (all slog messages, interactive prompts like approve `y/N`), include concrete examples of each category

**Checkpoint**: Output routing rules are clear and independently reviewable

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Final validation and consistency

- [x] T005 Review complete documents/backend/error-handling-and-logging.md against spec success criteria — verify SC-001 (log level determinable in 30 seconds from table), SC-002 (all Issue #19 acceptance criteria covered), SC-003 (every rule uses MUST/MUST NOT with at least one example)
- [x] T006 Verify no code changes were made in this branch (FR-011) — run `git diff main --stat` and confirm only documentation files changed

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — can start immediately
- **US1 Error Handling (Phase 2)**: Depends on Setup (T001 removes conflicting file)
- **US2 Logging (Phase 3)**: Depends on T002 (writes to same file, Logging section follows Error Handling section)
- **US3 Output Routing (Phase 4)**: Depends on T003 (Output Routing is a subsection within Logging)
- **Polish (Phase 5)**: Depends on all user stories complete

### Parallel Opportunities

- T002 and T003 write to the same file but different sections — sequential execution recommended to avoid conflicts
- T005 and T006 in Polish phase can run in parallel

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (delete redundant file)
2. Complete Phase 2: US1 Error Handling section
3. **STOP and VALIDATE**: Error handling policy stands alone as useful documentation

### Incremental Delivery

1. Setup → delete coding-standard.md
2. US1 → Error Handling section complete → reviewable
3. US2 → Logging section complete → reviewable
4. US3 → Output Routing subsection complete → reviewable
5. Polish → validate all success criteria

---

## Notes

- All tasks write to a single file: `documents/backend/error-handling-and-logging.md`
- The existing content in that file must be fully replaced, not amended (survey S1)
- Constitution is referenced as authority, not duplicated (research R5)
- No code changes permitted (FR-011)
- Commit after each phase for incremental review
