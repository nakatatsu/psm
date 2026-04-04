# Implementation Plan: デプロイ方式の変更

**Branch**: `060-change-deploy-method` | **Date**: 2026-04-04 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/060-change-deploy-method/spec.md`

## Summary

main ブランチへの PR マージをトリガーに、リリースブランチ名（`release-X.Y.Z` / `hotfix-X.Y.Z`）からバージョンを自動抽出し、Git タグ作成と GoReleaser によるリリースを一連のワークフローとして実行する。既存の手動タグ作成トリガー（`on: push: tags`）を置き換える。

## Technical Context

**Language/Version**: GitHub Actions workflow (YAML)
**Primary Dependencies**: goreleaser/goreleaser-action@v7, actions/checkout@v6, actions/setup-go@v6
**Storage**: N/A
**Testing**: GitHub Actions ワークフローのため go test 対象外。手動テスト（実際の PR マージ）で検証。
**Target Platform**: GitHub Actions (ubuntu-latest)
**Project Type**: CI/CD workflow
**Performance Goals**: N/A
**Constraints**: GITHUB_TOKEN でのタグ作成は他ワークフローをトリガーしない（survey S2）
**Scale/Scope**: 単一ファイル（`.github/workflows/release.yml`）の変更

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| # | Principle | Gate Question | Pass? |
|---|-----------|---------------|-------|
| I | Simplicity First | Is this the simplest viable design? No unnecessary abstractions or dependencies? | Yes — 単一ワークフローファイルの変更のみ。新規依存なし。 |
| II | YAGNI | Does every element serve a present, concrete need? No speculative features? | Yes — pre-release 対応等は含めない。必要最小限のブランチ名判定とタグ作成。 |
| III | Test-First (NON-NEGOTIABLE) | Are all tests written before implementation? | N/A — Go コード変更なし。GitHub Actions ワークフローは go test の対象外。 |

## Project Structure

### Documentation (this feature)

```text
specs/060-change-deploy-method/
├── spec.md
├── survey.md
├── plan.md              # This file
├── research.md          # Phase 0 output
└── quickstart.md        # Phase 1 output
```

### Source Code (repository root)

```text
.github/workflows/
└── release.yml          # 変更対象（唯一のファイル）
```

**Structure Decision**: 本変更は `.github/workflows/release.yml` のみの変更。Go ソースコードの変更はなし。data-model.md と contracts/ はデータモデルや外部インターフェースの変更がないため不要。
