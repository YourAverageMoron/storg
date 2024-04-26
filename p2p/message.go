package p2p

import "net"

// RPC hold data being sent over each transport between two nodes
type RPC struct {
	From    net.Addr
	Payload []byte
}
