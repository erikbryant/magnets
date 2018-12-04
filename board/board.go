package board

import (
	"../common"
	"fmt"
)

type Board struct {
	width  int
	height int
	cells  [][]rune
}

type Coord struct {
	Row int
	Col int
}

var (
	Adjacents []Coord = []Coord{{Row: -1, Col: 0}, {Row: 0, Col: -1}, {Row: 0, Col: +1}, {Row: +1, Col: 0}}
)

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

// Unpack() returns the row/col values from the struct.
func (c *Coord) Unpack() (int, int) {
	return c.Row, c.Col
}

// Cells() iterates over every cell in the layer, pushing the row/col to a channel.
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
					if l.Get(row, col) == val {
						c <- Coord{Row: row, Col: col}
					}
				}
			}
		}
	}()
	return c
}

// Width() returns the width of the layer.
func (l *Board) Width() int {
	return l.width
}

// Height() returns the height of the layer.
func (l *Board) Height() int {
	return l.height
}

// CountRow() returns the number of cells of type 'r' that are in the row.
func (l *Board) CountRow(row int, r rune) int {
	count := 0
	for col := 0; col < l.width; col++ {
		if l.cells[row][col] == r {
			count++
		}
	}
	return count
}

// CountCol() returns the number of cells of type 'r' that are in the col.
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

func (l *Board) Get(row, col int) rune {
	if row < 0 || row >= l.height || col < 0 || col >= l.width {
		return common.Border
	}
	return l.cells[row][col]
}

func (l *Board) Set(row, col int, r rune) {
	if row < 0 || row >= l.height || col < 0 || col >= l.width {
		return
	}
	l.cells[row][col] = r
}

func (l *Board) FloodFill() {
	changed := true

	for changed {
		changed = false
		for cell := range l.Cells(common.Positive, common.Negative) {
			row, col := cell.Unpack()
			grid := l.Get(row, col)
			for _, mod := range Adjacents {
				r := row + mod.Row
				c := col + mod.Col
				if l.Get(r, c) == common.Empty {
					l.Set(r, c, common.Negate(grid))
					changed = true
				}
			}
		}
	}
}

func (l *Board) Equal(l2 Board) bool {
	if l.Height() != l2.Height() || l.Width() != l2.Width() {
		return false
	}

	for cell := range l.Cells() {
		row, col := cell.Unpack()
		if l.Get(row, col) != l2.Get(row, col) {
			return false
		}
	}

	return true
}

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
