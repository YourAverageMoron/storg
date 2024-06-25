package raft

type RPC struct {
	From    string
	Payload any
}

type RegisterPeerRPC struct {
	AdvertisedAddr string
	Network        string
}

type LogEntry struct {
 //TODO: WHAT SHOULD BE STORED IN THIS?
}

type AppendEntriesRPCRequest struct {
 Term int32
 LeaderId net.Addr
 PrevLogIndex int32
 PrevLogTerm int32
 Entries []LogEntry
 LeaderCommit int32
}

type AppendEntriesRPCResponse struct {
 Term int32
 Success bool
}

type RequestVoteRPCRequest struct {
 Term int32
 CandidateId net.Addr
 LastLogIndex int32
 LastLogTerm int32
}

type RequestVoteRPCResponse struct {
 Term int32
 VoteGranted bool
}