package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fujiwara/go-zabbix-get/zabbix"
)

var (
	Version string
)

func main() {
	var (
		port        int
		server      string
		key         string
		timeout     int
		showVersion bool
		format      string
		outputKey   string
	)
	flag.IntVar(&port, "p", 10050, "port")
	flag.StringVar(&server, "s", "127.0.0.1", "hostname or IP")
	flag.StringVar(&key, "k", "", "key")
	flag.IntVar(&timeout, "t", 30, "timeout")
	flag.BoolVar(&showVersion, "V", false, "show Version")
	flag.StringVar(&format, "f", "zabbix", "output format (zabbix or sensu)")
	flag.StringVar(&outputKey, "o", "", "output key string (format=sensu only)")
	flag.Parse()

	if showVersion {
		fmt.Printf("go-zabbix-get version %s\n", Version)
		os.Exit(255)
	}

	if key == "" {
		flag.PrintDefaults()
		os.Exit(255)
	}

	address := fmt.Sprintf("%s:%d", server, port)
	value, err := zabbix.Get(address, key, time.Duration(timeout)*time.Second)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	switch format {
	case "sensu":
		if outputKey == "" {
			printFormatSensu(key, value)
		} else {
			printFormatSensu(outputKey, value)
		}
	default:
		printFormatZabbix(value)
	}
}

func printFormatZabbix(value string) {
	fmt.Printf("%s\n", value)
}

func printFormatSensu(key string, value string) {
	fmt.Printf("%s\t%s\t%d\n", key, value, time.Now().Unix())
}
