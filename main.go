package main

import (
	// "fmt"
	// "os"
	"sync"
	nw "typeracer/cmd/networking"
	// ui "typeracer/cmd/ui"
	// tea "github.com/charmbracelet/bubbletea"
	// "github.com/muesli/termenv"
)

func main() {

	var wg sync.WaitGroup
	go func() {
		nw.Server()
		wg.Done()
	}()
	go func() {
		nw.SendPlayerInfoToServer()
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
