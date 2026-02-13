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
		// 실제 AI 분석
		foodItems, analysisErr = s.analyzer.AnalyzeFood(ctx, imageURL)
	} else {
		// 시뮬레이션 모드 (개발/테스트용)
		foodItems = s.simulateAnalysis(imageURL)
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

// simulateAnalysis는 AI 분석기가 없을 때 시뮬레이션 결과를 반환합니다.
func (s *VisionService) simulateAnalysis(imageURL string) []FoodItem {
	_ = imageURL
	return []FoodItem{
		{
			Name:        "김치찌개",
			Confidence:  0.92,
			CalorieKcal: 180,
			PortionG:    300,
			Nutrients: []NutrientInfo{
				{Name: "탄수화물", Amount: 12.0, Unit: "g", DV: 0.04},
				{Name: "단백질", Amount: 15.0, Unit: "g", DV: 0.30},
				{Name: "지방", Amount: 8.0, Unit: "g", DV: 0.12},
				{Name: "나트륨", Amount: 1200.0, Unit: "mg", DV: 0.52},
			},
		},
		{
			Name:        "현미밥",
			Confidence:  0.95,
			CalorieKcal: 310,
			PortionG:    210,
			Nutrients: []NutrientInfo{
				{Name: "탄수화물", Amount: 67.0, Unit: "g", DV: 0.22},
				{Name: "단백질", Amount: 6.5, Unit: "g", DV: 0.13},
				{Name: "지방", Amount: 1.8, Unit: "g", DV: 0.03},
				{Name: "식이섬유", Amount: 3.0, Unit: "g", DV: 0.12},
			},
		},
	}
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
