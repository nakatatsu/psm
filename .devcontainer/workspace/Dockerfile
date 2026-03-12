FROM debian:bookworm-slim

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

ARG TZ
ENV TZ="$TZ"

ARG CLAUDE_CODE_VERSION=latest
ARG AWS_CLI_VERSION=2.33.31
ARG GO_VERSION=1.26.1
ARG GOFUMPT_VERSION=0.9.2
ARG GOIMPORTS_VERSION=0.42.0
ARG GOLANGCI_LINT_VERSION=2.9.0
ARG GOVULNCHECK_VERSION=1.1.4
ARG GOSEC_VERSION=2.23.0
ARG OSV_SCANNER_VERSION=2.3.3
ARG STATICCHECK_VERSION=2026.1
ARG GIT_DELTA_VERSION=0.18.2
ARG ZSH_IN_DOCKER_VERSION=1.2.0

# Create node user (matches devcontainer.json remoteUser)
RUN groupadd --gid 1000 node && \
  useradd --uid 1000 --gid node --shell /bin/zsh --create-home node

# Install basic development tools and iptables/ipset
RUN apt-get update && apt-get install -y --no-install-recommends \
  less \
  git \
  procps \
  sudo \
  fzf \
  zsh \
  man-db \
  unzip \
  gnupg2 \
  gh \
  iptables \
  ipset \
  iproute2 \
  dnsutils \
  aggregate \
  jq \
  nano \
  vim \
  curl \
  wget \
  ca-certificates \
  locales \
  && apt-get clean && rm -rf /var/lib/apt/lists/*

# Locale
RUN sed -i 's/^# *\(en_US.UTF-8 UTF-8\)/\1/' /etc/locale.gen && locale-gen en_US.UTF-8
ENV LANG=en_US.UTF-8 \
    LC_ALL=en_US.UTF-8

# AWS CLI
RUN ARCH=$(uname -m) && \
  if [ "$ARCH" = "x86_64" ]; then \
    curl -fsSL "https://awscli.amazonaws.com/awscli-exe-linux-x86_64-${AWS_CLI_VERSION}.zip" -o "awscliv2.zip"; \
  elif [ "$ARCH" = "aarch64" ]; then \
    curl -fsSL "https://awscli.amazonaws.com/awscli-exe-linux-aarch64-${AWS_CLI_VERSION}.zip" -o "awscliv2.zip"; \
  fi && \
  unzip awscliv2.zip && ./aws/install && rm -rf aws awscliv2.zip

# Go
RUN ARCH="$(dpkg --print-architecture)" && \
  wget -q "https://go.dev/dl/go${GO_VERSION}.linux-${ARCH}.tar.gz" && \
  tar -C /usr/local -xzf "go${GO_VERSION}.linux-${ARCH}.tar.gz" && \
  rm "go${GO_VERSION}.linux-${ARCH}.tar.gz"

ENV PATH=$PATH:/usr/local/go/bin
ENV GOPATH=/home/node/go
ENV PATH=$PATH:/home/node/go/bin

RUN mkdir -p /home/node/go/bin && chown -R node:node /home/node/go

# staticcheck (binary from GitHub Releases — avoids proxy.golang.org dependency)
RUN ARCH=$(dpkg --print-architecture) && \
  if [ "$ARCH" = "amd64" ]; then SC_ARCH="amd64"; \
  elif [ "$ARCH" = "arm64" ]; then SC_ARCH="arm64"; fi && \
  wget -q "https://github.com/dominikh/go-tools/releases/download/${STATICCHECK_VERSION}/staticcheck_linux_${SC_ARCH}.tar.gz" -O /tmp/staticcheck.tar.gz && \
  tar xzf /tmp/staticcheck.tar.gz -C /tmp/ && \
  mv /tmp/staticcheck/staticcheck /usr/local/bin/staticcheck && \
  chmod +x /usr/local/bin/staticcheck && \
  rm -rf /tmp/staticcheck /tmp/staticcheck.tar.gz

# osv-scanner
RUN ARCH=$(dpkg --print-architecture) && \
  wget -q "https://github.com/google/osv-scanner/releases/download/v${OSV_SCANNER_VERSION}/osv-scanner_linux_${ARCH}" -O /usr/local/bin/osv-scanner && \
  chmod +x /usr/local/bin/osv-scanner

# git-delta
RUN ARCH=$(dpkg --print-architecture) && \
  wget -q "https://github.com/dandavison/delta/releases/download/${GIT_DELTA_VERSION}/git-delta_${GIT_DELTA_VERSION}_${ARCH}.deb" && \
  dpkg -i "git-delta_${GIT_DELTA_VERSION}_${ARCH}.deb" && \
  rm "git-delta_${GIT_DELTA_VERSION}_${ARCH}.deb"

# Persist bash history
RUN mkdir /commandhistory && \
  touch /commandhistory/.bash_history && \
  chown -R node /commandhistory

# Create workspace and config directories
RUN mkdir -p /workspace /home/node/.claude && \
  chown -R node:node /workspace /home/node/.claude

ENV DEVCONTAINER=true \
    SHELL=/bin/zsh \
    EDITOR=nano \
    VISUAL=nano \
    PATH="/home/node/.local/bin:${PATH}"

WORKDIR /workspace

# Switch to non-root user
USER node

# zsh + powerlevel10k
RUN sh -c "$(wget -O- https://github.com/deluan/zsh-in-docker/releases/download/v${ZSH_IN_DOCKER_VERSION}/zsh-in-docker.sh)" -- \
  -p git \
  -a "export PROMPT_COMMAND='history -a' && export HISTFILE=/commandhistory/.bash_history" \
  -a "alias yolo='claude --dangerously-skip-permissions'" \
  -x

# Claude Code (native installer, no Node.js required)
RUN curl -fsSL https://claude.ai/install.sh | bash -s ${CLAUDE_CODE_VERSION}

# Go tools
RUN go install mvdan.cc/gofumpt@v${GOFUMPT_VERSION} && \
  go install golang.org/x/tools/cmd/goimports@v${GOIMPORTS_VERSION} && \
  go install golang.org/x/vuln/cmd/govulncheck@v${GOVULNCHECK_VERSION} && \
  go install github.com/securego/gosec/v2/cmd/gosec@v${GOSEC_VERSION} && \
  go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v${GOLANGCI_LINT_VERSION}

# Firewall
COPY init-firewall.sh /usr/local/bin/
USER root
RUN chmod +x /usr/local/bin/init-firewall.sh && \
  echo "node ALL=(root) NOPASSWD: /usr/local/bin/init-firewall.sh" > /etc/sudoers.d/node-firewall && \
  chmod 0440 /etc/sudoers.d/node-firewall
USER node
