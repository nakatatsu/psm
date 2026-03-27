# Quickstart: Read Approval Prompt from /dev/tty

## What Changed

The approval prompt (`Proceed? [y/N]`) now reads from `/dev/tty` when stdin is piped, enabling the primary SOPS workflow:

```bash
sops -d secrets.enc.yaml | psm sync --store ssm /dev/stdin
```

## Key Design Decisions

1. **`TtyOpener` field on `IOStreams`** — Dependency-injected function that returns a reader for the terminal. Production uses `os.Open("/dev/tty")`, tests inject a `strings.Reader`.

2. **Approval flow logic** (in `runSync`):
   - `--skip-approve` set → skip prompt (unchanged)
   - stdin is terminal → read from stdin (unchanged)
   - stdin is piped → open `/dev/tty` via `TtyOpener`, read from it
   - `/dev/tty` unavailable → auto-cancel with message to stderr

3. **No changes to `promptApprove`** — It already takes `io.Reader`; only the caller changes which reader to pass.

## Files Modified

| File | Change |
|------|--------|
| `store.go` | Add `TtyOpener func() (io.ReadCloser, error)` to `IOStreams` |
| `main.go` | Wire `TtyOpener` in production; update approval flow in `runSync` |
| `sync_test.go` | Add tests for piped-stdin-with-tty and no-tty-available scenarios |

## Testing

```bash
# Run all tests
go test ./...

# Manual verification
echo "test: value" | go run . sync --store ssm /dev/stdin
# Should show plan and prompt "Proceed? [y/N]" via terminal
```
