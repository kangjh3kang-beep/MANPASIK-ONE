// Package handler는 auth-service의 gRPC 핸들러입니다.
package handler

import (
	"context"

	"github.com/manpasik/backend/services/auth-service/internal/service"
	apperrors "github.com/manpasik/backend/shared/errors"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthHandler는 AuthService gRPC 서버를 구현합니다.
type AuthHandler struct {
	v1.UnimplementedAuthServiceServer
	auth *service.AuthService
	log  *zap.Logger
}

// NewAuthHandler는 AuthHandler를 생성합니다.
func NewAuthHandler(auth *service.AuthService, log *zap.Logger) *AuthHandler {
	return &AuthHandler{auth: auth, log: log}
}

// Register는 사용자 등록 RPC입니다.
func (h *AuthHandler) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	if req == nil || req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email과 password는 필수입니다")
	}

	user, err := h.auth.Register(ctx, req.Email, req.Password, req.DisplayName)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.RegisterResponse{
		UserId:      user.ID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		Role:        user.Role,
	}, nil
}

// Login은 로그인 RPC입니다.
func (h *AuthHandler) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {
	if req == nil || req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email과 password는 필수입니다")
	}

	tokens, err := h.auth.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
		TokenType:    tokens.TokenType,
	}, nil
}

// RefreshToken은 토큰 갱신 RPC입니다.
func (h *AuthHandler) RefreshToken(ctx context.Context, req *v1.RefreshTokenRequest) (*v1.LoginResponse, error) {
	if req == nil || req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh_token은 필수입니다")
	}

	tokens, err := h.auth.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
		TokenType:    tokens.TokenType,
	}, nil
}

// Logout은 로그아웃 RPC입니다.
func (h *AuthHandler) Logout(ctx context.Context, req *v1.LogoutRequest) (*v1.LogoutResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id는 필수입니다")
	}

	if err := h.auth.Logout(ctx, req.UserId); err != nil {
		return nil, toGRPC(err)
	}

	return &v1.LogoutResponse{}, nil
}

// ValidateToken은 Access Token 검증 RPC입니다.
func (h *AuthHandler) ValidateToken(ctx context.Context, req *v1.ValidateTokenRequest) (*v1.ValidateTokenResponse, error) {
	if req == nil || req.AccessToken == "" {
		return nil, status.Error(codes.InvalidArgument, "access_token은 필수입니다")
	}

	claims, err := h.auth.ValidateToken(req.AccessToken)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &v1.ValidateTokenResponse{
		UserId: claims.UserID,
		Email:  claims.Email,
		Role:   claims.Role,
	}, nil
}

// toGRPC는 AppError를 gRPC status로 변환합니다.
func toGRPC(err error) error {
	if err == nil {
		return nil
	}
	if ae, ok := err.(*apperrors.AppError); ok {
		return ae.ToGRPC()
	}
	if s, ok := status.FromError(err); ok {
		return s.Err()
	}
	return status.Error(codes.Internal, "내부 오류가 발생했습니다")
}
