# psm Docker Reference Image

A reference Dockerfile bundling **psm**, **SOPS**, and **AWS CLI v2** for CI/CD pipelines.

> This is a build template — no prebuilt image is published. Users build the image themselves.

## Build

```bash
docker build -t psm-tools docker/
```

Override tool versions:

```bash
docker build \
  --build-arg PSM_VERSION=0.1.0 \
  --build-arg SOPS_VERSION=3.12.1 \
  --build-arg AWS_CLI_VERSION=2.34.9 \
  -t psm-tools docker/
```

## Usage

```bash
# Decrypt and sync
docker run --rm \
  -v ~/.aws:/home/psm/.aws:ro \
  -v $(pwd):/work \
  psm-tools \
  sh -c "sops -d secrets.enc.yaml | psm sync --store ssm --profile myprofile /dev/stdin"

# Check versions
docker run --rm psm-tools psm sync --help
docker run --rm psm-tools sops --version
docker run --rm psm-tools aws --version
```

## Included Tools

| Tool | License | Source |
|------|---------|--------|
| psm | (project license) | [GitHub](https://github.com/nakatatsu/psm) |
| SOPS | Mozilla Public License 2.0 | [GitHub](https://github.com/getsops/sops) |
| AWS CLI v2 | Apache License 2.0 | [GitHub](https://github.com/aws/aws-cli) |

## Platform Support

- linux/amd64
- linux/arm64

## Notes

- No credentials are baked into the image. Mount `~/.aws` or pass environment variables at runtime.
- SOPS requires an encryption backend (AWS KMS, age, PGP, etc.). Install additional tools as needed.
- All tool versions are pinned via build ARGs. Check upstream releases for updates.
