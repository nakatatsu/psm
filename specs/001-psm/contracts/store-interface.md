# Store Interface Contract

psm は SSM と Secrets Manager を共通インターフェースで操作する。

## Interface: Store

```go
type Store interface {
    // GetAll returns all current key-value pairs in the account.
    GetAll(ctx context.Context) (map[string]string, error)

    // Put creates or updates a single parameter/secret.
    Put(ctx context.Context, key, value string) error

    // Delete removes parameters/secrets by key.
    // Multiple keys may be passed for batch deletion (implementation may batch internally).
    Delete(ctx context.Context, keys []string) error
}
```

## SSM Implementation

| Method  | AWS API Call         | Bulk | Notes                                              |
|---------|----------------------|------|----------------------------------------------------|
| GetAll  | GetParametersByPath  | Yes  | Path=`/`, Recursive=true, WithDecryption=true, paginate |
| Put     | PutParameter         | No   | 1 件ずつ。goroutine で並行実行（並行数 10）        |
| Delete  | DeleteParameters     | Yes  | max 10 keys/req で一括削除                         |

## Secrets Manager Implementation

| Method  | AWS API Call                        | Bulk | Notes                                    |
|---------|--------------------------------------|------|------------------------------------------|
| GetAll  | ListSecrets + BatchGetSecretValue    | Yes  | ListSecrets → BatchGetSecretValue (max 20/req) |
| Put     | CreateSecret or PutSecretValue       | No   | 1 件ずつ。goroutine で並行実行（並行数 10） |
| Delete  | DeleteSecret                         | No   | 1 件ずつ。goroutine で並行実行。ForceDeleteWithoutRecovery=true |

## Concurrency

Bulk write API が存在しないため、Put/Delete（SM のみ）は goroutine で並行実行する。

- 並行数上限: 10（SSM rate limit 40 TPS に対して余裕を持たせる）
- 制御方式: buffered channel によるセマフォ
- エラー処理: 1 件失敗しても残りは続行。エラーは集約して呼び出し元に返す
