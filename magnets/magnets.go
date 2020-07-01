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
	"math/rand"
	"strconv"
	"strings"
	"time"
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

// init sets the random seed.
func init() {
	rand.Seed(time.Now().UnixNano())
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
		fmt.Printf("WARNING: assigning '%c' to non-empty cell %d, %d = '%c'\n", r, row, col, l.Get(row, col))
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

// setFrameMagnet sets the polarities of a given frame to follow its neighbors.
// If there are no neighbors, use a random sign.
func (game *Game) setFrameMagnet(row, col int) {
	// Choose a random sign for the new frame.
	choices := []rune{common.Positive, common.Negative}
	sign := choices[rand.Intn(len(choices))]

	// If there is a neighbor that is already set, follow its polarity instead.
	for _, mod := range board.Adjacents {
		modR, modC := mod.Unpack()
		cell := game.grid.Get(row+modR, col+modC)
		switch cell {
		case common.Positive:
			sign = common.Negate(cell)
		case common.Negative:
			sign = common.Negate(cell)
		}
	}

	game.SetDomino(game.grid, row, col, sign)

	game.grid.FloodFill()
}

// placeFrames attempts to fill a given board with frames. It keeps
// trying until it gets a valid solution.
func (game *Game) placeFrames() {
	// This algorithm may sometimes generate an invalid
	// board frame. Loop until it generates a valid one.
	for {
		for cell := range game.frames.Cells() {
			row, col := cell.Unpack()
			if game.frames.Get(row, col) != common.Empty {
				continue
			}

			orient := common.Empty

			if game.frames.Get(row, col+1) == common.Empty {
				if game.frames.Get(row+1, col) == common.Empty {
					choices := []rune{common.Right, common.Down}
					orient = choices[rand.Intn(len(choices))]
				} else {
					orient = common.Right
				}
			} else {
				if game.frames.Get(row+1, col) == common.Empty {
					orient = common.Down
				}
			}

			// Horizontal
			if orient == common.Right {
				game.frames.Set(row, col, common.Left)
				game.frames.Set(row, col+1, common.Right)
			}
			// Vertical
			if orient == common.Down {
				game.frames.Set(row, col, common.Up)
				game.frames.Set(row+1, col, common.Down)
			}
		}

		// If this is an odd-sized board, there will be (at least)
		// one blank space. Replace it with a wall.
		if (game.frames.Width()*game.frames.Height())%2 == 1 {
			for cell := range game.frames.Cells() {
				row, col := cell.Unpack()
				if game.frames.Get(row, col) == common.Empty {
					game.frames.Set(row, col, common.Wall)
					game.grid.Set(row, col, common.Wall)
					game.Guess.Set(row, col, common.Wall)
				}
			}
		}

		// Is this board valid? If so, ship it! :-)
		if game.Valid() {
			break
		}

		// Reset the layers and try again.
		for cell := range game.frames.Cells() {
			row, col := cell.Unpack()
			game.frames.Set(row, col, common.Empty)
			game.grid.Set(row, col, common.Empty)
		}
	}
}

// placePieces puts the neutrals and magnets randomly on the board.
func (game *Game) placePieces() {
	// Place all of the neutrals before placing any magnets.
	// The neutrals have a chance to form walls that bound
	// disconnected areas. Placing the magnets calls flood
	// fill and would simply take up all the rest of the
	// board.
	for frame := range game.Frames() {
		// Random chance to add a neutral.
		if rand.Intn(10) != 0 {
			continue
		}
		row, col := frame.Unpack()
		if game.grid.Get(row, col) != common.Empty {
			fmt.Printf("ERROR: Frame at %d, %d was not empty\n", row, col)
			continue
		}
		game.SetDomino(game.grid, row, col, common.Neutral)
	}

	// Place the magnets in the remaining frames.
	for frame := range game.Frames() {
		row, col := frame.Unpack()
		if game.grid.Get(row, col) != common.Empty {
			continue
		}
		game.setFrameMagnet(row, col)
	}

	// Record how many magnets are in each row/col
	for row := 0; row < game.grid.Height(); row++ {
		game.rowPos[row] = game.grid.CountRow(row, common.Positive)
		game.rowNeg[row] = game.grid.CountRow(row, common.Negative)
	}
	for col := 0; col < game.grid.Width(); col++ {
		game.colPos[col] = game.grid.CountCol(col, common.Positive)
		game.colNeg[col] = game.grid.CountCol(col, common.Negative)
	}
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
	}

	for col := 0; col < game.Guess.Width(); col++ {
		if game.Guess.CountCol(col, common.Positive) != game.colPos[col] {
			return false
		}
		if game.Guess.CountCol(col, common.Negative) != game.colNeg[col] {
			return false
		}
	}

	return true
}

// makeGame creates an empty game state.
func makeGame(width, height int) Game {
	var game Game

	game.grid = board.New(width, height)
	game.frames = board.New(width, height)
	game.Guess = board.New(width, height)

	game.rowPos = make([]int, height)
	game.rowNeg = make([]int, height)
	game.colPos = make([]int, width)
	game.colNeg = make([]int, width)

	return game
}

// New creates and populates all of the layers that make up a game.
func New(width, height int) Game {
	game := makeGame(width, height)
	game.placeFrames()
	game.placePieces()

	if !game.Valid() {
		fmt.Println("ERROR: New() board is not valid.")
		game.Print()
	}

	return game
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

// CountCol counts tthe number of occurrences of the given rune in a column.
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

// Saving and loading (serializing and de-serializing) of a game.
//
// [Comment taken from Simon's magnets.c file (with corrections
//  based on the actual behavior of the iPhone version)]
//
// Puzzle definition is just the size, and then the list of + (across then
// down) and - (across then down) present, then domino edges.
//
// An example:
//
//  + 2 0 1
//   +-----+
//  1|+ -| |1
//   |-+-+ |
//  0|-|#| |1
//   | +-+-|
//  2|+|- +|1
//   +-----+
//    1 2 0 -
//
// 3x3:201,102,120,111,LRTT*BBLR
//
// 'Zotmeister' example:
// 5x5:.2..1,3..1.,.2..2,2..2.,LRLRTTLRTBBT*BTTBLRBBLRLR
//
// Janko 6x6 with solution:
// 6x6:322223,323132,232223,232223,LRTLRTTTBLRBBBTTLRLRBBLRTTLRTTBBLRBB
//
// iPhone 2x50 (note the use of lowercase letters for column counts):
// 2x50:jh,01111110101001111000011110111110011111011101111111,me,01111101101001110100101110111101011111011011111111,LRTTBBLRLRLRTTBBLRLRLRTTBBTTBBLRTTBBLRLRTTBBLRTTBBLRLRLRTTBBTTBBLRLRTTBBTTBBLRTTBBTTBBTTBBTTBBLRTTBB
//
// For a game of size w*h the game description is:
//
// wxh
// colon
// w-sized string of column positive count (L-R), or '.' for none
// comma
// h-sized string of row positive count (T-B), or '.'
// comma
// w-sized string of column negative count (L-R), or '.'
// comma
// h-sized string of row negative count (T-B), or '.'
// comma
// w*h-sized string of 'L', 'R', 'T', 'B' for domino associations,
//   or '*' for a black singleton square.
//
// There is only one character position allocated for the count. So, if a
// count is greater than 9 it rolls to alpha characters. First lowercase,
// then uppercase.

// countToRune returns the base-62 form of an int. Valid input is 0-61.
func countToRune(count int) rune {
	if count < 0 {
		return '-'
	}

	// 0-9
	if count <= 9 {
		return rune('0' + count)
	}

	// a-z
	if count <= 9+26 {
		return rune('a' + (count - 10))
	}

	// A-Z
	if count <= 9+26+26 {
		return rune('A' + (count - 10 - 26))
	}

	// Overflow!
	return '!'
}

// runeToCount returns the base 10 form of a base-62. Valid input is 0-9, a-z, and A-Z.
func runeToCount(r rune) int {
	// Underflow!
	if r < '0' {
		return -1
	}

	// 0-9
	if r <= '9' {
		return int(r - '0')
	}

	// A-Z
	if r <= 'Z' {
		return int(r - 'A' + 10 + 26)
	}

	// a-z
	if r <= 'z' {
		return int(r - 'a' + 10)
	}

	// Overflow!
	return -1
}

// Serialize returns a representation of the game in string form.
func (game *Game) Serialize() (string, bool) {
	if !game.Valid() {
		return "", false
	}

	var serial string
	valid := true

	// WxH:
	serial += fmt.Sprintf("%dx%d:", game.grid.Width(), game.grid.Height())

	// Col positive count
	for col := 0; col < game.grid.Width(); col++ {
		count := game.grid.CountCol(col, common.Positive)
		serial += string(countToRune(count))
	}
	serial += ","

	// Row positive count
	for row := 0; row < game.grid.Height(); row++ {
		count := game.grid.CountRow(row, common.Positive)
		serial += string(countToRune(count))
	}
	serial += ","

	// Col negative count
	for col := 0; col < game.grid.Width(); col++ {
		count := game.grid.CountCol(col, common.Negative)
		serial += string(countToRune(count))
	}
	serial += ","

	// Row negative count
	for row := 0; row < game.grid.Height(); row++ {
		count := game.grid.CountRow(row, common.Negative)
		serial += string(countToRune(count))
	}
	serial += ","

	// LRTB*
	for cell := range game.frames.Cells() {
		switch game.frames.Get(cell.Unpack()) {
		case common.Left:
			serial += "L"
		case common.Right:
			serial += "R"
		case common.Up:
			serial += "T"
		case common.Down:
			serial += "B"
		case common.Wall:
			serial += "*"
		default:
			serial += "!"
			valid = false
		}
	}

	return serial, valid
}

// Deserialize takes a serial representation of a game and unpacks it,
// returning a game and whether or not the unpacking was successful.
func Deserialize(s string) (Game, bool) {
	valid := true

	xPos := strings.IndexRune(s, 'x')
	colonPos := strings.IndexRune(s, ':')

	width, _ := strconv.Atoi(s[0:xPos])
	height, _ := strconv.Atoi(s[xPos+1 : colonPos])
	game := makeGame(width, height)
	game.serial = s
	s = s[colonPos+1:]

	// Col positive count
	commaPos := strings.IndexRune(s, ',')
	for i, r := range s[:commaPos] {
		game.colPos[i] = runeToCount(r)
	}
	s = s[commaPos+1:]

	// Row positive count
	commaPos = strings.IndexRune(s, ',')
	for i, r := range s[:commaPos] {
		game.rowPos[i] = runeToCount(r)
	}
	s = s[commaPos+1:]

	// Col negative count
	commaPos = strings.IndexRune(s, ',')
	for i, r := range s[:commaPos] {
		game.colNeg[i] = runeToCount(r)
	}
	s = s[commaPos+1:]

	// Row negative count
	commaPos = strings.IndexRune(s, ',')
	for i, r := range s[:commaPos] {
		game.rowNeg[i] = runeToCount(r)
	}
	s = s[commaPos+1:]

	// Place frames
	row := 0
	col := 0
	for _, cell := range s {
		switch cell {
		case 'L':
			game.frames.Set(row, col, common.Left)
		case 'R':
			game.frames.Set(row, col, common.Right)
		case 'T':
			game.frames.Set(row, col, common.Up)
		case 'B':
			game.frames.Set(row, col, common.Down)
		case '*':
			game.frames.Set(row, col, common.Wall)
			game.grid.Set(row, col, common.Wall)
			game.Guess.Set(row, col, common.Wall)
		case '!':
			game.frames.Set(row, col, common.Empty)
			valid = false
		default:
			game.frames.Set(row, col, common.Empty)
			valid = false
		}
		col++
		if col >= width {
			col = 0
			row++
		}
	}

	return game, valid
}
