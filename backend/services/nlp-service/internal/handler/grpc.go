// Package handler는 nlp-service의 gRPC 핸들러를 구현합니다.
package handler

import (
	"context"

	"github.com/manpasik/backend/services/nlp-service/internal/service"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// NLPHandler는 gRPC NLPService 핸들러입니다.
type NLPHandler struct {
	v1.UnimplementedNLPServiceServer
	svc *service.NLPService
	log *zap.Logger
}

// NewNLPHandler는 NLPHandler를 생성합니다.
func NewNLPHandler(svc *service.NLPService, log *zap.Logger) *NLPHandler {
	return &NLPHandler{svc: svc, log: log}
}

// ParseHealthQuery는 건강 질의를 파싱하는 RPC입니다.
func (h *NLPHandler) ParseHealthQuery(ctx context.Context, req *v1.ParseHealthQueryRequest) (*v1.ParseHealthQueryResponse, error) {
	if req == nil || req.Text == "" {
		return nil, status.Error(codes.InvalidArgument, "text는 필수입니다")
	}
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	query, err := h.svc.ParseHealthQuery(ctx, req.UserId, req.Text)
	if err != nil {
		return nil, toGRPC(err)
	}

	// 제안 자동 생성
	_, sugErr := h.svc.GenerateSuggestions(ctx, query)
	if sugErr != nil {
		h.log.Warn("제안 생성 실패", zap.Error(sugErr))
	}

	return &v1.ParseHealthQueryResponse{
		Id:         query.ID,
		UserId:     query.UserID,
		RawText:    query.RawText,
		Intent:     query.Intent,
		Entities:   query.Entities,
		Confidence: query.Confidence,
		Urgency:    query.Urgency,
	}, nil
}

// ExtractSymptoms는 텍스트에서 증상을 추출하는 RPC입니다.
func (h *NLPHandler) ExtractSymptoms(ctx context.Context, req *v1.ExtractSymptomsRequest) (*v1.ExtractSymptomsResponse, error) {
	if req == nil || req.Text == "" {
		return nil, status.Error(codes.InvalidArgument, "text는 필수입니다")
	}

	extraction, err := h.svc.ExtractSymptoms(ctx, req.Text)
	if err != nil {
		return nil, toGRPC(err)
	}

	pbSymptoms := make([]*v1.NLPSymptom, 0, len(extraction.Symptoms))
	for _, s := range extraction.Symptoms {
		pbSymptoms = append(pbSymptoms, &v1.NLPSymptom{
			Name:       s.Name,
			Severity:   s.Severity,
			BodyPart:   s.BodyPart,
			Confidence: s.Confidence,
		})
	}

	return &v1.ExtractSymptomsResponse{
		Id:          extraction.ID,
		Text:        extraction.Text,
		Symptoms:    pbSymptoms,
		Urgency:     extraction.Urgency,
		ProcessedAt: timestamppb.New(extraction.ProcessedAt),
	}, nil
}

// GetSuggestions는 질의에 대한 제안 목록을 반환하는 RPC입니다.
func (h *NLPHandler) GetSuggestions(ctx context.Context, req *v1.GetSuggestionsRequest) (*v1.GetSuggestionsResponse, error) {
	if req == nil || req.QueryId == "" {
		return nil, status.Error(codes.InvalidArgument, "query_id는 필수입니다")
	}

	suggestions, err := h.svc.GetSuggestions(ctx, req.QueryId)
	if err != nil {
		return nil, toGRPC(err)
	}

	pbSuggestions := make([]*v1.NLPSuggestion, 0, len(suggestions))
	for _, s := range suggestions {
		pbSuggestions = append(pbSuggestions, &v1.NLPSuggestion{
			Id:       s.ID,
			QueryId:  s.QueryID,
			Text:     s.Text,
			Category: s.Category,
			Priority: int32(s.Priority),
		})
	}

	return &v1.GetSuggestionsResponse{
		Suggestions: pbSuggestions,
	}, nil
}

// AnalyzeKoreanText는 한국어 텍스트를 종합 분석하는 RPC입니다.
func (h *NLPHandler) AnalyzeKoreanText(_ context.Context, req *v1.AnalyzeKoreanTextRequest) (*v1.AnalyzeKoreanTextResponse, error) {
	if req == nil || req.Text == "" {
		return nil, status.Error(codes.InvalidArgument, "text는 필수입니다")
	}

	analysis := h.svc.AnalyzeKorean(req.Text)

	return &v1.AnalyzeKoreanTextResponse{
		Text:       analysis.Text,
		Intent:     analysis.Intent,
		Confidence: analysis.Confidence,
		Severity:   int32(analysis.Severity),
		Keywords:   analysis.Keywords,
		Urgency:    service.SeverityToUrgency(analysis.Severity),
	}, nil
}

func toGRPC(err error) error {
	if err == nil {
		return nil
	}
	return status.Error(codes.Internal, err.Error())
}
