package oauth

import (
	"context"
	"testing"
)

func TestNoopVerifier_VerifyToken_Success(t *testing.T) {
	v := NewNoopVerifier("test")
	info, err := v.VerifyToken(context.Background(), "mytoken123", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Provider != "test" {
		t.Errorf("expected provider test, got %s", info.Provider)
	}
	if info.Email != "test_mytoken1@social.manpasik.com" {
		t.Errorf("unexpected email: %s", info.Email)
	}
}

func TestNoopVerifier_VerifyToken_EmptyToken(t *testing.T) {
	v := NewNoopVerifier("test")
	_, err := v.VerifyToken(context.Background(), "", "")
	if err == nil {
		t.Fatal("NoopVerifier should return error for empty tokens")
	}
}

func TestNoopVerifier_Provider(t *testing.T) {
	v := NewNoopVerifier("google")
	if v.Provider() != "google" {
		t.Errorf("expected google, got %s", v.Provider())
	}
}

func TestOAuthUserInfo_Fields(t *testing.T) {
	info := &OAuthUserInfo{
		ProviderID:   "123",
		Email:        "test@example.com",
		DisplayName:  "Test",
		ProfileImage: "https://img.test/photo.jpg",
		Provider:     "google",
	}
	if info.Email != "test@example.com" {
		t.Errorf("email mismatch")
	}
	if info.Provider != "google" {
		t.Errorf("provider mismatch")
	}
}
