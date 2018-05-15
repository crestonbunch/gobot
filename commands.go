package gobot

import (
	"errors"
	"fmt"
)

// Command issues a command and returns a response
type Command interface {
	Execute(string, GameStore, VoteStore, SchedulerStore) (*Response, error)
}

// StartCommand is a command to start a new game.
type StartCommand struct {
	Players  Players
	Settings Settings
}

// Execute a start command to begin a new game
func (c *StartCommand) Execute(
	player string, games GameStore, votes VoteStore, schedulers SchedulerStore,
) (*Response, error) {
	game, err := games.New(c.Players, c.Settings)
	if err != nil {
		return nil, err
	}
	votes.New(game)
	scheduler := schedulers.New(game)
	go scheduler.StartVoting(game, &votes)
	return NewGameResponse(game, ""), nil
}

// MoveCommand is a command to make a move.
type MoveCommand struct {
	Coordinates Coords
	Pass        bool
	Locator     GameLocator
}

// Execute a move command to make a move
func (c *MoveCommand) Execute(
	player string, games GameStore, votes VoteStore, schedulers SchedulerStore,
) (*Response, error) {
	game, err := c.Locator.Find(games.Sorted(), player)
	if err != nil {
		return nil, err
	}
	if game.Settings.Vote {
		return nil, errors.New("game requires voting")
	}
	if c.Pass {
		err := game.Pass(player)
		if err != nil {
			return nil, err
		}
		err = games.Save(game)
		if err != nil {
			return nil, err
		}
		return NewGameResponse(game, "passed"), nil
	}
	err = game.Move(player, c.Coordinates)
	if err != nil {
		return nil, err
	}
	err = games.Save(game)
	if err != nil {
		return nil, err
	}
	return NewGameResponse(game, ""), nil
}

// VoteCommand is a command to vote for a move
type VoteCommand MoveCommand

// Execute a move command to vote
func (c *VoteCommand) Execute(
	player string, games GameStore, votes VoteStore, schedulers SchedulerStore,
) (*Response, error) {
	game, err := c.Locator.Find(games.Sorted(), player)
	if err != nil {
		return nil, err
	}
	if !game.Settings.Vote {
		return nil, errors.New("voting not allowed")
	}
	if !game.Authorized(player) {
		return nil, errors.New("not your turn")
	}
	if c.Pass {
		votes.Pass(game, player)
		return NewTextResponse("thanks for voting"), nil
	}
	if !game.Valid(player, c.Coordinates) {
		return nil, errors.New("invalid move")
	}
	votes.Move(game, player, c.Coordinates)
	return NewTextResponse("thanks for voting"), nil
}

// PlayCommand is a command to pick a vote and play it
type PlayCommand struct {
	Locator GameLocator
}

// Execute a play command to make a move
func (c *PlayCommand) Execute(
	player string, games GameStore, votes VoteStore, schedulers SchedulerStore,
) (*Response, error) {
	game, err := c.Locator.Find(games.Sorted(), player)
	if err != nil {
		return nil, err
	}
	vote, player := votes.Random(game)
	if vote == nil {
		return nil, errors.New("no votes cast")
	}
	err = game.Vote(player, vote)
	if err != nil {
		return nil, err
	}
	err = games.Save(game)
	if err != nil {
		return nil, err
	}
	votes.Reset(game)
	schedulers.ResetVoting(game)
	if vote.Pass {
		return NewGameResponse(game, "voted to pass"), nil
	}
	details := fmt.Sprintf("voted to move at %s", vote.Move.String())
	return NewGameResponse(game, details), nil
}
