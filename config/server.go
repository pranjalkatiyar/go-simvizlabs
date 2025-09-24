package config

import (
	"fmt"
	"log"
	"os"
)

func ServerConfig() string {
	host := os.Getenv("SERVER_HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8000"
	}

	appServer := fmt.Sprintf("%s:%s", host, port)
	log.Printf("Server running at: %s", appServer)
	return appServer
}
