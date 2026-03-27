# CLI Schema: psm (after --store removal)

## Commands

### sync

```
psm sync [--profile <name>] [--dry-run] [--skip-approve] [--debug] [--delete <file>] <sync-file>
```

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| --profile | string | (none) | AWS profile name |
| --dry-run | bool | false | Show plan without executing |
| --skip-approve | bool | false | Skip approval prompt |
| --debug | bool | false | Enable debug logging |
| --delete | string | (none) | YAML file with delete regex patterns |

### export

```
psm export [--profile <name>] [--debug] <file>
```

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| --profile | string | (none) | AWS profile name |
| --debug | bool | false | Enable debug logging |

### Other

```
psm --version
```

## Removed Flags

| Flag | Removed In | Error Message |
|------|-----------|---------------|
| --store | #37 | `--store has been removed. SSM is now the default store. Remove --store from your command.` |
| --prune | (earlier) | `--prune has been removed. Use --delete <file> with regex patterns instead` |
