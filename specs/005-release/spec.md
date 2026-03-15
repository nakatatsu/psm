# Feature Specification: GitHub Releases Binary Distribution

**Feature Branch**: `005-release`
**Created**: 2026-03-15
**Status**: Draft
**Input**: タグ push 時に psm バイナリをクロスコンパイルして GitHub Releases に配布する GHA ワークフローを構築する

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Tag Push でバイナリリリース (Priority: P1)

開発者がセマンティックバージョニングのタグ（例: `v1.0.0`）を push すると、Linux/macOS/Windows 向けのバイナリが自動でビルドされ、GitHub Releases に公開される。ユーザーは Releases ページからバイナリをダウンロードしてすぐに使える。

**Why this priority**: バイナリ配布の自動化が本フィーチャーの唯一の目的。これ単体で完結する。

**Independent Test**: `v0.0.1-test` のようなテストタグを push し、GitHub Releases にバイナリが公開されることを確認する。

**Acceptance Scenarios**:

1. **Given** `v*` 形式のタグが push された, **When** リリースワークフローがトリガーされる, **Then** Linux (amd64, arm64), macOS (arm64) のバイナリがビルドされ、GitHub Releases に公開される
2. **Given** リリースが公開された, **When** ユーザーがバイナリをダウンロードして実行する, **Then** `psm sync --help` が正常に動作する
3. **Given** タグ以外の push（ブランチ push 等）, **When** リリースワークフローのトリガー条件を評価する, **Then** ワークフローは実行されない

### Edge Cases

- タグ形式が不正な場合（`v` prefix なし等）はワークフローがトリガーされない
- 同じタグを再 push した場合の挙動（既存リリースの上書き or エラー）
- ビルドが 1 つのプラットフォームで失敗した場合、他のプラットフォームのバイナリはどうなるか

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: `v*` 形式のタグ push でリリースワークフローが自動実行されなければならない
- **FR-002**: Linux (amd64, arm64), macOS (arm64) のバイナリをビルドしなければならない
- **FR-003**: ビルドされたバイナリを GitHub Releases に自動公開しなければならない
- **FR-004**: リリースノートにはタグ名とビルド対象のコミット SHA が含まれなければならない
- **FR-005**: バイナリは静的リンクで、実行環境に追加の依存関係を要求してはならない
- **FR-006**: リリースワークフローは CI ワークフロー（003-gha-ci）と独立して動作しなければならない
- **FR-007**: `v*` タグの作成は認可されたユーザーのみに制限しなければならない（tag protection rules または rulesets による）。write 権限だけでタグを作成しリリースをトリガーできてはならない

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: テストタグを push した際、5 分以内にすべてのプラットフォーム向けバイナリが GitHub Releases に公開される
- **SC-002**: 公開されたバイナリが対象プラットフォームで正常に実行できる

## Assumptions

- GoReleaser 等のリリースツールの採否は Plan で決定する
- バイナリ名は `psm-{os}-{arch}` 形式とする（具体的な命名規則は Plan で決定）
- チェックサムファイルの生成は Plan で必要性を判断する
