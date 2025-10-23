package storage

import (
	"context"
	"fmt"
	"sort"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/google/uuid"
	"github.com/pedroctbd/WhatsappClone/internal/domain"
)

type CassandraRepository struct {
	Session *gocql.Session
}

func NewCassandraRepo(session *gocql.Session) ChatRepository {
	return &CassandraRepository{Session: session}
}

func sortUserIDs(user1, user2 uuid.UUID) (uuid.UUID, uuid.UUID) {
	ids := []string{user1.String(), user2.String()}
	sort.Strings(ids)
	return uuid.MustParse(ids[0]), uuid.MustParse(ids[1])
}

func (r *CassandraRepository) GetOrCreateOneOnOneChat(ctx context.Context, user1, user2 uuid.UUID) (uuid.UUID, []uuid.UUID, error) {

	userA, userB := sortUserIDs(user1, user2)
	participantsList := []uuid.UUID{user1, user2}

	var chatID uuid.UUID
	err := r.Session.Query(`
		SELECT chat_id 
		FROM one_on_one_lookup 
		WHERE user_a_id = ? AND user_b_id = ?
		LIMIT 1`, userA, userB).Scan(&chatID)

	if err == gocql.ErrNotFound {
		newChatID := uuid.New()

		batch := r.Session.Batch(gocql.LoggedBatch)
		batch.Query(`INSERT INTO one_on_one_lookup (user_a_id, user_b_id, chat_id) VALUES (?, ?, ?)`, userA, userB, newChatID)
		batch.Query(`INSERT INTO chat_metadata (chat_id, chat_type, participants, created_at) VALUES (?, ?, ?, ?)`, newChatID, "one_on_one", participantsList, time.Now())

		now := gocql.UUIDFromTime(time.Now())
		batch.Query(`INSERT INTO user_chats (user_id, last_message_time, chat_id, chat_type) VALUES (?, ?, ?, ?)`, user1, now, newChatID, "one_on_one")
		batch.Query(`INSERT INTO user_chats (user_id, last_message_time, chat_id, chat_type) VALUES (?, ?, ?, ?)`, user2, now, newChatID, "one_on_one")

		if err := r.Session.Batch(gocql.LoggedBatch).ExecContext(ctx); err != nil {
			return uuid.Nil, nil, fmt.Errorf("failed to create chat batch: %w", err)
		}

		return newChatID, participantsList, nil

	} else if err != nil {
		return uuid.Nil, nil, fmt.Errorf("failed to query chat: %w", err)
	}

	// --- FIX: Handle case where chat *was* found ---
	// (err == nil)
	return chatID, participantsList, nil
}

func (r *CassandraRepository) SaveMessage(ctx context.Context, msg domain.Message) error {
	err := r.Session.Query(`
		INSERT INTO messages_by_chat 
		(chat_id, sent_at, message_id, sender_id, sender_name, content) 
		VALUES (?, ?, ?, ?, ?, ?)`,
		msg.ChatID, msg.SentAt, msg.MessageID, msg.SenderID, msg.SenderName, msg.Content).ExecContext(ctx)

	if err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}
	return nil
}
