package transport

type Message struct {
	Command 
	Data    []byte
}


type Command byte

const (
    RegisterPeer Command = 0x1 
)
