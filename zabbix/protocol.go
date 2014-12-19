package zabbix

import (
	"bytes"
	"encoding/binary"
	"errors"
)

const (
	HeaderString     = "ZBXD"
	HeaderLength     = len(HeaderString)
	HeaderVersion    = uint8(1)
	DataLengthOffset = int64(HeaderLength + 1)
	DataLengthSize   = int64(8)
	DataOffset       = int64(DataLengthOffset + DataLengthSize)
	ErrorMessage     = "ZBX_NOTSUPPORTED"
)

var (
	ErrorMessageBytes = []byte(ErrorMessage)
	Terminator        = []byte("\n")
	HeaderBytes       = []byte(HeaderString)
)

type ConnReader interface {
	Read(b []byte) (n int, err error)
}

func Data2Packet(data []byte) []byte {
	buf := new(bytes.Buffer)
	buf.Write(HeaderBytes)
	binary.Write(buf, binary.LittleEndian, HeaderVersion)
	binary.Write(buf, binary.LittleEndian, int64(len(data)))
	buf.Write(data)
	return buf.Bytes()
}

func Packet2Data(packet []byte) (data []byte, err error) {
	var dataLength int64
	if len(packet) < int(DataOffset) {
		err = errors.New("zabbix protocol packet too short")
		return
	}

	// read header
	headBuf := bytes.NewReader(packet[0:DataLengthOffset])
	head := make([]byte, DataLengthOffset)
	_, err = headBuf.Read(head)
	if !bytes.Equal(head[0:HeaderLength], HeaderBytes) || head[HeaderLength] != byte(HeaderVersion) {
		err = errors.New("invalid packet header")
		return
	}

	// read data
	buf := bytes.NewReader(packet[DataLengthOffset:DataOffset])
	err = binary.Read(buf, binary.LittleEndian, &dataLength)
	if err != nil {
		return
	}
	data = packet[DataOffset : DataOffset+dataLength]
	return
}

func Stream2Data(conn ConnReader) (rdata []byte, err error) {
	// read header "ZBXD\x01"
	head := make([]byte, DataLengthOffset)
	_, err = conn.Read(head)
	if err != nil {
		return
	}
	if bytes.Equal(head[0:HeaderLength], HeaderBytes) && head[HeaderLength] == byte(HeaderVersion) {
		rdata, err = parseBinary(conn)
	} else {
		rdata, err = parseText(conn, head)
	}
	return
}

func parseBinary(conn ConnReader) (rdata []byte, err error) {
	// read data length
	var dataLength int64
	err = binary.Read(conn, binary.LittleEndian, &dataLength)
	if err != nil {
		return
	}
	// read data body
	buf := make([]byte, 1024)
	data := new(bytes.Buffer)
	total := 0
	size := 0
	for total < int(dataLength) {
		size, err = conn.Read(buf)
		if err != nil {
			return
		}
		if size == 0 {
			break
		}
		total = total + size
		data.Write(buf[0:size])
	}
	rdata = data.Bytes()
	return
}

func parseText(conn ConnReader, head []byte) (rdata []byte, err error) {
	data := new(bytes.Buffer)
	data.Write(head)
	buf := make([]byte, 1024)
	size := 0
	for {
		// read data while "\n" found
		size, err = conn.Read(buf)
		if err != nil {
			return
		}
		if size == 0 {
			break
		}
		i := bytes.Index(buf[0:size], Terminator)
		if i == -1 {
			// terminator not found
			data.Write(buf[0:size])
			continue
		}
		// terminator found
		data.Write(buf[0 : i+1])
		break
	}
	rdata = data.Bytes()
	return
}
