package solver

import (
	"../common"
	"../magnets"
	"testing"
)

func TestNew(t *testing.T) {
	game := magnets.New(5, 6)
	cbs := new(game)

	answer := len(cbs)
	if answer != 6 {
		t.Errorf("ERROR: Expected 6, got %d", answer)
	}
	answer = len(cbs[0])
	if answer != 5 {
		t.Errorf("ERROR: Expected 5, got %d", answer)
	}
	answer = len(cbs[0][0])
	if answer != 3 {
		t.Errorf("ERROR: Expected 2, got %d", answer)
	}
}

func TestSetFrame(t *testing.T) {
	testCases := []struct {
		r        rune
		expected rune
	}{
		{common.Positive, common.Negative},
		{common.Negative, common.Positive},
		{common.Neutral, common.Neutral},
	}

	game := magnets.New(2, 1)
	cbs := new(game)

	for _, testCase := range testCases {
		setFrame(game, cbs, 0, 0, testCase.r)
		answer := game.Guess.Get(0, 1)
		if answer != testCase.expected {
			t.Errorf("ERROR: Expected '%c' got '%c'", testCase.expected, answer)
		}
	}
}
