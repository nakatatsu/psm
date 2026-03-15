# Survey: GitHub Actions CI

**Date**: 2026-03-15
**Spec**: specs/003-gha-ci/spec.md

## Summary

The spec's direction (GHA CI workflow) is correct — branch protection with required status check `ci` is already configured, so the infrastructure decision is made. However, the spec's proposed 11-step pipeline has significant redundancy: golangci-lint v2 already subsumes `go vet`, `staticcheck`, and `gosec`, making 3 steps and 2 tool installs unnecessary. Additionally, `govulncheck` as a blocking CI step is a reliability risk (external dependency on vuln.go.dev). The workflow can be simplified from 11 steps to ~7, installing 1-2 tools instead of 5.

## A. Problem Reframing

### S1: Problem Definition

**Category**: Problem Definition
**Finding**: The spec frames the problem as "run mechanical checks in GHA." The actual goal is "prevent bad code from being merged." GHA is the right mechanism (branch protection is already in place), but the real question is what goes inside the workflow — not whether to have one.
**Recommendation**: The problem is correctly identified. No reframing needed.
**Evidence**: Branch protection confirmed — required status check `ci` is set, direct push is blocked.

### S2: Hidden Assumptions

**Category**: Hidden Assumptions
**Finding**: The spec assumes all tools must be installed and run separately. In reality:

1. **golangci-lint v2 subsumes go vet, staticcheck, and gosec.** The `.golangci.yml` already enables gosec. golangci-lint runs `govet` by default. staticcheck can be enabled as a linter in golangci-lint. Running these three as separate steps is fully redundant.

2. **"All tool versions must match DevContainer"** — Only `gofumpt` version parity actually matters (format check would produce false failures otherwise). For linters, minor version differences are acceptable. For govulncheck, you want the latest vulnerability database anyway.

3. **"govulncheck should block PRs"** — A vulnerability in a dependency is not caused by the PR author. It is a pre-existing condition better handled by a scheduled scan, not a PR blocker. Also, vuln.go.dev is an external dependency that can timeout.

**Recommendation**: Eliminate redundant steps. Reduce version-sync surface from 8 tools to 2-3.
**Evidence**: `.golangci.yml` enables gosec. golangci-lint v2 docs confirm govet runs by default.

### S3: Lateral Thinking

**Category**: Lateral Thinking
**Finding**: Two alternative approaches considered:

1. **Use DevContainer image as CI runner** — The Dockerfile has all tools pre-installed. Could push to GHCR and use `container:` in workflow. However, the image is bloated (~1-2GB with Claude Code, AWS CLI, zsh, vim, etc.) and couples CI to the dev environment. A dev convenience change would trigger CI image rebuild. Bad coupling direction. **Rejected.**

2. **Use `golangci-lint-action` (first-party GHA action)** — Handles its own binary download and caching. Eliminates the need to install staticcheck, gosec, or run go vet separately. This is the right approach for CI — purpose-built action vs. manual `go install`. **Recommended.**

**Recommendation**: Use `golangci-lint-action` instead of manual `go install` for golangci-lint.
**Evidence**: golangci-lint/golangci-lint-action is the official GHA action, widely adopted.

## B. Solution Evaluation

### S4: Redundancy in Check Steps

**Category**: Approach Alternatives
**Finding**: The spec lists 7 check steps (FR-003 through FR-010). Of these:
- `go vet` (FR-006) — redundant, golangci-lint runs govet by default
- `staticcheck` (FR-007) — redundant if enabled in golangci-lint
- `gosec` — already covered by golangci-lint via `.golangci.yml`

Eliminating these removes 3 steps and 2 tool installs (staticcheck, separate gosec).

**Recommendation**: Update spec to consolidate: golangci-lint covers vet + staticcheck + gosec. Keep gofumpt/goimports as separate checks (they are formatters, not linters).
**Evidence**: golangci-lint v2 default linters include govet. staticcheck is available as a linter.

### S5: govulncheck Reliability

**Category**: Risk & Failure Modes
**Finding**: govulncheck queries vuln.go.dev on every run. If vuln.go.dev is slow or unavailable, the CI run fails and merges are blocked. This is a reliability risk for something unrelated to the PR's code quality.
**Recommendation**: Either run govulncheck with `continue-on-error: true`, or move it to a separate scheduled workflow (e.g., weekly cron). Dependency vulnerabilities are a pre-existing condition, not a PR-introduced problem.
**Evidence**: vuln.go.dev is an external service with no SLA for CI usage.

### S6: Version Sync Maintenance

**Category**: Cost & Complexity
**Finding**: The Dockerfile pins 8 tool versions. If CI independently pins the same 8, that's a maintenance burden and DRY violation. With the proposed simplification (golangci-lint subsumes 3 tools), CI only needs to sync: Go version, gofumpt version, and golangci-lint version — 3 instead of 8.
**Recommendation**: Accept the reduced duplication (3 versions). A `.tool-versions` file is overkill for a solo project.
**Evidence**: Dockerfile analysis.

### S7: goimports — Needed?

**Category**: Necessity
**Finding**: `goimports` checks import ordering and unused imports. golangci-lint has `goimports` as an available linter. If enabled in `.golangci.yml`, the separate install and step become redundant.
**Recommendation**: Verify during plan/research whether golangci-lint's goimports linter is sufficient, or if the standalone tool is still needed for formatting.
**Evidence**: Needs verification in plan Phase 0.

## Items Requiring PoC

None. All findings are based on documented tool behavior.

## Items Requiring Verification in Plan Phase 0

1. Can golangci-lint v2 run staticcheck as an internal linter? (likely yes, but confirm)
2. Can golangci-lint v2 replace standalone goimports check? (may differ in behavior)

## Constitution Impact

No amendments required. CI is explicitly mentioned in the constitution's Development Workflow section: "Run `go test ./...` and `go vet ./...` before every commit."

## Recommendation

Proceed to `/speckit.plan`. Update spec to reflect the simplified workflow before or during planning:

- Remove standalone go vet, staticcheck, gosec steps (subsumed by golangci-lint)
- Add golangci-lint staticcheck linter to `.golangci.yml`
- Consider moving govulncheck to a scheduled workflow or `continue-on-error`
- Reduce version-sync surface to Go + gofumpt + golangci-lint
