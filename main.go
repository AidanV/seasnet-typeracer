package main

import (
	"fmt"

	"os"
	nw "typeracer/cmd/networking"
	ui "typeracer/cmd/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

var Prog *tea.Program

func main() {
	// serverPtr := flag.Bool("server", false, "run the server")
	// numPtr := flag.Int("port", 8000, "port number")
	// var name string
	// flag.StringVar(&name, "name", "guest"+fmt.Sprint(rand.Intn(5000)), "your name")
	// flag.Parse()
	// if *serverPtr {
	// 	go nw.InitServer(*numPtr)
	// }
	Prog = tea.NewProgram(ui.InitialModel(termenv.ANSI256, termenv.ANSIWhite), tea.WithAltScreen())
	nw.Prog = Prog
	if _, err := Prog.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
