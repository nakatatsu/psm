# Tasks: Distinguish dry-run output from actual execution

**Input**: Design documents from `/specs/008-dryrun-output/`
**Prerequisites**: plan.md (required), spec.md (required)

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1)

---

## Phase 1: User Story 1 - dry-run 出力を区別可能にする (Priority: P1) 🎯 MVP

**Goal**: `--dry-run` の出力に `(dry-run)` 表記を追加し、通常実行と明確に区別する

**Independent Test**: `go test -run TestExecuteDryRun` で dry-run 表記の有無を検証

### Tests (Test-First)

- [ ] T001 [US1] Write test: dry-run 時のアクション行に `(dry-run)` プレフィックスが含まれることを検証する in `sync_test.go`
- [ ] T002 [US1] Write test: dry-run 時のサマリー行に `(dry-run)` が含まれることを検証する in `sync_test.go`
- [ ] T003 [US1] Write test: 通常実行時の出力に `(dry-run)` が含まれないことを検証する in `sync_test.go`
- [ ] T004 [US1] Confirm tests fail (Red phase)

### Implementation

- [ ] T005 [US1] Modify `execute` function: dry-run 時のアクション行に `(dry-run)` プレフィックスを追加する in `sync.go`
- [ ] T006 [US1] Modify `execute` function: dry-run 時のサマリー行に `(dry-run)` サフィックスを追加する in `sync.go`
- [ ] T007 [US1] Confirm all tests pass (Green phase)

**Checkpoint**: `go test ./...` が全てパス、dry-run の出力に `(dry-run)` が表示される

---

## Dependencies & Execution Order

### Phase Dependencies

- T001-T003: テスト作成（並列可能だが同一ファイルのため順次推奨）
- T004: T001-T003 完了後に実行
- T005-T006: T004 完了後（Red 確認後）に実装
- T007: T005-T006 完了後に実行

---

## Notes

- Total tasks: 7
- Single user story — Test-First サイクル 1 回で完了
- 変更対象ファイル: `sync.go`, `sync_test.go` の 2 ファイルのみ
