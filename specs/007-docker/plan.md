# Implementation Plan: psm Example Project (Starter Template)

**Branch**: `007-docker` | **Date**: 2026-03-16 | **Spec**: specs/007-docker/spec.md
**Input**: Feature specification from `/specs/007-docker/spec.md`

## Summary

`example/` ディレクトリとして、psm による Parameter Store 管理のスターターテンプレートを提供する。DevContainer 定義を含み、コピーするだけで psm + SOPS + age + AWS CLI が揃った環境が手に入る。鍵生成から暗号化、SSO ログイン、psm sync までの一連の手順を README で案内する。

## Technical Context

**Language/Version**: N/A（Go コードなし。Dockerfile + 設定ファイルのみ）
**Primary Dependencies**: psm (GitHub Releases), SOPS v3.12.1, age v1.2.1, AWS CLI v2.34.9
**Storage**: N/A
**Testing**: 手動検証（DevContainer 起動 → ツール動作確認 → E2E フロー）
**Target Platform**: linux/amd64, linux/arm64（DevContainer 内）
**Project Type**: Template / Reference files
**Performance Goals**: N/A
**Constraints**: DevContainer 起動が軽量であること（Claude Code / Go ツールチェイン含めない）
**Scale/Scope**: ファイル 5-6 個の小規模成果物

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| # | Principle | Gate Question | Pass? |
|---|-----------|---------------|-------|
| I | Simplicity First | Is this the simplest viable design? No unnecessary abstractions or dependencies? | Yes — 単一コンテナ DevContainer、docker-compose 不要 |
| II | YAGNI | Does every element serve a present, concrete need? No speculative features? | Yes — 最小限のツールセットのみ。セキュリティサイドカー等は含めない |
| III | Test-First (NON-NEGOTIABLE) | N/A — Go コード変更なし。Dockerfile と設定ファイルのみ | N/A |

## Project Structure

### Documentation (this feature)

```text
specs/007-docker/
├── spec.md
├── survey.md
├── plan.md              # This file
└── tasks.md             # /speckit.tasks で生成
```

### Source Code (repository root)

```text
example/
├── .devcontainer/
│   ├── Dockerfile       # psm + SOPS + age + AWS CLI v2
│   └── devcontainer.json
├── .sops.yaml.example   # SOPS 設定テンプレート（公開鍵のプレースホルダー）
├── secrets.yaml         # サンプルシークレット（ダミー値）
├── test.sh              # 手動検証用コマンド集
└── README.md            # 鍵生成 → 暗号化 → SSO → psm sync の手順
```

**Structure Decision**: 単一の `example/` ディレクトリに全ファイルを配置。DevContainer は単一コンテナ構成（docker-compose.yml 不要）。`example/` をコピーすればそのまま独立リポジトリとして機能する。

### Design Decisions

| Decision | Choice | Rationale (Survey ref) |
|----------|--------|----------------------|
| コンテナ構成 | 単一 Dockerfile | docker-compose 不要。Simplicity First（S3） |
| SOPS 設定 | `.sops.yaml.example` | ダミー鍵が .sops.yaml に残るリスク回避（S4） |
| AWS 認証 | コンテナ内で独立して SSO 設定・ログイン | ホスト権限に依存しない完全独立環境 |
| ベースイメージ | debian:bookworm-slim | Alpine は AWS CLI 非互換。実績あり（Survey S2, 前回） |
| secrets サンプル | 平文 YAML のみ | 暗号化済みは age キーなしで使えない（S5） |

### Dockerfile 設計

既存の `docker/Dockerfile` をベースに、DevContainer 向けに調整:

- **含めるもの**: psm, SOPS, age, age-keygen, AWS CLI v2
- **含めないもの**: Claude Code, Go, lint ツール, iptables/ファイアウォール, gh-token-sidecar
- **ユーザー**: デフォルトの非 root ユーザー（devcontainer.json の remoteUser で指定）
- **バージョン**: すべて ARG で変更可能

### devcontainer.json 設計

```jsonc
{
  "name": "psm Environment",
  "build": { "dockerfile": "Dockerfile" },
  "remoteUser": "psm"
  // ホストの ~/.aws はマウントしない。コンテナ内で aws configure sso → aws sso login する。
}
```

### test.sh 設計

手動実行するコマンドの羅列。スクリプトではなくコピペ用:

1. DevContainer 内ツール確認（psm, sops, age, aws）
2. age キー生成
3. `.sops.yaml` 設定
4. sops 暗号化
5. sops 復号確認
6. AWS SSO ログイン
7. psm sync 実行
8. SSM Parameter Store 確認
9. 後片付け

## Complexity Tracking

なし。Constitution 違反なし。

## Phase 0/1 Artifacts

このフィーチャーは Go コードを含まないため、以下は生成しない:
- research.md: survey.md で十分カバー済み
- data-model.md: データモデルなし
- contracts/: 外部インターフェースなし
- quickstart.md: test.sh がこの役割を担う
