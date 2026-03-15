# Tasks: psm

**Input**: Design documents from `/specs/001-psm/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Constitution Principle III (Test-First) により全タスクで Red-Green cycle 必須。テストを書く → 失敗確認 → 実装 → 成功確認。

**Testing Strategy**:
- **ユニットテスト** (AWS 不要): YAML パース・バリデーション、CLI パース、plan 関数（純粋なデータ比較）、YAML 書き出し
- **統合テスト** (Sandbox AWS): Store 実装 (SSMStore/SMStore)、execute 関数、sync/export E2E、dry-run、prune
- テストデータは `/psm-test/` プレフィックス（SSM）、`psm-test/` プレフィックス（SM）で隔離。各テストケースで setup/teardown 実施

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup

**Purpose**: Go モジュール初期化と依存解決

- [x] T001 Initialize Go module: `go mod init`, add dependencies (aws-sdk-go-v2 v1.41.3, aws-sdk-go-v2/config v1.32.11, aws-sdk-go-v2/service/ssm v1.68.2, aws-sdk-go-v2/service/secretsmanager v1.41.3, gopkg.in/yaml.v3 v3.0.1), run `go mod tidy`. Verify with `go build ./...` in go.mod

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 全 User Story が依存する型定義・YAML パース・CLI パースの実装

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T002 Define types (Config, Entry, Action with ActionType enum, Summary) and Store interface (GetAll, Put, Delete) per data-model.md and store-interface.md contract in store.go
- [x] T003 [P] Write tests for YAML parsing in yaml_test.go: valid key-values, sops key filtering (FR-006), duplicate key error, null value error, map/array value error, empty key error, zero keys error, non-string scalar extraction via ScalarNode.Value (int→string, bool→string), empty string value as valid. Use table-driven tests
- [x] T004 Implement YAML parsing with yaml.Node API in yaml.go: parse YAML, filter sops key first (FR-006), then validate per FR-020 (duplicates, null, map/array, empty key, zero keys). Return []Entry. Error messages must include problematic key names
- [x] T005 [P] Write tests for CLI parsing in main_test.go: `psm sync --store ssm file.yaml`, `psm sync --store sm --profile prod file.yaml`, `psm export --store ssm out.yaml`, no subcommand → error + exit code 1, no file arg → error + exit code 1, no --store → error + exit code 1, invalid --store value → error + exit code 1, --prune/--dry-run flags for sync (FR-012). Use table-driven tests
- [x] T006 Implement CLI subcommand dispatch and flag parsing in main.go: parse os.Args[1] as subcommand (sync/export), use flag.NewFlagSet per subcommand for --store (required, common), --profile (optional, common), --prune/--dry-run (sync only). Return Config. Show usage on error (FR-012)

**Checkpoint**: Types defined, YAML parsing works, CLI parsing works. `go test ./...` passes.

---

## Phase 3: User Story 1 — SSM Parameter Store への同期 (Priority: P1) 🎯 MVP

**Goal**: YAML → SSM Parameter Store への差分同期。create/update/skip を正しく判定し、stdout に差分行+サマリー、stderr にエラーを出力する

**Independent Test**: `psm sync --store ssm secrets.yaml` で SSM に正しく同期される

### Tests for User Story 1

> **Write these tests FIRST, ensure they FAIL before implementation**

- [x] T007 [US1] Write unit tests for plan function in sync_test.go: given []Entry and map[string]string (既存 AWS 状態を表す素のデータ), verify correct Action types — new key → create, changed value → update, same value → skip. AWS 不要（純粋なデータ比較）。Table-driven tests
- [x] T009 [US1] Write integration tests for SSMStore + execute in ssm_test.go: Sandbox AWS に接続し、テストデータ (`/psm-test/` プレフィックス) を setup/teardown。テスト項目 — SSMStore.GetAll のページネーション、SSMStore.Put (SecureString, Overwrite)、SSMStore.Delete (batch)、execute 経由の create/update/skip の stdout 出力形式 (`create: {key}`, `update: {key}`)、summary 行形式 (`N created, N updated, N deleted, N unchanged, N failed`)、values が出力に含まれないこと、部分失敗時に残りのキーが正常処理され終了コード 1 が返ること (FR-010, FR-011)。テストヘルパー (setupSSMTestData/cleanupSSMTestData) もこのタスクで作成

### Implementation for User Story 1

- [x] T008 [US1] Implement plan function in sync.go: accept []Entry and map[string]string (AWS state from Store.GetAll), return []Action. Compare YAML entries against AWS state per data-model.md state transition logic (create/update/skip)
- [x] T010 [US1] Implement execute function in sync.go: accept []Action, Store, and io.Writer (stdout/stderr), call Store.Put for create/update actions, format diff lines to stdout, errors to stderr (key + error message, never values), compute and print Summary line. Return exit code (0 = all success, 1 = any failure)
- [x] T011 [US1] Implement SSMStore in ssm.go: GetAll via GetParametersByPath (Path="/", Recursive=true, WithDecryption=true, paginate), Put via PutParameter (Type=SecureString, Overwrite=true, goroutine concurrency max 10 via buffered channel semaphore), Delete via DeleteParameters (batch max 10 keys/req). Constructor accepts aws.Config. T009 の統合テストが Green になることを確認
- [x] T012 [US1] Wire sync subcommand end-to-end in main.go: parse CLI → read YAML file → create AWS config (use --profile if specified, os.Unsetenv("AWS_PROFILE") to ignore env var per FR-002) → create SSMStore → call plan(entries, store.GetAll()) → call execute(actions, store) → os.Exit with result code

**Checkpoint**: `psm sync --store ssm secrets.yaml` で SSM 同期が動作する。create/update/skip を正しく判定し、差分とサマリーを出力する。

---

## Phase 4: User Story 2 — Secrets Manager への同期 (Priority: P2)

**Goal**: --store sm 指定で Secrets Manager に同期。sync ロジックは US1 と共有、Store 実装のみ追加

**Independent Test**: `psm sync --store sm secrets.yaml` で Secrets Manager に正しく同期される

### Tests for User Story 2

- [x] T013 [US2] Write integration tests for SMStore in sm_test.go: Sandbox AWS に接続し、テストデータ (`psm-test/` プレフィックス) を setup/teardown。テスト項目 — SMStore.GetAll (ListSecrets + BatchGetSecretValue のページネーション)、SMStore.Put (CreateSecret → ResourceExistsException → PutSecretValue)、SMStore.Delete (ForceDeleteWithoutRecovery)、execute 経由の sync 動作。テストヘルパー (setupSMTestData/cleanupSMTestData) もこのタスクで作成
### Implementation for User Story 2

- [x] T014 [US2] Implement SMStore in sm.go: GetAll via ListSecrets (paginate) + BatchGetSecretValue (max 20/req), Put via CreateSecret (catch ResourceExistsException → PutSecretValue, goroutine concurrency max 10), Delete via DeleteSecret (ForceDeleteWithoutRecovery=true, goroutine concurrency max 10). Constructor accepts aws.Config. T013 の統合テストが Green になることを確認
- [x] T015 [US2] Wire SM store selection in main.go: when --store sm, create SMStore instead of SSMStore. Verify sync logic works unchanged via Store interface

**Checkpoint**: `psm sync --store sm secrets.yaml` で Sandbox AWS の SM に同期が動作する。

---

## Phase 5: User Story 3 — Dry-run で差分を事前確認 (Priority: P3)

**Goal**: --dry-run 指定時に AWS への変更なしで差分を表示する。出力形式は通常実行と同一

**Independent Test**: `psm sync --store ssm --dry-run secrets.yaml` で差分表示のみ、AWS 変更なし

### Tests for User Story 3

- [x] T016 [US3] Write integration tests for dry-run mode in ssm_test.go: Sandbox AWS にテストデータを setup し、--dry-run で sync 実行。検証項目 — Store.GetAll は実行される（plan 計算のため）、stdout 出力は通常実行と同一形式、AWS 上のデータが変更されていないこと（teardown 前に GetAll で確認）、summary の failed=0

### Implementation for User Story 3

- [x] T017 [US3] Implement dry-run mode in sync.go: add dryRun bool parameter to execute (or top-level run function). When true, output planned actions to stdout in same format but skip Store.Put/Delete calls. Wire --dry-run flag in main.go. T016 の統合テストが Green になることを確認

**Checkpoint**: `psm sync --store ssm --dry-run` で差分表示のみ、Sandbox AWS への書き込みゼロ。

---

## Phase 6: User Story 4 — Prune で不要キーを削除 (Priority: P4)

**Goal**: --prune 指定時に YAML にないキーを AWS から削除する。アカウント全体が対象

**Independent Test**: `psm sync --store ssm --prune secrets.yaml` で YAML にないパラメータが削除される

### Tests for User Story 4

- [x] T018 [US4] Write integration tests for prune in ssm_test.go: Sandbox AWS に key1,key2,key3 を setup し、YAML に key1,key2 のみ指定。検証項目 — prune=true で key3 が AWS から削除される、prune=false で key3 が残存、stdout に `delete: {key}` 表示、summary の deleted カウント正確。SM でも同等のテストを sm_test.go に追加

### Implementation for User Story 4

- [x] T019 [US4] Implement prune in plan function in sync.go: when prune=true, iterate AWS keys not in YAML entries → generate delete Actions. Wire --prune flag in main.go. execute function handles delete Actions by calling Store.Delete. T018 の統合テストが Green になることを確認

**Checkpoint**: `psm sync --store ssm --prune` で Sandbox AWS 側の不要キーが削除される。`--dry-run --prune` で削除予定のみ表示。

---

## Phase 7: User Story 5 — Export で既存パラメータを YAML に書き出し (Priority: P5)

**Goal**: AWS 上の全パラメータ/シークレットを key-value YAML ファイルとして書き出す

**Independent Test**: `psm export --store ssm output.yaml` で全パラメータが YAML ファイルに出力される

### Tests for User Story 5

- [x] T020 [P] [US5] Write unit tests for YAML writing in yaml_test.go: given []Entry, verify output is valid YAML with correct key-value pairs, no metadata. AWS 不要（純粋なデータ変換）。Table-driven tests
- [x] T021 [P] [US5] Write integration tests for export logic in export_test.go: Sandbox AWS にテストデータを setup。検証項目 — Store.GetAll で取得 → YAML ファイル書き出し、出力ファイルが既存 → error (FR-022)、Store.GetAll が 0 件 → error (FR-023)、書き出した YAML を再 sync で差分 0 件 (SC-006 round-trip)

### Implementation for User Story 5

- [x] T022 [US5] Implement YAML writing function in yaml.go: accept []Entry, return []byte (YAML formatted). Use yaml.v3 Marshal. T020 のユニットテストが Green になることを確認
- [x] T023 [US5] Implement export logic in export.go: check output file not exists (FR-022), call Store.GetAll, check non-empty (FR-023), convert to []Entry, call YAML write, write to file. T021 の統合テストが Green になることを確認
- [x] T024 [US5] Wire export subcommand in main.go: parse CLI → create AWS config → create Store → call export → os.Exit

**Checkpoint**: `psm export --store ssm out.yaml` で Sandbox AWS の全パラメータが YAML に書き出される。書き出した YAML をそのまま `psm sync --store ssm` に渡すと差分 0 件 (SC-006 round-trip)。

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: 品質確認

- [x] T025 Run `go vet ./...` and `gofmt -l .` on all source files, fix any issues
- [x] T026 Run `go test ./...` to verify all tests pass (ユニットテスト + Sandbox AWS 統合テスト)
- [x] T027 Validate quickstart.md scenarios against Sandbox AWS: build → export → sync → update → dry-run → prune → SM sync

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Phase 1 — BLOCKS all user stories
- **US1 (Phase 3)**: Depends on Phase 2 — MVP, complete first
- **US2 (Phase 4)**: Depends on Phase 2 + SSMStore pattern from Phase 3
- **US3 (Phase 5)**: Depends on Phase 3 (sync must exist to add dry-run)
- **US4 (Phase 6)**: Depends on Phase 3 (plan function must exist to add prune)
- **US5 (Phase 7)**: Depends on Phase 2 (Store interface + YAML) — independent of US1-4
- **Polish (Phase 8)**: Depends on all user stories complete

### Within Each User Story

1. Write tests → Confirm they fail (Red)
2. Implement → Confirm tests pass (Green)
3. `go test ./...` and `go vet ./...` before moving on

### Parallel Opportunities

- Phase 2: T003 (yaml tests) and T005 (CLI tests) can run in parallel — different files
- Phase 7: T020 (YAML write tests) and T021 (export tests) can run in parallel — different files
- US5 (Phase 7) can start as soon as Phase 2 completes — independent of US1-4

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL — blocks all stories)
3. Complete Phase 3: User Story 1 (SSM sync)
4. **STOP and VALIDATE**: `psm sync --store ssm` works end-to-end
5. Deploy/demo if ready

### Incremental Delivery

1. Setup + Foundational → Foundation ready
2. Add US1 (SSM sync) → Test → **MVP!**
3. Add US2 (SM sync) → Test → Two stores supported
4. Add US3 (dry-run) → Test → Safe preview
5. Add US4 (prune) → Test → Full sync
6. Add US5 (export) → Test → Complete tool

### Recommended Order

US1 → US2 → US3 → US4 → US5 (priority order, sequential delivery)

US5 (export) は US1-4 と独立しているため、Phase 2 完了後に並行着手も可能。

---

## Notes

- Constitution Principle III: 全タスクで Red-Green cycle 厳守
- `go test` のみ使用（testify/mockgen 禁止）
- **Mock は使用しない**: psm の本質は AWS API とのやり取りであり、Store 実装がアプリの主体。Mock では何もテストしていないのと同じ
- **テスト二分類**: 純粋ロジック（YAML パース、CLI パース、plan 関数、YAML 書き出し）はユニットテスト、Store 経由の処理はすべて Sandbox AWS 統合テスト
- **テストデータ隔離**: SSM は `/psm-test/` プレフィックス、SM は `psm-test/` プレフィックス。各テストで setup/teardown
- フラットパッケージ構成（全ファイル `package main`）
- AWS_PROFILE 環境変数は os.Unsetenv で常に無視 (FR-002)
- Sandbox AWS のクレデンシャルは環境変数で管理（CI ではシークレット）
