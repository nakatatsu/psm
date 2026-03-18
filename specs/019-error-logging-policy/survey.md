# Survey: Define Error Handling and Logging Policy

**Date**: 2026-03-17
**Spec**: specs/019-error-logging-policy/spec.md

## Summary

The spec proposes creating a policy document at `documents/backend/error-handling-and-logging.md`. Three significant findings emerged:

1. **The existing document contains rules from a different project** (web server with constructor injection, user actions like "payment", "download goods"). These must be stripped and replaced with CLI-appropriate rules — not appended to.
2. **The constitution already defines error handling conventions** (line 54). The policy document must complement, not duplicate or contradict, the constitution. The relationship between the two needs to be explicit.
3. **slog is not currently used anywhere in the codebase.** All output is raw `fmt.Fprintf`. The policy defines the target state; implementation is a separate issue.

No fundamental problems with the spec's direction. The approach (define policy first, implement later) is sound.

## S1: Existing Document Contains Foreign Rules

**Category**: Problem Reframing — Hidden Assumptions
**Finding**: `documents/backend/error-handling-and-logging.md` currently contains rules from a web server project:
- "Use constructor injection, and use method `XXXContext`" — web server pattern (request-scoped loggers). Not applicable to a CLI tool.
- "Always log important user's action. e.g., payment, add point, download goods" — e-commerce domain. Irrelevant to psm.
- "Never log full request/response bodies" — HTTP context. psm has no request/response.

**Recommendation**: The policy document should be written from scratch for psm's CLI context, not built on top of the existing content. The existing content should be fully replaced, not amended.
**Evidence**: Reading `documents/backend/error-handling-and-logging.md` — the rules reference web server patterns (constructor injection, XXXContext, request bodies, payment actions) that have no equivalent in a CLI tool.

## S2: Constitution Overlap

**Category**: Integration & Governance — Constitution Compliance
**Finding**: The constitution (v3.0.0) already specifies at line 54:
> Error handling follows Go conventions: return `error`, wrap with `fmt.Errorf("context: %w", err)`. Do not panic in library code.

The spec's FR-006 (return `(..., error)`), FR-007 (no panic), and FR-008 (wrapping with `%w`) overlap directly with this existing rule.

**Recommendation**: The policy document should reference the constitution as the authority for these rules, then extend with details the constitution doesn't cover (log levels, output routing, sensitive data, exit codes). The policy document is a detailed operational guide; the constitution is the governing principle. No constitution amendment is needed.
**Evidence**: `.specify/memory/constitution.md` line 54, compared with spec FR-006/007/008.

## S3: slog Constructor Injection vs. Default Logger

**Category**: Solution Evaluation — Approach Alternatives
**Finding**: The existing doc mandates "constructor injection" and "XXXContext" methods for slog. For a web server with request-scoped loggers carrying trace IDs, this is appropriate. For psm (a short-lived CLI tool with no concurrent requests), constructor injection adds complexity with no benefit.

Two approaches for a CLI tool:

| Approach | Pros | Cons |
|----------|------|------|
| **Default logger** (`slog.SetDefault` in main, use `slog.Info` etc. globally) | Simple, matches Constitution Principle I, sufficient for CLI | Cannot have per-request context (not needed) |
| **Constructor injection** (pass `*slog.Logger` to every function) | Testable logger output, per-context attributes | Over-engineered for CLI, violates Simplicity First |

**Recommendation**: Use `slog.SetDefault()` in main. Use the package-level functions (`slog.Info`, `slog.Error`, etc.) throughout. No constructor injection. This aligns with Constitution Principle I (Simplicity First). If testing log output is needed, tests can replace the default handler temporarily.
**Evidence**: Constitution Principle I — "Every design decision defaults to the simplest viable option." psm has no concurrent request handling, no need for request-scoped loggers.

## S4: Log Format — Text vs. JSON

**Category**: Solution Evaluation — Cost & Complexity
**Finding**: The spec mentions logs should be "machine-parseable" (US2) but does not specify the format. Two options:

| Format | Pros | Cons |
|--------|------|------|
| **Text** (`slog.TextHandler`) | Human-readable in terminal, sufficient for CLI | Harder to parse programmatically |
| **JSON** (`slog.JSONHandler`) | Machine-parseable, grep-friendly | Noisy for interactive use |

For a CLI tool that humans run in a terminal, text format is the natural default. JSON is useful when output is consumed by log aggregators — unlikely for a local CLI tool.

**Recommendation**: Default to text format (`slog.TextHandler`). Do not add a `--json-log` flag unless a concrete need arises (YAGNI). The policy should specify text format as the default.
**Evidence**: Constitution Principle II (YAGNI) — no current need for JSON log output. psm is run interactively or in simple CI pipelines.

## S5: Scope Boundary — coding-standard.md

**Category**: Scope Boundaries
**Finding**: The spec assumes the policy applies to psm only, and `documents/backend/coding-standard.md` "may reference this policy but is not updated in this issue." However, `coding-standard.md` was recently stripped to contain only logging and error handling content — making it effectively a duplicate stub of the very document this issue creates.

**Recommendation**: Either delete `coding-standard.md` entirely (it now contains a subset of what `error-handling-and-logging.md` will cover), or repurpose it for non-error/logging standards in the future. The policy document should not need to worry about keeping two files in sync. This cleanup can happen in this issue since it's documentation-only work.
**Evidence**: `documents/backend/coding-standard.md` was stripped earlier in this conversation to contain only logging and error handling rules — the same scope as `error-handling-and-logging.md`.

## Items Requiring PoC

None. This is a policy/documentation issue — no technical verification needed.

## Constitution Impact

No amendments required. The constitution already covers error handling at the principle level (line 54). The policy document extends these principles with operational detail (log levels, output routing, sensitive data) without contradicting them.

## Recommendation

Proceed to `/speckit.plan`. Address the following during planning:

1. **S1**: Replace existing document content entirely (do not amend web server rules)
2. **S2**: Reference constitution as authority, extend with operational details
3. **S3**: Specify default logger pattern, not constructor injection
4. **S4**: Specify text format as default
5. **S5**: Decide whether to delete or repurpose `coding-standard.md`
