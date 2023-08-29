package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/xxuejie/chrome-in-snap-messaging-relayer/common"
)

var port = flag.Int("port", 21212, "Port to connect to")

func main() {
	flag.Parse()

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("127.0.0.1:%d", *port))
	if err != nil {
		log.Fatalf("Error resolving TCP address: %v", err)
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Fatalf("Error creating tcp connection: %v", err)
	}
	err = conn.SetNoDelay(true)
	if err != nil {
		log.Fatalf("Error setting no delay: %v", err)
	}
	defer conn.Close()

	err = common.RelayData(os.Stdin, os.Stdout, conn, conn)
	if err != nil {
		log.Printf("Error relaying data: %v", err)
	}
}
