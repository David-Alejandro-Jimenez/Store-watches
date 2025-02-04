package config

import (
	"log"

	"github.com/spf13/viper"
)

func LoadConfig() error {
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