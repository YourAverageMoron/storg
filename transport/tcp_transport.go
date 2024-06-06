package transport

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
)

type TCPTransportOpts struct {
	Addr           string
	HandlePeer     func(*TCPPeer) error
	AdvertisedAddr string
	Encoder
}

type TCPTransport struct {
	TCPTransportOpts
	peers map[net.Addr]*TCPPeer
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	if opts.HandlePeer == nil {
		opts.HandlePeer = func(peer *TCPPeer) error { return nil }
	}
	peers := make(map[net.Addr]*TCPPeer)
	t := &TCPTransport{opts, peers}
	return t
}

func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	peer, err := t.newPeer(conn, true)

	if err != nil {
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
	ln, err := net.Listen("tcp", t.Addr)
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
	defer conn.Close()
	for {
		m := TCPMessage{}
		m.UnmarshalBinary(conn)
		switch m.Command {
		case RegisterPeer:
			return t.handleRegisterPeer(m.Payload, conn)
		}
	}
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
    _, ok := t.peers[addr]
    if ok {
        fmt.Printf("[local: %s] [peer: %s] peer already exists in peer map\n", t.Addr, addr.String())
        return nil
    }
	peer, err := t.newPeer(conn, false)
	if err != nil {
		return err
	}
	t.peers[addr] = peer
	return nil
}

func (t *TCPTransport) newPeer(conn net.Conn, outbound bool) (*TCPPeer, error) {
	peer := NewTCPPeer(conn, outbound)
	t.peers[peer.RemoteAddr()] = peer
	if err := t.HandlePeer(peer); err != nil {
		return nil, err
	}
	fmt.Printf("[local: %s] [peer: %s] new peer added\n", t.Addr, peer.RemoteAddr())
	return peer, nil
}

func init() {
	gob.Register(RegisterPeerPayload{})
}
