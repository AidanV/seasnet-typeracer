package cmd

import (
	"fmt"
	"net"
	"os"
)

var playerInfos map[*net.UDPAddr]PlayerInfo

func InitServer(port int) {
	playerInfos = map[*net.UDPAddr]PlayerInfo{}
	go server(port)
}

func server(port int) {
	addr := net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: port,
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
		nn, addr, err := server.ReadFromUDP(p)
		if err != nil {
			fmt.Printf("Read err  %v", err)
			continue
		}
		go handlePacket(p, nn, addr, server)
	}
}

func handlePacket(p []byte, nn int, addr *net.UDPAddr, conn *net.UDPConn) {
	msg := p[:nn]
	playerInfo, err := DeSerialize(msg)
	if err != nil {
		fmt.Println("Failed to deserialize")
		os.Exit(-1)
	}
	playerInfos[addr] = playerInfo
	pis := []PlayerInfo{}
	for _, pi := range playerInfos {
		pis = append(pis, pi)
	}
	msg, err = Serialize(pis)
	if err != nil {
		fmt.Println("Failed to serialize response")
		os.Exit(-1)
	}
	_, err = conn.WriteToUDP(msg, addr)
	if err != nil {
		fmt.Println("Failed to send response")
		os.Exit(-1)
	}
}
