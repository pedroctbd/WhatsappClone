package chat

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/pedroctbd/WhatsappClone/internal/domain"
	"github.com/pedroctbd/WhatsappClone/internal/storage"
)

type ChatService struct {
	Repo storage.ChatRepository
}

func (s *ChatService) ProcessMessage(ctx context.Context, senderIDStr string, rawMessage []byte) (*domain.TargetedMessage, error) {
	var userMsg domain.UserMessage
	if err := json.Unmarshal(rawMessage, &userMsg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	senderID, err := uuid.Parse(senderIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid sender ID: %w", err)
	}

	recipientID, err := uuid.Parse(userMsg.RecipientID)
	if err != nil {
		return nil, fmt.Errorf("invalid recipient ID: %w", err)
	}

	chatID, participantsUUIDs, err := s.Repo.GetOrCreateOneOnOneChat(ctx, senderID, recipientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create chat: %w", err)
	}

	sentAtUUID, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("failed to generate V1 uuid: %w", err)
	}
	newMsg := domain.Message{
		ChatID:     chatID,
		MessageID:  uuid.MustParse(userMsg.ClientMessageID),
		SenderID:   senderID,
		SenderName: "Tester1",
		SentAt:     sentAtUUID,
		Content:    userMsg.Content,
	}

	if err := s.Repo.SaveMessage(ctx, newMsg); err != nil {
		return nil, fmt.Errorf("failed to save message: %w", err)
	}

	var participantIDs []string
	for _, id := range participantsUUIDs {
		participantIDs = append(participantIDs, id.String())
	}

	deliveryMessage := &domain.TargetedMessage{
		Content:      rawMessage,
		RecipientIDs: participantIDs,
	}
	return deliveryMessage, nil
}
