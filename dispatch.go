package gobot

import (
	"fmt"
)

// Response is a response to a command. It can contain text or a game state.
type Response struct {
	Text string
	Game *Game
}

// NewTextResponse builds a text response
func NewTextResponse(text string) *Response {
	return &Response{Text: text}
}

// NewGameResponse builds a game response
func NewGameResponse(game *Game) *Response {
	return &Response{Game: game}
}

// Dispatch a command by a player to a game
func Dispatch(game *Game, player string, command Command) (*Response, error) {
	switch command := command.(type) {
	case *MoveCommand:
		var err error
		var response *Response
		switch game.Settings.Vote {
		case false:
			err = game.Move(player, command.Coordinates)
			response = NewGameResponse(game)
		case true:
			err = game.VoteForMove(player, command.Coordinates)
			response = NewTextResponse("thanks for voting")
		}
		return response, err
	case *PassCommand:
		err := game.Pass(player)
		var response *Response
		switch game.Settings.Vote {
		case false:
			err = game.Pass(player)
			response = NewGameResponse(game)
		case true:
			err = game.VoteForPass(player)
			response = NewTextResponse("thanks for voting")
		}
		return response, err
	}
	return nil, fmt.Errorf("command not implemented")
}
