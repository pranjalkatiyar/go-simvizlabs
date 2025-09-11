package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	err         error
	MongoClient *mongo.Client
)

// DbConnection creates MongoDB connection
func DbConnection(uri string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)

	// Set client options based on configuration
	// If you want to add a CommandMonitor, implement it here.
	// For now, this block is removed to avoid errors.
	if viper.GetBool("DB_LOG_MODE") {
		// Add command monitor implementation here if needed

	}

	MongoClient, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
		return err
	}

	// Ping the database to verify connection
	err = MongoClient.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("MongoDB ping error: %v", err)
		return err
	}

	fmt.Println("Connected to MongoDB!")
	return nil
}

// GetMongoClient returns the MongoDB client
func GetMongoClient() *mongo.Client {
	return MongoClient
}
