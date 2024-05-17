package transport

import "net"

type TCPPeer struct {
    // 
    net.Conn
}

func NewTCPPeer(conn net.Conn) *TCPPeer {
	return &TCPPeer{conn}
}
