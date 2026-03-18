# Feature Specification: Add Approve Flow and Replace --prune with --delete

**Feature Branch**: `018-approve-delete-flow`
**Created**: 2026-03-17
**Status**: Draft
**Input**: GitHub Issue #18 — Add approve flow and replace --prune with --delete

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Operator reviews changes before execution (Priority: P1)

An operator runs `psm sync` to apply parameter changes to AWS. Before any changes are executed, the tool displays a full action plan (creates, updates, deletes) and asks for confirmation. The operator reviews the plan and confirms or cancels. This prevents accidental changes from going through unreviewed.

**Why this priority**: Safety is the core motivation for this feature. Without approval, any sync mistake is immediately applied with no undo.

**Independent Test**: Run `psm sync` with a YAML file that triggers creates and updates. Verify the action plan is displayed and execution only proceeds after typing `y`. Verify typing `N` or pressing Enter cancels without changes.

**Acceptance Scenarios**:

1. **Given** a YAML file with changes pending, **When** the operator runs `psm sync`, **Then** the full action plan is displayed and a `y/N` prompt appears before any changes are executed.
2. **Given** the approve prompt is shown, **When** the operator types `y`, **Then** all planned changes are executed.
3. **Given** the approve prompt is shown, **When** the operator types `N` or presses Enter (default), **Then** no changes are executed and the tool exits with code 0.
4. **Given** the operator passes `--skip-approve`, **When** `psm sync` runs, **Then** changes execute immediately without a prompt (for CI/automation).
5. **Given** the operator passes `--dry-run`, **When** `psm sync` runs, **Then** the action plan is displayed but no prompt is shown and no changes are executed.

---

### User Story 2 - Operator deletes obsolete keys using pattern file (Priority: P1)

An operator wants to clean up obsolete AWS parameters/secrets. They create a YAML file listing regex patterns of keys to delete, and pass it via `--delete`. The tool identifies matching keys in AWS, cross-checks against the sync YAML to prevent conflicts, and deletes them after approval.

**Why this priority**: Equal to US1 — replacing the dangerous `--prune` flag is the other core motivation. Without scoped deletion, operators risk deleting keys managed by other teams.

**Independent Test**: Create a delete patterns file targeting specific key prefixes. Run `psm sync --delete needless.yml params.yml` and verify only matching keys are proposed for deletion, conflicts are detected, and unmanaged keys trigger warnings.

**Acceptance Scenarios**:

1. **Given** a delete file with regex patterns and AWS contains matching keys not in sync YAML, **When** the operator runs `psm sync --delete needless.yml params.yml`, **Then** only keys matching the patterns AND not in the sync YAML are proposed for deletion.
2. **Given** a delete file where a matching key also exists in the sync YAML, **When** the operator runs the command, **Then** the entire operation aborts with an error before any changes are executed (all-or-nothing).
3. **Given** a delete file is provided, **When** AWS contains keys not in sync YAML and not matching any delete pattern, **Then** those keys are displayed as warnings (unmanaged keys).
4. **Given** `--delete` is used without a sync file, **When** the operator runs the command, **Then** the tool returns an error requiring the sync file.

---

### User Story 3 - Operator enables debug logging for troubleshooting (Priority: P3)

An operator encounters unexpected behavior and wants to see detailed diagnostic output. They pass `--debug` to enable Debug-level logging, which shows API call details, regex match results, and other internal operations.

**Why this priority**: Useful for troubleshooting but not core functionality. Normal operations work without it.

**Independent Test**: Run `psm sync --debug` and verify Debug-level messages appear in stderr output. Run without `--debug` and verify Debug messages are hidden.

**Acceptance Scenarios**:

1. **Given** the operator passes `--debug`, **When** `psm sync` runs, **Then** Debug-level slog messages appear in stderr (e.g., API call details, regex match results).
2. **Given** the operator does not pass `--debug`, **When** `psm sync` runs, **Then** only Error, Warn, and Info messages appear.

---

### User Story 4 - Remove --prune flag (Priority: P1)

The existing `--prune` flag deletes all AWS parameters/secrets not in the YAML file, with no scope restriction. This is too dangerous and must be removed entirely. Users who relied on `--prune` must migrate to `--delete <file>`.

**Why this priority**: Removing a dangerous feature is a safety prerequisite.

**Acceptance Scenarios**:

1. **Given** a user passes `--prune`, **When** `psm sync` runs, **Then** the tool returns an error indicating `--prune` has been removed and suggests using `--delete <file>` instead.

---

### Edge Cases

- What happens when the delete file contains an invalid regex pattern?
- What happens when the delete file is empty (no patterns)?
- What happens when all AWS keys match delete patterns (mass deletion)?
- What happens when stdin is not a terminal (piped input) and `--skip-approve` is not set?
- What happens when the approve prompt receives unexpected input (not `y` or `N`)?
- What happens when `--delete` and `--dry-run` are combined?

## Requirements *(mandatory)*

### Functional Requirements

#### Approve Flow

- **FR-001**: Before executing any changes (create, update, delete), the tool MUST display the full action plan and prompt the user with `y/N` for confirmation. Default MUST be `N` (no changes).
- **FR-002**: `--skip-approve` flag MUST bypass the confirmation prompt and execute immediately.
- **FR-003**: `--dry-run` MUST display the action plan without showing a prompt and without executing changes.
- **FR-004**: When the user declines (types `N` or presses Enter), the tool MUST exit with code 0 and make no changes.
- **FR-015**: When no changes are needed (all keys are unchanged), the tool MUST display the summary only and skip the approve prompt.

#### Delete with Pattern File

- **FR-005**: `--prune` flag MUST be removed. If passed, the tool MUST return an error with a migration message suggesting `--delete <file>`.
- **FR-006**: `--delete <file>` flag MUST accept a YAML file containing a list of regex patterns.
- **FR-007**: Deletion candidates MUST be AWS keys that match any regex in the delete file AND are not present in the sync YAML.
- **FR-008**: The sync file MUST be required when `--delete` is used.
- **FR-009**: Before executing any operation (create, update, or delete), the tool MUST validate all deletion candidates against the sync YAML. If any candidate also exists in the sync YAML, the entire operation MUST abort with an error — no partial execution.
- **FR-010**: Invalid regex patterns in the delete file MUST cause an immediate error before any AWS operations.

#### Unmanaged Key Warnings

- **FR-011**: When `--delete` is specified, AWS keys that are not in the sync YAML and do not match any delete pattern MUST be displayed as warnings.
- **FR-012**: Unmanaged key warnings MUST NOT block execution.

#### Debug Logging

- **FR-013**: `--debug` flag MUST enable Debug-level slog output to stderr.
- **FR-014**: Default log level MUST be Info (Error + Warn + Info visible; Debug hidden).

### Key Entities

- **Action Plan**: The computed list of operations (create, update, delete, skip) with key names, displayed to the user before approval.
- **Delete Pattern File**: A YAML file containing a list of regex strings used to match AWS keys for deletion.
- **Deletion Candidate**: An AWS key that matches a delete regex AND is not present in the sync YAML.
- **Unmanaged Key**: An AWS key that is not in the sync YAML and does not match any delete regex.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: No changes are ever applied without explicit user approval (via prompt or `--skip-approve`).
- **SC-002**: `--prune` is no longer accepted — the tool returns an error with migration guidance.
- **SC-003**: Conflict detection catches 100% of cases where a deletion candidate exists in the sync YAML, aborting before any changes.
- **SC-004**: Unmanaged keys are surfaced as warnings, enabling operators to detect forgotten keys.
- **SC-005**: Debug logging provides sufficient detail to troubleshoot regex matching and API interactions without reading source code.

## Clarifications

### Session 2026-03-17

- Q: When there are no changes (all skip), should the approve prompt still appear? -> A: No — skip the prompt and just show the summary.

## Assumptions

- The approve prompt writes to stderr (not stdout) per the logging policy, so it does not interfere with piped output.
- When stdin is not a terminal and `--skip-approve` is not set, the tool should treat it as if the user declined (exit 0, no changes). This prevents accidental execution in non-interactive environments.
- The delete file format is a simple YAML list of strings (regex patterns). No metadata, no nesting.
- `--delete` can be used together with sync (create/update) in the same invocation. Documentation should present them as separate operations to encourage safer usage patterns.
- The `--debug` flag applies to all subcommands (sync, export), not just sync.
