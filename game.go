package gobot

import (
	"errors"
	"strconv"
	"time"
)

// Players defines who is allowed to play the game.
type Players struct {
	// A list of user IDs who are allowed to play as black
	Black []string `json:"black"`
	// A list of user IDs who are allowed to play as white
	White []string `json:"white"`
	// If true, anyone can play as any color
	Anyone bool `json:"anyone"`
}

// Settings dictates how the game is played
type Settings struct {
	// Whether to allow voting on moves or not
	Vote bool `json:"vote"`
	// The period of time in seconds to vote between moves
	Timer time.Duration `json:"timer"`
}

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

// Coords represents a board coordinate in (x, y) values
type Coords [2]int

// String implements the stringer interface
func (c Coords) String() string {
	letter := string(c[0] + 'A')
	number := strconv.Itoa(c[1] + 1)
	return letter + number
}

// History stores the entire history of a game
type History []Board

// Ko checks if the given board state is the same as the previous board state
func (h History) Ko(b Board) bool {
	if len(h) < 2 {
		return false
	}
	return b.Equals(h[len(h)-2])
}

// Game stores the entire game state for a single game
type Game struct {
	ID        int64     `json:"-"`
	History   History   `json:"history"`
	Next      Stone     `json:"next"`
	Players   Players   `json:"players"`
	Settings  Settings  `json:"settings"`
	Captures  Captures  `json:"captures"`
	Passes    Passes    `json:"passes"`
	Finished  bool      `json:"finished"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Board gets the current board state
func (g *Game) Board() Board {
	return g.History[len(g.History)-1]
}

// IsPlayer checks if the user is a player of this game
func (g *Game) IsPlayer(user string) bool {
	return g.Players.Anyone || g.IsPlayerWhite(user) || g.IsPlayerBlack(user)
}

// IsPlayerWhite checks if user is authorized to play as white
func (g *Game) IsPlayerWhite(user string) bool {
	if g.Players.Anyone {
		return true
	}
	for _, player := range g.Players.White {
		if player == user {
			return true
		}
	}
	return false
}

// IsPlayerBlack checks if user is authorized to play as black
func (g *Game) IsPlayerBlack(user string) bool {
	if g.Players.Anyone {
		return true
	}
	for _, player := range g.Players.Black {
		if player == user {
			return true
		}
	}
	return false
}

// Authorized checks if the user can make the next move
func (g *Game) Authorized(user string) bool {
	switch g.Next {
	case BlackStone:
		return g.IsPlayerBlack(user)
	case WhiteStone:
		return g.IsPlayerWhite(user)
	}
	return false
}

// Valid checks whether a move can be made or not
func (g *Game) Valid(player string, coords Coords) bool {
	current := g.Board()
	next, _, err := current.Play(coords[0], coords[1], g.Next)
	if err != nil {
		return false
	}
	return !g.History.Ko(next)
}

// Move for a particular coordinate
func (g *Game) Move(player string, coords Coords) error {
	if !g.Authorized(player) {
		return errors.New("not your turn")
	}
	if !g.Valid(player, coords) {
		return errors.New("invalid move")
	}
	current := g.Board()
	// skip err check because we already made it
	next, captures, _ := current.Play(coords[0], coords[1], g.Next)
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

// Pass for a player
func (g *Game) Pass(player string) error {
	if !g.Authorized(player) {
		return errors.New("not your turn")
	}
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

// Vote plays the given vote for a player
func (g *Game) Vote(player string, vote *Vote) error {
	if vote.Pass {
		return g.Pass(player)
	}
	return g.Move(player, vote.Move)
}
