package p2p

import "net"

// Peer is represents a remote node
type Peer interface {
	net.Conn
	Send([]byte) error
	CloseStream()
}

// Transport handles commincation between nodes
// E.g TCP, UDP, Websockets
type Transport interface {
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
	Addr() string
}
