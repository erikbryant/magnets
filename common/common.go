package common

// https://en.wikipedia.org/wiki/List_of_Unicode_characters#Mathematical_symbols
const (
	Empty    = ' '
	Positive = '+'
	Negative = '-'
	Neutral  = '#'
	Wall     = '⊠'
	Up       = '⋂'
	Down     = '⋃'
	Left     = '⊂'
	Right    = '⊃'
	Marker   = '?'
	Border   = 'X'
)

// Negate returns the opposite value of the given rune.
func Negate(r rune) rune {
	switch r {
	case Positive:
		return Negative
	case Negative:
		return Positive
	case Up:
		return Down
	case Down:
		return Up
	case Left:
		return Right
	case Right:
		return Left
	}

	// All others negate to themselves.
	return r
}
