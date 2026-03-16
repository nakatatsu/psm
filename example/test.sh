# psm Example — Manual Verification Commands
# Run these inside the DevContainer, one by one.

# --- 0. Set age key path (required for all sops operations) ---
export SOPS_AGE_KEY_FILE=$(pwd)/age-key.txt

# --- 1. Tool checks ---
psm sync --help
sops --version
age --version
age-keygen --help
aws --version

# --- 2. age key generation ---
age-keygen -o age-key.txt
# Note the public key (age1...) from the output

# --- 3. SOPS configuration ---
cp .sops.yaml.example .sops.yaml
# Edit .sops.yaml: replace placeholder with your age public key

# --- 4. Encrypt ---
sops -e secrets.yaml > secrets.enc.yaml

# --- 5. Decrypt (verify round-trip) ---
sops -d secrets.enc.yaml

# --- 6. AWS SSO setup & login ---
# ~/.aws is a named volume (persists across rebuilds).
# First time only (--use-device-code required in DevContainer):
aws configure sso --use-device-code
# Subsequent logins (replace PROFILE with your profile name):
aws sso login --use-device-code --profile PROFILE
aws sts get-caller-identity --profile PROFILE

# --- 7. psm sync (replace PROFILE with your profile name) ---
sops -d secrets.enc.yaml | psm sync --store ssm --profile PROFILE /dev/stdin

# --- 8. Verify in Parameter Store ---
aws ssm get-parameter --name /myapp/database/host --profile PROFILE --query 'Parameter.Value' --output text
aws ssm get-parameter --name /myapp/api/key --profile PROFILE --query 'Parameter.Value' --output text

# --- 9. Cleanup (optional) ---
aws ssm delete-parameter --name /myapp/database/host --profile PROFILE
aws ssm delete-parameter --name /myapp/database/port --profile PROFILE
aws ssm delete-parameter --name /myapp/database/password --profile PROFILE
aws ssm delete-parameter --name /myapp/api/key --profile PROFILE
rm -f age-key.txt secrets.enc.yaml .sops.yaml
