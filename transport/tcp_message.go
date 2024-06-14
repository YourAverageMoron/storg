package transport

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	VERSION     byte = 1
	HEADER_SIZE      = 4
)

type TCPRPC struct {
	RPC
}

func (t *TCPRPC) MarshalBinary() (data []byte, err error) {
	length := uint16(len(t.Payload))
	lengthData := make([]byte, 2)
	binary.BigEndian.PutUint16(lengthData, length)
	b := make([]byte, 0, HEADER_SIZE+length)
	b = append(b, VERSION)
	b = append(b, byte(t.Command))
	b = append(b, lengthData...)
	b = append(b, t.Payload...)
	return b, nil

}

func (t *TCPRPC) UnmarshalBinary(r io.Reader) error {
	// TODO: THIS WILL NEED TO HANDLE STREAMING DIFFERENTLY
	h := make([]byte, HEADER_SIZE)
	if _, err := r.Read(h); err != nil {
		return err
	}
	if h[0] != VERSION {
		return fmt.Errorf("version mismatch %d != %d\n", h[0], VERSION)
	}
	length := int(binary.BigEndian.Uint16(h[2:]))
	payload := make([]byte, length)
	r.Read(payload)
	t.Command = Command(h[1])
	t.Payload = payload[:]
	return nil
}

func NewTcpRegisterPeerPayload(addr string) *RegisterPeerPayload {
	return &RegisterPeerPayload{
		Addr:    addr,
		Network: "tcp",
	}
}



