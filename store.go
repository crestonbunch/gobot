package gobot

import (
	"database/sql"
	"encoding/json"
	"errors"
	"sort"
	"time"
)

// StateList is a list of states in some order
type StateList []*State

// Len is part of sort.Interface
func (l StateList) Len() int {
	return len(l)
}

// Swap is part of sort.Interface
func (l StateList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

// Less is part of sort.Interface
func (l StateList) Less(i, j int) bool {
	return l[i].UpdatedAt.After(l[j].UpdatedAt)
}

// StateStore stores a map of game IDs to corresponding games states in SQLite
type StateStore struct {
	DB *sql.DB
}

// NewGameStore creates a game store connected to an SQLite database for
//  games.
func NewGameStore(db *sql.DB) *StateStore {
	return &StateStore{
		DB: db,
	}
}

// Load persistent storage
func (s *StateStore) Load() error {
	stmt := `
	CREATE TABLE IF NOT EXISTS
	games
	(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		blob BLOB NOT NULL
	)
	`
	_, err := s.DB.Exec(stmt)
	return err
}

// New creates a new Game and add it to the store
func (s *StateStore) New(bp Blueprint) (*Session, error) {
	game := &State{
		History:   History([]Board{New19by19Board()}),
		Next:      BlackStone,
		Players:   Players(bp.Players),
		Voting:    Voting(bp.Voting),
		Captures:  Captures{0, 0},
		Passes:    Passes{},
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
	game.id = id
	return NewSession(game, game, game, game), nil
}

// Get a game by id
func (s *StateStore) Get(id int64) (*Session, error) {
	stmt, err := s.DB.Prepare(`SELECT blob FROM games WHERE id = ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	var blob []byte
	err = row.Scan(&blob)
	if err != nil {
		return nil, err
	}
	game := &State{}
	err = json.Unmarshal(blob, game)
	if err != nil {
		return nil, err
	}
	game.id = id
	return NewSession(game, game, game, game), nil
}

// Last returns the last state played
func (s *StateStore) Last() (*Session, error) {
	sessions, err := s.List(false)
	if err != nil {
		return nil, err
	}
	if len(sessions) == 0 {
		return nil, errors.New("no active games")
	}
	return sessions[0], nil
}

// Save a game to persistent storage
func (s *StateStore) Save(storable Storable) error {
	storable.(*State).UpdatedAt = time.Now()
	blob, err := storable.Save()
	if err != nil {
		return err
	}
	stmt, err := s.DB.Prepare(`UPDATE games SET blob = ? WHERE id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(blob, storable.ID())
	return err
}

// List gets a list of the games in the store
func (s *StateStore) List(all bool) ([]*Session, error) {
	list, err := s.listAll()
	if err != nil {
		return nil, err
	}
	sort.Sort(list)
	output := []*Session{}
	for _, state := range list {
		if all || !state.Finished() {
			sess := NewSession(state, state, state, state)
			output = append(output, sess)
		}
	}
	return output, nil
}

func (s *StateStore) listAll() (StateList, error) {
	list := StateList([]*State{})
	rows, err := s.DB.Query("SELECT id, blob FROM games")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var blob []byte
		err = rows.Scan(&id, &blob)
		if err != nil {
			return nil, err
		}
		game := &State{}
		err := json.Unmarshal(blob, game)
		if err != nil {
			return nil, err
		}
		game.id = id
		list = append(list, game)
	}
	return StateList(list), nil
}

// Close closes the store
func (s *StateStore) Close() error {
	return s.DB.Close()
}
