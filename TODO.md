# TODO: AWS Params/Secrets Sync CLI Tool

## やること
- [ ] `spec.md` を書く（要件定義＋技術要望、簡潔に）
- [ ] 実装

## 決まっていること
- **言語**: Go
- **用途**: CI/CDで使うCLIツール
- **入力**: SOPSデフォルト準拠形式のテキストファイル（YAML or JSON、形式未確定）
- **出力先**: AWS SSM Parameter Store / Secrets Manager に反映
- **方針**: 無駄に書きすぎない。「あったほうがいい」程度のものは不要

## 調査済み: AWS API要点
- **SSM PutParameter**: `Overwrite=true`でupsert可能。型は String / StringList / SecureString
- **SSM パス規約**: `/{app}/{env}/{key}`
- **Secrets Manager**: upsertなし。CreateSecret → ResourceExistsException → PutSecretValue のフロー必要
- **Secrets Manager**: 値は通常JSON。常に暗号化。$0.40/secret/月
- **SSM**: Standard tier無料、4KB制限
- **Secrets Manager**: 64KBまで

## 未調査
- SOPSのデフォルトファイル形式の詳細（YAML/JSON構造、metadataブロック）

## ネットワーク制約
- この環境から外部接続不可（curl Yahoo確認済み）
