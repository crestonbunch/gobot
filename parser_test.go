package gobot_test

import (
	"reflect"
	"testing"

	. "github.com/crestonbunch/gobot"
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
				Anyone: true,
			},
		}, {
			input: "start USER1 USER2",
			command: &StartCommand{
				Black:  []string{"USER1"},
				White:  []string{"USER2"},
				Anyone: false,
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
				Move:    &Move{Coords: [2]int{3, 0}},
				Locator: Locator{Auto: true},
			},
		}, {
			input: "move 12 D14",
			command: &MoveCommand{
				Move:    &Move{Coords: [2]int{13, 3}},
				Locator: Locator{ID: 12},
			},
		}, {
			input: "move pass",
			command: &MoveCommand{
				Move:    &Move{Pass: true},
				Locator: Locator{Auto: true},
			},
		}, {
			input: "move 12 pass",
			command: &MoveCommand{
				Move:    &Move{Pass: true},
				Locator: Locator{ID: 12},
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
				Move:    &Move{Coords: [2]int{3, 0}},
				Locator: Locator{Auto: true},
			},
		}, {
			input: "vote 12 D14",
			command: &VoteCommand{
				Move:    &Move{Coords: [2]int{13, 3}},
				Locator: Locator{ID: 12},
			},
		}, {
			input: "vote pass",
			command: &VoteCommand{
				Move:    &Move{Pass: true},
				Locator: Locator{Auto: true},
			},
		}, {
			input: "vote 12 pass",
			command: &VoteCommand{
				Move:    &Move{Pass: true},
				Locator: Locator{ID: 12},
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

func TestParsePlayCommand(t *testing.T) {
	cases := []struct {
		input   string
		command *PlayCommand
		err     bool
	}{
		{
			input: "play",
			command: &PlayCommand{
				Locator: Locator{Auto: true},
			},
		}, {
			input: "play 12",
			command: &PlayCommand{
				Locator: Locator{ID: 12},
			},
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

func TestParseShowCommand(t *testing.T) {
	cases := []struct {
		input   string
		command *ShowCommand
		err     bool
	}{
		{
			input: "show",
			command: &ShowCommand{
				Locator: Locator{Auto: true},
			},
		}, {
			input: "show 12",
			command: &ShowCommand{
				Locator: Locator{ID: 12},
			},
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

func TestParseListCommand(t *testing.T) {
	cases := []struct {
		input   string
		command *ListCommand
		err     bool
	}{
		{
			input:   "list",
			command: &ListCommand{},
		}, {
			input:   "list all",
			command: &ListCommand{All: true},
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
