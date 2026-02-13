// Package service는 auth-service의 비즈니스 로직을 구현합니다.
package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	apperrors "github.com/manpasik/backend/shared/errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// AuthService는 인증 서비스의 비즈니스 로직입니다.
type AuthService struct {
	logger     *zap.Logger
	userRepo   UserRepository
	tokenRepo  TokenRepository
	jwtSecret  []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
	issuer     string
}

// UserRepository는 사용자 데이터 저장소 인터페이스입니다.
type UserRepository interface {
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) error
	UpdatePassword(ctx context.Context, id string, hashedPassword string) error
}

// TokenRepository는 토큰 저장소 인터페이스입니다 (Redis).
type TokenRepository interface {
	StoreRefreshToken(ctx context.Context, userID, tokenID string, ttl time.Duration) error
	ValidateRefreshToken(ctx context.Context, userID, tokenID string) (bool, error)
	RevokeRefreshToken(ctx context.Context, userID, tokenID string) error
	RevokeAllUserTokens(ctx context.Context, userID string) error
}

// User는 사용자 엔티티입니다.
type User struct {
	ID             string
	Email          string
	HashedPassword string
	DisplayName    string
	Role           string // "user", "admin", "medical_staff", "researcher", "family_member"
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// TokenPair는 Access/Refresh 토큰 쌍입니다.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64 // Access Token 만료까지 남은 초
	TokenType    string
}

// CustomClaims는 JWT 커스텀 클레임입니다.
type CustomClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

// NewAuthService는 새 AuthService를 생성합니다.
func NewAuthService(
	logger *zap.Logger,
	userRepo UserRepository,
	tokenRepo TokenRepository,
	jwtSecret string,
	accessTTL, refreshTTL time.Duration,
	issuer string,
) *AuthService {
	return &AuthService{
		logger:     logger,
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		jwtSecret:  []byte(jwtSecret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		issuer:     issuer,
	}
}

// Register는 새 사용자를 등록합니다.
func (s *AuthService) Register(ctx context.Context, email, password, displayName string) (*User, error) {
	// 이메일 중복 확인
	existing, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil && existing != nil {
		return nil, apperrors.New(apperrors.ErrAlreadyExists, "이미 등록된 이메일입니다")
	}

	// 비밀번호 해싱 (bcrypt, cost=12)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		s.logger.Error("비밀번호 해싱 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "사용자 등록에 실패했습니다")
	}

	user := &User{
		ID:             uuid.New().String(),
		Email:          email,
		HashedPassword: string(hashedPassword),
		DisplayName:    displayName,
		Role:           "user",
		IsActive:       true,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.Error("사용자 생성 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "사용자 등록에 실패했습니다")
	}

	s.logger.Info("사용자 등록 완료", zap.String("user_id", user.ID), zap.String("email", email))
	return user, nil
}

// Login은 이메일/비밀번호로 로그인하고 토큰 쌍을 반환합니다.
func (s *AuthService) Login(ctx context.Context, email, password string) (*TokenPair, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil || user == nil {
		// 보안: 이메일 존재 여부를 노출하지 않음
		return nil, apperrors.New(apperrors.ErrUnauthorized, "이메일 또는 비밀번호가 올바르지 않습니다")
	}

	if !user.IsActive {
		return nil, apperrors.New(apperrors.ErrForbidden, "비활성화된 계정입니다")
	}

	// 비밀번호 검증
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return nil, apperrors.New(apperrors.ErrUnauthorized, "이메일 또는 비밀번호가 올바르지 않습니다")
	}

	// 토큰 쌍 생성
	tokenPair, err := s.generateTokenPair(ctx, user)
	if err != nil {
		return nil, err
	}

	s.logger.Info("로그인 성공", zap.String("user_id", user.ID))
	return tokenPair, nil
}

// RefreshToken은 Refresh Token으로 새 토큰 쌍을 발급합니다.
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	// Refresh Token 파싱 및 검증
	claims, err := s.parseToken(refreshToken)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrInvalidToken, "유효하지 않은 리프레시 토큰")
	}

	// Redis에서 토큰 유효성 확인 (revoke 여부)
	valid, err := s.tokenRepo.ValidateRefreshToken(ctx, claims.UserID, claims.ID)
	if err != nil || !valid {
		return nil, apperrors.New(apperrors.ErrTokenExpired, "만료되거나 철회된 리프레시 토큰")
	}

	// 기존 Refresh Token 철회 (Rotation)
	_ = s.tokenRepo.RevokeRefreshToken(ctx, claims.UserID, claims.ID)

	// 사용자 정보 조회
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil || user == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "사용자를 찾을 수 없습니다")
	}

	// 새 토큰 쌍 생성
	return s.generateTokenPair(ctx, user)
}

// Logout은 사용자의 모든 토큰을 철회합니다.
func (s *AuthService) Logout(ctx context.Context, userID string) error {
	if err := s.tokenRepo.RevokeAllUserTokens(ctx, userID); err != nil {
		s.logger.Error("토큰 철회 실패", zap.String("user_id", userID), zap.Error(err))
		return apperrors.New(apperrors.ErrInternal, "로그아웃에 실패했습니다")
	}

	s.logger.Info("로그아웃 완료", zap.String("user_id", userID))
	return nil
}

// ValidateToken은 Access Token을 검증하고 클레임을 반환합니다.
func (s *AuthService) ValidateToken(token string) (*CustomClaims, error) {
	return s.parseToken(token)
}

// generateTokenPair는 Access + Refresh 토큰 쌍을 생성합니다.
func (s *AuthService) generateTokenPair(ctx context.Context, user *User) (*TokenPair, error) {
	now := time.Now().UTC()
	tokenID := uuid.New().String()

	// Access Token (15분)
	accessClaims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   user.ID,
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        tokenID,
		},
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenStr, err := accessToken.SignedString(s.jwtSecret)
	if err != nil {
		s.logger.Error("Access Token 서명 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "토큰 생성에 실패했습니다")
	}

	// Refresh Token (7일)
	refreshID := uuid.New().String()
	refreshClaims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   user.ID,
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        refreshID,
		},
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenStr, err := refreshToken.SignedString(s.jwtSecret)
	if err != nil {
		s.logger.Error("Refresh Token 서명 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "토큰 생성에 실패했습니다")
	}

	// Refresh Token을 Redis에 저장
	if err := s.tokenRepo.StoreRefreshToken(ctx, user.ID, refreshID, s.refreshTTL); err != nil {
		s.logger.Error("Refresh Token 저장 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "토큰 생성에 실패했습니다")
	}

	return &TokenPair{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
		ExpiresIn:    int64(s.accessTTL.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// parseToken은 JWT 토큰을 파싱하고 검증합니다.
func (s *AuthService) parseToken(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		// 서명 알고리즘 검증 (HS256만 허용)
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("예상하지 못한 서명 알고리즘: %v", t.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("유효하지 않은 토큰 클레임")
	}

	return claims, nil
}

// GenerateSecureRandom은 암호학적으로 안전한 랜덤 문자열을 생성합니다.
func GenerateSecureRandom(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
