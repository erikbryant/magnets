package board

import (
	"fmt"

	"github.com/erikbryant/magnets/common"
)

// Board implements a widthxheight grid of runes.
type Board struct {
	width  int
	height int
	cells  [][]rune
}

// Coord represents a single row/col address.
type Coord struct {
	Row int
	Col int
}

var (
	// Adjacents contains the offsets for each orthogonal neigbor on a grid.
	Adjacents = []Coord{
		{Row: -1, Col: 0},
		{Row: 0, Col: -1},
		{Row: 0, Col: +1},
		{Row: +1, Col: 0},
	}
)

// New creates a new board, populated with empty squares.
func New(width, height int) Board {
	var l Board

	l.width = width
	l.height = height
	l.cells = make([][]rune, height)
	for row := 0; row < height; row++ {
		l.cells[row] = make([]rune, width)
	}

	for cell := range l.Cells() {
		row, col := cell.Unpack()
		l.cells[row][col] = common.Empty
	}

	return l
}

// Unpack returns the row/col values from the struct.
func (c *Coord) Unpack() (int, int) {
	return c.Row, c.Col
}

// Cells iterates over every cell in the layer, pushing the row/col to a channel.
// If you provide a filter it will return only those cells matching those values.
// Be careful with the filter as you can get non-deterministic behavior if you change
// values in the layer while iterating over the layer.
func (l *Board) Cells(r ...rune) <-chan Coord {
	c := make(chan Coord, l.Width())

	go func() {
		defer close(c)
		for row := 0; row < l.Height(); row++ {
			for col := 0; col < l.Width(); col++ {
				if r == nil {
					c <- Coord{Row: row, Col: col}
					continue
				}
				for _, val := range r {
					if l.Get(row, col, false) == val {
						c <- Coord{Row: row, Col: col}
					}
				}
			}
		}
	}()

	return c
}

// Width returns the width of the layer.
func (l *Board) Width() int {
	return l.width
}

// Height returns the height of the layer.
func (l *Board) Height() int {
	return l.height
}

// CountRow returns the number of cells of type 'r' that are in the row.
func (l *Board) CountRow(row int, r rune) int {
	count := 0
	for col := 0; col < l.width; col++ {
		if l.Get(row, col, false) == r {
			count++
		}
	}
	return count
}

// CountCol returns the number of cells of type 'r' that are in the col.
// TODO: merge this with CountRow() and just flip the matrix.
func (l *Board) CountCol(col int, r rune) int {
	count := 0
	for row := 0; row < l.height; row++ {
		if l.cells[row][col] == r {
			count++
		}
	}
	return count
}

// Get returns the value at row,col unless flip is true in which case it
// returns the value at col,row.
func (l *Board) Get(row, col int, flip bool) rune {
	if flip {
		row, col = col, row
	}

	if row < 0 || row >= l.height || col < 0 || col >= l.width {
		return common.Border
	}
	return l.cells[row][col]
}

// Set sets the value at row,col unless flip is true in which case it
// sets the value at col,row.
func (l *Board) Set(row, col int, r rune, flip bool) {
	if flip {
		row, col = col, row
	}

	if row < 0 || row >= l.height || col < 0 || col >= l.width {
		return
	}
	l.cells[row][col] = r
}

// FloodFill fills in a board (or bounded region on a board).
func (l *Board) FloodFill() {
	changed := true

	for changed {
		changed = false
		for cell := range l.Cells(common.Positive, common.Negative) {
			row, col := cell.Unpack()
			grid := l.Get(row, col, false)
			for _, mod := range Adjacents {
				r := row + mod.Row
				c := col + mod.Col
				if l.Get(r, c, false) == common.Empty {
					l.Set(r, c, common.Negate(grid), false)
					changed = true
				}
			}
		}
	}
}

// Equal returns true if the two boards are identical, false otherwise.
func (l *Board) Equal(l2 Board) bool {
	if l.Height() != l2.Height() || l.Width() != l2.Width() {
		return false
	}

	for cell := range l.Cells() {
		row, col := cell.Unpack()
		if l.Get(row, col, false) != l2.Get(row, col, false) {
			return false
		}
	}

	return true
}

// Print prints a representation of the board state to the console.
func (l *Board) Print(name string, rowPos, rowNeg, colPos, colNeg []int) {
	fmt.Printf("%s (%dx%d)\n", name, l.width, l.height)

	fmt.Printf("     ")
	for i := 0; i < l.width; i++ {
		count := 0
		if len(colPos) > 0 {
			count = colPos[i]
		} else {
			count = l.CountCol(i, common.Positive)
		}
		fmt.Printf("%1d", count)
	}
	fmt.Printf("\n")

	fmt.Printf("   + ")
	for i := 0; i < l.width; i++ {
		fmt.Printf("―")
	}
	fmt.Printf("\n")

	for row := 0; row < l.height; row++ {
		count := 0
		if len(rowPos) > 0 {
			count = rowPos[row]
		} else {
			count = l.CountRow(row, common.Positive)
		}
		fmt.Printf("%2d | ", count)
		for _, cell := range l.cells[row] {
			fmt.Printf("%c", cell)
		}
		count = 0
		if len(rowNeg) > 0 {
			count = rowNeg[row]
		} else {
			count = l.CountRow(row, common.Negative)
		}
		fmt.Printf(" | %2d\n", count)
	}

	fmt.Printf("     ")
	for i := 0; i < l.width; i++ {
		fmt.Printf("―")
	}
	fmt.Printf(" -\n")

	fmt.Printf("     ")
	for i := 0; i < l.width; i++ {
		count := 0
		if len(colNeg) > 0 {
			count = colNeg[i]
		} else {
			count = l.CountCol(i, common.Negative)
		}
		fmt.Printf("%1d", count)
	}
	fmt.Printf("\n")

	fmt.Printf("\n")
}
