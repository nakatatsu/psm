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

AWS SSM Parameter Storeとの実通信を伴う操作。

| テスト                     | 検証内容                                  |
| -------------------------- | ----------------------------------------- |
| `TestSSMStoreGetAll`       | パラメータ一括取得                        |
| `TestSSMStorePutAndDelete` | パラメータ作成・削除のライフサイクル      |
| `TestSSMDryRun`            | dry-run時にAWS側が変更されないこと        |
| `TestSSMSyncExecute`       | sync実行後のAWS状態（create/update/skip） |
| `TestExportRoundTrip`      | export→再parseで差分ゼロ（SC-006）        |

### 実行方法

```bash
PSM_INTEGRATION_TEST=1 go test -race -v ./...
```

`PSM_TEST_PROFILE` で任意のAWSプロファイルを指定可能。

### 実行条件

- **リリース前**: 手動で実行する
- CI���は `PSM_INTEGRATION_TEST` が未設定のため自動スキップされる

### 方針

- テスト用パラメータは `/psm-test/` プレフィックス配下に限定し、本番データと分離する
- 各テストは `cleanAllSSMTestParams` で前処理し、`defer cleanupSSMTestData` で後処理する。テスト間の状態汚染を防止する
- 強整合性が必要な検証には `GetParameter`（強整合性保証）を使用し、`GetParametersByPath`（結果整合性）に依存しない

### 前提条件

- AWSクレデンシャルが設定済みであること
- 対象アカウントのSSMに書き込み・削除権限があること
- `/psm-test/` 配下を自由に使える��と

---

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

---

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

---

## 5. 品質ゲートまとめ

[テスト原則](../shared/standards/test-standard.md)の品質ゲート定義に対する、本プロジェクトでの実施内容。

### マージ前

| 項目           | 実施内容                               | 自動/手動  |
| -------------- | -------------------------------------- | ---------- |
| 自動チェック   | gofumpt, golangci-lint, go build       | 自動（CI） |
| ユニットテスト | `go test -race ./...`                  | 自動（CI） |
| コードレビュー | PRレビュー（main向けは1 approval必須） | 手動       |

### リリース前

| 項目                 | 実���内容                                                           | 自動/手��� |
| -------------------- | ------------------------------------------------------------------- | ---------- |
| 結合テスト           | `PSM_INTEGRATION_TEST=1 go test -race -v ./...`                     | 手動       |
| セキュリティスキャン | govulncheck, CodeQL                                                 | 自動（CI） |
| パフォー��ンステスト | 現時点では対象外（CLIツールのため応答時間が問題になる規模ではない） | —          |
| スモークテスト       | E2Eの主���シナリオ（sync dry-run, export）                          | 手動       |

### リリース後

| 項目           | 実施内容                                              | 自動/手動 |
| -------------- | ----------------------------------------------------- | --------- |
| ヘルスチェック | `psm --version`、dry-run実行                          | 手動      |
| メトリクス検証 | 現時点では対象外（CLIツールのためメトリクス基盤なし） | —         |
