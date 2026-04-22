# Quickstart: テスト実行基盤の整備

## 結合テストをローカルで実行する

```bash
# psmをビルド
go build -o psm .

# tests/integration/ で実行
PSM_BIN=./psm PSM_TEST_PROFILE=psm-sandbox bash tests/integration/test.sh
```

## 結合テストのCI実行

developまたはrelease-*ブランチにpushすると `.github/workflows/integration-test.yml` が自動実行される。

## E2Eバイナリの取得

リリースブランチへのpush時にCIがバイナリをビルドし、GitHub Actions artifactとしてアップロードする。Actions画面からダウンロード可能。

## AWS OIDC設定（前提条件）

テスト用AWSアカウントに以下が設定されていること:
1. GitHub OIDC プロバイダ (`token.actions.githubusercontent.com`)
2. IAMロール（信頼ポリシーでdevelop + release-*ブランチに制限）
3. IAMロールにSSM read/write/delete権限（テスト用パス配下のみ）
