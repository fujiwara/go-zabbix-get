package zabbix

import (
	"log"
	"net"
	"strings"
)

func RunAgent(addr string, callback func(string) (string, error)) error {
	addr = FillDefaultPort(addr, AgentDefaultPort)
	log.Println("Starting agent on", addr)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("Can't lesten", addr, err)
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
		go handleAgentConn(conn, callback)
	}
	return nil
}

func handleAgentConn(conn net.Conn, callback func(string) (string, error)) {
	var value string
	defer conn.Close()

	key, err := Stream2Data(conn)
	if err != nil {
		log.Println("request error:", err)
		conn.Write(Data2Packet(ErrorMessageBytes))
		return
	}
	keyStr := strings.TrimRight(string(key), "\n")
	value, err = callback(keyStr)
	if err != nil {
		log.Println("process callback error:", err)
		value = ErrorMessage
		conn.Write(Data2Packet(ErrorMessageBytes))
		return
	}
	conn.Write(Data2Packet([]byte(value)))
}
