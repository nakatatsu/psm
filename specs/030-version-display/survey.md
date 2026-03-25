# Survey: Version Display Command

**Date**: 2026-03-25
**Spec**: [spec.md](spec.md)

## Summary

この機能は非常にシンプルで、方向性に根本的な問題はない。ただし、既存のリリースパイプライン（GoReleaser）との連携が spec で考慮されていない点が最大の発見。GoReleaser はデフォルトで `main.version`, `main.commit`, `main.date` を `ldflags` 経由で注入するため、これを活用すれば追加のビルド設定は不要。ただし `.goreleaser.yaml` の `ldflags` に `-X main.version={{.Version}}` を明示的に追加する必要がある（現在は `-s -w` のみ）。

## S1: GoReleaser との連携

**Category**: Integration & External Dependencies
**Finding**: リリースパイプラインは GoReleaser v2 を使用している（`.goreleaser.yaml`）。現在の `ldflags` は `-s -w`（デバッグ情報の削除）のみで、バージョン注入は設定されていない。GoReleaser は `{{.Version}}`, `{{.Commit}}`, `{{.Date}}` テンプレート変数を提供しており、`ldflags` に追加するだけでタグからバージョンが自動注入される。
**Recommendation**: `.goreleaser.yaml` の `ldflags` に `-X main.version={{.Version}}` を追加する。手動ビルド時のフォールバックは `dev` でよい。
**Evidence**: `.goreleaser.yaml` L16-17 で `ldflags: ["-s -w"]` のみ確認。GoReleaser v2 ドキュメントでテンプレート変数によるバージョン注入がサポートされている。

## S2: `--version` フラグの解析位置

**Category**: Integration Impact
**Finding**: 現在の `parseArgs()` は `args[1]` をサブコマンド（`sync` / `export`）として解釈し、それ以外はエラーとなる（`main.go:177`）。`--version` を追加するには、サブコマンド解析の**前**に `args[1]` を検査する必要がある。これは既存のコードフローへの変更が最小限で済む。
**Recommendation**: `parseArgs()` の先頭で `args[1] == "--version"` をチェックし、早期リターンする。既存のサブコマンド解析には影響しない。
**Evidence**: `main.go:171-229` の `parseArgs` 実装を確認。

## S3: 問題の妥当性

**Category**: Problem Reframing
**Finding**: バージョン表示機能は CLI ツールの基本機能であり、必要性に疑問はない。代替アプローチ（例: `psm version` サブコマンド）もあり得るが、Go CLI の慣例として `--version` フラグが最も一般的であり、spec の方向性は妥当。`psm version` サブコマンドを追加すると、既存のサブコマンド解析と一貫性の問題が生じるため、`--version` フラグのみで十分。
**Recommendation**: spec 通り `--version` フラグのみで進める。
**Evidence**: Go 標準ツール群（`go version`）はサブコマンド形式だが、多くの CLI ツール（kubectl, terraform, docker）は `--version` フラグをサポート。psm の既存アーキテクチャ（フラグベース）に合致するのは `--version` フラグ。

## S4: スコープの確認

**Category**: Scope Boundaries
**Finding**: spec はバージョン表示のみに限定されており、適切にスコープされている。`commit` や `date` の追加表示は YAGNI 原則に基づき不要（constitution I, II に合致）。将来必要になれば変数を追加するだけで済む。
**Recommendation**: spec 通り、バージョン文字列のみ表示。
**Evidence**: Constitution I (Simplicity First), II (YAGNI) を確認。

## Items Requiring PoC

なし。すべて既存のコードとドキュメントで検証可能。

## Constitution Impact

修正不要。この機能は constitution のすべての原則に合致する：
- 標準ライブラリのみで実装可能（Principle I）
- 現在必要な機能のみ実装（Principle II）
- テストファースト開発が可能（Principle III）

## Recommendation

問題なし。`/speckit.plan` に進んでよい。plan 時に S1（GoReleaser ldflags 設定）を忘れずに含めること。
