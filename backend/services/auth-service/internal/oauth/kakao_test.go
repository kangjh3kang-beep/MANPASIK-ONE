package oauth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestKakaoVerifier_VerifyToken_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Authorization 헤더 검증
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-access-token" {
			t.Errorf("expected Authorization=Bearer test-access-token, got %s", auth)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 12345678,
			"kakao_account": {
				"email": "test@kakao.com",
				"profile": {
					"nickname": "카카오유저"
				}
			}
		}`))
	}))
	defer ts.Close()

	v := NewKakaoVerifier()
	v.SetAPIURL(ts.URL)

	info, err := v.VerifyToken(context.Background(), "", "test-access-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.ProviderID != "12345678" {
		t.Errorf("expected ProviderID=12345678, got %s", info.ProviderID)
	}
	if info.Email != "test@kakao.com" {
		t.Errorf("expected email=test@kakao.com, got %s", info.Email)
	}
	if info.DisplayName != "카카오유저" {
		t.Errorf("expected DisplayName=카카오유저, got %s", info.DisplayName)
	}
	if info.Provider != "kakao" {
		t.Errorf("expected Provider=kakao, got %s", info.Provider)
	}
}

func TestKakaoVerifier_VerifyToken_EmptyAccessToken(t *testing.T) {
	v := NewKakaoVerifier()

	_, err := v.VerifyToken(context.Background(), "", "")
	if err == nil {
		t.Fatal("expected error for empty access token")
	}
	if !strings.Contains(err.Error(), "비어있습니다") {
		t.Errorf("expected empty token error, got: %v", err)
	}
}

func TestKakaoVerifier_VerifyToken_HTTPError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"msg":"this access token does not exist","code":-401}`))
	}))
	defer ts.Close()

	v := NewKakaoVerifier()
	v.SetAPIURL(ts.URL)

	_, err := v.VerifyToken(context.Background(), "", "invalid-token")
	if err == nil {
		t.Fatal("expected error for HTTP 401")
	}
	if !strings.Contains(err.Error(), "HTTP 401") {
		t.Errorf("expected HTTP 401 error, got: %v", err)
	}
}

func TestKakaoVerifier_VerifyToken_IDZero(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 0,
			"kakao_account": {
				"email": "test@kakao.com",
				"profile": {
					"nickname": "유저"
				}
			}
		}`))
	}))
	defer ts.Close()

	v := NewKakaoVerifier()
	v.SetAPIURL(ts.URL)

	_, err := v.VerifyToken(context.Background(), "", "some-token")
	if err == nil {
		t.Fatal("expected error for id=0")
	}
	if !strings.Contains(err.Error(), "ID가 0") {
		t.Errorf("expected ID zero error, got: %v", err)
	}
}

func TestKakaoVerifier_VerifyToken_NoEmail_FallbackGenerated(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 99887766,
			"kakao_account": {
				"profile": {
					"nickname": "이메일없는유저"
				}
			}
		}`))
	}))
	defer ts.Close()

	v := NewKakaoVerifier()
	v.SetAPIURL(ts.URL)

	info, err := v.VerifyToken(context.Background(), "", "some-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedEmail := "kakao_99887766@social.manpasik.com"
	if info.Email != expectedEmail {
		t.Errorf("expected fallback email=%s, got %s", expectedEmail, info.Email)
	}
}

func TestKakaoVerifier_VerifyToken_InvalidJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{invalid json`))
	}))
	defer ts.Close()

	v := NewKakaoVerifier()
	v.SetAPIURL(ts.URL)

	_, err := v.VerifyToken(context.Background(), "", "some-token")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	if !strings.Contains(err.Error(), "파싱 실패") {
		t.Errorf("expected parse error, got: %v", err)
	}
}

func TestKakaoVerifier_Provider(t *testing.T) {
	v := NewKakaoVerifier()
	if v.Provider() != "kakao" {
		t.Errorf("expected kakao, got %s", v.Provider())
	}
}
