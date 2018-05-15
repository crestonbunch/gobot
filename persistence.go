package gobot

import (
	"encoding/json"
	"time"
)

// Init persistent storage
func (s GameStore) Init() error {
	stmt := `
	CREATE TABLE IF NOT EXISTS
	games
	(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		blob BLOB NOT NULL,
		finished BOOLEAN NOT NULL DEFAULT FALSE
	)
	`
	_, err := s.DB.Exec(stmt)
	return err
}

// Load games from persistent storage
func (s GameStore) Load() error {
	rows, err := s.DB.Query("SELECT id, blob FROM games WHERE finished = false")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var blob []byte
		err = rows.Scan(&id, &blob)
		if err != nil {
			return err
		}
		game := &Game{}
		err := json.Unmarshal(blob, game)
		if err != nil {
			return err
		}
		game.ID = id
		s.Games[id] = game
	}
	return nil
}

// New creates a new game and adds it to the store
func (s GameStore) New(players Players, settings Settings) (*Game, error) {
	game := &Game{
		History:   History([]Board{New19by19Board()}),
		Next:      BlackStone,
		Players:   players,
		Settings:  settings,
		Captures:  Captures{0, 0},
		Passes:    Passes{},
		Finished:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	blob, err := json.Marshal(game)
	if err != nil {
		return nil, err
	}
	stmt, err := s.DB.Prepare(`
		INSERT INTO games (blob) VALUES(?)
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(blob)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	game.ID = id
	s.Games[id] = game
	return game, nil
}

// Save a game to persistent storage
func (s GameStore) Save(game *Game) error {
	game.UpdatedAt = time.Now()
	blob, err := json.Marshal(game)
	if err != nil {
		return err
	}
	stmt, err := s.DB.Prepare(`
		UPDATE games SET
			blob = ?,
			finished = ?
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(blob, game.Finished, game.ID)
	return err
}
