package orchestrator

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// Mock implementations
type mockMeasurementProvider struct{}

func (m *mockMeasurementProvider) GetSessionResults(ctx context.Context, sessionID string) (*MeasurementResult, error) {
	return &MeasurementResult{
		SessionID: sessionID,
		UserID:    "user-001",
		DeviceID:  "dev-001",
		Biomarkers: []BiomarkerReading{
			{Name: "blood_glucose", Value: 126.0, Unit: "mg/dL"},
			{Name: "cholesterol_total", Value: 240.0, Unit: "mg/dL"},
			{Name: "blood_pressure_systolic", Value: 145.0, Unit: "mmHg"},
		},
		CompletedAt: time.Now(),
	}, nil
}

type mockHealthAnalyzer struct{}

func (m *mockHealthAnalyzer) AnalyzeBiomarkers(ctx context.Context, userID string, readings []BiomarkerReading) (*AnalysisResult, error) {
	anomalies := []string{}
	for _, r := range readings {
		switch r.Name {
		case "blood_glucose":
			if r.Value > 125 {
				anomalies = append(anomalies, "혈당 높음: "+fmt.Sprintf("%.0f", r.Value)+" mg/dL")
			}
		case "cholesterol_total":
			if r.Value > 200 {
				anomalies = append(anomalies, "콜레스테롤 높음: "+fmt.Sprintf("%.0f", r.Value)+" mg/dL")
			}
		case "blood_pressure_systolic":
			if r.Value > 140 {
				anomalies = append(anomalies, "수축기혈압 높음: "+fmt.Sprintf("%.0f", r.Value)+" mmHg")
			}
		}
	}
	risk := "low"
	if len(anomalies) >= 2 {
		risk = "high"
	}
	if len(anomalies) >= 3 {
		risk = "critical"
	}
	score := 100.0 - float64(len(anomalies)*15)
	return &AnalysisResult{
		HealthScore:     score,
		RiskLevel:       risk,
		Anomalies:       anomalies,
		Recommendations: []string{"식이 조절 권장", "정기 검진 필요"},
	}, nil
}

type mockCoachingGenerator struct{}

func (m *mockCoachingGenerator) GenerateFromAnalysis(ctx context.Context, userID string, analysis *AnalysisResult) (*CoachingAdvice, error) {
	priority := "normal"
	if analysis.RiskLevel == "high" || analysis.RiskLevel == "critical" {
		priority = "urgent"
	}
	return &CoachingAdvice{
		Message:     fmt.Sprintf("건강 점수 %.0f점입니다. %d개의 주의 항목이 있습니다.", analysis.HealthScore, len(analysis.Anomalies)),
		ActionItems: analysis.Recommendations,
		Priority:    priority,
	}, nil
}

type mockNotifier struct {
	notifications []string
}

func (m *mockNotifier) Notify(ctx context.Context, userID, title, body, priority string) error {
	m.notifications = append(m.notifications, fmt.Sprintf("[%s] %s: %s", priority, title, body))
	return nil
}

func TestHealthFlowOrchestrator_NormalMeasurement(t *testing.T) {
	notifier := &mockNotifier{}
	orch := NewHealthFlowOrchestrator(
		&mockMeasurementProvider{},
		&mockHealthAnalyzer{},
		&mockCoachingGenerator{},
		notifier,
	)

	result, err := orch.ProcessMeasurementCompleted(context.Background(), "session-001")
	if err != nil {
		t.Fatalf("ProcessMeasurementCompleted failed: %v", err)
	}

	if result.UserID != "user-001" {
		t.Errorf("expected user-001, got %s", result.UserID)
	}
	if result.HealthScore <= 0 {
		t.Errorf("expected positive health score, got %.1f", result.HealthScore)
	}
	if result.RiskLevel == "" {
		t.Error("expected non-empty risk level")
	}
	if len(result.Anomalies) != 3 {
		t.Errorf("expected 3 anomalies, got %d", len(result.Anomalies))
	}
	if result.Coaching == "" {
		t.Error("expected non-empty coaching message")
	}
	if len(notifier.notifications) != 1 {
		t.Errorf("expected 1 notification, got %d", len(notifier.notifications))
	}

	t.Logf("✅ Health flow complete: score=%.0f, risk=%s, anomalies=%d, coaching=%q",
		result.HealthScore, result.RiskLevel, len(result.Anomalies), result.Coaching)
	t.Logf("✅ Notification sent: %s", notifier.notifications[0])
}

func TestHealthFlowOrchestrator_CriticalRisk(t *testing.T) {
	notifier := &mockNotifier{}
	orch := NewHealthFlowOrchestrator(
		&mockMeasurementProvider{},
		&mockHealthAnalyzer{},
		&mockCoachingGenerator{},
		notifier,
	)

	result, err := orch.ProcessMeasurementCompleted(context.Background(), "session-critical")
	if err != nil {
		t.Fatalf("ProcessMeasurementCompleted failed: %v", err)
	}

	// With 3 anomalies from mock, risk should be "critical"
	if result.RiskLevel != "critical" {
		t.Errorf("expected critical risk level, got %s", result.RiskLevel)
	}

	// Notification should be urgent
	if len(notifier.notifications) == 0 {
		t.Fatal("expected notification")
	}
	if notifier.notifications[0][:8] != "[urgent]" {
		t.Errorf("expected urgent notification, got: %s", notifier.notifications[0])
	}

	t.Logf("✅ Critical risk flow: risk=%s, notification=%s", result.RiskLevel, notifier.notifications[0])
}
