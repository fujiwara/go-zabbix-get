package zabbix_test

import (
	"log"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	runTestTrapper()
	runTestAgent()
	time.Sleep(2 * time.Second) // wait for trapper and agent to start
	code := m.Run()
	log.Println("exit", code)
}
