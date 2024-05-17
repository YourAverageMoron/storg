package transport

type Transport interface {
	ListenAndAccept() error
    // TODO: CONSUME - THAT READS FROM A CHANNEL
    // TODO: DIAL
    // TODO CLOSE
    // TODO ADDR()
}
