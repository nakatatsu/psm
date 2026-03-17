# psm

> **Status: Under Development**

YAML files to AWS SSM Parameter Store / Secrets Manager sync tool.

## Install

Download prebuilt binaries from [GitHub Releases](https://github.com/nakatatsu/psm/releases).

## Usage

### Sync: push YAML to AWS

```bash
psm sync --store <ssm|sm> [--profile <name>] [--dry-run] [--prune] <file>
```

| Flag | Required | Description |
|------|----------|-------------|
| `--store <ssm\|sm>` | Yes | Target store: `ssm` (Parameter Store) or `sm` (Secrets Manager) |
| `--profile <name>` | No | AWS profile name (default: SDK default credentials) |
| `--dry-run` | No | Show planned changes without executing |
| `--prune` | No | Delete keys in AWS that are not in the YAML file |

Example:

```bash
# Preview changes
psm sync --store ssm --dry-run secrets.yaml

# Apply changes
psm sync --store ssm secrets.yaml

# Apply with a specific AWS profile
psm sync --store ssm --profile myprofile secrets.yaml

# Sync and remove keys not in YAML
psm sync --store ssm --prune secrets.yaml
```

### Export: pull AWS parameters to YAML

```bash
psm export --store <ssm|sm> [--profile <name>] <file>
```

Example:

```bash
psm export --store ssm output.yaml
```

### With SOPS (encrypted secrets)

Decrypt with SOPS and pipe directly to psm:

```bash
sops -d secrets.enc.yaml | psm sync --store ssm --dry-run /dev/stdin
sops -d secrets.enc.yaml | psm sync --store ssm /dev/stdin
```

See [example/README.md](example/README.md) for a full walkthrough including key generation and encryption setup.

### YAML format

Keys map directly to AWS parameter/secret names. Values must be scalars (string, int, bool, float).

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
