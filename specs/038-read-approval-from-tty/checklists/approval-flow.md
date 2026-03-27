# Checklist: Approval Flow Requirements Quality

**Purpose**: Validate that the /dev/tty approval prompt requirements are complete, clear, and implementation-ready.
**Created**: 2026-03-27
**Feature**: 038-read-approval-from-tty

## Requirement Completeness

- [ ] CHK001 Are all three input source scenarios (terminal stdin, piped stdin with tty, no tty at all) explicitly defined with expected behavior? [Completeness, Spec User Stories 1-3]
- [ ] CHK002 Is the exact auto-cancel message content specified for the no-tty scenario? [Completeness, Spec FR-005]
- [ ] CHK003 Are requirements defined for closing the `/dev/tty` file descriptor after use? [Completeness, Gap]

## Requirement Clarity

- [ ] CHK004 Is "clear message" in FR-002 and FR-005 quantified with specific content or format? [Clarity, Spec FR-002]
- [ ] CHK005 Is the decision logic for when to use `/dev/tty` vs stdin unambiguous (exactly when `IsTerminal()` is false)? [Clarity, Spec FR-001]

## Requirement Consistency

- [ ] CHK006 Are FR-002 (auto-cancel with message) and FR-005 (display clear message) consistent and non-duplicative? [Consistency, Spec FR-002/FR-005]
- [ ] CHK007 Is the exit code behavior consistent between auto-cancel (no tty) and user-decline (responds "N")? [Consistency, Gap]

## Acceptance Criteria Quality

- [ ] CHK008 Are acceptance scenarios for User Story 1 testable without manual terminal interaction (via DI)? [Measurability, Spec US1]
- [ ] CHK009 Is the success criterion SC-003 ("all existing tests pass without modification") realistic given new `TtyOpener` field on `IOStreams`? [Measurability, Spec SC-003]

## Edge Case Coverage

- [ ] CHK010 Is fallback behavior specified when `/dev/tty` open succeeds but read fails (e.g., permission error mid-read)? [Edge Case, Spec Edge Cases]
- [ ] CHK011 Is behavior defined for EOF (Ctrl+D) at the approval prompt from `/dev/tty`? [Edge Case, Spec Edge Cases]
- [ ] CHK012 Is behavior specified when `/dev/tty` exists but returns permission denied on open? [Edge Case, Spec Edge Cases]

## Dependencies and Assumptions

- [ ] CHK013 Is the Unix-only assumption (no Windows `/dev/tty`) explicitly scoped as out-of-scope in requirements, not just assumptions? [Completeness, Spec Assumptions]
- [ ] CHK014 Is the assumption "stdin is terminal → read from stdin" explicitly stated as a requirement, not just an assumption? [Clarity, Spec Assumptions]
