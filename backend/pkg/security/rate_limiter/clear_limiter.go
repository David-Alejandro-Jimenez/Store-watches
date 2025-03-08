package ratelimiter

import "time"

// The StartCleanupRoutine function starts a background routine that periodically cleans inactive limiters. This helps free up memory and keep the system efficient.
// This technique is common in applications that use rate limiting mechanisms to avoid unnecessary memory consumption and maintain system efficiency.
func StartCleanupRoutine(expirationDuration, cleanupInterval time.Duration) {
	go func() {
		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()
		for range ticker.C {
			CleanupInactiveLimiters(expirationDuration)
		}
	}()
}