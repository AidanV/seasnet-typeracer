package main

import (
	"fmt"
	"os"
	"typeracer/cmd"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(cmd.InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
