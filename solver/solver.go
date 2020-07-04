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

// setFrame takes a coordinate and a polarity, sets that, and sets the other end of the frame to correspond.
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

func (cbs CBS) unsetPossibility(row, col int, r rune) {
	if val, ok := cbs[row][col][r]; ok && val {
		dirty = true
	}

	delete(cbs[row][col], r)
}

// justOne iterates through all empty cells. For any that have just one possibility left in the cbs, it sets that frame.
func (cbs CBS) justOne(game magnets.Game) {
	for cell := range game.Guess.Cells(common.Empty) {
		row, col := cell.Unpack()

		if len(cbs[row][col]) != 1 {
			continue
		}

		var r rune
		for key := range cbs[row][col] {
			r = key
		}

		cbs.setFrame(game, row, col, r)
	}
}

// rowNeeds calculates how many of a given polarity are still needed in order to be complete.
func rowNeeds(game magnets.Game, row int, r rune) int {
	needs := game.CountRow(row, r)
	has := game.Guess.CountRow(row, r)
	return needs - has
}

// colNeeds calculates how many of a given polarity are still needed in order to be complete.
func colNeeds(game magnets.Game, col int, r rune) int {
	needs := game.CountCol(col, r)
	has := game.Guess.CountCol(col, r)
	return needs - has
}

// rowHasSpaceForTotal counts how many *possible* locations are present for the given polarity.
// This includes cells that have already been solved.
func (cbs CBS) rowHasSpaceForTotal(game magnets.Game, row int, r rune) int {
	count := 0
	for col := 0; col < game.Guess.Width(); col++ {
		if cbs[row][col][r] {
			count++
		}
	}
	return count
}

// colHasSpaceForTotal counts how many *possible* locations are present for the given polarity.
// This includes cells that have already been solved.
func (cbs CBS) colHasSpaceForTotal(game magnets.Game, col int, r rune) int {
	count := 0
	for row := 0; row < game.Guess.Height(); row++ {
		if cbs[row][col][r] {
			count++
		}
	}
	return count
}

// rowHasSpaceForRemaining counts how many *possible* locations are present for the given polarity.
// This DOES NOT INCLUDE cells that have already been solved.
func (cbs CBS) rowHasSpaceForRemaining(game magnets.Game, row int, r rune) int {
	count := 0
	for col := 0; col < game.Guess.Width(); col++ {
		if game.Guess.Get(row, col) == common.Empty && cbs[row][col][r] {
			count++
		}
	}
	return count
}

// colHasSpaceForRemaining counts how many *possible* locations are present for the given polarity.
// This DOES NOT INCLUDE cells that have already been solved.
func (cbs CBS) colHasSpaceForRemaining(game magnets.Game, col int, r rune) int {
	count := 0
	for row := 0; row < game.Guess.Height(); row++ {
		if game.Guess.Get(row, col) == common.Empty && cbs[row][col][r] {
			count++
		}
	}
	return count
}

// satisfied looks at each row/col to see if there are exactly as many spaces to put a given polarity as there are needed.
func (cbs CBS) satisfied(game magnets.Game) {
	for _, category := range []rune{common.Positive, common.Negative, common.Neutral} {
		// Row is satisfied in this category? Set those frames. Clear this possibility elsewhere.
		for row := 0; row < game.Guess.Height(); row++ {
			if rowNeeds(game, row, category) == cbs.rowHasSpaceForTotal(game, row, category) {
				for col := 0; col < game.Guess.Width(); col++ {
					if cbs[row][col][category] {
						cbs.setFrame(game, row, col, category)
					}
				}
			}
		}
		// Col is satisfied in this category? Set those frames. Clear this possibility elsewhere.
		for col := 0; col < game.Guess.Width(); col++ {
			if colNeeds(game, col, category) == cbs.colHasSpaceForTotal(game, col, category) {
				for row := 0; row < game.Guess.Height(); row++ {
					if cbs[row][col][category] {
						cbs.setFrame(game, row, col, category)
					}
				}
			}
		}
	}

	return
}

// needAll checks to see if the number of pos+neg needed is equal to the number of frames that
// are still undecided. If so, none of those frames can be neutral.
// NOTE: Once doubleSingle() is written this function will no longer be needed.
func (cbs CBS) needAll(game magnets.Game) {
	// If there are any that we know what they must be, but have not set them
	// yet, do that now. Otherwise, the count will be off.
	cbs.justOne(game)

	for _, category := range []rune{common.Positive, common.Negative} {
		// Row (#frames remaining that can be category) == (#squares needed).
		for row := 0; row < game.Guess.Height(); row++ {
			needs := rowNeeds(game, row, category)
			if needs == 0 {
				continue
			}
			provides := 0
			for col := 0; col < game.Guess.Width(); col++ {
				direction := game.GetFrame(row, col)
				switch direction {
				case common.Right:
					continue
				case common.Left:
					if cbs[row][col][category] || cbs[row][col+1][category] {
						provides++
					}
				default:
					if cbs[row][col][category] {
						provides++
					}
				}
			}
			if needs == provides {
				// All are needed for signs. None can be neutral.
				for col := 0; col < game.Guess.Width(); col++ {
					cbs.unsetPossibility(row, col, common.Neutral)
				}
			}
		}

		// Col (#frames remaining that can be category) == (#squares needed).
		for col := 0; col < game.Guess.Width(); col++ {
			needs := colNeeds(game, col, category)
			if needs == 0 {
				continue
			}
			provides := 0
			for row := 0; row < game.Guess.Height(); row++ {
				direction := game.GetFrame(row, col)
				switch direction {
				case common.Down:
					continue
				case common.Up:
					if cbs[row][col][category] || cbs[row+1][col][category] {
						provides++
					}
				default:
					if cbs[row][col][category] {
						provides++
					}
				}
			}
			if needs == provides {
				// All are needed for signs. None can be neutral.
				for row := 0; row < game.Guess.Height(); row++ {
					cbs.unsetPossibility(row, col, common.Neutral)
				}
			}
		}
	}
}

// doubleSingle() looks for cases where, based on the length of the frame (1 or 2 cells in this row/col) we know the frame can/cannot be a magnet. For instance if we need 2 polarities (1 plus and 1 minus) and there is 1 horizontal and 1 vertical frame we know the vertical frame cannot have a polarity.
func (cbs CBS) doubleSingle(game magnets.Game) {

	// Enumerate each of the combinations of frames (that are undecided) in the
	// row/col that will satisfy the pos+neg count conditions. If there is a
	// frame that is not in any of those combinations then that frame must not
	// be a magnet (and therefore is neutral). Once we mark it as neutral (and the
	// remaining frames as !neutral) the rest of the CBS will figure out whether
	// we now know the other frames, so we do not need to also do that here.
	// ALSO: If there is a frame that is in *every* one of those combinations then
	// it must be non-neutral.

	// If all of the frames are needed to make the solution, then none of them
	// can be neutral. This doesn't find us a solution directly, but might be
	// enough info to let the CBS tease out a solution.
	// NOTE: This turned out to be true. For small boards it works fine.
	// See the needAll() function. But, it will no longer be needed once this
	// function is complete.

}

// resolveNeighbors() propagates any constraint a cell has (like it can only be negative) to its neighbor (which can then only be positive).
func (cbs CBS) resolveNeighbors(game magnets.Game) {
	for cell := range game.Guess.Cells() {
		row, col := cell.Unpack()
		rowEnd, colEnd := game.GetFrameEnd(row, col)
		if rowEnd == -1 && colEnd == -1 {
			continue
		}
		for _, r := range []rune{common.Positive, common.Negative, common.Neutral} {
			// If r is missing from this end, its opposite cannot be in the other end.
			if !cbs[row][col][r] {
				cbs.unsetPossibility(rowEnd, colEnd, common.Negate(r))
			}
		}
	}

	// If a cell borders one that is already identified, update the cbs.
	for cell := range game.Guess.Cells() {
		row, col := cell.Unpack()
		for _, adj := range board.Adjacents {
			r, c := adj.Unpack()
			switch game.Guess.Get(row+r, col+c) {
			case common.Positive:
				cbs.unsetPossibility(row, col, common.Positive)
			case common.Negative:
				cbs.unsetPossibility(row, col, common.Negative)
			}
		}
	}
}

// zeroInRow looks for rows that have no positives or that have no negatives and removes those possibilities from the cbs.
func (cbs CBS) zeroInRow(game magnets.Game) {
	for _, category := range []rune{common.Positive, common.Negative} {
		for row := 0; row < game.Guess.Height(); row++ {
			if game.CountRow(row, category) == 0 {
				// Remove all instances of 'category' from the row in the cbs
				for col := 0; col < game.Guess.Width(); col++ {
					cbs.unsetPossibility(row, col, category)
				}
			}
		}
	}
}

// zeroInCol looks for columns that have no positives or that have no negatives and removes those possibilities from the cbs.
func (cbs CBS) zeroInCol(game magnets.Game) {
	for _, category := range []rune{common.Positive, common.Negative} {
		for col := 0; col < game.Guess.Width(); col++ {
			if game.CountCol(col, category) == 0 {
				// Remove all instances of 'category' from the column in the cbs
				for row := 0; row < game.Guess.Height(); row++ {
					cbs.unsetPossibility(row, col, category)
				}
			}
		}
	}
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

// Solve attempts to find a solution for the game, or gives up if it cannot.
func Solve(game magnets.Game) {
	cbs := new(game)

	cbs.zeroInRow(game)
	cbs.zeroInCol(game)

	attempts := 0
	for {
		dirty = false

		err := cbs.validate(game)
		if err != nil {
			fmt.Println(err)
			game.Print()
			break
		}

		cbs.satisfied(game)
		cbs.resolveNeighbors(game)
		cbs.doubleSingle(game)
		cbs.needAll(game)
		cbs.justOne(game)

		if !dirty {
			break
		}

		attempts++
		if attempts > 500 {
			break
		}
	}
}
