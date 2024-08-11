package cmd

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

type serverState struct {
	ready     bool
	startTime time.Time
}

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
	state := serverState{
		ready:     false,
		startTime: time.Now(),
	}
	go server(playerInfos, conn, &state)
	go sendPlayerInfosOnInterval(500*time.Millisecond, playerInfos, conn, &state)
}

func server(playerInfos *sync.Map, conn *net.UDPConn, state *serverState) {
	defer conn.Close()
	for {
		p := make([]byte, 1024)
		nn, addr, err := conn.ReadFromUDP(p)
		if err != nil {
			fmt.Printf("Read err  %v", err)
			continue
		}
		go handlePacket(playerInfos, p, nn, addr, state)
	}
}

func handlePacket(playerInfos *sync.Map, p []byte, nn int, addr *net.UDPAddr, state *serverState) {
	msg := p[:nn]
	playerInfo, err := DeSerialize[PlayerInfo](msg)
	if err != nil {
		fmt.Println("Failed to deserialize incoming packet")
		os.Exit(-1)
	}
	playerInfos.Store(addr.String(), playerInfo)
	readyToStart := true
	playerInfos.Range(func(key any, val any) bool {
		pi := val.(PlayerInfo)
		if !pi.ReadyToStart {
			readyToStart = false
		}
		return true
	})
	if !state.ready && readyToStart {
		state.ready = true
		state.startTime = time.Now()
	}
}

func sendPlayerInfosOnInterval(tick time.Duration, playerInfos *sync.Map, conn *net.UDPConn, state *serverState) {
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
		bcast := Broadcast{
			Done:        false,
			Started:     state.ready,
			StartTime:   state.startTime,
			Paragraph:   "A Lion lay asleep in the forest, his great head resting on his paws. A timid little Mouse came upon him unexpectedly, and in her fright and haste to get away, ran across the Lion's nose. Roused from his nap, the Lion laid his huge paw angrily on the tiny creature to kill her. \"Spare me!\"",
			PlayerInfos: pis,
		}
		msg, err := Serialize(bcast)
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
