package magnets

//
// Definitions:
//   Game - Each of the layers that makes up a playable game.
//   Board - A single layer. A grid r x c in dimension.
//   Cell - A single square in a given board.
//

// TODO: Replace the printf's with err.

import (
	"fmt"
	"github.com/erikbryant/magnets/board"
	"github.com/erikbryant/magnets/common"
)

// Game contains all of the representations to hold state for a game of magnets.
type Game struct {
	frames board.Board
	grid   board.Board
	Guess  board.Board

	// The row/col count cache
	colPos []int
	rowPos []int
	colNeg []int
	rowNeg []int

	// The string, if deserialized
	serial string
}

// Print prints an ASCII representation of the board.
func (game *Game) Print() {
	fmt.Printf("\n")
	game.frames.Print("Frames", game.rowPos, game.rowNeg, game.colPos, game.colNeg)
	game.grid.Print("Grid", game.rowPos, game.rowNeg, game.colPos, game.colNeg)
	game.Guess.Print("Guess", []int{}, []int{}, []int{}, []int{})
}

// Valid returns true if the game state is valid, false otherwise.
func (game *Game) Valid() bool {
	// Validate board size bounds.
	if game.grid.Width() <= 0 || game.grid.Height() <= 0 {
		fmt.Println("ERROR: dimensions out of bounds.", game.grid.Width(), "x", game.grid.Height())
		return false
	}
	if game.grid.Width()*game.grid.Height() <= 1 {
		fmt.Println("ERROR: dimensions are too small.", game.grid.Width(), "x", game.grid.Height())
		return false
	}

	// Validate that in frames, every cell is a frame or a wall.
	// Validate that in frames, each frame has both of its ends.
	for cell := range game.frames.Cells() {
		row, col := cell.Unpack()
		cell := game.frames.Get(row, col)
		if cell == common.Wall {
			// A wall has no other end.
			continue
		}
		if cell == common.Empty {
			// The board was not completely filled.
			return false
		}
		rowEnd, colEnd := game.GetFrameEnd(row, col)
		if rowEnd == -1 && colEnd == -1 {
			fmt.Printf("ERROR: At frames %d, %d found unexpected '%c'\n", row, col, cell)
			return false
		}
		adjacent := game.frames.Get(rowEnd, colEnd)
		if adjacent != common.Negate(cell) {
			fmt.Printf("ERROR: At frames %d, %d expected '%c', found '%c'\n", rowEnd, colEnd, common.Negate(cell), adjacent)
			return false
		}
	}

	// Validate that everything in grid is: positive, negative, neutral, or a wall.
	for cell := range game.grid.Cells() {
		row, col := cell.Unpack()
		grid := game.grid.Get(row, col)
		switch grid {
		case common.Positive:
		case common.Negative:
		case common.Neutral:
		case common.Wall:
		case common.Empty:
			// Valid() is called before the board is populated,
			// so Empty can also be a valid case. Would be nice
			// to fix that.
		default:
			fmt.Printf("ERROR: unexpected grid cell: '%c' at %d x %d\n", grid, row, col)
			return false
		}
	}

	// Validate that each frame has its expected values in grid.
	for frame := range game.Frames() {
		row, col := frame.Unpack()
		grid := game.grid.Get(row, col)
		rowEnd, colEnd := game.GetFrameEnd(row, col)
		found := game.grid.Get(rowEnd, colEnd)
		if common.Negate(grid) != found {
			fmt.Printf("ERROR: wrong sign at %d, %d! Expected '%c' got '%c'\n", row, col, common.Negate(grid), found)
			return false
		}
	}

	// Validate that there are no two identical signs next to each other.
	for cell := range game.grid.Cells(common.Positive, common.Negative) {
		row, col := cell.Unpack()
		grid := game.grid.Get(row, col)
		for _, adj := range board.Adjacents {
			r, c := adj.Unpack()
			if game.grid.Get(row+r, col+c) == grid {
				fmt.Printf("ERROR: '%c' sign at %d, %d is not consistent\n", grid, row, col)
				return false
			}
		}
	}

	return true
}

// SetDomino sets both ends of a domino, given one end.
func (game *Game) SetDomino(l board.Board, row, col int, r rune) {
	if l.Get(row, col) != common.Empty {
		fmt.Printf("WARNING: assigning '%c' to non-empty cell %d, %d = '%c'\n", r, row, col, l.Get(row, col))
	}

	// Set this end of the frame.
	l.Set(row, col, r)

	// Set the other end of the frame.
	rowEnd, colEnd := game.GetFrameEnd(row, col)
	if rowEnd == -1 && colEnd == -1 {
		return
	}
	if l.Get(rowEnd, colEnd) != common.Empty {
		fmt.Printf("WARNING: assigning '%c' to non-empty sister cell %d, %d = '%c'\n", r, row, col, l.Get(rowEnd, colEnd))
	}
	l.Set(rowEnd, colEnd, common.Negate(r))
}

// GetFrame returns the rune at this coordinate from the frames board.
func (game *Game) GetFrame(row, col int) rune {
	return game.frames.Get(row, col)
}

// GetFrameEnd returns the coordinates of the other end of the frame. If the cell is a wall (not a frame) it returns -1, -1.
func (game *Game) GetFrameEnd(row, col int) (int, int) {
	r := 0
	c := 0
	switch game.frames.Get(row, col) {
	case common.Up:
		r = 1
		c = 0
	case common.Down:
		r = -1
		c = 0
	case common.Left:
		r = 0
		c = 1
	case common.Right:
		r = 0
		c = -1
	case common.Wall:
		return -1, -1
	default:
		panic("asdf")
		fmt.Printf("ERROR: frame at %d, %d unexpected '%c'\n", row, col, game.frames.Get(row, col))
	}

	return row + r, col + c
}

// Frames iterates over every frame in the layer, pushing the row/col to a channel.
// It uses a filter, so if you change the contents of the frame layer while it is
// iterating it may return inconsistent results.
func (game *Game) Frames() <-chan board.Coord {
	return game.frames.Cells(common.Up, common.Left)
}

// Solved checks to see if the guess board has a valid solution.
func (game *Game) Solved() bool {
	for row := 0; row < game.Guess.Height(); row++ {
		if game.Guess.CountRow(row, common.Positive) != game.rowPos[row] {
			return false
		}
		if game.Guess.CountRow(row, common.Negative) != game.rowNeg[row] {
			return false
		}
		if game.Guess.CountRow(row, common.Neutral)+game.Guess.CountRow(row, common.Wall) != game.Guess.Width()-(game.rowPos[row]+game.rowNeg[row]) {
			return false
		}
	}

	for col := 0; col < game.Guess.Width(); col++ {
		if game.Guess.CountCol(col, common.Positive) != game.colPos[col] {
			return false
		}
		if game.Guess.CountCol(col, common.Negative) != game.colNeg[col] {
			return false
		}
		if game.Guess.CountCol(col, common.Neutral)+game.Guess.CountCol(col, common.Wall) != game.Guess.Height()-(game.colPos[col]+game.colNeg[col]) {
			return false
		}
	}

	// Validate that there are no two identical signs next to each other.
	for cell := range game.Guess.Cells(common.Positive, common.Negative) {
		row, col := cell.Unpack()
		grid := game.Guess.Get(row, col)
		for _, adj := range board.Adjacents {
			r, c := adj.Unpack()
			if game.Guess.Get(row+r, col+c) == grid {
				return false
			}
		}
	}

	return true
}

// CountRow counts the number of occurrences of the given rune in a row.
func (game *Game) CountRow(row int, r rune) int {
	if r == common.Positive {
		return game.rowPos[row]
	}
	if r == common.Negative {
		return game.rowNeg[row]
	}
	if r == common.Neutral {
		return game.grid.Width() - (game.rowPos[row] + game.rowNeg[row])
	}
	return game.grid.CountRow(row, r)
}

// CountCol counts the number of occurrences of the given rune in a column.
func (game *Game) CountCol(col int, r rune) int {
	if r == common.Positive {
		return game.colPos[col]
	}
	if r == common.Negative {
		return game.colNeg[col]
	}
	if r == common.Neutral {
		return game.grid.Height() - (game.colPos[col] + game.colNeg[col])
	}
	return game.grid.CountCol(col, r)
}
