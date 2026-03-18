# workspace Development Guidelines

## Important

- Always save `.claude/` files (skills, settings, etc.) in the repository (`/workspace/.claude/`), never in `~/.claude/`.

## Active Technologies
- Go 1.26 + AWS SDK for Go v2, gopkg.in/yaml.v3, `regexp` (stdlib), `log/slog` (stdlib) (018-approve-delete-flow)
- AWS SSM Parameter Store, AWS Secrets Manager (via Store interface) (018-approve-delete-flow)

- Go + 標準ライブラリ
- SOPS, age, AWS CLI

## Recent Changes
- 018-approve-delete-flow: Added Go 1.26 + AWS SDK for Go v2, gopkg.in/yaml.v3, `regexp` (stdlib), `log/slog` (stdlib)
