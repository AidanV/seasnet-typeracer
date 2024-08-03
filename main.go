package main

import (
	"fmt"
	"os"
	"typeracer/cmd"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

func main() {
	p := tea.NewProgram(cmd.InitialModel(termenv.ANSI256, termenv.ANSIWhite, 100, 50), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
