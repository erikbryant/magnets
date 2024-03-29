package magnets

import (
	"testing"

	"github.com/erikbryant/magnets/board"
	"github.com/erikbryant/magnets/common"
)

func TestCountToRune(t *testing.T) {
	testCases := []struct {
		n        int
		expected rune
	}{
		{-1, '-'},
		{0, '0'},
		{1, '1'},
		{9, '9'},
		{10, 'a'},
		{11, 'b'},
		{35, 'z'},
		{36, 'A'},
		{37, 'B'},
		{61, 'Z'},
		{62, '!'},
		{63, '!'},
	}

	for _, testCase := range testCases {
		answer := countToRune(testCase.n)
		if answer != testCase.expected {
			t.Errorf("ERROR: For %d expected '%c' got '%c'", testCase.n, testCase.expected, answer)
		}
	}
}

func TestRuneToCount(t *testing.T) {
	testCases := []struct {
		expected int
		n        rune
	}{
		{-1, '-'},
		{0, '0'},
		{1, '1'},
		{9, '9'},
		{10, 'a'},
		{11, 'b'},
		{35, 'z'},
		{36, 'A'},
		{37, 'B'},
		{61, 'Z'},
		{-1, '!'},
	}

	for _, testCase := range testCases {
		answer := runeToCount(testCase.n)
		if answer != testCase.expected {
			t.Errorf("ERROR: For '%c' expected %d got %d", testCase.n, testCase.expected, answer)
		}
	}
}

func TestSerialize(t *testing.T) {
	game := Game{
		grid:   board.New(2, 4),
		frames: board.New(2, 4),
		Guess:  board.New(2, 4),
		colPos: []int{2, 1},
		rowPos: []int{1, 0, 1, 1},
		colNeg: []int{1, 2},
		rowNeg: []int{1, 0, 1, 1},
	}

	game.frames.Set(0, 0, common.Left, false)
	game.frames.Set(0, 1, common.Right, false)
	game.frames.Set(1, 0, common.Wall, false)
	game.frames.Set(1, 1, common.Wall, false)
	game.frames.Set(2, 0, common.Up, false)
	game.frames.Set(2, 1, common.Up, false)
	game.frames.Set(3, 0, common.Down, false)
	game.frames.Set(3, 1, common.Down, false)

	game.grid.Set(0, 0, common.Positive, false)
	game.grid.Set(0, 1, common.Negative, false)
	game.grid.Set(1, 0, common.Wall, false)
	game.grid.Set(1, 1, common.Wall, false)
	game.grid.Set(2, 0, common.Positive, false)
	game.grid.Set(2, 1, common.Negative, false)
	game.grid.Set(3, 0, common.Negative, false)
	game.grid.Set(3, 1, common.Positive, false)

	expected := "2x4:21,1011,12,1011,LR**TTBB"

	answer, ok := game.Serialize()
	if !ok {
		t.Errorf("ERROR: game.Serialize() failed")
	}
	if answer != expected {
		t.Errorf("ERROR: Expected: %s Got: %s", expected, answer)
	}
}

func TestSerialize_2(t *testing.T) {
	testCases := []string{
		"2x6:13,101011,22,100111,LRLRTTBBTTBB",
		"3x5:132,12111,222,21210,LRTTTBBBTTTBBB*",
	}

	for _, testCase := range testCases {
		game, ok := Deserialize(testCase)
		if !ok {
			t.Errorf("ERROR: Unable to deserialize %s", testCase)
		}
		answer, ok := game.Serialize()
		if !ok {
			t.Errorf("ERROR: Unable to serialize %s", testCase)
		}
		if answer != testCase {
			t.Errorf("ERROR: serialize failed to match original. Expected %s got %s", testCase, answer)
		}
	}
}

func TestDeserialize(t *testing.T) {
	testCases := []struct {
		s     string
		valid bool
	}{
		// Empty string
		{"", false},
		// Missing the commas
		{"1x2:110101TB", false},
		// List of negatives is short
		{"5x2:11011,22,1101,22,LRTLRLRBLR", false},

		// Valid
		{"4x5:2022,12021,2013,21201,TTTTBBBBTLRTBLRBLRLR", true},
	}

	for _, testCase := range testCases {
		_, valid := Deserialize(testCase.s)
		if valid != testCase.valid {
			t.Errorf("For %s expected %t, got %t", testCase.s, testCase.valid, valid)
		}
	}
}
