package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
)

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

	// //cassandra
	// cluster := gocql.NewCluster("127.0.0.1")
	// cluster.Keyspace = "whatsappclone"
	// cluster.ProtoVersion = 4
	// session, err := cluster.CreateSession()
	// sslOpts := &gocql.SslOptions{
	// 	// For local dev, you can disable host verification.
	// 	// For production, you'd configure your certificates here.
	// 	EnableHostVerification: false,
	// }
	// cluster.SslOpts = sslOpts

	if err != nil {
		log.Fatalf("Failed to create cassandra sessions : %s", err)
	}
	fmt.Println("cassandra init done")

	// defer session.Close()
	defer rdb.Close()

	hub := newHub()
	go hub.run()
	app := &Application{
		RD:  rdb,
		HUB: hub,
		// CS:  session,
	}

	print("Starting server")
	log.Fatal(http.ListenAndServe(":3000", app.routes()))

}
