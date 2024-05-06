package main

import (
	"fmt"
	"io"
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

	time.Sleep(1 * time.Second)
	go func() {
		log.Fatal(server2.Start())
	}()

	time.Sleep(1 * time.Second)

	// for i := 0; i < 10; i++ {
	// 	key := fmt.Sprintf("key_here_%d", i)
	// 	data := bytes.NewReader([]byte("some infomoation here"))
	// 	if err := server.Store(key, data); err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	time.Sleep(5 * time.Millisecond)
	// }

	// key := "key_here"
	// data := bytes.NewReader([]byte("some infomoation here"))
	// if err := server.Store(key, data); err != nil {
	// 	fmt.Println(err)
	// }

	r, err := server2.Get("key_here")
	if err != nil {
	    log.Fatal(err)
	}

	b, err := io.ReadAll(r)
	if err != nil {
	    log.Fatal(err)
	}
	fmt.Println(string(b))
}
