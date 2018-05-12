package gobot

import (
	"fmt"
	"regexp"
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

// Command is an interface that represents a parsed bot command
type Command interface{}

// PlayCommand initiates a new game with the desired settings among the
// desired players.
type PlayCommand struct {
	Players  Players
	Settings Settings
}

// MoveCommand plays a move in a game
type MoveCommand struct {
	Coordinates [2]int
}

// PassCommand plays a pass in a game
type PassCommand struct{}

// ShowCommand shows the current game state
type ShowCommand struct{}

// ScoreCommand scores the current game state
type ScoreCommand struct{}

// ListCommand lists games
type ListCommand struct {
	All bool
}

func isCommand(input string, botID string) bool {
	return strings.HasPrefix(input, "<@"+botID+">")
}

func stripName(input string, botID string) string {
	return strings.Replace(input, "<@"+botID+"> ", "", 1)
}

func isStartCommand(input string) bool {
	return strings.HasPrefix(input, PlayToken)
}

func isGameCommand(input string) bool {
	return strings.HasPrefix(input, MoveToken) ||
		strings.HasPrefix(input, PassToken) ||
		strings.HasPrefix(input, ShowToken) ||
		strings.HasPrefix(input, ScoreToken)
}

func isInfoCommand(input string) bool {
	return strings.HasPrefix(input, ListToken)
}

// ParseStartCommand parses a command that starts a new game
func ParseStartCommand(command string) (Command, error) {
	args := strings.Split(command, " ")
	if len(args) < 2 {
		return nil, fmt.Errorf("could not understand %s", command)
	}
	switch args[1] {
	case PlayToken:
		return parsePlayCommand(args[2:])
	default:
		return nil, fmt.Errorf("could not understand %s", command)
	}
}

// ParseGameCommand parses a command that targets a specific game, and returns
// a game locator that can be used to find the target game.
func ParseGameCommand(command string) (Command, *GameLocator, error) {
	args := strings.Split(command, " ")
	if len(args) < 2 {
		return nil, nil, fmt.Errorf("could not understand %s", command)
	}
	switch args[1] {
	case MoveToken:
		return parseMoveCommand(args[2:])
	case PassToken:
		return parsePassCommand(args[2:])
	case ShowToken:
		return parseShowCommand(args[2:])
	case ScoreToken:
		return parseScoreCommand(args[2:])
	default:
		return nil, nil, fmt.Errorf("could not understand %s", command)
	}
}

// ParseInfoCommand parses a command that does not operate on a game, but might
// list games, rank users, etc.
func ParseInfoCommand(command string) (Command, error) {
	args := strings.Split(command, " ")
	if len(args) < 2 {
		return nil, fmt.Errorf("could not understand %s", command)
	}
	switch args[1] {
	case ListToken:
		return parseListCommand(args[2:])
	default:
		return nil, fmt.Errorf("could not understand %s", command)
	}
}

func parseUserID(token string) (string, error) {
	re := regexp.MustCompile("<@([^>]+)>")
	matches := re.FindStringSubmatch(token)
	if len(matches) < 2 {
		return "", fmt.Errorf("%s is not a valid user", token)
	}
	return matches[1], nil
}

func parsePlayCommand(players []string) (*PlayCommand, error) {
	switch len(players) {
	// Two players only
	case 2:
		p1, err := parseUserID(players[0])
		p2, err := parseUserID(players[1])
		if err != nil {
			return nil, err
		}
		return &PlayCommand{
			Players: Players{
				Black: []string{p1},
				White: []string{p2},
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
	if letter < 'A' || letter > 'S' {
		return result, fmt.Errorf("%s is out of range", string(letter))
	}
	if err != nil {
		return result, fmt.Errorf("%s is not a number", coords[1:])
	}
	result[0] = number - 1        // x
	result[1] = int(letter - 'A') // y
	return result, nil
}

func parseMoveCommand(args []string) (*MoveCommand, *GameLocator, error) {
	if len(args) < 1 {
		return nil, nil, fmt.Errorf("need to play a move")
	}
	gameID, err := strconv.Atoi(args[0])
	if err == nil {
		if len(args) < 2 {
			return nil, nil, fmt.Errorf("need to play a move")
		}
		coords, err := parseCoordinates(args[1])
		if err != nil {
			return nil, nil, err
		}
		return &MoveCommand{
				Coordinates: coords,
			}, &GameLocator{
				GameID: gameID,
			}, nil
	}
	coords, err := parseCoordinates(args[0])
	if err != nil {
		return nil, nil, err
	}
	return &MoveCommand{
			Coordinates: coords,
		}, &GameLocator{
			Auto: true,
		}, nil
}

func parsePassCommand(args []string) (*PassCommand, *GameLocator, error) {
	if len(args) < 1 {
		return &PassCommand{}, &GameLocator{Auto: true}, nil
	}
	gameID, err := strconv.Atoi(args[0])
	if err == nil {
		return &PassCommand{}, &GameLocator{GameID: gameID}, nil
	}
	return nil, nil, fmt.Errorf("invalid game id %s", args[0])
}

func parseShowCommand(args []string) (*ShowCommand, *GameLocator, error) {
	if len(args) < 1 {
		return &ShowCommand{}, &GameLocator{Auto: true}, nil
	}
	gameID, err := strconv.Atoi(args[0])
	if err == nil {
		return &ShowCommand{}, &GameLocator{GameID: gameID}, nil
	}
	return nil, nil, fmt.Errorf("invalid game id %s", args[0])
}

func parseScoreCommand(args []string) (*ScoreCommand, *GameLocator, error) {
	if len(args) < 1 {
		return &ScoreCommand{}, &GameLocator{Auto: true}, nil
	}
	gameID, err := strconv.Atoi(args[0])
	if err == nil {
		return &ScoreCommand{}, &GameLocator{GameID: gameID}, nil
	}
	return nil, nil, fmt.Errorf("invalid game id %s", args[0])
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
