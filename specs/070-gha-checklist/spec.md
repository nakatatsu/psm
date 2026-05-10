# Feature Specification: GHA ワークフローの CI/CD チェックリスト準拠化

**Feature Branch**: `070-gha-checklist`
**Created**: 2026-05-09
**Status**: Draft
**Input**: https://github.com/nakatatsu/psm/issues/70

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: 全ワークフロー（`ci.yml` / `codeql.yml` / `integration-test.yml` / `release.yml`）に `concurrency` ブロックを定義し、`cancel-in-progress` を真偽値で明示する。副作用なし（`ci.yml` / `codeql.yml`）は `true`、副作用あり（`integration-test.yml` / `release.yml`）は `false` とする。グループキーは「ワークフロー名 + 対象 ref」を基本とする。
- **FR-002**: 全ワークフローの全ジョブに `timeout-minutes` を設定する。値は正常時所要時間に安全マージンを加えたもの。
- **FR-003**: 全ワークフローでワークフロー全体の既定権限（`permissions:`）を最小（`contents: read` 相当）として明示する。書き込み相当の権限（`id-token: write` / `contents: write` / `security-events: write` 等）は必要なジョブでのみ個別に上書きする。
- **FR-004**: 利用する全ての外部 Action のうち、特権または副作用があり、かつ信頼境界の外側のものは、コミット SHA で固定する。バージョンの可読性を保つため、SHA の右側にバージョンタグをコメントとして併記する。
- **FR-005**: 補助的な外部 Action であっても `latest` や可変メジャータグ（例: `@v3`）への依存を持たない。
- **FR-006**: 各ワークフロー内に複数箇所現れるバージョン値（Go バージョンおよび利用ツールのバージョン）は、ファイル内の `env:` ブロックに集約し、参照側で展開する。
- **FR-007**: 自動ゲートとして組み込まれている品質チェックは、失敗時にジョブを失敗させる。`continue-on-error: true` および `if: always()` による警告化を行わない。
- **FR-008**: SHA で固定した依存先について、Dependabot 等の更新ボットを設定し、定期的な更新経路を担保する。対象は GitHub Actions と Go modules の双方を含む。
- **FR-009**: チェックリストに準拠した模範ワークフローをリポジトリ内に配置する。GitHub Actions に読み込まれない拡張子とすることで、誤実行を防ぐ。

## Constraints *(mandatory)*

- **C-001**: 修正範囲は CI/CD チェックリスト（`central-library/documents/procedures/cicd-checklist.md`）に明記されたゲート水準項目に限定する。チェックリストに無い改善（Nitpicking 系）は本作業に含めない。
- **C-002**: ワークフロー名（`name:`）は変更しない。required status checks やブランチ保護設定への影響を避けるため。
- **C-003**: トリガ（`on:`）の構成は変更しない。本作業の範囲はトリガの再設計ではないため。
- **C-004**: 既存ジョブが依存している外部リソース（OIDC のロール ARN、SOPS / age バイナリのリリース URL、AWS リージョン等）の参照方法は変更しない。
- **C-005**: 修正前後でワークフローの機能的挙動は等価であること。ただし FR-007 により vulnerability check のみは失敗時の挙動が「警告」から「ジョブ失敗」へ変わる。

## Success Criteria *(mandatory)*

- **SC-001**: CI/CD チェックリストの全 6 章のうち、強制ゲート水準項目（推奨項目を除く）について、4 ワークフロー全てが充足判定となる。
- **SC-002**: 連続 push を行った際、副作用なしワークフローでは古い実行がキャンセルされ、副作用ありワークフローでは古い実行が完了するまで新規実行が待機する。
- **SC-003**: 主要な外部 Action（特権／副作用ありかつ信頼境界外）を git grep した際、`@vN`（可変メジャータグ）形式の参照が 0 件である。
- **SC-004**: Dependabot により外部 Action と Go modules の更新 PR が、設定週次間隔で自動的に作成される。
- **SC-005**: 模範ワークフローファイルは GitHub Actions のワークフロー一覧に表示されない。

## Edge Cases

- 外部 Action の SHA 固定後、Dependabot による自動更新 PR がマージされた場合に、SHA とコメント併記の整合性が保たれること（Dependabot は両方を同時に更新する）。
- リリース副作用ワークフローで `cancel-in-progress: false` としたため、誤った release-* / hotfix- ブランチへの連続 push 時にキューが詰まる。タイムアウトの上限に依存することを前提とする。
- vulnerability check（`govulncheck`）が脆弱性 DB の一時的不整合等で誤検知した場合、CI 全体が失敗する。例外運用（一時的な無効化）が必要になった際は CI/CD チェックリスト §6.3 に従い理由と昇格要件を記録する。

## Assumptions

- 本リポジトリの GitHub Actions ランナーは GitHub-hosted ランナー（`ubuntu-latest`）を継続利用する。
- ブランチ保護および required status checks の設定はリポジトリ管理者によって別途維持されており、本作業のスコープ外とする。
- Dependabot を有効化するためのリポジトリ／組織レベルの権限は確保されている。
- CI/CD チェックリスト（`central-library/documents/procedures/cicd-checklist.md`）の現行版を判定基準として採用する。
