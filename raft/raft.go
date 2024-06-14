package raft

import (
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

type MessageRegisterPeer struct {
	AdvertisedAddr string
}

func (r *RaftNode) Broadcast() error {
	for _, peer := range r.peers {
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

func (r *RaftNode) OnPeer(addr net.Addr, p transport.Peer) error {
	fmt.Println(p.Outbound())
	fmt.Println(p.RemoteAddr())
	if p.Outbound() {
		r.peers[p.RemoteAddr()] = p
		// TODO: REFACTOR SEND THE MessageRegisterPeer message
	}
	fmt.Println(r.peers)
	return nil
}

func (r *RaftNode) Start() {
	go r.ListenAndAccept()
	r.consumeLoop()
}

func (r *RaftNode) consumeLoop() {
	for {
		select {
		case rpc := <-r.Consume():
			fmt.Println("received message", rpc)
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
