# Implementation Plan: --store フラグの除去

**Branch**: `037-remove-store-flag` | **Date**: 2026-03-27 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/037-remove-store-flag/spec.md`

## Summary

`--store ssm` フラグを CLI から除去し、SSM をデフォルトストアとしてハードコードする。Config 構造体から Store フィールドを削除し、`main.go` の switch 文を直接呼び出しに変更する。テスト・ドキュメントから `--store` 関連の記述をすべて除去する。

## Technical Context

**Language/Version**: Go (latest stable)
**Primary Dependencies**: aws-sdk-go-v2, standard library
**Storage**: AWS SSM Parameter Store (via SSMStore)
**Testing**: `go test` (standard library のみ)
**Target Platform**: Linux/macOS CLI
**Project Type**: CLI tool
**Constraints**: 変更はフラグ除去のみ。Store interface は保持。

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| # | Principle | Gate Question | Pass? |
|---|-----------|---------------|-------|
| I | Simplicity First | Is this the simplest viable design? No unnecessary abstractions or dependencies? | YES — 不要なフラグを除去してシンプル化する変更そのもの |
| II | YAGNI | Does every element serve a present, concrete need? No speculative features? | YES — 選択肢が1つしかないフラグを除去。YAGNI の実践 |
| III | Test-First (NON-NEGOTIABLE) | Are all tests written before implementation (Red-Green cycle)? Using `go test` only? | YES — テストを先に更新し、Red-Green サイクルで実装 |

## Project Structure

### Documentation (this feature)

```text
specs/037-remove-store-flag/
├── spec.md
├── survey.md
├── plan.md              # This file
├── research.md
├── data-model.md
└── contracts/
    └── cli-schema.md
```

### Source Code (repository root)

```text
main.go          # parseArgs: --store 除去, run: switch 除去
store.go         # Config: Store フィールド除去
main_test.go     # テストケース更新
README.md        # コマンド例更新
example/README.md  # CLI Reference 更新
example/test.sh  # 統合テスト更新
```

**Structure Decision**: 既存のフラットなパッケージ構造をそのまま維持。新規ファイルは不要。

## Complexity Tracking

該当なし。Constitution 違反なし。
