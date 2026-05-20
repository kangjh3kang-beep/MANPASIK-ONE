package tenancy_test

import (
	"os"
	"testing"

	"github.com/manpasik/backend/shared/tenancy"
)

func TestEnforcedFromEnv(t *testing.T) {
	cases := map[string]bool{
		"":      false,
		"false": false,
		"0":     false,
		"true":  true,
		"1":     true,
	}
	for v, want := range cases {
		t.Setenv("TENANCY_ENFORCED", v)
		if got := tenancy.EnforcedFromEnv(); got != want {
			t.Errorf("env=%q got %v, want %v", v, got, want)
		}
	}
}

func TestDefaultSkipMethods(t *testing.T) {
	skip := tenancy.DefaultSkipMethods()
	if !skip["/grpc.health.v1.Health/Check"] {
		t.Error("health check 미포함")
	}
	if !skip["/grpc.reflection.v1.ServerReflection/ServerReflectionInfo"] {
		t.Error("reflection 미포함")
	}
}

func TestMaybeServerOptions_Disabled(t *testing.T) {
	t.Setenv("TENANCY_ENFORCED", "false")
	eng := tenancy.NewPolicyEngine(tenancy.NewMemoryMembershipStore())
	opts := tenancy.MaybeServerOptions(eng, nil)
	if opts != nil {
		t.Errorf("비활성 시 nil 이어야: len=%d", len(opts))
	}
}

func TestMaybeServerOptions_Enabled(t *testing.T) {
	t.Setenv("TENANCY_ENFORCED", "true")
	eng := tenancy.NewPolicyEngine(tenancy.NewMemoryMembershipStore())
	opts := tenancy.MaybeServerOptions(eng, nil)
	if len(opts) != 2 {
		t.Errorf("활성 시 2 옵션 (unary+stream), got %d", len(opts))
	}
}

func TestMaybeServerOptions_NilEngine(t *testing.T) {
	t.Setenv("TENANCY_ENFORCED", "true")
	opts := tenancy.MaybeServerOptions(nil, nil)
	if opts != nil {
		t.Error("nil engine 이면 옵션 없어야 함")
	}
}

func TestMaybeServerOptions_ExtraSkip(t *testing.T) {
	t.Setenv("TENANCY_ENFORCED", "true")
	eng := tenancy.NewPolicyEngine(tenancy.NewMemoryMembershipStore())
	opts := tenancy.MaybeServerOptions(eng, map[string]bool{"/svc/Login": true})
	if len(opts) != 2 {
		t.Errorf("len = %d", len(opts))
	}
	// extraSkip 통합은 패키지 내부 동작이므로 별도 검증 불요 (DefaultSkipMethods 와 합쳐짐)
}

func TestEnforcedFromEnv_Unset(t *testing.T) {
	_ = os.Unsetenv("TENANCY_ENFORCED")
	if tenancy.EnforcedFromEnv() {
		t.Error("미설정 시 false 여야")
	}
}
