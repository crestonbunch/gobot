package gobot

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
