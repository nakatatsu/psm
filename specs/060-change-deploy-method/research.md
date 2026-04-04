# Research: デプロイ方式の変更

**Date**: 2026-04-04
**Phase**: 0 (Outline & Research)

## Overview

survey.md でトリガーイベントの選択（S1）、GITHUB_TOKEN の制約（S2）、既存ワークフローとの共存（S3）を調査済み。本 research.md では survey で確定した方針に基づく具体的な実装決定を記録する。

## R1: マージ元ブランチ名の取得方法

**Decision**: `pull_request` イベント（`types: [closed]`）+ `github.event.pull_request.merged == true` で検知し、`github.event.pull_request.head.ref` でブランチ名を取得する。

**Rationale**: push イベントではマージ元ブランチ名を直接取得できず、マージコミットメッセージのパースが必要になる。`pull_request` イベントなら `head.ref` で確実に取得可能。main への直接 push は Ruleset で禁止されているため、全マージは PR 経由であることが保証される。

**Alternatives considered**:
- `on: push` + マージコミットメッセージパース → フォーマット依存で脆弱
- `on: push` + GitHub API でマージ PR を逆引き → 不必要に複雑

## R2: タグ作成の方法

**Decision**: ワークフロー内で `git tag` + `git push --tags` を実行する。GITHUB_TOKEN の `contents: write` 権限を使用。

**Rationale**: GitHub Actions の GITHUB_TOKEN は `contents: write` でタグの作成・push が可能。別途 PAT や GitHub App token は不要。GoReleaser は同一ジョブ内でタグが存在すれば動作する。

**Alternatives considered**:
- GitHub API (`POST /repos/{owner}/{repo}/git/refs`) → 可能だが git コマンドの方がシンプル
- 別ワークフローでタグ作成 → GITHUB_TOKEN のタグ push は他ワークフローをトリガーしないため不可

## R3: GoReleaser へのタグの渡し方

**Decision**: `git tag` でタグを作成した後、`goreleaser-action` の `args: release --clean` をそのまま使用。GoReleaser は直近のタグを自動検出する。

**Rationale**: GoReleaser は `git describe --tags` でタグを取得する。同一ジョブ内でタグを作成すれば、追加設定なしで動作する。既存の `.goreleaser.yaml` の `ldflags: -X main.version={{.Version}}` もタグから自動解決される。

**Alternatives considered**:
- `GORELEASER_CURRENT_TAG` 環境変数で明示指定 → 動作するが不要な冗長さ

## R4: セマンティックバージョン検証の範囲

**Decision**: `X.Y.Z` 形式（数字3つ + ドット2つ）のみ許可。pre-release（`-alpha.1`）やビルドメタデータ（`+build.1`）は対象外。

**Rationale**: 既存タグ（v0.0.1〜v0.0.4）はすべて `X.Y.Z` 形式。pre-release の使用実績がなく、YAGNI 原則に従い対応しない。正規表現: `^[0-9]+\.[0-9]+\.[0-9]+$`

**Alternatives considered**:
- 完全な semver 正規表現（pre-release + metadata 対応）→ YAGNI
