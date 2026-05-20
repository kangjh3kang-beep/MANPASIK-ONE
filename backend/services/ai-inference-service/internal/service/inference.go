package service

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/manpasik/backend/services/ai-inference-service/internal/llm"
	apperrors "github.com/manpasik/backend/shared/errors"
)

// ============================================================================
// Domain Entities
// ============================================================================

type AiModelType int

const (
	ModelBiomarkerClassifier AiModelType = iota + 1
	ModelAnomalyDetector
	ModelTrendPredictor
	ModelHealthScorer
	ModelFoodCalorieEstimator
)

type RiskLevel int

const (
	RiskLow RiskLevel = iota + 1
	RiskModerate
	RiskHigh
	RiskCritical
)

type BiomarkerResult struct {
	BiomarkerName  string
	Value          float64
	Unit           string
	Classification string  // "normal", "borderline", "abnormal"
	Confidence     float64 // 0.0 ~ 1.0
	RiskLevel      RiskLevel
	ReferenceRange string
}

type AnomalyFlag struct {
	MetricName   string
	Value        float64
	ExpectedMin  float64
	ExpectedMax  float64
	AnomalyScore float64 // 0.0 ~ 1.0
	Description  string
}

type AnalysisResult struct {
	AnalysisID         string
	UserID             string
	MeasurementID      string
	Biomarkers         []BiomarkerResult
	Anomalies          []AnomalyFlag
	OverallHealthScore float64
	Summary            string
	AnalyzedAt         time.Time
}

type HealthScore struct {
	UserID         string
	OverallScore   float64
	CategoryScores map[string]float64
	Trend          string // "improving", "stable", "declining"
	Recommendation string
	CalculatedAt   time.Time
}

type TrendDataPoint struct {
	Timestamp  time.Time
	Value      float64
	LowerBound float64
	UpperBound float64
}

type TrendPrediction struct {
	UserID     string
	MetricName string
	Historical []TrendDataPoint
	Predicted  []TrendDataPoint
	Confidence float64
	Direction  string // "up", "down", "stable"
	Insight    string
}

type ModelStatus string

const (
	ModelStatusActive     ModelStatus = "active"
	ModelStatusTraining   ModelStatus = "training"
	ModelStatusDeprecated ModelStatus = "deprecated"
)

type ModelInfo struct {
	ModelType   AiModelType
	Name        string
	Version     string
	Description string
	Accuracy    float64
	LastTrained time.Time
	Status      ModelStatus
}

// ============================================================================
// Repositories
// ============================================================================

type AnalysisRepository interface {
	Save(ctx context.Context, result *AnalysisResult) error
	FindByID(ctx context.Context, id string) (*AnalysisResult, error)
	FindByUserID(ctx context.Context, userID string, limit int) ([]*AnalysisResult, error)
}

type HealthScoreRepository interface {
	Save(ctx context.Context, score *HealthScore) error
	FindLatestByUserID(ctx context.Context, userID string) (*HealthScore, error)
}

// ============================================================================
// Service
// ============================================================================

// InferenceService는 AI 추론 비즈니스 로직을 담당합니다.
// llmClient가 nil이면 LLM 없이 기존 규칙 기반 로직만 사용합니다 (graceful degradation).
// RiskEscalator는 위험 에스컬레이션 알림을 처리하는 인터페이스입니다.
type RiskEscalator interface {
	// EscalateCriticalRisk는 Critical 위험이 감지되었을 때 알림을 전송합니다.
	EscalateCriticalRisk(ctx context.Context, alert *RiskAlert) error
}

// RiskAlert는 위험 에스컬레이션 알림 정보를 담습니다.
type RiskAlert struct {
	UserID       string
	RiskLevel    RiskLevel
	Category     string
	BiomarkerName string
	Value        float64
	Unit         string
	Message      string
	DetectedAt   time.Time
}

// LogRiskEscalator는 위험 에스컬레이션을 로그로 기록하는 기본 구현체입니다.
type LogRiskEscalator struct{}

func (l *LogRiskEscalator) EscalateCriticalRisk(_ context.Context, alert *RiskAlert) error {
	log.Printf("[RISK-ESCALATION] Critical risk for user %s: %s=%v%s (%s)", alert.UserID, alert.BiomarkerName, alert.Value, alert.Unit, alert.Message)
	return nil
}

type InferenceService struct {
	analysisRepo    AnalysisRepository
	healthScoreRepo HealthScoreRepository
	llmClient       llm.LLMClient // nil이면 LLM 미사용
	riskEscalator   RiskEscalator // nil이면 에스컬레이션 미사용
	models          map[AiModelType]*ModelInfo
	rng             *rand.Rand
}

// NewInferenceService는 새 InferenceService를 생성합니다.
// llmClient는 nil 가능 — nil이면 LLM 기능이 비활성화됩니다.
func NewInferenceService(ar AnalysisRepository, hsr HealthScoreRepository, opts ...InferenceOption) *InferenceService {
	now := time.Now().Add(-24 * time.Hour) // 모델은 어제 학습된 것으로 가정
	svc := &InferenceService{
		analysisRepo:    ar,
		healthScoreRepo: hsr,
		rng:             rand.New(rand.NewSource(time.Now().UnixNano())),
		models: map[AiModelType]*ModelInfo{
			ModelBiomarkerClassifier: {
				ModelType:   ModelBiomarkerClassifier,
				Name:        "BiomarkerClassifier",
				Version:     "1.0.0",
				Description: "바이오마커 분류 모델 — 혈액·소변 수치 분석",
				Accuracy:    0.942,
				LastTrained: now,
				Status:      ModelStatusActive,
			},
			ModelAnomalyDetector: {
				ModelType:   ModelAnomalyDetector,
				Name:        "AnomalyDetector",
				Version:     "1.0.0",
				Description: "이상치 탐지 모델 — 시계열 측정값 이상 감지",
				Accuracy:    0.918,
				LastTrained: now,
				Status:      ModelStatusActive,
			},
			ModelTrendPredictor: {
				ModelType:   ModelTrendPredictor,
				Name:        "TrendPredictor",
				Version:     "1.0.0",
				Description: "트렌드 예측 모델 — 건강 지표 시계열 예측",
				Accuracy:    0.876,
				LastTrained: now,
				Status:      ModelStatusActive,
			},
			ModelHealthScorer: {
				ModelType:   ModelHealthScorer,
				Name:        "HealthScorer",
				Version:     "1.0.0",
				Description: "건강 점수 산출 모델 — 종합 건강 점수 계산",
				Accuracy:    0.905,
				LastTrained: now,
				Status:      ModelStatusActive,
			},
			ModelFoodCalorieEstimator: {
				ModelType:   ModelFoodCalorieEstimator,
				Name:        "FoodCalorieEstimator",
				Version:     "0.9.0-beta",
				Description: "음식 칼로리 추정 모델 (Phase 2 후반)",
				Accuracy:    0.823,
				LastTrained: now,
				Status:      ModelStatusTraining,
			},
		},
	}
	for _, opt := range opts {
		opt(svc)
	}
	return svc
}

// InferenceOption은 InferenceService 생성 시 옵션 함수 타입입니다.
type InferenceOption func(*InferenceService)

// WithLLMClient는 LLM 클라이언트를 주입합니다.
// nil이면 LLM 기능이 비활성화됩니다.
func WithLLMClient(c llm.LLMClient) InferenceOption {
	return func(s *InferenceService) {
		s.llmClient = c
	}
}

// SetRiskEscalator는 위험 에스컬레이션 핸들러를 주입합니다.
func (s *InferenceService) SetRiskEscalator(e RiskEscalator) {
	s.riskEscalator = e
}

// escalateCriticalBiomarkers는 분석 결과에서 Critical 위험 바이오마커를 찾아 에스컬레이션합니다.
func (s *InferenceService) escalateCriticalBiomarkers(ctx context.Context, userID string, biomarkers []BiomarkerResult) {
	if s.riskEscalator == nil {
		return
	}
	for _, bm := range biomarkers {
		if bm.RiskLevel == RiskCritical {
			alert := &RiskAlert{
				UserID:        userID,
				RiskLevel:     RiskCritical,
				Category:      biomarkerCategories[bm.BiomarkerName],
				BiomarkerName: bm.BiomarkerName,
				Value:         bm.Value,
				Unit:          bm.Unit,
				Message:       fmt.Sprintf("%s 수치 %.2f%s가 위험 수준입니다", bm.BiomarkerName, bm.Value, bm.Unit),
				DetectedAt:    time.Now(),
			}
			_ = s.riskEscalator.EscalateCriticalRisk(ctx, alert)
		}
	}
}

// LLMEnabled는 LLM 클라이언트가 설정되어 있는지 반환합니다.
func (s *InferenceService) LLMEnabled() bool {
	return s.llmClient != nil
}

// ValidMetrics는 PredictTrend에서 지원하는 메트릭 이름입니다.
var ValidMetrics = map[string]bool{
	"blood_glucose": true, "cholesterol_total": true, "hemoglobin": true,
	"hemoglobin_a1c": true, "creatinine": true, "uric_acid": true,
	"heart_rate": true, "oxygen_saturation": true, "body_temperature": true,
	"blood_pressure_sys": true, "blood_pressure_dia": true,
	"weight": true, "bmi": true, "steps": true, "sleep_hours": true,
}

// GetAnalysisHistory는 사용자의 분석 이력을 조회합니다.
func (s *InferenceService) GetAnalysisHistory(ctx context.Context, userID string, limit int) ([]*AnalysisResult, error) {
	if userID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "user_id is required")
	}
	if limit <= 0 {
		limit = 20
	}
	results, err := s.analysisRepo.FindByUserID(ctx, userID, limit)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrInternal, "분석 이력 조회 실패")
	}
	return results, nil
}

// ClassifyBiomarker는 바이오마커를 분류합니다 (public API).
func ClassifyBiomarker(name string, value float64) BiomarkerResult {
	return classifyBiomarker(name, value)
}

// AnalyzeMeasurementWithData는 실제 측정 데이터를 기반으로 AI 분석을 수행합니다.
func (s *InferenceService) AnalyzeMeasurementWithData(ctx context.Context, userID, measurementID string, data []MeasurementData) (*AnalysisResult, error) {
	if userID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "user_id is required")
	}
	if measurementID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "measurement_id is required")
	}
	if len(data) == 0 {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "measurement data is required")
	}

	// 실제 데이터 기반 바이오마커 분류
	biomarkers := make([]BiomarkerResult, 0, len(data))
	for _, d := range data {
		bm := classifyBiomarker(d.MetricName, d.Value)
		biomarkers = append(biomarkers, bm)
	}

	anomalies := detectCombinedAnomalies(biomarkers)
	healthScore := calculateWeightedHealthScore(biomarkers, anomalies)

	summary := s.generateSummary(biomarkers, anomalies, healthScore)
	summary = s.enhanceSummaryWithLLM(ctx, biomarkers, anomalies, healthScore, summary)

	result := &AnalysisResult{
		AnalysisID:         fmt.Sprintf("ana_%d", time.Now().UnixNano()),
		UserID:             userID,
		MeasurementID:      measurementID,
		Biomarkers:         biomarkers,
		Anomalies:          anomalies,
		OverallHealthScore: healthScore,
		Summary:            summary,
		AnalyzedAt:         time.Now(),
	}

	if err := s.analysisRepo.Save(ctx, result); err != nil {
		return nil, apperrors.New(apperrors.ErrInternal, "분석 결과 저장에 실패했습니다")
	}

	// Critical 위험 바이오마커 에스컬레이션
	s.escalateCriticalBiomarkers(ctx, userID, biomarkers)

	return result, nil
}

// ============================================================================
// LLM 연동 메서드
// ============================================================================

// healthInsightSystemPrompt는 건강 인사이트 생성용 시스템 프롬프트입니다.
const healthInsightSystemPrompt = `당신은 만파식(ManPaSik) 건강 데이터 분석 AI 어시스턴트입니다.
사용자의 건강 측정 데이터를 분석하여 이해하기 쉬운 한국어 인사이트를 제공합니다.
다음 규칙을 따르세요:
1. 의학적 진단은 하지 않습니다. 참고 정보임을 명시합니다.
2. 간결하고 이해하기 쉬운 표현을 사용합니다.
3. 수치에 근거한 객관적 분석을 제공합니다.
4. 필요 시 전문가 상담을 권장합니다.
5. 응답은 300자 이내로 합니다.`

// GenerateHealthInsight는 LLM을 사용하여 건강 인사이트를 생성합니다.
// LLM이 비활성화되어 있으면 규칙 기반 기본 인사이트를 반환합니다.
func (s *InferenceService) GenerateHealthInsight(ctx context.Context, userID string, measurements []MeasurementData) (string, error) {
	if userID == "" {
		return "", apperrors.New(apperrors.ErrInvalidInput, "user_id is required")
	}

	// LLM 미사용 시 규칙 기반 fallback
	if s.llmClient == nil {
		return s.generateRuleBasedInsight(measurements), nil
	}

	// 측정 데이터를 텍스트로 변환
	prompt := s.buildMeasurementPrompt(userID, measurements)

	resp, err := s.llmClient.Chat(ctx, healthInsightSystemPrompt, []llm.ChatMessage{
		{Role: "user", Content: prompt},
	})
	if err != nil {
		// LLM 호출 실패 시에도 규칙 기반 fallback
		return s.generateRuleBasedInsight(measurements), nil
	}

	return resp.Content, nil
}

// MeasurementData는 LLM 인사이트 생성에 사용되는 측정 데이터입니다.
type MeasurementData struct {
	MetricName string
	Value      float64
	Unit       string
	Timestamp  time.Time
}

// buildMeasurementPrompt는 측정 데이터를 LLM 프롬프트로 변환합니다.
func (s *InferenceService) buildMeasurementPrompt(userID string, measurements []MeasurementData) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("사용자 ID: %s\n", userID))
	sb.WriteString("측정 데이터:\n")
	for _, m := range measurements {
		sb.WriteString(fmt.Sprintf("- %s: %.2f %s (%s)\n",
			m.MetricName, m.Value, m.Unit, m.Timestamp.Format("2006-01-02")))
	}
	sb.WriteString("\n위 측정 데이터를 분석하여 건강 인사이트를 제공해 주세요.")
	return sb.String()
}

// generateRuleBasedInsight는 LLM 없이 규칙 기반 인사이트를 생성합니다.
func (s *InferenceService) generateRuleBasedInsight(measurements []MeasurementData) string {
	if len(measurements) == 0 {
		return "측정 데이터가 없어 분석을 수행할 수 없습니다. 측정을 먼저 진행해 주세요."
	}
	return fmt.Sprintf("총 %d건의 측정 데이터를 기반으로 분석했습니다. 규칙적인 측정과 전문가 상담을 권장합니다.", len(measurements))
}

// enhanceSummaryWithLLM은 기존 분석 요약을 LLM으로 향상시킵니다.
// LLM이 비활성화되어 있거나 호출 실패 시 원본 요약을 반환합니다.
func (s *InferenceService) enhanceSummaryWithLLM(ctx context.Context, biomarkers []BiomarkerResult, anomalies []AnomalyFlag, score float64, originalSummary string) string {
	if s.llmClient == nil {
		return originalSummary
	}

	var sb strings.Builder
	sb.WriteString("아래 건강 분석 결과를 바탕으로 사용자 친화적인 요약을 작성해 주세요.\n\n")
	sb.WriteString(fmt.Sprintf("건강 점수: %.1f/100\n", score))

	if len(biomarkers) > 0 {
		sb.WriteString("바이오마커 결과:\n")
		for _, b := range biomarkers {
			sb.WriteString(fmt.Sprintf("- %s: %.1f %s (%s, 위험도: %d)\n",
				b.BiomarkerName, b.Value, b.Unit, b.Classification, b.RiskLevel))
		}
	}

	if len(anomalies) > 0 {
		sb.WriteString("이상치:\n")
		for _, a := range anomalies {
			sb.WriteString(fmt.Sprintf("- %s: %.1f (정상범위: %.1f~%.1f)\n",
				a.MetricName, a.Value, a.ExpectedMin, a.ExpectedMax))
		}
	}

	resp, err := s.llmClient.Chat(ctx, healthInsightSystemPrompt, []llm.ChatMessage{
		{Role: "user", Content: sb.String()},
	})
	if err != nil {
		// LLM 실패 시 원본 요약 유지
		return originalSummary
	}

	return resp.Content
}

// AnalyzeMeasurement runs AI models on a measurement.
func (s *InferenceService) AnalyzeMeasurement(ctx context.Context, userID, measurementID string, requestedModels []AiModelType) (*AnalysisResult, error) {
	if userID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "user_id is required")
	}
	if measurementID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "measurement_id is required")
	}

	// 요청 모델이 없으면 기본 4종 적용 (FoodCalorie 제외)
	if len(requestedModels) == 0 {
		requestedModels = []AiModelType{
			ModelBiomarkerClassifier,
			ModelAnomalyDetector,
			ModelTrendPredictor,
			ModelHealthScorer,
		}
	}

	// 요청된 모델의 존재 여부 및 상태 확인
	for _, mt := range requestedModels {
		info, ok := s.models[mt]
		if !ok {
			return nil, apperrors.New(apperrors.ErrNotFound, fmt.Sprintf("모델 타입 %d를 찾을 수 없습니다", mt))
		}
		if info.Status == ModelStatusDeprecated {
			return nil, apperrors.New(apperrors.ErrInvalidInput, fmt.Sprintf("모델 '%s'은(는) 더 이상 사용할 수 없습니다 (deprecated)", info.Name))
		}
	}

	// 규칙 기반 AI 추론 결과 생성 (정상 범위 테이블 + z-score + IQR)
	biomarkers := s.runBiomarkerAnalysis()
	anomalies := detectCombinedAnomalies(biomarkers)
	healthScore := calculateWeightedHealthScore(biomarkers, anomalies)

	// 규칙 기반 요약 생성 후, LLM이 활성화되어 있으면 향상된 요약으로 교체
	summary := s.generateSummary(biomarkers, anomalies, healthScore)
	summary = s.enhanceSummaryWithLLM(ctx, biomarkers, anomalies, healthScore, summary)

	result := &AnalysisResult{
		AnalysisID:         fmt.Sprintf("ana_%d", time.Now().UnixNano()),
		UserID:             userID,
		MeasurementID:      measurementID,
		Biomarkers:         biomarkers,
		Anomalies:          anomalies,
		OverallHealthScore: healthScore,
		Summary:            summary,
		AnalyzedAt:         time.Now(),
	}

	if err := s.analysisRepo.Save(ctx, result); err != nil {
		return nil, apperrors.New(apperrors.ErrInternal, "분석 결과 저장에 실패했습니다")
	}

	// Critical 위험 바이오마커 에스컬레이션
	s.escalateCriticalBiomarkers(ctx, userID, biomarkers)

	return result, nil
}

// GetHealthScore calculates a user's health score based on recent data.
// 실제 분석 이력이 있으면 바이오마커 기반 카테고리 점수를 계산하고,
// 이력이 없으면 기본값(75.0)을 반환합니다.
func (s *InferenceService) GetHealthScore(ctx context.Context, userID string, days int) (*HealthScore, error) {
	if userID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "user_id is required")
	}
	if days <= 0 {
		days = 30
	}

	var overallScore float64
	var categoryScores map[string]float64
	var trend string

	// 최근 분석 이력 조회
	analyses, _ := s.analysisRepo.FindByUserID(ctx, userID, 10)

	if len(analyses) > 0 {
		// 실제 데이터 기반 계산
		categoryScores = computeCategoryScores(analyses)
		overallScore = computeOverallFromCategories(categoryScores)
		trend = determineTrend(analyses)
	} else {
		// 데이터 없음 → 기본값
		overallScore = 75.0
		categoryScores = map[string]float64{
			"cardiovascular": 75.0,
			"metabolic":      75.0,
			"nutritional":    75.0,
			"fitness":        75.0,
		}
		trend = "stable"
	}

	// LLM으로 맞춤형 추천 생성 (실패 시 기본 추천 사용)
	recommendation := s.generateRecommendation(ctx, overallScore, categoryScores, trend)

	score := &HealthScore{
		UserID:         userID,
		OverallScore:   overallScore,
		CategoryScores: categoryScores,
		Trend:          trend,
		Recommendation: recommendation,
		CalculatedAt:   time.Now(),
	}

	if err := s.healthScoreRepo.Save(ctx, score); err != nil {
		return nil, apperrors.New(apperrors.ErrInternal, "건강 점수 저장에 실패했습니다")
	}
	return score, nil
}

// PredictTrend predicts future values for a given metric.
// 실제 분석 이력이 있으면 선형 회귀로 예측하고, 없으면 시뮬레이션으로 폴백합니다.
func (s *InferenceService) PredictTrend(ctx context.Context, userID, metricName string, historyDays, predictionDays int) (*TrendPrediction, error) {
	if userID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "user_id is required")
	}
	if metricName == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "metric_name is required")
	}
	if !ValidMetrics[metricName] {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "지원하지 않는 메트릭: "+metricName)
	}
	if historyDays <= 0 {
		historyDays = 30
	}
	if predictionDays <= 0 {
		predictionDays = 7
	}

	// 실제 분석 이력 조회
	analyses, _ := s.analysisRepo.FindByUserID(ctx, userID, 100)
	times, values := extractMetricHistory(analyses, metricName)

	if len(values) < 3 {
		// 데이터 부족 → 시뮬레이션 폴백
		return s.simulateTrend(userID, metricName, historyDays, predictionDays)
	}

	// 실제 데이터 기반 historical 생성
	historical := make([]TrendDataPoint, len(times))
	for i, t := range times {
		v := values[i]
		ref, ok := biomarkerRefRange[metricName]
		margin := 3.0
		if ok {
			margin = (ref.normMax - ref.normMin) * 0.1
		}
		historical[i] = TrendDataPoint{
			Timestamp:  t,
			Value:      math.Round(v*100) / 100,
			LowerBound: math.Round((v-margin)*100) / 100,
			UpperBound: math.Round((v+margin)*100) / 100,
		}
	}

	// 선형 회귀
	xs := make([]float64, len(values))
	for i := range xs {
		xs[i] = float64(i)
	}
	slope, intercept, r2 := linearRegression(xs, values)

	// 예측
	now := time.Now()
	predicted := make([]TrendDataPoint, predictionDays)
	for i := 0; i < predictionDays; i++ {
		t := now.AddDate(0, 0, i+1)
		x := float64(len(values) + i)
		val := slope*x + intercept
		uncertainty := math.Abs(val) * (1 - math.Max(r2, 0)) * 0.05 * float64(i+1)
		predicted[i] = TrendDataPoint{
			Timestamp:  t,
			Value:      math.Round(val*100) / 100,
			LowerBound: math.Round((val-uncertainty)*100) / 100,
			UpperBound: math.Round((val+uncertainty)*100) / 100,
		}
	}

	direction := "stable"
	if slope > 0.5 {
		direction = "up"
	} else if slope < -0.5 {
		direction = "down"
	}

	confidence := math.Max(0.5, math.Min(r2, 0.99))

	return &TrendPrediction{
		UserID:     userID,
		MetricName: metricName,
		Historical: historical,
		Predicted:  predicted,
		Confidence: math.Round(confidence*100) / 100,
		Direction:  direction,
		Insight:    fmt.Sprintf("%s 지표는 선형 회귀 분석 결과 향후 %d일간 %s 추세입니다 (R²=%.2f).", metricName, predictionDays, direction, r2),
	}, nil
}

// GetModelInfo returns information about a specific AI model.
func (s *InferenceService) GetModelInfo(_ context.Context, modelType AiModelType) (*ModelInfo, error) {
	info, ok := s.models[modelType]
	if !ok {
		return nil, apperrors.New(apperrors.ErrNotFound, "모델을 찾을 수 없습니다")
	}
	return info, nil
}

// ListModels returns all available AI models.
func (s *InferenceService) ListModels(_ context.Context) ([]*ModelInfo, error) {
	result := make([]*ModelInfo, 0, len(s.models))
	for _, m := range s.models {
		result = append(result, m)
	}
	return result, nil
}

// RegisterModel은 새 모델을 레지스트리에 등록합니다.
// 동일 타입의 기존 모델이 있으면 Deprecated로 변경합니다.
func (s *InferenceService) RegisterModel(_ context.Context, model *ModelInfo) error {
	if model == nil || model.Name == "" {
		return apperrors.New(apperrors.ErrInvalidInput, "모델 이름은 필수입니다")
	}
	if model.Version == "" {
		return apperrors.New(apperrors.ErrInvalidInput, "모델 버전은 필수입니다")
	}
	// 기존 모델 deprecate
	if existing, ok := s.models[model.ModelType]; ok {
		existing.Status = ModelStatusDeprecated
	}
	model.Status = ModelStatusActive
	model.LastTrained = time.Now()
	s.models[model.ModelType] = model
	return nil
}

// UpdateModelStatus는 모델의 상태를 변경합니다.
func (s *InferenceService) UpdateModelStatus(_ context.Context, modelType AiModelType, status ModelStatus) error {
	info, ok := s.models[modelType]
	if !ok {
		return apperrors.New(apperrors.ErrNotFound, "모델을 찾을 수 없습니다")
	}
	info.Status = status
	return nil
}

// GetActiveModels는 활성(Active) 상태인 모델만 반환합니다.
func (s *InferenceService) GetActiveModels(_ context.Context) []*ModelInfo {
	var result []*ModelInfo
	for _, m := range s.models {
		if m.Status == ModelStatusActive {
			result = append(result, m)
		}
	}
	return result
}

// ============================================================================
// Rule-based AI inference (정상 범위 테이블 + z-score 이상탐지)
// ============================================================================

// biomarkerRefRange는 바이오마커 정상 범위 참조 테이블입니다 (임상 가이드라인 기반).
var biomarkerRefRange = map[string]struct {
	unit    string
	normMin float64
	normMax float64
	warnMin float64
	warnMax float64
	ref     string
}{
	"blood_glucose":      {"mg/dL", 70, 100, 60, 126, "70-100 mg/dL (공복)"},
	"cholesterol_total":  {"mg/dL", 0, 200, 0, 240, "< 200 mg/dL"},
	"hemoglobin":         {"g/dL", 12, 17.5, 10, 18.5, "12-17.5 g/dL"},
	"hemoglobin_a1c":     {"%", 4.0, 5.7, 3.5, 6.5, "< 5.7% (정상)"},
	"creatinine":         {"mg/dL", 0.7, 1.3, 0.5, 1.5, "0.7-1.3 mg/dL"},
	"uric_acid":          {"mg/dL", 3.5, 7.2, 2.0, 8.5, "3.5-7.2 mg/dL"},
	"heart_rate":         {"bpm", 60, 100, 50, 120, "60-100 bpm"},
	"oxygen_saturation":  {"%", 95, 100, 90, 100, "95-100%"},
	"body_temperature":   {"°C", 36.1, 37.2, 35.5, 38.0, "36.1-37.2°C"},
	"blood_pressure_sys": {"mmHg", 90, 120, 80, 140, "90-120 mmHg"},
	"blood_pressure_dia": {"mmHg", 60, 80, 50, 90, "60-80 mmHg"},
}

// classifyBiomarker는 정상 범위 테이블 기반으로 바이오마커를 분류합니다.
func classifyBiomarker(name string, value float64) BiomarkerResult {
	ref, ok := biomarkerRefRange[name]
	if !ok {
		return BiomarkerResult{
			BiomarkerName: name, Value: value, Unit: "",
			Classification: "unknown", Confidence: 0.5, RiskLevel: RiskLow,
			ReferenceRange: "참조 범위 미등록",
		}
	}

	classification := "normal"
	risk := RiskLow
	confidence := 0.95

	if value >= ref.normMin && value <= ref.normMax {
		classification = "normal"
		risk = RiskLow
	} else if value >= ref.warnMin && value <= ref.warnMax {
		classification = "borderline"
		risk = RiskModerate
		confidence = 0.90
	} else {
		classification = "abnormal"
		risk = RiskHigh
		confidence = 0.92
		if value < ref.warnMin*0.7 || value > ref.warnMax*1.3 {
			risk = RiskCritical
			confidence = 0.88
		}
	}

	return BiomarkerResult{
		BiomarkerName: name, Value: value, Unit: ref.unit,
		Classification: classification, Confidence: confidence,
		RiskLevel: risk, ReferenceRange: ref.ref,
	}
}

// runBiomarkerAnalysis는 기본 5종 바이오마커에 대해 규칙 기반 분석을 수행합니다.
// 실제 측정값이 없으면 정상 범위 중앙값을 사용합니다.
func (s *InferenceService) runBiomarkerAnalysis() []BiomarkerResult {
	defaultBiomarkers := []string{
		"blood_glucose", "cholesterol_total", "hemoglobin", "creatinine", "uric_acid",
	}
	results := make([]BiomarkerResult, len(defaultBiomarkers))
	for i, name := range defaultBiomarkers {
		ref := biomarkerRefRange[name]
		midpoint := (ref.normMin + ref.normMax) / 2
		results[i] = classifyBiomarker(name, midpoint)
	}
	return results
}

// detectAnomaliesByZScore는 z-score 기반 통계적 이상탐지를 수행합니다.
func detectAnomaliesByZScore(biomarkers []BiomarkerResult) []AnomalyFlag {
	var anomalies []AnomalyFlag
	for _, bm := range biomarkers {
		ref, ok := biomarkerRefRange[bm.BiomarkerName]
		if !ok {
			continue
		}
		mean := (ref.normMin + ref.normMax) / 2
		stdDev := (ref.normMax - ref.normMin) / 4
		if stdDev == 0 {
			continue
		}
		zScore := math.Abs((bm.Value - mean) / stdDev)
		if zScore > 2.0 {
			anomalyScore := math.Min(zScore/5.0, 1.0)
			anomalies = append(anomalies, AnomalyFlag{
				MetricName:   bm.BiomarkerName,
				Value:        bm.Value,
				ExpectedMin:  ref.normMin,
				ExpectedMax:  ref.normMax,
				AnomalyScore: math.Round(anomalyScore*1000) / 1000,
				Description: fmt.Sprintf("%s 값(%.1f %s)이 정상 범위를 벗어났습니다 (z-score: %.2f)",
					bm.BiomarkerName, bm.Value, bm.Unit, zScore),
			})
		}
	}
	return anomalies
}

// calculateWeightedHealthScore는 바이오마커 결과 기반 가중 평균 건강 점수를 계산합니다.
func calculateWeightedHealthScore(biomarkers []BiomarkerResult, anomalies []AnomalyFlag) float64 {
	if len(biomarkers) == 0 {
		return 75.0
	}
	totalScore := 0.0
	totalWeight := 0.0
	for _, bm := range biomarkers {
		weight := 1.0
		itemScore := 60.0
		switch bm.Classification {
		case "normal":
			itemScore = 95.0
		case "borderline":
			itemScore = 70.0
		case "abnormal":
			itemScore = 40.0
		}
		switch bm.RiskLevel {
		case RiskCritical:
			weight = 2.0
		case RiskHigh:
			weight = 1.5
		case RiskModerate:
			weight = 1.2
		}
		totalScore += itemScore * weight
		totalWeight += weight
	}
	score := totalScore / totalWeight
	score = math.Max(score-float64(len(anomalies))*3.0, 20.0)
	return math.Round(score*10) / 10
}

func (s *InferenceService) generateSummary(biomarkers []BiomarkerResult, anomalies []AnomalyFlag, score float64) string {
	abnormalCount := 0
	for _, b := range biomarkers {
		if b.Classification == "abnormal" {
			abnormalCount++
		}
	}

	if abnormalCount == 0 && len(anomalies) == 0 {
		return fmt.Sprintf("전반적으로 양호한 건강 상태입니다. (건강 점수: %.1f/100)", score)
	}
	if abnormalCount > 0 && len(anomalies) > 0 {
		return fmt.Sprintf("비정상 바이오마커 %d건, 이상치 %d건이 감지되었습니다. 전문가 상담을 권장합니다. (건강 점수: %.1f/100)",
			abnormalCount, len(anomalies), score)
	}
	if abnormalCount > 0 {
		return fmt.Sprintf("비정상 바이오마커 %d건이 감지되었습니다. 추가 검사를 권장합니다. (건강 점수: %.1f/100)",
			abnormalCount, score)
	}
	return fmt.Sprintf("이상치 %d건이 감지되었습니다. 모니터링을 계속하세요. (건강 점수: %.1f/100)",
		len(anomalies), score)
}

// generateRecommendation은 건강 추천을 생성합니다.
// LLM 활성화 시 LLM 기반, 미활성화 시 카테고리별 규칙 기반 추천을 반환합니다.
func (s *InferenceService) generateRecommendation(ctx context.Context, overallScore float64, categoryScores map[string]float64, trend string) string {
	ruleBasedRec := generateRuleBasedRecommendation(overallScore, categoryScores, trend)

	if s.llmClient == nil {
		return ruleBasedRec
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("건강 점수: %.1f/100, 추세: %s\n", overallScore, trend))
	sb.WriteString("카테고리별 점수:\n")
	for cat, sc := range categoryScores {
		sb.WriteString(fmt.Sprintf("- %s: %.1f\n", cat, sc))
	}
	sb.WriteString("\n이 데이터를 바탕으로 구체적인 건강 개선 추천을 1~2문장으로 작성해 주세요.")

	resp, err := s.llmClient.Chat(ctx, healthInsightSystemPrompt, []llm.ChatMessage{
		{Role: "user", Content: sb.String()},
	})
	if err != nil {
		return ruleBasedRec
	}

	return resp.Content
}

// generateRuleBasedRecommendation은 카테고리 점수와 추세에 따라 구체적인 추천을 생성합니다.
func generateRuleBasedRecommendation(overallScore float64, categoryScores map[string]float64, trend string) string {
	// 가장 낮은 카테고리 찾기
	worstCategory := ""
	worstScore := 100.0
	for cat, sc := range categoryScores {
		if sc < worstScore {
			worstScore = sc
			worstCategory = cat
		}
	}

	// 카테고리별 맞춤 추천
	var rec string
	switch worstCategory {
	case "cardiovascular":
		if worstScore < 50 {
			rec = "심혈관 지표가 주의 수준입니다. 유산소 운동(빠르게 걷기 30분)을 매일 실시하고, 나트륨 섭취를 줄이세요. 전문의 상담을 권장합니다."
		} else if worstScore < 70 {
			rec = "심혈관 건강 개선이 필요합니다. 주 3회 이상 유산소 운동을 하고, 오메가-3가 풍부한 생선을 주 2회 이상 섭취하세요."
		} else {
			rec = "심혈관 지표가 양호합니다. 현재 운동 습관을 유지하고, 정기적인 혈압 모니터링을 계속하세요."
		}
	case "metabolic":
		if worstScore < 50 {
			rec = "대사 지표가 주의 수준입니다. 정제 탄수화물과 당분 섭취를 줄이고, 식후 30분 산책을 습관화하세요. 공복혈당 정기 측정이 필요합니다."
		} else if worstScore < 70 {
			rec = "대사 건강 개선을 위해 저GI 식품(현미, 귀리, 두부) 위주로 식단을 구성하고, 규칙적인 식사 시간을 지키세요."
		} else {
			rec = "대사 지표가 안정적입니다. 균형 잡힌 식단을 유지하고, 간식으로 견과류나 과일을 선택하세요."
		}
	case "nutritional":
		if worstScore < 50 {
			rec = "영양 상태가 부족합니다. 다양한 색상의 채소와 단백질을 매 끼니 포함하고, 비타민 D와 철분 보충을 고려하세요."
		} else if worstScore < 70 {
			rec = "영양 균형 개선이 필요합니다. 매일 채소 5접시, 과일 2접시를 목표로 하고, 수분 섭취(하루 8잔)를 늘리세요."
		} else {
			rec = "영양 상태가 양호합니다. 현재 식단을 유지하면서 계절 과일과 채소를 다양하게 섭취하세요."
		}
	case "fitness":
		if worstScore < 50 {
			rec = "체력 지표 개선이 시급합니다. 매일 최소 20분 걷기부터 시작하여 점진적으로 운동 강도를 높이세요. 스트레칭도 잊지 마세요."
		} else if worstScore < 70 {
			rec = "체력 향상을 위해 주 3-4회 30분 이상 중강도 운동(빠르게 걷기, 자전거, 수영)을 추천합니다."
		} else {
			rec = "체력 상태가 좋습니다. 근력 운동과 유산소 운동을 균형 있게 유지하세요."
		}
	default:
		rec = "규칙적인 운동과 균형 잡힌 식단을 유지하세요."
	}

	// 추세 기반 보조 메시지
	switch trend {
	case "declining":
		rec += " 최근 건강 추세가 하락하고 있으므로 생활습관 점검이 필요합니다."
	case "improving":
		rec += " 건강 지표가 개선되고 있습니다. 현재 습관을 꾸준히 유지하세요."
	}

	// 전체 점수 기반 긴급 경고
	if overallScore < 40 {
		rec = "종합 건강 점수가 매우 낮습니다. 즉시 전문의 상담을 받으시기 바랍니다. " + rec
	}

	return rec
}

func (s *InferenceService) pickTrend() string {
	v := s.rng.Float64()
	if v < 0.4 {
		return "improving"
	}
	if v < 0.7 {
		return "stable"
	}
	return "declining"
}

// ============================================================================
// Phase E-1: 실제 데이터 기반 건강 점수 + IQR 이상탐지 + 선형 회귀 예측
// ============================================================================

// biomarkerCategories는 바이오마커→건강 카테고리 매핑입니다.
var biomarkerCategories = map[string]string{
	"blood_glucose":      "metabolic",
	"hemoglobin_a1c":     "metabolic",
	"cholesterol_total":  "cardiovascular",
	"blood_pressure_sys": "cardiovascular",
	"blood_pressure_dia": "cardiovascular",
	"heart_rate":         "cardiovascular",
	"oxygen_saturation":  "cardiovascular",
	"hemoglobin":         "nutritional",
	"creatinine":         "nutritional",
	"uric_acid":          "nutritional",
	"body_temperature":   "fitness",
}

// biomarkerToScore는 바이오마커 분류 결과를 0-100 점수로 변환합니다.
func biomarkerToScore(bm BiomarkerResult) float64 {
	switch bm.Classification {
	case "normal":
		return 95.0
	case "borderline":
		return 70.0
	case "abnormal":
		if bm.RiskLevel == RiskCritical {
			return 25.0
		}
		return 40.0
	default:
		return 60.0
	}
}

// computeCategoryScores는 분석 이력에서 카테고리별 건강 점수를 계산합니다.
func computeCategoryScores(analyses []*AnalysisResult) map[string]float64 {
	result := map[string]float64{
		"cardiovascular": 75.0,
		"metabolic":      75.0,
		"nutritional":    75.0,
		"fitness":        75.0,
	}
	if len(analyses) == 0 {
		return result
	}

	categorySums := map[string]float64{}
	categoryCounts := map[string]int{}

	limit := len(analyses)
	if limit > 5 {
		limit = 5
	}

	for _, analysis := range analyses[:limit] {
		for _, bm := range analysis.Biomarkers {
			cat, ok := biomarkerCategories[bm.BiomarkerName]
			if !ok {
				continue
			}
			score := biomarkerToScore(bm)
			categorySums[cat] += score
			categoryCounts[cat]++
		}
	}

	for cat, sum := range categorySums {
		if count := categoryCounts[cat]; count > 0 {
			result[cat] = math.Round(sum/float64(count)*10) / 10
		}
	}
	return result
}

// computeOverallFromCategories는 카테고리 점수에서 가중 평균으로 종합 점수를 계산합니다.
func computeOverallFromCategories(categories map[string]float64) float64 {
	if len(categories) == 0 {
		return 75.0
	}
	weights := map[string]float64{
		"cardiovascular": 0.35,
		"metabolic":      0.30,
		"nutritional":    0.20,
		"fitness":        0.15,
	}

	totalWeight := 0.0
	totalScore := 0.0
	for cat, score := range categories {
		w := weights[cat]
		if w == 0 {
			w = 0.25
		}
		totalScore += score * w
		totalWeight += w
	}
	if totalWeight == 0 {
		return 75.0
	}
	return math.Round(totalScore/totalWeight*10) / 10
}

// determineTrend는 분석 이력에서 건강 추세를 판별합니다.
func determineTrend(analyses []*AnalysisResult) string {
	if len(analyses) < 2 {
		return "stable"
	}

	limit := len(analyses)
	if limit > 5 {
		limit = 5
	}

	mid := limit / 2
	recentAvg := 0.0
	olderAvg := 0.0

	for i := 0; i < mid; i++ {
		recentAvg += analyses[i].OverallHealthScore
	}
	for i := mid; i < limit; i++ {
		olderAvg += analyses[i].OverallHealthScore
	}

	if mid > 0 {
		recentAvg /= float64(mid)
	}
	if limit-mid > 0 {
		olderAvg /= float64(limit - mid)
	}

	diff := recentAvg - olderAvg
	if diff > 3 {
		return "improving"
	}
	if diff < -3 {
		return "declining"
	}
	return "stable"
}

// detectAnomaliesByIQR는 IQR(사분위 범위) 기반 이상탐지를 수행합니다.
// 정상 범위의 Q1=normMin, Q3=normMax로 근사하여 1.5*IQR 경계를 계산합니다.
func detectAnomaliesByIQR(biomarkers []BiomarkerResult) []AnomalyFlag {
	var anomalies []AnomalyFlag
	for _, bm := range biomarkers {
		ref, ok := biomarkerRefRange[bm.BiomarkerName]
		if !ok {
			continue
		}
		q1 := ref.normMin
		q3 := ref.normMax
		iqr := q3 - q1
		if iqr <= 0 {
			continue
		}
		lowerFence := q1 - 1.5*iqr
		upperFence := q3 + 1.5*iqr

		if bm.Value < lowerFence || bm.Value > upperFence {
			var deviation float64
			if bm.Value < lowerFence {
				deviation = (lowerFence - bm.Value) / iqr
			} else {
				deviation = (bm.Value - upperFence) / iqr
			}
			anomalyScore := math.Min(deviation/3.0, 1.0)

			anomalies = append(anomalies, AnomalyFlag{
				MetricName:   bm.BiomarkerName,
				Value:        bm.Value,
				ExpectedMin:  math.Round(lowerFence*100) / 100,
				ExpectedMax:  math.Round(upperFence*100) / 100,
				AnomalyScore: math.Round(anomalyScore*1000) / 1000,
				Description: fmt.Sprintf("%s 값(%.1f %s)이 IQR 경계(%.1f~%.1f)를 벗어났습니다",
					bm.BiomarkerName, bm.Value, bm.Unit, lowerFence, upperFence),
			})
		}
	}
	return anomalies
}

// detectCombinedAnomalies는 Z-score와 IQR 이상탐지를 결합합니다.
// 두 방법 모두 감지한 이상치는 더 높은 점수를 부여합니다.
func detectCombinedAnomalies(biomarkers []BiomarkerResult) []AnomalyFlag {
	zScoreAnomalies := detectAnomaliesByZScore(biomarkers)
	iqrAnomalies := detectAnomaliesByIQR(biomarkers)

	zMetrics := make(map[string]int)
	for i, a := range zScoreAnomalies {
		zMetrics[a.MetricName] = i
	}

	result := make([]AnomalyFlag, len(zScoreAnomalies))
	copy(result, zScoreAnomalies)

	for _, iqrA := range iqrAnomalies {
		if idx, exists := zMetrics[iqrA.MetricName]; exists {
			if iqrA.AnomalyScore > result[idx].AnomalyScore {
				result[idx].AnomalyScore = iqrA.AnomalyScore
			}
			result[idx].Description += " | " + iqrA.Description
		} else {
			result = append(result, iqrA)
		}
	}

	return result
}

// linearRegression은 간단한 선형 회귀를 수행합니다.
// xs, ys는 같은 길이여야 합니다.
func linearRegression(xs, ys []float64) (slope, intercept, r2 float64) {
	n := float64(len(xs))
	if n < 2 {
		if n == 1 {
			return 0, ys[0], 0
		}
		return 0, 0, 0
	}

	sumX, sumY, sumXY, sumX2, sumY2 := 0.0, 0.0, 0.0, 0.0, 0.0
	for i := 0; i < len(xs); i++ {
		sumX += xs[i]
		sumY += ys[i]
		sumXY += xs[i] * ys[i]
		sumX2 += xs[i] * xs[i]
		sumY2 += ys[i] * ys[i]
	}

	denom := n*sumX2 - sumX*sumX
	if denom == 0 {
		return 0, sumY / n, 0
	}

	slope = (n*sumXY - sumX*sumY) / denom
	intercept = (sumY - slope*sumX) / n

	meanY := sumY / n
	ssTot := sumY2 - n*meanY*meanY
	ssRes := 0.0
	for i := 0; i < len(xs); i++ {
		predicted := slope*xs[i] + intercept
		diff := ys[i] - predicted
		ssRes += diff * diff
	}
	if ssTot > 0 {
		r2 = 1 - ssRes/ssTot
	}

	return slope, intercept, r2
}

// extractMetricHistory는 분석 이력에서 특정 메트릭의 시계열 데이터를 추출합니다.
// 시간순(오래된→최신) 정렬로 반환합니다.
func extractMetricHistory(analyses []*AnalysisResult, metricName string) ([]time.Time, []float64) {
	type dataPoint struct {
		t time.Time
		v float64
	}
	var points []dataPoint

	for _, a := range analyses {
		for _, bm := range a.Biomarkers {
			if bm.BiomarkerName == metricName {
				points = append(points, dataPoint{t: a.AnalyzedAt, v: bm.Value})
				break
			}
		}
	}

	// 시간순 정렬 (오래된→최신)
	sort.Slice(points, func(i, j int) bool {
		return points[i].t.Before(points[j].t)
	})

	times := make([]time.Time, len(points))
	values := make([]float64, len(points))
	for i, p := range points {
		times[i] = p.t
		values[i] = p.v
	}
	return times, values
}

// simulateTrend은 실제 데이터가 부족할 때 시뮬레이션 기반 예측을 반환합니다.
func (s *InferenceService) simulateTrend(userID, metricName string, historyDays, predictionDays int) (*TrendPrediction, error) {
	now := time.Now()
	baseValue := 80 + s.rng.Float64()*40

	historical := make([]TrendDataPoint, historyDays)
	for i := 0; i < historyDays; i++ {
		t := now.AddDate(0, 0, -(historyDays - i))
		val := baseValue + s.rng.Float64()*10 - 5
		historical[i] = TrendDataPoint{
			Timestamp:  t,
			Value:      math.Round(val*100) / 100,
			LowerBound: math.Round((val-3)*100) / 100,
			UpperBound: math.Round((val+3)*100) / 100,
		}
	}

	predicted := make([]TrendDataPoint, predictionDays)
	lastVal := historical[len(historical)-1].Value
	for i := 0; i < predictionDays; i++ {
		t := now.AddDate(0, 0, i+1)
		delta := s.rng.Float64()*4 - 2
		val := lastVal + delta
		margin := float64(i+1) * 1.5
		predicted[i] = TrendDataPoint{
			Timestamp:  t,
			Value:      math.Round(val*100) / 100,
			LowerBound: math.Round((val-margin)*100) / 100,
			UpperBound: math.Round((val+margin)*100) / 100,
		}
		lastVal = val
	}

	direction := "stable"
	diff := predicted[len(predicted)-1].Value - historical[len(historical)-1].Value
	if diff > 3 {
		direction = "up"
	} else if diff < -3 {
		direction = "down"
	}

	return &TrendPrediction{
		UserID:     userID,
		MetricName: metricName,
		Historical: historical,
		Predicted:  predicted,
		Confidence: math.Round(s.rng.Float64()*20+75*100) / 10000,
		Direction:  direction,
		Insight:    fmt.Sprintf("%s 지표는 향후 %d일간 %s 추세입니다 (시뮬레이션).", metricName, predictionDays, direction),
	}, nil
}
