package main

import (
	"html/template"
	"net/http"
	"sync"

	"github.com/teris-io/shortid"
)

var tpl = template.Must(template.ParseFiles("templates/form.html"))

type App struct {
	storage *Storage
	mu      sync.Mutex
}

func NewApp(storage *Storage) *App {
	return &App{storage: storage}
}

func (app *App) shortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		url := r.FormValue("url")

		// Generate a short key
		app.mu.Lock()
		shortKey, err := shortid.Generate()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		app.storage.SaveEntry(shortKey, url)
		app.mu.Unlock()

		tpl.Execute(w, map[string]string{"ShortURL": shortKey})
	} else {
		tpl.Execute(w, nil)
	}
}

func (app *App) redirectURL(w http.ResponseWriter, r *http.Request) {
	shortKey := r.URL.Path[1:]

	url := app.storage.Get(shortKey)
	if url == "" {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}
