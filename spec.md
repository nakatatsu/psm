# spec.md - parasync

## 概要

SOPS復号済みYAMLファイルを読み、AWS SSM Parameter Store または Secrets Manager に同期するCLIツール。

## CLI

```
parasync [flags] <file>
```

### フラグ

| フラグ | 必須 | 説明 |
|--------|------|------|
| `--store` | ○ | `ssm` または `sm` |
| `--app` | ○ | アプリ名 |
| `--env` | ○ | 環境名 |
| `--prune` | - | ファイルにないキーを削除 |
| `--dry-run` | - | 実行せず差分表示 |

`--region` は AWS_REGION / AWS_DEFAULT_REGION 環境変数に従う（AWS SDK デフォルト挙動）。

## 入力ファイル

復号済みYAML。フラットな key: value。

```yaml
DATABASE_URL: "postgres://localhost:5432/mydb"
REDIS_HOST: "redis.example.com"
```

`sops` キーが残っていれば無視する。

## AWS マッピング

### SSM Parameter Store (`--store ssm`)

- パス: `/{app}/{env}/{key}` (例: `/myapp/prod/DATABASE_URL`)
- 型: SecureString
- PutParameter with Overwrite=true (upsert)

### Secrets Manager (`--store sm`)

- 名前: `{app}/{env}/{key}` (例: `myapp/prod/API_KEY`)
- CreateSecret → ResourceExistsException → PutSecretValue

## 同期ルール

- デフォルト: upsert のみ（追加・更新）
- `--prune`: `/{app}/{env}/` 配下でファイルにないキーを削除
  - SSM: DeleteParameter
  - SM: DeleteSecret (ForceDeleteWithoutRecovery=true)
- `--dry-run`: create / update / delete / no-change の予定を表示、実行しない

## エラーハンドリング

- 1件失敗しても残りは続行
- 終了コード: 0=全成功, 1=一部失敗

## 依存ライブラリ

- AWS SDK for Go v2
- gopkg.in/yaml.v3
