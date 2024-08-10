package cmd

import (
	"fmt"
	"net"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

var Prog *tea.Program

func InitClient(playerInfo PlayerInfo, port int) net.Conn {
	conn, err := net.Dial("udp", fmt.Sprintf("127.0.0.1:%v", port))
	if err != nil {
		fmt.Printf("Dial err %v", err)
		os.Exit(-1)
	}
	go readPlayerInfosOnInterval(500*time.Millisecond, conn)
	go PublishPlayerInfo(
		playerInfo,
		conn,
	)
	return conn
}
func PublishPlayerInfo(playerInfo PlayerInfo, conn net.Conn) {
	msg, err := Serialize(playerInfo)
	if err != nil {
		fmt.Println("Failed to serialize")
		os.Exit(-1)
	}
	if _, err = conn.Write(msg); err != nil {
		fmt.Printf("Write err %v\n", err)
		os.Exit(-2)
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
		bcast, err := DeSerialize[Broadcast](p[:nn])
		if err != nil {
			fmt.Println("Failed to deserialize player infos response")
			os.Exit(-1)
		}
		if Prog == nil {
			fmt.Println("Program was nil")
			os.Exit(15)
		}
		Prog.Send(bcast)
	}
}
