package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

// MovieServer is an HTTP interface
type MovieServer struct {
	storage Storage
}

type getRequest struct {
	Id uuid.UUID `json:"id"`
}

func (s *MovieServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.storeMovie(w, r)
	case http.MethodGet:
		s.getMovie(w, r)
	}
}

func (ms *MovieServer) storeMovie(w http.ResponseWriter, r *http.Request) {
	var movie Movie
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to parse JSON body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if movie.Name == "" {
		http.Error(w, `"name" field is mandatory`, http.StatusBadRequest)
		return
	}

	err = ms.storage.Store(&movie)
	if err != nil {
		log.Println("movie could not be stored:", err)
		http.Error(w, "movie could not be stored", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(movie)
	if err != nil {
		log.Println("cannot encode response:", err)
		http.Error(w, "cannot encode response", http.StatusInternalServerError)
		return
	}
}

func (ms *MovieServer) getMovie(w http.ResponseWriter, r *http.Request) {
	var req getRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to parse JSON body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if req.Id == uuid.Nil {
		http.Error(w, `"id" must be different than 0`, http.StatusBadRequest)
		return
	}

	movie, err := ms.storage.Get(req.Id)
	switch {
	case err == sql.ErrNoRows:
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	case err != nil:
		log.Println("cannot retrieve movie:", err)
		http.Error(w, "cannot retrieve movie", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(movie)
	if err != nil {
		log.Println("cannot encode response:", err)
		http.Error(w, "cannot encode response", http.StatusInternalServerError)
		return
	}
}
