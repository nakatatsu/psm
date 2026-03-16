# Implementation Plan: Distinguish dry-run output from actual execution

**Branch**: `008-dryrun-output` | **Date**: 2026-03-16 | **Spec**: specs/008-dryrun-output/spec.md

## Summary

`psm sync --dry-run` の出力に `(dry-run)` プレフィックスを追加し、通常実行の出力と明確に区別できるようにする。変更対象は `sync.go` の `execute` 関数内の `fmt.Fprintf` 呼び出しのみ。

## Technical Context

**Language/Version**: Go 1.26
**Primary Dependencies**: 標準ライブラリのみ（fmt, io）
**Storage**: N/A
**Testing**: `go test`（Test-First）
**Target Platform**: CLI（Linux, macOS）
**Project Type**: CLI tool
**Performance Goals**: N/A
**Constraints**: プレーンテキスト出力（ANSI カラー不使用）
**Scale/Scope**: `sync.go` 1ファイル + テスト

## Constitution Check

| # | Principle | Gate Question | Pass? |
|---|-----------|---------------|-------|
| I | Simplicity First | Is this the simplest viable design? | Yes — `fmt.Fprintf` のフォーマット文字列変更のみ |
| II | YAGNI | Does every element serve a present, concrete need? | Yes — Issue #10 で報告された具体的な問題の修正 |
| III | Test-First (NON-NEGOTIABLE) | Tests written before implementation? | Yes — dry-run 表記の有無をテストしてから実装 |

## Project Structure

### Documentation (this feature)

```text
specs/008-dryrun-output/
├── spec.md
├── survey.md
├── plan.md              # This file
└── tasks.md
```

### Source Code (repository root)

```text
sync.go          # execute 関数の fmt.Fprintf 修正
sync_test.go     # dry-run 表記テスト追加
```

**Structure Decision**: 既存ファイルの修正のみ。新規ファイル不要。

## Design

### 出力フォーマット

**dry-run 時:**
```
(dry-run) create: /myapp/database/host
(dry-run) update: /myapp/api/key
(dry-run) delete: /myapp/old/key
4 created, 0 updated, 1 deleted, 0 unchanged, 0 failed (dry-run)
```

**通常実行時（変更なし）:**
```
create: /myapp/database/host
update: /myapp/api/key
delete: /myapp/old/key
4 created, 0 updated, 1 deleted, 0 unchanged, 0 failed
```

### 変更箇所

`sync.go` の `execute` 関数内:
- アクション行: `dryRun` が true の場合、`"(dry-run) %s: %s\n"` フォーマットを使用
- サマリー行: `dryRun` が true の場合、末尾に ` (dry-run)` を追加
