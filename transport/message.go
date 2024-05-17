package transport

type Message struct {
	Command byte
	Data    []byte
}
