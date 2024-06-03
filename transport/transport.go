package transport

import "net"

type Transport interface {
	ListenAndAccept() error
	Dial(addr string) error
	// TODO: CONSUME - THAT READS FROM A CHANNEL
	// TODO CLOSE
	// TODO ADDR()
}

type Peer interface {
	net.Conn
	Send(b []byte) error
}
