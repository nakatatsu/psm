# Implementation Plan: Integration Test Script

**Branch**: `021-verify-quickstart-scenarios` | **Date**: 2026-03-25 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/021-verify-quickstart-scenarios/spec.md`

## Summary

`example/test.sh` として、psm の全動作シナリオ（dry-run, sync, delete, conflict, debug, --prune 拒否, 非端末 auto-decline, no-changes）を AWS sandbox に対して自動実行するシェルスクリプトを作成する。

## Technical Context

**Language/Version**: Bash (POSIX-compatible shell script)
**Primary Dependencies**: psm, aws CLI, sops, age (all pre-installed in DevContainer)
**Storage**: AWS SSM Parameter Store (sandbox account)
**Testing**: Shell script self-test (exit codes + output assertions)
**Target Platform**: DevContainer (Debian bookworm-slim)
**Project Type**: CLI integration test script
**Performance Goals**: Full test suite completes in under 2 minutes
**Constraints**: Requires active AWS SSO session; no internet-free execution
**Scale/Scope**: 8 test scenarios, single script file

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| # | Principle | Gate Question | Pass? |
|---|-----------|---------------|-------|
| I | Simplicity First | Is this the simplest viable design? No unnecessary abstractions or dependencies? | YES — single bash script, no frameworks |
| II | YAGNI | Does every element serve a present, concrete need? No speculative features? | YES — each scenario maps to a requirement in specs/requirements.md |
| III | Test-First (NON-NEGOTIABLE) | Are all tests written before implementation? | N/A — this IS the test artifact. No Go code is introduced. |

## Project Structure

### Documentation (this feature)

```text
specs/021-verify-quickstart-scenarios/
├── plan.md
├── research.md
├── quickstart.md
└── checklists/
    └── requirements.md
```

### Source Code (repository root)

```text
example/
└── test.sh              # Integration test script (new file)
```

**Structure Decision**: Single script file in `example/` alongside the existing DevContainer and README. No new directories or test frameworks.
