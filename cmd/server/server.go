package cmd

import (
	"fmt"
	"net"
	"os"
)

func Server() {
	addr := net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 8000,
	}
	server, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Listen err %v\n", err)
		os.Exit(-1)
	}
	fmt.Printf("Listen at %v\n", addr.String())

	for {
		p := make([]byte, 1024)
		nn, raddr, err := server.ReadFromUDP(p)
		if err != nil {
			fmt.Printf("Read err  %v", err)
			continue
		}

		msg := p[:nn]
		fmt.Printf("Received %v %s\n", raddr, msg)

		go func(conn *net.UDPConn, raddr *net.UDPAddr, msg []byte) {
			_, err := conn.WriteToUDP([]byte(fmt.Sprintf("Pong: %s", msg)), raddr)
			if err != nil {
				fmt.Printf("Response err %v", err)
			}
		}(server, raddr, msg)
	}
}
