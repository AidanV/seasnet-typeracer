package cmd

import (
	"time"

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
	test   Test
	styles Styles
	width  int
	height int
	prog   progress.Model
}

type Test struct {
	stopwatch     TestStopwatch
	wpmEachSecond []float64
	inputBuffer   []rune
	wordsToEnter  []rune
	results       Results
	cursor        int
	completed     bool
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

type Results struct {
	wpm      int
	accuracy float64
	time     time.Duration
}
