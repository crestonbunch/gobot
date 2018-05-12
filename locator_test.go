package gobot

import (
	"reflect"
	"testing"
)

func TestLocatorFind(t *testing.T) {
	cases := []struct {
		desc    string
		user    string
		locator *GameLocator
		games   []*Game
		expect  *Game
		err     bool
	}{
		{
			desc:    "first game with anyone",
			user:    "foo",
			locator: &GameLocator{Auto: true},
			games: []*Game{
				{ID: 1, Players: Players{Anyone: true}},
				{ID: 2, Players: Players{Anyone: true}},
			},
			expect: &Game{ID: 1, Players: Players{Anyone: true}},
		}, {
			desc:    "second game with anyone",
			user:    "foo",
			locator: &GameLocator{Auto: true},
			games: []*Game{
				{ID: 1, Players: Players{Anyone: false}},
				{ID: 2, Players: Players{Anyone: true}},
			},
			expect: &Game{ID: 2, Players: Players{Anyone: true}},
		}, {
			desc:    "first game with players as white",
			user:    "foo",
			locator: &GameLocator{Auto: true},
			games: []*Game{
				{ID: 1, Players: Players{
					White: []string{"foo"}, Black: []string{"bar"}},
				},
				{ID: 2, Players: Players{Anyone: true}},
			},
			expect: &Game{ID: 1, Players: Players{
				White: []string{"foo"}, Black: []string{"bar"}},
			},
		}, {
			desc:    "second game with players as black",
			user:    "bar",
			locator: &GameLocator{Auto: true},
			games: []*Game{
				{ID: 2, Players: Players{Anyone: false}},
				{ID: 1, Players: Players{
					White: []string{"foo"}, Black: []string{"bar"}},
				},
			},
			expect: &Game{ID: 1, Players: Players{
				White: []string{"foo"}, Black: []string{"bar"}},
			},
		}, {
			desc:    "pick your game with anyone",
			user:    "baz",
			locator: &GameLocator{GameID: 2},
			games: []*Game{
				{ID: 2, Players: Players{Anyone: true}},
				{ID: 1, Players: Players{
					White: []string{"foo"}, Black: []string{"bar"}},
				},
			},
			expect: &Game{ID: 2, Players: Players{Anyone: true}},
		}, {
			desc:    "pick your game with players",
			user:    "bar",
			locator: &GameLocator{GameID: 1},
			games: []*Game{
				{ID: 2, Players: Players{Anyone: true}},
				{ID: 1, Players: Players{
					White: []string{"foo"}, Black: []string{"bar"}},
				},
			},
			expect: &Game{ID: 1, Players: Players{
				White: []string{"foo"}, Black: []string{"bar"}},
			},
		}, {
			desc:    "no game",
			user:    "baz",
			locator: &GameLocator{Auto: true},
			games: []*Game{
				{ID: 2, Players: Players{Anyone: false}},
				{ID: 1, Players: Players{
					White: []string{"foo"}, Black: []string{"bar"}},
				},
			},
			expect: nil,
			err:    true,
		}, {
			desc:    "not your game",
			user:    "baz",
			locator: &GameLocator{GameID: 1},
			games: []*Game{
				{ID: 2, Players: Players{Anyone: false}},
				{ID: 1, Players: Players{
					White: []string{"foo"}, Black: []string{"bar"}},
				},
			},
			expect: nil,
			err:    true,
		},
	}
	for _, test := range cases {
		game, err := test.locator.Find(test.games, test.user)
		if test.err && err == nil {
			t.Errorf("%s expected error", test.desc)
		}
		if !test.err && err != nil {
			t.Errorf("%s unexpected error %s", test.desc, err.Error())
		}
		if !reflect.DeepEqual(game, test.expect) {
			t.Errorf("%s\nexpected\n%#v\nbut got\n%#v\n",
				test.desc, test.expect, game)
		}
	}
}
