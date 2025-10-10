package main

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	redisClient *redis.Client
}

func newHub(redisC *redis.Client) *Hub {
	return &Hub{
		clients:     make(map[*Client]bool),
		broadcast:   make(chan []byte),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		redisClient: redisC,
	}
}

func (h *Hub) run() {
	ctx := context.Background()

	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {

				delete(h.clients, client)
				close(client.send)

				if err := h.redisClient.Del(ctx, client.userID).Err(); err != nil {
					log.Printf("Failed to delete user %s from Redis: %v", client.userID, err)
				} else {
					log.Printf("User %s disconnected and removed from Redis", client.userID)
				}
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
