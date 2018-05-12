package gobot

import (
	"reflect"
	"testing"
)

func TestParseBadCommands(t *testing.T) {
	cases := []string{
		"gobot",
		"gobot asdf",
	}
	for _, test := range cases {
		_, _, err := ParseGameCommand(test)
		if err == nil {
			t.Errorf("expected %s to produce an error", test)
		}
		_, err = ParseInfoCommand(test)
		if err == nil {
			t.Errorf("expected %s to produce an error", test)
		}
		_, err = ParseStartCommand(test)
		if err == nil {
			t.Errorf("expected %s to produce an error", test)
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
			input: "gobot play",
			command: &PlayCommand{
				Players:  Players{Anyone: true},
				Settings: Settings{Vote: true, Timer: 3600},
			},
		}, {
			input: "gobot play <@USER1> <@USER2>",
			command: &PlayCommand{
				Players: Players{
					Black:  []string{"USER1"},
					White:  []string{"USER2"},
					Anyone: false,
				},
				Settings: Settings{},
			},
		}, {
			input:   "gobot play user1 user2 user3",
			command: nil,
			err:     true,
		}, {
			input:   "gobot play user1 user2",
			command: nil,
			err:     true,
		},
	}

	for _, test := range cases {
		actual, err := ParseStartCommand(test.input)
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
		command *MoveCommand
		locator *GameLocator
		err     bool
	}{
		{
			input: "gobot move A4",
			command: &MoveCommand{
				Coordinates: [2]int{3, 0},
			},
			locator: &GameLocator{Auto: true},
		}, {
			input: "gobot move 12 D14",
			command: &MoveCommand{
				Coordinates: [2]int{13, 3},
			},
			locator: &GameLocator{GameID: 12},
		}, {
			input:   "gobot move 12 Z14",
			command: nil,
			locator: nil,
			err:     true,
		}, {
			input:   "gobot move BBZ",
			command: nil,
			locator: nil,
			err:     true,
		},
	}

	for _, test := range cases {
		actual, locator, err := ParseGameCommand(test.input)
		if err == nil && test.err {
			t.Errorf("expected %s to make an error", test.input)
		} else if err != nil && !test.err {
			t.Errorf(
				"%s triggered unexpected error %s", test.input, err.Error(),
			)
		} else if actual == nil && test.command != nil {
			t.Errorf("%s returned unexepected nil", test.input)
		} else {
			if !reflect.DeepEqual(actual, test.command) {
				t.Errorf(
					"%s return\n%#v\nbut expected\n%#v\n",
					test.input, actual, test.command,
				)
			}
			if !reflect.DeepEqual(locator, test.locator) {
				t.Errorf(
					"%s return\n%#v\nbut expected\n%#v\n",
					test.input, locator, test.locator,
				)
			}
		}
	}
}

func TestParsePassCommand(t *testing.T) {
	cases := []struct {
		input   string
		command *PassCommand
		locator *GameLocator
		err     bool
	}{
		{
			input:   "gobot pass",
			command: &PassCommand{},
			locator: &GameLocator{Auto: true},
		}, {
			input:   "gobot pass 13",
			command: &PassCommand{},
			locator: &GameLocator{GameID: 13},
		}, {
			input:   "gobot pass A",
			command: nil,
			locator: nil,
			err:     true,
		},
	}

	for _, test := range cases {
		actual, locator, err := ParseGameCommand(test.input)
		if err == nil && test.err {
			t.Errorf("expected %s to make an error", test.input)
		} else if err != nil && !test.err {
			t.Errorf(
				"%s triggered unexpected error %s", test.input, err.Error(),
			)
		} else if actual == nil && test.command != nil {
			t.Errorf("%s returned unexepected nil", test.input)
		} else {
			if !reflect.DeepEqual(actual, test.command) {
				t.Errorf(
					"%s return\n%#v\nbut expected\n%#v\n",
					test.input, actual, test.command,
				)
			}
			if !reflect.DeepEqual(locator, test.locator) {
				t.Errorf(
					"%s return\n%#v\nbut expected\n%#v\n",
					test.input, locator, test.locator,
				)
			}
		}
	}
}

func TestParseShowCommand(t *testing.T) {
	cases := []struct {
		input   string
		command *ShowCommand
		locator *GameLocator
		err     bool
	}{
		{
			input:   "gobot show",
			command: &ShowCommand{},
			locator: &GameLocator{Auto: true},
		}, {
			input:   "gobot show 13",
			command: &ShowCommand{},
			locator: &GameLocator{GameID: 13},
		}, {
			input:   "gobot show A",
			command: nil,
			locator: nil,
			err:     true,
		},
	}

	for _, test := range cases {
		parsed, locator, err := ParseGameCommand(test.input)
		actual := parsed.(*ShowCommand)
		if err == nil && test.err {
			t.Errorf("expected %s to make an error", test.input)
		} else if err != nil && !test.err {
			t.Errorf(
				"%s triggered unexpected error %s", test.input, err.Error(),
			)
		} else if actual == nil && test.command != nil {
			t.Errorf("%s returned unexepected nil", test.input)
		} else {
			if !reflect.DeepEqual(actual, test.command) {
				t.Errorf(
					"%s return\n%#v\nbut expected\n%#v\n",
					test.input, actual, test.command,
				)
			}
			if !reflect.DeepEqual(locator, test.locator) {
				t.Errorf(
					"%s return\n%#v\nbut expected\n%#v\n",
					test.input, locator, test.locator,
				)
			}
		}
	}
}

func TestParseScoreCommand(t *testing.T) {
	cases := []struct {
		input   string
		command *ScoreCommand
		locator *GameLocator
		err     bool
	}{
		{
			input:   "gobot score",
			command: &ScoreCommand{},
			locator: &GameLocator{Auto: true},
		}, {
			input:   "gobot score 13",
			command: &ScoreCommand{},
			locator: &GameLocator{GameID: 13},
		}, {
			input:   "gobot score A",
			command: nil,
			locator: nil,
			err:     true,
		},
	}

	for _, test := range cases {
		parsed, locator, err := ParseGameCommand(test.input)
		actual := parsed.(*ScoreCommand)
		if err == nil && test.err {
			t.Errorf("expected %s to make an error", test.input)
		} else if err != nil && !test.err {
			t.Errorf(
				"%s triggered unexpected error %s", test.input, err.Error(),
			)
		} else if actual == nil && test.command != nil {
			t.Errorf("%s returned unexepected nil", test.input)
		} else {
			if !reflect.DeepEqual(actual, test.command) {
				t.Errorf(
					"%s return\n%#v\nbut expected\n%#v\n",
					test.input, actual, test.command,
				)
			}
			if !reflect.DeepEqual(locator, test.locator) {
				t.Errorf(
					"%s return\n%#v\nbut expected\n%#v\n",
					test.input, locator, test.locator,
				)
			}
		}
	}
}

func TestParseListCommand(t *testing.T) {
	cases := []struct {
		input  string
		expect *ListCommand
		err    bool
	}{
		{
			input:  "gobot list",
			expect: &ListCommand{},
		}, {
			input: "gobot list all",
			expect: &ListCommand{
				All: true,
			},
		}, {
			input:  "gobot list asdf",
			expect: nil,
			err:    true,
		},
	}

	for _, test := range cases {
		actual, err := ParseInfoCommand(test.input)
		if err == nil && test.err {
			t.Errorf("expected %s to make an error", test.input)
		} else if err != nil && !test.err {
			t.Errorf(
				"%s triggered unexpected error %s", test.input, err.Error(),
			)
		} else if actual == nil && test.expect != nil {
			t.Errorf("%s returned unexepected nil", test.input)
		} else {
			if !reflect.DeepEqual(actual, test.expect) {
				t.Errorf(
					"%s return\n%#v\nbut expected\n%#v\n",
					test.input, actual, test.expect,
				)
			}
		}
	}
}
