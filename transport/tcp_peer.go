package transport

import "net"

type TCPPeer struct {
	net.Conn
	Outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		Outbound: outbound,
	}
}

func (p *TCPPeer) Send(message Message) error {
	tcp_message := TCPMessage{Message: message}
	b, err := tcp_message.MarshalBinary()
	if err != nil {
		return err
	}
	_, err = p.Conn.Write(b)
	return err
}
