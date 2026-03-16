# Survey: Distinguish dry-run output from actual execution

**Date**: 2026-03-16
**Spec**: specs/008-dryrun-output/spec.md

## Summary

正しい問題を解いている。変更範囲は `sync.go` の `execute` 関数内の `fmt.Fprintf` 呼び出し数箇所のみで、リスクは極めて低い。既存テスト（`sync_test.go`, `ssm_test.go`）が出力文字列をアサートしているため、dry-run 表記を追加すると既存テストの修正が必要。これは意図された変更であり問題ない。

## S1: 正しい問題か？

**Category**: Problem Definition
**Finding**: 問題は「dry-run の出力が通常実行と区別できない」であり、spec はそれを直接解決する。問題を再定義する余地はない。
**Recommendation**: そのまま進行。

## S2: 代替アプローチの検討

**Category**: Approach Alternatives
**Finding**: 3つのアプローチを検討:

1. **各行にプレフィックス追加** — `(dry-run) create: /key` のように各行に表記。一目瞭然。grep でフィルタ可能。
2. **サマリー行のみに表記** — `4 created, 0 updated ... (dry-run)` のようにサマリーのみ。行ごとのフォーマットは変わらないので、スクリプトで出力をパースしている場合の後方互換性が高い。
3. **ヘッダー行を追加** — 出力の冒頭に `=== DRY RUN ===` を出す。各行は変えない。

**Recommendation**: spec は FR-001（サマリー）と FR-002（各行）の両方を要求しており、アプローチ 1 が最も合致する。サマリー行にも `(dry-run)` を付与すれば FR-001 も満たせる。

## S3: 既存テストへの影響

**Category**: Integration Impact
**Finding**: `sync_test.go:53` が `"create: k1"` を `strings.Contains` でアサートしている。dry-run プレフィックスを追加すると `"(dry-run) create: k1"` になるため、`Contains` チェックは引き続きパスする（部分文字列一致）。ただし、dry-run 表記の存在自体をテストする新規テストケースが必要。

`ssm_test.go` の dry-run テストも同様に `strings.Contains` ベースなので互換性あり。

**Evidence**: `sync_test.go:53`: `strings.Contains(out, "create: k1")` — `"(dry-run) create: k1"` にも一致する。

**Recommendation**: 既存テストは壊れないが、dry-run 表記の有無を明示的にテストする新規アサーションを追加すべき。

## S4: CI/CD パイプラインでの出力パース

**Category**: Risk & Failure Modes
**Finding**: spec の Assumptions に「プレーンテキストで判別できる必要がある」とある。CI/CD でスクリプトが `psm sync --dry-run` の出力をパースしている場合、フォーマット変更が破壊的になる可能性がある。ただし psm v0.0.1 でユーザーベースは極めて小さく、出力フォーマットの安定性保証はない。
**Recommendation**: YAGNI。破壊的変更を恐れる段階ではない。プレフィックス方式で進行。

## Items Requiring PoC

なし。変更箇所が明確で、実装前に検証が必要な未知要素はない。

## Constitution Impact

なし。Go コードの変更であり、Test-First が適用される。

## Recommendation

Proceed to plan. 実装は `sync.go` の `execute` 関数内の `fmt.Fprintf` 呼び出しにプレフィックスを追加するだけ。テストは `sync_test.go` に dry-run 表記のアサーションを追加。
