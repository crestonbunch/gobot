package gobot

import "fmt"

// Manager handles dispatching commands to the various games in play
type Manager struct {
	Games []*Game
}

// Receive performs an action based on a command
func (m *Manager) Receive(command interface{}) string {
	switch command.(type) {
	case *PlayCommand:
		return m.NewGame(command.(*PlayCommand))
	default:
		return "unsupported command"
	}
}

// NewGame creates a new game to manage
func (m *Manager) NewGame(command *PlayCommand) string {
	id := len(m.Games)
	m.Games = append(m.Games, NewGame(command.Players, command.Settings))
	return fmt.Sprintf("Started game %d", id)
}
