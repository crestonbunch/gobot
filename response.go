package gobot

// Response is a response to a command. It can contain text or a game state.
type Response struct {
	Text    string
	Game    *Game
	Details string
}

// NewTextResponse builds a text response
func NewTextResponse(text string) *Response {
	return &Response{Text: text}
}

// NewGameResponse builds a game response
func NewGameResponse(game *Game, details string) *Response {
	return &Response{Game: game, Details: details}
}
