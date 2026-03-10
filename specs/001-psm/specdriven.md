# Spec-Driven Development 手順書: psm

このプロジェクトを `.claude/skills/speckit.*` スキル群を使って仕様駆動で実装するための手順。

## 前提

- 現在 `main` ブランチにいる
- `spec.md`（手書きの仕様）がプロジェクトルートにある
- `.specify/` のインフラ（テンプレート、スクリプト）は準備済み
- ネットワーク接続不可（TODO.md記載の通り）

## ワークフロー全体像

```
constitution (任意)
    |
 specify ── spec.md をテンプレート形式に変換、フィーチャーブランチ作成
    |
 clarify ── 仕様の曖昧箇所を質問→spec.mdに反映
    |
  plan ── 技術設計（research.md, data-model.md, contracts/）
    |
 checklist (任意) ── 仕様品質チェックリスト
    |
  tasks ── タスク一覧生成（tasks.md）
    |
 analyze (任意) ── spec/plan/tasks の整合性チェック
    |
implement ── タスクに沿って実装
```

各ステップは前のステップの成果物に依存する。スキップ可能なものは「任意」と表記。

---

## Step 0: Constitution（任意）

**目的**: プロジェクト全体の非交渉原則を定める。小規模CLIツールなら省略してよい。

**やるなら**:
```
/speckit.constitution Go言語のCLIツール。シンプルさ最優先、YAGNI原則。テストは標準のgo testのみ。
```

**何が起きるか**:
- `.specify/memory/constitution.md` のプレースホルダが具体的な原則に置き換わる
- 以降のスキルがこの原則を参照して整合性チェックを行う

**あなたがやること**: 原則の内容を確認し、不要なものがあれば指摘する。

**省略した場合**: 以降のスキルはconstitutionチェックをスキップする。問題なし。

---

## Step 1: Specify（必須）

**目的**: 手書きの `spec.md` をテンプレート形式の仕様書に変換し、フィーチャーブランチを作成する。

```
/speckit.specify SOPS復号済みYAMLファイルを読み、AWS SSM Parameter Store または Secrets Manager に同期するCLIツール psm を作る。Go言語。詳細はプロジェクトルートの spec.md を参照。
```

**何が起きるか**:
1. フィーチャーブランチが作成される（例: `001-psm`）
2. `specs/001-psm/` ディレクトリが作られる
3. テンプレートに沿った `specs/001-psm/spec.md` が生成される
   - ユーザーストーリー（P1, P2, P3...）に分解
   - 受入シナリオ（Given/When/Then）
   - 機能要件（FR-001, FR-002...）
   - 成功基準
4. `specs/001-psm/checklists/requirements.md` で品質チェック
5. 不明点があれば `[NEEDS CLARIFICATION]` マーカー付きで質問される（最大3問）

**あなたがやること**:
- 質問に回答する（選択肢から選ぶか、短い回答を入力）
- 生成された仕様書を確認する
- プロジェクトルートの `spec.md` に書いた意図と齟齬がないか確認する

**完了の目安**: ブランチ名と `specs/*/spec.md` のパスが報告される。

---

## Step 2: Clarify（推奨）

**目的**: 仕様の曖昧箇所を洗い出し、最大5問の質問で解消する。

```
/speckit.clarify
```

**何が起きるか**:
1. 仕様書を10カテゴリ（機能スコープ、データモデル、UX、非機能要件、外部依存、エッジケース等）で走査
2. 曖昧・未定義の箇所を特定
3. 1問ずつ質問が表示される（推奨回答付き）
4. 回答が `spec.md` の該当セクションに直接反映される

**あなたがやること**:
- 各質問に回答する。推奨で良ければ「yes」と答えるだけ
- 「done」と言えば途中で終了できる

**psmで想定される質問例**:
- SOPSメタデータキーの扱い方（無視? エラー?）
- `--prune` 時の確認プロンプトの有無
- 同時実行時のリトライ戦略
- ログ出力のフォーマット（構造化? プレーンテキスト?）

**省略した場合**: plan 以降で曖昧さが問題になりうる。手戻りリスクが上がる。

---

## Step 3: Plan（必須）

**目的**: 技術設計。使うライブラリ、データ構造、インターフェース契約を決める。

```
/speckit.plan Go 1.22, AWS SDK for Go v2, gopkg.in/yaml.v3, CLI引数はflagパッケージ
```

**何が起きるか**:
1. `specs/*/plan.md` に技術コンテキストが埋まる
2. 不明点があれば `research.md` で調査・解決される
3. `data-model.md` にエンティティ定義が生成される
4. `contracts/` にインターフェース契約が生成される（CLIツールならコマンドスキーマ）
5. `quickstart.md` に利用シナリオが書かれる
6. CLAUDE.md 等のエージェントコンテキストが更新される

**あなたがやること**:
- 技術的な選択肢について質問があれば回答する
- 生成されたplan.mdを確認する（特にProject Structure）
- data-model.md の エンティティが仕様と合っているか確認する

**完了の目安**: plan.md のパスと生成された成果物一覧が報告される。

**注意**: このコマンドは設計までで止まる。実装には進まない。

---

## Step 4: Checklist（任意）

**目的**: 仕様品質のチェックリストを生成する。仕様の見落としを発見するのに有用。

```
/speckit.checklist CLIツールのセキュリティとエッジケース
```

**何が起きるか**:
- `specs/*/checklists/security.md` などにチェックリストが生成される
- 各項目は「要件が書かれているか」を問う（実装の動作確認ではない）

**あなたがやること**: チェックリストを読み、未定義の要件に気づいたら `spec.md` を更新する。

**省略した場合**: Step 1 で基本的な `requirements.md` チェックリストは既に作られている。追加の観点が不要なら省略可。

---

## Step 5: Tasks（必須）

**目的**: 実装タスクを依存関係付きで生成する。

```
/speckit.tasks
```

**何が起きるか**:
1. plan.md と spec.md からタスク一覧が生成される
2. ユーザーストーリー単位でフェーズ分け:
   - Phase 1: Setup（プロジェクト初期化）
   - Phase 2: Foundational（共通基盤）
   - Phase 3+: 各ユーザーストーリー（P1, P2, P3...）
   - Final: Polish
3. 各タスクに ID、並列マーカー [P]、ストーリーラベル [US1] が付く
4. `specs/*/tasks.md` に出力される

**あなたがやること**:
- タスク一覧を確認する
- タスクの粒度や順序に問題があれば修正を依頼する
- MVPスコープ（通常 User Story 1 まで）が適切か確認する

**完了の目安**: タスク総数、ストーリーごとの数、並列実行可能箇所が報告される。

---

## Step 6: Analyze（推奨）

**目的**: spec.md、plan.md、tasks.md の三者間の整合性を検証する。

```
/speckit.analyze
```

**何が起きるか**:
1. 要件とタスクのカバレッジマッピング
2. 重複、曖昧さ、未カバー要件の検出
3. constitution違反チェック（設定済みの場合）
4. 深刻度付きの分析レポートが出力される（ファイル変更なし）

**あなたがやること**:
- CRITICAL/HIGH の問題があれば、指示に従って修正する
- 修正が必要な場合、該当するスキルを再実行する（例: `/speckit.specify` で仕様修正）
- 問題なければ実装に進む

**省略した場合**: 整合性の問題が実装中に発覚するリスクがある。小さいプロジェクトなら許容範囲。

---

## Step 7: Implement（必須）

**目的**: tasks.md に従ってプロジェクトを実装する。

```
/speckit.implement
```

**何が起きるか**:
1. チェックリストの完了状況を確認（未完了があれば確認される）
2. tasks.md, plan.md, data-model.md 等を読み込み
3. .gitignore 等の ignore ファイルを技術スタックに合わせて生成/検証
4. フェーズ順にタスクを実行:
   - Setup: go mod init, ディレクトリ構造作成
   - Foundational: AWS クライアント初期化、YAML パーサー等
   - User Stories: 各機能の実装
   - Polish: テスト、ドキュメント、リファクタリング
5. 完了したタスクは tasks.md 内で `[X]` にマークされる

**あなたがやること**:
- 基本的には自動で進む
- エラーが出た場合はデバッグを依頼するか、手動で修正する
- 各フェーズ完了時にチェックポイントの報告を確認する

---

## Step 8: Tasks to Issues（任意）

**目的**: tasks.md の各タスクを GitHub Issues に変換する。

```
/speckit.taskstoissues
```

**前提**: リモートが GitHub URL であること。チーム開発やタスク管理が必要な場合のみ。

---

## 実行順序の早見表

| Step | スキル | 必須? | 入力 | 出力 |
|------|--------|-------|------|------|
| 0 | `/speckit.constitution` | 任意 | 原則の方針 | `.specify/memory/constitution.md` |
| 1 | `/speckit.specify` | 必須 | 機能説明テキスト | ブランチ + `specs/*/spec.md` |
| 2 | `/speckit.clarify` | 推奨 | (自動) | `spec.md` 更新 |
| 3 | `/speckit.plan` | 必須 | 技術スタック | `plan.md`, `research.md`, `data-model.md`, `contracts/` |
| 4 | `/speckit.checklist` | 任意 | ドメイン指定 | `checklists/*.md` |
| 5 | `/speckit.tasks` | 必須 | (自動) | `tasks.md` |
| 6 | `/speckit.analyze` | 推奨 | (自動) | 分析レポート（画面出力のみ） |
| 7 | `/speckit.implement` | 必須 | (自動) | 実装コード |
| 8 | `/speckit.taskstoissues` | 任意 | (自動) | GitHub Issues |

## psm 固有のメモ

- **ネットワーク制約**: この環境から外部接続不可。AWS SDK のモック/スタブが必要になる可能性がある。plan フェーズで考慮すること。
- **既存 spec.md**: プロジェクトルートの `spec.md` は手書きの仕様。Step 1 で `specs/*/spec.md` にテンプレート形式で再生成される。元の `spec.md` はリファレンスとして残る。
- **Go 固有**: `go mod init`, `go test ./...` 等は implement フェーズで自動的にセットアップされる。
- **TODO.md**: 既存の TODO.md は Step 1 以降で実質的に置き換えられる。完了後に整理を検討。
