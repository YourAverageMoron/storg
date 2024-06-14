package transport

import "net"

type TCPPeer struct {
	AdvertAddr net.Addr
	net.Conn
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		outbound: outbound,
	}
}

func (p *TCPPeer) AdvertisedAddr() net.Addr {
	return p.AdvertAddr
}

func (p *TCPPeer) SetAdvertisedAddr(addr net.Addr) {
    p.AdvertAddr = addr
}

func (p *TCPPeer) Outbound() bool {
	return p.outbound
}

func (p *TCPPeer) Send(rpc RPC) error {
	tcp_message := TCPRPC{RPC: rpc}
	b, err := tcp_message.MarshalBinary()
	if err != nil {
		return err
	}
	_, err = p.Conn.Write(b)
	return err
}
