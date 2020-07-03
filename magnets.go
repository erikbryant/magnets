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

	fmt.Println("Solutions:", game.CountSolutions(0, 0))

	solver.Solve(game)
	if game.Solved() {
		fmt.Println("Solved!", s)
		game.Print()
	} else {
		fmt.Println("Could not solve!", s)
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
		}

		solutions := game.CountSolutions(0, 0)
		if solutions == 1 {
			append("test_success", testCase)
		} else {
			append("test_fail", testCase)
		}

		// solver.Solve(game)
		// if game.Solved() != expected {
		// 	fmt.Println("ERROR: For %s expected solved to be %t", testCase, expected)
		// }
	}
}

func testSolve() {
	helper("./solver/test_solve.txt", true)
	// helper("./solver/test_solve_fail.txt", false)
}

func main() {
	start := time.Now()

	// createCorpus()
	// stressTest()
	testSolve()

	// deserializer("1x2:1,10,1,01,TB")
	// deserializer("1x3:1,100,1,010,TB*")
	// deserializer("2x2:11,11,11,11,TTBB")
	// deserializer("6x6:233322,123333,323232,033333,LRTLRTTTBTTBBBTBBTTTBTTBBBTBBTLRBLRB")
	// deserializer("16x3:2121212121212120,878,1212121212121211,887,TTTLRTLRLRTTTTLRBBBLRBTTLRBBBBTTLRLRLRBBLRLRLRBB")

	fmt.Println("testSolve took", time.Since(start))
}
