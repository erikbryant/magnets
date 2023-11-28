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

		if len(cbs[row][col]) == 1 {
			cbs.setFrame(game, row, col, cbs.getOnlyPossibility(row, col))
		}
	}
}

// satisfied looks at each row/col to see if there are exactly as many spaces
// to put a given polarity as there are needed.
func (cbs CBS) satisfied(game magnets.Game) {
	for _, category := range []rune{common.Positive, common.Negative, common.Neutral} {
		// Row is satisfied in this category? Set those frames. Clear this
		// possibility elsewhere.
		for row := 0; row < game.Guess.Height(); row++ {
			if rowNeeds(game, row, category) == cbs.rowHasSpaceForTotal(game, row, category) {
				for col := 0; col < game.Guess.Width(); col++ {
					if cbs[row][col][category] {
						if category == common.Neutral {
							cbs.setFrame(game, row, col, category)
						} else {
							direction := game.GetFrame(row, col)
							switch direction {
							case common.Up:
								cbs[row][col] = map[rune]bool{category: true}
							case common.Down:
								cbs[row][col] = map[rune]bool{category: true}
							case common.Left:
								cbs.unsetPossibility(game, row, col, common.Neutral)
							case common.Right:
								cbs.unsetPossibility(game, row, col, common.Neutral)
							}
						}
					}
				}
			}
		}

		// Col is satisfied in this category? Set those frames. Clear this
		// possibility elsewhere.
		for col := 0; col < game.Guess.Width(); col++ {
			if colNeeds(game, col, category) == cbs.colHasSpaceForTotal(game, col, category) {
				for row := 0; row < game.Guess.Height(); row++ {
					if cbs[row][col][category] {
						if category == common.Neutral {
							cbs.setFrame(game, row, col, category)
						} else {
							direction := game.GetFrame(row, col)
							switch direction {
							case common.Up:
								cbs[row][col] = map[rune]bool{category: true}
							case common.Down:
								cbs[row][col] = map[rune]bool{category: true}
							case common.Left:
								cbs.unsetPossibility(game, row, col, common.Neutral)
							case common.Right:
								cbs.unsetPossibility(game, row, col, common.Neutral)
							}
						}
					}
				}
			}
		}
	}
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
			provides := cbs.rowHasSpaceForTotal(game, row, category)
			if needs == provides {
				// All that remain undecided are needed for signs. None can be neutral.
				cbs.unsetPossibilityRow(game, row, common.Neutral)
			}
		}

		// Col (#frames remaining that can be category) == (#squares needed).
		for col := 0; col < game.Guess.Width(); col++ {
			needs := colNeeds(game, col, category)
			if needs == 0 {
				continue
			}
			provides := cbs.colHasSpaceForTotal(game, col, category)
			if needs == provides {
				// All that remain undecided are needed for signs. None can be neutral.
				cbs.unsetPossibilityCol(game, col, common.Neutral)
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
				game.Guess.Set(row, col, polarity, false)
				polarity = common.Negate(polarity)
				cbs.unsetPossibility(game, row, col, polarity)
				cbs.unsetPossibility(game, row, col, common.Neutral)
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
				game.Guess.Set(row, col, polarity, false)
				polarity = common.Negate(polarity)
				cbs.unsetPossibility(game, row, col, polarity)
				cbs.unsetPossibility(game, row, col, common.Neutral)
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
	// If a cell borders one whose polarity is already identified, update the cbs.
	for cell := range game.Guess.Cells() {
		row, col := cell.Unpack()
		rowEnd, colEnd := game.GetFrameEnd(row, col)
		if rowEnd == -1 || colEnd == -1 {
			continue
		}
		for _, adj := range board.Adjacents {
			r, c := adj.Unpack()
			switch game.Guess.Get(row+r, col+c, false) {
			case common.Positive:
				cbs.unsetPossibility(game, row, col, common.Positive)
				cbs.unsetPossibility(game, rowEnd, colEnd, common.Negative)
			case common.Negative:
				cbs.unsetPossibility(game, row, col, common.Negative)
				cbs.unsetPossibility(game, rowEnd, colEnd, common.Positive)
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
				cbs.unsetPossibilityRow(game, row, category)
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
				cbs.unsetPossibilityCol(game, col, category)
			}
		}
	}
}

func (cbs CBS) checker(game magnets.Game, msg string) {
	// fmt.Printf("\n\n\n")
	// fmt.Println("----> State coming out of", msg, "<-----")
	// game.Print()
	// cbs.print()
	err := cbs.validate(game)
	if err != nil {
		game.Print()
		cbs.print()
		panic(err)
	}
}

// Solve attempts to find a solution for the game, or gives up if it cannot.
func Solve(game magnets.Game) {
	cbs := new(game)

	cbs.zeroInRow(game)
	cbs.checker(game, "zeroInRow")

	cbs.zeroInCol(game)
	cbs.checker(game, "zeroInCol")

	cbs.oddRowAllMagnets(game)
	cbs.checker(game, "oddRowAllMagnets")

	cbs.oddColAllMagnets(game)
	cbs.checker(game, "oddColAllMagnets")

	attempts := 0
	for {
		dirty = false

		// cbs.satisfied(game) // This is definitely buggy
		cbs.checker(game, "satisfied")

		cbs.resolveNeighbors(game)
		cbs.checker(game, "resolveNeighbors")

		cbs.doubleSingle(game)
		cbs.checker(game, "doubleSingle")

		// cbs.needAll(game) // This appears to be buggy
		cbs.checker(game, "needAll")

		cbs.justOne(game)
		cbs.checker(game, "justOne")

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
