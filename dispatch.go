package gobot

import (
	"errors"
)

// Dispatch sends actions to the appropriate game instance
func Dispatch(games []*Game, action *Action) (*Game, error) {
	if action.Command.Game() > len(games)-1 {
		return nil, errors.New("game does not exist")
	} else if action.Command.Game() < 0 {
		// find the most recent game that the player can move in
		for i := len(games) - 1; i >= 0; i-- {
			game := games[i]
			if game.IsPlayer(action.User) {
				return game, nil
			}
		}
		return nil, errors.New("could not find game")
	}
	return games[action.Command.Game()], nil
}
