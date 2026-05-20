package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/manpasik/backend/services/auth-service/internal/oauth"
	"go.uber.org/zap"
)

// =============================================================================
// 목(Mock) 저장소
// =============================================================================

type mockUserRepo struct {
	users map[string]*User // email → User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*User)}
}

func (m *mockUserRepo) GetByID(ctx context.Context, id string) (*User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*User, error) {
	if u, ok := m.users[email]; ok {
		return u, nil
	}
	return nil, nil
}

func (m *mockUserRepo) Create(ctx context.Context, user *User) error {
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepo) UpdatePassword(ctx context.Context, id string, hashedPassword string) error {
	for _, u := range m.users {
		if u.ID == id {
			u.HashedPassword = hashedPassword
			return nil
		}
	}
	return nil
}

type mockTokenRepo struct {
	tokens map[string]bool // "userID:tokenID" → valid
}

func newMockTokenRepo() *mockTokenRepo {
	return &mockTokenRepo{tokens: make(map[string]bool)}
}

func (m *mockTokenRepo) StoreRefreshToken(ctx context.Context, userID, tokenID string, ttl time.Duration) error {
	m.tokens[userID+":"+tokenID] = true
	return nil
}

func (m *mockTokenRepo) ValidateRefreshToken(ctx context.Context, userID, tokenID string) (bool, error) {
	valid, ok := m.tokens[userID+":"+tokenID]
	return ok && valid, nil
}

func (m *mockTokenRepo) RevokeRefreshToken(ctx context.Context, userID, tokenID string) error {
	delete(m.tokens, userID+":"+tokenID)
	return nil
}

func (m *mockTokenRepo) RevokeAllUserTokens(ctx context.Context, userID string) error {
	for key := range m.tokens {
		if len(key) > len(userID) && key[:len(userID)] == userID {
			delete(m.tokens, key)
		}
	}
	return nil
}

// =============================================================================
// 헬퍼
// =============================================================================

func newTestService() *AuthService {
	logger, _ := zap.NewDevelopment()
	if logger == nil {
		logger = zap.NewNop()
	}
	return NewAuthService(
		logger,
		newMockUserRepo(),
		newMockTokenRepo(),
		"test-jwt-secret-key-32bytes-long!",
		15*time.Minute,
		7*24*time.Hour,
		"test-issuer",
	)
}

// =============================================================================
// 테스트
// =============================================================================

func TestRegister_성공(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	user, err := svc.Register(ctx, "test@manpasik.com", "SecurePass123!", "테스트 사용자")
	if err != nil {
		t.Fatalf("회원가입 실패: %v", err)
	}

	if user.Email != "test@manpasik.com" {
		t.Errorf("이메일 불일치: got %s, want test@manpasik.com", user.Email)
	}

	if user.DisplayName != "테스트 사용자" {
		t.Errorf("이름 불일치: got %s, want 테스트 사용자", user.DisplayName)
	}

	if user.Role != "user" {
		t.Errorf("역할 불일치: got %s, want user", user.Role)
	}

	if user.ID == "" {
		t.Error("사용자 ID가 비어있습니다")
	}
}

func TestRegister_중복_이메일(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	// 첫 번째 등록
	_, err := svc.Register(ctx, "dup@manpasik.com", "Pass123!", "사용자1")
	if err != nil {
		t.Fatalf("첫 등록 실패: %v", err)
	}

	// 중복 등록 시도
	_, err = svc.Register(ctx, "dup@manpasik.com", "Pass456!", "사용자2")
	if err == nil {
		t.Error("중복 이메일 등록이 허용됨")
	}
}

func TestLogin_성공(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	// 회원가입
	user, err := svc.Register(ctx, "login@manpasik.com", "MyPassword123!", "로그인 테스트")
	if err != nil {
		t.Fatalf("회원가입 실패: %v", err)
	}

	// 로그인
	tokens, err := svc.Login(ctx, "login@manpasik.com", "MyPassword123!")
	if err != nil {
		t.Fatalf("로그인 실패: %v", err)
	}

	if tokens.AccessToken == "" {
		t.Error("Access Token이 비어있습니다")
	}
	if tokens.RefreshToken == "" {
		t.Error("Refresh Token이 비어있습니다")
	}
	if tokens.UserID != user.ID {
		t.Errorf("사용자 ID 불일치: got %s, want %s", tokens.UserID, user.ID)
	}
	if tokens.TokenType != "Bearer" {
		t.Errorf("토큰 타입 불일치: got %s, want Bearer", tokens.TokenType)
	}
	if tokens.ExpiresIn != 900 { // 15분 = 900초
		t.Errorf("만료 시간 불일치: got %d, want 900", tokens.ExpiresIn)
	}
}

func TestLogin_잘못된_비밀번호(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, _ = svc.Register(ctx, "wrong@manpasik.com", "CorrectPass!", "테스트")

	_, err := svc.Login(ctx, "wrong@manpasik.com", "WrongPass!")
	if err == nil {
		t.Error("잘못된 비밀번호로 로그인 성공됨")
	}
}

func TestLogin_존재하지_않는_이메일(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, err := svc.Login(ctx, "noone@manpasik.com", "AnyPass!")
	if err == nil {
		t.Error("존재하지 않는 이메일로 로그인 성공됨")
	}
}

func TestValidateToken_성공(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, _ = svc.Register(ctx, "validate@manpasik.com", "Pass123!", "검증 테스트")
	tokens, _ := svc.Login(ctx, "validate@manpasik.com", "Pass123!")

	claims, err := svc.ValidateToken(tokens.AccessToken)
	if err != nil {
		t.Fatalf("토큰 검증 실패: %v", err)
	}

	if claims.Email != "validate@manpasik.com" {
		t.Errorf("이메일 불일치: got %s", claims.Email)
	}
	if claims.Role != "user" {
		t.Errorf("역할 불일치: got %s", claims.Role)
	}
}

func TestValidateToken_잘못된_토큰(t *testing.T) {
	svc := newTestService()

	_, err := svc.ValidateToken("invalid.token.here")
	if err == nil {
		t.Error("잘못된 토큰이 검증 통과됨")
	}
}

func TestGenerateSecureRandom(t *testing.T) {
	random1, err := GenerateSecureRandom(32)
	if err != nil {
		t.Fatalf("랜덤 생성 실패: %v", err)
	}

	random2, _ := GenerateSecureRandom(32)

	if random1 == random2 {
		t.Error("두 랜덤 값이 동일합니다")
	}

	if len(random1) == 0 {
		t.Error("랜덤 값이 비어있습니다")
	}
}

// =============================================================================
// Phase F 테스트 보강
// =============================================================================

func TestRefreshToken_성공(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	user, _ := svc.Register(ctx, "refresh@manpasik.com", "Pass123!", "리프레시 테스트")
	tokens, _ := svc.Login(ctx, "refresh@manpasik.com", "Pass123!")

	newTokens, err := svc.RefreshToken(ctx, tokens.RefreshToken)
	if err != nil {
		t.Fatalf("리프레시 토큰 갱신 실패: %v", err)
	}
	if newTokens.AccessToken == "" {
		t.Error("새 Access Token이 비어있습니다")
	}
	if newTokens.RefreshToken == "" {
		t.Error("새 Refresh Token이 비어있습니다")
	}
	if newTokens.UserID != user.ID {
		t.Errorf("갱신 토큰 사용자 ID 불일치: got %s, want %s", newTokens.UserID, user.ID)
	}
	if newTokens.TokenType != "Bearer" {
		t.Errorf("토큰 타입 불일치: got %s", newTokens.TokenType)
	}
}

func TestRefreshToken_잘못된_토큰(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, err := svc.RefreshToken(ctx, "invalid.refresh.token")
	if err == nil {
		t.Fatal("잘못된 리프레시 토큰으로 갱신이 성공됨")
	}
}

func TestRefreshToken_철회된_토큰(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, _ = svc.Register(ctx, "revoked@manpasik.com", "Pass123!", "철회 테스트")
	tokens, _ := svc.Login(ctx, "revoked@manpasik.com", "Pass123!")

	// 첫 갱신 성공 (기존 토큰 철회됨)
	_, err := svc.RefreshToken(ctx, tokens.RefreshToken)
	if err != nil {
		t.Fatalf("첫 갱신 실패: %v", err)
	}

	// 같은 토큰으로 재갱신 시도 — 이미 철회됨
	_, err = svc.RefreshToken(ctx, tokens.RefreshToken)
	if err == nil {
		t.Fatal("철회된 리프레시 토큰으로 갱신이 성공됨")
	}
}

func TestLogout_성공(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	user, _ := svc.Register(ctx, "logout@manpasik.com", "Pass123!", "로그아웃 테스트")
	_, _ = svc.Login(ctx, "logout@manpasik.com", "Pass123!")

	err := svc.Logout(ctx, user.ID)
	if err != nil {
		t.Fatalf("로그아웃 실패: %v", err)
	}
}

func TestLogin_비활성_계정(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	// 사용자 등록 후 비활성화
	user, _ := svc.Register(ctx, "inactive@manpasik.com", "Pass123!", "비활성 테스트")
	user.IsActive = false

	_, err := svc.Login(ctx, "inactive@manpasik.com", "Pass123!")
	if err == nil {
		t.Fatal("비활성 계정으로 로그인이 성공됨")
	}
}

func TestSocialLogin_신규_사용자(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	tokens, err := svc.SocialLogin(ctx, "kakao", "test-id-token-12345678", "test-access-token")
	if err != nil {
		t.Fatalf("소셜 로그인 실패: %v", err)
	}
	if tokens.AccessToken == "" {
		t.Error("Access Token이 비어있습니다")
	}
	if tokens.UserID == "" {
		t.Error("소셜 로그인 사용자 ID가 비어있습니다")
	}
	if tokens.TokenType != "Bearer" {
		t.Errorf("토큰 타입 불일치: got %s", tokens.TokenType)
	}
}

func TestSocialLogin_기존_사용자(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	// 첫 소셜 로그인 — 자동 가입
	tokens1, _ := svc.SocialLogin(ctx, "kakao", "same-token-12345678", "access")
	// 같은 토큰으로 재로그인
	tokens2, err := svc.SocialLogin(ctx, "kakao", "same-token-12345678", "access")
	if err != nil {
		t.Fatalf("기존 사용자 소셜 로그인 실패: %v", err)
	}
	if tokens1.AccessToken == tokens2.AccessToken {
		t.Error("재로그인 시 새 토큰이 발급되어야 합니다")
	}
}

func TestSocialLogin_빈_제공자(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, err := svc.SocialLogin(ctx, "", "token123", "access")
	if err == nil {
		t.Fatal("빈 provider에 에러가 반환되어야 함")
	}
}

func TestSocialLogin_빈_토큰(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, err := svc.SocialLogin(ctx, "kakao", "", "")
	if err == nil {
		t.Fatal("빈 토큰에 에러가 반환되어야 함")
	}
}

func TestResetPassword_존재하는_이메일(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, _ = svc.Register(ctx, "reset@manpasik.com", "Pass123!", "리셋 테스트")

	success, msg, err := svc.ResetPassword(ctx, "reset@manpasik.com")
	if err != nil {
		t.Fatalf("비밀번호 재설정 실패: %v", err)
	}
	if !success {
		t.Error("성공이어야 합니다")
	}
	if msg == "" {
		t.Error("메시지가 비어있습니다")
	}
}

func TestResetPassword_존재하지_않는_이메일(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	// 보안: 존재하지 않는 이메일도 성공 응답
	success, msg, err := svc.ResetPassword(ctx, "noone@manpasik.com")
	if err != nil {
		t.Fatalf("비밀번호 재설정 실패: %v", err)
	}
	if !success {
		t.Error("보안상 성공이어야 합니다")
	}
	if msg == "" {
		t.Error("메시지가 비어있습니다")
	}
}

func TestResetPassword_빈_이메일(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, _, err := svc.ResetPassword(ctx, "")
	if err == nil {
		t.Fatal("빈 이메일에 에러가 반환되어야 함")
	}
}

// =============================================================================
// Phase F 추가 보강 — Failing Mock + 에러 경로 커버리지
// =============================================================================

type failingCreateUserRepo struct {
	*mockUserRepo
}

func (m *failingCreateUserRepo) Create(_ context.Context, _ *User) error {
	return fmt.Errorf("DB 연결 실패")
}

type failingRevokeTokenRepo struct {
	*mockTokenRepo
}

func (m *failingRevokeTokenRepo) RevokeAllUserTokens(_ context.Context, _ string) error {
	return fmt.Errorf("Redis 연결 실패")
}

type failingStoreTokenRepo struct {
	*mockTokenRepo
}

func (m *failingStoreTokenRepo) StoreRefreshToken(_ context.Context, _, _ string, _ time.Duration) error {
	return fmt.Errorf("Redis 저장 실패")
}

func TestKakaoLogin_빈_토큰(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, _, err := svc.KakaoLogin(ctx, "", "", "")
	if err == nil {
		t.Fatal("빈 카카오 토큰에 에러가 반환되어야 함")
	}
}

func TestRegister_저장소_실패(t *testing.T) {
	svc := NewAuthService(
		zap.NewNop(),
		&failingCreateUserRepo{mockUserRepo: newMockUserRepo()},
		newMockTokenRepo(),
		"test-jwt-secret-key-32bytes-long!",
		15*time.Minute, 7*24*time.Hour, "test-issuer",
	)
	ctx := context.Background()

	_, err := svc.Register(ctx, "fail@test.com", "Pass123!", "실패 테스트")
	if err == nil {
		t.Fatal("저장소 실패 시 에러가 반환되어야 함")
	}
}

func TestLogout_철회_실패(t *testing.T) {
	svc := NewAuthService(
		zap.NewNop(),
		newMockUserRepo(),
		&failingRevokeTokenRepo{mockTokenRepo: newMockTokenRepo()},
		"test-jwt-secret-key-32bytes-long!",
		15*time.Minute, 7*24*time.Hour, "test-issuer",
	)
	ctx := context.Background()

	err := svc.Logout(ctx, "any-user-id")
	if err == nil {
		t.Fatal("토큰 철회 실패 시 에러가 반환되어야 함")
	}
}

func TestLogin_토큰_저장_실패(t *testing.T) {
	userRepo := newMockUserRepo()
	svc := NewAuthService(
		zap.NewNop(),
		userRepo,
		&failingStoreTokenRepo{mockTokenRepo: newMockTokenRepo()},
		"test-jwt-secret-key-32bytes-long!",
		15*time.Minute, 7*24*time.Hour, "test-issuer",
	)
	ctx := context.Background()

	// Register는 tokenRepo를 사용하지 않으므로 성공 (bcrypt만 사용)
	_, err := svc.Register(ctx, "store-fail@test.com", "Pass123!", "저장 실패")
	if err != nil {
		t.Fatalf("등록 실패: %v", err)
	}

	// Login → generateTokenPair → StoreRefreshToken 실패
	_, err = svc.Login(ctx, "store-fail@test.com", "Pass123!")
	if err == nil {
		t.Fatal("토큰 저장 실패 시 로그인 에러가 반환되어야 함")
	}
}

func TestRefreshToken_사용자_삭제됨(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	user, _ := svc.Register(ctx, "deleted@manpasik.com", "Pass123!", "삭제 테스트")
	tokens, _ := svc.Login(ctx, "deleted@manpasik.com", "Pass123!")

	// mock에서 사용자 직접 삭제
	userRepo := svc.userRepo.(*mockUserRepo)
	delete(userRepo.users, user.Email)

	_, err := svc.RefreshToken(ctx, tokens.RefreshToken)
	if err == nil {
		t.Fatal("삭제된 사용자의 토큰 갱신은 실패해야 함")
	}
}

func TestSocialLogin_짧은_토큰(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	// idToken이 8자보다 짧은 경우 (min() 경로)
	tokens, err := svc.SocialLogin(ctx, "google", "abc", "access-token")
	if err != nil {
		t.Fatalf("짧은 idToken 소셜 로그인 실패: %v", err)
	}
	if tokens.AccessToken == "" {
		t.Error("Access Token이 비어있음")
	}
}

func TestSocialLogin_생성_실패(t *testing.T) {
	svc := NewAuthService(
		zap.NewNop(),
		&failingCreateUserRepo{mockUserRepo: newMockUserRepo()},
		newMockTokenRepo(),
		"test-jwt-secret-key-32bytes-long!",
		15*time.Minute, 7*24*time.Hour, "test-issuer",
	)
	ctx := context.Background()

	_, err := svc.SocialLogin(ctx, "kakao", "new-user-token-123", "access")
	if err == nil {
		t.Fatal("사용자 생성 실패 시 에러가 반환되어야 함")
	}
}

// =============================================================================
// Phase C-1: ProviderLogin 테스트
// =============================================================================

// mockFailVerifier는 항상 에러를 반환하는 OAuthVerifier입니다.
type mockFailVerifier struct {
	provider string
}

func (m *mockFailVerifier) VerifyToken(_ context.Context, _, _ string) (*oauth.OAuthUserInfo, error) {
	return nil, fmt.Errorf("인증 실패")
}

func (m *mockFailVerifier) Provider() string {
	return m.provider
}

func TestProviderLogin_NoopVerifier_성공(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	verifiers := map[string]oauth.OAuthVerifier{
		"google": oauth.NewNoopVerifier("google"),
	}
	svc.SetOAuthVerifiers(verifiers)

	tokens, user, err := svc.ProviderLogin(ctx, "google", "test-token-12345678", "")
	if err != nil {
		t.Fatalf("ProviderLogin 실패: %v", err)
	}
	if tokens.AccessToken == "" {
		t.Error("Access Token이 비어있습니다")
	}
	if user == nil {
		t.Error("사용자 정보가 nil입니다")
	}
}

func TestProviderLogin_미등록_제공자_폴백(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	// 빈 verifier map → SocialLogin 폴백
	svc.SetOAuthVerifiers(map[string]oauth.OAuthVerifier{})

	tokens, _, err := svc.ProviderLogin(ctx, "kakao", "test-token-12345678", "access")
	if err != nil {
		t.Fatalf("폴백 로그인 실패: %v", err)
	}
	if tokens.AccessToken == "" {
		t.Error("Access Token이 비어있습니다")
	}
}

func TestProviderLogin_빈_제공자(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, _, err := svc.ProviderLogin(ctx, "", "token", "")
	if err == nil {
		t.Fatal("빈 provider에 에러가 반환되어야 함")
	}
}

func TestProviderLogin_빈_토큰(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, _, err := svc.ProviderLogin(ctx, "google", "", "")
	if err == nil {
		t.Fatal("빈 토큰에 에러가 반환되어야 함")
	}
}

func TestProviderLogin_검증_실패(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	verifiers := map[string]oauth.OAuthVerifier{
		"google": &mockFailVerifier{provider: "google"},
	}
	svc.SetOAuthVerifiers(verifiers)

	_, _, err := svc.ProviderLogin(ctx, "google", "bad-token", "")
	if err == nil {
		t.Fatal("검증 실패 시 에러가 반환되어야 함")
	}
}

func TestProviderLogin_사용자_생성_실패(t *testing.T) {
	svc := NewAuthService(
		zap.NewNop(),
		&failingCreateUserRepo{mockUserRepo: newMockUserRepo()},
		newMockTokenRepo(),
		"test-jwt-secret-key-32bytes-long!",
		15*time.Minute, 7*24*time.Hour, "test-issuer",
	)

	verifiers := map[string]oauth.OAuthVerifier{
		"google": oauth.NewNoopVerifier("google"),
	}
	svc.SetOAuthVerifiers(verifiers)
	ctx := context.Background()

	_, _, err := svc.ProviderLogin(ctx, "google", "test-token-12345678", "")
	if err == nil {
		t.Fatal("사용자 생성 실패 시 에러가 반환되어야 함")
	}
}

func TestMin(t *testing.T) {
	if min(3, 5) != 3 {
		t.Error("min(3,5) should be 3")
	}
	if min(7, 2) != 2 {
		t.Error("min(7,2) should be 2")
	}
	if min(4, 4) != 4 {
		t.Error("min(4,4) should be 4")
	}
}

// SetJWTSecret 핫리로드 — 새 시크릿으로 발급한 토큰이 유효한지 검증.
func TestSetJWTSecret_HotReload(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	// 초기 시크릿으로 사용자 등록 + 로그인 → 토큰 발급
	_, err := svc.Register(ctx, "user@test.com", "password123!", "Test")
	if err != nil {
		t.Fatal(err)
	}
	pair1, err := svc.Login(ctx, "user@test.com", "password123!")
	if err != nil {
		t.Fatal(err)
	}

	// 시크릿 회전
	svc.SetJWTSecret("rotated-jwt-secret-key-32bytes-long!")

	// 새 토큰 발급 (회전된 키로 서명)
	pair2, err := svc.Login(ctx, "user@test.com", "password123!")
	if err != nil {
		t.Fatalf("회전 후 로그인 실패: %v", err)
	}

	// 새 토큰은 검증 통과
	if _, err := svc.ValidateToken(pair2.AccessToken); err != nil {
		t.Errorf("새 키 토큰 검증 실패: %v", err)
	}
	// 기존 토큰은 검증 실패해야 함 (다른 키로 서명되었으므로)
	if _, err := svc.ValidateToken(pair1.AccessToken); err == nil {
		t.Error("회전 전 토큰이 새 키에서 통과")
	}
}

func TestSetJWTSecret_EmptyIgnored(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()
	_, _ = svc.Register(ctx, "u@test.com", "password123!", "Test")
	pair, _ := svc.Login(ctx, "u@test.com", "password123!")

	// 빈 시크릿은 무시되어야 함 (기존 토큰이 여전히 유효)
	svc.SetJWTSecret("")
	if _, err := svc.ValidateToken(pair.AccessToken); err != nil {
		t.Errorf("빈 시크릿이 적용되어 토큰 검증 실패: %v", err)
	}
}
