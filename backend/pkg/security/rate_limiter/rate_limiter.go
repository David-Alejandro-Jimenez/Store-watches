package ratelimiter

import (
	"sync"
	"time"

	"github.com/David-Alejandro-Jimenez/sale-watches/internal/models"
	"golang.org/x/time/rate"
)

type RateLimiterManager interface {
	GetRateLimiterForIP(ipAddress string) *rate.Limiter
	CleanupInactiveLimiters(expirationDuration time.Duration)
	SetDefaultLimiterConfig(config models.LimiterConfig)
}

type DefaultRateLimiterManager struct {
	defaultConfig  models.LimiterConfig
	ipLimiterCache sync.Map	
}

func NewRateLimiterManager() *DefaultRateLimiterManager {
	return &DefaultRateLimiterManager{
		defaultConfig: models.LimiterConfig{
			RequestPerSecond: 1.5,
			Burst:            5,
		},
	}
}

func (m *DefaultRateLimiterManager) SetDefaultLimiterConfig(config models.LimiterConfig) {
	m.defaultConfig = config
}

func (m *DefaultRateLimiterManager) GetRateLimiterForIP(ipAddress string) *rate.Limiter {
	currentTime := time.Now()
	record, exists := m.ipLimiterCache.Load(ipAddress)
	if exists {
		limiterRecord := record.(*LimiterEntry)
		limiterRecord.lastSeen = currentTime
		return limiterRecord.limiter
	}

	newLimiter := rate.NewLimiter(rate.Limit(m.defaultConfig.RequestPerSecond), m.defaultConfig.Burst)
	newRecord := &LimiterEntry{
		limiter:  newLimiter,
		lastSeen: currentTime,
	}

	actualRecord, loaded := m.ipLimiterCache.LoadOrStore(ipAddress, newRecord)
	if loaded {
		return actualRecord.(*LimiterEntry).limiter
	}
	return newLimiter
}

func (m *DefaultRateLimiterManager) CleanupInactiveLimiters(expirationDuration time.Duration) {
	currentTime := time.Now()
	m.ipLimiterCache.Range(func(key, value interface{}) bool {
		limiterRecord := value.(*LimiterEntry)
		if currentTime.Sub(limiterRecord.lastSeen) > expirationDuration {
			m.ipLimiterCache.Delete(key)
		}
		return true	
	})
}

type LimiterEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}
	