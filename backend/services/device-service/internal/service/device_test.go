package service

import (
	"context"
	"testing"
	"time"

	"go.uber.org/zap"
)

// =============================================================================
// 목(Mock) 저장소
// =============================================================================

type mockDeviceRepo struct {
	devices map[string]*Device
}

func newMockDeviceRepo() *mockDeviceRepo {
	return &mockDeviceRepo{devices: make(map[string]*Device)}
}

func (r *mockDeviceRepo) Create(_ context.Context, device *Device) error {
	r.devices[device.ID] = device
	return nil
}

func (r *mockDeviceRepo) GetByID(_ context.Context, deviceID string) (*Device, error) {
	d, ok := r.devices[deviceID]
	if !ok {
		return nil, nil
	}
	return d, nil
}

func (r *mockDeviceRepo) ListByUser(_ context.Context, userID string) ([]*Device, error) {
	var result []*Device
	for _, d := range r.devices {
		if d.UserID == userID {
			result = append(result, d)
		}
	}
	return result, nil
}

func (r *mockDeviceRepo) UpdateStatus(_ context.Context, deviceID string, st DeviceStatus, battery int, lastSeen time.Time) error {
	d, ok := r.devices[deviceID]
	if !ok {
		return nil
	}
	d.Status = st
	d.BatteryPercent = battery
	d.LastSeen = lastSeen
	return nil
}

func (r *mockDeviceRepo) CountByUser(_ context.Context, userID string) (int, error) {
	count := 0
	for _, d := range r.devices {
		if d.UserID == userID {
			count++
		}
	}
	return count, nil
}

func (r *mockDeviceRepo) Delete(_ context.Context, deviceID string) error {
	delete(r.devices, deviceID)
	return nil
}

type mockEventRepo struct {
	events []*DeviceEvent
}

func newMockEventRepo() *mockEventRepo {
	return &mockEventRepo{}
}

func (r *mockEventRepo) LogEvent(_ context.Context, event *DeviceEvent) error {
	r.events = append(r.events, event)
	return nil
}

type mockSubChecker struct {
	maxDevices int
}

func newMockSubChecker(max int) *mockSubChecker {
	return &mockSubChecker{maxDevices: max}
}

func (c *mockSubChecker) GetMaxDevices(_ context.Context, _ string) (int, error) {
	return c.maxDevices, nil
}

// =============================================================================
// 헬퍼
// =============================================================================

func newTestDeviceService(maxDevices int) (*DeviceService, *mockDeviceRepo, *mockEventRepo) {
	logger, _ := zap.NewDevelopment()
	deviceRepo := newMockDeviceRepo()
	eventRepo := newMockEventRepo()
	subChecker := newMockSubChecker(maxDevices)
	svc := NewDeviceService(logger, deviceRepo, eventRepo, subChecker)
	return svc, deviceRepo, eventRepo
}

// =============================================================================
// RegisterDevice 테스트
// =============================================================================

func TestRegisterDevice_성공(t *testing.T) {
	svc, _, eventRepo := newTestDeviceService(10)
	ctx := context.Background()

	device, regToken, err := svc.RegisterDevice(ctx, "BLE-AA:BB:CC", "SN-001234", "1.0.0", "user-1")
	if err != nil {
		t.Fatalf("디바이스 등록 실패: %v", err)
	}
	if device.DeviceID != "BLE-AA:BB:CC" {
		t.Errorf("디바이스 ID 불일치: got %s", device.DeviceID)
	}
	if device.SerialNumber != "SN-001234" {
		t.Errorf("시리얼 불일치: got %s", device.SerialNumber)
	}
	if regToken == "" {
		t.Error("등록 토큰이 비어있습니다")
	}
	if device.Status != StatusOnline {
		t.Errorf("초기 상태가 online이어야 합니다: got %s", device.Status)
	}
	if len(eventRepo.events) != 1 {
		t.Errorf("등록 이벤트가 1개 기록되어야 합니다: got %d", len(eventRepo.events))
	}
	if eventRepo.events[0].EventType != "registered" {
		t.Errorf("이벤트 타입이 'registered'여야 합니다: got %s", eventRepo.events[0].EventType)
	}
}

func TestRegisterDevice_빈_입력(t *testing.T) {
	svc, _, _ := newTestDeviceService(10)
	ctx := context.Background()

	_, _, err := svc.RegisterDevice(ctx, "", "SN-001", "1.0.0", "user-1")
	if err == nil {
		t.Fatal("빈 device_id에 대해 에러가 발생해야 합니다")
	}

	_, _, err = svc.RegisterDevice(ctx, "BLE-01", "SN-001", "1.0.0", "")
	if err == nil {
		t.Fatal("빈 user_id에 대해 에러가 발생해야 합니다")
	}
}

func TestRegisterDevice_디바이스_제한_초과(t *testing.T) {
	svc, _, _ := newTestDeviceService(2) // 최대 2대
	ctx := context.Background()

	// 2대 등록 성공
	_, _, err := svc.RegisterDevice(ctx, "BLE-01", "SN-0001", "1.0.0", "user-1")
	if err != nil {
		t.Fatalf("첫 번째 디바이스 등록 실패: %v", err)
	}
	_, _, err = svc.RegisterDevice(ctx, "BLE-02", "SN-0002", "1.0.0", "user-1")
	if err != nil {
		t.Fatalf("두 번째 디바이스 등록 실패: %v", err)
	}

	// 3번째 등록 → 제한 초과
	_, _, err = svc.RegisterDevice(ctx, "BLE-03", "SN-0003", "1.0.0", "user-1")
	if err == nil {
		t.Fatal("디바이스 제한 초과 시 에러가 발생해야 합니다")
	}
}

// =============================================================================
// ListDevices 테스트
// =============================================================================

func TestListDevices_성공(t *testing.T) {
	svc, _, _ := newTestDeviceService(10)
	ctx := context.Background()

	svc.RegisterDevice(ctx, "BLE-01", "SN-0001", "1.0.0", "user-1")
	svc.RegisterDevice(ctx, "BLE-02", "SN-0002", "1.0.0", "user-1")
	svc.RegisterDevice(ctx, "BLE-03", "SN-0003", "1.0.0", "user-2") // 다른 사용자

	devices, err := svc.ListDevices(ctx, "user-1")
	if err != nil {
		t.Fatalf("디바이스 목록 조회 실패: %v", err)
	}
	if len(devices) != 2 {
		t.Errorf("user-1의 디바이스는 2개여야 합니다: got %d", len(devices))
	}
}

func TestListDevices_빈_유저ID(t *testing.T) {
	svc, _, _ := newTestDeviceService(10)
	ctx := context.Background()

	_, err := svc.ListDevices(ctx, "")
	if err == nil {
		t.Fatal("빈 user_id에 대해 에러가 발생해야 합니다")
	}
}

// =============================================================================
// UpdateDeviceStatus 테스트
// =============================================================================

func TestUpdateDeviceStatus_성공(t *testing.T) {
	svc, deviceRepo, eventRepo := newTestDeviceService(10)
	ctx := context.Background()

	device, _, _ := svc.RegisterDevice(ctx, "BLE-01", "SN-0001", "1.0.0", "user-1")

	err := svc.UpdateDeviceStatus(ctx, device.ID, StatusMeasuring, 85)
	if err != nil {
		t.Fatalf("상태 업데이트 실패: %v", err)
	}

	updated := deviceRepo.devices[device.ID]
	if updated.Status != StatusMeasuring {
		t.Errorf("상태 불일치: got %s, want measuring", updated.Status)
	}
	if updated.BatteryPercent != 85 {
		t.Errorf("배터리 불일치: got %d, want 85", updated.BatteryPercent)
	}
	// 등록 이벤트 + 상태 변경 이벤트 = 2
	if len(eventRepo.events) != 2 {
		t.Errorf("이벤트 2개 기록되어야 합니다: got %d", len(eventRepo.events))
	}
}

func TestUpdateDeviceStatus_빈_디바이스ID(t *testing.T) {
	svc, _, _ := newTestDeviceService(10)
	ctx := context.Background()

	err := svc.UpdateDeviceStatus(ctx, "", StatusOnline, 100)
	if err == nil {
		t.Fatal("빈 device_id에 대해 에러가 발생해야 합니다")
	}
}

// =============================================================================
// RequestOtaUpdate 테스트
// =============================================================================

func TestRequestOtaUpdate_성공(t *testing.T) {
	svc, _, _ := newTestDeviceService(10)
	ctx := context.Background()

	device, _, _ := svc.RegisterDevice(ctx, "BLE-01", "SN-0001", "1.0.0", "user-1")

	updateID, downloadURL, checksum, err := svc.RequestOtaUpdate(ctx, device.ID, "2.0.0")
	if err != nil {
		t.Fatalf("OTA 업데이트 요청 실패: %v", err)
	}
	if updateID == "" {
		t.Error("updateID가 비어있습니다")
	}
	if downloadURL == "" {
		t.Error("downloadURL이 비어있습니다")
	}
	if checksum == "" {
		t.Error("checksum이 비어있습니다")
	}
}

func TestRequestOtaUpdate_이미_최신_버전(t *testing.T) {
	svc, _, _ := newTestDeviceService(10)
	ctx := context.Background()

	device, _, _ := svc.RegisterDevice(ctx, "BLE-01", "SN-0001", "1.0.0", "user-1")

	_, _, _, err := svc.RequestOtaUpdate(ctx, device.ID, "1.0.0")
	if err == nil {
		t.Fatal("동일 버전 업데이트 시 에러가 발생해야 합니다")
	}
}

func TestRequestOtaUpdate_존재하지_않는_디바이스(t *testing.T) {
	svc, _, _ := newTestDeviceService(10)
	ctx := context.Background()

	_, _, _, err := svc.RequestOtaUpdate(ctx, "nonexistent", "2.0.0")
	if err == nil {
		t.Fatal("존재하지 않는 디바이스에 대해 에러가 발생해야 합니다")
	}
}
