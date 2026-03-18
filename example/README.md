# psm Parameter Store Management

Manage AWS SSM Parameter Store secrets with [psm](https://github.com/nakatatsu/psm), encrypted by SOPS + age.

## Prerequisites

- Docker (for DevContainer)
- VS Code with [Dev Containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers), or another DevContainer-compatible editor
- An AWS account with SSO configured

## Quick Start

### 1. Open DevContainer

Open this directory in VS Code and select **"Reopen in Container"**. The container includes:

| Tool       | Purpose                          | License      |
| ---------- | -------------------------------- | ------------ |
| psm        | Sync secrets to Parameter Store  | MIT          |
| SOPS       | Encrypt/decrypt secrets files    | MPL 2.0      |
| age        | Encryption backend for SOPS      | BSD 3-Clause |
| AWS CLI v2 | AWS authentication and debugging | Apache 2.0   |

### 2. Generate an age key

```bash
age-keygen -o age-key.txt
```

Save the public key (`age1...`) from the output. Keep `age-key.txt` safe — this is your private key.

> Add `age-key.txt` to `.gitignore` to avoid committing your private key.

### 3. Configure SOPS

```bash
cp .sops.example.yaml .sops.yaml
```

Edit `.sops.yaml` and replace the placeholder with your age public key:

```yaml
creation_rules:
  - age: "age1your-actual-public-key"
```

### 4. Encrypt secrets

```bash
cp secrets.example.yaml secrets.yaml

export SOPS_AGE_KEY_FILE=$(pwd)/age-key.txt
sops -e secrets.yaml > secrets.enc.yaml
```

Decrypt:

```bash
sops -d secrets.enc.yaml
```

From now on, commit `secrets.enc.yaml` (encrypted) and never commit `secrets.yaml` (plaintext).

### 5. Configure AWS SSO (first time only)

AWS configuration is stored in a Docker named volume (`psm-aws-config`) and persists across container rebuilds.

```bash
aws configure sso --use-device-code
```

Follow the prompts to set up your SSO session. You will need:

- SSO start URL (e.g., `https://your-org.awsapps.com/start`)
- SSO region
- Account ID and role name

### 6. Log in and sync

```bash
aws sso login --use-device-code --profile <your-profile>

# if need
export SOPS_AGE_KEY_FILE=$(pwd)/age-key.txt
```

> `--use-device-code` is required inside DevContainers because the default OAuth callback to localhost does not work.

#### Preview changes (dry-run)

No prompt, no execution — just shows what would happen:

```bash
sops -d secrets.enc.yaml | psm sync --store ssm --profile <your-profile> --dry-run /dev/stdin

# (dry-run) create: /myapp/prod/NEW_KEY
# (dry-run) update: /myapp/prod/CHANGED_KEY
# 1 created, 1 updated, 0 deleted, 3 unchanged, 0 failed (dry-run)
```

#### Apply changes (interactive)

Displays the action plan and prompts for confirmation before executing:

```bash
sops -d secrets.enc.yaml | psm sync --store ssm --profile <your-profile> --skip-approve /dev/stdin

# create: /myapp/prod/NEW_KEY
# update: /myapp/prod/CHANGED_KEY
# 1 created, 1 updated, 0 deleted, 3 unchanged, 0 failed
```

> When piping via SOPS, stdin is not a terminal, so `--skip-approve` is required to execute changes.
> Without `--skip-approve`, psm automatically declines and exits without making changes.

#### Apply changes (CI / automation)

For non-interactive environments such as CI/CD pipelines, use `--skip-approve`:

```bash
psm sync --store ssm --skip-approve secrets.yaml
```

### 7. Delete obsolete keys

Create a YAML file listing regex patterns for keys to delete:

```yaml
# needless.yml
- "^/myapp/legacy/"
- "^/myapp/deprecated-.*"
```

```bash
sops -d secrets.enc.yaml | psm sync --store ssm --profile <your-profile> --skip-approve --delete needless.yml /dev/stdin
```

Only keys matching the patterns **and not present in the sync YAML** are deleted.

**Safety features:**

- **Conflict detection** — If a key matches a delete pattern but also exists in the sync YAML, the entire operation aborts before any changes are made.
- **Unmanaged key warnings** — Keys in AWS that are not in the sync YAML and don't match any delete pattern are shown as warnings.
- **Approval required** — Deletions are included in the action plan and require the same approval as creates/updates.

> **Migrating from `--prune`**: The `--prune` flag has been removed. Use `--delete <file>` with explicit regex patterns instead. To replicate the old `--prune` behavior (delete everything not in YAML), use a delete file containing `- ".*"`.

### 8. Debug logging

Add `--debug` to any command to see detailed slog output on stderr:

```bash
psm sync --store ssm --debug secrets.yaml

# stderr:
# time=... level=DEBUG msg="API call" operation=GetParametersByPath path=/
# time=... level=DEBUG msg="regex match" pattern="^/myapp/legacy/" key=/myapp/legacy/OLD matched=true
```

### 9. Verify

```bash
aws ssm get-parameter --name /myapp/database/host --profile <your-profile> --query 'Parameter.Value' --output text
```

## File Structure

```
.devcontainer/       DevContainer definition (Dockerfile + config)
.sops.yaml.example   SOPS config template (copy to .sops.yaml)
secrets.yaml         Sample secrets (replace with real values)
README.md            This file
test.sh              Manual verification commands
```

## Customizing Tool Versions

Rebuild the DevContainer with custom versions:

## CLI Reference

```bash
psm sync --store <ssm|sm> [--profile <name>] [--dry-run] [--skip-approve] [--debug] [--delete <file>] <sync-file>
psm export --store <ssm|sm> [--profile <name>] [--debug] <file>
```

| Flag | Description |
|------|-------------|
| `--store <ssm\|sm>` | **(Required)** Target store: `ssm` (Parameter Store) or `sm` (Secrets Manager) |
| `--profile <name>` | AWS profile name (default: SDK default credentials) |
| `--dry-run` | Show planned changes without executing or prompting |
| `--skip-approve` | Skip approval prompt and execute immediately |
| `--debug` | Enable debug-level logging |
| `--delete <file>` | YAML file with regex patterns for key deletion |

## Notes

- Production environments should use AWS KMS instead of age keys.
- age-key.txt must never be committed to the repository.
- AWS configuration (`~/.aws`) is stored in a Docker named volume (`psm-aws-config`), independent of the host. Run `aws configure sso` once inside the container; the settings persist across rebuilds.
