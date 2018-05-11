package gobot

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// PlayToken indicates a play command
	PlayToken string = "play"
	// MoveToken indicates a move command
	MoveToken string = "move"
	// PassToken indicates a pass command
	PassToken string = "pass"
	// ShowToken indicates a show command
	ShowToken string = "show"
	// ScoreToken indicates a score command
	ScoreToken string = "score"
	// ListToken indicates a list command
	ListToken string = "list"
)

// A GameCommand is any command that corresponds to a single game
type GameCommand interface {
	Game() int
}

// PlayCommand initiates a new game
type PlayCommand struct {
	Players  Players
	Settings Settings
}

// MoveCommand plays a move in a game
type MoveCommand struct {
	GameID      int
	Coordinates [2]int
}

// Game returns the GameID target of the command
func (c *MoveCommand) Game() int {
	return c.GameID
}

// PassCommand plays a pass in a game
type PassCommand struct {
	GameID int
}

// Game returns the GameID target of the command
func (c *PassCommand) Game() int {
	return c.GameID
}

// ShowCommand shows the current game state
type ShowCommand struct {
	GameID int
}

// Game returns the GameID target of the command
func (c *ShowCommand) Game() int {
	return c.GameID
}

// ScoreCommand scores the current game state
type ScoreCommand struct {
	GameID int
}

// Game returns the GameID target of the command
func (c *ScoreCommand) Game() int {
	return c.GameID
}

// ListCommand lists games
type ListCommand struct {
	All bool
}

// Parse a string into a gobot command
func Parse(command string) (interface{}, error) {
	args := strings.Split(command, " ")
	if len(args) < 2 {
		return nil, fmt.Errorf("could not understand %s", command)
	}
	switch args[1] {
	case PlayToken:
		return parsePlayCommand(args[2:])
	case MoveToken:
		return parseMoveCommand(args[2:])
	case PassToken:
		return parsePassCommand(args[2:])
	case ShowToken:
		return parseShowCommand(args[2:])
	case ScoreToken:
		return parseScoreCommand(args[2:])
	case ListToken:
		return parseListCommand(args[2:])
	default:
		return nil, fmt.Errorf("could not understand %s", command)
	}
}

func parsePlayCommand(players []string) (*PlayCommand, error) {
	switch len(players) {
	// Two players only
	case 2:
		return &PlayCommand{
			Players: Players{
				Black: []string{players[0]},
				White: []string{players[1]},
			},
			Settings: Settings{
				Vote: false,
			},
		}, nil
	// Allow anyone to vote for moves
	case 0:
		return &PlayCommand{
			Players: Players{
				Anyone: true,
			},
			Settings: Settings{
				Vote:  true,
				Timer: 3600, // 1 hour
			},
		}, nil
	default:
		return nil, fmt.Errorf("incorrect number of players")
	}

}

func parseCoordinates(coords string) ([2]int, error) {
	// coords will be something like A12
	result := [2]int{0, 0}
	if len(coords) < 2 {
		return result, fmt.Errorf("invalid coordinate")
	}
	letter := coords[0]
	number, err := strconv.Atoi(coords[1:])
	if letter < 'A' || letter > 'P' {
		return result, fmt.Errorf("%b is out of range", letter)
	}
	if err != nil {
		return result, fmt.Errorf("%s is not a number", coords[1:])
	}
	result[0] = number - 1        // x
	result[1] = int(letter - 'A') // y
	return result, nil
}

func parseMoveCommand(args []string) (*MoveCommand, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("need to play a move")
	}
	gameID, err := strconv.Atoi(args[0])
	if err == nil {
		if len(args) < 2 {
			return nil, fmt.Errorf("need to play a move")
		}
		coords, err := parseCoordinates(args[1])
		if err != nil {
			return nil, err
		}
		return &MoveCommand{
			GameID:      gameID,
			Coordinates: coords,
		}, nil
	}
	coords, err := parseCoordinates(args[0])
	if err != nil {
		return nil, err
	}
	return &MoveCommand{
		GameID:      -1,
		Coordinates: coords,
	}, nil
}

func parsePassCommand(args []string) (*PassCommand, error) {
	if len(args) < 1 {
		return &PassCommand{GameID: -1}, nil
	}
	gameID, err := strconv.Atoi(args[0])
	if err == nil {
		return &PassCommand{
			GameID: gameID,
		}, nil
	}
	return nil, fmt.Errorf("invalid game id %s", args[0])
}

func parseShowCommand(args []string) (*ShowCommand, error) {
	if len(args) < 1 {
		return &ShowCommand{GameID: -1}, nil
	}
	gameID, err := strconv.Atoi(args[0])
	if err == nil {
		return &ShowCommand{
			GameID: gameID,
		}, nil
	}
	return nil, fmt.Errorf("invalid game id %s", args[0])
}

func parseScoreCommand(args []string) (*ScoreCommand, error) {
	if len(args) < 1 {
		return &ScoreCommand{GameID: -1}, nil
	}
	gameID, err := strconv.Atoi(args[0])
	if err == nil {
		return &ScoreCommand{
			GameID: gameID,
		}, nil
	}
	return nil, fmt.Errorf("invalid game id %s", args[0])
}

func parseListCommand(args []string) (*ListCommand, error) {
	if len(args) < 1 {
		return &ListCommand{}, nil
	}
	if args[0] == "all" {
		return &ListCommand{
			All: true,
		}, nil
	}
	return nil, fmt.Errorf("unknown value %s", args[0])
}
