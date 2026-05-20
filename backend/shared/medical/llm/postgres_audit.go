package llm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"
)

// AuditExecer 는 PostgresAuditLog 가 의존하는 최소 SQL 인터페이스.
//
// pgxpool.Pool 또는 database/sql 둘 다 만족하도록 추상화 (직접 의존 회피).
type AuditExecer interface {
	ExecContext(ctx context.Context, sql string, args ...interface{}) (AuditExecResult, error)
}

// AuditExecResult 는 ExecContext 결과.
type AuditExecResult interface {
	RowsAffected() (int64, error)
}

// AuditQuerier 는 통계 조회용 (옵션 — TotalTokensByTenant/Day 등).
type AuditQuerier interface {
	QueryRowContext(ctx context.Context, sql string, args ...interface{}) AuditRow
}

// AuditRow 는 단일 행 결과.
type AuditRow interface {
	Scan(dest ...interface{}) error
}

// PostgresAuditLog 는 TenancyAuditLog 인터페이스의 PostgreSQL 영속 구현.
//
// 매 호출마다 INSERT 가 발생하므로 고빈도 LLM 호출 환경에서는
// BatchAuditLog 로 wrapping 권장.
type PostgresAuditLog struct {
	exec    AuditExecer
	querier AuditQuerier // 옵션
}

// NewPostgresAuditLog 생성.
func NewPostgresAuditLog(exec AuditExecer, querier AuditQuerier) (*PostgresAuditLog, error) {
	if exec == nil {
		return nil, errors.New("exec 필수")
	}
	return &PostgresAuditLog{exec: exec, querier: querier}, nil
}

// RecordLLMCall 은 single-row INSERT.
func (l *PostgresAuditLog) RecordLLMCall(ctx context.Context, entry LLMAuditEntry) error {
	if entry.TenantID == "" {
		entry.TenantID = "personal" // NOT NULL 컬럼 보호
	}
	_, err := l.exec.ExecContext(ctx, `
		INSERT INTO llm_audit
			(tenant_id, user_id, provider, model,
			 prompt_tokens, completion_tokens, total_tokens,
			 latency_ms, success, error_message, recorded_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW())
	`,
		entry.TenantID, nullable(entry.UserID), nullable(entry.Provider),
		nullable(entry.Model),
		entry.PromptTokens, entry.CompletionTokens, entry.TotalTokens,
		entry.LatencyMs, entry.Success, nullable(entry.ErrorMessage),
	)
	if err != nil {
		return fmt.Errorf("INSERT llm_audit: %w", err)
	}
	return nil
}

// TotalTokensByTenant 는 tenant 별 누적 토큰 (성공한 호출만).
//
// querier 미설정 시 0 반환.
func (l *PostgresAuditLog) TotalTokensByTenant(ctx context.Context, tenantID string) (int64, error) {
	if l.querier == nil {
		return 0, errors.New("querier 미설정")
	}
	var total int64
	err := l.querier.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(total_tokens), 0)
		FROM llm_audit
		WHERE tenant_id = $1 AND success = TRUE
	`, tenantID).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

// TokensInWindow 는 tenant 의 최근 N 시간 토큰 사용량.
func (l *PostgresAuditLog) TokensInWindow(ctx context.Context, tenantID string, hours int) (int64, error) {
	if l.querier == nil {
		return 0, errors.New("querier 미설정")
	}
	if hours <= 0 {
		hours = 24
	}
	var total int64
	err := l.querier.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(total_tokens), 0)
		FROM llm_audit
		WHERE tenant_id = $1 AND success = TRUE
		  AND recorded_at >= NOW() - ($2 || ' hours')::INTERVAL
	`, tenantID, fmt.Sprintf("%d", hours)).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func nullable(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

// AuditFailureEntry 는 단일 실패 기록 (Phase AN-3).
type AuditFailureEntry struct {
	TenantID     string
	UserID       string
	Provider     string
	Model        string
	ErrorMessage string
	RecordedAt   time.Time
}

// RecentFailures 는 tenant 의 최근 실패 호출 목록 (Phase AN-3).
//
// limit=0 이면 10. limit > 100 이면 100 으로 클램프. querier 가 QueryContext
// 를 지원해야 함 (PgxAdapter 같은 QuotaQuerier 호환 구현).
func (l *PostgresAuditLog) RecentFailures(ctx context.Context, tenantID string, limit int) ([]AuditFailureEntry, error) {
	if l.querier == nil {
		return nil, errors.New("querier 미설정")
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	rowQuerier, ok := l.querier.(quotaQuerierForFailures)
	if !ok {
		return nil, errors.New("querier does not support QueryContext")
	}
	rows, err := rowQuerier.QueryContext(ctx, `
		SELECT tenant_id, user_id, provider, model, error_message, recorded_at
		FROM llm_audit
		WHERE tenant_id = $1 AND success = FALSE
		ORDER BY recorded_at DESC
		LIMIT $2
	`, tenantID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []AuditFailureEntry
	for rows.Next() {
		var e AuditFailureEntry
		var userID, provider, model, errMsg sql.NullString
		if err := rows.Scan(&e.TenantID, &userID, &provider, &model, &errMsg, &e.RecordedAt); err != nil {
			continue
		}
		if userID.Valid {
			e.UserID = userID.String
		}
		if provider.Valid {
			e.Provider = provider.String
		}
		if model.Valid {
			e.Model = model.String
		}
		if errMsg.Valid {
			e.ErrorMessage = errMsg.String
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

// quotaQuerierForFailures 는 QuotaQuerier 인터페이스 별칭.
type quotaQuerierForFailures interface {
	QueryContext(ctx context.Context, sql string, args ...interface{}) (QuotaRows, error)
}

// AuditTenantStats 는 tenant 별 운영 통계 스냅샷 (Phase AM-2).
type AuditTenantStats struct {
	TenantID        string
	TotalCalls      int64
	SuccessfulCalls int64
	FailedCalls     int64
	TotalTokens     int64
	TokensInWindow  int64
	WindowHours     int
}

// StatsByTenant 는 tenant 의 통합 통계 (윈도우 + 누적). 단일 SELECT 로 산출.
//
// windowHours 가 0 이면 24 사용.
func (l *PostgresAuditLog) StatsByTenant(ctx context.Context, tenantID string, windowHours int) (*AuditTenantStats, error) {
	if l.querier == nil {
		return nil, errors.New("querier 미설정")
	}
	if windowHours <= 0 {
		windowHours = 24
	}
	stats := &AuditTenantStats{
		TenantID:    tenantID,
		WindowHours: windowHours,
	}
	err := l.querier.QueryRowContext(ctx, `
		SELECT
			COUNT(*),
			COUNT(*) FILTER (WHERE success = TRUE),
			COUNT(*) FILTER (WHERE success = FALSE),
			COALESCE(SUM(total_tokens) FILTER (WHERE success = TRUE), 0),
			COALESCE(SUM(total_tokens) FILTER (
				WHERE success = TRUE AND recorded_at >= NOW() - ($2 || ' hours')::INTERVAL
			), 0)
		FROM llm_audit
		WHERE tenant_id = $1
	`, tenantID, fmt.Sprintf("%d", windowHours)).Scan(
		&stats.TotalCalls, &stats.SuccessfulCalls, &stats.FailedCalls,
		&stats.TotalTokens, &stats.TokensInWindow,
	)
	if err != nil {
		return nil, fmt.Errorf("StatsByTenant: %w", err)
	}
	return stats, nil
}

// ============================================================================
// BatchAuditLog — 버퍼링 + 주기적 flush (고빈도 환경)
// ============================================================================

// BatchAuditLog 는 LLM audit 호출을 메모리 버퍼에 모아 주기적으로 flush.
//
// 사용 예:
//
//	pgLog, _ := NewPostgresAuditLog(pool, pool)
//	batched := NewBatchAuditLog(pgLog, BatchAuditConfig{
//	    MaxSize: 1000, FlushInterval: 5*time.Second,
//	})
//	batched.Start(ctx)
//	adapter := NewTenancyAdapter(inner, batched)
//	defer batched.Stop()
type BatchAuditLog struct {
	delegate TenancyAuditLog
	cfg      BatchAuditConfig

	mu     sync.Mutex
	buf    []LLMAuditEntry
	stopCh chan struct{}
	doneCh chan struct{}
}

// BatchAuditConfig 는 배치 동작 설정.
type BatchAuditConfig struct {
	// MaxSize 버퍼 가득 차면 즉시 flush (기본 1000).
	MaxSize int
	// FlushInterval 주기적 flush 간격 (기본 5초).
	FlushInterval time.Duration
	// OnFlushError 는 flush 실패 콜백 (옵션).
	OnFlushError func(err error)
}

// NewBatchAuditLog 생성. delegate=nil 이면 nil 반환.
func NewBatchAuditLog(delegate TenancyAuditLog, cfg BatchAuditConfig) *BatchAuditLog {
	if delegate == nil {
		return nil
	}
	if cfg.MaxSize <= 0 {
		cfg.MaxSize = 1000
	}
	if cfg.FlushInterval <= 0 {
		cfg.FlushInterval = 5 * time.Second
	}
	return &BatchAuditLog{
		delegate: delegate,
		cfg:      cfg,
	}
}

// RecordLLMCall 은 버퍼에 추가; MaxSize 도달 시 즉시 flush.
func (b *BatchAuditLog) RecordLLMCall(ctx context.Context, entry LLMAuditEntry) error {
	b.mu.Lock()
	b.buf = append(b.buf, entry)
	shouldFlush := len(b.buf) >= b.cfg.MaxSize
	b.mu.Unlock()
	if shouldFlush {
		return b.Flush(ctx)
	}
	return nil
}

// Flush 는 버퍼를 모두 delegate 로 전달.
func (b *BatchAuditLog) Flush(ctx context.Context) error {
	b.mu.Lock()
	if len(b.buf) == 0 {
		b.mu.Unlock()
		return nil
	}
	toFlush := b.buf
	b.buf = nil
	b.mu.Unlock()

	for _, entry := range toFlush {
		if err := b.delegate.RecordLLMCall(ctx, entry); err != nil {
			b.reportErr(err)
		}
	}
	return nil
}

// Start 는 백그라운드 주기적 flush goroutine 시작.
func (b *BatchAuditLog) Start(ctx context.Context) {
	b.mu.Lock()
	if b.stopCh != nil {
		b.mu.Unlock()
		return
	}
	stopCh := make(chan struct{})
	doneCh := make(chan struct{})
	b.stopCh = stopCh
	b.doneCh = doneCh
	b.mu.Unlock()

	go b.run(ctx, stopCh, doneCh)
}

// Stop 은 flush 후 종료.
func (b *BatchAuditLog) Stop() {
	b.mu.Lock()
	if b.stopCh == nil {
		b.mu.Unlock()
		return
	}
	close(b.stopCh)
	doneCh := b.doneCh
	b.stopCh = nil
	b.mu.Unlock()
	if doneCh != nil {
		<-doneCh
	}
	_ = b.Flush(context.Background())
}

// PendingCount 는 현재 버퍼 크기.
func (b *BatchAuditLog) PendingCount() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.buf)
}

func (b *BatchAuditLog) run(parentCtx context.Context, stopCh, doneCh chan struct{}) {
	defer close(doneCh)
	ticker := time.NewTicker(b.cfg.FlushInterval)
	defer ticker.Stop()
	for {
		select {
		case <-stopCh:
			return
		case <-parentCtx.Done():
			return
		case <-ticker.C:
			_ = b.Flush(parentCtx)
		}
	}
}

func (b *BatchAuditLog) reportErr(err error) {
	if b.cfg.OnFlushError != nil {
		b.cfg.OnFlushError(err)
	}
}
