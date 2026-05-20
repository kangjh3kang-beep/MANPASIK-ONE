package handler

import (
	"context"

	"github.com/manpasik/backend/services/digital-twin-service/internal/service"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TwinHandler는 디지털 트윈 gRPC 핸들러입니다.
// MeasurementService의 SyncDigitalTwin RPC를 위임받아 처리합니다.
type TwinHandler struct {
	v1.UnimplementedMeasurementServiceServer
	svc *service.TwinService
	log *zap.Logger
}

// NewTwinHandler는 새 TwinHandler를 생성합니다.
func NewTwinHandler(svc *service.TwinService, log *zap.Logger) *TwinHandler {
	return &TwinHandler{svc: svc, log: log}
}

// SyncDigitalTwin은 디지털 트윈 동기화 RPC입니다.
func (h *TwinHandler) SyncDigitalTwin(ctx context.Context, req *v1.SyncDigitalTwinRequest) (*v1.SyncDigitalTwinResponse, error) {
	if req == nil || req.SessionId == "" || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id와 user_id는 필수입니다")
	}

	state, action, err := h.svc.SyncTwin(ctx,
		req.SessionId,
		req.UserId,
		req.DeviceId,
		req.Residuals,
		req.EwmaValue,
		req.CusumPos,
		req.CusumNeg,
		req.HealthState,
		req.DriftScore,
		int(req.RemainingMeasurements),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &v1.SyncDigitalTwinResponse{
		Accepted:             true,
		TwinId:               state.ID,
		RecommendedAction:    action,
		NextCalibrationDrift: 0.5,
		SyncedAt:             timestamppb.New(state.LastSyncedAt),
	}, nil
}

// GetCalibrationStatus는 보정 상태 조회 RPC입니다.
func (h *TwinHandler) GetCalibrationStatus(ctx context.Context, req *v1.GetCalibrationStatusRequest) (*v1.GetCalibrationStatusResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "요청이 비어있습니다")
	}

	state, _ := h.svc.GetTwinState(ctx, req.SessionId)

	calState := "valid"
	driftScore := 0.0
	measurementsSinceCal := int32(0)

	if state != nil {
		driftScore = state.DriftScore
		measurementsSinceCal = int32(state.MeasurementCount)
		switch state.HealthState {
		case "recalibrate":
			calState = "needs_recalibration"
		case "drift":
			calState = "expiring"
		case "warning":
			calState = "valid"
		}
	}

	return &v1.GetCalibrationStatusResponse{
		CalibrationState:         calState,
		DriftScore:               driftScore,
		MeasurementsSinceCal:     measurementsSinceCal,
		MaxMeasurementsBeforeCal: 500,
	}, nil
}
