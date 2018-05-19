package gobot_test

import (
	"reflect"
	"testing"
	"time"

	. "github.com/crestonbunch/gobot"
)

func TestStateCanMove(t *testing.T) {
	cases := []struct {
		state  *State
		user   string
		expect bool
	}{
		{
			state: &State{
				Players: Players{Anyone: true},
				Next:    BlackStone,
			},
			user:   "dummy",
			expect: true,
		}, {
			state: &State{
				Players: Players{Anyone: true},
				Next:    WhiteStone,
			},
			user:   "dummy",
			expect: true,
		}, {
			state: &State{
				Players: Players{
					Black: []string{"foo"},
					White: []string{"bar"},
				},
				Next: BlackStone,
			},
			user:   "foo",
			expect: true,
		}, {
			state: &State{
				Players: Players{
					Black: []string{"foo"},
					White: []string{"bar"},
				},
				Next: WhiteStone,
			},
			user:   "foo",
			expect: false,
		}, {
			state: &State{
				Players: Players{
					Black: []string{"foo"},
					White: []string{"bar"},
				},
				Next: WhiteStone,
			},
			user:   "bar",
			expect: true,
		}, {
			state: &State{
				Players: Players{
					Black: []string{"foo"},
					White: []string{"bar"},
				},
				Next: BlackStone,
			},
			user:   "bar",
			expect: false,
		}, {
			state: &State{
				Players: Players{Anyone: true},
				Next:    EmptyStone,
			},
			user:   "bar",
			expect: false,
		},
	}

	for _, test := range cases {
		result := test.state.CanMove(test.user)

		if result != test.expect {
			t.Errorf(
				"bad authorization for user %s in state %v",
				test.user, test.state,
			)
		}
	}
}

func TestStateMove(t *testing.T) {
	cases := []struct {
		desc   string
		state  *State
		move   *Move
		expect *State
		err    bool
	}{
		{
			desc: "move at A1",
			state: &State{
				History: History([]Board{[][]Stone{
					{EmptyStone, EmptyStone},
					{EmptyStone, EmptyStone},
				}}),
				Next:     BlackStone,
				Voting:   Voting{Required: false},
				Captures: Captures{},
				Passes:   Passes{},
			},
			expect: &State{
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
				Voting:   Voting{},
				Captures: Captures{},
				Passes:   Passes{},
			},
			move: &Move{Coords: [2]int{0, 0}},
			err:  false,
		}, {
			desc: "white captures two stones after pass",
			state: &State{
				History: History([]Board{[][]Stone{
					{BlackStone, EmptyStone},
					{WhiteStone, BlackStone},
				}}),
				Next:    WhiteStone,
				Players: Players{Anyone: true},
				Passes:  Passes{Black: true},
			},
			expect: &State{
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
			move: &Move{Coords: [2]int{1, 0}},
			err:  false,
		}, {
			desc: "black captures two stones after pass",
			state: &State{
				History: History([]Board{[][]Stone{
					{WhiteStone, EmptyStone},
					{BlackStone, WhiteStone},
				}}),
				Next:    BlackStone,
				Players: Players{Anyone: true},
				Passes:  Passes{White: true},
			},
			expect: &State{
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
			move: &Move{Coords: [2]int{1, 0}},
			err:  false,
		}, {
			desc: "black pass",
			state: &State{
				History: History([]Board{[][]Stone{
					{WhiteStone, EmptyStone},
					{BlackStone, WhiteStone},
				}}),
				Next:   BlackStone,
				Passes: Passes{},
			},
			expect: &State{
				History: History([]Board{
					[][]Stone{
						{WhiteStone, EmptyStone},
						{BlackStone, WhiteStone},
					}}),
				Next:   WhiteStone,
				Passes: Passes{Black: true},
			},
			move: &Move{Pass: true},
			err:  false,
		}, {
			desc: "white pass",
			state: &State{
				History: History([]Board{[][]Stone{
					{WhiteStone, EmptyStone},
					{BlackStone, WhiteStone},
				}}),
				Next:   WhiteStone,
				Passes: Passes{},
			},
			expect: &State{
				History: History([]Board{
					[][]Stone{
						{WhiteStone, EmptyStone},
						{BlackStone, WhiteStone},
					}}),
				Next:   BlackStone,
				Passes: Passes{White: true},
			},
			move: &Move{Pass: true},
			err:  false,
		}, {
			desc: "illegal move at A1",
			state: &State{
				History: History([]Board{
					[][]Stone{
						{EmptyStone, BlackStone},
						{BlackStone, EmptyStone},
					}}),
				Next:    WhiteStone,
				Players: Players{Anyone: true},
			},
			expect: &State{
				History: History([]Board{
					[][]Stone{
						{EmptyStone, BlackStone},
						{BlackStone, EmptyStone},
					}}),
				Next:    WhiteStone,
				Players: Players{Anyone: true},
			},
			move: &Move{Coords: [2]int{0, 0}},
			err:  true,
		},
	}
	for _, test := range cases {
		err := test.state.Move(test.move)
		if err != nil && !test.err {
			t.Errorf("unexpected error %s for %s", err.Error(), test.desc)
		} else if err == nil && test.err {
			t.Errorf("expected error for %s", test.desc)
		} else if !reflect.DeepEqual(test.state, test.expect) {
			t.Errorf(
				"%s\nexpected\n%#v\nbut got\n%#v\n",
				test.desc, test.expect, test.state,
			)
		}
	}
}
func TestStateValidate(t *testing.T) {
	cases := []struct {
		desc   string
		state  *State
		move   *Move
		expect bool
	}{
		{
			desc: "lega move at A1",
			state: &State{
				History: History([]Board{
					[][]Stone{
						{EmptyStone, EmptyStone},
						{BlackStone, EmptyStone},
					}}),
				Next:    WhiteStone,
				Players: Players{Anyone: true},
			},
			expect: true,
			move:   &Move{Coords: [2]int{0, 0}},
		}, {
			desc: "illegal move at A1",
			state: &State{
				History: History([]Board{
					[][]Stone{
						{EmptyStone, BlackStone},
						{BlackStone, EmptyStone},
					}}),
				Next:    WhiteStone,
				Players: Players{Anyone: true},
			},
			expect: false,
			move:   &Move{Coords: [2]int{0, 0}},
		}, {
			desc: "ko at A1",
			state: &State{
				History: History([]Board{[][]Stone{
					{WhiteStone, EmptyStone, WhiteStone},
					{BlackStone, WhiteStone, EmptyStone},
					{EmptyStone, EmptyStone, EmptyStone},
				}, [][]Stone{
					{EmptyStone, BlackStone, WhiteStone},
					{BlackStone, WhiteStone, EmptyStone},
					{EmptyStone, EmptyStone, EmptyStone},
				}}),
				Next:    WhiteStone,
				Players: Players{Anyone: true},
			},
			expect: false,
			move:   &Move{Coords: [2]int{0, 0}},
		},
	}
	for _, test := range cases {
		result := test.state.Validate(test.move)
		if result != test.expect {
			t.Errorf("expected %v but got %v for %s",
				test.expect, result, test.desc)
		}
	}
}

func TestStateFinished(t *testing.T) {
	cases := []struct {
		state  *State
		expect bool
	}{
		{
			state: &State{
				Passes: Passes{},
			},
			expect: false,
		}, {
			state: &State{
				Passes: Passes{White: true},
			},
			expect: false,
		}, {
			state: &State{
				Passes: Passes{Black: true},
			},
			expect: false,
		}, {
			state: &State{
				Passes: Passes{White: true, Black: true},
			},
			expect: true,
		},
	}
	for _, test := range cases {
		result := test.state.Finished()
		if result != test.expect {
			t.Errorf("expected %v but got %v",
				test.expect, result)
		}
	}
}

func TestStateVote(t *testing.T) {
	cases := []struct {
		state  *State
		vote   *Move
		expect *State
		err    bool
	}{
		{
			state: &State{
				Votes:  []*Move{},
				Voting: Voting{Required: true},
			},
			expect: &State{
				Votes:  []*Move{{Pass: true}},
				Voting: Voting{Required: true},
			},
			vote: &Move{Pass: true},
			err:  false,
		}, {
			state: &State{
				Votes:  []*Move{},
				Voting: Voting{Required: true},
			},
			expect: &State{
				Votes:  []*Move{{Coords: Coords{3, 3}}},
				Voting: Voting{Required: true},
			},
			vote: &Move{Coords: Coords{3, 3}},
			err:  false,
		},
	}
	for _, test := range cases {
		err := test.state.Vote(test.vote)
		if err != nil && !test.err {
			t.Errorf("unexpected error %s", err.Error())
		}
		if err == nil && test.err {
			t.Errorf("expected error")
		}
		if !reflect.DeepEqual(test.state, test.expect) {
			t.Errorf("expected %#v but got %#v",
				test.expect, test.state)
		}
	}
}

func TestStateSchedule(t *testing.T) {
	cases := []struct {
		state  *State
		expect *time.Timer
	}{
		{
			state:  &State{Voting: Voting{Duration: 0}},
			expect: time.NewTimer(0),
		},
	}
	for _, test := range cases {
		timer := test.state.Schedule()
		if timer == nil {
			t.Errorf("timer was nil")
		}
	}
}

func TestStateRandom(t *testing.T) {
	cases := []struct {
		state  *State
		expect *Move
		err    bool
	}{
		{
			state: &State{
				Votes: []*Move{},
			},
			expect: nil,
			err:    true,
		}, {
			state: &State{
				Votes: []*Move{{Pass: true}},
			},
			expect: &Move{Pass: true},
			err:    false,
		},
	}
	for _, test := range cases {
		move, err := test.state.Random()
		if err != nil && !test.err {
			t.Errorf("unexpected error %s", err.Error())
		}
		if err == nil && test.err {
			t.Errorf("expected error")
		}
		if !reflect.DeepEqual(move, test.expect) {
			t.Errorf("expected %#v but got %#v",
				test.expect, move)
		}
	}
}

func TestStateEmpty(t *testing.T) {
	cases := []struct {
		state  *State
		expect bool
	}{
		{
			state: &State{
				Votes: []*Move{},
			},
			expect: true,
		}, {
			state: &State{
				Votes: []*Move{{Pass: true}},
			},
			expect: false,
		},
	}
	for _, test := range cases {
		empty := test.state.Empty()
		if empty != test.expect {
			t.Errorf("expected %#v but got %#v",
				test.expect, empty)
		}
	}
}

func TestStateReset(t *testing.T) {
	cases := []struct {
		state  *State
		expect *State
		err    bool
	}{
		{
			state: &State{
				Votes: []*Move{},
			},
			expect: &State{
				Votes: []*Move{},
			},
			err: false,
		}, {
			state: &State{
				Votes: []*Move{{Pass: true}},
			},
			expect: &State{
				Votes: []*Move{},
			},
			err: false,
		},
	}
	for _, test := range cases {
		err := test.state.Reset()
		if err != nil && !test.err {
			t.Errorf("unexpected error %s", err.Error())
		}
		if err == nil && test.err {
			t.Errorf("expected error")
		}
		if !reflect.DeepEqual(test.state, test.expect) {
			t.Errorf("expected %#v but got %#v",
				test.expect, test.state)
		}
	}
}
