package tenancy_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/manpasik/backend/shared/tenancy"
)

// fixedTokens 는 결정적 토큰 생성기 (테스트용).
func fixedTokens(tokens ...string) func() (string, error) {
	idx := 0
	return func() (string, error) {
		if idx >= len(tokens) {
			return "", errors.New("토큰 소진")
		}
		t := tokens[idx]
		idx++
		return t, nil
	}
}

func newInviteSvc(t *testing.T) (*tenancy.InvitationService,
	*tenancy.MemoryInvitationStore, *tenancy.MemoryMembershipStore) {
	t.Helper()
	invStore := tenancy.NewMemoryInvitationStore()
	memStore := tenancy.NewMemoryMembershipStore()
	svc, err := tenancy.NewInvitationService(invStore, memStore, tenancy.InvitationServiceConfig{
		TokenGenerator: fixedTokens("tok-1", "tok-2", "tok-3"),
	})
	if err != nil {
		t.Fatal(err)
	}
	return svc, invStore, memStore
}

func TestInvitation_IsExpired(t *testing.T) {
	inv := tenancy.Invitation{ExpiresAt: time.Now().Add(-time.Hour)}
	if !inv.IsExpired() {
		t.Error("과거 만료 미감지")
	}
	inv2 := tenancy.Invitation{ExpiresAt: time.Now().Add(time.Hour)}
	if inv2.IsExpired() {
		t.Error("미래 만료를 만료로 판정")
	}
}

func TestInvitation_IsActive(t *testing.T) {
	inv := tenancy.Invitation{
		Status:    tenancy.InvitationPending,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	if !inv.IsActive() {
		t.Error("pending + 미만료 = active 여야")
	}

	inv.Status = tenancy.InvitationAccepted
	if inv.IsActive() {
		t.Error("accepted 인데 active 로 판정")
	}
}

func TestInvitationService_Invite_Basic(t *testing.T) {
	svc, _, _ := newInviteSvc(t)
	inv, err := svc.Invite(tenancy.InviteRequest{
		InviterID: "doc1",
		TenantID:  "hospA",
		Role:      tenancy.TenantRoleMember,
	})
	if err != nil {
		t.Fatal(err)
	}
	if inv.Token == "" || inv.Status != tenancy.InvitationPending {
		t.Errorf("inv = %+v", inv)
	}
}

func TestInvitationService_Invite_UnknownRole(t *testing.T) {
	svc, _, _ := newInviteSvc(t)
	if _, err := svc.Invite(tenancy.InviteRequest{
		InviterID: "u", TenantID: "t", Role: "ghost",
	}); err == nil {
		t.Error("unknown role 통과")
	}
}

func TestInvitationService_Invite_MissingFields(t *testing.T) {
	svc, _, _ := newInviteSvc(t)
	if _, err := svc.Invite(tenancy.InviteRequest{TenantID: "t", Role: tenancy.TenantRoleMember}); err == nil {
		t.Error("InviterID 누락 통과")
	}
}

func TestInvitationService_Invite_PolicyDeny(t *testing.T) {
	invStore := tenancy.NewMemoryInvitationStore()
	memStore := tenancy.NewMemoryMembershipStore()
	// inviter 가 t1 의 멤버가 아님 → admin 권한 없음
	policy := tenancy.NewPolicyEngine(memStore)
	svc, _ := tenancy.NewInvitationService(invStore, memStore, tenancy.InvitationServiceConfig{})
	svc.SetPolicyEngine(policy)

	if _, err := svc.Invite(tenancy.InviteRequest{
		InviterID: "outsider", TenantID: "t1", Role: tenancy.TenantRoleMember,
	}); err == nil {
		t.Error("권한 없는 inviter 통과")
	}
}

func TestInvitationService_Accept_Success(t *testing.T) {
	svc, _, memStore := newInviteSvc(t)
	inv, _ := svc.Invite(tenancy.InviteRequest{
		InviterID: "doc1", TenantID: "hospA", Role: tenancy.TenantRoleMedicalStaff,
	})

	m, err := svc.Accept(inv.Token, "newDoc")
	if err != nil {
		t.Fatal(err)
	}
	if m.UserID != "newDoc" || m.Role != tenancy.TenantRoleMedicalStaff {
		t.Errorf("m = %+v", m)
	}

	// 멤버십 영속 확인
	got, err := memStore.Get("newDoc", "hospA")
	if err != nil {
		t.Fatal(err)
	}
	if got.Role != tenancy.TenantRoleMedicalStaff {
		t.Errorf("Role = %q", got.Role)
	}
}

func TestInvitationService_Accept_AlreadyConsumed(t *testing.T) {
	svc, _, _ := newInviteSvc(t)
	inv, _ := svc.Invite(tenancy.InviteRequest{
		InviterID: "doc1", TenantID: "hospA", Role: tenancy.TenantRoleMember,
	})
	_, _ = svc.Accept(inv.Token, "user1")
	// 동일 토큰 재수락 시도
	if _, err := svc.Accept(inv.Token, "user2"); err != tenancy.ErrInvitationConsumed {
		t.Errorf("err = %v", err)
	}
}

func TestInvitationService_Accept_Expired(t *testing.T) {
	invStore := tenancy.NewMemoryInvitationStore()
	memStore := tenancy.NewMemoryMembershipStore()
	now := time.Now().UTC()
	svc, _ := tenancy.NewInvitationService(invStore, memStore, tenancy.InvitationServiceConfig{
		TokenGenerator: fixedTokens("tok-1"),
		Now:            func() time.Time { return now },
	})

	inv, _ := svc.Invite(tenancy.InviteRequest{
		InviterID: "u", TenantID: "t", Role: tenancy.TenantRoleMember,
		TTL: 1 * time.Hour,
	})
	// 시간이 만료 후로 진행했다고 가정 (직접 inv 만료 시각 검증)
	inv.ExpiresAt = now.Add(-time.Hour)
	_ = invStore.Update(*inv)

	if _, err := svc.Accept(inv.Token, "user1"); err != tenancy.ErrInvitationExpired {
		t.Errorf("err = %v", err)
	}
}

func TestInvitationService_Accept_NotFound(t *testing.T) {
	svc, _, _ := newInviteSvc(t)
	if _, err := svc.Accept("ghost", "u"); err != tenancy.ErrInvitationNotFound {
		t.Errorf("err = %v", err)
	}
}

func TestInvitationService_Accept_AlreadyMember(t *testing.T) {
	svc, _, memStore := newInviteSvc(t)
	_ = memStore.Add(tenancy.Membership{
		UserID: "u1", TenantID: "t1", Role: tenancy.TenantRoleAdmin,
	})

	inv, _ := svc.Invite(tenancy.InviteRequest{
		InviterID: "doc1", TenantID: "t1", Role: tenancy.TenantRoleMember,
	})
	if _, err := svc.Accept(inv.Token, "u1"); err == nil {
		t.Error("기존 멤버 재가입 허용됨")
	}
}

func TestInvitationService_Revoke_ByInviter(t *testing.T) {
	svc, _, _ := newInviteSvc(t)
	inv, _ := svc.Invite(tenancy.InviteRequest{
		InviterID: "doc1", TenantID: "t", Role: tenancy.TenantRoleMember,
	})
	if err := svc.Revoke(inv.Token, "doc1"); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Accept(inv.Token, "u"); err != tenancy.ErrInvitationConsumed {
		t.Errorf("취소 후 수락 err = %v", err)
	}
}

func TestInvitationService_Revoke_AlreadyConsumed(t *testing.T) {
	svc, _, _ := newInviteSvc(t)
	inv, _ := svc.Invite(tenancy.InviteRequest{
		InviterID: "doc1", TenantID: "t", Role: tenancy.TenantRoleMember,
	})
	_ = svc.Revoke(inv.Token, "doc1")
	if err := svc.Revoke(inv.Token, "doc1"); err != tenancy.ErrInvitationConsumed {
		t.Errorf("재취소 err = %v", err)
	}
}

func TestInvitationService_ListPendingByTenant(t *testing.T) {
	svc, _, _ := newInviteSvc(t)
	_, _ = svc.Invite(tenancy.InviteRequest{
		InviterID: "doc1", TenantID: "t1", Role: tenancy.TenantRoleMember,
	})
	inv2, _ := svc.Invite(tenancy.InviteRequest{
		InviterID: "doc1", TenantID: "t1", Role: tenancy.TenantRoleViewer,
	})
	_ = svc.Revoke(inv2.Token, "doc1")
	_, _ = svc.Invite(tenancy.InviteRequest{
		InviterID: "doc1", TenantID: "t2", Role: tenancy.TenantRoleAdmin,
	})

	pending := svc.ListPendingByTenant("t1")
	if len(pending) != 1 {
		t.Errorf("pending = %d, want 1", len(pending))
	}
}

func TestInvitationService_DefaultTTL(t *testing.T) {
	svc, _, _ := newInviteSvc(t)
	inv, _ := svc.Invite(tenancy.InviteRequest{
		InviterID: "u", TenantID: "t", Role: tenancy.TenantRoleMember,
	})
	expectedTTL := 7 * 24 * time.Hour
	gotTTL := inv.ExpiresAt.Sub(inv.IssuedAt)
	if gotTTL < expectedTTL-time.Minute || gotTTL > expectedTTL+time.Minute {
		t.Errorf("TTL = %v, want ~%v", gotTTL, expectedTTL)
	}
}

func TestNewInvitationService_Validation(t *testing.T) {
	if _, err := tenancy.NewInvitationService(nil, nil, tenancy.InvitationServiceConfig{}); err == nil {
		t.Error("nil store 통과")
	}
}

func TestMemoryInvitationStore_DuplicateToken(t *testing.T) {
	s := tenancy.NewMemoryInvitationStore()
	inv := tenancy.Invitation{Token: "dup", TenantID: "t", Role: tenancy.TenantRoleMember}
	if err := s.Add(inv); err != nil {
		t.Fatal(err)
	}
	if err := s.Add(inv); err == nil {
		t.Error("토큰 중복 통과")
	}
}

func TestMemoryInvitationStore_GetNotFound(t *testing.T) {
	s := tenancy.NewMemoryInvitationStore()
	if _, err := s.Get("missing"); err != tenancy.ErrInvitationNotFound {
		t.Errorf("err = %v", err)
	}
}

func TestMemoryInvitationStore_ListByTenant(t *testing.T) {
	s := tenancy.NewMemoryInvitationStore()
	_ = s.Add(tenancy.Invitation{Token: "t1", TenantID: "A", Role: tenancy.TenantRoleMember})
	_ = s.Add(tenancy.Invitation{Token: "t2", TenantID: "A", Role: tenancy.TenantRoleViewer})
	_ = s.Add(tenancy.Invitation{Token: "t3", TenantID: "B", Role: tenancy.TenantRoleMember})

	list := s.ListByTenant("A")
	if len(list) != 2 {
		t.Errorf("len = %d", len(list))
	}
}

// recordingNotifier 는 NotifyInvitation 호출을 기록.
type recordingNotifier struct {
	mu     sync.Mutex
	called int32
	last   tenancy.Invitation
	failBy error
}

func (n *recordingNotifier) NotifyInvitation(_ context.Context, inv tenancy.Invitation) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	atomic.AddInt32(&n.called, 1)
	n.last = inv
	return n.failBy
}

func TestInvitationService_NotifierCalledOnInvite(t *testing.T) {
	svc, _, _ := newInviteSvc(t)
	notifier := &recordingNotifier{}
	svc.SetNotifier(notifier)

	_, err := svc.Invite(tenancy.InviteRequest{
		InviterID: "doc1", TenantID: "hospA", Role: tenancy.TenantRoleMember,
		InviteeHint: "user@example.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	// 비동기 호출 — 잠시 대기
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if atomic.LoadInt32(&notifier.called) > 0 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if atomic.LoadInt32(&notifier.called) != 1 {
		t.Errorf("notifier 호출 = %d", notifier.called)
	}
	notifier.mu.Lock()
	if notifier.last.InviteeHint != "user@example.com" {
		t.Errorf("hint 전달 안됨: %q", notifier.last.InviteeHint)
	}
	notifier.mu.Unlock()
}

func TestInvitationService_NotifierFailureNonFatal(t *testing.T) {
	svc, _, _ := newInviteSvc(t)
	notifier := &recordingNotifier{failBy: errors.New("smtp down")}
	svc.SetNotifier(notifier)

	// 알림 실패해도 invite 발급 자체는 성공해야
	inv, err := svc.Invite(tenancy.InviteRequest{
		InviterID: "doc1", TenantID: "hospA", Role: tenancy.TenantRoleMember,
	})
	if err != nil {
		t.Fatalf("발급 실패: %v", err)
	}
	if inv.Token == "" {
		t.Error("발급 token 없음")
	}
}

func TestInvitationService_NoNotifier(t *testing.T) {
	// SetNotifier 미호출 시에도 정상 동작
	svc, _, _ := newInviteSvc(t)
	if _, err := svc.Invite(tenancy.InviteRequest{
		InviterID: "doc1", TenantID: "hospA", Role: tenancy.TenantRoleMember,
	}); err != nil {
		t.Fatal(err)
	}
}

// recordingWebhookServer 는 webhook 호출을 모두 기록하는 httptest 서버.
type recordingWebhookServer struct {
	mu       sync.Mutex
	bodies   [][]byte
	server   *httptest.Server
}

func newRecordingWebhookServer(t *testing.T) *recordingWebhookServer {
	t.Helper()
	r := &recordingWebhookServer{}
	r.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		body, _ := readBody(req)
		r.mu.Lock()
		r.bodies = append(r.bodies, body)
		r.mu.Unlock()
		w.WriteHeader(200)
	}))
	t.Cleanup(r.server.Close)
	return r
}

func (r *recordingWebhookServer) URL() string { return r.server.URL }
func (r *recordingWebhookServer) Count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.bodies)
}
func (r *recordingWebhookServer) Bodies() [][]byte {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([][]byte, len(r.bodies))
	copy(out, r.bodies)
	return out
}
func (r *recordingWebhookServer) WaitForCount(count int, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if r.Count() >= count {
			return true
		}
		time.Sleep(10 * time.Millisecond)
	}
	return false
}

func TestInvitationService_WebhookOnInvite(t *testing.T) {
	srv := newRecordingWebhookServer(t)
	dispatcher, _ := tenancy.NewWebhookDispatcher(tenancy.WebhookConfig{
		URL: srv.URL(), MaxRetries: 0,
	}, nil)

	svc, _, _ := newInviteSvc(t)
	svc.SetWebhookDispatcher(dispatcher)

	_, err := svc.Invite(tenancy.InviteRequest{
		InviterID: "doc1", TenantID: "hospA", Role: tenancy.TenantRoleMember,
	})
	if err != nil {
		t.Fatal(err)
	}

	if !srv.WaitForCount(1, 2*time.Second) {
		t.Errorf("webhook 미발송: count=%d", srv.Count())
	}
}

func TestInvitationService_WebhookOnAccept(t *testing.T) {
	srv := newRecordingWebhookServer(t)
	dispatcher, _ := tenancy.NewWebhookDispatcher(tenancy.WebhookConfig{
		URL: srv.URL(), MaxRetries: 0,
	}, nil)

	svc, _, _ := newInviteSvc(t)
	svc.SetWebhookDispatcher(dispatcher)

	inv, _ := svc.Invite(tenancy.InviteRequest{
		InviterID: "doc1", TenantID: "hospA", Role: tenancy.TenantRoleMember,
	})
	// invite 발송 + accept 시 invitation.accepted + membership.created = 3개
	_, err := svc.Accept(inv.Token, "newUser")
	if err != nil {
		t.Fatal(err)
	}

	if !srv.WaitForCount(3, 2*time.Second) {
		t.Errorf("webhook count = %d, want >= 3", srv.Count())
	}
}

func TestInvitationService_WebhookOnRevoke(t *testing.T) {
	srv := newRecordingWebhookServer(t)
	dispatcher, _ := tenancy.NewWebhookDispatcher(tenancy.WebhookConfig{
		URL: srv.URL(), MaxRetries: 0,
	}, nil)

	svc, _, _ := newInviteSvc(t)
	svc.SetWebhookDispatcher(dispatcher)

	inv, _ := svc.Invite(tenancy.InviteRequest{
		InviterID: "doc1", TenantID: "hospA", Role: tenancy.TenantRoleMember,
	})
	if err := svc.Revoke(inv.Token, "doc1"); err != nil {
		t.Fatal(err)
	}

	// invite + revoke = 2개
	if !srv.WaitForCount(2, 2*time.Second) {
		t.Errorf("webhook count = %d, want >= 2", srv.Count())
	}
}

func TestInvitationService_NoWebhookNoOp(t *testing.T) {
	// SetWebhookDispatcher 미호출 시에도 정상 동작
	svc, _, _ := newInviteSvc(t)
	if _, err := svc.Invite(tenancy.InviteRequest{
		InviterID: "doc1", TenantID: "hospA", Role: tenancy.TenantRoleMember,
	}); err != nil {
		t.Fatal(err)
	}
}
