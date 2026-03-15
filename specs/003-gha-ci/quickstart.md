# Quickstart: GitHub Actions CI

## What it does

PR to main or push to main triggers automated checks:
build → format (gofumpt) → lint (golangci-lint: govet + staticcheck + gosec + goimports) → test (go test -race) → vulnerability (govulncheck, non-blocking)

## Files

- `.github/workflows/ci.yml` — CI workflow
- `.golangci.yml` — Updated to add goimports formatter

## Validation

1. Create a PR to main → CI should run and pass
2. Push a gofumpt violation → CI should fail at format check
3. Push a failing test → CI should fail at test step
