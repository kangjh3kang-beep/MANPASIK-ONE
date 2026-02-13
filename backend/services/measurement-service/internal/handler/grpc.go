// Package handler는 measurement-service의 gRPC 핸들러입니다.
package handler

import (
	"context"
	"strings"
	"time"

	"github.com/manpasik/backend/services/measurement-service/internal/service"
	apperrors "github.com/manpasik/backend/shared/errors"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// MeasurementHandler는 MeasurementService gRPC 서버를 구현합니다.
type MeasurementHandler struct {
	v1.UnimplementedMeasurementServiceServer
	svc *service.MeasurementService
	log *zap.Logger
}

// NewMeasurementHandler는 MeasurementHandler를 생성합니다.
func NewMeasurementHandler(svc *service.MeasurementService, log *zap.Logger) *MeasurementHandler {
	return &MeasurementHandler{svc: svc, log: log}
}

// StartSession은 측정 세션 시작 RPC입니다.
func (h *MeasurementHandler) StartSession(ctx context.Context, req *v1.StartSessionRequest) (*v1.StartSessionResponse, error) {
	if req == nil || req.DeviceId == "" || req.CartridgeId == "" || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "device_id, cartridge_id, user_id는 필수입니다")
	}

	session, err := h.svc.StartSession(ctx, req.DeviceId, req.CartridgeId, req.UserId)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.StartSessionResponse{
		SessionId: session.ID,
		StartedAt: timestamppb.New(session.StartedAt),
	}, nil
}

// EndSession은 측정 세션 종료 RPC입니다.
func (h *MeasurementHandler) EndSession(ctx context.Context, req *v1.EndSessionRequest) (*v1.EndSessionResponse, error) {
	if req == nil || req.SessionId == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id는 필수입니다")
	}

	result, err := h.svc.EndSession(ctx, req.SessionId)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.EndSessionResponse{
		SessionId:         result.SessionID,
		TotalMeasurements: int32(result.TotalMeasurements),
		EndedAt:           timestamppb.New(result.EndedAt),
	}, nil
}

// GetMeasurementHistory는 측정 기록 조회 RPC입니다.
func (h *MeasurementHandler) GetMeasurementHistory(ctx context.Context, req *v1.GetHistoryRequest) (*v1.GetHistoryResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	var start, end time.Time
	if req.StartTime != nil {
		start = req.StartTime.AsTime()
	}
	if req.EndTime != nil {
		end = req.EndTime.AsTime()
	} else {
		end = time.Now().UTC()
	}

	limit := int(req.Limit)
	offset := int(req.Offset)

	summaries, total, err := h.svc.GetHistory(ctx, req.UserId, start, end, limit, offset)
	if err != nil {
		return nil, toGRPC(err)
	}

	var pbSummaries []*v1.MeasurementSummary
	for _, s := range summaries {
		pbSummaries = append(pbSummaries, &v1.MeasurementSummary{
			SessionId:     s.SessionID,
			CartridgeType: s.CartridgeType,
			PrimaryValue:  s.PrimaryValue,
			Unit:          s.Unit,
			MeasuredAt:    timestamppb.New(s.MeasuredAt),
		})
	}

	return &v1.GetHistoryResponse{
		Measurements: pbSummaries,
		TotalCount:   int32(total),
	}, nil
}

// ExportSingleMeasurement는 단일 측정 세션의 FHIR 내보내기 RPC입니다.
func (h *MeasurementHandler) ExportSingleMeasurement(ctx context.Context, req *v1.ExportSingleMeasurementRequest) (*v1.ExportFHIRResponse, error) {
	if req == nil || req.SessionId == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id는 필수입니다")
	}

	bundleJSON, err := h.svc.ExportSingleMeasurement(ctx, req.SessionId)
	if err != nil {
		return nil, toGRPC(err)
	}

	// Count resources by counting "resourceType" occurrences (rough estimate)
	count := strings.Count(bundleJSON, "\"resourceType\"") - 1 // subtract 1 for Bundle itself
	if count < 0 {
		count = 0
	}

	return &v1.ExportFHIRResponse{
		FhirBundleJson: bundleJSON,
		ResourceCount:  int32(count),
	}, nil
}

// ExportToFHIRObservations는 사용자 전체 측정 결과의 FHIR Observation 내보내기 RPC입니다.
func (h *MeasurementHandler) ExportToFHIRObservations(ctx context.Context, req *v1.ExportToFHIRObservationsRequest) (*v1.ExportFHIRResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	bundleJSON, _, err := h.svc.ExportToFHIRObservations(ctx, req.UserId, nil, nil, nil)
	if err != nil {
		return nil, toGRPC(err)
	}

	count := strings.Count(bundleJSON, "\"resourceType\"") - 1
	if count < 0 {
		count = 0
	}

	return &v1.ExportFHIRResponse{
		FhirBundleJson: bundleJSON,
		ResourceCount:  int32(count),
	}, nil
}

// toGRPC는 AppError를 gRPC status로 변환합니다.
func toGRPC(err error) error {
	if err == nil {
		return nil
	}
	if ae, ok := err.(*apperrors.AppError); ok {
		return ae.ToGRPC()
	}
	if s, ok := status.FromError(err); ok {
		return s.Err()
	}
	return status.Error(codes.Internal, "내부 오류가 발생했습니다")
}
