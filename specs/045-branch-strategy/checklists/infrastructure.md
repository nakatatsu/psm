# Infrastructure Checklist: ブランチ戦略の決定と実装

**Purpose**: ブランチ保護ルール、リポジトリ設定、CodeQL、git hooks、SKILL の要件品質を検証する
**Created**: 2026-04-03
**Feature**: [spec.md](../spec.md)

## Requirement Completeness

- [ ] CHK001 Are protection rules specified for all four branch types (develop, release-\*, hotfix-\*, main)? [Completeness, Spec FR-001〜FR-005, FR-009]
- [ ] CHK002 Are all PR sub-options (approval count, stale review dismissal, last push approval, comment resolution) explicitly specified for each branch type? [Completeness, Survey S1 #8]
- [ ] CHK003 Is the bypass actor configuration (Repository admin, bypass mode) defined for all new Rulesets? [Completeness, Spec FR-007]
- [ ] CHK004 Are all `github_repository` settings (merge strategy, delete_branch_on_merge, vulnerability_alerts, security_and_analysis) explicitly specified? [Completeness, Spec FR-011〜FR-014]
- [ ] CHK005 Is the CodeQL workflow trigger configuration (push, PR, schedule) specified? [Completeness, Spec FR-015]
- [ ] CHK006 Is the commit-msg hook pattern (`#<issue-no>`) and exception cases (merge commits) specified? [Completeness, Spec FR-016, Research R3]
- [ ] CHK007 Is the SKILL content scope defined (GitFlow flow, naming rules, commit message format)? [Completeness, Spec FR-017]

## Requirement Clarity

- [ ] CHK008 Is "main と同等の保護ルール" in FR-004/FR-005 quantified with specific settings rather than referencing another branch? [Clarity, Spec FR-004/FR-005]
- [ ] CHK009 Is the `required_code_scanning` fallback strategy clear — what specific action triggers the fallback to CI status checks? [Clarity, Spec FR-015]
- [ ] CHK010 Is "CI ステータスチェック" defined with the exact check context name(s) to require? [Clarity, Spec FR-003]
- [ ] CHK011 Is the `core.hooksPath` setup responsibility clearly assigned (developer manual setup vs automated)? [Clarity, Spec FR-016]

## Requirement Consistency

- [ ] CHK012 Are the protection rules for `release-*` and `hotfix-*` consistently specified as identical in both spec and survey? [Consistency, Spec FR-004/FR-005, Survey S1 #8]
- [ ] CHK013 Is the `required_linear_history = false` change in FR-009 consistent with the merge commit decision in FR-011? [Consistency, Spec FR-009/FR-011]
- [ ] CHK014 Is the deletion protection in FR-006 consistent with `delete_branch_on_merge = true` in FR-012 (protected branches won't be auto-deleted)? [Consistency, Spec FR-006/FR-012]
- [ ] CHK015 Is the `for_each = var.protected_repositories` pattern consistently applied across all new Terraform resources? [Consistency, Research R5]

## Acceptance Criteria Quality

- [ ] CHK016 Do acceptance scenarios cover all four protected branch types (develop, release-\*, hotfix-\*, main)? [Coverage, Spec User Story 1]
- [ ] CHK017 Is there an acceptance criterion for the `terraform import` step before `terraform apply`? [Gap, Research R1]
- [ ] CHK018 Is there an acceptance criterion for the CodeQL fallback scenario (Provider bug → CI status check)? [Gap, Spec FR-015]
- [ ] CHK019 Is there an acceptance criterion for the commit-msg hook rejecting commits without issue numbers? [Gap, Spec FR-016]

## Scenario Coverage

- [ ] CHK020 Are requirements defined for the scenario where `terraform apply` partially fails (some resources created, others not)? [Exception scenario, Gap]
- [ ] CHK021 Are requirements defined for the scenario where an existing repository setting conflicts with the new Terraform-managed setting? [Exception scenario, Gap]
- [ ] CHK022 Are rollback requirements defined if branch protection rules cause unexpected workflow failures? [Recovery scenario, Gap]
- [ ] CHK023 Are requirements defined for repositories other than `psm` in `protected_repositories`? [Alternate scenario, Spec Assumptions]

## Edge Case Coverage

- [ ] CHK024 Is the behavior specified when a branch name matches multiple Ruleset patterns (e.g., a branch named `release-hotfix-v1`)? [Edge Case, Gap]
- [ ] CHK025 Is the behavior specified when the `ci` status check doesn't exist yet on a new branch? [Edge Case, Gap]
- [ ] CHK026 Is the interaction between `secret_scanning_push_protection` and SOPS-encrypted files specified? [Edge Case, Gap]

## Dependencies and Assumptions

- [ ] CHK027 Is the assumption "Terraform の実行環境と権限は既に整備されている" verified against the actual `.tmp/mynote/infrastructure/github/` setup? [Assumption, Spec Assumptions]
- [ ] CHK028 Is the assumption "`psm` リポジトリは `protected_repositories` 変数に含まれている" verified? [Assumption, Spec Assumptions]
- [ ] CHK029 Is the Terraform GitHub Provider version requirement (`~> 6.0`) explicitly stated as a dependency? [Dependency, Plan Technical Context]
- [ ] CHK030 Is the dependency on GitHub Actions (CodeQL Action version) specified? [Dependency, Spec FR-015]

## Notes

- Check items off as completed: `[x]`
- Items CHK001-CHK007 cover completeness of all 17 FRs
- Items CHK016-CHK019 highlight gaps in acceptance criteria that should be addressed before implementation
- Items CHK020-CHK026 identify scenarios not yet covered in the spec
