package transport

import "net"

type Transport interface {
	ListenAndAccept() error
	Dial(addr net.Addr) error
	// TODO: CONSUME - THAT READS FROM A CHANNEL
	// TODO CLOSE
	Addr() string
}

type Peer interface {
	net.Conn
	Send(Message) error
}

type Addr struct {
	Addr string
	Net     string
}

func (addr Addr) String() string {
	return addr.Addr
}

func (addr Addr) Network() string {
	return addr.Net
}
