package storage

import (
	"context"

	"github.com/pedroctbd/WhatsappClone/internal/domain"
)

type ChatRepository interface {
	GetOrCreateOneOnOneChat(ctx context.Context, userID1, userID2 string) error
	SaveMessage(ctx context.Context, msg domain.Message) error
}
