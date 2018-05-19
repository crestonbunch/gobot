package gobot

import (
	"errors"
	"fmt"
)

// PipelineFunc handles a pipeline action.
type PipelineFunc = func(*Session, string, *Move) (*Response, error)

// A Pipeline executes a sequence of actions
type Pipeline []PipelineFunc

// Run a pipeline
func (p Pipeline) Run(s *Session, player string, m *Move) (*Response, error) {
	var r *Response
	var err error
	for _, f := range p {
		r, err = f(s, player, m)
		if err != nil {
			return nil, err
		}
	}
	return r, err
}

// InitializerPipeline executes the steps to initialize a game
var InitializerPipeline = Pipeline{
	handleShow,
}

// MovePipeline executes the steps to make a move in a game
var MovePipeline = Pipeline{
	requireUnfinished,
	requireMoving,
	requireAuth,
	requireValid,
	handleMove,
}

// VotePipeline executes the steps to vote for a move in a game
var VotePipeline = Pipeline{
	requireUnfinished,
	requireVoting,
	requireAuth,
	requireValid,
	handleVote,
}

// PlayPipeline executes the steps to pick a random vote in a game
var PlayPipeline = Pipeline{
	handleSchedule,
	handlePlay,
}

// ShowPipeline executes the steps to show a game
var ShowPipeline = Pipeline{
	handleShow,
}

func requireAuth(s *Session, player string, m *Move) (*Response, error) {
	if !s.Playable.CanMove(player) {
		return nil, errors.New("not your turn")
	}
	return nil, nil
}

func requireValid(s *Session, player string, m *Move) (*Response, error) {
	if !s.Game.Validate(m) {
		return nil, errors.New("invalid move")
	}
	return nil, nil
}

func requireVoting(s *Session, player string, m *Move) (*Response, error) {
	if !s.Votable.Required() {
		return nil, errors.New("voting not allowed")
	}
	return nil, nil
}

func requireMoving(s *Session, player string, m *Move) (*Response, error) {
	if s.Votable.Required() {
		return nil, errors.New("please vote for a move")
	}
	return nil, nil
}

func requireUnfinished(s *Session, player string, m *Move) (*Response, error) {
	if s.Game.Finished() {
		return nil, errors.New("game is over")
	}
	return nil, nil
}

func handleSchedule(s *Session, player string, m *Move) (*Response, error) {
	s.Votable.Schedule()
	return nil, nil
}

func handleMove(s *Session, player string, m *Move) (*Response, error) {
	return NewSessionResponse(s, m.String()), s.Game.Move(m)
}

func handleVote(s *Session, player string, m *Move) (*Response, error) {
	return NewTextResponse("thanks for voting"), s.Votable.Vote(m)
}

func handlePlay(s *Session, player string, m *Move) (*Response, error) {
	if s.Votable.Empty() {
		return nil, nil
	}
	vote, err := s.Votable.Random()
	if err != nil {
		return nil, err
	}
	err = s.Votable.Reset()
	if err != nil {
		return nil, err
	}
	details := fmt.Sprintf("voted to %s", vote.String())
	return NewSessionResponse(s, details), s.Game.Move(vote)
}

func handleShow(s *Session, player string, m *Move) (*Response, error) {
	return NewSessionResponse(s, ""), nil
}
