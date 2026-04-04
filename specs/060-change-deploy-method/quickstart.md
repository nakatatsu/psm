# Quickstart: デプロイ方式の変更

## What Changes

`.github/workflows/release.yml` を以下のように変更する:

- **トリガー**: `on: push: tags: ['v*']` → `on: pull_request: types: [closed]` (main ブランチ対象)
- **ブランチ判定**: マージ元が `release-*` または `hotfix-*` かチェック
- **バージョン抽出**: ブランチ名からプレフィックスを除去し semver 検証
- **タグ作成**: `git tag v{version}` + `git push --tags`
- **リリース**: 既存の GoReleaser ステップをそのまま実行

## Workflow Structure

```yaml
on:
  pull_request:
    types: [closed]
    branches: [main]

jobs:
  release:
    if: >-
      github.event.pull_request.merged == true &&
      (startsWith(github.event.pull_request.head.ref, 'release-') ||
       startsWith(github.event.pull_request.head.ref, 'hotfix-'))
    steps:
      # 1. checkout
      # 2. extract version from branch name
      # 3. validate semver format
      # 4. check tag doesn't already exist
      # 5. create and push tag
      # 6. setup-go
      # 7. goreleaser release
```

## Validation Steps

1. `release-X.Y.Z` ブランチから main への PR をマージ → タグ作成 + リリース実行を確認
2. feature ブランチから main への PR をマージ → ワークフローがスキップされることを確認
3. 既存タグと重複するバージョンのリリースブランチ → エラーで停止を確認

## Files Changed

| File | Action |
|------|--------|
| `.github/workflows/release.yml` | Modify (トリガーとステップの変更) |
