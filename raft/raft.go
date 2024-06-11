package raft

import (
	"fmt"
	"net"

	"ryan-jones.io/gastore/transport"
)

type RaftServerOpts struct {
	transport.Transport
	RaftNodes []net.Addr
}

type RaftNode struct {
	RaftServerOpts
	peers map[net.Addr]transport.Peer
}

func NewRaftServer(opts RaftServerOpts) *RaftNode {
	peers := make(map[net.Addr]transport.Peer)
	return &RaftNode{opts, peers}
}

func (r *RaftNode) OnPeer(addr net.Addr, p transport.Peer) error {
	fmt.Println(p.RemoteAddr())
	peer, ok := r.peers[addr]
	if ok {
		fmt.Printf("[local: %s] [peer: %s] peer already exists in peer map updating peer\n", r.Addr(), addr.String())
		peer.Close()
	}
	r.peers[addr] = p
	fmt.Printf("[local: %s] [peer: %s] new peer added\n", r.Addr(), addr.String())
	return nil
}

func (r *RaftNode) Start() {
	go r.ListenAndAccept()
}

func (r *RaftNode) messagePeer(p transport.Peer) error {
	fmt.Printf("[local: %s] [peer: %s] sending message\n", r.Addr(), p.RemoteAddr())
	b := []byte("some message")
	message := transport.Message{
		Command: transport.AnotherCommand,
		Payload: b,
	}
	return p.Send(message)
}

func (r *RaftNode) Broadcast() error {
	for _, addr := range r.RaftNodes {
		peer, ok := r.peers[addr]
		if !ok {
			if err := r.Dial(addr); err != nil {
				return err
			}
			peer, ok = r.peers[addr]
			if !ok {
				return fmt.Errorf("unable to connect to peer")
			}
		}
        if err := r.messagePeer(peer); err != nil{
            panic(err)
        }
	}
	return nil
}

// QUESTION: HOW DO WE STORE THE LOG? -> LSM TREE -> BINARY FORMAT?

// TODO: STEP 2 -> TRANSPORT NEEDS TO HAVE PEER DESCOVERY

// TODO: WRITE TO LOG

//TODO: STEP 3 -> HEARTBREAT

// TODO: LEADER ELECTION

// TODO: SPLIT VOTE

// TODO: MESSAGE LOGGING AND DISTRIBUTION (LOG REPLICATION)
