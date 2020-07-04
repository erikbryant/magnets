package main

import (
	"bufio"
	"fmt"
	"github.com/erikbryant/magnets/magnets"
	"github.com/erikbryant/magnets/solver"
	"math/rand"
	"os"
	"strings"
	"time"
)

func append(file, content string) {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer f.Close()

	_, err = f.WriteString(content + "\n")
	if err != nil {
		fmt.Println(err)
		return
	}
}

// createCorups creates games and tries to solve them. The ones it can solve it writes to
// one file and the ones it cannot solve it writes to another file.
func createCorpus() {
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

		if game.Solved() {
			solved++
			append("solved", serial)
		} else {
			append("unsolved", serial)
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

	solutions := game.CountSolutions(0, 0)
	if solutions != 1 {
		fmt.Println("Solutions:", solutions)
		return
	}

	solver.Solve(game)

	if game.Solved() {
		fmt.Println("Solved!", s)
	} else {
		fmt.Println("Could not solve!", s)
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

// helper runs solver tests against a given file.
func helper(file string, expected bool) {
	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("Unable to open testcases %s %s\n", file, err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		testCase := scanner.Text()

		testCase = strings.TrimSpace(testCase)

		if len(testCase) == 0 {
			continue
		}

		if strings.HasPrefix(testCase, "//") {
			continue
		}

		game, ok := magnets.Deserialize(testCase)
		if !ok {
			fmt.Printf("ERROR: Unable to deserialize %s\n", testCase)
			continue
		}

		solver.Solve(game)

		if game.Solved() != expected {
			fmt.Printf("ERROR: For %s expected solved to be %t\n", testCase, expected)
			append("error", testCase)
		}
	}
}

func testSolve() {
	helper("./solver/test_solve.txt", true)
	helper("./solver/test_solve_fail.txt", false)
}

func main() {
	start := time.Now()

	deserializer("2x3:12,111,21,111,LRLRLR")

	// createCorpus()
	// stressTest()
	testSolve()

	fmt.Println("Elapsed time:", time.Since(start))
}
