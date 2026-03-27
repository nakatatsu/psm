# test.sh 改善タスク

レビュー指摘を受けて、テストの検証品質を引き上げる。
方針: 「CLI 出力の文字列確認」から「SSM 実状態の遷移確認」へ軸足を移す。

---

## Phase 1: 状態検証の基盤 (最優先)

### T1: assert_state helper の追加

現状の `get_param` / `put_param` / `param_exists` に加え、
**期待する全キー→値を一括比較する helper** を追加する。

```
assert_state "テスト名" \
  "/myapp/database/host=localhost" \
  "/myapp/database/port=5432" \
  "/myapp/database/password=do-not-look-at-me" \
  "/myapp/api/key=come-on-do-not-look-at-me"
```

- 期待キーの値が一致するか
- 期待にないキー（`/myapp/` 配下）が存在しないか

これを作ることで Phase 2 の各シナリオ強化が簡潔に書ける。

### T2: cleanup の強化

現在は固定5キーを `delete-parameters` しているが、
`/myapp/` 配下を `get-parameters-by-path` で全取得→全削除する方式に変更。
テストデータの漏れ残りを防ぐ。

---

## Phase 2: 既存シナリオの検証強化 (最優先)

### T3: Scenario 1 (Dry-run) — 全キー不在を確認

現状: `/myapp/database/host` の存在だけ確認
改善: cleanup 後に dry-run → `assert_state` で `/myapp/` 配下が空であることを確認

### T4: Scenario 2 (Sync) — 全パラメータの実値を検証

現状: `/myapp/database/host` = `localhost` だけ確認
改善: sync 後に `assert_state` で全4キーの値を確認

### T5: Scenario 3 (Delete) — 削除対象だけ消え、他は残ることを確認

現状: `/myapp/legacy/old-key` が消えたかだけ確認
改善:
- 削除対象が消えていること
- sync YAML のキーが残っていること（値も一致）
- delete pattern に一致しない既存キーが残っていること

### T6: Scenario 4 (Conflict) — 部分適用されていないことを確認

現状: exit code 1 だけ確認
改善: conflict 前後で SSM 状態が不変であることを `assert_state` で確認

### T7: Scenario 7 (No changes) — 実状態が不変であることを確認

現状: summary 文字列の grep だけ
改善: 実行前後で `assert_state` が同じ結果を返すことを確認

---

## Phase 3: 失敗系シナリオの追加 (次点)

### T8: 不正入力系テスト

- 対象ファイル不存在 → exit code != 0、SSM 変更なし
- 不正な YAML → exit code != 0、SSM 変更なし
- 空入力 → exit code or 挙動を定義して検証

### T9: debug logging の強化

現状: `level=DEBUG` が含まれるかだけ
改善:
- `--debug` なしで `level=DEBUG` が出ないことを確認
- `--debug` ありで出ることを確認
（デバッグログポリシー策定後に、出力内容の検証を追加）

---

## Phase 4: 対話フロー検証 (次点)

### T10: approve フローの yes/no テスト

- TTY 上で `y` → 変更が適用される
- TTY 上で `n` → 変更が適用されない
- TTY 上で空 Enter → 変更が適用されない

※ `expect` や `script` コマンドで擬似 TTY を使う必要あり。
  実装コスト高めなので Phase 4 に分類。

---

## 対象外 (今回はやらない)

- AWS API 権限不足テスト — sandbox の IAM ポリシー変更が必要で、テストスクリプト単体では制御困難
- put 成功 / delete 失敗の部分障害 — AWS 側のエラー注入手段がない
- これらは将来的にモック Store を使った Go 単体テストで対応する方が適切
