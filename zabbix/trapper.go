package zabbix

import (
	"encoding/json"
	"log"
	"net"
)

type TrapperDatum struct {
	Host  string `json:"host"`
	Key   string `json:"key"`
	Value string `json:"value"`
	Clock int64  `json:"clock"`
}

type TrapperRequest struct {
	Request string `json:"request"`
	Clock   int64  `json:"clock"`
	Data    []TrapperDatum
}

type TrapperResponse struct {
	Response string `json:"response"`
	Info     string `json:"info"`
}

func RunTrapperServer(addr string, callback func(TrapperRequest) (TrapperResponse, error)) error {
	log.Println("Starting trapper server on", addr)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("Can't lesten tcp:10051", err)
	}
	log.Println("Ready for connection")
	var conn net.Conn
	for {
		conn, err = listener.Accept()
		if err != nil {
			log.Println("Error accept:", err)
			continue
		}
		log.Println("Accept connction from", conn.RemoteAddr())
		go handleTrapperConn(conn, callback)
	}
	return nil
}

func handleTrapperConn(conn net.Conn, callback func(TrapperRequest) (TrapperResponse, error)) {
	defer conn.Close()
	var request TrapperRequest
	input, err := Stream2Data(conn)
	if err != nil {
		log.Println("request error:", err)
		return
	}
	err = json.Unmarshal(input, &request)
	log.Printf("request: %#v", request)
	if err != nil {
		log.Println("decode request error:", err)
		return
	}

	res, err := callback(request)
	if err != nil {
		log.Println("process callback error", err)
		return
	}
	responseJson, _ := json.Marshal(res)
	conn.Write(Data2Packet(responseJson))
}
