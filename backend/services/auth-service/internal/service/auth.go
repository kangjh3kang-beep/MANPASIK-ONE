// Package service는 auth-service의 비즈니스 로직을 구현합니다.
package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/manpasik/backend/services/auth-service/internal/oauth"
	apperrors "github.com/manpasik/backend/shared/errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// AuthService는 인증 서비스의 비즈니스 로직입니다.
//
// jwtSecret 은 atomic.Pointer 로 보관되어 런타임 중 SetJWTSecret 으로 핫리로드
// 가능 (Vault 시크릿 회전 대응). 토큰 서명/검증 시 항상 최신 값 사용.
type AuthService struct {
	logger         *zap.Logger
	userRepo       UserRepository
	tokenRepo      TokenRepository
	jwtSecret      atomic.Pointer[[]byte]
	accessTTL      time.Duration
	refreshTTL     time.Duration
	issuer         string
	oauthVerifiers map[string]oauth.OAuthVerifier
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
	UserID       string
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
	s := &AuthService{
		logger:     logger,
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		issuer:     issuer,
	}
	secret := []byte(jwtSecret)
	s.jwtSecret.Store(&secret)
	return s
}

// SetJWTSecret 은 런타임 중 JWT 서명 키를 교체합니다 (Vault 회전 대응).
//
// 호출 시점 이후 발급되는 새 토큰은 새 키로 서명되며, 검증 시에도 새 키로
// 처리됨. 기존 발급된 토큰은 만료 시까지는 검증 실패할 수 있으므로
// 호출 전에 grace period 또는 token rotation 정책을 함께 운영해야 함.
func (s *AuthService) SetJWTSecret(newSecret string) {
	if newSecret == "" {
		return
	}
	secret := []byte(newSecret)
	s.jwtSecret.Store(&secret)
}

// loadJWTSecret 은 atomic.Pointer 에서 현재 시크릿을 읽어 반환.
func (s *AuthService) loadJWTSecret() []byte {
	p := s.jwtSecret.Load()
	if p == nil {
		return nil
	}
	return *p
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
	accessTokenStr, err := accessToken.SignedString(s.loadJWTSecret())
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
	refreshTokenStr, err := refreshToken.SignedString(s.loadJWTSecret())
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
		UserID:       user.ID,
	}, nil
}

// parseToken은 JWT 토큰을 파싱하고 검증합니다.
func (s *AuthService) parseToken(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		// 서명 알고리즘 검증 (HS256만 허용)
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("예상하지 못한 서명 알고리즘: %v", t.Header["alg"])
		}
		return s.loadJWTSecret(), nil
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

// SocialLogin은 OAuth2 소셜 로그인을 처리합니다.
// 소셜 제공자의 ID 토큰을 검증하고, 기존 사용자면 로그인/신규면 자동 가입합니다.
func (s *AuthService) SocialLogin(ctx context.Context, provider, idToken, accessToken string) (*TokenPair, error) {
	if provider == "" || (idToken == "" && accessToken == "") {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "provider와 토큰은 필수입니다")
	}

	// 소셜 제공자별 이메일 추출 (실제 환경에서는 각 OAuth2 API 호출)
	email := fmt.Sprintf("%s_%s@social.manpasik.com", provider, idToken[:min(8, len(idToken))])

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil || user == nil {
		// 신규 사용자 자동 가입
		hashedPw, _ := bcrypt.GenerateFromPassword([]byte(uuid.New().String()), 10)
		user = &User{
			ID:             uuid.New().String(),
			Email:          email,
			HashedPassword: string(hashedPw),
			DisplayName:    provider + " 사용자",
			Role:           "user",
			IsActive:       true,
			CreatedAt:      time.Now().UTC(),
			UpdatedAt:      time.Now().UTC(),
		}
		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, apperrors.New(apperrors.ErrInternal, "소셜 로그인 사용자 생성 실패")
		}
		s.logger.Info("소셜 로그인 신규 가입", zap.String("provider", provider), zap.String("user_id", user.ID))
	}

	tokenPair, err := s.generateTokenPair(ctx, user)
	if err != nil {
		return nil, err
	}

	s.logger.Info("소셜 로그인 성공", zap.String("provider", provider), zap.String("user_id", user.ID))
	return tokenPair, nil
}

// KakaoUserInfo는 카카오 사용자 정보 API 응답에서 추출한 데이터입니다.
type KakaoUserInfo struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
}

// KakaoLogin은 카카오 Access Token을 검증하고 만파식 JWT를 발급합니다.
//
// 1. 카카오 사용자 정보 API(kapi.kakao.com/v2/user/me) 호출로 토큰 검증
// 2. 카카오 ID 기반 기존 사용자 조회 → 없으면 자동 가입
// 3. 만파식 JWT(Access/Refresh Token) 발급
func (s *AuthService) KakaoLogin(ctx context.Context, kakaoAccessToken, email, nickname string) (*TokenPair, *User, error) {
	if kakaoAccessToken == "" {
		return nil, nil, apperrors.New(apperrors.ErrInvalidInput, "kakao_access_token은 필수입니다")
	}

	// 1단계: 카카오 사용자 정보 API 호출로 토큰 검증
	kakaoUser, err := s.verifyKakaoToken(ctx, kakaoAccessToken)
	if err != nil {
		s.logger.Error("카카오 토큰 검증 실패", zap.Error(err))
		return nil, nil, apperrors.New(apperrors.ErrUnauthorized, "카카오 인증에 실패했습니다")
	}

	s.logger.Info("카카오 사용자 검증 완료",
		zap.Int64("kakao_id", kakaoUser.ID),
		zap.String("email", kakaoUser.Email),
	)

	// 2단계: 카카오 ID 기반 이메일로 기존 사용자 조회
	kakaoEmail := kakaoUser.Email
	if kakaoEmail == "" {
		kakaoEmail = fmt.Sprintf("kakao_%d@social.manpasik.com", kakaoUser.ID)
	}
	displayName := kakaoUser.Nickname
	if displayName == "" {
		displayName = nickname
	}
	if displayName == "" {
		displayName = "카카오 사용자"
	}

	user, err := s.userRepo.GetByEmail(ctx, kakaoEmail)
	if err != nil || user == nil {
		// 3단계: 신규 사용자 자동 가입
		hashedPw, _ := bcrypt.GenerateFromPassword([]byte(uuid.New().String()), 10)
		user = &User{
			ID:             uuid.New().String(),
			Email:          kakaoEmail,
			HashedPassword: string(hashedPw),
			DisplayName:    displayName,
			Role:           "user",
			IsActive:       true,
			CreatedAt:      time.Now().UTC(),
			UpdatedAt:      time.Now().UTC(),
		}
		if createErr := s.userRepo.Create(ctx, user); createErr != nil {
			s.logger.Error("카카오 사용자 생성 실패", zap.Error(createErr))
			return nil, nil, apperrors.New(apperrors.ErrInternal, "카카오 로그인 사용자 생성 실패")
		}
		s.logger.Info("카카오 신규 가입 완료",
			zap.String("user_id", user.ID),
			zap.Int64("kakao_id", kakaoUser.ID),
		)
	}

	// 4단계: 만파식 JWT 발급
	tokenPair, err := s.generateTokenPair(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	s.logger.Info("카카오 로그인 성공",
		zap.String("user_id", user.ID),
		zap.Int64("kakao_id", kakaoUser.ID),
	)
	return tokenPair, user, nil
}

// verifyKakaoToken은 카카오 Access Token으로 사용자 정보 API를 호출하여
// 토큰 유효성을 검증하고 사용자 고유 번호(ID)를 반환합니다.
func (s *AuthService) verifyKakaoToken(ctx context.Context, accessToken string) (*KakaoUserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://kapi.kakao.com/v2/user/me", nil)
	if err != nil {
		return nil, fmt.Errorf("카카오 API 요청 생성 실패: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("카카오 API 호출 실패: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("카카오 API 응답 오류 (HTTP %d): %s", resp.StatusCode, string(body))
	}

	var kakaoResp struct {
		ID           int64 `json:"id"`
		KakaoAccount struct {
			Email   string `json:"email"`
			Profile struct {
				Nickname string `json:"nickname"`
			} `json:"profile"`
		} `json:"kakao_account"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&kakaoResp); err != nil {
		return nil, fmt.Errorf("카카오 API 응답 파싱 실패: %w", err)
	}

	if kakaoResp.ID == 0 {
		return nil, fmt.Errorf("카카오 사용자 ID가 0입니다 — 유효하지 않은 토큰")
	}

	return &KakaoUserInfo{
		ID:       kakaoResp.ID,
		Email:    kakaoResp.KakaoAccount.Email,
		Nickname: kakaoResp.KakaoAccount.Profile.Nickname,
	}, nil
}

// ResetPassword는 비밀번호 재설정 요청을 처리합니다.
// 이메일로 재설정 토큰을 발송합니다 (현재는 성공 응답만 반환).
func (s *AuthService) ResetPassword(ctx context.Context, email string) (bool, string, error) {
	if email == "" {
		return false, "", apperrors.New(apperrors.ErrInvalidInput, "이메일은 필수입니다")
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil || user == nil {
		// 보안: 이메일 존재 여부를 노출하지 않음
		return true, "비밀번호 재설정 링크가 이메일로 전송되었습니다.", nil
	}

	// 재설정 토큰 생성
	resetToken, err := GenerateSecureRandom(32)
	if err != nil {
		return false, "", apperrors.New(apperrors.ErrInternal, "토큰 생성에 실패했습니다")
	}

	// 토큰 저장 (30분 TTL)
	_ = s.tokenRepo.StoreRefreshToken(ctx, user.ID, "reset:"+resetToken, 30*time.Minute)

	s.logger.Info("비밀번호 재설정 요청", zap.String("email", email), zap.String("user_id", user.ID))
	return true, "비밀번호 재설정 링크가 이메일로 전송되었습니다.", nil
}

// GenerateSecureRandom은 암호학적으로 안전한 랜덤 문자열을 생성합니다.
func GenerateSecureRandom(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// SetOAuthVerifiers는 OAuth 제공자별 토큰 검증기를 설정합니다.
func (s *AuthService) SetOAuthVerifiers(verifiers map[string]oauth.OAuthVerifier) {
	s.oauthVerifiers = verifiers
}

// ProviderLogin은 OAuth 제공자 토큰을 검증하고 JWT를 발급합니다.
// 등록된 검증기가 없으면 기존 SocialLogin으로 폴백합니다.
func (s *AuthService) ProviderLogin(ctx context.Context, provider, idToken, accessToken string) (*TokenPair, *User, error) {
	if provider == "" {
		return nil, nil, apperrors.New(apperrors.ErrInvalidInput, "provider는 필수입니다")
	}
	if idToken == "" && accessToken == "" {
		return nil, nil, apperrors.New(apperrors.ErrInvalidInput, "토큰은 필수입니다")
	}

	// 등록된 검증기 조회
	verifier, ok := s.oauthVerifiers[provider]
	if !ok {
		// 폴백: 기존 SocialLogin 방식
		tokenPair, err := s.SocialLogin(ctx, provider, idToken, accessToken)
		if err != nil {
			return nil, nil, err
		}
		return tokenPair, nil, nil
	}

	// OAuth 토큰 검증
	userInfo, err := verifier.VerifyToken(ctx, idToken, accessToken)
	if err != nil {
		s.logger.Error("OAuth 토큰 검증 실패",
			zap.String("provider", provider),
			zap.Error(err),
		)
		return nil, nil, apperrors.New(apperrors.ErrUnauthorized, fmt.Sprintf("%s 인증에 실패했습니다", provider))
	}

	// 이메일로 기존 사용자 조회
	user, err := s.userRepo.GetByEmail(ctx, userInfo.Email)
	if err != nil || user == nil {
		// 신규 사용자 자동 가입
		hashedPw, _ := bcrypt.GenerateFromPassword([]byte(uuid.New().String()), 10)
		displayName := userInfo.DisplayName
		if displayName == "" {
			displayName = provider + " 사용자"
		}
		user = &User{
			ID:             uuid.New().String(),
			Email:          userInfo.Email,
			HashedPassword: string(hashedPw),
			DisplayName:    displayName,
			Role:           "user",
			IsActive:       true,
			CreatedAt:      time.Now().UTC(),
			UpdatedAt:      time.Now().UTC(),
		}
		if createErr := s.userRepo.Create(ctx, user); createErr != nil {
			s.logger.Error("OAuth 사용자 생성 실패", zap.Error(createErr))
			return nil, nil, apperrors.New(apperrors.ErrInternal, "사용자 생성 실패")
		}
		s.logger.Info("OAuth 신규 가입",
			zap.String("provider", provider),
			zap.String("user_id", user.ID),
		)
	}

	// JWT 토큰 쌍 발급
	tokenPair, err := s.generateTokenPair(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	s.logger.Info("OAuth 로그인 성공",
		zap.String("provider", provider),
		zap.String("user_id", user.ID),
	)
	return tokenPair, user, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
