// Package service는 iot-gateway-service의 비즈니스 로직을 구현합니다.
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ============================================================================
// 도메인 상수
// ============================================================================

// 지원 프로토콜 목록
var validProtocols = map[string]bool{
	"mqtt":      true,
	"ble":       true,
	"coap":      true,
	"http":      true,
	"websocket": true,
}

// 디바이스 상태 상수
const (
	DeviceStatusRegistered    = "registered"
	DeviceStatusOnline        = "online"
	DeviceStatusOffline       = "offline"
	DeviceStatusError         = "error"
	DeviceStatusDecommissioned = "decommissioned"
)

// 명령 상태 상수
const (
	CommandStatusPending      = "pending"
	CommandStatusSent         = "sent"
	CommandStatusAcknowledged = "acknowledged"
	CommandStatusCompleted    = "completed"
	CommandStatusFailed       = "failed"
	CommandStatusExpired      = "expired"
)

// 허용되는 디바이스 상태 전이
var validDeviceTransitions = map[string]map[string]bool{
	DeviceStatusRegistered:    {DeviceStatusOnline: true, DeviceStatusDecommissioned: true},
	DeviceStatusOnline:        {DeviceStatusOffline: true, DeviceStatusError: true, DeviceStatusDecommissioned: true},
	DeviceStatusOffline:       {DeviceStatusOnline: true, DeviceStatusDecommissioned: true},
	DeviceStatusError:         {DeviceStatusOnline: true, DeviceStatusOffline: true, DeviceStatusDecommissioned: true},
}

// 허용되는 명령 상태 전이
var validCommandTransitions = map[string]map[string]bool{
	CommandStatusPending:      {CommandStatusSent: true, CommandStatusFailed: true, CommandStatusExpired: true},
	CommandStatusSent:         {CommandStatusAcknowledged: true, CommandStatusFailed: true, CommandStatusExpired: true},
	CommandStatusAcknowledged: {CommandStatusCompleted: true, CommandStatusFailed: true},
}

// 지원 데이터 타입 목록
var validDataTypes = map[string]bool{
	"temperature":   true,
	"heart_rate":    true,
	"spo2":          true,
	"blood_glucose": true,
	"blood_pressure": true,
	"weight":        true,
	"humidity":      true,
	"battery":       true,
	"signal":        true,
}

// 에러 정의
var (
	ErrInvalidProtocol   = errors.New("지원하지 않는 프로토콜입니다")
	ErrInvalidTransition = errors.New("허용되지 않는 상태 전이입니다")
	ErrInvalidDataType   = errors.New("지원하지 않는 데이터 타입입니다")
	ErrDeviceNotFound    = errors.New("디바이스를 찾을 수 없습니다")
	ErrCommandNotFound   = errors.New("명령을 찾을 수 없습니다")
	ErrDeviceNotOnline   = errors.New("디바이스가 온라인 상태가 아닙니다")
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
	RegisterDevice(ctx context.Context, device *IoTDevice) error
	GetDevice(ctx context.Context, deviceID string) (*IoTDevice, error)
	UpdateDevice(ctx context.Context, device *IoTDevice) error
	SendCommand(ctx context.Context, cmd *IoTCommand) error
	GetCommand(ctx context.Context, commandID string) (*IoTCommand, error)
	UpdateCommand(ctx context.Context, cmd *IoTCommand) error
	ListCommandsByDevice(ctx context.Context, deviceID string, limit int) ([]*IoTCommand, error)
	ReceiveData(ctx context.Context, data *IoTData) error
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
	if !validProtocols[protocol] {
		return nil, fmt.Errorf("%w: %s", ErrInvalidProtocol, protocol)
	}

	// 중복 등록 방지
	existing, _ := s.repo.GetDevice(ctx, deviceID)
	if existing != nil {
		return nil, fmt.Errorf("이미 등록된 디바이스입니다: %s", deviceID)
	}

	now := time.Now().UTC()
	device := &IoTDevice{
		ID:         uuid.New().String(),
		DeviceID:   deviceID,
		Protocol:   protocol,
		LastPingAt: now,
		Status:     DeviceStatusRegistered,
		Metadata:   metadata,
	}

	if err := s.repo.RegisterDevice(ctx, device); err != nil {
		return nil, errors.New("디바이스 등록에 실패했습니다: " + err.Error())
	}

	return device, nil
}

// UpdateDeviceStatus는 디바이스 상태를 변경합니다.
func (s *IoTGatewayService) UpdateDeviceStatus(ctx context.Context, deviceID, newStatus string) (*IoTDevice, error) {
	if deviceID == "" {
		return nil, errors.New("device_id는 필수입니다")
	}

	device, err := s.repo.GetDevice(ctx, deviceID)
	if err != nil {
		return nil, err
	}
	if device == nil {
		return nil, ErrDeviceNotFound
	}

	// 상태 전이 검증
	allowed, ok := validDeviceTransitions[device.Status]
	if !ok || !allowed[newStatus] {
		return nil, fmt.Errorf("%w: %s → %s", ErrInvalidTransition, device.Status, newStatus)
	}

	device.Status = newStatus
	if newStatus == DeviceStatusOnline {
		device.LastPingAt = time.Now().UTC()
	}

	if err := s.repo.UpdateDevice(ctx, device); err != nil {
		return nil, err
	}

	return device, nil
}

// Heartbeat는 디바이스 핑을 갱신하고 상태를 online으로 전환합니다.
func (s *IoTGatewayService) Heartbeat(ctx context.Context, deviceID string) (*IoTDevice, error) {
	if deviceID == "" {
		return nil, errors.New("device_id는 필수입니다")
	}

	device, err := s.repo.GetDevice(ctx, deviceID)
	if err != nil {
		return nil, err
	}
	if device == nil {
		return nil, ErrDeviceNotFound
	}

	device.LastPingAt = time.Now().UTC()
	if device.Status == DeviceStatusRegistered || device.Status == DeviceStatusOffline || device.Status == DeviceStatusError {
		device.Status = DeviceStatusOnline
	}

	if err := s.repo.UpdateDevice(ctx, device); err != nil {
		return nil, err
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

	device, err := s.repo.GetDevice(ctx, deviceID)
	if err != nil {
		return nil, errors.New("디바이스 조회에 실패했습니다: " + err.Error())
	}
	if device == nil {
		return nil, fmt.Errorf("%w: %s", ErrDeviceNotFound, deviceID)
	}

	// decommissioned 디바이스에는 명령 불가
	if device.Status == DeviceStatusDecommissioned {
		return nil, errors.New("폐기된 디바이스에는 명령을 전송할 수 없습니다")
	}

	now := time.Now().UTC()
	cmd := &IoTCommand{
		ID:          uuid.New().String(),
		DeviceID:    deviceID,
		CommandType: commandType,
		Payload:     payload,
		Status:      CommandStatusPending,
		CreatedAt:   now,
	}

	if err := s.repo.SendCommand(ctx, cmd); err != nil {
		return nil, errors.New("명령 전송에 실패했습니다: " + err.Error())
	}

	return cmd, nil
}

// UpdateCommandStatus는 명령 상태를 변경합니다.
func (s *IoTGatewayService) UpdateCommandStatus(ctx context.Context, commandID, newStatus string) (*IoTCommand, error) {
	if commandID == "" {
		return nil, errors.New("command_id는 필수입니다")
	}

	cmd, err := s.repo.GetCommand(ctx, commandID)
	if err != nil {
		return nil, err
	}
	if cmd == nil {
		return nil, ErrCommandNotFound
	}

	allowed, ok := validCommandTransitions[cmd.Status]
	if !ok || !allowed[newStatus] {
		return nil, fmt.Errorf("%w: %s → %s", ErrInvalidTransition, cmd.Status, newStatus)
	}

	cmd.Status = newStatus
	if err := s.repo.UpdateCommand(ctx, cmd); err != nil {
		return nil, err
	}

	return cmd, nil
}

// ListCommandsByDevice는 디바이스의 명령 목록을 조회합니다.
func (s *IoTGatewayService) ListCommandsByDevice(ctx context.Context, deviceID string, limit int) ([]*IoTCommand, error) {
	if deviceID == "" {
		return nil, errors.New("device_id는 필수입니다")
	}
	if limit <= 0 {
		limit = 20
	}
	return s.repo.ListCommandsByDevice(ctx, deviceID, limit)
}

// ReceiveData는 디바이스로부터 데이터를 수신하여 저장합니다.
func (s *IoTGatewayService) ReceiveData(ctx context.Context, deviceID, dataType string, value float64, unit string) (*IoTData, error) {
	if deviceID == "" {
		return nil, errors.New("device_id는 필수입니다")
	}
	if dataType == "" {
		return nil, errors.New("data_type은 필수입니다")
	}
	if !validDataTypes[dataType] {
		return nil, fmt.Errorf("%w: %s", ErrInvalidDataType, dataType)
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
		return nil, fmt.Errorf("%w: %s", ErrDeviceNotFound, deviceID)
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
		return nil, fmt.Errorf("%w: %s", ErrCommandNotFound, commandID)
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
