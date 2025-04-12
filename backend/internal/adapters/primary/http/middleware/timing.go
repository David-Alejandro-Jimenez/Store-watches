// Package middleware provides HTTP middleware utilities as part of the application's
// infrastructure layer, following principles of a hexagonal architecture. This file
// contains a timing middleware that measures and logs the response duration for HTTP requests.
package middleware

import (
	"log"
	"net/http"
	"time"
)

/// TimingConfig contains configuration options for the timing middleware.
// It defines the duration threshold beyond which a performance warning is issued,
// and a function to log the performance metrics.
type TimingConfig struct {
	// WarningThreshold is the duration (time.Duration) beyond which a performance
	// warning is recorded. For example, if set to 500ms, any request taking longer
	// than 500ms will trigger a warning in the logs.
	WarningThreshold time.Duration
	// LogFunc is a function used to log performance information. It receives the HTTP method,
	// request path, the duration it took to process the request, and a boolean flag indicating
	// whether the duration exceeded the WarningThreshold (true if a warning should be logged).
	LogFunc func(method, path string, duration time.Duration, warning bool)
}

// DefaultTimingConfig creates and returns a default configuration for the timing middleware.

// By default, it sets a warning threshold of 500ms and defines a logging function that uses the standard log package to output performance information. If the duration exceeds the threshold, the log entry is prefixed with "[PERF WARNING]"; otherwise, it is prefixed with "[PERF]".
func DefaultTimingConfig() *TimingConfig {
	return &TimingConfig{
		WarningThreshold: 500 * time.Millisecond,  // 500ms threshold
		LogFunc: func(method, path string, duration time.Duration, warning bool) {
			if warning {
				log.Printf("[PERF WARNING] %s %s - %s", method, path, duration)
			} else {
				log.Printf("[PERF] %s %s - %s", method, path, duration)
			}
		},
	}
}

// TimingMiddleware returns a middleware that measures the response time of HTTP requests.

// This middleware records the start time before executing the next handler in the chain, calculates the duration once the response has been served, and invokes the LogFunc defined in the provided TimingConfig. It logs a performance warning if the duration exceeds the configured WarningThreshold.
func TimingMiddleware(config *TimingConfig) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Record the start time before processing the request.
			start := time.Now()

			// Execute the next handler in the chain.
			next.ServeHTTP(w, r)

			// Calculate the elapsed time after processing the request.
			duration := time.Since(start)

			// Determine whether the duration exceeds the warning threshold.
			warning := duration > config.WarningThreshold

			// Log the timing information using the provided logging function.
			config.LogFunc(r.Method, r.URL.Path, duration, warning)
		})
	}
}

// SimpleTimingMiddleware is a convenience function that applies the TimingMiddleware using a default configuration.

// This allows for easy integration of timing metrics without custom configuration making it ideal for rapid development or environments where the default settings are acceptable.
func SimpleTimingMiddleware(next http.Handler) http.Handler {
	return TimingMiddleware(DefaultTimingConfig())(next)
}
