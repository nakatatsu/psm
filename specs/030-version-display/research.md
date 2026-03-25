# Research: Version Display Command

**Date**: 2026-03-25

## R1: バージョン注入方式

**Decision**: `go build -ldflags "-X main.version=..."` によるビルド時注入
**Rationale**: Go の標準的な方法であり、外部依存なし。GoReleaser が `{{.Version}}` テンプレートで自動注入をサポートしている。
**Alternatives considered**:
- `runtime/debug.ReadBuildInfo()` — Go モジュールのバージョン情報を取得できるが、`go install` 以外のビルドでは空になることがあり、制御が難しい
- 定数ファイル生成 — `go generate` でバージョンファイルを生成する方法。ビルドステップが増え、constitution I (Simplicity First) に反する

## R2: GoReleaser 設定変更

**Decision**: `.goreleaser.yaml` の `ldflags` に `-X main.version={{.Version}}` を追加
**Rationale**: 既存のリリースパイプラインに自然に統合できる。現在の `ldflags` は `-s -w` のみで、バージョン注入が欠けている。
**Alternatives considered**:
- GoReleaser のデフォルト ldflags に任せる — GoReleaser v2 はデフォルトで `main.version` 等を注入するが、明示的に指定する方が確実で可読性が高い

## R3: `--version` フラグの解析位置

**Decision**: `parseArgs()` 内でサブコマンド解析の前に `args[1] == "--version"` をチェック
**Rationale**: 既存の `flag.FlagSet` はサブコマンドごとに作成されるため、グローバルフラグとして `--version` を追加するのは複雑。`args[1]` の直接比較が最もシンプル。
**Alternatives considered**:
- `flag.FlagSet` にグローバル `--version` を追加 — 現在のアーキテクチャ（サブコマンド別 FlagSet）と相性が悪い
- `psm version` サブコマンドとして追加 — 動作するが、`sync`/`export` と同列に並ぶのは意味的に不自然。`--version` フラグの方が CLI 慣例に合致
