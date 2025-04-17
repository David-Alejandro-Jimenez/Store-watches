// Package http implements HTTP handlers and the routing configuration for the sale-watches application.
// This file contains the RouterConfig implementation and the NewRouter factory function, which wires up all the necessary dependencies (handlers, middlewares, rate limiters, etc.) to create a fully configured router.
package http

import (
	"net/http"

	"github.com/David-Alejandro-Jimenez/sale-watches/internal/adapters/primary/http/middleware"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/ports/output"
	ratelimiter "github.com/David-Alejandro-Jimenez/sale-watches/pkg/security/rate_limiter"
	"github.com/gorilla/mux"
)

// RouterConfiguration defines the interface for configuring routes in the application.
type RouterConfiguration interface {
	SetupRoutes(router *mux.Router)
}

// RouterConfig implements the RouterConfiguration interface.
// It holds all the dependencies required to configure the application's routes.
type RouterConfig struct {
	IPExtractor       ratelimiter.IPExtractor
	RateLimiter       ratelimiter.RateLimiterHandler
	LoginHandler      *LoginHandler
	RegisterHandler   *RegisterHandler
	CommentsHandler   *CommentsHandler
	MainPageHandler   *MainPageHandler
	StaticFileHandler *StaticFileHandler
	MiddlewareManager *middleware.MiddlewareManager
}

// SetupRoutes configures all application routes on the given mux.Router.

// It registers routes for static files, public endpoints (main page, register, login) and protected routes (e.g., comments). Additionally, it applies the rate limiting middleware to each route via the MiddlewareManager.
func (c *RouterConfig) SetupRoutes(router *mux.Router) {
	// Configure routes for static files.
	c.StaticFileHandler.RegisterRoutes(router)

	// Create the rate limiting middleware.
	rateLimitMiddleware := middleware.RateLimitMiddleware(c.IPExtractor, c.RateLimiter)

	// Configure public routes.
	router.Handle("/", c.MiddlewareManager.Apply(
		http.HandlerFunc(c.MainPageHandler.Handle),
		rateLimitMiddleware,
	)).Methods("GET")

	router.Handle("/register", c.MiddlewareManager.Apply(
		http.HandlerFunc(c.RegisterHandler.Handle),
		rateLimitMiddleware,
	)).Methods("POST")

	router.Handle("/login", c.MiddlewareManager.Apply(
		http.HandlerFunc(c.LoginHandler.Handle),
		rateLimitMiddleware,
	)).Methods("POST")

	// Configure protected routes (with authentication and rate limiting).
	router.Handle("/comments", c.MiddlewareManager.Apply(
		http.HandlerFunc(c.CommentsHandler.Handle),
		rateLimitMiddleware,
	)).Methods("GET")
}

// NewRouter creates and returns a new mux.Router with all dependencies injected and configured.

// It instantiates the necessary HTTP handlers (LoginHandler, RegisterHandler, CommentsHandler, MainPageHandler, StaticFileHandler), configures middleware (logging, timing, CORS, rate limiting), and sets up all the routes via a RouterConfig. The function returns a ready-to-use router for the application.
func NewRouter(
	userServiceLogin input.UserServiceLogin,
	userServiceRegister input.UserServiceRegister,
	commentService input.CommentService,
	rateHandler ratelimiter.RateLimiterHandler,
	staticFileService output.StaticFilePort,
) *mux.Router {
	router := mux.NewRouter()

	// Create handler instances.
	loginHandler := NewLoginHandler(userServiceLogin)
	registerHandler := NewRegisterHandler(userServiceRegister)
	commentsHandler := NewCommentsHandler(commentService)
	mainPageHandler := NewMainPageHandler()
	staticFileHandler := NewStaticFileHandler(staticFileService)

	// Configure the main page handler with the static directory.
	mainPageHandler.SetStaticDir(staticFileService.GetStaticDir())

	// Create and configure the MiddlewareManager.
	middlewareManager := middleware.NewMiddlewareManager()

	// Create custom timing configuration.
	timingConfig := middleware.DefaultTimingConfig()

	// Set the timing warning threshold to 200ms.
	timingConfig.WarningThreshold = 200 * 1000 * 1000 // 200ms in nanoseconds

	// Create custom CORS configuration.
	corsConfig := middleware.DefaultCORSConfig()
	// Additional CORS settings can be customized here, e.g.:
	// corsConfig.AllowedOrigins = []string{"https://domain.com"}

	// Add global middlewares.
	middlewareManager.AddGlobal(middleware.LoggingMiddleware)
	middlewareManager.AddGlobal(middleware.TimingMiddleware(timingConfig))
	middlewareManager.AddGlobal(middleware.CORSMiddleware(corsConfig))

	// Apply global middlewares to the router.
	middlewareManager.ApplyToRouter(router)

	// Create RouterConfig with all dependencies.
	config := &RouterConfig{
		IPExtractor:       &ratelimiter.DefaultIPExtractor{},
		RateLimiter:       rateHandler,
		LoginHandler:      loginHandler,
		RegisterHandler:   registerHandler,
		CommentsHandler:   commentsHandler,
		MainPageHandler:   mainPageHandler,
		StaticFileHandler: staticFileHandler,
		MiddlewareManager: middlewareManager,
	}

	// Configure routes on the router.
	config.SetupRoutes(router)
	return router
}
