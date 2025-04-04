// Package main is the entry point for the watch store API application.
// It initializes all necessary components and starts the HTTP server.
package main

import (
	"database/sql"
	"log"
	"net/http"

	primaryHttp "github.com/David-Alejandro-Jimenez/sale-watches/internal/adapters/primary/http"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/adapters/secondary/repository"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/adapters/secondary/static"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/config"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/domain/services"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/ports/output"
	ratelimiter "github.com/David-Alejandro-Jimenez/sale-watches/pkg/security/rate_limiter"
	securityAuth "github.com/David-Alejandro-Jimenez/sale-watches/pkg/security/security_auth"
	_ "github.com/mattn/go-sqlite3"
)

// main is the application entry point. It initializes all components,
// configures the dependency injection, and starts the HTTP server.
func main() {
	// Load configuration
	appConfig := config.NewAppConfig()
	appConfig.ValidateConfig()

	// Initialize common services
	initializeCommonServices(appConfig)

	// Initialize database connection
	db, err := setupDatabase(appConfig)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Initialize services
	userRepo := setupUserRepository(db)
	userServiceLogin := setupLoginService(userRepo)
	userServiceRegister := setupRegisterService(userRepo)
	commentService := setupCommentService(db)
	rateHandler := setupRateLimiter(appConfig)

	// Initialize static file adapter
	staticFileAdapter := setupStaticFileAdapter(appConfig)

	// Configure router
	router := primaryHttp.NewRouter(
		userServiceLogin,
		userServiceRegister,
		commentService,
		rateHandler,
		staticFileAdapter,
	)

	// Start server
	port := appConfig.GetPort()
	log.Printf("Server started at http://localhost:%s", port)
	log.Printf("Serving static files from: %s", staticFileAdapter.GetStaticDir())
	log.Fatal(http.ListenAndServe(":"+port, router))
}

// initializeCommonServices sets up services that are shared across the application.
// Currently, it initializes the JWT service with the configured secret key.
func initializeCommonServices(appConfig *config.AppConfig) {
	// Initialize JWT service
	securityAuth.SetDefaultJWTService(appConfig.GetJWTSecret())
}

// setupDatabase establishes a connection to the SQLite database using the path
// specified in the application configuration.
// Returns a pointer to the database connection and an error if the connection fails.
func setupDatabase(appConfig *config.AppConfig) (*sql.DB, error) {
	dbPath := appConfig.GetDBPath()
	return sql.Open("sqlite3", dbPath)
}

// setupUserRepository creates and returns a new user repository implementation
// that interacts with the database. It also initializes the necessary security
// components for user authentication.
func setupUserRepository(db *sql.DB) output.UserRepository {
	saltGenerator := securityAuth.RandomSaltGenerator{}
	hasher := securityAuth.BcryptHasher{}
	return repository.NewSQLUserRepository(db, saltGenerator, hasher)
}

// setupLoginService creates and returns a new user login service.
// It initializes username and password validators and injects the user repository.
func setupLoginService(userRepo output.UserRepository) input.UserServiceLogin {
	userNameValidator := &services.UserNameValidator{}
	passwordValidator := &services.PasswordValidator{}
	return services.NewUserLoginService(userRepo, userNameValidator, passwordValidator)
}

// setupRegisterService creates and returns a new user registration service.
// It uses the same validators as the login service and injects the user repository.
func setupRegisterService(userRepo output.UserRepository) input.UserServiceRegister {
	userNameValidator := &services.UserNameValidator{}
	passwordValidator := &services.PasswordValidator{}
	return services.NewUserRegisterService(userRepo, userNameValidator, passwordValidator)
}

// setupCommentService creates and returns a comment service implementation.
// This function currently returns nil as the comment service is not yet implemented.
// In the future, this will be fully implemented to provide comment functionality.
func setupCommentService(db *sql.DB) input.CommentService {
	// Comment service implementation will be added in the future
	return nil
}

// setupRateLimiter creates and configures a rate limiter handler based on the
// application configuration to protect against DOS attacks.
// It uses requests per second and burst settings from the configuration.
func setupRateLimiter(appConfig *config.AppConfig) ratelimiter.RateLimiterHandler {
	// Get rate limiter configuration
	limiterConfig := appConfig.GetRateLimitConfig()

	// Create rate limiter with the specified configuration
	return ratelimiter.NewDefaultRateLimiter(limiterConfig.RequestPerSecond, limiterConfig.Burst)
}

// setupStaticFileAdapter creates an adapter for serving static files from the
// directory specified in the application configuration.
func setupStaticFileAdapter(appConfig *config.AppConfig) output.StaticFilePort {
	// Get the static directory path from configuration
	staticDir := appConfig.GetStaticDir()

	// Create static file adapter
	return static.NewStaticFileAdapter(staticDir)
}
