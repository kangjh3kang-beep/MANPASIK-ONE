package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/manpasik/backend/services/nlp-service/internal/classifier"
	"github.com/manpasik/backend/services/nlp-service/internal/repository/memory"
	"github.com/manpasik/backend/services/nlp-service/internal/service"
)

// newTestService는 테스트용 NLPService를 생성합니다.
func newTestService() *service.NLPService {
	repo := memory.NewNLPRepository()
	return service.NewNLPService(repo)
}

// TestParseHealthQuery_Success는 정상적인 건강 질의 파싱을 테스트합니다.
func TestParseHealthQuery_Success(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	query, err := svc.ParseHealthQuery(ctx, "user-001", "I have a headache and fever since yesterday")
	if err != nil {
		t.Fatalf("ParseHealthQuery 실패: %v", err)
	}

	if query == nil {
		t.Fatal("query가 nil입니다")
	}
	if query.ID == "" {
		t.Error("query ID가 비어 있습니다")
	}
	if query.UserID != "user-001" {
		t.Errorf("UserID = %q, want %q", query.UserID, "user-001")
	}
	if query.Intent != "health_inquiry" {
		t.Errorf("Intent = %q, want %q", query.Intent, "health_inquiry")
	}
	if query.RawText != "I have a headache and fever since yesterday" {
		t.Errorf("RawText가 원본 텍스트와 다릅니다")
	}

	// "headache"와 "fever" 키워드가 추출되어야 함
	if len(query.Entities) < 2 {
		t.Errorf("Entities 수 = %d, 최소 2개 예상 (headache, fever)", len(query.Entities))
	}

	foundHeadache := false
	foundFever := false
	for _, e := range query.Entities {
		if e == "headache" {
			foundHeadache = true
		}
		if e == "fever" {
			foundFever = true
		}
	}
	if !foundHeadache {
		t.Error("Entities에 'headache'가 포함되어야 합니다")
	}
	if !foundFever {
		t.Error("Entities에 'fever'가 포함되어야 합니다")
	}

	if query.Confidence <= 0.5 {
		t.Errorf("Confidence = %f, 키워드 매칭 시 0.5 초과 예상", query.Confidence)
	}
	if query.CreatedAt.IsZero() {
		t.Error("CreatedAt가 zero입니다")
	}
}

// TestParseHealthQuery_MissingText는 빈 텍스트에 대한 에러 처리를 테스트합니다.
func TestParseHealthQuery_MissingText(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	// 빈 텍스트
	_, err := svc.ParseHealthQuery(ctx, "user-001", "")
	if err == nil {
		t.Fatal("빈 텍스트에 대해 에러가 반환되어야 합니다")
	}

	// 공백만 있는 텍스트
	_, err = svc.ParseHealthQuery(ctx, "user-001", "   ")
	if err == nil {
		t.Fatal("공백만 있는 텍스트에 대해 에러가 반환되어야 합니다")
	}

	// 빈 userID
	_, err = svc.ParseHealthQuery(ctx, "", "some text")
	if err == nil {
		t.Fatal("빈 userID에 대해 에러가 반환되어야 합니다")
	}
}

// TestExtractSymptoms_WithKeywords는 증상 키워드가 포함된 텍스트에서
// 올바르게 증상이 추출되는지 테스트합니다.
func TestExtractSymptoms_WithKeywords(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	extraction, err := svc.ExtractSymptoms(ctx, "I have a terrible headache and high fever with nausea")
	if err != nil {
		t.Fatalf("ExtractSymptoms 실패: %v", err)
	}

	if extraction == nil {
		t.Fatal("extraction이 nil입니다")
	}
	if extraction.ID == "" {
		t.Error("extraction ID가 비어 있습니다")
	}
	if extraction.ProcessedAt.IsZero() {
		t.Error("ProcessedAt가 zero입니다")
	}

	// headache, fever, nausea 최소 3개 증상이 추출되어야 함
	if len(extraction.Symptoms) < 3 {
		t.Errorf("Symptoms 수 = %d, 최소 3개 예상 (headache, fever, nausea)", len(extraction.Symptoms))
	}

	symptomNames := make(map[string]bool)
	for _, s := range extraction.Symptoms {
		symptomNames[s.Name] = true
		if s.Confidence <= 0 {
			t.Errorf("증상 %q의 Confidence가 0 이하입니다", s.Name)
		}
		if s.BodyPart == "" {
			t.Errorf("증상 %q의 BodyPart가 비어 있습니다", s.Name)
		}
		if s.Severity == "" {
			t.Errorf("증상 %q의 Severity가 비어 있습니다", s.Name)
		}
	}

	for _, expected := range []string{"headache", "fever", "nausea"} {
		if !symptomNames[expected] {
			t.Errorf("Symptoms에 %q가 포함되어야 합니다", expected)
		}
	}
}

// TestExtractSymptoms_NoSymptoms는 증상 키워드가 없는 텍스트에서
// 빈 증상 목록이 반환되는지 테스트합니다.
func TestExtractSymptoms_NoSymptoms(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	extraction, err := svc.ExtractSymptoms(ctx, "I feel great today and everything is wonderful")
	if err != nil {
		t.Fatalf("ExtractSymptoms 실패: %v", err)
	}

	if extraction == nil {
		t.Fatal("extraction이 nil입니다")
	}
	if len(extraction.Symptoms) != 0 {
		t.Errorf("Symptoms 수 = %d, 0개 예상 (키워드 없는 텍스트)", len(extraction.Symptoms))
	}
	if extraction.Text != "I feel great today and everything is wonderful" {
		t.Errorf("Text가 원본과 다릅니다")
	}
}

// TestGetSuggestions_Empty는 제안이 없는 질의에 대해 빈 결과를 반환하는지 테스트합니다.
func TestGetSuggestions_Empty(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	// 먼저 질의를 생성하여 유효한 queryID를 얻음
	query, err := svc.ParseHealthQuery(ctx, "user-001", "My blood pressure is high")
	if err != nil {
		t.Fatalf("ParseHealthQuery 실패: %v", err)
	}

	// 해당 질의에 대한 제안 조회 — 저장된 제안이 없으므로 nil 반환
	suggestions, err := svc.GetSuggestions(ctx, query.ID)
	if err != nil {
		t.Fatalf("GetSuggestions 실패: %v", err)
	}

	if suggestions != nil && len(suggestions) != 0 {
		t.Errorf("Suggestions 수 = %d, 0개 예상 (초기 제안 없음)", len(suggestions))
	}

	// 빈 queryID에 대해 에러 반환 확인
	_, err = svc.GetSuggestions(ctx, "")
	if err == nil {
		t.Fatal("빈 queryID에 대해 에러가 반환되어야 합니다")
	}
}

// ============================================================================
// Phase B 신규 테스트
// ============================================================================

// TestParseHealthQuery_Korean은 한국어 건강 질의 파싱을 테스트합니다.
func TestParseHealthQuery_Korean(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	query, err := svc.ParseHealthQuery(ctx, "user-002", "두통이 심하고 혈당 수치가 궁금합니다")
	if err != nil {
		t.Fatalf("ParseHealthQuery 한국어 실패: %v", err)
	}

	// "두통" 증상 엔티티 추출 확인
	foundHeadache := false
	foundBloodSugar := false
	for _, e := range query.Entities {
		if e == "두통" {
			foundHeadache = true
		}
		if e == "blood_sugar" {
			foundBloodSugar = true
		}
	}
	if !foundHeadache {
		t.Error("한국어 '두통' 엔티티가 추출되어야 합니다")
	}
	if !foundBloodSugar {
		t.Error("'혈당' → 'blood_sugar' 인텐트 엔티티가 추출되어야 합니다")
	}

	// 인텐트에 blood_sugar가 포함되어야 함
	if query.Intent != "health_inquiry.blood_sugar" {
		t.Errorf("Intent = %q, want %q", query.Intent, "health_inquiry.blood_sugar")
	}
}

// TestExtractSymptoms_Korean은 한국어 텍스트에서 증상 추출을 테스트합니다.
func TestExtractSymptoms_Korean(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	extraction, err := svc.ExtractSymptoms(ctx, "흉통과 호흡곤란이 있습니다")
	if err != nil {
		t.Fatalf("ExtractSymptoms 한국어 실패: %v", err)
	}

	if len(extraction.Symptoms) < 2 {
		t.Fatalf("Symptoms 수 = %d, 최소 2개 예상 (흉통, 호흡곤란)", len(extraction.Symptoms))
	}

	names := make(map[string]bool)
	for _, s := range extraction.Symptoms {
		names[s.Name] = true
	}
	if !names["흉통"] {
		t.Error("'흉통' 증상이 추출되어야 합니다")
	}
	if !names["호흡곤란"] {
		t.Error("'호흡곤란' 증상이 추출되어야 합니다")
	}

	// 두 증상 모두 severe이므로 urgency는 emergency
	if extraction.Urgency != "emergency" {
		t.Errorf("Urgency = %q, want %q (severe 증상 존재)", extraction.Urgency, "emergency")
	}
}

// TestDetectUrgency_Emergency는 severe 증상 시 emergency 반환을 테스트합니다.
func TestDetectUrgency_Emergency(t *testing.T) {
	symptoms := []service.Symptom{
		{Name: "chest pain", Severity: "severe", BodyPart: "chest", Confidence: 0.95},
	}
	urgency := service.DetectUrgency(symptoms)
	if urgency != "emergency" {
		t.Errorf("DetectUrgency = %q, want %q (severe 증상)", urgency, "emergency")
	}
}

// TestDetectUrgency_Routine은 mild 증상만 있을 때 routine 반환을 테스트합니다.
func TestDetectUrgency_Routine(t *testing.T) {
	symptoms := []service.Symptom{
		{Name: "fatigue", Severity: "mild", BodyPart: "systemic", Confidence: 0.80},
		{Name: "cough", Severity: "mild", BodyPart: "respiratory", Confidence: 0.87},
	}
	urgency := service.DetectUrgency(symptoms)
	if urgency != "routine" {
		t.Errorf("DetectUrgency = %q, want %q (mild 증상만)", urgency, "routine")
	}

	// moderate 2개 이상이면 urgent
	symptoms = []service.Symptom{
		{Name: "headache", Severity: "moderate", BodyPart: "head", Confidence: 0.90},
		{Name: "fever", Severity: "moderate", BodyPart: "systemic", Confidence: 0.92},
	}
	urgency = service.DetectUrgency(symptoms)
	if urgency != "urgent" {
		t.Errorf("DetectUrgency = %q, want %q (moderate 2개)", urgency, "urgent")
	}
}

// TestParseHealthQuery_MixedLanguage는 영한 혼합 텍스트에서 엔티티 추출을 테스트합니다.
func TestParseHealthQuery_MixedLanguage(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	query, err := svc.ParseHealthQuery(ctx, "user-003", "I have a headache and 복통이 있습니다")
	if err != nil {
		t.Fatalf("ParseHealthQuery 혼합 실패: %v", err)
	}

	foundEnglish := false
	foundKorean := false
	for _, e := range query.Entities {
		if e == "headache" {
			foundEnglish = true
		}
		if e == "복통" {
			foundKorean = true
		}
	}
	if !foundEnglish {
		t.Error("영문 'headache' 엔티티가 추출되어야 합니다")
	}
	if !foundKorean {
		t.Error("한국어 '복통' 엔티티가 추출되어야 합니다")
	}

	if query.Confidence <= 0.5 {
		t.Errorf("Confidence = %f, 키워드 매칭 시 0.5 초과 예상", query.Confidence)
	}
}

// TestGenerateSuggestions는 인텐트 기반 제안 생성을 테스트합니다.
func TestGenerateSuggestions(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	// blood_sugar 인텐트 질의 생성
	query, err := svc.ParseHealthQuery(ctx, "user-004", "혈당 수치를 확인하고 싶어요")
	if err != nil {
		t.Fatalf("ParseHealthQuery 실패: %v", err)
	}

	suggestions, err := svc.GenerateSuggestions(ctx, query)
	if err != nil {
		t.Fatalf("GenerateSuggestions 실패: %v", err)
	}

	if len(suggestions) == 0 {
		t.Fatal("제안이 최소 1개 이상 생성되어야 합니다")
	}

	// blood_sugar 관련 제안이 있어야 함
	foundMeasurement := false
	for _, s := range suggestions {
		if s.Category == "measurement" {
			foundMeasurement = true
		}
		if s.QueryID != query.ID {
			t.Errorf("Suggestion QueryID = %q, want %q", s.QueryID, query.ID)
		}
	}
	if !foundMeasurement {
		t.Error("blood_sugar 인텐트에 measurement 카테고리 제안이 있어야 합니다")
	}

	// 저장된 제안 조회
	saved, err := svc.GetSuggestions(ctx, query.ID)
	if err != nil {
		t.Fatalf("GetSuggestions 실패: %v", err)
	}
	if len(saved) != len(suggestions) {
		t.Errorf("저장된 제안 수 = %d, 생성된 수 = %d", len(saved), len(suggestions))
	}
}

// ============================================================================
// Phase C-4 분류기 통합 테스트
// ============================================================================

// mockClassifier는 테스트용 인텐트 분류기입니다.
type mockClassifier struct {
	result *classifier.ClassificationResult
	err    error
}

func (m *mockClassifier) Classify(_ context.Context, _ string) (*classifier.ClassificationResult, error) {
	return m.result, m.err
}

func TestParseHealthQuery_WithClassifier_Success(t *testing.T) {
	svc := newTestService()
	svc.SetIntentClassifier(&mockClassifier{
		result: &classifier.ClassificationResult{
			Intent:     "health_inquiry.blood_sugar",
			Entities:   []string{"blood_sugar", "혈당"},
			Confidence: 0.95,
			Urgency:    "routine",
		},
	})

	query, err := svc.ParseHealthQuery(context.Background(), "user1", "혈당 수치 알려줘")
	if err != nil {
		t.Fatalf("ParseHealthQuery 실패: %v", err)
	}
	if query.Intent != "health_inquiry.blood_sugar" {
		t.Errorf("Intent = %q, want %q", query.Intent, "health_inquiry.blood_sugar")
	}
	if query.Confidence != 0.95 {
		t.Errorf("Confidence = %f, want 0.95", query.Confidence)
	}
}

func TestParseHealthQuery_WithClassifier_Fallback(t *testing.T) {
	svc := newTestService()
	svc.SetIntentClassifier(&mockClassifier{
		err: fmt.Errorf("API 오류"),
	})

	query, err := svc.ParseHealthQuery(context.Background(), "user1", "혈당 수치 알려줘")
	if err != nil {
		t.Fatalf("ParseHealthQuery 실패: %v", err)
	}
	// API 실패 시 키워드 기반 폴백
	if query.Intent != "health_inquiry.blood_sugar" {
		t.Errorf("Intent = %q, want %q (키워드 폴백)", query.Intent, "health_inquiry.blood_sugar")
	}
}

func TestParseHealthQuery_WithNoopClassifier(t *testing.T) {
	svc := newTestService()
	svc.SetIntentClassifier(classifier.NewNoopClassifier())

	query, err := svc.ParseHealthQuery(context.Background(), "user1", "두통이 심합니다")
	if err != nil {
		t.Fatalf("ParseHealthQuery 실패: %v", err)
	}
	// Noop은 nil 반환 → 키워드 기반 폴백
	foundHeadache := false
	for _, e := range query.Entities {
		if e == "두통" {
			foundHeadache = true
		}
	}
	if !foundHeadache {
		t.Error("키워드 폴백에서 '두통' 엔티티가 추출되어야 합니다")
	}
}

func TestParseHealthQuery_WithClassifier_EmergencyUrgency(t *testing.T) {
	svc := newTestService()
	svc.SetIntentClassifier(&mockClassifier{
		result: &classifier.ClassificationResult{
			Intent:     "health_inquiry",
			Entities:   []string{"흉통", "호흡곤란"},
			Confidence: 0.98,
			Urgency:    "emergency",
		},
	})

	query, err := svc.ParseHealthQuery(context.Background(), "user1", "흉통이 심합니다")
	if err != nil {
		t.Fatalf("ParseHealthQuery 실패: %v", err)
	}
	if query.Urgency != "emergency" {
		t.Errorf("Urgency = %q, want %q", query.Urgency, "emergency")
	}
}
