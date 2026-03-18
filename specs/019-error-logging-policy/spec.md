# Feature Specification: Define Error Handling and Logging Policy

**Feature Branch**: `019-error-logging-policy`
**Created**: 2026-03-17
**Status**: Draft
**Input**: GitHub Issue #19 — Define error handling and logging policy

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Developer writes new code with consistent error handling (Priority: P1)

A developer adding a new feature or fixing a bug needs clear, documented rules for how to handle errors and produce log output. Without a policy, each developer (or AI agent) invents their own conventions, leading to inconsistent error messages, missing context, and difficulty debugging production issues.

**Why this priority**: Error handling consistency is the foundation — without it, logging rules have no errors to log.

**Independent Test**: Can be verified by reviewing the policy document and confirming it covers all error handling conventions (return patterns, wrapping, panic prohibition, exit codes).

**Acceptance Scenarios**:

1. **Given** a developer reads the policy, **When** they write a function that can fail, **Then** they know to return `(..., error)` and how to wrap errors with context using `fmt.Errorf` and `%w`.
2. **Given** a developer reads the policy, **When** they encounter an error in a CLI entry point, **Then** they know to exit with code 1 and what to output to stderr.
3. **Given** a developer reads the policy, **When** they consider using `panic`, **Then** the policy clearly states it is prohibited and explains the alternative (return error).

---

### User Story 2 - Developer adds logging to a feature (Priority: P1)

A developer implementing a feature (e.g., the approve flow in #18) needs to know which logger to use, what log level to assign, and where output should be routed (stdout vs stderr). The policy removes guesswork and ensures all log output is consistent and machine-parseable.

**Why this priority**: Equal to US1 — logging and error handling are two sides of the same coin, and #18 is blocked on both.

**Independent Test**: Can be verified by reviewing the policy document and confirming it covers logger choice, log levels with usage examples, output routing rules, and sensitive data prohibitions.

**Acceptance Scenarios**:

1. **Given** a developer reads the policy, **When** they need to log an operational message, **Then** they know to use `slog` at the appropriate level and that log output goes to stderr.
2. **Given** a developer reads the policy, **When** they need to choose between Info and Warning, **Then** the level guidelines table provides clear, unambiguous criteria with examples.
3. **Given** a developer reads the policy, **When** they log an error involving an AWS parameter, **Then** the policy reminds them to never include secret values in log messages.

---

### User Story 3 - Developer distinguishes program output from log output (Priority: P2)

A developer working on CLI output (diff lines, summaries, approve prompts) needs to understand the boundary between "program output" (stdout) and "log messages" (stderr via slog). Mixing these breaks piping and scripting.

**Why this priority**: Important for CLI tool correctness but secondary to establishing the core conventions.

**Independent Test**: Can be verified by confirming the policy defines clear routing rules for stdout (program results only) vs stderr (all logs and interactive prompts).

**Acceptance Scenarios**:

1. **Given** a developer reads the policy, **When** they add a new output line (e.g., a diff line), **Then** they know it belongs on stdout as program output, not as a log message.
2. **Given** a developer reads the policy, **When** they add an interactive prompt (e.g., approve `y/N`), **Then** they know it belongs on stderr so it doesn't pollute piped output.

---

### Edge Cases

- What if a function has multiple error paths — does the policy specify wrapping at each level or only at boundaries?
- How should errors be handled when calling external services (AWS APIs) that return structured errors — wrap or pass through?
- What happens when a log message would need to include a key name that itself is sensitive (e.g., `/prod/secrets/db-password`) — is the key path considered sensitive?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The policy document MUST specify `log/slog` as the only permitted logger; all other logging methods (e.g., `log.Println`, `fmt.Fprintf(stderr, ...)` for log-like messages) MUST be prohibited.
- **FR-002**: The policy document MUST define four log levels (Error, Warning, Info, Debug) with usage guidelines and concrete examples for each level.
- **FR-003**: The policy document MUST specify that the default log level is Info (Error + Warning + Info visible; Debug hidden).
- **FR-004**: The policy document MUST define output routing: program output (diff lines, summaries) to stdout; all log messages and interactive prompts to stderr.
- **FR-005**: The policy document MUST prohibit logging of sensitive data (secret values, API keys, passwords, tokens). Key paths/names are permitted.
- **FR-006**: The policy document MUST specify that functions which can fail return `(..., error)` as the last return value.
- **FR-007**: The policy document MUST prohibit the use of `panic`.
- **FR-008**: The policy document MUST define error wrapping conventions using `fmt.Errorf` with `%w` verb, including guidance on adding context at each call site.
- **FR-009**: The policy document MUST define exit code conventions: 0 for success, 1 for failure.
- **FR-010**: The policy document MUST be written in the file `documents/backend/error-handling-and-logging.md`.
- **FR-011**: This issue MUST NOT include any code changes — implementation follows in a separate issue.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A developer reading only the policy document can determine the correct log level for any given message within 30 seconds, using the level guidelines table.
- **SC-002**: The policy document covers 100% of the conventions listed in the Issue #19 acceptance criteria (logger choice, log levels, output routing, sensitive data, return conventions, panic, wrapping, exit codes).
- **SC-003**: No ambiguity remains — every rule uses MUST/MUST NOT language and includes at least one concrete example.

## Assumptions

- This project is a CLI tool, not a web server. The policy is tailored accordingly (no HTTP status codes, no request-scoped logging).
- The policy applies to the `psm` codebase only. `documents/backend/coding-standard.md` is deleted as redundant (its content is a subset of this policy document).
- Key paths (e.g., `/myapp/prod/API_KEY`) are not considered sensitive data — only the values stored at those keys are sensitive.
- The policy will be enforced by code review and AI agent compliance, not by automated tooling (linters, etc.) in this issue.
