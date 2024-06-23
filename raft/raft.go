package raft

import (
	"bytes"
	"fmt"
	"net"
	"sync"

	"ryan-jones.io/gastore/transport"
	"ryan-jones.io/gastore/utils"
)

type RaftState int64

const (
  Leader RaftState = iota
  Follower
  Candidate
)

//TODO: RENAME TO RPC
type Message struct {
	From    string
	Payload any
}

//TODO: RENAME TO RPC
type MessageRegisterPeer struct {
	AdvertisedAddr string
	Network        string
}

//TODO: RENAME TO RPC
type MessageHeartbeat struct {
	Foo string
	Bar string
}

//TODO: ADD RPCs for requests and responses in the doc

type RaftServerOpts struct {
	Transport transport.Transport
	RaftNodes *utils.Set[net.Addr]
	transport.Encoder
}

type RaftNode struct {
	RaftServerOpts
	peers    map[net.Addr]transport.Peer
	peerLock sync.Mutex
	mch      chan Message
}

func NewRaftServer(opts RaftServerOpts) *RaftNode {
	peers := make(map[net.Addr]transport.Peer)
	mch := make(chan Message)
	return &RaftNode{
		RaftServerOpts: opts,
		peers:          peers,
		mch:            mch,
	}
}

func (r *RaftNode) Broadcast(m any) error {
	for _, addr := range r.RaftNodes.Iterate() {
		peer, err := r.getPeer(addr)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if err := r.messagePeer(peer, transport.IncomingMessage, m); err != nil {
			return err
		}
	}
	return nil
}

func (r *RaftNode) getPeer(addr net.Addr) (transport.Peer, error) {
	peer, ok := r.peers[addr]
	if !ok {
		fmt.Printf("[local: %s] [peer: %s] node not connected to peer attempting to dial\n", r.Transport.Addr(), addr.String())
		if err := r.Transport.Dial(addr); err != nil {
			return nil, err
		}
		peer, ok = r.peers[addr]
		if !ok {
			return nil, fmt.Errorf("[local: %s] [peer: %s] unable to connect to peer after dialing\n", r.Transport.Addr(), addr.String())
		}
	}
	return peer, nil
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
	go r.Transport.ListenAndAccept()
	r.consumeLoop()
}

func (r *RaftNode) Consume() <-chan Message {
	return r.mch
}

func (r *RaftNode) handleOutboundPeer(p transport.Peer) {
	r.peers[p.AdvertisedAddr()] = p
	m := MessageRegisterPeer{
		AdvertisedAddr: r.Transport.Addr(),
		Network:        r.Transport.Network(),
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
			fmt.Printf("[local: %s] [peer: %s] closing peer - peer does not exist in raft nodes\n", r.Transport.Addr(), addr)
			return p.Close()
		}
		p.SetAdvertisedAddr(addr)
		oldPeer, ok := r.peers[addr]
		if ok {
			fmt.Printf("[local: %s] [old peer: %s] closing old peer before inserting new peer", r.Transport.Addr(), oldPeer.AdvertisedAddr())
			oldPeer.Close()
		}
		r.peers[addr] = p
		fmt.Println("peers:", r.Transport.Addr(), p)
	default:
		fmt.Printf("[local: %s] [peer: %s] closing peer - invalid register peer message\n", r.Transport.Addr(), p.RemoteAddr())
		return p.Close()
	}
	return nil
}

func (r *RaftNode) consumeLoop() {
	for {
		select {
		case rpc := <-r.Transport.Consume():
			var m Message
			r.Encoder.Decode(bytes.NewReader(rpc.Payload), &m)
			r.handleMessage(m)
		}
	}
}

func (r *RaftNode) handleMessage(m Message) {
	switch payload := m.Payload.(type) {
	case MessageHeartbeat:
		r.handleHeartbeat(payload)
	default:
		r.handleNoMessageMatch(m)
	}
}

func (r *RaftNode) handleHeartbeat(m MessageHeartbeat) {
	fmt.Println("heartbeat -", m.Foo, m.Bar)
}

func (r *RaftNode) handleNoMessageMatch(m Message) {
    // NOTE: WE CAN CONSUME THIS FROM THE NEXT LAYER UP (E.G THE FILESEVER) TO SEND AND RECIEVE MESSAGES
	r.mch <- m
}

func (r *RaftNode) messagePeer(p transport.Peer, command transport.Command, m any) error {
	message := Message{
		From:    r.Transport.Addr(),
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

// TODO: PERSISTENT STATE (WHERE TO WRITE THIS - FILE?)
//  currentTerm int
//  votedFor net.Addr
//  log
//  currentTerm int

// TODO: VOLITILE STATE (params in RaftNode object?)
// commitIndex into
// lastApplied

// TODO: VOLITILE STATE (leader only)
// nextIndex
// matchIndex

// QUESTION: HOW DO WE STORE THE LOG? -> LSM TREE -> BINARY FORMAT?

// QUESTION: what data should the RaftNode store?

// QUESTION: what functions/RPCs does the raft node need to implement (AppendEntries, RequestVote...)



 