# Research: ブランチ戦略の決定と実装

**Date**: 2026-04-03
**Phase**: 0 (Outline & Research)

## Overview

survey.md で GitHub Repository Ruleset の全オプション調査、`github_repository` リソースの設定調査が完了済み。本 research.md では survey で未解決だった具体的な実装上の決定事項を記録する。

## R1: Terraform `github_repository` リソースの import 戦略

**Decision**: `terraform import` で既存リポジトリを取り込んでから設定を適用する

**Rationale**: `psm` リポジトリは既に存在するため、`github_repository` リソースを新規作成すると Terraform が新しいリポジトリを作ろうとしてエラーになる。`terraform import github_repository.psm psm` で既存リソースを state に取り込む必要がある。

**Alternatives considered**:
- `data` ソースのみ使用 → 設定変更ができないため不採用
- 手動で GitHub UI から設定 → IaC の原則に反するため不採用

## R2: CodeQL ワークフローの設定

**Decision**: GitHub 公式の CodeQL Action（`github/codeql-action`）を使用し、Go 言語を対象にする

**Rationale**: Go はビルド型言語のため `autobuild` が使えるが、シンプルな Go プロジェクトなので自動検出で十分。スケジュール実行（weekly）と push/PR トリガーの両方を設定する。

**Alternatives considered**:
- CodeQL CLI を直接使用 → Action の方がメンテナンスが容易
- サードパーティのセキュリティスキャナー → Constitution の「標準ツール優先」に反する

## R3: commit-msg hook のパターン設計

**Decision**: コミットメッセージに `#<数字>` が含まれることを検証する。マージコミット（Merge branch 等）は除外する。

**Rationale**: Issue 番号はコミットメッセージの任意の位置に `#123` 形式で含まれていればよい。タイトル行に強制するとフォーマットが窮屈になる。マージコミットは Git が自動生成するため除外が必要。

**Alternatives considered**:
- Conventional Commits 形式の強制（`feat(#123): ...`）→ 過剰な制約
- タイトル行のみチェック → 本文に書くケースを排除してしまう

## R4: `release-*` / `hotfix-*` の Ruleset パターンマッチ

**Decision**: `refs/heads/release-*` と `refs/heads/hotfix-*` をそれぞれ別の Ruleset として作成する

**Rationale**: `release-*` と `hotfix-*` は同一の保護ルールだが、別リソースとして管理することで将来的にルールを分離可能にする。Terraform の `for_each` で `protected_repositories` を回す既存パターンに合わせる。

**Alternatives considered**:
- 1つの Ruleset で `include = ["refs/heads/release-*", "refs/heads/hotfix-*"]` → 可能だが、将来の分離が困難
- Organization Ruleset → GitHub Free では利用不可

## R5: 既存 Terraform コードとの整合性

**Decision**: 既存の `main.tf` のリソース構造（`for_each = var.protected_repositories`）を踏襲する

**Rationale**: 既存コードとの一貫性を維持し、`psm` 以外のリポジトリへの展開も容易にする。新しい Ruleset も同じパターンで `for_each` を使う。

**Alternatives considered**:
- `psm` 専用のハードコードされたリソース → 拡張性がなく既存パターンと不整合
