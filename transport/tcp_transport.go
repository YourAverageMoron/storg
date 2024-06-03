package transport

import (
	"fmt"
	"net"
)

type TCPTransportOpts struct {
	Addr       string
	HandlePeer func(*TCPPeer) error
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
	fmt.Print("sdfsdfsd\n")
	conn, err := net.Dial("tcp", addr)
	peer, err := t.newPeer(conn, true)

	if err != nil {
		return err
	}
	go t.handleConn(peer)

	m := Message{Command: RegisterPeer, Data: []byte("some infomoation here")}
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
		// TODO: THIS NEEDS TO BREAK ON END OF MESSAGE
		// I THINK THIS WILL BE HANDLED BY THE STREAM...
		buf := make([]byte, 1028)
		_, err := conn.Read(buf)
		if err != nil {
			// TODO: HOW ARE WE HANDLING THESE ERRORS
			fmt.Printf("[local: %s] [peer %s] - error: %v \n", t.Addr, conn.RemoteAddr(), err)
		}
		m := TCPMessage{}
		m.UnmarshalBinary(buf)

		switch m.Command {
		case RegisterPeer:
			return t.handleRegisterPeer(m.Data, conn)
		}
	}
}

func (t *TCPTransport) handleRegisterPeer(payload []byte, conn net.Conn) error {
    fmt.Println(string(payload[:]))
    fmt.Println(conn.RemoteAddr())
    fmt.Println(conn.LocalAddr())
	return nil
}

func (t *TCPTransport) newPeer(conn net.Conn, outbound bool) (*TCPPeer, error) {
	peer := NewTCPPeer(conn, outbound)
	t.peers[peer.RemoteAddr()] = peer
	if err := t.HandlePeer(peer); err != nil {
		return nil, err
	}
	fmt.Printf("[local: %s] [peer: %s] new peer added \n", t.Addr, peer.RemoteAddr())
	return peer, nil
}
