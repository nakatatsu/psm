# Reference Dockerfile: psm + SOPS + AWS CLI v2
#
# This is a convenience template for CI/CD pipelines.
# Users build this image themselves — no prebuilt image is published.
#
# Licenses:
#   psm      — MIT (or project license)
#   SOPS     — Mozilla Public License 2.0 (https://github.com/getsops/sops)
#   AWS CLI  — Apache License 2.0 (https://github.com/aws/aws-cli)
#
# Usage:
#   docker build -t psm-tools docker/
#   docker run --rm -v ~/.aws:/home/psm/.aws:ro psm-tools psm sync --store ssm secrets.yaml
#
# Override versions at build time:
#   docker build --build-arg PSM_VERSION=0.1.0 --build-arg SOPS_VERSION=3.12.1 -t psm-tools docker/

FROM debian:bookworm-slim

ARG PSM_VERSION=0.0.1
ARG SOPS_VERSION=3.12.1
ARG AWS_CLI_VERSION=2.34.9
ARG TARGETARCH

RUN apt-get update && apt-get install -y --no-install-recommends \
  curl \
  unzip \
  ca-certificates \
  && apt-get clean && rm -rf /var/lib/apt/lists/*

# AWS CLI v2
RUN ARCH=$(uname -m) && \
  curl -fsSL "https://awscli.amazonaws.com/awscli-exe-linux-${ARCH}-${AWS_CLI_VERSION}.zip" -o awscliv2.zip && \
  unzip -q awscliv2.zip && ./aws/install && rm -rf aws awscliv2.zip

# SOPS
RUN curl -fsSL -o /usr/local/bin/sops \
  "https://github.com/getsops/sops/releases/download/v${SOPS_VERSION}/sops-v${SOPS_VERSION}.linux.${TARGETARCH}" && \
  chmod +x /usr/local/bin/sops

# psm
RUN curl -fsSL -o /tmp/psm.tar.gz \
  "https://github.com/nakatatsu/psm/releases/download/v${PSM_VERSION}/psm_${PSM_VERSION}_linux_${TARGETARCH}.tar.gz" && \
  tar xzf /tmp/psm.tar.gz -C /usr/local/bin psm && \
  chmod +x /usr/local/bin/psm && \
  rm /tmp/psm.tar.gz

# Non-root user
RUN groupadd --gid 1000 psm && \
  useradd --uid 1000 --gid psm --create-home psm
USER psm
WORKDIR /work
