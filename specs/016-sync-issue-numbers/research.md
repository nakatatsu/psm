# Research: Issue 番号と SpecKit 機能番号の同期

**Date**: 2026-03-17

## Summary

Survey で主要な調査は完了済み。Phase 0 では survey で発見された残課題（S3: FR-006 の実装範囲明確化）と、実装に必要な具体的手順を確認する。

## R1: `gh issue view` の出力形式

**Decision**: `gh issue view {number} --repo {owner}/{repo} --json title,body` を使用
**Rationale**: JSON 出力で title と body を構造化取得できる。既にこのプロジェクトの speckit.specify スキルで Issue 取得に使用実績がある（016-sync-issue-numbers の specify 時に実行済み）
**Alternatives considered**:
- `gh api repos/{owner}/{repo}/issues/{number}` — 同等だが `gh issue view --json` の方がシンプル
- WebFetch で Issue ページをスクレイピング — 認証不要だが構造化データ取得が困難

## R2: Issue URL のパースパターン

**Decision**: `https://github.com/{owner}/{repo}/issues/{number}` の正規表現パターンマッチ
**Rationale**: GitHub Issue URL の形式は安定しており、このパターンで十分。スキル定義はプロンプトなので、AI がパターンマッチを行う
**Alternatives considered**:
- `gh issue view` に URL を直接渡す — `gh issue view {URL}` はサポートされているが、番号抽出が別途必要なため URL パースが先

## R3: FR-006 の実装範囲（Survey S3 の解決）

**Decision**: 「ユーザー向け手順としての `--number` 指定を SKILL.md から削除する」と「`create-new-feature.sh` の `--number` オプションはスクリプトの内部インターフェースとして維持する」の二層構造
**Rationale**: Survey S3 で指摘された通り、FR-006 の「廃止」はスキル定義レベルの変更であり、スクリプトレベルでの削除は不要。スクリプトの `--number` は他の自動化やデバッグで有用
**Alternatives considered**:
- スクリプトからも `--number` を完全削除 — 過剰。内部ツールとしての柔軟性を失う

## R4: 衝突チェックの実装方法

**Decision**: 既存の `create-new-feature.sh` の衝突検出（`git checkout -b` の失敗）をそのまま利用。スキル定義側では事前に `git branch --list` で確認し、衝突時はスクリプト実行前にエラーを出す
**Rationale**: 二重チェック（スキル側 + スクリプト側）で確実性を確保しつつ、スクリプトの変更は不要
**Alternatives considered**:
- スクリプトにのみ任せる — スクリプトのエラーメッセージが不親切な場合がある
- スキル側のみで完結 — スクリプトの安全網がなくなる
