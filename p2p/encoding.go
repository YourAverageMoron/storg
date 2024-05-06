package p2p

import (
	"encoding/gob"
	"errors"
	"io"
)

var ErrInvalidHandshake = errors.New("invalid handshake")

type Decoder interface {
	Decode(io.Reader, *RPC) error
}

type GOBDecoder struct{}

func (dec GOBDecoder) Decode(r io.Reader, rpc *RPC) error {
	return gob.NewDecoder(r).Decode(rpc)
}

type DefaultDecoder struct {
}

func (dec DefaultDecoder) Decode(r io.Reader, rpc *RPC) error {
	peekBuf := make([]byte, 1)
	if _, err := r.Read(peekBuf); err != nil {
		return err
	}

	rpc.Stream = peekBuf[0] == byte(IncomingStream)
	if rpc.Stream {
		return nil
	}

	buf := make([]byte, 1028)
	n, err := r.Read(buf)
	if err != nil {
		return err
	}
	rpc.Payload = buf[:n]
	return nil
}
