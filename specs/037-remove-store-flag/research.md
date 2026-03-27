# Research: --store フラグの除去

**Date**: 2026-03-27

## Summary

本機能に未解決の技術的不明点はない。すべての変更は既存コードの削除・簡素化であり、新技術の導入や外部調査は不要。

## R1: --store フラグの除去方法

**Decision**: `flag.NewFlagSet` から `store` フラグ定義を削除し、`--store` が引数に含まれる場合は廃止エラーを返す。
**Rationale**: Go の `flag` パッケージでは未定義フラグが渡されると `fs.Parse` がエラーを返す。ただしエラーメッセージが汎用的なため、`--store` 専用の分かりやすいエラーメッセージを事前チェックで返す方が UX が良い。`--prune` の廃止と同じパターン（`main.go:209-213`）を踏襲する。
**Alternatives considered**:
- `fs.Parse` のデフォルトエラーに任せる → メッセージが不親切（「flag provided but not defined: -store」）
- フラグを残してデフォルト値を設定 → YAGNI 違反、不要な複雑性

## R2: Config.Store フィールドの扱い

**Decision**: Config 構造体から `Store string` フィールドを削除する。
**Rationale**: フラグを除去した以上、Config に Store フィールドを残す理由がない。`run()` 関数で直接 `NewSSMStore` を呼ぶ。
**Alternatives considered**:
- Store フィールドを残してデフォルト値 "ssm" を設定 → 使われないフィールドは YAGNI 違反

## R3: 後方互換性エラーメッセージ

**Decision**: `--prune` 廃止と同じパターンで、`args` を事前スキャンして `--store` を検出し、廃止メッセージを返す。
**Rationale**: `fs.Parse` の前にチェックすることで、他のフラグとの干渉を避け、明確なエラーメッセージを出せる。
**Alternatives considered**:
- `fs.Parse` 後にチェック → 未定義フラグとして先にエラーになるため不可
