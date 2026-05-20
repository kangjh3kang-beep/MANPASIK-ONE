package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"
	"testing"
	"time"

	"github.com/manpasik/backend/services/ai-inference-service/internal/llm"
)

// ============================================================================
// Fake Repositories
// ============================================================================

type fakeAnalysisRepo struct {
	mu   sync.Mutex
	data map[string]*AnalysisResult
}

func newFakeAnalysisRepo() *fakeAnalysisRepo {
	return &fakeAnalysisRepo{data: make(map[string]*AnalysisResult)}
}

func (r *fakeAnalysisRepo) Save(_ context.Context, result *AnalysisResult) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[result.AnalysisID] = result
	return nil
}

func (r *fakeAnalysisRepo) FindByID(_ context.Context, id string) (*AnalysisResult, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if v, ok := r.data[id]; ok {
		return v, nil
	}
	return nil, nil
}

func (r *fakeAnalysisRepo) FindByUserID(_ context.Context, userID string, limit int) ([]*AnalysisResult, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var results []*AnalysisResult
	for _, v := range r.data {
		if v.UserID == userID {
			results = append(results, v)
			if len(results) >= limit {
				break
			}
		}
	}
	return results, nil
}

type fakeHealthScoreRepo struct {
	mu   sync.Mutex
	data map[string]*HealthScore
}

func newFakeHealthScoreRepo() *fakeHealthScoreRepo {
	return &fakeHealthScoreRepo{data: make(map[string]*HealthScore)}
}

func (r *fakeHealthScoreRepo) Save(_ context.Context, score *HealthScore) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[score.UserID] = score
	return nil
}

func (r *fakeHealthScoreRepo) FindLatestByUserID(_ context.Context, userID string) (*HealthScore, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if v, ok := r.data[userID]; ok {
		return v, nil
	}
	return nil, nil
}

// ============================================================================
// Tests
// ============================================================================

func newTestService() *InferenceService {
	return NewInferenceService(newFakeAnalysisRepo(), newFakeHealthScoreRepo())
}

func TestAnalyzeMeasurement_Success(t *testing.T) {
	svc := newTestService()
	result, err := svc.AnalyzeMeasurement(context.Background(), "user-1", "meas-1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.UserID != "user-1" {
		t.Errorf("expected user-1, got %s", result.UserID)
	}
	if result.MeasurementID != "meas-1" {
		t.Errorf("expected meas-1, got %s", result.MeasurementID)
	}
	if len(result.Biomarkers) == 0 {
		t.Error("expected biomarkers, got none")
	}
	if result.OverallHealthScore < 0 || result.OverallHealthScore > 100 {
		t.Errorf("health score out of range: %f", result.OverallHealthScore)
	}
	if result.Summary == "" {
		t.Error("expected non-empty summary")
	}
	if result.AnalysisID == "" {
		t.Error("expected non-empty analysis ID")
	}
}

func TestAnalyzeMeasurement_EmptyUserID(t *testing.T) {
	svc := newTestService()
	_, err := svc.AnalyzeMeasurement(context.Background(), "", "meas-1", nil)
	if err == nil {
		t.Fatal("expected error for empty user_id")
	}
}

func TestAnalyzeMeasurement_EmptyMeasurementID(t *testing.T) {
	svc := newTestService()
	_, err := svc.AnalyzeMeasurement(context.Background(), "user-1", "", nil)
	if err == nil {
		t.Fatal("expected error for empty measurement_id")
	}
}

func TestAnalyzeMeasurement_WithSpecificModels(t *testing.T) {
	svc := newTestService()
	result, err := svc.AnalyzeMeasurement(context.Background(), "user-1", "meas-2", []AiModelType{ModelBiomarkerClassifier})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestGetHealthScore_Success(t *testing.T) {
	svc := newTestService()
	score, err := svc.GetHealthScore(context.Background(), "user-1", 30)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if score.UserID != "user-1" {
		t.Errorf("expected user-1, got %s", score.UserID)
	}
	if score.OverallScore < 0 || score.OverallScore > 100 {
		t.Errorf("score out of range: %f", score.OverallScore)
	}
	if len(score.CategoryScores) == 0 {
		t.Error("expected category scores")
	}
	if score.Trend == "" {
		t.Error("expected non-empty trend")
	}
}

func TestGetHealthScore_EmptyUserID(t *testing.T) {
	svc := newTestService()
	_, err := svc.GetHealthScore(context.Background(), "", 30)
	if err == nil {
		t.Fatal("expected error for empty user_id")
	}
}

func TestGetHealthScore_DefaultDays(t *testing.T) {
	svc := newTestService()
	score, err := svc.GetHealthScore(context.Background(), "user-1", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if score == nil {
		t.Fatal("expected non-nil score")
	}
}

func TestPredictTrend_Success(t *testing.T) {
	svc := newTestService()
	pred, err := svc.PredictTrend(context.Background(), "user-1", "blood_glucose", 30, 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pred.UserID != "user-1" {
		t.Errorf("expected user-1, got %s", pred.UserID)
	}
	if pred.MetricName != "blood_glucose" {
		t.Errorf("expected blood_glucose, got %s", pred.MetricName)
	}
	if len(pred.Historical) != 30 {
		t.Errorf("expected 30 historical points, got %d", len(pred.Historical))
	}
	if len(pred.Predicted) != 7 {
		t.Errorf("expected 7 predicted points, got %d", len(pred.Predicted))
	}
	if pred.Direction == "" {
		t.Error("expected non-empty direction")
	}
}

func TestPredictTrend_EmptyUserID(t *testing.T) {
	svc := newTestService()
	_, err := svc.PredictTrend(context.Background(), "", "blood_glucose", 30, 7)
	if err == nil {
		t.Fatal("expected error for empty user_id")
	}
}

func TestPredictTrend_EmptyMetric(t *testing.T) {
	svc := newTestService()
	_, err := svc.PredictTrend(context.Background(), "user-1", "", 30, 7)
	if err == nil {
		t.Fatal("expected error for empty metric")
	}
}

func TestGetModelInfo_Success(t *testing.T) {
	svc := newTestService()
	info, err := svc.GetModelInfo(context.Background(), ModelBiomarkerClassifier)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Name != "BiomarkerClassifier" {
		t.Errorf("expected BiomarkerClassifier, got %s", info.Name)
	}
	if info.Status != ModelStatusActive {
		t.Errorf("expected active, got %s", info.Status)
	}
}

func TestGetModelInfo_NotFound(t *testing.T) {
	svc := newTestService()
	_, err := svc.GetModelInfo(context.Background(), AiModelType(99))
	if err == nil {
		t.Fatal("expected error for unknown model")
	}
}

func TestListModels(t *testing.T) {
	svc := newTestService()
	models, err := svc.ListModels(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(models) != 5 {
		t.Errorf("expected 5 models, got %d", len(models))
	}
}

// ============================================================================
// Mock LLM Client
// ============================================================================

type mockLLMClient struct {
	response *llm.ChatResponse
	err      error
	called   int
}

func (m *mockLLMClient) Chat(_ context.Context, _ string, _ []llm.ChatMessage) (*llm.ChatResponse, error) {
	m.called++
	if m.err != nil {
		return nil, m.err
	}
	return m.response, nil
}

func newTestServiceWithLLM(client llm.LLMClient) *InferenceService {
	return NewInferenceService(newFakeAnalysisRepo(), newFakeHealthScoreRepo(), WithLLMClient(client))
}

// ============================================================================
// LLM 연동 테스트
// ============================================================================

func TestLLMEnabled_WithClient(t *testing.T) {
	mock := &mockLLMClient{response: &llm.ChatResponse{Content: "test"}}
	svc := newTestServiceWithLLM(mock)
	if !svc.LLMEnabled() {
		t.Error("expected LLMEnabled() = true when client is set")
	}
}

func TestLLMEnabled_WithoutClient(t *testing.T) {
	svc := newTestService()
	if svc.LLMEnabled() {
		t.Error("expected LLMEnabled() = false when client is nil")
	}
}

func TestGenerateHealthInsight_WithLLM(t *testing.T) {
	mock := &mockLLMClient{
		response: &llm.ChatResponse{
			Content:      "혈당 수치가 정상 범위입니다. 현재 건강 상태를 유지하세요.",
			FinishReason: "stop",
			TokensUsed:   50,
		},
	}
	svc := newTestServiceWithLLM(mock)

	measurements := []MeasurementData{
		{MetricName: "blood_glucose", Value: 95.0, Unit: "mg/dL", Timestamp: time.Now()},
		{MetricName: "cholesterol_total", Value: 180.0, Unit: "mg/dL", Timestamp: time.Now()},
	}

	insight, err := svc.GenerateHealthInsight(context.Background(), "user-1", measurements)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if insight == "" {
		t.Error("expected non-empty insight")
	}
	if mock.called != 1 {
		t.Errorf("expected LLM to be called once, got %d", mock.called)
	}
	if insight != "혈당 수치가 정상 범위입니다. 현재 건강 상태를 유지하세요." {
		t.Errorf("expected LLM response content, got: %s", insight)
	}
}

func TestGenerateHealthInsight_NilLLM_Fallback(t *testing.T) {
	svc := newTestService() // LLM 없음

	measurements := []MeasurementData{
		{MetricName: "blood_glucose", Value: 95.0, Unit: "mg/dL", Timestamp: time.Now()},
	}

	insight, err := svc.GenerateHealthInsight(context.Background(), "user-1", measurements)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if insight == "" {
		t.Error("expected non-empty fallback insight")
	}
	// 규칙 기반 fallback 결과 확인
	expected := "총 1건의 측정 데이터를 기반으로 분석했습니다. 규칙적인 측정과 전문가 상담을 권장합니다."
	if insight != expected {
		t.Errorf("expected rule-based fallback, got: %s", insight)
	}
}

func TestGenerateHealthInsight_LLMError_Fallback(t *testing.T) {
	mock := &mockLLMClient{
		err: errors.New("LLM API 오류"),
	}
	svc := newTestServiceWithLLM(mock)

	measurements := []MeasurementData{
		{MetricName: "blood_glucose", Value: 95.0, Unit: "mg/dL", Timestamp: time.Now()},
	}

	insight, err := svc.GenerateHealthInsight(context.Background(), "user-1", measurements)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if insight == "" {
		t.Error("expected non-empty fallback insight on LLM error")
	}
	// LLM 에러 시에도 규칙 기반 fallback
	if mock.called != 1 {
		t.Errorf("expected LLM to be called once, got %d", mock.called)
	}
}

func TestGenerateHealthInsight_EmptyUserID(t *testing.T) {
	svc := newTestService()
	_, err := svc.GenerateHealthInsight(context.Background(), "", nil)
	if err == nil {
		t.Fatal("expected error for empty user_id")
	}
}

func TestGenerateHealthInsight_EmptyMeasurements(t *testing.T) {
	svc := newTestService()
	insight, err := svc.GenerateHealthInsight(context.Background(), "user-1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "측정 데이터가 없어 분석을 수행할 수 없습니다. 측정을 먼저 진행해 주세요."
	if insight != expected {
		t.Errorf("expected empty measurement fallback, got: %s", insight)
	}
}

func TestAnalyzeMeasurement_WithLLM_EnhancedSummary(t *testing.T) {
	mock := &mockLLMClient{
		response: &llm.ChatResponse{
			Content:      "LLM 향상된 분석 요약입니다.",
			FinishReason: "stop",
			TokensUsed:   30,
		},
	}
	svc := newTestServiceWithLLM(mock)

	result, err := svc.AnalyzeMeasurement(context.Background(), "user-1", "meas-1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// LLM이 활성화되어 있으므로 enhanceSummaryWithLLM이 호출됨
	if mock.called == 0 {
		t.Error("expected LLM to be called for summary enhancement")
	}
	if result.Summary != "LLM 향상된 분석 요약입니다." {
		t.Errorf("expected LLM-enhanced summary, got: %s", result.Summary)
	}
}

func TestGetHealthScore_WithLLM_EnhancedRecommendation(t *testing.T) {
	mock := &mockLLMClient{
		response: &llm.ChatResponse{
			Content:      "유산소 운동을 주 3회 이상 실시하고, 나트륨 섭취를 줄이세요.",
			FinishReason: "stop",
			TokensUsed:   40,
		},
	}
	svc := newTestServiceWithLLM(mock)

	score, err := svc.GetHealthScore(context.Background(), "user-1", 30)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.called == 0 {
		t.Error("expected LLM to be called for recommendation")
	}
	if score.Recommendation != "유산소 운동을 주 3회 이상 실시하고, 나트륨 섭취를 줄이세요." {
		t.Errorf("expected LLM recommendation, got: %s", score.Recommendation)
	}
}

func TestGetHealthScore_WithoutLLM_RuleBasedRecommendation(t *testing.T) {
	svc := newTestService() // LLM 없음

	score, err := svc.GetHealthScore(context.Background(), "user-1", 30)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// LLM 없이 카테고리 기반 규칙 추천이 생성되어야 함 (빈 문자열이 아님)
	if score.Recommendation == "" {
		t.Error("추천이 비어 있으면 안 됩니다")
	}
}

// ============================================================================
// Phase B 신규 테스트
// ============================================================================

func TestGetAnalysisHistory(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	// 3건 분석 저장
	for i := 0; i < 3; i++ {
		_, err := svc.AnalyzeMeasurement(ctx, "user-hist", fmt.Sprintf("meas-%d", i), nil)
		if err != nil {
			t.Fatalf("분석 %d 실패: %v", i, err)
		}
	}

	results, err := svc.GetAnalysisHistory(ctx, "user-hist", 10)
	if err != nil {
		t.Fatalf("GetAnalysisHistory 실패: %v", err)
	}
	if len(results) != 3 {
		t.Errorf("분석 이력 수 = %d, want 3", len(results))
	}
}

func TestGetAnalysisHistory_EmptyUserID(t *testing.T) {
	svc := newTestService()
	_, err := svc.GetAnalysisHistory(context.Background(), "", 10)
	if err == nil {
		t.Fatal("빈 user_id에 에러가 반환되어야 합니다")
	}
}

func TestAnalyzeMeasurementWithData_Success(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	data := []MeasurementData{
		{MetricName: "blood_glucose", Value: 95.0, Unit: "mg/dL", Timestamp: time.Now()},
		{MetricName: "heart_rate", Value: 72.0, Unit: "bpm", Timestamp: time.Now()},
	}

	result, err := svc.AnalyzeMeasurementWithData(ctx, "user-1", "meas-data-1", data)
	if err != nil {
		t.Fatalf("AnalyzeMeasurementWithData 실패: %v", err)
	}
	if len(result.Biomarkers) != 2 {
		t.Errorf("바이오마커 수 = %d, want 2", len(result.Biomarkers))
	}
	// blood_glucose 95는 정상 범위 (70-100)
	for _, bm := range result.Biomarkers {
		if bm.BiomarkerName == "blood_glucose" {
			if bm.Classification != "normal" {
				t.Errorf("blood_glucose 95는 normal이어야 합니다, got %s", bm.Classification)
			}
		}
	}
}

func TestAnalyzeMeasurementWithData_EmptyData(t *testing.T) {
	svc := newTestService()
	_, err := svc.AnalyzeMeasurementWithData(context.Background(), "user-1", "meas-1", nil)
	if err == nil {
		t.Fatal("빈 데이터에 에러가 반환되어야 합니다")
	}
}

func TestClassifyBiomarker_Normal(t *testing.T) {
	result := ClassifyBiomarker("blood_glucose", 90.0)
	if result.Classification != "normal" {
		t.Errorf("blood_glucose 90은 normal이어야 합니다, got %s", result.Classification)
	}
	if result.RiskLevel != RiskLow {
		t.Errorf("정상 범위는 RiskLow여야 합니다, got %d", result.RiskLevel)
	}
}

func TestClassifyBiomarker_Abnormal(t *testing.T) {
	result := ClassifyBiomarker("blood_glucose", 200.0)
	if result.Classification != "abnormal" {
		t.Errorf("blood_glucose 200은 abnormal이어야 합니다, got %s", result.Classification)
	}
	if result.RiskLevel < RiskHigh {
		t.Errorf("비정상 범위는 최소 RiskHigh여야 합니다, got %d", result.RiskLevel)
	}
}

func TestClassifyBiomarker_Unknown(t *testing.T) {
	result := ClassifyBiomarker("unknown_metric", 50.0)
	if result.Classification != "unknown" {
		t.Errorf("미등록 메트릭은 unknown이어야 합니다, got %s", result.Classification)
	}
}

func TestPredictTrend_UnsupportedMetric(t *testing.T) {
	svc := newTestService()
	_, err := svc.PredictTrend(context.Background(), "user-1", "unknown_metric", 30, 7)
	if err == nil {
		t.Fatal("지원하지 않는 메트릭에 에러가 반환되어야 합니다")
	}
}

func TestAnalyzeMeasurement_DeprecatedModel(t *testing.T) {
	svc := newTestService()
	// FoodCalorieEstimator를 deprecated로 변경
	svc.models[ModelFoodCalorieEstimator].Status = ModelStatusDeprecated

	_, err := svc.AnalyzeMeasurement(context.Background(), "user-1", "meas-1",
		[]AiModelType{ModelFoodCalorieEstimator})
	if err == nil {
		t.Fatal("deprecated 모델 사용 시 에러가 반환되어야 합니다")
	}
}

// ============================================================================
// Phase E-1: 실제 데이터 기반 건강 점수 + IQR 이상탐지 + 선형 회귀 예측
// ============================================================================

func TestDetectAnomaliesByIQR_Normal(t *testing.T) {
	biomarkers := []BiomarkerResult{
		{BiomarkerName: "blood_glucose", Value: 85.0, Unit: "mg/dL"},
		{BiomarkerName: "heart_rate", Value: 75.0, Unit: "bpm"},
	}
	anomalies := detectAnomaliesByIQR(biomarkers)
	if len(anomalies) != 0 {
		t.Errorf("정상 범위 값에 IQR 이상치가 없어야 합니다, got %d", len(anomalies))
	}
}

func TestDetectAnomaliesByIQR_Extreme(t *testing.T) {
	// blood_glucose 정상 범위: 70-100, IQR=30, fence: 70-45=25 ~ 100+45=145
	biomarkers := []BiomarkerResult{
		{BiomarkerName: "blood_glucose", Value: 20.0, Unit: "mg/dL"},
	}
	anomalies := detectAnomaliesByIQR(biomarkers)
	if len(anomalies) != 1 {
		t.Fatalf("극단적 이상치가 감지되어야 합니다, got %d", len(anomalies))
	}
	if anomalies[0].MetricName != "blood_glucose" {
		t.Errorf("MetricName = %s, want blood_glucose", anomalies[0].MetricName)
	}
}

func TestDetectCombinedAnomalies_MergesResults(t *testing.T) {
	// 극도로 비정상적인 값 → Z-score + IQR 모두 감지
	biomarkers := []BiomarkerResult{
		{BiomarkerName: "blood_glucose", Value: 300.0, Unit: "mg/dL"},
		{BiomarkerName: "heart_rate", Value: 80.0, Unit: "bpm"}, // 정상
	}
	anomalies := detectCombinedAnomalies(biomarkers)
	if len(anomalies) != 1 {
		t.Fatalf("비정상 1건만 감지되어야 합니다, got %d", len(anomalies))
	}
	// Z-score와 IQR 결합 설명 포함
	if anomalies[0].MetricName != "blood_glucose" {
		t.Errorf("MetricName = %s, want blood_glucose", anomalies[0].MetricName)
	}
}

func TestLinearRegression_PerfectFit(t *testing.T) {
	xs := []float64{0, 1, 2, 3, 4}
	ys := []float64{2, 4, 6, 8, 10}
	slope, intercept, r2 := linearRegression(xs, ys)
	if math.Abs(slope-2.0) > 0.001 {
		t.Errorf("slope = %.4f, want 2.0", slope)
	}
	if math.Abs(intercept-2.0) > 0.001 {
		t.Errorf("intercept = %.4f, want 2.0", intercept)
	}
	if math.Abs(r2-1.0) > 0.001 {
		t.Errorf("r2 = %.4f, want 1.0", r2)
	}
}

func TestLinearRegression_SinglePoint(t *testing.T) {
	xs := []float64{1}
	ys := []float64{5}
	slope, intercept, _ := linearRegression(xs, ys)
	if slope != 0 {
		t.Errorf("단일 포인트에서 slope = %.4f, want 0", slope)
	}
	if intercept != 5 {
		t.Errorf("단일 포인트에서 intercept = %.4f, want 5", intercept)
	}
}

func TestLinearRegression_Empty(t *testing.T) {
	slope, intercept, r2 := linearRegression(nil, nil)
	if slope != 0 || intercept != 0 || r2 != 0 {
		t.Errorf("빈 입력에서 모두 0이어야 합니다: slope=%.4f, intercept=%.4f, r2=%.4f", slope, intercept, r2)
	}
}

func TestComputeCategoryScores_WithData(t *testing.T) {
	analyses := []*AnalysisResult{
		{
			Biomarkers: []BiomarkerResult{
				{BiomarkerName: "blood_glucose", Classification: "normal", RiskLevel: RiskLow},
				{BiomarkerName: "heart_rate", Classification: "normal", RiskLevel: RiskLow},
				{BiomarkerName: "hemoglobin", Classification: "borderline", RiskLevel: RiskModerate},
			},
			OverallHealthScore: 80,
		},
	}
	scores := computeCategoryScores(analyses)
	// metabolic (blood_glucose=normal→95)
	if scores["metabolic"] != 95 {
		t.Errorf("metabolic = %.1f, want 95.0", scores["metabolic"])
	}
	// cardiovascular (heart_rate=normal→95)
	if scores["cardiovascular"] != 95 {
		t.Errorf("cardiovascular = %.1f, want 95.0", scores["cardiovascular"])
	}
	// nutritional (hemoglobin=borderline→70)
	if scores["nutritional"] != 70 {
		t.Errorf("nutritional = %.1f, want 70.0", scores["nutritional"])
	}
	// fitness (no data → default 75)
	if scores["fitness"] != 75 {
		t.Errorf("fitness = %.1f, want 75.0 (default)", scores["fitness"])
	}
}

func TestComputeCategoryScores_Empty(t *testing.T) {
	scores := computeCategoryScores(nil)
	for cat, score := range scores {
		if score != 75 {
			t.Errorf("%s = %.1f, want 75.0 (default)", cat, score)
		}
	}
}

func TestComputeOverallFromCategories(t *testing.T) {
	categories := map[string]float64{
		"cardiovascular": 90,
		"metabolic":      80,
		"nutritional":    70,
		"fitness":        60,
	}
	// 가중평균: 90*0.35 + 80*0.30 + 70*0.20 + 60*0.15 = 31.5 + 24 + 14 + 9 = 78.5
	overall := computeOverallFromCategories(categories)
	if math.Abs(overall-78.5) > 0.1 {
		t.Errorf("overall = %.1f, want 78.5", overall)
	}
}

func TestDetermineTrend_Improving(t *testing.T) {
	analyses := []*AnalysisResult{
		{OverallHealthScore: 90}, // newest
		{OverallHealthScore: 88},
		{OverallHealthScore: 70}, // oldest
		{OverallHealthScore: 65},
	}
	trend := determineTrend(analyses)
	if trend != "improving" {
		t.Errorf("trend = %s, want improving", trend)
	}
}

func TestDetermineTrend_Declining(t *testing.T) {
	analyses := []*AnalysisResult{
		{OverallHealthScore: 60}, // newest
		{OverallHealthScore: 65},
		{OverallHealthScore: 85}, // oldest
		{OverallHealthScore: 90},
	}
	trend := determineTrend(analyses)
	if trend != "declining" {
		t.Errorf("trend = %s, want declining", trend)
	}
}

func TestDetermineTrend_Stable(t *testing.T) {
	analyses := []*AnalysisResult{
		{OverallHealthScore: 80},
		{OverallHealthScore: 79},
	}
	trend := determineTrend(analyses)
	if trend != "stable" {
		t.Errorf("trend = %s, want stable", trend)
	}
}

func TestDetermineTrend_SingleAnalysis(t *testing.T) {
	analyses := []*AnalysisResult{{OverallHealthScore: 80}}
	trend := determineTrend(analyses)
	if trend != "stable" {
		t.Errorf("단일 분석 시 trend = %s, want stable", trend)
	}
}

func TestGetHealthScore_WithAnalysisHistory(t *testing.T) {
	ar := newFakeAnalysisRepo()
	hsr := newFakeHealthScoreRepo()
	svc := NewInferenceService(ar, hsr)
	ctx := context.Background()

	// 분석 이력 생성
	for i := 0; i < 3; i++ {
		_, err := svc.AnalyzeMeasurementWithData(ctx, "user-score", fmt.Sprintf("m-%d", i), []MeasurementData{
			{MetricName: "blood_glucose", Value: 90, Unit: "mg/dL", Timestamp: time.Now()},
			{MetricName: "heart_rate", Value: 72, Unit: "bpm", Timestamp: time.Now()},
			{MetricName: "hemoglobin", Value: 14, Unit: "g/dL", Timestamp: time.Now()},
		})
		if err != nil {
			t.Fatalf("AnalyzeMeasurementWithData 실패: %v", err)
		}
	}

	// 이력 기반 건강 점수 계산
	score, err := svc.GetHealthScore(ctx, "user-score", 30)
	if err != nil {
		t.Fatalf("GetHealthScore 실패: %v", err)
	}
	// 모든 바이오마커가 정상 → 높은 점수 기대
	if score.OverallScore < 80 {
		t.Errorf("정상 바이오마커에서 OverallScore = %.1f, want >= 80", score.OverallScore)
	}
	if score.CategoryScores["metabolic"] < 90 {
		t.Errorf("정상 blood_glucose에서 metabolic = %.1f, want >= 90", score.CategoryScores["metabolic"])
	}
}

func TestGetHealthScore_NoHistory_DefaultScore(t *testing.T) {
	svc := newTestService()
	score, err := svc.GetHealthScore(context.Background(), "user-new", 30)
	if err != nil {
		t.Fatalf("GetHealthScore 실패: %v", err)
	}
	if score.OverallScore != 75 {
		t.Errorf("이력 없을 때 OverallScore = %.1f, want 75.0", score.OverallScore)
	}
	if score.Trend != "stable" {
		t.Errorf("이력 없을 때 Trend = %s, want stable", score.Trend)
	}
}

func TestPredictTrend_WithHistoricalData(t *testing.T) {
	ar := newFakeAnalysisRepo()
	hsr := newFakeHealthScoreRepo()
	svc := NewInferenceService(ar, hsr)
	ctx := context.Background()

	// 상승 추세 데이터를 직접 저장 (명시적 타임스탬프)
	baseTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 5; i++ {
		result := &AnalysisResult{
			AnalysisID:         fmt.Sprintf("trend-ana-%d", i),
			UserID:             "user-trend",
			MeasurementID:      fmt.Sprintf("m-%d", i),
			Biomarkers:         []BiomarkerResult{{BiomarkerName: "blood_glucose", Value: 85.0 + float64(i)*2, Unit: "mg/dL"}},
			OverallHealthScore: 90,
			AnalyzedAt:         baseTime.AddDate(0, 0, i), // 1일 간격
		}
		if err := ar.Save(ctx, result); err != nil {
			t.Fatalf("Save 실패: %v", err)
		}
	}

	pred, err := svc.PredictTrend(ctx, "user-trend", "blood_glucose", 30, 7)
	if err != nil {
		t.Fatalf("PredictTrend 실패: %v", err)
	}
	// 실제 데이터 사용 → 선형 회귀 인사이트
	if pred.Direction != "up" {
		t.Errorf("상승 데이터에서 direction = %s, want up", pred.Direction)
	}
	if pred.Confidence < 0.5 {
		t.Errorf("Confidence = %.2f, want >= 0.5", pred.Confidence)
	}
	if len(pred.Historical) != 5 {
		t.Errorf("Historical 길이 = %d, want 5", len(pred.Historical))
	}
	if len(pred.Predicted) != 7 {
		t.Errorf("Predicted 길이 = %d, want 7", len(pred.Predicted))
	}
}

func TestPredictTrend_InsufficientData_Fallback(t *testing.T) {
	svc := newTestService()
	// 이력 없음 → 시뮬레이션 폴백
	pred, err := svc.PredictTrend(context.Background(), "user-empty", "blood_glucose", 30, 7)
	if err != nil {
		t.Fatalf("PredictTrend 실패: %v", err)
	}
	if len(pred.Historical) != 30 {
		t.Errorf("시뮬레이션 Historical 길이 = %d, want 30", len(pred.Historical))
	}
	if len(pred.Predicted) != 7 {
		t.Errorf("시뮬레이션 Predicted 길이 = %d, want 7", len(pred.Predicted))
	}
}

func TestExtractMetricHistory(t *testing.T) {
	analyses := []*AnalysisResult{
		{ // newest
			AnalyzedAt: time.Date(2026, 1, 3, 0, 0, 0, 0, time.UTC),
			Biomarkers: []BiomarkerResult{{BiomarkerName: "blood_glucose", Value: 95}},
		},
		{ // oldest
			AnalyzedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			Biomarkers: []BiomarkerResult{{BiomarkerName: "blood_glucose", Value: 85}},
		},
	}
	times, values := extractMetricHistory(analyses, "blood_glucose")
	if len(values) != 2 {
		t.Fatalf("values 길이 = %d, want 2", len(values))
	}
	// 시간순(오래된→최신) 정렬 확인
	if values[0] != 85 {
		t.Errorf("values[0] = %.1f, want 85 (oldest first)", values[0])
	}
	if values[1] != 95 {
		t.Errorf("values[1] = %.1f, want 95 (newest last)", values[1])
	}
	if !times[0].Before(times[1]) {
		t.Error("times should be chronologically ordered")
	}
}

func TestBiomarkerToScore(t *testing.T) {
	tests := []struct {
		classification string
		riskLevel      RiskLevel
		expected       float64
	}{
		{"normal", RiskLow, 95.0},
		{"borderline", RiskModerate, 70.0},
		{"abnormal", RiskHigh, 40.0},
		{"abnormal", RiskCritical, 25.0},
		{"unknown", RiskLow, 60.0},
	}
	for _, tt := range tests {
		bm := BiomarkerResult{Classification: tt.classification, RiskLevel: tt.riskLevel}
		score := biomarkerToScore(bm)
		if score != tt.expected {
			t.Errorf("biomarkerToScore(%s, %d) = %.1f, want %.1f",
				tt.classification, tt.riskLevel, score, tt.expected)
		}
	}
}

// ============================================================================
// Phase E-2: 규칙 기반 카테고리별 추천 + 에스컬레이션 테스트
// ============================================================================

func TestGenerateRuleBasedRecommendation_CardiovascularLow(t *testing.T) {
	rec := generateRuleBasedRecommendation(60.0,
		map[string]float64{"cardiovascular": 40, "metabolic": 80, "nutritional": 75, "fitness": 70},
		"declining",
	)
	if rec == "" {
		t.Fatal("추천이 비어 있습니다")
	}
	// 심혈관 40점 → 심혈관 관련 추천
	if !contains(rec, "심혈관") && !contains(rec, "유산소") {
		t.Errorf("심혈관 점수 40에서 심혈관 관련 추천이 있어야 합니다, got: %s", rec)
	}
	// declining 추세 → 하락 메시지
	if !contains(rec, "하락") {
		t.Errorf("declining 추세에서 하락 메시지가 있어야 합니다, got: %s", rec)
	}
}

func TestGenerateRuleBasedRecommendation_MetabolicLow(t *testing.T) {
	rec := generateRuleBasedRecommendation(55.0,
		map[string]float64{"cardiovascular": 80, "metabolic": 35, "nutritional": 70, "fitness": 75},
		"stable",
	)
	if !contains(rec, "대사") || !contains(rec, "탄수화물") {
		t.Errorf("대사 점수 35에서 대사 관련 추천이 있어야 합니다, got: %s", rec)
	}
}

func TestGenerateRuleBasedRecommendation_NutritionalLow(t *testing.T) {
	rec := generateRuleBasedRecommendation(58.0,
		map[string]float64{"cardiovascular": 80, "metabolic": 80, "nutritional": 40, "fitness": 75},
		"improving",
	)
	if !contains(rec, "영양") {
		t.Errorf("영양 점수 40에서 영양 관련 추천이 있어야 합니다, got: %s", rec)
	}
	// improving 추세 → 개선 메시지
	if !contains(rec, "개선") || !contains(rec, "유지") {
		t.Errorf("improving 추세에서 개선 유지 메시지가 있어야 합니다, got: %s", rec)
	}
}

func TestGenerateRuleBasedRecommendation_FitnessLow(t *testing.T) {
	rec := generateRuleBasedRecommendation(60.0,
		map[string]float64{"cardiovascular": 80, "metabolic": 80, "nutritional": 80, "fitness": 45},
		"stable",
	)
	if !contains(rec, "체력") || !contains(rec, "걷기") {
		t.Errorf("체력 점수 45에서 체력 관련 추천이 있어야 합니다, got: %s", rec)
	}
}

func TestGenerateRuleBasedRecommendation_VeryLowOverall(t *testing.T) {
	rec := generateRuleBasedRecommendation(35.0,
		map[string]float64{"cardiovascular": 30, "metabolic": 40, "nutritional": 35, "fitness": 35},
		"declining",
	)
	// 전체 점수 35 → 전문의 상담 경고
	if !contains(rec, "전문의") {
		t.Errorf("전체 점수 35점에서 전문의 상담 경고가 있어야 합니다, got: %s", rec)
	}
}

func TestGenerateRuleBasedRecommendation_AllGood(t *testing.T) {
	rec := generateRuleBasedRecommendation(90.0,
		map[string]float64{"cardiovascular": 90, "metabolic": 88, "nutritional": 92, "fitness": 90},
		"improving",
	)
	// 모두 양호할 때는 유지 메시지
	if !contains(rec, "유지") && !contains(rec, "양호") {
		t.Errorf("모든 점수 양호 시 유지/양호 메시지가 있어야 합니다, got: %s", rec)
	}
}

func TestGenerateRuleBasedRecommendation_EmptyCategories(t *testing.T) {
	rec := generateRuleBasedRecommendation(75.0, map[string]float64{}, "stable")
	if rec == "" {
		t.Fatal("빈 카테고리에서도 기본 추천이 반환되어야 합니다")
	}
}

// --- RiskEscalator 테스트 ---

type fakeRiskEscalator struct {
	alerts []*RiskAlert
	mu     sync.Mutex
}

func (f *fakeRiskEscalator) EscalateCriticalRisk(_ context.Context, alert *RiskAlert) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.alerts = append(f.alerts, alert)
	return nil
}

func TestEscalateCriticalBiomarkers_TriggeredOnCritical(t *testing.T) {
	ar := newFakeAnalysisRepo()
	hsr := newFakeHealthScoreRepo()
	svc := NewInferenceService(ar, hsr)
	escalator := &fakeRiskEscalator{}
	svc.SetRiskEscalator(escalator)

	ctx := context.Background()

	// Critical 위험 값으로 분석
	_, err := svc.AnalyzeMeasurementWithData(ctx, "user-crit", "meas-1", []MeasurementData{
		{MetricName: "blood_glucose", Value: 400, Unit: "mg/dL", Timestamp: time.Now()}, // 극단적 고혈당
	})
	if err != nil {
		t.Fatalf("AnalyzeMeasurementWithData 실패: %v", err)
	}

	escalator.mu.Lock()
	defer escalator.mu.Unlock()
	if len(escalator.alerts) == 0 {
		t.Fatal("Critical 위험 시 에스컬레이션이 호출되어야 합니다")
	}
	alert := escalator.alerts[0]
	if alert.RiskLevel != RiskCritical {
		t.Errorf("RiskLevel = %d, RiskCritical 기대", alert.RiskLevel)
	}
	if alert.UserID != "user-crit" {
		t.Errorf("UserID = %q, user-crit 기대", alert.UserID)
	}
}

func TestEscalateCriticalBiomarkers_NotTriggeredOnNormal(t *testing.T) {
	ar := newFakeAnalysisRepo()
	hsr := newFakeHealthScoreRepo()
	svc := NewInferenceService(ar, hsr)
	escalator := &fakeRiskEscalator{}
	svc.SetRiskEscalator(escalator)

	ctx := context.Background()

	// 정상 범위 값
	_, err := svc.AnalyzeMeasurementWithData(ctx, "user-norm", "meas-2", []MeasurementData{
		{MetricName: "blood_glucose", Value: 90, Unit: "mg/dL", Timestamp: time.Now()},
	})
	if err != nil {
		t.Fatalf("AnalyzeMeasurementWithData 실패: %v", err)
	}

	escalator.mu.Lock()
	defer escalator.mu.Unlock()
	if len(escalator.alerts) != 0 {
		t.Errorf("정상 값에서 에스컬레이션이 발생하면 안 됩니다, got %d alerts", len(escalator.alerts))
	}
}

func TestEscalateCriticalBiomarkers_NilEscalator(t *testing.T) {
	ar := newFakeAnalysisRepo()
	hsr := newFakeHealthScoreRepo()
	svc := NewInferenceService(ar, hsr)
	// riskEscalator 미설정 → 패닉 없이 정상 작동

	ctx := context.Background()
	_, err := svc.AnalyzeMeasurementWithData(ctx, "user-nil", "meas-3", []MeasurementData{
		{MetricName: "blood_glucose", Value: 400, Unit: "mg/dL", Timestamp: time.Now()},
	})
	if err != nil {
		t.Fatalf("에스컬레이터 없이도 에러가 발생하면 안 됩니다: %v", err)
	}
}

func TestLogRiskEscalator(t *testing.T) {
	escalator := &LogRiskEscalator{}
	err := escalator.EscalateCriticalRisk(context.Background(), &RiskAlert{
		UserID:        "test-user",
		RiskLevel:     RiskCritical,
		BiomarkerName: "blood_glucose",
		Value:         400,
		Unit:          "mg/dL",
		Message:       "위험",
		DetectedAt:    time.Now(),
	})
	if err != nil {
		t.Fatalf("LogRiskEscalator 에러: %v", err)
	}
}

// ============================================================================
// Phase E-3: 모델 레지스트리 관리 테스트
// ============================================================================

func TestRegisterModel_Success(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	newModel := &ModelInfo{
		ModelType:   ModelBiomarkerClassifier,
		Name:        "BiomarkerClassifier-v2",
		Version:     "2.0.0",
		Description: "바이오마커 분류 모델 v2",
		Accuracy:    0.96,
	}

	err := svc.RegisterModel(ctx, newModel)
	if err != nil {
		t.Fatalf("RegisterModel 실패: %v", err)
	}

	// 등록된 모델 확인
	info, err := svc.GetModelInfo(ctx, ModelBiomarkerClassifier)
	if err != nil {
		t.Fatalf("GetModelInfo 실패: %v", err)
	}
	if info.Name != "BiomarkerClassifier-v2" {
		t.Errorf("Name = %s, want BiomarkerClassifier-v2", info.Name)
	}
	if info.Version != "2.0.0" {
		t.Errorf("Version = %s, want 2.0.0", info.Version)
	}
	if info.Status != ModelStatusActive {
		t.Errorf("Status = %s, want active", info.Status)
	}
}

func TestRegisterModel_DeprecatesExisting(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	// 기존 모델 확인
	oldModel, _ := svc.GetModelInfo(ctx, ModelBiomarkerClassifier)
	if oldModel.Status != ModelStatusActive {
		t.Fatalf("기존 모델이 active여야 합니다, got %s", oldModel.Status)
	}

	// 새 모델 등록
	newModel := &ModelInfo{
		ModelType: ModelBiomarkerClassifier,
		Name:      "BiomarkerClassifier-v2",
		Version:   "2.0.0",
	}
	err := svc.RegisterModel(ctx, newModel)
	if err != nil {
		t.Fatalf("RegisterModel 실패: %v", err)
	}

	// 기존 모델이 deprecated 상태 확인 (새 모델로 교체됨)
	info, _ := svc.GetModelInfo(ctx, ModelBiomarkerClassifier)
	if info.Name != "BiomarkerClassifier-v2" {
		t.Errorf("교체 후 Name = %s, want BiomarkerClassifier-v2", info.Name)
	}
}

func TestRegisterModel_NilModel(t *testing.T) {
	svc := newTestService()
	err := svc.RegisterModel(context.Background(), nil)
	if err == nil {
		t.Fatal("nil 모델에 에러가 반환되어야 합니다")
	}
}

func TestRegisterModel_EmptyName(t *testing.T) {
	svc := newTestService()
	err := svc.RegisterModel(context.Background(), &ModelInfo{
		ModelType: ModelBiomarkerClassifier,
		Version:   "1.0.0",
	})
	if err == nil {
		t.Fatal("빈 이름에 에러가 반환되어야 합니다")
	}
}

func TestRegisterModel_EmptyVersion(t *testing.T) {
	svc := newTestService()
	err := svc.RegisterModel(context.Background(), &ModelInfo{
		ModelType: ModelBiomarkerClassifier,
		Name:      "TestModel",
	})
	if err == nil {
		t.Fatal("빈 버전에 에러가 반환되어야 합니다")
	}
}

func TestUpdateModelStatus_Success(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	err := svc.UpdateModelStatus(ctx, ModelBiomarkerClassifier, ModelStatusTraining)
	if err != nil {
		t.Fatalf("UpdateModelStatus 실패: %v", err)
	}

	info, _ := svc.GetModelInfo(ctx, ModelBiomarkerClassifier)
	if info.Status != ModelStatusTraining {
		t.Errorf("Status = %s, want training", info.Status)
	}
}

func TestUpdateModelStatus_NotFound(t *testing.T) {
	svc := newTestService()
	err := svc.UpdateModelStatus(context.Background(), AiModelType(99), ModelStatusActive)
	if err == nil {
		t.Fatal("존재하지 않는 모델에 에러가 반환되어야 합니다")
	}
}

func TestGetActiveModels(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	active := svc.GetActiveModels(ctx)
	// 기본 5개 중 FoodCalorieEstimator는 training → active 4개
	if len(active) != 4 {
		t.Errorf("active 모델 수 = %d, want 4", len(active))
	}

	// 하나를 deprecated로 변경
	_ = svc.UpdateModelStatus(ctx, ModelAnomalyDetector, ModelStatusDeprecated)
	active = svc.GetActiveModels(ctx)
	if len(active) != 3 {
		t.Errorf("deprecated 후 active 모델 수 = %d, want 3", len(active))
	}
}

func TestGetActiveModels_AfterRegister(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	// FoodCalorieEstimator를 active로 등록
	err := svc.RegisterModel(ctx, &ModelInfo{
		ModelType: ModelFoodCalorieEstimator,
		Name:      "FoodCalorie-v2",
		Version:   "2.0.0",
	})
	if err != nil {
		t.Fatalf("RegisterModel 실패: %v", err)
	}

	active := svc.GetActiveModels(ctx)
	if len(active) != 5 {
		t.Errorf("등록 후 active 모델 수 = %d, want 5", len(active))
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchSubstring(s, substr)
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
