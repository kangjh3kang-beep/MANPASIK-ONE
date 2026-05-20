package tenancy

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// InterceptorConfig 는 테넌시 인터셉터 동작 설정.
type InterceptorConfig struct {
	// Engine 은 멤버십 검증에 사용되는 정책 엔진.
	Engine *PolicyEngine
	// RequireForMethods 가 비어있지 않으면 해당 메서드만 테넌트 검증.
	// 비어있으면 모든 메서드 검증.
	RequireForMethods map[string]bool
	// SkipMethods 는 테넌트 검증을 건너뛸 메서드 (인증/등록 등).
	SkipMethods map[string]bool
}

// UnaryInterceptor 는 모든 gRPC 호출에서:
//  1. 테넌트 컨텍스트 (x-tenant-id) 추출
//  2. 사용자가 해당 테넌트의 활성 멤버인지 검증
//  3. 통과하면 다음 핸들러로
func UnaryInterceptor(cfg *InterceptorConfig) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if cfg == nil || cfg.Engine == nil {
			return handler(ctx, req)
		}
		if cfg.SkipMethods[info.FullMethod] {
			return handler(ctx, req)
		}
		if len(cfg.RequireForMethods) > 0 && !cfg.RequireForMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		userID, ok := UserFromContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "사용자 컨텍스트 없음")
		}
		tenantID, ok := TenantFromContext(ctx)
		if !ok {
			return nil, status.Error(codes.PermissionDenied, ErrNoTenantContext.Error())
		}
		if _, err := cfg.Engine.CheckMembership(userID, tenantID); err != nil {
			return nil, status.Errorf(codes.PermissionDenied, "테넌시 검증 실패: %v", err)
		}
		return handler(ctx, req)
	}
}

// StreamInterceptor 는 스트리밍 RPC 용 동일 검증.
func StreamInterceptor(cfg *InterceptorConfig) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if cfg == nil || cfg.Engine == nil {
			return handler(srv, ss)
		}
		if cfg.SkipMethods[info.FullMethod] {
			return handler(srv, ss)
		}
		if len(cfg.RequireForMethods) > 0 && !cfg.RequireForMethods[info.FullMethod] {
			return handler(srv, ss)
		}

		ctx := ss.Context()
		userID, ok := UserFromContext(ctx)
		if !ok {
			return status.Error(codes.Unauthenticated, "사용자 컨텍스트 없음")
		}
		tenantID, ok := TenantFromContext(ctx)
		if !ok {
			return status.Error(codes.PermissionDenied, ErrNoTenantContext.Error())
		}
		if _, err := cfg.Engine.CheckMembership(userID, tenantID); err != nil {
			return status.Errorf(codes.PermissionDenied, "테넌시 검증 실패: %v", err)
		}
		return handler(srv, ss)
	}
}
