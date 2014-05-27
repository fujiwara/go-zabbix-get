package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"github.com/fujiwara/zabbix-aggregate-agent/zabbix_aggregate_agent"
)

func main() {
	var (
		port int
		server string
		key string
		timeout int
	)
	flag.IntVar(&port, "p", 10050, "port")
	flag.StringVar(&server, "s", "127.0.0.1", "hostname or IP")
	flag.StringVar(&key, "k", "", "key")
	flag.IntVar(&timeout, "t", 30, "timeout")
	flag.Parse()

	if key == "" {
		flag.PrintDefaults()
		os.Exit(255)
	}

	address := fmt.Sprintf("%s:%d", server, port)
	value, err := zabbix_aggregate_agent.Get(address, key, timeout)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	fmt.Printf("%s\n", value)
}
