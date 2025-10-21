package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pedroctbd/WhatsappClone/internal/domain"
	"github.com/pedroctbd/WhatsappClone/internal/storage"
)

type ChatService struct {
	Repo storage.ChatRepository
}

func (s *ChatService) ProcessMessage(ctx context.Context, senderID string, rawMessage []byte) (*domain.TargetedMessage, error) {
	var userMsg domain.UserMessage
	if err := json.Unmarshal(rawMessage, &userMsg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	//TODO: Add group logic

	chatID := "1"

	messageToSave := &domain.Message{
		ID:       userMsg.ClientMessageID,
		ChatID:   chatID,
		SenderID: senderID,
		Content:  userMsg.Content,
		SentAt:   time.Now(),
	}
	if err := s.Repo.SaveMessage(ctx, *messageToSave); err != nil {
		return nil, err
	}

	deliveryMessage := &domain.TargetedMessage{
		Content:      rawMessage,
		RecipientIDs: nil,
	}
	return deliveryMessage, nil
}
