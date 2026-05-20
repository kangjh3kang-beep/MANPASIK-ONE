package canary_test

import (
	"testing"
	"time"

	"github.com/manpasik/backend/shared/ops/canary"
)

func newTestRollout(id string, phases []canary.Phase) *canary.Rollout {
	return &canary.Rollout{
		ID:          id,
		ServiceName: "test-service",
		NewVersion:  "v2.0.0",
		OldVersion:  "v1.0.0",
		Phases:      phases,
	}
}

func TestManager_Start(t *testing.T) {
	m := canary.NewManager()
	r := newTestRollout("r1", canary.DefaultPhases)
	if err := m.Start(r); err != nil {
		t.Fatalf("Start 실패: %v", err)
	}
	if r.State != canary.StateInProgress {
		t.Errorf("State = %q", r.State)
	}
	if r.CurrentTrafficPercent() != 5 {
		t.Errorf("Traffic = %d, want 5", r.CurrentTrafficPercent())
	}
}

func TestManager_StartDuplicate(t *testing.T) {
	m := canary.NewManager()
	r := newTestRollout("dup", canary.DefaultPhases)
	_ = m.Start(r)
	if err := m.Start(r); err == nil {
		t.Error("중복 Start 통과됨")
	}
}

func TestManager_StartRequiredFields(t *testing.T) {
	m := canary.NewManager()
	if err := m.Start(&canary.Rollout{}); err == nil {
		t.Error("ID 없이 통과")
	}
	if err := m.Start(&canary.Rollout{ID: "x"}); err == nil {
		t.Error("ServiceName 없이 통과")
	}
}

func TestManager_RecordMetricsHealthy(t *testing.T) {
	m := canary.NewManager()
	r := newTestRollout("metrics", []canary.Phase{
		{Name: "p1", TrafficPercent: 5, MinDuration: 0, MaxErrorRate: 0.01, MaxLatencyP99Ms: 1000, MinSuccessChecks: 3},
		{Name: "p2", TrafficPercent: 100, MaxErrorRate: 0.01, MaxLatencyP99Ms: 1000},
	})
	_ = m.Start(r)

	for i := 0; i < 3; i++ {
		_ = m.RecordMetrics("metrics", &canary.MetricSnapshot{
			ErrorRate:    0.005,
			LatencyP99Ms: 500,
		})
	}
	got, _ := m.Get("metrics")
	if got.SuccessChecks != 3 {
		t.Errorf("SuccessChecks = %d, want 3", got.SuccessChecks)
	}
}

func TestManager_RecordMetricsResetOnBreach(t *testing.T) {
	m := canary.NewManager()
	r := newTestRollout("breach", []canary.Phase{
		{Name: "p1", TrafficPercent: 5, MaxErrorRate: 0.01, MaxLatencyP99Ms: 1000, MinSuccessChecks: 3},
	})
	_ = m.Start(r)

	_ = m.RecordMetrics("breach", &canary.MetricSnapshot{ErrorRate: 0.005, LatencyP99Ms: 500})
	_ = m.RecordMetrics("breach", &canary.MetricSnapshot{ErrorRate: 0.005, LatencyP99Ms: 500})
	// 임계값 초과 → 리셋
	_ = m.RecordMetrics("breach", &canary.MetricSnapshot{ErrorRate: 0.05, LatencyP99Ms: 500})

	got, _ := m.Get("breach")
	if got.SuccessChecks != 0 {
		t.Errorf("SuccessChecks = %d, want 0 (reset)", got.SuccessChecks)
	}
}

func TestManager_AdvanceWhenReady(t *testing.T) {
	m := canary.NewManager()
	phases := []canary.Phase{
		{Name: "p1", TrafficPercent: 5, MinDuration: 0, MaxErrorRate: 0.01, MaxLatencyP99Ms: 1000, MinSuccessChecks: 1},
		{Name: "p2", TrafficPercent: 100, MinDuration: 0, MaxErrorRate: 0.01, MaxLatencyP99Ms: 1000, MinSuccessChecks: 0},
	}
	r := newTestRollout("ready", phases)
	_ = m.Start(r)

	_ = m.RecordMetrics("ready", &canary.MetricSnapshot{ErrorRate: 0.005, LatencyP99Ms: 500})

	advanced, err := m.AdvanceIfReady("ready")
	if err != nil {
		t.Fatalf("AdvanceIfReady 실패: %v", err)
	}
	if !advanced {
		t.Error("진행되지 않음")
	}

	got, _ := m.Get("ready")
	// 마지막 단계 진입 + MinDuration 0이므로 완료
	if got.State != canary.StateCompleted {
		t.Errorf("State = %q, want completed", got.State)
	}
}

func TestManager_AdvanceBlockedByDuration(t *testing.T) {
	m := canary.NewManager()
	phases := []canary.Phase{
		{Name: "p1", TrafficPercent: 5, MinDuration: 1 * time.Hour, MinSuccessChecks: 1, MaxErrorRate: 0.01, MaxLatencyP99Ms: 1000},
		{Name: "p2", TrafficPercent: 100},
	}
	r := newTestRollout("dur", phases)
	_ = m.Start(r)
	_ = m.RecordMetrics("dur", &canary.MetricSnapshot{ErrorRate: 0.001, LatencyP99Ms: 100})

	advanced, _ := m.AdvanceIfReady("dur")
	if advanced {
		t.Error("MinDuration 미달인데 진행됨")
	}
}

func TestManager_Rollback(t *testing.T) {
	m := canary.NewManager()
	r := newTestRollout("rb", canary.DefaultPhases)
	_ = m.Start(r)

	if err := m.Rollback("rb", "high error rate"); err != nil {
		t.Fatalf("Rollback 실패: %v", err)
	}
	got, _ := m.Get("rb")
	if got.State != canary.StateRolledBack {
		t.Errorf("State = %q, want rolled_back", got.State)
	}
	if got.RollbackReason == "" {
		t.Error("RollbackReason 미기록")
	}
}

func TestManager_AutoRollback_OnHighError(t *testing.T) {
	m := canary.NewManager()
	r := newTestRollout("auto", canary.DefaultPhases)
	_ = m.Start(r)

	_ = m.RecordMetrics("auto", &canary.MetricSnapshot{ErrorRate: 0.05, LatencyP99Ms: 500})

	threshold := &canary.MetricSnapshot{ErrorRate: 0.02, LatencyP99Ms: 1000}
	if err := m.AutoRollback("auto", threshold); err != nil {
		t.Fatalf("AutoRollback 실패: %v", err)
	}
	got, _ := m.Get("auto")
	if got.State != canary.StateRolledBack {
		t.Errorf("State = %q, want rolled_back", got.State)
	}
}

func TestManager_AutoRollback_OnHighLatency(t *testing.T) {
	m := canary.NewManager()
	r := newTestRollout("lat", canary.DefaultPhases)
	_ = m.Start(r)
	_ = m.RecordMetrics("lat", &canary.MetricSnapshot{ErrorRate: 0.001, LatencyP99Ms: 5000})

	threshold := &canary.MetricSnapshot{ErrorRate: 0.02, LatencyP99Ms: 1000}
	if err := m.AutoRollback("lat", threshold); err != nil {
		t.Fatalf("AutoRollback 실패: %v", err)
	}
	got, _ := m.Get("lat")
	if got.State != canary.StateRolledBack {
		t.Errorf("State = %q", got.State)
	}
}

func TestManager_PauseResume(t *testing.T) {
	m := canary.NewManager()
	r := newTestRollout("pr", canary.DefaultPhases)
	_ = m.Start(r)

	if err := m.Pause("pr", "manual hold"); err != nil {
		t.Fatalf("Pause 실패: %v", err)
	}
	got, _ := m.Get("pr")
	if got.State != canary.StatePaused {
		t.Errorf("State = %q, want paused", got.State)
	}

	if err := m.Resume("pr"); err != nil {
		t.Fatalf("Resume 실패: %v", err)
	}
	got, _ = m.Get("pr")
	if got.State != canary.StateInProgress {
		t.Errorf("State = %q, want in_progress", got.State)
	}
}

func TestManager_DefaultPhasesUsed(t *testing.T) {
	m := canary.NewManager()
	r := &canary.Rollout{
		ID: "default", ServiceName: "x", NewVersion: "v1",
	}
	_ = m.Start(r)

	got, _ := m.Get("default")
	if len(got.Phases) != len(canary.DefaultPhases) {
		t.Errorf("Phases = %d, want %d", len(got.Phases), len(canary.DefaultPhases))
	}
}

func TestRollout_CurrentTrafficPercent_OutOfRange(t *testing.T) {
	r := &canary.Rollout{Phases: canary.DefaultPhases, CurrentPhase: -1}
	if r.CurrentTrafficPercent() != 0 {
		t.Errorf("음수 phase에서 traffic != 0")
	}

	r.CurrentPhase = 999
	if r.CurrentTrafficPercent() != 0 {
		t.Errorf("범위 초과 phase에서 traffic != 0")
	}
}

func TestManager_ListByService(t *testing.T) {
	m := canary.NewManager()
	for i := 0; i < 3; i++ {
		_ = m.Start(&canary.Rollout{
			ID: "svc-a-" + string(rune('a'+i)), ServiceName: "service-a", NewVersion: "v1",
		})
	}
	_ = m.Start(&canary.Rollout{ID: "svc-b-1", ServiceName: "service-b", NewVersion: "v1"})

	a := m.ListByService("service-a")
	if len(a) != 3 {
		t.Errorf("service-a = %d, want 3", len(a))
	}
}

func TestManager_CountByState(t *testing.T) {
	m := canary.NewManager()
	_ = m.Start(&canary.Rollout{ID: "a", ServiceName: "s", NewVersion: "v"})
	_ = m.Start(&canary.Rollout{ID: "b", ServiceName: "s", NewVersion: "v"})
	_ = m.Rollback("b", "test")

	counts := m.CountByState()
	if counts[canary.StateInProgress] != 1 {
		t.Errorf("in_progress = %d, want 1", counts[canary.StateInProgress])
	}
	if counts[canary.StateRolledBack] != 1 {
		t.Errorf("rolled_back = %d, want 1", counts[canary.StateRolledBack])
	}
}

func TestManager_RollbackCompletedRejected(t *testing.T) {
	m := canary.NewManager()
	phases := []canary.Phase{
		{Name: "p", TrafficPercent: 100, MinDuration: 0, MaxErrorRate: 1, MaxLatencyP99Ms: 99999, MinSuccessChecks: 0},
	}
	r := newTestRollout("done", phases)
	_ = m.Start(r)
	_, _ = m.AdvanceIfReady("done")

	if err := m.Rollback("done", "x"); err == nil {
		t.Error("완료된 롤아웃 롤백이 통과됨")
	}
}
