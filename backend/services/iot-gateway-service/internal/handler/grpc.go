// Package handler는 iot-gateway-service의 gRPC 핸들러입니다.
package handler

import (
	"context"

	"github.com/manpasik/backend/services/iot-gateway-service/internal/service"
)

// IoTGatewayHandler는 IoTGatewayService를 래핑하는 핸들러입니다.
type IoTGatewayHandler struct {
	svc *service.IoTGatewayService
}

// NewIoTGatewayHandler는 IoTGatewayHandler를 생성합니다.
func NewIoTGatewayHandler(svc *service.IoTGatewayService) *IoTGatewayHandler {
	return &IoTGatewayHandler{svc: svc}
}

// RegisterDevice는 디바이스 등록을 처리합니다.
func (h *IoTGatewayHandler) RegisterDevice(ctx context.Context, deviceID, protocol string, meta map[string]string) (*service.IoTDevice, error) {
	return h.svc.RegisterDevice(ctx, deviceID, protocol, meta)
}

// SendCommand는 디바이스 명령 전송을 처리합니다.
func (h *IoTGatewayHandler) SendCommand(ctx context.Context, deviceID, commandType, payload string) (*service.IoTCommand, error) {
	return h.svc.SendCommand(ctx, deviceID, commandType, payload)
}

// ReceiveData는 디바이스 데이터 수신을 처리합니다.
func (h *IoTGatewayHandler) ReceiveData(ctx context.Context, deviceID, dataType string, value float64, unit string) (*service.IoTData, error) {
	return h.svc.ReceiveData(ctx, deviceID, dataType, value, unit)
}

// GetDevice는 디바이스 조회를 처리합니다.
func (h *IoTGatewayHandler) GetDevice(ctx context.Context, deviceID string) (*service.IoTDevice, error) {
	return h.svc.GetDevice(ctx, deviceID)
}

// GetCommand는 명령 조회를 처리합니다.
func (h *IoTGatewayHandler) GetCommand(ctx context.Context, commandID string) (*service.IoTCommand, error) {
	return h.svc.GetCommand(ctx, commandID)
}

// ListData는 디바이스 데이터 목록 조회를 처리합니다.
func (h *IoTGatewayHandler) ListData(ctx context.Context, deviceID string, limit, offset int) ([]*service.IoTData, error) {
	return h.svc.ListData(ctx, deviceID, limit, offset)
}
