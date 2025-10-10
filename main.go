package main

import (
	"context"
	"log"
	"net/http"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var serverID = uuid.New().String()

func main() {

	ctx := context.Background()

	//Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to start redis : %s", err)
	}

	//cassandra
	// Create a cluster configuration
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Port = 9042
	cluster.DisableInitialHostLookup = true
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("unable to connect to Cassandra: %v", err)
	}

	defer session.Close()
	defer rdb.Close()

	hub := newHub(rdb)
	go hub.run()
	app := &Application{
		RD:       rdb,
		HUB:      hub,
		CS:       session,
		ServerID: serverID,
	}

	print("Starting server")
	log.Fatal(http.ListenAndServe(":3000", app.routes()))

}
