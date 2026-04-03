# Feature Specification: ブランチ戦略の決定と実装

**Feature Branch**: `045-branch-strategy`  
**Created**: 2026-04-02  
**Status**: Draft  
**Input**: GitHub Issue #45 — ブランチ戦略の決定と実装

## User Scenarios & Testing *(mandatory)*

### User Story 1 - ブランチ保護ルールの設定 (Priority: P1)

開発者が保護対象ブランチに直接プッシュすることを防ぎ、PRベースの開発フローを強制する。GitFlow に基づき、以下のブランチにルールセットを適用する：

- **`develop`**: PR必須（承認不要）、CI ステータスチェック必須
- **`release-*`**: `main` と同等（PR必須・1名承認・CI必須）
- **`hotfix-*`**: `main` と同等（PR必須・1名承認・CI必須）

**Why this priority**: ブランチ保護は安全な開発フローの基盤であり、誤った直接プッシュによるコード品質低下やリリース事故を防ぐ最も重要な施策である。特に `release-*` / `hotfix-*` は本番前の品質ゲートとして `develop` より厳格にする必要がある。

**Independent Test**: Terraform で保護ルールを適用後、GitHub 上で対象ブランチへの直接プッシュが拒否されること、PR 経由でのみマージ可能であることを確認できる。

**Acceptance Scenarios**:

1. **Given** `develop` ブランチに保護ルールが適用されている状態, **When** 開発者が `develop` へ直接プッシュしようとする, **Then** プッシュが拒否される
2. **Given** `release-*` ブランチに保護ルールが適用されている状態, **When** 開発者が承認なしで PR をマージしようとする, **Then** マージが拒否される
3. **Given** `hotfix-*` ブランチに保護ルールが適用されている状態, **When** 開発者が承認なしで PR をマージしようとする, **Then** マージが拒否される
4. **Given** 保護ルールが適用されている状態, **When** リポジトリ管理者が緊急対応として直接プッシュする, **Then** バイパスルールにより許可される

---

### User Story 2 - GitHub への保護ルールの実装 (Priority: P1)

既存の Terraform コード（`.tmp/mynote/infrastructure/github`）にブランチ保護ルールを追加し、GitHub Repository Ruleset として適用する。`psm` リポジトリを対象にルールを実装する。

**Why this priority**: 設計したルールを実際に GitHub に適用しなければ保護は機能しない。設計と実装はセットで完了させる必要がある。

**Independent Test**: `terraform plan` で差分を確認し、`terraform apply` 後に GitHub リポジトリのルールセット画面で設定が反映されていることを確認できる。

**Acceptance Scenarios**:

1. **Given** Terraform コードに `develop` ブランチの保護ルールが追加されている, **When** `terraform apply` を実行する, **Then** GitHub 上に `develop` ブランチのルールセットが作成される
2. **Given** 既存の `main` ブランチ保護ルールがある, **When** Terraform コードを更新して `terraform apply` を実行する, **Then** 既存ルールに影響を与えず `develop` のルールが追加される

---

### User Story 3 - PR マージ戦略の決定 (Priority: P2)

GitFlow の各ブランチ間のマージ方法（squash / merge commit / rebase）を決定し、ルールセットまたはドキュメントに反映する。

**Why this priority**: マージ戦略はコミット履歴の可読性とリリース管理に影響するが、保護ルール自体が先に必要。

**Independent Test**: 決定したマージ戦略がドキュメントに記載され、必要に応じて GitHub の設定に反映されていることを確認できる。

**Acceptance Scenarios**:

1. **Given** マージ戦略が Merge commit のみに設定されている, **When** 開発者が PR をマージする, **Then** Merge commit でマージが行われる

---

### User Story 4 - ブランチ戦略の SKILL 化 (Priority: P2)

ブランチ命名規則やフローを Claude Code の SKILL として記載し、開発支援時に自動的に正しいブランチ運用が行われるようにする。

**Why this priority**: 開発者の日常作業での一貫性を確保するために重要だが、保護ルールの実装が先。

**Independent Test**: SKILL ファイルが存在し、Claude Code がブランチ作成時に正しい命名規則を適用することを確認できる。

**Acceptance Scenarios**:

1. **Given** ブランチ戦略 SKILL が設定されている, **When** Claude Code で新しい feature ブランチを作成する, **Then** `feature/<issue-no>-<short-description>` の命名規則に従う

---

### Edge Cases

- `develop` ブランチがまだ存在しない場合、保護ルール適用前にブランチを作成する必要がある
- 既存の `main` ブランチ保護ルール・タグ保護ルールとの競合が発生しないこと
- `release-*` / `hotfix-*` はワイルドカードパターンでマッチさせるため、命名規則に従わないブランチは保護対象外となる
- `psm` 以外のリポジトリにも将来的にルールを展開する際の拡張性

## Clarifications

### Session 2026-04-02

- Q: `develop` ブランチの PR レビュー要件 -> A: PR 必須だが承認なし（0名）
- Q: `develop` ブランチのステータスチェック要件 -> A: `main` と同じ（CI ステータスチェック必須）
- Q: PR マージ戦略 -> A: Merge commit のみ許可（GitFlow 正統派）
- Q: `release-*` / `hotfix-*` ブランチの保護ルール -> A: `main` と同等の保護（PR必須・1名承認・CI必須）。`develop` を緩めた分、本番前の品質ゲートとして締める必要がある

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: `develop` ブランチに対して、直接プッシュを禁止する保護ルールを設定しなければならない
- **FR-002**: `develop` ブランチへのマージは PR 経由のみとし、承認は不要（0名）とする
- **FR-003**: `develop` ブランチへのマージには CI ステータスチェックの通過を必須とする
- **FR-004**: `release-*` ブランチに対して、`main` と同等の保護ルール（PR必須・1名承認・CI必須）を設定しなければならない
- **FR-005**: `hotfix-*` ブランチに対して、`main` と同等の保護ルール（PR必須・1名承認・CI必須）を設定しなければならない
- **FR-006**: `develop`, `release-*`, `hotfix-*` ブランチの削除を禁止しなければならない（`develop` は恒久ブランチ、`release-*` / `hotfix-*` はリリース履歴の追跡に必要）
- **FR-007**: リポジトリ管理者は保護ルールをバイパスできなければならない
- **FR-008**: 保護ルールは Terraform（GitHub Provider）で管理しなければならない
- **FR-009**: 既存の `main` ブランチ保護ルールの `required_linear_history` を `false` に変更する（GitFlow の merge commit に統一するため）
- **FR-010**: 上記以外の既存 `main` ブランチ保護ルール・タグ保護ルールに影響を与えてはならない
- **FR-011**: PR マージ戦略は Merge commit のみ許可とする。`github_repository` リソースで `allow_merge_commit = true`, `allow_squash_merge = false`, `allow_rebase_merge = false` を設定しなければならない
- **FR-012**: `delete_branch_on_merge` を有効化しなければならない（保護ブランチは削除保護ルールにより自動削除されない）
- **FR-013**: `vulnerability_alerts` を有効化しなければならない
- **FR-014**: `security_and_analysis` で `advanced_security`, `secret_scanning`, `secret_scanning_push_protection` を有効化しなければならない
- **FR-015**: CodeQL ワークフローを GitHub Actions に追加し、Terraform の `required_code_scanning` ブロックでマージ条件に含めることを試みる。Provider のバグにより失敗した場合は CI ステータスチェック経由にフォールバックする
- **FR-016**: コミットメッセージに Issue 番号（`#<issue-no>`）を含めることを `commit-msg` hook で検証する仕組みを導入しなければならない（`.githooks/commit-msg` + `core.hooksPath`）
- **FR-017**: ブランチ戦略を Claude Code の SKILL として記載しなければならない（Issue 番号付きコミットメッセージの指示を含む）

### Assumptions

- Terraform の実行環境と権限は既に整備されている
- `psm` リポジトリは `protected_repositories` 変数に含まれている

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: `develop` ブランチへの直接プッシュが拒否される
- **SC-002**: `release-*` / `hotfix-*` ブランチへの直接プッシュが拒否され、1名承認が必要である
- **SC-003**: `terraform plan` で意図した差分のみが表示される（既存ルールへの影響なし）
- **SC-004**: ブランチ戦略の SKILL ファイルが存在し、命名規則が記載されている
- **SC-005**: PR マージ方法が Merge commit のみに制限されている
