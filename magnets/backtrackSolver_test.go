package magnets

import (
	"testing"
)

func TestSetCell(t *testing.T) {
}

func TestCountSolutions(t *testing.T) {
	testCases := []struct {
		game     string
		expected int
	}{
		{"2x2:11,00,11,00,TTBB", 0},
		{"1x2:1,10,1,01,TB", 1},
		{"2x2:11,11,11,11,TTBB", 2},
		{"3x2:111,21,111,12,TTTBBB", 1},
		{"5x2:11011,22,11011,22,LRTLRLRBLR", 4},
	}

	for _, testCase := range testCases {
		game, ok := Deserialize(testCase.game)
		if !ok {
			t.Errorf("ERROR: failed to deserialize %s", testCase.game)
		}

		answer := game.CountSolutions(0, 0)
		if answer != testCase.expected {
			t.Errorf("ERROR: for %s expected %d got %d", testCase.game, testCase.expected, answer)
		}
	}
}

func TestSingleSolution(t *testing.T) {
	testCases := []struct {
		game     string
		expected bool
	}{
		{"2x2:11,00,11,00,TTBB", false}, // 0 solutions
		{"1x2:1,10,1,01,TB", true},      // 1 solution
		{"2x2:11,11,11,11,TTBB", false}, // 2 solutions
	}

	for _, testCase := range testCases {
		game, ok := Deserialize(testCase.game)
		if !ok {
			t.Errorf("ERROR: failed to deserialize %s", testCase.game)
		}

		answer := game.singleSolution()
		if answer != testCase.expected {
			t.Errorf("ERROR: for %s expected %t got %t", testCase.game, testCase.expected, answer)
		}
	}
}
