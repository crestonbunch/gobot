package gobot

import (
	"errors"
)

const (
	// EmptyStone means no one has a stone in this space
	EmptyStone Stone = iota
	// BlackStone is a stone of black's
	BlackStone
	// WhiteStone is a stone of white's
	WhiteStone
	// BoundaryStone represents a stone that is outside the board dimensions
	BoundaryStone
)

// Stone is the color of stone on the board
type Stone int8

// Board is the state of the current game
type Board [][]Stone

// New19by19Board creates an empty 19x19 board
func New19by19Board() Board {
	stones := [][]Stone{}
	for i := 0; i < 19; i++ {
		stones = append(stones, []Stone{})
		for j := 0; j < 19; j++ {
			stones[i] = append(stones[i], EmptyStone)
		}
	}
	return Board(stones)
}

// Get the stone at a particular position. May return BoundaryStone for
// coordinates that are outside the game dimensions.
func (b Board) Get(x, y int) Stone {
	if x < 0 || y < 0 {
		return BoundaryStone
	}
	if y > len(b)-1 {
		return BoundaryStone
	}
	row := b[y]
	if x > len(row)-1 {
		return BoundaryStone
	}
	return row[x]
}

// Height is the height of the board
func (b Board) Height() int {
	return len(b)
}

// Width is the width of the board
func (b Board) Width() int {
	for _, row := range b {
		return len(row)
	}
	return 0
}

// Copy makes a copy of the board
func (b Board) Copy() Board {
	rows := make([][]Stone, len(b))
	for i := range rows {
		rows[i] = make([]Stone, len(b[i]))
		for j, stone := range b[i] {
			rows[i][j] = stone
		}
	}
	return rows
}

// Set returns a copy of the board with the new stone set at (x, y)
func (b Board) Set(x, y int, stone Stone) Board {
	copy := b.Copy()
	if x < 0 || y < 0 {
		return copy
	}
	if y > len(b)-1 {
		return copy
	}
	row := b[y]
	if x > len(row)-1 {
		return copy
	}
	copy[y][x] = stone
	return copy
}

// Liberties counts the number of liberties of a stone and its connected stones
// at (x, y)
func (b Board) Liberties(x, y int) int {
	var recurse func(int, int, Stone)
	liberties := 0
	visited := map[int]map[int]bool{}
	stone := b.Get(x, y)
	// recursively follow neighbors to count liberties
	recurse = func(rx, ry int, color Stone) {
		// Skip stones we've already checked
		if row, ok := visited[rx]; ok {
			if _, ok := row[ry]; ok {
				return
			}
		} else {
			visited[rx] = map[int]bool{}
		}
		// Visit this stone
		visited[rx][ry] = true
		// Count liberties we find
		if b.Get(rx, ry) == EmptyStone {
			liberties++
			return
		}
		// Stop when we find an enemy or boundary stone
		if b.Get(rx, ry) != color {
			return
		}
		// recursively search in all directions for liberties
		recurse(rx-1, ry, color)
		recurse(rx+1, ry, color)
		recurse(rx, ry-1, color)
		recurse(rx, ry+1, color)
	}
	recurse(x, y, stone)
	return liberties
}

// Capture the stone and connected stones at (x, y) and return a copy of the
// board state after the capture occurs and the number of stones captured.
func (b Board) Capture(x, y int) (Board, int) {
	var recurse func(int, int, Stone)
	copy := b.Copy()
	captures := 0
	stone := copy.Get(x, y)
	// recursively follow neighbors to find captures
	recurse = func(rx, ry int, color Stone) {
		if copy.Get(rx, ry) == color {
			copy[ry][rx] = EmptyStone
			captures++
		} else {
			return
		}
		// recursively search in all directions for captures
		recurse(rx-1, ry, color)
		recurse(rx+1, ry, color)
		recurse(rx, ry-1, color)
		recurse(rx, ry+1, color)
	}
	recurse(x, y, stone)
	return copy, captures
}

// Play a stone at (x, y) and return the new board state and the number of
// captures. Returns an error if there are not enough liberties to play the
// given move.
func (b Board) Play(x, y int, stone Stone) (Board, int, error) {
	// Build a map of (x, y) -> nearest 4 neighbors
	neighbors := map[int]map[int]Stone{
		x: {
			y - 1: b.Get(x, y-1),
			y + 1: b.Get(x, y+1),
		},
		x + 1: {
			y: b.Get(x+1, y),
		},
		x - 1: {
			y: b.Get(x-1, y),
		},
	}
	if b.Get(x, y) != EmptyStone {
		return nil, 0, errors.New("must play in an empty space")
	}
	copy := b.Set(x, y, stone)
	captures := 0
	for nx, whys := range neighbors {
		for ny, neighbor := range whys {
			if neighbor == EmptyStone || neighbor == BoundaryStone {
				continue
			}
			// Check if we can capture an opponent's group
			if neighbor != stone {
				liberties := copy.Liberties(nx, ny)
				// Capture a neighbor group if it has no liberties
				if liberties == 0 {
					var caps int
					copy, caps = copy.Capture(nx, ny)
					captures += caps
				}
			}
		}
	}
	// Check if we have enough liberties to play this stone
	if copy.Liberties(x, y) == 0 {
		return nil, 0, errors.New("not enough liberties")
	}
	return copy, captures, nil
}

// Equals checks if two board states are equivalent
func (b Board) Equals(o Board) bool {
	if len(b) != len(o) {
		return false
	}
	for i, row := range b {
		if len(row) != len(o[i]) {
			return false
		}
		for j, stone := range row {
			if o[i][j] != stone {
				return false
			}
		}
	}
	return true
}
