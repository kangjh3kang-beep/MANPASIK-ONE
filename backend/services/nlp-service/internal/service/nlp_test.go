package service_test

import (
	"context"
	"testing"

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
