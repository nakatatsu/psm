# Implementation Plan: git worktree セットアップ Skill

**Branch**: `026-worktree-setup-skill` | **Date**: 2026-03-23 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/026-worktree-setup-skill/spec.md`

## Summary

Claude Code Skill として `.claude/skills/worktree-setup/SKILL.md` を作成する。Skill は git worktree の作成、前提条件の自動設定、および Git 2.39 向けの相対パス変換 workaround を実行する手順を Claude Code に指示する Markdown ファイル。

## Technical Context

**Language/Version**: Markdown (Claude Code SKILL.md format)
**Primary Dependencies**: git 2.39+, bash
**Storage**: N/A
**Testing**: 手動動作確認（Go コードではないため `go test` 非対象）
**Target Platform**: DevContainer (Debian bookworm) + ホスト OS
**Project Type**: Claude Code Skill (Markdown instruction file)
**Performance Goals**: N/A
**Constraints**: Git 2.39 では `worktree.useRelativePaths` が無効、手動相対パス変換が必要
**Scale/Scope**: 単一ファイル (SKILL.md) の作成

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| # | Principle | Gate Question | Pass? |
|---|-----------|---------------|-------|
| I | Simplicity First | Is this the simplest viable design? No unnecessary abstractions or dependencies? | YES — 単一 Markdown ファイル。スクリプトや追加依存なし |
| II | YAGNI | Does every element serve a present, concrete need? No speculative features? | YES — spec の要件のみ実装。削除機能は案内のみ |
| III | Test-First (NON-NEGOTIABLE) | Are all tests written before implementation? | N/A — Go コードではなく SKILL.md (Markdown)。手動検証で対応 |

## Project Structure

### Documentation (this feature)

```text
specs/026-worktree-setup-skill/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── survey.md            # Survey output (completed)
├── checklists/
│   └── requirements.md  # Requirements checklist (completed)
└── tasks.md             # Phase 2 output (/speckit.tasks)
```

### Source Code (repository root)

```text
.claude/skills/worktree-setup/
└── SKILL.md             # The skill definition file (sole deliverable)
```

**Structure Decision**: 成果物は `.claude/skills/worktree-setup/SKILL.md` の 1 ファイルのみ。data-model、contracts は不要（データもインターフェースも存在しない）。
