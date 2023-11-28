package solver

import (
	"testing"

	"github.com/erikbryant/magnets/common"
	"github.com/erikbryant/magnets/magnets"
)

func TestNewNoWall(t *testing.T) {
	game, ok := magnets.Deserialize("5x4:02121,2112,20211,1212,LRLRTLRTTBLRBBTLRLRB")
	if !ok {
		t.Errorf("Unable to deserialize board")
	}

	cbs := new(game)

	answer := len(cbs)
	if answer != 4 {
		t.Errorf("ERROR: Expected 4, got %d", answer)
	}
	answer = len(cbs[0])
	if answer != 5 {
		t.Errorf("ERROR: Expected 5, got %d", answer)
	}
	answer = len(cbs[0][0])
	if answer != 3 {
		t.Errorf("ERROR: Expected 3, got %d", answer)
	}

	_, ok = magnets.Deserialize("5x5:03232,22222,12322,32221,LRLRTTLRTBBTTBTTBBTBBLRB*")
	if !ok {
		t.Errorf("Unable to deserialize board")
	}
}

func TestNewWithWall(t *testing.T) {
	game, ok := magnets.Deserialize("5x7:32323,3222211,33223,2312221,TLRTTBLRBBTTTLRBBBTTTLRBBBLRLR*LRLR")
	if !ok {
		t.Errorf("Unable to deserialize board")
	}

	cbs := new(game)

	answer := len(cbs)
	if answer != 7 {
		t.Errorf("ERROR: Expected 7, got %d", answer)
	}
	answer = len(cbs[0])
	if answer != 5 {
		t.Errorf("ERROR: Expected 5, got %d", answer)
	}
	answer = len(cbs[0][0])
	if answer != 3 {
		t.Errorf("ERROR: Expected 3, got %d", answer)
	}
	r := game.Guess.Get(6, 0, false)
	if r != common.Wall {
		t.Errorf("ERROR: Expected Wall, got %c", r)
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

	game, ok := magnets.Deserialize("2x1:10,1,01,1,LR")
	if !ok {
		t.Errorf("Unable to deserialize board")
	}

	cbs := new(game)

	for _, testCase := range testCases {
		dirty = false
		cbs.setFrame(game, 0, 0, testCase.r)
		answer := game.Guess.Get(0, 1, false)
		if answer != testCase.expected {
			t.Errorf("ERROR: Expected '%c' got '%c'", testCase.expected, answer)
		}
		if !dirty {
			t.Errorf("ERROR: Expected dirty = true, got dirty = %v", dirty)
		}
	}
}

func TestUnsetPossibility(t *testing.T) {
	game, ok := magnets.Deserialize("2x1:10,1,01,1,LR")
	if !ok {
		t.Errorf("Unable to deserialize board")
	}

	cbs := new(game)

	answer := len(cbs[0][0])
	if answer != 3 {
		t.Errorf("ERROR: Expected 3, got %d", answer)
	}

	dirty = false
	cbs.unsetPossibility(game, 0, 0, common.Positive)
	answer = len(cbs[0][0])
	if answer != 2 {
		t.Errorf("ERROR: Expected 2, got %d", answer)
	}
	if !dirty {
		t.Errorf("ERROR: Expected dirty = true, got dirty = %v", dirty)
	}

	dirty = false
	cbs.unsetPossibility(game, 0, 0, common.Positive)
	answer = len(cbs[0][0])
	if answer != 2 {
		t.Errorf("ERROR: Expected still to be 2, got %d", answer)
	}
	if dirty {
		t.Errorf("ERROR: Expected dirty = false, got dirty = %v", dirty)
	}

	dirty = false
	cbs.unsetPossibility(game, 0, 0, common.Negative)
	answer = len(cbs[0][0])
	if answer != 1 {
		t.Errorf("ERROR: Expected 1, got %d", answer)
	}
	if !dirty {
		t.Errorf("ERROR: Expected dirty = true, got dirty = %v", dirty)
	}
}

func TestRowNeeds(t *testing.T) {
	testCases := []struct {
		row             int
		expectedPos     int
		expectedNeg     int
		expectedNeutral int
	}{
		{0, 1, 2, 0},
		{1, 2, 1, 0},
		{2, 0, 1, 2},
		{3, 2, 1, 0},
	}

	game, ok := magnets.Deserialize("3x4:212,1202,122,2111,TTTBBBLRTLRB")
	if !ok {
		t.Errorf("Unable to deserialize board")
	}

	// Positive
	for _, testCase := range testCases {
		answer := rowNeeds(game, testCase.row, common.Positive)
		if answer != testCase.expectedPos {
			t.Errorf("ERROR: For row %d expected rowNeeds %d, got %d", testCase.row, testCase.expectedPos, answer)
		}
	}

	// Negative
	for _, testCase := range testCases {
		answer := rowNeeds(game, testCase.row, common.Negative)
		if answer != testCase.expectedNeg {
			t.Errorf("ERROR: For row %d expected rowNeeds %d, got %d", testCase.row, testCase.expectedNeg, answer)
		}
	}

	// Neutral
	for _, testCase := range testCases {
		answer := rowNeeds(game, testCase.row, common.Neutral)
		if answer != testCase.expectedNeutral {
			t.Errorf("ERROR: For row %d expected rowNeeds %d, got %d", testCase.row, testCase.expectedNeutral, answer)
		}
	}
}

func TestRowNeeds_PartiallySolved(t *testing.T) {
	testCases := []struct {
		row             int
		expectedPos     int
		expectedNeg     int
		expectedNeutral int
	}{
		{0, 1, 2, 0},
		{1, 2, 1, 0},
		{2, 0, 1, 0},
		{3, 1, 0, 0},
	}

	game, ok := magnets.Deserialize("3x4:212,1202,122,2111,TTTBBBLRTLRB")
	if !ok {
		t.Errorf("Unable to deserialize board")
	}

	cbs := new(game)

	// Solve a little bit of the board.
	cbs.setFrame(game, 3, 0, common.Positive)
	cbs.setFrame(game, 2, 0, common.Neutral)

	// Positive
	for _, testCase := range testCases {
		answer := rowNeeds(game, testCase.row, common.Positive)
		if answer != testCase.expectedPos {
			t.Errorf("ERROR: For row %d expected rowNeeds %d, got %d", testCase.row, testCase.expectedPos, answer)
		}
	}

	// Negative
	for _, testCase := range testCases {
		answer := rowNeeds(game, testCase.row, common.Negative)
		if answer != testCase.expectedNeg {
			t.Errorf("ERROR: For row %d expected rowNeeds %d, got %d", testCase.row, testCase.expectedNeg, answer)
		}
	}

	// Neutral
	for _, testCase := range testCases {
		answer := rowNeeds(game, testCase.row, common.Neutral)
		if answer != testCase.expectedNeutral {
			t.Errorf("ERROR: For row %d expected rowNeeds %d, got %d", testCase.row, testCase.expectedNeutral, answer)
		}
	}
}

func TestColNeeds(t *testing.T) {
	testCases := []struct {
		col             int
		expectedPos     int
		expectedNeg     int
		expectedNeutral int
	}{
		{0, 2, 1, 1},
		{1, 1, 2, 1},
		{2, 2, 2, 0},
	}

	game, ok := magnets.Deserialize("3x4:212,1202,122,2111,TTTBBBLRTLRB")
	if !ok {
		t.Errorf("Unable to deserialize board")
	}

	// Positive
	for _, testCase := range testCases {
		answer := colNeeds(game, testCase.col, common.Positive)
		if answer != testCase.expectedPos {
			t.Errorf("ERROR: For col %d expected colNeeds %d, got %d", testCase.col, testCase.expectedPos, answer)
		}
	}

	// Negative
	for _, testCase := range testCases {
		answer := colNeeds(game, testCase.col, common.Negative)
		if answer != testCase.expectedNeg {
			t.Errorf("ERROR: For col %d expected colNeeds %d, got %d", testCase.col, testCase.expectedNeg, answer)
		}
	}

	// Neutral
	for _, testCase := range testCases {
		answer := colNeeds(game, testCase.col, common.Neutral)
		if answer != testCase.expectedNeutral {
			t.Errorf("ERROR: For col %d expected colNeeds %d, got %d", testCase.col, testCase.expectedNeutral, answer)
		}
	}
}

func TestColNeeds_PartiallySolved(t *testing.T) {
	testCases := []struct {
		col             int
		expectedPos     int
		expectedNeg     int
		expectedNeutral int
	}{
		{0, 1, 1, 0},
		{1, 1, 1, 0},
		{2, 2, 2, 0},
	}

	game, ok := magnets.Deserialize("3x4:212,1202,122,2111,TTTBBBLRTLRB")
	if !ok {
		t.Errorf("Unable to deserialize board")
	}

	cbs := new(game)

	// Solve a little bit of the board.
	cbs.setFrame(game, 3, 0, common.Positive)
	cbs.setFrame(game, 2, 0, common.Neutral)

	// Positive
	for _, testCase := range testCases {
		answer := colNeeds(game, testCase.col, common.Positive)
		if answer != testCase.expectedPos {
			t.Errorf("ERROR: For col %d expected colNeeds %d, got %d", testCase.col, testCase.expectedPos, answer)
		}
	}

	// Negative
	for _, testCase := range testCases {
		answer := colNeeds(game, testCase.col, common.Negative)
		if answer != testCase.expectedNeg {
			t.Errorf("ERROR: For col %d expected colNeeds %d, got %d", testCase.col, testCase.expectedNeg, answer)
		}
	}

	// Neutral
	for _, testCase := range testCases {
		answer := colNeeds(game, testCase.col, common.Neutral)
		if answer != testCase.expectedNeutral {
			t.Errorf("ERROR: For col %d expected colNeeds %d, got %d", testCase.col, testCase.expectedNeutral, answer)
		}
	}
}

func TestRowHasSpaceForTotal(t *testing.T) {
	testCases := []struct {
		row             int
		expectedPos     int
		expectedNeg     int
		expectedNeutral int
	}{
		{0, 3, 3, 3},
		{1, 3, 3, 3},
		{2, 2, 2, 3},
		{3, 2, 2, 3},
	}

	game, ok := magnets.Deserialize("3x4:212,1202,122,2111,TTTBBBLRTLRB")
	if !ok {
		t.Errorf("Unable to deserialize board")
	}

	cbs := new(game)

	// Positive
	for _, testCase := range testCases {
		answer := cbs.rowHasSpaceForTotal(game, testCase.row, common.Positive)
		if answer != testCase.expectedPos {
			t.Errorf("ERROR: For row %d expected rowHasSpaceForTotal positive %d, got %d", testCase.row, testCase.expectedPos, answer)
		}
	}

	// Negative
	for _, testCase := range testCases {
		answer := cbs.rowHasSpaceForTotal(game, testCase.row, common.Negative)
		if answer != testCase.expectedNeg {
			t.Errorf("ERROR: For row %d expected rowHasSpaceForTotal negative %d, got %d", testCase.row, testCase.expectedNeg, answer)
		}
	}

	// Neutral
	for _, testCase := range testCases {
		answer := cbs.rowHasSpaceForTotal(game, testCase.row, common.Neutral)
		if answer != testCase.expectedNeutral {
			t.Errorf("ERROR: For row %d expected rowHasSpaceForTotal neutral %d, got %d", testCase.row, testCase.expectedNeutral, answer)
		}
	}
}

func TestRowHasSpaceForTotalPartiallySolved(t *testing.T) {
	testCases := []struct {
		row             int
		expectedPos     int
		expectedNeg     int
		expectedNeutral int
	}{
		{0, 3, 3, 3},
		{1, 3, 3, 3},
		{2, 1, 1, 3},
		{3, 2, 2, 1},
	}

	game, ok := magnets.Deserialize("3x4:212,1202,122,2111,TTTBBBLRTLRB")
	if !ok {
		t.Errorf("Unable to deserialize board")
	}

	cbs := new(game)

	// Solve a little bit of the board.
	cbs.setFrame(game, 3, 0, common.Positive)
	cbs.setFrame(game, 2, 0, common.Neutral)

	// Positive
	for _, testCase := range testCases {
		answer := cbs.rowHasSpaceForTotal(game, testCase.row, common.Positive)
		if answer != testCase.expectedPos {
			t.Errorf("ERROR: For row %d expected rowHasSpaceForTotal positive %d, got %d", testCase.row, testCase.expectedPos, answer)
		}
	}

	// Negative
	for _, testCase := range testCases {
		answer := cbs.rowHasSpaceForTotal(game, testCase.row, common.Negative)
		if answer != testCase.expectedNeg {
			t.Errorf("ERROR: For row %d expected rowHasSpaceForTotal negative %d, got %d", testCase.row, testCase.expectedNeg, answer)
		}
	}

	// Neutral
	for _, testCase := range testCases {
		answer := cbs.rowHasSpaceForTotal(game, testCase.row, common.Neutral)
		if answer != testCase.expectedNeutral {
			t.Errorf("ERROR: For row %d expected rowHasSpaceForTotal neutral %d, got %d", testCase.row, testCase.expectedNeutral, answer)
		}
	}
}

// func TestColHasSpaceForTotal(t *testing.T) {
// 	t.Errorf("Not implemented")
// }

// func TestRowHasSpaceForRemaining(t *testing.T) {
// 	t.Errorf("Not implemented")
// }

// func TestColHasSpaceForRemaining(t *testing.T) {
// 	t.Errorf("Not implemented")
// }

func TestValidate(t *testing.T) {
	// Test #1 - Invalid CBS state
	game, ok := magnets.Deserialize("3x4:212,1202,122,2111,TTTBBBLRTLRB")
	if !ok {
		t.Errorf("Unable to deserialize board")
	}

	cbs := new(game)

	err := cbs.validate(game)
	if err != nil {
		t.Error("Error validating CBS", err)
	}

	game.Guess.Set(1, 1, common.Positive, false)
	err = cbs.validate(game)
	if err == nil {
		t.Error("validate was supposed to find an error but did not")
	}

	cbs[1][1] = map[rune]bool{'%': true}
	err = cbs.validate(game)
	if err == nil {
		t.Error("validate was supposed to find an error but did not")
	}

	// Test #2 - Invalid game state
	game, ok = magnets.Deserialize("3x4:212,1202,122,2111,TTTBBBLRTLRB")
	if !ok {
		t.Errorf("Unable to deserialize board")
	}

	cbs = new(game)

	err = cbs.validate(game)
	if err != nil {
		t.Error("Error validating CBS", err)
	}

	game.Guess.Set(0, 0, common.Negative, false)
	game.Guess.Set(0, 1, common.Negative, false)

	err = cbs.validate(game)
	if err == nil {
		t.Error("validate was supposed to find an error but did not")
	}
}
