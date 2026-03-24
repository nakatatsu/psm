# Feature Specification: git worktree セットアップ Skill

**Feature Branch**: `026-worktree-setup-skill`
**Created**: 2026-03-23
**Status**: Draft
**Input**: GitHub Issue #26

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 新規ブランチで worktree を作成する (Priority: P1)

開発者が新しいフィーチャーブランチで並行作業を始めたいとき、Skill をワンコマンドで呼び出すだけで、worktree の作成・前提条件の設定がすべて自動で行われる。手動でのステップ漏れやミスがなくなる。

**Why this priority**: 最も頻度の高いユースケース。新機能開発は常に新規ブランチから始まるため、これが動けば MVP として成立する。

**Independent Test**: Skill を呼び出してブランチ名を指定し、`.worktrees/` 配下に worktree が作成されることを確認する。

**Acceptance Scenarios**:

1. **Given** リポジトリルートにいる開発者が新規ブランチで作業を始めたい, **When** Skill を呼び出してブランチ名を指定する, **Then** `.worktrees/<dir>/` に worktree が作成され、指定ブランチにチェックアウトされている
2. **Given** `worktree.useRelativePaths` が未設定のリポジトリ, **When** Skill を実行する, **Then** `worktree.useRelativePaths` が `true` に自動設定される
3. **Given** `.gitignore` に `.worktrees/` が含まれていない, **When** Skill を実行する, **Then** `.gitignore` に `.worktrees/` が自動追加される

---

### User Story 2 - 既存ブランチで worktree を作成する (Priority: P2)

開発者がリモートまたはローカルに既に存在するブランチに対して worktree を作成したいとき、Skill が既存ブランチを検出して適切なコマンドで worktree をセットアップする。

**Why this priority**: 既存ブランチでの作業再開やレビュー用途で使うケース。新規ブランチに比べ頻度は低いが重要。

**Independent Test**: 既存のブランチ名を指定して Skill を呼び出し、そのブランチの worktree が作成されることを確認する。

**Acceptance Scenarios**:

1. **Given** リモートに既存ブランチがある, **When** Skill にそのブランチ名を指定する, **Then** `.worktrees/<dir>/` にそのブランチの worktree が作成される
2. **Given** ローカルに既存ブランチがある, **When** Skill にそのブランチ名を指定する, **Then** `.worktrees/<dir>/` にそのブランチの worktree が作成される

---

### Edge Cases

- 指定されたブランチ名で worktree が既に存在する場合はどうなるか？ → エラーメッセージで既存 worktree の存在を通知し、中断する
- `.gitignore` ファイルが存在しない場合はどうなるか？ → 新規作成して `.worktrees/` を追加する
- ブランチ名にスラッシュ (`feature/example`) が含まれる場合のディレクトリ名はどうなるか？ → スラッシュをハイフンに置換してディレクトリ名とする（例: `feature/example` → `feature-example`）
- 相対パス化した worktree を `git worktree remove` で削除しようとするとどうなるか？ → Git 2.39 では絶対パスを要求するため失敗する。削除手順を案内に含める（パスを一旦絶対に戻してから remove するか、手動削除 + `git worktree prune`）

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Skill は新規ブランチ指定時に `git worktree add -b <branch> .worktrees/<dir> origin/main` を実行して worktree を作成しなければならない
- **FR-002**: Skill は既存ブランチ指定時に `git worktree add .worktrees/<dir> <branch>` を実行して worktree を作成しなければならない
- **FR-003**: Skill は worktree 作成前に `git config worktree.useRelativePaths true` が設定されていることを確認し、未設定の場合は自動設定しなければならない
- **FR-004**: Skill は worktree 作成前に `.gitignore` に `.worktrees/` が含まれていることを確認し、含まれていない場合は追加しなければならない
- **FR-005**: Skill は `.worktrees/` ディレクトリをリポジトリルート直下に配置しなければならない
- **FR-006**: Skill は Claude Code の Skill フォーマット（`.claude/skills/<name>/SKILL.md`）に準拠しなければならない
- **FR-007**: Skill は worktree 作成完了後、作成されたディレクトリのパスと次のステップ（`cd` して `claude` を起動する手順）をユーザーに案内しなければならない
- **FR-008**: Skill は worktree 作成後、`.worktrees/<dir>/.git` と `.git/worktrees/<name>/gitdir` の絶対パスを相対パスに変換しなければならない（Git 2.39 では `worktree.useRelativePaths` が機能しないための workaround）

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 開発者が Skill を呼び出してから worktree が利用可能になるまでのステップが 1 回の Skill 呼び出しで完結する
- **SC-002**: 前提条件（`worktree.useRelativePaths` 設定、`.gitignore` エントリ）のセットアップ漏れが 0 件になる
- **SC-003**: 新規ブランチ・既存ブランチの両方のケースで正しく worktree が作成される

## Assumptions

- Git 2.39+ がインストール済み（`worktree.useRelativePaths` は Git 2.48+ でのみネイティブ動作するため、2.39 環境では手動相対パス変換 workaround を使用する）
- リポジトリは DevContainer 内で利用され、ホストと DevContainer の双方で git 操作を行う（そのため相対パスが必須）
- ベースブランチは `origin/main` とする
- ブランチが新規か既存かは Skill 側で自動判定する（ローカル・リモートのブランチ一覧を参照）
