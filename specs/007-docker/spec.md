# Feature Specification: psm Example Project (Starter Template)

**Feature Branch**: `007-docker`
**Created**: 2026-03-15
**Revised**: 2026-03-16
**Status**: Draft
**Input**: `example/` ディレクトリをコピーするだけで psm による Parameter Store 管理リポジトリとして使えるスターターテンプレートを提供する

## Background

前回の spec は「ツール同梱 Docker イメージをワンショットで使う」前提だったが、テスト設計を詰める過程で以下が判明:

- SSO 認証が必要 → 対話的にコンテナに入る必要がある → DevContainer が自然
- 暗号化/復号/sync まで一気通貫で試したい → Dockerfile 単体では不十分
- 「コピーして即使える」のがゴール → テンプレートリポジトリ相当の example が必要

## User Scenarios & Testing _(mandatory)_

### User Story 1 - example/ をコピーして psm 管理リポジトリを立ち上げる (Priority: P1)

ユーザーが `example/` ディレクトリをまるごとコピーし、新しいリポジトリのルートとして使う。DevContainer を開けば psm + SOPS + age + AWS CLI が揃った環境に入れる。age キーを作り、シークレットを暗号化し、AWS SSO でログインして `psm sync` で Parameter Store に反映する、という一連の流れが動く。

**Why this priority**: psm 単体のインストールはできても、SOPS の鍵設定や AWS 認証との組み合わせで躓くユーザーが多い。一気通貫の動作例があれば導入障壁が大幅に下がる。

**Independent Test**: example/ をコピー → DevContainer を起動 → age キー生成 → sops 暗号化 → AWS SSO ログイン → psm sync → SSM Parameter Store に値が入ることを確認。

**Acceptance Scenarios**:

1. **Given** `example/` をコピーした新ディレクトリ, **When** DevContainer を開く, **Then** psm, sops, age, aws コマンドがすべて使える
2. **Given** DevContainer 内, **When** age キーを生成し `.sops.yaml` に公開鍵を設定して `sops -e` する, **Then** secrets ファイルが暗号化される
3. **Given** 暗号化された secrets と AWS SSO 認証, **When** `sops -d secrets.enc.yaml | psm sync --store ssm /dev/stdin` を実行する, **Then** Parameter Store にパラメータが反映される

### Edge Cases

- age キーが未生成の状態で `sops -e` した場合、わかりやすいエラーになるか
- AWS SSO 未ログイン状態で `psm sync` した場合、認証エラーが明確か
- 既存リポジトリに example/ の中身をマージして使えるか（構造の独立性）

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: `example/` は独立したリポジトリとして機能するディレクトリ構成とする
- **FR-002**: DevContainer 定義を含み、開くだけで psm, SOPS, age, AWS CLI v2 が使える環境になること
- **FR-003**: サンプルの secrets ファイル（ダミー値）と `.sops.yaml` のテンプレートを含むこと
- **FR-004**: ツールバージョンは Dockerfile の ARG で変更可能とすること
- **FR-005**: README に鍵生成から psm sync までの一連の手順を記載すること
- **FR-006**: 各ツールのライセンス情報を記載すること
- **FR-007**: psm は GitHub Releases からバイナリをダウンロードしてインストールすること

### Non-Functional Requirements

- **NFR-001**: DevContainer の起動が軽量であること（Claude Code や Go ツールチェインは含めない）

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: `example/` を別ディレクトリにコピーし、DevContainer を開き、全ツールが動作する
- **SC-002**: README の手順に従って鍵生成 → 暗号化 → SSO ログイン → psm sync が完走する
- **SC-003**: example/ 内のファイルが psm 開発リポジトリの他の部分に依存していない

## Assumptions

- ベースイメージは debian:bookworm-slim（survey で検証済み）
- マルチプラットフォーム対応: amd64/arm64（survey で検証済み）
- ライセンス: SOPS (MPL 2.0), age (BSD 3-Clause), AWS CLI (Apache 2.0), psm (自作)
- ユーザーは AWS アカウントと SSO 設定を持っている前提
