package common

import (
	"testing"
)

func TestNegate(t *testing.T) {
	testCases := []struct {
		n        rune
		expected rune
	}{
		{Positive, Negative},
		{Negative, Positive},
		{Up, Down},
		{Down, Up},
		{Left, Right},
		{Right, Left},
		{Neutral, Neutral},
	}

	for _, testCase := range testCases {
		answer := Negate(testCase.n)
		if answer != testCase.expected {
			t.Errorf("ERROR: For '%c' expected '%c', got '%c'", testCase.n, testCase.expected, answer)
		}
	}
}
