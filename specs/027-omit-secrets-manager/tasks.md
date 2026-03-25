# Tasks: Secrets Manager 対応オミット

**Input**: Design documents from `/specs/027-omit-secrets-manager/`
**Prerequisites**: plan.md, spec.md, research.md

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup

**Purpose**: テスト準備（コード削除の前に既存テストが通ることを確認）

- [x] T001 Run `go test ./...` and `go build ./...` to confirm current state passes in repository root

**Checkpoint**: 現在の状態が正常であることを確認

---

## Phase 2: User Story 1 - SSM ユーザーの既存ワークフロー維持 (Priority: P1)

**Goal**: SM コードを削除しても SSM の全機能が従来通り動作する

**Independent Test**: `go test ./...` が全件パスし、`go build ./...` がエラーなくビルドできる

### Tests

- [x] T002 [US1] Add test case for `--store sm` returning invalid store error in main_test.go

### Implementation

- [x] T003 [US1] Delete sm.go (SMStore implementation)
- [x] T004 [P] [US1] Delete sm_test.go (SMStore tests)
- [x] T005 [US1] Remove `case "sm"` branch from store selection in main.go
- [x] T006 [US1] Update `--store` validation in main.go to accept only `ssm` (remove `sm` from valid values)
- [x] T007 [US1] Run `go test ./...` and `go build ./...` to confirm SSM functionality is intact

**Checkpoint**: SSM の全機能が動作し、`--store sm` はバリデーションエラーになる

---

## Phase 3: User Story 2 - SM 指定時の明確なエラー (Priority: P2)

**Goal**: `--store sm` 指定時に適切なエラーメッセージが返る

**Independent Test**: `--store sm` のテストケースがパスする

### Implementation

- [x] T008 [US2] Verify T002 test passes after Phase 2 changes (no additional code needed — covered by updated validation)

**Checkpoint**: `--store sm` が明確なエラーで拒否される

---

## Phase 4: User Story 3 - ドキュメントの一貫性 (Priority: P3)

**Goal**: すべてのユーザー向けドキュメントから SM 関連記述を除去

**Independent Test**: ドキュメント内に `--store sm`、`--store <ssm|sm>`、`Secrets Manager` への言及がない

### Implementation

- [x] T009 [P] [US3] Update README.md: `--store <ssm|sm>` → `--store ssm`, remove SM descriptions and With SOPS SM note
- [x] T010 [P] [US3] Update README.ja.md: same changes as README.md (Japanese version)
- [x] T011 [P] [US3] Update example/README.md: remove `sm` from CLI Reference table
- [x] T012 [P] [US3] Update CLAUDE.md: remove `AWS Secrets Manager` from Active Technologies

**Checkpoint**: 全ドキュメントが SSM のみを記述

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: 最終確認

- [x] T013 Run `go test ./...` and `go build ./...` for final verification
- [x] T014 Run `go vet ./...` for static analysis

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: No dependencies — start immediately
- **Phase 2 (US1)**: Depends on Phase 1 — core code changes
- **Phase 3 (US2)**: Depends on Phase 2 — validation already handled by US1 changes
- **Phase 4 (US3)**: Can start after Phase 1 (independent of code changes), but recommended after Phase 2
- **Phase 5 (Polish)**: Depends on all previous phases

### Parallel Opportunities

- T003 and T004 can run in parallel (deleting different files)
- T009, T010, T011, T012 can all run in parallel (different documentation files)
- Phase 4 (docs) can run in parallel with Phase 2 (code) if needed

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Confirm baseline
2. Complete Phase 2: Delete SM code, update validation
3. **STOP and VALIDATE**: `go test ./...` passes, `--store sm` returns error

### Incremental Delivery

1. Phase 1 → Baseline confirmed
2. Phase 2 (US1) → SM code removed, SSM intact
3. Phase 3 (US2) → Error handling verified (implicit from US1)
4. Phase 4 (US3) → Docs updated
5. Phase 5 → Final verification

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- This is a deletion-focused chore — most tasks are removals, not additions
- spec ディレクトリ (`specs/001-psm/`) は設計記録として保持（削除対象外）
- Store interface (`store.go`) は保持（拡張ポイント）
