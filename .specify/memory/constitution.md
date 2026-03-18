# psm Constitution

## Core Principles

### I. Simplicity First

- Every design decision defaults to the simplest viable option.
- No abstractions until a concrete, present need demands one. Three duplicated lines are preferable to a premature helper.
- Flat package structure unless a package boundary is forced by a clear dependency cycle or API contract.
- If a feature can be achieved with the standard library, do not introduce a third-party dependency.
- Code must be readable without comments explaining _what_ it does; comments are reserved for _why_.

### II. YAGNI (You Aren't Gonna Need It)

- Do not build for hypothetical future requirements. Implement only what is needed today.
- No feature flags, configuration toggles, or plugin systems unless the current task explicitly requires them.
- No backward-compatibility shims. When something changes, change it directly.
- No speculative error handling for scenarios that cannot occur in the current design.
- If in doubt whether something is needed, leave it out.

### III. Test-First (NON-NEGOTIABLE)

- **Test-first is mandatory**: Write test → Confirm test fails → Implement → Confirm test passes (Red-Green cycle). This order must never be violated.
- All tests use `go test` from the Go standard library. No third-party test frameworks, assertion libraries, or test runners.
- Test files live alongside the code they test (`*_test.go` in the same package).
- Table-driven tests are the preferred pattern where multiple inputs/outputs are exercised.
- Unit tests (pure logic) must be deterministic and fast with no external dependencies.
- Stub tests (for CI) use an AWS emulator (e.g., moto) running as a Docker service. Endpoint URL is injected via environment variable. Stub tests exercise the same Store implementations as real AWS but against the emulator. When behavioral differences between the stub and real AWS are discovered, real AWS is authoritative; differences are documented as known limitations.
- Integration tests (real AWS) use a sandbox environment with dedicated test prefixes and setup/teardown. Real AWS is the authoritative source of truth.
- Test branches in production code (e.g., `if test then ...`) are prohibited. Use dependency injection to switch between stub and real AWS.

## Technology Stack

- **Language**: Go (latest stable release)
- **Project type**: CLI tool
- **Dependencies**: Standard library preferred. Third-party dependencies require justification against Principle I.
- **Build**: `go build`
- **Test**: `go test ./...`
- **Lint**: `go vet` and `staticcheck` (or `golangci-lint` with minimal config)
- **Formatting**: `gofmt` (non-negotiable, enforced by CI)

## Development Workflow

- Keep the main branch always in a buildable, passing-tests state.
- Each change should be a single, focused commit or PR that addresses one concern.
- Run `go test ./...` and `go vet ./...` before every commit.
- No generated code unless mandated by an external protocol (e.g., protobuf). If generated code is used, the generator invocation must be documented and reproducible.
- Error handling follows Go conventions: return `error`, wrap with `fmt.Errorf("context: %w", err)`. Do not use `panic`.

### Implementation Order (NON-NEGOTIABLE)

1. **Write test**: Write tests for the target function/method first. The signature (name, parameters, return values) is finalized at this point.
2. **Confirm test fails**: Run `go test` and confirm the test either fails to compile or fails assertions.
3. **Implement**: Write the minimum code necessary to make the test pass.
4. **Confirm test passes**: Run `go test` and confirm all tests pass.
5. **When tests cannot be written without scaffolding**: Create only the interface or type definitions (with empty implementations), write the test, then fill in the implementation.

Commits that violate this order (implementation without tests, proceeding to next feature before tests pass) are not permitted.

## Governance

- This constitution is the highest authority on project practices. All specs, plans, and code reviews must comply.
- Amendments require: (1) a written proposal stating the change and rationale, (2) an update to this file with version bump, (3) a review of dependent templates for consistency.
- Version follows semantic versioning: MAJOR for principle removals or redefinitions, MINOR for new principles or sections, PATCH for wording clarifications.
- Compliance is verified during spec review (checklist gate) and plan review (Constitution Check section).

**Version**: 3.0.1 | **Ratified**: 2026-03-08 | **Last Amended**: 2026-03-17
