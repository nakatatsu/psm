# Data Model: --store フラグの除去

**Date**: 2026-03-27

## Entity Changes

### Config (store.go)

**変更**: `Store string` フィールドを削除

```diff
 type Config struct {
     Subcommand  string
-    Store       string
     Profile     string
     DryRun      bool
     SkipApprove bool
     Debug       bool
     DeleteFile  string
     File        string
     ShowVersion bool
 }
```

### Store interface (store.go)

**変更なし** — interface は既存のまま保持。将来の拡張ポイントとして価値がある。

### SSMStore (ssm.go)

**変更なし** — 実装は変わらない。`run()` から直接 `NewSSMStore` が呼ばれるようになるだけ。
