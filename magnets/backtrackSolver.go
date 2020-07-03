package magnets

import (
	"github.com/erikbryant/magnets/common"
)

// exceedLimits returns true if adding this rune to the row will exceed
// the legal positive/negative count for this row, false otherwise.
func (game *Game) exceedsRowLimits(row, col int, r rune) bool {
	if game.Guess.CountRow(row, common.Positive) > game.CountRow(row, common.Positive) {
		return true
	}

	if game.Guess.CountRow(row, common.Negative) > game.CountRow(row, common.Negative) {
		return true
	}

	if game.Guess.CountRow(row, common.Neutral)+game.Guess.CountRow(row, common.Wall) > game.CountRow(row, common.Neutral) {
		return true
	}

	return false
}

// blankCell sets the cell and its other end to be empty.
func (game *Game) blankCell(row, col int) {
	rowEnd, colEnd := game.GetFrameEnd(row, col)
	game.Guess.Set(row, col, common.Empty)
	game.Guess.Set(rowEnd, colEnd, common.Empty)
}

// setCell attempts to set the given cell (and its other end). If it is a legal move
// it sets the cells and returns true, false otherwise.
func (game *Game) setCell(row, col int, r rune) bool {
	rowEnd, colEnd := game.GetFrameEnd(row, col)

	// Should we be concerned about the polarity of the neighbors?
	if common.Negate(r) != r {
		if game.Guess.Get(row-1, col) == r {
			return false
		}
		if game.Guess.Get(row, col-1) == r {
			return false
		}
		// TODO: There is an optimization to also check the other end
		// of the frame to see if it would be an invalid placement.
	}

	game.Guess.Set(row, col, r)
	game.Guess.Set(rowEnd, colEnd, common.Negate(r))

	if game.exceedsRowLimits(row, col, r) {
		game.blankCell(row, col)
		return false
	}

	return true
}

// CountSolutions returns the total number of valid solutions for the given game.
func (game *Game) CountSolutions(row, col int) int {
	solutions := 0

	for {
		if col >= game.Guess.Width() {
			col = 0
			row++
		}
		if row >= game.Guess.Height() {
			if game.Valid() && game.Solved() {
				solutions++
			}

			return solutions
		}

		if game.frames.Get(row, col) == common.Up || game.frames.Get(row, col) == common.Left {
			break
		}
		col++
	}

	if game.setCell(row, col, common.Positive) {
		solutions += game.CountSolutions(row, col+1)
	}

	if game.setCell(row, col, common.Negative) {
		solutions += game.CountSolutions(row, col+1)
	}

	if game.setCell(row, col, common.Neutral) {
		solutions += game.CountSolutions(row, col+1)
	}

	game.blankCell(row, col)

	return solutions
}

// singleSolution returns true if there is only one solution for the game, false otherwise.
func (game *Game) singleSolution() bool {
	return game.CountSolutions(0, 0) == 1
}
