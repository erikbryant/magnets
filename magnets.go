package main

import (
	"fmt"
	"github.com/erikbryant/magnets/magnets"
	"github.com/erikbryant/magnets/solver"
	"math/rand"
)

// stressTest loops forever, creating random boards and trying to solve them. At intervals it prints success/fail statistics.
func stressTest() {
	games := 0
	solved := 0

	for {
		gamez := magnets.New(rand.Intn(15)+2, rand.Intn(15)+2)

		// --- cut here ---

		// TODO: There is something wrong with board creation. magnets.New()
		// sometimes (like 9 in 100,000 calls) returns data that is not
		// consistent and causes the solver to go into an infinite loop.
		//
		// Debugging of the infinite loop states show that the solver
		// has put invalid magnet orientations and has moved where the
		// neutrals are supposed to be.
		//
		// Serializing and deserializing the new board before calling the
		// solver alleviates the problem.

		if !gamez.Valid() {
			fmt.Println("Game is not valid")
			gamez.Print()
			return
		}

		s, ok := gamez.Serialize()
		if !ok {
			fmt.Println("Unable to serialize")
			gamez.Print()
		}
		game, ok2 := magnets.Deserialize(s)
		if !ok2 {
			fmt.Println("Unable to deserialize:", s)
		}

		// --- cut here ---

		games++
		solver.Solve(game)
		if game.Solved() {
			solved++
		}
		if games%10000 == 0 {
			pctSolved := int(float64(solved) / float64(games) * 100.0)
			fmt.Printf("Played: %d Solved: %d (%d%%)\n", games, solved, pctSolved)
		}
	}
}

// deserializer takes a game in serial form and tries to solve it.
func deserializer(s string) {
	game, ok := magnets.Deserialize(s)
	if !ok {
		fmt.Println("Could not deserialize!", s)
		return
	}

	solver.Solve(game)
	if game.Solved() {
		fmt.Println("Solved!")
	} else {
		fmt.Println("Could not solve:", s)
		game.Print()
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
		fmt.Printf("\"%s\",\n", s)
		break
	}
}

func main() {
	// stressTest()

	i := 3
	for {
		solvable(i, i)
		i++
	}
}
