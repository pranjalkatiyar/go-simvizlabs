package config

import (
	"os"
)

func MongoDBUri() string {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		// Fallback to localhost if no URI is provided
		return "mongodb://localhost:27017/simvizlab"
	}
	return uri
}
