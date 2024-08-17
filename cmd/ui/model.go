package cmd

import (
	"net"
	"time"
	nw "typeracer/cmd/networking"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/stopwatch"
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
	conn       net.Conn
	playerInfo nw.PlayerInfo
}

type State interface {
}

type PlayerProg struct {
	prog             progress.Model
	name             string
	percentCompleted uint
}

type Test struct {
	started       bool
	startTime     time.Time
	wpmEachSecond []float64
	inputBuffer   []rune
	wordsToEnter  []rune
	results       nw.Results
	cursor        int
	completed     bool // local test completition
	mistakes      mistakes
	rawInputCnt   int
}

type mistakes struct {
	mistakesAt     map[int]bool
	rawMistakesCnt int
}

type TestStopwatch struct {
	stopwatch stopwatch.Model
	isRunning bool
}
