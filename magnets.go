package main

import (
	"fmt"
	"github.com/erikbryant/magnets/magnets"
	"github.com/erikbryant/magnets/solver"
	"math/rand"
	"os"
)

// createCorups creates games and tries to solve them. The ones it can solve it writes to
// one file and the ones it cannot solve it writes to another file.
func createCorpus() {
	s, err := os.Create("solved")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer s.Close()

	u, err := os.Create("unsolved")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer u.Close()

	games := 0

	for solved := 0; solved < 1000000; {
		game := magnets.New(rand.Intn(15)+2, rand.Intn(15)+2)
		games++

		solver.Solve(game)

		serial, ok := game.Serialize()
		if !ok {
			fmt.Println("Could not serialize game")
			game.Print()
			return
		}

		serial = fmt.Sprintf("%s\n", serial)

		if game.Solved() {
			solved++
			_, err = s.WriteString(serial)
		} else {
			_, err = u.WriteString(serial)
		}
		if err != nil {
			fmt.Println(err)
			return
		}

		if games%10000 == 0 {
			pctSolved := 100.0 * float64(solved) / float64(games)
			fmt.Printf("Played: %d Solved: %d (%.3f%%)\n", games, solved, pctSolved)
		}
	}
}

// stressTest loops forever, creating random boards and trying to solve them. At intervals
// it prints success/fail statistics.
func stressTest() {
	games := 0
	solved := 0

	for {
		game := magnets.New(rand.Intn(15)+2, rand.Intn(15)+2)
		games++

		solver.Solve(game)

		if game.Solved() {
			solved++
		}

		if games%10000 == 0 {
			pctSolved := 100.0 * float64(solved) / float64(games)
			fmt.Printf("Played: %d Solved: %d (%.3f%%)\n", games, solved, pctSolved)
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

// solvable loops forever trying random boards until it can solve one.
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
	// createCorpus()

	stressTest()
}
