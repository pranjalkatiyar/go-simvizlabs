package main

import (
	"log"
	"time"

	"simvizlab-backend/config"
	"simvizlab-backend/infra/database"
	"simvizlab-backend/infra/logger"
	"simvizlab-backend/routers"

	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

func main() {
	// Load .env file if present (for local dev)
	if err := gotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	// Set default timezone
	viper.SetDefault("SERVER_TIMEZONE", "Asia/Dhaka")
	loc, err := time.LoadLocation(viper.GetString("SERVER_TIMEZONE"))
	if err != nil {
		log.Printf("Failed to load timezone: %v", err)
	} else {
		time.Local = loc
	}

	log.Println("Initializing config...")
	if err := config.SetupConfig(); err != nil {
		logger.Fatalf("Config setup failed: %s", err)
	}

	mongoURI := config.MongoDBUri()
	if mongoURI == "" {
		logger.Fatalf("MongoDB URI not set in config")
	}

	log.Println("Connecting to MongoDB...")
	if err := database.DbConnection(mongoURI); err != nil {
		logger.Fatalf("MongoDB connection error: %s", err)
	}

	log.Println("Setting up router...")
	router := routers.SetupRoute()

	log.Println("Starting server...")
	if err := router.Run(config.ServerConfig()); err != nil {
		logger.Fatalf("Server failed to start: %s", err)
	}
}
