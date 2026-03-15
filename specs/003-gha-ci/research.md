# Research: GitHub Actions CI

**Date**: 2026-03-15

## R1: Tool Redundancy — golangci-lint v2 Coverage

**Decision**: golangci-lint v2 (2.9.0) subsumes go vet, staticcheck, and gosec. Run golangci-lint as the single linter step. Do not run go vet or staticcheck separately.

**Rationale**: golangci-lint v2 uses a "standard" default preset that includes `govet`, `staticcheck` (with merged `gosimple` + `stylecheck`), `errcheck`, `ineffassign`, and `unused`. The project's `.golangci.yml` adds `gosec` on top of this default. Running go vet or staticcheck separately is fully redundant.

**Alternatives considered**:
- Run all tools separately (current check.md approach) — Redundant, slower, more tools to install
- Run only golangci-lint without the standard preset — Would lose coverage

## R2: goimports in golangci-lint v2

**Decision**: Use golangci-lint v2's formatter feature for goimports check. Add `goimports` to the `formatters:` section of `.golangci.yml`.

**Rationale**: In golangci-lint v2, goimports moved from `linters:` to `formatters:`. Formatters report violations (non-zero exit code) when run via `golangci-lint run`, making them usable in CI for detection. This eliminates the need to install standalone goimports.

**Configuration**:
```yaml
version: "2"
formatters:
  enable:
    - goimports
linters:
  enable:
    - gosec
```

**Alternatives considered**:
- Standalone goimports install — Extra tool, extra step, redundant
- Skip goimports check — Risk of inconsistent import ordering

## R3: gofumpt in CI

**Decision**: Keep gofumpt as a standalone check. Install via `go install` with pinned version matching Dockerfile.

**Rationale**: gofumpt is a superset of gofmt with stricter rules (blank line enforcement, etc.). golangci-lint v2 has `gofumpt` as a formatter option, but version parity between CI and developer environment is critical for formatters — if versions differ, the check produces false positives. Pinning the standalone version is the safest approach.

**Alternatives considered**:
- Use golangci-lint's built-in gofumpt formatter — Version is tied to golangci-lint release, may not match Dockerfile's pinned version
- Skip gofumpt, use only gofmt — Loses the stricter formatting rules the project already uses

## R4: govulncheck Execution Strategy

**Decision**: Run govulncheck with `continue-on-error: true`. Do not block PR merges on vulnerability findings.

**Rationale**: govulncheck queries vuln.go.dev on every run. This is an external service with no SLA. If it is slow or unavailable, blocking PRs is disproportionate. Dependency vulnerabilities are pre-existing conditions, not introduced by the PR. The check still reports findings visibly in the CI log.

**Alternatives considered**:
- Block PRs on govulncheck failure — External service outage blocks all development
- Move to scheduled workflow only — Loses visibility on PRs entirely
- Skip govulncheck — Loses vulnerability detection

## R5: GitHub Actions Setup

**Decision**: Use `golangci/golangci-lint-action@v9` for golangci-lint. Use `actions/setup-go@v5` for Go with built-in module caching.

**Rationale**: The official golangci-lint action handles binary download, caching, and GitHub annotations. No need for `go install golangci-lint`. `actions/setup-go@v5` provides Go setup with automatic module caching via `go.sum`.

**Tools to install manually**: Only `gofumpt` (via `go install`, pinned version).

**Alternatives considered**:
- Install all tools via `go install` — Slower, no caching, no annotations
- Use DevContainer image as CI runner — Bloated (~1-2GB), bad coupling direction

## R6: CI Workflow Steps (Final)

**Decision**: 6 steps in the CI job:

| Step | Tool | Purpose |
|------|------|---------|
| 1 | `actions/checkout@v4` | Code checkout |
| 2 | `actions/setup-go@v5` | Go 1.26.1 + module cache |
| 3 | `go install gofumpt` | Formatter (pinned version) |
| 4 | `go build ./...` | Compile check |
| 5 | `gofumpt -l . \| diff /dev/null -` | Format violation detection |
| 6 | `golangci/golangci-lint-action@v9` | Linting (govet + staticcheck + gosec + goimports) |
| 7 | `go test -race ./...` | Unit tests + race detection |
| 8 | `govulncheck ./...` (continue-on-error) | Vulnerability check (non-blocking) |

**Alternatives considered**:
- 11 steps with individual tools — Redundant, more maintenance
