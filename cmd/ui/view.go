package cmd

import (
	"math"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/indent"
	"github.com/muesli/reflow/wordwrap"
)

func (m model) View() string {
	var s string

	var termWidth, termHeight = m.width, m.height
	switch state := m.state.(type) {
	case Lobby:
		if m.playerInfo.ReadyToStart {
			s += lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, "Waiting for others...")
		} else {
			s += lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, "Press Enter to ready up")
		}

	case Test:
		if state.completed {
			s += lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, "Done")
		} else {

			lineLenLimit := termWidth * 3 / 4

			coloredStopwatch := style(time.Since(state.startTime).Round(time.Second).String(), m.styles.runningTimer)

			paragraphView := state.paragraphView(lineLenLimit, m.styles)
			lines := strings.Split(paragraphView, "\n")
			cursorLine := findCursorLine(strings.Split(paragraphView, "\n"), state.cursor)

			linesAroundCursor := strings.Join(getLinesAroundCursor(lines, cursorLine), "\n")

			s += positionVerticaly(termHeight)
			avgLineLen := averageLineLen(lines)
			indentBy := uint(math.Max(0, float64(termWidth/2-avgLineLen/2)))

			s += m.indent(coloredStopwatch, indentBy) + "\n\n" + m.indent(linesAroundCursor, indentBy)

			s += "\n\n\n"
			s += lipgloss.PlaceHorizontal(termWidth, lipgloss.Center, style("ctrl+r to restart, ctrl+q to menu", m.styles.toEnter))
			s += "\n\n"
			for _, prog := range m.progresses {
				s += lipgloss.PlaceHorizontal(termWidth, lipgloss.Center, prog.name) + "\n" + lipgloss.PlaceHorizontal(termWidth, lipgloss.Center, prog.prog.ViewAs(float64(prog.percentCompleted)/100.0)) + "\n"
			}
		}
	case Results:
		s += lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, "The winner is "+state.results.Winner)
	}
	return s
}
func positionVerticaly(termHeight int) string {
	var acc strings.Builder

	for i := 0; i < termHeight/2-3; i++ {
		acc.WriteRune('\n')
	}

	return acc.String()
}
func getLinesAroundCursor(lines []string, cursorLine int) []string {
	cursor := cursorLine

	// 3 lines to show
	if cursorLine == 0 {
		cursor += 3
	} else {
		cursor += 2
	}

	low := int(math.Max(0, float64(cursorLine-1)))
	high := int(math.Min(float64(len(lines)), float64(cursor)))

	return lines[low:high]
}
func findCursorLine(lines []string, cursorAt int) int {
	lenAcc := 0
	cursorLine := 0

	for _, line := range lines {
		lineLen := len(dropAnsiCodes(line))

		lenAcc += lineLen

		if cursorAt <= lenAcc-1 {
			return cursorLine
		} else {
			cursorLine += 1
		}
	}

	return cursorLine
}
func averageLineLen(lines []string) int {
	linesLen := len(lines)
	if linesLen > 1 {
		lines = lines[:linesLen-1] //Drop last line, as it might skew up average length
	}

	return averageStringLen(lines)
}

func averageStringLen(strings []string) int {
	var totalLen int = 0
	var cnt int = 0

	for _, str := range strings {
		currentLen := len([]rune(dropAnsiCodes(str)))
		totalLen += currentLen
		cnt += 1
	}

	if cnt == 0 {
		cnt = 1
	}

	return totalLen / cnt
}

func dropAnsiCodes(colored string) string {
	m := regexp.MustCompile("\x1b\\[[0-9;]*m")

	return m.ReplaceAllString(colored, "")
}

func (m model) indent(block string, indentBy uint) string {
	indentedBlock := indent.String(block, indentBy) // this crashes on small windows

	return indentedBlock
}

func style(str string, style StringStyle) string {
	return style(str).String()
}

func (test *Test) paragraphView(lineLimit int, styles Styles) string {
	paragraph := test.colorInput(styles)
	paragraph += test.colorCursor(styles)
	paragraph += test.colorWordsToEnter(styles)

	wrapped := wrapStyledParagraph(paragraph, lineLimit)

	return wrapped
}

func (test *Test) colorInput(styles Styles) string {
	mistakes := toKeysSlice(test.mistakes.mistakesAt)
	sort.Ints(mistakes)

	var coloredInput strings.Builder

	if len(mistakes) == 0 {

		coloredInput.WriteString(styleAllRunes(test.inputBuffer, styles.correct))

	} else {

		previousMistake := -1

		for _, mistakeAt := range mistakes {
			sliceUntilMistake := test.inputBuffer[previousMistake+1 : mistakeAt]
			mistakeSlice := test.wordsToEnter[mistakeAt : mistakeAt+1]

			coloredInput.WriteString(styleAllRunes(sliceUntilMistake, styles.correct))
			coloredInput.WriteString(style(string(mistakeSlice), styles.mistakes))

			previousMistake = mistakeAt
		}

		inputAfterLastMistake := test.inputBuffer[previousMistake+1:]
		coloredInput.WriteString(styleAllRunes(inputAfterLastMistake, styles.correct))
	}

	return coloredInput.String()
}

func (test *Test) colorCursor(styles Styles) string {
	cursorLetter := test.wordsToEnter[len(test.inputBuffer) : len(test.inputBuffer)+1]

	return style(string(cursorLetter), styles.cursor)
}

func (test *Test) colorWordsToEnter(styles Styles) string {
	if len(test.inputBuffer) >= len(test.wordsToEnter) {
		return ""
	}

	wordsToEnter := test.wordsToEnter[len(test.inputBuffer)+1:] // without cursor

	return style(string(wordsToEnter), styles.toEnter)
}

func wrapStyledParagraph(paragraph string, lineLimit int) string {
	// XXX: Replace spaces, because wordwrap trims them out at the ends
	paragraph = strings.ReplaceAll(paragraph, " ", "·")

	f := wordwrap.NewWriter(lineLimit)
	f.Breakpoints = []rune{'·'}
	f.KeepNewlines = false
	f.Write([]byte(paragraph))
	f.Close()

	paragraph = strings.ReplaceAll(f.String(), "·", " ")

	return paragraph
}

func toKeysSlice(mp map[int]bool) []int {
	acc := []int{}
	for key := range mp {
		acc = append(acc, key)
	}
	return acc
}

func styleAllRunes(runes []rune, style StringStyle) string {
	var acc strings.Builder

	for idx, char := range runes {
		_ = idx
		acc.WriteString(style(string(char)).String())
	}

	return acc.String()
}
