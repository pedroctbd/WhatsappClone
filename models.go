package main

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                uuid.UUID `json:"id"`
	PhoneNumber       string    `json:"phone_number"`
	DisplayName       string    `json:"display_name"`
	ProfilePictureURL string    `json:"profile_picture_url,omitempty"`
	AboutText         string    `json:"about_text,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	LastSeenAt        time.Time `json:"last_seen_at,omitempty"`
}

type Chat struct {
	ID           uuid.UUID   `json:"id"`
	IsGroup      bool        `json:"is_group"`
	Participants []uuid.UUID `json:"participants"`
	CreatedAt    time.Time   `json:"created_at"`
	GroupName    string      `json:"group_name,omitempty"`
	GroupAdminID uuid.UUID   `json:"group_admin_id,omitempty"`
	GroupIconURL string      `json:"group_icon_url,omitempty"`
}

type Message struct {
	ID       uuid.UUID `json:"id"`
	ChatID   uuid.UUID `json:"chat_id"`
	SenderID uuid.UUID `json:"sender_id"`
	Content  string    `json:"content"`
	SentAt   time.Time `json:"sent_at"`
	Status   string    `json:"status"`
}
