package main

import (
	"fmt"
	"log"

	"ryan-jones.io/gas/p2p"
)


func main() {
	fmt.Println("Stuff")
    tr := p2p.NewTCPTransport(":8000")
    if err := tr.ListenAndAccept(); err != nil {
        log.Fatal(err)
    }
    select {}
}
