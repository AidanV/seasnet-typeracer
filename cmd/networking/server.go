package cmd

import (
	"fmt"
	"math/rand"
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
	paragraph   string
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
		paragraph:   "",
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
	s.handleStartRequest()
	s.handleDisconnectRequest(playerInfo, addr.String())
}

func (s *server) handleStartRequest() {
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
		s.paragraph = getRandomParagraph()
	}
}

func (s *server) handleDisconnectRequest(p PlayerInfo, addr string) {
	if p.Disconnecting {
		s.playerInfos.Delete(addr)
	}
}

func (s *server) sendPlayerInfosOnInterval(tick time.Duration) {
	for range time.Tick(tick) {
		pis, addrs := s.getOrderedPlayerInfosAndAddresses()
		bcast := Broadcast{
			Results:     s.getResults(),
			Started:     s.ready,
			StartTime:   s.startTime,
			Paragraph:   s.paragraph,
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

func getRandomParagraph() string {
	paragraphs := []string{
		"The ellipsis (...) is a punctuation mark of mystery and intrigue. It signifies an omission, a pause in thought, or a trailing off of words. It can create suspense, indicate a change in topic, or suggest that there's more to be said. The ellipsis is a versatile tool that can add depth and nuance to your writing.",
		"Hunt and peck (two-fingered typing), also known as Eagle Finger, is a common form of typing in which the typist presses each key individually. Instead of relying on the memorized position of keys, the typist must find each key by sight. Use of this method may also prevent the typist from being able to see what has been typed without glancing away from the keys. Although good accuracy may be achieved, any typing errors that are made may not be noticed immediately due to the user not looking at the screen. There is also the disadvantage that because fewer fingers are used, those that are used are forced to move a much greater distance.",
		"In their early 20s, Emily and Ben set an ambitious goal: to save $100,000 in 5 years for a down payment on their dream home. They created a detailed budget, cut back on discretionary spending, and started side hustles to increase their income. They even tracked their progress on a whiteboard in their living room, marking each $1,000 saved with a star. The journey wasn't easy. There were months when they barely scraped by, and unexpected expenses threatened to derail their plans. But they stayed focused, motivated by their shared dream. After 4 years and 11 months, they reached their goal, a testament to their unwavering commitment and financial discipline.",
		"Professional email communication is essential in today's business world. Use a clear and concise subject line, address the recipient appropriately, and proofread your message before sending. Avoid using excessive exclamation points or emojis in formal emails.",
		"Accountability is a critical component of teamwork. Each team member must take responsibility for their actions, their contributions, and their impact on the team's overall performance. When everyone is accountable, the team can function smoothly, with everyone pulling their weight and working towards a common goal. Accountability also fosters trust, as team members know that they can rely on each other to fulfill their commitments.",
		"The hyphen (-) is the punctuation of connection. It joins words to create compound adjectives (well-known actor), links prefixes to words (pre-existing condition), and even breaks words at the end of a line. It's a small mark that plays a big role in ensuring clarity and readability.",
		"A late 20th century trend in typing, primarily used with devices with small keyboards (such as PDAs and Smartphones), is thumbing or thumb typing. This can be accomplished using one or both thumbs. Similar to desktop keyboards and input devices, if a user overuses keys which need hard presses and/or have small and unergonomic layouts, it could cause thumb tendonitis or other repetitive strain injury.",
	}
	return paragraphs[rand.Intn(len(paragraphs))]
}
