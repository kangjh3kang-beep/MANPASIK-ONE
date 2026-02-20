package service_test

import (
	"context"
	"testing"

	"github.com/manpasik/backend/services/vision-service/internal/repository/memory"
	"github.com/manpasik/backend/services/vision-service/internal/service"
	"go.uber.org/zap"
)

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
