package magnets

import (
	"testing"
)

// TODO: write tests for ...
// setFrameMagnet()
// placeFrames()
// placePieces()
// makeGame()

func TestNew(t *testing.T) {
	testCases := []struct {
		width  int
		height int
	}{
		{1, 2},
		{2, 1},
		{2, 2},
		{5, 4},
	}

	for _, testCase := range testCases {
		answer := New(testCase.width, testCase.height)
		h := answer.grid.Height()
		if h != testCase.height {
			t.Errorf("ERROR: For %d, %d got %d", testCase.width, testCase.height, h)
		}
	}
}
