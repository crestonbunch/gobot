package gobot

import (
	"fmt"
	"strings"
)

// Command issues a command and returns a response
type Command interface {
	Execute(*Request) (*Response, error)
}

// StartCommand is a command to start a new game.
type StartCommand struct {
	Anyone bool
	White  []string
	Black  []string
}

// Execute a start command to begin a new game
func (c *StartCommand) Execute(r *Request) (*Response, error) {
	return InitializerPipeline.Run(r.Session, r.Player, nil)
}

// MoveCommand is a command to make a move.
type MoveCommand struct {
	Move    *Move
	Locator Locator
}

// Execute a move command to make a move
func (c *MoveCommand) Execute(r *Request) (*Response, error) {
	return MovePipeline.Run(r.Session, r.Player, c.Move)
}

// VoteCommand is a command to vote for a move
type VoteCommand MoveCommand

// Execute a move command to vote
func (c *VoteCommand) Execute(r *Request) (*Response, error) {
	return VotePipeline.Run(r.Session, r.Player, c.Move)
}

// PlayCommand is a command to pick a vote and play it
type PlayCommand struct {
	Locator Locator
}

// Execute a play command to make a move
func (c *PlayCommand) Execute(r *Request) (*Response, error) {
	return PlayPipeline.Run(r.Session, r.Player, nil)
}

// ShowCommand is a command to show the game board
type ShowCommand struct {
	Locator Locator
}

// Execute a play command to make a move
func (c *ShowCommand) Execute(r *Request) (*Response, error) {
	return ShowPipeline.Run(r.Session, r.Player, nil)
}

// ListCommand is a command to list available games
type ListCommand struct {
	All bool
}

// Execute a list command to show games
func (c *ListCommand) Execute(r *Request) (*Response, error) {
	list := []string{}
	for _, sess := range r.List {
		id := sess.Storable.ID()
		fin := sess.Game.Finished()
		list = append(list, fmt.Sprintf("%d: finished: %t", id, fin))
	}
	if len(list) == 0 {
		return NewTextResponse("no games found"), nil
	}
	return NewTextResponse(strings.Join(list, "\n")), nil
}
