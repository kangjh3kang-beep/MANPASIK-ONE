// Package handler??measurement-service??gRPC ?몃뱾?ъ엯?덈떎.
package handler

import (
	"context"
	"io"
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

// MeasurementHandler??MeasurementService gRPC ?쒕쾭瑜?援ы쁽?⑸땲??
type MeasurementHandler struct {
	v1.UnimplementedMeasurementServiceServer
	svc *service.MeasurementService
	log *zap.Logger
}

// NewMeasurementHandler??MeasurementHandler瑜??앹꽦?⑸땲??
func NewMeasurementHandler(svc *service.MeasurementService, log *zap.Logger) *MeasurementHandler {
	return &MeasurementHandler{svc: svc, log: log}
}

// StartSession? 痢≪젙 ?몄뀡 ?쒖옉 RPC?낅땲??
func (h *MeasurementHandler) StartSession(ctx context.Context, req *v1.StartSessionRequest) (*v1.StartSessionResponse, error) {
	if req == nil || req.DeviceId == "" || req.CartridgeId == "" || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "device_id, cartridge_id, user_id are required")
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

// EndSession? 痢≪젙 ?몄뀡 醫낅즺 RPC?낅땲??
// StreamMeasurement processes measurement frames and returns processed results.
func (h *MeasurementHandler) StreamMeasurement(stream v1.MeasurementService_StreamMeasurementServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return toGRPC(err)
		}
		if req == nil || req.SessionId == "" {
			return status.Error(codes.InvalidArgument, "session_id is required")
		}

		fingerprint := make([]float32, 0, len(req.RawChannels))
		for _, channel := range req.RawChannels {
			fingerprint = append(fingerprint, float32(channel))
		}

		data := &service.MeasurementData{
			Time:              time.Now().UTC(),
			SessionID:         req.SessionId,
			RawChannels:       req.RawChannels,
			SDet:              req.GetDifferential().GetSDet(),
			SRef:              req.GetDifferential().GetSRef(),
			Alpha:             req.GetDifferential().GetAlpha(),
			SCorrected:        req.GetDifferential().GetSCorrected(),
			FingerprintVector: fingerprint,
			TempC:             req.GetEnvMeta().GetTempC(),
			HumidityPct:       req.GetEnvMeta().GetHumidityPct(),
		}

		result, err := h.svc.ProcessMeasurement(stream.Context(), data)
		if err != nil {
			return toGRPC(err)
		}

		if err := stream.Send(&v1.MeasurementResult{
			SessionId:         result.SessionID,
			PrimaryValue:      result.PrimaryValue,
			Unit:              result.Unit,
			Confidence:        result.Confidence,
			FingerprintVector: fingerprint,
			ProcessedAt:       timestamppb.New(result.ProcessedAt),
			EvidenceStatus:    string(result.EvidenceStatus),
			DiagnosticReady:   result.DiagnosticReady,
			EvidenceGaps:      append([]string(nil), result.EvidenceGaps...),
		}); err != nil {
			return toGRPC(err)
		}
	}
}
func (h *MeasurementHandler) EndSession(ctx context.Context, req *v1.EndSessionRequest) (*v1.EndSessionResponse, error) {
	if req == nil || req.SessionId == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id is required")
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

// GetMeasurementHistory??痢≪젙 湲곕줉 議고쉶 RPC?낅땲??
func (h *MeasurementHandler) GetMeasurementHistory(ctx context.Context, req *v1.GetHistoryRequest) (*v1.GetHistoryResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
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
			SessionId:       s.SessionID,
			CartridgeType:   s.CartridgeType,
			PrimaryValue:    s.PrimaryValue,
			Unit:            s.Unit,
			MeasuredAt:      timestamppb.New(s.MeasuredAt),
			EvidenceStatus:  string(s.EvidenceStatus),
			DiagnosticReady: s.DiagnosticReady,
			EvidenceGaps:    append([]string(nil), s.EvidenceGaps...),
		})
	}

	return &v1.GetHistoryResponse{
		Measurements: pbSummaries,
		TotalCount:   int32(total),
	}, nil
}

// ExportSingleMeasurement???⑥씪 痢≪젙 ?몄뀡??FHIR ?대낫?닿린 RPC?낅땲??
func (h *MeasurementHandler) ExportSingleMeasurement(ctx context.Context, req *v1.ExportSingleMeasurementRequest) (*v1.ExportFHIRResponse, error) {
	if req == nil || req.SessionId == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id is required")
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

// ExportToFHIRObservations???ъ슜???꾩껜 痢≪젙 寃곌낵??FHIR Observation ?대낫?닿린 RPC?낅땲??
func (h *MeasurementHandler) ExportToFHIRObservations(ctx context.Context, req *v1.ExportToFHIRObservationsRequest) (*v1.ExportFHIRResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
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

// SyncDigitalTwin? ?붿????몄쐢 ?숆린??RPC?낅땲??
func (h *MeasurementHandler) SyncDigitalTwin(ctx context.Context, req *v1.SyncDigitalTwinRequest) (*v1.SyncDigitalTwinResponse, error) {
	if req == nil || req.SessionId == "" || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id and user_id are required")
	}

	state, err := h.svc.SyncDigitalTwin(
		ctx,
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
		req.FingerprintVector,
		int(req.FingerprintDim),
	)
	if err != nil {
		return nil, toGRPC(err)
	}

	recommendedAction := "continue"
	if state.DriftScore > 0.8 {
		recommendedAction = "replace_cartridge"
	} else if state.DriftScore > 0.5 {
		recommendedAction = "recalibrate"
	}

	return &v1.SyncDigitalTwinResponse{
		Accepted:             true,
		TwinId:               state.TwinID,
		RecommendedAction:    recommendedAction,
		NextCalibrationDrift: 0.5,
		SyncedAt:             timestamppb.New(state.SyncedAt),
	}, nil
}

// GetCalibrationStatus??蹂댁젙 ?곹깭 議고쉶 RPC?낅땲??
func (h *MeasurementHandler) GetCalibrationStatus(ctx context.Context, req *v1.GetCalibrationStatusRequest) (*v1.GetCalibrationStatusResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "?붿껌??鍮꾩뼱?덉뒿?덈떎")
	}

	calStatus, err := h.svc.GetCalibrationStatus(ctx, req.SessionId, req.DeviceId)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.GetCalibrationStatusResponse{
		CalibrationState:         calStatus.CalibrationState,
		DriftScore:               calStatus.DriftScore,
		MeasurementsSinceCal:     int32(calStatus.MeasurementsSinceCal),
		MaxMeasurementsBeforeCal: int32(calStatus.MaxMeasurementsBeforeCal),
		LastCalibratedAt:         timestamppb.New(calStatus.LastCalibratedAt),
		NextCalibrationAt:        timestamppb.New(calStatus.NextCalibrationAt),
	}, nil
}

// toGRPC??AppError瑜?gRPC status濡?蹂?섑빀?덈떎.
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
	return status.Error(codes.Internal, "?대? ?ㅻ쪟媛 諛쒖깮?덉뒿?덈떎")
}
