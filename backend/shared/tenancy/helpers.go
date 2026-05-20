package tenancy

import (
	"os"

	"google.golang.org/grpc"
)

// EnforcedFromEnv 는 TENANCY_ENFORCED 환경변수가 "true"/"1" 인지 반환.
//
// false 면 인터셉터를 등록하지 않거나 등록해도 검증을 건너뜀 (점진 도입용).
func EnforcedFromEnv() bool {
	v := os.Getenv("TENANCY_ENFORCED")
	return v == "true" || v == "1"
}

// DefaultSkipMethods 는 인터셉터에서 항상 통과시킬 시스템 RPC.
//
// 헬스체크/리플렉션은 노드 외부 모니터링/디버깅 도구가 호출하므로 테넌시 검증 면제.
func DefaultSkipMethods() map[string]bool {
	return map[string]bool{
		"/grpc.health.v1.Health/Check":                              true,
		"/grpc.health.v1.Health/Watch":                              true,
		"/grpc.reflection.v1.ServerReflection/ServerReflectionInfo": true,
		"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo": true,
	}
}

// MaybeServerOptions 는 ENV 토글에 따라 인터셉터 옵션을 반환.
//
// extraSkip 으로 서비스별 추가 skip 메서드를 합칠 수 있다.
//
// 사용 예:
//
//	store := tenancy.NewMemoryMembershipStore()
//	engine := tenancy.NewPolicyEngine(store)
//	opts := tenancy.MaybeServerOptions(engine, nil)
//	grpc.NewServer(append(existing, opts...)...)
func MaybeServerOptions(engine *PolicyEngine, extraSkip map[string]bool) []grpc.ServerOption {
	if !EnforcedFromEnv() || engine == nil {
		return nil
	}
	skip := DefaultSkipMethods()
	for k, v := range extraSkip {
		skip[k] = v
	}
	cfg := &InterceptorConfig{Engine: engine, SkipMethods: skip}
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(UnaryInterceptor(cfg)),
		grpc.StreamInterceptor(StreamInterceptor(cfg)),
	}
}
