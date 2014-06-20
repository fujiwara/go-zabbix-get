go-zabbix-get
=============

zabbix-get compatible command (Golang version)

Usage
------

```
$ go-zabbix-get -s <hostname or IP> -p <port> -k <key>

  -k="": key
  -p=10050: port
  -s="127.0.0.1": hostname or IP
  -t=30: timeout
  -f="zabbix": output format (zabbix or sensu)
```

Output format
-----

1. "zabbix" (default)
zabbix-get compatible format.
```
[value]\n
```

2. "sensu"
sensu plugin compatible format.
```
[key]\t[value]\t[unix time]
```

LICENCE
-------

The MIT License (MIT)
