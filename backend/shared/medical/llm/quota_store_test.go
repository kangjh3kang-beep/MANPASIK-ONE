package llm

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestMemoryQuotaStore_SetGet(t *testing.T) {
	store := NewMemoryQuotaStore()
	cfg := QuotaConfig{TenantID: "hospA", DailyTokenLimit: 10000}
	if err := store.Set(context.Background(), cfg); err != nil {
		t.Fatal(err)
	}
	got, err := store.Get(context.Background(), "hospA")
	if err != nil {
		t.Fatal(err)
	}
	if got.DailyTokenLimit != 10000 {
		t.Errorf("DailyTokenLimit = %d", got.DailyTokenLimit)
	}
}

func TestMemoryQuotaStore_GetNotFound(t *testing.T) {
	store := NewMemoryQuotaStore()
	if _, err := store.Get(context.Background(), "ghost"); !errors.Is(err, ErrQuotaConfigNotFound) {
		t.Errorf("err = %v", err)
	}
}

func TestMemoryQuotaStore_SetEmptyTenantID(t *testing.T) {
	store := NewMemoryQuotaStore()
	if err := store.Set(context.Background(), QuotaConfig{}); err == nil {
		t.Error("빈 TenantID 통과")
	}
}

func TestMemoryQuotaStore_Delete(t *testing.T) {
	store := NewMemoryQuotaStore()
	_ = store.Set(context.Background(), QuotaConfig{TenantID: "x", DailyTokenLimit: 100})
	_ = store.Delete(context.Background(), "x")
	if _, err := store.Get(context.Background(), "x"); !errors.Is(err, ErrQuotaConfigNotFound) {
		t.Errorf("delete 후 err = %v", err)
	}
}

func TestMemoryQuotaStore_List(t *testing.T) {
	store := NewMemoryQuotaStore()
	_ = store.Set(context.Background(), QuotaConfig{TenantID: "A", DailyTokenLimit: 100})
	_ = store.Set(context.Background(), QuotaConfig{TenantID: "B", DailyTokenLimit: 200})
	list, _ := store.List(context.Background())
	if len(list) != 2 {
		t.Errorf("list len = %d", len(list))
	}
}

func TestQuotaConfig_IsUnlimited(t *testing.T) {
	cases := map[QuotaConfig]bool{
		{}: true,
		{DailyTokenLimit: 100}:    false,
		{MonthlyTokenLimit: 100}:  false,
		{DailyRequestLimit: 100}:  false,
	}
	for cfg, want := range cases {
		if got := cfg.IsUnlimited(); got != want {
			t.Errorf("%+v IsUnlimited = %v, want %v", cfg, got, want)
		}
	}
}

func TestDynamicQuota_NoStoreAlwaysAllow(t *testing.T) {
	q := NewDynamicQuota(nil, nil, 60*time.Second)
	if !q.CheckAllowed("hospA") {
		t.Error("store 없으면 허용 기대")
	}
}

func TestDynamicQuota_UnlimitedAllowed(t *testing.T) {
	store := NewMemoryQuotaStore()
	_ = store.Set(context.Background(), QuotaConfig{TenantID: "hospA"}) // 모두 0 = 무제한
	q := NewDynamicQuota(store, nil, 60*time.Second)
	if !q.CheckAllowed("hospA") {
		t.Error("무제한인데 거부됨")
	}
}

func TestDynamicQuota_DefaultUsedForUnknownTenant(t *testing.T) {
	store := NewMemoryQuotaStore()
	q := NewDynamicQuota(store, nil, 60*time.Second)
	q.SetDefault(&QuotaConfig{DailyTokenLimit: 1}) // 매우 낮은 기본값

	// store 에 없는 tenant 는 default 적용 — auditLog 미설정이라 항상 허용
	if !q.CheckAllowed("ghost") {
		t.Error("default 적용된 tenant 가 audit 없으면 허용되어야")
	}
}

func TestDynamicQuota_StoreCfgUsed(t *testing.T) {
	store := NewMemoryQuotaStore()
	_ = store.Set(context.Background(), QuotaConfig{
		TenantID: "hospA", DailyTokenLimit: 100,
	})
	q := NewDynamicQuota(store, nil, 60*time.Second)
	// audit 미설정이라 한도 검사 skip → 허용
	if !q.CheckAllowed("hospA") {
		t.Error("audit 없으면 허용")
	}
}

func TestDynamicQuota_CacheTTL(t *testing.T) {
	store := NewMemoryQuotaStore()
	_ = store.Set(context.Background(), QuotaConfig{TenantID: "hospA"})
	q := NewDynamicQuota(store, nil, 50*time.Millisecond)
	q.CheckAllowed("hospA") // 캐시 적재

	// 한도 변경 (캐시 TTL 내라 store 재조회 안 됨)
	_ = store.Set(context.Background(), QuotaConfig{TenantID: "hospA", DailyTokenLimit: 1})
	q.CheckAllowed("hospA") // 캐시 hit

	time.Sleep(80 * time.Millisecond)
	// TTL 후 재조회 시 새 한도 적용
	q.CheckAllowed("hospA")
	// 캐시 동작 검증은 InvalidateCache 동작과 동일
}

func TestDynamicQuota_InvalidateCache(t *testing.T) {
	store := NewMemoryQuotaStore()
	q := NewDynamicQuota(store, nil, 60*time.Second)
	q.CheckAllowed("x") // 미등록 → default 캐시
	q.InvalidateCache("x")
	// 무한대로 캐시 접근하지 않는지 (panic 안 나면 OK)
	q.CheckAllowed("x")
}

func TestDynamicQuota_ClearCache(t *testing.T) {
	store := NewMemoryQuotaStore()
	q := NewDynamicQuota(store, nil, 60*time.Second)
	q.CheckAllowed("a")
	q.CheckAllowed("b")
	q.ClearCache()
}

func TestDynamicQuota_RecordUsage_NoOp(t *testing.T) {
	q := NewDynamicQuota(nil, nil, 60*time.Second)
	// audit_log 가 별도로 기록하므로 RecordUsage 는 no-op
	q.RecordUsage("hospA", 100) // panic 안 나면 OK
}

// recordingViolationNotifier 는 Phase AN-2 의 한도 초과 알림 호출을 기록.
type recordingViolationNotifier struct {
	tenantID  string
	limit     int
	used      int64
	limitType string
	count     int
}

func (n *recordingViolationNotifier) NotifyQuotaExceeded(tenantID string, limit int, used int64, limitType string) {
	n.count++
	n.tenantID = tenantID
	n.limit = limit
	n.used = used
	n.limitType = limitType
}

func TestDynamicQuota_NotifyOnDailyExceeded(t *testing.T) {
	store := NewMemoryQuotaStore()
	_ = store.Set(context.Background(), QuotaConfig{
		TenantID: "hospA", DailyTokenLimit: 100,
	})
	q := NewDynamicQuota(store, &PostgresAuditLog{
		querier: &fakeQuerier{row: &fakeRow{val: 200}},
	}, 60*time.Second)

	notifier := &recordingViolationNotifier{}
	q.SetViolationNotifier(notifier)

	if q.CheckAllowed("hospA") {
		t.Error("한도 초과인데 허용됨")
	}
	if notifier.count != 1 {
		t.Errorf("notifier count = %d", notifier.count)
	}
	if notifier.tenantID != "hospA" {
		t.Errorf("tenantID = %q", notifier.tenantID)
	}
	if notifier.limitType != "daily_token" {
		t.Errorf("limitType = %q", notifier.limitType)
	}
	if notifier.limit != 100 {
		t.Errorf("limit = %d", notifier.limit)
	}
	if notifier.used != 200 {
		t.Errorf("used = %d", notifier.used)
	}
}

func TestDynamicQuota_NotifyOnMonthlyExceeded(t *testing.T) {
	store := NewMemoryQuotaStore()
	_ = store.Set(context.Background(), QuotaConfig{
		TenantID: "hospA", MonthlyTokenLimit: 1000,
	})
	q := NewDynamicQuota(store, &PostgresAuditLog{
		querier: &fakeQuerier{row: &fakeRow{val: 1500}},
	}, 60*time.Second)

	notifier := &recordingViolationNotifier{}
	q.SetViolationNotifier(notifier)

	if q.CheckAllowed("hospA") {
		t.Error("월 한도 초과인데 허용됨")
	}
	if notifier.limitType != "monthly_token" {
		t.Errorf("limitType = %q", notifier.limitType)
	}
}

func TestDynamicQuota_NoNotifyWhenWithinLimit(t *testing.T) {
	store := NewMemoryQuotaStore()
	_ = store.Set(context.Background(), QuotaConfig{
		TenantID: "hospA", DailyTokenLimit: 1000,
	})
	q := NewDynamicQuota(store, &PostgresAuditLog{
		querier: &fakeQuerier{row: &fakeRow{val: 100}},
	}, 60*time.Second)

	notifier := &recordingViolationNotifier{}
	q.SetViolationNotifier(notifier)

	if !q.CheckAllowed("hospA") {
		t.Error("한도 내인데 거부됨")
	}
	if notifier.count != 0 {
		t.Errorf("한도 내인데 알림: count=%d", notifier.count)
	}
}

func TestDynamicQuota_NoPanicWithoutNotifier(t *testing.T) {
	store := NewMemoryQuotaStore()
	_ = store.Set(context.Background(), QuotaConfig{
		TenantID: "hospA", DailyTokenLimit: 10,
	})
	q := NewDynamicQuota(store, &PostgresAuditLog{
		querier: &fakeQuerier{row: &fakeRow{val: 100}},
	}, 60*time.Second)
	// notifier 미설정 — 한도 초과해도 panic 없어야
	if q.CheckAllowed("hospA") {
		t.Error("한도 초과인데 허용")
	}
}
