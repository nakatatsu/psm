# Feature Specification: --store フラグの除去

**Feature Branch**: `037-remove-store-flag`
**Created**: 2026-03-27
**Status**: Draft
**Input**: GitHub Issue #37 — SSM のみがサポート対象であるため、`--store` フラグ自体を除去して CLI をシンプルにする

## User Scenarios & Testing *(mandatory)*

### User Story 1 - フラグなしでの sync/export 実行 (Priority: P1)

開発者は `--store ssm` を指定せずに `psm sync` や `psm export` を実行でき、従来通り SSM Parameter Store に対して操作が行われる。

**Why this priority**: ユーザー体験の改善が本機能の核心であり、毎回の冗長な `--store ssm` 指定が不要になることが最大の価値。

**Independent Test**: `psm sync secrets.yaml` および `psm export output.yaml` がフラグなしで正常に動作し、SSM に対して操作が行われることを確認する。

**Acceptance Scenarios**:

1. **Given** `--store` フラグが除去された状態、**When** `psm sync secrets.yaml` を実行、**Then** SSM Parameter Store に対してパラメータの作成・更新が行われる
2. **Given** `--store` フラグが除去された状態、**When** `psm sync --dry-run secrets.yaml` を実行、**Then** dry-run 結果が表示される
3. **Given** `--store` フラグが除去された状態、**When** `psm export output.yaml` を実行、**Then** SSM からパラメータがエクスポートされる

---

### User Story 2 - --store 指定時の明確なエラー (Priority: P2)

過去のドキュメントやスクリプトを参照して `--store ssm` を指定したユーザーに対し、フラグが廃止されたことを伝えるエラーメッセージを表示する。

**Why this priority**: 移行期の混乱を防ぐために、旧オプション使用時のガイダンスが重要。

**Independent Test**: `--store` を指定したときにエラーメッセージが表示され、正しい使い方が案内されることを確認する。

**Acceptance Scenarios**:

1. **Given** `--store` フラグが除去された状態、**When** `psm sync --store ssm secrets.yaml` を実行、**Then** 「`--store` フラグは廃止されました。SSM がデフォルトで使用されます。`--store` を省略してください」に相当するエラーが表示され、終了コード 1 で終了する
2. **Given** `--store` フラグが除去された状態、**When** `psm export --store sm output.yaml` を実行、**Then** 同様のエラーが表示され、終了コード 1 で終了する

---

### User Story 3 - ドキュメントの一貫性 (Priority: P3)

README やサンプルを読む開発者は、`--store` フラグに関する記述がなく、シンプルなコマンド例のみが記載されていることを確認できる。

**Why this priority**: ドキュメントとコードの乖離は混乱の原因になるが、機能そのものには影響しない。

**Independent Test**: ユーザー向けドキュメントで `--store` フラグへの言及が除去されていることを確認する。

**Acceptance Scenarios**:

1. **Given** 変更完了後の状態、**When** README.md / README.ja.md を確認、**Then** `--store ssm` の指定が除去され、シンプルなコマンド例に更新されている
2. **Given** 変更完了後の状態、**When** example/README.md を確認、**Then** CLI Reference から `--store` が除去されている
3. **Given** 変更完了後の状態、**When** テストコードを確認、**Then** `--store` に依存するテストが更新されている

---

### Edge Cases

- `--store` に任意の値（`ssm`, `sm`, その他）を指定した場合、すべて同じ「フラグ廃止」エラーが表示される
- `--store` と他のフラグを組み合わせた場合でも、`--store` のエラーが優先的に表示される
- SM タイポチェック（`--store sm` → SSM の示唆）のバリデーションロジックも同時に除去される

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: `main.go` から `--store` フラグの定義・パースを除去しなければならない
- **FR-002**: `--store` が指定された場合、フラグ廃止を示すエラーメッセージを返さなければならない
- **FR-003**: SSM をデフォルトかつ唯一のストアとしてハードコードしなければならない
- **FR-004**: `--store sm` タイポチェックのバリデーションロジックを除去しなければならない
- **FR-005**: `README.md` および `README.ja.md` のコマンド例から `--store ssm` を除去しなければならない
- **FR-006**: `example/README.md` の CLI Reference から `--store` を除去しなければならない
- **FR-007**: テストコードを更新し、`--store` フラグを使用していたテストを修正しなければならない
- **FR-008**: 既存の全テスト（`go test ./...`）がパスしなければならない
- **FR-009**: `specs/` 配下の過去の仕様記述は設計記録として保持しなければならない

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: `go test ./...` が全件パスする
- **SC-002**: `go build ./...` がエラーなくビルドできる
- **SC-003**: `psm sync secrets.yaml` が `--store` 指定なしで SSM に対して正常動作する
- **SC-004**: `--store` を指定した場合にエラーメッセージが表示され終了コード 1 で終了する
- **SC-005**: ユーザー向けドキュメント（README.md, README.ja.md, example/README.md）から `--store` フラグの記述が除去されている

### Assumptions

- #27 により SM 対応は既に除去済みで、SSM のみがサポートされている
- 将来別のストアが必要になった場合は、新たなフラグ設計で対応する（YAGNI）
- Store interface は既に保持されているため、本変更では触れない
