package config

import (
	"log"

	"github.com/David-Alejandro-Jimenez/sale-watches/internal/models"
	ratelimiter "github.com/David-Alejandro-Jimenez/sale-watches/pkg/security/rate_limiter"
	"github.com/spf13/viper"
)

// The ConfigRateLimiter function is responsible for loading and applying the rate limiter configuration from a YAML configuration file. 
// 1. Configuration Loading: The function uses Viper to read a YAML file called config.yaml.
// 2. Specific Section: Extracts the specific configuration of the rate limiter (the ratelimiter section).
// 3. Deserialization and Application: Convert that section to a configuration structure (models.LimiterConfig) and apply it using ratelimiter.SetDefaultLimiterConfig.
// 4. Error Handling: If any error occurs in reading or deserialization, it is logged and the error is returned.
// This feature is essential to customize the behavior of the rate limiter without needing to recompile the application, allowing dynamic adjustments through the configuration file.
func ConfigRateLimiter() error {
	viper.New()
	viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    
    if err := viper.ReadInConfig(); err != nil {
        log.Fatalf("Error reading config file: %v", err)
		return err
    }
    
    var config models.LimiterConfig
    if err := viper.Sub("ratelimiter").Unmarshal(&config); err != nil {
        log.Fatalf("Unable to decode ratelimiter config: %v", err)
		return err
    }
	log.Println("Rate limiter configuration loaded successfully")
    ratelimiter.SetDefaultLimiterConfig(config)
	return nil
}