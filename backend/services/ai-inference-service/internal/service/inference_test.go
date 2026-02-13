package service

import (
	"context"
	"errors"
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

func TestGetHealthScore_WithoutLLM_DefaultRecommendation(t *testing.T) {
	svc := newTestService() // LLM 없음

	score, err := svc.GetHealthScore(context.Background(), "user-1", 30)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "규칙적인 운동과 균형 잡힌 식단을 유지하세요."
	if score.Recommendation != expected {
		t.Errorf("expected default recommendation, got: %s", score.Recommendation)
	}
}
