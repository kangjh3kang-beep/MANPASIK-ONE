package oauth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNaverVerifier_VerifyToken_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Authorization 헤더 검증
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-access-token" {
			t.Errorf("expected Authorization=Bearer test-access-token, got %s", auth)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"resultcode": "00",
			"message": "success",
			"response": {
				"id": "naver123",
				"email": "test@naver.com",
				"nickname": "네이버유저",
				"profile_image": "https://pic.test/naver.jpg"
			}
		}`))
	}))
	defer ts.Close()

	v := NewNaverVerifier()
	v.SetAPIURL(ts.URL)

	info, err := v.VerifyToken(context.Background(), "", "test-access-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.ProviderID != "naver123" {
		t.Errorf("expected ProviderID=naver123, got %s", info.ProviderID)
	}
	if info.Email != "test@naver.com" {
		t.Errorf("expected email=test@naver.com, got %s", info.Email)
	}
	if info.DisplayName != "네이버유저" {
		t.Errorf("expected DisplayName=네이버유저, got %s", info.DisplayName)
	}
	if info.ProfileImage != "https://pic.test/naver.jpg" {
		t.Errorf("expected ProfileImage=https://pic.test/naver.jpg, got %s", info.ProfileImage)
	}
	if info.Provider != "naver" {
		t.Errorf("expected Provider=naver, got %s", info.Provider)
	}
}

func TestNaverVerifier_VerifyToken_EmptyAccessToken(t *testing.T) {
	v := NewNaverVerifier()

	_, err := v.VerifyToken(context.Background(), "", "")
	if err == nil {
		t.Fatal("expected error for empty access token")
	}
	if !strings.Contains(err.Error(), "비어있습니다") {
		t.Errorf("expected empty token error, got: %v", err)
	}
}

func TestNaverVerifier_VerifyToken_HTTPError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`Internal Server Error`))
	}))
	defer ts.Close()

	v := NewNaverVerifier()
	v.SetAPIURL(ts.URL)

	_, err := v.VerifyToken(context.Background(), "", "some-token")
	if err == nil {
		t.Fatal("expected error for HTTP 500")
	}
	if !strings.Contains(err.Error(), "HTTP 500") {
		t.Errorf("expected HTTP 500 error, got: %v", err)
	}
}

func TestNaverVerifier_VerifyToken_ResultCodeNotOK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"resultcode": "024",
			"message": "Authentication failed",
			"response": {}
		}`))
	}))
	defer ts.Close()

	v := NewNaverVerifier()
	v.SetAPIURL(ts.URL)

	_, err := v.VerifyToken(context.Background(), "", "some-token")
	if err == nil {
		t.Fatal("expected error for resultcode != 00")
	}
	if !strings.Contains(err.Error(), "결과 오류") {
		t.Errorf("expected result error, got: %v", err)
	}
}

func TestNaverVerifier_VerifyToken_EmptyEmail_FallbackGenerated(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"resultcode": "00",
			"message": "success",
			"response": {
				"id": "naver456",
				"nickname": "비공개유저",
				"profile_image": ""
			}
		}`))
	}))
	defer ts.Close()

	v := NewNaverVerifier()
	v.SetAPIURL(ts.URL)

	info, err := v.VerifyToken(context.Background(), "", "some-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedEmail := "naver_naver456@social.manpasik.com"
	if info.Email != expectedEmail {
		t.Errorf("expected fallback email=%s, got %s", expectedEmail, info.Email)
	}
}

func TestNaverVerifier_VerifyToken_InvalidJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`not valid json at all`))
	}))
	defer ts.Close()

	v := NewNaverVerifier()
	v.SetAPIURL(ts.URL)

	_, err := v.VerifyToken(context.Background(), "", "some-token")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	if !strings.Contains(err.Error(), "파싱 실패") {
		t.Errorf("expected parse error, got: %v", err)
	}
}

func TestNaverVerifier_Provider(t *testing.T) {
	v := NewNaverVerifier()
	if v.Provider() != "naver" {
		t.Errorf("expected naver, got %s", v.Provider())
	}
}
