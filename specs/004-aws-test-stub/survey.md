# Survey: AWS Test Stub Infrastructure

**Date**: 2026-03-15
**Spec**: specs/004-aws-test-stub/spec.md

## Summary

spec は「moto で AWS テストをスタブ化する」ことを前提にしているが、そもそもの問題は「CI で回帰を検出したい」であり、moto はその手段の一つに過ぎない。既存テストのカバレッジを分析すると、AWS 依存テスト 9 件がカバーする固有のロジックは実は非常に少なく（SM の create-or-update 分岐と SSM のバッチ削除ロジック程度）、もっと軽い手段で同等のカバレッジが得られる可能性がある。

## A. Problem Reframing

### S1: Problem Definition — spec は正しい問題を解いているか？

**Category**: Problem Definition
**Finding**: spec が解こうとしている問題は「AWS 認証なしで統合テストを実行する」だが、本当の目的は「CI で回帰を検出する」こと。この二つは同じではない。CI で回帰を検出するために、必ずしも AWS エミュレータが必要とは限らない。
**Recommendation**: 「CI で回帰を検出する」をゴールとして再評価し、手段を比較する。
**Evidence**: 後述 S2, S3 の分析。

### S2: Hidden Assumptions — 疑われていない前提

**Category**: Hidden Assumptions
**Finding**: spec には以下の前提が暗黙に含まれている:

1. **「AWS テスト 9 件がなければ CI は不十分」** — 本当か？ 既存の非 AWS テスト 10 件が何をカバーしているか検証が必要。
2. **「Store 実装のテストには AWS（またはエミュレータ）が必要」** — Store 実装（ssm.go: ~75 行、sm.go: ~93 行）はほぼ SDK ラッパー。固有のロジックは SM の create-or-update 分岐（sm.go:60-79）と SSM のバッチ削除（ssm.go:61-74）のみ。
3. **「エミュレータは実 AWS と同等の信頼性を持つ」** — 001-psm の research.md で「エミュレータの挙動差異に起因するバグは許容できない」と明記されている。moto を入れてもこの懸念は残る。

**Recommendation**: 前提 1, 2 を検証し、エミュレータなしで CI カバレッジを上げられるか先に検討する。
**Evidence**: ssm.go, sm.go のコード分析。sync_test.go の fakeStore パターン。

### S3: Lateral Thinking — moto を使わずに同じゴールを達成する方法

**Category**: Lateral Thinking
**Finding**: 以下の代替アプローチが存在する:

**アプローチ A: 何もしない（skip のまま）**
- 非 AWS テスト 10 件は YAML パース、plan/diff ロジック、execute（fakeStore）、CLI パース、export をカバー済み
- AWS テスト 9 件の固有価値は薄い（SDK ラッパーのテスト）
- コスト: ゼロ。リスク: SM create-or-update 分岐のみ未カバー

**アプローチ B: SDK レベルの interface fake を追加（推奨候補）**
- SM の create-or-update 分岐と SSM のバッチ削除ロジックを、SDK interface の fake でテスト
- sync_test.go の `fakeStore` パターンの延長。プロダクションコード変更なし
- 2-3 個のテスト追加で、AWS テスト 9 件の固有ロジックの大部分をカバー
- コスト: 小。メンテナンスコスト: 低。Docker 依存なし

**アプローチ C: GitHub Actions OIDC で実 AWS を使う**
- GHA OIDC federation で AWS 認証情報なしに実 AWS テスト実行
- パブリックリポジトリでは fork PR からの悪用リスク
- Secrets Manager は $0.40/secret/月（テスト後即削除でも課金）
- コスト: 中。信頼性: AWS の一時的な障害に左右される

**アプローチ D: ハイブリッド（A + B）**
- CI では非 AWS テスト + SDK interface fake テストを実行
- 実 AWS テストはローカル/手動のプレリリースゲートとして維持
- コスト: 小。カバレッジ: アプリケーションロジック全体

**アプローチ E: moto エミュレータ（spec の提案）**
- Docker サービスとして moto を追加、endpoint URL 差し替え
- コスト: 中。docker-compose 追加、moto バージョン管理、API 互換性検証
- カバレッジ: SDK 呼び出し構造まで含む（ただし moto の忠実度に依存）

**Recommendation**: アプローチ D（ハイブリッド）が最もコスト対効果が高い。150 行の SDK ラッパーのためにコンテナ化されたエミュレータを導入するのは過剰。ただし、将来 Store 実装が複雑化した場合は moto 導入を再検討する価値がある。
**Evidence**: ssm.go (~75 行), sm.go (~93 行) のコード量。sync_test.go の fakeStore パターンが既に存在。

## B. Solution Evaluation

### S4: Necessity — moto は本当に必要か？

**Category**: Necessity
**Finding**: Store 実装の固有ロジックは SM の `CreateSecret` → `ResourceExistsException` → `PutSecretValue` フォールバック（sm.go:60-79）と SSM の 10 件バッチ削除（ssm.go:61-74）のみ。残りは直線的な SDK 呼び出し。これらは SDK レベルの interface fake で十分テスト可能。
**Recommendation**: moto の必要性を「コスト対効果」で判断する。現時点のコード量（150 行の AWS ラッパー）に対してコンテナ化エミュレータは過剰の可能性が高い。
**Evidence**: sm.go の `Put` メソッド、ssm.go の `Delete` メソッドのコード分析。

### S5: Prior Decisions — 過去の決定は今も有効か？

**Category**: Prior Decisions
**Finding**: 001-psm の research.md R5 で以下の決定がされた:
> "Store interface を mock しても SDK ラッパーの薄い層をスキップするだけで、テストとしての価値がない"

この決定は「実 AWS テストが実行できる環境」を前提としていた。CI 導入で前提が変わった。ただし、決定の核心（Store は薄いラッパーなので mock しても価値が低い）は今も正しい。

**重要**: この決定が正しいからこそ、moto も同様に「薄いラッパーのテスト」であり、導入コストに見合わない可能性がある。

**Recommendation**: 過去の決定は部分的に有効。「mock は価値が低い」→「だから moto も価値が低いかもしれない」という帰結に注意。
**Evidence**: specs/001-psm/research.md R5。

### S6: Cost & Complexity

**Category**: Cost & Complexity
**Finding**: moto 導入のコスト:
- docker-compose にサービス追加
- moto バージョンのピン止め・更新管理
- API 互換性の継続的な検証（moto の回帰バグリスク）
- GHA で docker compose を起動するワークフロー複雑化
- テストコードのヘルパー書き換え（aws.Config の endpoint 差し替え）
- constitution 改定（v2 → v3、既に実施済み）

対して得られるもの: 150 行の SDK ラッパーのテストカバレッジ。

**Recommendation**: コスト対効果を慎重に評価。アプローチ D（ハイブリッド）なら上記コストの大部分が不要。

## C. Risk & Feasibility

### S7: Risk — moto の API 忠実度

**Category**: Risk & Failure Modes
**Finding**: 001-psm の research.md で LocalStack の挙動差異リスクが指摘されている。moto も同様のリスクがある:
- `CreateSecret` の `ResourceExistsException` 回帰バグ（moto#9700、v5.1.11）
- `GetParametersByPath(Path: "/")` のルートパスバグ（moto#1700、2018 年修正済みだが他のエッジケースの可能性）
- moto のバージョンアップで新たな非互換が入るリスク

moto のテストが通っても実 AWS で失敗するケースが起きると、「テストの信頼性」という本来の目的が損なわれる。

**Recommendation**: moto を入れる場合でも、実 AWS テストは廃止せず並行維持すべき（spec の P2 で対応済み）。
**Evidence**: moto GitHub issues #9700, #1700。

### S8: Feasibility — aws-sdk-go-v2 の endpoint 差し替え

**Category**: Feasibility Verification
**Finding**: `config.WithBaseEndpoint()` による endpoint 差し替えは SDK ドキュメントに記載があるが、moto との組み合わせでの動作は PoC で確認すべき。特に:
- `GetParametersByPath(Path: "/")` が moto で期待通り動くか
- `SecureString` タイプの扱い
- `ResourceExistsException` のエラー型マッチング

**Recommendation**: moto を採用する場合は PoC 必須。
**Evidence**: aws-sdk-go-v2 ドキュメント。PoC なしでは確認不可。

## D. Integration & Governance

### S9: Constitution Impact

**Category**: Constitution Compliance
**Finding**: Constitution v3.0.0 への改定は既に実施済み（Stub tests を許容）。ただし、アプローチ D（ハイブリッド）を採用する場合、SDK interface fake は constitution の「Mocks are not used」の改定前の文言に近い。v3.0.0 は「emulator」を許容したが「interface mock」については明示的に言及していない。

一方で、sync_test.go の `fakeStore` は既に存在し、constitution 違反として問題になっていない。SDK レベルの interface fake も同じカテゴリ。

**Recommendation**: `fakeStore` パターンが既に許容されている以上、SDK interface fake も同じ扱いで問題ない。constitution の追加改定は不要。
**Evidence**: sync_test.go の fakeStore、export_test.go の emptyStore。

### S10: Scope — spec は広すぎないか？

**Category**: Scope Boundaries
**Finding**: spec は「全 19 テストが skip なく完走する」（SC-001）をゴールにしている。これは暗黙に「moto で全 AWS テストを実行する」を意味する。しかし、アプローチ D を採用する場合、ゴールは「CI でアプリケーションロジックの回帰を検出する」に変わり、SC-001 は修正が必要。
**Recommendation**: spec の成功基準を再検討。

## Items Requiring PoC

moto を採用する場合:
1. `GetParametersByPath(Path: "/", Recursive: true)` が moto で全パラメータを返すか
2. `CreateSecret` → `ResourceExistsException` のエラー型が aws-sdk-go-v2 の `errors.As` でマッチするか
3. `config.WithBaseEndpoint()` で moto に接続した SSM/SM クライアントが正常動作するか

アプローチ D を採用する場合:
- PoC 不要（Go の標準テストパターンのみ）

## Constitution Impact

- v3.0.0 への改定は実施済み
- アプローチ D を採用する場合、追加改定は不要（fakeStore パターンの延長）
- アプローチ D を採用する場合、v3.0.0 の stub tests 条項は「将来の拡張時に有効」として維持

## Recommendation

**アプローチ D（ハイブリッド）を推奨。** spec の方向修正を検討すべき。

理由:
1. 150 行の SDK ラッパーに対して Docker エミュレータは過剰
2. fakeStore パターンが既に存在し、その延長で SDK レベルの fake を追加するだけ
3. Docker 依存なし、PoC 不要、メンテナンスコスト最小
4. 実 AWS テストはローカルの手動ゲートとして維持
5. 将来 Store 実装が複雑化した場合に moto 導入を再検討

ただし最終判断はユーザーに委ねる。moto にも「SDK 呼び出し構造の正しさまでテストできる」という利点はある。

次のステップ: spec の方向性について判断を仰ぎ、必要に応じて spec を修正してから `/speckit.plan` へ。
