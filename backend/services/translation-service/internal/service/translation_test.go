package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/manpasik/backend/services/translation-service/internal/repository/memory"
	"github.com/manpasik/backend/services/translation-service/internal/service"
	"go.uber.org/zap"
)

// mockTranslator는 외부 번역 API 모의 구현입니다.
type mockTranslator struct {
	result string
	err    error
}

func (m *mockTranslator) Translate(_ context.Context, _, _, _ string) (string, error) {
	return m.result, m.err
}

func setupTranslationService() *service.TranslationService {
	logger := zap.NewNop()
	translationRepo := memory.NewTranslationRepository()
	usageRepo := memory.NewUsageRepository()
	return service.NewTranslationService(logger, translationRepo, usageRepo)
}

func TestTranslateText_KoToEn(t *testing.T) {
	svc := setupTranslationService()
	ctx := context.Background()

	translated, sourceLang, confidence, translationID, err := svc.TranslateText(ctx, "혈압 측정", "", "en", "medical", "user-1")
	if err != nil {
		t.Fatalf("TranslateText 실패: %v", err)
	}
	if sourceLang != "ko" {
		t.Fatalf("소스 언어 예상 ko, 실제: %s", sourceLang)
	}
	if translated == "" {
		t.Fatal("번역 결과가 비어 있습니다")
	}
	if confidence <= 0 {
		t.Fatalf("신뢰도가 0 이하입니다: %f", confidence)
	}
	if translationID == "" {
		t.Fatal("번역 ID가 비어 있습니다")
	}
}

func TestTranslateText_EnToKo(t *testing.T) {
	svc := setupTranslationService()
	ctx := context.Background()

	translated, _, _, _, err := svc.TranslateText(ctx, "blood pressure check", "en", "ko", "", "user-1")
	if err != nil {
		t.Fatalf("TranslateText 실패: %v", err)
	}
	if translated == "" {
		t.Fatal("번역 결과가 비어 있습니다")
	}
}

func TestTranslateText_EmptyText(t *testing.T) {
	svc := setupTranslationService()
	ctx := context.Background()

	_, _, _, _, err := svc.TranslateText(ctx, "", "", "en", "", "")
	if err == nil {
		t.Fatal("빈 텍스트에 에러가 발생해야 합니다")
	}
}

func TestTranslateText_EmptyTargetLang(t *testing.T) {
	svc := setupTranslationService()
	ctx := context.Background()

	_, _, _, _, err := svc.TranslateText(ctx, "hello", "en", "", "", "")
	if err == nil {
		t.Fatal("빈 target_language에 에러가 발생해야 합니다")
	}
}

func TestDetectLanguage_Korean(t *testing.T) {
	svc := setupTranslationService()
	ctx := context.Background()

	detected, primary, err := svc.DetectLanguage(ctx, "안녕하세요")
	if err != nil {
		t.Fatalf("DetectLanguage 실패: %v", err)
	}
	if primary != "ko" {
		t.Fatalf("주요 언어 예상 ko, 실제: %s", primary)
	}
	if len(detected) == 0 {
		t.Fatal("감지된 언어가 없습니다")
	}
}

func TestDetectLanguage_English(t *testing.T) {
	svc := setupTranslationService()
	ctx := context.Background()

	_, primary, err := svc.DetectLanguage(ctx, "Hello world")
	if err != nil {
		t.Fatalf("DetectLanguage 실패: %v", err)
	}
	if primary != "en" {
		t.Fatalf("주요 언어 예상 en, 실제: %s", primary)
	}
}

func TestDetectLanguage_EmptyText(t *testing.T) {
	svc := setupTranslationService()
	ctx := context.Background()

	_, _, err := svc.DetectLanguage(ctx, "")
	if err == nil {
		t.Fatal("빈 텍스트에 에러가 발생해야 합니다")
	}
}

func TestListSupportedLanguages(t *testing.T) {
	svc := setupTranslationService()
	ctx := context.Background()

	languages, err := svc.ListSupportedLanguages(ctx)
	if err != nil {
		t.Fatalf("ListSupportedLanguages 실패: %v", err)
	}
	if len(languages) < 8 {
		t.Fatalf("지원 언어 수 최소 8, 실제: %d", len(languages))
	}

	// 한국어 의료 지원 확인
	found := false
	for _, l := range languages {
		if l.Code == "ko" && l.SupportsMedical {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("한국어 의료 번역 지원이 필요합니다")
	}
}

func TestTranslateBatch(t *testing.T) {
	svc := setupTranslationService()
	ctx := context.Background()

	texts := []string{"혈압", "혈당", "체온"}
	translated, sourceLang, totalChars, err := svc.TranslateBatch(ctx, texts, "", "en", "user-1")
	if err != nil {
		t.Fatalf("TranslateBatch 실패: %v", err)
	}
	if len(translated) != 3 {
		t.Fatalf("번역 결과 수 예상 3, 실제: %d", len(translated))
	}
	if sourceLang != "ko" {
		t.Fatalf("소스 언어 예상 ko, 실제: %s", sourceLang)
	}
	if totalChars <= 0 {
		t.Fatalf("총 문자 수가 0 이하입니다: %d", totalChars)
	}
}

func TestTranslateBatch_EmptyTexts(t *testing.T) {
	svc := setupTranslationService()
	ctx := context.Background()

	_, _, _, err := svc.TranslateBatch(ctx, []string{}, "", "en", "")
	if err == nil {
		t.Fatal("빈 texts에 에러가 발생해야 합니다")
	}
}

func TestGetTranslationHistory(t *testing.T) {
	svc := setupTranslationService()
	ctx := context.Background()

	// 번역 2건 생성
	svc.TranslateText(ctx, "혈압", "", "en", "", "user-1")
	svc.TranslateText(ctx, "혈당", "", "en", "", "user-1")

	records, total, err := svc.GetTranslationHistory(ctx, "user-1", 10, 0)
	if err != nil {
		t.Fatalf("GetTranslationHistory 실패: %v", err)
	}
	if total != 2 {
		t.Fatalf("이력 수 예상 2, 실제: %d", total)
	}
	if len(records) != 2 {
		t.Fatalf("반환 수 예상 2, 실제: %d", len(records))
	}
}

func TestGetTranslationUsage(t *testing.T) {
	svc := setupTranslationService()
	ctx := context.Background()

	svc.TranslateText(ctx, "Hello world", "en", "ko", "", "user-1")

	usage, err := svc.GetTranslationUsage(ctx, "user-1")
	if err != nil {
		t.Fatalf("GetTranslationUsage 실패: %v", err)
	}
	if usage.TotalCharacters <= 0 {
		t.Fatal("총 문자 수가 0 이하입니다")
	}
	if usage.TotalRequests != 1 {
		t.Fatalf("총 요청 수 예상 1, 실제: %d", usage.TotalRequests)
	}
	if usage.MonthlyLimit <= 0 {
		t.Fatal("월간 한도가 설정되어야 합니다")
	}
}

// =============================================================================
// Phase C-1: 외부 번역기 통합 테스트
// =============================================================================

func TestTranslateText_외부번역기_성공(t *testing.T) {
	svc := setupTranslationService()
	svc.SetTranslator(&mockTranslator{result: "blood pressure measurement", err: nil})
	ctx := context.Background()

	translated, _, confidence, _, err := svc.TranslateText(ctx, "혈압 측정", "", "en", "", "user-1")
	if err != nil {
		t.Fatalf("TranslateText 실패: %v", err)
	}
	if translated != "blood pressure measurement" {
		t.Errorf("번역 결과 예상 'blood pressure measurement', 실제: '%s'", translated)
	}
	if confidence < 0.96 {
		t.Errorf("외부 번역 신뢰도 예상 >= 0.96, 실제: %f", confidence)
	}
}

func TestTranslateText_외부번역기_의료(t *testing.T) {
	svc := setupTranslationService()
	svc.SetTranslator(&mockTranslator{result: "blood pressure", err: nil})
	ctx := context.Background()

	_, _, confidence, _, err := svc.TranslateText(ctx, "혈압", "", "en", "medical", "user-1")
	if err != nil {
		t.Fatalf("TranslateText 실패: %v", err)
	}
	if confidence < 0.98 {
		t.Errorf("외부+의료 신뢰도 예상 >= 0.98, 실제: %f", confidence)
	}
}

func TestTranslateText_외부번역기_실패_내부폴백(t *testing.T) {
	svc := setupTranslationService()
	svc.SetTranslator(&mockTranslator{result: "", err: fmt.Errorf("API 오류")})
	ctx := context.Background()

	translated, _, confidence, _, err := svc.TranslateText(ctx, "혈압", "ko", "en", "", "user-1")
	if err != nil {
		t.Fatalf("폴백 번역도 실패: %v", err)
	}
	// 폴백은 내부 의료사전을 사용
	if translated != "blood pressure" {
		t.Errorf("폴백 번역 예상 'blood pressure', 실제: '%s'", translated)
	}
	if confidence != 0.92 {
		t.Errorf("폴백 신뢰도 예상 0.92, 실제: %f", confidence)
	}
}

func TestTranslateText_외부번역기_빈결과_내부폴백(t *testing.T) {
	svc := setupTranslationService()
	svc.SetTranslator(&mockTranslator{result: "", err: nil})
	ctx := context.Background()

	translated, _, _, _, err := svc.TranslateText(ctx, "혈압", "ko", "en", "", "user-1")
	if err != nil {
		t.Fatalf("빈 결과 폴백 실패: %v", err)
	}
	// 빈 결과도 폴백 사용
	if translated != "blood pressure" {
		t.Errorf("빈 결과 폴백 예상 'blood pressure', 실제: '%s'", translated)
	}
}

func TestTranslateRealtime_성공(t *testing.T) {
	svc := setupTranslationService()
	ctx := context.Background()

	translated, _, confidence, latencyMs, err := svc.TranslateRealtime(ctx, "혈압", "ko", "en", "medical", "session-1", true)
	if err != nil {
		t.Fatalf("TranslateRealtime 실패: %v", err)
	}
	if translated == "" {
		t.Fatal("번역 결과가 비어 있습니다")
	}
	if confidence <= 0 {
		t.Fatalf("신뢰도 0 이하: %f", confidence)
	}
	if latencyMs < 0 {
		t.Fatalf("지연시간 음수: %d", latencyMs)
	}
}

func TestMedicalTermTranslation(t *testing.T) {
	svc := setupTranslationService()
	ctx := context.Background()

	translated, _, _, _, err := svc.TranslateText(ctx, "혈압", "ko", "en", "medical", "user-1")
	if err != nil {
		t.Fatalf("TranslateText 실패: %v", err)
	}
	if translated != "blood pressure" {
		t.Fatalf("의료 용어 번역 예상 'blood pressure', 실제: '%s'", translated)
	}
}
