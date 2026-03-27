# CLI Flag Removal Checklist: --store フラグの除去

**Purpose**: --store フラグ除去に関する要件の品質・完全性・一貫性を検証する
**Created**: 2026-03-27
**Feature**: [spec.md](../spec.md)

## Requirement Completeness

- [ ] CHK001 エラーメッセージの具体的な文言が spec に定義されているか？ [Completeness, Spec FR-002]
- [ ] CHK002 `--store` 検出の対象パターンが網羅的に定義されているか？（`--store ssm`, `--store sm`, `--store=ssm`, `--store=sm` 等） [Completeness, Spec Edge Cases]
- [ ] CHK003 影響を受けるドキュメントファイルの完全なリストが spec に含まれているか？ [Completeness, Spec FR-005/FR-006]
- [ ] CHK004 `example/test.sh` の更新が要件として明記されているか？ [Completeness, Gap]

## Requirement Clarity

- [ ] CHK005 「廃止エラー」の終了コードが明確に指定されているか？ [Clarity, Spec User Story 2]
- [ ] CHK006 SSM ハードコードの具体的な意味（Config からの除去 vs デフォルト値設定）が spec で区別されているか？ [Clarity, Spec FR-003]

## Requirement Consistency

- [ ] CHK007 User Story 2 のエラーメッセージと contracts/cli-schema.md のエラーメッセージが一致しているか？ [Consistency]
- [ ] CHK008 spec の FR-007 と User Story 3 の Acceptance Scenarios でテスト更新の範囲が一致しているか？ [Consistency]

## Acceptance Criteria Quality

- [ ] CHK009 各ユーザーストーリーの Acceptance Scenarios が Given-When-Then 形式で記述されているか？ [Measurability, Spec User Stories]
- [ ] CHK010 Success Criteria が実装なしで検証可能な形式か？ [Measurability, Spec SC-001~SC-005]

## Scenario Coverage

- [ ] CHK011 `--store` と `--store=value`（= 記法）の両方が要件でカバーされているか？ [Coverage, Gap]
- [ ] CHK012 `--store` 以外のフラグとの組み合わせ時の振る舞いが定義されているか？ [Coverage, Spec Edge Cases]
- [ ] CHK013 `--store` なしでの正常動作シナリオが sync/export 両方に定義されているか？ [Coverage, Spec User Story 1]

## Dependencies and Assumptions

- [ ] CHK014 #27（SM 除去）完了が前提条件として明記されているか？ [Assumption, Spec Assumptions]
- [ ] CHK015 Store interface 保持の方針が spec と plan で一致しているか？ [Consistency, Spec FR-009/Assumptions]

## Notes

- Check items off as completed: `[x]`
- Items are numbered sequentially for easy reference
