# Implementation Plan: Issue 番号と SpecKit 機能番号の同期

**Branch**: `016-sync-issue-numbers` | **Date**: 2026-03-17 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/016-sync-issue-numbers/spec.md`

## Summary

`/speckit.specify` のスキル定義（SKILL.md）を修正し、GitHub Issue URL を必須入力とする。Issue URL から番号を自動抽出して `create-new-feature.sh` の `--number` に渡し、Issue のタイトル・説明文から short-name を生成する。Go コードの変更なし。変更対象は `.claude/skills/speckit.specify/SKILL.md` の 1 ファイル。

## Technical Context

**Language/Version**: N/A（Markdown プロンプトの変更のみ）
**Primary Dependencies**: `gh` CLI（GitHub API アクセス）、`create-new-feature.sh`（既存スクリプト）
**Storage**: N/A
**Testing**: 手動テスト（Issue URL あり/なし/不正 URL の 3 パターン）
**Target Platform**: Claude Code スキル定義
**Project Type**: CLI ツールの開発ワークフロー改善
**Performance Goals**: N/A
**Constraints**: N/A
**Scale/Scope**: 1 ファイルの変更

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| # | Principle | Gate Question | Pass? |
|---|-----------|---------------|-------|
| I | Simplicity First | Is this the simplest viable design? No unnecessary abstractions or dependencies? | YES — 既存の `--number` オプションを活用し、スキル定義のプロンプト修正のみ。新規スクリプトやツール不要 |
| II | YAGNI | Does every element serve a present, concrete need? No speculative features? | YES — Issue URL からの番号抽出と Issue 必須化のみ。将来的な拡張は含まない |
| III | Test-First (NON-NEGOTIABLE) | Are all tests written before implementation? | N/A — Go コードの変更なし。スキル定義（Markdown）の変更のため、手動テストで検証 |

## Project Structure

### Documentation (this feature)

```text
specs/016-sync-issue-numbers/
├── spec.md
├── survey.md
├── plan.md              # This file
├── research.md          # Phase 0 output
└── checklists/
    └── requirements.md
```

### Source Code (repository root)

```text
.claude/skills/speckit.specify/
└── SKILL.md             # 変更対象（唯一）
```

**Structure Decision**: 変更対象は SKILL.md の 1 ファイルのみ。`create-new-feature.sh` は変更不要（survey S4 で確認済み）。

## Design: SKILL.md の変更内容

### 変更 1: Issue URL の検出と必須化（FR-001, FR-002, FR-003）

現在のステップ 1「Generate a concise short name」の前に、Issue URL の検出・検証ステップを追加する。

**新しいフロー**:
1. ユーザー入力から GitHub Issue URL を検出する（`https://github.com/{owner}/{repo}/issues/{number}` パターン）
2. URL が見つからない場合 → エラー: 「GitHub Issue URL が必要です。先に Issue を作成してください」
3. URL が不正な形式の場合 → エラー: 「不正な Issue URL です。`https://github.com/{owner}/{repo}/issues/{number}` 形式で指定してください」
4. URL から Issue 番号を抽出する

### 変更 2: GitHub API による Issue 情報取得（FR-007, FR-008）

Issue URL 検出後、`gh issue view` で Issue の実在確認とタイトル・説明文の取得を行う。

**手順**:
1. `gh issue view {number} --repo {owner}/{repo} --json title,body` を実行
2. エラーの場合 → 「Issue #{number} が見つかりません。URL を確認してください」
3. 成功の場合 → タイトルと説明文を取得し、short-name 生成と spec 入力情報に使用

### 変更 3: 自動採番ロジックの置き換え（FR-005, FR-006）

現在のステップ 2「Check for existing branches before creating a new one」の採番ロジックを変更する。

**現在**: 既存ブランチの最大番号 + 1 で自動採番
**変更後**: Issue URL から抽出した番号を `--number` に渡す。自動採番ロジック（ステップ 2b, 2c）は削除。

既存のブランチ衝突チェック（ステップ 2a, 2b の一部）は維持する。衝突した場合はエラーとして中断（FR-004、Clarification Q1）。

### 変更 4: short-name 生成の入力変更（FR-008）

**現在**: ユーザーの feature description からキーワードを抽出
**変更後**: Issue のタイトルと説明文の両方から抽出。タイトルを主要ソースとし、説明文はタイトルだけでは不十分な場合の補助情報として使用。

既存の `--short-name` オプションによるスキル側での指定は維持（AI がタイトル・説明文から適切な名前を判断して渡す）。
