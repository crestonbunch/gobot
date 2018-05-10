package gobot

// Players defines who is allowed to play the game.
type Players struct {
	// A list of user IDs who are allowed to play as black
	Black []string
	// A list of user IDs who are allowed to play as white
	White []string
	// If true, anyone can play as any color
	Anyone bool
}

// Settings dictates how the game is played
type Settings struct {
	// Whether to allow voting on moves or not
	Vote bool
	// The period of time in seconds to vote between moves
	Period float64
}

// History stores the entire history of a game
type History []Board

// Ko checks if the given board state is the same as the previous board state
func (h History) Ko(b Board) bool {
	if len(h) < 1 {
		return false
	}
	return b.Equals(h[len(h)-2])
}

// Game stores the entire game state for a single game
type Game struct {
	Finished bool
	History  History
	Players  *Players
	Captures struct {
		Black int
		White int
	}
	Settings *Settings
}

// NewGame creates a new game with the given players and settings
func NewGame(players *Players, settings *Settings) *Game {
	return &Game{
		Finished: false,
		History:  History([]Board{New19by19Board()}),
		Players:  players,
		Captures: struct {
			Black int
			White int
		}{0, 0},
		Settings: settings,
	}
}
