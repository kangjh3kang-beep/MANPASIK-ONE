// Package handler는 calibration-service의 gRPC 핸들러입니다.
package handler

import (
	"context"
	"time"

	"github.com/manpasik/backend/services/calibration-service/internal/service"
	apperrors "github.com/manpasik/backend/shared/errors"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// CalibrationHandler는 CalibrationService gRPC 서버를 구현합니다.
type CalibrationHandler struct {
	v1.UnimplementedCalibrationServiceServer
	svc *service.CalibrationService
	log *zap.Logger
}

// NewCalibrationHandler는 CalibrationHandler를 생성합니다.
func NewCalibrationHandler(svc *service.CalibrationService, log *zap.Logger) *CalibrationHandler {
	return &CalibrationHandler{svc: svc, log: log}
}

// RegisterFactoryCalibration은 팩토리 보정 등록 RPC입니다.
func (h *CalibrationHandler) RegisterFactoryCalibration(ctx context.Context, req *v1.RegisterFactoryCalibrationRequest) (*v1.CalibrationRecord, error) {
	if req == nil || req.DeviceId == "" {
		return nil, status.Error(codes.InvalidArgument, "device_id는 필수입니다")
	}
	if req.CartridgeCategory <= 0 {
		return nil, status.Error(codes.InvalidArgument, "유효하지 않은 카트리지 카테고리입니다")
	}

	record, err := h.svc.RegisterFactoryCalibration(
		ctx,
		req.DeviceId,
		req.CartridgeCategory,
		req.CartridgeTypeIndex,
		req.Alpha,
		req.ChannelOffsets,
		req.ChannelGains,
		req.TempCoefficient,
		req.HumidityCoefficient,
		req.ReferenceStandard,
		req.CalibratedBy,
	)
	if err != nil {
		return nil, toGRPC(err)
	}

	return calibrationRecordToProto(record), nil
}

// PerformFieldCalibration은 현장 보정 수행 RPC입니다.
func (h *CalibrationHandler) PerformFieldCalibration(ctx context.Context, req *v1.PerformFieldCalibrationRequest) (*v1.CalibrationRecord, error) {
	if req == nil || req.DeviceId == "" {
		return nil, status.Error(codes.InvalidArgument, "device_id는 필수입니다")
	}
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}
	if len(req.ReferenceValues) == 0 || len(req.MeasuredValues) == 0 {
		return nil, status.Error(codes.InvalidArgument, "기준값과 측정값이 필요합니다")
	}

	record, err := h.svc.PerformFieldCalibration(
		ctx,
		req.DeviceId,
		req.UserId,
		req.CartridgeCategory,
		req.CartridgeTypeIndex,
		req.ReferenceValues,
		req.MeasuredValues,
		req.TemperatureC,
		req.HumidityPct,
	)
	if err != nil {
		return nil, toGRPC(err)
	}

	return calibrationRecordToProto(record), nil
}

// GetCalibration은 보정 데이터 조회 RPC입니다.
func (h *CalibrationHandler) GetCalibration(ctx context.Context, req *v1.GetCalibrationRequest) (*v1.CalibrationRecord, error) {
	if req == nil || req.DeviceId == "" {
		return nil, status.Error(codes.InvalidArgument, "device_id는 필수입니다")
	}

	record, err := h.svc.GetCalibration(ctx, req.DeviceId, req.CartridgeCategory, req.CartridgeTypeIndex)
	if err != nil {
		return nil, toGRPC(err)
	}

	return calibrationRecordToProto(record), nil
}

// ListCalibrationHistory는 보정 이력 조회 RPC입니다.
func (h *CalibrationHandler) ListCalibrationHistory(ctx context.Context, req *v1.ListCalibrationHistoryRequest) (*v1.ListCalibrationHistoryResponse, error) {
	if req == nil || req.DeviceId == "" {
		return nil, status.Error(codes.InvalidArgument, "device_id는 필수입니다")
	}

	records, total, err := h.svc.ListCalibrationHistory(ctx, req.DeviceId, req.Limit, req.Offset)
	if err != nil {
		return nil, toGRPC(err)
	}

	protoRecords := make([]*v1.CalibrationRecord, 0, len(records))
	for _, r := range records {
		protoRecords = append(protoRecords, calibrationRecordToProto(r))
	}

	return &v1.ListCalibrationHistoryResponse{
		Records:    protoRecords,
		TotalCount: total,
	}, nil
}

// CheckCalibrationStatus는 보정 상태 확인 RPC입니다.
func (h *CalibrationHandler) CheckCalibrationStatus(ctx context.Context, req *v1.CheckCalibrationStatusRequest) (*v1.CalibrationStatusResponse, error) {
	if req == nil || req.DeviceId == "" {
		return nil, status.Error(codes.InvalidArgument, "device_id는 필수입니다")
	}

	st, msg, latest, err := h.svc.CheckCalibrationStatus(ctx, req.DeviceId, req.CartridgeCategory, req.CartridgeTypeIndex)
	if err != nil {
		return nil, toGRPC(err)
	}

	resp := &v1.CalibrationStatusResponse{
		Status:   serviceStatusToProto(st),
		DeviceId: req.DeviceId,
		Message:  msg,
	}

	if latest != nil {
		resp.LastCalibratedAt = timestamppb.New(latest.CalibratedAt)
		resp.ExpiresAt = timestamppb.New(latest.ExpiresAt)
		resp.LatestRecord = calibrationRecordToProto(latest)

		// 만료까지 남은 일수 계산
		now := time.Now().UTC()
		var daysUntil int32
		if st != service.CalibrationStatusExpired {
			remaining := latest.ExpiresAt.Sub(now).Hours() / 24
			if remaining > 0 {
				daysUntil = int32(remaining)
			} else {
				daysUntil = 0
			}
		} else {
			daysUntil = 0
		}
		resp.DaysUntilExpiry = daysUntil
	}

	return resp, nil
}

// ListCalibrationModels는 보정 모델 목록 조회 RPC입니다.
func (h *CalibrationHandler) ListCalibrationModels(ctx context.Context, _ *v1.ListCalibrationModelsRequest) (*v1.ListCalibrationModelsResponse, error) {
	models, err := h.svc.ListCalibrationModels(ctx)
	if err != nil {
		return nil, toGRPC(err)
	}

	protoModels := make([]*v1.CalibrationModel, 0, len(models))
	for _, m := range models {
		protoModels = append(protoModels, &v1.CalibrationModel{
			ModelId:            m.ID,
			CartridgeCategory:  m.CartridgeCategory,
			CartridgeTypeIndex: m.CartridgeTypeIndex,
			Name:               m.Name,
			Version:            m.Version,
			DefaultAlpha:       m.DefaultAlpha,
			ValidityDays:       m.ValidityDays,
			Description:        m.Description,
			CreatedAt:          timestamppb.New(m.CreatedAt),
		})
	}

	return &v1.ListCalibrationModelsResponse{Models: protoModels}, nil
}

// ============================================================================
// 헬퍼 함수
// ============================================================================

// calibrationRecordToProto는 서비스 CalibrationRecord를 proto CalibrationRecord로 변환합니다.
func calibrationRecordToProto(r *service.CalibrationRecord) *v1.CalibrationRecord {
	return &v1.CalibrationRecord{
		CalibrationId:       r.ID,
		DeviceId:            r.DeviceID,
		CartridgeCategory:   r.CartridgeCategory,
		CartridgeTypeIndex:  r.CartridgeTypeIndex,
		CalibrationType:     serviceTypeToProto(r.CalibrationType),
		Alpha:               r.Alpha,
		ChannelOffsets:      r.ChannelOffsets,
		ChannelGains:        r.ChannelGains,
		TempCoefficient:     r.TempCoefficient,
		HumidityCoefficient: r.HumidityCoefficient,
		AccuracyScore:       r.AccuracyScore,
		ReferenceStandard:   r.ReferenceStandard,
		CalibratedBy:        r.CalibratedBy,
		CalibratedAt:        timestamppb.New(r.CalibratedAt),
		ExpiresAt:           timestamppb.New(r.ExpiresAt),
		Status:              serviceStatusToProto(r.Status),
	}
}

// serviceTypeToProto는 서비스 CalibrationType을 proto로 변환합니다.
func serviceTypeToProto(t service.CalibrationType) v1.CalibrationType {
	switch t {
	case service.CalibrationTypeFactory:
		return v1.CalibrationType_CALIBRATION_TYPE_FACTORY
	case service.CalibrationTypeField:
		return v1.CalibrationType_CALIBRATION_TYPE_FIELD
	case service.CalibrationTypeAuto:
		return v1.CalibrationType_CALIBRATION_TYPE_AUTO
	default:
		return v1.CalibrationType_CALIBRATION_TYPE_UNKNOWN
	}
}

// serviceStatusToProto는 서비스 CalibrationStatus를 proto로 변환합니다.
func serviceStatusToProto(s service.CalibrationStatus) v1.CalibrationStatus {
	switch s {
	case service.CalibrationStatusValid:
		return v1.CalibrationStatus_CALIBRATION_STATUS_VALID
	case service.CalibrationStatusExpiring:
		return v1.CalibrationStatus_CALIBRATION_STATUS_EXPIRING
	case service.CalibrationStatusExpired:
		return v1.CalibrationStatus_CALIBRATION_STATUS_EXPIRED
	case service.CalibrationStatusNeeded:
		return v1.CalibrationStatus_CALIBRATION_STATUS_NEEDED
	default:
		return v1.CalibrationStatus_CALIBRATION_STATUS_UNKNOWN
	}
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
