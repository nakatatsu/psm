# Quickstart: 016-sync-issue-numbers

## What to change

1 ファイル: `.claude/skills/speckit.specify/SKILL.md`

## Key changes

1. **Issue URL 必須化**: ユーザー入力から GitHub Issue URL を検出。見つからなければエラー
2. **Issue 情報取得**: `gh issue view` で Issue の実在確認 + タイトル・説明文取得
3. **番号抽出**: URL から Issue 番号を抽出し `--number` に渡す（自動採番ロジック削除）
4. **short-name 生成**: Issue タイトル・説明文から生成（feature description からの生成を置換）
5. **衝突チェック**: 既存ブランチとの番号衝突時はエラー中断

## What NOT to change

- `create-new-feature.sh` — 変更不要。`--number` オプションは内部インターフェースとして維持
- 他のスキル（clarify, survey, plan, analyze, tasks）— 影響なし
- Constitution — 修正不要

## Testing

手動テスト 3 パターン:
1. 有効な Issue URL → ブランチ番号が Issue 番号と一致
2. Issue URL なし → エラーメッセージ表示
3. 不正な URL / 存在しない Issue → エラーメッセージ表示
