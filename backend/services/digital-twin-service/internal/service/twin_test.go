package service

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

type mockTwinRepo struct {
	states map[string]*TwinState
}

func newMockTwinRepo() *mockTwinRepo {
	return &mockTwinRepo{states: make(map[string]*TwinState)}
}
func (r *mockTwinRepo) Store(ctx context.Context, state *TwinState) error {
	r.states[state.SessionID] = state
	return nil
}
func (r *mockTwinRepo) Get(ctx context.Context, sessionID string) (*TwinState, error) {
	return r.states[sessionID], nil
}
func (r *mockTwinRepo) ListByDevice(ctx context.Context, deviceID string, limit int) ([]*TwinState, error) {
	var result []*TwinState
	for _, s := range r.states {
		if s.DeviceID == deviceID {
			result = append(result, s)
		}
	}
	return result, nil
}

func TestSyncTwin_Nominal(t *testing.T) {
	repo := newMockTwinRepo()
	svc := NewTwinService(zap.NewNop(), repo)

	state, action, err := svc.SyncTwin(context.Background(),
		"sess-1", "user-1", "dev-1",
		[]float64{0.01, -0.02, 0.03}, 0.1, 0.5, -0.3,
		"nominal", 0.1, 450,
	)
	if err != nil {
		t.Fatalf("SyncTwin failed: %v", err)
	}
	if state.ID == "" {
		t.Fatal("twin ID should not be empty")
	}
	if action != "continue" {
		t.Fatalf("expected action continue, got %s", action)
	}
	if state.HealthState != "nominal" {
		t.Fatalf("expected nominal, got %s", state.HealthState)
	}
}

func TestSyncTwin_DriftDetection(t *testing.T) {
	repo := newMockTwinRepo()
	svc := NewTwinService(zap.NewNop(), repo)

	state, action, err := svc.SyncTwin(context.Background(),
		"sess-2", "user-1", "dev-1",
		[]float64{0.5, 0.6, 0.7}, 2.5, 3.0, -1.0,
		"warning", 0.6, 200,
	)
	if err != nil {
		t.Fatalf("SyncTwin failed: %v", err)
	}
	if action != "recalibrate" {
		t.Fatalf("expected recalibrate, got %s", action)
	}
	if state.HealthState != "drift" {
		t.Fatalf("expected drift (server override), got %s", state.HealthState)
	}
}

func TestSyncTwin_Recalibrate(t *testing.T) {
	repo := newMockTwinRepo()
	svc := NewTwinService(zap.NewNop(), repo)

	_, action, err := svc.SyncTwin(context.Background(),
		"sess-3", "user-1", "dev-1",
		[]float64{1.0, 1.2}, 3.0, 5.0, -2.0,
		"recalibrate", 0.9, 5,
	)
	if err != nil {
		t.Fatalf("SyncTwin failed: %v", err)
	}
	if action != "replace_cartridge" {
		t.Fatalf("expected replace_cartridge (remaining < 10), got %s", action)
	}
}

func TestSyncTwin_Update(t *testing.T) {
	repo := newMockTwinRepo()
	svc := NewTwinService(zap.NewNop(), repo)

	svc.SyncTwin(context.Background(), "sess-4", "user-1", "dev-1", nil, 0.1, 0.1, -0.1, "nominal", 0.05, 400)
	state, _, _ := svc.SyncTwin(context.Background(), "sess-4", "user-1", "dev-1", nil, 0.2, 0.2, -0.2, "nominal", 0.1, 399)

	if state.MeasurementCount != 2 {
		t.Fatalf("expected 2 measurements, got %d", state.MeasurementCount)
	}
}

func TestGetTwinState_Exists(t *testing.T) {
	repo := newMockTwinRepo()
	svc := NewTwinService(zap.NewNop(), repo)

	svc.SyncTwin(context.Background(), "sess-get", "user-1", "dev-1", nil, 0.1, 0.2, -0.1, "nominal", 0.05, 400)

	state, err := svc.GetTwinState(context.Background(), "sess-get")
	if err != nil {
		t.Fatalf("GetTwinState failed: %v", err)
	}
	if state == nil {
		t.Fatal("expected non-nil state")
	}
	if state.SessionID != "sess-get" {
		t.Errorf("expected sess-get, got %s", state.SessionID)
	}
}

func TestGetTwinState_NotFound(t *testing.T) {
	svc := NewTwinService(zap.NewNop(), newMockTwinRepo())

	state, err := svc.GetTwinState(context.Background(), "nonexistent")
	if err != nil {
		t.Fatalf("GetTwinState should not error for missing: %v", err)
	}
	if state != nil {
		t.Error("expected nil state for nonexistent session")
	}
}

func TestListDeviceTwins(t *testing.T) {
	repo := newMockTwinRepo()
	svc := NewTwinService(zap.NewNop(), repo)

	svc.SyncTwin(context.Background(), "sess-a", "user-1", "dev-1", nil, 0.1, 0.2, -0.1, "nominal", 0.05, 400)
	svc.SyncTwin(context.Background(), "sess-b", "user-1", "dev-1", nil, 0.2, 0.3, -0.2, "warning", 0.35, 300)
	svc.SyncTwin(context.Background(), "sess-c", "user-1", "dev-2", nil, 0.1, 0.1, -0.1, "nominal", 0.02, 500) // different device

	twins, err := svc.ListDeviceTwins(context.Background(), "dev-1", 10)
	if err != nil {
		t.Fatalf("ListDeviceTwins failed: %v", err)
	}
	if len(twins) != 2 {
		t.Errorf("expected 2 twins for dev-1, got %d", len(twins))
	}
}

func TestListDeviceTwins_DefaultLimit(t *testing.T) {
	svc := NewTwinService(zap.NewNop(), newMockTwinRepo())

	twins, err := svc.ListDeviceTwins(context.Background(), "dev-empty", 0)
	if err != nil {
		t.Fatalf("ListDeviceTwins failed: %v", err)
	}
	if twins == nil {
		// nil is OK for empty result
	}
}

func TestSyncTwin_ServerOverridesClient(t *testing.T) {
	repo := newMockTwinRepo()
	svc := NewTwinService(zap.NewNop(), repo)

	// 클라이언트: nominal, 서버: recalibrate (CUSUM 초과)
	state, action, err := svc.SyncTwin(context.Background(),
		"sess-override", "user-1", "dev-1",
		nil, 3.0, 5.0, -5.0,
		"nominal", 0.9, 100,
	)
	if err != nil {
		t.Fatalf("SyncTwin failed: %v", err)
	}
	if state.HealthState != "recalibrate" {
		t.Errorf("서버 오버라이드: expected recalibrate, got %s", state.HealthState)
	}
	if action != "recalibrate" {
		t.Errorf("expected action recalibrate, got %s", action)
	}
}

func TestSyncTwin_NegativeCUSUM(t *testing.T) {
	repo := newMockTwinRepo()
	svc := NewTwinService(zap.NewNop(), repo)

	// cusumNeg 절대값 > cusumH → recalibrate
	state, _, err := svc.SyncTwin(context.Background(),
		"sess-neg", "user-1", "dev-1",
		nil, 0.5, 1.0, -5.0,
		"warning", 0.2, 200,
	)
	if err != nil {
		t.Fatalf("SyncTwin failed: %v", err)
	}
	if state.HealthState != "recalibrate" {
		t.Errorf("expected recalibrate for high |cusumNeg|, got %s", state.HealthState)
	}
}

func TestDetermineAction_Drift(t *testing.T) {
	svc := NewTwinService(zap.NewNop(), newMockTwinRepo())
	state := &TwinState{HealthState: "drift", RemainingMeasurements: 100}
	action := svc.determineAction(state)
	if action != "recalibrate" {
		t.Errorf("drift → expected recalibrate, got %s", action)
	}
}

func TestDetermineAction_Warning(t *testing.T) {
	svc := NewTwinService(zap.NewNop(), newMockTwinRepo())
	state := &TwinState{HealthState: "warning", RemainingMeasurements: 100}
	action := svc.determineAction(state)
	if action != "continue" {
		t.Errorf("warning → expected continue, got %s", action)
	}
}

func TestSeverityOrder(t *testing.T) {
	svc := NewTwinService(zap.NewNop(), newMockTwinRepo())
	if svc.severityOrder("nominal") >= svc.severityOrder("warning") {
		t.Error("nominal should be less severe than warning")
	}
	if svc.severityOrder("warning") >= svc.severityOrder("drift") {
		t.Error("warning should be less severe than drift")
	}
	if svc.severityOrder("drift") >= svc.severityOrder("recalibrate") {
		t.Error("drift should be less severe than recalibrate")
	}
}

func TestEvaluateHealthState(t *testing.T) {
	svc := NewTwinService(nil, newMockTwinRepo())

	tests := []struct {
		name     string
		ewma     float64
		cusumPos float64
		cusumNeg float64
		drift    float64
		expected string
	}{
		{"nominal", 0.1, 0.5, -0.3, 0.1, "nominal"},
		{"warning_drift", 1.0, 1.0, -0.5, 0.35, "warning"},
		{"warning_ewma", 2.5, 0.5, -0.3, 0.1, "warning"},
		{"drift", 1.0, 2.0, -1.0, 0.6, "drift"},
		{"recalibrate_drift", 1.0, 2.0, -1.0, 0.85, "recalibrate"},
		{"recalibrate_cusum", 1.0, 4.5, -1.0, 0.3, "recalibrate"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.evaluateHealthState(tt.ewma, tt.cusumPos, tt.cusumNeg, tt.drift)
			if result != tt.expected {
				t.Fatalf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}
