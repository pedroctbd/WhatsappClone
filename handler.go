package main

import (
	"context"
	"log"
	"net/http"
	"time"

	gosql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type Application struct {
	RD       *redis.Client
	HUB      *Hub
	CS       *gosql.Session
	ServerID string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (app *Application) routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", app.serveHome)
	r.Get("/ws", app.handleConnections)
	return r
}

func (app *Application) serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func (app *Application) handleConnections(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	userID := chi.URLParam(r, "userId")
	if userID == "" {
		log.Println("User ID is missing, connection rejected")
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: app.HUB, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	if setErr := app.RD.Set(ctx, client.userID, app.ServerID, 5*time.Minute).Err(); setErr != nil {
		log.Printf("Failed to set server for user %s: %v", userID, setErr)
	} else {
		log.Printf("User %s connected to server %s", client.userID, app.ServerID)
	}
	go client.readPump()
	go client.writePump()

}
