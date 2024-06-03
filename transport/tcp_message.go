package transport

import (
	"encoding/binary"
	"fmt"
)

const (
	VERSION     byte = 1
	HEADER_SIZE      = 4
)


type TCPMessage struct {
    Message
}

func (t *TCPMessage) MarshalBinary() (data []byte, err error) {
	length := uint16(len(t.Data))
	lengthData := make([]byte, 2)
	binary.BigEndian.PutUint16(lengthData, length)
	b := make([]byte, 0, HEADER_SIZE+length)
	b = append(b, VERSION)
	b = append(b, t.Command)
	b = append(b, lengthData...)
	b = append(b, t.Data...)
	return b, nil

}

func (t *TCPMessage) UnmarshalBinary(bytes []byte) error {
	if bytes[0] != VERSION {
		return fmt.Errorf("version mismatch %d != %d\n", bytes[0], VERSION)
	}
	length := int(binary.BigEndian.Uint16(bytes[2:]))
	end := HEADER_SIZE + length
	if len(bytes) < end {
		return fmt.Errorf("not enough data to parse packet: expected %d, actual %d", HEADER_SIZE+length, len(bytes))
	}
	t.Command = bytes[1]
	t.Data = bytes[HEADER_SIZE:end]
	return nil
}
