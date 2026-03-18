# Feature Specification: Codex を DevContainer で有効化する

**Feature Branch**: `015-codex-devcontainer`
**Created**: 2026-03-18
**Status**: Draft
**Input**: GitHub Issue #15 — Codex が DevContainer 環境で動作するか検証し、必要な設定を追加する

## User Scenarios & Testing *(mandatory)*

### User Story 1 - DevContainer 内で Codex CLI を実行する (Priority: P1)

開発者として、DevContainer 環境を起動した後、ターミナルから Codex CLI を実行して AI コーディング支援を受けたい。現在の環境は Claude Code 向けに構成されているが、Codex も並行して利用できるようにしたい。

**Why this priority**: Codex を DevContainer で使えること自体がこの Issue の核心であり、これが動作しなければ他の要件も意味をなさない。

**Independent Test**: DevContainer を起動し、ターミナルで `codex` コマンドを実行して正常にプロンプトが表示されることを確認する。

**Acceptance Scenarios**:

1. **Given** DevContainer が起動済みの状態で、**When** 開発者がターミナルで Codex CLI を起動する、**Then** Codex が正常に起動しプロンプトが表示される
2. **Given** DevContainer が起動済みの状態で、**When** 開発者が Codex にコード生成を依頼する、**Then** Codex が応答を返し、コード生成が実行される
3. **Given** DevContainer が起動済みの状態で、**When** 開発者が Claude Code と Codex を同一セッション内で切り替えて使用する、**Then** 両ツールが干渉なく正常に動作する

---

### User Story 2 - /codex スキルを使って Codex を呼び出す (Priority: P2)

開発者として、Claude Code のスキル機能から `/codex` コマンドで Codex を呼び出し、特定のタスクを Codex に委任したい。

**Why this priority**: CLI 直接実行（P1）が前提だが、スキル統合によりワークフローがスムーズになる。

**Independent Test**: Claude Code 内で `/codex` スキルを呼び出し、Codex が期待通りに起動・実行されることを確認する。

**Acceptance Scenarios**:

1. **Given** Claude Code が起動している状態で、**When** 開発者が `/codex` コマンドを入力する、**Then** Codex が起動し指定されたタスクを処理する
2. **Given** `/codex` スキルを実行中に、**When** Codex がエラーを返す、**Then** エラー内容が開発者に明確に伝わる

---

### User Story 3 - DevContainer 再構築後も Codex が利用可能 (Priority: P3)

開発者として、DevContainer を再構築（rebuild）した後も追加の手動設定なしで Codex が利用可能であってほしい。

**Why this priority**: 環境の再現性は DevContainer の根幹であり、手動設定が必要だと価値が半減する。

**Independent Test**: `devcontainer rebuild` 後に Codex CLI を実行し、初回起動時と同様に動作することを確認する。

**Acceptance Scenarios**:

1. **Given** DevContainer を rebuild した直後の状態で、**When** 開発者がターミナルで Codex CLI を起動する、**Then** 追加設定なしで Codex が正常に動作する

---

### Edge Cases

- Codex CLI のバージョンアップ時に DevContainer の設定が追従できるか？（バージョン固定 vs latest の方針）
- プロキシ環境（outbound-filter）経由で Codex の API 通信が正常に行えるか？
- Codex と Claude Code が同時に実行された場合にリソース競合が発生しないか？
- Codex に必要な API キーが未設定の場合、わかりやすいエラーメッセージが表示されるか？

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: DevContainer 起動後、追加の手動インストールなしで Codex CLI が利用可能であること
- **FR-002**: Codex CLI が DevContainer のプロキシ設定（outbound-filter）経由で外部 API と通信できること
- **FR-003**: Codex の導入が既存の Claude Code 環境に影響を与えないこと（共存できること）
- **FR-004**: `/codex` スキルが Claude Code のスキル一覧に登録され、呼び出し可能であること
- **FR-005**: DevContainer の rebuild 後も Codex が自動的に利用可能な状態になること
- **FR-006**: Codex が必要とする外部ドメインが outbound-filter の許可リストに追加されていること

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: DevContainer 起動後、1 分以内に Codex CLI でコード生成タスクを開始できる
- **SC-002**: Claude Code と Codex の両方が同一 DevContainer セッション内で問題なく動作する
- **SC-003**: DevContainer rebuild 後、手動設定ゼロで Codex が利用可能になる
- **SC-004**: `/codex` スキル経由で Codex を呼び出し、タスクが正常に完了する

## Assumptions

- Codex CLI は npm パッケージとしてグローバルインストール可能と仮定（OpenAI の公式 CLI に準拠）
- Codex の API 通信先ドメインは OpenAI の標準エンドポイント（api.openai.com 等）と仮定し、outbound-filter に追加する
- API キーの管理方法は既存の環境変数パターン（docker-compose.yml の environment セクション）に従う
- Codex のバージョンは `latest` を使用し、Claude Code と同様のバージョン管理方針とする
