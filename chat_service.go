package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ChatService struct {
	Repo ChatRepository
}

func (s *ChatService) ProcessMessage(ctx context.Context, senderID string, rawMessage []byte) (*TargetedMessage, error) {
	var userMsg UserMessage
	if err := json.Unmarshal(rawMessage, &userMsg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	//TODO: Add group logic

	chatID, participants, err := s.Repo.GetOrCreateOneOnOneChat(ctx, senderID, userMsg.RecipientID)
	if err != nil {
		return nil, err
	}

	messageToSave := &Message{
		ID:       uuid.MustParse(userMsg.ClientMessageID),
		ChatID:   chatID,
		SenderID: uuid.MustParse(senderID),
		Content:  userMsg.Content,
		SentAt:   time.Now(),
	}
	if err := s.Repo.SaveMessage(ctx, messageToSave); err != nil {
		return nil, err
	}

	deliveryMessage := &TargetedMessage{
		Content:      rawMessage,
		RecipientIDs: participants,
	}
	return deliveryMessage, nil
}
