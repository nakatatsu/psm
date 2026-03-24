# psm

> **ステータス: 開発中**

YAML ファイルを AWS SSM Parameter Store に同期するツール。

> 本書は [README.md](README.md) の日本語訳です。情報が古い可能性があるため、最新の情報は英語版を参照してください。

## インストール

[GitHub Releases](https://github.com/nakatatsu/psm/releases) からビルド済みバイナリをダウンロードしてください。

## 使い方

### Sync: YAML を AWS に反映する

```bash
psm sync --store ssm [--profile <name>] [--dry-run] [--skip-approve] [--debug] [--delete <file>] <sync-file>
```

| フラグ | 必須 | 説明 |
|--------|------|------|
| `--store ssm` | はい | 対象ストア: `ssm` (Parameter Store) |
| `--profile <name>` | いいえ | AWS プロファイル名（デフォルト: SDK デフォルト認証） |
| `--dry-run` | いいえ | 変更内容を表示するのみで実行しない（プロンプトなし） |
| `--skip-approve` | いいえ | 承認プロンプトをスキップして即実行（CI/自動化向け） |
| `--debug` | いいえ | デバッグレベルのログを有効化 |
| `--delete <file>` | いいえ | 削除対象キーの正規表現パターンを記載した YAML ファイル |

デフォルトでは `psm sync` は実行計画を表示し、確認を求めます。`y` で実行、Enter でキャンセルです。

```bash
# 変更内容をプレビュー（プロンプトなし、実行なし）
psm sync --store ssm --dry-run secrets.yaml

# 変更を適用（計画表示後、承認を求める）
psm sync --store ssm secrets.yaml

# AWS プロファイルを指定して適用
psm sync --store ssm --profile myprofile secrets.yaml

# 承認をスキップ（CI/CD パイプライン向け）
psm sync --store ssm --skip-approve secrets.yaml
```

> stdin がターミナルでない場合（パイプ入力など）かつ `--skip-approve` が未指定の場合、自動的に実行を中止します。

### 不要キーの削除

不要な AWS キーを削除するには、正規表現パターンを記載した YAML ファイルを作成し `--delete` で指定します:

```yaml
# needless.yml — 削除対象キーのパターン
- "^/myapp/legacy/"
- "^/myapp/deprecated-.*"
```

```bash
psm sync --store ssm --delete needless.yml secrets.yaml
```

パターンに一致し**かつ同期 YAML に存在しない**キーのみが削除されます。

**安全機能:**
- **コンフリクト検出**: 削除パターンに一致するキーが同期 YAML にも存在する場合、変更を一切行わず中止します。
- **未管理キー警告**: 同期 YAML にもパターンにも該当しないキーは警告として表示されます。
- **承認必須**: 削除は作成・更新と同様に承認が必要です。

> **`--prune` からの移行**: `--prune` フラグは廃止されました。代わりに `--delete <file>` で正規表現パターンを指定してください。旧 `--prune` と同等の動作には `- ".*"` を含むファイルを使用します。

### Export: AWS のパラメータを YAML に書き出す

```bash
psm export --store ssm [--profile <name>] [--debug] <file>
```

```bash
psm export --store ssm output.yaml
```

### SOPS との組み合わせ（暗号化されたシークレット）

SOPS で復号し、パイプで psm に渡します:

```bash
sops -d secrets.enc.yaml | psm sync --store ssm --dry-run /dev/stdin
sops -d secrets.enc.yaml | psm sync --store ssm --skip-approve /dev/stdin
```

> パイプ使用時は stdin がターミナルではないため、変更を実行するには `--skip-approve` が必要です。

鍵の生成や暗号化設定を含む詳しい手順は [example/README.md](example/README.md) を参照してください。

### YAML フォーマット

キーがそのまま AWS のパラメータ名になります。値はスカラー値（string, int, bool, float）のみ使用できます。

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
