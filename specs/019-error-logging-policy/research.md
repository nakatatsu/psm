# Research: Define Error Handling and Logging Policy

**Date**: 2026-03-17

## R1: slog Handler Pattern for CLI Tools

**Decision**: Use `slog.SetDefault()` in `main()` with `slog.TextHandler` writing to stderr. Use package-level functions (`slog.Info`, `slog.Warn`, `slog.Error`, `slog.Debug`) throughout the codebase.
**Rationale**: psm is a short-lived CLI tool with no concurrent requests. Constructor injection (passing `*slog.Logger` to every function) adds complexity with no benefit. The default logger pattern is idiomatic for CLI tools and aligns with Constitution Principle I (Simplicity First). Survey S3 confirmed this.
**Alternatives**:
- Constructor injection (`*slog.Logger` parameter) — over-engineered for CLI, appropriate for web servers with request-scoped loggers. Rejected per Constitution Principle I.
- `log.Printf` — no log levels, no structured fields. Insufficient for #18's needs (warnings, debug output).

## R2: Log Output Format

**Decision**: Text format (`slog.TextHandler`) as default. No JSON format option.
**Rationale**: psm is run interactively in a terminal or in simple CI pipelines. Text format is human-readable. JSON format is useful for log aggregators, which psm does not target. Adding a `--json-log` flag would violate Constitution Principle II (YAGNI). Survey S4 confirmed this.
**Alternatives**:
- JSON format (`slog.JSONHandler`) — machine-parseable but noisy for interactive use. No current need.
- Configurable format flag — YAGNI. Can be added later if a concrete need arises.

## R3: Log Level Definitions

**Decision**: Four levels with the following usage guidelines:

| Level | Usage | Example |
|-------|-------|---------|
| Error | Operation failed, cannot continue for this item | `slog.Error("failed to put parameter", "key", key, "error", err)` |
| Warn | Non-fatal issue that deserves attention | `slog.Warn("unmanaged key detected", "key", key)` |
| Info | Normal operational milestones | `slog.Info("sync complete", "created", 3, "updated", 1)` |
| Debug | Diagnostic details for troubleshooting | `slog.Debug("API call", "operation", "PutParameter", "key", key)` |

Default level: Info (Error + Warn + Info visible; Debug hidden).
**Rationale**: Standard slog levels. The distinction between Error and Warn is critical for #18 (unmanaged keys = Warn, API failures = Error). Debug exists for troubleshooting but is hidden by default (YAGNI for normal usage).
**Alternatives**: None considered — these are the standard slog levels.

## R4: Output Routing Policy

**Decision**: Strict separation:
- **stdout**: Program output only — diff lines, summaries, action plans. This is the "result" of the command, pipeable to other tools.
- **stderr**: Everything else — slog messages (all levels), interactive prompts (approve `y/N`), progress indicators.

**Rationale**: CLI convention. stdout is for data, stderr is for diagnostics and interaction. Mixing them breaks piping (`psm sync ... | grep create`). The user confirmed this direction in conversation.
**Alternatives**: None — this is the standard Unix convention for CLI tools.

## R5: Relationship to Constitution

**Decision**: The policy document references the constitution as the governing authority for core error handling rules, then extends with operational details the constitution doesn't cover.

| Topic | Constitution | Policy Document |
|-------|-------------|-----------------|
| Return `error` | ✓ (line 54) | References constitution, adds no new rules |
| Wrap with `%w` | ✓ (line 54) | References constitution, adds wrapping guidance (context at each call site) |
| No panic | ✓ (line 54, "library code") | Extends: no panic anywhere, not just library code |
| Log levels | ✗ | Defines Error/Warn/Info/Debug with examples |
| Output routing | ✗ | Defines stdout vs stderr rules |
| Sensitive data | ✗ | Defines what must never be logged |
| Exit codes | ✗ | Defines 0 = success, 1 = failure |
| Logger choice | ✗ | Mandates slog |

**Rationale**: The constitution is the highest authority (Governance section). Duplicating its rules creates maintenance burden and risk of divergence. The policy document should say "per constitution" and focus on what the constitution doesn't cover.
**Alternatives**: Copy constitution rules into the policy doc — rejected because of duplication risk.

## R6: coding-standard.md Disposition

**Decision**: Delete `documents/backend/coding-standard.md`.
**Rationale**: It was stripped earlier in this conversation to contain only logging and error handling rules — the exact scope of `error-handling-and-logging.md`. Keeping both creates confusion about which is authoritative. Survey S5 identified this.
**Alternatives**: Keep it for future non-error/logging standards — rejected because no such standards exist or are planned (YAGNI). Can be recreated if needed.
