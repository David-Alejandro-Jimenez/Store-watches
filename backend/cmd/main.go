package main

import (
	"log"
	"net/http"
	"time"

	"github.com/David-Alejandro-Jimenez/sale-watches/internal"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/config"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/config/auth_config"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/repository/database"
	"github.com/David-Alejandro-Jimenez/sale-watches/pkg/security/rate_limiter"
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
		log.Println("Did not connect to the database")
	}
}

func startTheServer(router http.Handler) {
	var port = ":8080"
	log.Printf("Server listening on http://localhost%s", port)
	var err = http.ListenAndServe(port, router)
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}

func main() {
	loadConfigurationEnv()
	startDatabase()
	authConfig.InitializeHandlers()

	manager := ratelimiter.NewRateLimiterManager()
    cleaner := ratelimiter.NewRateLimiterCleaner(manager)
    cleaner.Start(30*time.Minute, 10*time.Minute)
    rateHandler := ratelimiter.NewDefaultRateLimiter(manager)
	router := internal.SetupRouter(database.DB, rateHandler)
	
	startTheServer(router)
	defer database.DB.Close()
}