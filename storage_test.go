package main

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func TestGet(t *testing.T) {
	db, mock := NewSQLMock(t)

	sut := SQLStorage{db}

	id := uuid.New()
	wantMovie := Movie{
		Id:          id,
		Name:        "test",
		Description: "test description",
		Image:       "test image",
	}

	queryGetMovie := "SELECT id, name, description, image FROM movies WHERE id=\\$1"

	rows := sqlmock.NewRows([]string{"id", "name", "description", "image"}).
		AddRow(wantMovie.Id, wantMovie.Name, wantMovie.Description, wantMovie.Image)

	mock.ExpectQuery(queryGetMovie).WithArgs(wantMovie.Id).WillReturnRows(rows)

	gotMovie, err := sut.Get(wantMovie.Id)
	if err != nil {
		t.Error("unexpected error when getting movie:", err)
	}

	if wantMovie != gotMovie {
		t.Errorf("got %v want %v", gotMovie, wantMovie)
	}
}

func TestStore(t *testing.T) {
	queryInsert := "INSERT INTO movies \\(id, name, description, image\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\)"

	t.Run("no errors in DB", func(t *testing.T) {
		db, mock := NewSQLMock(t)

		sut := SQLStorage{db}

		movie := Movie{
			Name:        "test",
			Description: "test description",
			Image:       "test image",
		}

		prep := mock.ExpectPrepare(queryInsert)
		prep.ExpectExec().WithArgs(sqlmock.AnyArg(), movie.Name, movie.Description, movie.Image).WillReturnResult(sqlmock.NewResult(0, 1))

		err := sut.Store(&movie)
		if err != nil {
			t.Error("unexpected error when storing movie:", err)
		}

		if movie.Id == uuid.Nil {
			t.Error("movie ID should not null")
		}
	})
}

func TestClose(t *testing.T) {
	t.Run("no error when closing the connection", func(t *testing.T) {
		db, mock := NewSQLMock(t)
		mock.ExpectClose()

		sut := SQLStorage{db}

		err := sut.Close()
		if err != nil {
			t.Error("unexpected error when closing db connection:", err)
		}
	})
	t.Run("error when closing the connection", func(t *testing.T) {
		db, mock := NewSQLMock(t)
		mock.ExpectClose().WillReturnError(errors.New("cannot close the connection"))

		sut := SQLStorage{db}

		err := sut.Close()
		if err == nil {
			t.Error("expected an error when closing db connection")
		}
	})
}

func NewSQLMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}
