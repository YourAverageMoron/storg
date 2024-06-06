package main

import (
	"ryan-jones.io/gastore/transport"
)

func main() {

	opts_1 := transport.TCPTransportOpts{
        Encoder: transport.GobEncoder{},
		Addr: ":3001",
	}
	s1 := transport.NewTCPTransport(opts_1)

	opts_2 := transport.TCPTransportOpts{
        Encoder: transport.GobEncoder{},
		Addr: ":3002",
	}
	s2 := transport.NewTCPTransport(opts_2)

	go s1.ListenAndAccept()

    go s2.ListenAndAccept()

    s1.Dial("localhost:3002")
    select {}
}

