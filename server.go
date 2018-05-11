package gobot

import (
	"strings"

	"github.com/nlopes/slack"
)

// Server is the main entrypoint to Gobot. It handles listening to messages
// and making moves in a game.
type Server struct {
	BotID     string
	Games     []*Game
	Responses chan *Response
}

// Request is a raw text command received from Slack. It will get parsed
// into a Command, which will become part of an Action.
type Request struct {
	Command string
	Channel string
	User    string
}

// Response is a Slack message to send to a channel.
type Response struct {
	Text       string
	Attachment *slack.Attachment
	Channel    string
}

// New creates a new gobot server to listen to messages
func New(botID string) *Server {
	return &Server{
		BotID:     botID,
		Games:     []*Game{},
		Responses: make(chan *Response),
	}
}

// Start starts the bot server
func (s *Server) Start() {
}

// RespondWithText to the client
func (s *Server) RespondWithText(text, channel string) {
	response := &Response{
		Text:    text,
		Channel: channel,
	}
	s.Responses <- response
}

// Receive receives a raw text message to parse
func (s *Server) Receive(req *Request) {
	// Only parse commands that start with @gobot
	if !strings.HasPrefix(req.Command, "<@"+s.BotID+">") {
		return
	}
	command, err := Parse(req.Command)
	if err != nil {
		s.RespondWithText(err.Error(), req.Channel)
		return
	}
	switch command := command.(type) {
	case *PlayCommand:
		s.Games = append(s.Games, NewGame(command.Players, command.Settings))
		s.RespondWithText("game started", req.Channel)
	case *ListCommand:
		s.RespondWithText("not implemented", req.Channel)
	case GameCommand:
		action := &Action{
			Command: command,
			Channel: req.Channel,
			User:    req.User,
		}
		game, err := Dispatch(s.Games, action)
		if err != nil {
			s.RespondWithText(err.Error(), req.Channel)
			return
		}
		response, err := action.Run(game)
		if err != nil {
			s.RespondWithText(err.Error(), req.Channel)
			return
		}
		s.RespondWithText(response, req.Channel)
	}
}
