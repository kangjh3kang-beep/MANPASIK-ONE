package events

import (
	"testing"
)

func TestKafkaEventBus_topic(t *testing.T) {
	tests := []struct {
		name        string
		topicPrefix string
		eventType   string
		want        string
	}{
		{
			name:        "default prefix",
			topicPrefix: "manpasik.",
			eventType:   "reservation.created",
			want:        "manpasik.reservation.created",
		},
		{
			name:        "custom prefix",
			topicPrefix: "test.",
			eventType:   "payment.completed",
			want:        "test.payment.completed",
		},
		{
			name:        "empty prefix",
			topicPrefix: "",
			eventType:   "order.created",
			want:        "order.created",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &KafkaEventBus{
				topicPrefix: tt.topicPrefix,
			}
			if got := k.topic(tt.eventType); got != tt.want {
				t.Errorf("topic() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNewKafkaEventBus_InvalidBroker(t *testing.T) {
	// Use a non-routable address to force a quick connection failure
	cfg := KafkaAdapterConfig{
		Brokers:     []string{"192.0.2.1:9092"}, // RFC 5737 TEST-NET, guaranteed unreachable
		GroupID:     "test-group",
		TopicPrefix: "test.",
	}

	_, err := NewKafkaEventBus(cfg)
	if err == nil {
		t.Fatal("expected error when connecting to invalid broker, got nil")
	}

	t.Logf("got expected error: %v", err)
}
