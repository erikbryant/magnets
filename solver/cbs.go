package solver

// This package implements a constraints-based solver.
//
// The CBS solver holds a matrix representation that mimics the playing grid.
//
// It first finds whether there are any walls on the playing grid. If there are,
// it marks those in its representation as being walls.
//
// Then it marks the rest of the cells as having all the possibilities (positive,
// negative, or neutral).
//
// As the solver works through the information it is given it eliminates possibilities
// as it proves them non-possible. When the list of possibilities for a cell drops to
// just one, the solver sets that in the Guess grid as a known entity.
//
// The CBS should never drop to zero possibilities. The end state is that there is
// exactly one possibility.

import (
	"fmt"

	"github.com/erikbryant/magnets/board"
	"github.com/erikbryant/magnets/common"
	"github.com/erikbryant/magnets/magnets"
)

// CBS is the constraint-based solver representation.
type CBS [][]map[rune]bool

var (
	dirty = false
)

// new takes a game and returns a new, initialized constraint-based solver object for that game.
func new(game magnets.Game) CBS {
	cbs := make(CBS, game.Guess.Height())

	for row := 0; row < game.Guess.Height(); row++ {
		cbs[row] = make([]map[rune]bool, game.Guess.Width())
		for col := 0; col < game.Guess.Width(); col++ {
			cbs[row][col] = make(map[rune]bool)
		}
	}

	for cell := range game.Guess.Cells() {
		row, col := cell.Unpack()
		r := game.Guess.Get(row, col, false)
		if r == common.Wall {
			cbs[row][col] = map[rune]bool{r: true}
			continue
		}
		if r != common.Empty {
			fmt.Printf("ERROR: %d, %d was already set to '%c'\n", row, col, r)
			continue
		}

		// Each cell is a set of possibilities. At the start, each case is possible.
		cbs[row][col][common.Positive] = true
		cbs[row][col][common.Negative] = true
		cbs[row][col][common.Neutral] = true
	}

	return cbs
}

// getOnlyPossibility returns the only remaining value in the CBS, or if there
// is not only one possibility, it panics.
func (cbs CBS) getOnlyPossibility(row, col int) (r rune) {
	if len(cbs[row][col]) != 1 {
		msg := fmt.Sprintf("There is not only one possibility %v", cbs[row][col])
		panic(msg)
	}

	for r = range cbs[row][col] {
	}

	return
}

// setFrame takes a coordinate and a polarity, sets that, and sets the other end
// of the frame to correspond. This is different from the other implementations
// in that it also keeps track of whether the board is dirty and updates the CBS.
func (cbs CBS) setFrame(game magnets.Game, row, col int, r rune) {
	rowEnd, colEnd := game.GetFrameEnd(row, col)

	if r != game.Guess.Get(row, col, false) || common.Negate(r) != game.Guess.Get(rowEnd, colEnd, false) {
		dirty = true
	}

	// Set this end of the frame.
	game.Guess.Set(row, col, r, false)
	cbs[row][col] = map[rune]bool{r: true}

	// Set the other end of the frame.
	game.Guess.Set(rowEnd, colEnd, common.Negate(r), false)
	cbs[rowEnd][colEnd] = map[rune]bool{common.Negate(r): true}
}

// decided returns true if the final value of the cell has been decided,
// false otherwise.
func (cbs CBS) decided(row, col int) bool {
	return len(cbs[row][col]) == 1
}

// possibility returns true if the given rune is still a possibility, false otherwise.
func (cbs CBS) possibility(row, col int, r rune) bool {
	if val, ok := cbs[row][col][r]; ok {
		return val
	}

	return false
}

// unsetHelper does the actual work of removing the given possibility from the cbs.
func (cbs CBS) unsetHelper(row, col int, r rune) {
	if cbs.possibility(row, col, r) {
		dirty = true
	}

	delete(cbs[row][col], r)

	if len(cbs[row][col]) == 0 {
		msg := fmt.Sprintf("All possibilities have been deleted from cbs %d, %d", row, col)
		panic(msg)
	}
}

// unsetPossibility removes the given rune from the CBS' list of potential
// cell values.
func (cbs CBS) unsetPossibility(game magnets.Game, row, col int, r rune) {
	cbs.unsetHelper(row, col, r)
	rowEnd, colEnd := game.GetFrameEnd(row, col)
	if rowEnd == -1 || colEnd == -1 {
		return
	}
	cbs.unsetHelper(rowEnd, colEnd, common.Negate(r))
}

// unsetPossibilityRow removes the given rune from the CBS' list of potential
// cell values for an entire row.
func (cbs CBS) unsetPossibilityRow(game magnets.Game, row int, r rune) {
	for col := 0; col < game.Guess.Width(); col++ {
		// Never remove the last possibility.
		if len(cbs[row][col]) > 1 {
			cbs.unsetPossibility(game, row, col, r)
		}
	}
}

// unsetPossibilityCol removes the given rune from the CBS' list of potential
// cell values for an entire col.
func (cbs CBS) unsetPossibilityCol(game magnets.Game, col int, r rune) {
	for row := 0; row < game.Guess.Height(); row++ {
		// Never remove the last possibility.
		if len(cbs[row][col]) > 1 {
			cbs.unsetPossibility(game, row, col, r)
		}
	}
}

// rowNeeds calculates how many of a given polarity are still needed in order
// to be complete.
func rowNeeds(game magnets.Game, row int, r rune) int {
	needs := game.CountRow(row, r)
	has := game.Guess.CountRow(row, r)
	return needs - has
}

// colNeeds calculates how many of a given polarity are still needed in order
// to be complete.
func colNeeds(game magnets.Game, col int, r rune) int {
	needs := game.CountCol(col, r)
	has := game.Guess.CountCol(col, r)
	return needs - has
}

// rowHasSpaceForTotal counts how many *possible* locations are present for the
// given polarity. This includes cells that have already been solved.
// Note that each horizontal frame in this row only adds one possibility
// since both ends of the magnet cannot be the same polarity.
func (cbs CBS) rowHasSpaceForTotal(game magnets.Game, row int, r rune) int {
	count := 0

	for col := 0; col < game.Guess.Width(); col++ {
		// For negatable possibilities the counting is a little more complicated.
		if r != common.Negate(r) {
			direction := game.GetFrame(row, col)
			switch direction {
			case common.Right:
				continue
			case common.Left:
				if cbs[row][col][r] || cbs[row][col+1][r] {
					count++
				}
			default:
				if cbs[row][col][r] {
					count++
				}
			}
		} else {
			if cbs[row][col][r] {
				count++
			}
		}
	}

	return count
}

// colHasSpaceForTotal counts how many *possible* locations are present for the
// given polarity. This includes cells that have already been solved.
// Note that each vertical frame in this row only adds one possibility
// since both ends of the magnet cannot be the same polarity.
func (cbs CBS) colHasSpaceForTotal(game magnets.Game, col int, r rune) int {
	count := 0

	for row := 0; row < game.Guess.Height(); row++ {
		// For negatable possibilities the counting is a little more complicated.
		if r != common.Negate(r) {
			direction := game.GetFrame(row, col)
			switch direction {
			case common.Down:
				continue
			case common.Up:
				if cbs[row][col][r] || cbs[row+1][col][r] {
					count++
				}
			default:
				if cbs[row][col][r] {
					count++
				}
			}
		} else {
			if cbs[row][col][r] {
				count++
			}
		}
	}

	return count
}

// rowHasSpaceForRemaining counts how many *possible* locations are present for
// the given polarity. This DOES NOT INCLUDE cells that have already been solved.
func (cbs CBS) rowHasSpaceForRemaining(game magnets.Game, row int, r rune) int {
	count := 0
	for col := 0; col < game.Guess.Width(); col++ {
		if game.Guess.Get(row, col, false) == common.Empty && cbs[row][col][r] {
			count++
		}
	}
	return count
}

// colHasSpaceForRemaining counts how many *possible* locations are present for
// the given polarity. This DOES NOT INCLUDE cells that have already been solved.
func (cbs CBS) colHasSpaceForRemaining(game magnets.Game, col int, r rune) int {
	count := 0
	for row := 0; row < game.Guess.Height(); row++ {
		if game.Guess.Get(row, col, false) == common.Empty && cbs[row][col][r] {
			count++
		}
	}
	return count
}

// validate returns an error if the game or the CBS is inconsistent.
func (cbs CBS) validate(game magnets.Game) error {
	if !game.Valid() {
		return fmt.Errorf("invalid game board state detected")
	}

	// Validate that there are no two identical signs next to each other.
	for cell := range game.Guess.Cells(common.Positive, common.Negative) {
		row, col := cell.Unpack()
		guess := game.Guess.Get(row, col, false)
		for _, adj := range board.Adjacents {
			r, c := adj.Unpack()
			if game.Guess.Get(row+r, col+c, false) == guess {
				return fmt.Errorf("ERROR: '%c' sign at %d, %d is not consistent", guess, row, col)
			}
		}
	}

	for cell := range game.Guess.Cells(common.Positive, common.Negative, common.Neutral) {
		row, col := cell.Unpack()
		r := game.Guess.Get(row, col, false)
		// This is already solved, so the CBS should only have r in it.
		for key := range cbs[row][col] {
			if key != r {
				return fmt.Errorf("ERROR: CBS %d, %d had extraneous '%c'", row, col, key)
			}
		}
		if len(cbs[row][col]) != 1 {
			return fmt.Errorf("ERROR: CBS %d, %d has wrong length. Expected 1 got %d", row, col, len(cbs[row][col]))
		}
	}

	// Validate that the CBS contains only expected possibilities.
	for row := range cbs {
		for col := range cbs[row] {
			for key := range cbs[row][col] {
				switch key {
				case common.Positive:
					continue
				case common.Negative:
					continue
				case common.Neutral:
					continue
				case common.Wall:
					continue
				default:
					return fmt.Errorf("ERROR: CBS %d, %d has unexpected '%c'", row, col, key)
				}
			}
		}
	}

	return nil
}

// print prints a formatted representation of the cbs.
func (cbs CBS) print() {
	fmt.Printf("CBS (%dx%d)\n", len(cbs[0]), len(cbs))

	fmt.Printf("   + ")
	for i := 0; i < len(cbs[0]); i++ {
		fmt.Printf("―")
	}
	fmt.Printf("\n")

	for row := range cbs {
		fmt.Printf("   | ")
		for col := range cbs[row] {
			if len(cbs[row][col]) == 1 {
				for r := range cbs[row][col] {
					fmt.Printf("%c", r)
				}
			} else {
				fmt.Printf("%d", len(cbs[row][col]))
			}
		}
		fmt.Println(" |")
	}

	fmt.Printf("     ")
	for i := 0; i < len(cbs[0]); i++ {
		fmt.Printf("―")
	}
	fmt.Printf(" -\n")
}
