package gobot

import "testing"

func TestKo(t *testing.T) {
	history := History([]Board{
		Board([][]Stone{
			{EmptyStone, BlackStone, WhiteStone},
			{EmptyStone, WhiteStone, EmptyStone},
			{EmptyStone, EmptyStone, WhiteStone},
		}),
		Board([][]Stone{
			{EmptyStone, BlackStone, EmptyStone},
			{EmptyStone, WhiteStone, BlackStone},
			{EmptyStone, EmptyStone, WhiteStone},
		}),
	})

	cases := []struct {
		board   Board
		history History
		expect  bool
	}{
		{
			Board([][]Stone{
				{EmptyStone, BlackStone},
				{EmptyStone, EmptyStone},
			}),
			History([]Board{}),
			false,
		}, {
			Board([][]Stone{
				{EmptyStone, BlackStone, WhiteStone},
				{EmptyStone, WhiteStone, EmptyStone},
				{EmptyStone, EmptyStone, WhiteStone},
			}),
			history,
			true,
		}, {
			Board([][]Stone{
				{EmptyStone, BlackStone, EmptyStone},
				{EmptyStone, WhiteStone, BlackStone},
				{EmptyStone, WhiteStone, WhiteStone},
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
