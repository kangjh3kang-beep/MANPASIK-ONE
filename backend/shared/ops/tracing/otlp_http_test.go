package tracing_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/manpasik/backend/shared/ops/tracing"
)

func TestOTLPHTTPTracer_DefaultEndpoint(t *testing.T) {
	tr := tracing.NewOTLPHTTPTracer(tracing.OTLPConfig{})
	if tr.Provider() != "otlp_http" {
		t.Errorf("Provider = %q", tr.Provider())
	}
}

func TestOTLPHTTPTracer_BatchAndFlush(t *testing.T) {
	var receivedSpans int
	var mu sync.Mutex
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var msg map[string]interface{}
		if err := json.Unmarshal(body, &msg); err != nil {
			t.Errorf("invalid JSON: %v", err)
		}
		// resourceSpans[0].scopeSpans[0].spans 카운트
		if rs, ok := msg["resourceSpans"].([]interface{}); ok && len(rs) > 0 {
			if rsmap, ok := rs[0].(map[string]interface{}); ok {
				if ss, ok := rsmap["scopeSpans"].([]interface{}); ok && len(ss) > 0 {
					if ssmap, ok := ss[0].(map[string]interface{}); ok {
						if spans, ok := ssmap["spans"].([]interface{}); ok {
							mu.Lock()
							receivedSpans += len(spans)
							mu.Unlock()
						}
					}
				}
			}
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	tr := tracing.NewOTLPHTTPTracer(tracing.OTLPConfig{
		ServiceName: "test-svc",
		Endpoint:    server.URL,
		BatchSize:   3,
	})

	ctx := context.Background()
	for i := 0; i < 5; i++ {
		_, span := tr.StartSpan(ctx, "op-"+string(rune('a'+i)), tracing.SpanKindInternal)
		span.SetAttribute("idx", string(rune('a'+i)))
		span.End(tr)
	}

	// 5 span → batchSize=3, 자동 1회 flush + 명시적 flush
	if err := tr.Flush(ctx); err != nil {
		t.Fatalf("Flush = %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if receivedSpans != 5 {
		t.Errorf("받은 span = %d, want 5", receivedSpans)
	}
}

func TestOTLPHTTPTracer_AuthHeader(t *testing.T) {
	gotAuth := ""
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	tr := tracing.NewOTLPHTTPTracer(tracing.OTLPConfig{
		Endpoint: server.URL,
		Headers:  tracing.AuthHeaders("bearer", "secret-token"),
	})

	_, span := tr.StartSpan(context.Background(), "x", tracing.SpanKindInternal)
	span.End(tr)
	_ = tr.Flush(context.Background())

	if !strings.Contains(gotAuth, "Bearer secret-token") {
		t.Errorf("Authorization = %q", gotAuth)
	}
}

func TestOTLPHTTPTracer_FlushEmpty(t *testing.T) {
	tr := tracing.NewOTLPHTTPTracer(tracing.OTLPConfig{Endpoint: "http://x"})
	if err := tr.Flush(context.Background()); err != nil {
		t.Errorf("빈 flush = %v", err)
	}
}

func TestOTLPHTTPTracer_FlushOn4xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	tr := tracing.NewOTLPHTTPTracer(tracing.OTLPConfig{Endpoint: server.URL})
	_, span := tr.StartSpan(context.Background(), "x", tracing.SpanKindInternal)
	span.End(tr)
	if err := tr.Flush(context.Background()); err == nil {
		t.Error("4xx에 에러 없음")
	}
}

func TestOTLPHTTPTracer_PendingCount(t *testing.T) {
	// 충분히 큰 batch (자동 flush 안됨)
	tr := tracing.NewOTLPHTTPTracer(tracing.OTLPConfig{
		Endpoint: "http://x", BatchSize: 100,
	})
	for i := 0; i < 5; i++ {
		_, span := tr.StartSpan(context.Background(), "x", tracing.SpanKindInternal)
		span.End(tr)
	}
	if tr.PendingCount() != 5 {
		t.Errorf("PendingCount = %d, want 5", tr.PendingCount())
	}
}

func TestOTLPHTTPTracer_Shutdown(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	tr := tracing.NewOTLPHTTPTracer(tracing.OTLPConfig{Endpoint: server.URL, BatchSize: 100})
	for i := 0; i < 3; i++ {
		_, span := tr.StartSpan(context.Background(), "x", tracing.SpanKindInternal)
		span.End(tr)
	}

	if err := tr.Shutdown(context.Background()); err != nil {
		t.Errorf("Shutdown = %v", err)
	}
	if tr.PendingCount() != 0 {
		t.Errorf("Shutdown 후 PendingCount = %d", tr.PendingCount())
	}
}

func TestAuthHeaders_Bearer(t *testing.T) {
	h := tracing.AuthHeaders("bearer", "tok")
	if h["Authorization"] != "Bearer tok" {
		t.Errorf("Authorization = %q", h["Authorization"])
	}
}

func TestAuthHeaders_APIKey(t *testing.T) {
	h := tracing.AuthHeaders("api-key", "k")
	if h["Api-Key"] != "k" {
		t.Errorf("Api-Key = %q", h["Api-Key"])
	}
}

func TestAuthHeaders_Empty(t *testing.T) {
	if tracing.AuthHeaders("", "x") != nil {
		t.Error("빈 type에 헤더 반환")
	}
	if tracing.AuthHeaders("bearer", "") != nil {
		t.Error("빈 value에 헤더 반환")
	}
}

func TestOTLPHTTPTracer_SpanAttributesEncoded(t *testing.T) {
	var lastBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		lastBody = string(body)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	tr := tracing.NewOTLPHTTPTracer(tracing.OTLPConfig{
		ServiceName: "manpasik-test",
		Endpoint:    server.URL,
	})

	_, span := tr.StartSpan(context.Background(), "encode-test", tracing.SpanKindServer)
	span.SetAttribute("user.id", "u-001")
	span.SetAttribute("session.id", "s-1")
	span.SetStatus(tracing.StatusOK, "")
	span.End(tr)

	_ = tr.Flush(context.Background())

	if !strings.Contains(lastBody, "encode-test") {
		t.Error("span name 미인코딩")
	}
	if !strings.Contains(lastBody, "u-001") {
		t.Error("attribute 미인코딩")
	}
	if !strings.Contains(lastBody, "manpasik-test") {
		t.Error("service.name 미인코딩")
	}
}

func TestOTLPHTTPTracer_FlushTimeout(t *testing.T) {
	// 3초 응답 지연
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	tr := tracing.NewOTLPHTTPTracer(tracing.OTLPConfig{
		Endpoint:     server.URL,
		FlushTimeout: 500 * time.Millisecond,
	})
	_, span := tr.StartSpan(context.Background(), "x", tracing.SpanKindInternal)
	span.End(tr)

	if err := tr.Flush(context.Background()); err == nil {
		t.Error("타임아웃 시 에러 없음")
	}
}

func TestOTLPHTTPTracer_ErrorStatusEncoding(t *testing.T) {
	var receivedBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		receivedBody = string(body)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	tr := tracing.NewOTLPHTTPTracer(tracing.OTLPConfig{Endpoint: server.URL})
	_, span := tr.StartSpan(context.Background(), "fail", tracing.SpanKindInternal)
	span.SetStatus(tracing.StatusError, "simulated failure")
	span.End(tr)
	_ = tr.Flush(context.Background())

	// status code = 2 (Error)
	if !strings.Contains(receivedBody, `"code":2`) {
		t.Errorf("body = %q (Error status 미인코딩)", receivedBody)
	}
}
