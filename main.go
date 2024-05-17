package main

import (
	"fmt"

	"ryan-jones.io/gastore/transport"
)

func main() {
	fmt.Println("something")

	opts := transport.TcpTransportOpts{Addr: ":3001"}
    s1 := transport.NewTcpTransport(opts)

    s1.ListenAndAccept()
}
