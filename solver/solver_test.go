package solver

import (
	"github.com/erikbryant/magnets/common"
	"github.com/erikbryant/magnets/magnets"
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

func TestValidate(t *testing.T) {
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

// This is becoming a regression test. If the run time gets too high, move out of the unit tests.
func TestSolve(t *testing.T) {
	testCases := []string{
		// Baby steps.
		"1x2:1,10,1,01,TB",
		"2x1:10,1,01,1,LR",
		"2x2:10,01,01,01,LRLR",

		// Generated and solved by this program. All are known solvable.
		"2x2:00,00,00,00,LRLR",

		"3x2:000,00,000,00,LRTLRB",

		"3x3:101,200,011,110,LRTTTBBB*",
		"3x3:011,110,101,020,LRTLRBLR*",
		"3x3:111,021,111,111,LRTLRBLR*",
		"3x3:111,021,111,111,LRTTTBBB*",
		"3x3:201,111,111,111,TTTBBBLR*",
		"3x3:111,111,111,021,LRTLRBLR*",

		"3x4:212,1202,122,2111,TTTBBBLRTLRB",
		"3x4:212,2021,212,0212,TTTBBBTTTBBB",
		"3x4:122,2111,212,1202,TTTBBBLRTLRB",
		"3x4:221,1121,212,0212,TLRBLRTLRBLR",

		"4x4:1112,1211,1112,1220,TLRTBLRBTLRTBLRB",
		"4x4:1212,0222,1122,1122,TLRTBTTBTBBTBLRB",
		"4x4:0212,0212,0122,1121,TLRTBTTBTBBTBLRB",
		"4x4:1012,0121,0112,1021,TTTTBBBBLRTTLRBB",
		"4x4:2211,2211,2220,2121,LRLRLRTTLRBBLRLR",
		"4x4:1102,2110,1012,2101,TLRTBTTBTBBTBLRB",

		"5x5:03232,22222,12322,32221,LRLRTTLRTBBTTBTTBBTBBLRB*",
		"5x5:31312,30322,13132,21232,LRLRTLRLRBTTLRTBBTTBLRBB*",
		"5x5:21222,13032,30312,30312,TTLRTBBTTBTTBBTBBTTBLRBB*",
		"5x5:12032,21212,12122,12311,TTTLRBBBLRLRLRTLRTTBLRBB*",

		"6x6:331232,222233,330332,232223,LRLRLRLRTTLRTTBBTTBBTTBBLRBBLRLRLRLR",
		"6x6:321223,113233,231313,031333,LRTLRTTTBTTBBBTBBTTTBTTBBBTBBTLRBLRB",
		"6x6:232323,033333,322233,123333,LRTLRTLRBLRBTLRTTTBLRBBBTLRLRTBLRLRB",
		"6x6:223033,323221,132133,322222,LRLRLRLRLRLRLRTTLRTTBBTTBBTTBBLRBBLR",
		"6x6:332032,323221,331132,233212,TTTTTTBBBBBBLRLRLRLRLRLRLRLRTTLRLRBB",
		"6x6:233203,122323,323113,122233,LRTTTTTTBBBBBBTLRTTTBTTBBBTBBTLRBLRB",

		"6x10:535454,3333323033,544445,3333331133,TTLRTTBBTTBBTTBBTTBBLRBBTLRTLRBTTBTTTBBTBBBLRBLRTTLRLRBBLRLR",
		"6x10:244535,3230333321,425345,3221333321,LRLRTTLRTTBBLRBBTTLRLRBBTTTLRTBBBTTBTTTBBTBBBTTBLRTBBTLRBLRB",
		"6x10:533455,3223333321,443545,3322333330,LRTLRTTTBTTBBBTBBTTTBTTBBBTBBTTTBTTBBBTBBTTTBTTBBBTBBTLRBLRB",
		"6x10:434554,3212333323,443545,3203333323,TLRTTTBTTBBBTBBTLRBTTBTTTBBTBBBTTBTTTBBTBBBLRBTTTLRTBBBLRBLR",

		"7x7:3243423,4343412,3333333,3434340,TLRTLRTBTTBTTBTBBTBBTBTTBTTBTBBTBBTBTTBTTB*BB*BB*",
		"7x7:2434223,3424322,4242314,4243403,TTTTLRTBBBBTTBLRTTBBTTTBBLRBBBTLRLRLRBLRTTLRLR*BB",
		"7x7:3341223,4332312,2423043,3323313,TTTTLRTBBBBTTBTTTTBBTBBBBTTBLRLRBBTTTLRLRBBBLRLR*",
		"7x7:3242423,4343303,3333242,3434222,TLRTLRTBTTBTTBTBBTBBTBTTBTTBTBBTBBTBTTBLRB*BBLRLR",
		"7x7:4242423,4340433,3333243,3423243,LRTTLRTTTBBTTBBBTTBBTTTBBLRBBBTTLRTTTBBTTBBBLRBB*",
		"7x7:3333233,3304343,4242323,4132433,TTLRLRTBBTLRTBTTBLRBTBBTTTTBTTBBBBTBBTLRTBLRBLRB*",
		"7x7:3322332,3324033,4231332,4232223,TTLRTLRBBTTBTTLRBBTBBTLRTBLRBTTBTLRTBBTBTTBLRB*BB",
		"7x7:3333233,3433142,4241324,4340423,TTTLRLRBBBLRTTLRTLRBBTTBLRTTBBLRTBBTTTTBTTBBBB*BB",
		"7x7:3231324,4230333,3322143,3322143,TTTLRLRBBBLRTTTLRLRBBBTTTLRTTBBBTTBBTTTBBT*BBBLRB",
		"7x7:2242321,3333103,1334221,3234112,LRLRLRTTTTTTTBBBBBBBTLRTTLRBTTBBLRTBBTLRTBLRBLRB*",
		"7x7:3333323,3424043,3342224,4341233,TTTLRLRBBBTLRTTTTBTTBBBBTBBTLRTBLRBLRBLRTTLRLR*BB",
		"7x7:3242323,4333303,1433233,3334123,LRTTLRTTTBBLRBBBTLRTTTTBLRBBBBTTTTTLRBBBBBLRLRLR*",
		"7x7:4333323,2243433,3333333,0434343,TLRLRTTBLRLRBBTTTTLRTBBBBTTBTLRTBBTBLRBLRBLRLRLR*",
		"7x7:3033343,2424142,2142433,3241432,LRTLRLRTTBTTTTBBTBBBBTTBTTLRBBTBBLRLRBLRLRLRLRLR*",
		"7x7:3242423,3343403,2324342,2434313,LRLRTTTLRLRBBBLRLRTLRTTLRBLRBBLRTLRLRLRBLRLRLRLR*",
		"7x7:0334243,3322333,3043333,4222423,LRTTTTTTTBBBBBBBLRTTTLRTTBBBLRBBTTTLRTTBBBLRBBLR*",
		"7x7:3323333,1343333,3323333,0433343,LRLRLRTTLRLRTBBTLRTBTTBTTBTBBTBBTBTTBLRBTBBLRLRB*",
		"7x7:2214342,3422223,3033333,4242222,TLRTTTTBTTBBBBTBBTTLRBLRBBTTLRLRTBBLRTTBTTLRBB*BB",
		"7x7:4242423,4340433,3334233,3413343,TTTTTTTBBBBBBBTTTLRTTBBBLRBBLRTTTTTTTBBBBBBBLRLR*",
		"7x7:0333343,3313333,1243333,4223332,LRLRLRTLRTLRTBLRBLRBTTTLRLRBBBLRLRTTLRTLRBBLRBLR*",
		"7x7:3233314,2243413,3233323,0434242,TLRLRTTBTLRTBBTBLRBLRBTLRLRTTBTTTTBBTBBBBT*BLRLRB",
		"7x7:3224323,1324243,4222414,3043423,TLRLRTTBTLRTBBTBTTBTTBTBBTBBTBLRBTTBLRTTBBLR*BBLR",
		"7x7:3331342,3424231,4240333,4242412,LRLRTLRLRTTBTTTTBBTBBBBTTBLRLRBBTLRTTTTBTTBBBB*BB",
		"7x7:3233333,2434340,2243423,3343412,LRTTLRTTTBBTTBBBTTBBTLRBBTTBTLRTBBTBTTBTTB*BB*BB*",
		"7x7:4132333,4241332,3403243,3322333,TLRTLRTBLRBTTBTLRTBBTBLRBTTBTTLRBBTBBTTTTBLRBBBB*",
		"7x7:3233333,1243433,2333333,0334343,TTTLRLRBBBTLRTLRTBLRBLRBTLRTLRTBLRBTTBTLRTBB*BLRB",
		"7x7:3404333,3323333,3322424,4331432,TTLRTTTBBTTBBBTTBBTTTBBLRBBBTLRLRLRBTTTTTT*BBBBBB",
		"7x7:3424322,3334232,4342304,4242413,LRLRTLRLRLRBLRLRLRTLRLRLRBTTLRLRTBBLRLRBTTLRLR*BB",

		"8x8:23144434,43344241,40334344,43434331,LRLRTLRTTLRTBTTBBTTBTBBTTBBTBTTBBLRBTBBTLRTTBTTBLRBBTBBTLRLRBLRB",
		"8x8:33143424,03443424,32324244,12443424,TTTLRLRTBBBTTTTBLRTBBBBTLRBLRTTBLRLRTBBTLRLRBLRBTLRLRLRTBLRLRLRB",
		"8x8:33432133,11233444,33433123,20234344,TTTLRLRTBBBLRTTBTTTLRBBTBBBTLRTBTTTBTTBTBBBTBBTBTTTBLRBTBBBLRLRB",
		"8x8:24343434,04444434,42334344,13444434,LRTLRLRTTTBLRTTBBBTTTBBTTTBBBTTBBBLRTBBTLRTTBTTBLRBBTBBTLRLRBLRB",
		"8x8:43433443,44444431,34343434,44444440,TTTLRTTTBBBLRBBBTTTLRTLRBBBLRBTTTTLRTTBBBBLRBBTTLRLRTTBBLRLRBBLR",
		"8x8:34343433,44444430,43433442,44444331,LRTLRTLRTTBTTBTTBBTBBTBBLRBLRBTTTTLRTTBBBBTTBBTTLRBBTTBBLRLRBBLR",
		"8x8:34343434,44444404,33434344,44444413,LRTTTTTTTTBBBBBBBBTLRTTTTTBTTBBBBBTBBTTTLRBLRBBBTLRLRLRTBLRLRLRB",

		"9x9:534353305,535133524,443444314,354405343,LRLRTTTTTTLRTBBBBBBTTBTLRLRTBBTBTTTTBLRBTBBBBTLRTBTTTTBLRBTBBBBTLRTBLRTTBLRBLR*BB",
		"9x9:434444345,235454543,444344444,054545453,TTLRLRLRTBBLRTTTTBTLRTBBBBTBTTBTTTTBTBBTBBBBTBTTBTTLRBTBBTBBTLRBTTBTTBTT*BB*BB*BB",
		"9x9:445423534,544442542,454504444,444444343,TLRTLRLRTBTTBTTTTBTBBTBBBBTBTTBTTTTBTBBTBBBBTBTTBLRTTBTBBLRTBBTBTTTTBTTB*BBBB*BB*",
		"9x9:454535314,453335353,545444404,533353534,LRLRLRLRTTTTLRLRTBBBBTTTTBTTLRBBBBTBBLRLRLRBTLRTLRLRTBLRBLRLRBTTLRLRLRTBBLRLRLRB*",
		"9x9:513453544,524344534,440535454,433444444,TLRTLRLRTBTTBLRTTBTBBTTTBBTBLRBBBTTBLRTTTTBBTTTBBBBTTBBBLRTTBBTTLRTBBTTBBLRBLRBB*",

		// Came from iPhone game, so guaranteed to be solvable.
		"4x5:0203,11102,2012,12011,LRLRLRTTTTBBBBTTLRBB",
		"4x5:2022,12021,2013,21201,TTTTBBBBTLRTBLRBLRLR",
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

// Cases we know are solvable, but that the solver cannot do yet. If the solver solves one, error so we know to move the test case to the TestSolve() function.
func TestSolve_Unsolvable(t *testing.T) {
	testCases := []string{
		// Generated and solved by this program. All are known solvable.
		"3x3:201,201,111,120,LRTTTBBB*",
		"3x3:112,112,121,121,*LRTLRBLR",
		"3x3:112,112,121,121,*LRTLRBLR",

		"5x4:02121,2112,20211,1212,LRLRTLRTTBLRBBTLRLRB",

		"5x5:23031,22212,22122,32121,TLRTTBTTBBTBBTTBTTBB*BBLR",
		"5x5:21222,31320,12132,22221,LRLRTLRTTBTTBBTBBTTBLRBB*",

		"5x7:32323,3222211,33223,2312221,TLRTTBLRBBTTTLRBBBTTTLRBBBLRLR*LRLR",

		"7x7:4222414,4233403,3233323,2424340,TTTLRTTBBBTTBBTLRBBLRBLRTTTTLRTBBBBTTBTTTTBB*BBBB",

		"9x9:334345444,353545441,433444444,534454540,TTTTTTTLRBBBBBBBTTLRLRLRTBBTLRTTTBTTBTTBBBTBBTBBLRTBTTBTTTTBTBBTBBBBTBLRBLRLRBLR*",

		"2x50:jh,01111110101001111000011110111110011111011101111111,me,01111101101001110100101110111101011111011011111111,LRTTBBLRLRLRTTBBLRLRLRTTBBTTBBLRTTBBLRLRTTBBLRTTBBLRLRLRTTBBTTBBLRLRTTBBTTBBLRTTBBTTBBTTBBTTBBLRTTBB",

		// Came from iPhone game, so guaranteed to be solvable.
		"4x5:3222,22122,2322,22122,LRTTTTBBBBLRTTTTBBBB",
		"4x5:3201,11211,2301,11121,LRTTTTBBBBTTTTBBBBLR",

		"5x5:11212,12211,12202,21112,TTTTTBBBBBTLRT*BLRBTLRLRB",
		"5x5:21223,31123,12232,22132,LRTLRT*BLRBLRTTLRTBBLRBLR",
		"5x5:22122,12312,22122,13212,LRLR*TTTLRBBBTTLRTBBLRBLR",
		"5x5:31222,22312,22132,23221,LRTT*TTBBTBBLRBTLRTTBLRBB",

		"6x5:122232,33132,112323,32232,LRLRLRTLRLRTBLRLRBLRTTTTLRBBBB",
		"6x5:122123,23312,221132,23312,TTLRTTBBLRBBLRLRLRTLRTLRBLRBLR",
		"6x5:022131,23202,121212,32202,TTLRLRBBLRLRTLRLRTBLRLRBLRLRLR",

		"7x9:1444133,330323222,2244242,232333211,LRLRLRTTLRLRTBBLRT*BTTTTBLRBBBBTLRTLRTBTTBTTBTBBTBBTBLRBLRBLRLR",
		"7x9:4442324,412332242,3343343,333312233,LRLRTLRTLR*BLRBTTTLRTTBBBTTBBTTTBBTTBBBTTBBTLRBBTTBLRLRBBLRLRLR",
		"7x9:3342325,331221433,2442334,331222342,TLRLRTTBLRLRBBTTLR*TTBBTTTBBTTBBBTTBBTLRBBTTBTTTTBBTBBBBLRBLRLR",

		"8x7:30122222,2232302,30311231,2321321,LRLRLRLRTLRLRLRTBTLRTLRBTBTTBLRTBTBBTLRBTBLRBLRTBLRLRLRB",
		"8x7:22423232,2223344,04243223,3231434,LRLRLRTTLRLRTTBBTTLRBBTTBBLRTTBBTTLRBBTTBBLRTTBBLRLRBBLR",

		"10x10:4444543143,4444043445,4435443324,5254135254,TTTTLRLRLRBBBBLRTTLRTLRLRTBBLRBLRLRBLRLRTTLRLRTLRTBBLRLRBLRBTTTTTLRLRTBBBBBLRLRBTLRTTLRTLRBLRBBLRBLR",

		"12x10:555213545454,5543546565,555215354544,5434465656,TTTLRLRTTLRTBBBTLRTBBTTBLRTBTTBLRBBTLRBTBBLRTLRBLRTBLRLRBLRTLRBLRTTLRLRBTLRTTBBTTLRTBLRBBLRBBLRBTTTTLRTLRTLRBBBBLRBLRBLR",
		"12x10:545232142444,5462542336,444423124453,5652451345,TTTTLRLRLRLRBBBBTLRTLRTTLRLRBTTBLRBBTLRLRBBTLRLRBLRTTLRBLRTTTLRBBTLRLRBBBLRTTBTLRTTTLRTBBTBLRBBBTTBLRBTLRLRTBBLRLRBLRLRB",

		"15x9:454544143332433,766636477,545453133524242,775455667,LRTTLRTTTLRTTLRLRBBTTBBB*TBBTTTTLRBBLRLRBLRBBBBTLRLRTTTTLRTTTTBLRLRBBBBTTBBBBLRLRTLRLRBBTTLRTTTTBLRTTLRBBLRBBBBLRTBBTTLRLRLRLRLRBLRBBLR",
		"15x9:333445434434343,786555657,442354433434434,768646575,TTLRLRLR*LRTLRTBBTLRLRTLRTBLRBLRBLRLRBTTBTLRTTLRTTTLRBBTBLRBBTTBBBTTLRBLRLRTBBTLRBBLRTTTTTBLRBLRLRLRBBBBBLRTTTLRTTTLRLRTLRBBBLRBBBLRLRB",
		"15x9:423444254334424,686756455,414443345334334,758766544,TLRLRLRLRLRLRLRB*TTTTLRTLRTTTTLRBBBBLRBTTBBBBLRTLRLRLRBBLRLRTTBLRLRLRLRTLRTBBTLRLRTLRTBLRBTTBTLRTBLRBLRLRBBTBLRBLRTTTLRTLRBLRLRLRBBBLRB",

		"15x15:858767847577878,778686866676768,686778755768787,787777775767677,TLRT*LRTTTLRTLRBLRBTLRBBBLRBLRLRLRBLRLRTTTTTTTLRLRLRLRBBBBBBBLRTLRTTLRLRTLRTLRBTTBBLRLRBLRBLRTBBTLRLRTTTTLRTBLRBLRTTBBBBLRBLRLRLRBBLRTTLRTLRLRTTTLRTBBTTBTTLRBBBLRBLRBBTBBTTTTTTTTTTLRBLRBBBBBBBBBBTTLRLRTTLRTLRLRBBLRLRBBLRBLRLR",
		"15x15:467686748786757,767787767776455,556777657878666,686877776865555,TLRTTLRTLRLR*LRBLRBBLRBTTLRTLRLRLRLRLRBBLRBTTTLRTLRLRLRTLRBBBTTBLRTTLRBLRTTTBBTLRBBLRTLRBBBLRBLRTTLRBTLRTTTTTLRBBTTTBTTBBBBBLRLRBBBTBBTLRTTLRLRTTTBLRBLRBBLRLRBBBLRTTTLRTLRLRTTTLRBBBLRBTTTTBBBLRLRLRTTBBBBTTTLRLRLRBBLRLRBBBLRLR",

		"17x17:88775759697188978,77876687774878868,87876667887189878,68784778684888688,TLRTLRLRLRLRTTLRTBLRBLRTLRTT*BBTTBTTTLRTBLRBBTLRBBTBBBLRBTLRLRBLRTTBTLRTLRBLRTTTLRBBTBLRBLRTLRBBBTTTTBTLRTTTBLRLRTBBBBTBLRBBBLRLRTBTTLRBTLRTLRTLRTBTBBTTTBLRBLRBLRBTBLRBBBLRLRLRLRTTBTTTTLRLRLRLRLRBBTBBBBLRTTTTTTLRLRBLRTLRTBBBBBBLRLRTLRBLRBTTLRLRLRTTBLRLRTTBBTLRLRTBBTTTTTBBLRBLRLRBLRBBBBBLR",
	}

	for _, testCase := range testCases {
		game, ok := magnets.Deserialize(testCase)
		if !ok {
			t.Errorf("ERROR: Unable to deserialize %s", testCase)
		}
		Solve(game)
		if game.Solved() {
			t.Errorf("ERROR: We can now solve this. Move to other test case. %s", testCase)
		}
	}
}
