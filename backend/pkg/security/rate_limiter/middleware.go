package ratelimiter

import (
	"net"
	"net/http"
)

//The extractClientIP function extracts the IP address from a string that usually includes both the IP and the port (for example, "192.168.1.100:8080").
// 1. The function attempts to separate the IP and port from the chain.
// 2. If the format is correct, it returns only the IP.
// 3. Otherwise, it returns the original string without modifications.
// This implementation is useful to securely obtain the client's IP, considering that the address could come in different formats.
func extractClientIP(remoteAddr string) string {
	if ip, _, err := net.SplitHostPort(remoteAddr); err == nil {
		return ip
	}
	return remoteAddr
}

// The RateLimitMiddleware function is a middleware that implements rate limiting of requests based on the client's IP address. Its objective is to prevent the same IP from making too many requests in a short period, protecting the server from possible abuse.
// 1. Extracts the client's IP: Ensures identifying the source of the request.
// 2. Gets the rate limiter associated with that IP: Allows individual control per IP address.
// 3. Checks if the request exceeds the allowed limit: If it exceeds, a 429 error is returned.
// 4. Request processing continues: If the request meets the constraints, the next handler is called.
// This middleware is essential to protect the server against abuse and prevent the same IP from generating an excessive number of requests in a short time.
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := extractClientIP(r.RemoteAddr)

		limiter := GetRateLimiterForIP(clientIP)
		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}