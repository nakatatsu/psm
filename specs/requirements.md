# psm — Behavioral Requirements

> アプリケーションの挙動に関する要件をまとめたドキュメント。コード・仕様書・テストから検証済みの事実のみを記載。(2026-03-25)

## 1. Commands

### 1.1 sync

YAML ファイルの内容を AWS SSM Parameter Store に同期する。

```
psm sync --store ssm [--profile <name>] [--dry-run] [--skip-approve] [--debug] [--delete <file>] <sync-file>
```

| Flag | Required | Description |
|------|----------|-------------|
| `--store ssm` | Yes | ストア種別（`ssm` のみ有効） |
| `--profile <name>` | No | AWS プロファイル名（未指定時は SDK デフォルト認証） |
| `--dry-run` | No | 変更計画を表示するのみ。プロンプトも実行も行わない |
| `--skip-approve` | No | 承認プロンプトをスキップして即時実行 |
| `--debug` | No | Debug レベルのログを stderr に出力 |
| `--delete <file>` | No | 削除対象の正規表現パターンを含む YAML ファイル |

**処理フロー**:

1. YAML ファイルを読み込み、エントリをパース・バリデーション
2. AWS SSM から既存パラメータを全件取得
3. YAML と既存パラメータを比較し、アクションプラン（create/update/skip）を生成
4. `--delete` 指定時は削除パターンを処理し、アクションプランに delete を追加
5. アクションプランを stdout に表示
6. 承認プロンプトを表示（条件による — 3.3 参照）
7. 変更を実行

### 1.2 export

AWS SSM Parameter Store の全パラメータを YAML ファイルとして書き出す。

```
psm export --store ssm [--profile <name>] [--debug] <file>
```

- 出力ファイルが既に存在する場合はエラー（上書き禁止）
- パラメータが 0 件の場合はエラー
- キーはアルファベット順にソート
- 値はダブルクォート形式
- ファイル権限は `0600`
- export した YAML をそのまま sync すると差分 0 件（ラウンドトリップ互換）

## 2. CLI Validation

以下の条件でエラー（終了コード 1）:

- サブコマンド未指定または不正な値
- `--store` フラグ未指定
- `--store` に `ssm` 以外の値を指定（`sm` を含む）
- `--prune` フラグを指定（削除済み。`--delete <file>` を案内するエラーメッセージ）
- ファイル引数が 0 個または 2 個以上

環境変数 `AWS_PROFILE` は常に無視される（`--profile` フラグのみ有効）。

## 3. Features

### 3.1 Sync（作成・更新）

- YAML のキー名がそのまま SSM パラメータ名になる（パス構造の組み立てはしない）
- 型は常に `SecureString`、`Overwrite=true` で upsert
- 同期前に既存値を取得し比較。値が同一のキーは操作も表示もしない（skip）
- YAML の非文字列値（整数、ブール値、浮動小数点数）は文字列に自動変換
- 最大 10 並行で Put 操作を実行（固定値、設定不可）
- 個別キーの失敗時も残りの処理を続行（部分失敗許容）
- 失敗したキーは stderr にキー名とエラーメッセージを出力（値は含めない）

### 3.2 Delete（パターンベース削除）

`--delete <file>` で削除対象の正規表現パターンを指定する。

**削除パターンファイル形式**:
```yaml
- "^/myapp/legacy/"
- "^/myapp/deprecated-.*"
```

**分類ロジック**:
- **削除対象**: パターンにマッチし、かつ sync YAML に存在しないキー
- **コンフリクト**: パターンにマッチし、かつ sync YAML にも存在するキー → 全操作を中止
- **Unmanaged**: sync YAML に存在せず、パターンにもマッチしないキー → 警告表示（実行はブロックしない）

**安全機能**:
- コンフリクト検出時は create/update/delete すべてを中止（all-or-nothing）。変更は一切行われない
- 不正な正規表現は AWS 操作前にエラー
- 削除は AWS API のバッチ処理（最大 10 件/リクエスト）

### 3.3 Approval Flow（承認フロー）

- デフォルトで変更前に `Proceed? [y/N]` プロンプトを stderr に表示
- `y` または `Y` のみが承認（それ以外はすべて拒否）
- 拒否した場合は終了コード 0（変更なし）
- `--skip-approve` でプロンプトをスキップし即時実行
- `--dry-run` ではプロンプトを表示しない
- 変更がない場合（全件 skip）はプロンプトを表示しない
- stdin が端末でない場合（パイプ入力）かつ `--skip-approve` 未指定時は自動的に拒否（終了コード 0）

### 3.4 Dry-run

- アクションプラン（create/update/delete）を stdout に表示
- サマリー行に `(dry-run)` を付与
- プロンプトは表示しない
- AWS への変更は一切行わない

### 3.5 SOPS 連携

- YAML 内の `sops` キー（SOPS メタデータ）は自動的に除外
- psm 自体は復号を行わない（事前に SOPS で復号された入力を想定）
- パイプで連携: `sops -d secrets.enc.yaml | psm sync --store ssm --skip-approve /dev/stdin`

### 3.6 Debug Logging

- `--debug` フラグで Debug レベルのログを有効化（両サブコマンド共通）
- デフォルトは Info レベル（Error + Warn + Info が表示。Debug は非表示）
- ログは stderr に出力（`slog.TextHandler` 形式）

## 4. Input/Output Format

### 4.1 入力 YAML（sync）

```yaml
/myapp/database/host: localhost
/myapp/database/port: 5432
/myapp/api/key: my-secret-password
```

**バリデーションルール**（AWS 通信前に一括実行）:
- トップレベルはマッピングのみ
- 値はスカラー値（string, int, bool, float）のみ。マップ・配列は不可
- null 値は不可
- 空キーは不可
- 重複キーは不可
- 空文字列の値は許可
- `sops` キーは除外後にバリデーション（`sops` の値がマップでも誤検出しない）
- `sops` 除外後にキーが 0 件の場合はエラー
- バリデーションエラー時は何が問題かを示すメッセージを表示し終了コード 1

### 4.2 stdout 出力

**アクションプラン**（1 行/キー）:
```
create: /myapp/prod/DB_URL
update: /myapp/prod/DB_PORT
delete: /myapp/legacy/OLD_KEY
```
- skip（変更なし）は表示しない
- 値は絶対に表示しない
- 通常実行と `--dry-run` で同一形式

**サマリー行**:
```
2 created, 1 updated, 1 deleted, 5 unchanged, 0 failed
2 created, 1 updated, 1 deleted, 5 unchanged, 0 failed (dry-run)
```

### 4.3 stderr 出力

- すべての slog メッセージ（Error, Warn, Info, Debug）
- 承認プロンプト（`Proceed? [y/N]`）
- 個別キーのエラー（例: `error: /myapp/prod/API_KEY: AccessDeniedException: ...`）

### 4.4 出力 YAML（export）

```yaml
/myapp/database/host: "localhost"
/myapp/database/port: "5432"
```
- キーはアルファベット順
- 値はダブルクォート形式

## 5. Exit Codes

| Code | Condition |
|------|-----------|
| 0 | 全件成功 |
| 0 | 承認プロンプトで拒否 |
| 0 | 非端末 stdin かつ `--skip-approve` 未指定（自動拒否） |
| 1 | CLI 引数の不正 |
| 1 | YAML バリデーションエラー |
| 1 | 削除パターンの正規表現不正 |
| 1 | コンフリクト検出（削除パターンと sync YAML の競合） |
| 1 | 1 件以上の同期失敗 |
| 1 | export 先ファイルが既存 |
| 1 | export 対象パラメータが 0 件 |
| 1 | AWS GetAll 自体の失敗（致命的エラー） |

## 6. Security Constraints

- パラメータの値はいかなる出力（stdout, stderr, ログ）にも絶対に含めない
- キーパス（例: `/myapp/prod/API_KEY`）はログに記録可能（値ではないため）
- export ファイルの権限は `0600`（オーナーのみ読み書き）

## 7. AWS Integration Details

- SSM PutParameter: `SecureString` 型、`Overwrite=true`、並行度 10 で個別実行
- SSM DeleteParameters: 1 リクエストあたり最大 10 件のバッチ処理
- SSM GetParametersByPath: `path=/`, `Recursive=true`, `WithDecryption=true`、ページネーション対応
- タイムアウト: AWS SDK のデフォルトに委ねる（psm 側での設定なし）
- リージョン: 環境変数 `AWS_REGION` / `AWS_DEFAULT_REGION` または SDK デフォルト
- キー名バリデーション: psm 側では行わない（AWS API エラーに委ねる）
- 値サイズ制限（SSM: 4KB standard / 8KB advanced）: psm 側では検証しない（AWS API エラーに委ねる）
