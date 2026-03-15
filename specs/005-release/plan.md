# Implementation Plan: GitHub Releases Binary Distribution

**Branch**: `005-release` | **Date**: 2026-03-15 | **Spec**: specs/005-release/spec.md
**Input**: Feature specification from `/specs/005-release/spec.md`

## Summary

Tag push (`v*`) で GoReleaser を使ってクロスコンパイルし、GitHub Releases にバイナリを公開する。3 プラットフォーム (linux/amd64, linux/arm64, darwin/arm64)。

## Technical Context

**Language/Version**: Go 1.26.1
**Primary Dependencies**: goreleaser/goreleaser-action@v7, GoReleaser v2.x (OSS)
**Storage**: N/A
**Testing**: Tag push → GitHub Releases に公開されることを確認
**Target Platform**: GitHub-hosted runner (ubuntu-latest)
**Project Type**: CI/CD configuration (YAML workflow + GoReleaser config)
**Constraints**: No CGO. Static binaries only.

## Constitution Check

| # | Principle | Gate Question | Pass? |
|---|-----------|---------------|-------|
| I | Simplicity First | Is this the simplest viable design? | Yes — GoReleaser is one tool, ~35 lines total config |
| II | YAGNI | Does every element serve a present, concrete need? | Yes — 3 targets (not 5), no Windows |
| III | Test-First | N/A — YAML config, not Go code | N/A |

## Project Structure

### Source Code (repository root)

```text
.github/
└── workflows/
    └── release.yml      # Release workflow (new)

.goreleaser.yaml         # GoReleaser config (new)
```

**Structure Decision**: Two config files at repository root. No Go code changes.
