# Implementation Plan: GitHub Actions CI

**Branch**: `003-gha-ci` | **Date**: 2026-03-15 | **Spec**: specs/003-gha-ci/spec.md
**Input**: Feature specification from `/specs/003-gha-ci/spec.md`

## Summary

PR and main-branch push に対して機械的チェックを自動実行する GHA ワークフローを構築する。survey の指摘を受けてツールの冗長を整理し、golangci-lint v2 に govet/staticcheck/gosec/goimports を集約。CI ステップ数は 8（うち govulncheck は non-blocking）。

## Technical Context

**Language/Version**: Go 1.26.1
**Primary Dependencies**: GitHub Actions (actions/checkout@v4, actions/setup-go@v5, golangci/golangci-lint-action@v9)
**Storage**: N/A
**Testing**: `go test -race ./...` (10 non-AWS tests run, 9 AWS tests skipped)
**Target Platform**: GitHub-hosted runner (ubuntu-latest)
**Project Type**: CI/CD configuration (YAML workflow file)
**Performance Goals**: CI completes within 5 minutes
**Constraints**: No AWS credentials in CI. No stored secrets for external services.
**Scale/Scope**: Single workflow file, single job

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| # | Principle | Gate Question | Pass? |
|---|-----------|---------------|-------|
| I | Simplicity First | Is this the simplest viable design? No unnecessary abstractions or dependencies? | Yes — single workflow, single job, redundant tools eliminated |
| II | YAGNI | Does every element serve a present, concrete need? No speculative features? | Yes — branch protection already requires `ci` status check |
| III | Test-First (NON-NEGOTIABLE) | Are all tests written before implementation (Red-Green cycle)? Using `go test` only? No third-party test frameworks? Stub/integration tests use DI (no test branches in production code)? | N/A — this feature is a YAML workflow file, not Go code. Validation is done by running the workflow against the existing test suite. |

## Project Structure

### Documentation (this feature)

```text
specs/003-gha-ci/
├── spec.md
├── survey.md
├── plan.md              # This file
├── research.md
└── checklists/
    └── requirements.md
```

### Source Code (repository root)

```text
.github/
└── workflows/
    └── ci.yml           # The single deliverable

.golangci.yml            # Modified: add goimports formatter
```

**Structure Decision**: Single YAML file under `.github/workflows/`. No additional source code. `.golangci.yml` is updated to add `goimports` as a formatter (research R2).

## CI Workflow Design

### Job: `ci`

```yaml
name: CI
on:
  pull_request:
    branches: [main]
  push:
    branches: [main]

jobs:
  ci:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.26.1'

      - name: Install gofumpt
        run: go install mvdan.cc/gofumpt@v0.9.2

      - name: Build
        run: go build ./...

      - name: Format check (gofumpt)
        run: |
          test -z "$(gofumpt -l .)"

      - name: Lint
        uses: golangci/golangci-lint-action@v9
        with:
          version: v2.9.0

      - name: Test
        run: go test -race ./...

      - name: Vulnerability check
        if: always()
        continue-on-error: true
        run: |
          go install golang.org/x/vuln/cmd/govulncheck@v1.1.4
          govulncheck ./...
```

### Tool Coverage Matrix

| Check | Tool in CI | Previously separate? |
|-------|-----------|---------------------|
| Build | `go build ./...` | Same |
| Format (gofumpt) | standalone `gofumpt` | Same |
| Format (goimports) | golangci-lint formatter | Was standalone |
| Static analysis (go vet) | golangci-lint (govet) | Was standalone |
| Static analysis (staticcheck) | golangci-lint (staticcheck) | Was standalone |
| Security (gosec) | golangci-lint (gosec) | Was via golangci-lint (unchanged) |
| Vulnerability | govulncheck (non-blocking) | Was blocking |
| Unit tests | `go test -race` | Same |

### `.golangci.yml` Update

```yaml
version: "2"

formatters:
  enable:
    - goimports

linters:
  enable:
    - gosec

linters-settings:
  gosec:
    excludes:
      - G703
      - G705
```

## Complexity Tracking

No constitution violations. No complexity justification needed.
