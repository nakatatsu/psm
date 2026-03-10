# Research: psm

**Date**: 2026-03-08

## R1: Go Version

**Decision**: Go 1.26 (latest stable: 1.26.1, released 2026-03-05)
**Rationale**: Latest stable release. AWS SDK for Go v2 requires minimum Go 1.24. No breaking changes to `flag` package.
**Alternatives**: Go 1.25 — still supported but no reason to use older version for a new project.

## R2: AWS SDK for Go v2

**Decision**: AWS SDK for Go v2 (released 2026-03-03)
- `github.com/aws/aws-sdk-go-v2` v1.41.3
- `github.com/aws/aws-sdk-go-v2/config` v1.32.11
- `github.com/aws/aws-sdk-go-v2/service/ssm` v1.68.2
- `github.com/aws/aws-sdk-go-v2/service/secretsmanager` v1.41.3

**Rationale**: 公式 SDK。v1 は非推奨。
**Alternatives**: v1 (deprecated) — 不適合。

### SSM API Usage

| 操作 | API | Bulk 対応 | Notes |
|------|-----|-----------|-------|
| 全件取得 | `GetParametersByPath` | Yes (paginate) | Path=`/`, Recursive=true, WithDecryption=true |
| 書き込み | `PutParameter` | **No** | 1 件ずつ。Overwrite=true, Type=SecureString。Rate limit: 40 TPS |
| 削除 | `DeleteParameters` | Yes (max 10/req) | 複数キーを一括削除可能 |

### Secrets Manager API Usage

| 操作 | API | Bulk 対応 | Notes |
|------|-----|-----------|-------|
| 全件取得 | `ListSecrets` + `BatchGetSecretValue` | Yes (max 20/req) | ListSecrets で名前一覧 → BatchGetSecretValue で値取得 |
| 書き込み | `CreateSecret` / `PutSecretValue` | **No** | 1 件ずつ。CreateSecret → ResourceExistsException → PutSecretValue |
| 削除 | `DeleteSecret` | **No** | 1 件ずつ。ForceDeleteWithoutRecovery=true |

### 並行実行戦略

Bulk write API が存在しないため、goroutine による並行実行で対応する。

- **方式**: セマフォ（`chan struct{}` または `golang.org/x/sync/errgroup`）で並行数を制限
- **並行数**: 10（SSM の rate limit 40 TPS に対して余裕を持たせる）
- **依存**: `golang.org/x/sync` は不要。標準ライブラリの `sync.WaitGroup` + channel で十分（Constitution Principle I）
- **エラー処理**: 1 件失敗しても残りは続行（FR-010）。エラーは集約して最後にまとめる

## R3: YAML Library

**Decision**: `gopkg.in/yaml.v3` v3.0.1
**Rationale**: フラットな key-value の Unmarshal のみ使用するため、基本機能で十分。`/` を含むキー名も正常にパース可能（YAML 仕様上、キーは任意の文字列）。
**Alternatives**:
- `go.yaml.in/yaml/v3` — 互換後継フォーク。移行メリットは薄い。
- `github.com/goccy/go-yaml` — オーバースペック。Simplicity First 原則に反する。

### YAML バリデーション戦略

`map[string]interface{}` への直接 Unmarshal は重複キーを黙って上書きするため不可。
`yaml.Node` API を使うことで、追加ライブラリなしに全バリデーションを実現できる。

```go
var doc yaml.Node
yaml.Unmarshal(data, &doc)
// doc.Content[0] = MappingNode
// 子ノードは [key, value, key, value, ...] の配列
```

| チェック項目 | yaml.Node での実現方法 |
|---|---|
| 重複キー | MappingNode の子キーを走査、`map[string]bool` で既出チェック |
| ネスト（マップ/配列） | 値ノードの `Kind` が `MappingNode` or `SequenceNode` ならエラー |
| null | 値ノードの `Tag` が `!!null` ならエラー |
| 空キー | キーノードの `Value` が `""` ならエラー |
| 型変換 | **不要**。`ScalarNode.Value` は元の文字列表現をそのまま保持（`5432` → `"5432"`, `true` → `"true"`） |

最後の点が重要: `yaml.Node` の `ScalarNode` は `Value` フィールドに元のリテラル文字列を持つため、int/bool → string の変換ロジックは一切不要。

## R4: CLI Argument Parsing

**Decision**: 標準ライブラリ `flag` パッケージ
**Rationale**: Constitution Principle I（Simplicity First）。必要なフラグは `--prune` と `--dry-run` の 2 つ + 位置引数 1 つ。`flag` パッケージで十分。
**Alternatives**: `cobra`, `urfave/cli` — 不要な依存追加。サブコマンド不要のため。

## R5: テスト戦略

**Decision**: `go test` + Sandbox AWS アカウントでの統合テスト
**Rationale**: psm の本質は AWS API とのやり取りそのものであり、Store 実装がアプリケーションの主体である（リポジトリパターンのリポジトリに相当）。Store interface を mock しても SDK ラッパーの薄い層をスキップするだけで、テストとしての価値がない。顧客のシークレットを扱うツールであり、エミュレータの挙動差異に起因するバグは許容できない。

**テスト分類**:
- **YAML パース・バリデーション**: 純粋なロジック。`go test` で完結（AWS 不要）
- **CLI パース**: 純粋なロジック。`go test` で完結（AWS 不要）
- **plan 関数**: 入力は `[]Entry` と `map[string]string`。`go test` で完結（AWS 不要）
- **Store 実装 (SSMStore/SMStore)**: Sandbox AWS アカウントで統合テスト
- **execute 関数**: Store を経由するため Sandbox AWS で統合テスト
- **E2E (sync/export)**: Sandbox AWS で統合テスト

**Sandbox AWS アカウント**:
- テスト専用の AWS アカウントを使用
- テストデータは `/psm-test/` プレフィックス（SSM）、`psm-test/` プレフィックス（SM）で隔離
- 各テストケースは setup/teardown でテストデータをクリーンアップ
- CI では AWS クレデンシャルを環境変数またはシークレットで管理

**却下した代替案**:
- **interface-based mocking**: Store がアプリの主体であり、mock では何もテストしていないのと同じ
- **LocalStack**: 2026年3月以降 auth token 必須（商用は有料）。`ResourceExistsException` の quirk、`PutSecretValue` が存在しない secret を黙って作成するバグ等、挙動差異のリスクあり。顧客シークレットを扱うツールには不適
- **testify, mockgen**: Constitution で禁止
