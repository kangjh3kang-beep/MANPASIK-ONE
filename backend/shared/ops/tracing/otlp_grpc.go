package tracing

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
)

// SpanExporter 는 송신 백엔드 추상화.
//
// 실 OTel collector 송신은 go.opentelemetry.io/proto/otlp 의 traces v1 메시지를
// gRPC 로 송신해야 하지만, 이 패키지는 외부 SDK 의존을 피하기 위해 인터페이스로
// 분리. 운영 환경에서는 thin adapter 가 OTel SDK 의 grpc.ClientConn 을 감싸고
// 이 인터페이스를 구현하여 주입.
//
// payload 는 buildOTLPBody 결과의 JSON 바이트 (HTTP 와 동일 스키마).
// gRPC 어댑터는 이 JSON 을 proto 로 변환하여 전송.
type SpanExporter interface {
	Export(ctx context.Context, payload []byte) error
	Shutdown(ctx context.Context) error
	Provider() string
}

// OTLPGRPCConfig 는 gRPC 트레이서 설정.
type OTLPGRPCConfig struct {
	ServiceName  string
	BatchSize    int
	FlushTimeout time.Duration
}

// OTLPGRPCTracer 는 SpanExporter 인터페이스를 통해 송신하는 gRPC 트레이서.
//
// 배치 + 플러시 동작은 OTLPHTTPTracer 와 동일. 차이점은 실제 송신을 외부
// SpanExporter 에 위임하므로 다양한 백엔드 (gRPC/Kafka/파일 등) 와 호환.
type OTLPGRPCTracer struct {
	mu           sync.Mutex
	serviceName  string
	exporter     SpanExporter
	memory       *MemoryTracer
	batch        []*Span
	batchSize    int
	flushTimeout time.Duration
}

// NewOTLPGRPCTracer 생성. exporter=nil 이면 NoopExporter 사용.
func NewOTLPGRPCTracer(cfg OTLPGRPCConfig, exporter SpanExporter) *OTLPGRPCTracer {
	if cfg.BatchSize <= 0 {
		cfg.BatchSize = 100
	}
	if cfg.FlushTimeout <= 0 {
		cfg.FlushTimeout = 5 * time.Second
	}
	if cfg.ServiceName == "" {
		cfg.ServiceName = "manpasik"
	}
	if exporter == nil {
		exporter = &NoopSpanExporter{}
	}
	return &OTLPGRPCTracer{
		serviceName:  cfg.ServiceName,
		exporter:     exporter,
		memory:       NewMemoryTracer(cfg.ServiceName),
		batchSize:    cfg.BatchSize,
		flushTimeout: cfg.FlushTimeout,
	}
}

// SetExporter 는 런타임에 exporter 교체 (테스트/페일오버용).
func (t *OTLPGRPCTracer) SetExporter(e SpanExporter) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if e != nil {
		t.exporter = e
	}
}

// Provider 이름.
func (t *OTLPGRPCTracer) Provider() string { return "otlp_grpc" }

// StartSpan 은 메모리 tracer 에 위임.
func (t *OTLPGRPCTracer) StartSpan(ctx context.Context, name, kind string) (context.Context, *Span) {
	return t.memory.StartSpan(ctx, name, kind)
}

// Export 는 batch 에 추가하고 batchSize 도달 시 flush.
func (t *OTLPGRPCTracer) Export(span *Span) error {
	if span == nil {
		return errors.New("span is nil")
	}
	_ = t.memory.Export(span)

	t.mu.Lock()
	t.batch = append(t.batch, span)
	shouldFlush := len(t.batch) >= t.batchSize
	t.mu.Unlock()

	if shouldFlush {
		return t.Flush(context.Background())
	}
	return nil
}

// Flush 는 누적된 batch 를 SpanExporter 로 송신.
func (t *OTLPGRPCTracer) Flush(ctx context.Context) error {
	t.mu.Lock()
	if len(t.batch) == 0 {
		t.mu.Unlock()
		return nil
	}
	toSend := t.batch
	t.batch = nil
	exporter := t.exporter
	t.mu.Unlock()

	if exporter == nil {
		return errors.New("exporter 미설정")
	}

	body := buildOTLPBody(t.serviceName, toSend)
	payload, err := marshalOTLPBody(body)
	if err != nil {
		return fmt.Errorf("OTLP marshal: %w", err)
	}

	if ctx == nil {
		ctx = context.Background()
	}
	flushCtx, cancel := context.WithTimeout(ctx, t.flushTimeout)
	defer cancel()

	if err := exporter.Export(flushCtx, payload); err != nil {
		return fmt.Errorf("OTLP gRPC export: %w", err)
	}
	return nil
}

// Shutdown 은 flush + exporter 종료.
func (t *OTLPGRPCTracer) Shutdown(ctx context.Context) error {
	if err := t.Flush(ctx); err != nil {
		return err
	}
	t.mu.Lock()
	exporter := t.exporter
	t.mu.Unlock()
	if exporter != nil {
		if err := exporter.Shutdown(ctx); err != nil {
			return err
		}
	}
	return t.memory.Shutdown(ctx)
}

// PendingCount 는 flush 대기 중인 span 수.
func (t *OTLPGRPCTracer) PendingCount() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.batch)
}

// marshalOTLPBody 는 buildOTLPBody 결과를 JSON 바이트로 직렬화.
//
// 별도 함수로 분리하여 테스트 용이성 확보. gRPC 어댑터는 이 JSON 을
// otelproto 메시지로 변환하여 송신할 수도 있고, JSON-over-gRPC 같은
// 비표준 전송에 그대로 사용할 수도 있다.
func marshalOTLPBody(body *otlpBody) ([]byte, error) {
	return json.Marshal(body)
}

// ============================================================================
// Noop / Memory exporters
// ============================================================================

// NoopSpanExporter 는 모든 호출이 즉시 성공하는 더미 exporter.
type NoopSpanExporter struct {
	mu        sync.Mutex
	callCount int
	bytesSent int64
}

func (n *NoopSpanExporter) Export(_ context.Context, payload []byte) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.callCount++
	n.bytesSent += int64(len(payload))
	return nil
}

func (n *NoopSpanExporter) Shutdown(_ context.Context) error { return nil }
func (n *NoopSpanExporter) Provider() string                 { return "noop" }

// CallCount 는 Export 호출 횟수 반환.
func (n *NoopSpanExporter) CallCount() int {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.callCount
}

// BytesSent 는 누적 송신 바이트 반환.
func (n *NoopSpanExporter) BytesSent() int64 {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.bytesSent
}

// MemorySpanExporter 는 송신된 페이로드를 메모리에 저장 (테스트용).
type MemorySpanExporter struct {
	mu       sync.Mutex
	payloads [][]byte
	failWith error
	closed   bool
}

func (m *MemorySpanExporter) Export(_ context.Context, payload []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return errors.New("exporter closed")
	}
	if m.failWith != nil {
		return m.failWith
	}
	cp := make([]byte, len(payload))
	copy(cp, payload)
	m.payloads = append(m.payloads, cp)
	return nil
}

func (m *MemorySpanExporter) Shutdown(_ context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closed = true
	return nil
}

func (m *MemorySpanExporter) Provider() string { return "memory" }

// SetFailWith 는 다음 Export 호출이 실패하도록 설정 (테스트용).
func (m *MemorySpanExporter) SetFailWith(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.failWith = err
}

// Payloads 는 지금까지 송신된 페이로드 사본 반환.
func (m *MemorySpanExporter) Payloads() [][]byte {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([][]byte, len(m.payloads))
	copy(out, m.payloads)
	return out
}
