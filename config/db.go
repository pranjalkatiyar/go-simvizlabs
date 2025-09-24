package config

import (
	"log"
	"os"
)

func MongoDBUri() string {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Println("MONGODB_URI not set, using default localhost URI")
		return "mongodb://localhost:27017/simvizlab"
	}
	log.Println("Using MongoDB URI from environment")
	return uri
}
