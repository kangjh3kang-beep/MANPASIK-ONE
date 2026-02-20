package middleware

import (
	"sync"
	"time"
)

// ipRateLimiter는 IP 기반 요청 제한기입니다.
type ipRateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rate     int           // 윈도우당 최대 요청 수
	window   time.Duration // 윈도우 크기
}

type visitor struct {
	count    int
	windowAt time.Time
}

func newIPRateLimiter() *ipRateLimiter {
	rl := &ipRateLimiter{
		visitors: make(map[string]*visitor),
		rate:     100,             // 분당 100 요청
		window:   1 * time.Minute, // 1분 윈도우
	}
	// 주기적 정리
	go rl.cleanup()
	return rl
}

func (rl *ipRateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	v, exists := rl.visitors[ip]
	if !exists || now.Sub(v.windowAt) > rl.window {
		rl.visitors[ip] = &visitor{count: 1, windowAt: now}
		return true
	}

	v.count++
	return v.count <= rl.rate
}

func (rl *ipRateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, v := range rl.visitors {
			if now.Sub(v.windowAt) > 2*rl.window {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}
