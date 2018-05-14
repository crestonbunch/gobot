package gobot

import (
	"reflect"
	"testing"
	"time"
)

func TestParseBadCommands(t *testing.T) {
	cases := []string{
		"",
		"asdf",
	}
	for _, test := range cases {
		_, err := ParseCommand(test)
		if err == nil {
			t.Errorf("expected %s to produce an error", test)
		}
	}
}

func TestParseStartCommand(t *testing.T) {
	cases := []struct {
		input   string
		command *StartCommand
		err     bool
	}{
		{
			input: "start",
			command: &StartCommand{
				Players:  Players{Anyone: true},
				Settings: Settings{Vote: true, Timer: 3600 * time.Second},
			},
		}, {
			input: "start USER1 USER2",
			command: &StartCommand{
				Players: Players{
					Black:  []string{"USER1"},
					White:  []string{"USER2"},
					Anyone: false,
				},
				Settings: Settings{},
			},
		}, {
			input:   "start user1 user2 user3",
			command: nil,
			err:     true,
		},
	}

	for _, test := range cases {
		actual, err := ParseCommand(test.input)
		if err == nil && test.err {
			t.Errorf("expected error for %s", test.input)
		} else if err != nil && !test.err {
			t.Errorf("unexpected error %s for %s", err.Error(), test.input)
		} else if actual != nil {
			if !reflect.DeepEqual(actual, test.command) {
				t.Errorf(
					"%s\nexpected:\n%#v\nbut got:\n%#v\n",
					test.input,
					test.command,
					actual,
				)
			}
		} else if actual == nil {
			if test.command != nil {
				t.Errorf("encountered unexpected nil")
			}
		}
	}
}

func TestParseMoveCommand(t *testing.T) {
	cases := []struct {
		input   string
		command Command
		err     bool
	}{
		{
			input: "move A4",
			command: &MoveCommand{
				Coordinates: [2]int{3, 0},
				Locator:     GameLocator{Auto: true},
			},
		}, {
			input: "move 12 D14",
			command: &MoveCommand{
				Coordinates: [2]int{13, 3},
				Locator:     GameLocator{GameID: 12},
			},
		}, {
			input: "move pass",
			command: &MoveCommand{
				Pass:    true,
				Locator: GameLocator{Auto: true},
			},
		}, {
			input: "move 12 pass",
			command: &MoveCommand{
				Pass:    true,
				Locator: GameLocator{GameID: 12},
			},
		}, {
			input:   "move 12 Z14",
			command: nil,
			err:     true,
		}, {
			input:   "move BBZ",
			command: nil,
			err:     true,
		},
	}

	for _, test := range cases {
		actual, err := ParseCommand(test.input)
		if err == nil && test.err {
			t.Errorf("expected %s to make an error", test.input)
		} else if err != nil && !test.err {
			t.Errorf(
				"%s triggered unexpected error %s", test.input, err.Error(),
			)
		} else if actual == nil && test.command != nil {
			t.Errorf("%s returned unexepected nil", test.input)
		} else if actual != nil && test.command != nil {
			if !reflect.DeepEqual(actual, test.command) {
				t.Errorf(
					"%s returned\n%#v\nbut expected\n%#v\n",
					test.input, actual, test.command,
				)
			}
		}
	}
}

func TestParseVoteCommand(t *testing.T) {
	cases := []struct {
		input   string
		command *VoteCommand
		err     bool
	}{
		{
			input: "vote A4",
			command: &VoteCommand{
				Coordinates: [2]int{3, 0},
				Locator:     GameLocator{Auto: true},
			},
		}, {
			input: "vote 12 D14",
			command: &VoteCommand{
				Coordinates: [2]int{13, 3},
				Locator:     GameLocator{GameID: 12},
			},
		}, {
			input: "vote pass",
			command: &VoteCommand{
				Pass:    true,
				Locator: GameLocator{Auto: true},
			},
		}, {
			input: "vote 12 pass",
			command: &VoteCommand{
				Pass:    true,
				Locator: GameLocator{GameID: 12},
			},
		}, {
			input:   "vote 12 Z14",
			command: nil,
			err:     true,
		}, {
			input:   "vote BBZ",
			command: nil,
			err:     true,
		},
	}

	for _, test := range cases {
		actual, err := ParseCommand(test.input)
		if err == nil && test.err {
			t.Errorf("expected %s to make an error", test.input)
		} else if err != nil && !test.err {
			t.Errorf(
				"%s triggered unexpected error %s", test.input, err.Error(),
			)
		} else if actual == nil && test.command != nil {
			t.Errorf("%s returned unexepected nil", test.input)
		} else if actual != nil && test.command != nil {
			if !reflect.DeepEqual(actual, test.command) {
				t.Errorf(
					"%s\n%#v\nbut expected\n%#v\n",
					test.input, actual, test.command,
				)
			}
		}
	}
}
