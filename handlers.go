package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func serveHome(w http.ResponseWriter, r *http.Request) {

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
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		hub:         app.Hub,
		conn:        conn,
		send:        make(chan []byte, 256),
		userID:      userID,
		chatService: app.ChatService,
	}
	client.hub.register <- client

	if err := client.hub.redisClient.Set(ctx, client.userID, app.ServerID, 5*time.Minute).Err(); err != nil {
		log.Printf("Failed to set server for user %s: %v", client.userID, err)
	} else {
		log.Printf("User %s connected to server %s", client.userID, app.ServerID)
	}

	go client.writePump()
	go client.readPump()
}
