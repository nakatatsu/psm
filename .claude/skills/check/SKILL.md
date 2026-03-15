---
name: check
description: >
  Go プロジェクトの機械的チェック（フォーマッタ、静的解析、リンタ、テスト）を実行するスキル。
  /check で呼び出す。引数でスコープを制御できる。
  コードを書いた後、コミット前、push 前など品質確認が必要な場面で使う。
  ユーザーが「チェックして」「lint して」「テスト回して」「ビルド通るか確認」と言った場合もこのスキルを使う。
---

# check — Go 機械的チェック

プロジェクトの Go コードに対してフォーマッタ・静的解析・リンタ・テストを実行する。

## 使い方

| 呼び出し | 実行内容 |
|---------|---------|
| `/check` | フォーマッタ + 静的解析 + リンタ + ユニットテスト |
| `/check all` | 上記 + govulncheck |
| `/check int` | AWS 統合テストのみ |
| `/check vuln` | govulncheck のみ |

## 実行手順

すべてのコマンドは `/workspace` ディレクトリで実行する。

### デフォルト (`/check`)

以下を順番に実行する。途中で失敗しても最後まで実行し、結果をまとめて報告する。

```bash
cd /workspace

# 1. フォーマッタ
gofumpt -l .
goimports -l .

# 2. 静的解析
go vet ./...
staticcheck ./...

# 3. リンタ（gosec を内包）
golangci-lint run ./...

# 4. ユニットテスト
go test ./...
```

各ステップの結果を以下の形式で報告する：

```
## チェック結果

| チェック | 結果 |
|---------|------|
| gofumpt | ✅ OK / ❌ 要修正 |
| goimports | ✅ OK / ❌ 要修正 |
| go vet | ✅ OK / ❌ 要修正 |
| staticcheck | ✅ OK / ❌ 要修正 |
| golangci-lint | ✅ OK / ❌ 要修正 |
| go test | ✅ OK / ❌ 失敗 |
```

失敗があった場合は、各項目の詳細（どのファイル・行で何が問題か）を続けて表示する。

### `/check all`

デフォルトのチェックをすべて実行した後、追加で govulncheck を実行する：

```bash
govulncheck ./...
```

govulncheck は `vuln.go.dev` への外部接続が必要。Squid プロキシ経由でアクセス可能。

### `/check int`

AWS 統合テストを実行する：

```bash
cd /workspace
PSM_INTEGRATION_TEST=1 PSM_TEST_PROFILE=psm go test -v -timeout 120s ./...
```

AWS SSO 認証が必要。認証エラーが出た場合は `aws sso login` の実行を促す。

### `/check vuln`

脆弱性チェックのみ実行する：

```bash
cd /workspace
govulncheck ./...
```

## 環境に関する注意

- **golangci-lint** は gosec を内包。`.golangci.yml` で G703（パストラバーサル）と G705（XSS）を除外済み。CLI ツールなので該当しないため。
- **staticcheck** は Dockerfile では GitHub Releases バイナリからインストール済み。Squid プロキシ経由で `go install` も利用可能。
- **govulncheck** は `vuln.go.dev` への外部接続が必要。Squid プロキシ（outbound-filter）経由でアクセス可能。
