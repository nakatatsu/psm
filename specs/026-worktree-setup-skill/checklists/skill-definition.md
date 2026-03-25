# Skill Definition Checklist: git worktree セットアップ Skill

**Purpose**: SKILL.md の要件定義が完全・明確・一貫しているかを検証する
**Created**: 2026-03-23
**Feature**: [spec.md](../spec.md)

## Requirement Completeness

- [ ] CHK001 新規ブランチ作成時のベースブランチ (origin/main) が明記されているか？ [Completeness, Spec FR-001]
- [ ] CHK002 既存ブランチ判定の対象（ローカル・リモート両方）が spec で定義されているか？ [Completeness, Spec Assumptions]
- [ ] CHK003 相対パス変換の対象ファイル 2 つ（`.worktrees/<dir>/.git`, `.git/worktrees/<name>/gitdir`）が FR-008 で明記されているか？ [Completeness, Spec FR-008]
- [ ] CHK004 Skill の呼び出し方（コマンド名、引数形式）が spec で定義されているか？ [Completeness, Gap]
- [ ] CHK005 `.gitignore` への追加時のエントリ形式（`.worktrees/` vs `.worktrees` vs `/.worktrees/`）が明確か？ [Completeness, Spec FR-004]

## Requirement Clarity

- [ ] CHK006 「ブランチ名にスラッシュが含まれる場合はハイフンに置換」のルールが一意に解釈できるか？（例: `feature/auth/v2` → `feature-auth-v2`？） [Clarity, Spec Edge Cases]
- [ ] CHK007 `worktree.useRelativePaths` 設定のスコープ（`--local` vs `--global`）が明記されているか？ [Clarity, Gap]
- [ ] CHK008 FR-008 の相対パス変換で使用する具体的な相対パス形式（`../../`）が明確か？ [Clarity, Spec FR-008]

## Requirement Consistency

- [ ] CHK009 Assumptions の「Git 2.39+」と FR-003 の `worktree.useRelativePaths` 設定は整合しているか？（2.39 では無効だが将来のため設定する旨が明確か） [Consistency, Spec FR-003 / Assumptions]
- [ ] CHK010 spec の Edge Cases「削除時の案内」と FR-007 の「次のステップ案内」は一貫した案内範囲か？ [Consistency, Spec FR-007 / Edge Cases]

## Scenario Coverage

- [ ] CHK011 正常系: 新規ブランチでの worktree 作成シナリオが定義されているか？ [Coverage, Spec US-1]
- [ ] CHK012 正常系: 既存ブランチでの worktree 作成シナリオが定義されているか？ [Coverage, Spec US-2]
- [ ] CHK013 異常系: 指定ブランチで worktree が既に存在する場合の要件が定義されているか？ [Coverage, Spec Edge Cases]
- [ ] CHK014 異常系: `git worktree add` 自体が失敗した場合（dirty tree、ロック等）の要件が定義されているか？ [Coverage, Gap]
- [ ] CHK015 リモートブランチが存在するがローカルにない場合の fetch 要否が定義されているか？ [Coverage, Gap]

## Edge Case Coverage

- [ ] CHK016 `.gitignore` が存在しない場合の新規作成が Edge Cases で定義されているか？ [Edge Case, Spec Edge Cases]
- [ ] CHK017 相対パス変換後の `git worktree remove` 非互換が Edge Cases で定義されているか？ [Edge Case, Spec Edge Cases]
- [ ] CHK018 現在のカレントディレクトリがリポジトリルートでない場合の動作が定義されているか？ [Edge Case, Gap]

## Dependencies and Assumptions

- [ ] CHK019 Git 2.39+ の前提が Assumptions に明記されているか？ [Assumption, Spec Assumptions]
- [ ] CHK020 DevContainer / ホスト間の相対パス必要性の根拠が記載されているか？ [Assumption, Spec Assumptions]
- [ ] CHK021 Claude Code Skill フォーマット（frontmatter, description）の準拠要件が定義されているか？ [Dependency, Spec FR-006]

## Notes

- Check items off as completed: `[x]`
- CHK004, CHK007, CHK014, CHK015, CHK018 は Gap として特定 — spec 更新の検討対象
- Items are numbered sequentially for easy reference
