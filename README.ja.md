# psm

> **ステータス: 開発中**

YAML ファイルを AWS SSM Parameter Store / Secrets Manager に同期するツール。

[English](README.md)

## インストール

[GitHub Releases](https://github.com/nakatatsu/psm/releases) からビルド済みバイナリをダウンロードしてください。

## 使い方

### Sync: YAML を AWS に反映する

```bash
psm sync --store <ssm|sm> [--profile <name>] [--dry-run] [--prune] <file>
```

| フラグ | 必須 | 説明 |
|--------|------|------|
| `--store <ssm\|sm>` | はい | 対象ストア: `ssm` (Parameter Store) または `sm` (Secrets Manager) |
| `--profile <name>` | いいえ | AWS プロファイル名（デフォルト: SDK デフォルト認証） |
| `--dry-run` | いいえ | 変更内容を表示するのみで実行しない |
| `--prune` | いいえ | YAML に存在しないキーを AWS から削除する |

例:

```bash
# 変更内容をプレビュー
psm sync --store ssm --dry-run secrets.yaml

# 変更を適用
psm sync --store ssm secrets.yaml

# AWS プロファイルを指定して適用
psm sync --store ssm --profile myprofile secrets.yaml

# YAML にないキーを削除しつつ同期
psm sync --store ssm --prune secrets.yaml
```

### Export: AWS のパラメータを YAML に書き出す

```bash
psm export --store <ssm|sm> [--profile <name>] <file>
```

例:

```bash
psm export --store ssm output.yaml
```

### SOPS との組み合わせ（暗号化されたシークレット）

SOPS で復号し、パイプで psm に渡します:

```bash
sops -d secrets.enc.yaml | psm sync --store ssm --dry-run /dev/stdin
sops -d secrets.enc.yaml | psm sync --store ssm /dev/stdin
```

鍵の生成や暗号化設定を含む詳しい手順は [example/README.md](example/README.md) を参照してください。

### YAML フォーマット

キーがそのまま AWS のパラメータ名/シークレット名になります。値はスカラー値（string, int, bool, float）のみ使用できます。

```yaml
/myapp/database/host: localhost
/myapp/database/port: "5432"
/myapp/database/password: my-secret-password
/myapp/api/key: my-api-key
```

> `sops` メタデータキーは同期時に自動的に除外されます。

## 開発

### DevContainer からの AWS アクセス

```
aws sso login --sso-session psm-sandbox --use-device-code
```
