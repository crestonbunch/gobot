package gobot

import (
	"database/sql"
	"log"
	"os"
)

// Server handles receiving requests and performing actions.
type Server struct {
	Store    Store
	Sessions map[int64]*Session
	Replies  chan *Response
	logger   *log.Logger
}

// NewServer creates a new gobot server to listen to messages
func NewServer(db *sql.DB) (*Server, error) {
	responses := make(chan *Response)
	return &Server{
		Store:    NewGameStore(db),
		Sessions: map[int64]*Session{},
		Replies:  responses,
		logger:   log.New(os.Stdout, "bot: ", log.Lshortfile),
	}, nil
}

// Start starts the bot server
func (s *Server) Start() error {
	return s.Load()
}

// Handle a command and return a response
func (s *Server) Handle(input, player string) error {
	s.logger.Printf("handling %s", input)
	cmd, err := ParseCommand(input)
	if err != nil {
		return err
	}
	req, err := NewRequest(cmd, player, s)
	if err != nil {
		return err
	}
	response, err := cmd.Execute(req)
	if err != nil {
		return err
	}
	s.Replies <- response
	if req.Session != nil {
		return s.Save(req.Session.Storable)
	}
	return nil
}

// Close implements the Storable interface
func (s *Server) Close() error {
	return s.Store.Close()
}

// Load implements the Storable interface
func (s *Server) Load() error {
	err := s.Store.Load()
	if err != nil {
		return err
	}
	sessions, err := s.Store.List(false)
	if err != nil {
		return err
	}
	for _, sess := range sessions {
		s.Sessions[sess.Storable.ID()] = sess
		go sess.Background(s, s.Replies, s.logger)
	}
	s.logger.Printf("loaded %d games from store", len(sessions))
	return nil
}

// Get implements the Storable interface
func (s *Server) Get(id int64) (*Session, error) {
	if sess, ok := s.Sessions[id]; ok {
		return sess, nil
	}
	return s.Store.Get(id)
}

// New implements the Storable interface
func (s *Server) New(bp Blueprint) (*Session, error) {
	sess, err := s.Store.New(bp)
	if err != nil {
		return nil, err
	}
	s.Sessions[sess.Storable.ID()] = sess
	go sess.Background(s, s.Replies, s.logger)
	return sess, nil
}

// Last implements the Storable interface
func (s *Server) Last() (*Session, error) {
	sess, err := s.Store.Last()
	if err != nil {
		return nil, err
	}
	return s.Get(sess.Storable.ID())
}

// Save implements the Storable interface
func (s *Server) Save(storable Storable) error {
	s.logger.Printf("saving %d", storable.ID())
	return s.Store.Save(storable)
}

// List implements the Storable interface
func (s *Server) List(all bool) ([]*Session, error) {
	return s.Store.List(all)
}
