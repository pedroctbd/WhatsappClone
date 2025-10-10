package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// CheckOrigin is not recommended for production
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (app *Application) routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/", serveHome)
	// Route now correctly captures the userID from the path.
	r.Get("/ws/{userId}", app.handleConnections)
	return r
}
