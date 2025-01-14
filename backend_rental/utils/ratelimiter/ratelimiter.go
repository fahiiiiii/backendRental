// utils/ratelimiter/ratelimiter.go
package ratelimiter

import (
    "context"
    "sync"
    "time"
    "golang.org/x/time/rate"
)

type APIRateLimiter struct {
    limiter *rate.Limiter
    mutex   sync.RWMutex
}

var (
    instance *APIRateLimiter
    once     sync.Once
)

// GetInstance returns the singleton instance of APIRateLimiter
func GetInstance() *APIRateLimiter {
    once.Do(func() {
        instance = &APIRateLimiter{
            // 1 request every 12 seconds, burst of 1
            limiter: rate.NewLimiter(rate.Every(12*time.Second), 1),
        }
    })
    return instance
}

// Wait blocks until the rate limiter allows an event to occur
func (rl *APIRateLimiter) Wait(ctx context.Context) error {
    rl.mutex.RLock()
    defer rl.mutex.RUnlock()
    return rl.limiter.Wait(ctx)
}

// UpdateLimit allows dynamically updating the rate limit if needed
func (rl *APIRateLimiter) UpdateLimit(requestsPerSecond rate.Limit, burst int) {
    rl.mutex.Lock()
    defer rl.mutex.Unlock()
    rl.limiter.SetLimit(requestsPerSecond)
    rl.limiter.SetBurst(burst)
}