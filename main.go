package main

import (
	"log"
	"net/http"
)

func main() {
	storage, err := NewSQLStorage(LoadConfig())
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}
	defer storage.Close()

	server := &MovieServer{&storage}

	mux := http.NewServeMux()
	mux.Handle("/movies/", server)

	log.Println("Starting server")
	err = http.ListenAndServe(":8000", server)
	if err != nil {
		log.Println("error serving:", err)
	}
}
