package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

const queryTimeout = 5 * time.Second

type SQLStorage struct {
	db *sql.DB
}

type Movie struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Image       string    `json:"image"`
	Description string    `json:"description"`
}

type Storage interface {
	Get(id uuid.UUID) (Movie, error)
	Store(*Movie) error
	Close() error
}

// NewSQLStorage instantiates an SQLStorage.
func NewSQLStorage(config Configuration) (SQLStorage, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)

	db, err := sql.Open(config.Driver, psqlInfo)

	if err != nil {
		return SQLStorage{}, err
	}

	err = db.Ping()
	if err != nil {
		return SQLStorage{}, err
	}

	return SQLStorage{db}, nil
}

// Get returns a movie from the storage given the ID.
func (s *SQLStorage) Get(id uuid.UUID) (Movie, error) {
	return getMovie(s.db, id)
}

// Store saves a movie in the storage, generating a new ID and updating this field in the movie pointer.
func (s *SQLStorage) Store(m *Movie) error {
	m.Id = uuid.New()
	return insertMovie(s.db, m)
}

// Close ends the connection to the database. This function should be used in order not to leave unclosed connections.
func (s *SQLStorage) Close() error {
	// TODO: Update movie ID
	return s.db.Close()
}

func insertMovie(db *sql.DB, m *Movie) error {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	query := "INSERT INTO movies (id, name, description, image) VALUES ($1, $2, $3, $4)"
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, m.Id, m.Name, m.Description, m.Image)
	return err
}

func getContextWithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), queryTimeout)
}

func getMovie(db *sql.DB, id uuid.UUID) (Movie, error) {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	movie := Movie{}
	query := "SELECT id, name, description, image FROM movies WHERE id=$1"
	err := db.QueryRowContext(ctx, query, id).Scan(&movie.Id, &movie.Name, &movie.Description, &movie.Image)
	return movie, err
}
