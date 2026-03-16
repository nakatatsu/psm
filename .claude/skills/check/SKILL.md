---
name: check
description: >
  Run mechanical checks (formatter, static analysis, linter, tests) on Go code.
  Invoke with /check. Control scope with arguments.
  Use after writing code, before committing, or before pushing for quality assurance.
  Also use when the user says "check", "lint", "run tests", "does it build", or similar.
---

# check — Go Mechanical Checks

Run formatter, linter, and tests against the project's Go code.

## Usage

| Command        | What it runs                              |
| -------------- | ----------------------------------------- |
| `/check`       | Formatter + linter + unit tests           |
| `/check all`   | Above + govulncheck                       |
| `/check int`   | AWS integration tests only                |
| `/check vuln`  | govulncheck only                          |

## Execution

All commands run in `/workspace`.

### Default (`/check`)

Run the following in order. Continue to the end even if a step fails, then report all results together.

```bash
cd /workspace

# 1. Formatter
gofumpt -l .

# 2. Linter (includes govet, staticcheck, gosec, goimports, errcheck, ineffassign, unused)
golangci-lint run ./...

# 3. Unit tests
go test ./...
```

Report results in this format:

```
## Check Results

| Check         | Result              |
| ------------- | ------------------- |
| gofumpt       | ✅ OK / ❌ Fix needed |
| golangci-lint | ✅ OK / ❌ Fix needed |
| go test       | ✅ OK / ❌ Failed     |
```

If any step fails, show details (file, line, what went wrong) after the table.

### `/check all`

Run the default checks, then additionally:

```bash
govulncheck ./...
```

govulncheck requires external access to `vuln.go.dev`, available via the Squid proxy.

### `/check int`

Run AWS integration tests:

```bash
cd /workspace
PSM_INTEGRATION_TEST=1 PSM_TEST_PROFILE=psm go test -v -timeout 120s ./...
```

Requires AWS SSO authentication. If auth errors occur, prompt the user to run `aws sso login`.

### `/check vuln`

Run vulnerability check only:

```bash
cd /workspace
govulncheck ./...
```

## Environment Notes

- **golangci-lint v2** bundles govet, staticcheck, gosec, goimports, errcheck, ineffassign, unused. Standalone versions of these tools are not needed (removed from Dockerfile).
- **golangci-lint** config is in `.golangci.yml`. gosec rules G703 (path traversal) and G705 (XSS) are excluded (not applicable to a CLI tool).
- **govulncheck** requires external access to `vuln.go.dev`, available via the Squid proxy (outbound-filter).
