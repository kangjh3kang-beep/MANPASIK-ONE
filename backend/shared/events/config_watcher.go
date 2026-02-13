// config_watcher.go는 시스템 설정 변경 이벤트를 구독하여 서비스에 전달합니다.
package events

import (
	"context"
	"encoding/json"
	"log"
)

// ConfigChangedEvent는 설정 변경 이벤트 페이로드입니다.
const EventConfigChanged = "config.changed"

// ConfigChangeHandler는 설정 변경을 처리하는 콜백 함수입니다.
type ConfigChangeHandler func(key, newValue string) error

// ConfigWatcher는 설정 변경 이벤트를 감시하는 인터페이스입니다.
type ConfigWatcher interface {
	// Watch는 지정된 서비스의 설정 변경을 감시합니다.
	// serviceFilter가 비어 있으면 모든 설정 변경을 수신합니다.
	Watch(ctx context.Context, serviceFilter string, handler ConfigChangeHandler) error
	// Close는 감시를 중지합니다.
	Close() error
}

// EventBusConfigWatcher는 EventPublisher 기반 ConfigWatcher 구현입니다.
// KafkaEventBus 또는 인메모리 EventBus를 모두 지원합니다.
type EventBusConfigWatcher struct {
	publisher     EventPublisher
	serviceFilter string
	handler       ConfigChangeHandler
}

// NewEventBusConfigWatcher는 EventPublisher 기반 ConfigWatcher를 생성합니다.
func NewEventBusConfigWatcher(publisher EventPublisher) *EventBusConfigWatcher {
	return &EventBusConfigWatcher{
		publisher: publisher,
	}
}

// Watch는 설정 변경 이벤트를 구독합니다.
func (w *EventBusConfigWatcher) Watch(_ context.Context, serviceFilter string, handler ConfigChangeHandler) error {
	w.serviceFilter = serviceFilter
	w.handler = handler

	w.publisher.Subscribe(EventConfigChanged, func(_ context.Context, event Event) error {
		key, _ := event.Payload["key"].(string)
		newValue, _ := event.Payload["new_value"].(string)
		svcName, _ := event.Payload["service_name"].(string)

		// 서비스 필터: 빈 문자열이면 전체, "*"도 전체
		if w.serviceFilter != "" && w.serviceFilter != "*" && svcName != w.serviceFilter && svcName != "*" {
			return nil
		}

		if w.handler != nil {
			if err := w.handler(key, newValue); err != nil {
				log.Printf("[config_watcher] 핸들러 오류 (key=%s): %v", key, err)
			}
		}
		return nil
	})

	return nil
}

// Close는 감시를 중지합니다 (EventBus 기반은 별도 정리 불필요).
func (w *EventBusConfigWatcher) Close() error {
	return nil
}

// NoopConfigWatcher는 개발/테스트용 빈 구현입니다.
type NoopConfigWatcher struct{}

// NewNoopConfigWatcher는 NoopConfigWatcher를 생성합니다.
func NewNoopConfigWatcher() *NoopConfigWatcher {
	return &NoopConfigWatcher{}
}

// Watch는 아무 작업도 하지 않습니다.
func (w *NoopConfigWatcher) Watch(_ context.Context, _ string, _ ConfigChangeHandler) error {
	return nil
}

// Close는 아무 작업도 하지 않습니다.
func (w *NoopConfigWatcher) Close() error {
	return nil
}

// PublishConfigChanged는 설정 변경 이벤트를 발행하는 헬퍼 함수입니다.
func PublishConfigChanged(ctx context.Context, publisher EventPublisher, key, newValue, oldValue, changedBy, serviceName string) error {
	payload := map[string]interface{}{
		"key":          key,
		"new_value":    newValue,
		"old_value":    oldValue,
		"changed_by":   changedBy,
		"service_name": serviceName,
	}

	return publisher.Publish(ctx, Event{
		Type:    EventConfigChanged,
		Payload: payload,
	})
}

// ParseConfigChangedPayload는 이벤트 페이로드에서 설정 변경 정보를 추출합니다.
func ParseConfigChangedPayload(data []byte) (key, newValue, serviceName string, err error) {
	var payload map[string]interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		return "", "", "", err
	}
	key, _ = payload["key"].(string)
	newValue, _ = payload["new_value"].(string)
	serviceName, _ = payload["service_name"].(string)
	return key, newValue, serviceName, nil
}
