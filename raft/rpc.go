package raft

type RPC struct {
	From    string
	Payload any
}

type RegisterPeerRPC struct {
	AdvertisedAddr string
	Network        string
}