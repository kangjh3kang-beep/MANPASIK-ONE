package oauth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAppleVerifier_VerifyToken_EmptyToken(t *testing.T) {
	v := NewAppleVerifier("com.manpasik.app")

	_, err := v.VerifyToken(context.Background(), "", "")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
	if !strings.Contains(err.Error(), "비어있습니다") {
		t.Errorf("expected empty token error, got: %v", err)
	}
}

func TestAppleVerifier_VerifyToken_InvalidJWTFormat(t *testing.T) {
	v := NewAppleVerifier("com.manpasik.app")

	_, err := v.VerifyToken(context.Background(), "not-a-valid-jwt", "")
	if err == nil {
		t.Fatal("expected error for invalid JWT format")
	}
	if !strings.Contains(err.Error(), "JWT 파싱 실패") {
		t.Errorf("expected JWT parse error, got: %v", err)
	}
}

func TestAppleVerifier_Provider(t *testing.T) {
	v := NewAppleVerifier("com.manpasik.app")
	if v.Provider() != "apple" {
		t.Errorf("expected apple, got %s", v.Provider())
	}
}

func TestAppleVerifier_SetKeysURL(t *testing.T) {
	v := NewAppleVerifier("com.manpasik.app")

	customURL := "https://custom.apple.test/auth/keys"
	v.SetKeysURL(customURL)

	if v.keysURL != customURL {
		t.Errorf("expected keysURL=%s, got %s", customURL, v.keysURL)
	}
}

func TestAppleVerifier_FetchKeys_MockJWKS(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"keys": [
				{
					"kty": "RSA",
					"kid": "test-kid",
					"n": "0vx7agoebGcQSuuPiLJXZptN9nndrQmbXEps2aiAFbWhM78LhWx4cbbfAAtVT86zwu1RK7aPFFxuhDR1L6tSoc_BJECPebWKRXjBZCiFV4n3oknjhMstn64tZ_2W-5JsGY4Hc5n9yBXArwl93lqt7_RN5w6Cf0h4QyQ5v-65YGjQR0_FDW2QvzqY368QQMicAtaSqzs8KJZgnYb9c7d0zgdAZHzu6qMQvRL5hajrn1n91CbOpbISD08qNLyrdkt-bFTWhAI4vMQFh6WeZu0fM4lFd2NcRwr3XPksINHaQ-G_xBniIqbw0Ls1jF44-csFCur-kEgU8awapJzKnqDKgw",
					"e": "AQAB",
					"alg": "RS256",
					"use": "sig"
				}
			]
		}`))
	}))
	defer ts.Close()

	v := NewAppleVerifier("com.manpasik.app")
	v.SetKeysURL(ts.URL)

	keys, err := v.fetchKeys(context.Background())
	if err != nil {
		t.Fatalf("unexpected error fetching keys: %v", err)
	}
	if len(keys) != 1 {
		t.Fatalf("expected 1 key, got %d", len(keys))
	}
	if keys[0].Kid != "test-kid" {
		t.Errorf("expected kid=test-kid, got %s", keys[0].Kid)
	}
	if keys[0].Kty != "RSA" {
		t.Errorf("expected kty=RSA, got %s", keys[0].Kty)
	}
	if keys[0].Alg != "RS256" {
		t.Errorf("expected alg=RS256, got %s", keys[0].Alg)
	}
}
