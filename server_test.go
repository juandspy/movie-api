package main

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

type testCase struct {
	name         string
	storageError error
	body         string
	wantStatus   int
}

func TestGetMovie(t *testing.T) {
	id := uuid.New()
	testCases := []testCase{
		{
			name:       "movie is found in storage",
			body:       fmt.Sprintf(`{"id":%q}`, id),
			wantStatus: http.StatusOK,
		},
		{
			name:         "movie is not in storage",
			storageError: sql.ErrNoRows,
			body:         fmt.Sprintf(`{"id":%q}`, id),
			wantStatus:   http.StatusNotFound,
		},
		{
			name:       "id is not set",
			body:       `{"not_an_id":1}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:         "unknown error",
			storageError: errors.New("unknown error"),
			body:         fmt.Sprintf(`{"id":%q}`, id),
			wantStatus:   http.StatusInternalServerError,
		},
		{
			name:       "bad JSON",
			body:       `id:1`,
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storage := MockStorage{
				err: tc.storageError,
			}
			server := MovieServer{&storage}

			r := httptest.NewRequest(http.MethodGet, "/movies", setBody(tc.body))
			w := httptest.NewRecorder()
			server.ServeHTTP(w, r)

			if w.Code != tc.wantStatus {
				t.Errorf("expected %d got %d", w.Code, tc.wantStatus)
			}
		})
	}
}

func TestStoreMovie(t *testing.T) {
	testCases := []testCase{
		{
			name:       "body is OK",
			body:       `{"name":"test", "description": "test description", "image": "path to image"}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "body is not OK",
			body:       `"name":"test"`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "name is not in JSON",
			body:       `{"not_a_name":"test", "description": "test description", "image": "path to image"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:         "error storing movie",
			body:         `{"name":"test"}`,
			storageError: errors.New("unknown error"),
			wantStatus:   http.StatusInternalServerError,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storage := MockStorage{
				err: tc.storageError,
			}
			server := MovieServer{&storage}

			r := httptest.NewRequest(http.MethodPost, "/movies", setBody(tc.body))
			w := httptest.NewRecorder()
			server.ServeHTTP(w, r)

			if w.Code != tc.wantStatus {
				t.Errorf("expected %d got %d", w.Code, tc.wantStatus)
			}
		})
	}
}

type MockStorage struct {
	movie Movie
	err   error
	newId uuid.UUID
}

func (s *MockStorage) Get(id uuid.UUID) (Movie, error) {
	return s.movie, s.err
}

func (s *MockStorage) Store(m *Movie) error {
	m.Id = s.newId
	return s.err
}

func (s *MockStorage) Close() error {
	return s.err
}

func setBody(body string) *bytes.Reader {
	jsonBody := []byte(body)
	bodyReader := bytes.NewReader(jsonBody)
	return bodyReader
}
