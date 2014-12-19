package zabbix_test

import (
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/fujiwara/go-zabbix-get/zabbix"
)

func TestAgent(t *testing.T) {
	done := make(chan bool)
	timeout := 3 * time.Second

	{
		value, err := zabbix.Get("localhost", "agent.ping", timeout)
		if err == nil {
			t.Errorf("agent is not runnig, but not error value:", value)
		}
	}

	go zabbix.RunAgent("localhost", func(key string) (string, error) {
		switch key {
		case "agent.ping":
			log.Println("key", key)
			return "1", nil
		case "agent.uptime":
			log.Println("key", key)
			return "123", nil
		case "timeout":
			log.Println("key", key, "sleeping...")
			time.Sleep(timeout + time.Second)
			log.Println("wake up. response ok!")
			return "ok", nil
		case "shutdown":
			done <- true
			return "", nil
		default:
			return "", fmt.Errorf("not supported")
		}
	})
	time.Sleep(1 * time.Second)

	{
		value, err := zabbix.Get("localhost", "agent.ping", timeout)
		if err != nil {
			t.Error("get agent.ping failed", err)
		}
		if value != "1" {
			t.Error("agent.ping value expected: 1, got:", value)
		}
	}

	{
		value, err := zabbix.Get("localhost", "agent.uptime", timeout)
		if err != nil {
			t.Error("get agent.ping failed", err)
		}
		if value != "123" {
			t.Error("agent.uptime value expected: 123, got:", value)
		}
	}

	{
		value, err := zabbix.Get("localhost", "xxx", timeout)
		if err != nil {
			t.Error("xxx failed", err)
		}
		if value != zabbix.ErrorMessage {
			t.Error("xxx value expected: ", zabbix.ErrorMessage, "got:", value)
		}
	}

	{
		_, err := zabbix.Get("localhost", "timeout", timeout)
		if err == nil {
			t.Error("timeout must be timeouted.", err)
		}
		if _err := err.(*net.OpError); !_err.Timeout() {
			t.Error("err expected i/o timeout. got:", err)
		}
		log.Println("client timeout")
	}

	zabbix.Get("localhost", "shutdown", timeout)
	<-done
}
