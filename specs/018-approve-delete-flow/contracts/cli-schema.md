# CLI Schema: psm sync (updated)

## Synopsis

```
psm sync --store <ssm|sm> [--profile <name>] [--delete <file>] [--dry-run] [--skip-approve] [--debug] <sync-file>
psm export --store <ssm|sm> [--profile <name>] [--debug] <output-file>
```

## Flags

### sync subcommand

| Flag | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `--store` | string | Yes | — | Store type: `ssm` or `sm` |
| `--profile` | string | No | — | AWS profile name |
| `--delete` | string | No | — | Path to YAML file with delete regex patterns |
| `--dry-run` | bool | No | false | Show plan without executing or prompting |
| `--skip-approve` | bool | No | false | Skip approval prompt, execute immediately |
| `--debug` | bool | No | false | Enable debug-level logging |

### export subcommand

| Flag | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `--store` | string | Yes | — | Store type: `ssm` or `sm` |
| `--profile` | string | No | — | AWS profile name |
| `--debug` | bool | No | false | Enable debug-level logging |

### Removed flags

| Flag | Migration |
|------|-----------|
| `--prune` | Use `--delete <file>` with regex patterns instead |

## stdout Format

```
create: /myapp/prod/NEW_KEY
update: /myapp/prod/CHANGED_KEY
delete: /myapp/legacy/OLD_KEY
3 created, 1 updated, 1 deleted, 5 unchanged, 0 failed
```

## stderr Format

```
# Approval prompt (interactive only)
Proceed? [y/N]

# Warnings (unmanaged keys)
time=... level=WARN msg="unmanaged key detected" key=/other/team/KEY

# Errors
time=... level=ERROR msg="failed to put parameter" key=/myapp/prod/KEY error="AccessDeniedException: ..."

# Debug (with --debug)
time=... level=DEBUG msg="API call" operation=PutParameter key=/myapp/prod/KEY

# Migration error (--prune used)
error: --prune has been removed. Use --delete <file> with regex patterns instead.
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success, or user declined approval |
| 1 | Failure (validation error, conflict detected, API error, invalid input) |

## Delete Pattern File Format

```yaml
# Each entry is a Go RE2 regular expression
- "^/myapp/legacy/"
- "^/myapp/deprecated-.*"
- "^/myapp/temp/test-"
```
