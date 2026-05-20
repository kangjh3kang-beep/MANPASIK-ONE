package tenancy_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/manpasik/backend/shared/tenancy"
)

func newHTTPSetup(t *testing.T) (*httptest.Server, *tenancy.MemoryMembershipStore,
	*tenancy.InvitationService) {
	t.Helper()
	invStore := tenancy.NewMemoryInvitationStore()
	memStore := tenancy.NewMemoryMembershipStore()
	svc, _ := tenancy.NewInvitationService(invStore, memStore, tenancy.InvitationServiceConfig{
		TokenGenerator: fixedTokens("invite-1", "invite-2", "invite-3"),
	})
	policy := tenancy.NewPolicyEngine(memStore)
	handler := tenancy.NewHTTPHandler(svc, memStore, policy)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv, memStore, svc
}

// withUser 는 ctx 에 사용자 ID 를 넣는 미들웨어 (테스트용 인증 우회).
func withUser(userID string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := tenancy.WithUser(r.Context(), userID)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// callAs 는 사용자 컨텍스트를 갖는 새 서버에서 요청을 수행하고 응답 반환.
func callAs(t *testing.T, userID, method, path string, body string) *http.Response {
	t.Helper()
	r := httptest.NewRequest(method, path, strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	ctx := tenancy.WithUser(context.Background(), userID)
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	// 서버 mux 를 새로 만들어 직접 서빙
	invStore := tenancy.NewMemoryInvitationStore()
	memStore := tenancy.NewMemoryMembershipStore()
	svc, _ := tenancy.NewInvitationService(invStore, memStore, tenancy.InvitationServiceConfig{})
	handler := tenancy.NewHTTPHandler(svc, memStore, tenancy.NewPolicyEngine(memStore))
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	mux.ServeHTTP(w, r)
	return w.Result()
}

func TestHTTPHandler_Invite_Unauthenticated(t *testing.T) {
	memStore := tenancy.NewMemoryMembershipStore()
	svc, _ := tenancy.NewInvitationService(tenancy.NewMemoryInvitationStore(), memStore, tenancy.InvitationServiceConfig{})
	handler := tenancy.NewHTTPHandler(svc, memStore, nil)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	r := httptest.NewRequest("POST", "/tenancy/invitations",
		strings.NewReader(`{"tenant_id":"t","role":"member"}`))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Code = %d", w.Code)
	}
}

func TestHTTPHandler_Invite_BadJSON(t *testing.T) {
	memStore := tenancy.NewMemoryMembershipStore()
	_ = memStore.Add(tenancy.Membership{UserID: "u1", TenantID: "t", Role: tenancy.TenantRoleAdmin})
	svc, _ := tenancy.NewInvitationService(tenancy.NewMemoryInvitationStore(), memStore, tenancy.InvitationServiceConfig{})
	handler := tenancy.NewHTTPHandler(svc, memStore, tenancy.NewPolicyEngine(memStore))
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	r := httptest.NewRequest("POST", "/tenancy/invitations", strings.NewReader("invalid json"))
	r = r.WithContext(tenancy.WithUser(context.Background(), "u1"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Code = %d", w.Code)
	}
}

func TestHTTPHandler_Invite_Created(t *testing.T) {
	memStore := tenancy.NewMemoryMembershipStore()
	_ = memStore.Add(tenancy.Membership{UserID: "doc1", TenantID: "hospA", Role: tenancy.TenantRoleAdmin})

	invStore := tenancy.NewMemoryInvitationStore()
	svc, _ := tenancy.NewInvitationService(invStore, memStore, tenancy.InvitationServiceConfig{
		TokenGenerator: fixedTokens("test-token-x"),
	})
	policy := tenancy.NewPolicyEngine(memStore)
	svc.SetPolicyEngine(policy)
	handler := tenancy.NewHTTPHandler(svc, memStore, policy)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	r := httptest.NewRequest("POST", "/tenancy/invitations",
		strings.NewReader(`{"tenant_id":"hospA","role":"medical_staff","invitee_hint":"new@hosp.kr"}`))
	r.Header.Set("Content-Type", "application/json")
	r = r.WithContext(tenancy.WithUser(context.Background(), "doc1"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("Code = %d, body = %s", w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["token"] != "test-token-x" {
		t.Errorf("token = %v", resp["token"])
	}
	if resp["status"] != "pending" {
		t.Errorf("status = %v", resp["status"])
	}
}

func TestHTTPHandler_Invite_PolicyDenied(t *testing.T) {
	memStore := tenancy.NewMemoryMembershipStore()
	// outsider 는 t 의 멤버 아님
	invStore := tenancy.NewMemoryInvitationStore()
	svc, _ := tenancy.NewInvitationService(invStore, memStore, tenancy.InvitationServiceConfig{})
	policy := tenancy.NewPolicyEngine(memStore)
	svc.SetPolicyEngine(policy)
	handler := tenancy.NewHTTPHandler(svc, memStore, policy)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	r := httptest.NewRequest("POST", "/tenancy/invitations",
		strings.NewReader(`{"tenant_id":"hospA","role":"member"}`))
	r = r.WithContext(tenancy.WithUser(context.Background(), "outsider"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusForbidden {
		t.Errorf("Code = %d, body = %s", w.Code, w.Body.String())
	}
}

func TestHTTPHandler_Accept_Success(t *testing.T) {
	srv, memStore, svc := newHTTPSetup(t)
	defer srv.Close()
	// admin 멤버십 부여 (invite 권한)
	_ = memStore.Add(tenancy.Membership{UserID: "adm", TenantID: "t1", Role: tenancy.TenantRoleAdmin})
	policy := tenancy.NewPolicyEngine(memStore)
	svc.SetPolicyEngine(policy)

	// 1. invite (직접 service 호출 — server 호출은 callAs 가 새 인스턴스를 만들기 때문)
	inv, err := svc.Invite(tenancy.InviteRequest{
		InviterID: "adm", TenantID: "t1", Role: tenancy.TenantRoleViewer,
	})
	if err != nil {
		t.Fatal(err)
	}

	// 2. 단일 mux 인스턴스에 accept 호출
	handler := tenancy.NewHTTPHandler(svc, memStore, policy)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	r := httptest.NewRequest("POST", "/tenancy/invitations/accept",
		strings.NewReader(`{"token":"`+inv.Token+`"}`))
	r = r.WithContext(tenancy.WithUser(context.Background(), "newUser"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Code = %d, body = %s", w.Code, w.Body.String())
	}
	if got, _ := memStore.Get("newUser", "t1"); got == nil {
		t.Error("멤버십 미생성")
	}
}

func TestHTTPHandler_Accept_NotFound(t *testing.T) {
	memStore := tenancy.NewMemoryMembershipStore()
	svc, _ := tenancy.NewInvitationService(tenancy.NewMemoryInvitationStore(), memStore, tenancy.InvitationServiceConfig{})
	handler := tenancy.NewHTTPHandler(svc, memStore, nil)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	r := httptest.NewRequest("POST", "/tenancy/invitations/accept",
		strings.NewReader(`{"token":"ghost"}`))
	r = r.WithContext(tenancy.WithUser(context.Background(), "u"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusNotFound {
		t.Errorf("Code = %d", w.Code)
	}
}

func TestHTTPHandler_Accept_Conflict_AlreadyMember(t *testing.T) {
	memStore := tenancy.NewMemoryMembershipStore()
	_ = memStore.Add(tenancy.Membership{UserID: "u1", TenantID: "t1", Role: tenancy.TenantRoleAdmin})
	_ = memStore.Add(tenancy.Membership{UserID: "u2", TenantID: "t1", Role: tenancy.TenantRoleMember})
	svc, _ := tenancy.NewInvitationService(tenancy.NewMemoryInvitationStore(), memStore, tenancy.InvitationServiceConfig{
		TokenGenerator: fixedTokens("k"),
	})
	policy := tenancy.NewPolicyEngine(memStore)
	svc.SetPolicyEngine(policy)
	inv, _ := svc.Invite(tenancy.InviteRequest{
		InviterID: "u1", TenantID: "t1", Role: tenancy.TenantRoleViewer,
	})

	handler := tenancy.NewHTTPHandler(svc, memStore, policy)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	r := httptest.NewRequest("POST", "/tenancy/invitations/accept",
		strings.NewReader(`{"token":"`+inv.Token+`"}`))
	r = r.WithContext(tenancy.WithUser(context.Background(), "u2"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusConflict {
		t.Errorf("Code = %d, body = %s", w.Code, w.Body.String())
	}
}

func TestHTTPHandler_MyMemberships(t *testing.T) {
	memStore := tenancy.NewMemoryMembershipStore()
	_ = memStore.Add(tenancy.Membership{UserID: "u1", TenantID: "t1", Role: tenancy.TenantRoleAdmin})
	_ = memStore.Add(tenancy.Membership{UserID: "u1", TenantID: "t2", Role: tenancy.TenantRoleViewer})
	_ = memStore.Add(tenancy.Membership{UserID: "u2", TenantID: "t1", Role: tenancy.TenantRoleMember})

	svc, _ := tenancy.NewInvitationService(tenancy.NewMemoryInvitationStore(), memStore, tenancy.InvitationServiceConfig{})
	handler := tenancy.NewHTTPHandler(svc, memStore, nil)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	r := httptest.NewRequest("GET", "/tenancy/me/memberships", nil)
	r = r.WithContext(tenancy.WithUser(context.Background(), "u1"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("Code = %d", w.Code)
	}
	var resp struct {
		Memberships []map[string]interface{} `json:"memberships"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if len(resp.Memberships) != 2 {
		t.Errorf("memberships = %d", len(resp.Memberships))
	}
}

func TestHTTPHandler_Revoke_ByInviter(t *testing.T) {
	memStore := tenancy.NewMemoryMembershipStore()
	_ = memStore.Add(tenancy.Membership{UserID: "doc1", TenantID: "t", Role: tenancy.TenantRoleAdmin})
	invStore := tenancy.NewMemoryInvitationStore()
	svc, _ := tenancy.NewInvitationService(invStore, memStore, tenancy.InvitationServiceConfig{
		TokenGenerator: fixedTokens("rev-1"),
	})
	policy := tenancy.NewPolicyEngine(memStore)
	svc.SetPolicyEngine(policy)
	inv, _ := svc.Invite(tenancy.InviteRequest{
		InviterID: "doc1", TenantID: "t", Role: tenancy.TenantRoleMember,
	})

	handler := tenancy.NewHTTPHandler(svc, memStore, policy)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	r := httptest.NewRequest("DELETE", "/tenancy/invitations/"+inv.Token, nil)
	r = r.WithContext(tenancy.WithUser(context.Background(), "doc1"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusNoContent {
		t.Errorf("Code = %d, body = %s", w.Code, w.Body.String())
	}
}

func TestHTTPHandler_RemoveMember_Self(t *testing.T) {
	memStore := tenancy.NewMemoryMembershipStore()
	_ = memStore.Add(tenancy.Membership{UserID: "u1", TenantID: "t", Role: tenancy.TenantRoleMember})
	svc, _ := tenancy.NewInvitationService(tenancy.NewMemoryInvitationStore(), memStore, tenancy.InvitationServiceConfig{})
	handler := tenancy.NewHTTPHandler(svc, memStore, tenancy.NewPolicyEngine(memStore))
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	// 자기 자신 제거 (탈퇴)
	r := httptest.NewRequest("DELETE", "/tenancy/tenants/t/members/u1", nil)
	r = r.WithContext(tenancy.WithUser(context.Background(), "u1"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusNoContent {
		t.Errorf("Code = %d", w.Code)
	}
}

func TestHTTPHandler_RemoveMember_AdminAuthorized(t *testing.T) {
	memStore := tenancy.NewMemoryMembershipStore()
	_ = memStore.Add(tenancy.Membership{UserID: "adm", TenantID: "t", Role: tenancy.TenantRoleAdmin})
	_ = memStore.Add(tenancy.Membership{UserID: "u1", TenantID: "t", Role: tenancy.TenantRoleMember})
	svc, _ := tenancy.NewInvitationService(tenancy.NewMemoryInvitationStore(), memStore, tenancy.InvitationServiceConfig{})
	policy := tenancy.NewPolicyEngine(memStore)
	handler := tenancy.NewHTTPHandler(svc, memStore, policy)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	r := httptest.NewRequest("DELETE", "/tenancy/tenants/t/members/u1", nil)
	r = r.WithContext(tenancy.WithUser(context.Background(), "adm"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusNoContent {
		t.Errorf("Code = %d", w.Code)
	}
}

func TestHTTPHandler_RemoveMember_Forbidden(t *testing.T) {
	memStore := tenancy.NewMemoryMembershipStore()
	_ = memStore.Add(tenancy.Membership{UserID: "u1", TenantID: "t", Role: tenancy.TenantRoleMember})
	_ = memStore.Add(tenancy.Membership{UserID: "u2", TenantID: "t", Role: tenancy.TenantRoleMember})
	svc, _ := tenancy.NewInvitationService(tenancy.NewMemoryInvitationStore(), memStore, tenancy.InvitationServiceConfig{})
	policy := tenancy.NewPolicyEngine(memStore)
	handler := tenancy.NewHTTPHandler(svc, memStore, policy)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	// u1 (admin 아님) 이 u2 제거 시도
	r := httptest.NewRequest("DELETE", "/tenancy/tenants/t/members/u2", nil)
	r = r.WithContext(tenancy.WithUser(context.Background(), "u1"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusForbidden {
		t.Errorf("Code = %d", w.Code)
	}
}

func TestHTTPHandler_PathPrefix(t *testing.T) {
	memStore := tenancy.NewMemoryMembershipStore()
	svc, _ := tenancy.NewInvitationService(tenancy.NewMemoryInvitationStore(), memStore, tenancy.InvitationServiceConfig{})
	handler := tenancy.NewHTTPHandler(svc, memStore, nil)
	handler.SetPathPrefix("/api/v1")
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	r := httptest.NewRequest("GET", "/api/v1/tenancy/me/memberships", nil)
	r = r.WithContext(tenancy.WithUser(context.Background(), "u"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("Code = %d", w.Code)
	}
}

func TestHTTPHandler_ListTenantMembers_AdminAuthorized(t *testing.T) {
	memStore := tenancy.NewMemoryMembershipStore()
	_ = memStore.Add(tenancy.Membership{UserID: "adm", TenantID: "t", Role: tenancy.TenantRoleAdmin})
	_ = memStore.Add(tenancy.Membership{UserID: "u1", TenantID: "t", Role: tenancy.TenantRoleMember})
	_ = memStore.Add(tenancy.Membership{UserID: "u2", TenantID: "t", Role: tenancy.TenantRoleViewer})

	svc, _ := tenancy.NewInvitationService(tenancy.NewMemoryInvitationStore(), memStore, tenancy.InvitationServiceConfig{})
	policy := tenancy.NewPolicyEngine(memStore)
	handler := tenancy.NewHTTPHandler(svc, memStore, policy)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	r := httptest.NewRequest("GET", "/tenancy/tenants/t/members", nil)
	r = r.WithContext(tenancy.WithUser(context.Background(), "adm"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("Code = %d", w.Code)
	}
	var resp struct {
		Members []map[string]interface{} `json:"members"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if len(resp.Members) != 3 {
		t.Errorf("members = %d, want 3", len(resp.Members))
	}
}

func TestHTTPHandler_ListTenantMembers_Forbidden(t *testing.T) {
	memStore := tenancy.NewMemoryMembershipStore()
	_ = memStore.Add(tenancy.Membership{UserID: "u1", TenantID: "t", Role: tenancy.TenantRoleMember})
	svc, _ := tenancy.NewInvitationService(tenancy.NewMemoryInvitationStore(), memStore, tenancy.InvitationServiceConfig{})
	handler := tenancy.NewHTTPHandler(svc, memStore, tenancy.NewPolicyEngine(memStore))
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	r := httptest.NewRequest("GET", "/tenancy/tenants/t/members", nil)
	r = r.WithContext(tenancy.WithUser(context.Background(), "u1"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusForbidden {
		t.Errorf("Code = %d", w.Code)
	}
}

func TestHTTPHandler_UpdateRole_Success(t *testing.T) {
	memStore := tenancy.NewMemoryMembershipStore()
	_ = memStore.Add(tenancy.Membership{UserID: "adm", TenantID: "t", Role: tenancy.TenantRoleAdmin})
	_ = memStore.Add(tenancy.Membership{UserID: "u1", TenantID: "t", Role: tenancy.TenantRoleViewer})

	svc, _ := tenancy.NewInvitationService(tenancy.NewMemoryInvitationStore(), memStore, tenancy.InvitationServiceConfig{})
	policy := tenancy.NewPolicyEngine(memStore)
	handler := tenancy.NewHTTPHandler(svc, memStore, policy)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	r := httptest.NewRequest("PATCH", "/tenancy/tenants/t/members/u1/role",
		strings.NewReader(`{"role":"medical_staff"}`))
	r.Header.Set("Content-Type", "application/json")
	r = r.WithContext(tenancy.WithUser(context.Background(), "adm"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("Code = %d, body = %s", w.Code, w.Body.String())
	}

	got, _ := memStore.Get("u1", "t")
	if got.Role != tenancy.TenantRoleMedicalStaff {
		t.Errorf("Role = %q", got.Role)
	}
}

func TestHTTPHandler_UpdateRole_UnknownRole(t *testing.T) {
	memStore := tenancy.NewMemoryMembershipStore()
	_ = memStore.Add(tenancy.Membership{UserID: "adm", TenantID: "t", Role: tenancy.TenantRoleAdmin})
	_ = memStore.Add(tenancy.Membership{UserID: "u1", TenantID: "t", Role: tenancy.TenantRoleViewer})

	svc, _ := tenancy.NewInvitationService(tenancy.NewMemoryInvitationStore(), memStore, tenancy.InvitationServiceConfig{})
	handler := tenancy.NewHTTPHandler(svc, memStore, tenancy.NewPolicyEngine(memStore))
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	r := httptest.NewRequest("PATCH", "/tenancy/tenants/t/members/u1/role",
		strings.NewReader(`{"role":"ghost"}`))
	r = r.WithContext(tenancy.WithUser(context.Background(), "adm"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Code = %d", w.Code)
	}
}

func TestHTTPHandler_UpdateRole_NotFound(t *testing.T) {
	memStore := tenancy.NewMemoryMembershipStore()
	_ = memStore.Add(tenancy.Membership{UserID: "adm", TenantID: "t", Role: tenancy.TenantRoleAdmin})
	svc, _ := tenancy.NewInvitationService(tenancy.NewMemoryInvitationStore(), memStore, tenancy.InvitationServiceConfig{})
	handler := tenancy.NewHTTPHandler(svc, memStore, tenancy.NewPolicyEngine(memStore))
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	r := httptest.NewRequest("PATCH", "/tenancy/tenants/t/members/ghost/role",
		strings.NewReader(`{"role":"member"}`))
	r = r.WithContext(tenancy.WithUser(context.Background(), "adm"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusNotFound {
		t.Errorf("Code = %d", w.Code)
	}
}

func TestHTTPHandler_UpdateRole_Forbidden(t *testing.T) {
	memStore := tenancy.NewMemoryMembershipStore()
	_ = memStore.Add(tenancy.Membership{UserID: "u1", TenantID: "t", Role: tenancy.TenantRoleMember})
	_ = memStore.Add(tenancy.Membership{UserID: "u2", TenantID: "t", Role: tenancy.TenantRoleMember})
	svc, _ := tenancy.NewInvitationService(tenancy.NewMemoryInvitationStore(), memStore, tenancy.InvitationServiceConfig{})
	handler := tenancy.NewHTTPHandler(svc, memStore, tenancy.NewPolicyEngine(memStore))
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	r := httptest.NewRequest("PATCH", "/tenancy/tenants/t/members/u2/role",
		strings.NewReader(`{"role":"admin"}`))
	r = r.WithContext(tenancy.WithUser(context.Background(), "u1"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusForbidden {
		t.Errorf("Code = %d", w.Code)
	}
}

func TestHTTPHandler_MethodNotAllowed(t *testing.T) {
	memStore := tenancy.NewMemoryMembershipStore()
	svc, _ := tenancy.NewInvitationService(tenancy.NewMemoryInvitationStore(), memStore, tenancy.InvitationServiceConfig{})
	handler := tenancy.NewHTTPHandler(svc, memStore, nil)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	r := httptest.NewRequest("PUT", "/tenancy/invitations", nil)
	r = r.WithContext(tenancy.WithUser(context.Background(), "u"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Code = %d", w.Code)
	}
}
