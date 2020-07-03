package solver

import (
	"bufio"
	"github.com/erikbryant/magnets/common"
	"github.com/erikbryant/magnets/magnets"
	"os"
	"strings"
	"testing"
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

	game, ok = magnets.Deserialize("5x5:03232,22222,12322,32221,LRLRTTLRTBBTTBTTBBTBBLRB*")
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
	r := game.Guess.Get(6, 0)
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
		answer := game.Guess.Get(0, 1)
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
	cbs.unsetPossibility(0, 0, common.Positive)
	answer = len(cbs[0][0])
	if answer != 2 {
		t.Errorf("ERROR: Expected 2, got %d", answer)
	}
	if !dirty {
		t.Errorf("ERROR: Expected dirty = true, got dirty = %v", dirty)
	}

	dirty = false
	cbs.unsetPossibility(0, 0, common.Positive)
	answer = len(cbs[0][0])
	if answer != 2 {
		t.Errorf("ERROR: Expected still to be 2, got %d", answer)
	}
	if dirty {
		t.Errorf("ERROR: Expected dirty = false, got dirty = %v", dirty)
	}

	dirty = false
	cbs.unsetPossibility(0, 0, common.Negative)
	answer = len(cbs[0][0])
	if answer != 1 {
		t.Errorf("ERROR: Expected 1, got %d", answer)
	}
	if !dirty {
		t.Errorf("ERROR: Expected dirty = true, got dirty = %v", dirty)
	}

	dirty = false
	cbs.unsetPossibility(0, 0, common.Neutral)
	answer = len(cbs[0][0])
	if answer != 0 {
		t.Errorf("ERROR: Expected 0, got %d", answer)
	}
	if !dirty {
		t.Errorf("ERROR: Expected dirty = true, got dirty = %v", dirty)
	}
}

func TestJustOne(t *testing.T) {
	game, ok := magnets.Deserialize("1x2:1,10,1,01,TB")
	if !ok {
		t.Errorf("Unable to deserialize board")
	}

	cbs := new(game)

	dirty = false
	cbs[0][0] = map[rune]bool{common.Positive: true}
	cbs.justOne(game)
	answer := game.Guess.Get(0, 0)
	if answer != common.Positive {
		t.Errorf("ERROR: Expected %c, got %c", common.Positive, answer)
	}
	if !dirty {
		t.Errorf("ERROR: Expected dirty = true, got dirty = %v", dirty)
	}

	// This is the other end of that frame, so it should already be set.
	answer = game.Guess.Get(1, 0)
	if answer != common.Negative {
		t.Errorf("ERROR: Expected %c, got %c", common.Negative, answer)
	}

	// And, setting it again should not change it.
	dirty = false
	cbs[1][0] = map[rune]bool{common.Negative: true}
	cbs.justOne(game)
	answer = game.Guess.Get(1, 0)
	if answer != common.Negative {
		t.Errorf("ERROR: Expected %c, got %c", common.Negative, answer)
	}
	if dirty {
		t.Errorf("ERROR: Expected dirty = false, got dirty = %v", dirty)
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
		{2, 3, 3, 3},
		{3, 3, 3, 3},
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
			t.Errorf("ERROR: For row %d expected rowHasSpaceForTotal %d, got %d", testCase.row, testCase.expectedPos, answer)
		}
	}

	// Negative
	for _, testCase := range testCases {
		answer := cbs.rowHasSpaceForTotal(game, testCase.row, common.Negative)
		if answer != testCase.expectedNeg {
			t.Errorf("ERROR: For row %d expected rowHasSpaceForTotal %d, got %d", testCase.row, testCase.expectedNeg, answer)
		}
	}

	// Neutral
	for _, testCase := range testCases {
		answer := cbs.rowHasSpaceForTotal(game, testCase.row, common.Neutral)
		if answer != testCase.expectedNeutral {
			t.Errorf("ERROR: For row %d expected rowHasSpaceForTotal %d, got %d", testCase.row, testCase.expectedNeutral, answer)
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

// func TestSatisfied(t *testing.T) {
// 	t.Errorf("Not implemented")
// }

// func TestNeedAll(t *testing.T) {
// 	t.Errorf("Not implemented")
// }

// func TestDoubleSingle(t *testing.T) {
// 	t.Errorf("Not implemented")
// }

// func TestResolveNeighbors(t *testing.T) {
// 	t.Errorf("Not implemented")
// }

func TestZeroInRow(t *testing.T) {
	game, ok := magnets.Deserialize("2x2:00,00,00,00,LRLR")
	if !ok {
		t.Errorf("Unable to deserialize board")
	}

	cbs := new(game)
	cbs.zeroInRow(game)

	for col := 0; col < game.Guess.Width(); col++ {
		for row := 0; row < game.Guess.Height(); row++ {
			if cbs[row][col][common.Negative] {
				t.Errorf("Unexpected negative at %dx%d", row, col)
			}
			if cbs[row][col][common.Positive] {
				t.Errorf("Unexpected positive at %dx%d", row, col)
			}
		}
	}
}

func TestZeroInCol(t *testing.T) {
	game, ok := magnets.Deserialize("2x2:00,00,00,00,LRLR")
	if !ok {
		t.Errorf("Unable to deserialize board")
	}

	cbs := new(game)
	cbs.zeroInCol(game)

	for col := 0; col < game.Guess.Width(); col++ {
		for row := 0; row < game.Guess.Height(); row++ {
			if cbs[row][col][common.Negative] {
				t.Errorf("Unexpected negative at %dx%d", row, col)
			}
			if cbs[row][col][common.Positive] {
				t.Errorf("Unexpected positive at %dx%d", row, col)
			}
		}
	}
}

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

	game.Guess.Set(1, 1, common.Positive)
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

	game.Guess.Set(0, 0, common.Negative)
	game.Guess.Set(0, 1, common.Negative)

	err = cbs.validate(game)
	if err == nil {
		t.Error("validate was supposed to find an error but did not")
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
			t.Errorf("ERROR: For row %d expected rowHasSpaceForTotal %d, got %d", testCase.row, testCase.expectedPos, answer)
		}
	}

	// Negative
	for _, testCase := range testCases {
		answer := cbs.rowHasSpaceForTotal(game, testCase.row, common.Negative)
		if answer != testCase.expectedNeg {
			t.Errorf("ERROR: For row %d expected rowHasSpaceForTotal %d, got %d", testCase.row, testCase.expectedNeg, answer)
		}
	}

	// Neutral
	for _, testCase := range testCases {
		answer := cbs.rowHasSpaceForTotal(game, testCase.row, common.Neutral)
		if answer != testCase.expectedNeutral {
			t.Errorf("ERROR: For row %d expected rowHasSpaceForTotal %d, got %d", testCase.row, testCase.expectedNeutral, answer)
		}
	}
}

// helper runs solver tests against a given file.
func helper(t *testing.T, file string, expected bool) {
	f, err := os.Open(file)
	if err != nil {
		t.Errorf("Unable to open testcases %s %s", file, err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		testCase := scanner.Text()

		testCase = strings.TrimSpace(testCase)

		if len(testCase) == 0 {
			continue
		}

		if strings.HasPrefix(testCase, "//") {
			continue
		}

		game, ok := magnets.Deserialize(testCase)
		if !ok {
			t.Errorf("ERROR: Unable to deserialize %s", testCase)
		}
		Solve(game)
		if game.Solved() != expected {
			t.Errorf("ERROR: For %s expected solved to be %t", testCase, expected)
		}
	}
}

// This is becoming a regression test. If the run time gets too high, move out of the unit tests.
func TestSolve(t *testing.T) {
	helper(t, "test_solve.txt", true)
	helper(t, "test_solve_fail.txt", false)
}
