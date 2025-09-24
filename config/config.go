package config

import (
	"errors"
	"os"
	"simvizlab-backend/infra/logger"

	"github.com/spf13/viper"
)

type Configuration struct {
	Server ServerConfiguration
}

// SetupConfig configuration
func SetupConfig() error {
	var configuration Configuration

	pwd, err := os.Getwd()
	if err != nil {
		logger.Errorf("error getting working directory: %s", err)
		return err
	}

	// Enable reading from environment variables
	viper.AutomaticEnv()

	// Try to read from .env file if present, but don't fail if it's missing
	viper.SetConfigFile(pwd + "/.env")
	if err := viper.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if errors.As(err, &notFound) {
			logger.Warnf(".env file not found; continuing with environment variables and defaults")
		} else {
			logger.Errorf("Error reading config file: %s", err)
			return err
		}
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		logger.Errorf("error to decode, %v", err)
		return err
	}

	return nil
}
