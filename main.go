// main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// Application holds all application-wide dependencies (the "toolbox").
type Application struct {
	Logger      *log.Logger
	Hub         *Hub
	ChatService *ChatService
	ServerID    string
}

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	serverID := uuid.New().String()
	ctx := context.Background()

	redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		logger.Fatalf("Failed to start redis: %s", err)
	}

	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Port = 9042
	cluster.DisableInitialHostLookup = true
	cassandraSession, err := cluster.CreateSession()
	if err != nil {
		logger.Fatalf("Unable to connect to Cassandra: %v", err)
	}
	defer cassandraSession.Close()
	defer redisClient.Close()

	chatRepository := &CassandraRepository{Session: cassandraSession}
	chatService := &ChatService{Repo: chatRepository}

	hub := newHub(redisClient)
	go hub.run()

	app := &Application{
		Logger:      logger,
		Hub:         hub,
		ChatService: chatService,
		ServerID:    serverID,
	}

	logger.Printf("Starting server %s on :3000", app.ServerID)
	log.Fatal(http.ListenAndServe(":3000", app.routes()))
}
