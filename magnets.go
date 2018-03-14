package main

import (
	"./magnets"
	"./solver"
	"fmt"
	"math/rand"
)

func stressTest() {
	games := 0
	solved := 0

	for {
		game := magnets.New(rand.Intn(20)+1, rand.Intn(20)+2)
		games++
		solver.Solve(game)
		if game.Solved() {
			// game.Print()
			solved++
		}
		if games%100 == 0 {
			pctSolved := int(float64(solved) / float64(games) * 100.0)
			fmt.Printf("Played: %d Solved: %d (%d%%)\n", games, solved, pctSolved)
		}
	}
}

func solvable() {
	for {
		game := magnets.New(3, 3)
		solver.Solve(game)
		if game.Solved() {
			// game.Print()
			s, ok := game.Serialize()
			fmt.Println(s)
			if !ok {
				fmt.Println("FAIL!")
				return
			}
		}
	}
}

func deserializer() {
	// s := "3x3:201,102,120,111,LRTT*BBLR"
	// s := "2x50:jh,01111110101001111000011110111110011111011101111111,me,01111101101001110100101110111101011111011011111111,LRTTBBLRLRLRTTBBLRLRLRTTBBTTBBLRTTBBLRLRTTBBLRTTBBLRLRLRTTBBTTBBLRLRTTBBTTBBLRTTBBTTBBTTBBTTBBLRTTBB"
	// s := "3x3:112,112,121,121,*LRTLRBLR"
	s := "3x3:111,021,111,111,LRTTTBBB*"

	game, ok := magnets.Deserialize(s)
	game.Print()
	if !ok {
		fmt.Printf("Game is not valid!")
		return
	}
	solver.Solve(game)
	game.Print()
}

func main() {
	solvable()
	// deserializer()
	// stressTest()
}
