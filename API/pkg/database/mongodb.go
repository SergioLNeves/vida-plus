// Package database provides database connection and management functions
package database

import (
	"context"
	"log/slog"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InitMongoDB initializes and returns a MongoDB client
func InitMongoDB() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/vida_plus"))
	if err != nil {
		slog.Error("error connecting to MongoDB", slog.Any("error", err))
		os.Exit(1)
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		slog.Error("error pinging MongoDB", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("successfully connected to MongoDB")
	return client
}

// DisconnectMongoDB gracefully disconnects from MongoDB
func DisconnectMongoDB(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		slog.Error("error disconnecting from MongoDB", slog.Any("error", err))
		return
	}

	slog.Info("successfully disconnected from MongoDB")
}

// GetDatabase returns a specific database from the MongoDB client
func GetDatabase(client *mongo.Client, dbName string) *mongo.Database {
	return client.Database(dbName)
}
