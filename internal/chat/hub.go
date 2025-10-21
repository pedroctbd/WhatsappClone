package chat

import (
	"context"
	"log"

	"github.com/pedroctbd/WhatsappClone/internal/domain"
	"github.com/redis/go-redis/v9"
)

// Hub maintains the set of active Clients and broadcasts messages.
type Hub struct {
	Clients     map[string]*Client
	Delivery    chan *domain.TargetedMessage
	Register    chan *Client
	Unregister  chan *Client
	RedisClient *redis.Client
}

func NewHub(redisC *redis.Client) *Hub {
	return &Hub{
		Clients:     make(map[string]*Client),
		Delivery:    make(chan *domain.TargetedMessage),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		RedisClient: redisC,
	}
}

func (h *Hub) Run() {
	ctx := context.Background()
	for {
		select {
		case client := <-h.Register:
			// Use the client's UserID as the key.
			h.Clients[client.UserID] = client
			log.Printf("User %s registered to hub", client.UserID)

		case client := <-h.Unregister:
			if _, ok := h.Clients[client.UserID]; ok {
				delete(h.Clients, client.UserID)
				close(client.Send)

				// Use the correct variable 'client' (lowercase).
				if err := h.RedisClient.Del(ctx, client.UserID).Err(); err != nil {
					log.Printf("Failed to delete user %s from Redis: %v", client.UserID, err)
				} else {
					log.Printf("User %s disconnected and removed from Redis", client.UserID)
				}
			}

		case message := <-h.Delivery:
			for _, recipientID := range message.RecipientIDs {
				if client, ok := h.Clients[recipientID]; ok {
					select {
					case client.Send <- message.Content:
					default:
						close(client.Send)
						delete(h.Clients, client.UserID)
					}
				}
			}
		}
	}
}
