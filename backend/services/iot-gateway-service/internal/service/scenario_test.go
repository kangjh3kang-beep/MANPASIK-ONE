package service_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/manpasik/backend/services/iot-gateway-service/internal/service"
)

// TestScenario_MultiProtocolFleet는 다중 프로토콜 디바이스 플릿 관리를 검증합니다.
func TestScenario_MultiProtocolFleet(t *testing.T) {
	router := service.NewProtocolRouter()
	mqtt := service.NewMQTTHandler()
	ble := service.NewBLEHandler()
	coap := service.NewCoAPHandler()

	router.Register(mqtt)
	router.Register(ble)
	router.Register(coap)

	// 3종 프로토콜 디바이스 10개씩 등록
	for i := 0; i < 10; i++ {
		mqttID := "mqtt-" + string(rune('a'+i))
		bleID := "ble-" + string(rune('a'+i))
		coapID := "coap-" + string(rune('a'+i))

		_ = router.AssignDevice(mqttID, "mqtt")
		_ = ble.Connect(bleID)
		_ = router.AssignDevice(bleID, "ble")
		_ = router.AssignDevice(coapID, "coap")
	}

	stats := router.AggregateStats()
	if len(stats) != 3 {
		t.Errorf("프로토콜 수 = %d, want 3", len(stats))
	}
}

// TestScenario_HighPriorityCommandRouting는 긴급 명령 우선 처리를 검증합니다.
func TestScenario_HighPriorityCommandRouting(t *testing.T) {
	queue := service.NewCommandQueue()

	// 다양한 우선순위 명령
	queue.Enqueue(&service.QueuedCommand{ID: "low-1", DeviceID: "d1", Priority: 1, Payload: []byte("status")})
	queue.Enqueue(&service.QueuedCommand{ID: "urgent-stop", DeviceID: "d1", Priority: 4, Payload: []byte("emergency_stop"), TTL: 5 * time.Second})
	queue.Enqueue(&service.QueuedCommand{ID: "high-1", DeviceID: "d1", Priority: 3, Payload: []byte("calibrate")})
	queue.Enqueue(&service.QueuedCommand{ID: "urgent-2", DeviceID: "d1", Priority: 4, Payload: []byte("alarm")})

	// urgent 두 개가 먼저 처리
	first := queue.DequeueByPriority()
	second := queue.DequeueByPriority()

	if first.Priority != 4 || second.Priority != 4 {
		t.Errorf("우선순위 정렬 위반: first=%d, second=%d", first.Priority, second.Priority)
	}
}

// TestScenario_BLECommandRetryFlow는 BLE 명령 재시도 흐름을 검증합니다.
func TestScenario_BLECommandRetryFlow(t *testing.T) {
	ble := service.NewBLEHandler()

	// 미연결 디바이스 → 송신 거부
	err := ble.Send(context.Background(), "unknown-device", []byte("data"))
	if err == nil {
		t.Error("미연결 디바이스 송신이 허용됨")
	}

	// 연결 후 정상 송신
	_ = ble.Connect("connected-1")
	if err := ble.Send(context.Background(), "connected-1", []byte("data")); err != nil {
		t.Errorf("연결된 디바이스 송신 실패: %v", err)
	}

	// 연결 해제 후 재시도
	_ = ble.Disconnect("connected-1")
	if ble.IsConnected("connected-1") {
		t.Error("Disconnect 후에도 연결 상태")
	}
}

// TestScenario_HeartbeatTimeoutRecovery는 heartbeat 타임아웃 + 복구 흐름입니다.
func TestScenario_HeartbeatTimeoutRecovery(t *testing.T) {
	monitor := service.NewConnectionMonitor(50 * time.Millisecond)

	devices := []string{"dev-1", "dev-2", "dev-3"}
	for _, d := range devices {
		monitor.RecordHeartbeat(d)
	}

	// 모두 alive 상태
	for _, d := range devices {
		if !monitor.IsAlive(d) {
			t.Errorf("%s 초기 alive 아님", d)
		}
	}

	// 타임아웃
	time.Sleep(60 * time.Millisecond)

	// 모두 dead 상태
	dead := monitor.FindDeadDevices()
	if len(dead) != 3 {
		t.Errorf("dead = %d, want 3", len(dead))
	}

	// dev-1만 재연결
	monitor.RecordHeartbeat("dev-1")
	if !monitor.IsAlive("dev-1") {
		t.Error("재연결 후 alive 아님")
	}
	if monitor.IsAlive("dev-2") {
		t.Error("dev-2는 여전히 dead여야 함")
	}
}

// TestScenario_CommandQueueTTLExpiration는 TTL 만료 명령 자동 제거를 검증합니다.
func TestScenario_CommandQueueTTLExpiration(t *testing.T) {
	queue := service.NewCommandQueue()

	queue.Enqueue(&service.QueuedCommand{ID: "fresh", Priority: 4, DeviceID: "d", TTL: 1 * time.Hour})
	queue.Enqueue(&service.QueuedCommand{ID: "stale-1", Priority: 4, DeviceID: "d", TTL: 1 * time.Nanosecond})
	queue.Enqueue(&service.QueuedCommand{ID: "stale-2", Priority: 4, DeviceID: "d", TTL: 1 * time.Nanosecond})

	time.Sleep(10 * time.Millisecond)

	got := queue.DequeueByPriority()
	if got == nil || got.ID != "fresh" {
		t.Errorf("got = %v, want fresh", got)
	}

	if queue.Size() != 0 {
		t.Errorf("Size = %d, want 0", queue.Size())
	}
}

// TestScenario_CoAPRESTOperations는 CoAP REST 매핑 동작을 검증합니다.
func TestScenario_CoAPRESTOperations(t *testing.T) {
	coap := service.NewCoAPHandler()

	// PUT
	if err := coap.PUT("dev-1", "/sensor/temp", []byte("36.5")); err != nil {
		t.Fatalf("PUT 실패: %v", err)
	}

	// GET
	got, err := coap.GET("dev-1", "/sensor/temp")
	if err != nil {
		t.Fatalf("GET 실패: %v", err)
	}
	if string(got) != "36.5" {
		t.Errorf("payload = %q", got)
	}

	// POST (PUT과 동일)
	_ = coap.POST("dev-1", "/sensor/temp", []byte("37.0"))
	got2, _ := coap.GET("dev-1", "/sensor/temp")
	if string(got2) != "37.0" {
		t.Errorf("POST 후 payload = %q", got2)
	}

	// DELETE
	_ = coap.DELETE("dev-1", "/sensor/temp")
	if _, err := coap.GET("dev-1", "/sensor/temp"); err == nil {
		t.Error("DELETE 후 GET이 성공함")
	}
}

// TestScenario_ConcurrentMQTTOperations는 동시 MQTT 발송을 검증합니다.
func TestScenario_ConcurrentMQTTOperations(t *testing.T) {
	mqtt := service.NewMQTTHandler()

	var wg sync.WaitGroup
	successes := 0
	mu := sync.Mutex{}

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := mqtt.Send(context.Background(), "device-x", []byte("ping"))
			if err == nil {
				mu.Lock()
				successes++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	if successes != 50 {
		t.Errorf("successes = %d, want 50", successes)
	}

	stats := mqtt.Stats()
	if stats.MessagesSent != 50 {
		t.Errorf("MessagesSent = %d, want 50", stats.MessagesSent)
	}
}
