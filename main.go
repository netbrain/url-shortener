package main

import (
	"log"
	"net/http"
)

func main() {
	storage, err := NewStorage("data.db")
	if err != nil {
		log.Fatalf("Could not initialize storage: %v", err)
	}
	defer storage.Close()

	app := NewApp(storage)

	http.HandleFunc("/", app.shortenURL)
	http.HandleFunc("/{shortKey}", app.redirectURL)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
