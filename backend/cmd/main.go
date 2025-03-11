package main

import (
	"log"
	"net/http"
	"time"

	"github.com/David-Alejandro-Jimenez/sale-watches/internal"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/config"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/repository/database"
	ratelimiter "github.com/David-Alejandro-Jimenez/sale-watches/pkg/security/rate_limiter"
)

// The loadConfigurationEnv function is responsible for initializing the application configuration by loading the environment file.
// 1. The function calls LoadConfig() to load the application configuration from an environment file.
// 2. If an error occurs, it is logged and execution terminates with a fatal message.
// 3. This ensures that the application does not run without the necessary configuration.
// This feature is essential to ensure that the application starts with the correct configuration parameters defined in the environment file.
func loadConfigurationEnv() {
	var errConfig = config.LoadConfig()
	if errConfig != nil {
		log.Fatalf("Error loading configuration: %v", errConfig)
	}
}

// The startDatabase function is responsible for starting the connection to the database by calling database.InitDB().
// 1. Purpose: Start the connection to the database.
// 2. Actions:
		// Call database.InitDB() to set up the connection.
		// Checks if any errors occurred and, if so, records them in the log.
// This feature is essential to ensure that the application has an active connection to the database before executing operations that depend on it.
func startDatabase() {
	var errdb = database.InitDB()
	if errdb != nil {
		log.Println("Did not connect to the database")
	}
}

// The startTheServer function is responsible for starting the application's HTTP server. 
// 1. Listening Port: The server starts on port 8080.
// 2. Router Configuration: The router is established with the routes defined in the SetupRouter function.
// 3. Start and Listen: Start the server using http.ListenAndServe.
// 4. Error Handling: If any error occurs, execution is stopped and the error is logged.
// This feature is essential to get the HTTP server up and running and begin serving application requests.
func startTheServer(router http.Handler) {
	var port = ":8080"
	log.Printf("Server listening on http://localhost%s", port)
	var err = http.ListenAndServe(port, router)
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}

// The main function is the entry point of the application and orchestrates the initialization and execution of all key components.
// 1. The function closes the connection to the database upon application termination.
// 2. It loads the necessary configuration from environment files and YAML.
// 3. It initializes the connection to the database.
// 4. It configures and launches a rate limiter cleanup routine.
// 5. Finally, it starts the HTTP server to start receiving and processing requests.
// This flow ensures that the application is properly configured, connected to the database, and ready to handle HTTP requests safely and efficiently.
func main() {
	loadConfigurationEnv()
	startDatabase()
	//loadConfigurationYaml()
	manager := ratelimiter.NewRateLimiterManager()
    cleaner := ratelimiter.NewRateLimiterCleaner(manager)
    cleaner.Start(30*time.Minute, 10*time.Minute)

    // Creamos el handler con el manager
    rateHandler := ratelimiter.NewDefaultRateLimiter(manager)

	router := internal.SetupRouter(rateHandler)
	startTheServer(router)
	defer database.DB.Close()
}