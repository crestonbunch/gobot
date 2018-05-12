package main

import (
	"log"
	"os"

	"github.com/crestonbunch/gobot"
	"github.com/nlopes/slack"
)

func main() {
	token := os.Getenv("SLACK_API_TOKEN")

	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)

	api := slack.New(token)
	api.SetDebug(true)

	i, err := gobot.NewSlackInterface(api)
	if err != nil {
		panic(err)
	}
	defer i.Close()

	bot := gobot.New()
	go bot.Start()
	go i.StartReceiving(bot)
	go i.StartSending()

	i.Block()
}
