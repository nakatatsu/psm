# Tasks: Version Display Command

**Input**: Design documents from `/specs/030-version-display/`
**Prerequisites**: plan.md, spec.md, research.md

**Tests**: Constitution III (Test-First) requires tests before implementation.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2)
- Include exact file paths in descriptions

---

## Phase 1: Setup

**Purpose**: Prepare version variable and GoReleaser configuration

- [x] T001 Add `version` package variable with `dev` default in main.go
- [x] T002 Add `-X main.version={{.Version}}` to ldflags in .goreleaser.yaml

---

## Phase 2: User Story 1 - Check installed version (Priority: P1) 🎯 MVP

**Goal**: `psm --version` でバージョン文字列を表示する

**Independent Test**: `psm --version` を実行し、出力が `psm version <version>` 形式であることを確認

### Tests for User Story 1

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T003 [US1] Add table-driven tests for `--version` flag in main_test.go: (1) `--version` returns version output and no error, (2) `--version` with extra args still returns version

### Implementation for User Story 1

- [x] T004 [US1] Add `--version` check in `parseArgs()` before subcommand parsing in main.go
- [x] T005 [US1] Add version display handling in `run()` or `main()` in main.go

**Checkpoint**: `psm --version` が `psm version dev` を出力し、exit code 0 で終了する

---

## Phase 3: User Story 2 - Version check in CI/CD pipelines (Priority: P2)

**Goal**: CI/CD パイプラインでバージョンを確認できる（出力が単一行・exit code 0）

**Independent Test**: `psm --version` の出力が単一行で、exit code が 0 であることをスクリプトで確認

### Tests for User Story 2

- [x] T006 [US2] Add test verifying `--version` output is single line and exit code is 0 in main_test.go

**Checkpoint**: US1 のテストと合わせて全テスト通過。US2 は US1 の実装で既にカバーされるため、テスト追加のみ

---

## Phase 4: Polish & Cross-Cutting Concerns

**Purpose**: 既存機能への回帰がないことを確認

- [x] T007 Run `go test ./...` and `go vet ./...` to verify no regressions

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: T001, T002 は並行可能（異なるファイル）
- **US1 (Phase 2)**: T001 完了後に T003 → T004 → T005 の順（Test-First）
- **US2 (Phase 3)**: T005 完了後に T006
- **Polish (Phase 4)**: 全タスク完了後

### Within Each User Story

- Tests MUST be written and FAIL before implementation (Constitution III)
- T003 (test) → T004, T005 (implementation) の順序は厳守

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001, T002)
2. Complete Phase 2: US1 tests then implementation (T003 → T004 → T005)
3. **STOP and VALIDATE**: `psm --version` outputs `psm version dev`

### Incremental Delivery

1. Setup → US1 → Validate (MVP)
2. Add US2 test (T006) → Validate
3. Polish (T007) → Done

---

## Notes

- Total tasks: 7
- US1: 3 tasks (1 test + 2 implementation)
- US2: 1 task (test only — implementation is covered by US1)
- Setup: 2 tasks
- Polish: 1 task
- Parallel opportunities: T001 and T002 (different files)
