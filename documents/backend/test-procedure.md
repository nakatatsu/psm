# テスト手順書

本文書は[テスト原則](../shared/standards/test-standard.md)および[開発憲章](../shared/policies/constitution.md)に準拠した、psmプロジェクト固有のテスト手順を定める。

## 1. ユニットテスト

### 対象

ビジネスロジックを持つ全関数に`go test`で自動テストをかける。

### 実行方法

```bash
go test -race ./...
```

### 実行条件

- **開発中**: 対象関数を変更した都度
- **マージ前（CI）**: PR作成時に自動実行（`.github/workflows/ci.yml`）

### 方針

- `Store` インターフェースに対する `fakeStore` を用いてAWS依存を排除する。モックは外部依存（AWS SDK）に限定する（テスト原則 Should準拠）
- 入力と期待値のバリエーションが主体のテストにはテーブル駆動テストを用いる
- テスト間で状態を共有しない（テスト原則 MUST準拠）
- ダミー文字列をValueに設定し、それがstdout/stderrに出力されないことを検証する（実際のシークレットをテストに含めてはならない）

## 2. 結合テスト

### 対象

コンポーネント間の契約をユースケースの1操作単位で検証する。検証対象のユースケースは[要件定義](../../specs/requirements.md)を参照すること。

### 実行方法

CI（GitHub Actions）により自動実行される。developブランチおよびリリースブランチへのpush時に `.github/workflows/integration-test.yml` がトリガーされる。

ローカルで手動実行する場合:

```bash
go build -o psm .
PSM_BIN=./psm PSM_TEST_PROFILE=psm-sandbox bash tests/integration/test.sh
```

### 実行条件

- **CI自動**: developブランチおよびリリースブランチへのpush時
- **手動**: 必要に応じてローカルで実行可能

### 前提条件

- CI: AWSアカウントにOIDCプロバイダとIAMロールが設定されていること
- CI: GitHubリポジトリにAWS_ROLE_ARNシークレットとAWS_REGION変数が設定されていること
- ローカル: AWSクレデンシャルが設定済みであること（`aws sso login`）
- ローカル: 対象アカウントのSSMに書き込み・削除権限があること

## 3. E2Eテスト

### 対象

CLIとしてのエンドツーエンドのユーザーシナリオ。

| シナリオ          | 手順                                          | 期待結果                                 |
| ----------------- | --------------------------------------------- | ---------------------------------------- |
| 初回sync          | SOPSで暗号化したYAMLを `psm sync` で適用      | パラメータがSSMに作成される              |
| 差分sync          | 値を変更して再度 `psm sync`                   | 変更分のみupdateされる                   |
| dry-run           | `psm sync --dry-run`                          | 計画が表示され、AWS側は変更されない      |
| 削除フロー        | `--delete` オプションで不要パラメータを削除   | 正規表現に一致するパラメータが削除される |
| コンフリクト検出  | sync YAMLに含まれるキーが削除パターンにも一致 | エラーで中断し、一切の変更が行われない   |
| export            | `psm export out.yaml`                         | SSM上の全パラメータがYAML出力される      |
| export→sync冪等性 | exportした結果をそのままsyncする              | 差分ゼロ                                 |
| 承認フロー拒否    | sync実行時に `N` を入力                       | 変更されずに終了する                     |
| バージョン表示    | `psm --version`                               | バージョン文字列が出力される             |

### 実行方法

手動実行。ビルド済みバイナリを使用する。

```bash
go build -o psm .
./psm sync --dry-run params.yaml        # dry-run確認
./psm sync params.yaml                  # 実行
./psm export out.yaml                   # export確認
```

### 実行条件

- **リリース前**: 主要シナリオ（初回sync、差分sync、dry-run、export）を最低限実施する
- **リリース後**: ヘルスチェックとして `psm --version` およびdry-runを実施する

### 方針

- 実際のSOPS暗号化ファイルを使用し、復号→sync→検証の一連のフローを確認する
- 本番環境ではなく、テスト用AWSアカウントまたは `/psm-test/` プレフィックスで実施する

## 4. 静的解析・セキュリティスキャン

品質ゲート「マージ前」の自動チェックに該当する。

| チェック       | ツール                                                    | 実行タイミング                      |
| -------------- | --------------------------------------------------------- | ----------------------------------- |
| ビルド         | `go build ./...`                                          | CI（全PR・push）                    |
| フォーマット   | `gofumpt`                                                 | CI（全PR・push）                    |
| リント         | `golangci-lint`（go vet, staticcheck, gosec, errcheck等） | CI（全PR・push）                    |
| 脆弱性スキャン | `govulncheck`                                             | CI（全PR・push、continue-on-error） |
| コード解析     | CodeQL                                                    | CI（GitHub Advanced Security）      |

### 実行方法（ローカル）

```bash
gofumpt -l .
golangci-lint run ./...
govulncheck ./...
```

## 5. パフォーマンステスト

現時点では対象外とする。psmはCLIツールであり、処理対象のパラメータ数も限定的なため、応答時間が問題になる規模ではない。

将来パラメータ数が大幅に増加した場合は、大量パラメータでのsync所要時間を計測するベンチマークテストの導入を検討する。

## 6. スモークテスト

### 対象

リリース候補のバイナリが主要機能を正常に実行できることの簡易確認。

### 実行方法

```bash
go build -o psm .
./psm --version
./psm sync --dry-run params.yaml
./psm export /tmp/smoke-test-export.yaml && rm /tmp/smoke-test-export.yaml
```

### 実行条件

- **リリース前**: 必ず実施する

### 判定基準

- 各コマンドがエラーなく終了すること（exit code 0）

## 7. ヘルスチェック

CLIツールのため対象外。常駐プロセスではないためヘルスチェックの概念が該当しない。

## 8. メトリクス検証

CLIツールのため対象外。常駐プロセスではないためメトリクス基盤を持たない。

## ブランチとテスト種別の対応

| ブランチ   | ユニットテスト   | 静的解析         | 結合テスト       | E2Eテスト                              |
| ---------- | ---------------- | ---------------- | ---------------- | -------------------------------------- |
| feature/\* | CI自動（PR時）   | CI自動（PR時）   | —                | —                                      |
| develop    | CI自動（push時） | CI自動（push時） | CI自動（push時） | —                                      |
| release-\* | CI自動（push時） | CI自動（push時） | CI自動（push時） | リリースブランチのビルドバイナリで実施 |
