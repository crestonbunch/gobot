package gobot

// Move represents a player placing a stone on the board
type Move struct {
	X     int
	Y     int
	Stone Stone
}
