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
		_, err := Parse(test)
		if err == nil {
			t.Errorf("expected %s to produce an error", test)
		}
	}
}

func TestParsePlayCommand(t *testing.T) {
	cases := []struct {
		command string
		expect  *PlayCommand
		err     bool
	}{
		{
			command: "gobot play",
			expect: &PlayCommand{
				Players:  Players{Anyone: true},
				Settings: Settings{Vote: true, Timer: 3600},
			},
		}, {
			command: "gobot play user1 user2",
			expect: &PlayCommand{
				Players: Players{
					Black:  []string{"user1"},
					White:  []string{"user2"},
					Anyone: false,
				},
				Settings: Settings{},
			},
		}, {
			command: "gobot play user1 user2 user3",
			expect:  nil,
			err:     true,
		},
	}

	for _, test := range cases {
		parsed, err := Parse(test.command)
		actual := parsed.(*PlayCommand)
		if err == nil && test.err {
			t.Errorf("expected error for %s", test.command)
		} else if err != nil && !test.err {
			t.Errorf("unexpected error %s for %s", err.Error(), test.command)
		} else if actual != nil {
			if !reflect.DeepEqual(actual, test.expect) {
				t.Errorf(
					"expected:\n%#v\nbut got:\n%#v\n",
					*test.expect,
					*actual,
				)
			}
		} else if actual == nil {
			if test.expect != nil {
				t.Errorf("encountered unexpected nil")
			}
		}
	}
}

func TestParseMoveCommand(t *testing.T) {
	cases := []struct {
		command string
		expect  *MoveCommand
		err     bool
	}{
		{
			command: "gobot move A4",
			expect: &MoveCommand{
				GameID:      -1,
				Coordinates: [2]int{3, 0},
			},
		}, {
			command: "gobot move 12 D14",
			expect: &MoveCommand{
				GameID:      12,
				Coordinates: [2]int{13, 3},
			},
		}, {
			command: "gobot move 12 Z14",
			expect:  nil,
			err:     true,
		}, {
			command: "gobot move BBZ",
			expect:  nil,
			err:     true,
		},
	}

	for _, test := range cases {
		parsed, err := Parse(test.command)
		actual := parsed.(*MoveCommand)
		if err == nil && test.err {
			t.Errorf("expected %s to make an error", test.command)
		} else if err != nil && !test.err {
			t.Errorf(
				"%s triggered unexpected error %s", test.command, err.Error(),
			)
		} else if actual == nil && test.expect != nil {
			t.Errorf("%s returned unexepected nil", test.command)
		} else {
			if !reflect.DeepEqual(actual, test.expect) {
				t.Errorf(
					"%s return\n%#v\nbut expected\n%#v\n",
					test.command, actual, test.expect,
				)
			}
		}
	}
}

func TestParsePassCommand(t *testing.T) {
	cases := []struct {
		command string
		expect  *PassCommand
		err     bool
	}{
		{
			command: "gobot pass",
			expect: &PassCommand{
				GameID: -1,
			},
		}, {
			command: "gobot pass 13",
			expect: &PassCommand{
				GameID: 13,
			},
		}, {
			command: "gobot pass A",
			expect:  nil,
			err:     true,
		},
	}

	for _, test := range cases {
		parsed, err := Parse(test.command)
		actual := parsed.(*PassCommand)
		if err == nil && test.err {
			t.Errorf("expected %s to make an error", test.command)
		} else if err != nil && !test.err {
			t.Errorf(
				"%s triggered unexpected error %s", test.command, err.Error(),
			)
		} else if actual == nil && test.expect != nil {
			t.Errorf("%s returned unexepected nil", test.command)
		} else {
			if !reflect.DeepEqual(actual, test.expect) {
				t.Errorf(
					"%s return\n%#v\nbut expected\n%#v\n",
					test.command, actual, test.expect,
				)
			}
		}
	}
}

func TestParseShowCommand(t *testing.T) {
	cases := []struct {
		command string
		expect  *ShowCommand
		err     bool
	}{
		{
			command: "gobot show",
			expect: &ShowCommand{
				GameID: -1,
			},
		}, {
			command: "gobot show 13",
			expect: &ShowCommand{
				GameID: 13,
			},
		}, {
			command: "gobot show A",
			expect:  nil,
			err:     true,
		},
	}

	for _, test := range cases {
		parsed, err := Parse(test.command)
		actual := parsed.(*ShowCommand)
		if err == nil && test.err {
			t.Errorf("expected %s to make an error", test.command)
		} else if err != nil && !test.err {
			t.Errorf(
				"%s triggered unexpected error %s", test.command, err.Error(),
			)
		} else if actual == nil && test.expect != nil {
			t.Errorf("%s returned unexepected nil", test.command)
		} else {
			if !reflect.DeepEqual(actual, test.expect) {
				t.Errorf(
					"%s return\n%#v\nbut expected\n%#v\n",
					test.command, actual, test.expect,
				)
			}
		}
	}
}

func TestParseScoreCommand(t *testing.T) {
	cases := []struct {
		command string
		expect  *ScoreCommand
		err     bool
	}{
		{
			command: "gobot score",
			expect: &ScoreCommand{
				GameID: -1,
			},
		}, {
			command: "gobot score 13",
			expect: &ScoreCommand{
				GameID: 13,
			},
		}, {
			command: "gobot score A",
			expect:  nil,
			err:     true,
		},
	}

	for _, test := range cases {
		parsed, err := Parse(test.command)
		actual := parsed.(*ScoreCommand)
		if err == nil && test.err {
			t.Errorf("expected %s to make an error", test.command)
		} else if err != nil && !test.err {
			t.Errorf(
				"%s triggered unexpected error %s", test.command, err.Error(),
			)
		} else if actual == nil && test.expect != nil {
			t.Errorf("%s returned unexepected nil", test.command)
		} else {
			if !reflect.DeepEqual(actual, test.expect) {
				t.Errorf(
					"%s return\n%#v\nbut expected\n%#v\n",
					test.command, actual, test.expect,
				)
			}
		}
	}
}

func TestParseListCommand(t *testing.T) {
	cases := []struct {
		command string
		expect  *ListCommand
		err     bool
	}{
		{
			command: "gobot list",
			expect:  &ListCommand{},
		}, {
			command: "gobot list all",
			expect: &ListCommand{
				All: true,
			},
		}, {
			command: "gobot list asdf",
			expect:  nil,
			err:     true,
		},
	}

	for _, test := range cases {
		parsed, err := Parse(test.command)
		actual := parsed.(*ListCommand)
		if err == nil && test.err {
			t.Errorf("expected %s to make an error", test.command)
		} else if err != nil && !test.err {
			t.Errorf(
				"%s triggered unexpected error %s", test.command, err.Error(),
			)
		} else if actual == nil && test.expect != nil {
			t.Errorf("%s returned unexepected nil", test.command)
		} else {
			if !reflect.DeepEqual(actual, test.expect) {
				t.Errorf(
					"%s return\n%#v\nbut expected\n%#v\n",
					test.command, actual, test.expect,
				)
			}
		}
	}
}
