package service_test

import (
	"context"
	"testing"

	"github.com/manpasik/backend/services/vision-service/internal/repository/memory"
	"github.com/manpasik/backend/services/vision-service/internal/service"
	"go.uber.org/zap"
)

func newTestVisionService(t *testing.T) *service.VisionService {
	t.Helper()
	logger := zap.NewNop()
	repo := memory.NewFoodAnalysisRepository()
	return service.NewVisionService(logger, repo)
}

// TestExtendedFoodDB_50Items는 영양소 DB 확장(50개)을 검증합니다.
func TestExtendedFoodDB_50Items(t *testing.T) {
	cases := []string{"전복죽", "감자탕", "설렁탕", "삼계탕", "파스타", "햄버거", "샐러드", "우유"}
	for _, name := range cases {
		item, ok := service.LookupFood(name)
		if !ok {
			t.Errorf("음식 %q 가 DB에 없음", name)
			continue
		}
		if item.CalorieKcal <= 0 {
			t.Errorf("음식 %q 칼로리가 0 이하: %f", name, item.CalorieKcal)
		}
		if len(item.Nutrients) == 0 {
			t.Errorf("음식 %q 영양소 정보 없음", name)
		}
	}
}

// TestCalculateGI_HighGI는 고혈당 식단의 GI 계산을 검증합니다.
func TestCalculateGI_HighGI(t *testing.T) {
	rice, _ := service.LookupFood("쌀밥")  // GI=73
	bread, _ := service.LookupFood("빵")   // GI=70
	noodle, _ := service.LookupFood("라면") // GI=73

	gi := service.CalculateGI([]service.FoodItem{rice, bread, noodle})
	if gi < 70 {
		t.Errorf("고GI 식단 평균 = %d, want >= 70", gi)
	}
	if service.GICategory(gi) != "high" {
		t.Errorf("Category = %q, want %q", service.GICategory(gi), "high")
	}
}

// TestCalculateGI_LowGI는 저혈당 식단의 GI 계산을 검증합니다.
func TestCalculateGI_LowGI(t *testing.T) {
	salad, _ := service.LookupFood("샐러드")    // GI=20
	salmon, _ := service.LookupFood("연어구이")  // GI=30
	tofu, _ := service.LookupFood("두부조림")    // GI=40

	gi := service.CalculateGI([]service.FoodItem{salad, salmon, tofu})
	cat := service.GICategory(gi)
	if cat == "high" {
		t.Errorf("저GI 식단인데 high로 분류됨: gi=%d", gi)
	}
}

// TestGICategory는 GI 카테고리 분류를 검증합니다.
func TestGICategory(t *testing.T) {
	cases := []struct {
		gi   int
		want string
	}{
		{20, "low"},
		{55, "low"},
		{60, "medium"},
		{75, "high"},
	}
	for _, c := range cases {
		got := service.GICategory(c.gi)
		if got != c.want {
			t.Errorf("GICategory(%d) = %q, want %q", c.gi, got, c.want)
		}
	}
}

// TestFilterByDietaryRestriction_Vegan은 비건 필터링을 검증합니다.
func TestFilterByDietaryRestriction_Vegan(t *testing.T) {
	rice, _ := service.LookupFood("쌀밥")
	beef, _ := service.LookupFood("불고기")
	salad, _ := service.LookupFood("샐러드")
	chicken, _ := service.LookupFood("치킨")

	items := []service.FoodItem{rice, beef, salad, chicken}
	filtered := service.FilterByDietaryRestriction(items, []string{"vegan"})

	for _, f := range filtered {
		if f.Name == "불고기" || f.Name == "치킨" {
			t.Errorf("비건 필터에 동물성 식품 %q 포함됨", f.Name)
		}
	}
}

// TestFilterByDietaryRestriction_GlutenFree는 글루텐프리 필터링을 검증합니다.
func TestFilterByDietaryRestriction_GlutenFree(t *testing.T) {
	rice, _ := service.LookupFood("쌀밥")
	pasta, _ := service.LookupFood("파스타")
	bread, _ := service.LookupFood("빵")

	items := []service.FoodItem{rice, pasta, bread}
	filtered := service.FilterByDietaryRestriction(items, []string{"gluten_free"})

	for _, f := range filtered {
		if f.Name == "파스타" || f.Name == "빵" {
			t.Errorf("글루텐프리 필터에 글루텐 식품 %q 포함됨", f.Name)
		}
	}
}

// TestFilterByDietaryRestriction_Halal은 할랄 필터링을 검증합니다.
func TestFilterByDietaryRestriction_Halal(t *testing.T) {
	rice, _ := service.LookupFood("쌀밥")
	pork, _ := service.LookupFood("삼겹살")

	items := []service.FoodItem{rice, pork}
	filtered := service.FilterByDietaryRestriction(items, []string{"halal"})

	for _, f := range filtered {
		if f.Name == "삼겹살" {
			t.Errorf("할랄 필터에 돼지고기 포함됨")
		}
	}
}

// TestGetAllergens는 알레르겐 조회를 검증합니다.
func TestGetAllergens(t *testing.T) {
	allergens := service.GetAllergens("우유")
	hasDairy := false
	for _, a := range allergens {
		if a == "dairy" {
			hasDairy = true
		}
	}
	if !hasDairy {
		t.Error("우유에 dairy 알레르겐이 없음")
	}
}

// TestRecommendMeal_HighGI는 고혈당 보고서에 대한 식단 추천을 검증합니다.
func TestRecommendMeal_HighGI(t *testing.T) {
	report := &service.DailyNutritionReport{
		KcalRatio:   1.0,
		AvgGI:       75,
		GICategory:  "high",
		SodiumRatio: 0.8,
		FiberRatio:  0.7,
		ProteinRatio: 0.8,
	}
	recs := service.RecommendMeal(report, "lunch")
	if len(recs) == 0 {
		t.Fatal("추천 식단이 비어 있음")
	}
	// 고GI일 때는 저GI 식품이 추천되어야 함
	hasLowGI := false
	for _, r := range recs {
		if r == "현미밥" || r == "샐러드" || r == "두부조림" || r == "고등어구이" {
			hasLowGI = true
			break
		}
	}
	if !hasLowGI {
		t.Errorf("고GI 보고서에 저GI 식품 미추천: %v", recs)
	}
}

// TestAnalyzeDailyNutrition_Empty은 빈 식단의 영양 분석을 검증합니다.
func TestAnalyzeDailyNutrition_Empty(t *testing.T) {
	svc := newTestVisionService(t)
	ctx := context.Background()

	report, err := svc.AnalyzeDailyNutrition(ctx, "user-empty")
	if err != nil {
		t.Fatalf("AnalyzeDailyNutrition 실패: %v", err)
	}
	if report.TotalKcal != 0 {
		t.Errorf("TotalKcal = %f, want 0", report.TotalKcal)
	}
	if report.UserID != "user-empty" {
		t.Errorf("UserID = %q, want user-empty", report.UserID)
	}
}

// TestAnalyzeDailyNutrition_AfterMeal은 식사 기록 후 영양 분석을 검증합니다.
func TestAnalyzeDailyNutrition_AfterMeal(t *testing.T) {
	svc := newTestVisionService(t)
	ctx := context.Background()
	userID := "user-meal-001"

	// 점심: 비빔밥(550kcal) + 김치
	if _, err := svc.LogMeal(ctx, userID, []string{"비빔밥", "김치"}, "lunch"); err != nil {
		t.Fatalf("LogMeal 실패: %v", err)
	}

	report, err := svc.AnalyzeDailyNutrition(ctx, userID)
	if err != nil {
		t.Fatalf("AnalyzeDailyNutrition 실패: %v", err)
	}

	if report.TotalKcal < 500 {
		t.Errorf("TotalKcal = %f, want >= 500", report.TotalKcal)
	}
	if report.CarbG <= 0 {
		t.Error("탄수화물 합산이 0 이하")
	}
	if report.AvgGI <= 0 || report.AvgGI > 100 {
		t.Errorf("AvgGI 비정상: %d", report.AvgGI)
	}
}

// TestAnalyzeDailyNutrition_HighSodium은 고나트륨 경고를 검증합니다.
func TestAnalyzeDailyNutrition_HighSodium(t *testing.T) {
	svc := newTestVisionService(t)
	ctx := context.Background()
	userID := "user-sodium"

	// 부대찌개(2200mg 나트륨) — 단일 식사로 권장량 96% 도달
	if _, err := svc.LogMeal(ctx, userID, []string{"부대찌개", "라면"}, "dinner"); err != nil {
		t.Fatalf("LogMeal 실패: %v", err)
	}

	report, _ := svc.AnalyzeDailyNutrition(ctx, userID)
	if report.SodiumRatio < 1.0 {
		t.Errorf("SodiumRatio = %f, want >= 1.0", report.SodiumRatio)
	}

	hasSodiumWarning := false
	for _, w := range report.Warnings {
		if contains(w, "나트륨") {
			hasSodiumWarning = true
			break
		}
	}
	if !hasSodiumWarning {
		t.Errorf("나트륨 경고 미생성. Warnings=%v", report.Warnings)
	}
}

// TestAnalyzeDailyNutrition_Suggestions는 개선 제안 생성을 검증합니다.
func TestAnalyzeDailyNutrition_Suggestions(t *testing.T) {
	svc := newTestVisionService(t)
	ctx := context.Background()
	userID := "user-suggestion"

	// 라면 1개만 (저단백, 고나트륨, 고GI)
	if _, err := svc.LogMeal(ctx, userID, []string{"라면"}, "lunch"); err != nil {
		t.Fatalf("LogMeal 실패: %v", err)
	}

	report, _ := svc.AnalyzeDailyNutrition(ctx, userID)
	if len(report.Suggestions) == 0 {
		t.Error("제안이 생성되지 않음")
	}
}

// TestRecommendMeal_DefaultByMealType은 기본 시간대별 추천을 검증합니다.
func TestRecommendMeal_DefaultByMealType(t *testing.T) {
	report := &service.DailyNutritionReport{
		KcalRatio:    0.8,
		ProteinRatio: 0.8,
		FiberRatio:   0.7,
		SodiumRatio:  0.5,
		GICategory:   "low",
	}

	cases := []struct {
		mealType string
	}{
		{"breakfast"},
		{"lunch"},
		{"dinner"},
	}

	for _, c := range cases {
		recs := service.RecommendMeal(report, c.mealType)
		if len(recs) == 0 {
			t.Errorf("%s 추천이 비어 있음", c.mealType)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && stringContains(s, substr)))
}

func stringContains(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
