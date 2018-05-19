package gobot

import (
	"encoding/json"
	"errors"
	"math/rand"
	"time"
)

// MoveRule dictates how moves should be made
type MoveRule uint8

const (
	// RequireVote requires players to vote
	RequireVote MoveRule = iota
	// RequireMove requires players to move
	RequireMove
)

// Captures is how many _opponent's_ stones White or Black has captured
type Captures struct {
	Black int `json:"black"`
	White int `json:"white"`
}

// Passes stores which players have passed
type Passes struct {
	Black bool `json:"black"`
	White bool `json:"white"`
}

// Players defines who is allowed to play the game.
type Players struct {
	// A list of user IDs who are allowed to play as black
	Black []string `json:"black"`
	// A list of user IDs who are allowed to play as white
	White []string `json:"white"`
	// If true, anyone can play as any color
	Anyone bool `json:"anyone"`
}

// Voting dictates how voting is done
type Voting struct {
	Required bool          `json:"required"`
	Duration time.Duration `json:"duration"`
}

// A State stores the game state for a game, and implements the Game
// interface. It can be serialized into JSON.
type State struct {
	History   History   `json:"history"`
	Next      Stone     `json:"next"`
	Players   Players   `json:"players"`
	Voting    Voting    `json:"voting"`
	Captures  Captures  `json:"captures"`
	Passes    Passes    `json:"passes"`
	Votes     []*Move   `json:"votes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	id        int64
	timer     *time.Timer
}

// ID implements the Storable interface
func (g *State) ID() int64 {
	return g.id
}

// Save implements the Storable interface
func (g *State) Save() ([]byte, error) {
	g.UpdatedAt = time.Now()
	return json.Marshal(g)
}

// Load implements the Storable interface
func (g *State) Load(blob []byte) error {
	return json.Unmarshal(blob, g)
}

// Board implements the Game interface
func (g *State) Board() Board {
	return g.History[len(g.History)-1]
}

// Validate implements the Game interface
func (g *State) Validate(m *Move) bool {
	if m.Pass {
		return true
	}
	current := g.Board()
	next, _, err := current.Play(m.Coords[0], m.Coords[1], g.Next)
	if err != nil {
		return false
	}
	return !g.History.Ko(next)
}

// Move implements the Game interface
func (g *State) Move(m *Move) error {
	if m.Pass {
		return g.pass()
	}
	return g.move(m)
}

// Finished implements the Game interface
func (g *State) Finished() bool {
	return g.Passes.White && g.Passes.Black
}

// IsPlaying implements the Playable interface
func (g *State) IsPlaying(p string) bool {
	return g.Players.Anyone || g.isPlayerWhite(p) || g.isPlayerBlack(p)
}

// CanMove implements the Playable interface
func (g *State) CanMove(p string) bool {
	switch g.Next {
	case BlackStone:
		return g.isPlayerBlack(p)
	case WhiteStone:
		return g.isPlayerWhite(p)
	}
	return false
}

// Vote implements the Votable interface
func (g *State) Vote(m *Move) error {
	g.Votes = append(g.Votes, m)
	return nil
}

// Schedule implements the Votable interface
func (g *State) Schedule() *time.Timer {
	if g.timer != nil {
		g.timer.Reset(g.Voting.Duration)
	} else {
		g.timer = time.NewTimer(g.Voting.Duration)
	}
	return g.timer
}

// Block implements the Votable interface
func (g *State) Block() {
	<-g.timer.C
}

// Random implements the Votable interface
func (g *State) Random() (*Move, error) {
	if len(g.Votes) == 0 {
		return nil, errors.New("no votes cast")
	}
	roll := rand.Intn(len(g.Votes))
	return g.Votes[roll], nil
}

// Empty implements the Votable interface
func (g *State) Empty() bool {
	return len(g.Votes) == 0
}

// Reset implements the Votable interface
func (g *State) Reset() error {
	g.Votes = []*Move{}
	return nil
}

// Required implements the Votable interface
func (g *State) Required() bool {
	return g.Voting.Required
}

func (g *State) isPlayerWhite(p string) bool {
	if g.Players.Anyone {
		return true
	}
	for _, player := range g.Players.White {
		if player == p {
			return true
		}
	}
	return false
}
func (g *State) isPlayerBlack(p string) bool {
	if g.Players.Anyone {
		return true
	}
	for _, player := range g.Players.Black {
		if p == player {
			return true
		}
	}
	return false
}

func (g *State) move(move *Move) error {
	current := g.Board()
	next, captures, err := current.Play(move.Coords[0], move.Coords[1], g.Next)
	if err != nil {
		return err
	}
	g.History = append(g.History, next)
	g.Passes.Black = false
	g.Passes.White = false
	switch g.Next {
	case BlackStone:
		g.Next = WhiteStone
		g.Captures.Black += captures
	case WhiteStone:
		g.Next = BlackStone
		g.Captures.White += captures
	}
	return nil
}

func (g *State) pass() error {
	switch g.Next {
	case BlackStone:
		g.Next = WhiteStone
		g.Passes.Black = true
	case WhiteStone:
		g.Next = BlackStone
		g.Passes.White = true
	}
	return nil
}
