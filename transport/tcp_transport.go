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
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	if opts.HandlePeer == nil {
		opts.HandlePeer = func(peer *TCPPeer) error { return nil }
	}
	return &TCPTransport{opts}
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

func (t *TCPTransport) handleConn(conn net.Conn) {
	defer conn.Close()
	peer := NewTCPPeer(conn)
	if err := t.HandlePeer(peer); err != nil {
		// TODO: HOW ARE WE HANDLING THESE ERRORS
		fmt.Println("Error: ", err)
	}
	for {
		// TODO: THIS NEEDS TO BREAK ON END OF MESSAGE
		// I THINK THIS WILL BE HANDLED BY THE STREAM...
		buf := make([]byte, 1028)
		n, err := peer.Read(buf)
		if err != nil {
			// TODO: HOW ARE WE HANDLING THESE ERRORS
			fmt.Println("Error: ", err)
		}
		fmt.Println(n)
		fmt.Printf("Received: %s\n", buf[:n])
	}
}
