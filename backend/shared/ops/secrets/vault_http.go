// Package secrets: HashiCorp Vault HTTP API 실 통합 어댑터
//
// KV v2 시크릿 엔진 + 토큰 자동 갱신 + 네임스페이스 지원.
// 외부 SDK 의존성 없이 net/http로 직접 구현하여 코어 가벼움 유지.
package secrets

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// ============================================================================
// HTTP 클라이언트 인터페이스 (테스트 가능성)
// ============================================================================

// HTTPDoer는 테스트 가능한 HTTP 클라이언트 인터페이스입니다.
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// ============================================================================
// VaultHTTPProvider — KV v2 실 통합
// ============================================================================

// VaultHTTPProvider는 Vault KV v2 HTTP API 어댑터입니다.
//
// 요청 형식 (KV v2):
//   GET    /v1/{mount}/data/{path}     → 읽기
//   POST   /v1/{mount}/data/{path}     → 쓰기
//   DELETE /v1/{mount}/metadata/{path} → 영구 삭제
//   LIST   /v1/{mount}/metadata/{path} → 목록 (HTTP 'LIST' 메서드)
type VaultHTTPProvider struct {
	mu          sync.RWMutex
	addr        string
	token       string
	namespace   string
	mount       string // 기본 "secret"
	httpClient  HTTPDoer
	tokenExpiry time.Time
}

// NewVaultHTTPProvider는 새 Vault HTTP Provider를 생성합니다.
//
// addr 예: "https://vault.example.com:8200"
// mount 예: "secret" (KV v2 기본 마운트)
func NewVaultHTTPProvider(addr, token, namespace, mount string) *VaultHTTPProvider {
	if mount == "" {
		mount = "secret"
	}
	return &VaultHTTPProvider{
		addr:       strings.TrimSuffix(addr, "/"),
		token:      token,
		namespace:  namespace,
		mount:      mount,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// SetHTTPClient는 테스트용 HTTP 클라이언트를 설정합니다.
func (p *VaultHTTPProvider) SetHTTPClient(c HTTPDoer) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.httpClient = c
}

// Provider는 이름을 반환합니다.
func (p *VaultHTTPProvider) Provider() string { return "vault_http" }

// HealthCheck는 Vault sys/health 엔드포인트를 호출합니다.
func (p *VaultHTTPProvider) HealthCheck(ctx context.Context) error {
	if p.addr == "" {
		return errors.New("vault addr not configured")
	}
	if p.token == "" {
		return errors.New("vault token not configured")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", p.addr+"/v1/sys/health", nil)
	if err != nil {
		return err
	}
	p.applyHeaders(req)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("vault health failed: %w", err)
	}
	defer resp.Body.Close()

	// Vault 활성: 200 (sealed=false, standby=false)
	// 200/429/472/473/501/503 모두 응답이지만 healthy 판단 기준은 200/standby OK
	if resp.StatusCode == 200 || resp.StatusCode == 429 {
		return nil
	}
	return fmt.Errorf("vault unhealthy: HTTP %d", resp.StatusCode)
}

// applyHeaders는 토큰/네임스페이스 헤더를 추가합니다.
func (p *VaultHTTPProvider) applyHeaders(req *http.Request) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	req.Header.Set("X-Vault-Token", p.token)
	if p.namespace != "" {
		req.Header.Set("X-Vault-Namespace", p.namespace)
	}
}

// ============================================================================
// KV v2 Read/Write
// ============================================================================

type kvReadResponse struct {
	Data struct {
		Data     map[string]interface{} `json:"data"`
		Metadata struct {
			CreatedTime  time.Time `json:"created_time"`
			DeletionTime time.Time `json:"deletion_time"`
			Destroyed    bool      `json:"destroyed"`
			Version      int       `json:"version"`
		} `json:"metadata"`
	} `json:"data"`
}

type kvWriteRequest struct {
	Data    map[string]interface{} `json:"data"`
	Options map[string]interface{} `json:"options,omitempty"`
}

// Get은 KV v2에서 시크릿을 조회합니다.
func (p *VaultHTTPProvider) Get(ctx context.Context, path string) (*Secret, error) {
	if path == "" {
		return nil, errors.New("path required")
	}
	endpoint := fmt.Sprintf("%s/v1/%s/data/%s", p.addr, p.mount, url.PathEscape(path))

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	p.applyHeaders(req)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("vault get failed: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 404:
		return nil, fmt.Errorf("secret %q not found", path)
	case 403:
		return nil, errors.New("vault: permission denied (token may be invalid/expired)")
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("vault get HTTP %d", resp.StatusCode)
	}

	var kvResp kvReadResponse
	if err := json.NewDecoder(resp.Body).Decode(&kvResp); err != nil {
		return nil, fmt.Errorf("vault decode: %w", err)
	}

	// "value" 키 우선, 없으면 첫 번째 문자열 값
	value := ""
	if v, ok := kvResp.Data.Data["value"].(string); ok {
		value = v
	} else {
		for _, v := range kvResp.Data.Data {
			if s, ok := v.(string); ok {
				value = s
				break
			}
		}
	}

	metadata := make(map[string]string)
	for k, v := range kvResp.Data.Data {
		if k == "value" {
			continue
		}
		if s, ok := v.(string); ok {
			metadata[k] = s
		}
	}

	return &Secret{
		Path:      path,
		Value:     value,
		Version:   kvResp.Data.Metadata.Version,
		Metadata:  metadata,
		CreatedAt: kvResp.Data.Metadata.CreatedTime,
	}, nil
}

// Set은 KV v2에 시크릿을 저장합니다.
func (p *VaultHTTPProvider) Set(ctx context.Context, secret *Secret) error {
	if secret == nil || secret.Path == "" {
		return errors.New("path required")
	}

	data := map[string]interface{}{"value": secret.Value}
	for k, v := range secret.Metadata {
		data[k] = v
	}

	body := kvWriteRequest{Data: data}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("%s/v1/%s/data/%s", p.addr, p.mount, url.PathEscape(secret.Path))
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	p.applyHeaders(req)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("vault set failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("vault set HTTP %d", resp.StatusCode)
	}
	return nil
}

// Delete는 KV v2에서 시크릿 메타데이터(모든 버전)를 영구 삭제합니다.
func (p *VaultHTTPProvider) Delete(ctx context.Context, path string) error {
	endpoint := fmt.Sprintf("%s/v1/%s/metadata/%s", p.addr, p.mount, url.PathEscape(path))
	req, err := http.NewRequestWithContext(ctx, "DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	p.applyHeaders(req)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("vault delete failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 && resp.StatusCode != 404 {
		return fmt.Errorf("vault delete HTTP %d", resp.StatusCode)
	}
	return nil
}

// List는 prefix 하위 경로 목록을 반환합니다.
//
// Vault LIST는 HTTP 'LIST' 메서드 또는 GET ?list=true.
func (p *VaultHTTPProvider) List(ctx context.Context, prefix string) ([]string, error) {
	endpoint := fmt.Sprintf("%s/v1/%s/metadata/%s?list=true",
		p.addr, p.mount, url.PathEscape(prefix))
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	p.applyHeaders(req)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("vault list failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return []string{}, nil // 빈 prefix
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("vault list HTTP %d", resp.StatusCode)
	}

	var listResp struct {
		Data struct {
			Keys []string `json:"keys"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, fmt.Errorf("vault list decode: %w", err)
	}

	// prefix를 붙여 절대 경로로 반환
	out := make([]string, 0, len(listResp.Data.Keys))
	for _, k := range listResp.Data.Keys {
		out = append(out, prefix+k)
	}
	return out, nil
}

// Rotate는 시크릿을 새 랜덤 값으로 회전합니다.
func (p *VaultHTTPProvider) Rotate(ctx context.Context, path string) (*Secret, error) {
	existing, err := p.Get(ctx, path)
	if err != nil {
		return nil, err
	}
	rotated := &Secret{
		Path:      path,
		Value:     generateRandomKey(32),
		Metadata:  existing.Metadata,
		CreatedAt: time.Now().UTC(),
	}
	if err := p.Set(ctx, rotated); err != nil {
		return nil, err
	}
	rotated.Version = existing.Version + 1
	return rotated, nil
}

// ============================================================================
// 토큰 자동 갱신 (TTL 기반)
// ============================================================================

// TokenAutoRenewer는 백그라운드에서 토큰을 자동 갱신합니다.
//
// Vault 토큰은 TTL을 가지며 만료되면 모든 호출이 403 실패.
// renewBeforeExpiry 안에서 sys/auth/token/renew-self 호출하여 갱신.
type TokenAutoRenewer struct {
	provider          *VaultHTTPProvider
	renewBeforeExpiry time.Duration
	stopCh            chan struct{}
	mu                sync.Mutex
	running           bool
}

// NewTokenAutoRenewer는 새 갱신기를 생성합니다.
func NewTokenAutoRenewer(provider *VaultHTTPProvider, renewBefore time.Duration) *TokenAutoRenewer {
	if renewBefore <= 0 {
		renewBefore = 5 * time.Minute
	}
	return &TokenAutoRenewer{
		provider:          provider,
		renewBeforeExpiry: renewBefore,
		stopCh:            make(chan struct{}),
	}
}

// RenewSelf는 토큰을 즉시 갱신합니다.
func (r *TokenAutoRenewer) RenewSelf(ctx context.Context) (time.Duration, error) {
	endpoint := r.provider.addr + "/v1/auth/token/renew-self"
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, nil)
	if err != nil {
		return 0, err
	}
	r.provider.applyHeaders(req)

	resp, err := r.provider.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return 0, fmt.Errorf("renew-self HTTP %d", resp.StatusCode)
	}

	var renewResp struct {
		Auth struct {
			LeaseDurationSec int  `json:"lease_duration"`
			Renewable        bool `json:"renewable"`
		} `json:"auth"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&renewResp); err != nil {
		return 0, err
	}

	leaseDur := time.Duration(renewResp.Auth.LeaseDurationSec) * time.Second
	r.provider.mu.Lock()
	r.provider.tokenExpiry = time.Now().UTC().Add(leaseDur)
	r.provider.mu.Unlock()
	return leaseDur, nil
}

// Start는 백그라운드 자동 갱신을 시작합니다.
func (r *TokenAutoRenewer) Start(ctx context.Context, renewIntervalSeconds int) error {
	r.mu.Lock()
	if r.running {
		r.mu.Unlock()
		return errors.New("already running")
	}
	r.running = true
	r.mu.Unlock()

	if renewIntervalSeconds <= 0 {
		renewIntervalSeconds = 1800 // 30분 기본
	}

	go func() {
		ticker := time.NewTicker(time.Duration(renewIntervalSeconds) * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-r.stopCh:
				return
			case <-ctx.Done():
				return
			case <-ticker.C:
				_, _ = r.RenewSelf(ctx)
			}
		}
	}()
	return nil
}

// Stop은 자동 갱신을 종료합니다.
func (r *TokenAutoRenewer) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.running {
		return
	}
	close(r.stopCh)
	r.running = false
}

// IsRunning은 자동 갱신 진행 여부를 반환합니다.
func (r *TokenAutoRenewer) IsRunning() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.running
}
