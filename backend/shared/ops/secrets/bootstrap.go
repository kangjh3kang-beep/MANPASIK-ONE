package secrets

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

// VaultBootstrapConfig 는 Vault + RotationWatcher 결합 설정.
type VaultBootstrapConfig struct {
	// Addr 은 Vault 주소 (예: "http://vault:8200"). 비어있으면 ENV.VAULT_ADDR 폴백.
	Addr string
	// Token 은 Vault 토큰. 비어있으면 ENV.VAULT_TOKEN 폴백.
	Token string
	// Namespace 는 Vault Enterprise namespace (선택).
	Namespace string
	// Mount 는 KV v2 마운트 경로 (기본 "secret"). 비어있으면 ENV.VAULT_KV_MOUNT.
	Mount string
	// WatchPaths 는 자동 감시할 시크릿 경로 목록.
	WatchPaths []string
	// WatchInterval 은 폴링 주기 (기본 60초).
	WatchInterval time.Duration
	// CompareBy 는 변경 판정 기준 ("version" 또는 "value", 기본 "version").
	CompareBy string
	// AutoRenewToken=true 시 TokenAutoRenewer 함께 시작 (기본 true).
	AutoRenewToken bool
	// TokenRenewIntervalSeconds 는 토큰 갱신 주기 (기본 1800 = 30분).
	TokenRenewIntervalSeconds int
	// TokenRenewBefore 는 만료 N 시간 전 갱신 (기본 5분).
	TokenRenewBefore time.Duration
	// OnError 는 시크릿 폴링/리스너 에러 콜백.
	OnError func(path string, err error)
}

// VaultBootstrap 은 Vault Provider + RotationWatcher + (선택) TokenAutoRenewer
// 를 묶은 운영용 헬퍼.
//
// 서비스 main.go 에서 한 번 생성하면 시크릿 폴링/토큰 갱신이 백그라운드로 실행되며,
// 등록된 RotationListener 가 시크릿 변경 시 자동 호출됨.
type VaultBootstrap struct {
	provider                  *VaultHTTPProvider
	watcher                   *RotationWatcher
	renewer                   *TokenAutoRenewer
	tokenRenewIntervalSeconds int
}

// NewVaultBootstrap 은 cfg 로 Vault Provider 와 RotationWatcher 를 생성.
// AutoRenewToken=true 면 TokenAutoRenewer 도 함께 생성.
//
// 모든 컴포넌트는 생성만 되고 Start() 호출 전에는 polling 안 함.
func NewVaultBootstrap(cfg VaultBootstrapConfig) (*VaultBootstrap, error) {
	addr := firstNonEmpty(cfg.Addr, os.Getenv("VAULT_ADDR"))
	token := firstNonEmpty(cfg.Token, os.Getenv("VAULT_TOKEN"))
	mount := firstNonEmpty(cfg.Mount, os.Getenv("VAULT_KV_MOUNT"), "secret")

	if addr == "" {
		return nil, errors.New("Vault Addr 미설정 (VaultBootstrapConfig.Addr 또는 VAULT_ADDR)")
	}

	provider := NewVaultHTTPProvider(addr, token, cfg.Namespace, mount)
	if err := provider.HealthCheck(context.Background()); err != nil {
		return nil, fmt.Errorf("Vault 헬스체크 실패: %w", err)
	}

	wcfg := RotationWatcherConfig{
		Interval:  cfg.WatchInterval,
		CompareBy: cfg.CompareBy,
		OnError:   cfg.OnError,
	}
	watcher, err := NewRotationWatcher(provider, wcfg)
	if err != nil {
		return nil, err
	}
	for _, p := range cfg.WatchPaths {
		watcher.Watch(p)
	}

	bs := &VaultBootstrap{
		provider: provider,
		watcher:  watcher,
	}
	if cfg.AutoRenewToken && token != "" {
		bs.renewer = NewTokenAutoRenewer(provider, cfg.TokenRenewBefore)
	}
	bs.tokenRenewIntervalSeconds = cfg.TokenRenewIntervalSeconds
	return bs, nil
}

// Provider 반환 — 직접 시크릿 Get/Set 호출용.
func (b *VaultBootstrap) Provider() *VaultHTTPProvider { return b.provider }

// Watcher 반환 — Watch/AddListener 추가 호출용.
func (b *VaultBootstrap) Watcher() *RotationWatcher { return b.watcher }

// AddListener 는 RotationWatcher 에 리스너 위임 등록.
func (b *VaultBootstrap) AddListener(fn RotationListener) {
	b.watcher.AddListener(fn)
}

// Watch 는 새 path 를 RotationWatcher 에 추가.
func (b *VaultBootstrap) Watch(path string) {
	b.watcher.Watch(path)
}

// Start 는 watcher + (옵션) renewer 백그라운드 시작.
//
// renewer.Start 는 에러를 반환하지만 (이미 실행 중일 때만 발생) 무시하고 진행.
// 호출자가 정확한 상태를 알아야 한다면 직접 b.Renewer().Start() 호출 권장.
func (b *VaultBootstrap) Start(ctx context.Context) {
	b.watcher.Start(ctx)
	if b.renewer != nil {
		_ = b.renewer.Start(ctx, b.tokenRenewIntervalSeconds)
	}
}

// Renewer 반환 — 직접 RenewSelf 호출 또는 IsRunning 검사용.
func (b *VaultBootstrap) Renewer() *TokenAutoRenewer { return b.renewer }

// Stop 은 모든 백그라운드 작업 종료.
func (b *VaultBootstrap) Stop() {
	b.watcher.Stop()
	if b.renewer != nil {
		b.renewer.Stop()
	}
}

// firstNonEmpty 는 첫 비어있지 않은 문자열 반환.
func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
