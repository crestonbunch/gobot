package gobot

// Response is a response to a command. It can contain text or a game state.
type Response struct {
	Session *Session
	Text    string
	Details string
}

// NewTextResponse builds a text response
func NewTextResponse(text string) *Response {
	return &Response{Text: text}
}

// NewSessionResponse builds a session response
func NewSessionResponse(s *Session, details string) *Response {
	return &Response{Session: s, Details: details}
}
