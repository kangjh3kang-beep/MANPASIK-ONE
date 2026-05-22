// Package handler는 iot-gateway-service의 gRPC 핸들러를 구현합니다.
package handler

import (
	"context"

	"github.com/manpasik/backend/services/iot-gateway-service/internal/service"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// IoTGatewayHandler는 gRPC IoTGatewayService 핸들러입니다.
type IoTGatewayHandler struct {
	v1.UnimplementedIoTGatewayServiceServer
	svc *service.IoTGatewayService
	log *zap.Logger
}

// NewIoTGatewayHandler는 IoTGatewayHandler를 생성합니다.
func NewIoTGatewayHandler(svc *service.IoTGatewayService, log *zap.Logger) *IoTGatewayHandler {
	return &IoTGatewayHandler{svc: svc, log: log}
}

// RegisterDevice는 새 IoT 디바이스를 등록하는 RPC입니다.
func (h *IoTGatewayHandler) RegisterDevice(ctx context.Context, req *v1.RegisterIoTDeviceRequest) (*v1.IoTDeviceResponse, error) {
	if req == nil || req.DeviceId == "" {
		return nil, status.Error(codes.InvalidArgument, "device_id는 필수입니다")
	}
	if req.Protocol == "" {
		return nil, status.Error(codes.InvalidArgument, "protocol은 필수입니다")
	}

	device, err := h.svc.RegisterDevice(ctx, req.DeviceId, req.Protocol, req.Metadata)
	if err != nil {
		return nil, toGRPC(err)
	}

	return deviceToProto(device), nil
}

// GetDevice는 디바이스를 조회하는 RPC입니다.
func (h *IoTGatewayHandler) GetDevice(ctx context.Context, req *v1.GetIoTDeviceRequest) (*v1.IoTDeviceResponse, error) {
	if req == nil || req.DeviceId == "" {
		return nil, status.Error(codes.InvalidArgument, "device_id는 필수입니다")
	}

	device, err := h.svc.GetDevice(ctx, req.DeviceId)
	if err != nil {
		return nil, toGRPC(err)
	}

	return deviceToProto(device), nil
}

// SendDeviceCommand는 디바이스에 명령을 전송하는 RPC입니다.
func (h *IoTGatewayHandler) SendDeviceCommand(ctx context.Context, req *v1.SendDeviceCommandRequest) (*v1.DeviceCommandResponse, error) {
	if req == nil || req.DeviceId == "" {
		return nil, status.Error(codes.InvalidArgument, "device_id는 필수입니다")
	}
	if req.CommandType == "" {
		return nil, status.Error(codes.InvalidArgument, "command_type은 필수입니다")
	}

	cmd, err := h.svc.SendCommand(ctx, req.DeviceId, req.CommandType, req.Payload)
	if err != nil {
		return nil, toGRPC(err)
	}

	return commandToProto(cmd), nil
}

// GetDeviceTelemetry는 디바이스 텔레메트리 데이터를 조회하는 RPC입니다.
func (h *IoTGatewayHandler) GetDeviceTelemetry(ctx context.Context, req *v1.GetDeviceTelemetryRequest) (*v1.DeviceTelemetryResponse, error) {
	if req == nil || req.DeviceId == "" {
		return nil, status.Error(codes.InvalidArgument, "device_id는 필수입니다")
	}

	limit := int(req.Limit)
	if limit <= 0 {
		limit = 20
	}
	offset := int(req.Offset)
	if offset < 0 {
		offset = 0
	}

	data, err := h.svc.ListData(ctx, req.DeviceId, limit, offset)
	if err != nil {
		return nil, toGRPC(err)
	}

	pbData := make([]*v1.IoTDataPoint, 0, len(data))
	for _, d := range data {
		pbData = append(pbData, dataToProto(d))
	}

	return &v1.DeviceTelemetryResponse{Data: pbData}, nil
}

// Heartbeat는 디바이스 하트비트를 처리하는 RPC입니다.
func (h *IoTGatewayHandler) Heartbeat(ctx context.Context, req *v1.HeartbeatRequest) (*v1.IoTDeviceResponse, error) {
	if req == nil || req.DeviceId == "" {
		return nil, status.Error(codes.InvalidArgument, "device_id는 필수입니다")
	}

	device, err := h.svc.Heartbeat(ctx, req.DeviceId)
	if err != nil {
		return nil, toGRPC(err)
	}

	return deviceToProto(device), nil
}

// ListDeviceData는 디바이스 데이터 목록을 조회하는 RPC입니다.
func (h *IoTGatewayHandler) ListDeviceData(ctx context.Context, req *v1.ListDeviceDataRequest) (*v1.ListDeviceDataResponse, error) {
	if req == nil || req.DeviceId == "" {
		return nil, status.Error(codes.InvalidArgument, "device_id는 필수입니다")
	}

	limit := int(req.Limit)
	if limit <= 0 {
		limit = 20
	}
	offset := int(req.Offset)
	if offset < 0 {
		offset = 0
	}

	data, err := h.svc.ListData(ctx, req.DeviceId, limit, offset)
	if err != nil {
		return nil, toGRPC(err)
	}

	pbData := make([]*v1.IoTDataPoint, 0, len(data))
	for _, d := range data {
		pbData = append(pbData, dataToProto(d))
	}

	return &v1.ListDeviceDataResponse{Data: pbData}, nil
}

// ============================================================================
// Helpers
// ============================================================================

func deviceToProto(d *service.IoTDevice) *v1.IoTDeviceResponse {
	return &v1.IoTDeviceResponse{
		Id:         d.ID,
		DeviceId:   d.DeviceID,
		Protocol:   d.Protocol,
		Status:     d.Status,
		Metadata:   d.Metadata,
		LastPingAt: timestamppb.New(d.LastPingAt),
	}
}

func commandToProto(c *service.IoTCommand) *v1.DeviceCommandResponse {
	return &v1.DeviceCommandResponse{
		Id:          c.ID,
		DeviceId:    c.DeviceID,
		CommandType: c.CommandType,
		Payload:     c.Payload,
		Status:      c.Status,
		CreatedAt:   timestamppb.New(c.CreatedAt),
	}
}

func dataToProto(d *service.IoTData) *v1.IoTDataPoint {
	return &v1.IoTDataPoint{
		Id:         d.ID,
		DeviceId:   d.DeviceID,
		DataType:   d.DataType,
		Value:      d.Value,
		Unit:       d.Unit,
		ReceivedAt: timestamppb.New(d.ReceivedAt),
	}
}

func toGRPC(err error) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	switch {
	case err == service.ErrDeviceNotFound || err == service.ErrCommandNotFound:
		return status.Error(codes.NotFound, msg)
	case err == service.ErrInvalidProtocol || err == service.ErrInvalidDataType:
		return status.Error(codes.InvalidArgument, msg)
	case err == service.ErrInvalidTransition || err == service.ErrDeviceNotOnline:
		return status.Error(codes.FailedPrecondition, msg)
	default:
		return status.Error(codes.Internal, msg)
	}
}
