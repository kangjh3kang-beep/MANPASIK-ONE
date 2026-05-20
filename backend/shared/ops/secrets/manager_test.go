package secrets_test

import (
	"context"
	"testing"
	"time"

	"github.com/manpasik/backend/shared/ops/secrets"
)

func TestMemoryProvider_SetGet(t *testing.T) {
	p := secrets.NewMemoryProvider()
	ctx := context.Background()

	s := &secrets.Secret{Path: "auth/jwt", Value: "secret-key"}
	if err := p.Set(ctx, s); err != nil {
		t.Fatalf("Set 실패: %v", err)
	}
	if s.Version != 1 {
		t.Errorf("Version = %d, want 1", s.Version)
	}

	got, err := p.Get(ctx, "auth/jwt")
	if err != nil {
		t.Fatalf("Get 실패: %v", err)
	}
	if got.Value != "secret-key" {
		t.Errorf("Value = %q", got.Value)
	}
}

func TestMemoryProvider_VersionIncrement(t *testing.T) {
	p := secrets.NewMemoryProvider()
	ctx := context.Background()

	_ = p.Set(ctx, &secrets.Secret{Path: "x", Value: "v1"})
	_ = p.Set(ctx, &secrets.Secret{Path: "x", Value: "v2"})
	_ = p.Set(ctx, &secrets.Secret{Path: "x", Value: "v3"})

	got, _ := p.Get(ctx, "x")
	if got.Version != 3 {
		t.Errorf("Version = %d, want 3", got.Version)
	}
	if got.Value != "v3" {
		t.Errorf("Value = %q, want v3", got.Value)
	}
}

func TestMemoryProvider_NotFound(t *testing.T) {
	p := secrets.NewMemoryProvider()
	_, err := p.Get(context.Background(), "missing")
	if err == nil {
		t.Error("미존재 시크릿이 통과됨")
	}
}

func TestMemoryProvider_Expiration(t *testing.T) {
	p := secrets.NewMemoryProvider()
	ctx := context.Background()
	past := time.Now().UTC().Add(-1 * time.Hour)

	_ = p.Set(ctx, &secrets.Secret{
		Path:      "expired",
		Value:     "x",
		ExpiresAt: &past,
	})

	_, err := p.Get(ctx, "expired")
	if err == nil {
		t.Error("만료된 시크릿이 조회됨")
	}
}

func TestMemoryProvider_Delete(t *testing.T) {
	p := secrets.NewMemoryProvider()
	ctx := context.Background()
	_ = p.Set(ctx, &secrets.Secret{Path: "del", Value: "x"})
	_ = p.Delete(ctx, "del")

	_, err := p.Get(ctx, "del")
	if err == nil {
		t.Error("삭제된 시크릿이 조회됨")
	}
}

func TestMemoryProvider_List(t *testing.T) {
	p := secrets.NewMemoryProvider()
	ctx := context.Background()
	_ = p.Set(ctx, &secrets.Secret{Path: "auth/jwt", Value: "x"})
	_ = p.Set(ctx, &secrets.Secret{Path: "auth/refresh", Value: "y"})
	_ = p.Set(ctx, &secrets.Secret{Path: "db/password", Value: "z"})

	authPaths, err := p.List(ctx, "auth/")
	if err != nil {
		t.Fatalf("List 실패: %v", err)
	}
	if len(authPaths) != 2 {
		t.Errorf("auth paths = %d, want 2", len(authPaths))
	}
}

func TestMemoryProvider_Rotate(t *testing.T) {
	p := secrets.NewMemoryProvider()
	ctx := context.Background()
	_ = p.Set(ctx, &secrets.Secret{Path: "rotate", Value: "old-value"})

	rotated, err := p.Rotate(ctx, "rotate")
	if err != nil {
		t.Fatalf("Rotate 실패: %v", err)
	}
	if rotated.Value == "old-value" {
		t.Error("로테이션 후에도 값이 동일")
	}
	if rotated.Version != 2 {
		t.Errorf("Version = %d, want 2", rotated.Version)
	}
}

func TestEnvProvider_GetSet(t *testing.T) {
	p := secrets.NewEnvProvider("MPSK_TEST_")
	ctx := context.Background()

	t.Setenv("MPSK_TEST_AUTH_JWT", "env-secret")
	got, err := p.Get(ctx, "auth/jwt")
	if err != nil {
		t.Fatalf("Get 실패: %v", err)
	}
	if got.Value != "env-secret" {
		t.Errorf("Value = %q", got.Value)
	}
}

func TestEnvProvider_NotSet(t *testing.T) {
	p := secrets.NewEnvProvider("MPSK_TEST2_")
	_, err := p.Get(context.Background(), "missing")
	if err == nil {
		t.Error("미설정 환경변수가 통과됨")
	}
}

func TestEnvProvider_Set(t *testing.T) {
	p := secrets.NewEnvProvider("MPSK_TEST3_")
	ctx := context.Background()

	if err := p.Set(ctx, &secrets.Secret{Path: "key/value", Value: "data"}); err != nil {
		t.Fatalf("Set 실패: %v", err)
	}
	got, _ := p.Get(ctx, "key/value")
	if got.Value != "data" {
		t.Errorf("Value = %q", got.Value)
	}
}

func TestVaultProvider_HealthCheck(t *testing.T) {
	p := secrets.NewVaultProvider("https://vault.example.com", "token", "")
	if err := p.HealthCheck(context.Background()); err != nil {
		t.Errorf("HealthCheck 실패: %v", err)
	}

	p2 := secrets.NewVaultProvider("", "", "")
	if err := p2.HealthCheck(context.Background()); err == nil {
		t.Error("자격증명 없이 통과됨")
	}
}

func TestAWSProvider_HealthCheck(t *testing.T) {
	p := secrets.NewAWSProvider("us-east-1")
	if err := p.HealthCheck(context.Background()); err != nil {
		t.Errorf("HealthCheck 실패: %v", err)
	}
}

func TestNewFromEnv_Default(t *testing.T) {
	t.Setenv("SECRETS_PROVIDER", "")
	p := secrets.NewFromEnv()
	if p.Provider() != "memory" {
		t.Errorf("Provider = %q, want memory", p.Provider())
	}
}

func TestNewFromEnv_Vault(t *testing.T) {
	t.Setenv("SECRETS_PROVIDER", "vault")
	t.Setenv("VAULT_ADDR", "https://x")
	t.Setenv("VAULT_TOKEN", "tok")
	p := secrets.NewFromEnv()
	if p.Provider() != "vault" {
		t.Errorf("Provider = %q, want vault", p.Provider())
	}
}

func TestRotationManager_Policy(t *testing.T) {
	provider := secrets.NewMemoryProvider()
	mgr := secrets.NewRotationManager(provider)

	mgr.SetPolicy("auth/jwt", secrets.RotationPolicy{
		Enabled:  true,
		Interval: 24 * time.Hour,
		MaxAge:   1 * time.Hour,
	})

	// 오래된 시크릿 시뮬레이션
	old := time.Now().UTC().Add(-2 * time.Hour)
	_ = provider.Set(context.Background(), &secrets.Secret{
		Path:      "auth/jwt",
		Value:     "old",
		CreatedAt: old,
	})

	rotated, err := mgr.CheckAndRotate(context.Background())
	if err != nil {
		t.Fatalf("CheckAndRotate 실패: %v", err)
	}
	if len(rotated) != 1 || rotated[0] != "auth/jwt" {
		t.Errorf("rotated = %v, want [auth/jwt]", rotated)
	}
}

func TestRotationManager_DisabledPolicy(t *testing.T) {
	provider := secrets.NewMemoryProvider()
	mgr := secrets.NewRotationManager(provider)

	mgr.SetPolicy("x", secrets.RotationPolicy{Enabled: false, MaxAge: 1 * time.Nanosecond})
	old := time.Now().UTC().Add(-1 * time.Hour)
	_ = provider.Set(context.Background(), &secrets.Secret{Path: "x", Value: "v", CreatedAt: old})

	rotated, _ := mgr.CheckAndRotate(context.Background())
	if len(rotated) != 0 {
		t.Errorf("disabled 정책이 로테이션됨: %v", rotated)
	}
}

func TestSecret_IsExpired(t *testing.T) {
	past := time.Now().UTC().Add(-1 * time.Hour)
	future := time.Now().UTC().Add(1 * time.Hour)

	if !(&secrets.Secret{ExpiresAt: &past}).IsExpired() {
		t.Error("과거 만료가 IsExpired=false")
	}
	if (&secrets.Secret{ExpiresAt: &future}).IsExpired() {
		t.Error("미래 만료가 IsExpired=true")
	}
	if (&secrets.Secret{}).IsExpired() {
		t.Error("ExpiresAt nil은 IsExpired=false")
	}
}
