# Tasks: --store フラグの除去

**Input**: Design documents from `/specs/037-remove-store-flag/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Constitution (Test-First NON-NEGOTIABLE) に従い、テストを先に更新する。

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup

**Purpose**: 準備作業

- [x] T001 現在のテストがすべてパスすることを確認 (`go test ./...`)

---

## Phase 2: Foundational (Config 構造体の変更)

**Purpose**: Config から Store フィールドを除去する（全ストーリーの前提）

- [x] T002 Config 構造体から `Store string` フィールドを削除する in store.go
- [x] T003 `run()` 関数の switch 文を削除し、`NewSSMStore(awsCfg)` を直接呼び出すように変更する in main.go

**Checkpoint**: ビルドが通ることを確認（テストは一部失敗してよい）

---

## Phase 3: User Story 1 - フラグなしでの sync/export 実行 (Priority: P1) MVP

**Goal**: `--store` フラグなしで `psm sync` / `psm export` が SSM に対して正常動作する

**Independent Test**: `psm sync file.yaml` および `psm export out.yaml` がフラグなしで正常にパースされることを確認

### Tests for User Story 1

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T004 [US1] main_test.go のテストケースを更新: `--store ssm` を除去し、フラグなしでの正常パースを期待するテストに変更する in main_test.go

### Implementation for User Story 1

- [x] T005 [US1] `parseArgs()` から `--store` フラグ定義 (`fs.String("store", ...)`) を削除する in main.go
- [x] T006 [US1] `parseArgs()` から `--store` バリデーション（空チェック・値チェック）を削除する in main.go
- [x] T007 [US1] `parseArgs()` の戻り値から `Store: *store` を削除する in main.go
- [x] T008 [US1] usage メッセージから `--store ssm` の記述を削除する in main.go
- [x] T009 [US1] `go test ./...` を実行し全テストがパスすることを確認

**Checkpoint**: `--store` なしでの CLI パースが正常動作

---

## Phase 4: User Story 2 - --store 指定時の廃止エラー (Priority: P2)

**Goal**: `--store` を指定した場合に分かりやすい廃止エラーメッセージを表示する

**Independent Test**: `--store ssm` や `--store sm` を指定した場合にエラーが返ることを確認

### Tests for User Story 2

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T010 [US2] main_test.go に `--store` 指定時のエラーテストケースを追加する（`--store ssm`, `--store sm` の両方）in main_test.go

### Implementation for User Story 2

- [x] T011 [US2] `parseArgs()` に `--store` 検出ロジックを追加する（`--prune` 廃止と同じパターン）in main.go
- [x] T012 [US2] `go test ./...` を実行し全テストがパスすることを確認

**Checkpoint**: `--store` 指定時に廃止エラーが表示される

---

## Phase 5: User Story 3 - ドキュメントの一貫性 (Priority: P3)

**Goal**: ドキュメントから `--store` フラグの記述を除去する

**Independent Test**: ドキュメント内に `--store` の記述がないことを確認

### Implementation for User Story 3

- [x] T013 [P] [US3] README.md から `--store ssm` の記述を除去し、コマンド例を更新する in README.md
- [x] T014 [P] [US3] example/README.md の CLI Reference とコマンド例から `--store` を除去する in example/README.md
- [x] T015 [P] [US3] example/test.sh の全シナリオから `--store ssm` を除去する in example/test.sh

**Checkpoint**: 全ドキュメントから `--store` が除去されている

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: 最終確認

- [x] T016 `go test ./...` で全テストがパスすることを確認
- [x] T017 `go build ./...` でビルドが成功することを確認
- [x] T018 `go vet ./...` で警告がないことを確認

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies
- **Foundational (Phase 2)**: Depends on Phase 1
- **US1 (Phase 3)**: Depends on Phase 2
- **US2 (Phase 4)**: Depends on Phase 3（`--store` 定義除去後に廃止チェックを追加）
- **US3 (Phase 5)**: Can start after Phase 2（ドキュメントはコード変更と独立）
- **Polish (Phase 6)**: Depends on all phases

### Parallel Opportunities

- T013, T014, T015（ドキュメント更新）は互いに並列実行可能
- Phase 5（ドキュメント）は Phase 3/4（コード変更）と並列実行可能

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Phase 1: Setup → 既存テスト確認
2. Phase 2: Config 変更
3. Phase 3: US1 → テスト先行で `--store` 除去
4. **STOP and VALIDATE**: `go test ./...` パス確認

### Incremental Delivery

1. Setup + Foundational → Config 準備完了
2. US1: フラグ除去 → テスト → 動作確認
3. US2: 廃止エラー追加 → テスト → 動作確認
4. US3: ドキュメント更新（US1/US2 と並列可）
5. Polish: 最終確認

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story
- Constitution の Test-First を厳守: テスト更新 → Red 確認 → 実装 → Green 確認
- `--store` 廃止検出は `--prune` 廃止と同じパターン（`main.go:209-213`）を踏襲
