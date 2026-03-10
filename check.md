# Mechanical Checks

本プロジェクトで使用する機械的チェックの一覧と実行方法。

## 一覧

| カテゴリ | ツール | 目的 | コマンド |
|---------|--------|------|---------|
| コンパイル | `go build` | ビルド可否 | `go build ./...` |
| フォーマッタ | `gofmt` | Go 標準フォーマット | `gofmt -l .` |
| フォーマッタ | `gofumpt` | gofmt の厳格版（空行ルール等） | `gofumpt -l .` |
| フォーマッタ | `goimports` | import 順序・不要 import 検出 | `goimports -l .` |
| 静的解析 | `go vet` | Go 標準静的解析 | `go vet ./...` |
| 静的解析 | `staticcheck` | go vet 補完（非推奨 API、未使用コード等） | `staticcheck ./...` |
| リンタ | `golangci-lint` | 統合リンタ（errcheck, gosec 等を含む） | `golangci-lint run ./...` |
| セキュリティ | `gosec` | セキュリティ脆弱性検出 | `golangci-lint` 経由で実行 |
| 脆弱性 | `govulncheck` | 依存ライブラリの既知脆弱性検出 | `govulncheck ./...` |
| テスト | `go test` | ユニットテスト | `go test ./...` |
| テスト | `go test` (統合) | AWS 統合テスト | `PSM_INTEGRATION_TEST=1 PSM_TEST_PROFILE=psm go test -v -timeout 120s ./...` |

## 一括実行

```bash
# フォーマッタ + 静的解析 + リンタ + ユニットテスト
gofumpt -l . && goimports -l . && go vet ./... && staticcheck ./... && golangci-lint run ./... && go test ./...

# 統合テスト（AWS SSO 認証済みの状態で）
PSM_INTEGRATION_TEST=1 PSM_TEST_PROFILE=psm go test -v -timeout 120s ./...

# 脆弱性チェック（vuln.go.dev へのネットワーク接続が必要）
govulncheck ./...
```

## 設定ファイル

| ファイル | 対象 | 内容 |
|---------|------|------|
| `.golangci.yml` | golangci-lint | gosec の G703/G705 を除外（CLI ツールにつき該当なし） |

## 除外ルールの根拠

| ルール | 説明 | 除外理由 |
|--------|------|---------|
| G703 (Path traversal) | ユーザー入力パスによるファイル操作 | CLI ツールでは引数のファイルパスを読み書きするのが本来の機能。実行者＝操作者 |
| G705 (XSS) | 外部データの出力 | 出力先はターミナルの stdout/stderr であり、ブラウザではない |

## ツール取得方法

| ツール | 取得方法 | 備考 |
|--------|---------|------|
| `gofumpt` | `go install mvdan.cc/gofumpt@latest` | |
| `goimports` | `go install golang.org/x/tools/cmd/goimports@latest` | |
| `staticcheck` | GitHub Releases バイナリ | `go install` は proxy.golang.org 依存のため、ファイアウォール環境では使えない |
| `golangci-lint` | `go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest` | gosec を内包 |
| `govulncheck` | `go install golang.org/x/vuln/cmd/govulncheck@latest` | vuln.go.dev への接続が必要 |

## 環境制約

- **govulncheck**: `vuln.go.dev` への外部接続が必要。DevContainer のファイアウォール環境では実行不可
- **staticcheck**: `proxy.golang.org` 経由の `go install` が使えないため、GitHub Releases からバイナリ取得（Dockerfile に組み込み済み）
