# Survey: psm Example Project (Starter Template)

**Date**: 2026-03-16
**Spec**: specs/007-docker/spec.md

## Summary

spec の方向転換（ワンショット Docker → テンプレートリポジトリ）は正しい。しかし「example/ に DevContainer を入れる」設計には、既存の開発用 DevContainer との構造的な重複と、example/ の独立性に関する検討事項がある。最大のリスクは DevContainer の複雑さが「コピーして即使える」というゴールを損なうこと。軽量に保つ設計判断が鍵。

## S1: Problem Reframing — 本当に DevContainer が必要か？

**Category**: Problem Definition
**Finding**: spec は「DevContainer を開けば全ツールが揃う」を FR-002 で要求している。しかし本当のゴールは「psm を使った Parameter Store 管理の一気通貫のサンプル」であり、DevContainer はその手段の一つ。代替手段として:

1. **DevContainer（spec の提案）**: 開くだけで環境完成。IDE 統合あり。ただし Docker + VS Code/対応エディタが前提。
2. **Makefile + シェルスクリプト**: ホストに直接ツールをインストール。DevContainer 不要。ただしツールバージョン管理が煩雑。
3. **docker-compose でワンショット**: `docker compose run` で対話シェルに入る。devcontainer.json 不要だが VS Code 統合なし。

**Recommendation**: DevContainer が最適。psm の主要ユースケースは「チームで secrets を管理」であり、VS Code + DevContainer は最も一般的な開発環境。Makefile 方式はツールインストールのサポート負荷が高い。docker-compose 方式は DevContainer のサブセットでしかなく、中途半端。

## S2: 既存 DevContainer との関係整理

**Category**: Integration Impact
**Finding**: 既存の `.devcontainer/` は psm **開発者**向けで、以下を含む:
- Go 1.26.1 + lint ツール群
- Claude Code
- Squid プロキシ（セキュリティ）
- gh-token-sidecar（GitHub App 認証）
- iptables ファイアウォール

example/ の DevContainer は psm **ユーザー**向けで、必要なのは:
- psm + SOPS + age + AWS CLI のみ
- セキュリティサイドカー不要（ユーザーの自由なネットワーク環境）
- Go ツールチェイン不要

**重複するもの**: Dockerfile のベースイメージ（debian:bookworm-slim）、AWS CLI インストール手順、AWS SSO 設定スクリプトのパターン。ただし共通化すべきではない — 目的が異なるため、独立して進化すべき。

**Recommendation**: 構造を参考にしつつ完全に独立したファイルとして作る。共通化の誘惑に負けない（YAGNI）。

## S3: example/ の DevContainer は単一コンテナで十分か？

**Category**: Approach Alternatives
**Finding**: 開発用 DevContainer は 3 サービス構成（workspace + gh-token-sidecar + outbound-filter）だが、example/ にはどちらも不要:
- gh-token-sidecar: psm ユーザーは GitHub App 認証を使わない
- outbound-filter: セキュリティ制約はユーザー環境の責任

**Recommendation**: example/ は **単一 Dockerfile + devcontainer.json** のみ。docker-compose.yml は使わない。Simplicity First。

## S4: `.sops.yaml` テンプレートの設計

**Category**: Underspecification
**Finding**: FR-003 は「サンプルの secrets ファイルと `.sops.yaml` のテンプレート」を要求。しかし `.sops.yaml` にはユーザー固有の age 公開鍵が必要で、テンプレートのまま動くものではない。

選択肢:
1. プレースホルダー入りの `.sops.yaml` を置く → ユーザーが書き換え必須
2. `.sops.yaml.example` として置く → ユーザーがコピー＆編集
3. README の手順のみで `.sops.yaml` は含めない

**Recommendation**: `.sops.yaml.example` として配置し、README で「公開鍵を書き換えてリネーム」の手順を示す。Git 管理対象の `.sops.yaml` にダミー鍵が残るリスクを避けるため。

## S5: サンプル secrets ファイルの形式

**Category**: Feasibility Verification
**Finding**: psm が受け付ける secrets ファイルの形式を確認する必要がある。psm は YAML の flat key-value を期待（research.md より）。サンプルは以下のような形になる:

```yaml
/myapp/database/host: localhost
/myapp/database/port: "5432"
/myapp/api/key: dummy-api-key-replace-me
```

暗号化前の平文ファイルをリポジトリに含めるべきか？→ サンプル用ダミー値なので問題ない。暗号化済みファイルは age キーがないとユーザーの環境で復号できないため含めない。

**Recommendation**: 平文のサンプル `secrets.yaml`（ダミー値）のみ含める。暗号化済みファイルはユーザーが手順に従って生成する。

## S6: AWS SSO 設定の扱い

**Category**: External Dependencies
**Finding**: 開発用 DevContainer は `setup-aws-config.sh` で環境変数から `~/.aws/config` を生成している。example/ でも同じパターンが使えるが、ユーザーの AWS 環境は多様:
- SSO（推奨）
- IAM ユーザー（レガシー）
- IAM ロール（CI/CD）

**Recommendation**: `setup-aws-config.sh` は含めない。代わりに README で「`aws configure sso` を実行してプロファイルを設定」の手順を書く。ユーザーの既存 AWS 設定をそのまま使えるよう、ホストの `~/.aws` をマウントする設計にする。

## S7: psm バージョンの固定

**Category**: Risk & Failure Modes
**Finding**: 現在 psm は v0.0.1。example/ の Dockerfile で `PSM_VERSION=0.0.1` を固定すると、psm のリリースごとに example/ も更新が必要。しかしこれは「reference」であり、最新追従は YAGNI。

**Recommendation**: 初期値を現在の最新版で固定。README に「最新版は GitHub Releases を確認」と記載。Dependabot 等の自動更新は YAGNI。

## S8: ディレクトリ名 — `example/` vs `example-repo/` vs `starter/`

**Category**: Scope Boundaries
**Finding**: `example/` は一般的すぎて、テストコードの例やサンプルコードと混同される可能性がある。内容は「コピーして独立リポジトリとして使うテンプレート」なので、その意図が名前から伝わるべき。

候補:
- `example/` — シンプルだが汎用的すぎる
- `starter/` — テンプレート感がある
- `example-repo/` — 明示的

**Recommendation**: ユーザー判断に委ねる。機能に影響しない。

## Items Requiring PoC

- DevContainer 内で `aws sso login` が正常に動作するか（ブラウザ起動のリダイレクト）。開発用 DevContainer で実績があるため問題ないと思われるが、ポート転送設定等が必要かもしれない。

## Constitution Impact

None. この feature は Go コードの変更を含まず、テンプレートファイルの追加のみ。Test-First は N/A（Dockerfile / 設定ファイル）。Simplicity First と YAGNI は設計判断の指針として適用。

## Recommendation

Proceed to plan. 主な設計判断:
- 単一コンテナ DevContainer（docker-compose 不要）
- `.sops.yaml.example` パターン
- ホストの `~/.aws` マウント（SSO スクリプト不要）
- S8 のディレクトリ名はユーザーに確認
