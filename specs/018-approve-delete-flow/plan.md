# Implementation Plan: Add Approve Flow and Replace --prune with --delete

**Branch**: `018-approve-delete-flow` | **Date**: 2026-03-17 | **Spec**: specs/018-approve-delete-flow/spec.md
**Input**: Feature specification from `/specs/018-approve-delete-flow/spec.md`

## Summary

Add a mandatory approval prompt before all sync operations (create/update/delete), replace the dangerous `--prune` flag with `--delete <file>` for regex-based scoped deletion, and add `--debug` for verbose logging. The existing plan/execute separation in the codebase naturally supports inserting an approval step. Terminal detection uses stdlib (`os.Stdin.Stat`), no new dependencies needed.

## Technical Context

**Language/Version**: Go 1.26
**Primary Dependencies**: AWS SDK for Go v2, gopkg.in/yaml.v3, `regexp` (stdlib), `log/slog` (stdlib)
**Storage**: AWS SSM Parameter Store, AWS Secrets Manager (via Store interface)
**Testing**: `go test` (unit tests + sandbox AWS integration tests)
**Target Platform**: Linux/macOS CLI
**Project Type**: CLI tool
**Performance Goals**: N/A (batch operations, not latency-sensitive)
**Constraints**: No new third-party dependencies (Constitution Principle I)
**Scale/Scope**: Single flat-package Go project, ~7 source files modified

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| # | Principle | Gate Question | Pass? |
|---|-----------|---------------|-------|
| I | Simplicity First | Is this the simplest viable design? No unnecessary abstractions or dependencies? | YES — stdlib for terminal detection, flat functions, no new packages |
| II | YAGNI | Does every element serve a present, concrete need? No speculative features? | YES — each feature (approve, delete, debug) addresses a concrete stated need |
| III | Test-First (NON-NEGOTIABLE) | Are all tests written before implementation? Using `go test` only? | YES — plan includes test-first approach for all new functions |

## Project Structure

### Documentation (this feature)

```text
specs/018-approve-delete-flow/
├── plan.md
├── spec.md
├── survey.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── cli-schema.md
├── checklists/
│   └── requirements.md
└── tasks.md
```

### Source Code (repository root)

```text
(flat package — all files at repository root)
├── main.go              # parseArgs rework (new flags), runSync rework (approve flow)
├── store.go             # Config struct update (new fields)
├── sync.go              # plan() signature change, new displayPlan(), approve prompt
├── delete.go            # NEW — parseDeletePatterns(), planDeletes(), conflict detection
├── main_test.go         # parseArgs tests updated
├── sync_test.go         # plan/displayPlan/approve tests updated
├── delete_test.go       # NEW — delete pattern parsing, conflict detection tests
├── ssm_test.go          # prune tests removed, delete integration tests added
├── sm_test.go           # prune tests removed, delete integration tests added
└── README.md            # updated usage documentation
```

**Structure Decision**: Flat package (existing pattern). New `delete.go` file for delete-specific logic (pattern parsing, conflict detection, unmanaged key identification). Keeps `sync.go` focused on plan/display/execute.
