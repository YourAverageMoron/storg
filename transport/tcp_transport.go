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
    // TODO: NEED TO WORK OUT A BETTER WAY OF BOOTSTRAPPING SERVERS...
    // CURRENLY - WHILE THIS SENDS A CONNECTION TO PEERS, THE PORT IT USES WILL BE DIFFERENT
    // THIS MEANS THAT SERVER1 -> SERVER2 MAY USE PORT 50001 EVEN THOUGH IT S1 IS ALSO LISTENING ON PORT 3001
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
    t.handleConn(conn, true)
	return nil
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
		go t.handleConn(conn, false)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {
	defer conn.Close()
	peer, err := t.newPeer(conn, outbound)
	if err != nil {
		fmt.Printf("[local: %s] [peer %s] - error: %v \n", t.Addr, peer.Conn.RemoteAddr(), err)
	}
	for {
		// TODO: THIS NEEDS TO BREAK ON END OF MESSAGE
		// I THINK THIS WILL BE HANDLED BY THE STREAM...
		buf := make([]byte, 1028)
		n, err := peer.Read(buf)
		if err != nil {
			// TODO: HOW ARE WE HANDLING THESE ERRORS
			fmt.Printf("[local: %s] [peer %s] - error: %v \n", t.Addr, peer.Conn.RemoteAddr(), err)
		}
		fmt.Println(n)
		fmt.Printf("[local: %s] [peer %s] - recieved %s \n", t.Addr, peer.Conn.RemoteAddr(), buf[:n])
	}
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
