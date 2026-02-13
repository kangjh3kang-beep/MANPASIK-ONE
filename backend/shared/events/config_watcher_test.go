package events

import (
	"context"
	"testing"
	"time"
)

func TestEventBusConfigWatcher_Watch(t *testing.T) {
	bus := NewEventBus()
	watcher := NewEventBusConfigWatcher(bus)

	received := make(chan string, 1)
	err := watcher.Watch(context.Background(), "payment-service", func(key, newValue string) error {
		received <- key + "=" + newValue
		return nil
	})
	if err != nil {
		t.Fatalf("Watch 실패: %v", err)
	}

	// 이벤트 발행: 대상 서비스 일치
	_ = PublishConfigChanged(context.Background(), bus, "toss.secret_key", "new_secret", "old_secret", "admin", "payment-service")

	select {
	case val := <-received:
		if val != "toss.secret_key=new_secret" {
			t.Errorf("예상과 다른 값: %s", val)
		}
	case <-time.After(time.Second):
		t.Error("1초 내 이벤트 수신 실패")
	}
}

func TestEventBusConfigWatcher_ServiceFilter(t *testing.T) {
	bus := NewEventBus()
	watcher := NewEventBusConfigWatcher(bus)

	received := make(chan string, 1)
	_ = watcher.Watch(context.Background(), "payment-service", func(key, newValue string) error {
		received <- key
		return nil
	})

	// 다른 서비스 대상 이벤트는 무시
	_ = PublishConfigChanged(context.Background(), bus, "fcm.server_key", "val", "", "admin", "notification-service")

	select {
	case val := <-received:
		t.Errorf("필터링되어야 하는 이벤트 수신: %s", val)
	case <-time.After(200 * time.Millisecond):
		// 예상 동작: 수신 안 함
	}
}

func TestEventBusConfigWatcher_WildcardService(t *testing.T) {
	bus := NewEventBus()
	watcher := NewEventBusConfigWatcher(bus)

	received := make(chan string, 1)
	_ = watcher.Watch(context.Background(), "any-service", func(key, newValue string) error {
		received <- key
		return nil
	})

	// service_name이 "*"이면 모든 서비스가 수신
	_ = PublishConfigChanged(context.Background(), bus, "maintenance_mode", "true", "", "admin", "*")

	select {
	case val := <-received:
		if val != "maintenance_mode" {
			t.Errorf("예상과 다른 키: %s", val)
		}
	case <-time.After(time.Second):
		t.Error("와일드카드 이벤트 수신 실패")
	}
}

func TestNoopConfigWatcher(t *testing.T) {
	watcher := NewNoopConfigWatcher()
	err := watcher.Watch(context.Background(), "test", func(key, newValue string) error {
		t.Error("NoopWatcher에서 핸들러가 호출되면 안 됩니다")
		return nil
	})
	if err != nil {
		t.Fatalf("NoopWatcher Watch 실패: %v", err)
	}
	if err := watcher.Close(); err != nil {
		t.Fatalf("NoopWatcher Close 실패: %v", err)
	}
}

func TestPublishConfigChanged(t *testing.T) {
	bus := NewEventBus()

	received := false
	bus.Subscribe(EventConfigChanged, func(_ context.Context, event Event) error {
		received = true
		if event.Payload["key"] != "test.key" {
			t.Errorf("key 불일치: %v", event.Payload["key"])
		}
		if event.Payload["new_value"] != "new_val" {
			t.Errorf("new_value 불일치: %v", event.Payload["new_value"])
		}
		if event.Payload["changed_by"] != "admin-001" {
			t.Errorf("changed_by 불일치: %v", event.Payload["changed_by"])
		}
		return nil
	})

	err := PublishConfigChanged(context.Background(), bus, "test.key", "new_val", "old_val", "admin-001", "test-service")
	if err != nil {
		t.Fatalf("PublishConfigChanged 실패: %v", err)
	}
	if !received {
		t.Error("이벤트가 수신되지 않았습니다")
	}
}
