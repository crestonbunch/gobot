package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/crestonbunch/gobot"
	_ "github.com/mattn/go-sqlite3"

	"github.com/nlopes/slack"
)

func main() {
	token := os.Getenv("SLACK_API_TOKEN")

	logger := log.New(os.Stdout, "slack: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)

	api := slack.New(token)
	api.SetDebug(false)

	i, err := gobot.NewSlackInterface(api)
	if err != nil {
		logger.Fatal(err)
	}
	defer i.Close()

	db, err := sql.Open("sqlite3", "./games.db")
	if err != nil {
		logger.Fatal(err)
	}
	bot, err := gobot.NewServer(db)
	if err != nil {
		logger.Fatal(err)
	}
	defer bot.Close()

	err = bot.Start()
	if err != nil {
		logger.Fatalf("error starting server: %s", err.Error())
	}
	go i.StartReceiving(bot)
	go i.StartSending(bot)

	i.Block()
}
