# CI/CD Checklist: デプロイ方式の変更

**Purpose**: リリースワークフロー変更の要件品質を検証する
**Created**: 2026-04-04
**Feature**: [spec.md](../spec.md)

## Requirement Completeness

- [ ] CHK001 Are trigger conditions fully specified for all PR merge scenarios (merged vs closed-without-merge)? [Completeness, Spec FR-001]
- [ ] CHK002 Are both `release-` and `hotfix-` branch name patterns explicitly listed in all relevant requirements? [Completeness, Spec FR-002]
- [ ] CHK003 Is the exact version extraction logic specified for both `release-X.Y.Z` and `hotfix-X.Y.Z` patterns? [Completeness, Spec FR-003]
- [ ] CHK004 Is the required tag format (`v` prefix + version) explicitly defined? [Completeness, Spec FR-003]
- [ ] CHK005 Are the permissions required by the workflow (e.g., `contents: write`) documented? [Gap]

## Requirement Clarity

- [ ] CHK006 Is the semver validation scope clearly bounded (X.Y.Z only, no pre-release/metadata)? [Clarity, Spec FR-004]
- [ ] CHK007 Is "error として終了" defined with specific expected behavior (exit code, error message format)? [Clarity, Spec FR-004, FR-005]
- [ ] CHK008 Is the distinction between "skip" (FR-007) and "error" (FR-004, FR-005) clearly defined? [Clarity]

## Requirement Consistency

- [ ] CHK009 Is FR-007 ("スキップ") consistent with User Story 2 ("失敗して即停止する")? These describe different behaviors for the same scenario. [Conflict, Spec FR-007 vs User Story 2]
- [ ] CHK010 Are `hotfix-` branches treated consistently across all requirements (FR-002, FR-003, FR-007)? [Consistency]
- [ ] CHK011 Is the User Story 1 acceptance scenario consistent with FR-003 (tag format `vX.Y.Z`)? [Consistency]

## Acceptance Criteria Quality

- [ ] CHK012 Does User Story 1 have acceptance scenarios covering both `release-` and `hotfix-` branches? [Coverage, Spec User Story 1]
- [ ] CHK013 Are error message content requirements specified for FR-004 (invalid semver) and FR-005 (duplicate tag)? [Measurability]
- [ ] CHK014 Is SC-001 measurable without ambiguity ("追加の手動操作なし" — does this include branch creation)? [Measurability, Spec SC-001]

## Scenario Coverage

- [ ] CHK015 Are requirements defined for what happens when a PR is closed without merging? [Gap]
- [ ] CHK016 Is the behavior specified when GoReleaser fails after tag creation (partial failure / rollback)? [Gap]
- [ ] CHK017 Are concurrent release scenarios addressed (two release PRs merged in quick succession)? [Gap]

## Edge Case Coverage

- [ ] CHK018 Is behavior defined for branch names with extra segments (e.g., `release-1.0.0-rc1`)? [Edge Case, Spec Edge Cases]
- [ ] CHK019 Is behavior defined for branch names with leading zeros (e.g., `release-01.00.00`)? [Edge Case]

## Dependencies and Assumptions

- [ ] CHK020 Is the assumption that all main merges go through PRs explicitly stated and justified? [Assumption]
- [ ] CHK021 Is the dependency on existing `.goreleaser.yaml` configuration explicitly documented? [Assumption, Spec Assumptions]
- [ ] CHK022 Is the GITHUB_TOKEN limitation (cannot trigger other workflows) documented as a constraint? [Assumption]

## Notes

- CHK009 is a critical finding: FR-007 says "skip" but User Story 2 says "fail immediately" — these need reconciliation
- Check items off as completed: `[x]`
