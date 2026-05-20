package llm

import (
	"context"
	"errors"
	"sync"
	"time"
)

// QuotaConfig 는 tenant 별 LLM 사용 한도 설정.
type QuotaConfig struct {
	TenantID           string
	DailyTokenLimit    int    // 0 = 무제한
	MonthlyTokenLimit  int    // 0 = 무제한
	DailyRequestLimit  int    // 0 = 무제한
	UpdatedAt          time.Time
}

// IsUnlimited 는 모든 한도가 0 인지.
func (c *QuotaConfig) IsUnlimited() bool {
	return c.DailyTokenLimit == 0 && c.MonthlyTokenLimit == 0 && c.DailyRequestLimit == 0
}

// QuotaStore 는 tenant 별 quota 설정 저장소.
type QuotaStore interface {
	Get(ctx context.Context, tenantID string) (*QuotaConfig, error)
	Set(ctx context.Context, cfg QuotaConfig) error
	Delete(ctx context.Context, tenantID string) error
	List(ctx context.Context) ([]*QuotaConfig, error)
}

// MemoryQuotaStore 는 sync.RWMutex 기반 인메모리 저장소.
type MemoryQuotaStore struct {
	mu      sync.RWMutex
	configs map[string]*QuotaConfig
}

// NewMemoryQuotaStore 생성.
func NewMemoryQuotaStore() *MemoryQuotaStore {
	return &MemoryQuotaStore{configs: make(map[string]*QuotaConfig)}
}

// Get 은 tenant config 반환. 없으면 nil + 에러 (호출자가 default 적용).
func (s *MemoryQuotaStore) Get(_ context.Context, tenantID string) (*QuotaConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.configs[tenantID]
	if !ok {
		return nil, ErrQuotaConfigNotFound
	}
	cp := *c
	return &cp, nil
}

// Set 은 config 저장 (UPSERT).
func (s *MemoryQuotaStore) Set(_ context.Context, cfg QuotaConfig) error {
	if cfg.TenantID == "" {
		return errors.New("TenantID 필수")
	}
	if cfg.UpdatedAt.IsZero() {
		cfg.UpdatedAt = time.Now()
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	cp := cfg
	s.configs[cfg.TenantID] = &cp
	return nil
}

// Delete 는 config 제거.
func (s *MemoryQuotaStore) Delete(_ context.Context, tenantID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.configs, tenantID)
	return nil
}

// List 는 모든 config.
func (s *MemoryQuotaStore) List(_ context.Context) ([]*QuotaConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*QuotaConfig, 0, len(s.configs))
	for _, c := range s.configs {
		cp := *c
		out = append(out, &cp)
	}
	return out, nil
}

// ErrQuotaConfigNotFound 는 store 에서 tenant 를 찾을 수 없음.
var ErrQuotaConfigNotFound = errors.New("quota config not found")

// ============================================================================
// DynamicQuota — store 기반 quota (TTL 캐시 포함)
// ============================================================================

// QuotaViolationNotifier 는 한도 초과 시 호출되는 인터페이스 (Phase AN-2).
//
// tenancy.WebhookDispatcher 같은 외부 모듈을 직접 참조하면 순환 의존이 발생하므로
// 인터페이스로 추상화. Notify 는 비동기 송신 가정 (블로킹 안 됨).
type QuotaViolationNotifier interface {
	NotifyQuotaExceeded(tenantID string, limit int, used int64, limitType string)
}

// DynamicQuota 는 QuotaStore 와 AuditLog 를 조합하여 tenant 별 동적 한도 적용.
//
// AuditLog.TokensInWindow 로 실시간 사용량 조회 + QuotaStore 의 한도와 비교.
type DynamicQuota struct {
	store    QuotaStore
	auditLog *PostgresAuditLog // 옵션 — TokensInWindow 사용
	cacheTTL time.Duration
	notifier QuotaViolationNotifier // 옵션 — 한도 초과 시 알림

	mu         sync.RWMutex
	cache      map[string]quotaCacheEntry
	defaultCfg *QuotaConfig // 등록되지 않은 tenant 기본값
}

type quotaCacheEntry struct {
	cfg       *QuotaConfig
	expiresAt time.Time
}

// NewDynamicQuota 생성. cacheTTL=0 이면 60초 기본.
func NewDynamicQuota(store QuotaStore, auditLog *PostgresAuditLog, cacheTTL time.Duration) *DynamicQuota {
	if cacheTTL <= 0 {
		cacheTTL = 60 * time.Second
	}
	return &DynamicQuota{
		store:    store,
		auditLog: auditLog,
		cacheTTL: cacheTTL,
		cache:    make(map[string]quotaCacheEntry),
	}
}

// SetDefault 는 미등록 tenant 의 기본 한도 (예: free tier).
func (q *DynamicQuota) SetDefault(cfg *QuotaConfig) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.defaultCfg = cfg
}

// SetViolationNotifier 는 한도 초과 시 호출될 알림기 등록 (Phase AN-2).
//
// 미설정 시 알림 비활성. notifier.NotifyQuotaExceeded 는 비동기 송신 가정.
func (q *DynamicQuota) SetViolationNotifier(n QuotaViolationNotifier) {
	q.notifier = n
}

// CheckAllowed 는 store 의 한도 + audit 의 사용량을 비교.
//
// store 미연결 / 한도 미설정 시 항상 허용 (운영 default = lenient).
// 한도 초과 시 SetViolationNotifier 등록된 notifier 자동 호출.
func (q *DynamicQuota) CheckAllowed(tenantID string) bool {
	cfg := q.loadConfig(tenantID)
	if cfg == nil || cfg.IsUnlimited() {
		return true
	}
	if q.auditLog == nil {
		return true // 사용량 추적 불가 → 허용
	}
	// 일일 한도 검사
	if cfg.DailyTokenLimit > 0 {
		used, err := q.auditLog.TokensInWindow(context.Background(), tenantID, 24)
		if err == nil && used >= int64(cfg.DailyTokenLimit) {
			q.notifyViolation(tenantID, cfg.DailyTokenLimit, used, "daily_token")
			return false
		}
	}
	// 월간 한도 검사 (24*30시간)
	if cfg.MonthlyTokenLimit > 0 {
		used, err := q.auditLog.TokensInWindow(context.Background(), tenantID, 24*30)
		if err == nil && used >= int64(cfg.MonthlyTokenLimit) {
			q.notifyViolation(tenantID, cfg.MonthlyTokenLimit, used, "monthly_token")
			return false
		}
	}
	return true
}

// notifyViolation 는 한도 초과 알림 (notifier 미설정 시 no-op).
func (q *DynamicQuota) notifyViolation(tenantID string, limit int, used int64, limitType string) {
	if q.notifier != nil {
		q.notifier.NotifyQuotaExceeded(tenantID, limit, used, limitType)
	}
}

// RecordUsage 는 TenancyQuota 인터페이스 호환 — DynamicQuota 는 audit_log 를
// 통한 추적이므로 no-op (감사 로그가 별도 기록).
func (q *DynamicQuota) RecordUsage(_ string, _ int) {
	// audit_log 가 별도로 기록함
}

// loadConfig 는 캐시 우선 → store → default 순으로 config 조회.
func (q *DynamicQuota) loadConfig(tenantID string) *QuotaConfig {
	q.mu.RLock()
	if entry, ok := q.cache[tenantID]; ok && time.Now().Before(entry.expiresAt) {
		q.mu.RUnlock()
		return entry.cfg
	}
	q.mu.RUnlock()

	if q.store == nil {
		return q.defaultCfg
	}
	cfg, err := q.store.Get(context.Background(), tenantID)
	if err != nil {
		// 미등록 → default 사용
		q.mu.Lock()
		q.cache[tenantID] = quotaCacheEntry{
			cfg:       q.defaultCfg,
			expiresAt: time.Now().Add(q.cacheTTL),
		}
		q.mu.Unlock()
		return q.defaultCfg
	}
	q.mu.Lock()
	q.cache[tenantID] = quotaCacheEntry{
		cfg:       cfg,
		expiresAt: time.Now().Add(q.cacheTTL),
	}
	q.mu.Unlock()
	return cfg
}

// InvalidateCache 는 특정 tenant 의 캐시 강제 만료 (운영 시 한도 변경 후 호출).
func (q *DynamicQuota) InvalidateCache(tenantID string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	delete(q.cache, tenantID)
}

// ClearCache 는 모든 캐시 제거.
func (q *DynamicQuota) ClearCache() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.cache = make(map[string]quotaCacheEntry)
}
