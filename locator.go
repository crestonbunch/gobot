package gobot

import "errors"

// GameLocator describes rules for picking which game a command should
// be sent to. Either pick a specific game, or pick the game automatically,
type GameLocator struct {
	GameID int64
	Auto   bool
}

// Find the first valid game in a list of games
func (l *GameLocator) Find(games []*Game, user string) (*Game, error) {
	if l.Auto {
		return l.findAuto(games, user)
	}
	return l.findFixed(games, user)
}

func (l *GameLocator) findFixed(games []*Game, user string) (*Game, error) {
	for _, g := range games {
		if g.ID == l.GameID {
			if g.IsPlayer(user) {
				return g, nil
			}
			return nil, errors.New("not your game")
		}
	}
	return nil, errors.New("game not found")
}

func (l *GameLocator) findAuto(games []*Game, user string) (*Game, error) {
	for _, g := range games {
		if g.IsPlayer(user) {
			return g, nil
		}
	}
	return nil, errors.New("game not found")
}
