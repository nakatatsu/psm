# Survey: GitHub Releases Binary Distribution

**Date**: 2026-03-15
**Spec**: specs/005-release/spec.md

## Summary

The spec's direction is valid but potentially over-scoped. The primary question is whether prebuilt binaries are needed at all — `go install` covers Go-capable users with zero effort. If binaries are needed, goreleaser-action is the standard tool (less total effort than hand-rolling). Windows support is likely unnecessary for an AWS tool. Starting with 2-3 targets is sufficient.

## S1: Problem Definition — Are prebuilt binaries needed?

**Category**: Problem Definition
**Finding**: `go install github.com/nakatatsu/psm@latest` covers all users with a Go toolchain. Prebuilt binaries are needed only for non-Go users or CI pipelines without Go. For an AWS parameter sync tool, the audience likely has Go or can install it.
**Recommendation**: Add `go install` to README regardless. Proceed with binary release as a convenience, but don't over-invest.

## S2: Platform targets — 5 is too many

**Category**: Hidden Assumptions
**Finding**: Windows support for an AWS CLI tool is unlikely to be used. Linux amd64 + macOS arm64 covers 90%+ of real usage. arm64 Linux (Graviton) is nice-to-have.
**Recommendation**: Start with linux/amd64, darwin/arm64, linux/arm64 (3 targets). Add Windows only if requested.

## S3: GoReleaser vs hand-rolled

**Category**: Approach Alternatives
**Finding**: goreleaser-action handles cross-compile, checksums, archives, release notes, and GitHub Releases upload in ~35 lines of config. A hand-rolled matrix build looks simpler but accumulates complexity. GoReleaser is the standard for Go projects.
**Recommendation**: Use goreleaser-action. It's the right level of tool for this job.

## S4: Build flags

**Category**: Feasibility
**Finding**: Always set `CGO_ENABLED=0` (guarantee static binary) and `-ldflags="-s -w"` (strip debug, ~30% smaller). GoReleaser handles version embedding (`-X main.version`) automatically.
**Recommendation**: Include in GoReleaser config.

## Items Requiring PoC

None.

## Constitution Impact

No amendments needed.

## Recommendation

Proceed to plan. Reduce targets to 3 (drop Windows, drop macOS amd64). Use goreleaser-action.
