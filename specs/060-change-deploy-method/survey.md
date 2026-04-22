# Survey: デプロイ方式の変更

**Date**: 2026-04-04
**Spec**: [spec.md](./spec.md)

## Summary

spec の方向性は妥当。手動タグ作成を排除し、main マージをトリガーにする方針はオペミス防止の本質的な解決策である。ただし、トリガーイベントの選択（push vs pull_request）が実装の信頼性に大きく影響するため、ここを精査した。また、GITHUB_TOKEN によるタグ作成が他ワークフローをトリガーしない制約を確認し、単一ワークフローへの統合が必要であることを確認した。

## S1: トリガーイベントの選択 — push vs pull_request

**Category**: Problem Reframing / Approach Alternatives

**Finding**: spec の FR-001 は「main への push イベント」をトリガーとしている。しかし push イベントではマージ元ブランチ名を直接取得できない。マージコミットメッセージのパース（例: `Merge pull request #N from owner/release-1.0.0`）が必要になり、フォーマット依存で脆弱になる。

一方、`pull_request` イベント（`types: [closed]` + `merged == true`）では `github.head_ref` でマージ元ブランチ名を直接取得でき、パースが不要。

**Recommendation**: `pull_request` イベント（closed + merged）を使用する。`github.head_ref` からブランチ名を取得する方が堅牢。

**Evidence**:
- 既存の main へのマージは全て PR 経由（git log 確認: `Merge pull request #57 from nakatatsu/release-0.0.4`）
- GitHub Branch Protection / Ruleset で main への直接 push は禁止されている（specs/045-branch-strategy で設定済み）
- GitHub Actions 公式ドキュメント: `pull_request` イベントでは `github.head_ref` が利用可能

## S2: GITHUB_TOKEN によるタグ push の制約

**Category**: Risk & Failure Modes

**Finding**: GitHub Actions の GITHUB_TOKEN で作成した tag push は、他の `on: push: tags` ワークフローをトリガーしない（無限ループ防止のための GitHub の仕様）。つまり、新ワークフローでタグを作成しても、現在の release.yml（`on: push: tags: ['v*']`）は起動しない。

**Recommendation**: release.yml を「タグトリガー」から「PR マージトリガー + タグ作成 + GoReleaser」に一本化する。旧トリガーは削除する。

**Evidence**: GitHub Docs — "When you use the repository's GITHUB_TOKEN to perform tasks, events triggered by the GITHUB_TOKEN will not create a new workflow run."

## S3: 既存ワークフローとの共存

**Category**: Integration Impact

**Finding**: 変更対象は `.github/workflows/release.yml` のみ。ci.yml の `check-source-branch` ジョブは main への PR 時にソースブランチを検証しており、本変更とは独立して機能する。GoReleaser 設定（`.goreleaser.yaml`）もタグが存在すれば動作するため変更不要。

**Recommendation**: release.yml のみ変更。ci.yml と .goreleaser.yaml は変更不要。

**Evidence**: ci.yml は `pull_request` と `push: branches: [main, develop]` トリガーで、tags には関与しない。.goreleaser.yaml は `{{.Version}}` を ldflags で使用するが、これは GoReleaser が git tag から自動取得する。

## S4: ブランチ名からのバージョン抽出の堅牢性

**Category**: Edge Cases / Feasibility

**Finding**: `release-1.0.0` → `1.0.0` の抽出はプレフィックス除去で単純に実現できる。semver 検証は正規表現 `^[0-9]+\.[0-9]+\.[0-9]+$` で十分（pre-release やメタデータは現在使用していない）。既存タグ（v0.0.1〜v0.0.4）はすべてこの形式。

**Recommendation**: シンプルな正規表現で十分。pre-release 対応は YAGNI。

**Evidence**: 既存タグ: v0.0.1, v0.0.2, v0.0.3, v0.0.4（すべて `X.Y.Z` 形式）

## S5: Constitution 適合性

**Category**: Constitution Compliance

**Finding**: 本変更は GitHub Actions ワークフロー（YAML）の変更のみで、Go コードの変更を伴わない。Constitution の Test-First 原則（Section III）は Go コードに適用されるもので、CI/CD ワークフローの変更には直接適用されない。Simplicity First 原則には合致（単純なブランチ名判定とタグ作成）。

**Recommendation**: Constitution の修正は不要。

## Items Requiring PoC

なし。すべての技術的要素（pull_request イベント、github.head_ref、GITHUB_TOKEN の制約）は公式ドキュメントで確認済み。

## Constitution Impact

修正不要。

## Recommendation

Proceed to `/speckit.plan`. ただし spec の FR-001 を「push イベント」から「pull_request closed + merged イベント」に更新することを推奨（S1 の知見を反映）。
