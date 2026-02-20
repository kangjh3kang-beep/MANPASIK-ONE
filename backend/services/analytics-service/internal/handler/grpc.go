// Package handler는 analytics-service의 핸들러를 정의합니다.
// 현재 gRPC proto 정의가 없으므로 서비스 래퍼 형태로 제공합니다.
package handler

import (
	"context"

	"github.com/manpasik/backend/services/analytics-service/internal/service"
)

// AnalyticsHandler는 분석 서비스 핸들러입니다.
type AnalyticsHandler struct {
	svc *service.AnalyticsService
}

// NewAnalyticsHandler는 새 AnalyticsHandler를 생성합니다.
func NewAnalyticsHandler(svc *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{svc: svc}
}

// TrackEvent는 분석 이벤트를 기록합니다.
func (h *AnalyticsHandler) TrackEvent(ctx context.Context, userID, eventType string, props map[string]string) (string, error) {
	return h.svc.TrackEvent(ctx, userID, eventType, props)
}

// GetUserAnalytics는 사용자 분석 요약을 반환합니다.
func (h *AnalyticsHandler) GetUserAnalytics(ctx context.Context, userID string) (*service.UserAnalytics, error) {
	return h.svc.GetUserAnalytics(ctx, userID)
}

// GetDailyStats는 일별 분석 통계를 반환합니다.
func (h *AnalyticsHandler) GetDailyStats(ctx context.Context, date string) (*service.DailyStats, error) {
	return h.svc.GetDailyStats(ctx, date)
}
