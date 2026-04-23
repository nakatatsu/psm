# Survey: テスト実行基盤の整備

**Date**: 2026-04-21
**Spec**: [spec.md](spec.md)

## Summary

specの方向性は妥当。結合テスト資材の分離、CI自動化、OIDC認証のいずれも合理的な判断。ただしconstitutionに「stub tests (for CI) use an AWS emulator (e.g., moto)」という記載があり、CIでの結合テストが実AWS接続である点はconstitution上の整合性を確認する必要がある。また、test.shはbashスクリプトでGoのテストフレームワークを使っていないため、constitutionの「All tests use `go test`」との関係を明確にすべき。

## S1: 問題定義の妥当性

**Category**: Problem Reframing
**Finding**: specが解決しようとしている問題は「結合テストが手動で品質ゲートとして機能していない」こと。これは正しい問題認識。手動テストはリリース前に実行されないリスクがあり、自動化は必然。
**Recommendation**: そのまま進めてよい。
**Evidence**: 現状example/test.shは手動実行のみ。CIには結合テスト相当のジョブがない（`.github/workflows/ci.yml` はユニットテストと静的解析のみ）。

## S2: constitutionとの整合性 — テスト方式

**Category**: Constitution Compliance
**Finding**: constitutionは「Stub tests (for CI) use an AWS emulator (e.g., moto)」と定めている。一方、今回CIで走らせるtest.shは実AWSに接続する結合テスト。constitutionの「Integration tests (real AWS) use a sandbox environment」の定義に合致するが、CIで実AWSを使うことがstub test方針と矛盾する可能性がある。
**Recommendation**: test.shは「integration tests (real AWS)」に分類される。constitutionはstub testsとintegration testsを区別しており、integration testsについてはCIでの実AWS使用を禁止していない。ただし、constitutionの意図は「CIではemulatorを使い、実AWSはサンドボックスで手動」だった可能性がある。今回はユーザーとの合意に基づきCIで実AWS接続する方針で進める。constitution amendmentの要否は実装後にユーザーに確認する。
**Evidence**: constitution Section III: "Stub tests (for CI) use an AWS emulator" / "Integration tests (real AWS) use a sandbox environment with dedicated test prefixes and setup/teardown"

## S3: constitutionとの整合性 — テストフレームワーク

**Category**: Constitution Compliance
**Finding**: constitutionは「All tests use `go test` from the Go standard library」と定めている。test.shはbashスクリプトであり、`go test`ではない。
**Recommendation**: test.shはCLIバイナリを実環境でテストする結合テストであり、`go test`の範囲外。constitutionの「All tests use `go test`」はGoコードのユニットテストを指していると解釈する。ただし将来的にconstitutionの文言を明確化すべき。
**Evidence**: test.shは `go test` ではなくbashで直接psmバイナリを実行し、`aws ssm` コマンドで結果を検証している。

## S4: GitHub Actions OIDCの実現可能性

**Category**: Feasibility Verification
**Finding**: GitHub Actions OIDCからAWSへの認証は広く採用されている方式。公開リポジトリでも利用可能。IAM信頼ポリシーで `StringLike` を使えば `release-*` ブランチのワイルドカード指定が可能。
**Recommendation**: 技術的に問題なし。`aws-actions/configure-aws-credentials` アクションを使えばワークフロー側の実装はシンプル。
**Evidence**: GitHub公式ドキュメント、AWS公式ドキュメントで確認済みの標準パターン。

## S5: テスト並行実行の競合リスク

**Category**: Risk & Failure Modes
**Finding**: Edge Casesに記載の通り、developとrelease-*で同時に結合テストが走ると `/myapp/` 配下のパラメータが競合する。specでは「失敗でよい」としているが、頻発するとCI信頼性が低下する。
**Recommendation**: 現時点では発生頻度が低い（developとreleaseが同時にpushされることは稀）ため、specの判断（失敗許容）で問題ない。将来的に頻発する場合はテストプレフィックスをブランチ名で分離する対応を検討。
**Evidence**: GitFlow運用上、developへのマージとリリースブランチの更新が同時刻に重なる可能性は低い。

## S6: example/からの分離アプローチ

**Category**: Approach Alternatives
**Finding**: test.shをexample/からtests/integration/にコピーする方針。example/は触らない制約がある。test.shはexample/のsecrets.example.yamlを参照しているため、テストデータも合わせてコピーする必要がある。
**Recommendation**: tests/integration/にtest.sh、secrets.example.yaml、必要なテストデータをコピーし、パスを調整する。example/は現状維持。
**Evidence**: test.shの `TEST_YAML="${SCRIPT_DIR}/secrets.example.yaml"` がローカルパス依存。

## Items Requiring PoC

- OIDC信頼ポリシーの `StringLike` + `release-*` パターンの動作確認（IAM設定後に実際のワークフローで検証が必要）

## Constitution Impact

- constitutionの「Stub tests (for CI) use an AWS emulator」と「All tests use `go test`」の文言について、結合テスト（bash + 実AWS）のCI実行との整合性を将来的に明確化すべき。ただし今回はamendment不要（constitutionのintegration test定義の範囲内と解釈）。

## Recommendation

specの方向性に問題なし。`/speckit.plan` に進めてよい。
