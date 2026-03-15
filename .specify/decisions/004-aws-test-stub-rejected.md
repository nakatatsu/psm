# ADR-004: AWS Test Stub Infrastructure — Rejected

**Date**: 2026-03-15
**Status**: Rejected

## Context

With CI introduction (003-gha-ci), 9 out of 19 tests are skipped because they require real AWS credentials (SSM Parameter Store and Secrets Manager integration tests). The proposal was to add moto (Python AWS emulator) as a Docker service so these tests could run against a stub in CI. Constitution v2.0.0 was amended to v3.0.0 to permit emulator-based stub tests.

A survey (specs/004-aws-test-stub/survey.md) reframed the problem: the actual goal is "detect regressions in CI," not "run AWS integration tests without credentials." These are not the same thing.

Key findings from the survey:

- The 10 non-AWS tests already cover YAML parsing, plan/diff logic, execute (via fakeStore), CLI parsing, and export — the layers where bugs are most likely to be introduced.
- The AWS Store implementations (ssm.go ~75 lines, sm.go ~93 lines) are thin SDK wrappers with minimal logic (SM create-or-update branching and SSM batch delete only).
- 001-psm research R5 concluded: "Mocking the Store interface just skips a thin SDK wrapper layer — there's no test value." This applies equally to emulators.
- Moto itself has had API fidelity regressions (e.g., moto#9700). For a tool handling customer secrets, emulator tests passing while real AWS fails is unacceptable.

## Decision

Rejected. Do not introduce moto or any AWS emulator. CI runs non-AWS tests only (9 AWS tests skipped via `PSM_INTEGRATION_TEST` guard). Real AWS integration tests remain as a local/manual pre-release gate.

## Consequences

- 003-gha-ci proceeds with AWS tests skipped — CI covers application logic, not SDK wrapper behavior
- Real AWS tests remain available locally (`PSM_INTEGRATION_TEST=1 PSM_TEST_PROFILE=psm go test -v ./...`)
- Constitution v3.0.0 stub tests clause is retained for future use if Store implementations grow more complex
- SDK interface fakes or GitHub Actions OIDC federation can be revisited if the need arises
