package gobot

import "database/sql"

// Server handles receiving requests and performing actions.
type Server struct {
	NextID     int
	Games      GameStore
	Votes      VoteStore
	Schedulers SchedulerStore
	Replies    chan *Response
}

// New creates a new gobot server to listen to messages
func New(db *sql.DB) (*Server, error) {
	responses := make(chan *Response)
	return &Server{
		Games:      NewGameStore(db),
		Votes:      NewVoteStore(),
		Schedulers: NewSchedulerStore(responses),
		Replies:    responses,
	}, nil
}

// Start starts the bot server
func (s *Server) Start() error {
	err := s.Games.Init()
	if err != nil {
		return err
	}
	err = s.Games.Load()
	if err != nil {
		return err
	}
	return nil
}

// Close closes the bot server connections
func (s *Server) Close() {
	s.Games.DB.Close()
}

// Handle a command and return a response
func (s *Server) Handle(input, user string) error {
	command, err := ParseCommand(input)
	if err != nil {
		return err
	}
	response, err := command.Execute(user, s.Games, s.Votes, s.Schedulers)
	if err != nil {
		return err
	}
	s.Replies <- response
	return nil
}
