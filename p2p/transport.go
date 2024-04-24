package p2p

// Peer is represents a remote node
type Peer interface {
}

// Transport handles commincation between nodes
// E.g TCP, UDP, Websockets
type Transport interface {
	ListenAndAccept() error
}
