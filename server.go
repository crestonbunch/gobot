package gobot

// Server handles receiving requests and performing actions.
type Server struct {
	NextID     int
	Games      GameStore
	Votes      VoteStore
	Schedulers SchedulerStore
	Replies    chan *Response
}

// New creates a new gobot server to listen to messages
func New() *Server {
	responses := make(chan *Response)
	return &Server{
		Games:      map[int]*Game{},
		Votes:      map[int]map[string]*Vote{},
		Schedulers: NewSchedulerStore(responses),
		Replies:    responses,
	}
}

// Start starts the bot server
func (s *Server) Start() {
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
