# Data Model: psm

## Entities

### Config

CLI フラグから構成された実行設定。

| Field      | Type   | Description                          |
|------------|--------|--------------------------------------|
| Subcommand | string | `"sync"` or `"export"` |
| Store      | string | `"ssm"` or `"sm"` |
| Profile    | string | AWS プロファイル名（任意） |
| Prune      | bool   | `--prune` フラグ（sync のみ） |
| DryRun     | bool   | `--dry-run` フラグ（sync のみ） |
| File       | string | 入力/出力ファイルパス |

### Entry

YAML から読み取った 1 つの key-value ペア。

| Field | Type   | Description    |
|-------|--------|----------------|
| Key   | string | YAML キー      |
| Value | string | 値（文字列）   |

### Action

個々の Entry に対する同期操作の計画と結果。

| Field | Type   | Description                              |
|-------|--------|------------------------------------------|
| Key   | string | 対象キー                                 |
| Type  | enum   | `create`, `update`, `delete`, `skip`     |
| Error | error  | 失敗時のエラー（成功時は nil）           |

`skip` は no-change（内部処理用、表示しない）。

### Summary

全 Action の処理結果集計。

| Field     | Type | Description      |
|-----------|------|------------------|
| Created   | int  | 新規作成件数     |
| Updated   | int  | 更新件数         |
| Deleted   | int  | 削除件数         |
| Unchanged | int  | 変更なし件数     |
| Failed    | int  | エラー件数       |

## Relationships

```
Config 1---* Entry     (Config.File を読んで Entry を生成)
Entry  1---1 Action    (各 Entry に対して 1 つの Action を計算)
Action *---1 Summary   (全 Action を集計して Summary を生成)
```

## State Transitions

### Action.Type の決定ロジック

```
Entry が YAML にある & AWS にない        → create
Entry が YAML にある & AWS と値が異なる  → update
Entry が YAML にある & AWS と値が同じ    → skip
Entry が AWS にある  & YAML にない & prune=true  → delete
Entry が AWS にある  & YAML にない & prune=false → (無視、Action 生成しない)
```
