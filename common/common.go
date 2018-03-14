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
	return r
}
