# Quickstart: ブランチ戦略の決定と実装

## 前提条件

- Terraform >= 1.0 がインストール済み
- GitHub Provider の認証が設定済み
- `.tmp/mynote/infrastructure/github/` にアクセス可能

## 実装手順

### 1. Terraform: Ruleset 追加

```bash
cd .tmp/mynote/infrastructure/github/

# 既存リポジトリを import（初回のみ）
terraform import github_repository.psm psm

# 変更を確認
terraform plan

# 適用
terraform apply
```

### 2. CodeQL ワークフロー追加

`.github/workflows/codeql.yml` をリポジトリに追加し、push する。

### 3. commit-msg hook 設定

```bash
# hooks パスを設定
git config core.hooksPath .githooks

# 動作確認（Issue 番号なしでコミットするとエラー）
echo "test" > /tmp/test.txt
git add /tmp/test.txt
git commit -m "test without issue number"  # → rejected
git commit -m "test with issue number #45" # → accepted
```

### 4. 検証

```bash
# develop への直接プッシュが拒否されることを確認
git checkout develop
echo "test" >> README.md
git add README.md
git commit -m "direct push test #45"
git push origin develop  # → rejected

# PR 経由でマージ可能であることを確認
# GitHub UI で PR を作成し、CI 通過後にマージ
```
