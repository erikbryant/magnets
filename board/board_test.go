package board

import (
	"github.com/erikbryant/magnets/common"
	"testing"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		width  int
		height int
	}{
		{1, 1},
		{5, 3},
		{20, 20},
	}

	for _, testCase := range testCases {
		answer := New(testCase.width, testCase.height)
		h := len(answer.cells)
		if h != testCase.height {
			t.Errorf("ERROR: For %d, %d got height: %d", testCase.width, testCase.height, h)
		}
		for i := 0; i < h; i++ {
			w := len(answer.cells[0])
			if w != testCase.width {
				t.Errorf("ERROR: For %d, %d got width: %d", testCase.width, testCase.height, w)
			}
		}
	}
}

func TestUnpack(t *testing.T) {
	testCases := []struct {
		row int
		col int
	}{
		{1, 2},
		{-1, 3},
		{12, 8},
	}

	for _, testCase := range testCases {
		coord := Coord{Row: testCase.row, Col: testCase.col}
		answerR, answerC := coord.Unpack()
		if answerR != testCase.row {
			t.Errorf("ERROR: For %d, %d expected %d got %d", testCase.row, testCase.col, testCase.row, answerR)
		}
		if answerC != testCase.col {
			t.Errorf("ERROR: For %d, %d expected %d got %d", testCase.row, testCase.col, testCase.col, answerC)
		}
	}
}

func TestCells(t *testing.T) {
	w := 4
	h := 5
	l := New(w, h)

	l.Set(0, 0, common.Positive)
	l.Set(0, 1, common.Negative)
	l.Set(0, 2, common.Neutral)
	l.Set(0, 3, common.Neutral)
	l.Set(1, 0, common.Negative)
	l.Set(1, 1, common.Positive)
	l.Set(2, 0, common.Positive)
	l.Set(3, 0, common.Negative)

	answer := 0
	for range l.Cells() {
		answer++
	}
	expected := w * h
	if answer != expected {
		t.Errorf("ERROR: Expected length %d, got length %d", expected, answer)
	}

	answer = 0
	for range l.Cells(common.Positive) {
		answer++
	}
	expected = 3
	if answer != expected {
		t.Errorf("ERROR: Expected length %d, got length %d", expected, answer)
	}

	answer = 0
	for range l.Cells(common.Negative, common.Neutral) {
		answer++
	}
	expected = 5
	if answer != expected {
		t.Errorf("ERROR: Expected length %d, got length %d", expected, answer)
	}
}

func TestWidth(t *testing.T) {
	testCases := []struct {
		width  int
		height int
	}{
		{1, 1},
		{5, 3},
		{20, 20},
	}

	for _, testCase := range testCases {
		l := New(testCase.width, testCase.height)
		answer := l.Width()
		if answer != testCase.width {
			t.Errorf("ERROR: For %d, %d got %d", testCase.width, testCase.height, answer)
		}
	}
}

func TestHeight(t *testing.T) {
	testCases := []struct {
		width  int
		height int
	}{
		{1, 1},
		{5, 3},
		{20, 20},
	}

	for _, testCase := range testCases {
		l := New(testCase.width, testCase.height)
		answer := l.Height()
		if answer != testCase.height {
			t.Errorf("ERROR: For %d, %d got %d", testCase.width, testCase.height, answer)
		}
	}
}

func TestCountRow(t *testing.T) {
	w := 10
	h := 15
	l := New(w, h)

	l.Set(0, 0, common.Positive)
	l.Set(0, 1, common.Negative)
	l.Set(0, 2, common.Neutral)
	l.Set(0, 3, common.Neutral)
	l.Set(1, 0, common.Negative)
	l.Set(1, 1, common.Positive)
	l.Set(2, 0, common.Positive)
	l.Set(3, 0, common.Negative)

	testCases := []struct {
		n         int
		expectedP int
		expectedN int
	}{
		{0, 1, 1},
		{1, 1, 1},
		{2, 1, 0},
		{3, 0, 1},
		{4, 0, 0},
	}

	for _, testCase := range testCases {
		answer := l.CountRow(testCase.n, common.Positive)
		if answer != testCase.expectedP {
			t.Errorf("ERROR: For row %d positive expected %d got %d", testCase.n, testCase.expectedP, answer)
		}
		answer = l.CountRow(testCase.n, common.Negative)
		if answer != testCase.expectedN {
			t.Errorf("ERROR: For row %d negative expected %d got %d", testCase.n, testCase.expectedN, answer)
		}
	}
}

func TestCountCol(t *testing.T) {
	w := 10
	h := 15
	l := New(w, h)

	l.Set(0, 0, common.Positive)
	l.Set(0, 1, common.Negative)
	l.Set(0, 2, common.Neutral)
	l.Set(0, 3, common.Neutral)
	l.Set(1, 0, common.Negative)
	l.Set(1, 1, common.Positive)
	l.Set(2, 0, common.Positive)
	l.Set(3, 0, common.Negative)

	testCases := []struct {
		n         int
		expectedP int
		expectedN int
	}{
		{0, 2, 2},
		{1, 1, 1},
		{2, 0, 0},
	}

	for _, testCase := range testCases {
		answer := l.CountCol(testCase.n, common.Positive)
		if answer != testCase.expectedP {
			t.Errorf("ERROR: For col %d positive expected %d got %d", testCase.n, testCase.expectedP, answer)
		}
		answer = l.CountCol(testCase.n, common.Negative)
		if answer != testCase.expectedN {
			t.Errorf("ERROR: For col %d negative expected %d got %d", testCase.n, testCase.expectedN, answer)
		}
	}
}

func TestGet(t *testing.T) {
	w := 10
	h := 15
	l := New(w, h)

	testCases := []struct {
		row      int
		col      int
		expected rune
	}{
		// Outside the top left corner
		{0, -1, common.Border},
		{-1, 0, common.Border},
		{-1, -1, common.Border},
		// Outside the bottom right corner
		{h - 1, w, common.Border},
		{h, w - 1, common.Border},
		{h, w, common.Border},
	}

	for _, testCase := range testCases {
		answer := l.Get(testCase.row, testCase.col)
		if answer != testCase.expected {
			t.Errorf("ERROR: For %d, %d expected '%c' got '%c'", testCase.row, testCase.col, testCase.expected, answer)
		}
	}
}

func TestSet(t *testing.T) {
	w := 10
	h := 15
	l := New(w, h)

	testCases := []struct {
		row      int
		col      int
		expected rune
	}{
		{0, 0, common.Empty},
		{0, w - 1, common.Positive},
		{h - 1, 0, common.Wall},
		{h - 1, w - 1, common.Neutral},
	}

	for _, testCase := range testCases {
		l.Set(testCase.row, testCase.col, testCase.expected)
		answer := l.Get(testCase.row, testCase.col)
		if answer != testCase.expected {
			t.Errorf("ERROR: For %d, %d expected '%c' got '%c'", testCase.row, testCase.col, testCase.expected, answer)
		}
	}
}

func TestFloodFill(t *testing.T) {
	l := New(2, 4)

	l.Set(0, 0, common.Positive)
	l.Set(0, 1, common.Negative)
	l.Set(1, 0, common.Empty)
	l.Set(1, 1, common.Empty)
	l.Set(2, 0, common.Neutral)
	l.Set(2, 1, common.Neutral)
	l.Set(3, 0, common.Empty)
	l.Set(3, 1, common.Empty)

	l.FloodFill()

	testCases := []struct {
		row      int
		col      int
		expected rune
	}{
		{0, 0, common.Positive},
		{0, 1, common.Negative},
		{1, 0, common.Negative},
		{1, 1, common.Positive},
		{2, 0, common.Neutral},
		{2, 1, common.Neutral},
		{3, 0, common.Empty},
		{3, 1, common.Empty},
	}

	for _, testCase := range testCases {
		answer := l.Get(testCase.row, testCase.col)
		if answer != testCase.expected {
			t.Errorf("ERROR: For %d, %d expected '%c' got '%c'", testCase.row, testCase.col, testCase.expected, answer)
		}
	}
}

func TestEqual(t *testing.T) {
	l := New(2, 4)

	l.Set(0, 0, common.Positive)
	l.Set(0, 1, common.Negative)
	l.Set(1, 0, common.Empty)
	l.Set(1, 1, common.Empty)
	l.Set(2, 0, common.Neutral)
	l.Set(2, 1, common.Neutral)
	l.Set(3, 0, common.Empty)
	l.Set(3, 1, common.Empty)

	l2 := New(2, 4)

	answer := l.Equal(l)
	expected := true
	if answer != expected {
		t.Errorf("ERROR: For l == l expected %t got %t", expected, answer)
	}

	answer = l.Equal(l2)
	expected = false
	if answer != expected {
		t.Errorf("ERROR: For l == l2 expected %t got %t", expected, answer)
	}
}
