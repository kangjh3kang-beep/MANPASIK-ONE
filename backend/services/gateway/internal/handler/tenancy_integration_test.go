package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/manpasik/backend/shared/tenancy"
)

// gateway main.go 와 동일한 구조로 tenancy HTTP 핸들러를 마운트하여
// 클라이언트 → gateway → tenancy 흐름이 정상 동작하는지 검증.

func setupGatewayTenancy(t *testing.T) http.Handler {
	t.Helper()
	mem := tenancy.NewMemoryMembershipStore()
	_ = mem.Add(tenancy.Membership{UserID: "adm", TenantID: "hospA", Role: tenancy.TenantRoleAdmin})

	invStore := tenancy.InvitationStore(tenancy.NewMemoryInvitationStore())
	engine := tenancy.NewPolicyEngine(mem)
	invSvc, _ := tenancy.NewInvitationService(invStore, mem, tenancy.InvitationServiceConfig{})
	invSvc.SetPolicyEngine(engine)

	h := tenancy.NewHTTPHandler(invSvc, mem, engine)
	h.SetPathPrefix("/api/v1")
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return mux
}

func TestGateway_Tenancy_InviteCreated(t *testing.T) {
	mux := setupGatewayTenancy(t)

	r := httptest.NewRequest("POST", "/api/v1/tenancy/invitations",
		strings.NewReader(`{"tenant_id":"hospA","role":"member"}`))
	r.Header.Set("Content-Type", "application/json")
	r = r.WithContext(tenancy.WithUser(context.Background(), "adm"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("Code = %d, body = %s", w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["status"] != "pending" {
		t.Errorf("status = %v", resp["status"])
	}
}

func TestGateway_Tenancy_AcceptFlow(t *testing.T) {
	mux := setupGatewayTenancy(t)

	// 1. invite 발급 (admin 으로)
	r1 := httptest.NewRequest("POST", "/api/v1/tenancy/invitations",
		strings.NewReader(`{"tenant_id":"hospA","role":"viewer"}`))
	r1 = r1.WithContext(tenancy.WithUser(context.Background(), "adm"))
	w1 := httptest.NewRecorder()
	mux.ServeHTTP(w1, r1)
	var inv map[string]interface{}
	_ = json.Unmarshal(w1.Body.Bytes(), &inv)
	token := inv["token"].(string)

	// 2. 새 사용자가 accept
	r2 := httptest.NewRequest("POST", "/api/v1/tenancy/invitations/accept",
		strings.NewReader(`{"token":"`+token+`"}`))
	r2 = r2.WithContext(tenancy.WithUser(context.Background(), "newUser"))
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, r2)
	if w2.Code != http.StatusOK {
		t.Fatalf("accept Code = %d, body = %s", w2.Code, w2.Body.String())
	}

	// 3. newUser 의 멤버십 조회
	r3 := httptest.NewRequest("GET", "/api/v1/tenancy/me/memberships", nil)
	r3 = r3.WithContext(tenancy.WithUser(context.Background(), "newUser"))
	w3 := httptest.NewRecorder()
	mux.ServeHTTP(w3, r3)
	if w3.Code != http.StatusOK {
		t.Fatalf("memberships Code = %d", w3.Code)
	}
	var mems struct {
		Memberships []map[string]interface{} `json:"memberships"`
	}
	_ = json.Unmarshal(w3.Body.Bytes(), &mems)
	if len(mems.Memberships) != 1 {
		t.Errorf("memberships len = %d", len(mems.Memberships))
	}
}

func TestGateway_Tenancy_PathPrefix(t *testing.T) {
	mux := setupGatewayTenancy(t)
	// 잘못된 prefix → 404 (default mux)
	r := httptest.NewRequest("GET", "/tenancy/me/memberships", nil)
	r = r.WithContext(tenancy.WithUser(context.Background(), "adm"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusNotFound {
		t.Errorf("미설정 prefix Code = %d, want 404", w.Code)
	}
}

func TestGateway_Tenancy_Unauthenticated(t *testing.T) {
	mux := setupGatewayTenancy(t)
	// ctx 에 사용자 없음 → 401
	r := httptest.NewRequest("GET", "/api/v1/tenancy/me/memberships", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Code = %d, want 401", w.Code)
	}
}
