# psm

> **Status: Under Development**

YAML files to AWS SSM Parameter Store sync tool.

## Install

Download prebuilt binaries from [GitHub Releases](https://github.com/nakatatsu/psm/releases).

## Usage

### Sync: push YAML to AWS

```bash
psm sync --store ssm [--profile <name>] [--dry-run] [--skip-approve] [--debug] [--delete <file>] <sync-file>
```

| Flag | Required | Description |
|------|----------|-------------|
| `--store ssm` | Yes | Target store: `ssm` (Parameter Store) |
| `--profile <name>` | No | AWS profile name (default: SDK default credentials) |
| `--dry-run` | No | Show planned changes without executing or prompting |
| `--skip-approve` | No | Skip approval prompt and execute immediately (for CI/automation) |
| `--debug` | No | Enable debug-level logging |
| `--delete <file>` | No | YAML file with regex patterns for key deletion (see [Deleting keys](#deleting-keys)) |

By default, `psm sync` displays the action plan and asks for confirmation before executing. Type `y` to proceed, or press Enter to cancel.

```bash
# Preview changes (no prompt, no execution)
psm sync --store ssm --dry-run secrets.yaml

# Apply changes (displays plan, asks for approval)
psm sync --store ssm secrets.yaml

# Apply with a specific AWS profile
psm sync --store ssm --profile myprofile secrets.yaml

# Skip approval (for CI/CD pipelines)
psm sync --store ssm --skip-approve secrets.yaml
```

> When stdin is not a terminal (e.g., piped input) and `--skip-approve` is not set, the tool automatically declines and exits without changes.

### Deleting keys

To delete obsolete AWS keys, create a YAML file listing regex patterns and pass it via `--delete`:

```yaml
# needless.yml — patterns for keys to delete
- "^/myapp/legacy/"
- "^/myapp/deprecated-.*"
```

```bash
psm sync --store ssm --delete needless.yml secrets.yaml
```

Only keys matching the patterns **and not present in the sync YAML** are deleted. The sync file is always required with `--delete`.

**Safety features:**
- **Conflict detection**: If a key matches a delete pattern but also exists in the sync YAML, the entire operation aborts before any changes are made.
- **Unmanaged key warnings**: Keys in AWS that are not in the sync YAML and don't match any delete pattern are shown as warnings.
- **Approval required**: Deletions are included in the action plan and require the same approval as creates/updates.

> **Migrating from `--prune`**: The `--prune` flag has been removed. Use `--delete <file>` with explicit regex patterns instead. To replicate the old `--prune` behavior (delete everything not in YAML), use a delete file containing `- ".*"`.

### Export: pull AWS parameters to YAML

```bash
psm export --store ssm [--profile <name>] [--debug] <file>
```

```bash
psm export --store ssm output.yaml
```

### With SOPS (encrypted secrets)

Decrypt with SOPS and pipe directly to psm:

```bash
sops -d secrets.enc.yaml | psm sync --store ssm --dry-run /dev/stdin
sops -d secrets.enc.yaml | psm sync --store ssm --skip-approve /dev/stdin
```

> Note: When piping, stdin is not a terminal, so `--skip-approve` is required to execute changes.

See [example/README.md](example/README.md) for a full walkthrough including key generation and encryption setup.

### YAML format

Keys map directly to AWS parameter names. Values must be scalars (string, int, bool, float).

```yaml
/myapp/database/host: localhost
/myapp/database/port: "5432"
/myapp/database/password: my-secret-password
/myapp/api/key: my-api-key
```

> The `sops` metadata key is automatically excluded during sync.

## Development

### Access to AWS from DevContainer

```
aws sso login --sso-session psm-sandbox --use-device-code
```
