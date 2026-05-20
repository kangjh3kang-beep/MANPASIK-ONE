package tenancy_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/manpasik/backend/shared/tenancy"
)

// flakeyDoer 는 `failTimes` 회 실패 후 성공.
type flakeyDoer struct {
	failTimes int32
	calls     int32
}

func (f *flakeyDoer) Do(_ *http.Request) (*http.Response, error) {
	calls := atomic.AddInt32(&f.calls, 1)
	if calls <= atomic.LoadInt32(&f.failTimes) {
		return nil, errors.New("flakey")
	}
	return &http.Response{StatusCode: 200, Body: http.NoBody}, nil
}

func TestWebhookDispatcher_DLQ_OnMaxRetriesExceeded(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(500)
	}))
	defer srv.Close()

	d, err := tenancy.NewWebhookDispatcher(tenancy.WebhookConfig{
		URL: srv.URL, MaxRetries: 1, RetryBaseDelay: 5 * time.Millisecond,
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	d.Start(context.Background())
	defer d.Stop()

	_ = d.Dispatch(context.Background(), tenancy.Event{
		Type: tenancy.EventInvitationCreated, TenantID: "hospA",
	})

	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if d.DLQCount() > 0 {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if d.DLQCount() == 0 {
		t.Fatal("DLQ 비어있음 — max retries 후 자동 보관 실패")
	}
	entries := d.DLQEntries()
	if entries[0].Event.TenantID != "hospA" {
		t.Errorf("DLQ event tenant = %s", entries[0].Event.TenantID)
	}
	if entries[0].LastError == "" {
		t.Errorf("LastError 누락")
	}
	if entries[0].Attempts < 2 {
		t.Errorf("Attempts = %d, MaxRetries=1 시 최소 2 기대", entries[0].Attempts)
	}
	stats := d.Stats()
	if stats.DLQSize == 0 {
		t.Errorf("Stats.DLQSize = 0")
	}
}

func TestWebhookDispatcher_DLQ_ReplaySuccess(t *testing.T) {
	flakey := &flakeyDoer{failTimes: 999} // 항상 실패하도록 시작
	d, _ := tenancy.NewWebhookDispatcher(tenancy.WebhookConfig{
		URL: "http://x", MaxRetries: 0, RetryBaseDelay: 5 * time.Millisecond,
	}, flakey)
	d.Start(context.Background())
	defer d.Stop()

	_ = d.Dispatch(context.Background(), tenancy.Event{Type: "x", TenantID: "t"})

	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if d.DLQCount() > 0 {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if d.DLQCount() == 0 {
		t.Fatal("DLQ 누적 안 됨")
	}

	// 이제 doer 가 성공하도록 변경
	atomic.StoreInt32(&flakey.failTimes, 0)
	id := d.DLQEntries()[0].ID

	if err := d.DLQReplay(context.Background(), id); err != nil {
		t.Fatalf("Replay 실패: %v", err)
	}
	if d.DLQCount() != 0 {
		t.Errorf("Replay 성공 후 DLQ 미제거: %d", d.DLQCount())
	}
}

func TestWebhookDispatcher_DLQ_ReplayFailureUpdatesEntry(t *testing.T) {
	flakey := &flakeyDoer{failTimes: 999}
	d, _ := tenancy.NewWebhookDispatcher(tenancy.WebhookConfig{
		URL: "http://x", MaxRetries: 0, RetryBaseDelay: 5 * time.Millisecond,
	}, flakey)
	d.Start(context.Background())
	defer d.Stop()

	_ = d.Dispatch(context.Background(), tenancy.Event{Type: "x", TenantID: "t"})

	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if d.DLQCount() > 0 {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	entries := d.DLQEntries()
	if len(entries) == 0 {
		t.Fatal("DLQ 비어있음")
	}
	id := entries[0].ID
	beforeAttempts := entries[0].Attempts

	if err := d.DLQReplay(context.Background(), id); err == nil {
		t.Error("Replay 가 실패해야 하는데 성공")
	}
	after := d.DLQEntries()
	if len(after) != 1 {
		t.Fatalf("DLQ 크기 = %d, 1 기대 (실패 시 보관)", len(after))
	}
	if after[0].Attempts <= beforeAttempts {
		t.Errorf("Attempts 미증가: %d → %d", beforeAttempts, after[0].Attempts)
	}
}

func TestWebhookDispatcher_DLQ_ReplayNotFound(t *testing.T) {
	d, _ := tenancy.NewWebhookDispatcher(tenancy.WebhookConfig{URL: "http://x"}, nil)
	if err := d.DLQReplay(context.Background(), 9999); !errors.Is(err, tenancy.ErrDLQNotFound) {
		t.Errorf("err = %v, ErrDLQNotFound 기대", err)
	}
}

func TestWebhookDispatcher_DLQ_DropAndClear(t *testing.T) {
	d, _ := tenancy.NewWebhookDispatcher(tenancy.WebhookConfig{
		URL: "http://x", MaxRetries: 0, RetryBaseDelay: 5 * time.Millisecond,
	}, &flakeyDoer{failTimes: 999})
	d.Start(context.Background())
	defer d.Stop()

	for i := 0; i < 3; i++ {
		_ = d.Dispatch(context.Background(), tenancy.Event{Type: "x", TenantID: "t"})
	}
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if d.DLQCount() >= 3 {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if d.DLQCount() < 3 {
		t.Fatalf("DLQ = %d, 3 기대", d.DLQCount())
	}

	id := d.DLQEntries()[1].ID
	if !d.DLQDrop(id) {
		t.Error("DLQDrop 실패")
	}
	if d.DLQCount() != 2 {
		t.Errorf("Drop 후 = %d", d.DLQCount())
	}

	cleared := d.DLQClear()
	if cleared != 2 {
		t.Errorf("Clear 반환 = %d, 2 기대", cleared)
	}
	if d.DLQCount() != 0 {
		t.Errorf("Clear 후 = %d", d.DLQCount())
	}
}

func TestWebhookDispatcher_DLQ_MaxSizeFIFO(t *testing.T) {
	d, _ := tenancy.NewWebhookDispatcher(tenancy.WebhookConfig{
		URL: "http://x", MaxRetries: 0, RetryBaseDelay: 5 * time.Millisecond,
		DLQMaxSize: 2,
	}, &flakeyDoer{failTimes: 999})
	d.Start(context.Background())
	defer d.Stop()

	for i := 0; i < 5; i++ {
		_ = d.Dispatch(context.Background(), tenancy.Event{
			Type: "x", TenantID: "t",
			Payload: map[string]string{"n": string(rune('0' + i))},
		})
	}
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if d.Stats().DroppedCount >= 5 {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if d.DLQCount() != 2 {
		t.Errorf("DLQMaxSize=2 인데 size = %d", d.DLQCount())
	}
	// FIFO 검증 — 가장 최근 두 개만 남음 (n="3", n="4")
	entries := d.DLQEntries()
	if len(entries) == 2 {
		if entries[0].Event.Payload["n"] != "3" {
			t.Errorf("FIFO 위반: 첫 항목 n = %s, '3' 기대", entries[0].Event.Payload["n"])
		}
		if entries[1].Event.Payload["n"] != "4" {
			t.Errorf("FIFO 위반: 둘째 항목 n = %s, '4' 기대", entries[1].Event.Payload["n"])
		}
	}
}

// =============================================================================
// HTTP Handler 테스트
// =============================================================================

func TestWebhookDLQHandler_GETList(t *testing.T) {
	d, _ := tenancy.NewWebhookDispatcher(tenancy.WebhookConfig{
		URL: "http://x", MaxRetries: 0, RetryBaseDelay: 5 * time.Millisecond,
	}, &flakeyDoer{failTimes: 999})
	d.Start(context.Background())
	defer d.Stop()

	_ = d.Dispatch(context.Background(), tenancy.Event{Type: "x", TenantID: "t"})
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if d.DLQCount() > 0 {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}

	h := tenancy.NewWebhookDLQHandler(d)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/webhook/dlq", nil)
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["count"].(float64) < 1 {
		t.Errorf("count = %v", body["count"])
	}
}

func TestWebhookDLQHandler_DropByID(t *testing.T) {
	d, _ := tenancy.NewWebhookDispatcher(tenancy.WebhookConfig{
		URL: "http://x", MaxRetries: 0, RetryBaseDelay: 5 * time.Millisecond,
	}, &flakeyDoer{failTimes: 999})
	d.Start(context.Background())
	defer d.Stop()

	_ = d.Dispatch(context.Background(), tenancy.Event{Type: "x", TenantID: "t"})
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if d.DLQCount() > 0 {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	id := d.DLQEntries()[0].ID

	h := tenancy.NewWebhookDLQHandler(d)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/webhook/dlq/"+itoa(id), nil)
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, body=%s", w.Code, w.Body.String())
	}
	if d.DLQCount() != 0 {
		t.Errorf("드롭 후 size = %d", d.DLQCount())
	}
}

func TestWebhookDLQHandler_ReplayNotFound(t *testing.T) {
	d, _ := tenancy.NewWebhookDispatcher(tenancy.WebhookConfig{URL: "http://x"}, nil)
	h := tenancy.NewWebhookDLQHandler(d)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/webhook/dlq/9999/replay", nil)
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, 404 기대", w.Code)
	}
}

func TestWebhookDLQHandler_BadID(t *testing.T) {
	d, _ := tenancy.NewWebhookDispatcher(tenancy.WebhookConfig{URL: "http://x"}, nil)
	h := tenancy.NewWebhookDLQHandler(d)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/webhook/dlq/notanumber", nil)
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, 400 기대", w.Code)
	}
}

func TestWebhookDLQHandler_NoDispatcher(t *testing.T) {
	h := tenancy.NewWebhookDLQHandler(nil)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/webhook/dlq", nil)
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("status = %d, 503 기대", w.Code)
	}
}

func TestWebhookDLQHandler_PathPrefix(t *testing.T) {
	d, _ := tenancy.NewWebhookDispatcher(tenancy.WebhookConfig{URL: "http://x"}, nil)
	h := tenancy.NewWebhookDLQHandler(d)
	h.SetPathPrefix("/ops/tenancy")
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/ops/tenancy/webhook/dlq", nil)
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("path prefix 적용 안 됨: status = %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), `"count":0`) {
		t.Errorf("body = %s", w.Body.String())
	}
}

func itoa(i int64) string {
	if i == 0 {
		return "0"
	}
	neg := false
	if i < 0 {
		neg = true
		i = -i
	}
	var buf [20]byte
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}
