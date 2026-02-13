// Package service는 device-service의 비즈니스 로직을 구현합니다.
//
// 기능: 디바이스 등록, 상태 관리, OTA 업데이트, 구독 기반 디바이스 수 제한
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	apperrors "github.com/manpasik/backend/shared/errors"
	"go.uber.org/zap"
)

// ============================================================================
// 이벤트 타입 (Kafka 발행용)
// ============================================================================

// DeviceRegisteredEvent는 디바이스 등록 이벤트입니다.
type DeviceRegisteredEvent struct {
	DeviceID        string
	UserID          string
	SerialNumber    string
	FirmwareVersion string
	RegisteredAt    time.Time
}

// DeviceStatusChangedEvent는 디바이스 상태 변경 이벤트입니다.
type DeviceStatusChangedEvent struct {
	DeviceID        string
	UserID          string
	SerialNumber    string
	PreviousStatus  string
	NewStatus       string
	BatteryPercent  int
	FirmwareVersion string
	LastSeen        time.Time
}

// KafkaEventPublisher는 Kafka 이벤트 발행 인터페이스입니다.
type KafkaEventPublisher interface {
	PublishDeviceRegistered(ctx context.Context, event *DeviceRegisteredEvent) error
	PublishDeviceStatusChanged(ctx context.Context, event *DeviceStatusChangedEvent) error
}

// DeviceService는 디바이스 관리 서비스입니다.
type DeviceService struct {
	logger         *zap.Logger
	deviceRepo     DeviceRepository
	eventRepo      DeviceEventRepository
	subChecker     SubscriptionChecker
	eventPublisher KafkaEventPublisher
}

// DeviceRepository는 디바이스 데이터 저장소 인터페이스입니다.
type DeviceRepository interface {
	Create(ctx context.Context, device *Device) error
	GetByID(ctx context.Context, deviceID string) (*Device, error)
	ListByUser(ctx context.Context, userID string) ([]*Device, error)
	UpdateStatus(ctx context.Context, deviceID string, status DeviceStatus, battery int, lastSeen time.Time) error
	CountByUser(ctx context.Context, userID string) (int, error)
	Delete(ctx context.Context, deviceID string) error
}

// DeviceEventRepository는 디바이스 이벤트 저장소입니다.
type DeviceEventRepository interface {
	LogEvent(ctx context.Context, event *DeviceEvent) error
}

// SubscriptionChecker는 구독 정보 확인 인터페이스입니다.
type SubscriptionChecker interface {
	GetMaxDevices(ctx context.Context, userID string) (int, error)
}

// DeviceStatus는 디바이스 상태입니다.
type DeviceStatus string

const (
	StatusUnknown   DeviceStatus = "unknown"
	StatusOnline    DeviceStatus = "online"
	StatusOffline   DeviceStatus = "offline"
	StatusMeasuring DeviceStatus = "measuring"
	StatusUpdating  DeviceStatus = "updating"
	StatusError     DeviceStatus = "error"
)

// Device는 디바이스 엔티티입니다.
type Device struct {
	ID              string
	DeviceID        string // 하드웨어 고유 ID
	UserID          string
	Name            string
	SerialNumber    string
	FirmwareVersion string
	Status          DeviceStatus
	BatteryPercent  int
	LastSeen        time.Time
	RegisteredAt    time.Time
}

// DeviceEvent는 디바이스 이벤트 로그입니다.
type DeviceEvent struct {
	ID        string
	DeviceID  string
	EventType string // "registered", "status_changed", "ota_started", "ota_completed", "error"
	Payload   map[string]interface{}
	CreatedAt time.Time
}

// NewDeviceService는 새 DeviceService를 생성합니다.
func NewDeviceService(
	logger *zap.Logger,
	deviceRepo DeviceRepository,
	eventRepo DeviceEventRepository,
	subChecker SubscriptionChecker,
) *DeviceService {
	return &DeviceService{
		logger:     logger,
		deviceRepo: deviceRepo,
		eventRepo:  eventRepo,
		subChecker: subChecker,
	}
}

// SetEventPublisher는 Kafka 이벤트 발행기를 설정합니다 (optional).
func (s *DeviceService) SetEventPublisher(ep KafkaEventPublisher) {
	s.eventPublisher = ep
}

// RegisterDevice는 새 디바이스를 등록합니다.
func (s *DeviceService) RegisterDevice(
	ctx context.Context,
	deviceHwID, serialNumber, firmwareVersion, userID string,
) (*Device, string, error) {
	// 입력 검증
	if deviceHwID == "" || userID == "" {
		return nil, "", apperrors.New(apperrors.ErrInvalidInput, "device_id와 user_id는 필수입니다")
	}

	// 구독 기반 디바이스 수 제한 확인
	currentCount, err := s.deviceRepo.CountByUser(ctx, userID)
	if err != nil {
		return nil, "", apperrors.New(apperrors.ErrInternal, "디바이스 수 확인에 실패했습니다")
	}

	maxDevices, err := s.subChecker.GetMaxDevices(ctx, userID)
	if err != nil {
		// 기본값 사용 (Free 티어 = 무제한, 만파식 정책)
		maxDevices = 999
	}

	if currentCount >= maxDevices {
		return nil, "", apperrors.New(apperrors.ErrDeviceLimitExceeded,
			fmt.Sprintf("디바이스 등록 한도 초과 (현재: %d, 최대: %d)", currentCount, maxDevices))
	}

	// 등록 토큰 생성
	registrationToken := uuid.New().String()

	device := &Device{
		ID:              uuid.New().String(),
		DeviceID:        deviceHwID,
		UserID:          userID,
		Name:            fmt.Sprintf("ManPaSik Reader %s", serialNumber[:4]),
		SerialNumber:    serialNumber,
		FirmwareVersion: firmwareVersion,
		Status:          StatusOnline,
		BatteryPercent:  100,
		LastSeen:        time.Now().UTC(),
		RegisteredAt:    time.Now().UTC(),
	}

	if err := s.deviceRepo.Create(ctx, device); err != nil {
		return nil, "", apperrors.New(apperrors.ErrInternal, "디바이스 등록에 실패했습니다")
	}

	// 등록 이벤트 기록
	_ = s.eventRepo.LogEvent(ctx, &DeviceEvent{
		ID:        uuid.New().String(),
		DeviceID:  device.ID,
		EventType: "registered",
		Payload: map[string]interface{}{
			"serial_number":    serialNumber,
			"firmware_version": firmwareVersion,
		},
		CreatedAt: time.Now().UTC(),
	})

	s.logger.Info("디바이스 등록 완료",
		zap.String("device_id", device.ID),
		zap.String("user_id", userID),
		zap.String("serial", serialNumber),
	)

	// 디바이스 등록 이벤트 발행 (Kafka, optional)
	if s.eventPublisher != nil {
		evt := &DeviceRegisteredEvent{
			DeviceID:        device.ID,
			UserID:          userID,
			SerialNumber:    serialNumber,
			FirmwareVersion: firmwareVersion,
			RegisteredAt:    device.RegisteredAt,
		}
		if err := s.eventPublisher.PublishDeviceRegistered(ctx, evt); err != nil {
			s.logger.Warn("디바이스 등록 이벤트 발행 실패 (비치명적)", zap.Error(err))
		}
	}

	return device, registrationToken, nil
}

// ListDevices는 사용자의 디바이스 목록을 조회합니다.
func (s *DeviceService) ListDevices(ctx context.Context, userID string) ([]*Device, error) {
	if userID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "user_id는 필수입니다")
	}
	return s.deviceRepo.ListByUser(ctx, userID)
}

// UpdateDeviceStatus는 디바이스 상태를 업데이트합니다.
func (s *DeviceService) UpdateDeviceStatus(
	ctx context.Context,
	deviceID string,
	status DeviceStatus,
	batteryPercent int,
) error {
	if deviceID == "" {
		return apperrors.New(apperrors.ErrInvalidInput, "device_id는 필수입니다")
	}

	now := time.Now().UTC()
	if err := s.deviceRepo.UpdateStatus(ctx, deviceID, status, batteryPercent, now); err != nil {
		return apperrors.New(apperrors.ErrInternal, "디바이스 상태 업데이트에 실패했습니다")
	}

	// 상태 변경 이벤트 기록
	_ = s.eventRepo.LogEvent(ctx, &DeviceEvent{
		ID:        uuid.New().String(),
		DeviceID:  deviceID,
		EventType: "status_changed",
		Payload: map[string]interface{}{
			"status":          string(status),
			"battery_percent": batteryPercent,
		},
		CreatedAt: now,
	})

	// 상태 변경 이벤트 발행 (Kafka, optional)
	if s.eventPublisher != nil {
		evt := &DeviceStatusChangedEvent{
			DeviceID:       deviceID,
			NewStatus:      string(status),
			BatteryPercent: batteryPercent,
			LastSeen:       now,
		}
		if err := s.eventPublisher.PublishDeviceStatusChanged(ctx, evt); err != nil {
			s.logger.Warn("디바이스 상태변경 이벤트 발행 실패 (비치명적)", zap.Error(err))
		}
	}

	return nil
}

// RequestOtaUpdate는 OTA 펌웨어 업데이트를 요청합니다.
func (s *DeviceService) RequestOtaUpdate(
	ctx context.Context,
	deviceID, targetVersion string,
) (updateID, downloadURL, checksum string, err error) {
	device, err := s.deviceRepo.GetByID(ctx, deviceID)
	if err != nil || device == nil {
		return "", "", "", apperrors.New(apperrors.ErrNotFound, "디바이스를 찾을 수 없습니다")
	}

	if device.FirmwareVersion == targetVersion {
		return "", "", "", apperrors.New(apperrors.ErrConflict, "이미 최신 버전입니다")
	}

	updateID = uuid.New().String()
	downloadURL = fmt.Sprintf("https://ota.manpasik.com/firmware/%s/%s.bin", targetVersion, device.SerialNumber)
	checksum = "sha256:pending" // TODO: 실제 체크섬 생성

	// OTA 시작 이벤트 기록
	_ = s.eventRepo.LogEvent(ctx, &DeviceEvent{
		ID:        uuid.New().String(),
		DeviceID:  deviceID,
		EventType: "ota_started",
		Payload: map[string]interface{}{
			"update_id":       updateID,
			"target_version":  targetVersion,
			"current_version": device.FirmwareVersion,
		},
		CreatedAt: time.Now().UTC(),
	})

	s.logger.Info("OTA 업데이트 요청",
		zap.String("device_id", deviceID),
		zap.String("target_version", targetVersion),
	)

	return updateID, downloadURL, checksum, nil
}
