# mackerel-plugin-pinging

指定したサーバーへ ICMP Ping を複数回送信し、RTT（往復遅延時間）と成功・失敗回数を Mackerel のカスタムメトリクスとして投稿するプラグインです。IPv4 と IPv6 に対応しています。

## インストール

リリースページからバイナリをダウンロードするか、`mkr` を利用してインストールします。

```sh
mkr plugin install kazeburo/mackerel-plugin-pinging
```

ICMP ソケットを作成するため、通常は root 権限が必要です。Mackerel Agent から実行する場合も、実行ユーザーが必要な権限を持つように設定してください。

## 使い方

コマンドラインから直接実行する例です。`--host` には送信先の IP アドレスまたはホスト名を、`--key-prefix` にはメトリクスを識別する任意の文字列を指定します。

```sh
sudo mackerel-plugin-pinging --host 8.8.8.8 --key-prefix googledns
```

主なオプションは次のとおりです。

```
Usage:
  mackerel-plugin-pinging [OPTIONS]

Application Options:
      --host=       Target IP Address to ping
      --timeout=    timeout millisec per ping (default: 1000)
      --interval=   sleep millisec after every ping (default: 10)
      --count=      Count Sending ping (default: 10)
      --key-prefix= Metric key prefix
  -v, --version     Show version

Help Options:
  -h, --help        Show this help message
```

- `--host`: Ping の送信先です。必須です。
- `--key-prefix`: メトリクス名へ含める識別子です。必須です。複数の送信先を監視する場合は、送信先ごとに重複しない値を指定します。
- `--timeout`: 1 回の Ping のタイムアウト時間（ミリ秒）です。既定値は `1000` です。
- `--interval`: 各 Ping の前に待機する時間（ミリ秒）です。既定値は `10` です。
- `--count`: 集計対象として送信する Ping の回数です。既定値は `10` です。

## Mackerel Agent の設定例

`/etc/mackerel-agent/mackerel-agent.conf` に次のように設定します。設定後は Mackerel Agent を再起動してください。

```toml
[plugin.metrics.pinging_google_dns]
command = "/path/to/mackerel-plugin-pinging --host 8.8.8.8 --key-prefix google_dns --count 10 --timeout 1000 --interval 10"
```

複数の宛先を監視する場合は、プラグイン設定を宛先ごとに追加し、それぞれ異なる `--key-prefix` を設定します。

```toml
[plugin.metrics.pinging_cloudflare_dns]
command = "/path/to/mackerel-plugin-pinging --host 1.1.1.1 --key-prefix cloudflare_dns"
```

mackerel-gentを root権限で実行していない環境では `sudo` を使用するか、コマンドにSUIDを設定してください。

## 出力メトリクス

`--key-prefix` に `google_dns` を指定した場合、次のメトリクスを出力します。すべての RTT は、preflight を除く `--count` 回の Ping の結果から算出されます。

| メトリクス名 | 意味 | 単位 |
| --- | --- | --- |
| `pinging.google_dns_rtt_count.success` | 成功した Ping の回数 | 回 |
| `pinging.google_dns_rtt_count.error` | タイムアウトや通信エラーになった Ping の回数 | 回 |
| `pinging.google_dns_rtt_ms.max` | 成功した Ping の RTT の最大値 | ミリ秒 |
| `pinging.google_dns_rtt_ms.min` | 成功した Ping の RTT の最小値 | ミリ秒 |
| `pinging.google_dns_rtt_ms.average` | 成功した Ping の RTT の平均値 | ミリ秒 |
| `pinging.google_dns_rtt_ms.90_percentile` | 成功した Ping の RTT の 90 パーセンタイル | ミリ秒 |

成功した Ping が 1 回もない場合は、`rtt_count.success` と `rtt_count.error` のみを出力し、RTT の統計メトリクスは出力しません。送信先の名前解決または ICMP ソケットの作成に失敗した場合は、`rtt_count.success` に `0`、`rtt_count.error` に `--count` の値を出力して異常終了します。

## preflight

集計を開始する前に、プラグインは送信先へ 1 回 Ping を送る preflight を実行します。これは ICMP 通信経路やソケットの状態を事前に確認するための Ping であり、RTT・成功回数・失敗回数のいずれのメトリクスにも含まれません。

preflight が失敗した場合は標準エラーにエラーを出力しますが、計測は中止せず、続けて `--count` 回の Ping を実行します。そのため、preflight の失敗自体はメトリクスの `rtt_count.error` に加算されません。

## Sample

```
$ sudo mackerel-plugin-pinging --host 8.8.8.8 --key-prefix googledns
pinging.googledns_rtt_count.success     10.000000       1556117540
pinging.googledns_rtt_count.error       0.000000        1556117540
pinging.googledns_rtt_ms.max    11.853529       1556117540
pinging.googledns_rtt_ms.min    9.001526        1556117540
pinging.googledns_rtt_ms.average        10.104696       1556117540
pinging.googledns_rtt_ms.90_percentile  10.919082       1556117540
```

## Install

Please download release page or `mkr plugin install kazeburo/mackerel-plugin-pinging`.