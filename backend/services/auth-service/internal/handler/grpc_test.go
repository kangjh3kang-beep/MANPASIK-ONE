package handler

import (
	"context"
	"testing"
	"time"

	"github.com/manpasik/backend/services/auth-service/internal/service"
	v1 "github.com/manpasik/backend/shared/gen/go/v1"
	"go.uber.org/zap"
)

type testUserRepo struct {
	users map[string]*service.User
}

func newTestUserRepo() *testUserRepo {
	return &testUserRepo{users: make(map[string]*service.User)}
}

func (r *testUserRepo) GetByID(_ context.Context, id string) (*service.User, error) {
	for _, user := range r.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, nil
}

func (r *testUserRepo) GetByEmail(_ context.Context, email string) (*service.User, error) {
	return r.users[email], nil
}

func (r *testUserRepo) Create(_ context.Context, user *service.User) error {
	r.users[user.Email] = user
	return nil
}

func (r *testUserRepo) UpdatePassword(_ context.Context, id string, hashedPassword string) error {
	for _, user := range r.users {
		if user.ID == id {
			user.HashedPassword = hashedPassword
			return nil
		}
	}
	return nil
}

type testTokenRepo struct {
	tokens map[string]bool
}

func newTestTokenRepo() *testTokenRepo {
	return &testTokenRepo{tokens: make(map[string]bool)}
}

func (r *testTokenRepo) StoreRefreshToken(_ context.Context, userID, tokenID string, _ time.Duration) error {
	r.tokens[userID+":"+tokenID] = true
	return nil
}

func (r *testTokenRepo) ValidateRefreshToken(_ context.Context, userID, tokenID string) (bool, error) {
	return r.tokens[userID+":"+tokenID], nil
}

func (r *testTokenRepo) RevokeRefreshToken(_ context.Context, userID, tokenID string) error {
	delete(r.tokens, userID+":"+tokenID)
	return nil
}

func (r *testTokenRepo) RevokeAllUserTokens(_ context.Context, userID string) error {
	for key := range r.tokens {
		if len(key) > len(userID) && key[:len(userID)] == userID {
			delete(r.tokens, key)
		}
	}
	return nil
}

func TestLoginAndRefreshResponsesIncludeUserID(t *testing.T) {
	ctx := context.Background()
	auth := service.NewAuthService(
		zap.NewNop(),
		newTestUserRepo(),
		newTestTokenRepo(),
		"test-jwt-secret-key-32bytes-long!",
		15*time.Minute,
		7*24*time.Hour,
		"test-issuer",
	)
	handler := NewAuthHandler(auth, zap.NewNop())

	user, err := auth.Register(ctx, "handler-login@manpasik.com", "Pass123!", "핸들러")
	if err != nil {
		t.Fatalf("회원가입 실패: %v", err)
	}

	loginResp, err := handler.Login(ctx, &v1.LoginRequest{
		Email:    "handler-login@manpasik.com",
		Password: "Pass123!",
	})
	if err != nil {
		t.Fatalf("로그인 핸들러 실패: %v", err)
	}
	if loginResp.UserId != user.ID {
		t.Fatalf("로그인 응답 사용자 ID 불일치: got %s, want %s", loginResp.UserId, user.ID)
	}

	refreshResp, err := handler.RefreshToken(ctx, &v1.RefreshTokenRequest{
		RefreshToken: loginResp.RefreshToken,
	})
	if err != nil {
		t.Fatalf("토큰 갱신 핸들러 실패: %v", err)
	}
	if refreshResp.UserId != user.ID {
		t.Fatalf("갱신 응답 사용자 ID 불일치: got %s, want %s", refreshResp.UserId, user.ID)
	}
}
