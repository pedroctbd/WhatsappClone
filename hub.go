// hub.go
package main

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

// Hub maintains the set of active clients and broadcasts messages.
type Hub struct {
	clients     map[string]*Client
	delivery    chan *TargetedMessage
	register    chan *Client
	unregister  chan *Client
	redisClient *redis.Client
}

func newHub(redisC *redis.Client) *Hub {
	return &Hub{
		clients:     make(map[string]*Client),
		delivery:    make(chan *TargetedMessage),
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
			// Use the client's userID as the key.
			h.clients[client.userID] = client
			log.Printf("User %s registered to hub", client.userID)

		case client := <-h.unregister:
			if _, ok := h.clients[client.userID]; ok {
				delete(h.clients, client.userID)
				close(client.send)

				// Use the correct variable 'client' (lowercase).
				if err := h.redisClient.Del(ctx, client.userID).Err(); err != nil {
					log.Printf("Failed to delete user %s from Redis: %v", client.userID, err)
				} else {
					log.Printf("User %s disconnected and removed from Redis", client.userID)
				}
			}

		case message := <-h.delivery:
			for _, recipientID := range message.RecipientIDs {
				if client, ok := h.clients[recipientID]; ok {
					select {
					case client.send <- message.Content:
					default:
						close(client.send)
						delete(h.clients, client.userID)
					}
				}
			}
		}
	}
}
