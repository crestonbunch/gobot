//go:generate mockgen -destination=mocks/store.go -package=mocks -source=interface.go

package gobot

import (
	"fmt"
	"strconv"
	"time"
)

// A Move is a move that can be made in a game
type Move struct {
	Pass   bool
	Coords Coords
}

// String implements the stringer interface
func (m Move) String() string {
	if m.Pass {
		return "pass"
	}
	return fmt.Sprintf("move at %s", m.Coords.String())
}

// Coords represents a board coordinate in (x, y) values
type Coords [2]int

// String implements the stringer interface
func (c Coords) String() string {
	letter := string(c[1] + 'A')
	number := strconv.Itoa(c[0] + 1)
	return letter + number
}

// Votable implements something that can save and recall votes
type Votable interface {
	// Vote for a move
	Vote(*Move) error
	// Schedule starts a vote timer, and resets any existing timer
	Schedule() *time.Timer
	// Block until the vote timer is up
	Block()
	// Random picks random vote
	Random() (*Move, error)
	// Empty returns true if no votes have been cast
	Empty() bool
	// Reset the votes made
	Reset() error
	// Whether or not voting is required
	Required() bool
}

// Storable implements something that can be serialized and loaded by ID
type Storable interface {
	// Return a unique identifier for this storable
	ID() int64
	// Load the contents of a byte array into this storable
	Load([]byte) error
	// Serialize this game to a byte array
	Save() ([]byte, error)
}

// Playable implements something that players play
type Playable interface {
	// Check if a player is participating in a game
	IsPlaying(playerID string) bool
	// Check if a player can make the next move
	CanMove(playerID string) bool
}

// A Game interface for a game
type Game interface {
	// Get the current game board
	Board() Board
	// Check if the game is finished
	Finished() bool
	// Play a move
	Move(*Move) error
	// Whether or not a move is valid to play next
	Validate(*Move) bool
}

// Store is an interface for something that can be used to store games.
type Store interface {
	// Close the store connection
	Close() error
	// Load the store from a storage backend
	Load() error
	// Get a particular session by id
	Get(id int64) (*Session, error)
	// Create a new session from a blueprint
	New(Blueprint) (*Session, error)
	// Return the last session played
	Last() (*Session, error)
	// Save a storable to storage
	Save(Storable) error
	// List active sessions, optionally listing all sessions
	List(all bool) ([]*Session, error)
}
