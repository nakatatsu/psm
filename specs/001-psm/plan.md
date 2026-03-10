# Implementation Plan: psm

**Branch**: `001-psm` | **Date**: 2026-03-08 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/001-psm/spec.md`

## Summary

SOPS 復号済み YAML ファイルを読み、AWS SSM Parameter Store または Secrets Manager に同期する Go CLI ツール。2 つのサブコマンド（`sync`, `export`）を持つ。`sync` は YAML → AWS 同期（差分検出、dry-run、prune 対応）。`export` は AWS → YAML エクスポート（初回セットアップ用）。Store interface で SSM/SM を抽象化する。

## Technical Context

**Language/Version**: Go 1.26 (1.26.1)
**Primary Dependencies**:
- `github.com/aws/aws-sdk-go-v2` v1.41.3
- `github.com/aws/aws-sdk-go-v2/config` v1.32.11
- `github.com/aws/aws-sdk-go-v2/service/ssm` v1.68.2
- `github.com/aws/aws-sdk-go-v2/service/secretsmanager` v1.41.3
- `gopkg.in/yaml.v3` v3.0.1
**Storage**: N/A (AWS managed services)
**Testing**: `go test` (standard library only, Constitution Principle III)
**Target Platform**: Linux / macOS CLI
**Project Type**: CLI tool
**Performance Goals**: 100 件の key-value を 1 回のコマンド実行で同期
**Constraints**: AWS API rate limits (SSM PutParameter: 40 TPS)。Bulk write API なし → goroutine 並行実行（並行数 10）で対応
**Scale/Scope**: 単一バイナリ、サブコマンド 2 つ（sync, export）

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| # | Principle | Gate Question | Pass? |
|---|-----------|---------------|-------|
| I | Simplicity First | Is this the simplest viable design? No unnecessary abstractions or dependencies? | Yes — CLI フラグは `--store`, `--profile`, `--prune`, `--dry-run` のみ。ストア種別は CLI フラグ `--store` で指定。1 interface（Store）、フラットパッケージ構成。依存は AWS SDK と yaml.v3 のみ（いずれも必須）。 |
| II | YAGNI | Does every element serve a present, concrete need? No speculative features? | Yes — YAML は純粋な key-value のみ（メタデータなし）。設定は CLI フラグで完結。YAML キー = AWS リソース名のダイレクトマッピング。 |
| III | Test-First (NON-NEGOTIABLE) | Are all tests written before implementation (Red-Green cycle)? Using `go test` only? No third-party test frameworks? | Yes — Mock 不使用。純粋ロジックはユニットテスト、AWS API はサンドボックス統合テスト。testify/mockgen 不使用。Red-Green cycle 遵守。 |

## Project Structure

### Documentation (this feature)

```text
specs/001-psm/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── cli-schema.md
│   └── store-interface.md
└── tasks.md
```

### Source Code (repository root)

```text
main.go                 # エントリポイント: サブコマンドディスパッチ
store.go                # Store interface 定義
ssm.go                  # SSM Store 実装
sm.go                   # Secrets Manager Store 実装
yaml.go                 # YAML ファイル読み込み・パース・書き出し
sync.go                 # 差分計算 + 同期ロジック
export.go               # エクスポートロジック
sync_test.go            # sync ロジックのユニットテスト + Sandbox AWS 統合テスト
export_test.go          # export ロジックの Sandbox AWS 統合テスト
yaml_test.go            # YAML パース・書き出しのテスト
main_test.go            # CLI サブコマンドパースのテスト
go.mod
go.sum
```

**Structure Decision**: Constitution Principle I（Simplicity First）に従い、フラットなパッケージ構成。全ファイルを `package main` に配置。ファイル数が 10 未満のため、ディレクトリ分割は不要。

## Complexity Tracking

> Constitution Check に違反なし。記載不要。
