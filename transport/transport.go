package transport

type Transport interface {
	ListenAndAccept() error
}
