# Quickstart: psm

## Build

```bash
git clone <repository-url> && cd psm
go build -o psm ./
```

## Prerequisites

- AWS 認証情報が設定済み（認証ルールは [spec.md FR-002](spec.md) 参照）
- `AWS_REGION` または `AWS_DEFAULT_REGION` が設定済み

## 初回セットアップ（export）

### 1. 既存パラメータを YAML に書き出し

```bash
psm export --store ssm --profile prod secrets.yaml
```

生成される YAML:
```yaml
/myapp/prod/DATABASE_URL: "postgres://db.example.com:5432/mydb"
/myapp/prod/API_KEY: "sk-xxxxxxxxxxxx"
/myapp/prod/REDIS_HOST: "redis.example.com"
```

### 2. SOPS で暗号化して保存

```bash
sops -e secrets.yaml > secrets.enc.yaml
rm secrets.yaml
```

## 同期（sync）

### 1. SOPS で復号した YAML ファイルを用意

```bash
sops -d secrets.enc.yaml > secrets.yaml
```

### 2. SSM Parameter Store に同期

```bash
psm sync --store ssm --profile prod secrets.yaml
```

出力:
```
create: /myapp/prod/DATABASE_URL
create: /myapp/prod/API_KEY
create: /myapp/prod/REDIS_HOST
3 created, 0 updated, 0 deleted, 0 unchanged, 0 failed
```

### 3. 値を変更して再同期

```yaml
# secrets.yaml (API_KEY を更新)
/myapp/prod/DATABASE_URL: "postgres://db.example.com:5432/mydb"
/myapp/prod/API_KEY: "sk-yyyyyyyyyyyy"
/myapp/prod/REDIS_HOST: "redis.example.com"
```

```bash
psm sync --store ssm --profile prod secrets.yaml
```

出力（変更のあるキーのみ表示）:
```
update: /myapp/prod/API_KEY
0 created, 1 updated, 0 deleted, 2 unchanged, 0 failed
```

### 4. 事前に差分を確認（dry-run）

```bash
psm sync --store ssm --profile prod --dry-run secrets.yaml
```

### 5. 不要なキーを削除（prune）

YAML から `/myapp/prod/REDIS_HOST` を削除した後:

```bash
psm sync --store ssm --profile prod --prune secrets.yaml
```

出力:
```
delete: /myapp/prod/REDIS_HOST
0 created, 0 updated, 1 deleted, 2 unchanged, 0 failed
```

### 6. Secrets Manager に同期

```yaml
# secrets-sm.yaml
myapp/prod/DATABASE_URL: "postgres://db.example.com:5432/mydb"
myapp/prod/API_KEY: "sk-xxxxxxxxxxxx"
```

```bash
psm sync --store sm secrets-sm.yaml
```

出力:
```
create: myapp/prod/DATABASE_URL
create: myapp/prod/API_KEY
2 created, 0 updated, 0 deleted, 0 unchanged, 0 failed
```

## エラー時の出力例

```
# stdout
create: /myapp/prod/DATABASE_URL
create: /myapp/prod/REDIS_HOST
2 created, 0 updated, 0 deleted, 0 unchanged, 1 failed

# stderr
error: /myapp/prod/API_KEY: AccessDeniedException: ...
```
