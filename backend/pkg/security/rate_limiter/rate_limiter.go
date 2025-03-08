package ratelimiter

import (
	"sync"
	"time"

	"github.com/David-Alejandro-Jimenez/sale-watches/internal/models"
	"golang.org/x/time/rate"
)


type LimiterEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var defaultLimiterConfig = models.LimiterConfig{
	RequestPerSecond: 1.5,
	Burst:       5,
}

var ipLimiterCache sync.Map

// The SetDefaultLimiterConfig function sets a default configuration for a rate limiting system.
// This feature is useful in systems that prevent abuse or DDoS attacks by limiting requests per user or IP.
func SetDefaultLimiterConfig(config models.LimiterConfig) {
	defaultLimiterConfig = config
}

// The GetRateLimiterForIP function manages a Rate Limiter per IP address. Its purpose is to control the number of requests that an IP can make in a given time, avoiding abuse or DDoS attacks.
// 1. Check if the IP already has a rate limiter in ipLimiterCache.
// 2. If it exists, it uses it and updates the last time the IP was seen.
// 3. If it does not exist, it creates a new limiter and saves it to ipLimiterCache.
// 4. Use LoadOrStore to prevent two processes from creating a limiter for the same IP at the same time.
// 5. Returns the corresponding limiter to manage the number of requests allowed.
// This system is useful to protect the application from abuse by limiting the frequency of requests for each IP without blocking legitimate users.
func GetRateLimiterForIP(ipAddress string) *rate.Limiter {
	currentTime := time.Now()
	record, exists := ipLimiterCache.Load(ipAddress) 

    if exists {
        limiterRecord := record.(*LimiterEntry)
        limiterRecord.lastSeen = currentTime
        return limiterRecord.limiter
    }

    newLimiter := rate.NewLimiter(rate.Limit(defaultLimiterConfig.RequestPerSecond), defaultLimiterConfig.Burst)
    newRecord := &LimiterEntry{
        limiter:  newLimiter,
        lastSeen: currentTime,
    }

    actualRecord, loaded := ipLimiterCache.LoadOrStore(ipAddress, newRecord)
    if loaded {
        return actualRecord.(*LimiterEntry).limiter
    }
    return newLimiter
}

// The CleanupInactiveLimiters function is responsible for removing inactive rate limiters from the cache. This is useful for freeing resources and keeping memory clean, eliminating those registers that have not been used in a defined period of time.
// 1. The function gets the current time.
// 2. Loops through all cached limiters.
// 3. Evaluates the inactivity of each limiter by comparing the difference between the current time and the last time it was used.
// 4. If a limiter has been idle for longer than the allowed time (expirationDuration), it is removed from the cache.
// 5. This helps free up resources and keep the system efficient by not retaining limiters that are no longer needed.
func CleanupInactiveLimiters(expirationDuration time.Duration) {
	currentTime := time.Now()
	ipLimiterCache.Range(func(key, value interface{}) bool {
		limiterRecord := value.(*LimiterEntry)
		if currentTime.Sub(limiterRecord.lastSeen) > expirationDuration {
			ipLimiterCache.Delete(key)
		}
		return true
	})
}
