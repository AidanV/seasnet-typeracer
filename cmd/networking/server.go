package cmd

import (
	"fmt"
	"net"
	"os"
	"sort"
	"sync"
	"time"
)

type server struct {
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
	s := server{
		ready:     false,
		startTime: time.Now(),
	}
	go s.server(playerInfos, conn)
	go s.sendPlayerInfosOnInterval(500*time.Millisecond, playerInfos, conn)
}

func (s *server) server(playerInfos *sync.Map, conn *net.UDPConn) {
	defer conn.Close()
	for {
		p := make([]byte, 1024)
		nn, addr, err := conn.ReadFromUDP(p)
		if err != nil {
			fmt.Printf("Read err  %v", err)
			continue
		}
		go s.handlePacket(playerInfos, p, nn, addr)
	}
}

func (s *server) handlePacket(playerInfos *sync.Map, p []byte, nn int, addr *net.UDPAddr) {
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
	if !s.ready && readyToStart {
		s.ready = true
		s.startTime = time.Now()
	}
}

func (s *server) sendPlayerInfosOnInterval(tick time.Duration, playerInfos *sync.Map, conn *net.UDPConn) {
	for range time.Tick(tick) {
		type player struct {
			addr *net.UDPAddr
			pi   PlayerInfo
		}
		players := []player{}
		playerInfos.Range(func(a any, pi any) bool {
			addr, err := net.ResolveUDPAddr("udp", a.(string))
			if err != nil {
				fmt.Println("Failed to resolve udp address")
			}
			players = append(players, player{
				addr: addr,
				pi:   pi.(PlayerInfo),
			})
			return true
		})
		sort.Slice(players, func(i, j int) bool {
			return players[i].pi.PercentCompleted > players[j].pi.PercentCompleted
		})
		pis := []PlayerInfo{}
		addrs := []*net.UDPAddr{}
		for _, p := range players {
			pis = append(pis, p.pi)
			addrs = append(addrs, p.addr)
		}
		bcast := Broadcast{
			Done:        false,
			Started:     s.ready,
			StartTime:   s.startTime,
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
