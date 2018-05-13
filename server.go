package gobot

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Server handles receiving requests and performing actions.
type Server struct {
	NextID  int
	Games   GameStore
	Timers  map[*Game]*time.Timer
	Replies chan *Response
}

// New creates a new gobot server to listen to messages
func New() *Server {
	return &Server{
		Games:   map[int]*Game{},
		Timers:  map[*Game]*time.Timer{},
		Replies: make(chan *Response),
	}
}

// Start starts the bot server
func (s *Server) Start() {
}

// AddTimer adds a vote timer for a game
func (s *Server) AddTimer(game *Game) {
	s.Timers[game] = time.AfterFunc(
		game.Settings.Timer*time.Second,
		func() {
			rng := rand.New(rand.NewSource(time.Now().Unix()))
			s.Replies <- game.PickRandomVote(rng)
			// Add the timer again to keep on going
			s.AddTimer(game)
		},
	)
}

// Handle a command and return a response
func (s *Server) Handle(input, user string) (*Response, error) {
	if isStartCommand(input) {
		// handle command to start a new game
		command, err := ParseStartCommand(input)
		if err != nil {
			return nil, err
		}
		players := command.(*PlayCommand).Players
		settings := command.(*PlayCommand).Settings
		game := NewGame(s.NextID, players, settings)
		s.Games.Add(game)
		s.NextID++
		if game.Settings.Vote {
			// add a vote timer
			s.AddTimer(game)
		}
		return NewGameResponse(game), nil
	} else if isGameCommand(input) {
		// handle command to modify a game
		command, locator, err := ParseGameCommand(input)
		if err != nil {
			return nil, err
		}
		game, err := locator.Find(s.Games.Sorted(), user)
		if err != nil {
			return nil, err
		}
		response, err := Dispatch(game, user, command)
		if err != nil {
			return nil, err
		}
		return response, nil
	} else if isInfoCommand(input) {
		// handle command to get info about games
		_, err := ParseInfoCommand(input)
		if err != nil {
			return nil, err
		}
		return nil, errors.New("command not implemented")
	}
	return nil, fmt.Errorf("command \"%s\" not found", input)
}
