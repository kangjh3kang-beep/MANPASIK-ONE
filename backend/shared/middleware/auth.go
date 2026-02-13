// Package middleware는 gRPC 인터셉터를 제공합니다.
package middleware

import (
	"context"
	"strings"

	apperrors "github.com/manpasik/backend/shared/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// 컨텍스트 키
type contextKey string

const (
	// UserIDKey는 인증된 사용자 ID의 컨텍스트 키입니다.
	UserIDKey contextKey = "user_id"
	// UserRoleKey는 사용자 역할의 컨텍스트 키입니다.
	UserRoleKey contextKey = "user_role"
	// UserEmailKey는 사용자 이메일의 컨텍스트 키입니다.
	UserEmailKey contextKey = "user_email"
)

// TokenValidator는 JWT 토큰 검증 인터페이스입니다.
type TokenValidator interface {
	ValidateToken(token string) (*TokenClaims, error)
}

// TokenClaims는 JWT 토큰의 클레임입니다.
type TokenClaims struct {
	UserID string
	Email  string
	Role   string
}

// 인증이 필요 없는 RPC 메서드 목록
var publicMethods = map[string]bool{
	"/grpc.health.v1.Health/Check":                         true,
	"/grpc.health.v1.Health/Watch":                         true,
	"/manpasik.health.v1.ExtendedHealth/Ready":             true,
	"/manpasik.health.v1.ExtendedHealth/Live":              true,
	"/manpasik.health.v1.ExtendedHealth/GetStatus":         true,
	"/manpasik.health.v1.ExtendedHealth/CheckDependencies": true,
	"/manpasik.v1.AuthService/Register":                    true,
	"/manpasik.v1.AuthService/Login":                       true,
	"/manpasik.v1.AuthService/RefreshToken":                true,
	"/manpasik.v1.AuthService/Logout":                      true,
	"/manpasik.v1.AuthService/ValidateToken":               true,
}

// AuthInterceptor는 JWT 인증 gRPC Unary 인터셉터입니다.
func AuthInterceptor(validator TokenValidator) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// 공개 메서드는 인증 건너뜀
		if publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		// 메타데이터에서 토큰 추출
		token, err := extractToken(ctx)
		if err != nil {
			return nil, err
		}

		// 토큰 검증
		claims, err := validator.ValidateToken(token)
		if err != nil {
			return nil, apperrors.New(apperrors.ErrInvalidToken, "유효하지 않은 인증 토큰").ToGRPC()
		}

		// 컨텍스트에 사용자 정보 저장
		ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserRoleKey, claims.Role)
		ctx = context.WithValue(ctx, UserEmailKey, claims.Email)

		return handler(ctx, req)
	}
}

// AuthStreamInterceptor는 JWT 인증 gRPC Stream 인터셉터입니다.
func AuthStreamInterceptor(validator TokenValidator) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if publicMethods[info.FullMethod] {
			return handler(srv, ss)
		}

		token, err := extractToken(ss.Context())
		if err != nil {
			return err
		}

		claims, err := validator.ValidateToken(token)
		if err != nil {
			return apperrors.New(apperrors.ErrInvalidToken, "유효하지 않은 인증 토큰").ToGRPC()
		}

		ctx := ss.Context()
		ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserRoleKey, claims.Role)
		ctx = context.WithValue(ctx, UserEmailKey, claims.Email)

		wrapped := &wrappedStream{ServerStream: ss, ctx: ctx}
		return handler(srv, wrapped)
	}
}

// extractToken은 gRPC 메타데이터에서 Bearer 토큰을 추출합니다.
func extractToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", apperrors.New(apperrors.ErrUnauthorized, "인증 메타데이터 없음").ToGRPC()
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return "", apperrors.New(apperrors.ErrUnauthorized, "Authorization 헤더 없음").ToGRPC()
	}

	authHeader := values[0]
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", apperrors.New(apperrors.ErrInvalidToken, "Bearer 토큰 형식이 아닙니다").ToGRPC()
	}

	return strings.TrimPrefix(authHeader, "Bearer "), nil
}

// UserIDFromContext는 컨텍스트에서 사용자 ID를 추출합니다.
func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}

// UserRoleFromContext는 컨텍스트에서 사용자 역할을 추출합니다.
func UserRoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(UserRoleKey).(string)
	return role, ok
}

// wrappedStream은 컨텍스트를 오버라이드하는 스트림 래퍼입니다.
type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedStream) Context() context.Context {
	return w.ctx
}
