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