package main

import (
	"log"
	"net/http"

	gosql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type Application struct {
	RD  *redis.Client
	HUB *Hub
	CS  *gosql.Session
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
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: app.HUB, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	go client.readPump()
	go client.writePump()

}
