# Research: テスト実行基盤の整備

## R1: GitHub Actions OIDC + AWS設定

**Decision**: `aws-actions/configure-aws-credentials` アクションでOIDC認証を行う
**Rationale**: GitHub公式推奨の方法。IAMアクセスキー不要でセキュア。公開リポジトリでも利用可能。
**Alternatives considered**:
- IAMアクセスキーをGitHub Secretsに保管 → ローテーション管理が煩雑、セキュリティ上劣る。spec/constraintsで禁止済み
- AWS STS AssumeRole with external ID → OIDCより複雑で利点なし

### IAM信頼ポリシー

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "arn:aws:iam::<ACCOUNT_ID>:oidc-provider/token.actions.githubusercontent.com"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "token.actions.githubusercontent.com:aud": "sts.amazonaws.com"
        },
        "StringLike": {
          "token.actions.githubusercontent.com:sub": [
            "repo:nakatatsu/psm:ref:refs/heads/develop",
            "repo:nakatatsu/psm:ref:refs/heads/release-*"
          ]
        }
      }
    }
  ]
}
```

`StringLike` で `release-*` のワイルドカードが使える。

## R2: ワークフロートリガー設計

**Decision**: `push` to develop/release-* でトリガー
**Rationale**: featureブランチからのPRではなく、マージ後のpushで実行する。FR-002の「マージされた時に」に合致。
**Alternatives considered**:
- `pull_request` イベント → マージ前に実行できるがOIDCのsub claimがPR用になりブランチ制限と合わない
- `workflow_dispatch` → 手動トリガーでは自動化にならない

### ワークフロー構成

```yaml
on:
  push:
    branches:
      - develop
      - 'release-*'
```

## R3: test.shの移設方針

**Decision**: example/test.shをtests/integration/にコピーし、パス参照を調整する
**Rationale**: example/は触らない制約（FR-001）。test.shはSCRIPT_DIR相対でファイル参照しているため、テストデータも合わせてコピーすれば動作する。
**Alternatives considered**:
- シンボリックリンク → CIのcheckoutで壊れる可能性
- test.shからexample/を参照 → FR-002「example/に依存しない」に違反（注: FR-002の文言は変更されたが、テスト資材の独立性は維持すべき）

### 必要なファイル

- `tests/integration/test.sh` — example/test.shのコピー、PSM_BIN等のパスを調整
- `tests/integration/secrets.example.yaml` — example/secrets.example.yamlのコピー

## R4: E2Eバイナリのビルドと命名

**Decision**: リリースブランチへのpush時にCIで `go build` してGitHub Actions artifactとしてアップロード
**Rationale**: GoReleaserはリリース用（mainマージ後のタグ作成時）。E2Eテスト用にはシンプルな `go build` で十分。
**Alternatives considered**:
- GoReleaserでビルド → リリースブランチ段階ではタグがないため不適
- 手動ビルド → FR-006で自動化が要件

### 命名規則

`psm_<git-sha-short>_linux_amd64` 形式。バージョンタグはリリースブランチ段階では未確定のため、git短縮ハッシュを使用。

## R5: FR-005 マージブロックの実現方法

**Decision**: GitHub branch protection rulesで `integration-test` ワークフローをrequired status checkに追加
**Rationale**: 既存のCIワークフローと同じ仕組み。追加の実装不要。
**Alternatives considered**: 特になし。これが標準的な方法。

注意: この設定はリポジトリの Settings > Branches で手動設定が必要（コードでは実現できない）。
