package cmd

import (
	"fmt"
	"net"
	"os"
)

func PostToServer() {
	conn, err := net.Dial("udp", "127.0.0.1:8000")
	if err != nil {
		fmt.Printf("Dial err %v", err)
		os.Exit(-1)
	}
	defer conn.Close()
	msg := "Hello, UDP server"
	fmt.Printf("Ping: %v\n", msg)
	if _, err = conn.Write([]byte(msg)); err != nil {
		fmt.Printf("Write err %v", err)
		os.Exit(-1)
	}

	p := make([]byte, 1024)
	nn, err := conn.Read(p)
	if err != nil {
		fmt.Printf("Read err %v\n", err)
		os.Exit(-1)
	}

	fmt.Printf("%v\n", string(p[:nn]))
}
