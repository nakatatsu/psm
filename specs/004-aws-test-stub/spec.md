# Feature Specification: AWS テストスタブ基盤

**Feature Branch**: `004-aws-test-stub`
**Created**: 2026-03-15
**Status**: Draft
**Input**: moto（or mock）で AWS 依存テストを実 AWS なしで実行できるようにする

## User Scenarios & Testing _(mandatory)_

### User Story 1 - AWS なしでの統合テスト実行 (Priority: P1)

開発者がローカル環境や CI 環境で、実 AWS アカウントや認証情報なしに SSM Parameter Store および Secrets Manager の統合テストを実行できる。現在 `PSM_INTEGRATION_TEST=1` でのみ実行される AWS 系テストが、スタブ環境に対して常時実行可能になる。

**Why this priority**: 現状テストの約半数（9件/19件）が AWS 認証なしでは skip される。これらをスタブで実行可能にすることが本フィーチャーの核心的価値。

**Independent Test**: AWS 認証情報を一切設定せずに `go test ./...` を実行し、全テストが pass することを確認する。

**Acceptance Scenarios**:

1. **Given** AWS 認証情報が未設定の環境, **When** `go test ./...` を実行する, **Then** SSM および SecretsManager のテストがスタブに対して実行され、skip されない
2. **Given** スタブ環境が起動している, **When** SSM の GetParametersByPath/PutParameter/DeleteParameters を呼び出すテストが実行される, **Then** 期待通りの結果が返り、テストが pass する
3. **Given** スタブ環境が起動している, **When** SecretsManager の ListSecrets/GetSecretValue/CreateSecret/PutSecretValue/DeleteSecret を呼び出すテストが実行される, **Then** 期待通りの結果が返り、テストが pass する

---

### User Story 2 - 実 AWS との切り替え (Priority: P2)

開発者が必要に応じて、同じテストを実 AWS に対しても実行できる。スタブと実 AWS の切り替えが明確な方法で行える。

**Why this priority**: スタブだけでは検出できない問題（API の振る舞いの差異等）を実 AWS で確認する手段を残しておく必要がある。

**Independent Test**: `PSM_INTEGRATION_TEST=1` を設定してテストを実行し、実 AWS に対してテストが走ることを確認する。

**Acceptance Scenarios**:

1. **Given** `PSM_INTEGRATION_TEST=1` が設定されている, **When** テストを実行する, **Then** スタブではなく実 AWS に対してテストが実行される
2. **Given** `PSM_INTEGRATION_TEST` が未設定, **When** テストを実行する, **Then** スタブに対してテストが実行される

---

### Edge Cases

- スタブと実 AWS の API 挙動に差異がある場合、テストはどちらの挙動を正とするか？（実 AWS を正とし、スタブ固有の問題は既知の制限として文書化する）
- スタブ環境（docker-compose サービス）が起動していない場合、テストはどうなるか？（接続エラーで失敗する。docker-compose up で起動を促すメッセージは出さない — テスト失敗自体が十分なシグナル）
- 複数テストが並行実行された場合、スタブのデータが競合しないか？（テストごとに一意のキー prefix を使用する既存のパターンを踏襲する）

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: AWS 認証情報なしの環境で SSM Parameter Store のテスト（GetParametersByPath, PutParameter, DeleteParameters）が実行できなければならない
- **FR-002**: AWS 認証情報なしの環境で Secrets Manager のテスト（ListSecrets, GetSecretValue, CreateSecret, PutSecretValue, DeleteSecret）が実行できなければならない
- **FR-003**: `PSM_INTEGRATION_TEST=1` 設定時は実 AWS に対してテストを実行しなければならない（既存の動作を維持）
- **FR-004**: `PSM_INTEGRATION_TEST` 未設定時はスタブに対してテストを実行しなければならない。スタブの接続先は環境変数（`AWS_ENDPOINT_URL`）で docker-compose から注入する
- **FR-005**: スタブ導入にあたり、プロダクションコードへのテスト分岐（if test then ...）の混入は禁止する。DI により切り替えること。現状 Store インターフェースと `aws.Config` による DI が成立しているため、原則プロダクションコードの変更は不要の見込み。ただしスタブの API 互換性（例: `GetParametersByPath(Path: "/")` の挙動差異等）によってはプロダクションコード側の最小限のリファクタリングが必要になる可能性があり、その場合は実地検証の結果をもって柔軟に対応する
- **FR-006**: スタブ環境は追加の外部サービス契約や API キーなしで利用できなければならない
- **FR-007**: テストの実行結果（pass/fail）は実 AWS での実行結果と同等の信頼性を持たなければならない（スタブが対応する API 範囲において）

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: AWS 認証なしの環境で `go test ./...` を実行し、全テスト（現在19件）が skip なく完走する
- **SC-002**: スタブ導入後も `PSM_INTEGRATION_TEST=1` での実 AWS テストが従来通り pass する
- **SC-003**: テスト全体の実行時間がスタブ使用時に 30 秒以内で完了する

## Clarifications

### Session 2026-03-15

- Q: スタブ環境のライフサイクル管理をどのレイヤーで行うか？ -> A: DevContainer の docker-compose にサービスとして常駐させる。テスト実行時には起動済みを前提とする
- Q: CI 環境でのスタブ起動方法は？ -> A: GHA でも同じ docker-compose.yml を使って起動する。GHA services ディレクティブ等の環境固有の方法は使わず、環境差異を排除する
- Q: スタブの接続先をテストに伝える方法は？ -> A: docker-compose.yml の workspace サービスに環境変数（例: `AWS_ENDPOINT_URL`）を設定する。テストコードは `os.Getenv` で取得する

## Assumptions

- スタブの第一候補は moto（Python、Docker サーバーモード）。moto が技術的に不適合な場合は Go interface mock にフォールバックする
- moto 使用時は DevContainer の docker-compose に service として追加し、テストコードからは endpoint URL の差し替えで接続する
- 既存テストのテストロジック（テーブル駆動パターン、キー prefix によるテストデータ分離）はそのまま活用する
- テスト用の AWS 認証情報（ダミー）はテストコード内でハードコードする（`AWS_ACCESS_KEY_ID=testing` 等）
