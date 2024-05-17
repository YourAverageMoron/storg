package main

import (
	"fmt"

	"ryan-jones.io/gastore/transport"
)

func main() {
	fmt.Println("something")

	opts := transport.TCPTransportOpts{Addr: ":3001"}
    s1 := transport.NewTCPTransport(opts)

    s1.ListenAndAccept()
}
