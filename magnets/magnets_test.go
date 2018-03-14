package magnets

import (
	"../board"
	"../common"
	"testing"
)

// Print() is trivial and does not need a test.
// Solved() is trivial and does not need a test.
// CountRow() is trivial and does not need a test.
// CountCol() is trivial and does not need a test.

// TODO: write tests for ...
// Valid()
// setFrame()
// SetGuess()
// setFrameMagnet()
// placeFrames()
// placePieces()
// Deserialize()

func TestNew(t *testing.T) {
	testCases := []struct {
		width  int
		height int
	}{
		{1, 2},
		{2, 1},
		{2, 2},
		{25, 25},
	}

	for _, testCase := range testCases {
		answer := New(testCase.width, testCase.height)
		h := answer.grid.Height()
		if h != testCase.height {
			t.Errorf("ERROR: For %d, %d got %d", testCase.width, testCase.height, h)
		}
	}
}

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
	}

	game.frames.Set(0, 0, common.Left)
	game.frames.Set(0, 1, common.Right)
	game.frames.Set(1, 0, common.Wall)
	game.frames.Set(1, 1, common.Wall)
	game.frames.Set(2, 0, common.Up)
	game.frames.Set(2, 1, common.Up)
	game.frames.Set(3, 0, common.Down)
	game.frames.Set(3, 1, common.Down)

	game.grid.Set(0, 0, common.Positive)
	game.grid.Set(0, 1, common.Negative)
	game.grid.Set(1, 0, common.Wall)
	game.grid.Set(1, 1, common.Wall)
	game.grid.Set(2, 0, common.Positive)
	game.grid.Set(2, 1, common.Negative)
	game.grid.Set(3, 0, common.Negative)
	game.grid.Set(3, 1, common.Positive)

	expected := "2x4:21,1011,12,1011,LR**TTBB"

	answer, ok := game.Serialize()
	if !ok {
		t.Errorf("ERROR: game.Serialize() failed")
	}
	if answer != expected {
		t.Errorf("ERROR: Expected: %s Got: %s", expected, answer)
	}
}
