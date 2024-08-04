package cmd

import "time"

func (m Test) calculateResults() Results {
	return Results{
		wpm:      10,
		accuracy: 10.0,
		time:     time.Second,
	}
}
