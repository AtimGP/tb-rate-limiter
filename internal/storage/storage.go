package storage

import (
	"context"
	"sync"
	"time"

	"tb-rate-limiter/internal/models"
)

type bucket struct {
	tokens 		float64
	lastRefill 	time.Time
}

type MemoryLimiter struct {
	mu sync.Mutex
	buckets map[string]*bucket
}

func NewMemoryLimiter() *MemoryLimiter {
	return &MemoryLimiter{
		buckets: make(map[string]*bucket),
	}
}

func (b *bucket) Result(now time.Time, limit float64, rateRefill float64) (bool, int, time.Duration) {
	duration := now.Sub(b.lastRefill).Seconds()
	b.lastRefill = now
	b.tokens += duration * rateRefill

	if b.tokens > limit {
		b.tokens = limit
	}

	if b.tokens >= 1.0 {
		b.tokens -= 1.0
		resetAfter := time.Duration((limit - b.tokens) / rateRefill * float64(time.Second))
		return true, int(b.tokens), resetAfter
	}

	wait := (1.0 - b.tokens) / rateRefill
	return false, 0, time.Duration(wait * float64(time.Second))
}

func (ml *MemoryLimiter) Verify(ctx context.Context, req models.LimiterRequest) (models.LimiterResponse, error) {
	now := time.Now()

	windowDur, err := time.ParseDuration(req.Window)
	if err != nil {
		return models.LimiterResponse{}, nil
	}

	ml.mu.Lock()
	defer ml.mu.Unlock()

	b, exists := ml.buckets[req.Key]
	if !exists {
		b = &bucket{
			tokens: float64(req.Limit),
			lastRefill: now,
		}

		ml.buckets[req.Key] = b
	}

	rateRefill := float64(req.Limit) / windowDur.Seconds()

	allowed, lastTokens, waits := b.Result(now, float64(req.Limit), rateRefill)

	return models.LimiterResponse{
		Allowed: allowed,
		Remaining: lastTokens,
		ResetAfter: waits.String(),
	}, nil
}