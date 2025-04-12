// Package middleware contains HTTP middleware utilities designed following a hexagonal
// architecture. It provides functions for chaining middlewares, managing global middlewares,
// and built-in examples such as a logging middleware for HTTP requests.
package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Middleware defines a function that wraps an http.Handler to perform pre-processing or post-processing actions on HTTP requests and responses.
type Middleware func(http.Handler) http.Handler

// Chain returns a Middleware that combines multiple middlewares into one.
// Chain accepts a variadic list of Middlewares and applies them in reverse order.
// This ensures that the middleware listed first is the outermost layer, and the last
// middleware is the innermost layer wrapping the actual handler.
func Chain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

// MiddlewareManager manages global and route-specific middlewares within the application.
// This manager is part of the hexagonal architecture, serving as an adapter layer that decouples middleware management from business logic.
type MiddlewareManager struct {
	globalMiddlewares []Middleware
}

// NewMiddlewareManager creates a new MiddlewareManager with an empty slice of global middlewares.

// Use this function to initialize the middleware manager before adding any middleware.
func NewMiddlewareManager() *MiddlewareManager {
	return &MiddlewareManager{
		globalMiddlewares: []Middleware{},
	}
}

// AddGlobal registers a middleware to be applied globally to all routes. 
// The provided middleware is appended to the manager's list of global middlewares.
func (m *MiddlewareManager) AddGlobal(middleware Middleware) {
	m.globalMiddlewares = append(m.globalMiddlewares, middleware)
}

// Apply wraps the given http.Handler with both the global middlewares and any additional middlewares provided as arguments.
// This method merges the global middlewares with route-specific ones using the Chain function.
func (m *MiddlewareManager) Apply(handler http.Handler, middlewares ...Middleware) http.Handler {
	// Combine global middlewares with the provided specific middlewares.
	allMiddlewares := append(m.globalMiddlewares, middlewares...)
	return Chain(allMiddlewares...)(handler)
}

// ApplyToRouter applies all global middlewares to every route of the given mux.Router.
// This function registers a middleware on the router that wraps each route handler with the global middlewares.
func (m *MiddlewareManager) ApplyToRouter(router *mux.Router) {
	if len(m.globalMiddlewares) > 0 {
		router.Use(func(next http.Handler) http.Handler {
			return Chain(m.globalMiddlewares...)(next)
		})
	}
}

// LoggingMiddleware is an HTTP middleware that logs details about each request.
// It records the HTTP method, request URL, response status code, and the time duration taken to process the request. If a response contains an error (status code >= 400), it logs the event as an error; otherwise, it logs it as an informational message.

// In production, consider using a structured logging library instead of the standard log package.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a custom ResponseWriter to capture the HTTP status code.
		rw := NewResponseWriter(w)

		// Invoke the next handler in the chain.
		next.ServeHTTP(rw, r)

		// Calculate the duration of the request.
		duration := time.Since(start)

		// Log the request information depending on the status code.
		if rw.statusCode >= 400 {
			// Log error if status code indicates failure.
			log.Printf(
				"[ERROR] %s %s %d %s",
				r.Method,
				r.URL.Path,
				rw.statusCode,
				duration,
			)
		} else {
			// Log as informational.
			log.Printf(
				"[INFO] %s %s %d %s",
				r.Method,
				r.URL.Path,
				rw.statusCode,
				duration,
			)
		}
	})
}

// ResponseWriter is a custom implementation of http.ResponseWriter that captures the HTTP response status code for logging purposes.

// It embeds the original http.ResponseWriter and overrides the WriteHeader method.
type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// NewResponseWriter creates a new ResponseWriter instance wrapping the provided http.ResponseWriter.

// The returned ResponseWriter is initialized with a default status code of http.StatusOK (200).
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, http.StatusOK}
}

// WriteHeader intercepts calls to write the HTTP status code, storing the provided code and then delegating the call to the embedded http.ResponseWriter.
func (rw *ResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
