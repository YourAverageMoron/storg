package transport

type Message struct {
	Command
	Payload []byte
}

type Command byte

const (
	RegisterPeer Command = 0x1
	AnotherCommand Command = 0x2
)

type RegisterPeerPayload struct {
	Network string
	Addr    string
}
