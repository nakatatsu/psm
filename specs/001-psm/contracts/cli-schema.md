# CLI Contract: psm

動作ルール・バリデーション・命名規則は [spec.md](../spec.md) を参照。本ドキュメントは CLI インターフェースの形式定義のみを扱う。

## Subcommands

| Subcommand | Description              |
| ---------- | ------------------------ |
| `sync`     | YAML → AWS 同期          |
| `export`   | AWS → YAML エクスポート  |

サブコマンドは必須。省略時は使い方を表示し終了コード 1。

## Common Flags

両サブコマンド共通:

| Flag        | Required | Type   | Default | Description              |
| ----------- | -------- | ------ | ------- | ------------------------ |
| `--store`   | Yes      | string | —       | `ssm` or `sm`            |
| `--profile` | No       | string | —       | AWS プロファイル名       |

---

## `psm sync`

```
psm sync --store <ssm|sm> [--profile <name>] [flags] <file>
```

### Arguments

| Position | Name | Required | Description                |
| -------- | ---- | -------- | -------------------------- |
| 1        | file | Yes      | 復号済み YAML ファイルパス |

### Flags (sync 固有)

| Flag        | Required | Type | Default | Description                                                                                 |
| ----------- | -------- | ---- | ------- | ------------------------------------------------------------------------------------------- |
| `--prune`   | No       | bool | false   | YAML にないキーを AWS から削除（スコープは spec.md FR-008 参照）                            |
| `--dry-run` | No       | bool | false   | 実行計画の表示のみ、変更を実行しない。--pruneの有無で実行計画が変わり、こちらの表示も変わる |

### stdout Format

差分行（変更があるキーのみ、no-change は出力しない）:

```
create: {key}
update: {key}
delete: {key}
```

サマリー行（最終行、常に表示）:

```
{N} created, {N} updated, {N} deleted, {N} unchanged, {N} failed
```

`--dry-run` 時も同一形式（failed は常に 0）。

### stderr Format

```
error: {key}: {error message}
```

---

## `psm export`

```
psm export --store <ssm|sm> [--profile <name>] <output-file>
```

### Arguments

| Position | Name        | Required | Description        |
| -------- | ----------- | -------- | ------------------ |
| 1        | output-file | Yes      | 出力先 YAML パス   |

### Output Format

出力ファイルは純粋な key-value YAML（メタデータなし）:

```yaml
/myapp/prod/DATABASE_URL: "postgres://..."
/myapp/prod/API_KEY: "sk-..."
```

---

## Exit Codes (共通)

| Code | Meaning                          |
| ---- | -------------------------------- |
| 0    | 成功                             |
| 1    | 失敗、または入力エラー           |
