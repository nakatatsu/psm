# AWS SSO 設定手順

## 概要

devcontainer 環境は揮発するため、SSO 再設定が必要になる場合がある。

## 設定ファイル

`~/.aws/config` （`setup-aws-config.sh` が環境変数から自動生成）:

```ini
[sso-session psm-sandbox]
sso_start_url = <AWS_SSO_START_URL>
sso_region = ap-northeast-1
sso_registration_scopes = sso:account:access

[profile psm]
sso_session = psm-sandbox
sso_account_id = <AWS_SSO_ACCOUNT_ID>
sso_role_name = ai-agent
region = ap-northeast-1
output = json
```

## ログイン手順

```bash
# setup-aws-config.sh が postStartCommand で自動実行済み
aws sso login --sso-session psm-sandbox --use-device-code
aws sts get-caller-identity --profile psm
```

## 環境変数

ホスト側に以下の環境変数をセットすること（`devcontainer.json` の `containerEnv` → `localEnv` 経由でコンテナに注入される）：

| 環境変数                     | 説明                     |
| ---------------------------- | ------------------------ |
| `AWS_SANDBOX_SSO_START_URL`  | SSO ポータル URL         |
| `AWS_SANDBOX_SSO_ACCOUNT_ID` | AWS アカウント ID (12桁) |

以下はリポジトリ内にハードコード済み（秘匿不要）：

- ロール名: `ai-agent`
- リージョン: `ap-northeast-1`
