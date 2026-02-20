package service_test

import (
	"context"
	"testing"

	"github.com/manpasik/backend/services/iot-gateway-service/internal/repository/memory"
	"github.com/manpasik/backend/services/iot-gateway-service/internal/service"
)

// newTestService는 테스트용 IoTGatewayService를 생성합니다.
func newTestService() *service.IoTGatewayService {
	repo := memory.NewIoTRepository()
	return service.NewIoTGatewayService(repo)
}

// ============================================================================
// 디바이스 등록 테스트
// ============================================================================

func TestRegisterDevice_Success(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	meta := map[string]string{
		"firmware": "v1.2.3",
		"model":    "MPS-100",
	}

	device, err := svc.RegisterDevice(ctx, "device-001", "mqtt", meta)
	if err != nil {
		t.Fatalf("RegisterDevice 실패: %v", err)
	}
	if device.ID == "" {
		t.Error("ID가 비어 있습니다")
	}
	if device.DeviceID != "device-001" {
		t.Errorf("DeviceID: got %s, want device-001", device.DeviceID)
	}
	if device.Protocol != "mqtt" {
		t.Errorf("Protocol: got %s, want mqtt", device.Protocol)
	}
	if device.Status != "registered" {
		t.Errorf("Status: got %s, want registered", device.Status)
	}
	if device.Metadata["firmware"] != "v1.2.3" {
		t.Errorf("Metadata[firmware]: got %s, want v1.2.3", device.Metadata["firmware"])
	}
	if device.LastPingAt.IsZero() {
		t.Error("LastPingAt가 zero입니다")
	}
}

func TestRegisterDevice_MissingDeviceID(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, err := svc.RegisterDevice(ctx, "", "mqtt", nil)
	if err == nil {
		t.Error("device_id 없이 등록이 허용되었습니다")
	}
}

func TestRegisterDevice_MissingProtocol(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, err := svc.RegisterDevice(ctx, "device-002", "", nil)
	if err == nil {
		t.Error("protocol 없이 등록이 허용되었습니다")
	}
}

// ============================================================================
// 명령 전송 테스트
// ============================================================================

func TestSendCommand_Success(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	// 먼저 디바이스 등록
	_, err := svc.RegisterDevice(ctx, "device-cmd-001", "mqtt", nil)
	if err != nil {
		t.Fatalf("디바이스 등록 실패: %v", err)
	}

	cmd, err := svc.SendCommand(ctx, "device-cmd-001", "reboot", `{"delay":5}`)
	if err != nil {
		t.Fatalf("SendCommand 실패: %v", err)
	}
	if cmd.ID == "" {
		t.Error("Command ID가 비어 있습니다")
	}
	if cmd.DeviceID != "device-cmd-001" {
		t.Errorf("DeviceID: got %s, want device-cmd-001", cmd.DeviceID)
	}
	if cmd.CommandType != "reboot" {
		t.Errorf("CommandType: got %s, want reboot", cmd.CommandType)
	}
	if cmd.Payload != `{"delay":5}` {
		t.Errorf("Payload: got %s, want {\"delay\":5}", cmd.Payload)
	}
	if cmd.Status != "pending" {
		t.Errorf("Status: got %s, want pending", cmd.Status)
	}
	if cmd.CreatedAt.IsZero() {
		t.Error("CreatedAt가 zero입니다")
	}
}

func TestSendCommand_DeviceNotFound(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, err := svc.SendCommand(ctx, "non-existent-device", "reboot", "")
	if err == nil {
		t.Error("존재하지 않는 디바이스에 명령 전송이 허용되었습니다")
	}
}

func TestSendCommand_MissingDeviceID(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, err := svc.SendCommand(ctx, "", "reboot", "")
	if err == nil {
		t.Error("device_id 없이 명령 전송이 허용되었습니다")
	}
}

func TestSendCommand_MissingCommandType(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, _ = svc.RegisterDevice(ctx, "device-cmd-002", "mqtt", nil)

	_, err := svc.SendCommand(ctx, "device-cmd-002", "", "payload")
	if err == nil {
		t.Error("command_type 없이 명령 전송이 허용되었습니다")
	}
}

// ============================================================================
// 데이터 수신 테스트
// ============================================================================

func TestReceiveData_Success(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	data, err := svc.ReceiveData(ctx, "device-data-001", "temperature", 36.5, "celsius")
	if err != nil {
		t.Fatalf("ReceiveData 실패: %v", err)
	}
	if data.ID == "" {
		t.Error("Data ID가 비어 있습니다")
	}
	if data.DeviceID != "device-data-001" {
		t.Errorf("DeviceID: got %s, want device-data-001", data.DeviceID)
	}
	if data.DataType != "temperature" {
		t.Errorf("DataType: got %s, want temperature", data.DataType)
	}
	if data.Value != 36.5 {
		t.Errorf("Value: got %f, want 36.5", data.Value)
	}
	if data.Unit != "celsius" {
		t.Errorf("Unit: got %s, want celsius", data.Unit)
	}
	if data.ReceivedAt.IsZero() {
		t.Error("ReceivedAt가 zero입니다")
	}
}

func TestReceiveData_MissingDeviceID(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, err := svc.ReceiveData(ctx, "", "temperature", 36.5, "celsius")
	if err == nil {
		t.Error("device_id 없이 데이터 수신이 허용되었습니다")
	}
}

func TestReceiveData_MissingDataType(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, err := svc.ReceiveData(ctx, "device-data-002", "", 36.5, "celsius")
	if err == nil {
		t.Error("data_type 없이 데이터 수신이 허용되었습니다")
	}
}

// ============================================================================
// 디바이스 조회 테스트
// ============================================================================

func TestGetDevice_Success(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, err := svc.RegisterDevice(ctx, "device-get-001", "ble", map[string]string{"version": "2.0"})
	if err != nil {
		t.Fatalf("디바이스 등록 실패: %v", err)
	}

	device, err := svc.GetDevice(ctx, "device-get-001")
	if err != nil {
		t.Fatalf("GetDevice 실패: %v", err)
	}
	if device.DeviceID != "device-get-001" {
		t.Errorf("DeviceID: got %s, want device-get-001", device.DeviceID)
	}
	if device.Protocol != "ble" {
		t.Errorf("Protocol: got %s, want ble", device.Protocol)
	}
}

func TestGetDevice_NotFound(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, err := svc.GetDevice(ctx, "non-existent")
	if err == nil {
		t.Error("존재하지 않는 디바이스 조회가 성공했습니다")
	}
}

// ============================================================================
// 데이터 목록 조회 테스트
// ============================================================================

func TestListData_Success(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	// 여러 데이터 수신
	_, _ = svc.ReceiveData(ctx, "device-list-001", "temperature", 36.5, "celsius")
	_, _ = svc.ReceiveData(ctx, "device-list-001", "heart_rate", 72.0, "bpm")
	_, _ = svc.ReceiveData(ctx, "device-list-001", "spo2", 98.0, "percent")
	_, _ = svc.ReceiveData(ctx, "device-list-002", "temperature", 37.0, "celsius") // 다른 디바이스

	dataList, err := svc.ListData(ctx, "device-list-001", 10, 0)
	if err != nil {
		t.Fatalf("ListData 실패: %v", err)
	}
	if len(dataList) != 3 {
		t.Errorf("데이터 수: got %d, want 3", len(dataList))
	}

	// 페이지네이션 테스트
	dataList, err = svc.ListData(ctx, "device-list-001", 2, 0)
	if err != nil {
		t.Fatalf("ListData(limit=2) 실패: %v", err)
	}
	if len(dataList) != 2 {
		t.Errorf("페이지네이션 데이터 수: got %d, want 2", len(dataList))
	}
}

// ============================================================================
// E2E 통합 테스트
// ============================================================================

func TestEndToEnd_IoTFlow(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	// 1. 디바이스 등록
	device, err := svc.RegisterDevice(ctx, "e2e-device", "mqtt", map[string]string{
		"model": "MPS-200",
	})
	if err != nil {
		t.Fatalf("E2E: 디바이스 등록 실패: %v", err)
	}
	if device.Status != "registered" {
		t.Errorf("E2E: Status: got %s, want registered", device.Status)
	}

	// 2. 디바이스 조회
	fetched, err := svc.GetDevice(ctx, "e2e-device")
	if err != nil {
		t.Fatalf("E2E: 디바이스 조회 실패: %v", err)
	}
	if fetched.Protocol != "mqtt" {
		t.Errorf("E2E: Protocol: got %s, want mqtt", fetched.Protocol)
	}

	// 3. 명령 전송
	cmd, err := svc.SendCommand(ctx, "e2e-device", "measure", `{"type":"blood_glucose"}`)
	if err != nil {
		t.Fatalf("E2E: 명령 전송 실패: %v", err)
	}
	if cmd.Status != "pending" {
		t.Errorf("E2E: Command Status: got %s, want pending", cmd.Status)
	}

	// 4. 명령 조회
	fetchedCmd, err := svc.GetCommand(ctx, cmd.ID)
	if err != nil {
		t.Fatalf("E2E: 명령 조회 실패: %v", err)
	}
	if fetchedCmd.CommandType != "measure" {
		t.Errorf("E2E: CommandType: got %s, want measure", fetchedCmd.CommandType)
	}

	// 5. 데이터 수신
	data, err := svc.ReceiveData(ctx, "e2e-device", "blood_glucose", 95.0, "mg/dL")
	if err != nil {
		t.Fatalf("E2E: 데이터 수신 실패: %v", err)
	}
	if data.Value != 95.0 {
		t.Errorf("E2E: Value: got %f, want 95.0", data.Value)
	}

	// 6. 데이터 목록 조회
	dataList, err := svc.ListData(ctx, "e2e-device", 10, 0)
	if err != nil {
		t.Fatalf("E2E: 데이터 목록 조회 실패: %v", err)
	}
	if len(dataList) != 1 {
		t.Errorf("E2E: 데이터 수: got %d, want 1", len(dataList))
	}
}
