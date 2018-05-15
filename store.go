package gobot

import (
	"database/sql"
	"fmt"
	"math/rand"
	"sort"
	"time"
)

// GameStore stores a map of game IDs to corresponding games
type GameStore struct {
	Games map[int64]*Game
	DB    *sql.DB
}

// NewGameStore creates a game store connected to a database
func NewGameStore(db *sql.DB) GameStore {
	return GameStore{
		Games: map[int64]*Game{},
		DB:    db,
	}
}

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

// ToList gets a list of the games in the store
func (s GameStore) ToList() GameList {
	list := []*Game{}
	for _, game := range s.Games {
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

// Vote is a vote for a move
type Vote struct {
	Move Coords
	Pass bool
}

// VoteStore stores the current set of votes per game
type VoteStore map[int64]map[string]*Vote

// NewVoteStore creates an empty vote store
func NewVoteStore() VoteStore {
	return map[int64]map[string]*Vote{}
}

// New creates a new vote store for a game
func (s VoteStore) New(game *Game) {
	s[game.ID] = map[string]*Vote{}
}

// Move votes to make a move
func (s VoteStore) Move(game *Game, user string, coords Coords) {
	if _, ok := s[game.ID]; !ok {
		s[game.ID] = map[string]*Vote{}
	}
	s[game.ID][user] = &Vote{Move: coords}
}

// Pass votes to pass
func (s VoteStore) Pass(game *Game, user string) {
	if _, ok := s[game.ID]; !ok {
		s[game.ID] = map[string]*Vote{}
	}
	s[game.ID][user] = &Vote{Pass: true}
}

// Random picks a random vote
func (s VoteStore) Random(game *Game) (*Vote, string) {
	votes := s[game.ID]
	users := []string{}
	for user := range votes {
		users = append(users, user)
	}
	if len(users) == 0 {
		return nil, ""
	}
	// sort usernames so we can deterministically know which user it will pick
	sort.Strings(users)
	roll := rand.Intn(len(users))
	player := users[roll]
	return votes[player], player
}

// Reset the store for a game
func (s VoteStore) Reset(game *Game) {
	s[game.ID] = map[string]*Vote{}
}

// Scheduler schedules regular events for games
type Scheduler struct {
	VoteTimer *time.Timer
	Response  chan *Response
}

// StartVoting schedules a vote timer and sends responses to a channel
func (s *Scheduler) StartVoting(game *Game, votes *VoteStore) {
	if !s.VoteTimer.Stop() {
		<-s.VoteTimer.C
	}
	s.VoteTimer.Reset(game.Settings.Timer)
	for {
		select {
		case <-s.VoteTimer.C:
			// reset the timer to vote again in x seconds
			s.VoteTimer.Reset(game.Settings.Timer)
			vote, player := votes.Random(game)
			if vote == nil {
				break
			}
			err := game.Vote(player, vote)
			if err != nil {
				break
			}
			votes.Reset(game)
			if vote.Pass {
				s.Response <- NewGameResponse(game, "voted to pass")
			} else {
				details := fmt.Sprintf(
					"voted to move at %s", vote.Move.String())
				s.Response <- NewGameResponse(game, details)
			}
		}
	}
}

// A SchedulerStore stores schedulers for games
type SchedulerStore struct {
	Schedulers map[int64]*Scheduler
	Responses  chan *Response
}

// NewSchedulerStore creates a new scheduler store
func NewSchedulerStore(responses chan *Response) SchedulerStore {
	return SchedulerStore{
		Schedulers: map[int64]*Scheduler{},
		Responses:  responses,
	}
}

// New creates a new scheduler for a game
func (s SchedulerStore) New(game *Game) *Scheduler {
	s.Schedulers[game.ID] = &Scheduler{
		VoteTimer: time.NewTimer(game.Settings.Timer),
		Response:  s.Responses,
	}
	return s.Schedulers[game.ID]
}

// ResetVoting resets a voting timer
func (s SchedulerStore) ResetVoting(game *Game) {
	t := s.Schedulers[game.ID].VoteTimer
	if !t.Stop() {
		<-t.C
	}
	t.Reset(game.Settings.Timer)
}
