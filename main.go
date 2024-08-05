package main

import (
	// "fmt"
	// "os"
	"sync"
	client "typeracer/cmd/client"
	server "typeracer/cmd/server"
	// ui "typeracer/cmd/ui"
	// tea "github.com/charmbracelet/bubbletea"
	// "github.com/muesli/termenv"
)

func main() {

	var wg sync.WaitGroup
	go func() {
		server.Server()
		wg.Done()
	}()
	go func() {
		client.PostToServer()
		wg.Done()
	}()
	wg.Add(2)
	// p := tea.NewProgram(ui.InitialModel(termenv.ANSI256, termenv.ANSIWhite), tea.WithAltScreen())
	// if _, err := p.Run(); err != nil {
	// 	fmt.Printf("Alas, there's been an error: %v", err)
	// 	os.Exit(1)
	// }
	wg.Wait()
}
