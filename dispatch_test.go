package gobot

import (
	"testing"
)

func TestMoveDispatch(t *testing.T) {
	one := NewGame(Players{
		Anyone: true,
	}, Settings{})
	two := NewGame(Players{
		Black: []string{"bar"},
		White: []string{"foo"},
	}, Settings{})
	cases := []struct {
		desc   string
		games  []*Game
		action *Action
		expect *Game
		err    bool
	}{
		{
			desc:   "latest game with two players, player 1",
			games:  []*Game{one, two},
			action: &Action{User: "foo", Command: &MoveCommand{GameID: -1}},
			expect: two,
			err:    false,
		}, {
			desc:   "single game with any players",
			games:  []*Game{one},
			action: &Action{User: "foo", Command: &MoveCommand{GameID: -1}},
			expect: one,
			err:    false,
		}, {
			desc:   "first game with any players",
			games:  []*Game{one, two},
			action: &Action{User: "baz", Command: &MoveCommand{GameID: -1}},
			expect: one,
			err:    false,
		}, {
			desc:   "last game with two players, player 2",
			games:  []*Game{two, one},
			action: &Action{User: "bar", Command: &MoveCommand{GameID: -1}},
			expect: one,
			err:    false,
		}, {
			desc:   "first game with fixed id",
			games:  []*Game{one, two},
			action: &Action{User: "foo", Command: &MoveCommand{GameID: 0}},
			expect: one,
			err:    false,
		}, {
			desc:   "no games with id -1",
			games:  []*Game{},
			action: &Action{User: "bar", Command: &MoveCommand{GameID: -1}},
			expect: nil,
			err:    true,
		}, {
			desc:   "one game with id out of bounds",
			games:  []*Game{one},
			action: &Action{User: "bar", Command: &MoveCommand{GameID: 2}},
			expect: nil,
			err:    true,
		},
	}
	for _, test := range cases {
		result, err := Dispatch(test.games, test.action)
		if !test.err && err != nil {
			t.Errorf("unexpected error %s for %s ", err.Error(), test.desc)
		} else if test.err && err == nil {
			t.Errorf("expected error for %s", test.desc)
		} else if result != test.expect {
			t.Errorf("wrong game selected %s", test.desc)
		}
	}
}
