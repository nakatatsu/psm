# psm User Environment: psm + SOPS + age + AWS CLI v2
#
# Lightweight DevContainer for managing AWS Parameter Store secrets with psm.
# Does NOT include development tools (Go, linters, Claude Code, etc.).
#
# Licenses:
#   psm      — MIT (https://github.com/nakatatsu/psm)
#   SOPS     — Mozilla Public License 2.0 (https://github.com/getsops/sops)
#   age      — BSD 3-Clause (https://github.com/FiloSottile/age)
#   AWS CLI  — Apache License 2.0 (https://github.com/aws/aws-cli)

FROM debian:bookworm-slim

ARG PSM_VERSION=0.0.1
ARG SOPS_VERSION=3.12.1
ARG AGE_VERSION=1.2.1
ARG AWS_CLI_VERSION=2.34.9
ARG TARGETARCH

RUN apt-get update && apt-get install -y --no-install-recommends \
  curl \
  unzip \
  ca-certificates \
  git \
  jq \
  less \
  && apt-get clean && rm -rf /var/lib/apt/lists/*

# AWS CLI v2
RUN ARCH=$(uname -m) && \
  curl -fsSL "https://awscli.amazonaws.com/awscli-exe-linux-${ARCH}-${AWS_CLI_VERSION}.zip" -o awscliv2.zip && \
  unzip -q awscliv2.zip && ./aws/install && rm -rf aws awscliv2.zip

# SOPS
RUN curl -fsSL -o /usr/local/bin/sops \
  "https://github.com/getsops/sops/releases/download/v${SOPS_VERSION}/sops-v${SOPS_VERSION}.linux.${TARGETARCH}" && \
  chmod +x /usr/local/bin/sops

# age
RUN curl -fsSL -o /tmp/age.tar.gz \
  "https://github.com/FiloSottile/age/releases/download/v${AGE_VERSION}/age-v${AGE_VERSION}-linux-${TARGETARCH}.tar.gz" && \
  tar xzf /tmp/age.tar.gz -C /usr/local/bin --strip-components=1 age/age age/age-keygen && \
  chmod +x /usr/local/bin/age /usr/local/bin/age-keygen && \
  rm /tmp/age.tar.gz

# psm
RUN curl -fsSL -o /tmp/psm.tar.gz \
  "https://github.com/nakatatsu/psm/releases/download/v${PSM_VERSION}/psm_${PSM_VERSION}_linux_${TARGETARCH}.tar.gz" && \
  tar xzf /tmp/psm.tar.gz -C /usr/local/bin psm && \
  chmod +x /usr/local/bin/psm && \
  rm /tmp/psm.tar.gz

# Non-root user
RUN groupadd --gid 1000 psm && \
  useradd --uid 1000 --gid psm --create-home psm && \
  mkdir -p /home/psm/.aws && chown psm:psm /home/psm/.aws
USER psm
WORKDIR /workspace
