package transport

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
)

type TCPTransportOpts struct {
	Port           string
	HandlePeer     func(net.Addr, Peer) error
	AdvertisedAddr string
	Encoder
}

type TCPTransport struct {
	TCPTransportOpts
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	if opts.HandlePeer == nil {
		opts.HandlePeer = func(addr net.Addr, peer Peer) error { return nil }
	}
	t := &TCPTransport{opts}
	return t
}

func (t *TCPTransport) Addr() string {
	return t.Port
}

func (t *TCPTransport) Dial(addr net.Addr) error {
	conn, err := net.Dial(addr.Network(), addr.String())
	if err != nil {
		return err
	}
	peer, err := t.newPeer(conn, true)
	if err != nil {
		return err
	}

	if err := t.HandlePeer(addr, peer); err != nil {
		return err
	}

	go t.handleConn(peer)

	payload := RegisterPeerPayload{
		Addr:    t.AdvertisedAddr,
		Network: "tcp",
	}
	var buf bytes.Buffer
	if err = t.Encoder.Encode(&buf, payload); err != nil {
		return err
	}
	m := Message{Command: RegisterPeer, Payload: buf.Bytes()}
	return peer.Send(m)
}

func (t *TCPTransport) ListenAndAccept() error {
	ln, err := net.Listen("tcp", t.Addr())
	if err != nil {
		return err
	}
	// defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			// TODO: SHOULD THIS FALL OVER OR JUST REGECT THE CONN
			return err
		}
		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) error {
    // TODO: SHOULD THIS BE CLOSED?
	defer conn.Close()
	for {
		m := TCPMessage{}
		m.UnmarshalBinary(conn)
		switch m.Command {
		case RegisterPeer:
			t.handleRegisterPeer(m.Payload, conn)
		case AnotherCommand:
			t.handleHandleAnotherCommand(m.Payload, conn)
		}
	}
}

func (t *TCPTransport) handleHandleAnotherCommand(payload []byte, conn net.Conn) error {
    fmt.Println(payload)
	return nil
}

func (t *TCPTransport) handleRegisterPeer(payload []byte, conn net.Conn) error {
	r := bytes.NewReader(payload)
	data := &RegisterPeerPayload{}
	if err := t.Encoder.Decode(r, data); err != nil {
		return err
	}
	addr := Addr{
		Addr: data.Addr,
		Net:  data.Network,
	}
	peer, err := t.newPeer(conn, false)
	if err != nil {
		return err
	}
	t.HandlePeer(addr, peer)
	return nil
}

func (t *TCPTransport) newPeer(conn net.Conn, outbound bool) (*TCPPeer, error) {
	peer := NewTCPPeer(conn, outbound)
	return peer, nil
}

func init() {
	gob.Register(RegisterPeerPayload{})
}
