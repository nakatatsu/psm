# Feature Specification: psm — SOPS-to-AWS Parameter Sync CLI

**Feature Branch**: `001-psm`
**Created**: 2026-03-08
**Status**: Draft
**Input**: User description: "SOPS復号済みYAMLファイルを読み、AWS SSM Parameter Store または Secrets Manager に同期するCLIツール psm を作る。Go言語。"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - SSM Parameter Store への同期 (Priority: P1)

運用者として、SOPS で復号済みの YAML ファイルに記載されたキーと値を、AWS SSM Parameter Store に一括で追加・更新（upsert）したい。YAML のキーがそのまま SSM パラメータ名になる。これにより、手動で Parameter Store を操作する手間とヒューマンエラーを排除できる。

**Why this priority**: Parameter Store への upsert はツールの最も基本的な価値提供であり、これだけで運用に使える MVP となる。

**Independent Test**: 復号済み YAML ファイルを用意し、`psm sync --store ssm secrets.yaml` を実行して、Parameter Store にキーと値が正しく同期されることを確認する。

**Acceptance Scenarios**:

1. **Given** 3 つの key-value を含む復号済み YAML ファイルがある, **When** `psm sync --store ssm secrets.yaml` を実行する, **Then** YAML のキー名がそのまま SSM パラメータ名として SecureString 型で 3 件 upsert される
2. **Given** Parameter Store に既存の値がある, **When** YAML 内の同名キーの値を変更して実行する, **Then** 既存パラメータが新しい値で上書きされる
3. **Given** YAML ファイル内に `sops` キーが残っている, **When** 実行する, **Then** `sops` は無視され、それ以外のキーのみ同期される
4. **Given** 3 件中 1 件の書き込みが失敗する, **When** 実行する, **Then** 残り 2 件は正常に処理され、終了コード 1 が返る

---

### User Story 2 - Secrets Manager への同期 (Priority: P2)

運用者として、同じ YAML ファイルの内容を Secrets Manager にも同期したい。ストアの選択は CLI フラグ `--store sm` で切り替えるだけで、ワークフローは SSM と同一にしたい。

**Why this priority**: SSM と並ぶもう一つの主要ストア。SSM の基盤（YAML パース、差分計算）を再利用でき、P1 完了後に効率的に実装できる。

**Independent Test**: `psm sync --store sm secrets.yaml` を実行して、Secrets Manager に正しい名前・値でシークレットが作成/更新されることを確認する。

**Acceptance Scenarios**:

1. **Given** 復号済み YAML ファイルがある, **When** `psm sync --store sm secrets.yaml` を実行する, **Then** YAML のキー名がそのまま Secrets Manager のシークレット名として作成される
2. **Given** 同名のシークレットが既に存在する, **When** 実行する, **Then** 既存シークレットの値が更新される

---

### User Story 3 - Dry-run で差分を事前確認 (Priority: P3)

運用者として、実際に AWS を変更する前に、何が追加・更新・削除されるかを確認したい。

**Why this priority**: 本番環境での誤操作防止に不可欠だが、同期機能が先に動いている必要がある。

**Independent Test**: `psm sync --store ssm --dry-run secrets.yaml` を実行し、標準出力に差分（create / update / delete のみ、no-change は非表示）が表示され、AWS への書き込みが一切行われないことを確認する。

**Acceptance Scenarios**:

1. **Given** AWS に既存キーがあり YAML に新しいキーが追加されている, **When** `--dry-run` で実行する, **Then** 新キーは `create` と表示され、値が変わっていない既存キーは表示されず、AWS への変更は行われない
2. **Given** `--prune` と `--dry-run` を同時指定, **When** 実行する, **Then** 削除予定のキーが `delete` と表示されるが実際には削除されない

---

### User Story 4 - Prune で不要キーを削除 (Priority: P4)

運用者として、YAML ファイルから削除したキーを AWS 側からも自動的に削除したい。YAML ファイルがそのアカウントの全パラメータ/シークレットの唯一の真実（single source of truth）である。

**Why this priority**: 完全同期の実現に必要だが、デフォルト動作（upsert のみ）より危険性が高く、明示的なフラグ指定を要するため優先度は低め。

**Independent Test**: YAML からキーを削除した状態で `psm sync --store ssm --prune secrets.yaml` を実行し、AWS 側の該当パラメータが削除されることを確認する。

**Acceptance Scenarios**:

1. **Given** AWS にパラメータが存在するが YAML にはない, **When** `--prune` 付きで実行する, **Then** そのパラメータが AWS から削除される
2. **Given** `--store sm --prune` を指定, **When** 実行する, **Then** Secrets Manager 側の該当シークレットが即時削除される（復旧期間なし）
3. **Given** `--prune` なしで実行, **When** YAML にないキーが AWS に存在する, **Then** 何も削除されない
4. **Given** `--prune` を指定, **When** 実行する, **Then** アカウント内の全パラメータ/シークレットが対象となる（プレフィックスによるスコープ制限はない）

---

### User Story 5 - Export で既存パラメータを YAML に書き出し (Priority: P5)

運用者として、AWS 上の既存パラメータ/シークレットを psm の YAML 形式でエクスポートしたい。初回導入時や、手動で追加されたパラメータを YAML 管理下に取り込む際に使用する。

**Why this priority**: 同期機能（P1-P4）が基本。エクスポートは初回セットアップの利便性向上であり、同期機能なしでは意味がない。

**Independent Test**: `psm export --store ssm output.yaml` を実行し、AWS 上の全パラメータが key-value の YAML ファイルとして出力されることを確認する。

**Acceptance Scenarios**:

1. **Given** SSM Parameter Store に 3 つのパラメータがある, **When** `psm export --store ssm output.yaml` を実行する, **Then** 3 つの key-value を含む YAML ファイルが生成される
2. **Given** 出力先ファイルが既に存在する, **When** エクスポートを実行する, **Then** 上書きせずエラーメッセージを表示し終了コード 1
3. **Given** Secrets Manager にシークレットがある, **When** `psm export --store sm output.yaml` を実行する, **Then** 全シークレットを含む YAML ファイルが生成される
4. **Given** AWS 上にパラメータ/シークレットが 0 件, **When** エクスポートを実行する, **Then** エラーメッセージを表示し終了コード 1

---

### Edge Cases

入力バリデーションエラー（重複キー、null 値、空キー、不正な型等）は FR-020 で一括定義。以下はそれ以外のエッジケース:

- YAML の値が空文字列の場合: 空文字列として正常に同期する
- 指定されたファイルが存在しない場合: エラーメッセージを表示し終了コード 1
- AWS 認証情報が未設定の場合: AWS SDK のエラーメッセージがそのまま表示される
- 複数ファイル引数: サポートしない。引数は 1 ファイルのみ。複数指定時はエラー（FR-012 に包含）
- AWS API エラー（認証期限切れ、権限不足等）が発生した場合: FR-010 に従い個別キーの失敗として扱い、残りの処理を続行する。GetAll 自体が失敗した場合は致命的エラーとし、即座に終了コード 1 を返す
- AWS が拒否するキー名（名前規則違反等）: AWS API エラーとして FR-010/FR-019 に従い処理される（psm 側でのキー名バリデーションは行わない）
- AWS の値サイズ制限（SSM: 8KB advanced/4KB standard、SM: 64KB）超過: AWS API エラーとして FR-010/FR-019 に従い処理される

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: ツールは 2 つのサブコマンドを持つ。`psm sync [flags] <file>` で同期、`psm export [flags] <file>` でエクスポート。サブコマンドは必須（省略不可）
- **FR-002**: `--store <ssm|sm>` フラグ（必須）でストア種別を指定する。`--profile <name>` フラグ（任意）で AWS プロファイルを指定する。`--profile` 未指定時は IAM Role 等の SDK デフォルト認証を使用する。環境変数 `AWS_PROFILE` は常に無視しなければならない。`--store` と `--profile` は両サブコマンド共通のフラグである
- **FR-003**: YAML のキー名がそのまま AWS 上のパラメータ名/シークレット名になる。パス構造の組み立てはツール側で行わない
- **FR-004**: SSM 同期時、型は SecureString、Overwrite=true で upsert しなければならない
- **FR-005**: SM 同期時、存在しなければ CreateSecret、存在すれば PutSecretValue で更新しなければならない
- **FR-006**: YAML 内の `sops` キー（SOPS メタデータ）は同期対象から除外しなければならない。予約キーの除外は FR-020 のバリデーションより先に実行する（`sops` の値はマップであり、除外前にバリデーションすると誤検出する）
- **FR-007**: `--dry-run` フラグ指定時は、差分（create / update / delete）を標準出力に表示し、AWS への変更を一切行ってはならない。変更のないキー（no-change）は表示しない
- **FR-008**: `--prune` フラグ指定時は、アカウント内の全パラメータ/シークレットのうち YAML に存在しないものを AWS から削除しなければならない。プレフィックスによるスコープ制限はない
- **FR-009**: `--prune` 未指定時は、追加・更新のみ行い、削除は行わない
- **FR-010**: 個々のキーの同期失敗時も残りのキーの処理を続行しなければならない
- **FR-011**: 全件成功時は終了コード 0、1 件以上失敗時は終了コード 1 を返さなければならない
- **FR-012**: サブコマンドが未指定、不正、または各サブコマンドで必須引数が不足している場合、使い方を表示し終了コード 1 を返さなければならない
- **FR-013**: YAML の値はすべてスカラー値（string / int / bool / float）でなければならない。マップ・配列はエラーとする
- **FR-014**: 通常実行時も差分（create / update / delete）を 1 行/キーで標準出力に表示しなければならない（例: `create: /myapp/prod/DB_URL`）。変更のないキー（no-change）は表示してはならない
- **FR-015**: 出力形式は通常実行と `--dry-run` で同一でなければならない
- **FR-016**: YAML の非文字列値（整数、ブール値等）は文字列に変換して同期しなければならない（例: `5432` → `"5432"`, `true` → `"true"`）
- **FR-017**: 同期前に AWS 上の既存値を取得し、YAML の値と比較しなければならない。値が一致するキーは操作も表示もしない（no-change）。値が異なるキーのみ `update` として処理・表示する
- **FR-018**: 全キーの処理完了後、サマリー行を標準出力に表示しなければならない（例: `2 created, 1 updated, 0 deleted, 5 unchanged, 0 failed`）
- **FR-019**: 個々のキーの同期エラーは stderr にキー名とエラーメッセージを出力しなければならない（例: `error: /myapp/prod/API_KEY: AccessDeniedException: ...`）。値は絶対に出力してはならない。すべての失敗キーが個別に表示されなければならない。stdout には差分情報とサマリーのみを出力する
- **FR-020**: YAML ファイルの読み込み直後（FR-006 の予約キー除外後）、AWS への通信前に入力バリデーションを一括実行しなければならない。以下のいずれかに該当する場合は、何が問題かをユーザーが理解できる具体的なエラーメッセージを表示し終了コード 1 を返す。AWS への通信は一切行わない:
  - キーが 0 件（`sops` を除いた後）
  - キーの重複がある（重複したキー名を表示）
  - 値がマップ、配列、または `null`（該当キー名を表示）
  - キーが空文字列
- **FR-021**: `psm export` は AWS 上の全パラメータ/シークレットを取得し、key-value のみの YAML ファイルとして書き出さなければならない。YAML にメタデータ（ストア種別等）は含めない
- **FR-022**: 出力先ファイルが既に存在する場合、上書きせずエラーメッセージを表示し終了コード 1 を返さなければならない
- **FR-023**: AWS 上にパラメータ/シークレットが 0 件の場合、エラーメッセージを表示し終了コード 1 を返さなければならない
- **FR-024**: `--help` フラグまたはサブコマンドに `--help` を指定した場合、使い方を標準出力に表示し終了コード 0 を返さなければならない。Go の `flag` パッケージのデフォルト動作に従う
- **FR-025**: SM の Delete は `ForceDeleteWithoutRecovery=true` で即時削除しなければならない。復旧期間を設けない
- **FR-026**: 予約キーは `sops` のみ。他のキープレフィックスは予約しない

### Key Entities

エンティティの詳細定義は [data-model.md](data-model.md) を参照。Config, Entry, Action, Summary の 4 エンティティで構成される。

### Assumptions

- 入力 YAML は SOPS によって事前に復号されている（psm は復号処理を行わない）
- AWS リージョンは環境変数（`AWS_REGION` / `AWS_DEFAULT_REGION`）または AWS SDK のデフォルト設定に従う
- AWS 認証情報・プロファイル選択は FR-002 に従う
- YAML ファイルはそのアカウント内の全パラメータ/シークレットの single source of truth である

## Clarifications

### Session 2026-03-08

- Q: 通常実行時（非 dry-run）の出力形式は？ -> A: 差分のみ 1 行/キーで表示（create / update / delete）。no-change は絶対に表示しない。何が変わるか正確にわかることが必須。
- Q: YAML の非文字列値（整数、ブール値等）の扱いは？ -> A: すべて文字列に自動変換して同期する。
- Q: update 判定のために既存値との比較は行うか？ -> A: 必須。同期前に既存値を取得し比較する。値が同一なら操作も表示もしない。比較なしの全上書きは許されない。
- Q: 処理完了後にサマリー行を表示するか？ -> A: 必須。差分行の後にサマリー 1 行を表示する。
- Q: エラー発生キーの出力先は？ -> A: stderr。stdout は差分+サマリーのみ。
- Q: --app/--env フラグは必要か？ -> A: 不要。YAML のキーがそのまま AWS リソース名。ストア種別・プロファイルは CLI 共通フラグ `--store` / `--profile` で指定（環境変数 `AWS_PROFILE` は常に無視）。
- Q: prune のスコープは？ -> A: アカウント内の全パラメータ/シークレット。プレフィックスによるセグメント分割は事故の元なので禁止。

### Session 2026-03-09

- Q: diff 出力行のソート順は？ -> A: 未定義。AWS API の返却順に依存する。アルファベット順の保証は不要。
- Q: 使い方メッセージ（FR-012/FR-024）の内容は？ -> A: Go の `flag` パッケージのデフォルト出力に従う。カスタムフォーマットは不要。
- Q: `--profile` に無効な値が渡された場合は？ -> A: AWS SDK がエラーを返すのでそのまま表示。psm 側でのバリデーションは不要。
- Q: SM で最近削除されたシークレットと同名の CreateSecret は？ -> A: `ForceDeleteWithoutRecovery=true` で即時削除するため復旧期間のコンフリクトは発生しない。新規 CreateSecret が AWS に拒否された場合は FR-010/FR-019 に従う。
- Q: 並行数 10 は設定可能か？ -> A: 固定値。設定可能にはしない（YAGNI）。実装詳細であり FR には含めない。
- Q: 100 件を超えるキーの上限は？ -> A: 上限は設けない。AWS API の rate limit に依存する。SC-001 は動作確認の基準値であり制限ではない。
- Q: タイムアウトは？ -> A: AWS SDK のデフォルトタイムアウトに委ねる。psm 側でのタイムアウト設定は不要。
- Q: SOPS が不正な YAML を出力した場合は？ -> A: YAML パースエラーとしてエラーメッセージを表示し終了コード 1。psm は入力が有効な YAML であることを前提とするが、パース失敗時は適切にエラーを返す。
- Q: `GetParametersByPath` で path=`/` は全パラメータを返すか？ -> A: はい。Recursive=true と合わせて全パラメータを返す。AWS ドキュメントで確認済み（research.md R2 参照）。

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100 件の key-value を含む YAML ファイルを 1 回のコマンド実行で SSM/SM に同期できる
- **SC-002**: `--dry-run` で表示される差分と、実際に実行した場合の結果が一致する
- **SC-003**: 1 件の失敗が発生しても、残りのキーはすべて正常に処理される
- **SC-004**: 運用者が初めて使う場合でも、ヘルプメッセージ（`--help`）だけで正しく実行できる
- **SC-005**: `--prune` による削除はアカウント内の全パラメータ/シークレットを対象とし、YAML にないものすべてを削除する
- **SC-006**: `psm export --store ssm` で生成した YAML をそのまま `psm sync --store ssm` に渡すと、差分 0 件（全件 unchanged）で同期が完了する
