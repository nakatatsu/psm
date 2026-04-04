# Survey: ブランチ戦略の決定と実装

**Date**: 2026-04-02
**Spec**: [spec.md](./spec.md)

## Summary

GitHub Repository Ruleset の全オプションを網羅的に調査し、今回のブランチ保護設計への採否を判断した。GitFlow + Terraform（GitHub Provider）による管理方針は妥当であり、根本的な方向転換は不要。ただし、Terraform Provider の一部機能（merge queue、push ruleset）に制約があるため注意が必要。また、パブリックリポジトリで GitHub Free プランを使用している前提のため、Enterprise 限定機能は使用不可。

---

## S1: GitHub Repository Ruleset — 全ルールオプションの網羅調査

**Category**: Solution Evaluation / External Dependencies

以下、GitHub Repository Ruleset で利用可能な全ルールを列挙し、今回の採否と理由を記載する。

### ブランチ/タグ ルール（Branch & Tag Rules）

#### 1. Restrict creations（作成の制限）

| 項目           | 内容                                                                                                                   |
| -------------- | ---------------------------------------------------------------------------------------------------------------------- |
| 概要           | バイパス権限を持つユーザーのみがマッチするブランチ/タグを作成できる                                                    |
| Terraform 属性 | `rules { creation = true }`                                                                                            |
| プラン         | Free（パブリック）で利用可                                                                                             |
| **今回の採否** | **不採用**                                                                                                             |
| 理由           | GitFlow では feature/release/hotfix ブランチを開発者が自由に作成する必要がある。作成を制限すると開発フローが阻害される |

#### 2. Restrict updates（更新の制限）

| 項目           | 内容                                                                                                       |
| -------------- | ---------------------------------------------------------------------------------------------------------- |
| 概要           | バイパス権限を持つユーザーのみがマッチするブランチ/タグにプッシュできる                                    |
| Terraform 属性 | `rules { update = true }`                                                                                  |
| プラン         | Free（パブリック）で利用可                                                                                 |
| **今回の採否** | **不採用**                                                                                                 |
| 理由           | PR 経由での更新を強制する `pull_request` ルールで十分。`update` は PR すら許可しない完全ロックであり、過剰 |

#### 3. Restrict deletions（削除の制限）

| 項目           | 内容                                                                                                                                                                                                                                 |
| -------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| 概要           | バイパス権限を持つユーザーのみがマッチするブランチ/タグを削除できる                                                                                                                                                                  |
| Terraform 属性 | `rules { deletion = true }`                                                                                                                                                                                                          |
| プラン         | Free（パブリック）で利用可                                                                                                                                                                                                           |
| **今回の採否** | **採用（main, develop, release-\*, hotfix-\*）**                                                                                                                                                                                     |
| 理由           | `main` の削除防止は必須（既存設定で有効）。`develop` は恒久ブランチであり削除は運用上あり得ないため保護必須。`release-*` はリリース履歴の追跡に不可欠であり、削除されると調査時に困るため保護必須。`hotfix-*` も同様の理由で保護する |

#### 4. Require linear history（リニアヒストリーの強制）

| 項目           | 内容                                                                                                                                                                                          |
| -------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 概要           | マージコミットを禁止し、squash または rebase のみを許可する                                                                                                                                   |
| Terraform 属性 | `rules { required_linear_history = true }`                                                                                                                                                    |
| プラン         | Free（パブリック）で利用可                                                                                                                                                                    |
| **今回の採否** | **不採用（全ブランチ）。既存の `main` 設定も `false` に変更する**                                                                                                                             |
| 理由           | GitFlow では merge commit を使うことが正統的なフローであり、全ブランチで merge commit に統一する。`main` の既存 `required_linear_history = true` は GitFlow と矛盾するため `false` に変更する |

#### 5. Block force pushes（フォースプッシュの禁止）

| 項目           | 内容                                                                       |
| -------------- | -------------------------------------------------------------------------- |
| 概要           | フォースプッシュを禁止する（デフォルトで有効）                             |
| Terraform 属性 | `rules { non_fast_forward = true }`                                        |
| プラン         | Free（パブリック）で利用可                                                 |
| **今回の採否** | **採用（全保護ブランチ）**                                                 |
| 理由           | フォースプッシュは履歴破壊のリスクがあり、すべての保護ブランチで禁止すべき |

#### 6. Require signed commits（署名付きコミットの強制）

| 項目           | 内容                                                                                                                   |
| -------------- | ---------------------------------------------------------------------------------------------------------------------- |
| 概要           | 署名済み・検証済みのコミットのみプッシュを許可する                                                                     |
| Terraform 属性 | `rules { required_signatures = true }`                                                                                 |
| プラン         | Free（パブリック）で利用可                                                                                             |
| **今回の採否** | **不採用**                                                                                                             |
| 理由           | GPG/SSH 署名の設定は全開発者に追加の運用負荷をかける。セキュリティ要件が高まった段階で導入を検討するが、現時点では過剰 |

#### 7. Require deployments to succeed（デプロイ成功の要求）

| 項目           | 内容                                                                                                                                                   |
| -------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------ |
| 概要           | 指定環境へのデプロイが成功してからマージを許可する                                                                                                     |
| Terraform 属性 | `rules { required_deployments { required_deployment_environments = [...] } }`                                                                          |
| プラン         | Free（パブリック）で利用可                                                                                                                             |
| **今回の採否** | **不採用**                                                                                                                                             |
| 理由           | そもそも現時点で staging/production のデプロイ環境が構築されていないため不可能だが、運用思想がまず異なるため（デプロイだけでよいわけではない）設定不要 |

### PR ルール（Pull Request Rules）

#### 8. Require a pull request before merging（PR 必須）

| 項目           | 内容                             |
| -------------- | -------------------------------- |
| 概要           | マージ前に PR を必須とする       |
| Terraform 属性 | `rules { pull_request { ... } }` |
| プラン         | Free（パブリック）で利用可       |
| **今回の採否** | **採用（全保護ブランチ）**       |

**PR ルールのサブオプション一覧：**

| サブオプション               | Terraform 属性                      | 型            | デフォルト | 今回の設定                              | 理由                                                                   |
| ---------------------------- | ----------------------------------- | ------------- | ---------- | --------------------------------------- | ---------------------------------------------------------------------- |
| 承認必要数                   | `required_approving_review_count`   | Number (0-10) | 0          | `develop`: 0, `release-*`/`hotfix-*`: 1 | develop は個人開発で承認不要、release/hotfix は品質ゲートとして1名承認 |
| プッシュ時の古いレビュー却下 | `dismiss_stale_reviews_on_push`     | Boolean       | false      | `release-*`/`hotfix-*`: true            | 新しいコミット後に古い承認が残ると意味がないため                       |
| 最後のプッシュ者以外の承認   | `require_last_push_approval`        | Boolean       | false      | `release-*`/`hotfix-*`: true            | 自分でプッシュして自分で承認する抜け穴を防ぐ                           |
| コードオーナーの承認         | `require_code_owner_review`         | Boolean       | false      | 不採用                                  | CODEOWNERS ファイル未設定。現時点では不要                              |
| コメント解決の要求           | `required_review_thread_resolution` | Boolean       | false      | `release-*`/`hotfix-*`: true            | レビューコメントの見落とし防止                                         |

#### 9. Require status checks to pass（ステータスチェック必須）

| 項目           | 内容                                                          |
| -------------- | ------------------------------------------------------------- |
| 概要           | CI テスト等のステータスチェックが通過してからマージを許可する |
| Terraform 属性 | `rules { required_status_checks { ... } }`                    |
| プラン         | Free（パブリック）で利用可                                    |
| **今回の採否** | **採用（全保護ブランチ）**                                    |

**サブオプション：**

| サブオプション | Terraform 属性                         | 型      | 今回の設定 | 理由                                                                 |
| -------------- | -------------------------------------- | ------- | ---------- | -------------------------------------------------------------------- |
| Strict モード  | `strict_required_status_checks_policy` | Boolean | true       | ベースブランチの最新変更を取り込んでからチェックを実行。安全性を優先 |
| 必須チェック   | `required_check { context = "ci" }`    | Block   | `ci`       | 既存の CI ワークフローに合わせる                                     |

### メタデータ制限ルール（Metadata Restrictions）

#### 10. Branch name pattern（ブランチ名パターン）

| 項目           | 内容                                                                 |
| -------------- | -------------------------------------------------------------------- |
| 概要           | ブランチ名が指定パターンに一致することを要求する                     |
| Terraform 属性 | `rules { branch_name_pattern { operator = "..." pattern = "..." } }` |
| プラン         | Enterprise 限定（GitHub Free では利用不可）                          |
| **今回の採否** | **不採用**                                                           |
| 理由           | GitHub Free プランでは利用不可。SKILL でのガイドで代替する           |

#### 11. Tag name pattern（タグ名パターン）

| 項目           | 内容                                                              |
| -------------- | ----------------------------------------------------------------- |
| 概要           | タグ名が指定パターンに一致することを要求する                      |
| Terraform 属性 | `rules { tag_name_pattern { operator = "..." pattern = "..." } }` |
| プラン         | Enterprise 限定                                                   |
| **今回の採否** | **不採用**                                                        |
| 理由           | GitHub Free プランでは利用不可                                    |

#### 12. Commit message pattern（コミットメッセージパターン）

| 項目           | 内容                                                                    |
| -------------- | ----------------------------------------------------------------------- |
| 概要           | コミットメッセージが指定パターンに一致することを要求する                |
| Terraform 属性 | `rules { commit_message_pattern { operator = "..." pattern = "..." } }` |
| プラン         | Enterprise 限定                                                         |
| **今回の採否** | **不採用（サーバーサイド）、ローカル hook + SKILL で代替採用**                                                                                                                                                                                                                                            |
| 理由           | GitHub Free プランでは利用不可。代替として `commit-msg` hook でコミットメッセージに Issue 番号（`#<issue-no>`）の含有を検証する。`.githooks/commit-msg` にスクリプトを配置し `core.hooksPath` で有効化する。Claude Code の SKILL でも Issue 番号付与を指示する。サーバーサイド強制は不可だがローカル運用で十分 |

#### 13. Commit author email pattern（コミット作者メールパターン）

| 項目           | 内容                                                                         |
| -------------- | ---------------------------------------------------------------------------- |
| 概要           | コミット作者のメールアドレスが指定パターンに一致することを要求する           |
| Terraform 属性 | `rules { commit_author_email_pattern { operator = "..." pattern = "..." } }` |
| プラン         | Enterprise 限定                                                              |
| **今回の採否** | **不採用**                                                                   |
| 理由           | GitHub Free プランでは利用不可                                               |

#### 14. Committer email pattern（コミッターメールパターン）

| 項目           | 内容                                                                     |
| -------------- | ------------------------------------------------------------------------ |
| 概要           | コミッターのメールアドレスが指定パターンに一致することを要求する         |
| Terraform 属性 | `rules { committer_email_pattern { operator = "..." pattern = "..." } }` |
| プラン         | Enterprise 限定                                                          |
| **今回の採否** | **不採用**                                                               |
| 理由           | GitHub Free プランでは利用不可                                           |

### コード品質ルール

#### 15. Require code scanning results（コードスキャン結果の要求）

| 項目           | 内容                                                                                                                 |
| -------------- | -------------------------------------------------------------------------------------------------------------------- |
| 概要           | コードスキャンの結果が指定の重大度以下であることを要求する                                                           |
| Terraform 属性 | `rules { required_code_scanning { required_code_scanning_tool { ... } } }`                                           |
| プラン         | GitHub Advanced Security が必要（Free プランのパブリックリポジトリでは CodeQL は利用可能）                           |
| **今回の採否** | **採用を試みる。失敗時は CI ステータスチェックにフォールバック**                                                                                                                                                                                                                                                                                                                               |
| 理由           | CodeQL はパブリックリポジトリで無料利用可。Terraform Provider の `required_code_scanning` ブロックにバグ報告あり（Issue #2599）だが、最新バージョンでは修正されている可能性がある。まず `required_code_scanning` ブロックでの設定を試み、失敗した場合は CodeQL を CI ステータスチェック経由でマージ条件にするフォールバック戦略をとる |

#### 16. Require code quality results（コード品質結果の要求）

| 項目           | 内容                                             |
| -------------- | ------------------------------------------------ |
| 概要           | コード品質分析の結果が基準を満たすことを要求する |
| Terraform 属性 | 未確認（比較的新しい機能）                       |
| プラン         | 不明（Enterprise の可能性が高い）                |
| **今回の採否** | **不採用**                                                                                                                                                                    |
| 理由           | Go のリンター（`go vet` / `staticcheck` / `gofmt`）は既に CI に含まれており、`required_status_checks` 経由でカバー済み。この GitHub 機能は別途の Code Quality 分析サービスを前提としており、プラン要件も不明なため現時点では不要 |

### プッシュルール（Push Rules）

#### 17. Restrict file paths（ファイルパスの制限）

| 項目           | 内容                                                              |
| -------------- | ----------------------------------------------------------------- |
| 概要           | 指定パスのファイル変更を含むコミットのプッシュを防止する          |
| Terraform 属性 | Terraform Provider 未対応（Issue #2394）                          |
| プラン         | Team/Enterprise（private/internal リポジトリ）                    |
| **今回の採否** | **不採用**                                                        |
| 理由           | パブリックリポジトリでの利用可否が不明。Terraform Provider 未対応 |

#### 18. Restrict file path length（ファイルパス長の制限）

| 項目           | 内容                                                  |
| -------------- | ----------------------------------------------------- |
| 概要           | ファイルパスの文字数上限を設定する                    |
| Terraform 属性 | Terraform Provider 未対応                             |
| プラン         | Team/Enterprise                                       |
| **今回の採否** | **不採用**                                            |
| 理由           | Terraform Provider 未対応。通常のプロジェクトでは不要 |

#### 19. Restrict file extensions（ファイル拡張子の制限）

| 項目           | 内容                                                   |
| -------------- | ------------------------------------------------------ |
| 概要           | 指定拡張子のファイルを含むコミットのプッシュを防止する |
| Terraform 属性 | Terraform Provider 未対応                              |
| プラン         | Team/Enterprise                                        |
| **今回の採否** | **不採用**                                             |
| 理由           | Terraform Provider 未対応。.gitignore で十分           |

#### 20. Restrict file size（ファイルサイズの制限）

| 項目           | 内容                                                           |
| -------------- | -------------------------------------------------------------- |
| 概要           | 指定サイズを超えるファイルのプッシュを防止する                 |
| Terraform 属性 | Terraform Provider 未対応                                      |
| プラン         | Team/Enterprise                                                |
| **今回の採否** | **不採用**                                                     |
| 理由           | Terraform Provider 未対応。必要に応じて GitHub UI から設定可能 |

### その他

#### 21. Merge queue（マージキュー）

| 項目           | 内容                                                                                  |
| -------------- | ------------------------------------------------------------------------------------- |
| 概要           | PR のマージをキューで管理し、順序制御や自動チェックを行う                             |
| Terraform 属性 | 既知のバグあり（Issue #2192, #2339 — apply 時にマージキューが無効化される）           |
| プラン         | パブリックリポジトリ（organization 所有）で利用可                                     |
| **今回の採否** | **不採用**                                                                            |
| 理由           | 個人プロジェクトでマージキューは不要。Terraform Provider にバグがあり、安定していない |

---

## S2: 既存 `main` ブランチ保護ルールの `required_linear_history` と GitFlow の矛盾

**Category**: Prior Decisions / Risk & Failure Modes

**Finding**: 既存の `main` ブランチ保護ルールでは `required_linear_history = true` が設定されている。しかし GitFlow では `release-*` や `hotfix-*` から `main` へ merge commit でマージするのが標準フローであり、リニアヒストリーと矛盾する。

**Recommendation**: `main` ブランチの `required_linear_history` を `false` に変更する。 **決定済み：merge commit に統一。spec.md FR-009 に反映済み。**

**Evidence**: 既存 `main.tf` の L23 で `required_linear_history = true` が設定されている。GitFlow の公式ドキュメントでは merge commit を使用する前提。Survey 確認時にユーザーが merge commit 統一を決定。**解決済み。**

---

## S3: `develop` ブランチの事前作成

**Category**: Edge Cases and Failure Handling

**Finding**: Terraform で `develop` ブランチのルールセットを作成しても、ブランチ自体が存在しない場合はルールセットは作成できるが実効性がない。

**Recommendation**: ~~Terraform apply の前に `develop` ブランチを `main` から作成する手順を計画に含める。~~ **解決済み：`develop` ブランチは作成済み。**

**Evidence**: GitHub Ruleset はブランチの存在に関係なくパターンマッチで動作するため、ルールセット自体は作成可能。`develop` ブランチはユーザーが手動で作成済み。

---

## S4: Terraform Provider の制約と注意点

**Category**: External Dependencies / Feasibility Verification

**Finding**: Terraform GitHub Provider（v6.x）にはいくつかの既知の問題がある：

1. **bypass_actors の削除が反映されない**（Issue #2269, #2179, #2952）— bypass_actors を削除しても GitHub 側に反映されないバグ
2. **merge queue の競合**（Issue #2339）— ruleset を apply するとマージキューが無効化される
3. **required_code_scanning の不具合**（Issue #2599）— 作成時にエラーになるケースがある
4. **Push Ruleset 未対応**（Issue #2394）— ファイルパス/サイズ/拡張子制限は未実装

**Recommendation**: 今回使用するルール（pull_request, required_status_checks, non_fast_forward, deletion）は安定しており、既知の問題の影響を受けない。bypass_actors の変更を行う場合は手動確認を推奨。

**Evidence**: GitHub Issues での報告を確認。既存の `main` ブランチ保護ルールが正常に動作していることから、基本的なルールは安定していると判断。

---

## S5: PR マージ戦略の GitHub 設定との整合性

**Category**: Integration Impact

**Finding**: Clarify で「Merge commit のみ許可」と決定したが、GitHub の Repository Ruleset の `pull_request` ブロックには直接的なマージタイプ制限のオプションがない（GitHub UI の「Allow merge commits / Allow squash merging / Allow rebase merging」はリポジトリ設定であり、Ruleset ではない）。

**Recommendation**: マージタイプの制限は Terraform の `github_repository` リソース（Ruleset ではない）で `allow_merge_commit = true`, `allow_squash_merge = false`, `allow_rebase_merge = false` として設定する。または GitHub UI のリポジトリ設定で手動設定する。Ruleset だけでは実現できない点に注意。

**Evidence**: Terraform Provider ドキュメントおよび GitHub ドキュメントの確認。`pull_request` ルール内にマージタイプ制限は存在しない。

### S5 追加調査: `github_repository` リソースの全設定オプション

Ruleset とは別に、`github_repository` リソースで管理すべきリポジトリ設定を網羅調査した。

#### マージ関連設定

| 設定 | Terraform 属性 | 型 | デフォルト | 今回の採否 | 理由 |
|---|---|---|---|---|---|
| マージコミット許可 | `allow_merge_commit` | Boolean | true | **採用（true）** | GitFlow で merge commit を使用するため |
| Squash マージ許可 | `allow_squash_merge` | Boolean | false | **採用（false）** | Merge commit に統一するため禁止 |
| Rebase マージ許可 | `allow_rebase_merge` | Boolean | false | **採用（false）** | Merge commit に統一するため禁止 |
| 自動マージ許可 | `allow_auto_merge` | Boolean | false | **不採用（false）** | 個人プロジェクトでは不要。承認後に手動でマージすれば十分 |
| マージ後にブランチ自動削除 | `delete_branch_on_merge` | Boolean | false | **採用（true）** | 削除保護ルール（`deletion = true`）が優先されるため、保護ブランチ（`develop`, `release-*`, `hotfix-*`）は自動削除されない。feature ブランチのみ自動削除され、ブランチの手動クリーンアップが不要になる |
| マージコミットタイトル | `merge_commit_title` | String | `"MERGE_MESSAGE"` | **不採用（デフォルト）** | デフォルトの挙動で問題なし |
| マージコミットメッセージ | `merge_commit_message` | String | `"PR_TITLE"` | **不採用（デフォルト）** | デフォルトの挙動で問題なし |
| Squash コミットタイトル | `squash_merge_commit_title` | String | — | **不採用** | Squash マージを禁止するため設定不要 |
| Squash コミットメッセージ | `squash_merge_commit_message` | String | — | **不採用** | Squash マージを禁止するため設定不要 |

#### セキュリティ関連設定

| 設定 | Terraform 属性 | 型 | デフォルト | 今回の採否 | 理由 |
|---|---|---|---|---|---|
| 脆弱性アラート | `vulnerability_alerts` | Boolean | false | **採用（true）** | 依存パッケージの脆弱性を検知する基本的なセキュリティ機能。無料で利用可能 |
| Web コミット署名要求 | `web_commit_signoff_required` | Boolean | false | **不採用（false）** | GitHub UI からのコミット時に署名を要求する機能。現時点では運用負荷に見合わない |
| Advanced Security | `security_and_analysis.advanced_security` | String | — | **採用（enabled）** | パブリックリポジトリでは無料。CodeQL 等の前提条件 |
| Secret Scanning | `security_and_analysis.secret_scanning` | String | — | **採用（enabled）** | シークレット（API キー等）の誤コミットを検出。パブリックリポジトリでは無料かつ必須級 |
| Secret Scanning Push Protection | `security_and_analysis.secret_scanning_push_protection` | String | — | **採用（enabled）** | シークレットを含むコミットのプッシュをブロック。パブリックリポジトリでは無料かつ必須級 |
| Secret Scanning AI Detection | `security_and_analysis.secret_scanning_ai_detection` | String | — | **不採用** | 比較的新しい機能でパブリックリポジトリでの利用可否が不明確。通常の Secret Scanning で十分 |

#### その他の設定

| 設定 | Terraform 属性 | 型 | 今回の採否 | 理由 |
|---|---|---|---|---|
| リポジトリ説明 | `description` | String | **スコープ外** | 既に設定済みのため今回は対象外 |
| Issues 有効化 | `has_issues` | Boolean | **スコープ外** | 既に有効 |
| Wiki 有効化 | `has_wiki` | Boolean | **スコープ外** | 既存設定を維持 |
| Projects 有効化 | `has_projects` | Boolean | **スコープ外** | 既存設定を維持 |
| Discussions 有効化 | `has_discussions` | Boolean | **スコープ外** | 既存設定を維持 |
| テンプレートリポジトリ | `is_template` | Boolean | **不採用** | テンプレートリポジトリではない |
| アーカイブ | `archived` | Boolean | **不採用** | アクティブなリポジトリ |
| Topics | `topics` | List | **スコープ外** | 既存設定を維持 |

---

## S6: パブリックリポジトリ + GitHub Free プランの制約

**Category**: Constraints and Tradeoffs

**Finding**: `psm` はパブリックリポジトリで GitHub Free プランを使用している。以下の制約がある：

- **利用可能**: Branch/Tag Ruleset の基本ルール（PR必須、ステータスチェック、フォースプッシュ禁止、削除制限等）
- **利用不可**: メタデータ制限（ブランチ名/コミットメッセージ/メールパターン）— Enterprise 限定
- **利用不可**: Push Ruleset（ファイルパス/サイズ/拡張子制限）— Team/Enterprise（private/internal のみ）
- **制限**: Organization Ruleset は Team 以上

**Recommendation**: 今回使用するルールはすべて Free プランで利用可能。メタデータ制限（ブランチ命名規則の強制等）は SKILL とローカルフック（pre-commit hook）で代替する。

**Evidence**: GitHub 公式ドキュメント「GitHub's plans」および Rulesets ドキュメントの確認。

---

## Items Requiring PoC

- **S2 の検証**: `main` の `required_linear_history` を `false` に変更した場合に、既存の保護ルール全体が意図通り動作するかの確認（`terraform plan` で確認可能）

## Constitution Impact

今回の変更に憲法改正は不要。ただし以下の注意点：

- Constitution の「Development Workflow」セクションに「Keep the main branch always in a buildable, passing-tests state」とあり、GitFlow のブランチ保護ルールはこれを強化する方向で整合的
- Test-First 原則は CI ステータスチェック必須と整合的

## Recommendation

**S2（`main` の `required_linear_history` と GitFlow の矛盾）を先に解決してから `/speckit.plan` に進むことを推奨。**

この決定は Terraform コードの設計に直接影響するため、plan フェーズ前に確定させるべき。選択肢：

1. `main` の `required_linear_history = false` に変更する（GitFlow 正統派に合わせる）
2. `main` への merge は squash merge を許可する（リニアヒストリーを維持するが GitFlow からの逸脱）
