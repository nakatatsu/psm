# Implementation Plan: Version Display Command

**Branch**: `030-version-display` | **Date**: 2026-03-25 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/030-version-display/spec.md`

## Summary

psm に `--version` フラグを追加し、ビルド時に注入されたバージョン文字列を表示する。GoReleaser がリリースビルド時にバージョンを自動注入し、手動ビルド時は `dev` をフォールバック表示する。

## Technical Context

**Language/Version**: Go 1.26.1
**Primary Dependencies**: 標準ライブラリのみ（`fmt`, `os`）
**Storage**: N/A
**Testing**: `go test`（標準ライブラリのみ）
**Target Platform**: Linux (amd64, arm64), macOS (arm64)
**Project Type**: CLI tool
**Performance Goals**: N/A（即時表示）
**Constraints**: N/A
**Scale/Scope**: 単一変数の追加と `parseArgs` の分岐追加

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| # | Principle | Gate Question | Pass? |
|---|-----------|---------------|-------|
| I | Simplicity First | Is this the simplest viable design? No unnecessary abstractions or dependencies? | Yes — パッケージ変数1つと条件分岐1つのみ |
| II | YAGNI | Does every element serve a present, concrete need? No speculative features? | Yes — バージョン表示のみ、commit/date は含めない |
| III | Test-First (NON-NEGOTIABLE) | Are all tests written before implementation (Red-Green cycle)? Using `go test` only? No third-party test frameworks? | Yes — table-driven test で `parseArgs` のバージョン分岐をテスト |

## Project Structure

### Documentation (this feature)

```text
specs/030-version-display/
├── spec.md
├── survey.md
├── plan.md              # This file
├── research.md
└── tasks.md             # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

```text
.goreleaser.yaml         # ldflags にバージョン注入を追加
main.go                  # version 変数と --version 分岐を追加
main_test.go             # --version のテストケースを追加
```

**Structure Decision**: 既存のフラットパッケージ構造を維持。新規ファイルは不要で、既存ファイルの修正のみ。
