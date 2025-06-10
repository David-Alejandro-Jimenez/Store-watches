// Package http implements HTTP handlers and the routing configuration for the sale-watches application.
// It provides the RouterConfig type and NewRouter factory function, which wire up handlers, middlewares, rate limiters, and static file serving to produce a fully configured *mux.Router* ready to handle API requests.
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

// RouterConfig holds dependencies required to set up application routes.
// Fields correspond to handlers for each endpoint, middleware manager, and rate limiter components.
//
// Fields:
//   - IPExtractor: extracts client IP from *http.Request* for rate limiting.
//   - RateLimiter: handles request rate limiting based on extracted IP.
//   - LoginHandler: processes user login requests.
//   - RegisterHandler: processes user registration requests.
//   - CommentsGetHandler: handles retrieval of comments.
//   - CommentsAddHandler: handles creation of new comments.
//   - MainPageHandler: serves the application's main page.
//   - StaticFileHandler: serves static assets like CSS/JS/images.
//   - MiddlewareManager: orchestrates application of global and route-specific middleware.
type RouterConfig struct {
	IPExtractor        ratelimiter.IPExtractor
	RateLimiter        ratelimiter.RateLimiterHandler
	LoginHandler       *LoginHandler
	RegisterHandler    *RegisterHandler
	CommentsGetHandler *CommentsGetHandler
	CommentsAddHandler *CommentsAddHandler
	MainPageHandler    *MainPageHandler
	StaticFileHandler  *StaticFileHandler
	MiddlewareManager  *middleware.MiddlewareManager
}

// SetupRoutes registers all application endpoints on the given router and applies route-specific middleware for authentication and rate limiting.
// Routes include:
//   - Static files (CSS, JS, images)
//   - Public endpoints: GET /, POST /register, POST /login
//   - Protected endpoints: GET /comments, POST /comments/newComments

// Each route is wrapped with authentication and rate limiting via the MiddlewareManager.Apply method.

// Parameters:
//   - router: *mux.Router instance to configure routes on.
func (c *RouterConfig) SetupRoutes(router *mux.Router) {
	// 1. Register static file serving routes
	c.StaticFileHandler.RegisterRoutes(router)

	// 2. Prepare middleware for rate limiting and authentication
	rateLimitMW := middleware.RateLimitMiddleware(c.IPExtractor, c.RateLimiter)
	authMW := middleware.AuthMiddleware(middleware.DefaultAuthOptions())

	// 3. Public routes
	router.Handle("/", c.MiddlewareManager.Apply(
		http.HandlerFunc(c.MainPageHandler.Handle),
		authMW, rateLimitMW,
	)).Methods("GET")

	router.Handle("/register", c.MiddlewareManager.Apply(
		http.HandlerFunc(c.RegisterHandler.Handle),
		authMW, rateLimitMW,
	)).Methods("POST")

	router.Handle("/login", c.MiddlewareManager.Apply(
		http.HandlerFunc(c.LoginHandler.Handle),
		authMW, rateLimitMW,
	)).Methods("POST")

	router.Handle("/comments", c.MiddlewareManager.Apply(
		http.HandlerFunc(c.CommentsGetHandler.Handle),
		authMW, rateLimitMW,
	)).Methods("GET")

	// 4. Protected routes
	router.Handle("/comments/newComments", c.MiddlewareManager.Apply(
		http.HandlerFunc(c.CommentsAddHandler.Handle),
		authMW, rateLimitMW,
	)).Methods("POST")
}

// NewRouter constructs and returns a *mux.Router configured with all application routes, handlers, and global middleware.
// It performs the following steps:
//  1. Instantiate handler objects for login, registration, comments, etc.
//  2. Set up the MainPageHandler with the static directory path.
//  3. Create and configure a MiddlewareManager, adding global middleware
//     for logging, timing, and CORS.
//  4. Build a RouterConfig with dependencies and call SetupRoutes.

// Parameters:
//   - userServiceLogin: service for authenticating users on login.
//   - userServiceRegister: service for registering new users.
//   - commentGetService: service for fetching existing comments.
//   - commentAddService: service for adding new comments.
//   - rateHandler: rate limiting handler middleware for DoS protection.
//   - staticFileService: adapter for serving static files from disk.

// Returns:
//   - *mux.Router: fully configured router ready to be passed to http.ListenAndServe.
func NewRouter(
	userServiceLogin input.UserServiceLogin,
	userServiceRegister input.UserServiceRegister,
	commentGetService input.CommentGetService,
	commentAddService input.CommentAddService,
	rateHandler ratelimiter.RateLimiterHandler,
	staticFileService output.StaticFilePort,
) *mux.Router {
	// 1. Initialize a new router
	router := mux.NewRouter()

	// 2. Instantiate HTTP handlers with injected domain services
	loginHandler := NewLoginHandler(userServiceLogin)
	registerHandler := NewRegisterHandler(userServiceRegister)
	commentsGetHandler := NewCommentsGetHandler(commentGetService)
	commentsAddHandler := NewCommentAddsHandler(commentAddService)
	mainPageHandler := NewMainPageHandler()
	staticFileHandler := NewStaticFileHandler(staticFileService)

	// 3. Configure main page handler with static directory
	mainPageHandler.SetStaticDir(staticFileService.GetStaticDir())

	// 4. Create and configure MiddlewareManager
	middlewareManager := middleware.NewMiddlewareManager()
	timingConfig := middleware.DefaultTimingConfig()
	timingConfig.WarningThreshold = 200 * 1000 * 1000 // 200 milliseconds

	corsConfig := middleware.DefaultCORSConfig()
	// corsCfg.AllowedOrigins = []string{"https://example.com"} // customize as needed

	// Add global middleware: logging, timing, CORS
	middlewareManager.AddGlobal(middleware.LoggingMiddleware)
	middlewareManager.AddGlobal(middleware.TimingMiddleware(timingConfig))
	middlewareManager.AddGlobal(middleware.CORSMiddleware(corsConfig))
	middlewareManager.ApplyToRouter(router)

	// 5. Build RouterConfig with dependencies
	config := &RouterConfig{
		IPExtractor:        &ratelimiter.DefaultIPExtractor{},
		RateLimiter:        rateHandler,
		LoginHandler:       loginHandler,
		RegisterHandler:    registerHandler,
		CommentsGetHandler: commentsGetHandler,
		CommentsAddHandler: commentsAddHandler,
		MainPageHandler:    mainPageHandler,
		StaticFileHandler:  staticFileHandler,
		MiddlewareManager:  middlewareManager,
	}

	// 6. Register routes on router
	config.SetupRoutes(router)
	return router
}
