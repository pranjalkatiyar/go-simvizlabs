package config

import (
	"fmt"
	"os"
)

type ServerConfiguration struct {
	Port string
}

// SetupConfig validates required environment variables
func SetupConfig() error {
	requiredVars := []string{
		"SERVER_PORT",
		"MONGODB_URI",
		"APPSTORE_KEY_ID",
		"APPSTORE_ISSUER_ID",
		"APPSTORE_PRIVATE_KEY",
		"APPSTORE_BUNDLE_ID",
		"BASE_URL",
	}

	for _, key := range requiredVars {
		if os.Getenv(key) == "" {
			return fmt.Errorf("missing required environment variable: %s", key)
		}
	}

	return nil
}
