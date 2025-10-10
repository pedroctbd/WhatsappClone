package main

import (
	"context"

	"github.com/google/uuid"
)

type ChatRepository interface {
	GetOrCreateOneOnOneChat(ctx context.Context, userID1, userID2 string) (chatID uuid.UUID, participants []string, err error)
	SaveMessage(ctx context.Context, msg *Message) error
}
