package main

import (
	"fmt"
	"log"

	"ryan-jones.io/gas/p2p"
)

func OnPeer(p p2p.Peer) error {
	fmt.Println("some logic here")
    p.Close()
	return nil
	// return fmt.Errorf("failed the openpeer func")
}

func main() {
	fmt.Println("Stuff")
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer:        OnPeer,
	}
	tr := p2p.NewTCPTransport(tcpOpts)

	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("%+v\n", msg)
		}
	}()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	select {}
}
