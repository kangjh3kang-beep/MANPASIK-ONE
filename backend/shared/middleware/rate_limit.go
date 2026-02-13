package middleware

import (
	"context"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RateLimiter implements a simple in-memory token bucket rate limiter
type RateLimiter struct {
	mu       sync.Mutex
	buckets  map[string]*bucket
	rate     int           // tokens per interval
	interval time.Duration // refill interval
	burst    int           // max tokens (burst capacity)
}

type bucket struct {
	tokens     int
	lastRefill time.Time
}

// NewRateLimiter creates a rate limiter
// rate: number of requests allowed per interval
// interval: time window
// burst: maximum burst capacity
func NewRateLimiter(rate int, interval time.Duration, burst int) *RateLimiter {
	return &RateLimiter{
		buckets:  make(map[string]*bucket),
		rate:     rate,
		interval: interval,
		burst:    burst,
	}
}

// Allow checks if a request from the given key is allowed
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	b, exists := rl.buckets[key]
	if !exists {
		rl.buckets[key] = &bucket{tokens: rl.burst - 1, lastRefill: time.Now()}
		return true
	}

	// Refill tokens based on elapsed time
	elapsed := time.Since(b.lastRefill)
	refills := int(elapsed/rl.interval) * rl.rate
	if refills > 0 {
		b.tokens += refills
		if b.tokens > rl.burst {
			b.tokens = rl.burst
		}
		b.lastRefill = time.Now()
	}

	if b.tokens <= 0 {
		return false
	}

	b.tokens--
	return true
}

// RateLimitInterceptor creates a gRPC interceptor that limits requests per user
func RateLimitInterceptor(limiter *RateLimiter) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Use user ID as rate limit key, or "anonymous" for unauthenticated
		key := "anonymous"
		if userID, ok := UserIDFromContext(ctx); ok {
			key = userID
		}

		if !limiter.Allow(key) {
			return nil, status.Error(codes.ResourceExhausted, "요청 한도 초과. 잠시 후 다시 시도해주세요.")
		}

		return handler(ctx, req)
	}
}
