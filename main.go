package main

import (
	"simvizlab-backend/config"
	"simvizlab-backend/infra/database"
	"simvizlab-backend/infra/logger"
	"simvizlab-backend/routers"
	"time"

	"github.com/spf13/viper"
)

func main() {
	//set timezone
	viper.SetDefault("SERVER_TIMEZONE", "Asia/Dhaka")
	loc, _ := time.LoadLocation(viper.GetString("SERVER_TIMEZONE"))
	time.Local = loc

	if err := config.SetupConfig(); err != nil {
		logger.Fatalf("config SetupConfig() error: %s", err)
	}

	// Get MongoDB URI from config
	mongoURI := config.MongoDBUri()

	// Initialize MongoDB connection
	if err := database.DbConnection(mongoURI); err != nil {
		logger.Fatalf("MongoDB connection error: %s", err)
	}

	router := routers.SetupRoute()
	logger.Fatalf("%v", router.Run(config.ServerConfig()))

}
