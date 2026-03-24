# Quickstart: git worktree セットアップ Skill

## 成果物

`.claude/skills/worktree-setup/SKILL.md` — Claude Code Skill 定義ファイル（1 ファイル）

## Skill が行うこと

1. 前提条件の確認・自動セットアップ
   - `git config worktree.useRelativePaths true` 設定
   - `.gitignore` に `.worktrees/` 追加
2. ブランチの新規/既存を自動判定
3. `git worktree add` で `.worktrees/<dir>` に worktree 作成
4. 絶対パスを相対パスに変換（Git 2.39 workaround）
5. 完了メッセージと次のステップを案内

## Skill の呼び出し方

```
/worktree-setup <branch-name>
```

## 実装の流れ

1. SKILL.md の frontmatter（name, description）を定義
2. 手順セクションを Markdown で記述（Claude Code が実行する bash コマンド群）
3. エラーケースと案内メッセージを記述
