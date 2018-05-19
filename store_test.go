package gobot_test

import (
	"database/sql"
	"testing"

	. "github.com/crestonbunch/gobot"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestSQLiteLoad(t *testing.T) {
	cases := []struct {
		setup func() (*sql.DB, sqlmock.Sqlmock)
	}{
		{
			setup: func() (*sql.DB, sqlmock.Sqlmock) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf(err.Error())
				}
				mock.ExpectExec("CREATE TABLE IF NOT EXISTS.+").
					WillReturnResult(sqlmock.NewResult(0, 0))
				return db, mock
			},
		},
	}
	for _, test := range cases {
		db, mock := test.setup()
		defer db.Close()
		store := NewGameStore(db)
		err := store.Load()
		if err != nil {
			t.Errorf(err.Error())
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf(err.Error())
		}
	}
}

func TestSQLiteNew(t *testing.T) {
	cases := []struct {
		bp    Blueprint
		setup func(Blueprint) (*sql.DB, sqlmock.Sqlmock)
	}{
		{
			setup: func(bp Blueprint) (*sql.DB, sqlmock.Sqlmock) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf(err.Error())
				}
				stmt := mock.ExpectPrepare("INSERT INTO.+")
				stmt.ExpectExec().
					WithArgs(sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				return db, mock
			},
		},
	}
	for _, test := range cases {
		db, mock := test.setup(test.bp)
		defer db.Close()
		store := NewGameStore(db)
		_, err := store.New(test.bp)
		if err != nil {
			t.Errorf(err.Error())
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf(err.Error())
		}
	}
}

func TestSQLiteGet(t *testing.T) {
	cases := []struct {
		id    int64
		setup func(int64) (*sql.DB, sqlmock.Sqlmock)
	}{
		{
			id: 1,
			setup: func(id int64) (*sql.DB, sqlmock.Sqlmock) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf(err.Error())
				}
				stmt := mock.ExpectPrepare("SELECT.+")
				stmt.ExpectQuery().WillReturnRows(
					sqlmock.NewRows([]string{"blob"}).AddRow("{}"),
				).WithArgs(id)
				return db, mock
			},
		},
	}
	for _, test := range cases {
		db, mock := test.setup(test.id)
		defer db.Close()
		store := NewGameStore(db)
		_, err := store.Get(test.id)
		if err != nil {
			t.Errorf(err.Error())
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf(err.Error())
		}
	}
}

func TestSQLiteSave(t *testing.T) {
	cases := []struct {
		state *State
		setup func(*State) (*sql.DB, sqlmock.Sqlmock)
	}{
		{
			state: &State{},
			setup: func(state *State) (*sql.DB, sqlmock.Sqlmock) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf(err.Error())
				}
				stmt := mock.ExpectPrepare("UPDATE.+")
				stmt.ExpectExec().
					WillReturnResult(sqlmock.NewResult(0, 1)).
					WithArgs(sqlmock.AnyArg(), state.ID())
				return db, mock
			},
		},
	}
	for _, test := range cases {
		db, mock := test.setup(test.state)
		defer db.Close()
		store := NewGameStore(db)
		err := store.Save(test.state)
		if err != nil {
			t.Errorf(err.Error())
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf(err.Error())
		}
	}
}

func TestSQLiteList(t *testing.T) {
	cases := []struct {
		all   bool
		setup func(bool) (*sql.DB, sqlmock.Sqlmock)
	}{
		{
			all: true,
			setup: func(all bool) (*sql.DB, sqlmock.Sqlmock) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf(err.Error())
				}
				mock.ExpectQuery("SELECT.+").
					WillReturnRows(
						sqlmock.
							NewRows([]string{"id", "blob"}).
							AddRow(1, "{}"),
					)
				return db, mock
			},
		},
	}
	for _, test := range cases {
		db, mock := test.setup(test.all)
		defer db.Close()
		store := NewGameStore(db)
		_, err := store.List(test.all)
		if err != nil {
			t.Errorf(err.Error())
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf(err.Error())
		}
	}
}
