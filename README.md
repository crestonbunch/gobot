Gobot
=====

Gobot is a Slack bot for playing go with other people instead of doing work.

Naturally it is written in Go.

Setup
-----

Set an environment variable with your Slack API token

    export SLACK_API_TOKEN=<token goes here>

Or put it in your .bash_rc file (or wherever you put env variables)

Run
-----

    go run ./slack/main.go

Todo
----

MVP

    * Basic commands (below)
    * Drawing the board on the screen
    * Persistence

Future

    * Player statistics / ranking
    * WWAGD (what would AlphaGo do?)
    * SGF export

Talking to the Bot
-----

1. Invite it into your channel

    > 👋 @gobot

2. Start a new game

    With everyone in the channel (vote per move):
    > @gobot play

    With two players (black and white respectively):
    > @gobot play @goseigen @shusaku

3. Make a move

    Respond to the last move played
    > @gobot move D4

    Play a move in a particular game (e.g. game 14)
    > @gobot move 14 D4

    Pass
    > @gobot pass

4. Show the board

    Show the last game played
    > @gobot show

    Show a particular game (e.g. game 14)
    > @gobot show 14

5. Estimate a game score

    Estimate the last game played
    > @gobot score

    Estimate a particular game (e.g. game 14)
    > @gobot score 14

6. List games

    Unfinished games
    > @gobot list

    All games
    > @gobot list all
