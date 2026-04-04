# Feature Specification: main へのマージ元ブランチ制限

**Feature Branch**: `048-branch-source-guard`
**Created**: 2026-04-03
**Status**: Draft
**Input**: GitHub Issue #48

## 意図

GitFlow において `main` へマージできるのは `release-*` と `hotfix-*` のみであるべきだが、GitHub Ruleset にはマージ元ブランチを制限するネイティブ機能がない。

GitHub Actions の `check-source-branch` ジョブで PR のソースブランチ名を検証し、Required Status Check として強制することで代替する。

## 変更内容

- `.github/workflows/ci.yml` に `check-source-branch` ジョブを追加
- `main` 向け PR のみ発火し、`release-*` または `hotfix-*` 以外からの PR を失敗させる
- CI トリガーを全 PR に広げる（`pull_request:` にブランチフィルタなし）
- main の Ruleset で `check-source-branch` を Required Status Check に追加
