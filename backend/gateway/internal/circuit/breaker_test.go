package circuit

import (
	"testing"
	"time"
)

func TestNewBreaker(t *testing.T) {
	b := New("test-svc", 5, 30*time.Second)
	if b.State() != Closed {
		t.Errorf("expected Closed, got %v", b.State())
	}
	if b.Name() != "test-svc" {
		t.Errorf("expected name test-svc, got %s", b.Name())
	}
}

func TestClosedAllowsRequests(t *testing.T) {
	b := New("test", 3, time.Second)
	if !b.Allow() {
		t.Error("closed breaker should allow requests")
	}
}

func TestOpenAfterThreshold(t *testing.T) {
	b := New("test", 3, 10*time.Second)

	b.RecordFailure()
	b.RecordFailure()
	if b.State() != Closed {
		t.Error("should still be closed after 2 failures")
	}

	b.RecordFailure()
	if b.State() != Open {
		t.Errorf("expected Open after 3 failures, got %v", b.State())
	}

	if b.Allow() {
		t.Error("open breaker should not allow requests")
	}
}

func TestOpenToHalfOpenAfterTimeout(t *testing.T) {
	b := New("test", 2, 50*time.Millisecond)

	b.RecordFailure()
	b.RecordFailure()
	if b.State() != Open {
		t.Fatal("expected Open")
	}

	time.Sleep(60 * time.Millisecond)

	if !b.Allow() {
		t.Error("should allow after recovery timeout (transition to HalfOpen)")
	}
	if b.State() != HalfOpen {
		t.Errorf("expected HalfOpen, got %v", b.State())
	}
}

func TestHalfOpenToClosedAfterSuccesses(t *testing.T) {
	b := New("test", 2, 10*time.Millisecond)

	// Open the breaker
	b.RecordFailure()
	b.RecordFailure()

	// Wait for recovery
	time.Sleep(15 * time.Millisecond)
	b.Allow() // triggers HalfOpen

	// 3 consecutive successes close the circuit
	b.RecordSuccess()
	b.RecordSuccess()
	b.RecordSuccess()

	if b.State() != Closed {
		t.Errorf("expected Closed after 3 successes in HalfOpen, got %v", b.State())
	}
}

func TestHalfOpenFailureReopens(t *testing.T) {
	b := New("test", 2, 10*time.Millisecond)

	b.RecordFailure()
	b.RecordFailure()

	time.Sleep(15 * time.Millisecond)
	b.Allow() // triggers HalfOpen

	b.RecordFailure()
	if b.State() != Open {
		t.Errorf("expected Open after failure in HalfOpen, got %v", b.State())
	}
}

func TestSuccessResetsFailureCount(t *testing.T) {
	b := New("test", 3, time.Second)

	b.RecordFailure()
	b.RecordFailure()
	b.RecordSuccess()
	b.RecordFailure()

	if b.State() != Closed {
		t.Error("success should reset failure count, preventing Open")
	}
}

func TestStats(t *testing.T) {
	b := New("auth", 5, 30*time.Second)
	b.RecordFailure()
	b.RecordFailure()

	stats := b.Stats()
	if stats.Name != "auth" {
		t.Errorf("expected name auth, got %s", stats.Name)
	}
	if stats.FailureCount != 2 {
		t.Errorf("expected 2 failures, got %d", stats.FailureCount)
	}
	if stats.State != "closed" {
		t.Errorf("expected closed, got %s", stats.State)
	}
}

func TestStateString(t *testing.T) {
	tests := []struct {
		state State
		want  string
	}{
		{Closed, "closed"},
		{Open, "open"},
		{HalfOpen, "half-open"},
		{State(99), "unknown"},
	}
	for _, tt := range tests {
		if got := tt.state.String(); got != tt.want {
			t.Errorf("State(%d).String() = %q, want %q", tt.state, got, tt.want)
		}
	}
}
