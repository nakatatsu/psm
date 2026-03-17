# Survey: README.md に Usage セクションを追加する

**Date**: 2026-03-16
**Spec**: [spec.md](spec.md)

## Summary

ドキュメントのみの小タスク。問題定義は明確で方向性に疑問なし。主な検討点は (1) 英語/日本語の2ファイル構成における情報の一貫性、(2) example/README.md との内容重複の扱い、(3) CLI ヘルプ出力との正確な一致。リスクは低い。

## S1: Problem Reframing — ドキュメント不足は本当に問題か

**Category**: A. Problem Reframing
**Finding**: README に Usage がないのは明らかな問題。`psm --help` は存在するが、SOPS パイプパターンや secrets.yaml フォーマットはヘルプに含まれない。example/README.md にはほぼ全情報があるが、本体READMEからの導線がない。問題定義は正しい。
**Recommendation**: spec 通り進めてよい。
**Evidence**: 現在の README.md は Install と AWS SSO login の2セクションのみ。`psm` バイナリが手元にないため `--help` の出力は未確認だが、spec 001-psm の FR-024 により Go flag パッケージのデフォルト出力に従う設計。

## S2: example/README.md との重複

**Category**: B. Solution Evaluation
**Finding**: example/README.md にはすでに psm sync の使い方（dry-run、SOPS パイプ）、secrets.yaml のフォーマット例が詳細に記載されている。本体 README に同じ内容を書くと二重メンテになる。
**Recommendation**: 本体 README は簡潔なリファレンス（コマンド構文 + 主要オプション一覧 + 最小限の例）に留め、example/README.md の詳細チュートリアルへリンクで誘導する。secrets.yaml フォーマット例は本体 README にも載せるが、SOPS 連携の詳細手順は example/ に委ねる。
**Evidence**: example/README.md L84-87 に SOPS パイプパターン、L1-4 に secrets.yaml のフォーマットが既存。

## S3: 英語/日本語 2ファイル構成

**Category**: C. Risk & Feasibility
**Finding**: Addition で README.md（英語）と README.ja.md（日本語）の2ファイルが要求されている。翻訳の同期がメンテコストになるリスクはあるが、OSSとしては標準的なアプローチ。現在の README.md は英語ベースのため、既存内容の書き換えは最小限で済む。
**Recommendation**: 英語を正とし、日本語は翻訳版として明示する。README.ja.md の冒頭に「This is a Japanese translation of [README.md](README.md)」のような注記を入れる。
**Evidence**: GitHub は README.md をデフォルト表示するため、英語版が自動的にメインとなる。

## S4: CLI オプション体系の正確性

**Category**: C. Risk & Feasibility
**Finding**: psm バイナリが DevContainer にインストールされていないため、`--help` 出力を直接確認できない。spec 001-psm の FR 定義と Go ソースコードが信頼できる情報源。
**Recommendation**: 実装時に Go ソースの flag 定義を読んでコマンド例を正確に記述する。可能であれば `go run . --help` で出力を確認する。
**Evidence**: spec 001-psm FR-001, FR-002 にフラグ体系が定義済み。

## S5: 既存セクションのリファクタリング

**Category**: D. Integration Impact
**Finding**: 現在の「Access to AWS from DevContainer」セクションは DevContainer 固有の情報であり、Usage とは別の関心事。Usage セクション追加に伴い、既存セクションの再配置が必要になる可能性がある。
**Recommendation**: セクション順は Install → Usage → (その他) とし、DevContainer 固有の情報は Development セクション等に移動するか、そのまま残す。spec の FR-006 に従う。
**Evidence**: 現在の README.md は Install → Access to AWS の2セクション構成。

## Items Requiring PoC

- CLI ヘルプ出力の確認: `go run . --help` および `go run . sync --help` で実際のフラグ名・説明を確認すべき（実装時に実施）

## Constitution Impact

本タスクはドキュメントのみのため、constitution への影響なし。テストファースト原則はコード変更がないため適用外。

## Recommendation

小タスクにつき survey で十分。直接実装に進んでよい。実装時の注意点:
1. Go ソースの flag 定義を確認してコマンド例の正確性を担保する
2. example/README.md との重複を避け、リンクで誘導する
3. 英語を正として日本語は翻訳版とする
