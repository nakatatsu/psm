# KPT: 001-psm (SOPS-to-AWS Parameter Sync CLI)

**Date**: 2026-03-10
**Branch**: `001-psm`

## Keep

### テスト駆動 / 契約ベース開発
- spec.md の FR → contracts/ → テスト → 実装の流れが、LLM 実装で真価を発揮した
- 人間なら「わかってるから省略」できるものが、LLM には契約の明示がないと暴走する。逆に契約があれば精度が劇的に上がる
- 早すぎる技術だったテスト駆動が、LLM 時代には必須技術になった

### SpecKit のノイズ低減
- 以前と比較してノイズが大幅に減少
- 27タスクが明確に定義され、迷いなく実装できた

### 設計判断の質
- Constitution "Simplicity First" が効いて不要な抽象化を防止
- 全ファイル `package main` のフラット構成が結果的に正解
- FR Coverage 表で要件→実装→テストの対応が追える状態を維持

## Problem

### 開発環境 (DevContainer + AWS)
- AWS マネージドサービス利用にあたって、DevContainer のネットワーク制限が信頼性を大きく犠牲にした
- staticcheck インストール不可（proxy.golang.org がファイアウォール外）
- SSM/SM の統合テストで eventual consistency 問題が顕在化
- SM の `BatchGetSecretValue` が IAM ポリシー不足で失敗
- 今回はロジックが薄くユニットテストの意味が薄い一方、結合テストの比重が非常に重かったため、ネットワーク問題の影響が強く出た

### スペック段階での考慮不足
- AWS API の整合性モデル（eventual consistency）が research に含まれておらず、テスト設計で後から全面修正
- IAM 必要権限の事前洗い出しが不足

## Try

### Sandbox 環境の構築
- ネットワーク制限を外した Sandbox を真面目に検討する（渡せる権限が減る問題はあるが、結合テスト主体のプロジェクトでは必要）

### スペックフェーズの強化
- research に「AWS API 整合性モデル」を含める
- contracts/ に IAM 必要権限一覧を含める
- Dockerfile のツール取得方針（`go install` vs GitHub Releases binary）を明文化

### 新機能追加の SpecDriven + SpecKit フロー
- 既存コードベースへの spec → plan → implement フローの確立
- 「既存コードの理解」「影響範囲の特定」をどのフェーズで吸収するかが鍵
