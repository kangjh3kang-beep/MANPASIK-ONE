package service

import (
	"context"
	"testing"
	"time"

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
	_, err := svc.Register(ctx, "login@manpasik.com", "MyPassword123!", "로그인 테스트")
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
