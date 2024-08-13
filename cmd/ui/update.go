package cmd

import (
	nw "typeracer/cmd/networking"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	// "github.com/muesli/termenv"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var commands []tea.Cmd

	switch msg := msg.(type) {

	// Update window size
	case tea.WindowSizeMsg:
		if msg.Width == 0 && msg.Height == 0 {
			return m, nil
		} else {
			m.width = msg.Width
			m.height = msg.Height
			return m, nil
		}

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {
		case "enter":
			if !m.playerInfo.ReadyToStart {
				m.playerInfo.ReadyToStart = true
			}

		// These keys should exit the program.
		case "ctrl+c", "esc":
			defer m.conn.Close()
			return m, tea.Quit
		case "tab":

		case "backspace", "ctrl+h":
			handleBackspace(&m.test)

		case "ctrl+w":
			handleCtrlW(&m.test)

		case " ":
			handleSpace(&m.test)

		default:
			switch msg.Type {
			case tea.KeyRunes:
				handleRunes(msg, &m.test)
			}
		}

	case nw.Broadcast:
		m.test.wordsToEnter = []rune(msg.Paragraph)
		m.progresses = []PlayerProg{}
		m.test.started = msg.Started
		m.test.startTime = msg.StartTime
		for _, pi := range msg.PlayerInfos {
			m.progresses = append(
				m.progresses,
				PlayerProg{
					prog:             progress.New(),
					name:             pi.Name,
					percentCompleted: pi.PercentCompleted,
				},
			)
		}
	}

	m.playerInfo.PercentCompleted = uint(100.0 * float64(m.test.calculateNumCorrect()) / float64(len(m.test.wordsToEnter)))
	m.playerInfo.Wpm = uint(m.test.calculateNormalizedWpm())

	nw.PublishPlayerInfo(m.playerInfo, m.conn)

	return m, tea.Batch(commands...)
}

func (m model) quitOn(msg tea.Msg, strokes ...string) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		for _, elem := range strokes {
			if elem == msg.String() {
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func handleBackspace(test *Test) {
	test.inputBuffer = dropLastRune(test.inputBuffer)

	//Delete mistakes
	inputLength := len(test.inputBuffer)
	_, ok := test.mistakes.mistakesAt[inputLength]
	if ok {
		delete(test.mistakes.mistakesAt, inputLength)
	}

	test.cursor = inputLength
}

func handleCtrlW(test *Test) {
	test.inputBuffer = dropUntilWsIdx(test.inputBuffer, test.findLatestWsIndex())
	bufferLen := len(test.inputBuffer)
	test.cursor = bufferLen

	//Delete mistakes
	var newMistakes map[int]bool = make(map[int]bool, 0)
	for at := range test.mistakes.mistakesAt {
		if at < bufferLen {
			newMistakes[at] = true
		}
	}
	test.mistakes.mistakesAt = newMistakes
}

func dropUntilWsIdx(input []rune, wsIdx int) []rune {
	if wsIdx == 0 {
		return make([]rune, 0)
	} else {
		return input[:wsIdx+1]
	}
}

func handleRunes(msg tea.KeyMsg, test *Test) {
	inputLetter := msg.Runes[len(msg.Runes)-1]

	inputLenDec := len(test.inputBuffer)
	letterToInput := test.wordsToEnter[inputLenDec]

	test.inputBuffer = append(test.inputBuffer, inputLetter)
	test.rawInputCnt += 1

	if letterToInput != inputLetter {
		test.mistakes.mistakesAt[inputLenDec] = true
		test.mistakes.rawMistakesCnt = test.mistakes.rawMistakesCnt + 1
	}

	lenAfterAppend := len(test.inputBuffer)

	// Set cursor
	test.cursor = lenAfterAppend
}

func handleSpace(test *Test) {
	if len(test.inputBuffer) > 0 {
		test.inputBuffer = append(test.inputBuffer, ' ')
		test.cursor = len(test.inputBuffer)
		test.rawInputCnt += 1

		letterToInput := test.wordsToEnter[test.cursor-1]
		inputLetter := test.inputBuffer[test.cursor-1]

		if letterToInput != inputLetter {
			test.mistakes.mistakesAt[test.cursor-1] = true
			test.mistakes.rawMistakesCnt += 1
		}

	}
}

func (test *Test) findLatestWsIndex() int {
	var wsIdx int = 0
	for idx, value := range test.wordsToEnter {
		if idx+1 >= test.cursor {
			break
		}
		if value == ' ' {
			wsIdx = idx
		}
	}

	return wsIdx
}

func dropLastRune(runes []rune) []rune {
	le := len(runes)
	if le != 0 {
		return runes[:le-1]
	} else {
		return runes
	}
}
