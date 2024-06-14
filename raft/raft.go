package raft

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"

	"ryan-jones.io/gastore/transport"
)

type RaftServerOpts struct {
	transport.Transport
	RaftNodes []net.Addr
	transport.Encoder
}

type RaftNode struct {
	RaftServerOpts
	peers map[net.Addr]transport.Peer
}

type Message struct {
	From    string
	Payload any
}

type MessageRegisterPeer struct {
	AdvertisedAddr string
	Network        string
}

func (r *RaftNode) Broadcast() error {
	for _, addr := range r.RaftNodes {
		peer, ok := r.peers[addr]
		if !ok {
			fmt.Println("unable to retrieve addr", addr, r.peers)
			continue
		}
		if err := r.messagePeer(peer); err != nil {
			panic(err)
		}
	}
	return nil
}

func NewRaftServer(opts RaftServerOpts) *RaftNode {
	peers := make(map[net.Addr]transport.Peer)
	return &RaftNode{opts, peers}
}

// TODO: SHOULD THIS BE A POINTER TO A PEER?
func (r *RaftNode) OnPeer(p transport.Peer, rpc *transport.RPC) error {
	if p.Outbound() {
		r.handleOutboundPeer(p)
	} else {
		r.handleInboundPeer(p, rpc)
	}
	return nil
}

func (r *RaftNode) handleOutboundPeer(p transport.Peer) {
	r.peers[p.AdvertisedAddr()] = p
	m := MessageRegisterPeer{
		AdvertisedAddr: r.Addr(),
        Network: "tcp", //TODO: THIS SHOULDNT BE HARDCODED
	}
	payload := Message{
		// TODO: THIS NEEDS TO BE ADVERT ADDR
		From:    r.Addr(),
		Payload: m,
	}

	b := new(bytes.Buffer)
	r.Encoder.Encode(b, payload)
	message := transport.RPC{
		Command: transport.RegisterPeer,
		Payload: b.Bytes(),
	}
	p.Send(message)
}

func (r *RaftNode) handleInboundPeer(p transport.Peer, rpc *transport.RPC) {
	// TODO: close peer if Advertised Addr is not in nodes map
	fmt.Println(p, rpc)
	var m Message
	r.Encoder.Decode(bytes.NewReader(rpc.Payload), &m)

	switch payload := m.Payload.(type) {
	case MessageRegisterPeer:
        addr := transport.Addr{
            Addr: payload.AdvertisedAddr,
            Net: payload.Network,
        }
        p.SetAdvertisedAddr(addr)
        r.peers[addr] = p
        fmt.Println(r.Addr(), p)
	default:
		fmt.Printf("[local: %s] [peer: %s] closing peer - invalid register peer message\n", r.Addr(), p.AdvertisedAddr())
		p.Close()
	}
}

func (r *RaftNode) Start() {
	go r.ListenAndAccept()
	r.consumeLoop()
}

func (r *RaftNode) consumeLoop() {
	for {
		select {
		case rpc := <-r.Consume():
			fmt.Println("received message", rpc.Payload)
		}
	}
}

func (r *RaftNode) messagePeer(p transport.Peer) error {
	fmt.Printf("[local: %s] [peer: %s] sending message\n", r.Addr(), p.RemoteAddr())
	b := []byte("some message")
	message := transport.RPC{
		Command: transport.IncomingMessage,
		Payload: b,
	}
	return p.Send(message)
}

// QUESTION: HOW DO WE STORE THE LOG? -> LSM TREE -> BINARY FORMAT?

// TODO: STEP 2 -> TRANSPORT NEEDS TO HAVE PEER DESCOVERY

// TODO: WRITE TO LOG

//TODO: STEP 3 -> HEARTBREAT

// TODO: LEADER ELECTION

// TODO: SPLIT VOTE

// TODO: MESSAGE LOGGING AND DISTRIBUTION (LOG REPLICATION)

func init() {
	// TODO: THIS SHOULD BE DEPENDENT ON THE ENCODER
	// ENCODER SHOULD HAVE A FUNCTION THAT TAKES IN A LIST OF STRUCTS TO REGISTER
	gob.Register(MessageRegisterPeer{})
}
