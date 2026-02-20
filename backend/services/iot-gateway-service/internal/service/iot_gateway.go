// Package service는 iot-gateway-service의 비즈니스 로직을 구현합니다.
package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ============================================================================
// 엔티티 (Entities)
// ============================================================================

// IoTDevice는 IoT 디바이스 엔티티입니다.
type IoTDevice struct {
	ID         string
	DeviceID   string
	Protocol   string
	LastPingAt time.Time
	Status     string
	Metadata   map[string]string
}

// IoTCommand는 디바이스에 전송하는 명령 엔티티입니다.
type IoTCommand struct {
	ID          string
	DeviceID    string
	CommandType string
	Payload     string
	Status      string
	CreatedAt   time.Time
}

// IoTData는 디바이스로부터 수신한 데이터 엔티티입니다.
type IoTData struct {
	ID         string
	DeviceID   string
	DataType   string
	Value      float64
	Unit       string
	ReceivedAt time.Time
}

// ============================================================================
// Repository 인터페이스
// ============================================================================

// IoTRepository는 IoT 데이터 저장소 인터페이스입니다.
type IoTRepository interface {
	// RegisterDevice는 디바이스를 등록합니다.
	RegisterDevice(ctx context.Context, device *IoTDevice) error
	// GetDevice는 디바이스 ID로 조회합니다.
	GetDevice(ctx context.Context, deviceID string) (*IoTDevice, error)
	// SendCommand는 명령을 저장합니다.
	SendCommand(ctx context.Context, cmd *IoTCommand) error
	// GetCommand는 명령 ID로 조회합니다.
	GetCommand(ctx context.Context, commandID string) (*IoTCommand, error)
	// ReceiveData는 수신 데이터를 저장합니다.
	ReceiveData(ctx context.Context, data *IoTData) error
	// ListData는 디바이스의 수신 데이터 목록을 조회합니다.
	ListData(ctx context.Context, deviceID string, limit, offset int) ([]*IoTData, error)
}

// ============================================================================
// IoTGatewayService
// ============================================================================

// IoTGatewayService는 IoT 게이트웨이 비즈니스 로직입니다.
type IoTGatewayService struct {
	repo IoTRepository
}

// NewIoTGatewayService는 새 IoTGatewayService를 생성합니다.
func NewIoTGatewayService(repo IoTRepository) *IoTGatewayService {
	return &IoTGatewayService{repo: repo}
}

// RegisterDevice는 새 IoT 디바이스를 등록합니다.
func (s *IoTGatewayService) RegisterDevice(ctx context.Context, deviceID, protocol string, metadata map[string]string) (*IoTDevice, error) {
	if deviceID == "" {
		return nil, errors.New("device_id는 필수입니다")
	}
	if protocol == "" {
		return nil, errors.New("protocol은 필수입니다")
	}

	now := time.Now().UTC()
	device := &IoTDevice{
		ID:         uuid.New().String(),
		DeviceID:   deviceID,
		Protocol:   protocol,
		LastPingAt: now,
		Status:     "registered",
		Metadata:   metadata,
	}

	if err := s.repo.RegisterDevice(ctx, device); err != nil {
		return nil, errors.New("디바이스 등록에 실패했습니다: " + err.Error())
	}

	return device, nil
}

// SendCommand는 디바이스에 명령을 전송합니다.
func (s *IoTGatewayService) SendCommand(ctx context.Context, deviceID, commandType, payload string) (*IoTCommand, error) {
	if deviceID == "" {
		return nil, errors.New("device_id는 필수입니다")
	}
	if commandType == "" {
		return nil, errors.New("command_type은 필수입니다")
	}

	// 디바이스 존재 여부 확인
	device, err := s.repo.GetDevice(ctx, deviceID)
	if err != nil {
		return nil, errors.New("디바이스 조회에 실패했습니다: " + err.Error())
	}
	if device == nil {
		return nil, errors.New("디바이스를 찾을 수 없습니다: " + deviceID)
	}

	now := time.Now().UTC()
	cmd := &IoTCommand{
		ID:          uuid.New().String(),
		DeviceID:    deviceID,
		CommandType: commandType,
		Payload:     payload,
		Status:      "pending",
		CreatedAt:   now,
	}

	if err := s.repo.SendCommand(ctx, cmd); err != nil {
		return nil, errors.New("명령 전송에 실패했습니다: " + err.Error())
	}

	return cmd, nil
}

// ReceiveData는 디바이스로부터 데이터를 수신하여 저장합니다.
func (s *IoTGatewayService) ReceiveData(ctx context.Context, deviceID, dataType string, value float64, unit string) (*IoTData, error) {
	if deviceID == "" {
		return nil, errors.New("device_id는 필수입니다")
	}
	if dataType == "" {
		return nil, errors.New("data_type은 필수입니다")
	}

	now := time.Now().UTC()
	data := &IoTData{
		ID:         uuid.New().String(),
		DeviceID:   deviceID,
		DataType:   dataType,
		Value:      value,
		Unit:       unit,
		ReceivedAt: now,
	}

	if err := s.repo.ReceiveData(ctx, data); err != nil {
		return nil, errors.New("데이터 수신에 실패했습니다: " + err.Error())
	}

	return data, nil
}

// GetDevice는 디바이스를 조회합니다.
func (s *IoTGatewayService) GetDevice(ctx context.Context, deviceID string) (*IoTDevice, error) {
	if deviceID == "" {
		return nil, errors.New("device_id는 필수입니다")
	}

	device, err := s.repo.GetDevice(ctx, deviceID)
	if err != nil {
		return nil, errors.New("디바이스 조회에 실패했습니다: " + err.Error())
	}
	if device == nil {
		return nil, errors.New("디바이스를 찾을 수 없습니다: " + deviceID)
	}

	return device, nil
}

// GetCommand는 명령을 조회합니다.
func (s *IoTGatewayService) GetCommand(ctx context.Context, commandID string) (*IoTCommand, error) {
	if commandID == "" {
		return nil, errors.New("command_id는 필수입니다")
	}

	cmd, err := s.repo.GetCommand(ctx, commandID)
	if err != nil {
		return nil, errors.New("명령 조회에 실패했습니다: " + err.Error())
	}
	if cmd == nil {
		return nil, errors.New("명령을 찾을 수 없습니다: " + commandID)
	}

	return cmd, nil
}

// ListData는 디바이스의 수신 데이터 목록을 조회합니다.
func (s *IoTGatewayService) ListData(ctx context.Context, deviceID string, limit, offset int) ([]*IoTData, error) {
	if deviceID == "" {
		return nil, errors.New("device_id는 필수입니다")
	}
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.ListData(ctx, deviceID, limit, offset)
}
