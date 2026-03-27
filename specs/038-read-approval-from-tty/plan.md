# Implementation Plan: Read Approval Prompt from /dev/tty

**Branch**: `038-read-approval-from-tty` | **Date**: 2026-03-27 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/038-read-approval-from-tty/spec.md`

## Summary

When stdin is piped, read the approval prompt from `/dev/tty` instead of stdin. This enables the primary use case (`sops -d ... | psm sync ...`) to work with interactive approval. Falls back to auto-cancel with a message when no terminal is available.

## Technical Context

**Language/Version**: Go (latest stable)
**Primary Dependencies**: Standard library only (`os`, `bufio`)
**Storage**: N/A
**Testing**: `go test` (standard library)
**Target Platform**: Linux/macOS (Unix systems with `/dev/tty`)
**Project Type**: CLI tool
**Performance Goals**: N/A (interactive prompt, not performance-sensitive)
**Constraints**: Must not break existing non-piped usage

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| # | Principle | Gate Question | Pass? |
|---|-----------|---------------|-------|
| I | Simplicity First | Is this the simplest viable design? No unnecessary abstractions or dependencies? | YES — `os.Open("/dev/tty")` is the simplest solution, stdlib only |
| II | YAGNI | Does every element serve a present, concrete need? No speculative features? | YES — solves the exact reported problem, nothing more |
| III | Test-First (NON-NEGOTIABLE) | Are all tests written before implementation (Red-Green cycle)? Using `go test` only? No third-party test frameworks? Stub/integration tests use DI (no test branches in production code)? | YES — existing `IOStreams` injection pattern supports test-first. New field for TTY reader injected via DI. |

## Project Structure

### Documentation (this feature)

```text
specs/038-read-approval-from-tty/
├── spec.md
├── survey.md
├── plan.md              # This file
├── research.md          # Phase 0 output
├── quickstart.md        # Phase 1 output
└── tasks.md             # Phase 2 output (created by /speckit.tasks)
```

### Source Code (repository root)

```text
# Flat package structure (existing)
main.go          # runSync approval flow — primary change site
store.go         # IOStreams struct — add TtyOpener field
sync.go          # promptApprove — no signature change needed
sync_test.go     # approval flow tests — add piped-stdin + tty tests
main_test.go     # CLI integration tests
```

**Structure Decision**: No new files or packages needed. All changes are within existing files, consistent with the flat package structure.
