// Package handler는 vision-service의 gRPC 핸들러를 구현합니다.
package handler

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/manpasik/backend/services/vision-service/internal/service"
	apperrors "github.com/manpasik/backend/shared/errors"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// VisionHandler는 gRPC VisionService 핸들러입니다.
type VisionHandler struct {
	v1.UnimplementedVisionServiceServer
	svc      *service.VisionService
	log      *zap.Logger
	mealLogs []mealLogEntry // 간단한 인메모리 식사 로그
}

type mealLogEntry struct {
	id          string
	userID      string
	mealType    string
	description string
	calories    float64
	analysisID  string
	loggedAt    time.Time
}

// NewVisionHandler는 VisionHandler를 생성합니다.
func NewVisionHandler(svc *service.VisionService, log *zap.Logger) *VisionHandler {
	return &VisionHandler{svc: svc, log: log}
}

func (h *VisionHandler) AnalyzeFood(ctx context.Context, req *v1.AnalyzeFoodRequest) (*v1.FoodAnalysisResult, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}
	analysis, err := h.svc.AnalyzeFood(ctx, req.UserId, req.ImageUrl, req.MealType)
	if err != nil {
		return nil, toGRPC(err)
	}
	return analysisToProto(analysis), nil
}

func (h *VisionHandler) GetFoodAnalysis(ctx context.Context, req *v1.GetFoodAnalysisRequest) (*v1.FoodAnalysisResult, error) {
	if req == nil || req.AnalysisId == "" {
		return nil, status.Error(codes.InvalidArgument, "analysis_id는 필수입니다")
	}
	analysis, err := h.svc.GetAnalysis(ctx, req.AnalysisId)
	if err != nil {
		return nil, toGRPC(err)
	}
	return analysisToProto(analysis), nil
}

func (h *VisionHandler) ListFoodAnalyses(ctx context.Context, req *v1.ListFoodAnalysesRequest) (*v1.ListFoodAnalysesResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}
	analyses, total, err := h.svc.ListAnalyses(ctx, req.UserId, req.Limit, req.Offset)
	if err != nil {
		return nil, toGRPC(err)
	}
	results := make([]*v1.FoodAnalysisResult, 0, len(analyses))
	for _, a := range analyses {
		results = append(results, analysisToProto(a))
	}
	return &v1.ListFoodAnalysesResponse{Analyses: results, Total: total}, nil
}

func (h *VisionHandler) GetDailyNutritionSummary(ctx context.Context, req *v1.GetDailyNutritionSummaryRequest) (*v1.DailyNutritionSummary, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}
	totalKcal, breakdown, err := h.svc.GetDailySummary(ctx, req.UserId)
	if err != nil {
		return nil, toGRPC(err)
	}

	meals := make([]*v1.MealLog, 0)
	for mealType, cal := range breakdown {
		if cal > 0 {
			meals = append(meals, &v1.MealLog{
				MealType:      mealType,
				TotalCalories: cal,
			})
		}
	}

	return &v1.DailyNutritionSummary{
		UserId:        req.UserId,
		Date:          time.Now().UTC().Format("2006-01-02"),
		TotalCalories: totalKcal,
		MealCount:     int32(len(meals)),
		Meals:         meals,
	}, nil
}

func (h *VisionHandler) LogMeal(_ context.Context, req *v1.LogMealRequest) (*v1.MealLog, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	now := time.Now().UTC()

	// Items의 칼로리 합산
	var totalCal float64
	for _, item := range req.Items {
		totalCal += item.Calories
	}

	entry := mealLogEntry{
		id:          uuid.New().String(),
		userID:      req.UserId,
		mealType:    req.MealType,
		description: req.Description,
		calories:    totalCal,
		analysisID:  req.AnalysisId,
		loggedAt:    now,
	}
	h.mealLogs = append(h.mealLogs, entry)

	return &v1.MealLog{
		Id:            entry.id,
		UserId:        req.UserId,
		MealType:      req.MealType,
		Description:   req.Description,
		Items:         req.Items,
		TotalCalories: totalCal,
		AnalysisId:    req.AnalysisId,
		LoggedAt:      timestamppb.New(now),
	}, nil
}

func (h *VisionHandler) GetMealHistory(_ context.Context, req *v1.GetMealHistoryRequest) (*v1.GetMealHistoryResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	meals := make([]*v1.MealLog, 0)
	for _, entry := range h.mealLogs {
		if entry.userID != req.UserId {
			continue
		}
		meals = append(meals, &v1.MealLog{
			Id:            entry.id,
			UserId:        entry.userID,
			MealType:      entry.mealType,
			Description:   entry.description,
			TotalCalories: entry.calories,
			AnalysisId:    entry.analysisID,
			LoggedAt:      timestamppb.New(entry.loggedAt),
		})
	}

	return &v1.GetMealHistoryResponse{Meals: meals, Total: int32(len(meals))}, nil
}

func analysisToProto(a *service.FoodAnalysis) *v1.FoodAnalysisResult {
	items := make([]*v1.FoodItem, 0, len(a.FoodItems))
	var totalCarbs, totalProtein, totalFat float64
	for _, item := range a.FoodItems {
		fi := &v1.FoodItem{
			Name:       item.Name,
			Calories:   item.CalorieKcal,
			PortionG:   item.PortionG,
			Confidence: item.Confidence,
		}
		for _, n := range item.Nutrients {
			switch n.Name {
			case "탄수화물":
				fi.CarbsG = n.Amount
				totalCarbs += n.Amount
			case "단백질":
				fi.ProteinG = n.Amount
				totalProtein += n.Amount
			case "지방":
				fi.FatG = n.Amount
				totalFat += n.Amount
			}
		}
		items = append(items, fi)
	}

	result := &v1.FoodAnalysisResult{
		Id:            a.ID,
		UserId:        a.UserID,
		FoodItems:     items,
		TotalCalories: a.TotalCalorieKcal,
		TotalCarbsG:   totalCarbs,
		TotalProteinG: totalProtein,
		TotalFatG:     totalFat,
		MealType:      a.MealType,
		ImageUrl:      a.ImageURL,
	}
	if a.AnalyzedAt != nil {
		result.AnalyzedAt = timestamppb.New(*a.AnalyzedAt)
	}
	return result
}

func toGRPC(err error) error {
	if err == nil {
		return nil
	}
	if ae, ok := err.(*apperrors.AppError); ok {
		return ae.ToGRPC()
	}
	return status.Error(codes.Internal, "내부 오류가 발생했습니다")
}
