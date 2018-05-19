package gobot

import (
	"fmt"
	"log"
)

// A Session gets created from a request and a Store in order to make actions.
type Session struct {
	Game     Game
	Playable Playable
	Storable Storable
	Votable  Votable
}

// NewSession creates a new session
func NewSession(g Game, p Playable, s Storable, v Votable) *Session {
	return &Session{g, p, s, v}
}

// Background runs background tasks for a session
func (sess *Session) Background(s Store, ch chan *Response, l *log.Logger) {
	for {
		l.Printf("scheduling vote for %d", sess.Storable.ID())
		sess.Votable.Schedule()
		sess.Votable.Block()
		if sess.Votable.Empty() {
			continue
		}
		move, err := sess.Votable.Random()
		if err != nil {
			ch <- NewTextResponse(err.Error())
			continue
		}
		err = sess.Votable.Reset()
		if err != nil {
			ch <- NewTextResponse(err.Error())
			continue
		}
		err = sess.Game.Move(move)
		if err != nil {
			ch <- NewTextResponse(err.Error())
			continue
		}
		s.Save(sess.Storable)
		details := fmt.Sprintf("voted to %s", move.String())
		ch <- NewSessionResponse(sess, details)
	}
}
