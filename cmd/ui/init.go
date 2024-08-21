package cmd

import (
	"os"
	nw "typeracer/cmd/networking"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
	"golang.org/x/term"
)

func (m model) Init() tea.Cmd {
	return nil
}

func InitialModel(profile termenv.Profile, fore termenv.Color) model {
	termWidth, termHeight, _ := term.GetSize(int(os.Stdin.Fd()))
	playerInfo := nw.PlayerInfo{
		Name:             "",
		PercentCompleted: 0,
		Wpm:              0,
		ReadyToStart:     false,
		Disconnecting:    false,
	}
	items := []list.Item{
		item{
			isServer: true,
			title:    "Create a lobby",
		},
		item{
			isServer: false,
			title:    "Join a game",
		},
	}
	l := list.New(items, itemDelegate{}, 20, 6)
	l.SetFilteringEnabled(false)
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	return model{
		width:  termWidth,
		height: termHeight,
		state: Setup{
			list: l,
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
		conn:       Disconnected{},
		playerInfo: playerInfo,
	}
}
