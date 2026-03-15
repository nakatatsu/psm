# Feature Specification: GitHub Actions CI

**Feature Branch**: `003-gha-ci`
**Created**: 2026-03-15
**Status**: Draft
**Input**: PR および main ブランチへの push 時に check.md 記載の機械的チェックを自動実行する CI ワークフローを GitHub Actions で構築する

## User Scenarios & Testing _(mandatory)_

### User Story 1 - PR 品質ゲート (Priority: P1)

開発者が PR を作成または更新したとき、コードの品質チェック（ビルド、フォーマット、静的解析、リンタ、脆弱性チェック、ユニットテスト）が自動で実行され、結果が PR 上で確認できる。

**Why this priority**: マージ前に問題を検出することが CI の最大の価値。これ単体で MVP として成立する。

**Independent Test**: PR を作成し、チェック結果が GitHub の Checks タブに表示されることを確認する。

**Acceptance Scenarios**:

1. **Given** main ブランチ向けの PR が作成された, **When** CI ワークフローがトリガーされる, **Then** すべてのチェックステップ（ビルド、フォーマット、静的解析、リンタ、脆弱性、テスト）が実行され、結果が PR に表示される
2. **Given** PR のコードにフォーマット違反がある, **When** CI が実行される, **Then** フォーマットチェックステップが失敗し、PR 上で失敗が明示される
3. **Given** PR のコードにユニットテストの失敗がある, **When** CI が実行される, **Then** テストステップが失敗し、どのテストが失敗したか確認できる
4. **Given** PR のコードに既知の脆弱性を持つ依存関係がある, **When** CI が実行される, **Then** 脆弱性チェックステップが失敗する

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: CI は main ブランチ向けの pull_request イベントで自動実行されなければならない
- **FR-002**: CI は main ブランチへの push イベントで自動実行されなければならない（ブランチ保護により直接 push は禁止されているが、マージ後のチェックとして機能する）
- **FR-003**: CI は Go 1.26.1 環境でコードをビルド（`go build ./...`）しなければならない
- **FR-004**: CI は gofumpt によるフォーマットチェックを実行し、違反があれば失敗しなければならない
- **FR-005**: CI は goimports による import 順序チェックを実行し、違反があれば失敗しなければならない
- **FR-006**: CI は `go vet ./...` による静的解析を実行しなければならない
- **FR-007**: CI は staticcheck による追加の静的解析を実行しなければならない
- **FR-008**: CI は golangci-lint（gosec 含む）をプロジェクトの `.golangci.yml` 設定で実行しなければならない
- **FR-009**: CI は govulncheck による依存ライブラリの脆弱性チェックを実行しなければならない
- **FR-010**: CI は `go test -race` によるユニットテストを実行し、データ競合も検出しなければならない
- **FR-011**: CI で使用するツールバージョンは DevContainer の Dockerfile と一致させなければならない（Go 1.26.1, gofumpt 0.9.2, goimports 0.42.0, golangci-lint 2.9.0, staticcheck 2026.1, govulncheck 1.1.4）
- **FR-012**: CI はいずれかのステップが失敗した場合、全体として失敗ステータスを返さなければならない
- **FR-013**: CI は AWS 統合テスト（`PSM_INTEGRATION_TEST=1` が必要なテスト）を実行してはならない
- **FR-014**: CI ジョブ名は `ci` としなければならない（リポジトリの required status check `ci` と一致させるため）

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: main 向け PR を作成した際、CI が自動で起動し、全チェック結果が PR 上で確認できる
- **SC-002**: フォーマット違反・テスト失敗・脆弱性検出のいずれかがある場合に CI が確実に失敗する
- **SC-003**: 正常なコードの PR では CI が成功ステータスを返す

## Assumptions

- GitHub-hosted runner（`ubuntu-latest`）を使用する。Self-hosted runner は不要
- AWS 統合テストはスコープ外。`PSM_INTEGRATION_TEST` 未設定により AWS 系テスト（ssm_test, sm_test）は自動 skip され、ロジック系テスト（yaml, sync, export, main）のみ実行される
- ツールのインストールは `go install` で行う（DevContainer も Squid プロキシ導入済みで同様に可能）
- 単一ジョブ構成とする（ステップ数が少なく、ジョブ分割のオーバーヘッドが見合わない）
- main ブランチ保護（直接 push 禁止、required status check `ci`）は設定済み。本フィーチャーではワークフロー定義のみが対象
