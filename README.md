# psm

SOPS-decrypted YAML files to AWS SSM Parameter Store / Secrets Manager sync tool.

## Install

```bash
go install github.com/nakatatsu/psm@latest
```

Or download prebuilt binaries from [GitHub Releases](https://github.com/nakatatsu/psm/releases).

## Access to AWS from DevContainer

```
aws sso login --sso-session psm-sandbox --use-device-code
```
