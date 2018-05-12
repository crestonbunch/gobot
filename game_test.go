package gobot

import (
	"reflect"
	"testing"
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

func TestGameAuthorization(t *testing.T) {
	cases := []struct {
		game   *Game
		user   string
		expect bool
	}{
		{
			game: &Game{
				Players: Players{Anyone: true},
				Next:    BlackStone,
			},
			user:   "dummy",
			expect: true,
		}, {
			game: &Game{
				Players: Players{Anyone: true},
				Next:    WhiteStone,
			},
			user:   "dummy",
			expect: true,
		}, {
			game: &Game{
				Players: Players{
					Black: []string{"foo"},
					White: []string{"bar"},
				},
				Next: BlackStone,
			},
			user:   "foo",
			expect: true,
		}, {
			game: &Game{
				Players: Players{
					Black: []string{"foo"},
					White: []string{"bar"},
				},
				Next: WhiteStone,
			},
			user:   "foo",
			expect: false,
		}, {
			game: &Game{
				Players: Players{
					Black: []string{"foo"},
					White: []string{"bar"},
				},
				Next: WhiteStone,
			},
			user:   "bar",
			expect: true,
		}, {
			game: &Game{
				Players: Players{
					Black: []string{"foo"},
					White: []string{"bar"},
				},
				Next: BlackStone,
			},
			user:   "bar",
			expect: false,
		}, {
			game: &Game{
				Players: Players{Anyone: true},
				Next:    EmptyStone,
			},
			user:   "bar",
			expect: false,
		},
	}

	for _, test := range cases {
		result := test.game.Authorized(test.user)

		if result != test.expect {
			t.Errorf(
				"bad authorization for user %s in game %v",
				test.user, test.game,
			)
		}
	}
}

func TestGameMove(t *testing.T) {
	cases := []struct {
		desc   string
		game   *Game
		player string
		coords [2]int
		expect *Game
		err    bool
	}{
		{
			desc: "anyone move at A1",
			game: &Game{
				History: History([]Board{[][]Stone{
					{EmptyStone, EmptyStone},
					{EmptyStone, EmptyStone},
				}}),
				Next:     BlackStone,
				Players:  Players{Anyone: true},
				Settings: Settings{},
				Captures: Captures{0, 0},
				Passes:   Passes{},
				Votes: Votes{
					Moves:  map[string][2]int{},
					Passes: []string{},
				},
				Finished: false,
			},
			expect: &Game{
				History: History([]Board{
					[][]Stone{
						{EmptyStone, EmptyStone},
						{EmptyStone, EmptyStone},
					},
					[][]Stone{
						{BlackStone, EmptyStone},
						{EmptyStone, EmptyStone},
					}}),
				Next:     WhiteStone,
				Players:  Players{Anyone: true},
				Settings: Settings{},
				Captures: Captures{0, 0},
				Passes:   Passes{},
				Votes: Votes{
					Moves:  map[string][2]int{},
					Passes: []string{},
				},
				Finished: false,
			},
			player: "foo",
			coords: [2]int{0, 0},
			err:    false,
		}, {
			desc: "white captures two stones after pass",
			game: &Game{
				History: History([]Board{[][]Stone{
					{BlackStone, EmptyStone},
					{WhiteStone, BlackStone},
				}}),
				Next:    WhiteStone,
				Players: Players{Anyone: true},
				Passes:  Passes{Black: true},
			},
			expect: &Game{
				History: History([]Board{
					[][]Stone{
						{BlackStone, EmptyStone},
						{WhiteStone, BlackStone},
					},
					[][]Stone{
						{EmptyStone, WhiteStone},
						{WhiteStone, EmptyStone},
					}}),
				Next:     BlackStone,
				Players:  Players{Anyone: true},
				Captures: Captures{Black: 0, White: 2},
			},
			player: "foo",
			coords: [2]int{1, 0},
			err:    false,
		}, {
			desc: "black captures two stones after pass",
			game: &Game{
				History: History([]Board{[][]Stone{
					{WhiteStone, EmptyStone},
					{BlackStone, WhiteStone},
				}}),
				Next:    BlackStone,
				Players: Players{Anyone: true},
				Passes:  Passes{White: true},
			},
			expect: &Game{
				History: History([]Board{
					[][]Stone{
						{WhiteStone, EmptyStone},
						{BlackStone, WhiteStone},
					},
					[][]Stone{
						{EmptyStone, BlackStone},
						{BlackStone, EmptyStone},
					}}),
				Next:     WhiteStone,
				Players:  Players{Anyone: true},
				Captures: Captures{Black: 2, White: 0},
			},
			player: "foo",
			coords: [2]int{1, 0},
			err:    false,
		}, {
			desc: "unauthorized move at A1",
			game: &Game{
				History: History([]Board{[][]Stone{
					{EmptyStone, EmptyStone},
					{EmptyStone, EmptyStone},
				}}),
				Next:    BlackStone,
				Players: Players{Anyone: false},
			},
			expect: &Game{
				History: History([]Board{[][]Stone{
					{EmptyStone, EmptyStone},
					{EmptyStone, EmptyStone},
				}}),
				Next:    BlackStone,
				Players: Players{Anyone: false},
			},
			player: "foo",
			coords: [2]int{0, 0},
			err:    true,
		}, {
			desc: "ko at A2",
			game: &Game{
				History: History([]Board{
					[][]Stone{
						{EmptyStone, BlackStone, WhiteStone},
						{BlackStone, WhiteStone, EmptyStone},
						{EmptyStone, EmptyStone, EmptyStone},
					}, [][]Stone{
						{WhiteStone, EmptyStone, WhiteStone},
						{BlackStone, WhiteStone, EmptyStone},
						{EmptyStone, EmptyStone, EmptyStone},
					}}),
				Next:    BlackStone,
				Players: Players{Anyone: true},
			},
			expect: &Game{
				History: History([]Board{
					[][]Stone{
						{EmptyStone, BlackStone, WhiteStone},
						{BlackStone, WhiteStone, EmptyStone},
						{EmptyStone, EmptyStone, EmptyStone},
					}, [][]Stone{
						{WhiteStone, EmptyStone, WhiteStone},
						{BlackStone, WhiteStone, EmptyStone},
						{EmptyStone, EmptyStone, EmptyStone},
					}}),
				Next:    BlackStone,
				Players: Players{Anyone: true},
			},

			player: "foo",
			coords: [2]int{1, 0},
			err:    true,
		}, {
			desc: "illegal move at A1",
			game: &Game{
				History: History([]Board{
					[][]Stone{
						{EmptyStone, BlackStone},
						{BlackStone, EmptyStone},
					}}),
				Next:    WhiteStone,
				Players: Players{Anyone: true},
			},
			expect: &Game{
				History: History([]Board{
					[][]Stone{
						{EmptyStone, BlackStone},
						{BlackStone, EmptyStone},
					}}),
				Next:    WhiteStone,
				Players: Players{Anyone: true},
			},
			player: "foo",
			coords: [2]int{0, 0},
			err:    true,
		},
	}
	for _, test := range cases {
		err := test.game.Move(test.player, test.coords)
		if err != nil && !test.err {
			t.Errorf("unexpected error %s for %s", err.Error(), test.desc)
		} else if err == nil && test.err {
			t.Errorf("expected error for %s", test.desc)
		} else if !reflect.DeepEqual(test.game, test.expect) {
			t.Errorf(
				"%s\nexpected\n%#v\nbut got\n%#v\n",
				test.desc, test.expect, test.game,
			)
		}
	}
}

func TestGameVoteForMove(t *testing.T) {
	cases := []struct {
		desc   string
		game   *Game
		player string
		coords [2]int
		expect *Game
		err    bool
	}{
		{
			desc: "anyone vote at A1",
			game: &Game{
				History: History([]Board{[][]Stone{
					{EmptyStone, EmptyStone},
					{EmptyStone, EmptyStone},
				}}),
				Next:     BlackStone,
				Players:  Players{Anyone: true},
				Settings: Settings{Vote: true},
				Captures: Captures{0, 0},
				Passes:   Passes{},
				Votes: Votes{
					Moves:  map[string][2]int{},
					Passes: []string{},
				},
				Finished: false,
			},
			expect: &Game{
				History: History([]Board{
					[][]Stone{
						{EmptyStone, EmptyStone},
						{EmptyStone, EmptyStone},
					},
				}),
				Next:     BlackStone,
				Players:  Players{Anyone: true},
				Settings: Settings{Vote: true},
				Captures: Captures{0, 0},
				Passes:   Passes{},
				Votes: Votes{
					Moves: map[string][2]int{
						"foo": {0, 0},
					},
					Passes: []string{},
				},
				Finished: false,
			},
			player: "foo",
			coords: [2]int{0, 0},
			err:    false,
		}, {
			desc: "unauthorized move at A1",
			game: &Game{
				History: History([]Board{[][]Stone{
					{EmptyStone, EmptyStone},
					{EmptyStone, EmptyStone},
				}}),
				Next:     BlackStone,
				Players:  Players{Anyone: false},
				Settings: Settings{Vote: true},
			},
			expect: &Game{
				History: History([]Board{[][]Stone{
					{EmptyStone, EmptyStone},
					{EmptyStone, EmptyStone},
				}}),
				Next:     BlackStone,
				Players:  Players{Anyone: false},
				Settings: Settings{Vote: true},
			},
			player: "foo",
			coords: [2]int{0, 0},
			err:    true,
		}, {
			desc: "ko at A2",
			game: &Game{
				History: History([]Board{
					[][]Stone{
						{EmptyStone, BlackStone, WhiteStone},
						{BlackStone, WhiteStone, EmptyStone},
						{EmptyStone, EmptyStone, EmptyStone},
					}, [][]Stone{
						{WhiteStone, EmptyStone, WhiteStone},
						{BlackStone, WhiteStone, EmptyStone},
						{EmptyStone, EmptyStone, EmptyStone},
					}}),
				Next:     BlackStone,
				Players:  Players{Anyone: true},
				Settings: Settings{Vote: true},
			},
			expect: &Game{
				History: History([]Board{
					[][]Stone{
						{EmptyStone, BlackStone, WhiteStone},
						{BlackStone, WhiteStone, EmptyStone},
						{EmptyStone, EmptyStone, EmptyStone},
					}, [][]Stone{
						{WhiteStone, EmptyStone, WhiteStone},
						{BlackStone, WhiteStone, EmptyStone},
						{EmptyStone, EmptyStone, EmptyStone},
					}}),
				Next:     BlackStone,
				Players:  Players{Anyone: true},
				Settings: Settings{Vote: true},
			},
			player: "foo",
			coords: [2]int{1, 0},
			err:    true,
		},
	}
	for _, test := range cases {
		err := test.game.VoteForMove(test.player, test.coords)
		if err != nil && !test.err {
			t.Errorf("unexpected error %s for %s", err.Error(), test.desc)
		} else if err == nil && test.err {
			t.Errorf("expected error for %s", test.desc)
		} else if !reflect.DeepEqual(test.game, test.expect) {
			t.Errorf(
				"%s\nexpected\n%#v\nbut got\n%#v\n",
				test.desc, test.expect, test.game,
			)
		}
	}
}

func TestGameVoteForPass(t *testing.T) {
	cases := []struct {
		desc   string
		game   *Game
		player string
		expect *Game
		err    bool
	}{
		{
			desc: "vote for pass",
			game: &Game{
				Next:     BlackStone,
				Players:  Players{Anyone: true},
				Settings: Settings{Vote: true},
				Passes:   Passes{},
				Votes: Votes{
					Moves:  map[string][2]int{},
					Passes: []string{},
				},
				Finished: false,
			},
			expect: &Game{
				Next:     BlackStone,
				Players:  Players{Anyone: true},
				Settings: Settings{Vote: true},
				Captures: Captures{0, 0},
				Passes:   Passes{},
				Votes: Votes{
					Moves:  map[string][2]int{},
					Passes: []string{"foo"},
				},
				Finished: false,
			},
			player: "foo",
			err:    false,
		}, {
			desc: "unauthorized pass",
			game: &Game{
				Next:    BlackStone,
				Players: Players{Anyone: false},
			},
			expect: &Game{
				Next:    BlackStone,
				Players: Players{Anyone: false},
			},
			player: "foo",
			err:    true,
		},
	}
	for _, test := range cases {
		err := test.game.VoteForPass(test.player)
		if err != nil && !test.err {
			t.Errorf("unexpected error %s for %s", err.Error(), test.desc)
		} else if err == nil && test.err {
			t.Errorf("expected error for %s", test.desc)
		} else if !reflect.DeepEqual(test.game, test.expect) {
			t.Errorf(
				"%s\nexpected\n%#v\nbut got\n%#v\n",
				test.desc, test.expect, test.game,
			)
		}
	}
}

func TestGamePass(t *testing.T) {
	cases := []struct {
		desc   string
		game   *Game
		player string
		expect *Game
		err    bool
	}{
		{
			desc: "black pass",
			game: &Game{
				Next:     BlackStone,
				Players:  Players{Anyone: true},
				Settings: Settings{},
				Passes:   Passes{},
				Votes: Votes{
					Moves:  map[string][2]int{},
					Passes: []string{},
				},
				Finished: false,
			},
			expect: &Game{
				Next:     WhiteStone,
				Players:  Players{Anyone: true},
				Settings: Settings{},
				Captures: Captures{0, 0},
				Passes:   Passes{Black: true},
				Votes: Votes{
					Moves:  map[string][2]int{},
					Passes: []string{},
				},
				Finished: false,
			},
			player: "foo",
			err:    false,
		}, {
			desc: "white pass",
			game: &Game{
				Next:     WhiteStone,
				Players:  Players{Anyone: true},
				Settings: Settings{},
				Passes:   Passes{},
				Votes: Votes{
					Moves:  map[string][2]int{},
					Passes: []string{},
				},
				Finished: false,
			},
			expect: &Game{
				Next:     BlackStone,
				Players:  Players{Anyone: true},
				Settings: Settings{},
				Captures: Captures{0, 0},
				Passes:   Passes{White: true},
				Votes: Votes{
					Moves:  map[string][2]int{},
					Passes: []string{},
				},
				Finished: false,
			},
			player: "foo",
			err:    false,
		},
	}
	for _, test := range cases {
		err := test.game.Pass(test.player)
		if err != nil && !test.err {
			t.Errorf("unexpected error %s for %s", err.Error(), test.desc)
		} else if err == nil && test.err {
			t.Errorf("expected error for %s", test.desc)
		} else if !reflect.DeepEqual(test.game, test.expect) {
			t.Errorf(
				"%s\nexpected\n%#v\nbut got\n%#v\n",
				test.desc, test.expect, test.game,
			)
		}
	}
}
