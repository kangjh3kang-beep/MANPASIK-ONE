package integration_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/manpasik/backend/shared/tenancy"
)

// 통합 시나리오: 멀티테넌트 격리 검증.
//
// 구성: 단일 mux 에 tenancy HTTP handler + memberStore (인메모리).
// 시나리오 검증:
//   1. A 조직 admin 이 A 의 멤버 목록 조회 → OK
//   2. A 조직 admin 이 B 조직 멤버 목록 조회 시도 → 거부 (cross-tenant)
//   3. A 조직 일반 멤버가 A 의 멤버 목록 조회 시도 → 거부 (admin 권한 부족)
//   4. 인증 안 된 요청 → 401

func setupIsolationFixture(t *testing.T) http.Handler {
	t.Helper()
	mem := tenancy.NewMemoryMembershipStore()
	// hospA: docA(admin), nurseA(member)
	_ = mem.Add(tenancy.Membership{UserID: "docA", TenantID: "hospA", Role: tenancy.TenantRoleAdmin})
	_ = mem.Add(tenancy.Membership{UserID: "nurseA", TenantID: "hospA", Role: tenancy.TenantRoleMember})
	// hospB: docB(admin)
	_ = mem.Add(tenancy.Membership{UserID: "docB", TenantID: "hospB", Role: tenancy.TenantRoleAdmin})

	invStore := tenancy.InvitationStore(tenancy.NewMemoryInvitationStore())
	engine := tenancy.NewPolicyEngine(mem)
	invSvc, _ := tenancy.NewInvitationService(invStore, mem, tenancy.InvitationServiceConfig{})
	invSvc.SetPolicyEngine(engine)

	h := tenancy.NewHTTPHandler(invSvc, mem, engine)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return mux
}

func TestTenancyIsolation_AdminListsOwnTenant_OK(t *testing.T) {
	mux := setupIsolationFixture(t)

	r := httptest.NewRequest("GET", "/tenancy/tenants/hospA/members", nil)
	r = r.WithContext(tenancy.WithUser(context.Background(), "docA"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("Code = %d, body = %s", w.Code, w.Body.String())
	}
}

func TestTenancyIsolation_AdminListsOtherTenant_Forbidden(t *testing.T) {
	mux := setupIsolationFixture(t)

	// docA (hospA admin) 가 hospB 멤버 조회 시도
	r := httptest.NewRequest("GET", "/tenancy/tenants/hospB/members", nil)
	r = r.WithContext(tenancy.WithUser(context.Background(), "docA"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusForbidden {
		t.Errorf("교차 조직 접근 허용됨: Code=%d", w.Code)
	}
}

func TestTenancyIsolation_NonAdminListsTenant_Forbidden(t *testing.T) {
	mux := setupIsolationFixture(t)

	// nurseA (hospA member) 가 hospA 멤버 조회 시도 (admin 권한 부족)
	r := httptest.NewRequest("GET", "/tenancy/tenants/hospA/members", nil)
	r = r.WithContext(tenancy.WithUser(context.Background(), "nurseA"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusForbidden {
		t.Errorf("admin 권한 부족 케이스 통과: Code=%d", w.Code)
	}
}

func TestTenancyIsolation_NonMemberListsTenant_Forbidden(t *testing.T) {
	mux := setupIsolationFixture(t)

	// outsider (어떤 조직 멤버도 아님) 가 hospA 멤버 조회 시도
	r := httptest.NewRequest("GET", "/tenancy/tenants/hospA/members", nil)
	r = r.WithContext(tenancy.WithUser(context.Background(), "outsider"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusForbidden {
		t.Errorf("비멤버 통과: Code=%d", w.Code)
	}
}

func TestTenancyIsolation_CrossTenantRoleChange_Forbidden(t *testing.T) {
	mux := setupIsolationFixture(t)

	// docA (hospA admin) 가 hospB 멤버의 역할을 변경 시도
	r := httptest.NewRequest("PATCH", "/tenancy/tenants/hospB/members/docB/role",
		strings.NewReader(`{"role":"viewer"}`))
	r.Header.Set("Content-Type", "application/json")
	r = r.WithContext(tenancy.WithUser(context.Background(), "docA"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusForbidden {
		t.Errorf("교차 조직 역할 변경 통과: Code=%d", w.Code)
	}
}

func TestTenancyIsolation_CrossTenantMemberRemoval_Forbidden(t *testing.T) {
	mux := setupIsolationFixture(t)

	// docA (hospA admin) 가 hospB 의 docB 를 제거 시도
	r := httptest.NewRequest("DELETE", "/tenancy/tenants/hospB/members/docB", nil)
	r = r.WithContext(tenancy.WithUser(context.Background(), "docA"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusForbidden {
		t.Errorf("교차 조직 제거 통과: Code=%d", w.Code)
	}
}

func TestTenancyIsolation_Unauthenticated_401(t *testing.T) {
	mux := setupIsolationFixture(t)

	// 사용자 컨텍스트 없음 → 401
	r := httptest.NewRequest("GET", "/tenancy/tenants/hospA/members", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("미인증 통과: Code=%d", w.Code)
	}
}

func TestTenancyIsolation_PolicyEngine_OwnerOverride(t *testing.T) {
	// PolicyEngine 직접 검증: 자기 데이터는 항상 접근 가능
	mem := tenancy.NewMemoryMembershipStore()
	_ = mem.Add(tenancy.Membership{UserID: "u1", TenantID: "hospA", Role: tenancy.TenantRoleMember})
	engine := tenancy.NewPolicyEngine(mem)

	res := &tenancy.Resource{
		TenantID: "hospA",
		OwnerID:  "u1", // 자기 데이터
		Type:     "measurement",
	}
	d := engine.Evaluate("u1", "hospA", res, tenancy.ActionWrite)
	if !d.Allowed {
		t.Errorf("자기 데이터 write 거부: %v", d)
	}
}

func TestTenancyIsolation_PolicyEngine_CrossTenant_Denied(t *testing.T) {
	mem := tenancy.NewMemoryMembershipStore()
	_ = mem.Add(tenancy.Membership{UserID: "docA", TenantID: "hospA", Role: tenancy.TenantRoleAdmin})
	engine := tenancy.NewPolicyEngine(mem)

	// hospA admin 이 hospB 의 리소스에 접근 시도
	res := &tenancy.Resource{
		TenantID: "hospB",
		OwnerID:  "patientB",
	}
	d := engine.Evaluate("docA", "hospA", res, tenancy.ActionRead)
	if d.Allowed {
		t.Errorf("교차 tenant 허용됨: %v", d)
	}
}
