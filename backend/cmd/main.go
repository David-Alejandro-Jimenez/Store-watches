package main

import (
	"log"
	"net/http"
	"time"

	"github.com/David-Alejandro-Jimenez/sale-watches/internal"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/config"
	authConfig "github.com/David-Alejandro-Jimenez/sale-watches/internal/config/auth_config"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/repository/database"
	ratelimiter "github.com/David-Alejandro-Jimenez/sale-watches/pkg/security/rate_limiter"
)

func loadConfigurationEnv() {
	var errConfig = config.LoadConfig()
	if errConfig != nil {
		log.Fatalf("Error loading configuration: %v", errConfig)
	}
}

func startDatabase() {
	var errdb = database.InitDB()
	if errdb != nil {
		log.Fatalf("Failed to connect to the database: %v", errdb)
	}
}

func startTheServer(router http.Handler) {
	var port = ":" + config.Config.Server.Port
	log.Printf("Server listening on http://localhost%s", port)
	var err = http.ListenAndServe(port, router)
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}

func configureRateLimiter() (*ratelimiter.DefaultRateLimiterHandler, *ratelimiter.RateLimiterCleaner) {
	manager := ratelimiter.NewRateLimiterManager()
	cleaner := ratelimiter.NewRateLimiterCleaner(manager)

	// Use configuration values for rate limiting
	cleanupTime := time.Duration(config.Config.RateLimit.CleanupInterval) * time.Minute
	expirationTime := time.Duration(config.Config.RateLimit.ExpirationTime) * time.Minute
	cleaner.Start(expirationTime, cleanupTime)

	rateHandler := ratelimiter.NewDefaultRateLimiter(manager)
	return rateHandler, cleaner
}

func main() {
	log.Println("Starting Store Watches API...")

	// Load configuration first
	loadConfigurationEnv()
	log.Println("Configuration loaded successfully")

	// Initialize database
	startDatabase()
	log.Println("Database initialized successfully")

	// Initialize auth handlers
	authConfig.InitializeHandlers()
	log.Println("Auth handlers initialized successfully")

	// Configure rate limiter
	rateHandler, _ := configureRateLimiter()
	log.Println("Rate limiter configured successfully")

	// Setup router
	router := internal.SetupRouter(database.DB, rateHandler)
	log.Println("Router setup successfully")

	// Start the server
	log.Printf("Starting server on port %s...", config.Config.Server.Port)
	startTheServer(router)

	// Cleanup
	defer database.DB.Close()
	log.Println("Server stopped")
}
