package tenancy

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// EventType 은 tenancy 이벤트 종류.
type EventType string

const (
	EventInvitationCreated     EventType = "invitation.created"
	EventInvitationAccepted    EventType = "invitation.accepted"
	EventInvitationRevoked     EventType = "invitation.revoked"
	EventMembershipCreated     EventType = "membership.created"
	EventMembershipRemoved     EventType = "membership.removed"
	EventMembershipRoleChanged EventType = "membership.role_changed"
	EventQuotaExceeded         EventType = "quota.exceeded" // Phase AN-2
)

// Event 는 webhook 으로 송신할 이벤트.
type Event struct {
	Type      EventType         `json:"type"`
	TenantID  string            `json:"tenant_id"`
	UserID    string            `json:"user_id,omitempty"`
	ActorID   string            `json:"actor_id,omitempty"` // 이벤트 트리거 사용자
	Timestamp time.Time         `json:"timestamp"`
	Payload   map[string]string `json:"payload,omitempty"`
}

// HTTPDoer 는 webhook 송신용 HTTP 클라이언트 (테스트 가능).
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// WebhookDispatcher 는 tenancy 이벤트를 외부 webhook 으로 비동기 송신.
//
// Slack 호환 (mode="slack") 또는 일반 webhook (mode="generic") 두 모드 지원.
// 송신 실패 시 OnError 콜백 + 재시도 (지수 백오프).
type WebhookDispatcher struct {
	cfg    WebhookConfig
	doer   HTTPDoer
	mu     sync.Mutex
	stopCh chan struct{}
	doneCh chan struct{}

	// 재시도 큐 (실패 이벤트)
	retryQueue []*pendingEvent

	// Dead Letter Queue (Phase AO-1) — max retries 초과 후 운영자 검사 대기
	dlq       []*DLQEntry
	dlqNextID int64 // 자동 증가 ID (atomic)

	// 통계 (Phase AK-3) — atomic 접근으로 race-safe
	sentCount    int64 // 송신 성공
	failedCount  int64 // 송신 실패 (단일 시도)
	retriedCount int64 // 재시도 시도 횟수 (누적)
	droppedCount int64 // max retries 초과 후 DLQ 로 이동 (영구 드롭 아님)
}

// DLQEntry 는 max retries 초과 후 dead letter queue 에 보관된 이벤트 (Phase AO-1).
//
// 운영자가 원인 조사 후 수동으로 재발송하거나 (DLQReplay) 폐기 (DLQDrop) 가능.
type DLQEntry struct {
	ID        int64     `json:"id"`
	Event     Event     `json:"event"`
	LastError string    `json:"last_error"`
	Attempts  int       `json:"attempts"` // 누적 시도 횟수 (dispatch + retries)
	AddedAt   time.Time `json:"added_at"`
}

type pendingEvent struct {
	Event   Event
	Attempt int
	NextAt  time.Time
}

// WebhookConfig 는 dispatcher 설정.
type WebhookConfig struct {
	URL            string        // 필수
	Mode           string        // "slack" | "generic" (기본 generic)
	Timeout        time.Duration // 단일 호출 타임아웃 (기본 5초)
	MaxRetries     int           // 재시도 횟수 (기본 3)
	RetryBaseDelay time.Duration // 첫 재시도 지연 (기본 1초, 지수 증가)
	OnError        func(ev Event, err error)

	// DLQMaxSize 는 dead letter queue 최대 크기 (Phase AO-1).
	// 초과 시 가장 오래된 항목부터 FIFO 제거. 0 이면 기본값 1000 사용.
	DLQMaxSize int
}

// NewWebhookDispatcher 생성. URL 미설정 시 nil + 에러.
func NewWebhookDispatcher(cfg WebhookConfig, doer HTTPDoer) (*WebhookDispatcher, error) {
	if cfg.URL == "" {
		return nil, errors.New("URL 필수")
	}
	if cfg.Mode == "" {
		cfg.Mode = "generic"
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 5 * time.Second
	}
	if cfg.MaxRetries < 0 {
		cfg.MaxRetries = 3
	}
	if cfg.RetryBaseDelay <= 0 {
		cfg.RetryBaseDelay = time.Second
	}
	if cfg.DLQMaxSize <= 0 {
		cfg.DLQMaxSize = 1000
	}
	if doer == nil {
		doer = &http.Client{Timeout: cfg.Timeout}
	}
	return &WebhookDispatcher{cfg: cfg, doer: doer}, nil
}

// Dispatch 는 이벤트 송신 (동기). 실패 시 재시도 큐에 추가.
func (d *WebhookDispatcher) Dispatch(ctx context.Context, ev Event) error {
	if ev.Timestamp.IsZero() {
		ev.Timestamp = time.Now().UTC()
	}
	if err := d.send(ctx, ev); err != nil {
		atomic.AddInt64(&d.failedCount, 1)
		d.enqueueRetry(ev)
		d.reportErr(ev, err)
		return err
	}
	atomic.AddInt64(&d.sentCount, 1)
	return nil
}

// DispatchAsync 는 비동기 송신 (호출자 블로킹 없음).
func (d *WebhookDispatcher) DispatchAsync(ev Event) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), d.cfg.Timeout)
		defer cancel()
		_ = d.Dispatch(ctx, ev)
	}()
}

// send 는 단일 HTTP POST.
func (d *WebhookDispatcher) send(ctx context.Context, ev Event) error {
	payload, err := d.encode(ev)
	if err != nil {
		return err
	}
	reqCtx, cancel := context.WithTimeout(ctx, d.cfg.Timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(reqCtx, "POST", d.cfg.URL, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := d.doer.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook HTTP %d", resp.StatusCode)
	}
	return nil
}

// encode 는 mode 별 페이로드 직렬화.
func (d *WebhookDispatcher) encode(ev Event) ([]byte, error) {
	switch d.cfg.Mode {
	case "slack":
		return json.Marshal(map[string]interface{}{
			"text": fmt.Sprintf("*[Manpasik Tenancy]* %s tenant=%s actor=%s",
				ev.Type, ev.TenantID, ev.ActorID),
			"attachments": []map[string]interface{}{
				{
					"color": slackColorFor(ev.Type),
					"fields": []map[string]string{
						{"title": "Type", "value": string(ev.Type), "short": "true"},
						{"title": "Tenant", "value": ev.TenantID, "short": "true"},
						{"title": "User", "value": ev.UserID, "short": "true"},
						{"title": "Actor", "value": ev.ActorID, "short": "true"},
						{"title": "Timestamp", "value": ev.Timestamp.Format(time.RFC3339), "short": "false"},
					},
				},
			},
		})
	default:
		return json.Marshal(ev)
	}
}

// slackColorFor 는 이벤트 타입별 색상.
func slackColorFor(t EventType) string {
	switch t {
	case EventQuotaExceeded:
		return "danger" // 빨강 (경고 — 즉시 운영자 확인)
	case EventInvitationCreated, EventMembershipCreated:
		return "good" // 초록
	case EventInvitationRevoked, EventMembershipRemoved:
		return "warning" // 주황
	default:
		return "#3AA3E3" // 파랑
	}
}

// enqueueRetry 는 실패 이벤트를 재시도 큐에 추가.
func (d *WebhookDispatcher) enqueueRetry(ev Event) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.retryQueue = append(d.retryQueue, &pendingEvent{
		Event:   ev,
		Attempt: 1,
		NextAt:  time.Now().Add(d.cfg.RetryBaseDelay),
	})
}

// processRetries 는 재시도 큐 처리.
func (d *WebhookDispatcher) processRetries(ctx context.Context) {
	d.mu.Lock()
	now := time.Now()
	pending := d.retryQueue
	d.retryQueue = nil
	var stillPending []*pendingEvent
	d.mu.Unlock()

	for _, p := range pending {
		if now.Before(p.NextAt) {
			stillPending = append(stillPending, p)
			continue
		}
		atomic.AddInt64(&d.retriedCount, 1)
		err := d.send(ctx, p.Event)
		if err == nil {
			atomic.AddInt64(&d.sentCount, 1)
			continue // 성공
		}
		p.Attempt++
		if p.Attempt > d.cfg.MaxRetries {
			atomic.AddInt64(&d.droppedCount, 1)
			// Phase AO-1: DLQ 보관 (즉시 폐기 대신 운영자 재발송/폐기 게이트)
			d.enqueueDLQ(p.Event, p.Attempt, err.Error())
			d.reportErr(p.Event, errors.New("max retries exceeded"))
			continue
		}
		// 지수 백오프
		backoff := d.cfg.RetryBaseDelay * time.Duration(1<<uint(p.Attempt-1))
		p.NextAt = time.Now().Add(backoff)
		stillPending = append(stillPending, p)
	}
	if len(stillPending) > 0 {
		d.mu.Lock()
		d.retryQueue = append(d.retryQueue, stillPending...)
		d.mu.Unlock()
	}
}

// Start 는 백그라운드 재시도 처리 goroutine 시작.
func (d *WebhookDispatcher) Start(ctx context.Context) {
	d.mu.Lock()
	if d.stopCh != nil {
		d.mu.Unlock()
		return
	}
	stopCh := make(chan struct{})
	doneCh := make(chan struct{})
	d.stopCh = stopCh
	d.doneCh = doneCh
	d.mu.Unlock()

	go func() {
		defer close(doneCh)
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-stopCh:
				return
			case <-ctx.Done():
				return
			case <-ticker.C:
				d.processRetries(ctx)
			}
		}
	}()
}

// Stop 은 재시도 처리 종료.
func (d *WebhookDispatcher) Stop() {
	d.mu.Lock()
	if d.stopCh == nil {
		d.mu.Unlock()
		return
	}
	close(d.stopCh)
	doneCh := d.doneCh
	d.stopCh = nil
	d.mu.Unlock()
	if doneCh != nil {
		<-doneCh
	}
}

// PendingRetryCount 는 현재 재시도 대기 중인 이벤트 수.
func (d *WebhookDispatcher) PendingRetryCount() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.retryQueue)
}

// enqueueDLQ 는 max retries 초과 이벤트를 DLQ 에 추가 (Phase AO-1).
//
// DLQMaxSize 초과 시 가장 오래된 항목 제거 (FIFO).
func (d *WebhookDispatcher) enqueueDLQ(ev Event, attempts int, lastErr string) {
	id := atomic.AddInt64(&d.dlqNextID, 1)
	entry := &DLQEntry{
		ID:        id,
		Event:     ev,
		LastError: lastErr,
		Attempts:  attempts,
		AddedAt:   time.Now().UTC(),
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	d.dlq = append(d.dlq, entry)
	if d.cfg.DLQMaxSize > 0 && len(d.dlq) > d.cfg.DLQMaxSize {
		drop := len(d.dlq) - d.cfg.DLQMaxSize
		d.dlq = d.dlq[drop:]
	}
}

// DLQEntries 는 현재 DLQ 스냅샷 반환 (복사본).
func (d *WebhookDispatcher) DLQEntries() []DLQEntry {
	d.mu.Lock()
	defer d.mu.Unlock()
	out := make([]DLQEntry, len(d.dlq))
	for i, e := range d.dlq {
		out[i] = *e
	}
	return out
}

// DLQCount 은 현재 DLQ 크기.
func (d *WebhookDispatcher) DLQCount() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.dlq)
}

// DLQReplay 는 DLQ 항목 재발송. 성공 시 DLQ 에서 제거, 실패 시 last_error/attempts 갱신.
//
// ID 가 존재하지 않으면 ErrDLQNotFound 반환.
func (d *WebhookDispatcher) DLQReplay(ctx context.Context, id int64) error {
	d.mu.Lock()
	var ev Event
	found := false
	for _, e := range d.dlq {
		if e.ID == id {
			ev = e.Event
			found = true
			break
		}
	}
	d.mu.Unlock()
	if !found {
		return ErrDLQNotFound
	}
	if err := d.send(ctx, ev); err != nil {
		// 재발송 실패 — DLQ 항목 갱신
		d.mu.Lock()
		for _, e := range d.dlq {
			if e.ID == id {
				e.LastError = err.Error()
				e.Attempts++
				break
			}
		}
		d.mu.Unlock()
		atomic.AddInt64(&d.failedCount, 1)
		return err
	}
	// 성공 — DLQ 에서 제거
	atomic.AddInt64(&d.sentCount, 1)
	d.mu.Lock()
	defer d.mu.Unlock()
	for i, e := range d.dlq {
		if e.ID == id {
			d.dlq = append(d.dlq[:i], d.dlq[i+1:]...)
			return nil
		}
	}
	return nil
}

// DLQDrop 은 ID 로 DLQ 항목 제거 (재발송 없음). 미존재 시 false.
func (d *WebhookDispatcher) DLQDrop(id int64) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	for i, e := range d.dlq {
		if e.ID == id {
			d.dlq = append(d.dlq[:i], d.dlq[i+1:]...)
			return true
		}
	}
	return false
}

// DLQClear 는 DLQ 전체 비움 (제거된 항목 수 반환).
func (d *WebhookDispatcher) DLQClear() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	n := len(d.dlq)
	d.dlq = nil
	return n
}

// ErrDLQNotFound 는 DLQReplay 시 ID 미존재.
var ErrDLQNotFound = errors.New("DLQ 항목 없음")

// WebhookStats 는 운영 통계 스냅샷 (Phase AK-3).
type WebhookStats struct {
	SentCount    int64 // 송신 성공 (재시도 성공 포함)
	FailedCount  int64 // 단일 시도 실패 (재시도 큐에 들어간 횟수)
	RetriedCount int64 // 재시도 시도 (누적)
	DroppedCount int64 // max retries 초과 후 DLQ 로 이동
	PendingRetry int   // 현재 재시도 큐 크기
	DLQSize      int   // 현재 DLQ 크기 (Phase AO-1)
}

// Stats 는 현재 통계 반환 (atomic 읽기).
func (d *WebhookDispatcher) Stats() WebhookStats {
	return WebhookStats{
		SentCount:    atomic.LoadInt64(&d.sentCount),
		FailedCount:  atomic.LoadInt64(&d.failedCount),
		RetriedCount: atomic.LoadInt64(&d.retriedCount),
		DroppedCount: atomic.LoadInt64(&d.droppedCount),
		PendingRetry: d.PendingRetryCount(),
		DLQSize:      d.DLQCount(),
	}
}

// CollectMetrics 는 dashboard.MetricsCollector 호환 — 운영 대시보드 노출.
//
// status 는 DLQSize > 0 또는 PendingRetry 가 100 이상이면 degraded, 그 외 healthy.
// (DroppedCount 는 누적 카운터라 시간이 지나도 0 으로 돌아오지 않으므로 status
// 게이팅에서 제외 — DLQ 가 비어 있으면 운영자가 처리 완료한 상태.)
func (d *WebhookDispatcher) CollectMetrics() map[string]interface{} {
	s := d.Stats()
	status := "healthy"
	if s.DLQSize > 0 {
		status = "degraded"
	} else if s.PendingRetry >= 100 {
		status = "degraded"
	}
	return map[string]interface{}{
		"status":        status,
		"sent_count":    s.SentCount,
		"failed_count":  s.FailedCount,
		"retried_count": s.RetriedCount,
		"dropped_count": s.DroppedCount,
		"pending_retry": s.PendingRetry,
		"dlq_size":      s.DLQSize,
	}
}

func (d *WebhookDispatcher) reportErr(ev Event, err error) {
	if d.cfg.OnError != nil {
		d.cfg.OnError(ev, err)
	}
}

// WriteWebhookPrometheusMetrics 는 dispatcher 통계를 Prometheus text 형식 출력.
//
// 메트릭:
//
//	manpasik_tenancy_webhook_sent_total      counter
//	manpasik_tenancy_webhook_failed_total    counter
//	manpasik_tenancy_webhook_retried_total   counter
//	manpasik_tenancy_webhook_dropped_total   counter
//	manpasik_tenancy_webhook_pending_retry   gauge
//	manpasik_tenancy_webhook_status          gauge (1=healthy 0=degraded)
func WriteWebhookPrometheusMetrics(w http.ResponseWriter, d *WebhookDispatcher) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if d == nil {
		_, _ = w.Write([]byte("# HELP manpasik_tenancy_webhook_status webhook dispatcher 상태\n"))
		_, _ = w.Write([]byte("# TYPE manpasik_tenancy_webhook_status gauge\n"))
		_, _ = w.Write([]byte("manpasik_tenancy_webhook_status 0\n"))
		return
	}
	s := d.Stats()
	type metric struct {
		name  string
		help  string
		typ   string
		value int64
	}
	metrics := []metric{
		{"manpasik_tenancy_webhook_sent_total", "송신 성공 누적", "counter", s.SentCount},
		{"manpasik_tenancy_webhook_failed_total", "단일 시도 실패 누적", "counter", s.FailedCount},
		{"manpasik_tenancy_webhook_retried_total", "재시도 시도 누적", "counter", s.RetriedCount},
		{"manpasik_tenancy_webhook_dropped_total", "max retries 초과 후 DLQ 이동 누적", "counter", s.DroppedCount},
		{"manpasik_tenancy_webhook_pending_retry", "현재 재시도 큐 크기", "gauge", int64(s.PendingRetry)},
		{"manpasik_tenancy_webhook_dlq_size", "현재 DLQ 크기", "gauge", int64(s.DLQSize)},
	}
	for _, m := range metrics {
		_, _ = w.Write([]byte("# HELP " + m.name + " " + m.help + "\n"))
		_, _ = w.Write([]byte("# TYPE " + m.name + " " + m.typ + "\n"))
		_, _ = w.Write([]byte(m.name + " " + intStr(int(m.value)) + "\n"))
	}
	statusVal := "1"
	if s.DLQSize > 0 || s.PendingRetry >= 100 {
		statusVal = "0"
	}
	_, _ = w.Write([]byte("# HELP manpasik_tenancy_webhook_status 1=healthy 0=degraded\n"))
	_, _ = w.Write([]byte("# TYPE manpasik_tenancy_webhook_status gauge\n"))
	_, _ = w.Write([]byte("manpasik_tenancy_webhook_status " + statusVal + "\n"))
}
