package solver

import (
	"github.com/erikbryant/magnets/common"
	"github.com/erikbryant/magnets/magnets"
	"testing"
)

func TestNewNoWall(t *testing.T) {
	game, _ := magnets.Deserialize("5x4:02121,2112,20211,1212,LRLRTLRTTBLRBBTLRLRB")
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

	game, _ = magnets.Deserialize("5x5:03232,22222,12322,32221,LRLRTTLRTBBTTBTTBBTBBLRB*")
}

func TestNewWithWall(t *testing.T) {
	game, _ := magnets.Deserialize("5x7:32323,3222211,33223,2312221,TLRTTBLRBBTTTLRBBBTTTLRBBBLRLR*LRLR")
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

	game, _ := magnets.Deserialize("2x1:10,1,01,1,LR")
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
	game, _ := magnets.Deserialize("2x1:10,1,01,1,LR")
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
	game, _ := magnets.Deserialize("1x2:1,10,1,01,TB")
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
	game, _ := magnets.Deserialize("3x4:212,1202,122,2111,TTTBBBLRTLRB")

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

func TestRowNeedsPartiallySolved(t *testing.T) {
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
	game, _ := magnets.Deserialize("3x4:212,1202,122,2111,TTTBBBLRTLRB")
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
	game, _ := magnets.Deserialize("3x4:212,1202,122,2111,TTTBBBLRTLRB")

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

func TestColNeedsPartiallySolved(t *testing.T) {
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
	game, _ := magnets.Deserialize("3x4:212,1202,122,2111,TTTBBBLRTLRB")
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
	game, _ := magnets.Deserialize("3x4:212,1202,122,2111,TTTBBBLRTLRB")
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
		{3, 3, 3, 3},
	}
	game, _ := magnets.Deserialize("3x4:212,1202,122,2111,TTTBBBLRTLRB")
	cbs := new(game)

	// Solve a little bit of the board.
	cbs.setFrame(game, 3, 0, common.Positive)
	cbs.setFrame(game, 2, 0, common.Neutral)

	game.Print()

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

// This is becoming more of a regression test. If the run time gets too high maybe it should be moved out of the unit tests.
func TestSolve(t *testing.T) {
	testCases := []string{
		// Baby steps.
		"1x2:1,10,1,01,TB",
		"2x1:10,1,01,1,LR",
		"2x2:10,01,01,01,LRLR",

		// Generated and solved by this program.
		"2x2:00,00,00,00,LRLR",
		"3x2:000,00,000,00,LRTLRB",

		"3x3:101,200,011,110,LRTTTBBB*",
		"3x3:011,110,101,020,LRTLRBLR*",
		// "3x3:201,201,111,120,LRTTTBBB*",
		"3x3:111,021,111,111,LRTLRBLR*",
		// "3x3:112,112,121,121,*LRTLRBLR",
		"3x3:111,021,111,111,LRTTTBBB*",
		// "3x3:112,112,121,121,*LRTLRBLR",
		"3x4:212,1202,122,2111,TTTBBBLRTLRB",
		"3x4:212,2021,212,0212,TTTBBBTTTBBB",
		"3x4:122,2111,212,1202,TTTBBBLRTLRB",
		"3x4:221,1121,212,0212,TLRBLRTLRBLR",
		"4x4:1112,1211,1112,1220,TLRTBLRBTLRTBLRB",
		"4x4:1212,0222,1122,1122,TLRTBTTBTBBTBLRB",
		"4x4:0212,0212,0122,1121,TLRTBTTBTBBTBLRB",
		"4x4:1012,0121,0112,1021,TTTTBBBBLRTTLRBB",
		// "5x4:02121,2112,20211,1212,LRLRTLRTTBLRBBTLRLRB",
		// "5x5:23031,22212,22122,32121,TLRTTBTTBBTBBTTBTTBB*BBLR",
		"5x5:03232,22222,12322,32221,LRLRTTLRTBBTTBTTBBTBBLRB*",
		"5x5:31312,30322,13132,21232,LRLRTLRLRBTTLRTBBTTBLRBB*",
		// "5x5:21222,31320,12132,22221,LRLRTLRTTBTTBBTBBTTBLRBB*",
		// "5x7:32323,3222211,33223,2312221,TLRTTBLRBBTTTLRBBBTTTLRBBBLRLR*LRLR",
		"6x6:331232,222233,330332,232223,LRLRLRLRTTLRTTBBTTBBTTBBLRBBLRLRLRLR",
		"6x6:321223,113233,231313,031333,LRTLRTTTBTTBBBTBBTTTBTTBBBTBBTLRBLRB",
		"6x6:232323,033333,322233,123333,LRTLRTLRBLRBTLRTTTBLRBBBTLRLRTBLRLRB",
		"6x6:223033,323221,132133,322222,LRLRLRLRLRLRLRTTLRTTBBTTBBTTBBLRBBLR",
		"7x7:3243423,4343412,3333333,3434340,TLRTLRTBTTBTTBTBBTBBTBTTBTTBTBBTBBTBTTBTTB*BB*BB*",
		"7x7:2434223,3424322,4242314,4243403,TTTTLRTBBBBTTBLRTTBBTTTBBLRBBBTLRLRLRBLRTTLRLR*BB",
		// "7x7:4222414,4233403,3233323,2424340,TTTLRTTBBBTTBBTLRBBLRBLRTTTTLRTBBBBTTBTTTTBB*BBBB",
		"7x7:3341223,4332312,2423043,3323313,TTTTLRTBBBBTTBTTTTBBTBBBBTTBLRLRBBTTTLRLRBBBLRLR*",
		"6x10:535454,3333323033,544445,3333331133,TTLRTTBBTTBBTTBBTTBBLRBBTLRTLRBTTBTTTBBTBBBLRBLRTTLRLRBBLRLR",
		"6x10:244535,3230333321,425345,3221333321,LRLRTTLRTTBBLRBBTTLRLRBBTTTLRTBBBTTBTTTBBTBBBTTBLRTBBTLRBLRB",
		"6x10:533455,3223333321,443545,3322333330,LRTLRTTTBTTBBBTBBTTTBTTBBBTBBTTTBTTBBBTBBTTTBTTBBBTBBTLRBLRB",
		"6x10:434554,3212333323,443545,3203333323,TLRTTTBTTBBBTBBTLRBTTBTTTBBTBBBTTBTTTBBTBBBLRBTTTLRTBBBLRBLR",
		// "9x9:334345444,353545441,433444444,534454540,TTTTTTTLRBBBBBBBTTLRLRLRTBBTLRTTTBTTBTTBBBTBBTBBLRTBTTBTTTTBTBBTBBBBTBLRBLRLRBLR*",
		// "2x50:jh,01111110101001111000011110111110011111011101111111,me,01111101101001110100101110111101011111011011111111,LRTTBBLRLRLRTTBBLRLRLRTTBBTTBBLRTTBBLRLRTTBBLRTTBBLRLRLRTTBBTTBBLRLRTTBBTTBBLRTTBBTTBBTTBBTTBBLRTTBB",

		// Came from iPhone game, so guaranteed to be solvable. Tricky, but solvable.
		// "7x9:1444133,330323222,2244242,232333211,LRLRLRTTLRLRTBBLRT*BTTTTBLRBBBBTLRTLRTBTTBTTBTBBTBBTBLRBLRBLRLR",
		// "7x9:4442324,412332242,3343343,333312233,LRLRTLRTLR*BLRBTTTLRTTBBBTTBBTTTBBTTBBBTTBBTLRBBTTBLRLRBBLRLRLR",
		// "7x9:3342325,331221433,2442334,331222342,TLRLRTTBLRLRBBTTLR*TTBBTTTBBTTBBBTTBBTLRBBTTBTTTTBBTBBBBLRBLRLR",
		// "15x9:454544143332433,766636477,545453133524242,775455667,LRTTLRTTTLRTTLRLRBBTTBBB*TBBTTTTLRBBLRLRBLRBBBBTLRLRTTTTLRTTTTBLRLRBBBBTTBBBBLRLRTLRLRBBTTLRTTTTBLRTTLRBBLRBBBBLRTBBTTLRLRLRLRLRBLRBBLR",
		// "15x9:333445434434343,786555657,442354433434434,768646575,TTLRLRLR*LRTLRTBBTLRLRTLRTBLRBLRBLRLRBTTBTLRTTLRTTTLRBBTBLRBBTTBBBTTLRBLRLRTBBTLRBBLRTTTTTBLRBLRLRLRBBBBBLRTTTLRTTTLRLRTLRBBBLRBBBLRLRB",
		// "15x9:423444254334424,686756455,414443345334334,758766544,TLRLRLRLRLRLRLRB*TTTTLRTLRTTTTLRBBBBLRBTTBBBBLRTLRLRLRBBLRLRTTBLRLRLRLRTLRTBBTLRLRTLRTBLRBTTBTLRTBLRBLRLRBBTBLRBLRTTTLRTLRBLRLRLRBBBLRB",
		// "12x10:555213545454,5543546565,555215354544,5434465656,TTTLRLRTTLRTBBBTLRTBBTTBLRTBTTBLRBBTLRBTBBLRTLRBLRTBLRLRBLRTLRBLRTTLRLRBTLRTTBBTTLRTBLRBBLRBBLRBTTTTLRTLRTLRBBBBLRBLRBLR",
		// "12x10:545232142444,5462542336,444423124453,5652451345,TTTTLRLRLRLRBBBBTLRTLRTTLRLRBTTBLRBBTLRLRBBTLRLRBLRTTLRBLRTTTLRBBTLRLRBBBLRTTBTLRTTTLRTBBTBLRBBBTTBLRBTLRLRTBBLRLRBLRLRB",
	}

	for _, testCase := range testCases {
		game, ok := magnets.Deserialize(testCase)
		if !ok {
			t.Errorf("ERROR: Unable to deserialize %s", testCase)
		}
		Solve(game)
		if !game.Solved() {
			t.Errorf("ERROR: Unable to solve %s", testCase)
		}
	}
}
