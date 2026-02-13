package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

// KafkaEventBus implements event publishing and subscribing via Kafka/Redpanda
type KafkaEventBus struct {
	mu          sync.RWMutex
	client      *kgo.Client
	handlers    map[string][]Handler
	topicPrefix string
	groupID     string
	cancel      context.CancelFunc
}

// KafkaAdapterConfig holds Kafka connection configuration for the event bus adapter
type KafkaAdapterConfig struct {
	Brokers     []string
	GroupID     string
	TopicPrefix string
}

// NewKafkaEventBus creates a Kafka-backed event bus
func NewKafkaEventBus(cfg KafkaAdapterConfig) (*KafkaEventBus, error) {
	if cfg.TopicPrefix == "" {
		cfg.TopicPrefix = "manpasik."
	}
	if cfg.GroupID == "" {
		cfg.GroupID = "manpasik"
	}

	client, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.Brokers...),
		kgo.ConsumerGroup(cfg.GroupID),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtEnd()),
		kgo.ProducerBatchMaxBytes(1024*1024), // 1MB
		kgo.RecordDeliveryTimeout(10*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka client creation failed: %w", err)
	}

	// Verify connectivity
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx); err != nil {
		client.Close()
		return nil, fmt.Errorf("kafka ping failed: %w", err)
	}

	return &KafkaEventBus{
		client:      client,
		handlers:    make(map[string][]Handler),
		topicPrefix: cfg.TopicPrefix,
		groupID:     cfg.GroupID,
	}, nil
}

// topic converts event type to Kafka topic name
func (k *KafkaEventBus) topic(eventType string) string {
	return k.topicPrefix + eventType
}

// Subscribe registers a handler for an event type
func (k *KafkaEventBus) Subscribe(eventType string, handler Handler) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.handlers[eventType] = append(k.handlers[eventType], handler)
}

// Publish sends an event to Kafka
func (k *KafkaEventBus) Publish(ctx context.Context, event Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("event marshal failed: %w", err)
	}

	record := &kgo.Record{
		Topic: k.topic(event.Type),
		Value: data,
	}

	k.client.Produce(ctx, record, func(_ *kgo.Record, err error) {
		if err != nil {
			log.Printf("[kafka] produce error for %s: %v", event.Type, err)
		}
	})

	return nil
}

// StartConsuming begins consuming events from subscribed topics
func (k *KafkaEventBus) StartConsuming(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	k.cancel = cancel

	// Collect all subscribed topics
	k.mu.RLock()
	topics := make([]string, 0, len(k.handlers))
	for eventType := range k.handlers {
		topics = append(topics, k.topic(eventType))
	}
	k.mu.RUnlock()

	if len(topics) == 0 {
		return
	}

	// Add topics to consume
	k.client.AddConsumeTopics(topics...)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				fetches := k.client.PollFetches(ctx)
				if fetches.IsClientClosed() {
					return
				}

				fetches.EachRecord(func(record *kgo.Record) {
					var event Event
					if err := json.Unmarshal(record.Value, &event); err != nil {
						log.Printf("[kafka] unmarshal error: %v", err)
						return
					}

					k.mu.RLock()
					handlers := k.handlers[event.Type]
					k.mu.RUnlock()

					for _, h := range handlers {
						if err := h(ctx, event); err != nil {
							log.Printf("[kafka] handler error for %s: %v", event.Type, err)
						}
					}
				})
			}
		}
	}()
}

// Close stops consuming and closes the client
func (k *KafkaEventBus) Close() {
	if k.cancel != nil {
		k.cancel()
	}
	k.client.Close()
}
