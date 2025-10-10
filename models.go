// models.go
package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// User websocket connection
type Client struct {
	hub         *Hub
	conn        *websocket.Conn
	send        chan []byte
	userID      string
	chatService *ChatService
}

// Message decoded from json
type UserMessage struct {
	ClientMessageID string `json:"clientMessageId"`
	RecipientID     string `json:"recipientId,omitempty"`
	ChatID          string `json:"chatId,omitempty"`
	Content         string `json:"content"`
}

// Message that client sends to hub
type TargetedMessage struct {
	Content      []byte
	RecipientIDs []string
}

// Stored in cassandra
type Message struct {
	ID       uuid.UUID
	ChatID   uuid.UUID
	SenderID uuid.UUID
	Content  string
	SentAt   time.Time
}
