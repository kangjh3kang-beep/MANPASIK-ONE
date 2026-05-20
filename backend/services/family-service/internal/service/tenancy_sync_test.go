package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/manpasik/backend/shared/tenancy"
)

func TestTenancySyncAdapter_OnGroupCreated(t *testing.T) {
	store := tenancy.NewMemoryMembershipStore()
	a := NewTenancySyncAdapter(store)
	if err := a.OnGroupCreated(context.Background(), "group-1", "user-1"); err != nil {
		t.Fatal(err)
	}
	got, err := store.Get("user-1", "group-1")
	if err != nil {
		t.Fatal(err)
	}
	if got.Role != tenancy.TenantRoleOwner {
		t.Errorf("Role = %q, want owner", got.Role)
	}
	if !got.Active {
		t.Error("멤버가 비활성")
	}
}

func TestTenancySyncAdapter_OnMemberAdded_Mapping(t *testing.T) {
	cases := []struct {
		familyRole string
		want       tenancy.TenantRole
	}{
		{"owner", tenancy.TenantRoleOwner},
		{"guardian", tenancy.TenantRoleAdmin},
		{"member", tenancy.TenantRoleMember},
		{"child", tenancy.TenantRoleMember},
		{"elderly", tenancy.TenantRoleMember},
		{"unknown", tenancy.TenantRoleViewer},
		{"", tenancy.TenantRoleViewer},
	}
	for _, tc := range cases {
		t.Run(tc.familyRole, func(t *testing.T) {
			store := tenancy.NewMemoryMembershipStore()
			a := NewTenancySyncAdapter(store)
			if err := a.OnMemberAdded(context.Background(), "g", "u-"+tc.familyRole, tc.familyRole); err != nil {
				t.Fatal(err)
			}
			got, err := store.Get("u-"+tc.familyRole, "g")
			if err != nil {
				t.Fatal(err)
			}
			if got.Role != tc.want {
				t.Errorf("got %q, want %q", got.Role, tc.want)
			}
		})
	}
}

func TestTenancySyncAdapter_OnMemberRemoved(t *testing.T) {
	store := tenancy.NewMemoryMembershipStore()
	_ = store.Add(tenancy.Membership{UserID: "u1", TenantID: "g", Role: tenancy.TenantRoleMember})
	a := NewTenancySyncAdapter(store)
	if err := a.OnMemberRemoved(context.Background(), "g", "u1"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Get("u1", "g"); err != tenancy.ErrNoMembership {
		t.Errorf("err = %v, want ErrNoMembership", err)
	}
}

func TestTenancySyncAdapter_OnMemberRoleChanged(t *testing.T) {
	store := tenancy.NewMemoryMembershipStore()
	_ = store.Add(tenancy.Membership{UserID: "u1", TenantID: "g", Role: tenancy.TenantRoleMember})
	a := NewTenancySyncAdapter(store)
	if err := a.OnMemberRoleChanged(context.Background(), "g", "u1", "guardian"); err != nil {
		t.Fatal(err)
	}
	got, _ := store.Get("u1", "g")
	if got.Role != tenancy.TenantRoleAdmin {
		t.Errorf("Role = %q, want admin", got.Role)
	}
}

func TestTenancySyncAdapter_NilStore(t *testing.T) {
	a := NewTenancySyncAdapter(nil)
	// 모든 메서드가 nil 에러로 안전 통과
	if err := a.OnGroupCreated(context.Background(), "g", "u"); err != nil {
		t.Errorf("OnGroupCreated err = %v", err)
	}
	if err := a.OnMemberAdded(context.Background(), "g", "u", "member"); err != nil {
		t.Errorf("OnMemberAdded err = %v", err)
	}
	if err := a.OnMemberRemoved(context.Background(), "g", "u"); err != nil {
		t.Errorf("OnMemberRemoved err = %v", err)
	}
	if err := a.OnMemberRoleChanged(context.Background(), "g", "u", "owner"); err != nil {
		t.Errorf("OnMemberRoleChanged err = %v", err)
	}
}

func TestFamilyRole_String(t *testing.T) {
	cases := map[FamilyRole]string{
		RoleOwner:    "owner",
		RoleGuardian: "guardian",
		RoleMember:   "member",
		RoleChild:    "child",
		RoleElderly:  "elderly",
		RoleUnknown:  "unknown",
	}
	for r, want := range cases {
		if got := r.String(); got != want {
			t.Errorf("FamilyRole(%d) = %q, want %q", r, got, want)
		}
	}
}

// recordingHook 는 webhook 호출 수신을 기록.
type recordingHook struct {
	mu    sync.Mutex
	count int32
}

func newRecordingHook(t *testing.T) (*recordingHook, *httptest.Server) {
	t.Helper()
	r := &recordingHook{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		r.mu.Lock()
		atomic.AddInt32(&r.count, 1)
		r.mu.Unlock()
		w.WriteHeader(200)
	}))
	t.Cleanup(srv.Close)
	return r, srv
}

func (h *recordingHook) waitFor(target int32, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if atomic.LoadInt32(&h.count) >= target {
			return true
		}
		time.Sleep(10 * time.Millisecond)
	}
	return false
}

func TestTenancySyncAdapter_WebhookOnGroupCreated(t *testing.T) {
	hook, srv := newRecordingHook(t)
	dispatcher, _ := tenancy.NewWebhookDispatcher(tenancy.WebhookConfig{
		URL: srv.URL, MaxRetries: 0,
	}, nil)

	store := tenancy.NewMemoryMembershipStore()
	a := NewTenancySyncAdapter(store)
	a.SetWebhookDispatcher(dispatcher)

	if err := a.OnGroupCreated(context.Background(), "fam-A", "u-owner"); err != nil {
		t.Fatal(err)
	}
	if !hook.waitFor(1, 2*time.Second) {
		t.Errorf("webhook 미발송: count=%d", atomic.LoadInt32(&hook.count))
	}
}

func TestTenancySyncAdapter_WebhookOnMemberAdded(t *testing.T) {
	hook, srv := newRecordingHook(t)
	dispatcher, _ := tenancy.NewWebhookDispatcher(tenancy.WebhookConfig{
		URL: srv.URL, MaxRetries: 0,
	}, nil)

	store := tenancy.NewMemoryMembershipStore()
	a := NewTenancySyncAdapter(store)
	a.SetWebhookDispatcher(dispatcher)

	if err := a.OnMemberAdded(context.Background(), "fam-A", "child", "child"); err != nil {
		t.Fatal(err)
	}
	if !hook.waitFor(1, 2*time.Second) {
		t.Errorf("webhook count=%d", atomic.LoadInt32(&hook.count))
	}
}

func TestTenancySyncAdapter_WebhookOnMemberRemoved(t *testing.T) {
	hook, srv := newRecordingHook(t)
	dispatcher, _ := tenancy.NewWebhookDispatcher(tenancy.WebhookConfig{
		URL: srv.URL, MaxRetries: 0,
	}, nil)

	store := tenancy.NewMemoryMembershipStore()
	_ = store.Add(tenancy.Membership{
		UserID: "u1", TenantID: "fam-A", Role: tenancy.TenantRoleMember,
	})
	a := NewTenancySyncAdapter(store)
	a.SetWebhookDispatcher(dispatcher)

	if err := a.OnMemberRemoved(context.Background(), "fam-A", "u1"); err != nil {
		t.Fatal(err)
	}
	if !hook.waitFor(1, 2*time.Second) {
		t.Errorf("webhook count=%d", atomic.LoadInt32(&hook.count))
	}
}

func TestTenancySyncAdapter_WebhookOnRoleChanged(t *testing.T) {
	hook, srv := newRecordingHook(t)
	dispatcher, _ := tenancy.NewWebhookDispatcher(tenancy.WebhookConfig{
		URL: srv.URL, MaxRetries: 0,
	}, nil)

	store := tenancy.NewMemoryMembershipStore()
	_ = store.Add(tenancy.Membership{
		UserID: "u1", TenantID: "fam-A", Role: tenancy.TenantRoleMember,
	})
	a := NewTenancySyncAdapter(store)
	a.SetWebhookDispatcher(dispatcher)

	if err := a.OnMemberRoleChanged(context.Background(), "fam-A", "u1", "guardian"); err != nil {
		t.Fatal(err)
	}
	if !hook.waitFor(1, 2*time.Second) {
		t.Errorf("webhook count=%d", atomic.LoadInt32(&hook.count))
	}
}

func TestTenancySyncAdapter_NoWebhookNoOp(t *testing.T) {
	// SetWebhookDispatcher 미호출 시에도 정상 동작
	store := tenancy.NewMemoryMembershipStore()
	a := NewTenancySyncAdapter(store)
	if err := a.OnGroupCreated(context.Background(), "g", "u"); err != nil {
		t.Fatal(err)
	}
}
