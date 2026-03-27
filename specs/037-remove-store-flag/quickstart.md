# Quickstart: --store フラグの除去

## 変更後の使い方

```bash
# sync (--store ssm は不要)
psm sync --dry-run secrets.yaml
psm sync --skip-approve secrets.yaml
psm sync --profile prod secrets.yaml

# export (--store ssm は不要)
psm export output.yaml
psm export --profile prod output.yaml

# SOPS との組み合わせ
sops -d secrets.enc.yaml | psm sync --dry-run /dev/stdin
```

## 移行ガイド

既存のスクリプトやコマンドから `--store ssm` を削除するだけ。

```bash
# Before
psm sync --store ssm --dry-run secrets.yaml

# After
psm sync --dry-run secrets.yaml
```

`--store` を指定するとエラーになります:

```
--store has been removed. SSM is now the default store. Remove --store from your command.
```
