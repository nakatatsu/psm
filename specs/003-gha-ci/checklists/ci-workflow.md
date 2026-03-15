# CI Workflow Requirements Checklist: GitHub Actions CI

**Purpose**: Validate that CI workflow requirements are complete, clear, and consistent before implementation
**Created**: 2026-03-15
**Feature**: specs/003-gha-ci/spec.md

## Requirement Completeness

- [ ] CHK001 Are all check categories (build, format, lint, vulnerability, test) explicitly listed with their expected behavior on failure? [Completeness, Spec FR-003~007]
- [ ] CHK002 Is the job name `ci` explicitly specified as a requirement, with the reason (required status check match) documented? [Completeness, Spec FR-011]
- [ ] CHK003 Are the trigger events (pull_request, push) and their branch filters specified? [Completeness, Spec FR-001~002]
- [ ] CHK004 Is the runner environment (ubuntu-latest) specified? [Completeness, Spec Assumptions]
- [ ] CHK005 Is the AWS integration test exclusion mechanism documented (PSM_INTEGRATION_TEST guard)? [Completeness, Spec FR-010]

## Requirement Clarity

- [ ] CHK006 Is it clear which tools are standalone vs. subsumed by golangci-lint? [Clarity, Plan R1]
- [ ] CHK007 Is the govulncheck execution strategy (blocking vs. non-blocking) explicitly defined? [Clarity, Plan R4]
- [ ] CHK008 Is "format violation" defined precisely enough to distinguish gofumpt failures from goimports failures? [Clarity, Spec FR-004]
- [ ] CHK009 Are version pinning requirements clear — which tools need exact version match and which don't? [Clarity, Spec FR-008]

## Requirement Consistency

- [ ] CHK010 Are the tool choices in plan.md consistent with the spec's functional requirements? [Consistency, Plan vs Spec]
- [ ] CHK011 Is the `.golangci.yml` update (adding goimports formatter) documented in both plan.md and the workflow design? [Consistency, Plan R2]
- [ ] CHK012 Are the survey findings (tool redundancy elimination) reflected in the final spec requirements? [Consistency, Survey vs Spec]

## Scenario Coverage

- [ ] CHK013 Is the behavior defined for when a CI step fails mid-pipeline — do subsequent steps still run? [Scenario Coverage, Gap]
- [ ] CHK014 Is the behavior defined for when govulncheck's external service (vuln.go.dev) is unavailable? [Scenario Coverage, Plan R4]
- [ ] CHK015 Is the behavior defined for when `go install gofumpt` fails (network issue, version not found)? [Scenario Coverage, Gap]
- [ ] CHK016 Are concurrent CI runs on the same PR (rapid push) addressed? [Scenario Coverage, Gap]

## Edge Case Coverage

- [ ] CHK017 Does the spec define behavior when the workflow YAML file itself is modified in a PR? [Edge Case, Spec Edge Cases]
- [ ] CHK018 Does the spec define behavior when no test files exist? [Edge Case, Spec Edge Cases]
- [ ] CHK019 Does the spec address what happens when a PR has no Go source changes (docs-only, YAML-only)? [Edge Case, Gap]

## Dependencies and Assumptions

- [ ] CHK020 Is the dependency on branch protection (required status check `ci`) documented as a prerequisite, not a deliverable? [Assumption, Spec Assumptions]
- [ ] CHK021 Is the dependency on GitHub-hosted runner availability acknowledged? [Assumption, Spec Assumptions]
- [ ] CHK022 Is the ADR-004 decision (no AWS stubs in CI) referenced or consistent with FR-010? [Assumption, ADR-004]
