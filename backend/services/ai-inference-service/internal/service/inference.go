package service

import (
	"context"
	"fmt"
	"math"
	"math/rand"
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
type InferenceService struct {
	analysisRepo    AnalysisRepository
	healthScoreRepo HealthScoreRepository
	llmClient       llm.LLMClient // nil이면 LLM 미사용
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

// LLMEnabled는 LLM 클라이언트가 설정되어 있는지 반환합니다.
func (s *InferenceService) LLMEnabled() bool {
	return s.llmClient != nil
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

	// 시뮬레이션된 AI 추론 결과 생성
	biomarkers := s.simulateBiomarkerAnalysis()
	anomalies := s.simulateAnomalyDetection()
	healthScore := s.simulateHealthScore()

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
	return result, nil
}

// GetHealthScore calculates a user's health score based on recent data.
func (s *InferenceService) GetHealthScore(ctx context.Context, userID string, days int) (*HealthScore, error) {
	if userID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "user_id is required")
	}
	if days <= 0 {
		days = 30
	}

	overallScore := math.Round(60+s.rng.Float64()*40*10) / 10 // 60~100
	categoryScores := map[string]float64{
		"cardiovascular": math.Round((70+s.rng.Float64()*30)*10) / 10,
		"metabolic":      math.Round((65+s.rng.Float64()*35)*10) / 10,
		"nutritional":    math.Round((60+s.rng.Float64()*40)*10) / 10,
		"fitness":        math.Round((55+s.rng.Float64()*45)*10) / 10,
	}
	trend := s.pickTrend()

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
func (s *InferenceService) PredictTrend(ctx context.Context, userID, metricName string, historyDays, predictionDays int) (*TrendPrediction, error) {
	if userID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "user_id is required")
	}
	if metricName == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "metric_name is required")
	}
	if historyDays <= 0 {
		historyDays = 30
	}
	if predictionDays <= 0 {
		predictionDays = 7
	}

	now := time.Now()
	baseValue := 80 + s.rng.Float64()*40 // 80~120 기본값

	// 시뮬레이션된 과거 데이터
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

	// 시뮬레이션된 예측 데이터
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
		Confidence: math.Round(s.rng.Float64()*20+75*100) / 10000, // 0.75~0.95
		Direction:  direction,
		Insight:    fmt.Sprintf("%s 지표는 향후 %d일간 %s 추세입니다.", metricName, predictionDays, direction),
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

// ============================================================================
// Internal helpers — simulated AI inference
// ============================================================================

func (s *InferenceService) simulateBiomarkerAnalysis() []BiomarkerResult {
	biomarkers := []struct {
		name string
		unit string
		min  float64
		max  float64
		ref  string
	}{
		{"blood_glucose", "mg/dL", 70, 140, "70-100 mg/dL"},
		{"cholesterol_total", "mg/dL", 120, 280, "< 200 mg/dL"},
		{"hemoglobin", "g/dL", 10, 18, "12-16 g/dL (여), 14-18 g/dL (남)"},
		{"creatinine", "mg/dL", 0.5, 2.0, "0.7-1.3 mg/dL"},
		{"uric_acid", "mg/dL", 2, 10, "3.5-7.2 mg/dL"},
	}

	results := make([]BiomarkerResult, len(biomarkers))
	for i, bm := range biomarkers {
		val := bm.min + s.rng.Float64()*(bm.max-bm.min)
		val = math.Round(val*10) / 10

		classification := "normal"
		risk := RiskLow
		if val > bm.max*0.85 {
			classification = "abnormal"
			risk = RiskHigh
		} else if val > bm.max*0.7 {
			classification = "borderline"
			risk = RiskModerate
		}

		results[i] = BiomarkerResult{
			BiomarkerName:  bm.name,
			Value:          val,
			Unit:           bm.unit,
			Classification: classification,
			Confidence:     math.Round((0.85+s.rng.Float64()*0.15)*1000) / 1000,
			RiskLevel:      risk,
			ReferenceRange: bm.ref,
		}
	}
	return results
}

func (s *InferenceService) simulateAnomalyDetection() []AnomalyFlag {
	// 20% 확률로 이상치 발생
	if s.rng.Float64() > 0.2 {
		return nil
	}
	return []AnomalyFlag{
		{
			MetricName:   "heart_rate_variability",
			Value:        30 + s.rng.Float64()*20,
			ExpectedMin:  50,
			ExpectedMax:  100,
			AnomalyScore: 0.7 + s.rng.Float64()*0.3,
			Description:  "심박변이도(HRV)가 정상 범위보다 낮습니다.",
		},
	}
}

func (s *InferenceService) simulateHealthScore() float64 {
	return math.Round((60+s.rng.Float64()*40)*10) / 10
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

// generateRecommendation은 LLM을 사용하여 맞춤형 건강 추천을 생성합니다.
// LLM이 비활성화되어 있거나 호출 실패 시 기본 추천을 반환합니다.
func (s *InferenceService) generateRecommendation(ctx context.Context, overallScore float64, categoryScores map[string]float64, trend string) string {
	defaultRec := "규칙적인 운동과 균형 잡힌 식단을 유지하세요."

	if s.llmClient == nil {
		return defaultRec
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
		return defaultRec
	}

	return resp.Content
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
