#!/usr/bin/env bash

#
# FOR TEST ONLY! DO NOT USE IN PRODUCTION!
#

# =============================================================================
# psm Integration Test Script
#
# Runs all behavioral scenarios against a real AWS sandbox and reports results.
#
# Prerequisites (complete these before running):
#   1. Open this directory in the DevContainer
#   2. Log in to AWS SSO:
#        aws configure sso --use-device-code
#        aws sso login --sso-session psm-sandbox --use-device-code
#
# Environment variables (required):
#   PSM_BIN           Path to psm binary
#   PSM_TEST_PROFILE  AWS profile name
#
# Usage:
#   PSM_BIN=./psm PSM_TEST_PROFILE=psm-sandbox bash test.sh
# =============================================================================
set -uo pipefail

# ---------------------------------------------------------------------------
# Configuration
# ---------------------------------------------------------------------------
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
TEST_YAML="${SCRIPT_DIR}/secrets.example.yaml"
if [[ -z "${PSM_BIN:-}" ]]; then
  echo "Error: PSM_BIN is required. Set it to the path of the psm binary."
  exit 1
fi
PSM="${PSM_BIN}"
if [[ -n "${PSM_TEST_PROFILE:-}" ]]; then
  PROFILE_FLAG="--profile ${PSM_TEST_PROFILE}"
else
  PROFILE_FLAG=""
fi

PASS_COUNT=0
FAIL_COUNT=0
SKIP_COUNT=0
SCENARIO_TOTAL=14

# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------
pass() {
  echo "PASS"
  PASS_COUNT=$((PASS_COUNT + 1))
}

fail() {
  echo "FAIL: $1"
  FAIL_COUNT=$((FAIL_COUNT + 1))
}

skip() {
  echo "SKIP: $1"
  SKIP_COUNT=$((SKIP_COUNT + 1))
}

# Get a single SSM parameter value (strongly consistent)
get_param() {
  aws ssm get-parameter \
    --name "$1" \
    --with-decryption \
    ${PROFILE_FLAG} \
    --query 'Parameter.Value' \
    --output text 2>/dev/null
}

# Put a single SSM parameter
put_param() {
  aws ssm put-parameter \
    --name "$1" \
    --value "$2" \
    --type SecureString \
    --overwrite \
    ${PROFILE_FLAG} \
    --output text
}

# Check if a parameter exists
param_exists() {
  aws ssm get-parameter \
    --name "$1" \
    --with-decryption \
    ${PROFILE_FLAG} \
    --output text >/dev/null 2>&1
}

# List all parameter names under /myapp/ (sorted)
get_all_param_names() {
  aws ssm get-parameters-by-path \
    --path "/myapp/" --recursive --with-decryption \
    ${PROFILE_FLAG} \
    --query 'Parameters[].Name' \
    --output text 2>/dev/null | tr '\t' '\n' | grep -v '^$' | sort
}

# Delete all parameters under /myapp/
cleanup_all() {
  local names
  names=$(get_all_param_names)
  if [[ -n "${names}" ]]; then
    # shellcheck disable=SC2086
    aws ssm delete-parameters \
      --names ${names} \
      ${PROFILE_FLAG} \
      --output text 2>/dev/null || true
  fi
}

# Assert SSM state matches expected key=value pairs exactly under /myapp/.
# Usage: assert_state "label" ["key=value" ...]
# - All expected keys must exist with matching values
# - No unexpected keys may exist under /myapp/
assert_state() {
  local label="$1"
  shift
  local errors=""

  # Check each expected key=value
  for pair in "$@"; do
    local key="${pair%%=*}"
    local expected_val="${pair#*=}"
    local actual_val
    actual_val=$(get_param "${key}" 2>/dev/null) || actual_val=""
    if [[ "${actual_val}" != "${expected_val}" ]]; then
      errors="${errors}\n  ${key}: expected '${expected_val}', got '${actual_val}'"
    fi
  done

  # Check no unexpected keys under /myapp/
  local actual_keys
  actual_keys=$(get_all_param_names)

  local actual_key
  for actual_key in ${actual_keys}; do
    local found=false
    for pair in "$@"; do
      if [[ "${pair%%=*}" == "${actual_key}" ]]; then
        found=true
        break
      fi
    done
    if [[ "${found}" == "false" ]]; then
      errors="${errors}\n  unexpected key: ${actual_key}"
    fi
  done

  if [[ -n "${errors}" ]]; then
    printf "  assert_state [%s]:%b\n" "${label}" "${errors}"
    return 1
  fi
  return 0
}

# Expected full state after a normal sync of secrets.example.yaml
EXPECTED_FULL_STATE=(
  "/myapp/database/host=localhost"
  "/myapp/database/port=5432"
  "/myapp/database/password=do-not-look-at-me"
  "/myapp/api/key=come-on-do-not-look-at-me"
)

# ---------------------------------------------------------------------------
# Prerequisites
# ---------------------------------------------------------------------------
echo "=== Prerequisites ==="

printf "aws credentials ... "
if aws sts get-caller-identity ${PROFILE_FLAG} --output text >/dev/null 2>&1; then
  echo "ok"
else
  echo "FAIL"
  echo "Error: AWS credentials not valid. Run: aws sso login --sso-session psm-sandbox --use-device-code"
  exit 1
fi
echo ""

# ---------------------------------------------------------------------------
# Safety gate: confirm the target AWS account before proceeding
# ---------------------------------------------------------------------------
CALLER_IDENTITY=$(aws sts get-caller-identity ${PROFILE_FLAG} --output json)
AWS_ACCOUNT_ID=$(echo "${CALLER_IDENTITY}" | grep -o '"Account": *"[^"]*"' | cut -d'"' -f4)
AWS_ARN=$(echo "${CALLER_IDENTITY}" | grep -o '"Arn": *"[^"]*"' | cut -d'"' -f4)

echo "╔══════════════════════════════════════════════════════════════"
echo "║  ⚠  FOR TEST ONLY! DO NOT USE IN PRODUCTION!  ⚠"
echo "║"
echo "║  This script WRITES and DELETES SSM parameters. "
echo "║  Running against the wrong account can cause data loss. "
echo "╠══════════════════════════════════════════════════════════════"
printf "║  Account : %-48s \n" "${AWS_ACCOUNT_ID}"
printf "║  ARN     : %-48s \n" "${AWS_ARN}"
echo "╚══════════════════════════════════════════════════════════════"
echo ""

if [[ "${CI:-}" == "true" ]]; then
  echo "CI mode: skipping interactive confirmation"
else
  if [[ ! -t 0 ]]; then
    echo "Error: this script must be run from an interactive terminal."
    exit 1
  fi

  printf "Continue? [y/N] "
  read -r answer
  if [[ "${answer}" != "y" && "${answer}" != "Y" ]]; then
    echo "Aborted."
    exit 1
  fi
fi
echo ""

# ---------------------------------------------------------------------------
# Setup: temporary directory, SOPS encryption, pattern files
# ---------------------------------------------------------------------------
echo "=== Setup ==="

TMPDIR_TEST=$(mktemp -d)
trap 'rm -rf "${TMPDIR_TEST}"' EXIT

# Generate ephemeral age key and encrypt secrets.example.yaml
AGE_KEY_FILE="${TMPDIR_TEST}/age-key.txt"
AGE_PUBLIC_KEY=$(age-keygen -o "${AGE_KEY_FILE}" 2>&1 | grep "Public key:" | awk '{print $3}')
export SOPS_AGE_KEY_FILE="${AGE_KEY_FILE}"

cat > "${TMPDIR_TEST}/.sops.yaml" <<SOPSEOF
creation_rules:
  - age: "${AGE_PUBLIC_KEY}"
SOPSEOF

(cd "${TMPDIR_TEST}" && sops -e "${TEST_YAML}" > secrets.enc.yaml)

# Delete pattern file (matches /myapp/legacy/)
cat > "${TMPDIR_TEST}/delete-patterns.yaml" <<DELEOF
- "^/myapp/legacy/"
DELEOF

# Conflict pattern file (matches a key that's also in sync YAML)
cat > "${TMPDIR_TEST}/conflict-patterns.yaml" <<CONFEOF
- "^/myapp/database/host$"
CONFEOF

# Cleanup: remove all test parameters under /myapp/
echo "cleaning up previous test data ..."
cleanup_all

# ---------------------------------------------------------------------------
# Scenario 1/14: Dry-run
# ---------------------------------------------------------------------------
echo "=== Scenario 1/${SCENARIO_TOTAL}: Dry-run ==="

stdout=$( (cd "${TMPDIR_TEST}" && sops -d secrets.enc.yaml) | \
  ${PSM} sync ${PROFILE_FLAG} --dry-run /dev/stdin )
echo "${stdout}"

if ! echo "${stdout}" | grep -q "(dry-run)"; then
  fail "stdout missing (dry-run) indicator"
elif ! assert_state "dry-run: SSM must be empty" ; then
  fail "dry-run modified AWS state"
else
  pass
fi
echo ""

# ---------------------------------------------------------------------------
# Scenario 2/14: Sync with --skip-approve
# ---------------------------------------------------------------------------
echo "=== Scenario 2/${SCENARIO_TOTAL}: Sync with --skip-approve ==="

exit_code=0
(cd "${TMPDIR_TEST}" && sops -d secrets.enc.yaml) | \
  ${PSM} sync ${PROFILE_FLAG} --skip-approve /dev/stdin || exit_code=$?

if [[ ${exit_code} -ne 0 ]]; then
  fail "exit code ${exit_code}, expected 0"
elif ! assert_state "sync: all 4 keys" "${EXPECTED_FULL_STATE[@]}"; then
  fail "SSM state does not match expected"
else
  pass
fi
echo ""

# ---------------------------------------------------------------------------
# Scenario 3/14: Delete with --delete
# ---------------------------------------------------------------------------
echo "=== Scenario 3/${SCENARIO_TOTAL}: Delete with --delete ==="

# Seed a parameter under /myapp/legacy/ (matches delete pattern, not in sync YAML)
put_param "/myapp/legacy/old-key" "to-be-deleted"

exit_code=0
${PSM} sync ${PROFILE_FLAG} --skip-approve \
  --delete "${TMPDIR_TEST}/delete-patterns.yaml" \
  "${TEST_YAML}" || exit_code=$?

if [[ ${exit_code} -ne 0 ]]; then
  fail "exit code ${exit_code}, expected 0"
elif ! assert_state "delete: legacy gone, sync keys intact" "${EXPECTED_FULL_STATE[@]}"; then
  fail "SSM state does not match expected (legacy key should be gone, sync keys intact)"
else
  pass
fi
echo ""

# ---------------------------------------------------------------------------
# Scenario 4/14: Conflict detection
# ---------------------------------------------------------------------------
echo "=== Scenario 4/${SCENARIO_TOTAL}: Conflict detection ==="

# Capture state before conflict attempt
exit_code=0
${PSM} sync ${PROFILE_FLAG} --skip-approve \
  --delete "${TMPDIR_TEST}/conflict-patterns.yaml" \
  "${TEST_YAML}" 2>&1 || exit_code=$?

if [[ ${exit_code} -ne 1 ]]; then
  fail "exit code ${exit_code}, expected 1"
elif ! assert_state "conflict: no partial apply" "${EXPECTED_FULL_STATE[@]}"; then
  fail "conflict caused partial state change"
else
  pass
fi
echo ""

# ---------------------------------------------------------------------------
# Scenario 5/14: Debug logging
# ---------------------------------------------------------------------------
echo "=== Scenario 5/${SCENARIO_TOTAL}: Debug logging ==="

# 5a: --debug produces level=DEBUG
output_debug=$(${PSM} sync ${PROFILE_FLAG} --debug --dry-run \
  "${TEST_YAML}" 2>&1)
echo "${output_debug}"

# 5b: without --debug, level=DEBUG must NOT appear
output_normal=$(${PSM} sync ${PROFILE_FLAG} --dry-run \
  "${TEST_YAML}" 2>&1)

if ! echo "${output_debug}" | grep -q "level=DEBUG"; then
  fail "--debug did not produce level=DEBUG"
elif echo "${output_normal}" | grep -q "level=DEBUG"; then
  fail "level=DEBUG appeared without --debug flag"
else
  pass
fi
echo ""

# ---------------------------------------------------------------------------
# Scenario 6/14: Piped input + /dev/tty approve
# ---------------------------------------------------------------------------
echo "=== Scenario 6/${SCENARIO_TOTAL}: Piped input + /dev/tty approve ==="

if ! command -v script >/dev/null 2>&1; then
  skip "'script' command not found"
else
  # Modify a value so sync has work to do
  put_param "/myapp/database/host" "old-value"

  # Pipe YAML via stdin; approval prompt reads from /dev/tty (via script pseudo-TTY)
  exit_code=0
  printf 'y\n' | script -qec \
    "cat ${TEST_YAML} | ${PSM} sync ${PROFILE_FLAG} /dev/stdin" /dev/null \
    || exit_code=$?

  if ! assert_state "piped+tty yes: sync applied" "${EXPECTED_FULL_STATE[@]}"; then
    fail "piped input + tty approve did not apply changes"
  else
    pass
  fi
fi
echo ""

# ---------------------------------------------------------------------------
# Scenario 7/14: No changes
# ---------------------------------------------------------------------------
echo "=== Scenario 7/${SCENARIO_TOTAL}: No changes ==="

# Ensure AWS state matches the YAML exactly
put_param "/myapp/database/host" "localhost"
put_param "/myapp/database/port" "5432"
put_param "/myapp/database/password" "do-not-look-at-me"
put_param "/myapp/api/key" "come-on-do-not-look-at-me"

stdout=$(${PSM} sync ${PROFILE_FLAG} --skip-approve \
  "${TEST_YAML}")

if ! echo "${stdout}" | grep -q "0 created, 0 updated, 0 deleted"; then
  fail "expected no-changes summary, got: ${stdout}"
elif ! assert_state "no-op: state unchanged" "${EXPECTED_FULL_STATE[@]}"; then
  fail "no-op run changed SSM state"
else
  pass
fi
echo ""

# ---------------------------------------------------------------------------
# Scenario 8/14: Missing file
# ---------------------------------------------------------------------------
echo "=== Scenario 8/${SCENARIO_TOTAL}: Missing file ==="

exit_code=0
${PSM} sync ${PROFILE_FLAG} --skip-approve \
  "/nonexistent/path/secrets.yaml" 2>/dev/null || exit_code=$?

if [[ ${exit_code} -eq 0 ]]; then
  fail "expected non-zero exit code for missing file"
elif ! assert_state "missing file: state unchanged" "${EXPECTED_FULL_STATE[@]}"; then
  fail "missing file error changed SSM state"
else
  pass
fi
echo ""

# ---------------------------------------------------------------------------
# Scenario 9/14: Invalid YAML
# ---------------------------------------------------------------------------
echo "=== Scenario 9/${SCENARIO_TOTAL}: Invalid YAML ==="

cat > "${TMPDIR_TEST}/invalid.yaml" <<INVALEOF
invalid: [broken
  not: valid: yaml: {{{}}}
INVALEOF

exit_code=0
${PSM} sync ${PROFILE_FLAG} --skip-approve \
  "${TMPDIR_TEST}/invalid.yaml" 2>/dev/null || exit_code=$?

if [[ ${exit_code} -eq 0 ]]; then
  fail "expected non-zero exit code for invalid YAML"
elif ! assert_state "invalid YAML: state unchanged" "${EXPECTED_FULL_STATE[@]}"; then
  fail "invalid YAML error changed SSM state"
else
  pass
fi
echo ""

# ---------------------------------------------------------------------------
# Scenario 10/14: Empty input
# ---------------------------------------------------------------------------
echo "=== Scenario 10/${SCENARIO_TOTAL}: Empty input ==="

exit_code=0
${PSM} sync ${PROFILE_FLAG} --skip-approve \
  /dev/null 2>/dev/null || exit_code=$?

if [[ ${exit_code} -eq 0 ]]; then
  fail "expected non-zero exit code for empty input"
elif ! assert_state "empty input: state unchanged" "${EXPECTED_FULL_STATE[@]}"; then
  fail "empty input error changed SSM state"
else
  pass
fi
echo ""

# ---------------------------------------------------------------------------
# Scenario 11/14: TTY approve yes
# ---------------------------------------------------------------------------
echo "=== Scenario 11/${SCENARIO_TOTAL}: TTY approve yes ==="

if ! command -v script >/dev/null 2>&1; then
  skip "'script' command not found"
else
  # Change a value so sync has work to do
  put_param "/myapp/database/host" "old-value"

  # script creates a pseudo-TTY; pipe "y" to approve
  exit_code=0
  printf 'y\n' | script -qec "${PSM} sync ${PROFILE_FLAG} ${TEST_YAML}" /dev/null \
    || exit_code=$?

  if ! assert_state "TTY yes: sync applied" "${EXPECTED_FULL_STATE[@]}"; then
    fail "approve yes did not apply changes"
  else
    pass
  fi
fi
echo ""

# ---------------------------------------------------------------------------
# Scenario 12/14: TTY approve no
# ---------------------------------------------------------------------------
echo "=== Scenario 12/${SCENARIO_TOTAL}: TTY approve no ==="

if ! command -v script >/dev/null 2>&1; then
  skip "'script' command not found"
else
  # Change a value so sync has work to do
  put_param "/myapp/database/host" "should-stay"

  # script creates a pseudo-TTY; pipe "n" to decline
  exit_code=0
  printf 'n\n' | script -qec "${PSM} sync ${PROFILE_FLAG} ${TEST_YAML}" /dev/null \
    || exit_code=$?

  val=$(get_param "/myapp/database/host" || echo "")
  if [[ "${val}" != "should-stay" ]]; then
    fail "approve no should not change state, but /myapp/database/host='${val}'"
  else
    # Also test empty input (just Enter) -> should decline (default N)
    printf '\n' | script -qec "${PSM} sync ${PROFILE_FLAG} ${TEST_YAML}" /dev/null \
      || true

    val2=$(get_param "/myapp/database/host" || echo "")
    if [[ "${val2}" != "should-stay" ]]; then
      fail "empty Enter should not change state, but /myapp/database/host='${val2}'"
    else
      pass
    fi
  fi

  # Restore for clean exit
  put_param "/myapp/database/host" "localhost"
fi
echo ""

# ---------------------------------------------------------------------------
# Scenario 13/14: Piped input + /dev/tty decline
# ---------------------------------------------------------------------------
echo "=== Scenario 13/${SCENARIO_TOTAL}: Piped input + /dev/tty decline ==="

if ! command -v script >/dev/null 2>&1; then
  skip "'script' command not found"
else
  # Change a value so sync has work to do
  put_param "/myapp/database/host" "should-stay"

  # Pipe YAML via stdin; decline approval via /dev/tty (pseudo-TTY from script)
  exit_code=0
  printf 'N\n' | script -qec \
    "cat ${TEST_YAML} | ${PSM} sync ${PROFILE_FLAG} /dev/stdin" /dev/null \
    || exit_code=$?

  val=$(get_param "/myapp/database/host" || echo "")
  if [[ "${val}" != "should-stay" ]]; then
    fail "piped input + tty decline should not change state, got '${val}'"
  else
    pass
  fi

  # Restore for next scenario
  put_param "/myapp/database/host" "localhost"
fi
echo ""

# ---------------------------------------------------------------------------
# Scenario 14/14: Piped input + no TTY → error
# ---------------------------------------------------------------------------
echo "=== Scenario 14/${SCENARIO_TOTAL}: Piped input + no TTY → error ==="

if ! command -v setsid >/dev/null 2>&1; then
  skip "'setsid' command not found"
else
  # setsid detaches from controlling terminal, so /dev/tty is unavailable
  exit_code=0
  cat "${TEST_YAML}" | \
    setsid --wait ${PSM} sync ${PROFILE_FLAG} /dev/stdin 2>/dev/null || exit_code=$?

  if [[ ${exit_code} -ne 1 ]]; then
    fail "expected exit code 1 (no tty), got ${exit_code}"
  elif ! assert_state "no-tty: state unchanged" "${EXPECTED_FULL_STATE[@]}"; then
    fail "no-tty error changed SSM state"
  else
    pass
  fi
fi
echo ""

# ---------------------------------------------------------------------------
# Results
# ---------------------------------------------------------------------------
echo "=== Results: ${PASS_COUNT} passed, ${FAIL_COUNT} failed, ${SKIP_COUNT} skipped ==="

if [[ ${FAIL_COUNT} -gt 0 ]]; then
  exit 1
fi
exit 0
