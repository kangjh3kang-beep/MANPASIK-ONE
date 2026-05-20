package canary

import (
	"context"
	"testing"
	"time"

	"github.com/manpasik/backend/shared/observability"
)

// 패키지 내부 테스트로 작성 (newObservabilityCollectorFromIface 접근 위해).

type fakeSnapshotter struct {
	snap observability.MetricsSnapshot
}

func (f *fakeSnapshotter) Snapshot() observability.MetricsSnapshot { return f.snap }

func TestObservabilityCollector_Collect(t *testing.T) {
	fake := &fakeSnapshotter{
		snap: observability.MetricsSnapshot{
			TotalRequests: 1000,
			TotalErrors:   12,
			ErrorRate:     0.012,
			LatencyP99:    250 * time.Millisecond,
		},
	}
	c := newObservabilityCollectorFromIface(fake)
	snap, err := c.Collect(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if snap.ErrorRate != 0.012 {
		t.Errorf("ErrorRate = %f", snap.ErrorRate)
	}
	if snap.LatencyP99Ms != 250 {
		t.Errorf("LatencyP99Ms = %d", snap.LatencyP99Ms)
	}
	if snap.RequestCount != 1000 {
		t.Errorf("RequestCount = %d", snap.RequestCount)
	}
}

func TestObservabilityCollector_NoMetrics(t *testing.T) {
	c := newObservabilityCollectorFromIface(nil)
	if _, err := c.Collect(context.Background()); err == nil {
		t.Error("nil snapshotter 통과")
	}
}

func TestNewObservabilityCollector_Real(t *testing.T) {
	m := observability.NewMetrics()
	m.RecordRequest("/test", 100*time.Millisecond, 200)
	m.RecordRequest("/test", 200*time.Millisecond, 200)
	m.RecordRequest("/test", 300*time.Millisecond, 500)

	c := NewObservabilityCollector(m)
	snap, err := c.Collect(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	// 3 요청 / 1 에러 → 33% 에러율
	if abs(snap.ErrorRate-0.333) > 0.01 {
		t.Errorf("ErrorRate = %f, want ~0.333", snap.ErrorRate)
	}
	if snap.RequestCount != 3 {
		t.Errorf("RequestCount = %d", snap.RequestCount)
	}
	if snap.LatencyP99Ms <= 0 {
		t.Errorf("LatencyP99Ms = %d", snap.LatencyP99Ms)
	}
}

func TestWindowedObservabilityCollector_FirstCallBaseline(t *testing.T) {
	m := observability.NewMetrics()
	m.RecordRequest("/x", 100*time.Millisecond, 200)
	m.RecordRequest("/x", 100*time.Millisecond, 500)

	w := NewWindowedObservabilityCollector(m)
	snap, err := w.Collect(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	// 첫 호출은 baseline 만 설정 → RequestCount=0
	if snap.RequestCount != 0 {
		t.Errorf("첫 호출 RequestCount = %d, want 0", snap.RequestCount)
	}
	if snap.ErrorRate != 0 {
		t.Errorf("첫 호출 ErrorRate = %f", snap.ErrorRate)
	}
}

func TestWindowedObservabilityCollector_DeltaWindow(t *testing.T) {
	m := observability.NewMetrics()
	w := NewWindowedObservabilityCollector(m)

	// 첫 사이클: 5 요청 / 1 에러
	for i := 0; i < 5; i++ {
		status := 200
		if i == 0 {
			status = 500
		}
		m.RecordRequest("/x", 100*time.Millisecond, status)
	}
	_, _ = w.Collect(context.Background()) // baseline

	// 두 번째 사이클: 추가 10 요청 / 2 에러
	for i := 0; i < 10; i++ {
		status := 200
		if i < 2 {
			status = 500
		}
		m.RecordRequest("/x", 100*time.Millisecond, status)
	}
	snap, err := w.Collect(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if snap.RequestCount != 10 {
		t.Errorf("delta RequestCount = %d, want 10", snap.RequestCount)
	}
	// 10 요청 / 2 에러 = 0.2
	if abs(snap.ErrorRate-0.2) > 0.001 {
		t.Errorf("delta ErrorRate = %f, want 0.2", snap.ErrorRate)
	}
}

func TestWindowedObservabilityCollector_ZeroRequestsNoDivByZero(t *testing.T) {
	m := observability.NewMetrics()
	w := NewWindowedObservabilityCollector(m)
	_, _ = w.Collect(context.Background())
	// 두 번째: 추가 요청 없음 → 0 에러율
	snap, err := w.Collect(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if snap.ErrorRate != 0 || snap.RequestCount != 0 {
		t.Errorf("snap = %+v", snap)
	}
}

func TestObservabilityMetrics_Snapshot(t *testing.T) {
	m := observability.NewMetrics()
	for i := 0; i < 100; i++ {
		m.RecordRequest("/test", time.Duration(i)*time.Millisecond, 200)
	}
	snap := m.Snapshot()
	if snap.TotalRequests != 100 {
		t.Errorf("TotalRequests = %d", snap.TotalRequests)
	}
	if snap.LatencyP99 <= 0 {
		t.Errorf("LatencyP99 = %v", snap.LatencyP99)
	}
	// p99 of 0..99ms 는 약 99ms 또는 인접값
	if snap.LatencyP99 < 90*time.Millisecond || snap.LatencyP99 > 100*time.Millisecond {
		t.Errorf("LatencyP99 = %v, want ~99ms", snap.LatencyP99)
	}
	if snap.LatencyP50 <= 0 {
		t.Errorf("LatencyP50 = %v", snap.LatencyP50)
	}
}

func TestObservabilityMetrics_SnapshotEmpty(t *testing.T) {
	m := observability.NewMetrics()
	snap := m.Snapshot()
	if snap.TotalRequests != 0 || snap.LatencyP99 != 0 {
		t.Errorf("빈 metrics snap = %+v", snap)
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
