# Tasks: ブランチ戦略の決定と実装

**Input**: Design documents from `/specs/045-branch-strategy/`
**Prerequisites**: plan.md, spec.md, survey.md, research.md

**Tests**: Go コードの変更なし。検証は `terraform plan` と GitHub UI で行う。

**Organization**: Tasks are grouped by user story to enable independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3, US4)

---

## Phase 1: Setup

**Purpose**: 実装に必要な前提条件の確認と準備

- [x] T001 既存 Terraform コードの現状確認（.tmp/mynote/infrastructure/github/ の全ファイル）
- [x] T002 `psm` リポジトリが `protected_repositories` 変数に含まれていることを確認（.tmp/mynote/infrastructure/github/variables.tf）

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 全 User Story の前提となる既存 Ruleset の修正

**⚠️ CRITICAL**: この修正が完了しないと GitFlow の merge commit フローが動かない

- [x] T003 既存 `main` ブランチ Ruleset の `required_linear_history` を `false` に変更（.tmp/mynote/infrastructure/github/main.tf L23）

**Checkpoint**: `terraform plan` で `main` Ruleset の `required_linear_history` 変更のみが差分として表示される

---

## Phase 3: User Story 1 & 2 - ブランチ保護ルールの設定と GitHub への実装 (Priority: P1) 🎯 MVP

**Goal**: `develop`, `release-*`, `hotfix-*` ブランチに保護ルールを適用し、GitFlow の品質ゲートを確立する

**Independent Test**: `terraform plan` で意図した Ruleset が作成され、`terraform apply` 後に GitHub 上で直接プッシュが拒否されることを確認

### Implementation

- [x] T004 [P] [US1] `develop` ブランチ Ruleset を追加（.tmp/mynote/infrastructure/github/main.tf）: PR必須（承認0名）、CI ステータスチェック必須（strict）、フォースプッシュ禁止、削除禁止、admin バイパス
- [x] T005 [P] [US1] `release-*` ブランチ Ruleset を追加（.tmp/mynote/infrastructure/github/main.tf）: PR必須（1名承認、stale review 却下、last push approval、コメント解決必須）、CI ステータスチェック必須（strict）、フォースプッシュ禁止、削除禁止、admin バイパス
- [x] T006 [P] [US1] `hotfix-*` ブランチ Ruleset を追加（.tmp/mynote/infrastructure/github/main.tf）: `release-*` と同一設定
- [x] T007 [US2] `terraform plan` で差分を確認し、既存の `main` Ruleset・タグ保護ルールに影響がないことを検証

**Checkpoint**: `terraform plan` で `develop`, `release-*`, `hotfix-*` の 3 Ruleset 追加 + `main` の `required_linear_history` 変更のみが差分

---

## Phase 4: User Story 3 - PR マージ戦略とリポジトリ設定 (Priority: P2)

**Goal**: Merge commit のみ許可、セキュリティ設定、ブランチ自動削除を適用する

**Independent Test**: `terraform plan` で `github_repository` リソースの設定変更が表示される

### Implementation

- [x] T008 [US3] `github_repository` リソースを新規作成（.tmp/mynote/infrastructure/github/repository.tf）: `allow_merge_commit = true`, `allow_squash_merge = false`, `allow_rebase_merge = false`, `delete_branch_on_merge = true`, `vulnerability_alerts = true`, `security_and_analysis` ブロック（`advanced_security = "enabled"`, `secret_scanning = "enabled"`, `secret_scanning_push_protection = "enabled"`）
- [x] T009 [US3] 既存 `psm` リポジトリを `terraform import github_repository.psm psm` で state に取り込む手順を文書化
- [x] T010 [US3] CodeQL ワークフローを作成（.github/workflows/codeql.yml）: Go 対象、push/PR/weekly スケジュールトリガー
- [x] T011 [US3] Terraform の `required_code_scanning` ブロックを Ruleset に追加を試みる（.tmp/mynote/infrastructure/github/main.tf）。Provider バグで失敗した場合は `required_status_checks` に CodeQL チェックを追加してフォールバック

**Checkpoint**: `terraform plan` でリポジトリ設定変更 + CodeQL 関連の差分が表示される

---

## Phase 5: User Story 4 - ブランチ戦略の SKILL 化 (Priority: P2)

**Goal**: commit-msg hook と Claude Code SKILL でブランチ戦略を開発フローに組み込む

**Independent Test**: commit-msg hook が Issue 番号なしのコミットを拒否し、SKILL ファイルが存在する

### Implementation

- [x] T012 [P] [US4] commit-msg hook スクリプトを作成（.githooks/commit-msg）: `#<数字>` パターン検証、マージコミット除外
- [x] T013 [P] [US4] ブランチ戦略 SKILL を作成（.claude/skills/branch-strategy/）: GitFlow フロー、ブランチ命名規則（`feature/<issue-no>-<short-description>`, `release-<version>`, `hotfix-<version>`）、Issue 番号付きコミットメッセージの指示

**Checkpoint**: `git commit -m "test"` が拒否され、`git commit -m "test #45"` が成功する

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: ADR 記録とドキュメント更新

- [x] T014 [P] ADR を記録（.specify/decisions/）: ブランチ戦略（GitFlow 採用）、ブランチ命名規則の 2 件
- [x] T015 [P] CLAUDE.md にブランチ戦略の要約を追記（`core.hooksPath` 設定手順を含む）
- [x] T016 quickstart.md の手順に沿って全体動作を検証

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup - BLOCKS all user stories
- **US1&2 (Phase 3)**: Depends on Foundational
- **US3 (Phase 4)**: Depends on Foundational. Can run in parallel with Phase 3
- **US4 (Phase 5)**: No dependency on Terraform phases. Can run in parallel with Phase 3 & 4
- **Polish (Phase 6)**: Depends on all previous phases

### Parallel Opportunities

- T004, T005, T006 can run in parallel (different Ruleset リソース)
- Phase 3 (Ruleset) と Phase 4 (リポジトリ設定) は並行可能
- Phase 5 (hook + SKILL) は Terraform 作業と完全に独立して並行可能
- T012, T013 can run in parallel (different files)
- T014, T015 can run in parallel (different files)

---

## Implementation Strategy

### MVP First (Phase 1-3)

1. Setup + Foundational を完了
2. US1&2: Ruleset 追加 + `terraform plan` で検証
3. **STOP and VALIDATE**: `terraform apply` は別リポジトリ（mynote）で実施

### Incremental Delivery

1. Phase 1-3: Ruleset（MVP）→ 検証
2. Phase 4: リポジトリ設定 + CodeQL → 検証
3. Phase 5: hook + SKILL → 検証
4. Phase 6: ADR + ドキュメント → 完了

---

## Notes

- Terraform コードは `.tmp/mynote/infrastructure/github/` に配置（別リポジトリ管理）
- `terraform apply` はこのリポジトリ（psm）からは実行しない
- CodeQL ワークフローと commit-msg hook は psm リポジトリに配置
- [P] tasks = different files, no dependencies
