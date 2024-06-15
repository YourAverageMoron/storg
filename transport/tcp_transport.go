package transport

import (
	"fmt"
	"net"
)

type TCPTransportOpts struct {
	Port           string
	OnPeer         func(Peer, *RPC) error
	AdvertisedAddr string
}

type TCPTransport struct {
	listener net.Listener
	rpcch    chan RPC
	TCPTransportOpts
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	if opts.OnPeer == nil {
		opts.OnPeer = func(peer Peer, data *RPC) error { return nil }
	}
	rpcch := make(chan RPC)
	t := &TCPTransport{TCPTransportOpts: opts, rpcch: rpcch}
	return t
}

func (t *TCPTransport) Addr() string {
	return t.AdvertisedAddr
}

func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

func (t *TCPTransport) Dial(addr net.Addr) error {
	conn, err := net.Dial(addr.Network(), addr.String())
	if err != nil {
		return err
	}
	peer, err := t.newPeer(conn, true)
	if err != nil {
		return err
	}
	peer.AdvertAddr = addr
	if err := t.OnPeer(peer, nil); err != nil {
		return err
	}
	go t.handleConn(peer)
	return nil
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.Addr())
	if err != nil {
		return err
	}
	for {
		conn, err := t.listener.Accept()
		fmt.Printf("[local: %s] [peer: %s] new connection \n", t.Addr(), conn.RemoteAddr())
		if err != nil {
			fmt.Printf("[local: %s] error - %+v \n", t.Addr(), err)
			conn.Close()
		}
		peer, err := t.newPeer(conn, false)
		if err != nil {
			return err
		}
		go t.handleConn(peer)
	}
}

func (t *TCPTransport) Network() string {
	return "tcp"
}

func (t *TCPTransport) handleConn(peer *TCPPeer) error {
	defer peer.Close()
	for {
		rpc := TCPRPC{}
		rpc.UnmarshalBinary(peer)
		switch rpc.Command {
		case IncomingMessage:
			t.handleIncomingMessage(rpc.RPC)
		case IncomingStream:
			t.handleIncomingStream(rpc.RPC, peer)
		case RegisterPeer:
			t.handleRegisterPeer(rpc.RPC, peer)
		}
	}
}

func (t *TCPTransport) handleIncomingMessage(rpc RPC) {
	t.rpcch <- rpc
}

func (t *TCPTransport) handleIncomingStream(rpc RPC, p *TCPPeer) error {
	// TODO: IMPLEMENT STREAMING
	// ALL THIS NEEDS TO DO IS PUT A LOCK ON THE PEER
	fmt.Println(rpc, p)
	return nil
}

func (t *TCPTransport) handleRegisterPeer(rpc RPC, p *TCPPeer) error {
	return t.OnPeer(p, &rpc)
}

func (t *TCPTransport) newPeer(conn net.Conn, outbound bool) (*TCPPeer, error) {
	peer := NewTCPPeer(conn, outbound)
	return peer, nil
}
