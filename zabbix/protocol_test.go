package zabbix

import (
	"github.com/fujiwara/go-zabbix-get/zabbix"
	"bytes"
	"testing"
)

func TestPacket(t *testing.T) {
	source := []byte("test.data")
	packet := zabbix.Data2Packet(source)
	compare := []byte{
		90, 66, 88, 68,         // "ZBXD"
		1,                      // \x01
		9, 0, 0, 0, 0, 0, 0, 0, // binary int64(9) == len("test.data")
		116, 101, 115, 116,     // "test"
		46,                     // "."
		100, 97, 116, 97,       // "data"
	}
	if ! bytes.Equal(packet, compare) {
		t.Errorf("invalid packet %v <=> %v", packet, compare)
	}
	restored, err := zabbix.Packet2Data(packet)
	if err != nil {
		t.Errorf("must be no error: %v", err)
	}
	if ! bytes.Equal(restored, source) {
		t.Errorf("invalid data %v <=> %v", restored, source)
	}
}

func TestCorruptedPacket(t *testing.T) {
	packet := []byte{
		90, 66, 88, 68,         // "ZBXD"
		0,                      // CORRUPTED!!
		9, 0, 0, 0, 0, 0, 0, 0, // binary int64(9) == len("test.data")
		116, 101, 115, 116,     // "test"
		46,                     // "."
		100, 97, 116, 97,       // "data"
	}
	restored, err := zabbix.Packet2Data(packet)
	if err == nil {
		t.Errorf("must be error")
	}
	t.Logf("restored: %v", restored)
}

func TestStreamBinary(t *testing.T) {
	source := []byte("test.data")
	packet := zabbix.Data2Packet(source)
	buf := bytes.NewReader(packet)
	restored, err := zabbix.Stream2Data(buf)
	if err != nil {
		t.Errorf("must be no error: %v", err)
	}
	if ! bytes.Equal(restored, source) {
		t.Errorf("invalid data %v <=> %v", restored, source)
	}
}

func TestStreamText(t *testing.T) {
	source := []byte("test.data\n")
	reader := bytes.NewReader(source)
	restored, err := zabbix.Stream2Data(reader)
	if err != nil {
		t.Errorf("must be no error: %v", err)
	}
	if ! bytes.Equal(restored, source) {
		t.Errorf("invalid data %v <=> %v", restored, source)
	}
}
