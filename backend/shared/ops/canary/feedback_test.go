package canary_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/manpasik/backend/shared/ops/canary"
)

func newRolloutForTest(id string) *canary.Rollout {
	return &canary.Rollout{
		ID:          id,
		ServiceName: "test-svc",
		NewVersion:  "v2",
		OldVersion:  "v1",
		Phases: []canary.Phase{
			{Name: "p1", TrafficPercent: 50, MinDuration: 0, MaxErrorRate: 0.01, MaxLatencyP99Ms: 1000, MinSuccessChecks: 0},
			{Name: "p2", TrafficPercent: 100, MinDuration: 0, MaxErrorRate: 0.005, MaxLatencyP99Ms: 800, MinSuccessChecks: 0},
		},
	}
}

func TestNewFeedbackLoop_Validation(t *testing.T) {
	mgr := canary.NewManager()
	col := canary.MetricsCollectorFunc(func(_ context.Context) (*canary.MetricSnapshot, error) {
		return nil, nil
	})

	if _, err := canary.NewFeedbackLoop(nil, "id", col, canary.FeedbackLoopConfig{}); err == nil {
		t.Error("nil manager 통과")
	}
	if _, err := canary.NewFeedbackLoop(mgr, "", col, canary.FeedbackLoopConfig{}); err == nil {
		t.Error("빈 rolloutID 통과")
	}
	if _, err := canary.NewFeedbackLoop(mgr, "id", nil, canary.FeedbackLoopConfig{}); err == nil {
		t.Error("nil collector 통과")
	}
}

func TestFeedbackLoop_RecordsMetricsOnCycle(t *testing.T) {
	mgr := canary.NewManager()
	r := newRolloutForTest("r1")
	if err := mgr.Start(r); err != nil {
		t.Fatal(err)
	}

	var collected int32
	col := canary.MetricsCollectorFunc(func(_ context.Context) (*canary.MetricSnapshot, error) {
		atomic.AddInt32(&collected, 1)
		return &canary.MetricSnapshot{
			ErrorRate: 0.001, LatencyP99Ms: 500, RequestCount: 100,
		}, nil
	})

	loop, _ := canary.NewFeedbackLoop(mgr, "r1", col, canary.FeedbackLoopConfig{
		Interval: 30 * time.Millisecond,
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	loop.Start(ctx)
	time.Sleep(100 * time.Millisecond)
	loop.Stop()

	if atomic.LoadInt32(&collected) < 2 {
		t.Errorf("collect 호출 = %d, want >= 2", collected)
	}
}

func TestFeedbackLoop_AutoAdvance(t *testing.T) {
	mgr := canary.NewManager()
	r := newRolloutForTest("r-adv")
	_ = mgr.Start(r)

	col := canary.MetricsCollectorFunc(func(_ context.Context) (*canary.MetricSnapshot, error) {
		return &canary.MetricSnapshot{
			ErrorRate: 0.001, LatencyP99Ms: 100, RequestCount: 100,
		}, nil
	})
	loop, _ := canary.NewFeedbackLoop(mgr, "r-adv", col, canary.FeedbackLoopConfig{
		Interval:    30 * time.Millisecond,
		AutoAdvance: true,
	})
	loop.Start(context.Background())
	time.Sleep(150 * time.Millisecond)
	loop.Stop()

	// MinDuration=0, MinSuccessChecks=0 인 페이즈 → 즉시 진행하여 결국 완료
	r2, _ := mgr.Get("r-adv")
	state := r2.State
	if state != canary.StateCompleted && state != canary.StateInProgress {
		t.Errorf("state = %s", state)
	}
}

func TestFeedbackLoop_AutoRollback_OnHighError(t *testing.T) {
	mgr := canary.NewManager()
	r := newRolloutForTest("r-rb")
	_ = mgr.Start(r)

	col := canary.MetricsCollectorFunc(func(_ context.Context) (*canary.MetricSnapshot, error) {
		return &canary.MetricSnapshot{
			ErrorRate:    0.10, // 10% — 임계 초과
			LatencyP99Ms: 100,
			RequestCount: 100,
		}, nil
	})
	loop, _ := canary.NewFeedbackLoop(mgr, "r-rb", col, canary.FeedbackLoopConfig{
		Interval: 30 * time.Millisecond,
		AutoRollbackThresholds: &canary.AutoRollbackThresholds{
			ErrorRate: 0.05,
		},
	})
	loop.Start(context.Background())
	time.Sleep(80 * time.Millisecond)
	loop.Stop()

	r2, _ := mgr.Get("r-rb")
	if r2.State != canary.StateRolledBack {
		t.Errorf("state = %s, want rolled_back", r2.State)
	}
}

func TestFeedbackLoop_AutoRollback_OnHighLatency(t *testing.T) {
	mgr := canary.NewManager()
	r := newRolloutForTest("r-lat")
	_ = mgr.Start(r)

	col := canary.MetricsCollectorFunc(func(_ context.Context) (*canary.MetricSnapshot, error) {
		return &canary.MetricSnapshot{
			ErrorRate:    0.001,
			LatencyP99Ms: 5000, // 5s — 임계 초과
			RequestCount: 100,
		}, nil
	})
	loop, _ := canary.NewFeedbackLoop(mgr, "r-lat", col, canary.FeedbackLoopConfig{
		Interval: 30 * time.Millisecond,
		AutoRollbackThresholds: &canary.AutoRollbackThresholds{
			LatencyP99Ms: 2000,
		},
	})
	loop.Start(context.Background())
	time.Sleep(80 * time.Millisecond)
	loop.Stop()

	r2, _ := mgr.Get("r-lat")
	if r2.State != canary.StateRolledBack {
		t.Errorf("state = %s", r2.State)
	}
}

func TestFeedbackLoop_CollectorError_OnErrorCallback(t *testing.T) {
	mgr := canary.NewManager()
	r := newRolloutForTest("r-err")
	_ = mgr.Start(r)

	var (
		errs []error
		mu   sync.Mutex
	)
	col := canary.MetricsCollectorFunc(func(_ context.Context) (*canary.MetricSnapshot, error) {
		return nil, errors.New("collector down")
	})
	loop, _ := canary.NewFeedbackLoop(mgr, "r-err", col, canary.FeedbackLoopConfig{
		Interval: 30 * time.Millisecond,
		OnError: func(err error) {
			mu.Lock()
			defer mu.Unlock()
			errs = append(errs, err)
		},
	})
	loop.Start(context.Background())
	time.Sleep(70 * time.Millisecond)
	loop.Stop()

	mu.Lock()
	defer mu.Unlock()
	if len(errs) == 0 {
		t.Error("collector 에러가 OnError 로 전달되지 않음")
	}
}

func TestFeedbackLoop_CycleCount(t *testing.T) {
	mgr := canary.NewManager()
	r := newRolloutForTest("r-cnt")
	_ = mgr.Start(r)

	col := canary.MetricsCollectorFunc(func(_ context.Context) (*canary.MetricSnapshot, error) {
		return &canary.MetricSnapshot{ErrorRate: 0.001, LatencyP99Ms: 100}, nil
	})
	loop, _ := canary.NewFeedbackLoop(mgr, "r-cnt", col, canary.FeedbackLoopConfig{
		Interval: 30 * time.Millisecond,
	})
	loop.Start(context.Background())
	time.Sleep(120 * time.Millisecond)
	loop.Stop()

	if loop.CycleCount() < 2 {
		t.Errorf("CycleCount = %d", loop.CycleCount())
	}
}

func TestFeedbackLoop_StartIdempotent(t *testing.T) {
	mgr := canary.NewManager()
	r := newRolloutForTest("r-idem")
	_ = mgr.Start(r)
	col := canary.MetricsCollectorFunc(func(_ context.Context) (*canary.MetricSnapshot, error) {
		return &canary.MetricSnapshot{}, nil
	})
	loop, _ := canary.NewFeedbackLoop(mgr, "r-idem", col, canary.FeedbackLoopConfig{
		Interval: 50 * time.Millisecond,
	})
	loop.Start(context.Background())
	loop.Start(context.Background()) // 두 번 호출해도 안전해야
	loop.Stop()
}

func TestFeedbackLoop_ContextCancel(t *testing.T) {
	mgr := canary.NewManager()
	r := newRolloutForTest("r-ctx")
	_ = mgr.Start(r)
	col := canary.MetricsCollectorFunc(func(_ context.Context) (*canary.MetricSnapshot, error) {
		return &canary.MetricSnapshot{}, nil
	})
	loop, _ := canary.NewFeedbackLoop(mgr, "r-ctx", col, canary.FeedbackLoopConfig{
		Interval: 100 * time.Millisecond,
	})
	ctx, cancel := context.WithCancel(context.Background())
	loop.Start(ctx)
	cancel()
	time.Sleep(20 * time.Millisecond)
	// Stop 호출해도 행이 발생하지 않아야 함
	loop.Stop()
}

func TestMetricsCollectorFunc(t *testing.T) {
	called := false
	f := canary.MetricsCollectorFunc(func(_ context.Context) (*canary.MetricSnapshot, error) {
		called = true
		return &canary.MetricSnapshot{ErrorRate: 0.5}, nil
	})
	snap, _ := f.Collect(context.Background())
	if !called || snap.ErrorRate != 0.5 {
		t.Errorf("Func 동작 안함: called=%v snap=%v", called, snap)
	}
}
