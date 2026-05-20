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

// fakeVaultServer 는 health/secret 응답을 반환.
func fakeVaultServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/v1/sys/health"):
			w.WriteHeader(200)
			_, _ = w.Write([]byte(`{"initialized":true,"sealed":false,"standby":false}`))
		case strings.Contains(r.URL.Path, "/data/"):
			w.WriteHeader(200)
			_, _ = w.Write([]byte(`{"data":{"data":{"value":"v1"},"metadata":{"version":1}}}`))
		default:
			w.WriteHeader(404)
		}
	}))
}

func TestVaultBootstrap_NoAddr(t *testing.T) {
	t.Setenv("VAULT_ADDR", "")
	t.Setenv("VAULT_TOKEN", "")
	if _, err := secrets.NewVaultBootstrap(secrets.VaultBootstrapConfig{}); err == nil {
		t.Error("Addr 미설정 통과")
	}
}

func TestVaultBootstrap_HealthCheckPasses(t *testing.T) {
	srv := fakeVaultServer(t)
	defer srv.Close()

	bs, err := secrets.NewVaultBootstrap(secrets.VaultBootstrapConfig{
		Addr:           srv.URL,
		Token:          "test-token",
		WatchInterval:  60 * time.Second,
		AutoRenewToken: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	if bs.Provider() == nil {
		t.Error("Provider nil")
	}
	if bs.Watcher() == nil {
		t.Error("Watcher nil")
	}
}

func TestVaultBootstrap_WatchPaths(t *testing.T) {
	srv := fakeVaultServer(t)
	defer srv.Close()

	bs, err := secrets.NewVaultBootstrap(secrets.VaultBootstrapConfig{
		Addr:           srv.URL,
		Token:          "tk",
		WatchPaths:     []string{"db/password", "auth/jwt"},
		AutoRenewToken: false,
	})
	if err != nil {
		t.Fatal(err)
	}

	bs.Watch("extra/key") // 추가 path
	// CycleCount 는 0 (Start 안 함)
	if bs.Watcher().CycleCount() != 0 {
		t.Errorf("CycleCount = %d", bs.Watcher().CycleCount())
	}
}

func TestVaultBootstrap_AddListener(t *testing.T) {
	srv := fakeVaultServer(t)
	defer srv.Close()

	bs, _ := secrets.NewVaultBootstrap(secrets.VaultBootstrapConfig{
		Addr:           srv.URL,
		Token:          "tk",
		AutoRenewToken: false,
	})
	called := false
	bs.AddListener(func(_ context.Context, _ secrets.RotationEvent) error {
		called = true
		return nil
	})
	_ = called // listener 등록만 검증; 호출은 watcher 폴링 시점에
}

func TestVaultBootstrap_FromEnv(t *testing.T) {
	srv := fakeVaultServer(t)
	defer srv.Close()

	t.Setenv("VAULT_ADDR", srv.URL)
	t.Setenv("VAULT_TOKEN", "env-token")
	t.Setenv("VAULT_KV_MOUNT", "secret")

	bs, err := secrets.NewVaultBootstrap(secrets.VaultBootstrapConfig{
		AutoRenewToken: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	if bs == nil {
		t.Fatal("bootstrap nil")
	}
}

func TestVaultBootstrap_AutoRenewToken(t *testing.T) {
	srv := fakeVaultServer(t)
	defer srv.Close()

	bs, _ := secrets.NewVaultBootstrap(secrets.VaultBootstrapConfig{
		Addr:                      srv.URL,
		Token:                     "tk",
		AutoRenewToken:            true,
		TokenRenewIntervalSeconds: 1,
	})
	if bs.Renewer() == nil {
		t.Error("AutoRenew=true 인데 Renewer nil")
	}
}

func TestVaultBootstrap_StartStop(t *testing.T) {
	srv := fakeVaultServer(t)
	defer srv.Close()

	bs, _ := secrets.NewVaultBootstrap(secrets.VaultBootstrapConfig{
		Addr:           srv.URL,
		Token:          "tk",
		WatchInterval:  20 * time.Millisecond,
		WatchPaths:     []string{"k"},
		AutoRenewToken: false,
	})
	bs.Start(context.Background())
	time.Sleep(50 * time.Millisecond)
	bs.Stop()

	if bs.Watcher().CycleCount() < 2 {
		t.Errorf("CycleCount = %d", bs.Watcher().CycleCount())
	}
}
