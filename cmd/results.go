package cmd

import (
	"math"
	"time"
)

func (m Test) calculateResults() Results {
	return Results{
		wpm:      10,
		accuracy: 10.0,
		time:     time.Second,
	}
}

func (t Test) calculateNormalizedWpm(elapsedMinutes float64) float64 {
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
