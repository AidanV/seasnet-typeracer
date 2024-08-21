package cmd

import (
	"fmt"
	"io"
	"strings"
	"time"
	nw "typeracer/cmd/networking"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

type StringStyle func(string) termenv.Style

type Styles struct {
	correct      StringStyle
	toEnter      StringStyle
	mistakes     StringStyle
	cursor       StringStyle
	runningTimer StringStyle
	stoppedTimer StringStyle
	greener      StringStyle
	faintGreen   StringStyle
}

type model struct {
	state      State
	styles     Styles
	width      int
	height     int
	progresses []PlayerProg
	conn       ConnOption
	playerInfo nw.PlayerInfo
}

type ConnOption interface {
}

type Disconnected struct {
}

type State interface {
}

type PlayerProg struct {
	prog             progress.Model
	name             string
	percentCompleted uint
}

type Setup struct {
	list list.Model
}

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item struct {
	isServer bool
	title    string
}

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.title)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func (i item) FilterValue() string { return "" }

type Lobby struct {
}

type Test struct {
	startTime     time.Time
	wpmEachSecond []float64
	inputBuffer   []rune
	wordsToEnter  []rune
	cursor        int
	completed     bool // local test completition
	mistakes      mistakes
	rawInputCnt   int
}

type Results struct {
	results nw.Results
}

type mistakes struct {
	mistakesAt     map[int]bool
	rawMistakesCnt int
}

type TestStopwatch struct {
	stopwatch stopwatch.Model
	isRunning bool
}
