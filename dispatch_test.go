package gobot

import (
	"reflect"
	"testing"
)

func TestDispatch(t *testing.T) {
	voteGame := NewGame(1, Players{Anyone: true}, Settings{Vote: true})
	moveGame := NewGame(1, Players{Anyone: true}, Settings{Vote: false})
	cases := []struct {
		desc     string
		game     *Game
		player   string
		command  Command
		response *Response
		err      bool
	}{
		{
			desc:     "vote for move",
			game:     voteGame,
			player:   "foo",
			command:  &MoveCommand{Coordinates: [2]int{2, 2}},
			response: NewTextResponse("thanks for voting"),
			err:      false,
		}, {
			desc:     "vote for pass",
			game:     voteGame,
			player:   "foo",
			command:  &PassCommand{},
			response: NewTextResponse("thanks for voting"),
			err:      false,
		}, {
			desc:     "move",
			game:     moveGame,
			player:   "foo",
			command:  &MoveCommand{Coordinates: [2]int{2, 2}},
			response: NewGameResponse(moveGame),
			err:      false,
		}, {
			desc:     "not implemented",
			game:     moveGame,
			player:   "foo",
			command:  &ListCommand{},
			response: nil,
			err:      true,
		},
	}
	for _, test := range cases {
		response, err := Dispatch(test.game, test.player, test.command)
		if err != nil && !test.err {
			t.Errorf("%s unexpected error %s", test.desc, err.Error())
		}
		if err == nil && test.err {
			t.Errorf("%s expected error", test.desc)
		}
		if !reflect.DeepEqual(response, test.response) {
			t.Errorf("%s wrong response", test.desc)
		}
	}
}
