package config

import (
	"github.com/spf13/viper"
)

func MongoDBUri() string {
	uri := viper.GetString("MONGODB_URI")
	if uri == "" {
		// Fallback to localhost if no URI is provided
		return "mongodb://localhost:27017/simvizlab"
	}
	return uri
}
