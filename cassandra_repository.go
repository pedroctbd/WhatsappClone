package main

import (
	"context"
	"fmt"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/google/uuid"
)

type CassandraRepository struct {
	Session *gocql.Session
}

func (r *CassandraRepository) GetOrCreateOneOnOneChat(ctx context.Context, userID1, userID2 string) (uuid.UUID, []string, error) {
	var chatID uuid.UUID
	participants := []string{userID1, userID2}

	err := r.Session.Query(`SELECT chat_id FROM chats WHERE participants CONTAINS ? AND participants CONTAINS ? AND is_group = ?`,
		userID1, userID2, false).Scan(&chatID)

	if err == gocql.ErrNotFound {
		chatID = uuid.New()
		err = r.Session.Query(`INSERT INTO chats (chat_id, is_group, participants) VALUES (?, ?, ?)`,
			chatID, false, participants).ExecContext(ctx)
		if err != nil {
			return uuid.Nil, nil, fmt.Errorf("failed to create new chat: %w", err)
		}
		return chatID, participants, nil
	} else if err != nil {
		return uuid.Nil, nil, fmt.Errorf("failed to retrieve chat: %w", err)
	}

	return chatID, participants, nil
}

func (r *CassandraRepository) SaveMessage(ctx context.Context, msg *Message) error {
	err := r.Session.Query(`INSERT INTO messages (message_id, chat_id, sender_id, content, sent_at) VALUES (?, ?, ?, ?, ?)`,
		msg.ID, msg.ChatID, msg.SenderID, msg.Content, msg.SentAt).ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}
	return nil
}
