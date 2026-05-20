package secrets_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/manpasik/backend/shared/ops/secrets"
)

func TestVaultHTTPProvider_HealthCheck_Healthy(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sys/health" {
			http.Error(w, "wrong path", http.StatusNotFound)
			return
		}
		if r.Header.Get("X-Vault-Token") != "test-token" {
			http.Error(w, "no token", http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"initialized":true,"sealed":false}`))
	}))
	defer server.Close()

	p := secrets.NewVaultHTTPProvider(server.URL, "test-token", "", "secret")
	if err := p.HealthCheck(context.Background()); err != nil {
		t.Errorf("HealthCheck = %v", err)
	}
}

func TestVaultHTTPProvider_HealthCheck_NoToken(t *testing.T) {
	p := secrets.NewVaultHTTPProvider("https://vault.example.com", "", "", "")
	if err := p.HealthCheck(context.Background()); err == nil {
		t.Error("토큰 없이 통과")
	}
}

func TestVaultHTTPProvider_HealthCheck_NoAddr(t *testing.T) {
	p := secrets.NewVaultHTTPProvider("", "tok", "", "")
	if err := p.HealthCheck(context.Background()); err == nil {
		t.Error("addr 없이 통과")
	}
}

func TestVaultHTTPProvider_Get_KVv2(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/v1/secret/data/auth/jwt") {
			t.Errorf("path = %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"data": {
				"data": {"value": "supersecret-key", "owner": "auth-team"},
				"metadata": {"version": 3, "created_time": "2026-04-30T12:00:00Z"}
			}
		}`))
	}))
	defer server.Close()

	p := secrets.NewVaultHTTPProvider(server.URL, "tok", "", "secret")
	got, err := p.Get(context.Background(), "auth/jwt")
	if err != nil {
		t.Fatalf("Get 실패: %v", err)
	}
	if got.Value != "supersecret-key" {
		t.Errorf("Value = %q", got.Value)
	}
	if got.Version != 3 {
		t.Errorf("Version = %d", got.Version)
	}
	if got.Metadata["owner"] != "auth-team" {
		t.Errorf("Metadata.owner = %q", got.Metadata["owner"])
	}
}

func TestVaultHTTPProvider_Get_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	p := secrets.NewVaultHTTPProvider(server.URL, "tok", "", "secret")
	_, err := p.Get(context.Background(), "missing")
	if err == nil {
		t.Error("404에 에러 없음")
	}
}

func TestVaultHTTPProvider_Get_PermissionDenied(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	p := secrets.NewVaultHTTPProvider(server.URL, "expired-tok", "", "secret")
	_, err := p.Get(context.Background(), "x")
	if err == nil || !strings.Contains(err.Error(), "permission") {
		t.Errorf("err = %v, want permission denied", err)
	}
}

func TestVaultHTTPProvider_Set(t *testing.T) {
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Method = %q", r.Method)
		}
		called = true
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"version":1}}`))
	}))
	defer server.Close()

	p := secrets.NewVaultHTTPProvider(server.URL, "tok", "", "secret")
	err := p.Set(context.Background(), &secrets.Secret{
		Path:  "test/path",
		Value: "secret-value",
		Metadata: map[string]string{"owner": "team-a"},
	})
	if err != nil {
		t.Fatalf("Set 실패: %v", err)
	}
	if !called {
		t.Error("Set HTTP 호출 안됨")
	}
}

func TestVaultHTTPProvider_Set_RequiresPath(t *testing.T) {
	p := secrets.NewVaultHTTPProvider("https://x", "tok", "", "")
	if err := p.Set(context.Background(), &secrets.Secret{Value: "x"}); err == nil {
		t.Error("path 없이 통과")
	}
}

func TestVaultHTTPProvider_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Method = %q", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/metadata/") {
			t.Errorf("path = %q (DELETE는 metadata 엔드포인트)", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	p := secrets.NewVaultHTTPProvider(server.URL, "tok", "", "secret")
	if err := p.Delete(context.Background(), "x"); err != nil {
		t.Errorf("Delete = %v", err)
	}
}

func TestVaultHTTPProvider_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.RawQuery, "list=true") {
			t.Errorf("query = %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"data":{"keys":["jwt","refresh","db_password"]}}`))
	}))
	defer server.Close()

	p := secrets.NewVaultHTTPProvider(server.URL, "tok", "", "secret")
	keys, err := p.List(context.Background(), "auth/")
	if err != nil {
		t.Fatalf("List = %v", err)
	}
	if len(keys) != 3 {
		t.Errorf("keys = %d, want 3", len(keys))
	}
	if keys[0] != "auth/jwt" {
		t.Errorf("[0] = %q (prefix 미적용)", keys[0])
	}
}

func TestVaultHTTPProvider_List_EmptyPrefix(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	p := secrets.NewVaultHTTPProvider(server.URL, "tok", "", "secret")
	keys, err := p.List(context.Background(), "missing/")
	if err != nil {
		t.Errorf("404는 빈 결과여야 함: %v", err)
	}
	if len(keys) != 0 {
		t.Errorf("keys = %d, want 0", len(keys))
	}
}

func TestVaultHTTPProvider_Rotate(t *testing.T) {
	getCalled := false
	setCalled := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getCalled = true
			_, _ = w.Write([]byte(`{"data":{"data":{"value":"old"},"metadata":{"version":1}}}`))
		case "POST":
			setCalled = true
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"data":{"version":2}}`))
		}
	}))
	defer server.Close()

	p := secrets.NewVaultHTTPProvider(server.URL, "tok", "", "secret")
	rotated, err := p.Rotate(context.Background(), "test")
	if err != nil {
		t.Fatalf("Rotate = %v", err)
	}
	if !getCalled || !setCalled {
		t.Errorf("get=%v set=%v (둘 다 true여야 함)", getCalled, setCalled)
	}
	if rotated.Value == "old" {
		t.Error("Rotate 후에도 같은 값")
	}
}

func TestVaultHTTPProvider_Namespace(t *testing.T) {
	gotNS := ""
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotNS = r.Header.Get("X-Vault-Namespace")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"initialized":true,"sealed":false}`))
	}))
	defer server.Close()

	p := secrets.NewVaultHTTPProvider(server.URL, "tok", "manpasik-prod", "secret")
	_ = p.HealthCheck(context.Background())
	if gotNS != "manpasik-prod" {
		t.Errorf("Namespace 헤더 = %q", gotNS)
	}
}

func TestTokenAutoRenewer_RenewSelf(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "auth/token/renew-self") {
			t.Errorf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"auth":{"lease_duration":3600,"renewable":true}}`))
	}))
	defer server.Close()

	p := secrets.NewVaultHTTPProvider(server.URL, "tok", "", "secret")
	r := secrets.NewTokenAutoRenewer(p, 5*time.Minute)

	dur, err := r.RenewSelf(context.Background())
	if err != nil {
		t.Fatalf("RenewSelf = %v", err)
	}
	if dur != time.Hour {
		t.Errorf("Lease = %v, want 1h", dur)
	}
}

func TestTokenAutoRenewer_StartStop(t *testing.T) {
	p := secrets.NewVaultHTTPProvider("https://x", "tok", "", "")
	r := secrets.NewTokenAutoRenewer(p, 5*time.Minute)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := r.Start(ctx, 3600); err != nil {
		t.Fatalf("Start = %v", err)
	}
	if !r.IsRunning() {
		t.Error("Start 후 IsRunning = false")
	}

	if err := r.Start(ctx, 3600); err == nil {
		t.Error("이중 Start 통과")
	}

	r.Stop()
	if r.IsRunning() {
		t.Error("Stop 후에도 IsRunning = true")
	}
}

func TestVaultHTTPProvider_Get_RequiresPath(t *testing.T) {
	p := secrets.NewVaultHTTPProvider("https://x", "tok", "", "")
	if _, err := p.Get(context.Background(), ""); err == nil {
		t.Error("빈 path 통과")
	}
}

func TestVaultHTTPProvider_Provider(t *testing.T) {
	p := secrets.NewVaultHTTPProvider("x", "y", "z", "w")
	if p.Provider() != "vault_http" {
		t.Errorf("Provider = %q", p.Provider())
	}
}
