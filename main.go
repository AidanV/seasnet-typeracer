package main

import (
	"fmt"
	"os"
	client "typeracer/cmd/client"
	server "typeracer/cmd/server"
	ui "typeracer/cmd/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

func main() {
	go server.Server()
	go client.PostToServer()
	p := tea.NewProgram(ui.InitialModel(termenv.ANSI256, termenv.ANSIWhite), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
