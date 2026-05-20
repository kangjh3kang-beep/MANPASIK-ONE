package llm

import (
	"context"
	"errors"

	"github.com/manpasik/backend/shared/tenancy"
)

// TenancyAuditLog 는 LLM 호출의 tenant/user/cost 감사 추적용 인터페이스.
//
// 운영 환경에서는 audit-service 또는 별도 감사 DB 에 영속화. 미주입 시
// 감사 비활성화 (운영 권장 X — 의료 LLM 호출은 항상 감사).
type TenancyAuditLog interface {
	RecordLLMCall(ctx context.Context, entry LLMAuditEntry) error
}

// LLMAuditEntry 는 단일 LLM 호출 감사 기록.
type LLMAuditEntry struct {
	TenantID         string // ctx 에서 추출 (없으면 "personal")
	UserID           string // req.UserID
	Provider         string
	Model            string
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
	LatencyMs        int64
	Success          bool
	ErrorMessage     string
}

// TenancyAdapter 는 다른 Adapter 를 감싸 tenant 컨텍스트 추적을 추가.
//
// 동작:
//  1. ctx 의 tenant_id 추출 (없으면 "personal")
//  2. inner adapter.Complete 호출
//  3. 결과/에러를 audit log 에 기록
//  4. tenant 별 토큰 사용량 카운터 누적 (선택, 옵저버블 메트릭용)
type TenancyAdapter struct {
	inner    Adapter
	auditLog TenancyAuditLog
	// 대안: tenant 별 quota 적용 (운영 옵션, 미설정 시 무제한)
	quota TenancyQuota
}

// TenancyQuota 는 tenant 별 호출 한도 검사기.
type TenancyQuota interface {
	// CheckAllowed 는 tenant 의 한도 초과 시 false 반환.
	// CheckAllowed 호출은 호출 전, RecordUsage 는 호출 후.
	CheckAllowed(tenantID string) bool
	RecordUsage(tenantID string, tokens int)
}

// NewTenancyAdapter 생성. inner 는 필수, auditLog/quota 는 옵션.
func NewTenancyAdapter(inner Adapter, auditLog TenancyAuditLog) *TenancyAdapter {
	return &TenancyAdapter{
		inner:    inner,
		auditLog: auditLog,
	}
}

// SetQuota 는 tenant 별 한도 검사기 등록.
func (a *TenancyAdapter) SetQuota(quota TenancyQuota) { a.quota = quota }

// Provider 는 inner 의 provider 명에 "+tenancy" 접미사 부착.
func (a *TenancyAdapter) Provider() string {
	if a.inner == nil {
		return "tenancy"
	}
	return a.inner.Provider() + "+tenancy"
}

// HealthCheck 는 inner 위임.
func (a *TenancyAdapter) HealthCheck(ctx context.Context) error {
	if a.inner == nil {
		return errors.New("inner adapter nil")
	}
	return a.inner.HealthCheck(ctx)
}

// Complete 는 ctx 의 tenant 추적 + audit + quota 검사 후 inner 호출.
func (a *TenancyAdapter) Complete(ctx context.Context, req *Request) (*Response, error) {
	if a.inner == nil {
		return nil, errors.New("inner adapter nil")
	}
	if req == nil {
		return nil, errors.New("request nil")
	}

	tenantID := "personal"
	if tid, ok := tenancy.TenantFromContext(ctx); ok && !tid.IsZero() {
		tenantID = string(tid)
	}

	// Quota 검사 (옵션)
	if a.quota != nil && !a.quota.CheckAllowed(tenantID) {
		entry := LLMAuditEntry{
			TenantID:     tenantID,
			UserID:       req.UserID,
			Success:      false,
			ErrorMessage: "quota exceeded",
		}
		_ = a.recordAudit(ctx, entry)
		return nil, errors.New("LLM quota exceeded for tenant: " + tenantID)
	}

	resp, err := a.inner.Complete(ctx, req)

	entry := LLMAuditEntry{
		TenantID: tenantID,
		UserID:   req.UserID,
	}
	if resp != nil {
		entry.Provider = resp.Provider
		entry.Model = resp.Model
		entry.PromptTokens = resp.PromptTokens
		entry.CompletionTokens = resp.CompletionTokens
		// TotalTokens 가 0 이면 prompt+completion 으로 보정 (provider 별 차이 흡수)
		if resp.TotalTokens > 0 {
			entry.TotalTokens = resp.TotalTokens
		} else {
			entry.TotalTokens = resp.PromptTokens + resp.CompletionTokens
		}
		entry.LatencyMs = resp.LatencyMs
	}
	if err != nil {
		entry.Success = false
		entry.ErrorMessage = err.Error()
	} else {
		entry.Success = true
	}
	_ = a.recordAudit(ctx, entry)

	if a.quota != nil && err == nil && resp != nil {
		a.quota.RecordUsage(tenantID, entry.TotalTokens)
	}

	return resp, err
}

func (a *TenancyAdapter) recordAudit(ctx context.Context, entry LLMAuditEntry) error {
	if a.auditLog == nil {
		return nil
	}
	return a.auditLog.RecordLLMCall(ctx, entry)
}

// ============================================================================
// 인메모리 구현 (테스트/개발용)
// ============================================================================

// MemoryAuditLog 는 LLM 호출 기록을 메모리에 저장하는 감사 로그.
type MemoryAuditLog struct {
	entries []LLMAuditEntry
}

// NewMemoryAuditLog 생성.
func NewMemoryAuditLog() *MemoryAuditLog { return &MemoryAuditLog{} }

// RecordLLMCall 은 entry 를 메모리에 추가.
func (l *MemoryAuditLog) RecordLLMCall(_ context.Context, entry LLMAuditEntry) error {
	l.entries = append(l.entries, entry)
	return nil
}

// Entries 는 누적된 모든 entry 반환.
func (l *MemoryAuditLog) Entries() []LLMAuditEntry {
	out := make([]LLMAuditEntry, len(l.entries))
	copy(out, l.entries)
	return out
}

// EntriesByTenant 는 tenant 별 필터링.
func (l *MemoryAuditLog) EntriesByTenant(tenantID string) []LLMAuditEntry {
	var out []LLMAuditEntry
	for _, e := range l.entries {
		if e.TenantID == tenantID {
			out = append(out, e)
		}
	}
	return out
}

// TotalTokensByTenant 는 tenant 별 누적 토큰.
func (l *MemoryAuditLog) TotalTokensByTenant(tenantID string) int {
	sum := 0
	for _, e := range l.entries {
		if e.TenantID == tenantID && e.Success {
			sum += e.TotalTokens
		}
	}
	return sum
}

// ============================================================================
// 단순 Quota 구현 (테스트/개발용)
// ============================================================================

// MemoryQuota 는 tenant 당 토큰 한도를 메모리에서 추적.
type MemoryQuota struct {
	limit      int            // tenant 당 최대 토큰 (0 = 무제한)
	usage      map[string]int // tenant_id → 누적 토큰
}

// NewMemoryQuota 생성. limit=0 이면 무제한.
func NewMemoryQuota(limit int) *MemoryQuota {
	return &MemoryQuota{limit: limit, usage: make(map[string]int)}
}

// CheckAllowed 는 quota 검사.
func (q *MemoryQuota) CheckAllowed(tenantID string) bool {
	if q.limit <= 0 {
		return true
	}
	return q.usage[tenantID] < q.limit
}

// RecordUsage 는 누적 토큰 기록.
func (q *MemoryQuota) RecordUsage(tenantID string, tokens int) {
	q.usage[tenantID] += tokens
}

// Usage 반환 (tenant 별).
func (q *MemoryQuota) Usage(tenantID string) int { return q.usage[tenantID] }
