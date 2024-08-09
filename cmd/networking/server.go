package cmd

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

func InitServer(port int) {
	playerInfos := new(sync.Map)
	addr := net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: port,
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Listen err %v\n", err)
		os.Exit(-1)
	}
	go server(playerInfos, conn)
	go sendPlayerInfosOnInterval(500*time.Millisecond, playerInfos, conn)
}

func server(playerInfos *sync.Map, conn *net.UDPConn) {
	defer conn.Close()
	for {
		p := make([]byte, 1024)
		nn, addr, err := conn.ReadFromUDP(p)
		if err != nil {
			fmt.Printf("Read err  %v", err)
			continue
		}
		go handlePacket(playerInfos, p, nn, addr)
	}
}

func handlePacket(playerInfos *sync.Map, p []byte, nn int, addr *net.UDPAddr) {
	msg := p[:nn]
	playerInfo, err := DeSerialize[PlayerInfo](msg)
	if err != nil {
		fmt.Println("Failed to deserialize incoming packet")
		os.Exit(-1)
	}
	playerInfos.Store(addr.String(), playerInfo)
}

func sendPlayerInfosOnInterval(tick time.Duration, playerInfos *sync.Map, conn *net.UDPConn) {
	for range time.Tick(tick) {
		addrs := []*net.UDPAddr{}
		pis := []PlayerInfo{}
		playerInfos.Range(func(a any, pi any) bool {
			addr, err := net.ResolveUDPAddr("udp", a.(string))
			if err != nil {
				fmt.Println("Failed to resolve udp address")
			}
			addrs = append(addrs, addr)
			pis = append(pis, pi.(PlayerInfo))
			return true
		})
		msg, err := Serialize(pis)
		if err != nil {
			fmt.Println("Failed to serialize response")
			os.Exit(-1)
		}
		for _, a := range addrs {
			_, err = conn.WriteToUDP(msg, a)
			if err != nil {
				fmt.Println("Failed to send response")
				os.Exit(-1)
			}
		}
	}
}
