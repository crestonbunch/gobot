package gobot

import (
	"testing"
	"time"
)

func TestStartCommand(t *testing.T) {
	cases := []struct {
		command        *StartCommand
		player         string
		gameStore      GameStore
		voteStore      VoteStore
		schedulerStore SchedulerStore
		err            bool
	}{
		{
			command: &StartCommand{
				Players:  Players{Anyone: true},
				Settings: Settings{Vote: false},
			},
			player:         "foo",
			gameStore:      map[int]*Game{},
			voteStore:      map[int]map[string]*Vote{},
			schedulerStore: NewSchedulerStore(make(chan *Response)),
		},
	}
	for _, test := range cases {
		response, err := test.command.Execute(
			test.player, test.gameStore, test.voteStore, test.schedulerStore,
		)
		if err == nil && test.err {
			t.Errorf("unexpected error for %#v", test)
		}
		if err != nil && !test.err {
			t.Errorf("expected error for %#v", test)
		}
		if response != nil && response.Game == nil {
			t.Errorf("unexpected nil game for %#v", test)
		}
	}
}

func TestMoveCommand(t *testing.T) {
	cases := []struct {
		desc           string
		command        *MoveCommand
		player         string
		gameStore      GameStore
		voteStore      VoteStore
		schedulerStore SchedulerStore
		err            bool
	}{
		{
			desc: "anyone move auto locate",
			command: &MoveCommand{
				Coordinates: Coords{0, 0},
				Locator:     GameLocator{Auto: true},
			},
			player: "foo",
			gameStore: map[int]*Game{
				1: {
					ID: 1,
					History: []Board{{
						[]Stone{EmptyStone, EmptyStone},
						[]Stone{EmptyStone, EmptyStone},
					}},
					Next:     BlackStone,
					Players:  Players{Anyone: true},
					Settings: Settings{Vote: false},
				},
			},
			err: false,
		}, {
			desc: "anyone move game id locate",
			command: &MoveCommand{
				Coordinates: Coords{0, 0},
				Locator:     GameLocator{GameID: 1},
			},
			player: "foo",
			gameStore: map[int]*Game{
				1: {
					ID: 1,
					History: []Board{{
						[]Stone{EmptyStone, EmptyStone},
						[]Stone{EmptyStone, EmptyStone},
					}},
					Next:     BlackStone,
					Players:  Players{Anyone: true},
					Settings: Settings{Vote: false},
				},
			},
			err: false,
		}, {
			desc: "require foo to move auto locate",
			command: &MoveCommand{
				Coordinates: Coords{0, 0},
				Locator:     GameLocator{GameID: 1},
			},
			player: "foo",
			gameStore: map[int]*Game{
				1: {
					ID: 1,
					History: []Board{{
						[]Stone{EmptyStone, EmptyStone},
						[]Stone{EmptyStone, EmptyStone},
					}},
					Next: BlackStone,
					Players: Players{
						Black: []string{"foo"}, White: []string{"bar"},
					},
					Settings: Settings{Vote: false},
				},
			},
			err: false,
		}, {
			desc: "wrong user to move auto locate",
			command: &MoveCommand{
				Coordinates: Coords{0, 0},
				Locator:     GameLocator{GameID: 1},
			},
			player: "bar",
			gameStore: map[int]*Game{
				1: {
					ID: 1,
					History: []Board{{
						[]Stone{EmptyStone, EmptyStone},
						[]Stone{EmptyStone, EmptyStone},
					}},
					Next: BlackStone,
					Players: Players{
						Black: []string{"foo"}, White: []string{"bar"},
					},
					Settings: Settings{Vote: false},
				},
			},
			err: true,
		}, {
			desc: "requires voting",
			command: &MoveCommand{
				Coordinates: Coords{0, 0},
				Locator:     GameLocator{GameID: 1},
			},
			player: "foo",
			gameStore: map[int]*Game{
				1: {
					ID: 1,
					History: []Board{{
						[]Stone{EmptyStone, EmptyStone},
						[]Stone{EmptyStone, EmptyStone},
					}},
					Next: BlackStone,
					Players: Players{
						Black: []string{"foo"}, White: []string{"bar"},
					},
					Settings: Settings{Vote: true},
				},
			},
			err: true,
		}, {
			desc: "passing",
			command: &MoveCommand{
				Pass:    true,
				Locator: GameLocator{GameID: 1},
			},
			player: "foo",
			gameStore: map[int]*Game{
				1: {
					ID: 1,
					History: []Board{{
						[]Stone{EmptyStone, EmptyStone},
						[]Stone{EmptyStone, EmptyStone},
					}},
					Next: BlackStone,
					Players: Players{
						Black: []string{"foo"}, White: []string{"bar"},
					},
					Settings: Settings{Vote: false},
				},
			},
			err: false,
		}, {
			desc: "passing not allowed",
			command: &MoveCommand{
				Pass:    true,
				Locator: GameLocator{GameID: 1},
			},
			player: "bar",
			gameStore: map[int]*Game{
				1: {
					ID: 1,
					History: []Board{{
						[]Stone{EmptyStone, EmptyStone},
						[]Stone{EmptyStone, EmptyStone},
					}},
					Next: BlackStone,
					Players: Players{
						Black: []string{"foo"}, White: []string{"bar"},
					},
					Settings: Settings{Vote: false},
				},
			},
			err: true,
		}, {
			desc: "locator error",
			command: &MoveCommand{
				Coordinates: Coords{0, 0},
				Locator:     GameLocator{GameID: 1},
			},
			player:    "bar",
			gameStore: map[int]*Game{},
			err:       true,
		},
	}
	for _, test := range cases {
		response, err := test.command.Execute(
			test.player, test.gameStore, test.voteStore, test.schedulerStore,
		)
		if err == nil && test.err {
			t.Errorf("expected error for %s", test.desc)
		}
		if err != nil && !test.err {
			t.Errorf("unexpected error %s for %s", err.Error(), test.desc)
		}
		if response != nil && response.Game == nil {
			t.Errorf("unexpected nil game for %s", test.desc)
		}
	}
}

func TestVoteCommand(t *testing.T) {
	cases := []struct {
		desc           string
		command        *VoteCommand
		player         string
		gameStore      GameStore
		voteStore      VoteStore
		schedulerStore SchedulerStore
		err            bool
	}{
		{
			desc: "anyone vote auto locate",
			command: &VoteCommand{
				Coordinates: Coords{0, 0},
				Locator:     GameLocator{Auto: true},
			},
			player: "foo",
			gameStore: map[int]*Game{
				1: {
					ID: 1,
					History: []Board{{
						[]Stone{EmptyStone, EmptyStone},
						[]Stone{EmptyStone, EmptyStone},
					}},
					Next:     BlackStone,
					Players:  Players{Anyone: true},
					Settings: Settings{Vote: true},
				},
			},
			voteStore: map[int]map[string]*Vote{},
			err:       false,
		}, {
			desc: "anyone vote game id locate",
			command: &VoteCommand{
				Coordinates: Coords{0, 0},
				Locator:     GameLocator{GameID: 1},
			},
			player: "foo",
			gameStore: map[int]*Game{
				1: {
					ID: 1,
					History: []Board{{
						[]Stone{EmptyStone, EmptyStone},
						[]Stone{EmptyStone, EmptyStone},
					}},
					Next:     BlackStone,
					Players:  Players{Anyone: true},
					Settings: Settings{Vote: true},
				},
			},
			voteStore: map[int]map[string]*Vote{},
			err:       false,
		}, {
			desc: "require foo to vote auto locate",
			command: &VoteCommand{
				Coordinates: Coords{0, 0},
				Locator:     GameLocator{GameID: 1},
			},
			player: "foo",
			gameStore: map[int]*Game{
				1: {
					ID: 1,
					History: []Board{{
						[]Stone{EmptyStone, EmptyStone},
						[]Stone{EmptyStone, EmptyStone},
					}},
					Next: BlackStone,
					Players: Players{
						Black: []string{"foo"}, White: []string{"bar"},
					},
					Settings: Settings{Vote: true},
				},
			},
			voteStore: map[int]map[string]*Vote{},
			err:       false,
		}, {
			desc: "wrong user to vote auto locate",
			command: &VoteCommand{
				Coordinates: Coords{0, 0},
				Locator:     GameLocator{GameID: 1},
			},
			player: "bar",
			gameStore: map[int]*Game{
				1: {
					ID: 1,
					History: []Board{{
						[]Stone{EmptyStone, EmptyStone},
						[]Stone{EmptyStone, EmptyStone},
					}},
					Next: BlackStone,
					Players: Players{
						Black: []string{"foo"}, White: []string{"bar"},
					},
					Settings: Settings{Vote: true},
				},
			},
			voteStore: map[int]map[string]*Vote{},
			err:       true,
		}, {
			desc: "voting not allowed",
			command: &VoteCommand{
				Coordinates: Coords{0, 0},
				Locator:     GameLocator{GameID: 1},
			},
			player: "foo",
			gameStore: map[int]*Game{
				1: {
					ID: 1,
					History: []Board{{
						[]Stone{EmptyStone, EmptyStone},
						[]Stone{EmptyStone, EmptyStone},
					}},
					Next: BlackStone,
					Players: Players{
						Black: []string{"foo"}, White: []string{"bar"},
					},
					Settings: Settings{Vote: false},
				},
			},
			voteStore: map[int]map[string]*Vote{},
			err:       true,
		}, {
			desc: "passing",
			command: &VoteCommand{
				Pass:    true,
				Locator: GameLocator{GameID: 1},
			},
			player: "foo",
			gameStore: map[int]*Game{
				1: {
					ID: 1,
					History: []Board{{
						[]Stone{EmptyStone, EmptyStone},
						[]Stone{EmptyStone, EmptyStone},
					}},
					Next: BlackStone,
					Players: Players{
						Black: []string{"foo"}, White: []string{"bar"},
					},
					Settings: Settings{Vote: true},
				},
			},
			voteStore: map[int]map[string]*Vote{},
			err:       false,
		}, {
			desc: "passing not allowed",
			command: &VoteCommand{
				Pass:    true,
				Locator: GameLocator{GameID: 1},
			},
			player: "bar",
			gameStore: map[int]*Game{
				1: {
					ID: 1,
					History: []Board{{
						[]Stone{EmptyStone, EmptyStone},
						[]Stone{EmptyStone, EmptyStone},
					}},
					Next: BlackStone,
					Players: Players{
						Black: []string{"foo"}, White: []string{"bar"},
					},
					Settings: Settings{Vote: true},
				},
			},
			voteStore: map[int]map[string]*Vote{},
			err:       true,
		}, {
			desc: "locator error",
			command: &VoteCommand{
				Coordinates: Coords{0, 0},
				Locator:     GameLocator{GameID: 1},
			},
			player:    "bar",
			gameStore: map[int]*Game{},
			err:       true,
		},
	}
	for _, test := range cases {
		response, err := test.command.Execute(
			test.player, test.gameStore, test.voteStore, test.schedulerStore,
		)
		if err == nil && test.err {
			t.Errorf("expected error for %s", test.desc)
		}
		if err != nil && !test.err {
			t.Errorf("unexpected error %s for %s", err.Error(), test.desc)
		}
		if response != nil && response.Text == "" {
			t.Errorf("unexpected empty response for %s", test.desc)
		}
	}
}

func TestPlayCommand(t *testing.T) {
	cases := []struct {
		desc           string
		command        *PlayCommand
		player         string
		gameStore      GameStore
		voteStore      VoteStore
		schedulerStore SchedulerStore
		err            bool
	}{
		{
			desc: "anyone play auto locate",
			command: &PlayCommand{
				Locator: GameLocator{Auto: true},
			},
			player: "foo",
			gameStore: map[int]*Game{
				1: {
					ID: 1,
					History: []Board{{
						[]Stone{EmptyStone, EmptyStone},
						[]Stone{EmptyStone, EmptyStone},
					}},
					Next:     BlackStone,
					Players:  Players{Anyone: true},
					Settings: Settings{Vote: true},
				},
			},
			voteStore: map[int]map[string]*Vote{
				1: {"foo": {Move: Coords{0, 0}}},
			},
			schedulerStore: SchedulerStore{
				Schedulers: map[int]*Scheduler{
					1: {
						VoteTimer: time.NewTimer(100 * time.Second),
					},
				},
				Responses: make(chan *Response),
			},
			err: false,
		}, {
			desc: "no votes cast",
			command: &PlayCommand{
				Locator: GameLocator{Auto: true},
			},
			player: "foo",
			gameStore: map[int]*Game{
				1: {
					ID: 1,
					History: []Board{{
						[]Stone{EmptyStone, EmptyStone},
						[]Stone{EmptyStone, EmptyStone},
					}},
					Next:     BlackStone,
					Players:  Players{Anyone: true},
					Settings: Settings{Vote: true},
				},
			},
			voteStore: map[int]map[string]*Vote{},
			schedulerStore: SchedulerStore{
				Schedulers: map[int]*Scheduler{
					1: {
						VoteTimer: time.NewTimer(100 * time.Second),
					},
				},
				Responses: make(chan *Response),
			},
			err: true,
		}, {
			desc: "not authorized to move",
			command: &PlayCommand{
				Locator: GameLocator{GameID: 1},
			},
			player: "foo",
			gameStore: map[int]*Game{
				1: {
					ID: 1,
					History: []Board{{
						[]Stone{EmptyStone, EmptyStone},
						[]Stone{EmptyStone, EmptyStone},
					}},
					Next: WhiteStone,
					Players: Players{
						Black: []string{"foo"},
						White: []string{"bar"},
					},
					Settings: Settings{Vote: true},
				},
			},
			voteStore: map[int]map[string]*Vote{
				1: {"foo": {Move: Coords{0, 0}}},
			},
			schedulerStore: SchedulerStore{
				Schedulers: map[int]*Scheduler{
					1: {
						VoteTimer: time.NewTimer(100 * time.Second),
					},
				},
				Responses: make(chan *Response),
			},
			err: true,
		}, {
			desc: "vote to pass",
			command: &PlayCommand{
				Locator: GameLocator{GameID: 1},
			},
			player: "bar",
			gameStore: map[int]*Game{
				1: {
					ID: 1,
					History: []Board{{
						[]Stone{EmptyStone, EmptyStone},
						[]Stone{EmptyStone, EmptyStone},
					}},
					Next: WhiteStone,
					Players: Players{
						Black: []string{"foo"},
						White: []string{"bar"},
					},
					Settings: Settings{Vote: true},
				},
			},
			voteStore: map[int]map[string]*Vote{
				1: {"bar": {Pass: true}},
			},
			schedulerStore: SchedulerStore{
				Schedulers: map[int]*Scheduler{
					1: {
						VoteTimer: time.NewTimer(100 * time.Second),
					},
				},
				Responses: make(chan *Response),
			},
			err: false,
		},
	}
	for _, test := range cases {
		response, err := test.command.Execute(
			test.player, test.gameStore, test.voteStore, test.schedulerStore,
		)
		if err == nil && test.err {
			t.Errorf("expected error for %s", test.desc)
		}
		if err != nil && !test.err {
			t.Errorf("unexpected error %s for %s", err.Error(), test.desc)
		}
		if response != nil && response.Game == nil {
			t.Errorf("unexpected nil game for %#v", test)
		}
	}
}
