package p2p

var (
	IncomingMessage = 0x1
	IncomingStream  = 0x2
)

// RPC hold data being sent over each transport between two nodes
type RPC struct {
	From    string
	Payload []byte
	Stream  bool
}
