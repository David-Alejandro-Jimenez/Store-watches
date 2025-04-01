package ratelimiter

import (
	"net"
	"net/http"

	"github.com/David-Alejandro-Jimenez/sale-watches/pkg/errors"
	"golang.org/x/time/rate"
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
	GetRateLimiterForIP(ipAddress string) *rate.Limiter
}

type DefaultRateLimiterHandler struct {
	rateLimiterManager RateLimiterManager
}

func NewDefaultRateLimiterHandler() *DefaultRateLimiterHandler {
	return &DefaultRateLimiterHandler{
		rateLimiterManager: NewRateLimiterManager(),
	}
}

func NewDefaultRateLimiter(manager RateLimiterManager) *DefaultRateLimiterHandler {
    return &DefaultRateLimiterHandler{
        rateLimiterManager: manager,
    }
}

func (d *DefaultRateLimiterHandler) Allow(clientIP string) bool {
	limiter := d.rateLimiterManager.GetRateLimiterForIP(clientIP)
	return limiter.Allow()
}

func (d *DefaultRateLimiterHandler) GetRateLimiterForIP(ipAddress string) *rate.Limiter {
	return d.rateLimiterManager.GetRateLimiterForIP(ipAddress)
}


func RateLimitMiddleware(next http.Handler, ipExtractor IPExtractor, rateLimiter RateLimiterHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := ipExtractor.Extract(r.RemoteAddr)

		limiter := rateLimiter.GetRateLimiterForIP(clientIP)
		if !limiter.Allow() {
			errors.NewTooManyRequestsError("Too many Requests")
			return
		}
		next.ServeHTTP(w, r)
	})
}