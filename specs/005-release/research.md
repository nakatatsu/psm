# Research: GitHub Releases Binary Distribution

**Date**: 2026-03-15

## R1: Release Tool

**Decision**: Use goreleaser-action@v7 with GoReleaser v2.x (OSS edition).
**Rationale**: GoReleaser handles cross-compile, checksums, archives, release notes, and GitHub Releases upload in one tool. OSS edition fully supports GitHub Releases. ~35 lines of total config.
**Alternatives considered**:
- Hand-rolled matrix build + `softprops/action-gh-release` — simpler upfront but accumulates complexity (no checksums, no archives, manual release notes)
- `go install` only — zero effort but no prebuilt binaries for non-Go users

## R2: Target Platforms

**Decision**: linux/amd64, linux/arm64, darwin/arm64 (3 targets).
**Rationale**: Survey identified Windows as unlikely for an AWS CLI tool. macOS amd64 (Intel) is declining. These 3 cover 90%+ of real usage.
**Alternatives considered**:
- 5 targets (spec original) — over-scoped, Windows unlikely
- 2 targets (linux/amd64 + darwin/arm64) — viable but arm64 Linux (Graviton) is worth including

## R3: Build Flags

**Decision**: `CGO_ENABLED=0`, `-ldflags="-s -w"`, version embedding via GoReleaser `ldflags` template.
**Rationale**: CGO_ENABLED=0 guarantees static binary. `-s -w` strips debug info (~30% smaller). GoReleaser auto-embeds version from git tag.

## R4: GoReleaser Config

**Decision**: Minimal `.goreleaser.yaml`:

```yaml
version: 2

builds:
  - env:
      - CGO_ENABLED=0
    goos:
      - linux
      - darwin
    goarch:
      - amd64
      - arm64
    ignore:
      - goos: darwin
        goarch: amd64
    ldflags:
      - -s -w

archives:
  - format: tar.gz

checksum:
  name_template: 'checksums.txt'
```

## R5: Workflow File

**Decision**: `.github/workflows/release.yml`, triggers on `v*` tag push only.

```yaml
on:
  push:
    tags:
      - 'v*'
```

Uses `goreleaser/goreleaser-action@v7` with `version: '~> v2'`.
