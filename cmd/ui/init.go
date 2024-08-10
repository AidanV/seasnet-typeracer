package cmd

import (
	"os"
	nw "typeracer/cmd/networking"

	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
	"golang.org/x/term"
)

func (m model) Init() tea.Cmd {
	return nil
}

func InitialModel(profile termenv.Profile, fore termenv.Color, name string, port int) model {
	termWidth, termHeight, _ := term.GetSize(int(os.Stdin.Fd()))
	playerInfo := nw.PlayerInfo{
		Name:             name,
		PercentCompleted: 0,
		Wpm:              0,
	}
	return model{
		width:  termWidth,
		height: termHeight,
		test: Test{
			stopwatch: TestStopwatch{
				stopwatch: stopwatch.New(),
				isRunning: false,
			},
			wpmEachSecond: []float64{},
			inputBuffer:   []rune{},
			wordsToEnter:  []rune(" "),
			cursor:        0,
			completed:     false,
			mistakes: mistakes{
				mistakesAt:     map[int]bool{},
				rawMistakesCnt: 0,
			},
		},
		styles: Styles{
			correct: func(str string) termenv.Style {
				return termenv.String(str).Foreground(fore)
			},
			toEnter: func(str string) termenv.Style {
				return termenv.String(str).Foreground(fore).Faint()
			},
			mistakes: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("1")).Underline()
			},
			cursor: func(str string) termenv.Style {
				return termenv.String(str).Reverse().Bold()
			},
			runningTimer: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("2"))
			},
			stoppedTimer: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("2")).Faint()
			},
			greener: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("6")).Faint()
			},
			faintGreen: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("10")).Faint()
			},
		},
		progresses: []PlayerProg{},
		conn: nw.InitClient(
			playerInfo,
			port,
		),
		playerInfo: playerInfo,
	}
}
