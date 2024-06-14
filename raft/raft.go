package raft

import (
	"bytes"
	"fmt"
	"net"
	"sync"

	"ryan-jones.io/gastore/transport"
	"ryan-jones.io/gastore/utils"
)

type RaftServerOpts struct {
	transport.Transport
	RaftNodes *utils.Set[net.Addr]
	transport.Encoder
}

type RaftNode struct {
	RaftServerOpts
	peers    map[net.Addr]transport.Peer
	peerLock sync.Mutex
}

type Message struct {
	From    string
	Payload any
}

type MessageRegisterPeer struct {
	AdvertisedAddr string
	Network        string
}

type MessageHeartbeat struct {
	Foo string
	Bar string
}

func (r *RaftNode) Broadcast() error {
	for _, addr := range r.RaftNodes.Iterate() {
		peer, ok := r.peers[addr]
		if !ok {
			fmt.Printf("[local: %s] [peer: %s] node not connected to peer attempting to dial\n", r.Addr(), addr.String())
			if err := r.Dial(addr); err != nil {
				return err
			}
			peer, ok = r.peers[addr]
			if !ok {
				fmt.Printf("[local: %s] [peer: %s] unable to connect to peer after dialing\n", r.Addr(), addr.String())
				continue
			}
		}
		m := MessageHeartbeat{
			Foo: "my advert addr",
			Bar:        "some network here",
		}
		if err := r.messagePeer(peer, transport.IncomingMessage, m); err != nil {
			return err
		}
	}
	return nil
}

func NewRaftServer(opts RaftServerOpts) *RaftNode {
	peers := make(map[net.Addr]transport.Peer)
	return &RaftNode{
		RaftServerOpts: opts,
		peers:          peers,
	}
}

func (r *RaftNode) OnPeer(p transport.Peer, rpc *transport.RPC) error {
	r.peerLock.Lock()
	defer r.peerLock.Unlock()
	if p.Outbound() {
		r.handleOutboundPeer(p)
	} else {
		r.handleInboundPeer(p, rpc)
	}
	return nil
}

func (r *RaftNode) Start() {
	r.registerMessages()
	go r.ListenAndAccept()
	r.consumeLoop()
}

func (r *RaftNode) handleOutboundPeer(p transport.Peer) {
	r.peers[p.AdvertisedAddr()] = p
	m := MessageRegisterPeer{
		AdvertisedAddr: r.Addr(),
		Network:        r.Network(),
	}
	r.messagePeer(p, transport.RegisterPeer, m)
}

func (r *RaftNode) handleInboundPeer(p transport.Peer, rpc *transport.RPC) error {
	var m Message
	r.Encoder.Decode(bytes.NewReader(rpc.Payload), &m)
	switch payload := m.Payload.(type) {
	case MessageRegisterPeer:
		addr := transport.Addr{
			Addr: payload.AdvertisedAddr,
			Net:  payload.Network,
		}
		if !r.RaftNodes.Has(addr) {
			fmt.Printf("[local: %s] [peer: %s] closing peer - peer does not exist in raft nodes\n", r.Addr(), addr)
			return p.Close()
		}
		p.SetAdvertisedAddr(addr)
		oldPeer, ok := r.peers[addr]
		if ok {
			fmt.Printf("[local: %s] [old peer: %s] closing old peer before inserting new peer", r.Addr(), oldPeer.AdvertisedAddr())
			oldPeer.Close()
		}
		r.peers[addr] = p
		fmt.Println("peers:", r.Addr(), p)
	default:
		fmt.Printf("[local: %s] [peer: %s] closing peer - invalid register peer message\n", r.Addr(), p.RemoteAddr())
		return p.Close()
	}
	return nil
}

func (r *RaftNode) consumeLoop() {
	for {
		select {
		case rpc := <-r.Consume():
			var m Message
			r.Encoder.Decode(bytes.NewReader(rpc.Payload), &m)
			r.handleMessage(m)
		}
	}
}

func (r *RaftNode) handleMessage(m Message) {
	switch payload := m.Payload.(type) {
	case MessageHeartbeat:
		// TODO HANDLE THESE
		fmt.Println("heartbeat -", payload.Foo, payload.Bar)
	}
}

func (r *RaftNode) messagePeer(p transport.Peer, command transport.Command, m any) error {
	message := Message{
		From:    r.Addr(),
		Payload: m,
	}
	b := new(bytes.Buffer)
	r.Encoder.Encode(b, message)
	rpc := transport.RPC{
		Command: command,
		Payload: b.Bytes(),
	}
	return p.Send(rpc)
}

func (r *RaftNode) registerMessages() {
	r.Encoder.Register(
		MessageRegisterPeer{},
        MessageHeartbeat{},
	)
}

// QUESTION: HOW DO WE STORE THE LOG? -> LSM TREE -> BINARY FORMAT?

// TODO: WRITE TO LOG

//TODO: HEARTBREAT

// TODO: LEADER ELECTION

// TODO: SPLIT VOTE

// TODO: MESSAGE LOGGING AND DISTRIBUTION (LOG REPLICATION)
