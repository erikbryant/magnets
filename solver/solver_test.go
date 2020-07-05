package solver

import (
	"bufio"
	"github.com/erikbryant/magnets/common"
	"github.com/erikbryant/magnets/magnets"
	"os"
	"strings"
	"testing"
)

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
			continue
		}

		Solve(game)

		if game.Solved() != expected {
			t.Errorf("ERROR: For %s expected solved to be %t", testCase, expected)
		}
	}
}

// This is becoming a regression test. If the run time gets too high, move out of the unit tests.
func TestSolve(t *testing.T) {
	// helper(t, "testcases_solve.txt", true)
	// helper(t, "testcases_solve_fail.txt", false)
}
