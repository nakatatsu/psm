# Implementation Plan: Secrets Manager 対応オミット

**Branch**: `027-omit-secrets-manager` | **Date**: 2026-03-18 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/027-omit-secrets-manager/spec.md`

## Summary

Secrets Manager (`--store sm`) のコードとドキュメントを削除し、SSM のみをサポートする状態にする。Store interface は拡張ポイントとして保持する。コード削除が主体の chore タスクのため、テストは既存テストの通過確認が中心。

## Technical Context

**Language/Version**: Go 1.26
**Primary Dependencies**: AWS SDK for Go v2, gopkg.in/yaml.v3, regexp, log/slog
**Storage**: AWS SSM Parameter Store（SM は削除対象）
**Testing**: `go test ./...`
**Target Platform**: Linux (CLI)
**Project Type**: CLI tool
**Performance Goals**: N/A（機能削除のため）
**Constraints**: N/A
**Scale/Scope**: ファイル削除 2 件、ファイル修正 4-5 件

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| # | Principle | Gate Question | Pass? |
|---|-----------|---------------|-------|
| I | Simplicity First | Is this the simplest viable design? No unnecessary abstractions or dependencies? | YES — コード削除により複雑性が減少する |
| II | YAGNI | Does every element serve a present, concrete need? No speculative features? | YES — 使わない SM コードを削除する YAGNI 適用そのもの |
| III | Test-First (NON-NEGOTIABLE) | Are all tests written before implementation? | YES — 既存テストの通過確認。SM テスト削除と `--store sm` エラーのテスト追加は Red-Green サイクルに従う |

## Project Structure

### Documentation (this feature)

```text
specs/027-omit-secrets-manager/
├── plan.md              # This file
├── research.md          # Phase 0 output (minimal — no unknowns)
├── spec.md              # Feature specification
└── checklists/
    └── requirements.md  # Validation checklist
```

### Source Code (repository root)

```text
# 削除対象
sm.go          # SMStore 実装 → 削除
sm_test.go     # SMStore テスト → 削除

# 修正対象
main.go        # --store バリデーション変更
main_test.go   # --store sm エラーテスト追加
README.md      # SM 記述削除
README.ja.md   # SM 記述削除
example/README.md  # CLI Reference から sm 削除
CLAUDE.md      # Active Technologies から SM 削除
```

**Structure Decision**: 既存のフラットなプロジェクト構造を維持。ファイル削除と修正のみで新規ファイルは作成しない。
