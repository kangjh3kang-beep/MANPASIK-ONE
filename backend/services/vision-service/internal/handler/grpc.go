// Package handler는 vision-service의 gRPC 핸들러를 구현합니다.
//
// NOTE: VisionService Proto 정의가 manpasik.proto에 추가된 후
// v1.RegisterVisionServiceServer로 연결됩니다.
// 현재는 서비스 레이어만 사용 가능합니다.
package handler

import (
	"github.com/manpasik/backend/services/vision-service/internal/service"
	"go.uber.org/zap"
)

// VisionHandler는 gRPC VisionService 핸들러입니다.
// Proto 정의 추가 후 gRPC 메서드를 구현합니다.
type VisionHandler struct {
	svc *service.VisionService
	log *zap.Logger
}

// NewVisionHandler는 VisionHandler를 생성합니다.
func NewVisionHandler(svc *service.VisionService, log *zap.Logger) *VisionHandler {
	return &VisionHandler{svc: svc, log: log}
}

// TODO: Proto 확장 후 구현
// - AnalyzeFood(ctx, req) → v1.FoodAnalysis
// - GetFoodAnalysis(ctx, req) → v1.FoodAnalysis
// - ListFoodAnalyses(ctx, req) → v1.ListFoodAnalysesResponse
// - GetDailyNutritionSummary(ctx, req) → v1.DailyNutritionSummary
