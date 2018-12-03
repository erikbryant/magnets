package main

import (
	"./magnets"
	"./solver"
	"fmt"
	"math/rand"
)

// stressTest loops forever, creating random boards and trying to solve them. At intervals it will print its success/fail statistics.
func stressTest() {
	games := 0
	solved := 0

	for {
		game := magnets.New(rand.Intn(6)+1, rand.Intn(6)+2)
		games++
		solver.Solve(game)
		if game.Solved() {
			solved++
		}
		if games%100 == 0 {
			pctSolved := int(float64(solved) / float64(games) * 100.0)
			fmt.Printf("Played: %d Solved: %d (%d%%)\n", games, solved, pctSolved)
		}
	}
}

// solvable loops forever trying to solve random boards. It will print any it can solve.
func solvable(width, height int) {
	for {
		game := magnets.New(width, height)
		solver.Solve(game)
		if !game.Solved() {
			continue
		}
		s, _ := game.Serialize()
		fmt.Println(s)
	}
}

// deserializer takes a game in serial form and tries to solve it.
func deserializer(s string) {
	// s = "3x3:201,102,120,111,LRTT*BBLR"
	// s = "2x50:jh,01111110101001111000011110111110011111011101111111,me,01111101101001110100101110111101011111011011111111,LRTTBBLRLRLRTTBBLRLRLRTTBBTTBBLRTTBBLRLRTTBBLRTTBBLRLRLRTTBBTTBBLRLRTTBBTTBBLRTTBBTTBBTTBBTTBBLRTTBB"
	// s = "3x3:112,112,121,121,*LRTLRBLR"
	s = "3x3:111,021,111,111,LRTTTBBB*"

	game, ok := magnets.Deserialize(s)
	if !ok {
		fmt.Printf("Could not deserialized!", s)
		return
	}
	game.Print()
	solver.Solve(game)
	game.Print()
}

func main() {
	stressTest()
	// solvable(7, 7)
	// deserializer("3x3:112,112,121,121,*LRTLRBLR")
}
