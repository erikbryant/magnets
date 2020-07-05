package solver

import (
	"fmt"
	"github.com/erikbryant/magnets/board"
	"github.com/erikbryant/magnets/common"
	"github.com/erikbryant/magnets/magnets"
)

// justOne iterates through all empty cells. For any that have just one
// possibility left in the cbs, it sets that frame.
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

// satisfied looks at each row/col to see if there are exactly as many spaces
// to put a given polarity as there are needed.
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

// needAll checks to see if the number of pos+neg needed is equal to the number
// of frames that are still undecided. If so, none of those frames can be neutral.
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

// oddRowAllMagnets checks to see if the entire row is full of magnets. If it
// is, and if the row length is odd, then we know the pattern of the magnets.
func (cbs CBS) oddRowAllMagnets(game magnets.Game) {
	// If the length is not odd then there is nothing we can determine.
	if game.Guess.Width()%2 == 0 {
		return
	}

	for row := 0; row < game.Guess.Height(); row++ {
		if game.CountRow(row, common.Positive)+game.CountRow(row, common.Negative) == game.Guess.Width() {
			// There is only one way the magnets can be arranged.
			// The polarity with the larger count goes first.
			var polarity rune

			if game.CountRow(row, common.Positive) > game.CountRow(row, common.Negative) {
				polarity = common.Positive
			} else {
				polarity = common.Negative
			}

			for col := 0; col < game.Guess.Width(); col++ {
				game.Guess.Set(row, col, polarity)
				polarity = common.Negate(polarity)
				cbs.unsetPossibility(row, col, polarity)
				cbs.unsetPossibility(row, col, common.Neutral)
			}

			dirty = true
		}
	}
}

// oddColAllMagnets checks to see if the entire col is full of magnets. If it
// is, and if the col length is odd, then we know the pattern of the magnets.
func (cbs CBS) oddColAllMagnets(game magnets.Game) {
	// If the length is not odd then there is nothing we can determine.
	if game.Guess.Height()%2 == 0 {
		return
	}

	for col := 0; col < game.Guess.Width(); col++ {
		if game.CountCol(col, common.Positive)+game.CountCol(col, common.Negative) == game.Guess.Height() {
			// There is only one way the magnets can be arranged.
			// The polarity with the larger count goes first.
			var polarity rune

			if game.CountCol(col, common.Positive) > game.CountCol(col, common.Negative) {
				polarity = common.Positive
			} else {
				polarity = common.Negative
			}

			for row := 0; row < game.Guess.Height(); row++ {
				game.Guess.Set(row, col, polarity)
				polarity = common.Negate(polarity)
				cbs.unsetPossibility(row, col, polarity)
				cbs.unsetPossibility(row, col, common.Neutral)
			}

			dirty = true
		}
	}
}

// doubleSingle() looks for cases where, based on the length of the frame
// (1 or 2 cells in this row/col) we know the frame can/cannot be a magnet.
// For instance if we need 2 polarities (1 plus and 1 minus) and there is
// 1 horizontal and 1 vertical frame we know the vertical frame cannot have
// a polarity.
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

// resolveNeighbors() propagates any constraint a cell has (like it can only be
// negative) to its neighbor (which can then only be positive).
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

// zeroInRow looks for rows that have no positives or that have no negatives
// and removes those possibilities from the cbs.
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

// zeroInCol looks for columns that have no positives or that have no negatives
// and removes those possibilities from the cbs.
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

// Solve attempts to find a solution for the game, or gives up if it cannot.
func Solve(game magnets.Game) {
	cbs := new(game)

	cbs.zeroInRow(game)
	cbs.zeroInCol(game)

	cbs.oddRowAllMagnets(game)
	cbs.oddColAllMagnets(game)

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
			fmt.Println("WARNING: unable to solve game after >500 attempts")
			game.Print()
			break
		}
	}
}
