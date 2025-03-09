package ratelimiter

import (
	"net"
	"net/http"
)

type IPExtractor interface {
	Extract(remoteAddr string) string
}

type DefaultIPExtractor struct{}



func (e *DefaultIPExtractor) Extract(remoteAddr string) string {
	if ip, _, err := net.SplitHostPort(remoteAddr); err == nil {
		return ip
	}
	return remoteAddr
}

type RateLimiterHandler interface {
	Allow(clientIP string) bool
}

type DefaultRateLimiterHandler struct {}

func (d *DefaultRateLimiterHandler) Allow(clientIP string) bool {
	limiter := GetRateLimiterForIP(clientIP)
	return limiter.Allow()
}

func RateLimitMiddleware(next http.Handler, ipExtractor IPExtractor, rateLimiter RateLimiterHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := ipExtractor.Extract(r.RemoteAddr)

		limiter := GetRateLimiterForIP(clientIP)
		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}