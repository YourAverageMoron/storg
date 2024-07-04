package main

import (
	// "fmt"
	// "net"
	"fmt"
	"time"

	"ryan-jones.io/storg/raft"
	//
	// "ryan-jones.io/storg/raft"
	// "ryan-jones.io/storg/transport"
	// "ryan-jones.io/storg/utils"
)

//	func make_server(port string, addrs ...net.Addr) *raft.RaftNode {
//		t_opts := transport.TCPTransportOpts{
//			Port:           fmt.Sprintf(":%s", port),
//			AdvertisedAddr: fmt.Sprintf("localhost:%s", port),
//		}
//		nodes := utils.NewSet[net.Addr]()
//		nodes.AddMulti(addrs...)
//
//		t := transport.NewTCPTransport(t_opts)
//		rs_opts := raft.RaftServerOpts{
//			Encoder:   transport.GobEncoder{},
//			RaftNodes: nodes,
//			Transport: t,
//		}
//		rs := raft.NewRaftServer(rs_opts)
//		t.OnPeer = rs.OnPeer
//		return rs
//	}
func main() {
	d := func() time.Duration {
		return time.Duration(10) * time.Second
	}
	timeout := raft.NewTimeout(d)

	timeout.Start()
	<-timeout.Consume()
	fmt.Println("something else")
	go func() {
		time.Sleep(time.Second * time.Duration(8))
		fmt.Println("resetting")
		timeout.Reset()
	}()
	<-timeout.Consume()
	fmt.Println("another thing")
	// 	addr_1 := transport.Addr{
	// 		Net:  "tcp",
	// 		Addr: "localhost:3001",
	// 	}
	// 	addr_2 := transport.Addr{
	// 		Net:  "tcp",
	// 		Addr: "localhost:3002",
	// 	}
	//
	// 	rs_1 := make_server("3001", addr_2)
	// 	rs_2 := make_server("3002", addr_1)
	// 	go rs_2.Start()
	// 	go rs_1.Start()
	//
	select {}

}
