package zabbix

import (
	"net"
	"time"
)

func FillDefaultPort (addr string) string {
	_, _, err := net.SplitHostPort(addr)
	if err != nil {
		return addr + ":10050"
	}
	return addr
}

func Get(host string, key string, timeout int) (value []byte, err error) {
	host = FillDefaultPort(host)
	conn, err := net.DialTimeout("tcp", host, time.Duration(timeout)*time.Second)
	if err != nil {
		return
	}
	defer conn.Close()
	msg := Data2Packet([]byte(key))
	_, err = conn.Write(msg)
	if err != nil {
		return
	}
	value, err = Stream2Data(conn)
	return
}
