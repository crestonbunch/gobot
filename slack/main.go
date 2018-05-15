package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/crestonbunch/gobot"
	"github.com/nlopes/slack"
)

func main() {
	token := os.Getenv("SLACK_API_TOKEN")

	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)

	api := slack.New(token)
	api.SetDebug(false)

	i, err := gobot.NewSlackInterface(api)
	if err != nil {
		panic(err)
	}
	defer i.Close()

	db, err := sql.Open("sqlite3", "./games.db")
	if err != nil {
		panic(err)
	}
	bot, err := gobot.New(db)
	if err != nil {
		panic(err)
	}
	defer bot.Close()

	err = bot.Start()
	if err != nil {
		panic(err)
	}
	go i.StartReceiving(bot)
	go i.StartSending(bot)

	i.Block()
}
