package cmd

import (
	"fmt"
	"net"
	"os"
)

func InitClient(port int) {
	go client(port)
}

func client(port int) {
	conn, err := net.Dial("udp", fmt.Sprintf("127.0.0.1:%v", port))
	if err != nil {
		fmt.Printf("Dial err %v", err)
		os.Exit(-1)
	}
	playerInfo := PlayerInfo{
		Name:             "aidan",
		PercentCompleted: 88,
		Wpm:              13,
	}
	defer conn.Close()
	msg, err := Serialize(playerInfo)
	if err != nil {
		fmt.Println("Failed to serialize")
		os.Exit(-1)
	}
	if _, err = conn.Write(msg); err != nil {
		fmt.Printf("Write err %v", err)
		os.Exit(-1)
	}

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
