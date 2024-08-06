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
	defer server.Close()
	fmt.Printf("Listen at %v\n", addr.String())

	for {
		p := make([]byte, 1024)
		nn, _, err := server.ReadFromUDP(p)
		if err != nil {
			fmt.Printf("Read err  %v", err)
			continue
		}

		msg := p[:nn]
		playerInfo, err := DeSerialize(msg)
		if err != nil {
			fmt.Println("Failed to deserialize")
			os.Exit(-1)
		}
		fmt.Printf("Name: %s\nPercent Completed: %d\nWPM: %d", playerInfo.Name, playerInfo.PercentCompleted, playerInfo.Wpm)
	}
}
