package cmd

import (
	"math"
	"time"
)

func (t Test) calculateNormalizedWpm() float64 {
	elapsedMinutes := time.Since(t.startTime).Minutes()
	return t.calculateWpm(len(t.inputBuffer)/5, elapsedMinutes)
}

func (t Test) calculateWpm(wordCnt int, elapsedMinutes float64) float64 {
	if elapsedMinutes == 0 {
		return 0
	} else {
		grossWpm := float64(wordCnt) / elapsedMinutes
		netWpm := grossWpm - float64(len(t.mistakes.mistakesAt))/elapsedMinutes
		return math.Max(0, netWpm)
	}
}

func (t Test) calculateNumCorrect() int {
	totalCorrect := 0
	for i, r := range t.wordsToEnter {
		if i >= len(t.inputBuffer) {
			return totalCorrect
		}
		if r == t.inputBuffer[i] {
			totalCorrect++
		}
	}
	return totalCorrect
}
