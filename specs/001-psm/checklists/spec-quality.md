# Spec Quality Checklist: psm

**Purpose**: Validate requirement completeness, clarity, consistency, and coverage across all spec and design documents
**Created**: 2026-03-08
**Feature**: `specs/001-psm/spec.md`

## Requirement Consistency

- [x] CHK001 Is the handling of `null` YAML values consistent between FR-020 (skip) and FR-021 (exit code 1)? FR-020 says "skip", FR-021 says terminate. [Conflict, FR-020 vs FR-021] — **FIXED**: FR-020 removed (merged into FR-021). Edge Cases updated to exit code 1.
- [x] CHK002 Is the handling of empty key strings consistent between Edge Cases ("skip") and FR-021 (exit code 1)? Edge Cases says skip, FR-021 says terminate. [Conflict, Edge Cases vs FR-021] — **FIXED**: Edge Cases updated to exit code 1.
- [x] CHK003 Is the `AWS_PROFILE` policy consistent across all documents? quickstart.md line 14 says "対象アカウントの切り替えは `AWS_PROFILE` で行う", contradicting FR-002/Assumptions. [Conflict, quickstart.md vs spec.md FR-002] — **FIXED**: quickstart.md updated to reference `_psm.profile`.
- [x] CHK004 Is Clarifications Q6 consistent with current spec? Q6 says "AWS_PROFILE で対象アカウントを切り替える", contradicting FR-002. [Conflict, Clarifications Q6 vs FR-002] — **FIXED**: Clarifications Q6 updated.
- [x] CHK005 Are output format requirements consistent between cli-schema.md and spec.md FR-014/FR-015/FR-018? — **VERIFIED**: cli-schema.md stdout format (`create: {key}`, summary `{N} created, ...`, dry-run 同一形式) は FR-014/015/018 と完全一致。
- [x] CHK006 Are prune scope requirements consistent between spec.md FR-008, cli-schema.md Prune Scope, and User Story 4 Scenario 4? — **VERIFIED**: cli-schema.md は `spec.md FR-008 参照` と明記。FR-008 = アカウント全体、プレフィックス制限なし。US4 Scenario 4 も同一。
- [x] CHK007 Are error output format requirements consistent between spec.md FR-019, cli-schema.md stderr Format, and data-model.md Summary.Failed? — **VERIFIED**: cli-schema.md stderr `error: {key}: {error message}` = FR-019 の例と一致。data-model.md Summary.Failed はカウントのみで形式矛盾なし。

## Requirement Completeness

- [x] CHK008 Is there a functional requirement for the `--help` flag? SC-004 references it but no FR defines its behavior. — **FIXED**: FR-024 追加。Go `flag` パッケージのデフォルト動作に従い、終了コード 0。
- [x] CHK009 Is there a functional requirement for goroutine concurrency (max 10)? store-interface.md specifies it but no FR captures it. — **RESOLVED**: 並行数は実装詳細であり FR に含めない。Clarifications Session 2026-03-09 で「固定値、設定不可、YAGNI」と明記。
- [x] CHK010 Are SSM-specific parameter type requirements (SecureString) documented in cli-schema.md? FR-004 specifies it but cli-schema.md SSM section doesn't mention it. — **RESOLVED**: cli-schema.md は CLI インターフェース形式定義のみ。SSM の内部動作は FR-004 + store-interface.md で定義済み。cli-schema.md への追記は不要（責務外）。
- [x] CHK011 Is the SM deletion behavior (ForceDeleteWithoutRecovery=true) specified in spec.md? User Story 4 Scenario 2 says "即時削除" but no FR captures this explicitly. — **FIXED**: FR-025 追加。ForceDeleteWithoutRecovery=true を明記。
- [x] CHK012 Are the YAML reserved keys (`_psm`, `sops`) exhaustively listed? Is any other key prefix reserved? — **FIXED**: FR-026 追加。予約キーは `sops` のみ。`_psm` は現仕様に存在しない（旧仕様の残骸）。
- [x] CHK013 Is the behavior when multiple YAML files are passed as arguments defined? — **FIXED**: Edge Cases に追記。1 ファイルのみ、複数指定は FR-012 エラー。
- [x] CHK014 Is the maximum file size or key count limitation specified? SC-001 mentions 100 keys but no upper bound is defined. — **RESOLVED**: Clarifications Session 2026-03-09 で「上限なし、AWS rate limit に依存、SC-001 は基準値であり制限ではない」と明記。

## Requirement Clarity

- [x] CHK015 Is "エラーメッセージを表示しそのキーをスキップ" (Edge Cases for null and empty key) quantified? Does "skip" mean the key is silently ignored, logged to stderr, or counted as failed? — **RESOLVED**: CHK001/002 で修正済み。null/空キーは FR-020 バリデーションエラーとして即終了（スキップではない）。
- [x] CHK016 Is the sort order of diff output lines specified? Are they alphabetical, insertion-order, or undefined? — **RESOLVED**: Clarifications Session 2026-03-09 で「未定義、AWS API 返却順に依存」と明記。
- [x] CHK017 Is "使い方を表示" (FR-012) specified? What exactly constitutes the usage message? — **RESOLVED**: Clarifications Session 2026-03-09 で「Go flag パッケージのデフォルト出力、カスタム不要」と明記。FR-024 も参照。
- [x] CHK018 Is the `_psm.profile` validation defined? What happens if profile name is empty string or invalid? — **RESOLVED**: `_psm` は現仕様に存在しない（旧仕様の残骸、CHK012 参照）。`--profile` の無効値は AWS SDK エラーに委ねる（Clarifications Session 2026-03-09）。
- [x] CHK019 Is the SM CreateSecret behavior on name collision with a recently deleted secret defined? (Secrets Manager has a recovery window concept.) — **RESOLVED**: FR-025 で ForceDeleteWithoutRecovery=true を規定。復旧期間のコンフリクトは発生しない。Clarifications Session 2026-03-09 で詳細記載。

## Acceptance Criteria Quality

- [x] CHK020 Do User Story 1 acceptance scenarios cover the `_psm.profile` flow? — **RESOLVED**: `_psm` は現仕様に存在しない（CHK012 参照）。`--profile` の動作は FR-002 で定義済み。US1 は SSM 同期の基本動作に集中しており、プロファイル指定はフラグパースのテスト（T005）でカバー。
- [x] CHK021 Does User Story 3 have a scenario for dry-run with update operations (not just create and delete)? — **RESOLVED**: FR-015 が「出力形式は通常実行と --dry-run で同一」と規定。update の表示形式は FR-014 で定義済み。US3 Scenario 1 の「値が変わっていない既存キーは表示されず」は暗黙に update ケースを含む（変わっているキーは表示される）。統合テスト T016 でカバー。
- [x] CHK022 Are acceptance scenarios for input validation errors (FR-020) defined? No user story covers the validation-first gate. — **RESOLVED**: FR-020 自体が具体的なエラー条件と期待動作を定義（0件、重複、null/map/array、空キー → エラーメッセージ + 終了コード 1）。ユニットテスト T003 で全条件をテーブルドリブンテスト。User Story 横断の品質ゲートであり特定ストーリーに属さない。
- [x] CHK023 Does User Story 2 have a scenario for partial failure (like User Story 1 Scenario 4)? — **RESOLVED**: US2 は US1 の sync ロジックを再利用（Store 実装のみ異なる）。部分失敗の動作は FR-010/FR-011 で規定済みであり、T009 の統合テストでカバー。SM 固有の部分失敗シナリオは T013 でテスト。

## Edge Case Coverage

- [x] CHK024 Is the behavior defined when AWS credentials expire mid-sync (after some keys succeed)? — **RESOLVED**: Edge Cases に追記。FR-010 に従い個別キーの失敗として扱い、残りの処理を続行する。
- [x] CHK025 Is the behavior defined when SSM `GetParametersByPath` returns parameters that the caller lacks permission to decrypt? — **RESOLVED**: Edge Cases に追記。GetAll 自体が失敗した場合は致命的エラーとして即終了。個別キーの権限エラーは FR-010/FR-019 に従う。
- [x] CHK026 Is the behavior defined for YAML keys with characters that AWS rejects (e.g., SSM parameter name validation rules)? — **RESOLVED**: Edge Cases に追記。AWS API エラーとして FR-010/FR-019 に従い処理。psm 側でのキー名バリデーションは行わない。
- [x] CHK027 Is the behavior defined when `GetAll` fails entirely (e.g., no permissions)? Is it a fatal error or does it fall back to create-all? — **RESOLVED**: Edge Cases に追記。致命的エラーとし即座に終了コード 1。create-all へのフォールバックは行わない。
- [x] CHK028 Is the behavior defined for extremely long YAML values exceeding AWS limits (SSM: 8KB for advanced, 4KB standard; SM: 64KB)? — **RESOLVED**: Edge Cases に追記。AWS API エラーとして FR-010/FR-019 に従い処理。

## Non-Functional Requirements

- [x] CHK029 Are performance requirements specified beyond SC-001 (100 keys)? What is the expected behavior at 1000+ keys? — **RESOLVED**: Clarifications Session 2026-03-09 で「上限なし、AWS rate limit に依存」と明記。SC-001 は動作確認基準であり制限ではない。
- [x] CHK030 Is the concurrency limit (10) configurable or fixed? Is this a requirement or an implementation detail? — **RESOLVED**: Clarifications Session 2026-03-09 で「固定値、設定不可、YAGNI。実装詳細であり FR に含めない」と明記。
- [x] CHK031 Are timeout requirements for individual AWS API calls specified? — **RESOLVED**: Clarifications Session 2026-03-09 で「AWS SDK デフォルトに委ねる、psm 側でのタイムアウト設定不要」と明記。

## Dependencies and Assumptions

- [x] CHK032 Is the assumption that `GetParametersByPath` with path `/` returns ALL parameters in the account validated? Does this work across all AWS regions? — **RESOLVED**: Clarifications Session 2026-03-09 で「Recursive=true と合わせて全パラメータを返す、AWS ドキュメントで確認済み」と明記。research.md R2 参照。
- [x] CHK033 Is the assumption that SOPS always produces valid YAML documented and tested? What if SOPS output is malformed? — **RESOLVED**: Clarifications Session 2026-03-09 で「psm は有効な YAML を前提とするが、パース失敗時はエラーメッセージ + 終了コード 1」と明記。Assumptions に「入力 YAML は SOPS によって事前に復号されている」と記載済み。
- [x] CHK034 Is the dependency on `yaml.v3` (for `yaml.Node` API) documented in spec-level requirements, or only in research.md? — **RESOLVED**: plan.md Primary Dependencies に `gopkg.in/yaml.v3 v3.0.1` を明記。research.md R3 で選定理由を詳述。spec はツール要件を定義するものであり、ライブラリ選定は plan/research の責務。

## Notes

- CHK001 and CHK002 are critical contradictions that were resolved before implementation
- CHK003 and CHK004 are documentation inconsistencies from the spec revision history
- CHK005-007 were consistency checks that passed verification
- CHK008-034 were resolved in Session 2026-03-09: gaps filled with FR-024/025/026, edge cases documented, ambiguities clarified
