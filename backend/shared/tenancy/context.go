package tenancy

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// 메타데이터 키.
const (
	MetadataTenantKey = "x-tenant-id"
	MetadataUserKey   = "x-user-id"
)

type ctxKey int

const (
	ctxKeyTenant ctxKey = iota
	ctxKeyUser
)

// WithTenant 는 ctx 에 테넌트 ID 주입.
func WithTenant(ctx context.Context, t TenantID) context.Context {
	return context.WithValue(ctx, ctxKeyTenant, t)
}

// TenantFromContext 는 ctx 에서 테넌트 ID 추출. 없으면 (zero, false).
//
// 우선순위:
//  1. context.Value 직접 주입값
//  2. gRPC metadata 의 x-tenant-id
func TenantFromContext(ctx context.Context) (TenantID, bool) {
	if v := ctx.Value(ctxKeyTenant); v != nil {
		if t, ok := v.(TenantID); ok && !t.IsZero() {
			return t, true
		}
	}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		vals := md.Get(MetadataTenantKey)
		if len(vals) > 0 && vals[0] != "" {
			return TenantID(vals[0]), true
		}
	}
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		vals := md.Get(MetadataTenantKey)
		if len(vals) > 0 && vals[0] != "" {
			return TenantID(vals[0]), true
		}
	}
	return "", false
}

// WithUser 는 ctx 에 사용자 ID 주입.
func WithUser(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, ctxKeyUser, userID)
}

// UserFromContext 는 ctx 에서 사용자 ID 추출.
func UserFromContext(ctx context.Context) (string, bool) {
	if v := ctx.Value(ctxKeyUser); v != nil {
		if s, ok := v.(string); ok && s != "" {
			return s, true
		}
	}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		vals := md.Get(MetadataUserKey)
		if len(vals) > 0 && vals[0] != "" {
			return vals[0], true
		}
	}
	return "", false
}

// AppendTenantToOutgoing 은 outgoing gRPC 메타데이터에 테넌트 ID 부착.
func AppendTenantToOutgoing(ctx context.Context, t TenantID) context.Context {
	if t.IsZero() {
		return ctx
	}
	return metadata.AppendToOutgoingContext(ctx, MetadataTenantKey, string(t))
}

// AppendUserToOutgoing 은 outgoing gRPC 메타데이터에 사용자 ID 부착.
func AppendUserToOutgoing(ctx context.Context, userID string) context.Context {
	if userID == "" {
		return ctx
	}
	return metadata.AppendToOutgoingContext(ctx, MetadataUserKey, userID)
}
