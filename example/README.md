# psm Parameter Store Management

Manage AWS SSM Parameter Store secrets with [psm](https://github.com/nakatatsu/psm), encrypted by SOPS + age.

## Prerequisites

- Docker (for DevContainer)
- VS Code with [Dev Containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers), or another DevContainer-compatible editor
- An AWS account with SSO configured

## Quick Start

### 1. Open DevContainer

Open this directory in VS Code and select **"Reopen in Container"**. The container includes:

| Tool | Purpose | License |
|------|---------|---------|
| psm | Sync secrets to Parameter Store | MIT |
| SOPS | Encrypt/decrypt secrets files | MPL 2.0 |
| age | Encryption backend for SOPS | BSD 3-Clause |
| AWS CLI v2 | AWS authentication and debugging | Apache 2.0 |

### 2. Generate an age key

```bash
age-keygen -o age-key.txt
```

Save the public key (`age1...`) from the output. Keep `age-key.txt` safe — this is your private key.

> Add `age-key.txt` to `.gitignore` to avoid committing your private key.

### 3. Configure SOPS

```bash
cp .sops.yaml.example .sops.yaml
```

Edit `.sops.yaml` and replace the placeholder with your age public key:

```yaml
creation_rules:
  - age: "age1your-actual-public-key"
```

### 4. Encrypt secrets

```bash
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
sops -d secrets.enc.yaml | psm sync --store ssm --profile <your-profile> /dev/stdin
```

> `--use-device-code` is required inside DevContainers because the default OAuth callback to localhost does not work.

### 7. Verify

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

```jsonc
// .devcontainer/devcontainer.json
{
  "build": {
    "dockerfile": "Dockerfile",
    "args": {
      "PSM_VERSION": "0.1.0",
      "SOPS_VERSION": "3.12.1",
      "AGE_VERSION": "1.2.1",
      "AWS_CLI_VERSION": "2.34.9"
    }
  }
}
```

## Platform Support

- linux/amd64
- linux/arm64

## Notes

- **Production environments should use AWS KMS instead of age keys.** KMS provides key rotation, IAM-based access control, and CloudTrail audit logging. With KMS, no key file distribution is needed.
- age-key.txt must never be committed to the repository.
- AWS configuration (`~/.aws`) is stored in a Docker named volume (`psm-aws-config`), independent of the host. Run `aws configure sso` once inside the container; the settings persist across rebuilds.
