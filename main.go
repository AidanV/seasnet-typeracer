package main

import (
	"flag"
	"fmt"

	"os"
	nw "typeracer/cmd/networking"
	ui "typeracer/cmd/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

func main() {
	serverPtr := flag.Bool("server", false, "a bool")
	flag.Parse()
	p := tea.NewProgram(ui.InitialModel(termenv.ANSI256, termenv.ANSIWhite), tea.WithAltScreen())
	if *serverPtr {
		go func() {
			nw.InitServer(8000)
		}()
	}
	go func() {
		nw.InitClient(8000, p)
	}()
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
