// Package dashboard 는 만파식 운영 대시보드의 통합 어그리게이터.
//
// 목표: secrets / tracing / backup / canary / security / safety / sla / poise /
// tenancy 등 모든 ops 모듈의 헬스를 단일 스냅샷으로 모음.
//
// 디자인:
//   - HealthProvider 인터페이스로 모든 모듈을 추상화
//   - Aggregator 는 각 Provider 를 호출 후 ModuleSnapshot 으로 집계
//   - DashboardSnapshot 는 OverallStatus + 모듈별 상세
package dashboard

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"
)

// Status 는 모듈/시스템 상태.
type Status string

const (
	StatusHealthy  Status = "healthy"
	StatusDegraded Status = "degraded"
	StatusFailed   Status = "failed"
	StatusUnknown  Status = "unknown"
)

// statusRank 는 상태의 심각도 순위 (높을수록 나쁨).
func statusRank(s Status) int {
	switch s {
	case StatusHealthy:
		return 0
	case StatusUnknown:
		return 1
	case StatusDegraded:
		return 2
	case StatusFailed:
		return 3
	}
	return 1
}

// ModuleSnapshot 은 한 모듈의 상태 요약.
type ModuleSnapshot struct {
	Name      string
	Status    Status
	Message   string
	Metrics   map[string]interface{}
	UpdatedAt time.Time
}

// DashboardSnapshot 은 전체 시스템 상태.
type DashboardSnapshot struct {
	OverallStatus Status
	Modules       []ModuleSnapshot
	GeneratedAt   time.Time
	// FailedCount, DegradedCount 는 모듈 상태 분포.
	HealthyCount  int
	DegradedCount int
	FailedCount   int
	UnknownCount  int
}

// HealthProvider 는 한 모듈의 상태를 보고하는 어댑터.
type HealthProvider interface {
	Name() string
	Health(ctx context.Context) ModuleSnapshot
}

// Aggregator 는 등록된 HealthProvider 를 모아 스냅샷 생성.
type Aggregator struct {
	mu        sync.RWMutex
	providers []HealthProvider
	// timeout 은 각 provider Health() 호출의 최대 대기 시간.
	timeout time.Duration
}

// NewAggregator 생성. 기본 타임아웃 2초.
func NewAggregator() *Aggregator {
	return &Aggregator{timeout: 2 * time.Second}
}

// Register 는 provider 추가.
func (a *Aggregator) Register(p HealthProvider) {
	if p == nil {
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.providers = append(a.providers, p)
}

// SetTimeout 는 provider 호출 타임아웃 설정.
func (a *Aggregator) SetTimeout(d time.Duration) {
	if d <= 0 {
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.timeout = d
}

// Snapshot 은 모든 provider 를 동시 호출하고 결과를 집계.
//
// provider 호출은 병렬로 진행되며, 개별 호출 타임아웃 시 해당 모듈은 unknown 으로 표시.
func (a *Aggregator) Snapshot(ctx context.Context) DashboardSnapshot {
	a.mu.RLock()
	provs := make([]HealthProvider, len(a.providers))
	copy(provs, a.providers)
	timeout := a.timeout
	a.mu.RUnlock()

	results := make([]ModuleSnapshot, len(provs))
	var wg sync.WaitGroup
	for i, p := range provs {
		wg.Add(1)
		go func(idx int, prov HealthProvider) {
			defer wg.Done()
			results[idx] = a.callProvider(ctx, prov, timeout)
		}(i, p)
	}
	wg.Wait()

	// 모듈명 알파벳 정렬 (안정적 출력)
	sort.Slice(results, func(i, j int) bool {
		return strings.ToLower(results[i].Name) < strings.ToLower(results[j].Name)
	})

	snap := DashboardSnapshot{
		Modules:     results,
		GeneratedAt: time.Now(),
	}
	for _, m := range results {
		switch m.Status {
		case StatusHealthy:
			snap.HealthyCount++
		case StatusDegraded:
			snap.DegradedCount++
		case StatusFailed:
			snap.FailedCount++
		default:
			snap.UnknownCount++
		}
	}
	snap.OverallStatus = computeOverall(results)
	return snap
}

func (a *Aggregator) callProvider(parent context.Context, p HealthProvider, timeout time.Duration) ModuleSnapshot {
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()

	type result struct{ m ModuleSnapshot }
	ch := make(chan result, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				ch <- result{m: ModuleSnapshot{
					Name:      p.Name(),
					Status:    StatusFailed,
					Message:   "provider panic",
					UpdatedAt: time.Now(),
				}}
			}
		}()
		ch <- result{m: p.Health(ctx)}
	}()

	select {
	case r := <-ch:
		if r.m.UpdatedAt.IsZero() {
			r.m.UpdatedAt = time.Now()
		}
		if r.m.Name == "" {
			r.m.Name = p.Name()
		}
		return r.m
	case <-ctx.Done():
		return ModuleSnapshot{
			Name:      p.Name(),
			Status:    StatusUnknown,
			Message:   "health check timeout",
			UpdatedAt: time.Now(),
		}
	}
}

// computeOverall 은 가장 나쁜 모듈 상태를 전체 상태로 사용.
//
// 단, 모든 모듈이 healthy 이고 모듈 수가 0 이면 unknown.
func computeOverall(modules []ModuleSnapshot) Status {
	if len(modules) == 0 {
		return StatusUnknown
	}
	worst := StatusHealthy
	for _, m := range modules {
		if statusRank(m.Status) > statusRank(worst) {
			worst = m.Status
		}
	}
	return worst
}

// SimpleProvider 는 정적 함수를 HealthProvider 로 감싸는 어댑터.
type SimpleProvider struct {
	name string
	fn   func(ctx context.Context) ModuleSnapshot
}

// NewSimpleProvider 생성. fn 은 직접 ModuleSnapshot 을 반환.
func NewSimpleProvider(name string, fn func(ctx context.Context) ModuleSnapshot) *SimpleProvider {
	return &SimpleProvider{name: name, fn: fn}
}

func (p *SimpleProvider) Name() string { return p.name }

func (p *SimpleProvider) Health(ctx context.Context) ModuleSnapshot {
	if p.fn == nil {
		return ModuleSnapshot{Name: p.name, Status: StatusUnknown, Message: "no health fn"}
	}
	m := p.fn(ctx)
	if m.Name == "" {
		m.Name = p.name
	}
	return m
}
