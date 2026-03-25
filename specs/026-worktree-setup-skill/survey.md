# Survey: git worktree セットアップ Skill

**Date**: 2026-03-23
**Spec**: [spec.md](./spec.md)

## Summary

最も重大な発見は **Git バージョンの不一致**。DevContainer の Git は 2.39.5 であり、spec が前提とする 2.48+ の `worktree.useRelativePaths` ネイティブサポートがない。ただし、worktree 作成後にパスファイルを手動で相対パスに書き換える workaround が機能することを PoC で確認済み。この workaround には副作用（`git worktree remove` が失敗する）があり、削除時の対処が必要。

Claude Code 組み込みの `EnterWorktree`/`ExitWorktree` ツールの存在も確認したが、用途が異なるため代替にはならない。

## S1: Git バージョンと worktree.useRelativePaths の非互換

**Category**: Risk & Feasibility (C-10)
**Finding**: DevContainer は Debian bookworm の Git 2.39.5 を使用。`worktree.useRelativePaths` は Git 2.48 で導入された機能で、2.39 では config を設定してもサイレントに無視される（エラーにならない）。apt リポジトリにも 2.48+ は存在しない。
**Recommendation**: `git worktree add` 実行後に `.worktrees/<dir>/.git` と `.git/worktrees/<name>/gitdir` の絶対パスを相対パスに書き換える workaround を採用する。`git config worktree.useRelativePaths true` は将来 Git がアップグレードされたときのために設定しておくが、現時点では workaround に依存する。
**Evidence**:
- `git --version` → `git version 2.39.5`
- `apt-cache policy git` → 候補は `1:2.39.5-0+deb12u3` のみ
- PoC: 手動で相対パスに書き換え後、`git status` と `git log` が正常動作することを確認

## S2: 相対パス化による git worktree remove の非互換

**Category**: Risk & Feasibility (C-8)
**Finding**: Git 2.39 の `git worktree remove` は gitdir ファイルに絶対パスを要求する。相対パスに変換済みの worktree に対して実行すると `fatal: validation failed, cannot remove working tree` エラーになる。
**Recommendation**: Skill で worktree 削除機能もサポートするか、削除手順を案内に含める。削除時は (a) パスを一旦絶対に戻して `git worktree remove` するか、(b) 手動で `rm -rf .worktrees/<dir>` + `rm -rf .git/worktrees/<name>` + `git worktree prune` する。
**Evidence**: PoC で相対パス状態の worktree に対して `git worktree remove` を実行し、エラーを確認。絶対パスに戻してからの `git worktree remove` は成功。

## S3: Claude Code 組み込み EnterWorktree との関係

**Category**: Necessity (B-4)
**Finding**: Claude Code には `EnterWorktree`/`ExitWorktree` ツールが組み込まれている。ただし、これは `.claude/worktrees/` に worktree を作成しセッション内で使うもので、(1) ディレクトリが `.claude/worktrees/` 固定、(2) HEAD ベースの新規ブランチのみ、(3) 既存ブランチ指定不可、(4) セッションスコープ。ユーザーの要件（`.worktrees/` 配置、origin/main ベース、既存ブランチ対応、ホスト/DevContainer 間の相対パス）とは合致しない。
**Recommendation**: 組み込みツールとは別物として Skill を作成する。Skill 名で混乱しないよう、`worktree` ではなく `worktree-setup` のような名前にすることを推奨。
**Evidence**: `EnterWorktree` ツール定義を確認。`.claude/worktrees/` 固定、`name` パラメータのみ、ブランチ指定パラメータなし。

## S4: 問題の再定義 — 本当に解決すべきこと

**Category**: Problem Reframing (A-1)
**Finding**: spec は「worktree セットアップの自動化」と定義しているが、本質的な問題は「ホストと DevContainer で同じ worktree を使うための相対パス設定が手動では漏れやすい」こと。worktree 作成自体は `git worktree add` 一行で済むが、相対パス変換と前提条件チェックが手作業でのミスの元。Skill の価値はこの「漏れやすい前提条件の自動化」にある。
**Recommendation**: spec の方向性は正しい。ただし Git 2.48 に依存せず、workaround ベースの相対パス変換を中核機能として明確に位置づける。

## Items Requiring PoC

なし（S1, S2 の PoC は survey 中に完了済み）。

## Constitution Impact

- **Test-First (III)**: この Skill は SKILL.md（Markdown ファイル）の作成であり、Go コードではない。Constitution の test-first 原則は Go コードに対するもの。SKILL.md に対して `go test` は適用不可。手動での動作確認が検証手段となる。
- 憲法の修正は不要。

## Recommendation

以下を反映した上で `/speckit.plan` に進むことを推奨:

1. **spec 修正**: Git 2.48+ 前提を撤回し、Git 2.39+ での workaround（手動相対パス変換）を前提とする
2. **spec 修正**: worktree 削除時の注意事項（相対パスと `git worktree remove` の非互換）を Edge Cases に追加
3. **FR 追加**: worktree 作成後の相対パス変換を明示的な要件として追加
