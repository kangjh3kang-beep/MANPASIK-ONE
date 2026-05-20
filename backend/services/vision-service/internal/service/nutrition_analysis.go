// Package service: 영양소 분석 + GI 계산 + 식이 추천
//
// 한국 식품영양성분표 기반 음식 DB 확장 (총 50개)과
// 혈당지수(GI), 일일 영양소 분석, 알레르기/식이제한 필터링을 제공합니다.
package service

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	apperrors "github.com/manpasik/backend/shared/errors"
)

// ============================================================================
// 영양소 권장 일일 섭취량 (성인 기준, 한국영양학회)
// ============================================================================

const (
	DailyCalorieKcal    = 2000.0
	DailyCarbohydrateG  = 300.0
	DailyProteinG       = 50.0
	DailyFatG           = 65.0
	DailySodiumMg       = 2300.0
	DailyFiberG         = 25.0
	DailySugarG         = 50.0
)

// ============================================================================
// 한국 음식 영양 DB 확장 (28개 추가, 총 50개)
// ============================================================================

// extendedKoreanFoodDB는 추가 한국 음식 영양 정보입니다.
// init()에서 koreanFoodDB와 병합됩니다.
var extendedKoreanFoodDB = map[string]FoodItem{
	// 죽/이유식
	"전복죽": {Name: "전복죽", Confidence: 0.93, CalorieKcal: 280, PortionG: 350, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 45.0, Unit: "g", DV: 0.15}, {Name: "단백질", Amount: 15.0, Unit: "g", DV: 0.30}, {Name: "지방", Amount: 4.0, Unit: "g", DV: 0.06}}},
	"호박죽": {Name: "호박죽", Confidence: 0.94, CalorieKcal: 220, PortionG: 350, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 48.0, Unit: "g", DV: 0.16}, {Name: "단백질", Amount: 4.0, Unit: "g", DV: 0.08}, {Name: "지방", Amount: 1.5, Unit: "g", DV: 0.02}, {Name: "식이섬유", Amount: 5.0, Unit: "g", DV: 0.20}}},

	// 추가 찌개/탕
	"부대찌개":   {Name: "부대찌개", Confidence: 0.91, CalorieKcal: 420, PortionG: 400, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 25.0, Unit: "g", DV: 0.08}, {Name: "단백질", Amount: 22.0, Unit: "g", DV: 0.44}, {Name: "지방", Amount: 25.0, Unit: "g", DV: 0.38}, {Name: "나트륨", Amount: 2200.0, Unit: "mg", DV: 0.96}}},
	"감자탕":    {Name: "감자탕", Confidence: 0.92, CalorieKcal: 380, PortionG: 500, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 30.0, Unit: "g", DV: 0.10}, {Name: "단백질", Amount: 28.0, Unit: "g", DV: 0.56}, {Name: "지방", Amount: 18.0, Unit: "g", DV: 0.28}, {Name: "나트륨", Amount: 1500.0, Unit: "mg", DV: 0.65}}},
	"설렁탕":    {Name: "설렁탕", Confidence: 0.93, CalorieKcal: 350, PortionG: 500, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 20.0, Unit: "g", DV: 0.07}, {Name: "단백질", Amount: 30.0, Unit: "g", DV: 0.60}, {Name: "지방", Amount: 16.0, Unit: "g", DV: 0.25}, {Name: "나트륨", Amount: 1100.0, Unit: "mg", DV: 0.48}}},
	"갈비탕":    {Name: "갈비탕", Confidence: 0.92, CalorieKcal: 400, PortionG: 500, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 18.0, Unit: "g", DV: 0.06}, {Name: "단백질", Amount: 32.0, Unit: "g", DV: 0.64}, {Name: "지방", Amount: 22.0, Unit: "g", DV: 0.34}}},
	"북엇국":    {Name: "북엇국", Confidence: 0.90, CalorieKcal: 110, PortionG: 300, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 5.0, Unit: "g", DV: 0.02}, {Name: "단백질", Amount: 14.0, Unit: "g", DV: 0.28}, {Name: "지방", Amount: 4.0, Unit: "g", DV: 0.06}, {Name: "나트륨", Amount: 850.0, Unit: "mg", DV: 0.37}}},

	// 반찬 추가
	"두부조림":   {Name: "두부조림", Confidence: 0.91, CalorieKcal: 130, PortionG: 100, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 8.0, Unit: "g", DV: 0.03}, {Name: "단백질", Amount: 12.0, Unit: "g", DV: 0.24}, {Name: "지방", Amount: 6.0, Unit: "g", DV: 0.09}}},
	"콩나물무침":  {Name: "콩나물무침", Confidence: 0.93, CalorieKcal: 35, PortionG: 80, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 4.0, Unit: "g", DV: 0.01}, {Name: "단백질", Amount: 4.0, Unit: "g", DV: 0.08}, {Name: "식이섬유", Amount: 2.0, Unit: "g", DV: 0.08}}},
	"시금치나물":  {Name: "시금치나물", Confidence: 0.94, CalorieKcal: 30, PortionG: 70, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 4.0, Unit: "g", DV: 0.01}, {Name: "단백질", Amount: 3.0, Unit: "g", DV: 0.06}, {Name: "식이섬유", Amount: 2.5, Unit: "g", DV: 0.10}}},
	"무생채":    {Name: "무생채", Confidence: 0.92, CalorieKcal: 25, PortionG: 60, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 5.0, Unit: "g", DV: 0.02}, {Name: "단백질", Amount: 1.0, Unit: "g", DV: 0.02}, {Name: "식이섬유", Amount: 2.0, Unit: "g", DV: 0.08}}},
	"오이무침":   {Name: "오이무침", Confidence: 0.92, CalorieKcal: 20, PortionG: 60, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 4.0, Unit: "g", DV: 0.01}, {Name: "식이섬유", Amount: 1.5, Unit: "g", DV: 0.06}}},
	"멸치볶음":   {Name: "멸치볶음", Confidence: 0.91, CalorieKcal: 90, PortionG: 30, Nutrients: []NutrientInfo{{Name: "단백질", Amount: 12.0, Unit: "g", DV: 0.24}, {Name: "지방", Amount: 3.0, Unit: "g", DV: 0.05}, {Name: "나트륨", Amount: 600.0, Unit: "mg", DV: 0.26}}},
	"장조림":    {Name: "장조림", Confidence: 0.91, CalorieKcal: 110, PortionG: 80, Nutrients: []NutrientInfo{{Name: "단백질", Amount: 18.0, Unit: "g", DV: 0.36}, {Name: "지방", Amount: 3.0, Unit: "g", DV: 0.05}, {Name: "나트륨", Amount: 800.0, Unit: "mg", DV: 0.35}}},

	// 면류 추가
	"잔치국수":   {Name: "잔치국수", Confidence: 0.90, CalorieKcal: 380, PortionG: 500, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 70.0, Unit: "g", DV: 0.23}, {Name: "단백질", Amount: 12.0, Unit: "g", DV: 0.24}, {Name: "지방", Amount: 5.0, Unit: "g", DV: 0.08}, {Name: "나트륨", Amount: 1400.0, Unit: "mg", DV: 0.61}}},
	"비빔국수":   {Name: "비빔국수", Confidence: 0.90, CalorieKcal: 470, PortionG: 400, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 80.0, Unit: "g", DV: 0.27}, {Name: "단백질", Amount: 12.0, Unit: "g", DV: 0.24}, {Name: "지방", Amount: 8.0, Unit: "g", DV: 0.12}}},

	// 분식 추가
	"순대":     {Name: "순대", Confidence: 0.93, CalorieKcal: 350, PortionG: 200, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 50.0, Unit: "g", DV: 0.17}, {Name: "단백질", Amount: 12.0, Unit: "g", DV: 0.24}, {Name: "지방", Amount: 10.0, Unit: "g", DV: 0.15}}},
	"오뎅":     {Name: "오뎅", Confidence: 0.91, CalorieKcal: 180, PortionG: 200, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 25.0, Unit: "g", DV: 0.08}, {Name: "단백질", Amount: 10.0, Unit: "g", DV: 0.20}, {Name: "지방", Amount: 4.0, Unit: "g", DV: 0.06}, {Name: "나트륨", Amount: 1000.0, Unit: "mg", DV: 0.43}}},
	"튀김":     {Name: "튀김", Confidence: 0.90, CalorieKcal: 280, PortionG: 100, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 30.0, Unit: "g", DV: 0.10}, {Name: "단백질", Amount: 5.0, Unit: "g", DV: 0.10}, {Name: "지방", Amount: 16.0, Unit: "g", DV: 0.25}}},

	// 고기/생선
	"삼계탕":    {Name: "삼계탕", Confidence: 0.93, CalorieKcal: 600, PortionG: 800, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 30.0, Unit: "g", DV: 0.10}, {Name: "단백질", Amount: 50.0, Unit: "g", DV: 1.00}, {Name: "지방", Amount: 25.0, Unit: "g", DV: 0.38}}},
	"고등어구이":  {Name: "고등어구이", Confidence: 0.92, CalorieKcal: 280, PortionG: 150, Nutrients: []NutrientInfo{{Name: "단백질", Amount: 22.0, Unit: "g", DV: 0.44}, {Name: "지방", Amount: 20.0, Unit: "g", DV: 0.31}, {Name: "나트륨", Amount: 350.0, Unit: "mg", DV: 0.15}}},
	"갈치구이":   {Name: "갈치구이", Confidence: 0.91, CalorieKcal: 220, PortionG: 150, Nutrients: []NutrientInfo{{Name: "단백질", Amount: 25.0, Unit: "g", DV: 0.50}, {Name: "지방", Amount: 12.0, Unit: "g", DV: 0.18}}},
	"연어구이":   {Name: "연어구이", Confidence: 0.94, CalorieKcal: 280, PortionG: 150, Nutrients: []NutrientInfo{{Name: "단백질", Amount: 30.0, Unit: "g", DV: 0.60}, {Name: "지방", Amount: 18.0, Unit: "g", DV: 0.28}}},

	// 양식
	"파스타":    {Name: "파스타", Confidence: 0.93, CalorieKcal: 450, PortionG: 250, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 70.0, Unit: "g", DV: 0.23}, {Name: "단백질", Amount: 15.0, Unit: "g", DV: 0.30}, {Name: "지방", Amount: 12.0, Unit: "g", DV: 0.18}}},
	"스테이크":   {Name: "스테이크", Confidence: 0.93, CalorieKcal: 480, PortionG: 200, Nutrients: []NutrientInfo{{Name: "단백질", Amount: 45.0, Unit: "g", DV: 0.90}, {Name: "지방", Amount: 32.0, Unit: "g", DV: 0.49}}},
	"햄버거":    {Name: "햄버거", Confidence: 0.94, CalorieKcal: 540, PortionG: 220, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 40.0, Unit: "g", DV: 0.13}, {Name: "단백질", Amount: 25.0, Unit: "g", DV: 0.50}, {Name: "지방", Amount: 30.0, Unit: "g", DV: 0.46}}},
	"샐러드":    {Name: "샐러드", Confidence: 0.92, CalorieKcal: 150, PortionG: 200, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 12.0, Unit: "g", DV: 0.04}, {Name: "단백질", Amount: 5.0, Unit: "g", DV: 0.10}, {Name: "지방", Amount: 8.0, Unit: "g", DV: 0.12}, {Name: "식이섬유", Amount: 4.0, Unit: "g", DV: 0.16}}},

	// 간식/음료
	"빵":      {Name: "빵", Confidence: 0.91, CalorieKcal: 250, PortionG: 100, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 45.0, Unit: "g", DV: 0.15}, {Name: "단백질", Amount: 8.0, Unit: "g", DV: 0.16}, {Name: "지방", Amount: 4.0, Unit: "g", DV: 0.06}}},
	"우유":     {Name: "우유", Confidence: 0.96, CalorieKcal: 130, PortionG: 200, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 10.0, Unit: "g", DV: 0.03}, {Name: "단백질", Amount: 7.0, Unit: "g", DV: 0.14}, {Name: "지방", Amount: 6.0, Unit: "g", DV: 0.09}}},
}

// extendedDBOnce는 확장 음식 DB 병합을 한 번만 실행하도록 보호합니다.
var extendedDBOnce sync.Once

// init은 확장 음식 DB를 메인 DB로 병합합니다.
//
// init 순서가 vision.go보다 늦어도 sync.Once로 안전하게 1회 병합됩니다.
// 동일 키 충돌 시 extended가 우선합니다.
func init() {
	extendedDBOnce.Do(func() {
		if koreanFoodDB == nil {
			return
		}
		for k, v := range extendedKoreanFoodDB {
			koreanFoodDB[k] = v
		}
	})
}

// ============================================================================
// 혈당지수 (GI, Glycemic Index)
// ============================================================================

// foodGITable은 음식별 혈당지수입니다 (0~100, 70+ = High).
// 출처: 호주 시드니대학교 GI Database 한국 식품 표준화
var foodGITable = map[string]int{
	"쌀밥":    73,
	"현미밥":   55,
	"잡곡밥":   50,
	"라면":    73,
	"칼국수":   65,
	"잔치국수":  68,
	"비빔국수":  60,
	"비빔밥":   65,
	"파스타":   45,
	"피자":    60,
	"햄버거":   55,
	"빵":     70,
	"떡국":    85,
	"떡볶이":   80,
	"김밥":    65,
	"치킨":    45,
	"불고기":   45,
	"갈비찜":   45,
	"삼겹살":   30,
	"제육볶음":  50,
	"고등어구이": 25,
	"연어구이":  30,
	"샐러드":   20,
	"우유":    30,
	"김치":    15,
	"두부조림":  40,
}

// CalculateGI는 음식 항목들의 평균 혈당지수를 계산합니다 (탄수화물 가중).
func CalculateGI(items []FoodItem) int {
	totalCarb := 0.0
	weightedGI := 0.0

	for _, item := range items {
		gi, ok := foodGITable[item.Name]
		if !ok {
			gi = 50 // 미상 식품 기본값
		}
		carb := getCarbAmount(item)
		totalCarb += carb
		weightedGI += float64(gi) * carb
	}

	if totalCarb == 0 {
		return 50
	}
	return int(weightedGI / totalCarb)
}

func getCarbAmount(item FoodItem) float64 {
	for _, n := range item.Nutrients {
		if n.Name == "탄수화물" {
			return n.Amount
		}
	}
	return 0
}

// GICategory는 GI 점수를 카테고리로 분류합니다.
func GICategory(gi int) string {
	switch {
	case gi >= 70:
		return "high" // 고GI: 혈당 급상승
	case gi >= 56:
		return "medium" // 중GI
	default:
		return "low" // 저GI: 혈당 안정
	}
}

// ============================================================================
// 일일 영양소 분석
// ============================================================================

// DailyNutritionReport는 일일 영양소 섭취 보고서입니다.
type DailyNutritionReport struct {
	UserID       string
	TotalKcal    float64
	KcalRatio    float64 // 0~1, 권장 대비
	CarbG        float64
	CarbRatio    float64
	ProteinG     float64
	ProteinRatio float64
	FatG         float64
	FatRatio     float64
	SodiumMg     float64
	SodiumRatio  float64
	FiberG       float64
	FiberRatio   float64
	AvgGI        int
	GICategory   string
	Warnings     []string // 영양 경고
	Suggestions  []string // 개선 제안
}

// AnalyzeDailyNutrition은 사용자의 일일 영양 섭취를 분석합니다.
func (s *VisionService) AnalyzeDailyNutrition(ctx context.Context, userID string) (*DailyNutritionReport, error) {
	if userID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "user_id는 필수입니다")
	}

	totalKcal, _, err := s.GetDailySummary(ctx, userID)
	if err != nil {
		return nil, err
	}

	analyses, _, err := s.repo.FindByUserID(ctx, userID, 100, 0)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrInternal, "분석 이력 조회 실패")
	}

	report := &DailyNutritionReport{
		UserID:    userID,
		TotalKcal: totalKcal,
		KcalRatio: totalKcal / DailyCalorieKcal,
	}

	// 모든 오늘 음식 아이템 수집
	var allItems []FoodItem
	today := truncateToDay(timeNow())
	for _, a := range analyses {
		if a.CreatedAt.Before(today) || a.Status != AnalysisStatusCompleted {
			continue
		}
		allItems = append(allItems, a.FoodItems...)
	}

	// 영양소 합산
	for _, item := range allItems {
		for _, n := range item.Nutrients {
			switch n.Name {
			case "탄수화물":
				report.CarbG += n.Amount
			case "단백질":
				report.ProteinG += n.Amount
			case "지방":
				report.FatG += n.Amount
			case "나트륨":
				report.SodiumMg += n.Amount
			case "식이섬유":
				report.FiberG += n.Amount
			}
		}
	}

	report.CarbRatio = report.CarbG / DailyCarbohydrateG
	report.ProteinRatio = report.ProteinG / DailyProteinG
	report.FatRatio = report.FatG / DailyFatG
	report.SodiumRatio = report.SodiumMg / DailySodiumMg
	report.FiberRatio = report.FiberG / DailyFiberG

	// GI 분석
	report.AvgGI = CalculateGI(allItems)
	report.GICategory = GICategory(report.AvgGI)

	// 경고 + 제안 생성
	report.Warnings, report.Suggestions = generateNutritionAdvice(report)

	return report, nil
}

func generateNutritionAdvice(r *DailyNutritionReport) (warnings, suggestions []string) {
	// 칼로리
	if r.KcalRatio > 1.2 {
		warnings = append(warnings, "일일 칼로리 권장량을 20% 이상 초과했습니다")
		suggestions = append(suggestions, "다음 식사는 야채 위주의 저칼로리 식단을 권장합니다")
	} else if r.KcalRatio < 0.6 && r.TotalKcal > 0 {
		warnings = append(warnings, "일일 칼로리 권장량의 60% 미만 섭취했습니다")
		suggestions = append(suggestions, "균형잡힌 식사를 추가로 섭취하세요")
	}

	// 나트륨
	if r.SodiumRatio > 1.0 {
		warnings = append(warnings, "나트륨 일일 권장량을 초과했습니다 — 고혈압 주의")
		suggestions = append(suggestions, "내일은 국/찌개 없이 담백한 식단을 권장합니다")
	}

	// 단백질
	if r.ProteinRatio < 0.7 && r.TotalKcal > 1000 {
		warnings = append(warnings, "단백질 섭취가 부족합니다")
		suggestions = append(suggestions, "닭가슴살, 두부, 계란 등 단백질 식품을 추가하세요")
	}

	// 식이섬유
	if r.FiberRatio < 0.5 && r.TotalKcal > 1000 {
		warnings = append(warnings, "식이섬유 섭취가 부족합니다")
		suggestions = append(suggestions, "채소/과일/잡곡류를 늘리세요")
	}

	// GI
	if r.GICategory == "high" {
		warnings = append(warnings, "오늘 식단의 평균 혈당지수가 높습니다 (당뇨 주의)")
		suggestions = append(suggestions, "현미/잡곡밥, 두부, 채소 등 저GI 식품으로 대체하세요")
	}

	return warnings, suggestions
}

// ============================================================================
// 알레르기/식이제한 필터링
// ============================================================================

// FoodAllergens는 음식별 알레르겐 매핑입니다.
var FoodAllergens = map[string][]string{
	"우유":     {"dairy"},
	"빵":      {"gluten", "wheat"},
	"파스타":    {"gluten", "wheat"},
	"피자":     {"gluten", "wheat", "dairy"},
	"햄버거":    {"gluten", "wheat"},
	"라면":     {"gluten", "wheat"},
	"칼국수":    {"gluten", "wheat"},
	"잔치국수":   {"gluten", "wheat"},
	"비빔국수":   {"gluten", "wheat"},
	"고등어구이":  {"fish"},
	"갈치구이":   {"fish"},
	"연어구이":   {"fish"},
	"멸치볶음":   {"fish"},
	"북엇국":    {"fish"},
	"전복죽":    {"shellfish"},
	"미역국":    {"seafood"},
	"불고기":    {"beef"},
	"갈비찜":    {"beef"},
	"갈비탕":    {"beef"},
	"설렁탕":    {"beef"},
	"스테이크":   {"beef"},
	"제육볶음":   {"pork"},
	"삼겹살":    {"pork"},
	"감자탕":    {"pork"},
	"부대찌개":   {"pork"},
	"치킨":     {"chicken"},
	"삼계탕":    {"chicken"},
	"순대":     {"pork", "egg"},
	"계란말이":   {"egg"},
}

// FilterByDietaryRestriction은 식이 제한에 맞지 않는 음식을 필터링합니다.
//
// restrictions: ["vegan", "vegetarian", "halal", "kosher", "gluten_free", "dairy_free", "nut_free"]
func FilterByDietaryRestriction(items []FoodItem, restrictions []string) []FoodItem {
	if len(restrictions) == 0 {
		return items
	}

	filtered := make([]FoodItem, 0, len(items))
	for _, item := range items {
		if isAllowed(item.Name, restrictions) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func isAllowed(foodName string, restrictions []string) bool {
	allergens := FoodAllergens[foodName]

	for _, r := range restrictions {
		switch strings.ToLower(r) {
		case "vegan":
			for _, a := range allergens {
				if a == "dairy" || a == "egg" || a == "beef" || a == "pork" || a == "chicken" || a == "fish" || a == "seafood" || a == "shellfish" {
					return false
				}
			}
		case "vegetarian":
			for _, a := range allergens {
				if a == "beef" || a == "pork" || a == "chicken" || a == "fish" || a == "seafood" || a == "shellfish" {
					return false
				}
			}
		case "halal":
			for _, a := range allergens {
				if a == "pork" {
					return false
				}
			}
		case "gluten_free":
			for _, a := range allergens {
				if a == "gluten" || a == "wheat" {
					return false
				}
			}
		case "dairy_free":
			for _, a := range allergens {
				if a == "dairy" {
					return false
				}
			}
		}
	}
	return true
}

// GetAllergens는 음식의 알레르겐 목록을 반환합니다.
func GetAllergens(foodName string) []string {
	return FoodAllergens[foodName]
}

// ============================================================================
// 식단 추천
// ============================================================================

// RecommendMeal은 일일 영양 보고서에 기반한 다음 식사를 추천합니다.
func RecommendMeal(report *DailyNutritionReport, mealType string) []string {
	if report == nil {
		return []string{"비빔밥", "샐러드"}
	}

	var recommended []string

	// 칼로리 미달 → 고에너지 식품
	if report.KcalRatio < 0.6 {
		recommended = append(recommended, "비빔밥", "갈비탕", "삼계탕")
	}

	// 단백질 부족 → 고단백
	if report.ProteinRatio < 0.7 {
		recommended = append(recommended, "연어구이", "삼계탕", "스테이크")
	}

	// GI 높음 → 저GI 식품
	if report.GICategory == "high" {
		recommended = append(recommended, "현미밥", "샐러드", "두부조림", "고등어구이")
	}

	// 나트륨 초과 → 저염
	if report.SodiumRatio > 1.0 {
		recommended = append(recommended, "샐러드", "현미밥", "연어구이")
	}

	// 식이섬유 부족
	if report.FiberRatio < 0.5 {
		recommended = append(recommended, "잡곡밥", "샐러드", "시금치나물", "콩나물무침")
	}

	// 기본 추천
	if len(recommended) == 0 {
		switch mealType {
		case "breakfast":
			recommended = []string{"호박죽", "계란말이", "우유"}
		case "lunch":
			recommended = []string{"비빔밥", "김치", "미역국"}
		case "dinner":
			recommended = []string{"잡곡밥", "고등어구이", "시금치나물"}
		default:
			recommended = []string{"샐러드", "우유"}
		}
	}

	// 중복 제거
	return dedupStringList(recommended)
}

func dedupStringList(in []string) []string {
	seen := make(map[string]bool)
	out := make([]string, 0, len(in))
	for _, s := range in {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	sort.Strings(out)
	return out
}

// ============================================================================
// 헬퍼 (시간 모킹 가능)
// ============================================================================

// timeNow는 테스트 가능한 현재 시간 반환자입니다.
var timeNow = func() time.Time { return time.Now().UTC() }

func truncateToDay(t time.Time) time.Time {
	return t.Truncate(24 * time.Hour)
}
