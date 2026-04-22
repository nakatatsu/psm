# Implementation Plan: テスト実行基盤の整備

**Branch**: `068-test-infra` | **Date**: 2026-04-21 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/068-test-infra/spec.md`

## Summary

結合テスト資材をexample/から専用ディレクトリに分離し、GitHub Actions CIで自動実行する。AWS認証はOIDCを使用。E2Eテスト用バイナリはリリースブランチからCIで自動ビルドする。

## Technical Context

**Language/Version**: Go 1.26.1
**Primary Dependencies**: GitHub Actions, AWS SSM, SOPS, age
**Storage**: N/A
**Testing**: bash (test.sh — CLIバイナリの結合テスト), `go test` (ユニットテスト、既存)
**Target Platform**: GitHub Actions (ubuntu-latest)
**Project Type**: CLI tool (CI/CD infrastructure)
**Performance Goals**: N/A
**Constraints**: OIDC認証のみ（IAMアクセスキー禁止）。ブランチ制限でdevelop + release-*のみ許可
**Scale/Scope**: 14シナリオの結合テスト、1バイナリビルド

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| # | Principle | Gate Question | Pass? |
|---|-----------|---------------|-------|
| I | Simplicity First | Is this the simplest viable design? No unnecessary abstractions or dependencies? | Yes — 既存test.shをコピーしてCIワークフローを追加するだけ |
| II | YAGNI | Does every element serve a present, concrete need? No speculative features? | Yes — specのFRに対応する要素のみ |
| III | Test-First (NON-NEGOTIABLE) | Are all tests written before implementation (Red-Green cycle)? Using `go test` only? | Partial — test.shはbashスクリプトでありgo testではない。ただしこれはCLIバイナリの結合テストであり、Goコードのテストではない。constitution上のgo test要件はGoコードのユニットテストに適用されると解釈（survey S3参照） |

## Project Structure

### Documentation (this feature)

```text
specs/068-test-infra/
├── plan.md
├── research.md
├── survey.md
├── checklists/
│   └── requirements.md
└── spec.md
```

### Source Code (repository root)

```text
tests/
└── integration/
    ├── test.sh                  # 結合テストスクリプト（example/test.shから移設）
    └── secrets.example.yaml     # テストデータ（example/secrets.example.yamlから移設）

.github/
└── workflows/
    ├── ci.yml                   # 既存（ユニットテスト・静的解析）
    ├── integration-test.yml     # 新規（結合テスト）
    └── release.yml              # 既存（リリース）
```

**Structure Decision**: `tests/integration/` にテスト資材を配置。CIワークフローは `.github/workflows/integration-test.yml` として新規追加。既存ワークフローは変更しない。
