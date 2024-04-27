package main

import (
	"fmt"
	"log"

	"ryan-jones.io/gastore/p2p"
)

func makeServer(listenAddr, root string, nodes ...string) *FileServer {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}
	transport := p2p.NewTCPTransport(tcpTransportOpts).(*p2p.TCPTransport)
	fileServerOpts := FileServerOpts{
		StorageRoot:       root,
		PathTransformFunc: CASPathTransformFunction,
		Transport:         transport,
		BootstrapNodes:    nodes,
	}
	server := NewFileServer(fileServerOpts)
	transport.OnPeer = server.OnPeer
	return server
}

func OnPeer(p p2p.Peer) error {
	fmt.Println("some logic here")
	p.Close()
	return nil
	// return fmt.Errorf("failed the openpeer func")
}

func main() {
	server := makeServer(":3000", "file_root")
	server2 := makeServer(":4000", "file_root", ":3000")
	go func() {
		log.Fatal(server.Start())
	}()

	server2.Start()
}
