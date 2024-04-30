package main

import (
	"bytes"
	"fmt"
	"log"
	"time"

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
	fmt.Println("on peer logic here")
	p.Close()
	return nil
	// return fmt.Errorf("failed the openpeer func") TODO
}

func main() {
	server := makeServer(":3000", "file_root/server1")
	server2 := makeServer(":4000", "file_root/server2", ":3000")
	go func() {
		log.Fatal(server.Start())
	}()

    time.Sleep(3 * time.Second)
	go func() {
		log.Fatal(server2.Start())
	}()

    time.Sleep(3 * time.Second)

	data := bytes.NewReader([]byte("some infomoation here"))
    if err := server2.StoreData("key_data_here", data); err != nil {
        fmt.Println(err)
    }

    select {}
}
