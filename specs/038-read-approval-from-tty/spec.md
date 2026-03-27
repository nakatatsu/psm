# Feature Specification: Read Approval Prompt from /dev/tty

**Feature Branch**: `038-read-approval-from-tty`
**Created**: 2026-03-27
**Status**: Draft
**Input**: GitHub Issue #38

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Pipe Input with Interactive Approval (Priority: P1)

As a user, I want to pipe SOPS-decrypted secrets into `psm sync` and still be prompted for approval, so that I can use the standard workflow without bypassing safety checks.

**Why this priority**: This is the core use case — piping is the primary way users invoke psm with SOPS, and the current behavior either silently cancels or requires skipping approval entirely.

**Independent Test**: Can be fully tested by running `echo "test" | psm sync --store ssm /dev/stdin` in a terminal and verifying that the approval prompt appears and accepts user input.

**Acceptance Scenarios**:

1. **Given** stdin is redirected (piped input), **When** the user runs `psm sync` without `--skip-approve`, **Then** the approval prompt is displayed and reads the user's response interactively.
2. **Given** stdin is redirected and the user responds "y" to the approval prompt, **When** the sync executes, **Then** changes are applied as expected.
3. **Given** stdin is redirected and the user responds "N" or presses Enter to the approval prompt, **When** the sync is cancelled, **Then** no changes are made.

---

### User Story 2 - Non-Interactive Failsafe (Priority: P2)

As a CI/CD pipeline operator, I want `psm sync` to auto-cancel when no terminal is available and `--skip-approve` is not set, so that unattended runs do not accidentally apply changes.

**Why this priority**: Safety mechanism — prevents accidental changes in environments like cron jobs or headless CI where no human can respond.

**Independent Test**: Can be tested by running `psm sync` in an environment with no controlling terminal (no `/dev/tty`) and without `--skip-approve`, verifying that it auto-cancels.

**Acceptance Scenarios**:

1. **Given** no terminal is available (no `/dev/tty`), **When** the user runs `psm sync` without `--skip-approve`, **Then** the command auto-cancels with a clear message indicating no terminal is available.
2. **Given** no terminal is available, **When** the user runs `psm sync` with `--skip-approve`, **Then** the command proceeds without prompting (existing behavior preserved).

---

### User Story 3 - Direct Terminal Usage Unchanged (Priority: P3)

As a user running psm directly in a terminal (no pipe), I want the approval prompt to continue working as before.

**Why this priority**: Ensures backward compatibility — existing non-piped workflows must not break.

**Independent Test**: Can be tested by running `psm sync` directly in a terminal without piping, verifying the prompt works as before.

**Acceptance Scenarios**:

1. **Given** stdin is a terminal (not piped), **When** the user runs `psm sync`, **Then** the approval prompt appears and works identically to the current behavior.

---

### Edge Cases

- What happens when `/dev/tty` exists but is not accessible (permission denied)?
- What happens when the user pipes input and the piped data is empty?
- What happens when the terminal is available but the user sends EOF (Ctrl+D) at the prompt?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST read the approval prompt response from `/dev/tty` when stdin is not a terminal, following the standard Unix pattern used by tools like `rm -i` and `git add -p`.
- **FR-002**: System MUST auto-cancel the operation with a clear message when `/dev/tty` is not available and `--skip-approve` is not set.
- **FR-003**: System MUST preserve existing behavior when stdin is a terminal (non-piped usage).
- **FR-004**: System MUST preserve existing `--skip-approve` behavior regardless of input source.
- **FR-005**: System MUST display a clear error or informational message when auto-cancelling due to unavailable terminal, so the user understands why the operation was cancelled.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can successfully run `sops -d ... | psm sync ...` and respond to the approval prompt interactively.
- **SC-002**: Non-interactive environments without `--skip-approve` auto-cancel with a descriptive message rather than silently doing nothing.
- **SC-003**: All existing tests continue to pass without modification (backward compatibility).
- **SC-004**: The approval prompt behavior matches established Unix conventions (`rm -i`, `git add -p` patterns).

## Assumptions

- The target platform supports `/dev/tty` (Unix/Linux/macOS). Windows compatibility is not in scope.
- The existing `--skip-approve` flag behavior is well-defined and does not change.
- When stdin is a terminal, the system continues to read from stdin as before (no unnecessary `/dev/tty` open).
