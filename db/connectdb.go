package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectMongo connects to MongoDB and returns a Database instance.
func ConnectMongo(uri, dbname string) (*mongo.Client, *mongo.Database, error) {
	// create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// connect to client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect mongo: %w", err)
	}

	// ping to check connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, nil, fmt.Errorf("mongo ping failed: %w", err)
	}

	fmt.Println("✅ Mongo connection successful")
	return client, client.Database(dbname), nil
}
