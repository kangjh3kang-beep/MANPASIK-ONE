package tracing_test

import (
	"context"
	"testing"
	"time"

	"github.com/manpasik/backend/shared/ops/tracing"
)

func TestNoopTracer_DoesNothing(t *testing.T) {
	tracer := tracing.NewNoopTracer()
	ctx := context.Background()

	_, span := tracer.StartSpan(ctx, "test", tracing.SpanKindInternal)
	span.End(tracer)

	if tracer.Provider() != "noop" {
		t.Errorf("Provider = %q", tracer.Provider())
	}
}

func TestMemoryTracer_StartAndEnd(t *testing.T) {
	tracer := tracing.NewMemoryTracer("test-service")
	ctx := context.Background()

	_, span := tracer.StartSpan(ctx, "operation", tracing.SpanKindServer)
	span.SetAttribute("user_id", "u1")
	span.SetStatus(tracing.StatusOK, "")
	span.End(tracer)

	if tracer.Count() != 1 {
		t.Errorf("Count = %d, want 1", tracer.Count())
	}

	spans := tracer.Spans()
	if spans[0].Name != "operation" {
		t.Errorf("Name = %q", spans[0].Name)
	}
	if spans[0].Service != "test-service" {
		t.Errorf("Service = %q", spans[0].Service)
	}
	if spans[0].Attributes["user_id"] != "u1" {
		t.Errorf("Attribute 미저장")
	}
	if spans[0].Duration() <= 0 {
		t.Error("Duration이 0 이하")
	}
}

func TestMemoryTracer_ParentChildHierarchy(t *testing.T) {
	tracer := tracing.NewMemoryTracer("test")
	ctx := context.Background()

	parentCtx, parent := tracer.StartSpan(ctx, "parent", tracing.SpanKindServer)
	_, child := tracer.StartSpan(parentCtx, "child", tracing.SpanKindClient)

	if child.TraceID != parent.TraceID {
		t.Error("Trace ID가 부모와 다름")
	}
	if child.ParentID != parent.SpanID {
		t.Error("ParentID가 부모 SpanID와 다름")
	}

	child.End(tracer)
	parent.End(tracer)

	traceSpans := tracer.SpansByTraceID(parent.TraceID)
	if len(traceSpans) != 2 {
		t.Errorf("Trace 스팬 = %d, want 2", len(traceSpans))
	}
}

func TestMemoryTracer_Events(t *testing.T) {
	tracer := tracing.NewMemoryTracer("test")
	_, span := tracer.StartSpan(context.Background(), "x", tracing.SpanKindInternal)

	span.AddEvent("checkpoint", map[string]string{"step": "1"})
	span.AddEvent("checkpoint", map[string]string{"step": "2"})
	span.End(tracer)

	if len(span.Events) != 2 {
		t.Errorf("Events = %d, want 2", len(span.Events))
	}
}

func TestMemoryTracer_ErrorStatus(t *testing.T) {
	tracer := tracing.NewMemoryTracer("test")
	_, span := tracer.StartSpan(context.Background(), "failing", tracing.SpanKindInternal)
	span.SetStatus(tracing.StatusError, "database connection failed")
	span.End(tracer)

	spans := tracer.Spans()
	if spans[0].Status != tracing.StatusError {
		t.Errorf("Status = %q", spans[0].Status)
	}
	if spans[0].StatusMsg == "" {
		t.Error("StatusMsg 미설정")
	}
}

func TestMemoryTracer_DefaultStatusOK(t *testing.T) {
	tracer := tracing.NewMemoryTracer("t")
	_, span := tracer.StartSpan(context.Background(), "x", "internal")
	span.End(tracer)

	if span.Status != tracing.StatusOK {
		t.Errorf("기본 Status = %q, want ok", span.Status)
	}
}

func TestMemoryTracer_Clear(t *testing.T) {
	tracer := tracing.NewMemoryTracer("t")
	for i := 0; i < 5; i++ {
		_, span := tracer.StartSpan(context.Background(), "x", "internal")
		span.End(tracer)
	}
	if tracer.Count() != 5 {
		t.Errorf("Count = %d", tracer.Count())
	}
	tracer.Clear()
	if tracer.Count() != 0 {
		t.Errorf("Clear 후 Count = %d", tracer.Count())
	}
}

func TestContextWithSpan_Roundtrip(t *testing.T) {
	tracer := tracing.NewMemoryTracer("t")
	ctx, span := tracer.StartSpan(context.Background(), "x", "internal")

	got := tracing.SpanFromContext(ctx)
	if got != span {
		t.Error("Context에서 Span 복원 실패")
	}
}

func TestSpanFromContext_Empty(t *testing.T) {
	ctx := context.Background()
	if tracing.SpanFromContext(ctx) != nil {
		t.Error("빈 컨텍스트에서 Span이 추출됨")
	}
}

func TestEncodeW3C(t *testing.T) {
	tc := &tracing.TraceContext{
		TraceID: "0123456789abcdef0123456789abcdef",
		SpanID:  "0123456789abcdef",
		Sampled: true,
	}
	encoded := tracing.EncodeW3C(tc)
	expected := "00-0123456789abcdef0123456789abcdef-0123456789abcdef-01"
	if encoded != expected {
		t.Errorf("encoded = %q, want %q", encoded, expected)
	}
}

func TestEncodeW3C_NotSampled(t *testing.T) {
	tc := &tracing.TraceContext{
		TraceID: "abc123abc123abc123abc123abc12345",
		SpanID:  "abc123abc123abcd",
		Sampled: false,
	}
	encoded := tracing.EncodeW3C(tc)
	if encoded[len(encoded)-2:] != "00" {
		t.Errorf("flags = %q, want 00", encoded[len(encoded)-2:])
	}
}

func TestDecodeW3C_Valid(t *testing.T) {
	header := "00-0123456789abcdef0123456789abcdef-0123456789abcdef-01"
	tc, err := tracing.DecodeW3C(header)
	if err != nil {
		t.Fatalf("DecodeW3C 실패: %v", err)
	}
	if tc.TraceID != "0123456789abcdef0123456789abcdef" {
		t.Errorf("TraceID = %q", tc.TraceID)
	}
	if !tc.Sampled {
		t.Error("Sampled = false")
	}
}

func TestDecodeW3C_Invalid(t *testing.T) {
	cases := []string{
		"invalid",
		"00-tooshort-0123456789abcdef-01",
		"01-0123456789abcdef0123456789abcdef-0123456789abcdef-01", // wrong version
	}
	for _, c := range cases {
		if _, err := tracing.DecodeW3C(c); err == nil {
			t.Errorf("invalid header %q가 통과됨", c)
		}
	}
}

func TestStdoutTracer(t *testing.T) {
	tracer := tracing.NewStdoutTracer("test")
	_, span := tracer.StartSpan(context.Background(), "x", "internal")
	span.End(tracer)

	if tracer.Provider() != "stdout" {
		t.Errorf("Provider = %q", tracer.Provider())
	}
	_ = tracer.Shutdown(context.Background())
}

func TestNewFromEnv_Disabled(t *testing.T) {
	t.Setenv("TRACING_ENABLED", "")
	tracer := tracing.NewFromEnv()
	if tracer.Provider() != "noop" {
		t.Errorf("Provider = %q, want noop", tracer.Provider())
	}
}

func TestNewFromEnv_OTLP(t *testing.T) {
	t.Setenv("TRACING_ENABLED", "true")
	t.Setenv("OTEL_EXPORTER", "otlp")
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "https://otel.example.com")

	tracer := tracing.NewFromEnv()
	if tracer.Provider() != "otlp" {
		t.Errorf("Provider = %q", tracer.Provider())
	}
}

func TestNewFromEnv_Stdout(t *testing.T) {
	t.Setenv("TRACING_ENABLED", "true")
	t.Setenv("OTEL_EXPORTER", "stdout")

	tracer := tracing.NewFromEnv()
	if tracer.Provider() != "stdout" {
		t.Errorf("Provider = %q", tracer.Provider())
	}
}

func TestSpan_Duration_NotEnded(t *testing.T) {
	span := &tracing.Span{StartTime: time.Now()}
	if span.Duration() != 0 {
		t.Error("미종료 스팬의 Duration이 0이 아님")
	}
}

func TestSpan_SetAttributeMultiple(t *testing.T) {
	span := &tracing.Span{}
	span.SetAttribute("k1", "v1")
	span.SetAttribute("k2", "v2")
	span.SetAttribute("k1", "v1-updated")

	if span.Attributes["k1"] != "v1-updated" {
		t.Errorf("k1 = %q", span.Attributes["k1"])
	}
	if span.Attributes["k2"] != "v2" {
		t.Errorf("k2 = %q", span.Attributes["k2"])
	}
}
