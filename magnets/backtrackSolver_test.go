package magnets

import (
	"testing"
)

func TestSetCell(t *testing.T) {
}

func TestCountSolutions(t *testing.T) {
	testCases := []struct {
		game     string
		expected int
	}{
		{"2x2:11,00,11,00,TTBB", 0},
		{"1x2:1,10,1,01,TB", 1},
		{"1x3:1,100,1,010,TB*", 1},
		{"2x2:11,11,11,11,TTBB", 2},
		{"3x2:111,21,111,12,TTTBBB", 1},
		{"5x2:11011,22,11011,22,LRTLRLRBLR", 4},

		// Came from the iPhone. Guaranteed to have only one solution.
		{"4x5:3222,22122,2322,22122,LRTTTTBBBBLRTTTTBBBB", 1},
		{"4x5:3201,11211,2301,11121,LRTTTTBBBBTTTTBBBBLR", 1},
		{"5x5:11212,12211,12202,21112,TTTTTBBBBBTLRT*BLRBTLRLRB", 1},
		{"5x5:21223,31123,12232,22132,LRTLRT*BLRBLRTTLRTBBLRBLR", 1},
		{"5x5:22122,12312,22122,13212,LRLR*TTTLRBBBTTLRTBBLRBLR", 1},
		{"5x5:31222,22312,22132,23221,LRTT*TTBBTBBLRBTLRTTBLRBB", 1},
		{"6x5:122232,33132,112323,32232,LRLRLRTLRLRTBLRLRBLRTTTTLRBBBB", 1},
		{"6x5:122123,23312,221132,23312,TTLRTTBBLRBBLRLRLRTLRTLRBLRBLR", 1},
		{"6x5:022131,23202,121212,32202,TTLRLRBBLRLRTLRLRTBLRLRBLRLRLR", 1},
		// Too slow
		// {"7x9:1444133,330323222,2244242,232333211,LRLRLRTTLRLRTBBLRT*BTTTTBLRBBBBTLRTLRTBTTBTTBTBBTBBTBLRBLRBLRLR", 1},
		// {"7x9:4442324,412332242,3343343,333312233,LRLRTLRTLR*BLRBTTTLRTTBBBTTBBTTTBBTTBBBTTBBTLRBBTTBLRLRBBLRLRLR", 1},
		// {"7x9:3342325,331221433,2442334,331222342,TLRLRTTBLRLRBBTTLR*TTBBTTTBBTTBBBTTBBTLRBBTTBTTTTBBTBBBBLRBLRLR", 1},
		// {"8x7:30122222,2232302,30311231,2321321,LRLRLRLRTLRLRLRTBTLRTLRBTBTTBLRTBTBBTLRBTBLRBLRTBLRLRLRB", 1},
		// {"8x7:22423232,2223344,04243223,3231434,LRLRLRTTLRLRTTBBTTLRBBTTBBLRTTBBTTLRBBTTBBLRTTBBLRLRBBLR", 1},
		// {"10x10:4444543143,4444043445,4435443324,5254135254,TTTTLRLRLRBBBBLRTTLRTLRLRTBBLRBLRLRBLRLRTTLRLRTLRTBBLRLRBLRBTTTTTLRLRTBBBBBLRLRBTLRTTLRTLRBLRBBLRBLR", 1},
		// {"12x10:555213545454,5543546565,555215354544,5434465656,TTTLRLRTTLRTBBBTLRTBBTTBLRTBTTBLRBBTLRBTBBLRTLRBLRTBLRLRBLRTLRBLRTTLRLRBTLRTTBBTTLRTBLRBBLRBBLRBTTTTLRTLRTLRBBBBLRBLRBLR", 1},
		// {"12x10:545232142444,5462542336,444423124453,5652451345,TTTTLRLRLRLRBBBBTLRTLRTTLRLRBTTBLRBBTLRLRBBTLRLRBLRTTLRBLRTTTLRBBTLRLRBBBLRTTBTLRTTTLRTBBTBLRBBBTTBLRBTLRLRTBBLRLRBLRLRB", 1},
		// {"15x9:454544143332433,766636477,545453133524242,775455667,LRTTLRTTTLRTTLRLRBBTTBBB*TBBTTTTLRBBLRLRBLRBBBBTLRLRTTTTLRTTTTBLRLRBBBBTTBBBBLRLRTLRLRBBTTLRTTTTBLRTTLRBBLRBBBBLRTBBTTLRLRLRLRLRBLRBBLR", 1},
		// {"15x9:333445434434343,786555657,442354433434434,768646575,TTLRLRLR*LRTLRTBBTLRLRTLRTBLRBLRBLRLRBTTBTLRTTLRTTTLRBBTBLRBBTTBBBTTLRBLRLRTBBTLRBBLRTTTTTBLRBLRLRLRBBBBBLRTTTLRTTTLRLRTLRBBBLRBBBLRLRB", 1},
		// {"15x9:423444254334424,686756455,414443345334334,758766544,TLRLRLRLRLRLRLRB*TTTTLRTLRTTTTLRBBBBLRBTTBBBBLRTLRLRLRBBLRLRTTBLRLRLRLRTLRTBBTLRLRTLRTBLRBTTBTLRTBLRBLRLRBBTBLRBLRTTTLRTLRBLRLRLRBBBLRB", 1},
		// {"15x15:858767847577878,778686866676768,686778755768787,787777775767677,TLRT*LRTTTLRTLRBLRBTLRBBBLRBLRLRLRBLRLRTTTTTTTLRLRLRLRBBBBBBBLRTLRTTLRLRTLRTLRBTTBBLRLRBLRBLRTBBTLRLRTTTTLRTBLRBLRTTBBBBLRBLRLRLRBBLRTTLRTLRLRTTTLRTBBTTBTTLRBBBLRBLRBBTBBTTTTTTTTTTLRBLRBBBBBBBBBBTTLRLRTTLRTLRLRBBLRLRBBLRBLRLR", 1},
		// {"15x15:467686748786757,767787767776455,556777657878666,686877776865555,TLRTTLRTLRLR*LRBLRBBLRBTTLRTLRLRLRLRLRBBLRBTTTLRTLRLRLRTLRBBBTTBLRTTLRBLRTTTBBTLRBBLRTLRBBBLRBLRTTLRBTLRTTTTTLRBBTTTBTTBBBBBLRLRBBBTBBTLRTTLRLRTTTBLRBLRBBLRLRBBBLRTTTLRTLRLRTTTLRBBBLRBTTTTBBBLRLRLRTTBBBBTTTLRLRLRBBLRLRBBBLRLR", 1},
		// {"17x17:88775759697188978,77876687774878868,87876667887189878,68784778684888688,TLRTLRLRLRLRTTLRTBLRBLRTLRTT*BBTTBTTTLRTBLRBBTLRBBTBBBLRBTLRLRBLRTTBTLRTLRBLRTTTLRBBTBLRBLRTLRBBBTTTTBTLRTTTBLRLRTBBBBTBLRBBBLRLRTBTTLRBTLRTLRTLRTBTBBTTTBLRBLRBLRBTBLRBBBLRLRLRLRTTBTTTTLRLRLRLRLRBBTBBBBLRTTTTTTLRLRBLRTLRTBBBBBBLRLRTLRBLRBTTLRLRLRTTBLRLRTTBBTLRLRTBBTTTTTBBLRBLRLRBLRBBBBBLR", 1},
	}

	for _, testCase := range testCases {
		game, ok := Deserialize(testCase.game)
		if !ok {
			t.Errorf("ERROR: failed to deserialize %s", testCase.game)
		}

		answer := game.CountSolutions(0, 0)
		if answer != testCase.expected {
			t.Errorf("ERROR: for %s expected %d got %d", testCase.game, testCase.expected, answer)
		}
	}
}

func TestSingleSolution(t *testing.T) {
	testCases := []struct {
		game     string
		expected bool
	}{
		{"2x2:11,00,11,00,TTBB", false}, // 0 solutions
		{"1x2:1,10,1,01,TB", true},      // 1 solution
		{"2x2:11,11,11,11,TTBB", false}, // 2 solutions
	}

	for _, testCase := range testCases {
		game, ok := Deserialize(testCase.game)
		if !ok {
			t.Errorf("ERROR: failed to deserialize %s", testCase.game)
		}

		answer := game.singleSolution()
		if answer != testCase.expected {
			t.Errorf("ERROR: for %s expected %t got %t", testCase.game, testCase.expected, answer)
		}
	}
}
