# Tasks: git worktree セットアップ Skill

**Input**: Design documents from `/specs/026-worktree-setup-skill/`
**Prerequisites**: plan.md, spec.md, research.md, quickstart.md, survey.md

**Tests**: テスト不要（Go コードではなく SKILL.md の作成のため）

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2)
- Include exact file paths in descriptions

---

## Phase 1: Setup

**Purpose**: Skill ディレクトリ構造の作成と SKILL.md の骨格

- [x] T001 Create skill directory at `.claude/skills/worktree-setup/`
- [x] T002 Create SKILL.md frontmatter (name, description) in `.claude/skills/worktree-setup/SKILL.md`

---

## Phase 2: Foundational (前提条件セットアップ)

**Purpose**: worktree 作成前に必ず実行される前提条件チェック・自動設定の手順を SKILL.md に記述

**⚠️ CRITICAL**: この手順が US1/US2 両方で共通利用される

- [x] T003 Write prerequisite check section: `git config worktree.useRelativePaths` 確認・設定手順を `.claude/skills/worktree-setup/SKILL.md` に記述
- [x] T004 Write prerequisite check section: `.gitignore` に `.worktrees/` が含まれるか確認・追加手順を `.claude/skills/worktree-setup/SKILL.md` に記述

**Checkpoint**: 前提条件セクション完成

---

## Phase 3: User Story 1 - 新規ブランチで worktree を作成する (Priority: P1) 🎯 MVP

**Goal**: 新規ブランチ名を指定して `.worktrees/<dir>` に worktree を作成し、相対パス変換まで完了する

**Independent Test**: `/worktree-setup feature/new-branch` を実行し、`.worktrees/feature-new-branch/` が作成され、`.git` ファイルが相対パスになっていることを確認

### Implementation for User Story 1

- [x] T005 [US1] Write branch existence check logic (ローカル・リモートのブランチ一覧を参照して新規/既存を判定) in `.claude/skills/worktree-setup/SKILL.md`
- [x] T006 [US1] Write new branch worktree creation section (`git worktree add -b <branch> .worktrees/<dir> origin/main`) in `.claude/skills/worktree-setup/SKILL.md`
- [x] T007 [US1] Write directory name derivation rule (スラッシュをハイフンに置換) in `.claude/skills/worktree-setup/SKILL.md`
- [x] T008 [US1] Write relative path conversion section (`.worktrees/<dir>/.git` と `.git/worktrees/<name>/gitdir` の絶対→相対パス変換) in `.claude/skills/worktree-setup/SKILL.md`
- [x] T009 [US1] Write completion message section (作成されたパスと `cd .worktrees/<dir> && claude` の案内) in `.claude/skills/worktree-setup/SKILL.md`

**Checkpoint**: 新規ブランチでの worktree 作成が一通り動作する

---

## Phase 4: User Story 2 - 既存ブランチで worktree を作成する (Priority: P2)

**Goal**: 既存ブランチ名を指定して `.worktrees/<dir>` に worktree を作成する

**Independent Test**: 既存のリモートブランチ名を指定して `/worktree-setup existing-branch` を実行し、`.worktrees/existing-branch/` が作成されることを確認

### Implementation for User Story 2

- [x] T010 [US2] Write existing branch worktree creation section (`git worktree add .worktrees/<dir> <branch>`) in `.claude/skills/worktree-setup/SKILL.md`
- [x] T011 [US2] Integrate branch existence check (T005) with existing branch path to create unified flow in `.claude/skills/worktree-setup/SKILL.md`

**Checkpoint**: 新規・既存ブランチの両方で worktree 作成が動作する

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: エラーハンドリング、案内の充実

- [x] T012 Write error handling section (worktree 既存、git worktree add 失敗時) in `.claude/skills/worktree-setup/SKILL.md`
- [x] T013 Write worktree deletion as separate skill in `.claude/skills/worktree-remove/SKILL.md`
- [x] T014 Review and finalize complete SKILL.md for consistency and completeness in `.claude/skills/worktree-setup/SKILL.md`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Phase 1 (T001, T002)
- **User Story 1 (Phase 3)**: Depends on Phase 2 (T003, T004)
- **User Story 2 (Phase 4)**: Depends on Phase 3 completion (T005 の判定ロジックを共有)
- **Polish (Phase 5)**: Depends on Phase 4 completion

### User Story Dependencies

- **User Story 1 (P1)**: Phase 2 完了後に開始可能。他ストーリーへの依存なし
- **User Story 2 (P2)**: US1 の判定ロジック (T005) に依存。US1 完了後に開始

### Within Each User Story

- 全タスクが同一ファイル (SKILL.md) への書き込みのため、並列実行不可
- 各タスクは順序通り実行

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001-T002)
2. Complete Phase 2: Foundational (T003-T004)
3. Complete Phase 3: User Story 1 (T005-T009)
4. **STOP and VALIDATE**: 新規ブランチで worktree 作成をテスト
5. 動作確認 OK なら次へ

### Incremental Delivery

1. Setup + Foundational → SKILL.md 骨格完成
2. User Story 1 → 新規ブランチ対応 → 動作確認 (MVP)
3. User Story 2 → 既存ブランチ対応 → 動作確認
4. Polish → エラーハンドリング・削除案内 → 最終確認

---

## Notes

- 全タスクが単一ファイル `.claude/skills/worktree-setup/SKILL.md` への追記・編集
- [P] マークなし（同一ファイルのため並列不可）
- 各 Phase 完了後に手動で動作確認を推奨
