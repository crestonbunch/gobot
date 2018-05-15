package gobot

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// StartRegex matches a start command
var StartRegex = regexp.MustCompile("^start$")

// TwoPlayerStartRegex matches a start command with two players
var TwoPlayerStartRegex = regexp.MustCompile("^start ([^ ]+) ([^ ]+)$")

// MoveRegex matches a move command
var MoveRegex = regexp.MustCompile("^move (pass|[A-Z][0-9]+)$")

// GameMoveRegex matches a move command for a specific game
var GameMoveRegex = regexp.MustCompile("^move ([0-9]+) (pass|[A-Z][0-9]+)$")

// VoteRegex matches a vote command
var VoteRegex = regexp.MustCompile("^vote (pass|[A-Z][0-9]+)$")

// GameVoteRegex matches a vote command for a specific game
var GameVoteRegex = regexp.MustCompile("^vote ([0-9]+) (pass|[A-Z][0-9]+)$")

// PlayRegex matches a play command
var PlayRegex = regexp.MustCompile("^play$")

// GamePlayRegex matches a play command for a specific game
var GamePlayRegex = regexp.MustCompile("^play ([0-9]+)$")

// ParseCommand parses a command from an input string
func ParseCommand(input string) (Command, error) {
	if StartRegex.MatchString(input) {
		matches := StartRegex.FindStringSubmatch(input)
		return parseStartCommand(matches[1:])
	}
	if TwoPlayerStartRegex.MatchString(input) {
		matches := TwoPlayerStartRegex.FindStringSubmatch(input)
		return parseStartCommand(matches[1:])
	}
	if MoveRegex.MatchString(input) {
		matches := MoveRegex.FindStringSubmatch(input)
		return parseMoveCommand(matches[1:])
	}
	if GameMoveRegex.MatchString(input) {
		matches := GameMoveRegex.FindStringSubmatch(input)
		return parseGameMoveCommand(matches[1:])
	}
	if VoteRegex.MatchString(input) {
		matches := VoteRegex.FindStringSubmatch(input)
		return parseVoteCommand(matches[1:])
	}
	if GameVoteRegex.MatchString(input) {
		matches := GameVoteRegex.FindStringSubmatch(input)
		return parseGameVoteCommand(matches[1:])
	}
	if PlayRegex.MatchString(input) {
		matches := PlayRegex.FindStringSubmatch(input)
		return parsePlayCommand(matches[1:])
	}
	if GamePlayRegex.MatchString(input) {
		matches := GamePlayRegex.FindStringSubmatch(input)
		return parsePlayCommand(matches[1:])
	}
	return nil, fmt.Errorf("%s not understood", input)
}

func parseStartCommand(players []string) (*StartCommand, error) {
	switch len(players) {
	// Two players only
	case 2:
		return &StartCommand{
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
		return &StartCommand{
			Players: Players{
				Anyone: true,
			},
			Settings: Settings{
				Vote:  true,
				Timer: 3600 * time.Second, // 1 hour
			},
		}, nil
	default:
		return nil, fmt.Errorf("incorrect number of players")
	}

}

func parseCoordinates(coords string) (Coords, error) {
	// coords will be something like A12
	result := Coords{0, 0}
	if len(coords) < 2 {
		return result, fmt.Errorf("invalid coordinate")
	}
	letter := coords[0]
	number, err := strconv.Atoi(coords[1:])
	if letter < 'A' || letter > 'S' {
		return result, fmt.Errorf("%s is out of range", string(letter))
	}
	if err != nil {
		return result, fmt.Errorf("%s is not a number", coords[1:])
	}
	result[0] = number - 1        // x
	result[1] = int(letter - 'A') // y
	return Coords(result), nil
}

func parseMoveCommand(args []string) (*MoveCommand, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("need to play a move")
	}
	if args[0] == "pass" {
		return &MoveCommand{
			Pass:    true,
			Locator: GameLocator{Auto: true},
		}, nil
	}
	coords, err := parseCoordinates(args[0])
	if err != nil {
		return nil, err
	}
	return &MoveCommand{
		Coordinates: coords,
		Locator:     GameLocator{Auto: true},
	}, nil
}

func parseGameMoveCommand(args []string) (*MoveCommand, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("need to play a game and move id")
	}
	gameID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return nil, err
	}
	if args[1] == "pass" {
		return &MoveCommand{
			Pass:    true,
			Locator: GameLocator{GameID: gameID},
		}, nil
	}
	coords, err := parseCoordinates(args[1])
	if err != nil {
		return nil, err
	}
	return &MoveCommand{
		Coordinates: coords,
		Locator:     GameLocator{GameID: gameID},
	}, nil
}

func parseVoteCommand(args []string) (*VoteCommand, error) {
	move, err := parseMoveCommand(args)
	if err != nil {
		return nil, err
	}
	vote := VoteCommand(*move)
	return &vote, nil
}

func parseGameVoteCommand(args []string) (*VoteCommand, error) {
	move, err := parseGameMoveCommand(args)
	if err != nil {
		return nil, err
	}
	vote := VoteCommand(*move)
	return &vote, nil
}

func parsePlayCommand(args []string) (*PlayCommand, error) {
	return &PlayCommand{
		Locator: GameLocator{Auto: true},
	}, nil
}

func parseGamePlayCommand(args []string) (*PlayCommand, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("missing game id")
	}
	gameID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return nil, err
	}
	return &PlayCommand{
		Locator: GameLocator{GameID: gameID},
	}, nil
}
