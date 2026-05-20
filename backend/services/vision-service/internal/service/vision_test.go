package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/manpasik/backend/services/vision-service/internal/repository/memory"
	"github.com/manpasik/backend/services/vision-service/internal/service"
	"go.uber.org/zap"
)

// mockVisionAnalyzer는 VisionAnalyzer 모의 구현입니다.
type mockVisionAnalyzer struct {
	items []service.FoodItem
	err   error
}

func (m *mockVisionAnalyzer) AnalyzeFood(_ context.Context, _ string) ([]service.FoodItem, error) {
	return m.items, m.err
}

func setupTestService() *service.VisionService {
	logger := zap.NewNop()
	repo := memory.NewFoodAnalysisRepository()
	return service.NewVisionService(logger, repo)
}

func TestAnalyzeFood_SimulationMode(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	result, err := svc.AnalyzeFood(ctx, "user-1", "https://example.com/food.jpg", "lunch")
	if err != nil {
		t.Fatalf("음식 분석 실패: %v", err)
	}
	if result == nil {
		t.Fatal("분석 결과가 nil이면 안 됨")
	}
	if result.Status != service.AnalysisStatusCompleted {
		t.Fatalf("상태 불일치: got %d, want %d", result.Status, service.AnalysisStatusCompleted)
	}
	if result.TotalCalorieKcal <= 0 {
		t.Fatalf("총 칼로리가 0보다 커야 함: got %f", result.TotalCalorieKcal)
	}
	if len(result.FoodItems) == 0 {
		t.Fatal("음식 항목이 비어 있으면 안 됨")
	}
	if result.MealType != "lunch" {
		t.Fatalf("MealType 불일치: got %s, want lunch", result.MealType)
	}
}

func TestAnalyzeFood_MissingUserID(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	_, err := svc.AnalyzeFood(ctx, "", "https://example.com/food.jpg", "lunch")
	if err == nil {
		t.Fatal("빈 user_id에 에러가 반환되어야 함")
	}
}

func TestAnalyzeFood_MissingImageURL(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	_, err := svc.AnalyzeFood(ctx, "user-1", "", "lunch")
	if err == nil {
		t.Fatal("빈 image_url에 에러가 반환되어야 함")
	}
}

func TestGetAnalysis_Success(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	// 먼저 분석 생성
	created, err := svc.AnalyzeFood(ctx, "user-1", "https://example.com/food.jpg", "dinner")
	if err != nil {
		t.Fatalf("음식 분석 실패: %v", err)
	}

	// 결과 조회
	result, err := svc.GetAnalysis(ctx, created.ID)
	if err != nil {
		t.Fatalf("분석 결과 조회 실패: %v", err)
	}
	if result.ID != created.ID {
		t.Fatalf("ID 불일치: got %s, want %s", result.ID, created.ID)
	}
}

func TestGetAnalysis_NotFound(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	_, err := svc.GetAnalysis(ctx, "nonexistent-id")
	if err == nil {
		t.Fatal("존재하지 않는 분석에 에러가 반환되어야 함")
	}
}

func TestListAnalyses_Success(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	// 3건 분석 생성
	for i := 0; i < 3; i++ {
		_, err := svc.AnalyzeFood(ctx, "user-list", "https://example.com/food.jpg", "lunch")
		if err != nil {
			t.Fatalf("음식 분석 %d 실패: %v", i, err)
		}
	}

	analyses, total, err := svc.ListAnalyses(ctx, "user-list", 10, 0)
	if err != nil {
		t.Fatalf("분석 이력 조회 실패: %v", err)
	}
	if total != 3 {
		t.Fatalf("총 분석 수 불일치: got %d, want 3", total)
	}
	if len(analyses) != 3 {
		t.Fatalf("반환 분석 수 불일치: got %d, want 3", len(analyses))
	}
}

func TestListAnalyses_DefaultLimit(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	_, err := svc.AnalyzeFood(ctx, "user-default", "https://example.com/food.jpg", "breakfast")
	if err != nil {
		t.Fatalf("음식 분석 실패: %v", err)
	}

	// limit <= 0이면 기본값 20
	analyses, _, err := svc.ListAnalyses(ctx, "user-default", 0, 0)
	if err != nil {
		t.Fatalf("분석 이력 조회 실패: %v", err)
	}
	if len(analyses) != 1 {
		t.Fatalf("반환 분석 수 불일치: got %d, want 1", len(analyses))
	}
}

func TestListAnalyses_Empty(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	analyses, total, err := svc.ListAnalyses(ctx, "user-no-data", 10, 0)
	if err != nil {
		t.Fatalf("분석 이력 조회 실패: %v", err)
	}
	if total != 0 {
		t.Fatalf("총 분석 수가 0이어야 함: got %d", total)
	}
	if len(analyses) != 0 {
		t.Fatalf("반환 분석이 비어야 함: got %d", len(analyses))
	}
}

func TestGetDailySummary_NoData(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	totalKcal, mealBreakdown, err := svc.GetDailySummary(ctx, "user-no-summary")
	if err != nil {
		t.Fatalf("일일 요약 조회 실패: %v", err)
	}
	if totalKcal != 0 {
		t.Fatalf("데이터 없는 사용자의 totalKcal이 0이어야 함: got %f", totalKcal)
	}
	if mealBreakdown == nil {
		t.Fatal("mealBreakdown가 nil이면 안 됨")
	}
}

func TestGetDailySummary_WithData(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	// 오늘 데이터 생성
	_, err := svc.AnalyzeFood(ctx, "user-summary", "https://example.com/food.jpg", "lunch")
	if err != nil {
		t.Fatalf("음식 분석 실패: %v", err)
	}

	totalKcal, mealBreakdown, err := svc.GetDailySummary(ctx, "user-summary")
	if err != nil {
		t.Fatalf("일일 요약 조회 실패: %v", err)
	}
	if totalKcal <= 0 {
		t.Fatalf("오늘 데이터가 있으므로 totalKcal > 0이어야 함: got %f", totalKcal)
	}
	if mealBreakdown["lunch"] <= 0 {
		t.Fatalf("lunch 칼로리가 0보다 커야 함: got %f", mealBreakdown["lunch"])
	}
}

func TestFoodAnalysisStatusToString(t *testing.T) {
	tests := []struct {
		status service.FoodAnalysisStatus
		want   string
	}{
		{service.AnalysisStatusPending, "pending"},
		{service.AnalysisStatusProcessing, "processing"},
		{service.AnalysisStatusCompleted, "completed"},
		{service.AnalysisStatusFailed, "failed"},
		{service.AnalysisStatusUnknown, "unknown(0)"},
	}

	for _, tt := range tests {
		got := service.FoodAnalysisStatusToString(tt.status)
		if got != tt.want {
			t.Errorf("FoodAnalysisStatusToString(%d) = %s, want %s", tt.status, got, tt.want)
		}
	}
}

func TestSetAnalyzer(t *testing.T) {
	svc := setupTestService()
	// SetAnalyzer에 nil을 설정해도 panic이 발생하지 않아야 함
	svc.SetAnalyzer(nil)
}

// ============================================================================
// Phase B 신규 테스트
// ============================================================================

func TestAnalyzeFood_InvalidMealType(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	_, err := svc.AnalyzeFood(ctx, "user-1", "https://example.com/food.jpg", "brunch")
	if err == nil {
		t.Fatal("지원하지 않는 식사 타입에 에러가 반환되어야 합니다")
	}
}

func TestAnalyzeFood_EmptyMealTypeDefaultsToSnack(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	result, err := svc.AnalyzeFood(ctx, "user-1", "https://example.com/food.jpg", "")
	if err != nil {
		t.Fatalf("빈 식사 타입 분석 실패: %v", err)
	}
	if result.MealType != "snack" {
		t.Errorf("기본 MealType = %s, want snack", result.MealType)
	}
}

func TestLogMeal_Success(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	result, err := svc.LogMeal(ctx, "user-1", []string{"쌀밥", "김치찌개", "김치"}, "lunch")
	if err != nil {
		t.Fatalf("LogMeal 실패: %v", err)
	}
	if result.Status != service.AnalysisStatusCompleted {
		t.Errorf("Status = %d, want completed", result.Status)
	}
	if len(result.FoodItems) != 3 {
		t.Errorf("FoodItems 수 = %d, want 3", len(result.FoodItems))
	}
	if result.TotalCalorieKcal <= 0 {
		t.Error("총 칼로리가 0보다 커야 합니다")
	}
	if result.MealType != "lunch" {
		t.Errorf("MealType = %s, want lunch", result.MealType)
	}
}

func TestLogMeal_UnknownFood(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	_, err := svc.LogMeal(ctx, "user-1", []string{"알수없는음식", "존재하지않는요리"}, "dinner")
	if err == nil {
		t.Fatal("인식할 수 없는 음식만 있을 때 에러가 반환되어야 합니다")
	}
}

func TestLogMeal_EmptyFoods(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	_, err := svc.LogMeal(ctx, "user-1", []string{}, "lunch")
	if err == nil {
		t.Fatal("빈 음식 목록에 에러가 반환되어야 합니다")
	}
}

func TestDeleteAnalysis(t *testing.T) {
	svc := setupTestService()
	ctx := context.Background()

	created, _ := svc.AnalyzeFood(ctx, "user-1", "https://example.com/food.jpg", "dinner")

	err := svc.DeleteAnalysis(ctx, created.ID)
	if err != nil {
		t.Fatalf("DeleteAnalysis 실패: %v", err)
	}

	_, err = svc.GetAnalysis(ctx, created.ID)
	if err == nil {
		t.Fatal("삭제 후 조회에 에러가 반환되어야 합니다")
	}
}

// =============================================================================
// Phase C-1: VisionAnalyzer 통합 테스트
// =============================================================================

func TestAnalyzeFood_WithMockAnalyzer_성공(t *testing.T) {
	svc := setupTestService()
	svc.SetAnalyzer(&mockVisionAnalyzer{
		items: []service.FoodItem{
			{Name: "비빔밥", Confidence: 0.95, CalorieKcal: 550, PortionG: 400},
		},
	})
	ctx := context.Background()

	result, err := svc.AnalyzeFood(ctx, "user-1", "https://example.com/food.jpg", "lunch")
	if err != nil {
		t.Fatalf("AI 분석 실패: %v", err)
	}
	if result.Status != service.AnalysisStatusCompleted {
		t.Errorf("상태 불일치: got %d, want completed", result.Status)
	}
	if len(result.FoodItems) != 1 {
		t.Fatalf("음식 항목 수 불일치: got %d, want 1", len(result.FoodItems))
	}
	if result.FoodItems[0].Name != "비빔밥" {
		t.Errorf("음식 이름 불일치: got %s", result.FoodItems[0].Name)
	}
	if result.TotalCalorieKcal != 550 {
		t.Errorf("총 칼로리 불일치: got %f, want 550", result.TotalCalorieKcal)
	}
}

func TestAnalyzeFood_WithMockAnalyzer_실패(t *testing.T) {
	svc := setupTestService()
	svc.SetAnalyzer(&mockVisionAnalyzer{
		err: fmt.Errorf("Vision API 호출 실패"),
	})
	ctx := context.Background()

	result, err := svc.AnalyzeFood(ctx, "user-1", "https://example.com/food.jpg", "lunch")
	if err != nil {
		t.Fatalf("서비스 에러가 아니라 분석 실패 결과를 반환해야 함: %v", err)
	}
	if result.Status != service.AnalysisStatusFailed {
		t.Errorf("상태가 failed여야 함: got %d", result.Status)
	}
	if result.ErrorMessage == "" {
		t.Error("에러 메시지가 비어있으면 안 됨")
	}
}

func TestAnalyzeFood_WithMockAnalyzer_빈결과(t *testing.T) {
	svc := setupTestService()
	svc.SetAnalyzer(&mockVisionAnalyzer{
		items: []service.FoodItem{},
	})
	ctx := context.Background()

	result, err := svc.AnalyzeFood(ctx, "user-1", "https://example.com/food.jpg", "dinner")
	if err != nil {
		t.Fatalf("빈 결과 분석 실패: %v", err)
	}
	if result.Status != service.AnalysisStatusCompleted {
		t.Errorf("상태 불일치: got %d", result.Status)
	}
	if result.TotalCalorieKcal != 0 {
		t.Errorf("빈 결과의 칼로리는 0이어야 함: got %f", result.TotalCalorieKcal)
	}
}

func TestLookupFood(t *testing.T) {
	item, ok := service.LookupFood("비빔밥")
	if !ok {
		t.Fatal("비빔밥이 DB에 존재해야 합니다")
	}
	if item.Name != "비빔밥" {
		t.Errorf("Name = %s, want 비빔밥", item.Name)
	}
	if item.CalorieKcal <= 0 {
		t.Error("칼로리가 0보다 커야 합니다")
	}

	_, ok = service.LookupFood("존재하지않는음식")
	if ok {
		t.Error("존재하지 않는 음식은 false를 반환해야 합니다")
	}
}
