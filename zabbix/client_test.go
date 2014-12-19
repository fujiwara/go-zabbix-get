package zabbix_test

import (
	"testing"

	"github.com/fujiwara/go-zabbix-get/zabbix"
)

func TestIncludesPortNumber(t *testing.T) {
	var addrWithPort = []string{
		"example.com:10050",
		"192.0.2.1:10050",
		"[2001:db8::dead:beef]:10050",
	}
	for _, addr := range addrWithPort {
		if zabbix.FillDefaultPort(addr, zabbix.AgentDefaultPort) != addr {
			t.Errorf("%s includes port number", addr)
		}
	}
}

func TestNotIncludesPortNumber(t *testing.T) {
	var addrWithoutPort = []string{
		"example.com",
		"192.0.2.1",
		"[2001:db8::dead:beef]",
	}
	for _, addr := range addrWithoutPort {
		if zabbix.FillDefaultPort(addr, zabbix.AgentDefaultPort) != addr+":10050" {
			t.Errorf("%s not includes port number", addr)
		}
	}
}
