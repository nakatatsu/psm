# Survey: Read Approval Prompt from /dev/tty

**Date**: 2026-03-27
**Spec**: [spec.md](spec.md)

## Summary

The spec correctly identifies the problem and proposes the standard Unix solution. The current code (`main.go:162-164`) silently returns success when stdin is not a terminal and `--skip-approve` is not set — the user gets no feedback and no changes are applied. The `/dev/tty` approach is the established Unix convention. The change is small, well-scoped, and low-risk. No alternative approaches offer meaningful advantages.

## S1: Problem Validation

**Category**: A. Problem Reframing
**Finding**: The problem is real and correctly framed. The primary use case (`sops -d ... | psm sync ...`) is broken because the approval flow reads from stdin, which contains pipe data. The spec correctly identifies this as a problem with the input source, not with the approval mechanism itself.
**Recommendation**: Proceed as specified.
**Evidence**: `main.go:70-76` creates IOStreams with `os.Stdin` for both data reading and approval prompts. `main.go:163-164` silently exits when `IsTerminal()` returns false.

## S2: Approach Alternatives

**Category**: B. Solution Evaluation
**Finding**: Three approaches considered:
1. **`/dev/tty` (proposed)** — Standard Unix pattern. Used by `rm -i`, `git add -p`, `ssh`, `sudo`. Zero dependencies. Works on all Unix platforms.
2. **Separate fd/flag for approval input** — User passes approval source via flag (e.g., `--approve-from /dev/tty`). Unnecessary complexity; violates Simplicity First principle.
3. **Always require `--skip-approve` for piped input** — Current behavior, which the Issue explicitly identifies as broken.
**Recommendation**: `/dev/tty` is the correct and only sensible approach.
**Evidence**: Go stdlib `os.Open("/dev/tty")` is well-documented and reliable. No third-party dependency needed.

## S3: Implementation Impact

**Category**: D. Integration & Governance
**Finding**: The change is minimal and well-contained:
- `IOStreams` struct needs a new field or the `promptApprove` caller needs to open `/dev/tty` as the reader.
- `IsTerminal()` check at `main.go:163` changes meaning: instead of "is stdin a terminal?" it becomes "can we get interactive input?" (try `/dev/tty`).
- Existing tests use `testIOStreams()` helper (`sync_test.go:140`) which injects a reader — this pattern supports the change naturally.
- No changes needed outside the approval flow.
**Recommendation**: Modify the approval flow in `main.go` (the `runSync` function) to open `/dev/tty` when `IsTerminal()` is false. Keep `promptApprove` signature unchanged — just pass it a different reader.
**Evidence**: `sync_test.go:140-147` shows testable IOStreams injection. `promptApprove` at `sync.go:67` takes `io.Reader` — any reader works.

## S4: Risk Assessment

**Category**: C. Risk & Feasibility
**Finding**: Low risk. The only failure mode is `/dev/tty` not being available, which the spec already addresses (auto-cancel with message). Current behavior for this case is worse (silent no-op). The change is fully backward-compatible: when stdin IS a terminal, behavior is unchanged.
**Recommendation**: No PoC needed. The pattern is well-established and the Go implementation is trivial.
**Evidence**: Go stdlib `os.Open("/dev/tty")` returns a clear error when unavailable. `rm -i` has used this pattern for decades.

## S5: Constitution Compliance

**Category**: D. Constitution Compliance
**Finding**: Fully compliant.
- Standard library only (Principle I: Simplicity First)
- No speculative features (Principle II: YAGNI)
- Test-first workflow applies naturally (Principle III)
- Small, focused change suitable for single PR
**Recommendation**: No constitution amendments needed.

## Items Requiring PoC

None. The `/dev/tty` pattern is well-established and requires no experimentation.

## Constitution Impact

No amendments required. The change aligns with all constitution principles.

## Recommendation

Proceed to `/speckit.plan`. This is a clean, low-risk change with a well-understood solution.
