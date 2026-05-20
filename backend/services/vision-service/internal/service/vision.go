// Package service는 vision-service의 비즈니스 로직을 구현합니다.
//
// 기능:
// - 음식 이미지 분석 (카메라/갤러리)
// - 칼로리 추정
// - 영양소 분석
// - 음식 분석 이력 조회
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	apperrors "github.com/manpasik/backend/shared/errors"
	"go.uber.org/zap"
)

// ============================================================================
// 상수 및 검증
// ============================================================================

// ValidMealTypes는 허용되는 식사 타입입니다.
var ValidMealTypes = map[string]bool{
	"breakfast": true,
	"lunch":     true,
	"dinner":    true,
	"snack":     true,
}

// ============================================================================
// 도메인 모델
// ============================================================================

// FoodAnalysisStatus는 음식 분석 상태입니다.
type FoodAnalysisStatus int32

const (
	AnalysisStatusUnknown    FoodAnalysisStatus = 0
	AnalysisStatusPending    FoodAnalysisStatus = 1
	AnalysisStatusProcessing FoodAnalysisStatus = 2
	AnalysisStatusCompleted  FoodAnalysisStatus = 3
	AnalysisStatusFailed     FoodAnalysisStatus = 4
)

// NutrientInfo는 영양소 정보입니다.
type NutrientInfo struct {
	Name   string  // 영양소 이름 (예: "탄수화물", "단백질", "지방")
	Amount float64 // 양 (g)
	Unit   string  // 단위 (g, mg, kcal)
	DV     float64 // 일일 권장량 대비 비율 (0.0~1.0)
}

// FoodItem은 인식된 음식 항목입니다.
type FoodItem struct {
	Name       string         // 음식 이름 (예: "김치찌개", "현미밥")
	Confidence float64        // 인식 신뢰도 (0.0~1.0)
	CalorieKcal float64       // 칼로리 (kcal)
	PortionG   float64        // 1인분 기준 무게 (g)
	Nutrients  []NutrientInfo // 영양소 목록
}

// FoodAnalysis는 음식 분석 결과 엔티티입니다.
type FoodAnalysis struct {
	ID              string
	UserID          string
	ImageURL        string             // S3/MinIO 이미지 경로
	Status          FoodAnalysisStatus
	TotalCalorieKcal float64
	FoodItems       []FoodItem
	MealType        string // "breakfast", "lunch", "dinner", "snack"
	AnalyzedAt      *time.Time
	CreatedAt       time.Time
	ErrorMessage    string
}

// ============================================================================
// 리포지토리 인터페이스
// ============================================================================

// FoodAnalysisRepository는 음식 분석 결과 저장소 인터페이스입니다.
type FoodAnalysisRepository interface {
	Save(ctx context.Context, analysis *FoodAnalysis) error
	FindByID(ctx context.Context, id string) (*FoodAnalysis, error)
	FindByUserID(ctx context.Context, userID string, limit, offset int32) ([]*FoodAnalysis, int32, error)
	Update(ctx context.Context, analysis *FoodAnalysis) error
	Delete(ctx context.Context, id string) error
}

// VisionAnalyzer는 AI 비전 분석 인터페이스입니다.
// 실제 구현은 TFLite, Cloud Vision API, 또는 자체 모델이 됩니다.
type VisionAnalyzer interface {
	AnalyzeFood(ctx context.Context, imageURL string) ([]FoodItem, error)
}

// ============================================================================
// 서비스
// ============================================================================

// VisionService는 음식 비전 분석 서비스입니다.
type VisionService struct {
	logger   *zap.Logger
	repo     FoodAnalysisRepository
	analyzer VisionAnalyzer // optional: nil이면 시뮬레이션 모드
}

// NewVisionService는 새 VisionService를 생성합니다.
func NewVisionService(logger *zap.Logger, repo FoodAnalysisRepository) *VisionService {
	return &VisionService{
		logger: logger,
		repo:   repo,
	}
}

// SetAnalyzer는 AI 비전 분석기를 설정합니다 (optional).
func (s *VisionService) SetAnalyzer(a VisionAnalyzer) {
	s.analyzer = a
}

// AnalyzeFood는 음식 이미지를 분석합니다.
func (s *VisionService) AnalyzeFood(ctx context.Context, userID, imageURL, mealType string) (*FoodAnalysis, error) {
	if userID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "user_id는 필수입니다")
	}
	if imageURL == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "image_url은 필수입니다")
	}
	if mealType == "" {
		mealType = "snack"
	}
	if !ValidMealTypes[mealType] {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "지원하지 않는 식사 타입: "+mealType)
	}

	now := time.Now().UTC()
	analysis := &FoodAnalysis{
		ID:        uuid.New().String(),
		UserID:    userID,
		ImageURL:  imageURL,
		Status:    AnalysisStatusProcessing,
		MealType:  mealType,
		CreatedAt: now,
	}

	// 분석 실행
	var foodItems []FoodItem
	var analysisErr error

	if s.analyzer != nil {
		// 실제 AI 분석 (Cloud Vision / TFLite)
		foodItems, analysisErr = s.analyzer.AnalyzeFood(ctx, imageURL)
	} else {
		// 룩업 테이블 기반 기본 분석 (한식 1인분 구성)
		foodItems = s.fallbackAnalysis(imageURL)
	}

	if analysisErr != nil {
		analysis.Status = AnalysisStatusFailed
		analysis.ErrorMessage = analysisErr.Error()
		if err := s.repo.Save(ctx, analysis); err != nil {
			return nil, apperrors.New(apperrors.ErrInternal, "분석 결과 저장 실패")
		}
		return analysis, nil
	}

	// 총 칼로리 계산
	var totalCal float64
	for _, item := range foodItems {
		totalCal += item.CalorieKcal
	}

	analysis.Status = AnalysisStatusCompleted
	analysis.FoodItems = foodItems
	analysis.TotalCalorieKcal = totalCal
	analysis.AnalyzedAt = &now

	if err := s.repo.Save(ctx, analysis); err != nil {
		return nil, apperrors.New(apperrors.ErrInternal, "분석 결과 저장 실패")
	}

	s.logger.Info("음식 분석 완료",
		zap.String("analysis_id", analysis.ID),
		zap.String("user_id", userID),
		zap.Float64("total_kcal", totalCal),
		zap.Int("food_items", len(foodItems)),
	)

	return analysis, nil
}

// GetAnalysis는 음식 분석 결과를 조회합니다.
func (s *VisionService) GetAnalysis(ctx context.Context, analysisID string) (*FoodAnalysis, error) {
	analysis, err := s.repo.FindByID(ctx, analysisID)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrInternal, "분석 결과 조회 실패")
	}
	if analysis == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "분석 결과를 찾을 수 없습니다")
	}
	return analysis, nil
}

// ListAnalyses는 사용자의 음식 분석 이력을 조회합니다.
func (s *VisionService) ListAnalyses(ctx context.Context, userID string, limit, offset int32) ([]*FoodAnalysis, int32, error) {
	if limit <= 0 {
		limit = 20
	}
	analyses, total, err := s.repo.FindByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, apperrors.New(apperrors.ErrInternal, "분석 이력 조회 실패")
	}
	return analyses, total, nil
}

// GetDailySummary는 사용자의 일일 칼로리 섭취 요약을 조회합니다.
func (s *VisionService) GetDailySummary(ctx context.Context, userID string) (totalKcal float64, mealBreakdown map[string]float64, err error) {
	analyses, _, listErr := s.repo.FindByUserID(ctx, userID, 100, 0)
	if listErr != nil {
		return 0, nil, apperrors.New(apperrors.ErrInternal, "일일 요약 조회 실패")
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	mealBreakdown = map[string]float64{
		"breakfast": 0,
		"lunch":     0,
		"dinner":    0,
		"snack":     0,
	}

	for _, a := range analyses {
		if a.CreatedAt.Before(today) {
			continue
		}
		if a.Status != AnalysisStatusCompleted {
			continue
		}
		totalKcal += a.TotalCalorieKcal
		if a.MealType != "" {
			mealBreakdown[a.MealType] += a.TotalCalorieKcal
		}
	}

	return totalKcal, mealBreakdown, nil
}

// ============================================================================
// 한국 음식 영양성분 룩업 테이블 (국가표준식품성분표 기반)
// ============================================================================

// koreanFoodDB는 한국 음식 영양성분 데이터베이스입니다.
// 1인분 기준. 출처: 국립농업과학원 국가표준식품성분표
var koreanFoodDB = map[string]FoodItem{
	// 밥류
	"현미밥":  {Name: "현미밥", Confidence: 0.95, CalorieKcal: 310, PortionG: 210, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 67.0, Unit: "g", DV: 0.22}, {Name: "단백질", Amount: 6.5, Unit: "g", DV: 0.13}, {Name: "지방", Amount: 1.8, Unit: "g", DV: 0.03}, {Name: "식이섬유", Amount: 3.0, Unit: "g", DV: 0.12}}},
	"쌀밥":   {Name: "쌀밥", Confidence: 0.96, CalorieKcal: 300, PortionG: 210, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 66.0, Unit: "g", DV: 0.22}, {Name: "단백질", Amount: 5.5, Unit: "g", DV: 0.11}, {Name: "지방", Amount: 0.5, Unit: "g", DV: 0.01}}},
	"잡곡밥":  {Name: "잡곡밥", Confidence: 0.94, CalorieKcal: 295, PortionG: 210, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 63.0, Unit: "g", DV: 0.21}, {Name: "단백질", Amount: 7.0, Unit: "g", DV: 0.14}, {Name: "지방", Amount: 2.0, Unit: "g", DV: 0.03}, {Name: "식이섬유", Amount: 4.0, Unit: "g", DV: 0.16}}},
	// 찌개/국류
	"김치찌개": {Name: "김치찌개", Confidence: 0.92, CalorieKcal: 180, PortionG: 300, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 12.0, Unit: "g", DV: 0.04}, {Name: "단백질", Amount: 15.0, Unit: "g", DV: 0.30}, {Name: "지방", Amount: 8.0, Unit: "g", DV: 0.12}, {Name: "나트륨", Amount: 1200.0, Unit: "mg", DV: 0.52}}},
	"된장찌개": {Name: "된장찌개", Confidence: 0.91, CalorieKcal: 120, PortionG: 300, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 10.0, Unit: "g", DV: 0.03}, {Name: "단백질", Amount: 8.0, Unit: "g", DV: 0.16}, {Name: "지방", Amount: 5.0, Unit: "g", DV: 0.08}, {Name: "나트륨", Amount: 1100.0, Unit: "mg", DV: 0.48}}},
	"순두부찌개": {Name: "순두부찌개", Confidence: 0.90, CalorieKcal: 150, PortionG: 350, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 8.0, Unit: "g", DV: 0.03}, {Name: "단백질", Amount: 12.0, Unit: "g", DV: 0.24}, {Name: "지방", Amount: 9.0, Unit: "g", DV: 0.14}, {Name: "나트륨", Amount: 950.0, Unit: "mg", DV: 0.41}}},
	"미역국":  {Name: "미역국", Confidence: 0.93, CalorieKcal: 90, PortionG: 300, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 3.0, Unit: "g", DV: 0.01}, {Name: "단백질", Amount: 10.0, Unit: "g", DV: 0.20}, {Name: "지방", Amount: 4.0, Unit: "g", DV: 0.06}, {Name: "나트륨", Amount: 900.0, Unit: "mg", DV: 0.39}}},
	"떡국":   {Name: "떡국", Confidence: 0.91, CalorieKcal: 380, PortionG: 400, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 52.0, Unit: "g", DV: 0.17}, {Name: "단백질", Amount: 15.0, Unit: "g", DV: 0.30}, {Name: "지방", Amount: 10.0, Unit: "g", DV: 0.15}, {Name: "나트륨", Amount: 1300.0, Unit: "mg", DV: 0.57}}},
	// 반찬류
	"불고기":  {Name: "불고기", Confidence: 0.93, CalorieKcal: 250, PortionG: 150, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 12.0, Unit: "g", DV: 0.04}, {Name: "단백질", Amount: 22.0, Unit: "g", DV: 0.44}, {Name: "지방", Amount: 13.0, Unit: "g", DV: 0.20}}},
	"갈비찜":  {Name: "갈비찜", Confidence: 0.90, CalorieKcal: 320, PortionG: 200, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 15.0, Unit: "g", DV: 0.05}, {Name: "단백질", Amount: 25.0, Unit: "g", DV: 0.50}, {Name: "지방", Amount: 18.0, Unit: "g", DV: 0.28}}},
	"제육볶음": {Name: "제육볶음", Confidence: 0.92, CalorieKcal: 280, PortionG: 150, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 10.0, Unit: "g", DV: 0.03}, {Name: "단백질", Amount: 20.0, Unit: "g", DV: 0.40}, {Name: "지방", Amount: 18.0, Unit: "g", DV: 0.28}}},
	"잡채":   {Name: "잡채", Confidence: 0.91, CalorieKcal: 220, PortionG: 150, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 30.0, Unit: "g", DV: 0.10}, {Name: "단백질", Amount: 8.0, Unit: "g", DV: 0.16}, {Name: "지방", Amount: 8.0, Unit: "g", DV: 0.12}}},
	"계란말이": {Name: "계란말이", Confidence: 0.94, CalorieKcal: 150, PortionG: 120, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 2.0, Unit: "g", DV: 0.01}, {Name: "단백질", Amount: 12.0, Unit: "g", DV: 0.24}, {Name: "지방", Amount: 10.0, Unit: "g", DV: 0.15}}},
	"김치":   {Name: "김치", Confidence: 0.96, CalorieKcal: 20, PortionG: 50, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 3.0, Unit: "g", DV: 0.01}, {Name: "단백질", Amount: 1.5, Unit: "g", DV: 0.03}, {Name: "나트륨", Amount: 350.0, Unit: "mg", DV: 0.15}, {Name: "식이섬유", Amount: 1.5, Unit: "g", DV: 0.06}}},
	// 면류
	"비빔밥":  {Name: "비빔밥", Confidence: 0.92, CalorieKcal: 550, PortionG: 400, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 75.0, Unit: "g", DV: 0.25}, {Name: "단백질", Amount: 20.0, Unit: "g", DV: 0.40}, {Name: "지방", Amount: 18.0, Unit: "g", DV: 0.28}}},
	"라면":   {Name: "라면", Confidence: 0.95, CalorieKcal: 500, PortionG: 550, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 70.0, Unit: "g", DV: 0.23}, {Name: "단백질", Amount: 10.0, Unit: "g", DV: 0.20}, {Name: "지방", Amount: 18.0, Unit: "g", DV: 0.28}, {Name: "나트륨", Amount: 1800.0, Unit: "mg", DV: 0.78}}},
	"냉면":   {Name: "냉면", Confidence: 0.91, CalorieKcal: 430, PortionG: 500, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 80.0, Unit: "g", DV: 0.27}, {Name: "단백질", Amount: 12.0, Unit: "g", DV: 0.24}, {Name: "지방", Amount: 5.0, Unit: "g", DV: 0.08}}},
	"칼국수":  {Name: "칼국수", Confidence: 0.90, CalorieKcal: 450, PortionG: 500, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 65.0, Unit: "g", DV: 0.22}, {Name: "단백질", Amount: 15.0, Unit: "g", DV: 0.30}, {Name: "지방", Amount: 12.0, Unit: "g", DV: 0.18}}},
	// 기타
	"삼겹살":  {Name: "삼겹살", Confidence: 0.94, CalorieKcal: 330, PortionG: 150, Nutrients: []NutrientInfo{{Name: "단백질", Amount: 17.0, Unit: "g", DV: 0.34}, {Name: "지방", Amount: 28.0, Unit: "g", DV: 0.43}}},
	"치킨":   {Name: "치킨", Confidence: 0.93, CalorieKcal: 350, PortionG: 200, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 15.0, Unit: "g", DV: 0.05}, {Name: "단백질", Amount: 25.0, Unit: "g", DV: 0.50}, {Name: "지방", Amount: 20.0, Unit: "g", DV: 0.31}}},
	"피자":   {Name: "피자", Confidence: 0.94, CalorieKcal: 270, PortionG: 120, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 30.0, Unit: "g", DV: 0.10}, {Name: "단백질", Amount: 12.0, Unit: "g", DV: 0.24}, {Name: "지방", Amount: 12.0, Unit: "g", DV: 0.18}}},
	"떡볶이":  {Name: "떡볶이", Confidence: 0.93, CalorieKcal: 380, PortionG: 300, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 65.0, Unit: "g", DV: 0.22}, {Name: "단백질", Amount: 8.0, Unit: "g", DV: 0.16}, {Name: "지방", Amount: 10.0, Unit: "g", DV: 0.15}, {Name: "나트륨", Amount: 1100.0, Unit: "mg", DV: 0.48}}},
	"김밥":   {Name: "김밥", Confidence: 0.94, CalorieKcal: 310, PortionG: 250, Nutrients: []NutrientInfo{{Name: "탄수화물", Amount: 45.0, Unit: "g", DV: 0.15}, {Name: "단백질", Amount: 10.0, Unit: "g", DV: 0.20}, {Name: "지방", Amount: 10.0, Unit: "g", DV: 0.15}}},
}

// lookupFoodByName는 한국 음식 DB에서 음식을 검색합니다.
func lookupFoodByName(name string) (FoodItem, bool) {
	item, ok := koreanFoodDB[name]
	return item, ok
}

// fallbackAnalysis는 AI 분석기 없이 기본 음식 조합을 반환합니다.
// 향후 Cloud Vision API 또는 LLM 연동 시 교체됩니다.
func (s *VisionService) fallbackAnalysis(imageURL string) []FoodItem {
	_ = imageURL
	// 기본 한식 1인분 구성 (밥 + 국 + 반찬)
	items := []FoodItem{
		koreanFoodDB["쌀밥"],
		koreanFoodDB["된장찌개"],
		koreanFoodDB["김치"],
	}
	return items
}

// LogMeal은 이미지 분석 없이 음식 이름으로 식사를 기록합니다.
func (s *VisionService) LogMeal(ctx context.Context, userID string, foodNames []string, mealType string) (*FoodAnalysis, error) {
	if userID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "user_id는 필수입니다")
	}
	if len(foodNames) == 0 {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "음식 이름이 비어 있습니다")
	}
	if mealType == "" {
		mealType = "snack"
	}
	if !ValidMealTypes[mealType] {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "지원하지 않는 식사 타입: "+mealType)
	}

	var items []FoodItem
	var unknownFoods []string
	for _, name := range foodNames {
		if item, ok := LookupFood(name); ok {
			items = append(items, item)
		} else {
			unknownFoods = append(unknownFoods, name)
		}
	}

	if len(items) == 0 {
		return nil, apperrors.New(apperrors.ErrNotFound, "인식할 수 있는 음식이 없습니다: "+fmt.Sprintf("%v", unknownFoods))
	}

	var totalCal float64
	for _, item := range items {
		totalCal += item.CalorieKcal
	}

	now := time.Now().UTC()
	analysis := &FoodAnalysis{
		ID:               uuid.New().String(),
		UserID:           userID,
		ImageURL:         "",
		Status:           AnalysisStatusCompleted,
		TotalCalorieKcal: totalCal,
		FoodItems:        items,
		MealType:         mealType,
		AnalyzedAt:       &now,
		CreatedAt:        now,
	}

	if err := s.repo.Save(ctx, analysis); err != nil {
		return nil, apperrors.New(apperrors.ErrInternal, "식사 기록 저장 실패")
	}

	s.logger.Info("식사 기록 완료",
		zap.String("user_id", userID),
		zap.String("meal_type", mealType),
		zap.Float64("total_kcal", totalCal),
		zap.Int("food_items", len(items)),
	)

	return analysis, nil
}

// DeleteAnalysis는 음식 분석 결과를 삭제합니다.
func (s *VisionService) DeleteAnalysis(ctx context.Context, analysisID string) error {
	if analysisID == "" {
		return apperrors.New(apperrors.ErrInvalidInput, "analysis_id는 필수입니다")
	}
	existing, err := s.repo.FindByID(ctx, analysisID)
	if err != nil {
		return apperrors.New(apperrors.ErrInternal, "분석 결과 조회 실패")
	}
	if existing == nil {
		return apperrors.New(apperrors.ErrNotFound, "분석 결과를 찾을 수 없습니다")
	}
	return s.repo.Delete(ctx, analysisID)
}

// LookupFood는 한국 음식 DB에서 음식을 검색합니다 (public).
func LookupFood(name string) (FoodItem, bool) {
	return lookupFoodByName(name)
}

// FoodAnalysisStatusToString은 상태를 문자열로 변환합니다.
func FoodAnalysisStatusToString(s FoodAnalysisStatus) string {
	switch s {
	case AnalysisStatusPending:
		return "pending"
	case AnalysisStatusProcessing:
		return "processing"
	case AnalysisStatusCompleted:
		return "completed"
	case AnalysisStatusFailed:
		return "failed"
	default:
		return fmt.Sprintf("unknown(%d)", s)
	}
}
