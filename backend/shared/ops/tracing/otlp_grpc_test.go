package tracing_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/manpasik/backend/shared/ops/tracing"
)

func TestOTLPGRPCTracer_Provider(t *testing.T) {
	tr := tracing.NewOTLPGRPCTracer(tracing.OTLPGRPCConfig{}, nil)
	if tr.Provider() != "otlp_grpc" {
		t.Errorf("Provider = %q", tr.Provider())
	}
}

func TestOTLPGRPCTracer_DefaultExporter(t *testing.T) {
	tr := tracing.NewOTLPGRPCTracer(tracing.OTLPGRPCConfig{}, nil)
	_, span := tr.StartSpan(context.Background(), "x", tracing.SpanKindInternal)
	span.End(tr)
	if err := tr.Flush(context.Background()); err != nil {
		t.Errorf("Flush err = %v", err)
	}
}

func TestOTLPGRPCTracer_ExporterReceivesPayload(t *testing.T) {
	exporter := &tracing.MemorySpanExporter{}
	tr := tracing.NewOTLPGRPCTracer(tracing.OTLPGRPCConfig{ServiceName: "svc-a"}, exporter)

	for i := 0; i < 3; i++ {
		_, span := tr.StartSpan(context.Background(), "op", tracing.SpanKindInternal)
		span.SetAttribute("idx", string(rune('a'+i)))
		span.End(tr)
	}
	if err := tr.Flush(context.Background()); err != nil {
		t.Fatal(err)
	}

	payloads := exporter.Payloads()
	if len(payloads) != 1 {
		t.Fatalf("payloads = %d, want 1", len(payloads))
	}

	// 페이로드는 OTLP/HTTP JSON 과 동일 스키마
	var body map[string]interface{}
	if err := json.Unmarshal(payloads[0], &body); err != nil {
		t.Fatal(err)
	}
	if rs, ok := body["resourceSpans"].([]interface{}); !ok || len(rs) == 0 {
		t.Errorf("resourceSpans 누락: %v", body)
	}
}

func TestOTLPGRPCTracer_BatchSizeAutoFlush(t *testing.T) {
	exporter := &tracing.MemorySpanExporter{}
	tr := tracing.NewOTLPGRPCTracer(tracing.OTLPGRPCConfig{BatchSize: 2}, exporter)

	for i := 0; i < 5; i++ {
		_, span := tr.StartSpan(context.Background(), "x", tracing.SpanKindInternal)
		span.End(tr)
	}
	_ = tr.Flush(context.Background())

	// 5 spans / batchSize=2 → 자동 flush 2회 + 명시 flush 1회 = 최대 3 호출 (꼬리 1개 포함)
	if len(exporter.Payloads()) < 2 {
		t.Errorf("payloads = %d, want >= 2", len(exporter.Payloads()))
	}
}

func TestOTLPGRPCTracer_ExporterFailure(t *testing.T) {
	exporter := &tracing.MemorySpanExporter{}
	exporter.SetFailWith(errors.New("network down"))
	tr := tracing.NewOTLPGRPCTracer(tracing.OTLPGRPCConfig{}, exporter)
	_, span := tr.StartSpan(context.Background(), "x", tracing.SpanKindInternal)
	span.End(tr)
	if err := tr.Flush(context.Background()); err == nil {
		t.Error("실패에 에러 없음")
	}
}

func TestOTLPGRPCTracer_FlushEmpty(t *testing.T) {
	tr := tracing.NewOTLPGRPCTracer(tracing.OTLPGRPCConfig{}, &tracing.MemorySpanExporter{})
	if err := tr.Flush(context.Background()); err != nil {
		t.Errorf("빈 flush err = %v", err)
	}
}

func TestOTLPGRPCTracer_PendingCount(t *testing.T) {
	tr := tracing.NewOTLPGRPCTracer(tracing.OTLPGRPCConfig{BatchSize: 100}, &tracing.MemorySpanExporter{})
	for i := 0; i < 5; i++ {
		_, span := tr.StartSpan(context.Background(), "x", tracing.SpanKindInternal)
		span.End(tr)
	}
	if tr.PendingCount() != 5 {
		t.Errorf("PendingCount = %d", tr.PendingCount())
	}
}

func TestOTLPGRPCTracer_Shutdown(t *testing.T) {
	exporter := &tracing.MemorySpanExporter{}
	tr := tracing.NewOTLPGRPCTracer(tracing.OTLPGRPCConfig{BatchSize: 100}, exporter)
	for i := 0; i < 3; i++ {
		_, span := tr.StartSpan(context.Background(), "x", tracing.SpanKindInternal)
		span.End(tr)
	}
	if err := tr.Shutdown(context.Background()); err != nil {
		t.Errorf("Shutdown err = %v", err)
	}
	if tr.PendingCount() != 0 {
		t.Errorf("Shutdown 후 PendingCount = %d", tr.PendingCount())
	}
	// shutdown 후 추가 export 시도 시 에러
	if err := exporter.Export(context.Background(), []byte("after shutdown")); err == nil {
		t.Error("Shutdown 후 Export 가 성공함")
	}
}

func TestOTLPGRPCTracer_SetExporter(t *testing.T) {
	first := &tracing.MemorySpanExporter{}
	second := &tracing.MemorySpanExporter{}
	tr := tracing.NewOTLPGRPCTracer(tracing.OTLPGRPCConfig{}, first)

	tr.SetExporter(second)
	_, span := tr.StartSpan(context.Background(), "x", tracing.SpanKindInternal)
	span.End(tr)
	_ = tr.Flush(context.Background())

	if len(first.Payloads()) != 0 {
		t.Error("교체 후에도 first exporter 가 호출됨")
	}
	if len(second.Payloads()) != 1 {
		t.Errorf("second exporter payloads = %d", len(second.Payloads()))
	}
}

func TestNoopSpanExporter_CountAndBytes(t *testing.T) {
	exp := &tracing.NoopSpanExporter{}
	_ = exp.Export(context.Background(), []byte("hello"))
	_ = exp.Export(context.Background(), []byte("world!"))
	if exp.CallCount() != 2 {
		t.Errorf("CallCount = %d", exp.CallCount())
	}
	if exp.BytesSent() != 11 {
		t.Errorf("BytesSent = %d", exp.BytesSent())
	}
	if exp.Provider() != "noop" {
		t.Errorf("Provider = %q", exp.Provider())
	}
}

func TestOTLPGRPCTracer_FlushTimeout(t *testing.T) {
	// slow exporter
	slow := &slowExporter{delay: 200 * time.Millisecond}
	tr := tracing.NewOTLPGRPCTracer(tracing.OTLPGRPCConfig{
		FlushTimeout: 50 * time.Millisecond,
	}, slow)
	_, span := tr.StartSpan(context.Background(), "x", tracing.SpanKindInternal)
	span.End(tr)
	if err := tr.Flush(context.Background()); err == nil {
		t.Error("타임아웃에 에러 없음")
	}
}

type slowExporter struct {
	delay time.Duration
}

func (s *slowExporter) Export(ctx context.Context, _ []byte) error {
	select {
	case <-time.After(s.delay):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *slowExporter) Shutdown(_ context.Context) error { return nil }
func (s *slowExporter) Provider() string                 { return "slow" }
