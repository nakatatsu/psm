# Feature Specification: Issue 番号と SpecKit 機能番号の同期

**Feature Branch**: `016-sync-issue-numbers`
**Created**: 2026-03-17
**Status**: Draft
**Input**: User description: "Issue 番号と SpecKit 機能番号を合わせる方法を整える" (GitHub Issue #16)

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Issue URL から SpecKit 機能を作成する (Priority: P1)

開発者が GitHub Issue を起点に `/speckit.specify` を実行する際、Issue 番号と SpecKit 機能番号が自動的に一致する。現在は `--number` オプションを毎回手動指定する必要があり、忘れると番号がずれてしまう。

**Why this priority**: これが本機能の核心。番号の不一致は、Issue とブランチ・spec の紐付けを困難にし、プロジェクト管理の効率を下げる。

**Independent Test**: `/speckit.specify` に Issue URL を渡した際、生成されるブランチ番号が Issue 番号と一致することを確認する。

**Acceptance Scenarios**:

1. **Given** GitHub Issue #20 が存在する, **When** `/speckit.specify https://github.com/nakatatsu/psm/issues/20` を実行する, **Then** ブランチ `020-<short-name>` が作成され、spec ディレクトリも `specs/020-<short-name>/` となる
2. **Given** GitHub Issue URL が `/speckit.specify` に渡された, **When** ブランチ作成処理が実行される, **Then** `--number` オプションの手動指定なしに Issue 番号が自動的に使用される
3. **Given** GitHub Issue #5 の URL が渡された, **When** ローカルにブランチ `005-other-feature` が既に存在する, **Then** エラーとして処理が中断され、ブランチ・spec は作成されない

---

### User Story 2 - Issue URL なしの実行を拒否する (Priority: P2)

`/speckit.specify` を Issue URL なしで実行した場合、エラーとして処理を中断する。SpecKit を使う規模の作業は必ず Issue を起点とするルールを強制し、番号の衝突を根本的に防ぐ。

**Why this priority**: Issue 必須化により番号体系が一元化され、「自動採番が将来の Issue 番号と衝突する」問題を根本的に排除する。SpecKit を使うレベルの作業であれば Issue を立てるコストは十分低い。

**Independent Test**: Issue URL を渡さずに `/speckit.specify "Some feature"` を実行し、エラーメッセージが表示されて処理が中断されることを確認する。

**Acceptance Scenarios**:

1. **Given** `/speckit.specify` が Issue URL なしで呼ばれた, **When** 処理が開始される, **Then** 「Issue URL が必要です。先に GitHub Issue を作成してください」という旨のエラーメッセージが表示され、ブランチ・spec は作成されない
2. **Given** `--number` オプションのみが指定された（Issue URL なし）, **When** 処理が開始される, **Then** 同様にエラーとなる（`--number` による手動指定も廃止）

---

### Edge Cases

- Issue URL の形式が不正な場合（例: `https://github.com/nakatatsu/psm/pull/10`）、エラーメッセージを表示して処理を中断する
- 既に同じ番号のブランチが存在する場合、エラーとして処理を中断する（ブランチ・spec は作成しない）
- Issue が存在しない番号（例: 削除済み）の URL が渡された場合、GitHub API で存在確認を行い、エラーメッセージを表示して処理を中断する
- Issue URL のリポジトリ部分が現在のリポジトリと異なる場合、警告を表示する

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: `/speckit.specify` に GitHub Issue URL が渡された場合、システムは Issue 番号を自動的に抽出し、SpecKit 機能番号として使用しなければならない
- **FR-002**: `/speckit.specify` は Issue URL を必須とし、Issue URL なしで実行された場合はエラーとして処理を中断しなければならない
- **FR-003**: Issue URL の形式を検証し、不正な形式の場合はわかりやすいエラーメッセージを表示しなければならない
- **FR-004**: 抽出した Issue 番号が既存のブランチ番号と衝突する場合、エラーとして処理を中断しなければならない
- **FR-005**: `/speckit.specify` のスキル定義（プロンプト）において、Issue URL から番号を抽出し `--number` へ自動的に渡すよう手順を更新しなければならない
- **FR-006**: `--number` オプションによる手動番号指定を廃止し、番号の決定は Issue URL からの自動抽出に一元化しなければならない
- **FR-007**: GitHub API を使用して Issue の実在を確認し、存在しない場合はエラーメッセージを表示して処理を中断しなければならない
- **FR-008**: Issue のタイトルと説明文を GitHub API から取得し、short-name の生成および spec の入力情報として使用しなければならない

### Assumptions

- GitHub Issue URL のフォーマットは `https://github.com/{owner}/{repo}/issues/{number}` で安定している
- 番号の同期は「Issue 起点で spec を作る」フローに限定し、spec から逆に Issue を作成するフローは本機能のスコープ外とする
- 番号の「穴」（欠番）は許容する。例: Issue #16 の次に Issue #20 が作られた場合、017〜019 は欠番となるが問題ない
- `create-new-feature.sh` の `--number` オプション自体はスクリプトの内部インターフェースとして残してよい（スキル定義側で Issue URL からの自動抽出を強制する）
- GitHub API へのアクセスは既存の認証基盤（gh-token-sidecar）を利用する

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Issue URL を渡して `/speckit.specify` を実行した際、手動で `--number` を指定せずとも、100% の確率で Issue 番号と一致するブランチが作成される
- **SC-002**: Issue URL なしで `/speckit.specify` を実行した場合、100% の確率でエラーとなりブランチが作成されない
- **SC-003**: 不正な URL や衝突のケースで、ユーザーが次に何をすべきか理解できるエラーメッセージが表示される

## Clarifications

### Session 2026-03-17

- Q: 番号衝突時の動作は警告のみか、エラー中断か？ -> A: エラーとして中断する（ブランチ・spec は作成しない）
- Q: Issue の存在確認は URL フォーマット検証のみか、GitHub API で実在確認するか？ -> A: GitHub API で Issue の存在を明示的に確認する
- Q: short-name の生成元は何を使うか？ -> A: Issue のタイトルと説明文の両方を GitHub API から取得して使用する
