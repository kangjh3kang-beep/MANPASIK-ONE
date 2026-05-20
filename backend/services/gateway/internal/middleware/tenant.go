package middleware

import (
	"context"
	"net/http"

	"github.com/manpasik/backend/shared/tenancy"
)

// TenantPropagation 미들웨어는 클라이언트가 보낸 `X-Tenant-ID` / `X-User-ID`
// 헤더를 읽어 request context 에 보관한다. 이후 gRPC 핸들러에서
// `WithTenantMetadata(ctx)` 헬퍼를 통해 outgoing 메타데이터로 자동 전파.
//
// 헤더 미존재 시 컨텍스트에 아무것도 추가하지 않음 (기존 동작 보존).
func TenantPropagation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if tid := r.Header.Get("X-Tenant-ID"); tid != "" {
			ctx = tenancy.WithTenant(ctx, tenancy.TenantID(tid))
		}
		if uid := r.Header.Get("X-User-ID"); uid != "" {
			ctx = tenancy.WithUser(ctx, uid)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// WithTenantMetadata 는 ctx 의 테넌트/유저 ID 를 outgoing gRPC 메타데이터로 부착.
//
// gateway REST 핸들러가 백엔드 gRPC 호출 직전에 이 함수를 호출하여
// 컨텍스트 → 메타데이터 변환.
//
//	ctx = WithTenantMetadata(r.Context())
//	clients.Telemedicine.CreateConsultation(ctx, req)
func WithTenantMetadata(ctx context.Context) context.Context {
	if tid, ok := tenancy.TenantFromContext(ctx); ok {
		ctx = tenancy.AppendTenantToOutgoing(ctx, tid)
	}
	if uid, ok := tenancy.UserFromContext(ctx); ok {
		ctx = tenancy.AppendUserToOutgoing(ctx, uid)
	}
	return ctx
}
