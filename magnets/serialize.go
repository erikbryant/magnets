package magnets

import (
	"fmt"
	"github.com/erikbryant/magnets/common"
	"strconv"
	"strings"
)

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
		count := game.CountCol(col, common.Positive)
		serial += string(countToRune(count))
	}
	serial += ","

	// Row positive count
	for row := 0; row < game.grid.Height(); row++ {
		count := game.CountRow(row, common.Positive)
		serial += string(countToRune(count))
	}
	serial += ","

	// Col negative count
	for col := 0; col < game.grid.Width(); col++ {
		count := game.CountCol(col, common.Negative)
		serial += string(countToRune(count))
	}
	serial += ","

	// Row negative count
	for row := 0; row < game.grid.Height(); row++ {
		count := game.CountRow(row, common.Negative)
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
	xPos := strings.IndexRune(s, 'x')
	if xPos == -1 {
		return makeGame(0, 0), false
	}

	colonPos := strings.IndexRune(s, ':')
	if colonPos == -1 {
		return makeGame(0, 0), false
	}

	width, _ := strconv.Atoi(s[0:xPos])
	height, _ := strconv.Atoi(s[xPos+1 : colonPos])
	game := makeGame(width, height)
	game.serial = s
	s = s[colonPos+1:]

	// Col positive count
	for i, r := range s[:width] {
		game.colPos[i] = runeToCount(r)
	}
	s = s[width:]
	if s[0] != ',' {
		return game, false
	}
	s = s[1:]

	// Row positive count
	for i, r := range s[:height] {
		game.rowPos[i] = runeToCount(r)
	}
	s = s[height:]
	if s[0] != ',' {
		return game, false
	}
	s = s[1:]

	// Col negative count
	for i, r := range s[:width] {
		game.colNeg[i] = runeToCount(r)
	}
	s = s[width:]
	if s[0] != ',' {
		return game, false
	}
	s = s[1:]

	// Row negative count
	for i, r := range s[:height] {
		game.rowNeg[i] = runeToCount(r)
	}
	s = s[height:]
	if s[0] != ',' {
		return game, false
	}
	s = s[1:]

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
			return game, false
		default:
			game.frames.Set(row, col, common.Empty)
			return game, false
		}
		col++
		if col >= width {
			col = 0
			row++
		}
	}

	return game, true
}

// Print prints an ASCII representation of the board.
func (game *Game) Print() {
	fmt.Printf("\n")
	s, ok := game.Serialize()
	if !ok {
		s = "ERROR: Unable to serialize game!"
	}
	fmt.Println(s)
	game.frames.Print("Frames", game.rowPos, game.rowNeg, game.colPos, game.colNeg)
	game.grid.Print("Grid", game.rowPos, game.rowNeg, game.colPos, game.colNeg)
	game.Guess.Print("Guess", []int{}, []int{}, []int{}, []int{})
}
