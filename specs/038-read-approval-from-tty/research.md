# Research: Read Approval Prompt from /dev/tty

**Date**: 2026-03-27

## R1: /dev/tty Pattern in Go

**Decision**: Use `os.Open("/dev/tty")` to get an `io.Reader` for interactive input when stdin is piped.

**Rationale**: This is the standard Unix approach. Go's `os.Open` works directly with `/dev/tty`. The opened file implements `io.Reader`, which is exactly what `promptApprove` expects. No third-party dependencies needed.

**Alternatives considered**:
- `golang.org/x/term` package — unnecessary dependency for this use case; we only need to read a line, not raw terminal control.
- Reading from `/proc/self/fd/0` — non-portable, doesn't solve the problem (still stdin).

## R2: IOStreams Injection Strategy

**Decision**: Add a `TtyOpener` function field to `IOStreams` that returns an `io.ReadCloser` for the terminal. In production, this opens `/dev/tty`. In tests, it returns the injected reader.

**Rationale**: Follows the existing DI pattern (`IsTerminal func() bool`). Keeps `promptApprove` signature unchanged — only the caller decides which reader to pass. No test branches in production code (constitution compliance).

**Alternatives considered**:
- Global function variable — violates DI pattern, harder to test.
- Interface with TTY method — over-abstraction for a single function call.
- Passing `/dev/tty` path as config — unnecessary indirection.

## R3: Fallback Behavior

**Decision**: When `/dev/tty` cannot be opened (error), auto-cancel with a message to stderr: "No terminal available for approval prompt. Use --skip-approve for non-interactive usage."

**Rationale**: This is safer than the current silent no-op. Users get actionable feedback. Matches the behavior of `sudo` when no tty is available.

**Alternatives considered**:
- Exit with non-zero code — too aggressive for "no changes made" scenario.
- Silent cancel (current behavior) — the exact problem being fixed.

## R4: When to Use /dev/tty vs stdin

**Decision**: Use `/dev/tty` only when `IsTerminal()` returns false (stdin is piped). When stdin is already a terminal, continue reading from stdin as before.

**Rationale**: Avoids unnecessary file descriptor usage. Preserves exact backward compatibility for the common non-piped case. Opening `/dev/tty` when stdin is already a terminal would work but is redundant.
