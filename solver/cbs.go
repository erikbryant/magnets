package solver

// This package implements a constraints-based solver.
//
// The CBS solver holds a matrix representation that mimics the playing grid.
//
// It first finds whether there are any walls on the playing grid. If there are,
// it marks those in its representation as being walls.
//
// Then it marks the rest of the cells as having all of the possibilities (positive,
// negative, or neutral).
//
// As the solver works through the information it is given it eliminates possibilities
// as it proves them non-possible. When the list of possibilities for a cell drops to
// just one, the solver sets that in the Guess grid as a known entity.

import (
	"fmt"
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
		r := game.Guess.Get(row, col)
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

// setFrame takes a coordinate and a polarity, sets that, and sets the other end
// of the frame to correspond. This is different from the other implementations
// in that it also keeps track of whether the board is dirty and updates the CBS.
func (cbs CBS) setFrame(game magnets.Game, row, col int, r rune) {
	rowEnd, colEnd := game.GetFrameEnd(row, col)

	if r != game.Guess.Get(row, col) || common.Negate(r) != game.Guess.Get(rowEnd, colEnd) {
		dirty = true
	}

	// Set this end of the frame.
	game.Guess.Set(row, col, r)
	cbs[row][col] = map[rune]bool{r: true}

	// Set the other end of the frame.
	game.Guess.Set(rowEnd, colEnd, common.Negate(r))
	cbs[rowEnd][colEnd] = map[rune]bool{common.Negate(r): true}
}

// unsetPossibility removes the given rune from the CBS' list of potential
// cell values.
func (cbs CBS) unsetPossibility(row, col int, r rune) {
	if val, ok := cbs[row][col][r]; ok && val {
		dirty = true
	}

	delete(cbs[row][col], r)
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
func (cbs CBS) rowHasSpaceForTotal(game magnets.Game, row int, r rune) int {
	count := 0
	for col := 0; col < game.Guess.Width(); col++ {
		if cbs[row][col][r] {
			count++
		}
	}
	return count
}

// colHasSpaceForTotal counts how many *possible* locations are present for the
// given polarity. This includes cells that have already been solved.
func (cbs CBS) colHasSpaceForTotal(game magnets.Game, col int, r rune) int {
	count := 0
	for row := 0; row < game.Guess.Height(); row++ {
		if cbs[row][col][r] {
			count++
		}
	}
	return count
}

// rowHasSpaceForRemaining counts how many *possible* locations are present for
// the given polarity. This DOES NOT INCLUDE cells that have already been solved.
func (cbs CBS) rowHasSpaceForRemaining(game magnets.Game, row int, r rune) int {
	count := 0
	for col := 0; col < game.Guess.Width(); col++ {
		if game.Guess.Get(row, col) == common.Empty && cbs[row][col][r] {
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
		if game.Guess.Get(row, col) == common.Empty && cbs[row][col][r] {
			count++
		}
	}
	return count
}

// validate returns an error if the game or the CBS is inconsistent.
func (cbs CBS) validate(game magnets.Game) error {
	if !game.Valid() {
		return fmt.Errorf("Invalid game board state detected")
	}

	for cell := range game.Guess.Cells(common.Positive, common.Negative, common.Neutral) {
		row, col := cell.Unpack()
		r := game.Guess.Get(row, col)
		// This is already solved, so the CBS should only have r in it.
		for key := range cbs[row][col] {
			if key != r {
				return fmt.Errorf("ERROR: CBS %d, %d had extraneous '%c'", row, col, key)
			}
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