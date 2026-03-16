# Feature Specification: Reference Dockerfile for psm

**Feature Branch**: `007-docker`
**Created**: 2026-03-15
**Status**: Draft
**Input**: SOPS + AWS CLI + psm を同梱した参考用 Dockerfile を提供する

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Dockerfile で psm 実行環境を構築 (Priority: P1)

ユーザーが Dockerfile をビルドするだけで SOPS + AWS CLI + psm が揃った環境を手に入れられる。CI/CD パイプラインや手元のマシンで即座に psm を使い始められる。

**Why this priority**: psm 単体では SOPS による復号と AWS 認証が別途必要。3 ツールをまとめた参考 Dockerfile があれば、導入の手間が大幅に減る。

**Independent Test**: Dockerfile をビルドし、コンテナ内で `psm sync --help`、`sops --version`、`aws --version` が全て動作することを確認する。

**Acceptance Scenarios**:

1. **Given** リポジトリの Dockerfile, **When** `docker build` する, **Then** イメージが正常にビルドされる
2. **Given** ビルドされたイメージ, **When** コンテナ内で `psm sync --help` を実行する, **Then** ヘルプが表示される
3. **Given** ビルドされたイメージ, **When** コンテナ内で `sops --version` と `aws --version` を実行する, **Then** それぞれバージョンが表示される

### Edge Cases

- psm の GitHub Releases からのダウンロードが失敗した場合、ビルドエラーが明確に出るか
- ARM64 環境でもビルドできるか

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: Dockerfile は psm、SOPS、AWS CLI v2 を含むイメージをビルドできなければならない
- **FR-002**: Dockerfile は参考用としてリポジトリに配置し、ユーザーが自分でビルドする形式とする（プリビルドイメージの公開はしない）
- **FR-003**: Dockerfile 内のツールバージョンは ARG で指定でき、ユーザーが変更可能でなければならない
- **FR-004**: psm は GitHub Releases からバイナリをダウンロードしてインストールする
- **FR-005**: イメージは軽量なベースイメージを使用すること
- **FR-006**: 各ツールのライセンス情報をコメントまたはドキュメントに記載すること

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: `docker build` が成功し、psm / sops / aws の全コマンドがコンテナ内で動作する
- **SC-002**: README に Dockerfile の使い方と前提条件が記載されている

## Assumptions

- ベースイメージ、SOPS のバージョン等は Plan で決定する
- マルチプラットフォーム対応（amd64/arm64）の要否は Survey で検討する
- ライセンス上の制約は確認済み（SOPS: MPL 2.0、AWS CLI: Apache 2.0、psm: 自作）
