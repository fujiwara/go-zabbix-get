package zabbix_test

import (
	"log"
	"net"
	"testing"
	"time"

	"github.com/fujiwara/go-zabbix-get/zabbix"
)

func TestTrapper(t *testing.T) {
	done = make(chan bool)
	go zabbix.RunTrapper("localhost", func(req zabbix.TrapperRequest) (res zabbix.TrapperResponse, err error) {
		switch req.Request {
		case "timeout":
			log.Println("request timeout sleeping...")
			time.Sleep(timeout + time.Second)
			log.Println("wake up")
		case "shutdown":
			done <- true
		}
		for _, data := range req.Data {
			log.Println(data)
		}
		res.Proceeded = len(req.Data)
		return res, nil
	})
	time.Sleep(1 * time.Second)
}

func TestTrapperCannotConnect(t *testing.T) {
	value, err := zabbix.Send(
		zabbix.TrapperData{Host: "localhost", Key: "foo", Value: "bar"},
		"localhost:10049",
		timeout,
	)
	if err == nil {
		t.Errorf("trapper is not runnig, but not error value:", value)
	}
}

func TestTrapperSend(t *testing.T) {
	res, err := zabbix.Send(
		zabbix.TrapperData{Host: "localhost", Key: "foo", Value: "bar"},
		"localhost",
		timeout,
	)
	if err != nil {
		t.Errorf("send failed", err)
	}
	if res.Proceeded != 1 {
		t.Errorf("proceeded expected 1 got", res.Proceeded)
	}
	if res.Failed != 0 {
		t.Errorf("failed expected 0 got", res.Failed)
	}
	if res.Total != 1 {
		t.Errorf("total expected 1 got", res.Total)
	}
}

func TestTrapperSendBulk(t *testing.T) {
	res, err := zabbix.SendBulk(
		zabbix.TrapperRequest{
			Data: []zabbix.TrapperData{
				zabbix.TrapperData{Host: "localhost", Key: "foo", Value: "bar"},
				zabbix.TrapperData{Host: "localhost", Key: "xxx", Value: "yyy"},
			},
		},
		"localhost",
		timeout,
	)
	if err != nil {
		t.Errorf("send failed", err)
	}
	if res.Proceeded != 2 {
		t.Errorf("proceeded expected 2 got", res.Proceeded)
	}
	if res.Failed != 0 {
		t.Errorf("failed expected 0 got", res.Failed)
	}
	if res.Total != 2 {
		t.Errorf("total expected 2 got", res.Total)
	}
}

func TestTrapperSendTimeout(t *testing.T) {
	res, err := zabbix.SendBulk(
		zabbix.TrapperRequest{
			Request: "timeout",
			Data: []zabbix.TrapperData{
				zabbix.TrapperData{Host: "localhost", Key: "foo", Value: "bar"},
			},
		},
		"localhost",
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

func TestTrapperShutdown(t *testing.T) {
	zabbix.SendBulk(
		zabbix.TrapperRequest{Request: "shutdown"},
		"localhost",
		timeout,
	)
	<-done
}
