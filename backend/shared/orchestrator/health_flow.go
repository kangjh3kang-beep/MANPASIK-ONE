package orchestrator

import (
	"context"
	"fmt"
	"log"
	"time"
)

// HealthFlowOrchestrator coordinates the health measurement → analysis → coaching flow
type HealthFlowOrchestrator struct {
	measurementProvider MeasurementProvider
	healthAnalyzer      HealthAnalyzer
	coachingGenerator   CoachingGenerator
	notifier            Notifier
}

// MeasurementProvider provides measurement data
type MeasurementProvider interface {
	GetSessionResults(ctx context.Context, sessionID string) (*MeasurementResult, error)
}

// MeasurementResult from a measurement session
type MeasurementResult struct {
	SessionID   string
	UserID      string
	DeviceID    string
	Biomarkers  []BiomarkerReading
	CompletedAt time.Time
}

// BiomarkerReading is a single biomarker measurement
type BiomarkerReading struct {
	Name  string
	Value float64
	Unit  string
}

// HealthAnalyzer provides AI analysis
type HealthAnalyzer interface {
	AnalyzeBiomarkers(ctx context.Context, userID string, readings []BiomarkerReading) (*AnalysisResult, error)
}

// AnalysisResult from AI inference
type AnalysisResult struct {
	HealthScore     float64
	RiskLevel       string // "low", "medium", "high", "critical"
	Anomalies       []string
	Recommendations []string
}

// CoachingGenerator generates coaching messages
type CoachingGenerator interface {
	GenerateFromAnalysis(ctx context.Context, userID string, analysis *AnalysisResult) (*CoachingAdvice, error)
}

// CoachingAdvice from the coaching engine
type CoachingAdvice struct {
	Message     string
	ActionItems []string
	Priority    string
}

// Notifier sends notifications
type Notifier interface {
	Notify(ctx context.Context, userID, title, body, priority string) error
}

// NewHealthFlowOrchestrator creates a new orchestrator
func NewHealthFlowOrchestrator(
	mp MeasurementProvider,
	ha HealthAnalyzer,
	cg CoachingGenerator,
	n Notifier,
) *HealthFlowOrchestrator {
	return &HealthFlowOrchestrator{
		measurementProvider: mp,
		healthAnalyzer:      ha,
		coachingGenerator:   cg,
		notifier:            n,
	}
}

// ProcessMeasurementCompleted handles the full flow after a measurement session ends
func (o *HealthFlowOrchestrator) ProcessMeasurementCompleted(ctx context.Context, sessionID string) (*FlowResult, error) {
	log.Printf("[HealthFlow] Processing session: %s", sessionID)

	// Step 1: Get measurement results
	measurement, err := o.measurementProvider.GetSessionResults(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("측정 결과 조회 실패: %w", err)
	}
	log.Printf("[HealthFlow] Step 1: Got %d biomarker readings", len(measurement.Biomarkers))

	// Step 2: AI analysis
	analysis, err := o.healthAnalyzer.AnalyzeBiomarkers(ctx, measurement.UserID, measurement.Biomarkers)
	if err != nil {
		return nil, fmt.Errorf("AI 분석 실패: %w", err)
	}
	log.Printf("[HealthFlow] Step 2: Health score=%.1f, risk=%s, anomalies=%d",
		analysis.HealthScore, analysis.RiskLevel, len(analysis.Anomalies))

	// Step 3: Generate coaching
	coaching, err := o.coachingGenerator.GenerateFromAnalysis(ctx, measurement.UserID, analysis)
	if err != nil {
		return nil, fmt.Errorf("코칭 생성 실패: %w", err)
	}
	log.Printf("[HealthFlow] Step 3: Coaching generated, priority=%s", coaching.Priority)

	// Step 4: Send notification based on risk level
	notifTitle := "측정 완료"
	notifBody := coaching.Message
	notifPriority := "normal"
	if analysis.RiskLevel == "high" || analysis.RiskLevel == "critical" {
		notifTitle = "건강 이상 감지"
		notifPriority = "urgent"
	}

	if err := o.notifier.Notify(ctx, measurement.UserID, notifTitle, notifBody, notifPriority); err != nil {
		log.Printf("[HealthFlow] Notification failed (non-fatal): %v", err)
	}
	log.Printf("[HealthFlow] Step 4: Notification sent")

	return &FlowResult{
		SessionID:   sessionID,
		UserID:      measurement.UserID,
		HealthScore: analysis.HealthScore,
		RiskLevel:   analysis.RiskLevel,
		Anomalies:   analysis.Anomalies,
		Coaching:    coaching.Message,
		ActionItems: coaching.ActionItems,
	}, nil
}

// FlowResult is the combined result of the health flow
type FlowResult struct {
	SessionID   string
	UserID      string
	HealthScore float64
	RiskLevel   string
	Anomalies   []string
	Coaching    string
	ActionItems []string
}
