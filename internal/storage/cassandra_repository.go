package storage

import (
	"context"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/pedroctbd/WhatsappClone/internal/domain"
)

type CassandraRepository struct {
	Session *gocql.Session
}

func NewCassandraRepo(session *gocql.Session) ChatRepository {
	return &CassandraRepository{Session: session}
}

func (c *CassandraRepository) GetOrCreateOneOnOneChat(ctx context.Context, userID1 string, userID2 string) error {
	panic("unimplemented")
}

func (c *CassandraRepository) SaveMessage(ctx context.Context, msg domain.Message) error {
	panic("unimplemented")
}
