# Research: git worktree セットアップ Skill

## R1: 相対パス変換の具体的手順

**Decision**: worktree 作成後に 2 つのファイルを書き換える
**Rationale**: survey の PoC で動作確認済み。Git 2.39 環境で唯一の workaround
**Alternatives considered**:
- Git 2.48+ にアップグレード → Debian bookworm リポジトリに存在せず、ソースビルドが必要で DevContainer の保守コスト増大
- Claude Code 組み込み EnterWorktree を使う → ディレクトリ固定 (`.claude/worktrees/`)、既存ブランチ非対応、相対パス非対応

**変換対象**:
1. `.worktrees/<dir>/.git` — `gitdir: /absolute/path/.git/worktrees/<name>` → `gitdir: ../../.git/worktrees/<name>`
2. `.git/worktrees/<name>/gitdir` — `/absolute/path/.worktrees/<dir>/.git` → `../../.worktrees/<dir>/.git`

## R2: Skill 名とトリガー条件

**Decision**: `worktree-setup` を Skill 名とする
**Rationale**: Claude Code 組み込みの `EnterWorktree` との名前衝突を避ける。`setup` が付くことで「セットアップ作業の自動化」という目的が明確
**Alternatives considered**:
- `worktree` → 組み込みツールと混同リスク
- `wt` → 略称すぎて discoverability が低い

## R3: ブランチ新規/既存の自動判定方法

**Decision**: ローカル・リモートのブランチ一覧を `git branch -a` で取得し、指定ブランチが存在するかで判定
**Rationale**: 追加依存なし、git 標準コマンドのみで判定可能
**Alternatives considered**:
- ユーザーに新規/既存を明示指定させる → 手順が増えて UX 低下

## R4: worktree 削除時の対処

**Decision**: Skill の案内に削除手順を含める（手動削除 + `git worktree prune`）
**Rationale**: 削除は作成に比べ頻度が低く、Skill に削除機能を組み込むとスコープが広がる。YAGNI 原則に従い案内のみ
**Alternatives considered**:
- Skill に削除サブコマンドを追加 → YAGNI 違反、スコープ拡大
- パスを絶対に戻してから `git worktree remove` → 自動化は可能だが削除頻度が低いため過剰
