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

	bot := gobot.New()
	api := slack.New(token)
	api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MemberJoinedChannelEvent:
		case *slack.MessageEvent:
			bot.Receive(ev.Text, ev.Channel)
		}
	}

	for msg := range bot.OutgoingMessages {
		api.PostMessage(msg.Channel, msg.Text, slack.PostMessageParameters{})
	}
}
