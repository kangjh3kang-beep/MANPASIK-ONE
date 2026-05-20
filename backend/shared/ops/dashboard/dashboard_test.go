package dashboard_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/manpasik/backend/shared/ops/dashboard"
)

// stubProvider 는 정해진 ModuleSnapshot 을 반환.
type stubProvider struct {
	name string
	snap dashboard.ModuleSnapshot
	// delay 가 있으면 호출이 지연됨 (timeout 테스트용).
	delay time.Duration
}

func (s *stubProvider) Name() string { return s.name }
func (s *stubProvider) Health(ctx context.Context) dashboard.ModuleSnapshot {
	if s.delay > 0 {
		select {
		case <-time.After(s.delay):
		case <-ctx.Done():
			return dashboard.ModuleSnapshot{
				Name:    s.name,
				Status:  dashboard.StatusUnknown,
				Message: "ctx canceled",
			}
		}
	}
	if s.snap.Name == "" {
		s.snap.Name = s.name
	}
	return s.snap
}

func TestAggregator_Empty(t *testing.T) {
	a := dashboard.NewAggregator()
	snap := a.Snapshot(context.Background())
	if snap.OverallStatus != dashboard.StatusUnknown {
		t.Errorf("빈 aggregator status = %q", snap.OverallStatus)
	}
}

func TestAggregator_AllHealthy(t *testing.T) {
	a := dashboard.NewAggregator()
	a.Register(&stubProvider{name: "secrets", snap: dashboard.ModuleSnapshot{Status: dashboard.StatusHealthy}})
	a.Register(&stubProvider{name: "tracing", snap: dashboard.ModuleSnapshot{Status: dashboard.StatusHealthy}})

	snap := a.Snapshot(context.Background())
	if snap.OverallStatus != dashboard.StatusHealthy {
		t.Errorf("OverallStatus = %q, want healthy", snap.OverallStatus)
	}
	if snap.HealthyCount != 2 {
		t.Errorf("HealthyCount = %d", snap.HealthyCount)
	}
}

func TestAggregator_WorstWins(t *testing.T) {
	a := dashboard.NewAggregator()
	a.Register(&stubProvider{name: "a", snap: dashboard.ModuleSnapshot{Status: dashboard.StatusHealthy}})
	a.Register(&stubProvider{name: "b", snap: dashboard.ModuleSnapshot{Status: dashboard.StatusDegraded}})
	a.Register(&stubProvider{name: "c", snap: dashboard.ModuleSnapshot{Status: dashboard.StatusFailed}})

	snap := a.Snapshot(context.Background())
	if snap.OverallStatus != dashboard.StatusFailed {
		t.Errorf("OverallStatus = %q, want failed", snap.OverallStatus)
	}
	if snap.FailedCount != 1 || snap.DegradedCount != 1 || snap.HealthyCount != 1 {
		t.Errorf("count 분포 = h%d/d%d/f%d", snap.HealthyCount, snap.DegradedCount, snap.FailedCount)
	}
}

func TestAggregator_SortedByName(t *testing.T) {
	a := dashboard.NewAggregator()
	a.Register(&stubProvider{name: "Zoo", snap: dashboard.ModuleSnapshot{Status: dashboard.StatusHealthy}})
	a.Register(&stubProvider{name: "alpha", snap: dashboard.ModuleSnapshot{Status: dashboard.StatusHealthy}})
	a.Register(&stubProvider{name: "Mike", snap: dashboard.ModuleSnapshot{Status: dashboard.StatusHealthy}})

	snap := a.Snapshot(context.Background())
	if snap.Modules[0].Name != "alpha" || snap.Modules[1].Name != "Mike" || snap.Modules[2].Name != "Zoo" {
		t.Errorf("정렬 실패: %v", []string{
			snap.Modules[0].Name, snap.Modules[1].Name, snap.Modules[2].Name,
		})
	}
}

func TestAggregator_Timeout(t *testing.T) {
	a := dashboard.NewAggregator()
	a.SetTimeout(50 * time.Millisecond)
	a.Register(&stubProvider{name: "slow", delay: 200 * time.Millisecond,
		snap: dashboard.ModuleSnapshot{Status: dashboard.StatusHealthy}})

	snap := a.Snapshot(context.Background())
	if len(snap.Modules) != 1 {
		t.Fatalf("Modules len = %d", len(snap.Modules))
	}
	if snap.Modules[0].Status != dashboard.StatusUnknown {
		t.Errorf("타임아웃 모듈 status = %q", snap.Modules[0].Status)
	}
}

func TestAggregator_ProviderPanic(t *testing.T) {
	a := dashboard.NewAggregator()
	panicker := dashboard.NewSimpleProvider("panic", func(ctx context.Context) dashboard.ModuleSnapshot {
		panic("intentional")
	})
	a.Register(panicker)

	snap := a.Snapshot(context.Background())
	if len(snap.Modules) != 1 {
		t.Fatalf("Modules len = %d", len(snap.Modules))
	}
	if snap.Modules[0].Status != dashboard.StatusFailed {
		t.Errorf("패닉 모듈 status = %q", snap.Modules[0].Status)
	}
}

type stubHealthChecker struct{ err error }

func (s *stubHealthChecker) HealthCheck(ctx context.Context) error { return s.err }

func TestAdapterFromHealthChecker_OK(t *testing.T) {
	hc := &stubHealthChecker{err: nil}
	prov := dashboard.AdapterFromHealthChecker("svc", hc)
	m := prov.Health(context.Background())
	if m.Status != dashboard.StatusHealthy {
		t.Errorf("status = %q", m.Status)
	}
}

func TestAdapterFromHealthChecker_Failed(t *testing.T) {
	hc := &stubHealthChecker{err: errors.New("conn refused")}
	prov := dashboard.AdapterFromHealthChecker("svc", hc)
	m := prov.Health(context.Background())
	if m.Status != dashboard.StatusFailed {
		t.Errorf("status = %q, want failed", m.Status)
	}
	if m.Message != "conn refused" {
		t.Errorf("message = %q", m.Message)
	}
}

type stubPending struct{ count int }

func (s *stubPending) PendingCount() int { return s.count }

func TestAdapterFromPendingCounter(t *testing.T) {
	tests := []struct {
		count int
		want  dashboard.Status
	}{
		{0, dashboard.StatusHealthy},
		{50, dashboard.StatusHealthy},
		{150, dashboard.StatusDegraded},
		{500, dashboard.StatusFailed},
	}
	for _, tt := range tests {
		prov := dashboard.AdapterFromPendingCounter("queue", &stubPending{count: tt.count}, 100, 300)
		m := prov.Health(context.Background())
		if m.Status != tt.want {
			t.Errorf("count=%d status = %q, want %q", tt.count, m.Status, tt.want)
		}
		if m.Metrics["pending"] != tt.count {
			t.Errorf("metrics pending = %v", m.Metrics["pending"])
		}
	}
}

type stubMetrics struct{ data map[string]interface{} }

func (s *stubMetrics) CollectMetrics() map[string]interface{} { return s.data }

func TestAdapterFromMetricsCollector(t *testing.T) {
	mc := &stubMetrics{data: map[string]interface{}{
		"status":  "degraded",
		"message": "high latency",
		"latency_p99": 2500,
	}}
	prov := dashboard.AdapterFromMetricsCollector("api", mc)
	m := prov.Health(context.Background())
	if m.Status != dashboard.StatusDegraded {
		t.Errorf("status = %q", m.Status)
	}
	if m.Metrics["latency_p99"] != 2500 {
		t.Error("메트릭 누락")
	}
}

func TestSimpleProvider_NoFn(t *testing.T) {
	p := dashboard.NewSimpleProvider("nil", nil)
	m := p.Health(context.Background())
	if m.Status != dashboard.StatusUnknown {
		t.Errorf("status = %q", m.Status)
	}
}

func TestAggregator_RegisterNil(t *testing.T) {
	a := dashboard.NewAggregator()
	a.Register(nil) // 안전해야 함
	snap := a.Snapshot(context.Background())
	if len(snap.Modules) != 0 {
		t.Error("nil provider 가 등록됨")
	}
}

func TestAggregator_AutoFillName(t *testing.T) {
	a := dashboard.NewAggregator()
	// SimpleProvider 가 빈 이름의 ModuleSnapshot 을 반환해도 자동으로 채워야 함
	a.Register(dashboard.NewSimpleProvider("xyz", func(ctx context.Context) dashboard.ModuleSnapshot {
		return dashboard.ModuleSnapshot{Status: dashboard.StatusHealthy}
	}))
	snap := a.Snapshot(context.Background())
	if snap.Modules[0].Name != "xyz" {
		t.Errorf("Name = %q", snap.Modules[0].Name)
	}
}
