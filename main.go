package main

import (
	"fmt"
	"net"
	"time"

	"ryan-jones.io/gastore/raft"
	"ryan-jones.io/gastore/transport"
	"ryan-jones.io/gastore/utils"
)

func make_server(port string, addrs ...net.Addr) *raft.RaftNode {

	t_opts := transport.TCPTransportOpts{
		Port:           fmt.Sprintf(":%s", port),
		AdvertisedAddr: fmt.Sprintf("localhost:%s", port),
	}
	nodes := utils.NewSet[net.Addr]()
	nodes.AddMulti(addrs...)

	t := transport.NewTCPTransport(t_opts)
	rs_opts := raft.RaftServerOpts{
		Encoder:   transport.GobEncoder{},
		RaftNodes: nodes,
		Transport: t,
	}
	rs := raft.NewRaftServer(rs_opts)
	t.OnPeer = rs.OnPeer
	return rs
}

func main() {
	addr_1 := transport.Addr{
		Net:  "tcp",
		Addr: "localhost:3001",
	}
	addr_2 := transport.Addr{
		Net:  "tcp",
		Addr: "localhost:3002",
	}

	rs_1 := make_server("3001", addr_2)
	rs_2 := make_server("3002", addr_1)
	go rs_2.Start()
	go rs_1.Start()


    err := rs_1.Broadcast()
	if err != nil {
		panic(err)
	}
	time.Sleep(1 * time.Second)
	err = rs_2.Broadcast()
	if err != nil {
		panic(err)
	}

	// err := rs_1.Broadcast()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// // time.Sleep(1 * time.Second)
	// time.Sleep(1 * time.Second)
	// err = rs_1.Broadcast()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	select {}

}
