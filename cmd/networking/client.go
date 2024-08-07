package cmd

import (
	"fmt"
	"net"
	"os"
	"time"
)

func InitClient(port int) {
	conn, err := net.Dial("udp", fmt.Sprintf("127.0.0.1:%v", port))
	if err != nil {
		fmt.Printf("Dial err %v", err)
		os.Exit(-1)
	}
	go client(conn)
	go readPlayerInfosOnInterval(500*time.Millisecond, conn)
}

func client(conn net.Conn) {
	defer conn.Close()
	for {
		msg, err := Serialize(getPlayerInfo())
		if err != nil {
			fmt.Println("Failed to serialize")
			os.Exit(-1)
		}
		if _, err = conn.Write(msg); err != nil {
			fmt.Printf("Write err %v", err)
			os.Exit(-1)
		}
	}
}

func getPlayerInfo() PlayerInfo {
	return PlayerInfo{
		Name:             "aidan",
		PercentCompleted: 88,
		Wpm:              13,
	}
}

func readPlayerInfosOnInterval(tick time.Duration, conn net.Conn) {
	for range time.Tick(tick) {
		p := make([]byte, 1024)
		nn, err := conn.Read(p)
		if err != nil {
			fmt.Printf("Read err %v\n", err)
			os.Exit(-1)
		}
		playerInfos, err := DeSerializeList(p[:nn])
		if err != nil {
			fmt.Println("Failed to deserialize player infos response")
			os.Exit(-1)
		}
		for _, pi := range playerInfos {
			fmt.Printf("---\nName: %s\nPercent Completed: %d\nWPM: %d\n---\n", pi.Name, pi.PercentCompleted, pi.Wpm)
		}
	}
}
