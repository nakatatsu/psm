# フォワードプロキシ

## Squid を導入した背景

DevContainer の外向き通信を制限するため、以前は `init-firewall.sh` 内で許可ドメインを `dig` で DNS 解決し、得られた IP を `ipset` に登録していた。

しかし `proxy.golang.org` など GCP (`storage.googleapis.com`) 背後のドメインは IP が頻繁にローテーションするため、ファイアウォール初期化時の IP と実際のアクセス時の IP が乖離し、通信が失敗する問題が発生していた。

代替案（cron 更新、GCP IP レンジ全開放）にも限界があり、**ドメインベースで判断できるレイヤー**として Squid フォワードプロキシを導入した。

## 構成

```
DevContainer (workspace)
  ├── HTTP_PROXY=http://outbound-filter:3128
  ├── HTTPS_PROXY=http://outbound-filter:3128
  ├── NO_PROXY=localhost,127.0.0.1,gh-token-sidecar
  ├── iptables: OUTPUT → outbound-filter:3128 のみ ACCEPT（+ DNS + localhost）
  └── それ以外 DROP

outbound-filter (docker-compose 別サービス, Squid)
  ├── allowed-domains.txt に基づきドメイン単位で allow/deny
  └── 外部への通信は制限なし（Docker デフォルト）
```

- workspace コンテナの iptables は outbound-filter 宛の通信だけを許可する（IP ベースの許可リストは不要）
- どのドメインを通すかの判断は outbound-filter 側に委譲する
- gh-token-sidecar は `NO_PROXY` で除外し、プロキシを経由せず直接通信する

## 許可ドメインの管理

`allowed-domains.txt` に許可ドメインを記載する。先頭の `.` はサブドメインを含むワイルドカードを意味する。

新しいドメインを許可する場合はこのファイルに追記し、Squid コンテナを再起動する。

## なぜ Squid か（他の候補との比較）

| 候補      | 評価                                                                                 |
| --------- | ------------------------------------------------------------------------------------ |
| **Squid** | CONNECT プロキシとしてドメインベース制御が素直にできる。設定が簡潔で運用コストが低い |
| Tinyproxy | 軽量だが HTTPS CONNECT の制御が簡素すぎる。検証環境向き                              |
| Envoy     | 高機能だが設定が複雑。単純な許可リスト型プロキシには過剰                             |
