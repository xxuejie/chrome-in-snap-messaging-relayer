package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os/exec"

	"github.com/xxuejie/chrome-in-snap-messaging-relayer/common"
)

var port = flag.Int("port", 21212, "Port to listen to")
var program = flag.String("program", "", "Program to relay message to")

func main() {
	flag.Parse()

	if *program == "" {
		log.Fatalf("Please specify a program to use!")
	}

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("127.0.0.1:%d", *port))
	if err != nil {
		log.Fatalf("Error resolving TCP address: %v", err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalf("Error creating tcp listener: %v", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Fatalf("Error accepting connection: %v", err)
		}
		err = conn.SetNoDelay(true)
		if err != nil {
			log.Fatalf("Error setting no delay: %v", err)
		}
		go process(conn)
	}
}

func process(c *net.TCPConn) {
	defer c.Close()

	cmd := exec.Command(*program)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Printf("Error creating stdin pipe: %v", err)
		return
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Error creating stdout pipe: %v", err)
		return
	}
	err = cmd.Start()
	if err != nil {
		log.Printf("Error starting command: %v", err)
		return
	}

	err = common.RelayData(stdout, stdin, c, c)
	if err != nil {
		log.Printf("Error relaying data: %v", err)
	}

	if err := cmd.Process.Kill(); err != nil {
		log.Printf("Error killing process: %v", err)
	}
	if err := cmd.Wait(); err != nil {
		log.Printf("Error waiting process: %v", err)
	}
}
