package p2p

// RPC hold data being sent over each transport between two nodes
type RPC struct {
	From    string
	Payload []byte
}
