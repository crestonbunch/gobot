// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gobot "github.com/crestonbunch/gobot"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	time "time"
)

// MockVotable is a mock of Votable interface
type MockVotable struct {
	ctrl     *gomock.Controller
	recorder *MockVotableMockRecorder
}

// MockVotableMockRecorder is the mock recorder for MockVotable
type MockVotableMockRecorder struct {
	mock *MockVotable
}

// NewMockVotable creates a new mock instance
func NewMockVotable(ctrl *gomock.Controller) *MockVotable {
	mock := &MockVotable{ctrl: ctrl}
	mock.recorder = &MockVotableMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockVotable) EXPECT() *MockVotableMockRecorder {
	return m.recorder
}

// Vote mocks base method
func (m *MockVotable) Vote(arg0 *gobot.Move) error {
	ret := m.ctrl.Call(m, "Vote", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Vote indicates an expected call of Vote
func (mr *MockVotableMockRecorder) Vote(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Vote", reflect.TypeOf((*MockVotable)(nil).Vote), arg0)
}

// Schedule mocks base method
func (m *MockVotable) Schedule() *time.Timer {
	ret := m.ctrl.Call(m, "Schedule")
	ret0, _ := ret[0].(*time.Timer)
	return ret0
}

// Schedule indicates an expected call of Schedule
func (mr *MockVotableMockRecorder) Schedule() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Schedule", reflect.TypeOf((*MockVotable)(nil).Schedule))
}

// Random mocks base method
func (m *MockVotable) Random() (*gobot.Move, error) {
	ret := m.ctrl.Call(m, "Random")
	ret0, _ := ret[0].(*gobot.Move)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Random indicates an expected call of Random
func (mr *MockVotableMockRecorder) Random() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Random", reflect.TypeOf((*MockVotable)(nil).Random))
}

// Empty mocks base method
func (m *MockVotable) Empty() bool {
	ret := m.ctrl.Call(m, "Empty")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Empty indicates an expected call of Empty
func (mr *MockVotableMockRecorder) Empty() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Empty", reflect.TypeOf((*MockVotable)(nil).Empty))
}

// Reset mocks base method
func (m *MockVotable) Reset() error {
	ret := m.ctrl.Call(m, "Reset")
	ret0, _ := ret[0].(error)
	return ret0
}

// Reset indicates an expected call of Reset
func (mr *MockVotableMockRecorder) Reset() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reset", reflect.TypeOf((*MockVotable)(nil).Reset))
}

// MockStorable is a mock of Storable interface
type MockStorable struct {
	ctrl     *gomock.Controller
	recorder *MockStorableMockRecorder
}

// MockStorableMockRecorder is the mock recorder for MockStorable
type MockStorableMockRecorder struct {
	mock *MockStorable
}

// NewMockStorable creates a new mock instance
func NewMockStorable(ctrl *gomock.Controller) *MockStorable {
	mock := &MockStorable{ctrl: ctrl}
	mock.recorder = &MockStorableMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStorable) EXPECT() *MockStorableMockRecorder {
	return m.recorder
}

// ID mocks base method
func (m *MockStorable) ID() int64 {
	ret := m.ctrl.Call(m, "ID")
	ret0, _ := ret[0].(int64)
	return ret0
}

// ID indicates an expected call of ID
func (mr *MockStorableMockRecorder) ID() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ID", reflect.TypeOf((*MockStorable)(nil).ID))
}

// Load mocks base method
func (m *MockStorable) Load(arg0 []byte) error {
	ret := m.ctrl.Call(m, "Load", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Load indicates an expected call of Load
func (mr *MockStorableMockRecorder) Load(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Load", reflect.TypeOf((*MockStorable)(nil).Load), arg0)
}

// Save mocks base method
func (m *MockStorable) Save() ([]byte, error) {
	ret := m.ctrl.Call(m, "Save")
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Save indicates an expected call of Save
func (mr *MockStorableMockRecorder) Save() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockStorable)(nil).Save))
}

// MockPlayable is a mock of Playable interface
type MockPlayable struct {
	ctrl     *gomock.Controller
	recorder *MockPlayableMockRecorder
}

// MockPlayableMockRecorder is the mock recorder for MockPlayable
type MockPlayableMockRecorder struct {
	mock *MockPlayable
}

// NewMockPlayable creates a new mock instance
func NewMockPlayable(ctrl *gomock.Controller) *MockPlayable {
	mock := &MockPlayable{ctrl: ctrl}
	mock.recorder = &MockPlayableMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPlayable) EXPECT() *MockPlayableMockRecorder {
	return m.recorder
}

// IsPlaying mocks base method
func (m *MockPlayable) IsPlaying(playerID string) bool {
	ret := m.ctrl.Call(m, "IsPlaying", playerID)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsPlaying indicates an expected call of IsPlaying
func (mr *MockPlayableMockRecorder) IsPlaying(playerID interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsPlaying", reflect.TypeOf((*MockPlayable)(nil).IsPlaying), playerID)
}

// CanMove mocks base method
func (m *MockPlayable) CanMove(playerID string) bool {
	ret := m.ctrl.Call(m, "CanMove", playerID)
	ret0, _ := ret[0].(bool)
	return ret0
}

// CanMove indicates an expected call of CanMove
func (mr *MockPlayableMockRecorder) CanMove(playerID interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CanMove", reflect.TypeOf((*MockPlayable)(nil).CanMove), playerID)
}

// MockGame is a mock of Game interface
type MockGame struct {
	ctrl     *gomock.Controller
	recorder *MockGameMockRecorder
}

// MockGameMockRecorder is the mock recorder for MockGame
type MockGameMockRecorder struct {
	mock *MockGame
}

// NewMockGame creates a new mock instance
func NewMockGame(ctrl *gomock.Controller) *MockGame {
	mock := &MockGame{ctrl: ctrl}
	mock.recorder = &MockGameMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockGame) EXPECT() *MockGameMockRecorder {
	return m.recorder
}

// Board mocks base method
func (m *MockGame) Board() gobot.Board {
	ret := m.ctrl.Call(m, "Board")
	ret0, _ := ret[0].(gobot.Board)
	return ret0
}

// Board indicates an expected call of Board
func (mr *MockGameMockRecorder) Board() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Board", reflect.TypeOf((*MockGame)(nil).Board))
}

// Finished mocks base method
func (m *MockGame) Finished() bool {
	ret := m.ctrl.Call(m, "Finished")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Finished indicates an expected call of Finished
func (mr *MockGameMockRecorder) Finished() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Finished", reflect.TypeOf((*MockGame)(nil).Finished))
}

// Move mocks base method
func (m *MockGame) Move(arg0 *gobot.Move) error {
	ret := m.ctrl.Call(m, "Move", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Move indicates an expected call of Move
func (mr *MockGameMockRecorder) Move(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Move", reflect.TypeOf((*MockGame)(nil).Move), arg0)
}

// Validate mocks base method
func (m *MockGame) Validate(arg0 *gobot.Move) bool {
	ret := m.ctrl.Call(m, "Validate", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Validate indicates an expected call of Validate
func (mr *MockGameMockRecorder) Validate(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validate", reflect.TypeOf((*MockGame)(nil).Validate), arg0)
}

// MockStore is a mock of Store interface
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// Close mocks base method
func (m *MockStore) Close() error {
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockStoreMockRecorder) Close() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStore)(nil).Close))
}

// Load mocks base method
func (m *MockStore) Load() error {
	ret := m.ctrl.Call(m, "Load")
	ret0, _ := ret[0].(error)
	return ret0
}

// Load indicates an expected call of Load
func (mr *MockStoreMockRecorder) Load() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Load", reflect.TypeOf((*MockStore)(nil).Load))
}

// Get mocks base method
func (m *MockStore) Get(id int64) (*gobot.Session, error) {
	ret := m.ctrl.Call(m, "Get", id)
	ret0, _ := ret[0].(*gobot.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockStoreMockRecorder) Get(id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockStore)(nil).Get), id)
}

// New mocks base method
func (m *MockStore) New(arg0 gobot.Blueprint) (*gobot.Session, error) {
	ret := m.ctrl.Call(m, "New", arg0)
	ret0, _ := ret[0].(*gobot.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// New indicates an expected call of New
func (mr *MockStoreMockRecorder) New(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "New", reflect.TypeOf((*MockStore)(nil).New), arg0)
}

// Last mocks base method
func (m *MockStore) Last() (*gobot.Session, error) {
	ret := m.ctrl.Call(m, "Last")
	ret0, _ := ret[0].(*gobot.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Last indicates an expected call of Last
func (mr *MockStoreMockRecorder) Last() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Last", reflect.TypeOf((*MockStore)(nil).Last))
}

// Save mocks base method
func (m *MockStore) Save(arg0 gobot.Storable) error {
	ret := m.ctrl.Call(m, "Save", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save
func (mr *MockStoreMockRecorder) Save(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockStore)(nil).Save), arg0)
}

// List mocks base method
func (m *MockStore) List(all bool) ([]*gobot.Session, error) {
	ret := m.ctrl.Call(m, "List", all)
	ret0, _ := ret[0].([]*gobot.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockStoreMockRecorder) List(all interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockStore)(nil).List), all)
}
