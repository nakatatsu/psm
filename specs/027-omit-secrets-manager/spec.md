# Feature Specification: Secrets Manager 対応オミット

**Feature Branch**: `027-omit-secrets-manager`
**Created**: 2026-03-18
**Status**: Draft
**Input**: GitHub Issue #27 — Secrets Manager (`--store sm`) は現状使う予定がないため、メンテナンスコストを減らすためにオミットする

## User Scenarios & Testing *(mandatory)*

### User Story 1 - SSM ユーザーの既存ワークフロー維持 (Priority: P1)

psm を SSM Parameter Store に対して使っている開発者は、Secrets Manager コードの削除後もこれまでと全く同じコマンドで同じ結果を得られる。

**Why this priority**: SSM は psm の主要かつ唯一のアクティブユースケースであり、既存ユーザーの体験を損なわないことが最優先。

**Independent Test**: `psm sync --store ssm` による作成・更新・削除・dry-run・承認フローがすべて既存と同じ結果になることを確認する。

**Acceptance Scenarios**:

1. **Given** SM コードが削除された状態、**When** `psm sync --store ssm secrets.yaml` を実行、**Then** 従来通りパラメータの作成・更新が行われる
2. **Given** SM コードが削除された状態、**When** `psm sync --store ssm --dry-run secrets.yaml` を実行、**Then** dry-run 結果が従来通り表示される
3. **Given** SM コードが削除された状態、**When** `psm export --store ssm output.yaml` を実行、**Then** 従来通りエクスポートされる

---

### User Story 2 - SM 指定時の明確なエラー (Priority: P2)

開発者が誤って `--store sm` を指定した場合、何が起きたか即座に理解でき、対処方法がわかるエラーメッセージが表示される。

**Why this priority**: SM を削除した以上、過去の利用者や古いドキュメントを参照した人が混乱しないよう、明確なガイダンスが必要。

**Independent Test**: `--store sm` を指定したときに適切なエラーメッセージが表示され、プログラムが正常に終了することを確認する。

**Acceptance Scenarios**:

1. **Given** SM コードが削除された状態、**When** `psm sync --store sm secrets.yaml` を実行、**Then** 「`--store sm` は現在サポートされていません。`--store ssm` を使用してください」に相当するエラーが表示され、終了コード 1 で終了する
2. **Given** SM コードが削除された状態、**When** `psm export --store sm output.yaml` を実行、**Then** 同様のエラーが表示され、終了コード 1 で終了する

---

### User Story 3 - ドキュメントの一貫性 (Priority: P3)

README やプロジェクト設定を読む開発者は、SM に関する記述がなく、SSM のみがサポートされていることを一貫して理解できる。

**Why this priority**: ドキュメントとコードの乖離はメンテナンスコストと混乱の原因になるが、機能そのものには影響しない。

**Independent Test**: すべてのユーザー向けドキュメントで `sm` や Secrets Manager への言及が除去されていることを確認する。

**Acceptance Scenarios**:

1. **Given** 変更完了後の状態、**When** README.md を確認、**Then** `--store <ssm|sm>` が `--store ssm` に更新され、Secrets Manager の説明が除去されている
2. **Given** 変更完了後の状態、**When** example/README.md を確認、**Then** CLI Reference から `sm` が除去されている
3. **Given** 変更完了後の状態、**When** CLAUDE.md を確認、**Then** Active Technologies から `AWS Secrets Manager` が除去されている

---

### Edge Cases

- `--store` に `sm` でも `ssm` でもない値を指定した場合は、従来通りバリデーションエラーが表示される
- Store interface は残すため、将来の拡張でコンパイルが壊れないことを確認する
- spec ディレクトリ (`specs/001-psm/`) 内の SM 関連記述は設計記録として残すため、削除対象外

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: `sm.go`（SMStore 実装）を削除しなければならない
- **FR-002**: `sm_test.go`（SMStore テスト）を削除しなければならない
- **FR-003**: `main.go` の `--store` バリデーションを `ssm` のみに変更しなければならない
- **FR-004**: `--store sm` を指定した場合、サポート外であることを示すエラーメッセージを返さなければならない
- **FR-005**: `README.md` と `README.ja.md` の `--store <ssm|sm>` を `--store ssm` に更新し、SM 関連の記述を削除しなければならない
- **FR-006**: `example/README.md` の CLI Reference から `sm` を削除しなければならない
- **FR-007**: `CLAUDE.md` の Active Technologies から `AWS Secrets Manager` を削除しなければならない
- **FR-008**: Store interface は削除せず保持しなければならない（将来の拡張ポイント）
- **FR-009**: `specs/001-psm/` 配下の SM 関連記述は設計記録として保持しなければならない
- **FR-010**: 既存の SSM 関連の全テスト（`go test ./...`）がパスしなければならない

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: `go test ./...` が全件パスする
- **SC-002**: `go build ./...` がエラーなくビルドできる
- **SC-003**: `--store sm` 指定時にエラーメッセージが表示され終了コード 1 で終了する
- **SC-004**: ユーザー向けドキュメント（README.md, README.ja.md, example/README.md, CLAUDE.md）に Secrets Manager / `sm` への言及がない

### Assumptions

- Secrets Manager を利用しているユーザーは現時点でいない
- 将来再対応する場合は git history から `sm.go` を復元する想定
- Store interface は他のストアバックエンド（例: HashiCorp Vault）の追加にも対応できる拡張ポイントとして価値がある
