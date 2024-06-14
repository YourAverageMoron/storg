package transport

import "net"

type Transport interface {
	ListenAndAccept() error
	Dial(addr net.Addr) error
	Consume() <-chan RPC
	Close() error
	Addr() string
}

type Peer interface {
	net.Conn
	Send(RPC) error
	Outbound() bool
	AdvertisedAddr() net.Addr
	SetAdvertisedAddr(net.Addr)
}

type Addr struct {
	Addr string
	Net  string
}

func (addr Addr) String() string {
	return addr.Addr
}

func (addr Addr) Network() string {
	return addr.Net
}
