package tenancy_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/manpasik/backend/shared/tenancy"
)

func TestWebhookDispatcher_GenericMode(t *testing.T) {
	var (
		gotBody []byte
		mu      sync.Mutex
	)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		gotBody, _ = readBody(r)
		w.WriteHeader(200)
	}))
	defer srv.Close()

	d, err := tenancy.NewWebhookDispatcher(tenancy.WebhookConfig{
		URL: srv.URL, Mode: "generic",
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := d.Dispatch(context.Background(), tenancy.Event{
		Type:     tenancy.EventInvitationCreated,
		TenantID: "hospA",
		ActorID:  "doc1",
	}); err != nil {
		t.Fatal(err)
	}

	mu.Lock()
	defer mu.Unlock()
	var ev tenancy.Event
	if err := json.Unmarshal(gotBody, &ev); err != nil {
		t.Fatal(err)
	}
	if ev.Type != tenancy.EventInvitationCreated || ev.TenantID != "hospA" {
		t.Errorf("ev = %+v", ev)
	}
}

func TestWebhookDispatcher_SlackMode(t *testing.T) {
	var gotBody []byte
	var mu sync.Mutex
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		gotBody, _ = readBody(r)
		w.WriteHeader(200)
	}))
	defer srv.Close()

	d, _ := tenancy.NewWebhookDispatcher(tenancy.WebhookConfig{
		URL: srv.URL, Mode: "slack",
	}, nil)
	_ = d.Dispatch(context.Background(), tenancy.Event{
		Type:     tenancy.EventMembershipCreated,
		TenantID: "fam-A",
		UserID:   "child",
	})

	mu.Lock()
	defer mu.Unlock()
	if !strings.Contains(string(gotBody), "Manpasik Tenancy") {
		t.Errorf("Slack format 누락:\n%s", gotBody)
	}
	if !strings.Contains(string(gotBody), "fam-A") {
		t.Errorf("tenant 누락")
	}
	// attachments 색상 검증
	if !strings.Contains(string(gotBody), "good") {
		t.Errorf("color 누락 (membership.created → good 기대)")
	}
}

func TestWebhookDispatcher_NoURL(t *testing.T) {
	if _, err := tenancy.NewWebhookDispatcher(tenancy.WebhookConfig{}, nil); err == nil {
		t.Error("URL 없는데 통과")
	}
}

func TestWebhookDispatcher_5xxRetry(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.WriteHeader(500)
	}))
	defer srv.Close()

	d, _ := tenancy.NewWebhookDispatcher(tenancy.WebhookConfig{
		URL: srv.URL, MaxRetries: 2, RetryBaseDelay: 10 * time.Millisecond,
	}, nil)
	_ = d.Dispatch(context.Background(), tenancy.Event{Type: "x", TenantID: "t"})
	if atomic.LoadInt32(&attempts) != 1 {
		t.Errorf("첫 시도 = %d", attempts)
	}
	// 재시도 큐에 들어감
	if d.PendingRetryCount() == 0 {
		t.Error("재시도 큐 비어있음")
	}
}

func TestWebhookDispatcher_Async(t *testing.T) {
	var received int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(&received, 1)
		w.WriteHeader(200)
	}))
	defer srv.Close()

	d, _ := tenancy.NewWebhookDispatcher(tenancy.WebhookConfig{URL: srv.URL}, nil)
	d.DispatchAsync(tenancy.Event{Type: "x", TenantID: "t"})
	d.DispatchAsync(tenancy.Event{Type: "y", TenantID: "t"})

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if atomic.LoadInt32(&received) >= 2 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if atomic.LoadInt32(&received) < 2 {
		t.Errorf("async 송신 = %d", received)
	}
}

func TestWebhookDispatcher_OnError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(500)
	}))
	defer srv.Close()

	var errCount int32
	d, _ := tenancy.NewWebhookDispatcher(tenancy.WebhookConfig{
		URL: srv.URL, MaxRetries: 0,
		OnError: func(_ tenancy.Event, _ error) {
			atomic.AddInt32(&errCount, 1)
		},
	}, nil)
	_ = d.Dispatch(context.Background(), tenancy.Event{Type: "x", TenantID: "t"})
	if atomic.LoadInt32(&errCount) == 0 {
		t.Error("OnError 미호출")
	}
}

func TestWebhookDispatcher_StartStop(t *testing.T) {
	d, _ := tenancy.NewWebhookDispatcher(tenancy.WebhookConfig{
		URL: "http://x", RetryBaseDelay: 10 * time.Millisecond,
	}, nil)
	d.Start(context.Background())
	d.Start(context.Background()) // idempotent
	d.Stop()
}

func TestEvent_DefaultTimestamp(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ev tenancy.Event
		body, _ := readBody(r)
		_ = json.Unmarshal(body, &ev)
		if ev.Timestamp.IsZero() {
			t.Error("Dispatch 가 timestamp 자동 설정 안 함")
		}
		w.WriteHeader(200)
	}))
	defer srv.Close()

	d, _ := tenancy.NewWebhookDispatcher(tenancy.WebhookConfig{URL: srv.URL}, nil)
	_ = d.Dispatch(context.Background(), tenancy.Event{Type: "x", TenantID: "t"})
}

// failingDoer 는 모든 요청을 실패시킴.
type failingDoer struct{ err error }

func (f failingDoer) Do(_ *http.Request) (*http.Response, error) {
	return nil, f.err
}

func TestWebhookDispatcher_DoerError(t *testing.T) {
	d, _ := tenancy.NewWebhookDispatcher(tenancy.WebhookConfig{
		URL: "http://x", MaxRetries: 0,
	}, failingDoer{err: errors.New("connection refused")})
	if err := d.Dispatch(context.Background(), tenancy.Event{Type: "x", TenantID: "t"}); err == nil {
		t.Error("doer 에러 통과")
	}
}

func TestWebhookDispatcher_Stats_Sent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(200)
	}))
	defer srv.Close()
	d, _ := tenancy.NewWebhookDispatcher(tenancy.WebhookConfig{URL: srv.URL}, nil)

	for i := 0; i < 3; i++ {
		_ = d.Dispatch(context.Background(), tenancy.Event{Type: "x", TenantID: "t"})
	}
	stats := d.Stats()
	if stats.SentCount != 3 {
		t.Errorf("SentCount = %d, want 3", stats.SentCount)
	}
	if stats.FailedCount != 0 {
		t.Errorf("FailedCount = %d", stats.FailedCount)
	}
}

func TestWebhookDispatcher_Stats_Failed(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(500)
	}))
	defer srv.Close()
	d, _ := tenancy.NewWebhookDispatcher(tenancy.WebhookConfig{
		URL: srv.URL, MaxRetries: 0,
	}, nil)
	_ = d.Dispatch(context.Background(), tenancy.Event{Type: "x", TenantID: "t"})

	stats := d.Stats()
	if stats.FailedCount == 0 {
		t.Errorf("FailedCount = %d", stats.FailedCount)
	}
}

func TestWebhookDispatcher_CollectMetrics_HealthyByDefault(t *testing.T) {
	d, _ := tenancy.NewWebhookDispatcher(tenancy.WebhookConfig{URL: "http://x"}, nil)
	m := d.CollectMetrics()
	if m["status"] != "healthy" {
		t.Errorf("초기 status = %v", m["status"])
	}
	if m["sent_count"].(int64) != 0 {
		t.Errorf("sent_count = %v", m["sent_count"])
	}
}

func TestWriteWebhookPrometheusMetrics_NilDispatcher(t *testing.T) {
	w := httptest.NewRecorder()
	tenancy.WriteWebhookPrometheusMetrics(w, nil)
	body := w.Body.String()
	if !strings.Contains(body, "manpasik_tenancy_webhook_status 0") {
		t.Errorf("nil dispatcher 응답 누락:\n%s", body)
	}
}

func TestWriteWebhookPrometheusMetrics_HealthyDefault(t *testing.T) {
	d, _ := tenancy.NewWebhookDispatcher(tenancy.WebhookConfig{URL: "http://x"}, nil)
	w := httptest.NewRecorder()
	tenancy.WriteWebhookPrometheusMetrics(w, d)
	body := w.Body.String()
	for _, want := range []string{
		"manpasik_tenancy_webhook_sent_total 0",
		"manpasik_tenancy_webhook_failed_total 0",
		"manpasik_tenancy_webhook_status 1",
		"# HELP",
		"# TYPE",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("응답에 누락: %q\n%s", want, body)
		}
	}
}

func TestWriteWebhookPrometheusMetrics_AfterDispatch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(200)
	}))
	defer srv.Close()
	d, _ := tenancy.NewWebhookDispatcher(tenancy.WebhookConfig{URL: srv.URL}, nil)
	for i := 0; i < 5; i++ {
		_ = d.Dispatch(context.Background(), tenancy.Event{Type: "x", TenantID: "t"})
	}
	w := httptest.NewRecorder()
	tenancy.WriteWebhookPrometheusMetrics(w, d)
	body := w.Body.String()
	if !strings.Contains(body, "manpasik_tenancy_webhook_sent_total 5") {
		t.Errorf("sent_total 5 누락:\n%s", body)
	}
}

func TestWebhookDispatcher_CollectMetrics_DegradedOnDropped(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(500)
	}))
	defer srv.Close()
	d, _ := tenancy.NewWebhookDispatcher(tenancy.WebhookConfig{
		URL: srv.URL, MaxRetries: 0, RetryBaseDelay: 5 * time.Millisecond,
	}, nil)
	d.Start(context.Background())
	defer d.Stop()

	_ = d.Dispatch(context.Background(), tenancy.Event{Type: "x", TenantID: "t"})
	// retry 처리 대기 (MaxRetries=0 이므로 즉시 dropped)
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if d.Stats().DroppedCount > 0 {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}

	m := d.CollectMetrics()
	if m["status"] != "degraded" {
		t.Errorf("dropped > 0 인데 status = %v", m["status"])
	}
}

func readBody(r *http.Request) ([]byte, error) {
	defer r.Body.Close()
	const maxRead = 1 << 20
	buf := make([]byte, 0, 512)
	tmp := make([]byte, 512)
	for len(buf) < maxRead {
		n, err := r.Body.Read(tmp)
		if n > 0 {
			buf = append(buf, tmp[:n]...)
		}
		if err != nil {
			return buf, nil
		}
	}
	return buf, nil
}
