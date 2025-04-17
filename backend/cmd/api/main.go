// Package main is the entry point for the watch store API application.
// It initializes all necessary components and starts the HTTP server.
package main

import (
	"database/sql"
	"fmt"
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
	_ "github.com/go-sql-driver/mysql"
)

// main is the application entry point.

// It loads the application configuration, initializes all required components such as database connections, services, adapters, and middleware, and finally starts the HTTP server on the configured port.
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

// initializeCommonServices sets up services that are shared globally across the application.

// Currently, this function initializes the default JWT authentication service using the secret key from configuration.
func initializeCommonServices(appConfig *config.AppConfig) {
	securityAuth.SetDefaultJWTService(appConfig.GetJWTSecret())
}

// setupDatabase establishes a connection to the MySQL database.

// It uses configuration values such as username, password, host, and database name to construct the DSN string and open the connection. It returns a *sql.DB instance and an error if the connection fails.
func setupDatabase(appConfig *config.AppConfig) (*sql.DB, error) {
	cfg := appConfig.GetConfig()
	user := cfg.GetString("database.user")
	password := cfg.GetString("database.password")
	host := cfg.GetString("database.host")
	port := cfg.GetInt("database.port")
	dbName := cfg.GetString("database.name")
	
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", user, password, host, port, dbName)
	return sql.Open("mysql", dsn)
}

// setupUserRepository returns an implementation of the UserRepository interface.

// It sets up dependencies for user authentication such as salt generation and password hashing and injects them into the SQL-based repository.
func setupUserRepository(db *sql.DB) output.UserRepository {
	saltGenerator := securityAuth.RandomSaltGenerator{}
	hasher := securityAuth.BcryptHasher{}
	return repository.NewSQLUserRepository(db, saltGenerator, hasher)
}

// setupLoginService initializes and returns the user login service.

// This service validates credentials and authenticates users.
// It relies on validators for username and password and uses the user repository to query user data.
func setupLoginService(userRepo output.UserRepository) input.UserServiceLogin {
	userNameValidator := &services.UserNameValidator{}
	passwordValidator := &services.PasswordValidator{}
	return services.NewUserLoginService(userRepo, userNameValidator, passwordValidator)
}

// setupRegisterService initializes and returns the user registration service.
// It validates user input and stores new users in the database using the provided repository.
func setupRegisterService(userRepo output.UserRepository) input.UserServiceRegister {
	userNameValidator := &services.UserNameValidator{}
	passwordValidator := &services.PasswordValidator{}
	return services.NewUserRegisterService(userRepo, userNameValidator, passwordValidator)
}

// setupCommentService returns an implementation of the CommentService interface.

// Currently, this function returns nil. The comment service is a planned feature and will be implemented in the future to support user comments.
func setupCommentService(db *sql.DB) input.CommentService {
	// Comment service implementation will be added in the future
	return nil
}

// setupRateLimiter configures and returns a rate limiting handler.

// It uses rate limit settings (requests per second and burst) defined in the application configuration to protect the API against abuse or DoS attacks.
func setupRateLimiter(appConfig *config.AppConfig) ratelimiter.RateLimiterHandler {
	limiterConfig := appConfig.GetRateLimitConfig()
	return ratelimiter.NewDefaultRateLimiter(limiterConfig.RequestPerSecond, limiterConfig.Burst)
}

// setupStaticFileAdapter creates and returns an adapter for serving static files.
//
// The adapter serves assets such as images, stylesheets, or JavaScript files
// from a directory defined in the configuration.
func setupStaticFileAdapter(appConfig *config.AppConfig) output.StaticFilePort {
	staticDir := appConfig.GetStaticDir()
	return static.NewStaticFileAdapter(staticDir)
}
