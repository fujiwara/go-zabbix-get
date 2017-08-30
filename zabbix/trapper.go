package zabbix

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

const (
	TrapperResponseSuccess = "success"
	TrapperResponseInfoFmt = "processed: %d; failed: %d; total: %d; seconds spent: %f"
	TrapperRequestString   = "sender data"
)

type TrapperData struct {
	Host  string `json:"host"`
	Key   string `json:"key"`
	Value string `json:"value"`
	Clock int64  `json:"clock"`
}

type TrapperRequest struct {
	Request string        `json:"request"`
	Clock   int64         `json:"clock"`
	Data    []TrapperData `json:"data"`
}

func (r TrapperRequest) ToPacket() []byte {
	data, _ := json.Marshal(r)
	return Data2Packet(data)
}

type TrapperResponse struct {
	Response  string `json:"response"`
	Info      string `json:"info"`
	Proceeded int    `json:""`
	Failed    int    `json:""`
	Total     int    `json:""`
}

func (r TrapperResponse) ToPacket() []byte {
	data, _ := json.Marshal(r)
	return Data2Packet(data)
}

func SendBulk(addr string, req TrapperRequest, timeout time.Duration) (res TrapperResponse, err error) {
	packet := req.ToPacket()

	addr = FillDefaultPort(addr, ServerDefaultPort)
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(timeout))

	_, err = conn.Write(packet)
	if err != nil {
		return
	}

	data, err := Stream2Data(conn)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &res)
	if err != nil {
		return
	}
	return res, nil
}

func Send(addr string, data TrapperData, timeout time.Duration) (TrapperResponse, error) {
	req := TrapperRequest{
		Request: TrapperRequestString,
		Data:    []TrapperData{data},
	}
	return SendBulk(addr, req, timeout)
}

func RunTrapper(addr string, callback func(TrapperRequest) (TrapperResponse, error)) error {
	addr = FillDefaultPort(addr, ServerDefaultPort)
	log.Println("Starting trapper on", addr)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("Can't lesten tcp", addr, err)
	}
	log.Println("Ready for connection")
	var conn net.Conn
	for {
		conn, err = listener.Accept()
		if err != nil {
			log.Println("Error accept:", err)
			continue
		}
		// log.Println("Accept connection from", conn.RemoteAddr())
		go handleTrapperConn(conn, callback)
	}
	return nil
}

func handleTrapperConn(conn net.Conn, callback func(TrapperRequest) (TrapperResponse, error)) {
	defer conn.Close()
	start := time.Now()
	var request TrapperRequest
	input, err := Stream2Data(conn)
	if err != nil {
		log.Println("request error:", err)
		return
	}
	err = json.Unmarshal(input, &request)
	if err != nil {
		log.Println("decode request error:", err)
		return
	}
	if request.Request != TrapperRequestString {
		log.Println("invalid request.Request", request.Request)
		return
	}

	res, err := callback(request)
	if err != nil {
		log.Println("process callback error", err)
		return
	}
	res.Total = res.Proceeded + res.Failed
	if res.Response == "" {
		res.Response = TrapperResponseSuccess
	}
	if res.Info == "" {
		res.Info = fmt.Sprintf(
			TrapperResponseInfoFmt,
			res.Proceeded,
			res.Failed,
			res.Total,
			time.Now().Sub(start).Seconds(),
		)
	}
	conn.Write(res.ToPacket())
}
