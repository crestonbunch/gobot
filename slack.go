package gobot

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/nlopes/slack"
)

// SlackInterface controls the bot through a slack channel
type SlackInterface struct {
	BotID   string
	Channel string
	API     *slack.Client
	RTM     *slack.RTM
	Command chan string
	Stop    chan bool
}

// NewSlackInterface connects to slack and sets up an interface.
func NewSlackInterface(api *slack.Client) (*SlackInterface, error) {
	identity, err := api.AuthTest()
	if err != nil {
		return nil, err
	}
	rtm := api.NewRTM()

	return &SlackInterface{
		BotID:   identity.UserID,
		API:     api,
		RTM:     rtm,
		Command: make(chan string),
		Stop:    make(chan bool),
	}, nil
}

// Close go channels and open connections
func (i *SlackInterface) Close() {
}

// Block the current goroutine until the stop signal is received
func (i *SlackInterface) Block() bool {
	return <-i.Stop
}

func (i *SlackInterface) sendText(text string) {
	params := slack.PostMessageParameters{}
	i.API.PostMessage(i.Channel, text, params)
}

func (i *SlackInterface) sendImage(im image.Image, name, details string) {
	i.sendText("please hold...")
	temp, err := ioutil.TempFile("", "gobot")
	if err != nil {
		i.sendText("could not save image")
		return
	}
	defer os.Remove(temp.Name())
	err = png.Encode(temp, im)
	if err != nil {
		i.sendText("could encode png image")
		return
	}
	file := &slack.FileUploadParameters{
		Title:          name,
		File:           temp.Name(),
		Channels:       []string{i.Channel},
		InitialComment: details,
	}
	_, err = i.API.UploadFile(*file)
	if err != nil {
		i.sendText(fmt.Sprintf("error uploading image %s", err.Error()))
	}
}

func (i *SlackInterface) sendGame(id int64, g Game, details string) {
	im, _ := Render(g.Board())
	suffix := ""
	if g.Finished() {
		suffix = " (finished)"
	}
	name := fmt.Sprintf("Game %d%s", id, suffix)
	i.sendImage(im, name, details)
}

// IsSlackCommand checks if the command is for the slack bot
func (i *SlackInterface) IsSlackCommand(input string) bool {
	return strings.HasPrefix(input, "<@"+i.BotID+"> ")
}

// ConvertSlackCommand sanitizes input from slack into a standard format
// to be consumed by the gobot.
func (i *SlackInterface) ConvertSlackCommand(input string) string {
	re := regexp.MustCompile("<@([^>]+)>")
	// strip leading @gobot command
	output := strings.Replace(input, "<@"+i.BotID+"> ", "", 1)
	// sanitize @user substrings
	output = re.ReplaceAllString(output, "$1")
	return output
}

// StartSending replies received along the reply channel
func (i *SlackInterface) StartSending(server *Server) {
	for r := range server.Replies {
		if r == nil {
			continue
		}
		if r.Session != nil {
			i.sendGame(r.Session.Storable.ID(), r.Session.Game, r.Details)
		} else {
			i.sendText(r.Text)
		}
	}
}

// StartReceiving commands from the Slack client
func (i *SlackInterface) StartReceiving(server *Server) {
	go i.RTM.ManageConnection()

	for msg := range i.RTM.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MemberJoinedChannelEvent:
			i.Channel = ev.Channel
		case *slack.MessageEvent:
			if ev.SubType == "" && i.IsSlackCommand(ev.Text) {
				i.Channel = ev.Channel
				command := i.ConvertSlackCommand(ev.Text)
				err := server.Handle(command, ev.User)
				if err != nil {
					i.sendText(err.Error())
				}
			}
		}
	}
}
