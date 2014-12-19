package zabbix

import (
	"fmt"
	"net"
	"time"
)

const (
	AgentDefaultPort  = 10050
	ServerDefaultPort = 10051
)

func FillDefaultPort(addr string, port int) string {
	_, _, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Sprintf("%s:%d", addr, port)
	}
	return addr
}

func Get(addr string, key string, timeout time.Duration) (value string, err error) {
	addr = FillDefaultPort(addr, AgentDefaultPort)
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(timeout))

	msg := Data2Packet([]byte(key))
	_, err = conn.Write(msg)
	if err != nil {
		return
	}
	_value, err := Stream2Data(conn)
	return string(_value), err
}
