package main

import (
	"fmt"
	"log"
	"os"

	"github.com/crestonbunch/gobot"
	"github.com/nlopes/slack"
)

func main() {
	token := os.Getenv("SLACK_API_TOKEN")
	fmt.Println(token)

	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)

	api := slack.New(token)
	api.SetDebug(true)
	identity, err := api.AuthTest()
	if err != nil {
		panic(err)
	}

	bot := gobot.New(identity.UserID)
	rtm := api.NewRTM()

	go bot.Start()
	go rtm.ManageConnection()

	go func() {
		for msg := range rtm.IncomingEvents {
			switch ev := msg.Data.(type) {
			case *slack.MemberJoinedChannelEvent:
			case *slack.MessageEvent:
				bot.Receive(&gobot.Request{
					Command: ev.Text,
					Channel: ev.Channel,
					User:    ev.User,
				})
			}
		}
	}()

	for msg := range bot.Replies {
		params := slack.PostMessageParameters{}
		if msg.File != nil {
			_, err := api.UploadFile(*msg.File)
			if err != nil {
				api.PostMessage(
					msg.Channel,
					fmt.Sprintf("error uploading image %s", err.Error()),
					params,
				)
			}
			// Cleanup file we don't need anymore
			if msg.File.File != "" {
				os.Remove(msg.File.File)
			}
		} else {
			api.PostMessage(msg.Channel, msg.Text, params)
		}
	}
}
