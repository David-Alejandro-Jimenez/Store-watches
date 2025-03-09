package ratelimiter

import "time"

type CleanupHandler interface {
	Cleanup(expirationDuration time.Duration)
}

type RatelimiterCleanup struct {}

func (c *RatelimiterCleanup) Cleanup(expirationDuration time.Duration) {
	CleanupInactiveLimiters(expirationDuration)
}

func StartCleanupRoutine(handler CleanupHandler, expirationDuration, cleanupInterval time.Duration) {
	go func() {
		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()
		for range ticker.C {
			CleanupInactiveLimiters(expirationDuration)
		}
	}()
}