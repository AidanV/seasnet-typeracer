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

var Prog *tea.Program

func main() {
	serverPtr := flag.Bool("server", false, "a bool")
	flag.Parse()
	if *serverPtr {
		go nw.InitServer(8000)
	}
	Prog = tea.NewProgram(ui.InitialModel(termenv.ANSI256, termenv.ANSIWhite, 8000), tea.WithAltScreen())
	nw.Prog = Prog
	if _, err := Prog.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
