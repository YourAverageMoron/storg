package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"ryan-jones.io/gastore/p2p"
)

type FileServerOpts struct {
	EncKey            []byte
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
	BootstrapNodes    []string
}

type FileServer struct {
	FileServerOpts

	peerLock sync.Mutex
	peers    map[string]p2p.Peer
	store    *Store
	quitch   chan struct{}
}

type Message struct {
	From    string
	Payload any
}

type MessageGetFile struct {
	Key string
}

type MessageStoreFile struct {
	Key  string
	Size int64
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		Root:              opts.StorageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}
	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(storeOpts),
		quitch:         make(chan struct{}),
		peers:          make(map[string]p2p.Peer),
	}
}

func (s *FileServer) Get(key string) (io.Reader, error) {
	if s.store.Has(key) {
		fmt.Printf("[%s] serving file (%s) from local disk\n", s.Transport.Addr(), key)
		_, r, err := s.store.Read(key)
		return r, err
	}
	fmt.Printf("[%s] don't have file (%s) locally, attempting to fetch from network\n", s.Transport.Addr(), key)
	msg := Message{
		Payload: MessageGetFile{
			Key: key,
		},
	}

	if err := s.broadcast(&msg); err != nil {
		return nil, err
	}

	time.Sleep(time.Millisecond * 500)

	for _, peer := range s.peers {
		var fileSize int64
		binary.Read(peer, binary.LittleEndian, &fileSize)
		n, err := s.store.Write(key, io.LimitReader(peer, 21))
		if err != nil {
			return nil, err
		}
		fmt.Printf("[%s]revieved bytes (%d) over the network from (%s)\n", s.Transport.Addr(), n, peer.RemoteAddr().String())
		peer.CloseStream()
	}

	_, r, err := s.store.Read(key)
	return r, err
}

// TODO SDFS

// func (s *FileServer) Remove(key string) error {
//
// }

func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}
	s.bootstrapNetwork()
	s.loop()
	return nil
}

func (s *FileServer) Stop() {
	close(s.quitch)
}

func (s *FileServer) Store(key string, r io.Reader) error {
	// V1 this will store the file on every node on the network
	// TODO: on look into having replication configuration
	var (
		fileBuffer = new(bytes.Buffer)
		tee        = io.TeeReader(r, fileBuffer)
	)
	size, err := s.store.Write(key, tee)
	if err != nil {
		return err
	}

	msg := Message{
		Payload: MessageStoreFile{
			Key:  key,
			Size: size,
		},
	}

	if err := s.broadcast(&msg); err != nil {
		return err
	}

	time.Sleep(time.Millisecond * 1)

	// TODO: (@YourAverageMoron) use a multiwriter here
	for _, peer := range s.peers {
		if err := peer.Send([]byte{byte(p2p.IncomingStream)}); err != nil {
			return err
		}
		n, err := copyEncrypt(s.EncKey, fileBuffer, peer)
        if err != nil {
            return err
        }
		// n, err := io.Copy(peer, fileBuffer)
		// if err != nil {
		// 	return err
		// }
		fmt.Printf("recieved and writen bytes (%d) to (%s) \n", n, peer.RemoteAddr())
	}

	return nil
}

func (s *FileServer) OnPeer(peer p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()
	s.peers[peer.RemoteAddr().String()] = peer
	log.Printf("connected with remote %s", peer.RemoteAddr())
	return nil
}

func (s *FileServer) stream(msg *Message) error {
	peers := []io.Writer{}
	for _, peer := range s.peers {
		peers = append(peers, peer)
	}
	mw := io.MultiWriter(peers...)
	return gob.NewEncoder(mw).Encode(msg)
}

func (s *FileServer) broadcast(msg *Message) error {
	buf := new(bytes.Buffer)

	if err := gob.NewEncoder(buf).Encode(msg); err != nil {
		return err
	}
	for _, peer := range s.peers {
		if err := peer.Send([]byte{byte(p2p.IncomingMessage)}); err != nil {
			return err
		}
		if err := peer.Send(buf.Bytes()); err != nil {
			return err
		}
	}
	return nil
}

func (s *FileServer) bootstrapNetwork() error {
	for _, addr := range s.BootstrapNodes {
		if len(addr) == 0 {
			continue
		}
		log.Println("bootstraping node: ", addr)
		go func(addr string) {
			if err := s.Transport.Dial(addr); err != nil {
				log.Println("dial error: ", err)
			}
		}(addr)
	}
	return nil
}

func (s *FileServer) handleMessage(from string, msg *Message) error {
	fmt.Println(msg)
	switch payload := msg.Payload.(type) {
	case MessageStoreFile:
		return s.handleMessageStoreFile(from, payload)
	case MessageGetFile:
		return s.handleMessageGetFile(from, payload)
	}

	return nil
}

func (s *FileServer) handleMessageGetFile(from string, msg MessageGetFile) error {
	if !s.store.Has(msg.Key) {
		return fmt.Errorf("[%s] file (%s) does not exist on disk\n", s.Transport.Addr(), msg.Key)
	}
	fmt.Printf("[%s] got file (%s) serving over network\n", s.Transport.Addr(), msg.Key)
	fileSize, r, err := s.store.Read(msg.Key)
	if err != nil {
		return err
	}

	if rc, ok := r.(io.ReadCloser); ok {
		defer func() {
			fmt.Println("closing ReadCloser")
			rc.Close()
		}()
	}

	peer, ok := s.peers[from]
	if !ok {
		return fmt.Errorf("[%s] peer (%s) could not be found in peer list\n", s.Transport.Addr(), from)
	}

	peer.Send([]byte{byte(p2p.IncomingStream)})
	binary.Write(peer, binary.LittleEndian, fileSize)
	n, err := io.Copy(peer, r)
	if err != nil {
		return err
	}
	fmt.Printf("[%s] written (%d) bytes to peer (%s)\n", s.Transport.Addr(), n, from)
	return nil
}

func (s *FileServer) handleMessageStoreFile(from string, msg MessageStoreFile) error {
	peer, ok := s.peers[from]
	if !ok {
		return fmt.Errorf("peer (%s) could not be found in peer list", from)
	}
	n, err := s.store.Write(msg.Key, io.LimitReader(peer, msg.Size))
	if err != nil {
		return err
	}
	fmt.Printf("[%s] written (%d) bytes to disk\n", s.Transport.Addr(), n)
	peer.CloseStream()
	return nil
}

func (s *FileServer) loop() {
	defer func() {
		log.Println("file server stopped due to error or user quit aciton")
		s.Transport.Close()
	}()
	for {
		select {
		case rpc := <-s.Transport.Consume():
			var msg Message
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg); err != nil {
				log.Printf("%s - decoding error: %s\n", s.Transport.Addr(), err)
			}
			if err := s.handleMessage(rpc.From, &msg); err != nil {
				log.Printf("%s - handle message error: %s \n", s.Transport.Addr(), err)
			}
		case <-s.quitch:
			return
		}
	}
}

func init() {
	gob.Register(MessageStoreFile{})
	gob.Register(MessageGetFile{})
}
