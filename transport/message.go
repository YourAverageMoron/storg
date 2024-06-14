package transport

type RPC struct {
	Command
	Payload []byte
}

type Command byte

const (
	IncomingMessage Command = 0x1
	IncomingStream  Command = 0x2
)

type RegisterPeerPayload struct {
	Network string
	Addr    string
}

type Message struct {
	From    Addr
	Payload any
}
