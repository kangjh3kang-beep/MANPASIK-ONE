package tenancy_test

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/manpasik/backend/shared/tenancy"
)

func TestMembership_Validate(t *testing.T) {
	tests := []struct {
		name    string
		m       tenancy.Membership
		wantErr bool
	}{
		{"valid", tenancy.Membership{UserID: "u1", TenantID: "t1", Role: tenancy.TenantRoleMember}, false},
		{"no user", tenancy.Membership{TenantID: "t1", Role: tenancy.TenantRoleMember}, true},
		{"no tenant", tenancy.Membership{UserID: "u1", Role: tenancy.TenantRoleMember}, true},
		{"unknown role", tenancy.Membership{UserID: "u1", TenantID: "t1", Role: "ghost"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.m.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("err = %v, wantErr=%v", err, tt.wantErr)
			}
		})
	}
}

func TestMemoryMembershipStore_AddGet(t *testing.T) {
	s := tenancy.NewMemoryMembershipStore()
	if err := s.Add(tenancy.Membership{UserID: "u1", TenantID: "t1", Role: tenancy.TenantRoleAdmin}); err != nil {
		t.Fatal(err)
	}
	got, err := s.Get("u1", "t1")
	if err != nil {
		t.Fatal(err)
	}
	if got.Role != tenancy.TenantRoleAdmin {
		t.Errorf("Role = %q", got.Role)
	}
	if !got.Active {
		t.Error("새 멤버는 활성 상태여야 함")
	}
}

func TestMemoryMembershipStore_GetNotFound(t *testing.T) {
	s := tenancy.NewMemoryMembershipStore()
	if _, err := s.Get("ghost", "t1"); err != tenancy.ErrNoMembership {
		t.Errorf("err = %v, want ErrNoMembership", err)
	}
}

func TestMemoryMembershipStore_ListUserTenants(t *testing.T) {
	s := tenancy.NewMemoryMembershipStore()
	_ = s.Add(tenancy.Membership{UserID: "u1", TenantID: "t1", Role: tenancy.TenantRoleAdmin})
	_ = s.Add(tenancy.Membership{UserID: "u1", TenantID: "t2", Role: tenancy.TenantRoleViewer})
	_ = s.Add(tenancy.Membership{UserID: "u2", TenantID: "t1", Role: tenancy.TenantRoleMember})

	list := s.ListUserTenants("u1")
	if len(list) != 2 {
		t.Errorf("len = %d, want 2", len(list))
	}
}

func TestMemoryMembershipStore_UpdateRole(t *testing.T) {
	s := tenancy.NewMemoryMembershipStore()
	_ = s.Add(tenancy.Membership{UserID: "u1", TenantID: "t1", Role: tenancy.TenantRoleViewer})
	if err := s.UpdateRole("u1", "t1", tenancy.TenantRoleAdmin); err != nil {
		t.Fatal(err)
	}
	got, _ := s.Get("u1", "t1")
	if got.Role != tenancy.TenantRoleAdmin {
		t.Errorf("Role = %q, want admin", got.Role)
	}
}

func TestMemoryMembershipStore_SetActive(t *testing.T) {
	s := tenancy.NewMemoryMembershipStore()
	_ = s.Add(tenancy.Membership{UserID: "u1", TenantID: "t1", Role: tenancy.TenantRoleMember})
	_ = s.SetActive("u1", "t1", false)
	got, _ := s.Get("u1", "t1")
	if got.Active {
		t.Error("Active=false 적용 안됨")
	}
}

func TestPolicyEngine_OwnerOverride(t *testing.T) {
	s := tenancy.NewMemoryMembershipStore()
	eng := tenancy.NewPolicyEngine(s)

	// 멤버십 없어도 OwnerID 일치하면 허용
	res := &tenancy.Resource{TenantID: "t1", OwnerID: "u1"}
	d := eng.Evaluate("u1", "t1", res, tenancy.ActionWrite)
	if !d.Allowed {
		t.Errorf("owner override 미적용: %v", d)
	}
}

func TestPolicyEngine_OwnerOverrideDisabled(t *testing.T) {
	s := tenancy.NewMemoryMembershipStore()
	eng := tenancy.NewPolicyEngine(s)
	eng.SetAllowOwnerOverride(false)

	res := &tenancy.Resource{TenantID: "t1", OwnerID: "u1"}
	d := eng.Evaluate("u1", "t1", res, tenancy.ActionWrite)
	if d.Allowed {
		t.Error("owner override 비활성인데 허용됨")
	}
}

func TestPolicyEngine_CrossTenantBlocked(t *testing.T) {
	s := tenancy.NewMemoryMembershipStore()
	_ = s.Add(tenancy.Membership{UserID: "u1", TenantID: "t1", Role: tenancy.TenantRoleAdmin})
	eng := tenancy.NewPolicyEngine(s)

	// t1 멤버이지만 t2 리소스에 접근 시도
	res := &tenancy.Resource{TenantID: "t2", OwnerID: "other"}
	d := eng.Evaluate("u1", "t1", res, tenancy.ActionRead)
	if d.Allowed {
		t.Error("교차 조직 접근 허용됨")
	}
}

func TestPolicyEngine_RoleBasedAccess(t *testing.T) {
	s := tenancy.NewMemoryMembershipStore()
	_ = s.Add(tenancy.Membership{UserID: "viewer", TenantID: "t1", Role: tenancy.TenantRoleViewer})
	_ = s.Add(tenancy.Membership{UserID: "admin", TenantID: "t1", Role: tenancy.TenantRoleAdmin})
	eng := tenancy.NewPolicyEngine(s)

	res := &tenancy.Resource{TenantID: "t1", OwnerID: "other"}

	// viewer 는 read 가능
	if !eng.Evaluate("viewer", "t1", res, tenancy.ActionRead).Allowed {
		t.Error("viewer read 거부됨")
	}
	// viewer 는 write 불가
	if eng.Evaluate("viewer", "t1", res, tenancy.ActionWrite).Allowed {
		t.Error("viewer write 허용됨")
	}
	// admin 은 모두 가능
	if !eng.Evaluate("admin", "t1", res, tenancy.ActionDelete).Allowed {
		t.Error("admin delete 거부됨")
	}
}

func TestPolicyEngine_InactiveMember(t *testing.T) {
	s := tenancy.NewMemoryMembershipStore()
	_ = s.Add(tenancy.Membership{UserID: "u1", TenantID: "t1", Role: tenancy.TenantRoleAdmin})
	_ = s.SetActive("u1", "t1", false)
	eng := tenancy.NewPolicyEngine(s)

	res := &tenancy.Resource{TenantID: "t1", OwnerID: "other"}
	d := eng.Evaluate("u1", "t1", res, tenancy.ActionRead)
	if d.Allowed {
		t.Error("비활성 멤버 허용됨")
	}
}

func TestPolicyEngine_SharedResource(t *testing.T) {
	s := tenancy.NewMemoryMembershipStore()
	eng := tenancy.NewPolicyEngine(s)

	// 멤버 아니어도 공유받은 리소스는 read 가능
	res := &tenancy.Resource{
		TenantID:   "t1",
		OwnerID:    "other",
		SharedWith: []string{"u1"},
	}
	d := eng.Evaluate("u1", "t1", res, tenancy.ActionRead)
	if !d.Allowed {
		t.Errorf("공유 리소스 read 거부: %v", d)
	}

	// 그러나 write 는 불가
	dw := eng.Evaluate("u1", "t1", res, tenancy.ActionWrite)
	if dw.Allowed {
		t.Error("공유 리소스에 write 허용됨")
	}
}

func TestPolicyEngine_FilterResources(t *testing.T) {
	s := tenancy.NewMemoryMembershipStore()
	_ = s.Add(tenancy.Membership{UserID: "u1", TenantID: "t1", Role: tenancy.TenantRoleViewer})
	eng := tenancy.NewPolicyEngine(s)

	resources := []*tenancy.Resource{
		{TenantID: "t1", OwnerID: "other"}, // 멤버 → read 가능
		{TenantID: "t2", OwnerID: "other"}, // 교차 → 차단
		{TenantID: "t1", OwnerID: "u1"},    // 본인 → 허용
	}
	filtered := eng.FilterResources("u1", "t1", resources)
	if len(filtered) != 2 {
		t.Errorf("filtered = %d, want 2", len(filtered))
	}
}

func TestContext_TenantFromContext_Direct(t *testing.T) {
	ctx := tenancy.WithTenant(context.Background(), "t1")
	got, ok := tenancy.TenantFromContext(ctx)
	if !ok || got != "t1" {
		t.Errorf("TenantFromContext = %q, %v", got, ok)
	}
}

func TestContext_TenantFromContext_Metadata(t *testing.T) {
	md := metadata.Pairs(tenancy.MetadataTenantKey, "t-meta")
	ctx := metadata.NewIncomingContext(context.Background(), md)
	got, ok := tenancy.TenantFromContext(ctx)
	if !ok || got != "t-meta" {
		t.Errorf("TenantFromContext = %q, %v", got, ok)
	}
}

func TestContext_TenantFromContext_NotPresent(t *testing.T) {
	if _, ok := tenancy.TenantFromContext(context.Background()); ok {
		t.Error("빈 ctx 에서 테넌트 발견됨")
	}
}

func TestContext_UserFromContext(t *testing.T) {
	ctx := tenancy.WithUser(context.Background(), "u-ctx")
	got, ok := tenancy.UserFromContext(ctx)
	if !ok || got != "u-ctx" {
		t.Errorf("got %q, %v", got, ok)
	}
}

func TestUnaryInterceptor_Allow(t *testing.T) {
	store := tenancy.NewMemoryMembershipStore()
	_ = store.Add(tenancy.Membership{UserID: "u1", TenantID: "t1", Role: tenancy.TenantRoleAdmin})
	eng := tenancy.NewPolicyEngine(store)

	intc := tenancy.UnaryInterceptor(&tenancy.InterceptorConfig{Engine: eng})
	ctx := tenancy.WithUser(tenancy.WithTenant(context.Background(), "t1"), "u1")

	called := false
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		called = true
		return "ok", nil
	}
	resp, err := intc(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/svc/M"}, handler)
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if !called || resp != "ok" {
		t.Error("핸들러 호출 안됨")
	}
}

func TestUnaryInterceptor_NoTenant(t *testing.T) {
	eng := tenancy.NewPolicyEngine(tenancy.NewMemoryMembershipStore())
	intc := tenancy.UnaryInterceptor(&tenancy.InterceptorConfig{Engine: eng})
	ctx := tenancy.WithUser(context.Background(), "u1") // 테넌트 없음

	_, err := intc(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/svc/M"},
		func(ctx context.Context, req interface{}) (interface{}, error) { return nil, nil })

	if err == nil {
		t.Fatal("테넌트 없는데 통과")
	}
	if status.Code(err) != codes.PermissionDenied {
		t.Errorf("코드 = %v, want PermissionDenied", status.Code(err))
	}
}

func TestUnaryInterceptor_NoUser(t *testing.T) {
	eng := tenancy.NewPolicyEngine(tenancy.NewMemoryMembershipStore())
	intc := tenancy.UnaryInterceptor(&tenancy.InterceptorConfig{Engine: eng})
	ctx := tenancy.WithTenant(context.Background(), "t1")

	_, err := intc(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/svc/M"},
		func(ctx context.Context, req interface{}) (interface{}, error) { return nil, nil })
	if status.Code(err) != codes.Unauthenticated {
		t.Errorf("코드 = %v, want Unauthenticated", status.Code(err))
	}
}

func TestUnaryInterceptor_SkipMethod(t *testing.T) {
	eng := tenancy.NewPolicyEngine(tenancy.NewMemoryMembershipStore())
	intc := tenancy.UnaryInterceptor(&tenancy.InterceptorConfig{
		Engine:      eng,
		SkipMethods: map[string]bool{"/svc/Login": true},
	})
	// 테넌트/유저 없어도 skip 메서드는 통과
	called := false
	_, err := intc(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/svc/Login"},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			called = true
			return nil, nil
		})
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if !called {
		t.Error("skip 메서드 핸들러 호출 안됨")
	}
}

func TestUnaryInterceptor_NotMember(t *testing.T) {
	eng := tenancy.NewPolicyEngine(tenancy.NewMemoryMembershipStore())
	intc := tenancy.UnaryInterceptor(&tenancy.InterceptorConfig{Engine: eng})
	ctx := tenancy.WithUser(tenancy.WithTenant(context.Background(), "t1"), "ghost")

	_, err := intc(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/svc/M"},
		func(ctx context.Context, req interface{}) (interface{}, error) { return nil, nil })
	if status.Code(err) != codes.PermissionDenied {
		t.Errorf("코드 = %v, want PermissionDenied", status.Code(err))
	}
}

func TestAppendTenantToOutgoing(t *testing.T) {
	ctx := tenancy.AppendTenantToOutgoing(context.Background(), "t-out")
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		t.Fatal("outgoing metadata 없음")
	}
	vals := md.Get(tenancy.MetadataTenantKey)
	if len(vals) == 0 || vals[0] != "t-out" {
		t.Errorf("vals = %v", vals)
	}
}

func TestPolicyEngine_CustomPermission(t *testing.T) {
	store := tenancy.NewMemoryMembershipStore()
	_ = store.Add(tenancy.Membership{UserID: "u1", TenantID: "t1", Role: tenancy.TenantRoleViewer})
	eng := tenancy.NewPolicyEngine(store)

	// viewer 에게 write 권한 추가
	eng.SetRolePermission(tenancy.TenantRoleViewer, tenancy.ActionWrite, true)

	res := &tenancy.Resource{TenantID: "t1", OwnerID: "other"}
	d := eng.Evaluate("u1", "t1", res, tenancy.ActionWrite)
	if !d.Allowed {
		t.Errorf("커스텀 권한 미적용: %v", d)
	}
}

func TestDecision_String(t *testing.T) {
	d := tenancy.Decision{Allowed: true, Reason: "owner"}
	s := d.String()
	if len(s) == 0 || s[:5] != "ALLOW" {
		t.Errorf("String = %q", s)
	}
}
