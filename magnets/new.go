package magnets

import (
	"fmt"
	"github.com/erikbryant/magnets/board"
	"github.com/erikbryant/magnets/common"
	"math/rand"
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
		if game.Valid() && game.singleSolution() {
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
