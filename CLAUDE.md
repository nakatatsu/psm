# workspace Development Guidelines

## Important

- Always save `.claude/` files (skills, settings, etc.) in the repository (`/workspace/.claude/`), never in `~/.claude/`.

### Key Skills (always invoke when conditions match)

- `/gh-token` — GitHub token retrieval via sidecar. Use before any Git remote operation or on auth errors using `gh`.
- `/branch-strategy` — Branching strategy. Use when creating branches, merging, or committing.
- `/coding-workflow` — Coding workflow guide. Use when deciding what to do next or how to proceed.

## Active Technologies

- Go
- SOPS
- age
- AWS CLI
