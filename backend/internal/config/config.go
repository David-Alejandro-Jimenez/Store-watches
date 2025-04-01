package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Database  DatabaseConfig
	Server    ServerConfig
	Security  SecurityConfig
	RateLimit RateLimitConfig
}

type DatabaseConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
	DSN      string
}

type ServerConfig struct {
	Port           string
	AllowedOrigins []string
}

type SecurityConfig struct {
	JWTSecret            string
	JWTExpirationMinutes int
	PasswordMinLength    int
}

type RateLimitConfig struct {
	RequestsPerMinute int
	CleanupInterval   int
	ExpirationTime    int
}

var Config AppConfig

func LoadConfig() error {
	viper.Reset()
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

	Config.Database.User = viper.GetString("DB_USER")
	Config.Database.Password = viper.GetString("DB_PASSWORD")
	Config.Database.Host = viper.GetString("DB_HOST")
	Config.Database.Port = viper.GetString("DB_PORT")
	Config.Database.Name = viper.GetString("DB_NAME")
	Config.Database.DSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4",
		Config.Database.User,
		Config.Database.Password,
		Config.Database.Host,
		Config.Database.Port,
		Config.Database.Name)

	Config.Server.Port = viper.GetString("SERVER_PORT")
	if Config.Server.Port == "" {
		Config.Server.Port = "8080"
	}

	Config.Security.JWTSecret = viper.GetString("JWT_SECRET")
	Config.Security.JWTExpirationMinutes = viper.GetInt("JWT_EXPIRATION_MINUTES")
	if Config.Security.JWTExpirationMinutes == 0 {
		Config.Security.JWTExpirationMinutes = 60
	}
	Config.Security.PasswordMinLength = viper.GetInt("PASSWORD_MIN_LENGTH")
	if Config.Security.PasswordMinLength == 0 {
		Config.Security.PasswordMinLength = 8
	}

	Config.RateLimit.RequestsPerMinute = viper.GetInt("RATE_LIMIT_REQUESTS")
	if Config.RateLimit.RequestsPerMinute == 0 {
		Config.RateLimit.RequestsPerMinute = 100
	}
	Config.RateLimit.CleanupInterval = viper.GetInt("RATE_LIMIT_CLEANUP_MINUTES")
	if Config.RateLimit.CleanupInterval == 0 {
		Config.RateLimit.CleanupInterval = 10
	}
	Config.RateLimit.ExpirationTime = viper.GetInt("RATE_LIMIT_EXPIRATION_MINUTES")
	if Config.RateLimit.ExpirationTime == 0 {
		Config.RateLimit.ExpirationTime = 30
	}

	return validateConfig()
}

func validateConfig() error {
	if Config.Database.User == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if Config.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if Config.Database.Name == "" {
		return fmt.Errorf("DB_NAME is required")
	}

	if Config.Security.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required for security")
	}

	return nil
}
