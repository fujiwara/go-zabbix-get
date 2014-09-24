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
  -o="": output key string (format=sensu only)
```

Output format
-----

"zabbix" (default)

zabbix-get compatible format.
```
[value]\n
```

"sensu"

sensu plugin compatible format.
```
[key]\t[value]\t[unixtime]\n
```

When you use option `-o KEYNAME` with `-f sensu`, go-zabbix-get prints it with a result.

```
$ go-zabbix-get -o listen.http -f sensu -k "net.tcp.listen[80]" -s 127.0.0.1

listen.http\t[value]\t[unixtime]\n
```

LICENCE
-------

The MIT License (MIT)
