package oauth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGoogleVerifier_VerifyToken_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"sub": "12345",
			"email": "test@gmail.com",
			"name": "Test User",
			"picture": "https://pic.test/photo.jpg",
			"aud": "test-client-id",
			"iss": "accounts.google.com",
			"email_verified": "true"
		}`))
	}))
	defer ts.Close()

	v := NewGoogleVerifier("test-client-id")
	v.SetTokenURL(ts.URL)

	info, err := v.VerifyToken(context.Background(), "valid-id-token", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.ProviderID != "12345" {
		t.Errorf("expected ProviderID=12345, got %s", info.ProviderID)
	}
	if info.Email != "test@gmail.com" {
		t.Errorf("expected email=test@gmail.com, got %s", info.Email)
	}
	if info.DisplayName != "Test User" {
		t.Errorf("expected DisplayName=Test User, got %s", info.DisplayName)
	}
	if info.ProfileImage != "https://pic.test/photo.jpg" {
		t.Errorf("expected ProfileImage=https://pic.test/photo.jpg, got %s", info.ProfileImage)
	}
	if info.Provider != "google" {
		t.Errorf("expected Provider=google, got %s", info.Provider)
	}
}

func TestGoogleVerifier_VerifyToken_AudMismatch(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"sub": "12345",
			"email": "test@gmail.com",
			"name": "Test User",
			"picture": "https://pic.test/photo.jpg",
			"aud": "wrong-client-id",
			"iss": "accounts.google.com",
			"email_verified": "true"
		}`))
	}))
	defer ts.Close()

	v := NewGoogleVerifier("test-client-id")
	v.SetTokenURL(ts.URL)

	_, err := v.VerifyToken(context.Background(), "valid-id-token", "")
	if err == nil {
		t.Fatal("expected error for audience mismatch")
	}
	if !strings.Contains(err.Error(), "audience 불일치") {
		t.Errorf("expected audience mismatch error, got: %v", err)
	}
}

func TestGoogleVerifier_VerifyToken_HTTPError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "invalid_token"}`))
	}))
	defer ts.Close()

	v := NewGoogleVerifier("test-client-id")
	v.SetTokenURL(ts.URL)

	_, err := v.VerifyToken(context.Background(), "bad-token", "")
	if err == nil {
		t.Fatal("expected error for HTTP 401")
	}
	if !strings.Contains(err.Error(), "HTTP 401") {
		t.Errorf("expected HTTP 401 error, got: %v", err)
	}
}

func TestGoogleVerifier_VerifyToken_EmptyToken(t *testing.T) {
	v := NewGoogleVerifier("test-client-id")

	_, err := v.VerifyToken(context.Background(), "", "")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
	if !strings.Contains(err.Error(), "비어있습니다") {
		t.Errorf("expected empty token error, got: %v", err)
	}
}

func TestGoogleVerifier_VerifyToken_InvalidJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`not valid json`))
	}))
	defer ts.Close()

	v := NewGoogleVerifier("test-client-id")
	v.SetTokenURL(ts.URL)

	_, err := v.VerifyToken(context.Background(), "some-token", "")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	if !strings.Contains(err.Error(), "파싱 실패") {
		t.Errorf("expected parse error, got: %v", err)
	}
}

func TestGoogleVerifier_VerifyToken_WrongIssuer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"sub": "12345",
			"email": "test@gmail.com",
			"name": "Test User",
			"picture": "https://pic.test/photo.jpg",
			"aud": "test-client-id",
			"iss": "https://evil.example.com",
			"email_verified": "true"
		}`))
	}))
	defer ts.Close()

	v := NewGoogleVerifier("test-client-id")
	v.SetTokenURL(ts.URL)

	_, err := v.VerifyToken(context.Background(), "some-token", "")
	if err == nil {
		t.Fatal("expected error for wrong issuer")
	}
	if !strings.Contains(err.Error(), "issuer 불일치") {
		t.Errorf("expected issuer mismatch error, got: %v", err)
	}
}

func TestGoogleVerifier_Provider(t *testing.T) {
	v := NewGoogleVerifier("test-client-id")
	if v.Provider() != "google" {
		t.Errorf("expected google, got %s", v.Provider())
	}
}
