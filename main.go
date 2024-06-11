package main

import (
	"fmt"
	"net"
	"time"

	"ryan-jones.io/gastore/raft"
	"ryan-jones.io/gastore/transport"
)

func make_server(port string, addrs ...net.Addr) *raft.RaftNode {
	t_opts := transport.TCPTransportOpts{
		Encoder: transport.GobEncoder{},
		Port:    fmt.Sprintf(":%s", port),

		AdvertisedAddr: fmt.Sprintf("localhost:%s", port),
	}
	t := transport.NewTCPTransport(t_opts)
    rs_opts := raft.RaftServerOpts{
        RaftNodes: addrs, 
		Transport: t,
	}
	rs := raft.NewRaftServer(rs_opts)
	t.HandlePeer = rs.OnPeer
	return rs
}

func main() {

    addr_1 := transport.Addr{
        Net: "tcp",
        Addr: "localhost:3001",
    }
    addr_2 := transport.Addr{
        Net: "tcp",
        Addr: "localhost:3002",
    }

    rs_1 := make_server("3001", addr_2)	
    rs_2 := make_server("3002", addr_1)	
    rs_2.Start() 
    rs_1.Start()
    time.Sleep(1 * time.Second) 
    rs_1.Broadcast()
    // time.Sleep(1 * time.Second) 
    time.Sleep(1 * time.Second) 
    rs_2.Broadcast()
    
    select {}

}
