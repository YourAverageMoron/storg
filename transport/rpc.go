package transport

type RPC struct {
	Command
	Payload []byte
}

type Command byte

const (
	IncomingMessage Command = 0x1
	IncomingStream  Command = 0x2
	RegisterPeer    Command = 0x3
)

type RegisterPeerPayload struct {
	Network string
	Addr    string
}
