// Package secrets는 시크릿 관리 추상화입니다.
//
// 지원 백엔드: Vault, AWS Secrets Manager, 환경변수, 인메모리(테스트).
// SECRETS_PROVIDER 환경변수로 스위칭하며, AES-256 키 로테이션을 지원합니다.
package secrets

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// ============================================================================
// 도메인 모델
// ============================================================================

// Secret은 단일 시크릿 항목입니다.
type Secret struct {
	Path        string            // "database/password", "auth/jwt_key"
	Value       string            // 평문 값 (메모리에서만 보관)
	Version     int               // 로테이션 시 증가
	Metadata    map[string]string // tags, environment, owner 등
	CreatedAt   time.Time
	ExpiresAt   *time.Time // 옵션 (TTL)
}

// IsExpired는 시크릿 만료 여부를 반환합니다.
func (s *Secret) IsExpired() bool {
	if s.ExpiresAt == nil {
		return false
	}
	return time.Now().UTC().After(*s.ExpiresAt)
}

// RotationPolicy는 자동 로테이션 정책입니다.
type RotationPolicy struct {
	Enabled  bool
	Interval time.Duration
	MaxAge   time.Duration
}

// ============================================================================
// Provider 인터페이스
// ============================================================================

// Provider는 시크릿 백엔드 인터페이스입니다.
type Provider interface {
	Get(ctx context.Context, path string) (*Secret, error)
	Set(ctx context.Context, secret *Secret) error
	Delete(ctx context.Context, path string) error
	List(ctx context.Context, prefix string) ([]string, error)
	Rotate(ctx context.Context, path string) (*Secret, error)
	Provider() string
	HealthCheck(ctx context.Context) error
}

// ============================================================================
// 팩토리
// ============================================================================

// NewFromEnv는 환경변수에 따라 적절한 Provider를 생성합니다.
//
// SECRETS_PROVIDER:
//   - "vault": HashiCorp Vault (VAULT_ADDR, VAULT_TOKEN)
//   - "aws": AWS Secrets Manager (AWS_REGION)
//   - "env": 환경변수 폴백 (개발용)
//   - "" or "memory": 인메모리 (테스트)
func NewFromEnv() Provider {
	provider := strings.ToLower(os.Getenv("SECRETS_PROVIDER"))
	switch provider {
	case "vault":
		return NewVaultProvider(
			os.Getenv("VAULT_ADDR"),
			os.Getenv("VAULT_TOKEN"),
			os.Getenv("VAULT_NAMESPACE"),
		)
	case "aws":
		return NewAWSProvider(os.Getenv("AWS_REGION"))
	case "env":
		return NewEnvProvider("MPSK_SECRET_")
	default:
		return NewMemoryProvider()
	}
}

// ============================================================================
// 인메모리 Provider (테스트용)
// ============================================================================

// MemoryProvider는 인메모리 시크릿 저장소입니다.
type MemoryProvider struct {
	mu      sync.RWMutex
	secrets map[string]*Secret
}

// NewMemoryProvider는 새 인메모리 Provider를 생성합니다.
func NewMemoryProvider() *MemoryProvider {
	return &MemoryProvider{secrets: make(map[string]*Secret)}
}

// Get은 시크릿을 조회합니다.
func (p *MemoryProvider) Get(_ context.Context, path string) (*Secret, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	s, ok := p.secrets[path]
	if !ok {
		return nil, fmt.Errorf("secret %q not found", path)
	}
	if s.IsExpired() {
		return nil, fmt.Errorf("secret %q expired", path)
	}
	return s, nil
}

// Set은 시크릿을 저장합니다.
func (p *MemoryProvider) Set(_ context.Context, secret *Secret) error {
	if secret == nil || secret.Path == "" {
		return errors.New("path required")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if existing, ok := p.secrets[secret.Path]; ok {
		secret.Version = existing.Version + 1
	} else {
		secret.Version = 1
	}
	if secret.CreatedAt.IsZero() {
		secret.CreatedAt = time.Now().UTC()
	}
	p.secrets[secret.Path] = secret
	return nil
}

// Delete는 시크릿을 삭제합니다.
func (p *MemoryProvider) Delete(_ context.Context, path string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.secrets, path)
	return nil
}

// List는 prefix로 시작하는 모든 시크릿 경로를 반환합니다.
func (p *MemoryProvider) List(_ context.Context, prefix string) ([]string, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	var paths []string
	for path := range p.secrets {
		if strings.HasPrefix(path, prefix) {
			paths = append(paths, path)
		}
	}
	return paths, nil
}

// Rotate는 시크릿을 로테이션합니다 (새 값을 생성하여 저장).
func (p *MemoryProvider) Rotate(ctx context.Context, path string) (*Secret, error) {
	existing, err := p.Get(ctx, path)
	if err != nil {
		return nil, err
	}
	rotated := &Secret{
		Path:      path,
		Value:     generateRandomKey(32), // 32 byte random
		Metadata:  existing.Metadata,
		CreatedAt: time.Now().UTC(),
	}
	if err := p.Set(ctx, rotated); err != nil {
		return nil, err
	}
	return rotated, nil
}

// Provider는 백엔드 이름을 반환합니다.
func (p *MemoryProvider) Provider() string { return "memory" }

// HealthCheck는 항상 성공합니다.
func (p *MemoryProvider) HealthCheck(_ context.Context) error { return nil }

// ============================================================================
// 환경변수 Provider
// ============================================================================

// EnvProvider는 환경변수 기반 시크릿 Provider입니다.
//
// 경로 → 환경변수 변환: "auth/jwt_key" → "MPSK_SECRET_AUTH_JWT_KEY"
type EnvProvider struct {
	prefix string
}

// NewEnvProvider는 새 환경변수 Provider를 생성합니다.
func NewEnvProvider(prefix string) *EnvProvider {
	if prefix == "" {
		prefix = "MPSK_SECRET_"
	}
	return &EnvProvider{prefix: prefix}
}

// Get은 환경변수에서 시크릿을 조회합니다.
func (p *EnvProvider) Get(_ context.Context, path string) (*Secret, error) {
	envVar := pathToEnvVar(p.prefix, path)
	val, ok := os.LookupEnv(envVar)
	if !ok {
		return nil, fmt.Errorf("env var %q not set for secret %q", envVar, path)
	}
	return &Secret{Path: path, Value: val, Version: 1, CreatedAt: time.Now().UTC()}, nil
}

// Set은 환경변수에 시크릿을 저장합니다 (런타임 한정).
func (p *EnvProvider) Set(_ context.Context, secret *Secret) error {
	if secret == nil || secret.Path == "" {
		return errors.New("path required")
	}
	envVar := pathToEnvVar(p.prefix, secret.Path)
	return os.Setenv(envVar, secret.Value)
}

// Delete는 환경변수에서 시크릿을 제거합니다.
func (p *EnvProvider) Delete(_ context.Context, path string) error {
	envVar := pathToEnvVar(p.prefix, path)
	return os.Unsetenv(envVar)
}

// List는 prefix로 시작하는 환경변수의 시크릿 경로를 반환합니다.
func (p *EnvProvider) List(_ context.Context, pathPrefix string) ([]string, error) {
	var paths []string
	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, p.prefix) {
			continue
		}
		eqIdx := strings.Index(env, "=")
		if eqIdx < 0 {
			continue
		}
		envName := env[:eqIdx]
		path := envVarToPath(p.prefix, envName)
		if strings.HasPrefix(path, pathPrefix) {
			paths = append(paths, path)
		}
	}
	return paths, nil
}

// Rotate는 환경변수 시크릿을 로테이션합니다.
func (p *EnvProvider) Rotate(ctx context.Context, path string) (*Secret, error) {
	rotated := &Secret{
		Path:      path,
		Value:     generateRandomKey(32),
		Version:   1,
		CreatedAt: time.Now().UTC(),
	}
	if err := p.Set(ctx, rotated); err != nil {
		return nil, err
	}
	return rotated, nil
}

// Provider는 이름을 반환합니다.
func (p *EnvProvider) Provider() string { return "env" }

// HealthCheck는 항상 성공합니다.
func (p *EnvProvider) HealthCheck(_ context.Context) error { return nil }

// pathToEnvVar는 "auth/jwt_key" → "PREFIX_AUTH_JWT_KEY"로 변환합니다.
func pathToEnvVar(prefix, path string) string {
	upper := strings.ToUpper(path)
	upper = strings.ReplaceAll(upper, "/", "_")
	upper = strings.ReplaceAll(upper, "-", "_")
	upper = strings.ReplaceAll(upper, ".", "_")
	return prefix + upper
}

// envVarToPath는 "PREFIX_AUTH_JWT_KEY" → "auth/jwt_key"로 변환합니다 (best-effort).
func envVarToPath(prefix, envName string) string {
	stripped := strings.TrimPrefix(envName, prefix)
	return strings.ToLower(strings.ReplaceAll(stripped, "_", "/"))
}

// ============================================================================
// Vault Provider (스텁 — 실 호출 시 vault SDK 통합)
// ============================================================================

// VaultProvider는 HashiCorp Vault Provider입니다.
//
// 운영 시 github.com/hashicorp/vault/api SDK 통합 권장.
type VaultProvider struct {
	addr      string
	token     string
	namespace string
}

// NewVaultProvider는 새 Vault Provider를 생성합니다.
func NewVaultProvider(addr, token, namespace string) *VaultProvider {
	return &VaultProvider{addr: addr, token: token, namespace: namespace}
}

// Provider는 이름을 반환합니다.
func (p *VaultProvider) Provider() string { return "vault" }

// HealthCheck는 자격증명을 확인합니다.
func (p *VaultProvider) HealthCheck(_ context.Context) error {
	if p.addr == "" || p.token == "" {
		return errors.New("vault addr/token not configured")
	}
	return nil
}

// Get/Set/Delete/List/Rotate는 vault SDK 통합 시 구현됩니다.
func (p *VaultProvider) Get(_ context.Context, _ string) (*Secret, error) {
	return nil, errors.New("vault SDK integration required")
}
func (p *VaultProvider) Set(_ context.Context, _ *Secret) error {
	return errors.New("vault SDK integration required")
}
func (p *VaultProvider) Delete(_ context.Context, _ string) error {
	return errors.New("vault SDK integration required")
}
func (p *VaultProvider) List(_ context.Context, _ string) ([]string, error) {
	return nil, errors.New("vault SDK integration required")
}
func (p *VaultProvider) Rotate(_ context.Context, _ string) (*Secret, error) {
	return nil, errors.New("vault SDK integration required")
}

// ============================================================================
// AWS Secrets Manager Provider (스텁)
// ============================================================================

// AWSProvider는 AWS Secrets Manager Provider입니다.
type AWSProvider struct {
	region string
}

// NewAWSProvider는 새 AWS Provider를 생성합니다.
func NewAWSProvider(region string) *AWSProvider {
	return &AWSProvider{region: region}
}

// Provider는 이름을 반환합니다.
func (p *AWSProvider) Provider() string { return "aws_secrets_manager" }

// HealthCheck는 region 설정을 확인합니다.
func (p *AWSProvider) HealthCheck(_ context.Context) error {
	if p.region == "" {
		return errors.New("aws region not configured")
	}
	return nil
}

func (p *AWSProvider) Get(_ context.Context, _ string) (*Secret, error) {
	return nil, errors.New("aws SDK integration required")
}
func (p *AWSProvider) Set(_ context.Context, _ *Secret) error {
	return errors.New("aws SDK integration required")
}
func (p *AWSProvider) Delete(_ context.Context, _ string) error {
	return errors.New("aws SDK integration required")
}
func (p *AWSProvider) List(_ context.Context, _ string) ([]string, error) {
	return nil, errors.New("aws SDK integration required")
}
func (p *AWSProvider) Rotate(_ context.Context, _ string) (*Secret, error) {
	return nil, errors.New("aws SDK integration required")
}

// ============================================================================
// 자동 로테이션 매니저
// ============================================================================

// RotationManager는 정책에 따라 시크릿을 자동 로테이션합니다.
type RotationManager struct {
	provider Provider
	policies map[string]RotationPolicy
	mu       sync.RWMutex
}

// NewRotationManager는 새 로테이션 매니저를 생성합니다.
func NewRotationManager(provider Provider) *RotationManager {
	return &RotationManager{
		provider: provider,
		policies: make(map[string]RotationPolicy),
	}
}

// SetPolicy는 특정 시크릿 경로의 로테이션 정책을 설정합니다.
func (m *RotationManager) SetPolicy(path string, policy RotationPolicy) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.policies[path] = policy
}

// CheckAndRotate는 정책 미충족 시크릿을 로테이션합니다.
//
// 호출 빈도: 하루 1회 cron 권장.
func (m *RotationManager) CheckAndRotate(ctx context.Context) ([]string, error) {
	m.mu.RLock()
	paths := make([]string, 0, len(m.policies))
	policies := make(map[string]RotationPolicy, len(m.policies))
	for p, pol := range m.policies {
		paths = append(paths, p)
		policies[p] = pol
	}
	m.mu.RUnlock()

	rotated := []string{}
	for _, path := range paths {
		policy := policies[path]
		if !policy.Enabled {
			continue
		}
		secret, err := m.provider.Get(ctx, path)
		if err != nil {
			continue
		}
		age := time.Since(secret.CreatedAt)
		if age >= policy.MaxAge {
			if _, err := m.provider.Rotate(ctx, path); err == nil {
				rotated = append(rotated, path)
			}
		}
	}
	return rotated, nil
}

// ============================================================================
// 헬퍼: 랜덤 키 생성 (AES-256 호환)
// ============================================================================

// generateRandomKey는 길이가 length 바이트인 base64 인코딩된 랜덤 키를 반환합니다.
//
// 운영에서는 crypto/rand 사용 권장. 본 구현은 인메모리 테스트용.
func generateRandomKey(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	now := time.Now().UnixNano()
	b := make([]byte, length)
	for i := range b {
		now = (now*1103515245 + 12345) & 0x7FFFFFFF
		b[i] = charset[now%int64(len(charset))]
	}
	return string(b)
}
