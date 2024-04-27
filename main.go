package main

import (
	"fmt"
	"log"
	"time"

	"ryan-jones.io/gastore/p2p"
)

func OnPeer(p p2p.Peer) error {
	fmt.Println("some logic here")
	p.Close()
	return nil
	// return fmt.Errorf("failed the openpeer func")
}

func main() {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}
	transport := p2p.NewTCPTransport(tcpTransportOpts)
	fileServerOpts := FileServerOpts{
		StorageRoot:       "3000_network",
		PathTransformFunc: CASPathTransformFunction,
		Transport:         transport,
	}
	server := NewFileServer(fileServerOpts)

	go func() {
		time.Sleep(time.Second * 3)
		server.Stop()
	}()
	if err := server.Start(); err != nil {

		log.Fatal(err)
	}
}
