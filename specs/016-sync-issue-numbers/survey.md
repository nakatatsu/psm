# Survey: Issue 番号と SpecKit 機能番号の同期

**Date**: 2026-03-17
**Spec**: [spec.md](./spec.md)

## Summary

方向性は妥当。Issue 番号と SpecKit 機能番号の同期は、「Issue 起点のワークフロー強制」という手段で正しく解決できる。変更範囲はスキル定義（プロンプト）が主であり、シェルスクリプトへの変更は最小限。リスクは低い。

ただし 2 点注意が必要: (1) FR-006「`--number` 廃止」の表現が Assumptions の「内部インターフェースとして残す」と一見矛盾するため、実装時に混乱しないよう明確化が望ましい。(2) Issue のタイトルと説明文の両方から short-name を生成する際、説明文が長い場合の扱いを実装レベルで考慮する必要がある。

## S1: 問題定義の妥当性

**Category**: A. Problem Reframing — Problem Definition
**Finding**: spec は正しい問題を解いている。本質的な目標は「Issue と spec/ブランチの追跡可能性（traceability）」であり、番号の一致はそれを実現するもっともシンプルな手段。代替案（番号独立 + メタデータでの紐付け、Issue URL を spec.md に記録するだけ等）は追跡性を得るために追加のルックアップが必要になり、Simplicity First 原則に反する。
**Recommendation**: 現在のアプローチを維持。
**Evidence**: 既存の specs/ ディレクトリ（001, 003, 004, 005, 007, 008, 012, 016）を確認。番号がブランチ名・ディレクトリ名に直接埋め込まれており、番号一致が最も自然な追跡手段。

## S2: Issue 必須化の妥当性

**Category**: A. Problem Reframing — Hidden Assumptions
**Finding**: 「SpecKit を使う規模の作業なら Issue を立てるべき」という前提は、このプロジェクトの実態に合致する。既存の spec 付き機能（001, 003, 005, 007, 008, 012）はすべて対応する Issue が存在する。Issue なしで SpecKit を使うケースは実質ゼロ。
**Recommendation**: Issue 必須化を進める。ただし、既存の番号が Issue 番号と一致しない spec（001, 003 等）はレガシーとして放置し、マイグレーション不要。
**Evidence**: `gh issue list` と `ls specs/` の突合。既存の spec はすべて Issue 起点で作られている。

## S3: `--number` 廃止と内部インターフェースの整合性

**Category**: B. Solution Evaluation — Cost & Complexity
**Finding**: FR-006 は「`--number` による手動番号指定を廃止」、Assumptions は「`--number` 自体はスクリプトの内部インターフェースとして残す」。意図は「ユーザー（スキル利用者）に `--number` を直接使わせない。スキル定義が Issue URL から番号を抽出して内部的に `--number` を渡す」という二層構造。これは正しいが、FR-006 の文言だけ読むと「スクリプトから `--number` を削除する」と誤解される可能性がある。
**Recommendation**: plan フェーズで FR-006 の実装範囲を明確にする。「スキル定義から `--number` の直接指定手順を削除する」と「`create-new-feature.sh` の `--number` オプションは維持する」を分けて記述する。
**Evidence**: `create-new-feature.sh` を確認。`--number` は他のスクリプトや自動化から呼ばれる可能性があり、スクリプトレベルでの削除は過剰。

## S4: 変更対象の範囲確認

**Category**: D. Integration & Governance — Integration Impact
**Finding**: 変更対象は以下の 1 ファイルが主:
- `.claude/skills/speckit.specify/SKILL.md` — Issue URL の必須化、番号抽出ロジック、Issue 情報取得手順の追加

`create-new-feature.sh` への変更は不要（`--number` を受け取る既存インターフェースをそのまま利用）。他のスキル（clarify, survey, plan, analyze）は `speckit.specify` を直接呼ばないため影響なし。`coding-workflow` スキルは参照のみで影響なし。
**Recommendation**: 変更範囲を SKILL.md のプロンプト修正に集中させる。
**Evidence**: `grep -r "speckit.specify"` で参照箇所を確認。`create-new-feature.sh` は speckit.specify の SKILL.md からのみ呼ばれる。

## S5: Constitution 適合性

**Category**: D. Integration & Governance — Constitution Compliance
**Finding**: Constitution の Test-First 原則は Go コードに適用される。本機能の変更対象はスキル定義（Markdown プロンプト）であり、Go コードの変更はない。Test-First の直接適用外だが、動作確認（手動テスト）は必要。
**Recommendation**: 実装後に Issue URL あり/なしの両パターンで `/speckit.specify` を実行し、期待通りの動作を確認する。
**Evidence**: Constitution v3.0.0 の III. Test-First セクション。対象は `go test` と明記されている。

## Items Requiring PoC

なし。すべての構成要素（`gh issue view` による Issue 取得、`--number` によるブランチ番号指定、URL からの番号抽出）は既に動作実績がある。

## Constitution Impact

修正不要。本機能はスキル定義（プロンプト）の変更であり、Go コードや Constitution が管理する開発プラクティスに影響しない。

## Recommendation

問題なし。`/speckit.plan` に進んでよい。実装時に S3（FR-006 の文言明確化）を考慮すること。
