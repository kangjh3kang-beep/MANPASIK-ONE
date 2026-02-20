// Package handler는 nlp-service의 gRPC 핸들러입니다.
package handler

import (
	"context"

	"github.com/manpasik/backend/services/nlp-service/internal/service"
)

// NLPHandler는 NLPService를 래핑하는 핸들러입니다.
type NLPHandler struct {
	svc *service.NLPService
}

// NewNLPHandler는 NLPHandler를 생성합니다.
func NewNLPHandler(svc *service.NLPService) *NLPHandler {
	return &NLPHandler{svc: svc}
}

// ParseHealthQuery는 건강 질의 파싱을 위임합니다.
func (h *NLPHandler) ParseHealthQuery(ctx context.Context, userID, text string) (*service.HealthQuery, error) {
	return h.svc.ParseHealthQuery(ctx, userID, text)
}

// ExtractSymptoms는 증상 추출을 위임합니다.
func (h *NLPHandler) ExtractSymptoms(ctx context.Context, text string) (*service.SymptomExtraction, error) {
	return h.svc.ExtractSymptoms(ctx, text)
}

// GetSuggestions는 제안 조회를 위임합니다.
func (h *NLPHandler) GetSuggestions(ctx context.Context, queryID string) ([]service.Suggestion, error) {
	return h.svc.GetSuggestions(ctx, queryID)
}
