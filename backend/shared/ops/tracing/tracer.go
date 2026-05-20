// Package tracing은 OpenTelemetry 호환 분산 추적을 제공합니다.
//
// 환경변수:
//   - TRACING_ENABLED: "true" 활성화
//   - OTEL_EXPORTER: "otlp", "stdout", "noop" (기본 noop)
//   - OTEL_SERVICE_NAME: 서비스 이름
//   - OTEL_EXPORTER_OTLP_ENDPOINT: OTLP collector 주소
package tracing

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// ============================================================================
// 도메인 모델
// ============================================================================

// SpanKind는 스팬 종류입니다.
const (
	SpanKindClient   = "client"
	SpanKindServer   = "server"
	SpanKindInternal = "internal"
	SpanKindProducer = "producer"
	SpanKindConsumer = "consumer"
)

// Status는 스팬 상태입니다.
const (
	StatusOK    = "ok"
	StatusError = "error"
	StatusUnset = "unset"
)

// Span은 단일 추적 스팬입니다.
type Span struct {
	TraceID    string
	SpanID     string
	ParentID   string
	Name       string
	Kind       string
	StartTime  time.Time
	EndTime    time.Time
	Attributes map[string]string
	Events     []SpanEvent
	Status     string
	StatusMsg  string
	Service    string
}

// SpanEvent는 스팬 내 이벤트입니다.
type SpanEvent struct {
	Name       string
	Timestamp  time.Time
	Attributes map[string]string
}

// Duration은 스팬의 실행 시간을 반환합니다.
func (s *Span) Duration() time.Duration {
	if s.EndTime.IsZero() {
		return 0
	}
	return s.EndTime.Sub(s.StartTime)
}

// SetAttribute는 스팬 속성을 설정합니다.
func (s *Span) SetAttribute(key, value string) {
	if s.Attributes == nil {
		s.Attributes = make(map[string]string)
	}
	s.Attributes[key] = value
}

// AddEvent는 스팬에 이벤트를 추가합니다.
func (s *Span) AddEvent(name string, attrs map[string]string) {
	s.Events = append(s.Events, SpanEvent{
		Name:       name,
		Timestamp:  time.Now().UTC(),
		Attributes: attrs,
	})
}

// SetStatus는 스팬 상태를 설정합니다.
func (s *Span) SetStatus(status, msg string) {
	s.Status = status
	s.StatusMsg = msg
}

// End는 스팬을 종료하고 exporter로 전송합니다.
func (s *Span) End(t Tracer) {
	if s.EndTime.IsZero() {
		s.EndTime = time.Now().UTC()
	}
	if s.Status == "" {
		s.Status = StatusOK
	}
	t.Export(s)
}

// ============================================================================
// Context 통합
// ============================================================================

type spanContextKey struct{}

// ContextWithSpan은 컨텍스트에 스팬을 첨부합니다.
func ContextWithSpan(ctx context.Context, span *Span) context.Context {
	return context.WithValue(ctx, spanContextKey{}, span)
}

// SpanFromContext는 컨텍스트에서 스팬을 추출합니다.
func SpanFromContext(ctx context.Context) *Span {
	if v := ctx.Value(spanContextKey{}); v != nil {
		if s, ok := v.(*Span); ok {
			return s
		}
	}
	return nil
}

// ============================================================================
// Tracer 인터페이스
// ============================================================================

// Tracer는 분산 추적 인터페이스입니다.
type Tracer interface {
	StartSpan(ctx context.Context, name string, kind string) (context.Context, *Span)
	Export(span *Span) error
	Provider() string
	Shutdown(ctx context.Context) error
}

// ============================================================================
// 팩토리
// ============================================================================

// NewFromEnv는 환경변수에 따라 Tracer를 생성합니다.
func NewFromEnv() Tracer {
	enabled := strings.ToLower(os.Getenv("TRACING_ENABLED")) == "true"
	if !enabled {
		return NewNoopTracer()
	}

	exporter := strings.ToLower(os.Getenv("OTEL_EXPORTER"))
	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		serviceName = "manpasik"
	}

	switch exporter {
	case "otlp":
		return NewOTLPTracer(
			serviceName,
			os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		)
	case "stdout":
		return NewStdoutTracer(serviceName)
	default:
		return NewMemoryTracer(serviceName)
	}
}

// ============================================================================
// Noop Tracer (기본)
// ============================================================================

// NoopTracer는 추적하지 않는 Tracer입니다.
type NoopTracer struct{}

// NewNoopTracer는 새 Noop Tracer를 생성합니다.
func NewNoopTracer() *NoopTracer { return &NoopTracer{} }

// StartSpan은 빈 스팬을 반환합니다.
func (t *NoopTracer) StartSpan(ctx context.Context, name, kind string) (context.Context, *Span) {
	return ctx, &Span{Name: name, Kind: kind, StartTime: time.Now().UTC()}
}

// Export는 아무 동작도 하지 않습니다.
func (t *NoopTracer) Export(_ *Span) error { return nil }

// Provider는 이름을 반환합니다.
func (t *NoopTracer) Provider() string { return "noop" }

// Shutdown은 즉시 반환됩니다.
func (t *NoopTracer) Shutdown(_ context.Context) error { return nil }

// ============================================================================
// 인메모리 Tracer (테스트/디버깅용)
// ============================================================================

// MemoryTracer는 인메모리 스팬 저장소입니다.
type MemoryTracer struct {
	mu          sync.Mutex
	serviceName string
	spans       []*Span
}

// NewMemoryTracer는 새 인메모리 Tracer를 생성합니다.
func NewMemoryTracer(serviceName string) *MemoryTracer {
	return &MemoryTracer{serviceName: serviceName}
}

// StartSpan은 새 스팬을 시작합니다.
func (t *MemoryTracer) StartSpan(ctx context.Context, name, kind string) (context.Context, *Span) {
	parentSpan := SpanFromContext(ctx)
	traceID := generateTraceID()
	parentID := ""
	if parentSpan != nil {
		traceID = parentSpan.TraceID
		parentID = parentSpan.SpanID
	}

	span := &Span{
		TraceID:    traceID,
		SpanID:     generateSpanID(),
		ParentID:   parentID,
		Name:       name,
		Kind:       kind,
		Service:    t.serviceName,
		StartTime:  time.Now().UTC(),
		Attributes: make(map[string]string),
	}
	return ContextWithSpan(ctx, span), span
}

// Export는 스팬을 메모리에 저장합니다.
func (t *MemoryTracer) Export(span *Span) error {
	if span == nil {
		return errors.New("span is nil")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.spans = append(t.spans, span)
	return nil
}

// Provider는 이름을 반환합니다.
func (t *MemoryTracer) Provider() string { return "memory" }

// Shutdown은 메모리를 비웁니다.
func (t *MemoryTracer) Shutdown(_ context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.spans = nil
	return nil
}

// Spans는 저장된 스팬 목록을 반환합니다 (테스트용).
func (t *MemoryTracer) Spans() []*Span {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]*Span, len(t.spans))
	copy(out, t.spans)
	return out
}

// SpansByTraceID는 특정 trace의 모든 스팬을 반환합니다.
func (t *MemoryTracer) SpansByTraceID(traceID string) []*Span {
	t.mu.Lock()
	defer t.mu.Unlock()
	var out []*Span
	for _, s := range t.spans {
		if s.TraceID == traceID {
			out = append(out, s)
		}
	}
	return out
}

// Clear는 저장된 스팬을 비웁니다.
func (t *MemoryTracer) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.spans = nil
}

// Count는 저장된 스팬 수를 반환합니다.
func (t *MemoryTracer) Count() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.spans)
}

// ============================================================================
// stdout Tracer (개발 디버깅용)
// ============================================================================

// StdoutTracer는 표준 출력에 스팬을 기록하는 Tracer입니다.
type StdoutTracer struct {
	serviceName string
	memory      *MemoryTracer
}

// NewStdoutTracer는 새 stdout Tracer를 생성합니다.
func NewStdoutTracer(serviceName string) *StdoutTracer {
	return &StdoutTracer{
		serviceName: serviceName,
		memory:      NewMemoryTracer(serviceName),
	}
}

// StartSpan은 스팬을 시작합니다.
func (t *StdoutTracer) StartSpan(ctx context.Context, name, kind string) (context.Context, *Span) {
	return t.memory.StartSpan(ctx, name, kind)
}

// Export는 stdout에 출력합니다.
func (t *StdoutTracer) Export(span *Span) error {
	if span == nil {
		return errors.New("span is nil")
	}
	fmt.Fprintf(
		os.Stdout,
		"[trace] service=%s trace_id=%s span_id=%s name=%s kind=%s duration=%s status=%s\n",
		span.Service, span.TraceID, span.SpanID, span.Name, span.Kind,
		span.Duration(), span.Status,
	)
	return t.memory.Export(span)
}

// Provider는 이름을 반환합니다.
func (t *StdoutTracer) Provider() string { return "stdout" }

// Shutdown은 메모리 tracer를 종료합니다.
func (t *StdoutTracer) Shutdown(ctx context.Context) error {
	return t.memory.Shutdown(ctx)
}

// ============================================================================
// OTLP Tracer (스텁)
// ============================================================================

// OTLPTracer는 OTLP gRPC/HTTP exporter Tracer입니다.
//
// 운영에서는 go.opentelemetry.io/otel SDK 통합 권장.
type OTLPTracer struct {
	serviceName string
	endpoint    string
	memory      *MemoryTracer
}

// NewOTLPTracer는 새 OTLP Tracer를 생성합니다.
func NewOTLPTracer(serviceName, endpoint string) *OTLPTracer {
	return &OTLPTracer{
		serviceName: serviceName,
		endpoint:    endpoint,
		memory:      NewMemoryTracer(serviceName),
	}
}

// StartSpan은 스팬을 시작합니다.
func (t *OTLPTracer) StartSpan(ctx context.Context, name, kind string) (context.Context, *Span) {
	return t.memory.StartSpan(ctx, name, kind)
}

// Export는 OTLP collector로 전송합니다 (실제 구현은 SDK 통합).
func (t *OTLPTracer) Export(span *Span) error {
	return t.memory.Export(span)
}

// Provider는 이름을 반환합니다.
func (t *OTLPTracer) Provider() string { return "otlp" }

// Shutdown은 종료합니다.
func (t *OTLPTracer) Shutdown(ctx context.Context) error {
	return t.memory.Shutdown(ctx)
}

// ============================================================================
// Trace ID 전파 (W3C Trace Context)
// ============================================================================

// TraceContext는 HTTP 헤더 기반 trace 컨텍스트입니다.
type TraceContext struct {
	TraceID  string
	SpanID   string
	Sampled  bool
}

// W3CTraceParentHeader는 W3C Trace Context 헤더 이름입니다.
const W3CTraceParentHeader = "traceparent"

// EncodeW3C는 traceparent 헤더 값을 생성합니다.
//
// 형식: 00-{trace-id-32}-{span-id-16}-{flags-2}
func EncodeW3C(tc *TraceContext) string {
	if tc == nil || tc.TraceID == "" {
		return ""
	}
	flags := "00"
	if tc.Sampled {
		flags = "01"
	}
	return fmt.Sprintf("00-%s-%s-%s", tc.TraceID, tc.SpanID, flags)
}

// DecodeW3C는 traceparent 헤더 값을 파싱합니다.
func DecodeW3C(header string) (*TraceContext, error) {
	parts := strings.Split(header, "-")
	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid traceparent format")
	}
	if parts[0] != "00" {
		return nil, fmt.Errorf("unsupported version: %s", parts[0])
	}
	if len(parts[1]) != 32 || len(parts[2]) != 16 {
		return nil, fmt.Errorf("invalid trace_id or span_id length")
	}
	return &TraceContext{
		TraceID: parts[1],
		SpanID:  parts[2],
		Sampled: parts[3] == "01",
	}, nil
}

// ============================================================================
// 헬퍼: ID 생성
// ============================================================================

// generateTraceID는 32자 헥스 trace ID를 생성합니다.
func generateTraceID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// fallback: 타임스탬프 기반
		now := time.Now().UnixNano()
		for i := range b {
			b[i] = byte((now >> (i * 8)) & 0xFF)
		}
	}
	return hex.EncodeToString(b)
}

// generateSpanID는 16자 헥스 span ID를 생성합니다.
func generateSpanID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		now := time.Now().UnixNano()
		for i := range b {
			b[i] = byte((now >> (i * 8)) & 0xFF)
		}
	}
	return hex.EncodeToString(b)
}
