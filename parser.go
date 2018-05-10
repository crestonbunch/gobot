package gobot

import (
	"fmt"
	"strings"
)

const (
	// PlayToken indicates a play command
	PlayToken string = "play"
)

// PlayCommand initiates a new game
type PlayCommand struct {
	Players  *Players
	Settings *Settings
}

// Parse a string into a gobot command
func Parse(command string) (interface{}, error) {
	args := strings.Split(command, " ")
	if len(args) < 1 {
		return nil, fmt.Errorf("could not understand %s", command)
	}
	switch args[1] {
	case PlayToken:
		return parsePlayCommand(args[2:])
	default:
		return nil, fmt.Errorf("could not understand %s", command)
	}
}

func parsePlayCommand(players []string) (*PlayCommand, error) {
	switch len(players) {
	// Two players only
	case 2:
		return &PlayCommand{
			Players: &Players{
				Black: []string{players[0]},
				White: []string{players[1]},
			},
			Settings: &Settings{
				Vote: false,
			},
		}, nil
	// Allow anyone to vote for moves
	case 0:
		return &PlayCommand{
			Players: &Players{
				Anyone: true,
			},
			Settings: &Settings{
				Vote: true,
			},
		}, nil
	default:
		return nil, fmt.Errorf("incorrect number of players")
	}

}
