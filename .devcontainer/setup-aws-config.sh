#!/bin/bash
set -euo pipefail

# Require environment variables (injected via devcontainer.json containerEnv)
for var in AWS_SANDBOX_SSO_START_URL AWS_SANDBOX_SSO_ACCOUNT_ID AWS_SSO_ROLE_NAME AWS_SSO_REGION; do
    if [ -z "${!var:-}" ]; then
        echo "Missing $var environment variable, skipping AWS config setup"
        exit 0
    fi
done

mkdir -p /home/node/.aws

cat > /home/node/.aws/config <<AWSEOF
[sso-session psm-sandbox]
sso_start_url = ${AWS_SANDBOX_SSO_START_URL}
sso_region = ${AWS_SSO_REGION}
sso_registration_scopes = sso:account:access

[profile psm]
sso_session = psm-sandbox
sso_account_id = ${AWS_SANDBOX_SSO_ACCOUNT_ID}
sso_role_name = ${AWS_SSO_ROLE_NAME}
region = ${AWS_SSO_REGION}
output = json
AWSEOF

echo "AWS config generated at ~/.aws/config"
