# Feature Specification: Version Display Command

**Feature Branch**: `030-version-display`
**Created**: 2026-03-25
**Status**: Draft
**Input**: GitHub Issue #30: feat: add version display command

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Check installed version (Priority: P1)

A user wants to confirm which version of psm they are running, for example when reporting a bug or verifying that an upgrade was applied correctly. They run `psm --version` and see the version number printed to the terminal.

**Why this priority**: This is the core feature — without it, there is no way to identify the installed version.

**Independent Test**: Can be fully tested by building the binary with a version string and running `psm --version`, verifying the output matches the expected format.

**Acceptance Scenarios**:

1. **Given** psm was built with version `1.2.3`, **When** the user runs `psm --version`, **Then** the output is `psm version 1.2.3`
2. **Given** psm was built without a version string (e.g., via plain `go build`), **When** the user runs `psm --version`, **Then** the output is `psm version dev`

---

### User Story 2 - Version check in CI/CD pipelines (Priority: P2)

A CI/CD pipeline needs to programmatically verify the psm version before running parameter sync operations. The version output is a single line of predictable text that can be parsed by scripts.

**Why this priority**: Enables automated version verification, but depends on P1 being implemented first.

**Independent Test**: Can be tested by piping `psm --version` output to a script that parses the version string and asserts it matches the expected format.

**Acceptance Scenarios**:

1. **Given** psm is installed in a CI environment, **When** a script runs `psm --version`, **Then** the output is a single line matching the pattern `psm version <semver-or-dev>`
2. **Given** psm is invoked with `--version`, **When** the exit code is checked, **Then** it is `0`

---

### Edge Cases

- What happens when `--version` is combined with a subcommand (e.g., `psm sync --version`)? The version flag should only be recognized as the first argument (i.e., `psm --version`), not within subcommands.
- What happens when `--version` is combined with other flags (e.g., `psm --version --debug`)? Extra flags after `--version` should be ignored; the version is printed and the program exits.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST display the version string when invoked with `--version` as the first argument
- **FR-002**: System MUST output the version in the format `psm version <version-string>` followed by a newline
- **FR-003**: System MUST exit with code `0` after displaying the version
- **FR-004**: System MUST display `dev` as the version when no version was embedded at build time
- **FR-005**: System MUST support version embedding via build-time variable injection (e.g., linker flags)

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can determine their installed psm version in a single command invocation
- **SC-002**: The version output is a single, predictable line that can be reliably parsed by scripts
- **SC-003**: All existing commands and flags continue to work without regression after adding version support

## Assumptions

- The version string follows semantic versioning (e.g., `1.2.3`) but validation of the format is not enforced — whatever value is injected at build time is displayed as-is
- `--version` is handled before subcommand parsing, so it does not conflict with existing `sync` or `export` subcommands
