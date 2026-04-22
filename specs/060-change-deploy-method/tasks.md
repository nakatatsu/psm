# Tasks: デプロイ方式の変更

**Input**: Design documents from `/specs/060-change-deploy-method/`
**Prerequisites**: plan.md, spec.md, research.md, quickstart.md

**Tests**: Go コード変更なしのため、テストタスクは不要。ワークフローの検証は実際の PR マージで行う。

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup

**Purpose**: 既存ワークフローの理解と変更準備

- [x] T001 既存の `.github/workflows/release.yml` の内容を確認し、変���箇所を特定する

---

## Phase 2: User Story 1 - main ブランチへのマージで自動リリース (Priority: P1) MVP

**Goal**: `release-X.Y.Z` または `hotfix-X.Y.Z` ブランチを main にマージすると、自動的にタグが作成され GoReleaser でリリースされる

**Independent Test**: `release-X.Y.Z` ブランチから main への PR をマージし、`vX.Y.Z` タグと GitHub Release が自動作成されることを確認

### Implementation for User Story 1

- [x] T002 [US1] `.github/workflows/release.yml` のトリガーを `on: push: tags: ['v*']` から `on: pull_request: types: [closed], branches: [main]` に変更する
- [x] T003 [US1] `.github/workflows/release.yml` の release ジョブに `if` 条件を追加: `github.event.pull_request.merged == true` かつブランチ名が `release-` または `hotfix-` で始まることを検証する
- [x] T004 [US1] `.github/workflows/release.yml` にバージョン抽出ステップを追加: `github.event.pull_request.head.ref` からプレフィックス（`release-` / `hotfix-`）を除去してバージョン文字列を取得する
- [x] T005 [US1] `.github/workflows/release.yml` に semver 検証ステップを追加: 抽出したバージョンが `^[0-9]+\.[0-9]+\.[0-9]+$` に合致しない場合エラー終了する
- [x] T006 [US1] `.github/workflows/release.yml` にタグ作成ステップを追加: `git tag v{version}` と `git push origin v{version}` を実行する
- [x] T007 [US1] 既存の GoReleaser ステップ（checkout, setup-go, goreleaser-action）がタグ作成後に正しく動作する順序になっていることを確認・調整する

**Checkpoint**: この時点で `release-X.Y.Z` からの main マージで自動リリースが動作する

---

## Phase 3: User Story 2 - 不正なブランチからのマージではリリースしない (Priority: P1)

**Goal**: `release-` / `hotfix-` 以外のブランチからのマージではリリース処理が実行されない

**Independent Test**: feature ブランチから main への PR をマージし、リリースワークフローがスキップされることを確認

### Implementation for User Story 2

- [x] T008 [US2] T003 で追加した `if` 条件により、不正なブランチからのマージでジョブ全体がスキップされることを確認する（T003 で実装済み、ここでは動作確認のみ）

**Checkpoint**: 不正ブランチからのマージでリリースが実行されないことを確認

---

## Phase 4: User Story 3 - 既存タグとの重複防止 (Priority: P2)

**Goal**: 同一バージョンのタグが既に存在する場合、ワークフローがエラーで停止する

**Independent Test**: 既存タグと同じバージョンのリリースブランチを main にマージし、エラーが報告されることを確認

### Implementation for User Story 3

- [x] T009 [US3] `.github/workflows/release.yml` に重複タグチェックステップを追加: `git ls-remote --tags origin` で `v{version}` タグの存在を確認し、存在する場合エラー終了する（T005 の後、T006 の前に配置）

**Checkpoint**: 重複タグでリリースが阻止されることを確認

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: 最終確認と整理

- [x] T010 `.github/workflows/release.yml` 全体を見直し、ステップの順序・命名・エラーメッセージが明確であることを確認する
- [x] T011 quickstart.md の Validation Steps に従い、全シナリオの動作を確認する

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: No dependencies
- **Phase 2 (US1)**: Depends on Phase 1
- **Phase 3 (US2)**: Depends on Phase 2 (T003 の `if` 条件で実現済み)
- **Phase 4 (US3)**: Depends on Phase 2 (T005 の後に配置)
- **Phase 5 (Polish)**: Depends on all phases

### Within Phase 2 (User Story 1)

T002 → T003 → T004 → T005 → T006 → T007（すべて同一ファイルへの順次変更）

### Parallel Opportunities

- 全タスクが単一ファイル（`.github/workflows/release.yml`）への変更のため、並列実行の機会はない
- ただし T009 (US3) は Phase 2 完了後に独立して追加可能

---

## Implementation Strategy

### MVP First (User Story 1)

1. Phase 1: 既存ワークフロー確認
2. Phase 2: トリガー変更 → ブランチ判定 → バージョン抽出 → semver 検証 → タグ作成 → GoReleaser 連携
3. **STOP and VALIDATE**: 実際の PR マージでリリースが自動実行されることを確認

### Incremental Delivery

1. US1 完了 → 自動リリースが動作（MVP）
2. US2 確認 → 不正ブランチでスキップされることを確認（US1 の `if` 条件で既に対応済み）
3. US3 追加 → 重複タグチェックを追加

---

## Notes

- 全タスクが `.github/workflows/release.yml` の変更に集中
- Go コード変更なし、go test 対象なし
- ワークフローの検証は実際の PR マージで行う（GitHub Actions のローカルテストは非現実的）
