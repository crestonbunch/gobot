package gobot

// Message is a response that gobot sends back after a command
type Message struct {
	Text    string
	Channel string
}
