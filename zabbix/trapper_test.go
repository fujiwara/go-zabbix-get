package zabbix_test

import (
	"log"
	"net"
	"testing"
	"time"

	"github.com/fujiwara/go-zabbix-get/zabbix"
)

func runTestTrapper() {
	go zabbix.RunTrapper("localhost", func(req zabbix.TrapperRequest) (res zabbix.TrapperResponse, err error) {
		switch req.Request {
		case "timeout":
			log.Println("request timeout sleeping...")
			time.Sleep(timeout + time.Second)
			log.Println("wake up")
		}
		for _, data := range req.Data {
			log.Println(data)
		}
		res.Proceeded = len(req.Data)
		return res, nil
	})
}

func TestTrapperCannotConnect(t *testing.T) {
	value, err := zabbix.Send(
		"localhost:10049",
		zabbix.TrapperData{Host: "localhost", Key: "foo", Value: "bar"},
		timeout,
	)
	if err == nil {
		t.Errorf("trapper is not runnig, but not error value: %#v", value)
	}
}

func TestTrapperSend(t *testing.T) {
	res, err := zabbix.Send(
		"localhost",
		zabbix.TrapperData{Host: "localhost", Key: "foo", Value: "bar"},
		timeout,
	)
	if err != nil {
		t.Errorf("send failed %s", err)
	}
	if res.Proceeded != 1 {
		t.Errorf("proceeded expected 1 got %d", res.Proceeded)
	}
	if res.Failed != 0 {
		t.Errorf("failed expected 0 got %d", res.Failed)
	}
	if res.Total != 1 {
		t.Errorf("total expected 1 got %d", res.Total)
	}
}

func TestTrapperSendBulk(t *testing.T) {
	res, err := zabbix.SendBulk(
		"localhost",
		zabbix.TrapperRequest{
			Request: "sender data",
			Data: []zabbix.TrapperData{
				{Host: "localhost", Key: "foo", Value: "bar"},
				{Host: "localhost", Key: "xxx", Value: "yyy"},
			},
		},
		timeout,
	)
	if err != nil {
		t.Errorf("send failed %s", err)
	}
	if res.Proceeded != 2 {
		t.Errorf("proceeded expected 2 got %d", res.Proceeded)
	}
	if res.Failed != 0 {
		t.Errorf("failed expected 0 got %d", res.Failed)
	}
	if res.Total != 2 {
		t.Errorf("total expected 2 got %d", res.Total)
	}
}

func TestTrapperSendTimeout(t *testing.T) {
	res, err := zabbix.SendBulk(
		"localhost",
		zabbix.TrapperRequest{
			Request: "timeout",
			Data: []zabbix.TrapperData{
				{Host: "localhost", Key: "foo", Value: "bar"},
			},
		},
		timeout,
	)
	if err == nil {
		t.Error("timeout must be timeouted.", err)
	}
	if _err := err.(*net.OpError); !_err.Timeout() {
		t.Error("err expected i/o timeout. got:", err)
	}
	log.Println("client timeout", res)
}
