package gobot

import "sort"

// GameStore stores a map of game IDs to corresponding games
type GameStore map[int]*Game

// GameList is a list of games in some order
type GameList []*Game

// Len is part of sort.Interface
func (l GameList) Len() int {
	return len(l)
}

// Swap is part of sort.Interface
func (l GameList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

// Less is part of sort.Interface
func (l GameList) Less(i, j int) bool {
	return l[i].UpdatedAt.Before(l[j].UpdatedAt)
}

// Add adds a game to the store
func (s GameStore) Add(g *Game) {
	s[g.ID] = g
}

// ToList gets a list of the games in the store
func (s GameStore) ToList() GameList {
	list := []*Game{}
	for _, game := range s {
		list = append(list, game)
	}
	return list
}

// Sorted returns a list of games sorted by most recently updated first
func (s *GameStore) Sorted() GameList {
	list := s.ToList()
	sort.Sort(list)
	return list
}
