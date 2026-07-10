# mackerel-plugin-pinging

ICMP Ping RTT custom mackerel plugin

## Usage

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

Notice: Root privileges are required to use this plugin.

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