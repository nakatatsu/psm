---
name: gh-token
description: >
  GitHub トークンを gh-token-sidecar コンテナ経由で取得・更新するスキル。
  gh コマンドが認証エラー（401, 403, "auth required", "token expired", "bad credentials" 等）で
  失敗した場合は、必ずこのスキルを使ってトークンを再取得すること。
  gh auth, git push/pull/fetch, GitHub API 呼び出しなど GitHub 認証が絡む操作全般で
  権限切れが疑われる場合にも自動的に利用せよ。
---

# gh-token トークン取得スキル

このワークスペースでは GitHub App のインストールトークンを発行する
サイドカーコンテナ (`gh-token-sidecar`) が稼働している。
PAT は使わず、必ずこのサイドカー経由でトークンを取得する。

## トークン取得手順

```bash
export GH_TOKEN=$(curl -sf http://gh-token-sidecar/token | jq -r '.token')
```

取得後、同じシェルセッション内で `gh` コマンドや GitHub API を利用できる。

## サイドカーが応答しない場合

ヘルスチェックで生死を確認する:

```bash
curl -sf http://gh-token-sidecar/health
```

応答がなければサイドカーコンテナが停止している可能性がある。
ユーザーに「gh-token-sidecar コンテナが応答しません。`docker compose up gh-token-sidecar` で起動してください」と伝える。

## 重要な注意点

- トークンはリクエストごとに新規発行される（キャッシュなし）ので、期限切れを心配する必要はない。呼べば常に新鮮なトークンが返る。
- `GH_TOKEN` 環境変数にセットすれば `gh` CLI が自動的にそれを使う。`gh auth login` は不要。
- `.pem` 秘密鍵はサイドカーコンテナ内にのみ存在し、このコンテナからは一切アクセスできない。鍵を探したり読もうとしないこと。

## いつ使うか

1. `gh` コマンドが認証エラーで失敗したとき（最も典型的なケース）
2. 新しいシェルセッションを開始して GitHub 操作を行う前
3. トークンが失効した可能性があるとき（長時間経過後など）

認証エラーを検知したら、まずこのスキルでトークンを再取得し、それからコマンドをリトライする。
