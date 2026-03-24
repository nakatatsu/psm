# workspace Development Guidelines

## Important

- Always save `.claude/` files (skills, settings, etc.) in the repository (`/workspace/.claude/`), never in `~/.claude/`.

## Active Technologies
- Markdown (Claude Code SKILL.md format) + git 2.39+, bash (026-worktree-setup-skill)

- Go 1.26 + AWS SDK for Go v2, gopkg.in/yaml.v3, regexp, log/slog (027-omit-secrets-manager)
- AWS SSM Parameter Store (027-omit-secrets-manager)

- Go
- SOPS, age, AWS CLI

## Codex CLI

- The `/codex` skill runs with `--dangerously-bypass-approvals-and-sandbox` due to bwrap sandbox limitations in DevContainer. The skill requires explicit user permission by default, but the project owner has pre-approved its use in this repository.
