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

type RaftServerOpts struct {
	Transport transport.Transport
	RaftNodes *utils.Set[net.Addr]
	transport.Encoder
}

type RaftNode struct {
	RaftServerOpts
	peers    map[net.Addr]transport.Peer
	peerLock sync.Mutex
	rpcch    chan RPC
 timeout  *Timeout
}

func NewRaftServer(opts RaftServerOpts) *RaftNode {
	peers := make(map[net.Addr]transport.Peer)
	rpcch := make(chan RPC)
 timeout := NewTimeout(ElectionTimeoutFunc)
	return &RaftNode{
		RaftServerOpts: opts,
		peers:          peers,
		rpcch:          rpcch,
  timeout:        timeout,
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

func (r *RaftNode) Start() {
	r.registerMessages()
	go r.Transport.ListenAndAccept()
 r.timeout.Start()
 //TODO go listenForTimeout (r.timeout.Consume())
	r.consumeLoop()
}

func (r *RaftNode) Consume() <-chan RPC {
	return r.rpcch
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

func (r *RaftNode) handleOutboundPeer(p transport.Peer) {
	r.peers[p.AdvertisedAddr()] = p
	m := RegisterPeerRPC{
		AdvertisedAddr: r.Transport.Addr(),
		Network:        r.Transport.Network(),
	}
	r.messagePeer(p, transport.RegisterPeer, m)
}

func (r *RaftNode) handleInboundPeer(p transport.Peer, rpc *transport.RPC) error {
	var m RPC
	r.Encoder.Decode(bytes.NewReader(rpc.Payload), &m)
	switch payload := m.Payload.(type) {
	case RegisterPeerRPC:
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
			var m RPC
			r.Encoder.Decode(bytes.NewReader(rpc.Payload), &m)
			r.handleMessage(m)
		}
	}
}

func (r *RaftNode) handleMessage(rpc RPC) {
	switch payload := rpc.Payload.(type) {
	case AppendEntriesRPCRequest:
		r.handleAppendEntriesRequest(payload)
 case AppendEntriesRPCResponse:
		r.handleAppendEntriesResponse(payload)
 case RequestVoteRPCRequest:
  r.handleRequestVoteRequest(payload)
 case RequestVoteRPCResponse:
  r.handleRequestVoteResponse(payload)
	default:
		r.handleNoRPCMatch(rpc)
	}
}

func (r *RaftNode) handleAppendEntriesRequest(rpc AppendEntriesRPCRequest){
  //TODO: IMPLEMENT
  r.timeout.Reset()
}

func (r *RaftNode) handleAppendEntriesResponse(rpc AppendEntriesRPCResponse){
  //TODO: IMPLEMENT
}

func (r *RaftNode) handleRequestVoteRequest(rpc RequestVoteRPCRequest){
  //TODO: IMPLEMENT
}

func (r *RaftNode) handleRequestVoteResponse(rpc RequestVoteRPCResponse){
  //TODO: IMPLEMENT
}

func (r *RaftNode) handleNoRPCMatch(rpc RPC) {
    // NOTE: WE CAN CONSUME THIS FROM THE NEXT LAYER UP (E.G THE FILESEVER) TO SEND AND RECIEVE MESSAGES
	r.rpcch <- rpc
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

func (r *RaftNode) messagePeer(p transport.Peer, command transport.Command, m any) error {
	message := RPC{
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
 //TODO: SHOULD THIS RETURN AN ERROR?
	r.Encoder.Register(
		RegisterPeerRPC{},
	)
}

//TODO: ELECTION TIMEOUT
// IF FOLLOWER
// START TIMEOUT (SLEEP?) 150 - 300ms
// IF RECIEVES VALID APPEND ENTRIES RPC RESET TIMEOUT
// IF TIMEOUT EXCEEDED SWITCH TO CANDIDATE AND REQUEST VOTE
// HOW TO WRITE THE TIMEOUT THREAD?



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



 