package tenancy_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/manpasik/backend/shared/tenancy"
)

func TestInvitationCleaner_MarkExpired_Single(t *testing.T) {
	store := tenancy.NewMemoryInvitationStore()
	now := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)

	// 만료된 pending
	_ = store.Add(tenancy.Invitation{
		Token: "expired-1", TenantID: "t1", Role: tenancy.TenantRoleMember,
		Status: tenancy.InvitationPending,
		IssuedAt: now.Add(-2 * 24 * time.Hour), ExpiresAt: now.Add(-1 * time.Hour),
	})
	// 미만료 pending
	_ = store.Add(tenancy.Invitation{
		Token: "active-1", TenantID: "t1", Role: tenancy.TenantRoleMember,
		Status: tenancy.InvitationPending,
		IssuedAt: now, ExpiresAt: now.Add(7 * 24 * time.Hour),
	})
	// 이미 accepted (변경 안 됨)
	_ = store.Add(tenancy.Invitation{
		Token: "done-1", TenantID: "t1", Role: tenancy.TenantRoleMember,
		Status: tenancy.InvitationAccepted,
		IssuedAt: now.Add(-3 * 24 * time.Hour), ExpiresAt: now.Add(-2 * 24 * time.Hour),
	})

	cleaner, err := tenancy.NewInvitationCleaner(store, nil, tenancy.InvitationCleanerConfig{
		Now: func() time.Time { return now },
	})
	if err != nil {
		t.Fatal(err)
	}

	count := cleaner.MarkExpired("t1")
	if count != 1 {
		t.Errorf("expired = %d, want 1", count)
	}

	// 상태 확인
	expired, _ := store.Get("expired-1")
	if expired.Status != tenancy.InvitationExpired {
		t.Errorf("expired-1 status = %q", expired.Status)
	}
	active, _ := store.Get("active-1")
	if active.Status != tenancy.InvitationPending {
		t.Errorf("active-1 status changed: %q", active.Status)
	}
	done, _ := store.Get("done-1")
	if done.Status != tenancy.InvitationAccepted {
		t.Errorf("done-1 status changed: %q", done.Status)
	}
}

func TestInvitationCleaner_MarkExpired_Idempotent(t *testing.T) {
	store := tenancy.NewMemoryInvitationStore()
	now := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	_ = store.Add(tenancy.Invitation{
		Token: "x", TenantID: "t", Role: tenancy.TenantRoleMember,
		Status: tenancy.InvitationPending,
		ExpiresAt: now.Add(-time.Hour),
	})
	cleaner, _ := tenancy.NewInvitationCleaner(store, nil, tenancy.InvitationCleanerConfig{
		Now: func() time.Time { return now },
	})
	if cleaner.MarkExpired("t") != 1 {
		t.Error("첫 호출 expired != 1")
	}
	if cleaner.MarkExpired("t") != 0 {
		t.Error("두 번째 호출에서 다시 처리됨 (멱등 X)")
	}
}

func TestInvitationCleaner_MarkExpiredAll(t *testing.T) {
	store := tenancy.NewMemoryInvitationStore()
	now := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	for _, tid := range []string{"A", "A", "B"} {
		_ = store.Add(tenancy.Invitation{
			Token: tid + "-" + string(rune('a'+(now.Nanosecond()%26))) + tid,
			TenantID: tenancy.TenantID(tid), Role: tenancy.TenantRoleMember,
			Status: tenancy.InvitationPending,
			ExpiresAt: now.Add(-time.Hour),
		})
		now = now.Add(time.Microsecond) // 토큰 충돌 방지
	}

	tenantList := tenancy.TenantListerFunc(func() []tenancy.TenantID {
		return []tenancy.TenantID{"A", "B"}
	})
	cleaner, _ := tenancy.NewInvitationCleaner(store, tenantList, tenancy.InvitationCleanerConfig{
		Now: func() time.Time { return now },
	})
	count := cleaner.MarkExpiredAll()
	if count != 3 {
		t.Errorf("expired = %d, want 3", count)
	}
	if cleaner.ExpiredCount() != 3 {
		t.Errorf("ExpiredCount = %d", cleaner.ExpiredCount())
	}
	if cleaner.CycleCount() != 1 {
		t.Errorf("CycleCount = %d", cleaner.CycleCount())
	}
}

func TestInvitationCleaner_NilTenantLister(t *testing.T) {
	store := tenancy.NewMemoryInvitationStore()
	cleaner, _ := tenancy.NewInvitationCleaner(store, nil, tenancy.InvitationCleanerConfig{})
	// MarkExpiredAll 은 tenants=nil 시 0 반환 + 에러 없음
	if n := cleaner.MarkExpiredAll(); n != 0 {
		t.Errorf("nil tenants 에서 expired = %d", n)
	}
}

func TestInvitationCleaner_NilStore(t *testing.T) {
	if _, err := tenancy.NewInvitationCleaner(nil, nil, tenancy.InvitationCleanerConfig{}); err == nil {
		t.Error("nil store 통과")
	}
}

func TestInvitationCleaner_StartStop(t *testing.T) {
	store := tenancy.NewMemoryInvitationStore()
	now := time.Now().UTC()
	_ = store.Add(tenancy.Invitation{
		Token: "x", TenantID: "t", Role: tenancy.TenantRoleMember,
		Status: tenancy.InvitationPending,
		ExpiresAt: now.Add(-time.Hour),
	})
	tenants := tenancy.TenantListerFunc(func() []tenancy.TenantID {
		return []tenancy.TenantID{"t"}
	})
	cleaner, _ := tenancy.NewInvitationCleaner(store, tenants, tenancy.InvitationCleanerConfig{
		Interval: 30 * time.Millisecond,
	})

	cleaner.Start(context.Background())
	time.Sleep(80 * time.Millisecond)
	cleaner.Stop()

	if cleaner.CycleCount() < 2 {
		t.Errorf("CycleCount = %d", cleaner.CycleCount())
	}
	got, _ := store.Get("x")
	if got.Status != tenancy.InvitationExpired {
		t.Errorf("status = %q", got.Status)
	}
}

func TestInvitationCleaner_StartIdempotent(t *testing.T) {
	store := tenancy.NewMemoryInvitationStore()
	cleaner, _ := tenancy.NewInvitationCleaner(store, nil, tenancy.InvitationCleanerConfig{
		Interval: 50 * time.Millisecond,
	})
	cleaner.Start(context.Background())
	cleaner.Start(context.Background()) // 두 번째 호출 안전
	cleaner.Stop()
}

func TestInvitationCleaner_StopWithoutStart(t *testing.T) {
	store := tenancy.NewMemoryInvitationStore()
	cleaner, _ := tenancy.NewInvitationCleaner(store, nil, tenancy.InvitationCleanerConfig{})
	cleaner.Stop() // panic 안 되어야
}

func TestInvitationCleaner_OnError(t *testing.T) {
	// store.Update 가 항상 실패하는 가짜 store
	failing := &failingUpdateStore{
		base: tenancy.NewMemoryInvitationStore(),
		err:  errors.New("update failed"),
	}
	now := time.Now().UTC()
	_ = failing.base.Add(tenancy.Invitation{
		Token: "x", TenantID: "t", Role: tenancy.TenantRoleMember,
		Status: tenancy.InvitationPending,
		ExpiresAt: now.Add(-time.Hour),
	})

	var gotErrs int
	cleaner, _ := tenancy.NewInvitationCleaner(failing, nil, tenancy.InvitationCleanerConfig{
		Now:     func() time.Time { return now },
		OnError: func(err error) { gotErrs++ },
	})
	cleaner.MarkExpired("t")
	if gotErrs == 0 {
		t.Error("OnError 미호출")
	}
}

func TestInvitationCleaner_SetTenants(t *testing.T) {
	store := tenancy.NewMemoryInvitationStore()
	cleaner, _ := tenancy.NewInvitationCleaner(store, nil, tenancy.InvitationCleanerConfig{})
	cleaner.SetTenants(tenancy.TenantListerFunc(func() []tenancy.TenantID {
		return []tenancy.TenantID{"new-t"}
	}))
	// SetTenants 후 MarkExpiredAll 이 nil 이 아니어야 — 단순히 panic 안 하는지 검증
	cleaner.MarkExpiredAll()
}

func TestMembershipBackedTenantLister(t *testing.T) {
	store := tenancy.NewMemoryMembershipStore()
	_ = store.Add(tenancy.Membership{UserID: "u1", TenantID: "A", Role: tenancy.TenantRoleAdmin})
	_ = store.Add(tenancy.Membership{UserID: "u1", TenantID: "B", Role: tenancy.TenantRoleMember})
	_ = store.Add(tenancy.Membership{UserID: "u2", TenantID: "A", Role: tenancy.TenantRoleMember})

	lister := tenancy.NewMembershipBackedTenantLister(store, func() []string {
		return []string{"u1", "u2"}
	})
	tenants := lister.AllTenants()
	if len(tenants) != 2 {
		t.Errorf("tenants = %v", tenants)
	}
}

func TestMembershipBackedTenantLister_NilUsers(t *testing.T) {
	store := tenancy.NewMemoryMembershipStore()
	lister := tenancy.NewMembershipBackedTenantLister(store, nil)
	if got := lister.AllTenants(); got != nil {
		t.Errorf("nil allUserIDs 에서 = %v", got)
	}
}

// failingUpdateStore 는 Update 가 항상 실패하는 가짜 store.
type failingUpdateStore struct {
	base *tenancy.MemoryInvitationStore
	err  error
}

func (s *failingUpdateStore) Add(inv tenancy.Invitation) error {
	return s.base.Add(inv)
}
func (s *failingUpdateStore) Get(token string) (*tenancy.Invitation, error) {
	return s.base.Get(token)
}
func (s *failingUpdateStore) Update(_ tenancy.Invitation) error {
	return s.err
}
func (s *failingUpdateStore) ListByTenant(tid tenancy.TenantID) []*tenancy.Invitation {
	return s.base.ListByTenant(tid)
}
