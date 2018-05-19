package gobot

import "time"

// A Request connects a user command to a Store to perform an action.
type Request struct {
	Store   Store
	Command Command
	Session *Session
	Player  string
}

// NewRequest constructs a request from a user command and session store
func NewRequest(cmd Command, player string, str Store) (*Request, error) {
	var sess *Session
	var err error
	switch cmd := cmd.(type) {
	case *StartCommand:
		b := Blueprint{
			Players: Players{
				Anyone: cmd.Anyone,
				Black:  cmd.Black,
				White:  cmd.White,
			},
			Voting: Voting{
				Required: cmd.Anyone,         // require voting if anyone can play
				Duration: 3600 * time.Second, // select a vote every hour
			},
		}
		sess, err = str.New(b)
	case *MoveCommand:
		sess, err = cmd.Locator.Find(str)
	case *VoteCommand:
		sess, err = cmd.Locator.Find(str)
	case *PlayCommand:
		sess, err = cmd.Locator.Find(str)
	}
	if err != nil {
		return nil, err
	}

	return &Request{
		Store:   str,
		Command: cmd,
		Session: sess,
		Player:  player,
	}, nil
}

// Respond to a request
func (r *Request) Respond() (*Response, error) {
	err := r.Store.Save(r.Session.Storable)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
