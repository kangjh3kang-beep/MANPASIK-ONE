// Package handler는 translation-service의 gRPC 핸들러입니다.
package handler

import (
	"context"

	"github.com/manpasik/backend/services/translation-service/internal/service"
	apperrors "github.com/manpasik/backend/shared/errors"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TranslationHandler는 TranslationService gRPC 서버를 구현합니다.
type TranslationHandler struct {
	v1.UnimplementedTranslationServiceServer
	svc *service.TranslationService
	log *zap.Logger
}

// NewTranslationHandler는 TranslationHandler를 생성합니다.
func NewTranslationHandler(svc *service.TranslationService, log *zap.Logger) *TranslationHandler {
	return &TranslationHandler{svc: svc, log: log}
}

// TranslateText는 텍스트 번역 RPC입니다.
func (h *TranslationHandler) TranslateText(ctx context.Context, req *v1.TranslateTextRequest) (*v1.TranslateTextResponse, error) {
	if req == nil || req.Text == "" {
		return nil, status.Error(codes.InvalidArgument, "text는 필수입니다")
	}
	if req.TargetLanguage == "" {
		return nil, status.Error(codes.InvalidArgument, "target_language는 필수입니다")
	}

	translated, sourceLang, confidence, _, err := h.svc.TranslateText(
		ctx, req.Text, req.SourceLanguage, req.TargetLanguage, req.Context, "",
	)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.TranslateTextResponse{
		TranslatedText: translated,
		SourceLanguage: sourceLang,
		TargetLanguage: req.TargetLanguage,
		Confidence:     float64(confidence),
		IsMedicalTerm:  req.IsMedical,
		OriginalText:   req.Text,
	}, nil
}

// DetectLanguage는 언어 감지 RPC입니다.
func (h *TranslationHandler) DetectLanguage(ctx context.Context, req *v1.DetectLanguageRequest) (*v1.DetectLanguageResponse, error) {
	if req == nil || req.Text == "" {
		return nil, status.Error(codes.InvalidArgument, "text는 필수입니다")
	}

	detected, _, err := h.svc.DetectLanguage(ctx, req.Text)
	if err != nil {
		return nil, toGRPC(err)
	}

	var pbDetected []*v1.DetectedLanguage
	for _, d := range detected {
		pbDetected = append(pbDetected, &v1.DetectedLanguage{
			LanguageCode: d.Code,
			LanguageName: d.Name,
			Confidence:   float64(d.Confidence),
		})
	}

	return &v1.DetectLanguageResponse{
		Languages: pbDetected,
	}, nil
}

// ListSupportedLanguages는 지원 언어 목록 RPC입니다.
func (h *TranslationHandler) ListSupportedLanguages(ctx context.Context, _ *v1.ListSupportedLanguagesRequest) (*v1.ListSupportedLanguagesResponse, error) {
	languages, err := h.svc.ListSupportedLanguages(ctx)
	if err != nil {
		return nil, toGRPC(err)
	}

	var pb []*v1.SupportedLanguage
	for _, l := range languages {
		pb = append(pb, &v1.SupportedLanguage{
			LanguageCode:    l.Code,
			LanguageName:    l.Name,
			NativeName:      l.NativeName,
			SupportsMedical: l.SupportsMedical,
		})
	}

	return &v1.ListSupportedLanguagesResponse{
		Languages: pb,
	}, nil
}

// TranslateBatch는 배치 번역 RPC입니다.
func (h *TranslationHandler) TranslateBatch(ctx context.Context, req *v1.TranslateBatchRequest) (*v1.TranslateBatchResponse, error) {
	if req == nil || len(req.Texts) == 0 {
		return nil, status.Error(codes.InvalidArgument, "texts는 필수입니다")
	}

	translated, sourceLang, _, err := h.svc.TranslateBatch(
		ctx, req.Texts, req.SourceLanguage, req.TargetLanguage, "",
	)
	if err != nil {
		return nil, toGRPC(err)
	}

	var translations []*v1.TranslateTextResponse
	for i, t := range translated {
		translations = append(translations, &v1.TranslateTextResponse{
			TranslatedText: t,
			SourceLanguage: sourceLang,
			TargetLanguage: req.TargetLanguage,
			IsMedicalTerm:  req.IsMedical,
			OriginalText:   req.Texts[i],
		})
	}

	return &v1.TranslateBatchResponse{
		Translations: translations,
	}, nil
}

// GetTranslationHistory는 번역 이력 조회 RPC입니다.
func (h *TranslationHandler) GetTranslationHistory(ctx context.Context, req *v1.GetTranslationHistoryRequest) (*v1.GetTranslationHistoryResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	records, total, err := h.svc.GetTranslationHistory(ctx, req.UserId, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, toGRPC(err)
	}

	var pb []*v1.TranslationRecord
	for _, r := range records {
		pb = append(pb, &v1.TranslationRecord{
			RecordId:       r.ID,
			UserId:         req.UserId,
			SourceText:     r.SourceText,
			TranslatedText: r.TranslatedText,
			SourceLanguage: r.SourceLanguage,
			TargetLanguage: r.TargetLanguage,
			IsMedical:      false,
			CreatedAt:      timestamppb.New(r.CreatedAt),
		})
	}

	return &v1.GetTranslationHistoryResponse{
		Records:    pb,
		TotalCount: int32(total),
	}, nil
}

// GetTranslationUsage는 번역 사용량 조회 RPC입니다.
func (h *TranslationHandler) GetTranslationUsage(ctx context.Context, req *v1.GetTranslationUsageRequest) (*v1.GetTranslationUsageResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	usage, err := h.svc.GetTranslationUsage(ctx, req.UserId)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.GetTranslationUsageResponse{
		UserId:              req.UserId,
		TotalTranslations:   int32(usage.TotalRequests),
		MonthlyTranslations: int32(usage.MonthlyRequests),
		MonthlyLimit:        int32(usage.MonthlyLimit),
	}, nil
}

// TranslateRealtime은 실시간 번역 RPC입니다.
func (h *TranslationHandler) TranslateRealtime(ctx context.Context, req *v1.TranslateRealtimeRequest) (*v1.TranslateRealtimeResponse, error) {
	if req == nil || req.Text == "" {
		return nil, status.Error(codes.InvalidArgument, "text는 필수입니다")
	}
	if req.TargetLanguage == "" {
		return nil, status.Error(codes.InvalidArgument, "target_language는 필수입니다")
	}

	translated, sourceLang, confidence, latencyMs, err := h.svc.TranslateRealtime(
		ctx, req.Text, req.SourceLanguage, req.TargetLanguage, req.Context, req.SessionId, req.IsMedical,
	)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.TranslateRealtimeResponse{
		TranslatedText: translated,
		SourceLanguage: sourceLang,
		TargetLanguage: req.TargetLanguage,
		Confidence:     float64(confidence),
		IsMedicalTerm:  req.IsMedical,
		LatencyMs:      latencyMs,
	}, nil
}

func toGRPC(err error) error {
	if appErr, ok := err.(*apperrors.AppError); ok {
		return appErr.ToGRPC()
	}
	return status.Error(codes.Internal, err.Error())
}
