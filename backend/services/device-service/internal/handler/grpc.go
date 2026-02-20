// Package handler는 device-service의 gRPC 핸들러입니다.
package handler

import (
	"context"

	"github.com/manpasik/backend/services/device-service/internal/service"
	apperrors "github.com/manpasik/backend/shared/errors"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// DeviceHandler는 DeviceService gRPC 서버를 구현합니다.
type DeviceHandler struct {
	v1.UnimplementedDeviceServiceServer
	svc *service.DeviceService
	log *zap.Logger
}

// NewDeviceHandler는 DeviceHandler를 생성합니다.
func NewDeviceHandler(svc *service.DeviceService, log *zap.Logger) *DeviceHandler {
	return &DeviceHandler{svc: svc, log: log}
}

// RegisterDevice는 디바이스 등록 RPC입니다.
func (h *DeviceHandler) RegisterDevice(ctx context.Context, req *v1.RegisterDeviceRequest) (*v1.RegisterDeviceResponse, error) {
	if req == nil || req.DeviceId == "" || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "device_id와 user_id는 필수입니다")
	}

	device, regToken, err := h.svc.RegisterDevice(ctx, req.DeviceId, req.SerialNumber, req.FirmwareVersion, req.UserId)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.RegisterDeviceResponse{
		DeviceId:          device.ID,
		RegistrationToken: regToken,
		RegisteredAt:      timestamppb.New(device.RegisteredAt),
	}, nil
}

// ListDevices는 디바이스 목록 조회 RPC입니다.
func (h *DeviceHandler) ListDevices(ctx context.Context, req *v1.ListDevicesRequest) (*v1.ListDevicesResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	devices, err := h.svc.ListDevices(ctx, req.UserId)
	if err != nil {
		return nil, toGRPC(err)
	}

	var infos []*v1.DeviceInfo
	for _, d := range devices {
		infos = append(infos, &v1.DeviceInfo{
			DeviceId:        d.ID,
			Name:            d.Name,
			FirmwareVersion: d.FirmwareVersion,
			LastSeen:        timestamppb.New(d.LastSeen),
		})
	}

	return &v1.ListDevicesResponse{Devices: infos}, nil
}

// UpdateDeviceStatus는 디바이스 상태 업데이트 RPC입니다.
func (h *DeviceHandler) UpdateDeviceStatus(ctx context.Context, req *v1.UpdateDeviceStatusRequest) (*v1.UpdateDeviceStatusResponse, error) {
	if req == nil || req.DeviceId == "" {
		return nil, status.Error(codes.InvalidArgument, "device_id는 필수입니다")
	}

	err := h.svc.UpdateDeviceStatus(ctx, req.DeviceId, service.DeviceStatus(req.Status), 0)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.UpdateDeviceStatusResponse{
		Success: true,
		Message: "디바이스 상태가 업데이트되었습니다",
	}, nil
}

// RequestOtaUpdate는 OTA 펌웨어 업데이트 요청 RPC입니다.
func (h *DeviceHandler) RequestOtaUpdate(ctx context.Context, req *v1.OtaRequest) (*v1.OtaResponse, error) {
	if req == nil || req.DeviceId == "" || req.TargetVersion == "" {
		return nil, status.Error(codes.InvalidArgument, "device_id와 target_version은 필수입니다")
	}

	updateID, downloadURL, checksum, err := h.svc.RequestOtaUpdate(ctx, req.DeviceId, req.TargetVersion)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.OtaResponse{
		UpdateId:    updateID,
		DownloadUrl: downloadURL,
		Checksum:    checksum,
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
