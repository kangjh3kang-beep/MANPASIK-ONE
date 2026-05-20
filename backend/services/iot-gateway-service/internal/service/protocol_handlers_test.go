package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/manpasik/backend/services/iot-gateway-service/internal/service"
)

// TestMQTTHandler_SendAndStats는 MQTT 송신과 통계를 검증합니다.
func TestMQTTHandler_SendAndStats(t *testing.T) {
	h := service.NewMQTTHandler()
	ctx := context.Background()

	if err := h.Send(ctx, "dev-001", []byte("ping")); err != nil {
		t.Fatalf("Send 실패: %v", err)
	}

	stats := h.Stats()
	if stats.MessagesSent != 1 {
		t.Errorf("MessagesSent = %d, want 1", stats.MessagesSent)
	}
	if stats.Protocol != "mqtt" {
		t.Errorf("Protocol = %q, want mqtt", stats.Protocol)
	}
}

// TestMQTTHandler_PayloadTooLarge는 페이로드 크기 제한을 검증합니다.
func TestMQTTHandler_PayloadTooLarge(t *testing.T) {
	h := service.NewMQTTHandler()
	big := make([]byte, 300*1024) // 300KB

	err := h.Send(context.Background(), "dev", big)
	if err == nil {
		t.Error("256KB 초과 페이로드가 허용됨")
	}

	stats := h.Stats()
	if stats.Errors == 0 {
		t.Error("에러 통계가 증가하지 않음")
	}
}

// TestMQTTHandler_SubscribeReceive는 구독/수신을 검증합니다.
func TestMQTTHandler_SubscribeReceive(t *testing.T) {
	h := service.NewMQTTHandler()
	received := make(chan []byte, 1)

	_ = h.Subscribe("dev", func(_ string, payload []byte) {
		received <- payload
	})

	h.SimulateReceive("dev", []byte("data"))

	select {
	case got := <-received:
		if string(got) != "data" {
			t.Errorf("payload = %q, want data", got)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("메시지 수신 타임아웃")
	}
}

// TestBLEHandler_ConnectAndSend는 BLE 연결 + 송신을 검증합니다.
func TestBLEHandler_ConnectAndSend(t *testing.T) {
	h := service.NewBLEHandler()

	if err := h.Connect("ble-001"); err != nil {
		t.Fatalf("Connect 실패: %v", err)
	}
	if !h.IsConnected("ble-001") {
		t.Error("연결 상태 false")
	}

	if err := h.Send(context.Background(), "ble-001", []byte("data")); err != nil {
		t.Fatalf("Send 실패: %v", err)
	}
}

// TestBLEHandler_DisconnectedSend는 미연결 디바이스에 송신 거부를 검증합니다.
func TestBLEHandler_DisconnectedSend(t *testing.T) {
	h := service.NewBLEHandler()
	err := h.Send(context.Background(), "unknown", []byte("data"))
	if err == nil {
		t.Error("미연결 디바이스에 송신이 허용됨")
	}
}

// TestBLEHandler_PayloadLimit는 BLE MTU 제한을 검증합니다.
func TestBLEHandler_PayloadLimit(t *testing.T) {
	h := service.NewBLEHandler()
	_ = h.Connect("dev")

	big := make([]byte, 600) // 512 초과
	err := h.Send(context.Background(), "dev", big)
	if err == nil {
		t.Error("512바이트 초과 페이로드가 허용됨")
	}
}

// TestCoAPHandler_PUTGET은 CoAP PUT/GET을 검증합니다.
func TestCoAPHandler_PUTGET(t *testing.T) {
	h := service.NewCoAPHandler()

	if err := h.PUT("dev", "/sensor/temp", []byte("36.5")); err != nil {
		t.Fatalf("PUT 실패: %v", err)
	}

	got, err := h.GET("dev", "/sensor/temp")
	if err != nil {
		t.Fatalf("GET 실패: %v", err)
	}
	if string(got) != "36.5" {
		t.Errorf("payload = %q, want 36.5", got)
	}
}

// TestCoAPHandler_DELETE는 CoAP DELETE를 검증합니다.
func TestCoAPHandler_DELETE(t *testing.T) {
	h := service.NewCoAPHandler()
	_ = h.PUT("dev", "/x", []byte("data"))
	_ = h.DELETE("dev", "/x")

	_, err := h.GET("dev", "/x")
	if err == nil {
		t.Error("DELETE 후 GET이 성공함")
	}
}

// TestCoAPHandler_NotFound는 미등록 리소스 GET 시 4.04 응답을 검증합니다.
func TestCoAPHandler_NotFound(t *testing.T) {
	h := service.NewCoAPHandler()
	_, err := h.GET("dev", "/missing")
	if err == nil {
		t.Error("미등록 리소스 GET이 성공함")
	}
}

// TestCommandQueue_Priority는 우선순위 큐 정렬을 검증합니다.
func TestCommandQueue_Priority(t *testing.T) {
	q := service.NewCommandQueue()
	q.Enqueue(&service.QueuedCommand{ID: "low", Priority: 1, DeviceID: "d"})
	q.Enqueue(&service.QueuedCommand{ID: "urgent", Priority: 4, DeviceID: "d"})
	q.Enqueue(&service.QueuedCommand{ID: "high", Priority: 3, DeviceID: "d"})

	first := q.DequeueByPriority()
	if first == nil || first.ID != "urgent" {
		t.Errorf("first = %v, want urgent", first)
	}

	second := q.DequeueByPriority()
	if second == nil || second.ID != "high" {
		t.Errorf("second = %v, want high", second)
	}
}

// TestCommandQueue_TTL은 TTL 만료 명령을 자동 제거하는지 검증합니다.
func TestCommandQueue_TTL(t *testing.T) {
	q := service.NewCommandQueue()
	q.Enqueue(&service.QueuedCommand{ID: "fresh", Priority: 2, DeviceID: "d", TTL: 1 * time.Hour})
	q.Enqueue(&service.QueuedCommand{ID: "stale", Priority: 4, DeviceID: "d", TTL: 1 * time.Nanosecond})

	time.Sleep(2 * time.Millisecond)

	got := q.DequeueByPriority()
	if got == nil || got.ID != "fresh" {
		t.Errorf("got = %v, want fresh (stale should be expired)", got)
	}
}

// TestCommandQueue_PurgeExpired는 명시적 purge를 검증합니다.
func TestCommandQueue_PurgeExpired(t *testing.T) {
	q := service.NewCommandQueue()
	for i := 0; i < 5; i++ {
		q.Enqueue(&service.QueuedCommand{
			ID: "x", Priority: 1, DeviceID: "d", TTL: 1 * time.Nanosecond,
		})
	}
	time.Sleep(2 * time.Millisecond)

	purged := q.PurgeExpired()
	if purged != 5 {
		t.Errorf("purged = %d, want 5", purged)
	}
	if q.Size() != 0 {
		t.Errorf("Size = %d, want 0", q.Size())
	}
}

// TestCommandQueue_SizeByDevice는 디바이스별 큐 크기를 검증합니다.
func TestCommandQueue_SizeByDevice(t *testing.T) {
	q := service.NewCommandQueue()
	q.Enqueue(&service.QueuedCommand{ID: "1", DeviceID: "d1", Priority: 1})
	q.Enqueue(&service.QueuedCommand{ID: "2", DeviceID: "d1", Priority: 1})
	q.Enqueue(&service.QueuedCommand{ID: "3", DeviceID: "d2", Priority: 1})

	if got := q.SizeByDevice("d1"); got != 2 {
		t.Errorf("d1 size = %d, want 2", got)
	}
	if got := q.SizeByDevice("d2"); got != 1 {
		t.Errorf("d2 size = %d, want 1", got)
	}
}

// TestConnectionMonitor_Heartbeat는 heartbeat 기록과 alive 체크를 검증합니다.
func TestConnectionMonitor_Heartbeat(t *testing.T) {
	monitor := service.NewConnectionMonitor(100 * time.Millisecond)
	monitor.RecordHeartbeat("dev-1")

	if !monitor.IsAlive("dev-1") {
		t.Error("dev-1이 alive가 아님")
	}

	time.Sleep(150 * time.Millisecond)
	if monitor.IsAlive("dev-1") {
		t.Error("타임아웃 후에도 alive로 표시됨")
	}
}

// TestConnectionMonitor_FindDead는 dead 디바이스 탐지를 검증합니다.
func TestConnectionMonitor_FindDead(t *testing.T) {
	monitor := service.NewConnectionMonitor(50 * time.Millisecond)
	monitor.RecordHeartbeat("alive")
	monitor.RecordHeartbeat("dead-1")
	monitor.RecordHeartbeat("dead-2")

	time.Sleep(100 * time.Millisecond)
	monitor.RecordHeartbeat("alive") // alive만 갱신

	dead := monitor.FindDeadDevices()
	if len(dead) != 2 {
		t.Errorf("dead = %v, want 2", dead)
	}
}

// TestConnectionMonitor_Cleanup는 cleanup이 매우 오래된 데드 디바이스를 제거하는지 검증합니다.
func TestConnectionMonitor_Cleanup(t *testing.T) {
	monitor := service.NewConnectionMonitor(20 * time.Millisecond)
	monitor.RecordHeartbeat("d1")
	monitor.RecordHeartbeat("d2")

	time.Sleep(100 * time.Millisecond) // 5x timeout

	cleaned := monitor.Cleanup()
	if cleaned != 2 {
		t.Errorf("cleaned = %d, want 2", cleaned)
	}
}

// TestProtocolRouter_RouteByProtocol은 라우팅을 검증합니다.
func TestProtocolRouter_RouteByProtocol(t *testing.T) {
	router := service.NewProtocolRouter()
	mqtt := service.NewMQTTHandler()
	ble := service.NewBLEHandler()

	router.Register(mqtt)
	router.Register(ble)
	_ = router.AssignDevice("d-mqtt", "mqtt")
	_ = ble.Connect("d-ble")
	_ = router.AssignDevice("d-ble", "ble")

	if err := router.Send(context.Background(), "d-mqtt", []byte("hi")); err != nil {
		t.Errorf("MQTT 라우팅 실패: %v", err)
	}
	if err := router.Send(context.Background(), "d-ble", []byte("hi")); err != nil {
		t.Errorf("BLE 라우팅 실패: %v", err)
	}

	// 미할당 디바이스
	if err := router.Send(context.Background(), "unknown", []byte("hi")); err == nil {
		t.Error("미할당 디바이스 라우팅이 성공함")
	}
}

// TestProtocolRouter_AggregateStats는 통계 집계를 검증합니다.
func TestProtocolRouter_AggregateStats(t *testing.T) {
	router := service.NewProtocolRouter()
	router.Register(service.NewMQTTHandler())
	router.Register(service.NewCoAPHandler())

	stats := router.AggregateStats()
	if len(stats) != 2 {
		t.Errorf("프로토콜 수 = %d, want 2", len(stats))
	}
	if _, ok := stats["mqtt"]; !ok {
		t.Error("mqtt 통계 미집계")
	}
}
