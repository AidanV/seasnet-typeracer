package cmd

import (
	"os"

	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
	"golang.org/x/term"
)

// func InitialModel() model {
// 	return model{
// 		// Our to-do list is a grocery list
// 		// choices: []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},

// 		// A map which indicates which choices are selected. We're using
// 		// the  map like a mathematical set. The keys refer to the indexes
// 		// of the `choices` slice, above.
// 		// selected: make(map[int]struct{}),
// 	}
// }

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func InitialModel(profile termenv.Profile, fore termenv.Color) model {
	termWidth, termHeight, _ := term.GetSize(int(os.Stdin.Fd()))
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
			wordsToEnter:  []rune("this is my test string for testing. a typing test that is very tested and among other things."),
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
	}
}
