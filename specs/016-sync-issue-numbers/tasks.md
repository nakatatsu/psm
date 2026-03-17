# Tasks: Issue 番号と SpecKit 機能番号の同期

**Input**: Design documents from `/specs/016-sync-issue-numbers/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2)
- Include exact file paths in descriptions

## Phase 1: Setup

**Purpose**: 現行の SKILL.md を理解し、変更計画を確認する

- [x] T001 Read current skill definition at .claude/skills/speckit.specify/SKILL.md and identify all sections that reference branch numbering, `--number`, and auto-numbering logic

---

## Phase 2: User Story 1 - Issue URL から SpecKit 機能を作成する (Priority: P1) 🎯 MVP

**Goal**: Issue URL を渡すと Issue 番号が自動的に SpecKit 機能番号になる

**Independent Test**: `/speckit.specify https://github.com/nakatatsu/psm/issues/{N}` を実行し、ブランチ番号が Issue 番号と一致することを確認

### Implementation for User Story 1

- [x] T002 [US1] Add Issue URL detection and validation step before step 1 "Generate a concise short name" in .claude/skills/speckit.specify/SKILL.md — detect `https://github.com/{owner}/{repo}/issues/{number}` pattern, extract issue number, validate URL format (FR-001, FR-003)
- [x] T003 [US1] Add GitHub API validation step using `gh issue view {number} --repo {owner}/{repo} --json title,body` to verify Issue exists and retrieve title/body in .claude/skills/speckit.specify/SKILL.md (FR-007)
- [x] T004 [US1] Replace auto-numbering logic (steps 2b, 2c) with Issue-number-based `--number` in .claude/skills/speckit.specify/SKILL.md — remove "Find the highest feature number" and "Extract all numbers" steps, use extracted Issue number directly (FR-005, FR-006)
- [x] T005 [US1] Update short-name generation step to use Issue title and body instead of user's feature description in .claude/skills/speckit.specify/SKILL.md — title as primary source, body as supplementary context (FR-008)
- [x] T006 [US1] Add branch collision check using `git branch --list` and `git ls-remote --heads origin` before calling create-new-feature.sh, with error message on collision in .claude/skills/speckit.specify/SKILL.md (FR-004)

**Checkpoint**: Issue URL を渡して `/speckit.specify` を実行し、Issue 番号と一致するブランチが作成されることを確認

---

## Phase 3: User Story 2 - Issue URL なしの実行を拒否する (Priority: P2)

**Goal**: Issue URL なしの実行をエラーにして番号衝突を根本的に防ぐ

**Independent Test**: Issue URL なしで `/speckit.specify "Some feature"` を実行し、エラーが表示されることを確認

### Implementation for User Story 2

- [x] T007 [US2] Add Issue URL requirement check at the beginning of the Outline section in .claude/skills/speckit.specify/SKILL.md — if no GitHub Issue URL detected in user input, output error "GitHub Issue URL が必要です。先に Issue を作成してください" and stop (FR-002)
- [x] T008 [US2] Remove `--number` manual specification instructions from user-facing steps in .claude/skills/speckit.specify/SKILL.md — keep `--number` only as internal parameter passed automatically from Issue URL extraction (FR-006, Research R3)

**Checkpoint**: Issue URL なしで実行するとエラーで中断し、ブランチが作成されないことを確認

---

## Phase 4: Polish & Cross-Cutting Concerns

**Purpose**: エッジケース対応と最終確認

- [x] T009 Add edge case handling instructions in .claude/skills/speckit.specify/SKILL.md — invalid URL format (e.g., PR URL), non-existent Issue, cross-repo URL warning
- [x] T010 Review full .claude/skills/speckit.specify/SKILL.md for consistency — ensure no references to old auto-numbering behavior remain, all error messages are actionable, and the flow is coherent end-to-end

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **US1 (Phase 2)**: Depends on Phase 1
- **US2 (Phase 3)**: Can run in parallel with Phase 2 (different sections of SKILL.md), but logically depends on US1 being designed first
- **Polish (Phase 4)**: Depends on Phase 2 and Phase 3

### Within User Story 1

- T002 (URL detection) → T003 (API validation) → T004 (numbering replacement) — sequential, each builds on the previous
- T005 (short-name) can run after T003 (needs Issue title/body retrieval to be in place)
- T006 (collision check) can run after T004 (needs number to be determined)

### Within User Story 2

- T007 (URL requirement) and T008 (remove --number instructions) are independent within the same file but different sections

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Read and understand current SKILL.md
2. Complete Phase 2: Implement Issue URL → number extraction flow
3. **STOP and VALIDATE**: Test with a real Issue URL
4. If working, proceed to Phase 3

### Incremental Delivery

1. T001 → Understand current state
2. T002–T006 → Issue URL flow working (MVP)
3. T007–T008 → Issue URL required, old flow blocked
4. T009–T010 → Edge cases and cleanup

---

## Notes

- All changes are in a single file: `.claude/skills/speckit.specify/SKILL.md`
- No Go code, no shell script changes, no test code
- `create-new-feature.sh` の `--number` オプションはスクリプト内部インターフェースとして維持（Research R3）
- 手動テストで検証（Constitution Test-First は Go コード対象、本機能は対象外）
