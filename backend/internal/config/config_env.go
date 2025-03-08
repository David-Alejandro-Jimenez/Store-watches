package config

import (
	"log"

	"github.com/spf13/viper"
)

// The LoadConfig function is responsible for loading the application configuration from an environment file (.env) located in the internal/config directory using Viper.
// 1. Initialize Viper and configure the name and type of the file to load (.env of type env).
// 2. Sets the path where the configuration file is located (internal/config).
// 3. Reads the configuration file, handling possible errors.
// 4. Enables automatic reading of environment variables to allow configurations to be overwritten.
// 5. Confirms the successful loading with a message in the log and returns nil.
// This feature is essential for loading and centralizing application configuration, allowing it to be flexibly tuned using an environment file and environment variables.
func LoadConfig() error {
	viper.New()
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath("internal/config")

	var err = viper.ReadInConfig()
	if err != nil {
		log.Printf("Error loading configuration: %v", err)
		return err
	}

	viper.AutomaticEnv()
	log.Println(".env file loaded successfully")
	return nil
}