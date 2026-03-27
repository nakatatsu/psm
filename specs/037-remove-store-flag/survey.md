# Survey: --store フラグの除去

**Date**: 2026-03-27
**Spec**: [spec.md](spec.md)

## Summary

本機能は #27 で SM サポートを除去した後の自然な後続作業であり、冗長な `--store ssm` 指定を除去して CLI をシンプルにするもの。YAGNI 原則・Simplicity First 原則に完全に合致しており、問題定義・方向性ともに適切。リスクは低く、変更の影響範囲も明確に限定されている。

## S1: 問題定義の妥当性

**Category**: Problem Reframing
**Finding**: 問題定義は正しい。SSM のみがサポート対象であるにもかかわらず、毎回 `--store ssm` を要求するのは不要な認知負荷。#27 で SM を除去した時点でこのフラグは形骸化している。
**Recommendation**: そのまま進行。
**Evidence**: `main.go:219-224` のバリデーションで `ssm` のみ受け付けており、選択肢が 1 つしかないフラグは不要。constitution の YAGNI 原則（「No feature flags, configuration toggles, or plugin systems unless the current task explicitly requires them」）にも合致。

## S2: 代替アプローチの検討

**Category**: Approach Alternatives
**Finding**: 3 つのアプローチを検討した:
1. **フラグを完全除去（spec の提案）** — シンプル、YAGNI に合致
2. **フラグをオプショナル化（デフォルト ssm）** — 後方互換性は高いが、将来使わないフラグを残すことになる
3. **フラグを除去しエイリアスコマンドで対応** — 過剰設計

**Recommendation**: アプローチ 1（完全除去）が最適。アプローチ 2 は YAGNI 違反、アプローチ 3 は過剰設計。
**Evidence**: Constitution の Simplicity First: 「No abstractions until needed」。将来別ストアが必要になれば、その時点で新たに設計するのが正しい。

## S3: 後方互換性のリスク

**Category**: Risk & Failure Modes
**Finding**: `--store ssm` を含む既存スクリプトやワークフローが壊れるリスクがある。ただし spec では `--store` 指定時にエラーメッセージでガイダンスを表示する方針。
**Recommendation**: spec 通りエラーメッセージで対応。ただしエラーメッセージは「`--store` フラグは廃止されました。省略してください。SSM がデフォルトで使用されます」のように具体的な修正方法を含めること。
**Evidence**: psm は個人／小規模チーム向けツールであり、大規模な後方互換性リスクは低い。

## S4: 影響範囲の確認

**Category**: Integration Impact
**Finding**: 変更が必要なファイルは明確:
- コード: `main.go`（フラグ定義・バリデーション・switch 文）、`store.go`（Config 構造体の Store フィールド）
- テスト: `main_test.go`（20 テストケースすべて）
- ドキュメント: `README.md`, `example/README.md`, `example/test.sh`
- 影響しないファイル: `ssm.go`, `sync.go`, `export.go`, `ssm_test.go`（Store interface 経由で使用しているため）

**Recommendation**: Config 構造体から `Store` フィールドを除去し、`main.go` で直接 `NewSSMStore` を呼ぶように変更。
**Evidence**: `ssm_test.go` は `NewSSMStore(cfg)` を直接使用しており CLI パースに依存していない。`sync.go` と `export.go` は Store interface を受け取るのみ。

## S5: Constitution 整合性

**Category**: Constitution Compliance
**Finding**: 本変更は constitution の以下の原則に完全に合致:
- YAGNI: 選択肢が 1 つしかないフラグを除去
- Simplicity First: 不要な抽象化の除去
- Test-First: テスト更新が必要（main_test.go）

**Recommendation**: Constitution の修正は不要。
**Evidence**: `.specify/memory/constitution.md` の全原則を確認済み。

## Items Requiring PoC

なし。すべての変更は既存コードの削除・簡素化であり、検証不要。

## Constitution Impact

修正不要。本変更は constitution の原則を実現するもの。

## Recommendation

問題なし。`/speckit.plan` に進行可能。
