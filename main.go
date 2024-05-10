package main

import (
	"bytes"
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
		EncKey:            newEncryptionKey(),
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
	server1 := makeServer(":3000", "file_root/server1")
	server2 := makeServer(":4000", "file_root/server2", ":3000")
	server3 := makeServer(":6000", "file_root/server3", ":3000", ":4000")
	go func() {
		log.Fatal(server1.Start())
	}()
	time.Sleep(1 * time.Second)
	go func() {
		log.Fatal(server2.Start())
	}()

	time.Sleep(1 * time.Second)
	go func() {
		log.Fatal(server3.Start())
	}()

	time.Sleep(1 * time.Second)

	for i := 0; i < 30; i++ {
		key := fmt.Sprintf("key_here_%d", i)
		data := bytes.NewReader([]byte("some infomoation here"))
		if err := server1.Store(key, data); err != nil {
			fmt.Println(err)
		}
	}

	key := "key_here"
	data := bytes.NewReader([]byte("some infomoation heresdfsdfdsf"))
	if err := server1.Store(key, data); err != nil {
		fmt.Println(err)
	}

	if err := server1.store.Delete(server1.ID, key); err != nil {
		fmt.Println(err)
	}

	r, err := server1.Get("key_here")
	if err != nil {
		log.Fatal(err)
	}

	b, err := io.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
