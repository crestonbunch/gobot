package gobot

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"

	"github.com/nlopes/slack"
)

// Server is the main entrypoint to Gobot. It handles listening to messages
// and making moves in a game.
type Server struct {
	BotID   string
	NextID  int
	Games   GameStore
	Replies chan *Reply
}

// Request is a raw text command received from Slack. It will get parsed
// into a Command, which will become part of an Action.
type Request struct {
	Command string
	Channel string
	User    string
}

// Reply is a Slack message to send to a channel.
type Reply struct {
	Text    string
	File    *slack.FileUploadParameters
	Channel string
}

// New creates a new gobot server to listen to messages
func New(botID string) *Server {
	return &Server{
		BotID:   botID,
		Games:   map[int]*Game{},
		Replies: make(chan *Reply),
	}
}

// Start starts the bot server
func (s *Server) Start() {
}

// ReplyWithText to the client
func (s *Server) ReplyWithText(text, channel string) {
	reply := &Reply{
		Text:    text,
		Channel: channel,
	}
	s.Replies <- reply
}

// ReplyWithImage to the client
func (s *Server) ReplyWithImage(im image.Image, name, channel string) {
	temp, err := ioutil.TempFile("", "gobot")
	if err != nil {
		s.ReplyWithText("could not save image", channel)
		return
	}
	png.Encode(temp, im)
	file := &slack.FileUploadParameters{
		Title:    name,
		File:     temp.Name(),
		Channels: []string{channel},
	}
	reply := &Reply{
		File:    file,
		Channel: channel,
	}
	s.Replies <- reply
}

// ReplyWithGame to the client
func (s *Server) ReplyWithGame(g *Game, channel string) {
	im, _ := Render(g.Board())
	name := fmt.Sprintf("Game %d", g.ID)
	s.ReplyWithImage(im, name, channel)
}

// ReplyWithResponse to the client
func (s *Server) ReplyWithResponse(r *Response, channel string) {
	if r.Game != nil {
		s.ReplyWithGame(r.Game, channel)
	} else {
		s.ReplyWithText(r.Text, channel)
	}
}

// Receive receives a raw text message to parse
func (s *Server) Receive(req *Request) {
	// Only parse commands that start with @gobot
	if !isCommand(req.Command, s.BotID) {
		return
	}
	input := stripName(req.Command, s.BotID)

	if isStartCommand(input) {
		// handle command to start a new game
		command, err := ParseStartCommand(req.Command)
		if err != nil {
			s.ReplyWithText(err.Error(), req.Channel)
			return
		}
		players := command.(*PlayCommand).Players
		settings := command.(*PlayCommand).Settings
		game := NewGame(s.NextID, players, settings)
		s.Games.Add(game)
		s.NextID++
		s.ReplyWithGame(game, req.Channel)
	} else if isGameCommand(input) {
		// handle command to modify a game
		command, locator, err := ParseGameCommand(req.Command)
		if err != nil {
			s.ReplyWithText(err.Error(), req.Channel)
			return
		}
		user := req.User
		game, err := locator.Find(s.Games.Sorted(), user)
		if err != nil {
			s.ReplyWithText(err.Error(), req.Channel)
			return
		}
		response, err := Dispatch(game, user, command)
		if err != nil {
			s.ReplyWithText(err.Error(), req.Channel)
			return
		}
		s.ReplyWithResponse(response, req.Channel)
	} else if isInfoCommand(input) {
		// handle command to get info about games
		_, err := ParseInfoCommand(req.Command)
		if err != nil {
			s.ReplyWithText(err.Error(), req.Channel)
			return
		}
	}
}
