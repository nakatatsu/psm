# Research: Secrets Manager 対応オミット

**Date**: 2026-03-18

## Findings

このタスクは既存コードの削除と既存ドキュメントの修正が主体であり、技術的な不明点はない。

### Decision: Store interface の保持

- **Decision**: `store.go` の Store interface は削除せず保持する
- **Rationale**: 将来的に別のストアバックエンド（HashiCorp Vault 等）を追加する際の拡張ポイントとして有用。interface 自体のメンテナンスコストはゼロに近い
- **Alternatives considered**: interface ごと削除して SSM 直接呼び出しに変更 → テスト時の DI が困難になるため却下

### Decision: `--store sm` のエラーハンドリング

- **Decision**: バリデーションで `--store ssm` のみ許可し、それ以外は既存の invalid store エラーで処理
- **Rationale**: SM 専用のマイグレーションメッセージは過剰。`--prune` のような移行パスが必要なケースではない（SM ユーザーが現時点でいないため）
- **Alternatives considered**: `--store sm` に対して専用の「removed」メッセージを出す → YAGNI。汎用の invalid store エラーで十分
