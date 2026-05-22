// Package circuit implements the circuit breaker pattern (H6: 실패 격리).
//
// 서비스 장애 시 전체 게이트웨이 중단을 방지합니다.
// 개별 gRPC 서비스별로 독립적인 Breaker를 유지하여,
// 하나의 백엔드 서비스 장애가 다른 서비스로 전파되지 않도록 합니다.
package circuit

import (
	"fmt"
	"sync"
	"time"
)

// State represents the current circuit breaker state.
type State int

const (
	// Closed — 정상 운영. 요청을 통과시킵니다.
	Closed State = iota
	// Open — 장애 감지. 요청을 차단합니다.
	Open
	// HalfOpen — 복구 시도. 제한된 요청만 허용합니다.
	HalfOpen
)

// String returns a human-readable state name.
func (s State) String() string {
	switch s {
	case Closed:
		return "closed"
	case Open:
		return "open"
	case HalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// Breaker implements the circuit breaker pattern with three states:
// Closed (normal) -> Open (failing) -> HalfOpen (recovery test) -> Closed.
type Breaker struct {
	mu               sync.Mutex
	name             string
	state            State
	failureCount     int
	successCount     int
	failureThreshold int
	recoveryTimeout  time.Duration
	lastFailure      time.Time
	halfOpenMax      int // HalfOpen에서 Closed로 전환에 필요한 연속 성공 수
}

// New creates a new circuit breaker.
//
// failureThreshold: Open 상태로 전환되기 전 허용되는 연속 실패 수.
// recoveryTimeout: Open 상태에서 HalfOpen으로 전환되기까지 대기 시간.
func New(name string, failureThreshold int, recoveryTimeout time.Duration) *Breaker {
	return &Breaker{
		name:             name,
		state:            Closed,
		failureThreshold: failureThreshold,
		recoveryTimeout:  recoveryTimeout,
		halfOpenMax:      3,
	}
}

// Allow checks whether a request should be permitted through the breaker.
//
// Returns true if the request is allowed, false if the circuit is open.
func (b *Breaker) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case Closed:
		return true
	case Open:
		if time.Since(b.lastFailure) > b.recoveryTimeout {
			b.state = HalfOpen
			b.successCount = 0
			return true
		}
		return false
	case HalfOpen:
		return true
	}
	return false
}

// RecordSuccess records a successful request.
//
// In HalfOpen state, after enough consecutive successes the circuit
// transitions back to Closed.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.failureCount = 0
	b.successCount++

	if b.state == HalfOpen && b.successCount >= b.halfOpenMax {
		b.state = Closed
		b.successCount = 0
	}
}

// RecordFailure records a failed request.
//
// When failures exceed the threshold the circuit transitions to Open.
// In HalfOpen state, a single failure immediately re-opens the circuit.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.failureCount++
	b.successCount = 0
	b.lastFailure = time.Now()

	if b.state == HalfOpen {
		// HalfOpen 상태에서 실패 시 즉시 Open으로 복귀
		b.state = Open
		return
	}

	if b.failureCount >= b.failureThreshold {
		b.state = Open
	}
}

// State returns the current breaker state.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}

// Name returns the breaker name (typically the service name).
func (b *Breaker) Name() string {
	return b.name
}

// Stats returns diagnostic information about the breaker.
func (b *Breaker) Stats() BreakerStats {
	b.mu.Lock()
	defer b.mu.Unlock()
	return BreakerStats{
		Name:             b.name,
		State:            b.state.String(),
		FailureCount:     b.failureCount,
		SuccessCount:     b.successCount,
		FailureThreshold: b.failureThreshold,
		LastFailure:      b.lastFailure,
	}
}

// BreakerStats holds diagnostic information for a circuit breaker.
type BreakerStats struct {
	Name             string    `json:"name"`
	State            string    `json:"state"`
	FailureCount     int       `json:"failure_count"`
	SuccessCount     int       `json:"success_count"`
	FailureThreshold int       `json:"failure_threshold"`
	LastFailure      time.Time `json:"last_failure,omitempty"`
}

// String returns a human-readable summary of the breaker stats.
func (s BreakerStats) String() string {
	return fmt.Sprintf(
		"[%s] state=%s failures=%d/%d successes=%d",
		s.Name, s.State, s.FailureCount, s.FailureThreshold, s.SuccessCount,
	)
}
