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
	ready       bool
	startTime   time.Time
	playerInfos *sync.Map
	conn        *net.UDPConn
}

func InitServer(port int) {
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
		ready:       false,
		startTime:   time.Now(),
		playerInfos: new(sync.Map),
		conn:        conn,
	}
	go s.server()
	go s.sendPlayerInfosOnInterval(500 * time.Millisecond)
}

func (s *server) server() {
	defer s.conn.Close()
	for {
		p := make([]byte, 1024)
		nn, addr, err := s.conn.ReadFromUDP(p)
		if err != nil {
			fmt.Printf("Read err  %v", err)
			continue
		}
		go s.handlePacket(p, nn, addr)
	}
}

func (s *server) handlePacket(p []byte, nn int, addr *net.UDPAddr) {
	msg := p[:nn]
	playerInfo, err := DeSerialize[PlayerInfo](msg)
	if err != nil {
		fmt.Println("Failed to deserialize incoming packet")
		os.Exit(-1)
	}
	s.playerInfos.Store(addr.String(), playerInfo)
	readyToStart := true
	s.playerInfos.Range(func(key any, val any) bool {
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

func (s *server) sendPlayerInfosOnInterval(tick time.Duration) {
	for range time.Tick(tick) {
		pis, addrs := s.getOrderedPlayerInfosAndAddresses()
		bcast := Broadcast{
			Results:     s.getResults(),
			Started:     s.ready,
			StartTime:   s.startTime,
			Paragraph:   "Testing string",
			PlayerInfos: pis,
		}
		msg, err := Serialize(bcast)
		if err != nil {
			fmt.Println("Failed to serialize response")
			os.Exit(-1)
		}
		for _, a := range addrs {
			_, err = s.conn.WriteToUDP(msg, a)
			if err != nil {
				fmt.Println("Failed to send response")
				os.Exit(-1)
			}
		}
	}
}

func (s *server) getResults() Results {
	results := Results{
		Done:   false,
		Winner: "",
	}
	s.playerInfos.Range(func(_ any, pi any) bool {
		if pi.(PlayerInfo).PercentCompleted == 100 {
			results = Results{
				Done:   true,
				Winner: pi.(PlayerInfo).Name,
			}
			return false
		}
		return true
	})
	return results
}

func (s *server) getOrderedPlayerInfosAndAddresses() ([]PlayerInfo, []*net.UDPAddr) {
	type player struct {
		addr *net.UDPAddr
		pi   PlayerInfo
	}
	players := []player{}
	s.playerInfos.Range(func(a any, pi any) bool {
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
	return pis, addrs
}
