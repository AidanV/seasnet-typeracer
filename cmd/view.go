package cmd

import (
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

func (m model) View() string {

	var termWidth, termHeight = m.width, m.height

	var coloredStopwatch string
	if m.state.stopwatch.isRunning {
		coloredStopwatch = style(state.stopwatch.stopwatch.View(), m.styles.runningTimer)
	} else {
		coloredStopwatch = style(state.stopwatch.stopwatch.View(), m.styles.stoppedTimer)
	}

	paragraphView := state.base.paragraphView(lineLenLimit, m.styles)
	lines := strings.Split(paragraphView, "\n")
	cursorLine := findCursorLine(strings.Split(paragraphView, "\n"), state.base.cursor)

	linesAroundCursor := strings.Join(getLinesAroundCursor(lines, cursorLine), "\n")

	s += positionVerticaly(termHeight)
	avgLineLen := averageLineLen(lines)
	indentBy := uint(math.Max(0, float64(termWidth/2-avgLineLen/2)))

	s += m.indent(coloredStopwatch, indentBy) + "\n\n" + m.indent(linesAroundCursor, indentBy)

	if !state.stopwatch.isRunning {
		s += "\n\n\n"
		s += lipgloss.PlaceHorizontal(termWidth, lipgloss.Center, style("ctrl+r to restart, ctrl+q to menu", m.styles.toEnter))
	}
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
