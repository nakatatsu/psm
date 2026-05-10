# workspace Development Guidelines

## Important

- Always save `.claude/` files (skills, settings, etc.) in the repository (`/workspace/.claude/`), never in `~/.claude/`.

### Pre-commit gate (MANDATORY)

`.githooks/pre-commit` runs `gofumpt` + `golangci-lint` + `go test -race ./...` whenever staged changes touch `*.go` / `go.mod` / `go.sum`. Setup is `git config core.hooksPath .githooks` (already configured).

- **Never bypass with `--no-verify`** unless the user has explicitly authorised it for that single commit. Pushing untested commits and waiting for CI is treated as a defect, not a tradeoff.
- If the hook fails, fix the cause and create a new commit. Do not amend around it.
- If you change a runtime / dependency / tool / build config (e.g. `Dockerfile`, workflow `GO_VERSION`) without touching `*.go`, the hook will not trigger automatically — run `/check` (or `/check all` for vulnerability scans) manually before committing in that case.

### Key Skills (always invoke when conditions match)

- `/gh-token` — GitHub token retrieval via sidecar. Use before any Git remote operation or on auth errors using `gh`.
- `/branch-strategy` — Branching strategy. Use when creating branches, merging, or committing.
- `/coding-workflow` — Coding workflow guide. Use when deciding what to do next or how to proceed.

## Active Technologies

- Go
- SOPS
- age
- AWS CLI
