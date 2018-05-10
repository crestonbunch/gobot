package gobot

// Server is the main entrypoint to Gobot. It handles listening to messages
// and making moves in a game.
type Server struct {
	Manager          *Manager
	OutgoingMessages chan *Message
}

// New creates a new gobot server to listen to messages
func New() *Server {
	return &Server{
		Manager:          &Manager{},
		OutgoingMessages: make(chan *Message),
	}
}

// Receive receives a raw text message to parse
func (s *Server) Receive(msg string, channel string) {
	command, err := Parse(msg)
	if err != nil {
		s.Send(&Message{"could not understand command", channel})
		return
	}
	response := s.Manager.Receive(command)
	s.Send(&Message{response, channel})
}

// Send receives a raw text message to parse
func (s *Server) Send(message *Message) {
	s.OutgoingMessages <- message
}
