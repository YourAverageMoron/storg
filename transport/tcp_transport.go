package transport

import (
	"fmt"
	"net"
)

type TcpTransportOpts struct {
	Addr string
}

type TcpTransport struct {
	TcpTransportOpts
}

func NewTcpTransport(opts TcpTransportOpts) *TcpTransport {
	return &TcpTransport{opts}
}

func (t *TcpTransport) handleConn(conn net.Conn) {
	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println(n)
	fmt.Println(buf)

}

func (t *TcpTransport) ListenAndAccept() error {
	ln, err := net.Listen("tcp", t.Addr)
	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
            // TODO: SHOULD THIS FALL OVER OR JUST REGECT THE CONN
            return err
		}
		go t.handleConn(conn)
	}
}
