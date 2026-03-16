# Smoke test — run inside DevContainer to verify all tools are installed.
# Usage: bash test.sh

psm sync --help
sops --version
age --version
age-keygen --help
aws --version
echo "All tools OK"
