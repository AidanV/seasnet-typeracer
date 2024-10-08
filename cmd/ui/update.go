package cmd

import (
	"net"
	nw "typeracer/cmd/networking"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var commands []tea.Cmd

	switch state := m.state.(type) {
	case Setup:
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			state.list.SetWidth(msg.Width)
			return m, nil

		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "q", "ctrl+c":
				return m, tea.Quit

			case "enter":
				i, ok := state.list.SelectedItem().(item)
				if ok {
					if i.isServer {
						go nw.InitServer(8000)
					}
					m.conn = nw.InitClient(m.playerInfo, 8000)
				}
				m.state = Lobby{}
				return m, nil
			}
		}

		var cmd tea.Cmd
		state.list, cmd = state.list.Update(msg)
		m.state = state
		return m, cmd
	case Lobby:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				if !m.playerInfo.ReadyToStart {
					m.playerInfo.ReadyToStart = true
				}
			case "ctrl+c", "esc":
				switch conn := m.conn.(type) {
				case Disconnected:
				case net.Conn:
					defer conn.Close()
					m.playerInfo.Disconnecting = true
					nw.PublishPlayerInfo(m.playerInfo, conn)
				}
				return m, tea.Quit
			}
		case nw.Broadcast:
			if msg.Started {
				m.state = Test{
					startTime:     msg.StartTime,
					wpmEachSecond: []float64{},
					inputBuffer:   []rune(""),
					wordsToEnter:  []rune(msg.Paragraph),
					cursor:        0,
					completed:     false,
					mistakes: mistakes{
						mistakesAt:     map[int]bool{},
						rawMistakesCnt: 0,
					},
					rawInputCnt: 0,
				}
			}
		}
	case Test:
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

		case tea.KeyMsg:

			switch msg.String() {
			case "enter":

			case "ctrl+c", "esc":
				switch conn := m.conn.(type) {
				case Disconnected:
				case net.Conn:
					defer conn.Close()
					m.playerInfo.Disconnecting = true
					nw.PublishPlayerInfo(m.playerInfo, conn)
				}
				return m, tea.Quit

			case "backspace", "ctrl+h":
				state.handleBackspace()

			case "ctrl+w":
				state.handleCtrlW()

			case " ":
				if len(state.inputBuffer) < len(state.wordsToEnter) {
					state.handleSpace()
				}

			default:
				switch msg.Type {
				case tea.KeyRunes:
					if len(state.inputBuffer) < len(state.wordsToEnter) {
						state.handleRunes(msg)
					}
				}
			}

		case nw.Broadcast:
			state.wordsToEnter = []rune(msg.Paragraph)
			m.progresses = []PlayerProg{}
			state.startTime = msg.StartTime
			if msg.Results.Done {
				m.state = Results{
					results: msg.Results,
				}
				return m, tea.Batch(commands...)
			}
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

		m.playerInfo.PercentCompleted = uint(100.0 * float64(state.calculateNumCorrect()) / float64(len(state.wordsToEnter)))
		m.playerInfo.Wpm = uint(state.calculateNormalizedWpm())
		state.completed = m.playerInfo.PercentCompleted == 100

		m.state = state
	case Results:
		switch msg := msg.(type) {
		case tea.KeyMsg:

			switch msg.String() { // TODO: fix unable to  leave
			case "ctrl+c", "esc":
				tea.Quit()
			}

		}
	}

	switch conn := m.conn.(type) {
	case Disconnected:
	case net.Conn:
		nw.PublishPlayerInfo(m.playerInfo, conn)

	}
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

func (test *Test) handleBackspace() {
	test.inputBuffer = dropLastRune(test.inputBuffer)

	//Delete mistakes
	inputLength := len(test.inputBuffer)
	_, ok := test.mistakes.mistakesAt[inputLength]
	if ok {
		delete(test.mistakes.mistakesAt, inputLength)
	}

	test.cursor = inputLength
}

func (test *Test) handleCtrlW() {
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

func (test *Test) handleRunes(msg tea.KeyMsg) {
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

func (test *Test) handleSpace() {
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
