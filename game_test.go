package gobot

import "testing"

func TestKo(t *testing.T) {
	history := History([]Board{
		Board([][]Stone{
			[]Stone{EmptyStone, BlackStone, WhiteStone},
			[]Stone{EmptyStone, WhiteStone, EmptyStone},
			[]Stone{EmptyStone, EmptyStone, WhiteStone},
		}),
		Board([][]Stone{
			[]Stone{EmptyStone, BlackStone, EmptyStone},
			[]Stone{EmptyStone, WhiteStone, BlackStone},
			[]Stone{EmptyStone, EmptyStone, WhiteStone},
		}),
	})

	cases := []struct {
		board   Board
		history History
		expect  bool
	}{
		{
			Board([][]Stone{
				[]Stone{EmptyStone, BlackStone},
				[]Stone{EmptyStone, EmptyStone},
			}),
			History([]Board{}),
			false,
		}, {
			Board([][]Stone{
				[]Stone{EmptyStone, BlackStone, WhiteStone},
				[]Stone{EmptyStone, WhiteStone, EmptyStone},
				[]Stone{EmptyStone, EmptyStone, WhiteStone},
			}),
			history,
			true,
		}, {
			Board([][]Stone{
				[]Stone{EmptyStone, BlackStone, EmptyStone},
				[]Stone{EmptyStone, WhiteStone, BlackStone},
				[]Stone{EmptyStone, WhiteStone, WhiteStone},
			}),
			history,
			false,
		},
	}

	for _, test := range cases {
		actual := test.history.Ko(test.board)

		if actual != test.expect {
			t.Errorf(
				"Expected %v -> %v to be a Ko but was not",
				history, test.board,
			)
		}
	}
}
