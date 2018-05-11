package gobot

// An Action is a command to be issued by a user in a channel that targets
// specific game.
type Action struct {
	Command GameCommand
	Channel string
	User    string
}

// Run performs the action on a game
func (a *Action) Run(game *Game) (string, error) {
	switch command := a.Command.(type) {
	case *MoveCommand:
		return game.Play(a.User, command.Coordinates)
	case *PassCommand:
		return game.Pass(a.User)
	case *ScoreCommand:
		return "not implemented", nil
	case *ShowCommand:
		return "not implemented", nil
	default:
		return "action did nothing", nil
	}
}
