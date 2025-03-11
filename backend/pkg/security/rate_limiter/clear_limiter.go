package ratelimiter

import "time"

type RateLimiterCleaner struct {
	manager RateLimiterManager
}

func NewRateLimiterCleaner(manager RateLimiterManager) *RateLimiterCleaner {
	return &RateLimiterCleaner{manager: manager}
}

func (c *RateLimiterCleaner) Start(expirationDuration, cleanupInterval time.Duration) {
	go func() {
		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()
		for range ticker.C {
			c.manager.CleanupInactiveLimiters(expirationDuration)
		}
	}()
}
