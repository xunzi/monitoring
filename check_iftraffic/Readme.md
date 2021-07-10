# check iftraffic plugin

This plugin reads the counters the linux kernel writes to `/proc/net/dev` similar to vnstat, allowing for 
interface monitoring on stock linux systems without the need for snmp.

Uses https://github.com/olorin/nagiosplugin for nagios compatible output formatting.

Documentation of the various fields available at kernel.org https://www.kernel.org/doc/html/latest/networking/statistics.html.

## Sample output of cat /proc/net/dev

```shell
Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
    lo:  137302    1500    0    0    0     0          0         0   137302    1500    0    0    0     0       0          0
enp3s0:       0       0    0    0    0     0          0         0        0       0    0    0    0     0       0          0
wlp2s0: 674584134  620442    0 14017    0     0          0         0 399538989  380073    0    0    0     0       0          0
virbr0:       0       0    0    0    0     0          0         0        0       0    0    0    0     0       0          0
virbr0-nic:       0       0    0    0    0     0          0         0        0       0    0    0    0     0       0          0
```
