package llm

import (
	"context"
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// fakeExec 은 PostgresAuditLog 의 SQL 호출을 추적.
type fakeExec struct {
	mu      sync.Mutex
	calls   []fakeCall
	failErr error
}

type fakeCall struct {
	sql  string
	args []interface{}
}

type fakeExecResult struct{ rowsAffected int64 }

func (r fakeExecResult) RowsAffected() (int64, error) { return r.rowsAffected, nil }

func (f *fakeExec) ExecContext(_ context.Context, sql string, args ...interface{}) (AuditExecResult, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.failErr != nil {
		return nil, f.failErr
	}
	f.calls = append(f.calls, fakeCall{sql: sql, args: append([]interface{}{}, args...)})
	return fakeExecResult{rowsAffected: 1}, nil
}

type fakeRow struct {
	val int64
	err error
}

func (r *fakeRow) Scan(dest ...interface{}) error {
	if r.err != nil {
		return r.err
	}
	if len(dest) > 0 {
		if p, ok := dest[0].(*int64); ok {
			*p = r.val
		}
	}
	return nil
}

type fakeQuerier struct{ row *fakeRow }

func (q *fakeQuerier) QueryRowContext(_ context.Context, _ string, _ ...interface{}) AuditRow {
	return q.row
}

// ============================================================================
// PostgresAuditLog 테스트
// ============================================================================

func TestPostgresAuditLog_RecordsBasic(t *testing.T) {
	exec := &fakeExec{}
	log, err := NewPostgresAuditLog(exec, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := log.RecordLLMCall(context.Background(), LLMAuditEntry{
		TenantID:     "hospA",
		UserID:       "doc1",
		Provider:     "openai",
		Model:        "gpt-4",
		PromptTokens: 100,
		TotalTokens:  150,
		Success:      true,
	}); err != nil {
		t.Fatal(err)
	}
	if len(exec.calls) != 1 {
		t.Fatalf("calls = %d", len(exec.calls))
	}
	c := exec.calls[0]
	if !strings.Contains(c.sql, "INSERT INTO llm_audit") {
		t.Errorf("SQL = %s", c.sql)
	}
	if c.args[0] != "hospA" {
		t.Errorf("tenant_id = %v", c.args[0])
	}
}

func TestPostgresAuditLog_DefaultsTenantToPersonal(t *testing.T) {
	exec := &fakeExec{}
	log, _ := NewPostgresAuditLog(exec, nil)
	_ = log.RecordLLMCall(context.Background(), LLMAuditEntry{
		// TenantID 빈 값
	})
	if exec.calls[0].args[0] != "personal" {
		t.Errorf("tenant_id = %v, want personal", exec.calls[0].args[0])
	}
}

func TestPostgresAuditLog_NilExec(t *testing.T) {
	if _, err := NewPostgresAuditLog(nil, nil); err == nil {
		t.Error("nil exec 통과")
	}
}

func TestPostgresAuditLog_ExecError(t *testing.T) {
	exec := &fakeExec{failErr: errors.New("conn lost")}
	log, _ := NewPostgresAuditLog(exec, nil)
	if err := log.RecordLLMCall(context.Background(), LLMAuditEntry{TenantID: "t"}); err == nil {
		t.Error("DB 에러 통과")
	}
}

func TestPostgresAuditLog_TotalTokensByTenant(t *testing.T) {
	exec := &fakeExec{}
	q := &fakeQuerier{row: &fakeRow{val: 5000}}
	log, _ := NewPostgresAuditLog(exec, q)
	got, err := log.TotalTokensByTenant(context.Background(), "hospA")
	if err != nil {
		t.Fatal(err)
	}
	if got != 5000 {
		t.Errorf("total = %d", got)
	}
}

func TestPostgresAuditLog_TokensInWindow(t *testing.T) {
	q := &fakeQuerier{row: &fakeRow{val: 1234}}
	log, _ := NewPostgresAuditLog(&fakeExec{}, q)
	got, err := log.TokensInWindow(context.Background(), "hospA", 24)
	if err != nil {
		t.Fatal(err)
	}
	if got != 1234 {
		t.Errorf("window = %d", got)
	}
}

func TestPostgresAuditLog_NoQuerier_Error(t *testing.T) {
	log, _ := NewPostgresAuditLog(&fakeExec{}, nil)
	if _, err := log.TotalTokensByTenant(context.Background(), "x"); err == nil {
		t.Error("querier 미설정 통과")
	}
}

// ============================================================================
// BatchAuditLog 테스트
// ============================================================================

type recordingDelegate struct {
	mu    sync.Mutex
	count int32
	last  LLMAuditEntry
}

func (r *recordingDelegate) RecordLLMCall(_ context.Context, entry LLMAuditEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	atomic.AddInt32(&r.count, 1)
	r.last = entry
	return nil
}

func TestBatchAuditLog_BufferUntilMaxSize(t *testing.T) {
	d := &recordingDelegate{}
	b := NewBatchAuditLog(d, BatchAuditConfig{
		MaxSize: 3, FlushInterval: time.Hour, // 자동 flush 안 됨
	})
	for i := 0; i < 2; i++ {
		_ = b.RecordLLMCall(context.Background(), LLMAuditEntry{})
	}
	if atomic.LoadInt32(&d.count) != 0 {
		t.Errorf("MaxSize 미도달인데 flush: %d", d.count)
	}
	// 3번째 → MaxSize 도달 → flush
	_ = b.RecordLLMCall(context.Background(), LLMAuditEntry{})
	if atomic.LoadInt32(&d.count) != 3 {
		t.Errorf("flush count = %d, want 3", d.count)
	}
}

func TestBatchAuditLog_PeriodicFlush(t *testing.T) {
	d := &recordingDelegate{}
	b := NewBatchAuditLog(d, BatchAuditConfig{
		MaxSize: 1000, FlushInterval: 30 * time.Millisecond,
	})
	b.Start(context.Background())
	defer b.Stop()

	for i := 0; i < 5; i++ {
		_ = b.RecordLLMCall(context.Background(), LLMAuditEntry{})
	}
	// 주기적 flush 대기
	time.Sleep(80 * time.Millisecond)

	if atomic.LoadInt32(&d.count) < 5 {
		t.Errorf("주기 flush 안됨: %d", d.count)
	}
}

func TestBatchAuditLog_StopFlushesPending(t *testing.T) {
	d := &recordingDelegate{}
	b := NewBatchAuditLog(d, BatchAuditConfig{
		MaxSize: 1000, FlushInterval: time.Hour,
	})
	b.Start(context.Background())
	for i := 0; i < 4; i++ {
		_ = b.RecordLLMCall(context.Background(), LLMAuditEntry{})
	}
	if b.PendingCount() != 4 {
		t.Errorf("pending = %d", b.PendingCount())
	}
	b.Stop()
	if atomic.LoadInt32(&d.count) != 4 {
		t.Errorf("Stop flush 안됨: %d", d.count)
	}
}

func TestBatchAuditLog_NilDelegate(t *testing.T) {
	if NewBatchAuditLog(nil, BatchAuditConfig{}) != nil {
		t.Error("nil delegate → nil 반환 기대")
	}
}

func TestBatchAuditLog_StartIdempotent(t *testing.T) {
	d := &recordingDelegate{}
	b := NewBatchAuditLog(d, BatchAuditConfig{
		MaxSize: 100, FlushInterval: 50 * time.Millisecond,
	})
	b.Start(context.Background())
	b.Start(context.Background()) // 두 번째 호출 안전
	b.Stop()
}

func TestBatchAuditLog_OnFlushError(t *testing.T) {
	failingDelegate := &failingAuditLog{err: errors.New("INSERT failed")}
	var errCount int32
	b := NewBatchAuditLog(failingDelegate, BatchAuditConfig{
		MaxSize: 1, FlushInterval: time.Hour,
		OnFlushError: func(_ error) {
			atomic.AddInt32(&errCount, 1)
		},
	})
	_ = b.RecordLLMCall(context.Background(), LLMAuditEntry{})
	if atomic.LoadInt32(&errCount) != 1 {
		t.Errorf("OnFlushError = %d", errCount)
	}
}

type failingAuditLog struct{ err error }

func (f *failingAuditLog) RecordLLMCall(_ context.Context, _ LLMAuditEntry) error {
	return f.err
}
