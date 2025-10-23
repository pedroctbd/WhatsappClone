package deliveryHttp

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/pedroctbd/WhatsappClone/internal/chat"
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

func (app Application) handleConnections(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	UserID := chi.URLParam(r, "userID")
	if UserID == "" {
		log.Println("User ID is missing, connection rejected")
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &chat.Client{
		Hub:         app.Hub,
		Conn:        conn,
		Send:        make(chan []byte, 256),
		UserID:      UserID,
		ChatService: app.ChatService,
	}
	client.Hub.Register <- client

	if err := client.Hub.RedisClient.Set(ctx, client.UserID, app.ServerID, 5*time.Minute).Err(); err != nil {
		log.Printf("Failed to set server for user %s: %v", client.UserID, err)
	} else {
		log.Printf("User %s connected to server %s", client.UserID, app.ServerID)
	}

	go client.WritePump()
	go client.ReadPump()
}
