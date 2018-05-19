package gobot_test

import (
	"testing"

	. "github.com/crestonbunch/gobot"
)

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

func TestBoardSizes(t *testing.T) {
	cases := []struct {
		constructor func() Board
		width       int
		height      int
	}{
		{New19by19Board, 19, 19},
	}

	for _, test := range cases {
		board := test.constructor()
		if test.height != len(board) {
			t.Errorf("expected height %d but got %d", test.height, len(board))
		}
		for _, row := range board {
			if test.width != len(row) {
				t.Errorf("expected width %d but got %d", test.width, len(row))
			}
		}
	}
}

func TestBoardGet(t *testing.T) {
	board := Board([][]Stone{
		{EmptyStone, EmptyStone},
		{BlackStone, WhiteStone},
	})

	cases := []struct {
		x      int
		y      int
		expect Stone
	}{
		{0, 0, EmptyStone},
		{0, 1, BlackStone},
		{1, 0, EmptyStone},
		{1, 1, WhiteStone},
		{-1, 0, BoundaryStone},
		{0, -1, BoundaryStone},
		{3, 1, BoundaryStone},
		{1, 3, BoundaryStone},
	}

	for _, test := range cases {
		stone := board.Get(test.x, test.y)

		if stone != test.expect {
			t.Errorf(
				"Expected stone at (%d, %d) to be %d but got %d",
				test.x, test.y, test.expect, stone,
			)
		}
	}
}

func TestBoardSet(t *testing.T) {
	board := Board([][]Stone{
		{WhiteStone, BlackStone},
		{EmptyStone, EmptyStone},
	})

	cases := []struct {
		x      int
		y      int
		stone  Stone
		expect Board
	}{
		{0, 0, EmptyStone, Board([][]Stone{
			{EmptyStone, BlackStone},
			{EmptyStone, EmptyStone},
		})},
		{0, 1, BlackStone, Board([][]Stone{
			{WhiteStone, BlackStone},
			{BlackStone, EmptyStone},
		})},
		{1, 0, EmptyStone, Board([][]Stone{
			{WhiteStone, EmptyStone},
			{EmptyStone, EmptyStone},
		})},
		{1, 1, WhiteStone, Board([][]Stone{
			{WhiteStone, BlackStone},
			{EmptyStone, WhiteStone},
		})},
		{-1, 0, WhiteStone, Board([][]Stone{
			{WhiteStone, BlackStone},
			{EmptyStone, EmptyStone},
		})},
		{0, -1, WhiteStone, Board([][]Stone{
			{WhiteStone, BlackStone},
			{EmptyStone, EmptyStone},
		})},
		{3, 0, WhiteStone, Board([][]Stone{
			{WhiteStone, BlackStone},
			{EmptyStone, EmptyStone},
		})},
		{0, 3, WhiteStone, Board([][]Stone{
			{WhiteStone, BlackStone},
			{EmptyStone, EmptyStone},
		})},
	}

	for _, test := range cases {
		result := board.Set(test.x, test.y, test.stone)

		if &board == &result {
			t.Errorf("Board and result are the same reference!")
		}

		for i, row := range result {
			for j, stone := range row {
				if test.expect[i][j] != stone {
					t.Errorf(
						"Expected stone at (%d %d) to be %d but got %d",
						test.x, test.y, test.expect[i][j], stone,
					)
				}
			}
		}
	}
}

func TestBoardEquals(t *testing.T) {
	board := Board([][]Stone{
		{WhiteStone, BlackStone},
		{EmptyStone, EmptyStone},
	})

	cases := []struct {
		board  Board
		expect bool
	}{
		{
			Board([][]Stone{
				{WhiteStone, BlackStone},
				{EmptyStone, EmptyStone},
			}),
			true,
		},
		{
			Board([][]Stone{
				{BlackStone, BlackStone},
				{EmptyStone, EmptyStone},
			}),
			false,
		},
		{
			Board([][]Stone{
				{WhiteStone, BlackStone},
			}),
			false,
		},
		{
			Board([][]Stone{
				{WhiteStone, BlackStone},
				{EmptyStone},
			}),
			false,
		},
	}

	for _, test := range cases {
		result := board.Equals(test.board)

		if result != test.expect {
			t.Errorf("%v does not equal %v", board, test.board)
		}
	}
}

func TestBoardLiberties(t *testing.T) {
	cases := []struct {
		board  Board
		x      int
		y      int
		expect int
	}{
		{
			Board([][]Stone{
				{WhiteStone, BlackStone},
				{EmptyStone, EmptyStone},
			}),
			0, 0, 1,
		},
		{
			Board([][]Stone{
				{WhiteStone, BlackStone},
				{EmptyStone, EmptyStone},
			}),
			1, 0, 1,
		},
		{
			Board([][]Stone{
				{WhiteStone, BlackStone},
				{EmptyStone, BlackStone},
			}),
			1, 0, 1,
		},
		{
			Board([][]Stone{
				{WhiteStone, BlackStone},
				{WhiteStone, BlackStone},
			}),
			1, 0, 0,
		},
		{
			Board([][]Stone{
				{BlackStone, BlackStone, EmptyStone},
				{BlackStone, BlackStone, BlackStone},
				{EmptyStone, BlackStone, BlackStone},
			}),
			1, 1, 2,
		},
	}

	for _, test := range cases {
		result := test.board.Liberties(test.x, test.y)

		if result != test.expect {
			t.Errorf(
				"counted %d liberties but expected %d", result, test.expect,
			)
		}
	}
}

func TestBoardCapture(t *testing.T) {
	cases := []struct {
		board    Board
		x        int
		y        int
		expect   Board
		captures int
	}{
		{
			Board([][]Stone{
				{WhiteStone, BlackStone},
				{EmptyStone, EmptyStone},
			}),
			0, 0,
			Board([][]Stone{
				{EmptyStone, BlackStone},
				{EmptyStone, EmptyStone},
			}),
			1,
		}, {
			Board([][]Stone{
				{WhiteStone, BlackStone, EmptyStone},
				{WhiteStone, WhiteStone, BlackStone},
				{WhiteStone, WhiteStone, WhiteStone},
			}),
			1, 1,
			Board([][]Stone{
				{EmptyStone, BlackStone, EmptyStone},
				{EmptyStone, EmptyStone, BlackStone},
				{EmptyStone, EmptyStone, EmptyStone},
			}),
			6,
		},
	}

	for _, test := range cases {
		board, captures := test.board.Capture(test.x, test.y)

		if !board.Equals(test.expect) {
			t.Errorf(
				"expected capture %v but got %v", test.expect, board,
			)
		}
		if captures != test.captures {
			t.Errorf(
				"expected %d captures but got %d", test.captures, captures,
			)
		}
	}
}

func TestBoardPlay(t *testing.T) {
	cases := []struct {
		board    Board
		expect   Board
		x        int
		y        int
		stone    Stone
		captures int
		err      bool
	}{
		{
			Board([][]Stone{
				{WhiteStone, BlackStone, EmptyStone},
				{WhiteStone, WhiteStone, BlackStone},
				{WhiteStone, WhiteStone, WhiteStone},
			}),
			nil, 0, 0, WhiteStone, 0, true,
		}, {
			Board([][]Stone{
				{EmptyStone, BlackStone, EmptyStone},
				{BlackStone, EmptyStone, BlackStone},
				{WhiteStone, WhiteStone, WhiteStone},
			}),
			nil, 0, 0, WhiteStone, 0, true,
		}, {
			Board([][]Stone{
				{WhiteStone, BlackStone, EmptyStone},
				{WhiteStone, WhiteStone, BlackStone},
				{WhiteStone, WhiteStone, WhiteStone},
			}),
			Board([][]Stone{
				{WhiteStone, EmptyStone, WhiteStone},
				{WhiteStone, WhiteStone, EmptyStone},
				{WhiteStone, WhiteStone, WhiteStone},
			}),
			2, 0, WhiteStone, 2, false,
		},
	}

	for _, test := range cases {
		board, captures, err := test.board.Play(test.x, test.y, test.stone)

		if test.expect == nil && board != nil {
			t.Errorf("expected nil but got %v", board)
		} else if board == nil && test.expect != nil {
			t.Errorf("expected %v but got nil", test.expect)
		} else if !board.Equals(test.expect) {
			t.Errorf(
				"expected %v but got %v", test.expect, board,
			)
		}
		if captures != test.captures {
			t.Errorf(
				"expected %d captures but got %d", test.captures, captures,
			)
		}
		if err == nil && test.err {
			t.Errorf(
				"expected err playing (%d, %d) with %v",
				test.x, test.y, test.board,
			)
		}
	}
}
