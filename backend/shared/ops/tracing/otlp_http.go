// Package tracing: OTLP HTTP exporter — 실 OTel collector 송신
//
// W3C Trace Context를 따르며 OTLP/HTTP JSON encoding(application/json)을 사용합니다.
// 외부 SDK 의존성 없이 net/http로 직접 구현.
package tracing

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// ============================================================================
// HTTPDoer (테스트 가능성)
// ============================================================================

// OTLPHTTPDoer는 테스트 가능한 HTTP 클라이언트입니다.
type OTLPHTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// ============================================================================
// OTLPHTTPTracer
// ============================================================================

// OTLPHTTPTracer는 OTLP/HTTP JSON 형식으로 OTel collector에 송신합니다.
//
// 엔드포인트 예: "https://otel-collector.example.com:4318/v1/traces"
type OTLPHTTPTracer struct {
	mu          sync.Mutex
	serviceName string
	endpoint    string
	headers     map[string]string
	httpClient  OTLPHTTPDoer
	memory      *MemoryTracer
	batch       []*Span
	batchSize   int
	flushTimeout time.Duration
}

// OTLPConfig는 exporter 설정입니다.
type OTLPConfig struct {
	ServiceName string
	Endpoint    string
	Headers     map[string]string // 인증 등
	BatchSize   int
	FlushTimeout time.Duration
}

// NewOTLPHTTPTracer는 새 OTLP/HTTP 트레이서를 생성합니다.
func NewOTLPHTTPTracer(cfg OTLPConfig) *OTLPHTTPTracer {
	if cfg.Endpoint == "" {
		cfg.Endpoint = "http://localhost:4318/v1/traces"
	}
	if cfg.BatchSize <= 0 {
		cfg.BatchSize = 100
	}
	if cfg.FlushTimeout <= 0 {
		cfg.FlushTimeout = 5 * time.Second
	}
	if cfg.ServiceName == "" {
		cfg.ServiceName = "manpasik"
	}
	return &OTLPHTTPTracer{
		serviceName:  cfg.ServiceName,
		endpoint:     cfg.Endpoint,
		headers:      cfg.Headers,
		httpClient:   &http.Client{Timeout: 10 * time.Second},
		memory:       NewMemoryTracer(cfg.ServiceName),
		batchSize:    cfg.BatchSize,
		flushTimeout: cfg.FlushTimeout,
	}
}

// SetHTTPClient는 테스트용 HTTP 클라이언트를 설정합니다.
func (t *OTLPHTTPTracer) SetHTTPClient(c OTLPHTTPDoer) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.httpClient = c
}

// Provider는 이름을 반환합니다.
func (t *OTLPHTTPTracer) Provider() string { return "otlp_http" }

// StartSpan은 새 span을 시작합니다 (memory tracer에 위임).
func (t *OTLPHTTPTracer) StartSpan(ctx context.Context, name, kind string) (context.Context, *Span) {
	return t.memory.StartSpan(ctx, name, kind)
}

// Export는 span을 batch에 추가하고, batchSize 도달 시 flush합니다.
func (t *OTLPHTTPTracer) Export(span *Span) error {
	if span == nil {
		return errors.New("span is nil")
	}
	// 메모리 tracer에도 동시 저장 (로컬 디버깅 + 보고서 생성용)
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

// Flush는 누적된 batch를 즉시 송신합니다.
func (t *OTLPHTTPTracer) Flush(ctx context.Context) error {
	t.mu.Lock()
	if len(t.batch) == 0 {
		t.mu.Unlock()
		return nil
	}
	toSend := t.batch
	t.batch = nil
	t.mu.Unlock()

	body := buildOTLPBody(t.serviceName, toSend)
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("OTLP marshal: %w", err)
	}

	if ctx == nil {
		ctx = context.Background()
	}
	flushCtx, cancel := context.WithTimeout(ctx, t.flushTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(flushCtx, "POST", t.endpoint, bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range t.headers {
		req.Header.Set(k, v)
	}

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("OTLP send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("OTLP HTTP %d", resp.StatusCode)
	}
	return nil
}

// Shutdown은 batch를 flush하고 메모리 정리.
func (t *OTLPHTTPTracer) Shutdown(ctx context.Context) error {
	if err := t.Flush(ctx); err != nil {
		return err
	}
	return t.memory.Shutdown(ctx)
}

// PendingCount는 flush 대기 중인 span 수를 반환합니다.
func (t *OTLPHTTPTracer) PendingCount() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.batch)
}

// ============================================================================
// OTLP/HTTP JSON 페이로드 빌더
// ============================================================================

// otlpResource는 OTel resource (서비스 메타데이터)입니다.
type otlpResource struct {
	Attributes []otlpKV `json:"attributes"`
}

// otlpKV는 OTel attribute key-value pair입니다.
type otlpKV struct {
	Key   string         `json:"key"`
	Value otlpAnyValue   `json:"value"`
}

// otlpAnyValue는 OTel AnyValue (다형성 값)입니다.
type otlpAnyValue struct {
	StringValue string `json:"stringValue,omitempty"`
	IntValue    string `json:"intValue,omitempty"`   // OTLP는 int64를 string으로
	BoolValue   bool   `json:"boolValue,omitempty"`
	DoubleValue float64 `json:"doubleValue,omitempty"`
}

// otlpSpan은 OTLP/HTTP 스팬 메시지 형식입니다.
type otlpSpan struct {
	TraceID           string   `json:"traceId"`
	SpanID            string   `json:"spanId"`
	ParentSpanID      string   `json:"parentSpanId,omitempty"`
	Name              string   `json:"name"`
	Kind              int      `json:"kind"`
	StartTimeUnixNano string   `json:"startTimeUnixNano"`
	EndTimeUnixNano   string   `json:"endTimeUnixNano"`
	Attributes        []otlpKV `json:"attributes,omitempty"`
	Events            []otlpEvent `json:"events,omitempty"`
	Status            otlpStatus `json:"status"`
}

type otlpEvent struct {
	TimeUnixNano string   `json:"timeUnixNano"`
	Name         string   `json:"name"`
	Attributes   []otlpKV `json:"attributes,omitempty"`
}

type otlpStatus struct {
	Code    int    `json:"code"` // 0=Unset 1=OK 2=Error
	Message string `json:"message,omitempty"`
}

// otlpBody는 최상위 OTLP/HTTP traces 메시지입니다.
type otlpBody struct {
	ResourceSpans []otlpResourceSpans `json:"resourceSpans"`
}

type otlpResourceSpans struct {
	Resource    otlpResource    `json:"resource"`
	ScopeSpans  []otlpScopeSpan `json:"scopeSpans"`
}

type otlpScopeSpan struct {
	Scope otlpScope  `json:"scope"`
	Spans []otlpSpan `json:"spans"`
}

type otlpScope struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

func buildOTLPBody(serviceName string, spans []*Span) *otlpBody {
	otlpSpans := make([]otlpSpan, 0, len(spans))
	for _, s := range spans {
		otlpSpans = append(otlpSpans, convertSpan(s))
	}

	return &otlpBody{
		ResourceSpans: []otlpResourceSpans{{
			Resource: otlpResource{
				Attributes: []otlpKV{{
					Key: "service.name",
					Value: otlpAnyValue{StringValue: serviceName},
				}},
			},
			ScopeSpans: []otlpScopeSpan{{
				Scope: otlpScope{Name: "manpasik/medical", Version: "1.0"},
				Spans: otlpSpans,
			}},
		}},
	}
}

func convertSpan(s *Span) otlpSpan {
	attrs := make([]otlpKV, 0, len(s.Attributes))
	for k, v := range s.Attributes {
		attrs = append(attrs, otlpKV{Key: k, Value: otlpAnyValue{StringValue: v}})
	}
	if s.Service != "" {
		attrs = append(attrs, otlpKV{Key: "service.instance", Value: otlpAnyValue{StringValue: s.Service}})
	}

	events := make([]otlpEvent, 0, len(s.Events))
	for _, e := range s.Events {
		evAttrs := make([]otlpKV, 0, len(e.Attributes))
		for k, v := range e.Attributes {
			evAttrs = append(evAttrs, otlpKV{Key: k, Value: otlpAnyValue{StringValue: v}})
		}
		events = append(events, otlpEvent{
			TimeUnixNano: fmt.Sprintf("%d", e.Timestamp.UnixNano()),
			Name:         e.Name,
			Attributes:   evAttrs,
		})
	}

	statusCode := 0 // Unset
	switch s.Status {
	case StatusOK:
		statusCode = 1
	case StatusError:
		statusCode = 2
	}

	return otlpSpan{
		TraceID:           s.TraceID,
		SpanID:            s.SpanID,
		ParentSpanID:      s.ParentID,
		Name:              s.Name,
		Kind:              kindToOTLP(s.Kind),
		StartTimeUnixNano: fmt.Sprintf("%d", s.StartTime.UnixNano()),
		EndTimeUnixNano:   fmt.Sprintf("%d", s.EndTime.UnixNano()),
		Attributes:        attrs,
		Events:            events,
		Status: otlpStatus{
			Code:    statusCode,
			Message: s.StatusMsg,
		},
	}
}

// kindToOTLP는 만파식 SpanKind를 OTLP 정수 코드로 변환합니다.
//
// OTLP: 0=Internal 1=Internal 2=Server 3=Client 4=Producer 5=Consumer
func kindToOTLP(kind string) int {
	switch kind {
	case SpanKindServer:
		return 2
	case SpanKindClient:
		return 3
	case SpanKindProducer:
		return 4
	case SpanKindConsumer:
		return 5
	default:
		return 1 // Internal
	}
}

// ============================================================================
// 헬퍼: 인증 헤더 빌더
// ============================================================================

// AuthHeaders는 OTel collector 인증 헤더를 생성합니다.
//
// 지원 패턴:
//   - Bearer 토큰: AuthHeaders("bearer", "token-string")
//   - API Key: AuthHeaders("api-key", "key-string")
func AuthHeaders(authType, value string) map[string]string {
	if authType == "" || value == "" {
		return nil
	}
	switch authType {
	case "bearer":
		return map[string]string{"Authorization": "Bearer " + value}
	case "api-key":
		return map[string]string{"Api-Key": value}
	default:
		return map[string]string{authType: value}
	}
}
