# Implementation Plan: ブランチ戦略の決定と実装

**Branch**: `045-branch-strategy` | **Date**: 2026-04-03 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/045-branch-strategy/spec.md`

## Summary

GitFlow ブランチ戦略に基づく保護ルールを Terraform で実装し、リポジトリ設定（マージ戦略・セキュリティ）を適用する。加えて CodeQL ワークフロー、commit-msg hook、Claude Code SKILL を整備する。

## Technical Context

**Language/Version**: HCL (Terraform), Bash, YAML (GitHub Actions)
**Primary Dependencies**: Terraform GitHub Provider ~> 6.11, GitHub Actions
**Storage**: N/A（GitHub API 経由でリポジトリ設定を管理）
**Testing**: `terraform plan`（差分確認）、GitHub UI での動作確認
**Target Platform**: GitHub（パブリックリポジトリ、Free プラン）
**Project Type**: Infrastructure as Code（リポジトリガバナンス）
**Performance Goals**: N/A
**Constraints**: GitHub Free プランの制約（Enterprise 限定機能は使用不可）
**Scale/Scope**: `psm` リポジトリ1つ（`protected_repositories` 変数で管理）

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

| #   | Principle                   | Gate Question                                                                    | Pass?                                                                                          |
| --- | --------------------------- | -------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------- |
| I   | Simplicity First            | Is this the simplest viable design? No unnecessary abstractions or dependencies? | Yes — 既存の Terraform コードに最小限のリソースを追加するのみ。新たな依存関係なし              |
| II  | YAGNI                       | Does every element serve a present, concrete need? No speculative features?      | Yes — 全 FR が spec/survey の決定に基づく。Enterprise 限定機能等は明示的に除外済み             |
| III | Test-First (NON-NEGOTIABLE) | Are all tests written before implementation?                                     | N/A — Go コードの変更なし。Terraform は `terraform plan` で検証、GitHub Actions は実環境テスト |

## Project Structure

### Documentation (this feature)

```text
specs/045-branch-strategy/
├── spec.md
├── survey.md
├── plan.md              # This file
├── research.md          # Phase 0 output
├── checklists/
│   └── requirements.md
└── tasks.md             # Phase 2 output (/speckit.tasks)
```

### Source Code (変更対象ファイル)

```text
# Terraform（.tmp/mynote/infrastructure/github/ — 別リポジトリ管理）
.tmp/mynote/infrastructure/github/
├── main.tf              # Ruleset リソース追加・既存 main ルール修正
├── repository.tf        # 新規: github_repository リソース（マージ戦略・セキュリティ設定）
├── variables.tf         # 必要に応じて変数追加
├── outputs.tf
├── providers.tf
└── versions.tf

# GitHub Actions（このリポジトリ）
.github/workflows/
├── ci.yml               # 既存 CI（変更なし）
├── release.yml          # 既存リリース（変更なし）
└── codeql.yml           # 新規: CodeQL 分析ワークフロー

# Git hooks（このリポジトリ）
.githooks/
└── commit-msg           # 新規: Issue 番号検証スクリプト

# Claude Code SKILL（このリポジトリ）
.claude/skills/
└── branch-strategy/     # 新規: ブランチ戦略 SKILL
```

**Structure Decision**: Terraform コードは `.tmp/mynote/infrastructure/github/` に配置（別リポジトリ管理）。GitHub Actions、git hooks、SKILL は `psm` リポジトリに配置。

## Implementation Phases

### Phase A: Terraform — Ruleset 追加（FR-001〜FR-010）

1. **既存 `main` ルールの修正**: `required_linear_history = false` に変更
2. **`develop` ブランチ Ruleset 追加**:
   - `pull_request`: `required_approving_review_count = 0`
   - `required_status_checks`: `ci` チェック必須（strict）
   - `non_fast_forward = true`
   - `deletion = true`
   - bypass_actors: Repository admin
3. **`release-*` ブランチ Ruleset 追加**:
   - `main` と同等の PR ルール（1名承認、stale review 却下、last push approval、コメント解決必須）
   - `required_status_checks`: `ci` チェック必須（strict）
   - `non_fast_forward = true`
   - `deletion = true`
   - bypass_actors: Repository admin
4. **`hotfix-*` ブランチ Ruleset 追加**: `release-*` と同一設定

### Phase B: Terraform — リポジトリ設定（FR-011〜FR-014）

1. **`github_repository` リソース追加**（`repository.tf`）:
   - `allow_merge_commit = true`
   - `allow_squash_merge = false`
   - `allow_rebase_merge = false`
   - `delete_branch_on_merge = true`
   - `vulnerability_alerts = true`
   - `security_and_analysis` ブロック（`advanced_security`, `secret_scanning`, `secret_scanning_push_protection`）

### Phase C: CodeQL ワークフロー（FR-015）

1. `.github/workflows/codeql.yml` を作成（Go 用 CodeQL 分析）
2. Terraform の `required_code_scanning` ブロックでマージ条件に含めることを試みる
3. Provider バグで失敗した場合は `required_status_checks` に CodeQL チェックを追加してフォールバック

### Phase D: commit-msg hook（FR-016）

1. `.githooks/commit-msg` スクリプトを作成（`#<issue-no>` パターン検証）
2. README または CLAUDE.md に `git config core.hooksPath .githooks` の設定手順を記載

### Phase E: SKILL 作成（FR-017）

1. `.claude/skills/branch-strategy/` に SKILL ファイルを作成
2. GitFlow フロー、ブランチ命名規則、Issue 番号付きコミットメッセージの指示を含む

## Complexity Tracking

憲法違反なし。追記不要。
