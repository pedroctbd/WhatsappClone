package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/pedroctbd/WhatsappClone/internal/domain"
)

type ChatRepository interface {
	GetOrCreateOneOnOneChat(ctx context.Context, userA uuid.UUID, userB uuid.UUID) (uuid.UUID, []uuid.UUID, error)
	SaveMessage(ctx context.Context, msg domain.Message) error
}
