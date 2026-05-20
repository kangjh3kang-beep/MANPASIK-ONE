package tenancy

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// InvitationCleaner 는 만료된 pending 초대를 자동으로 status=expired 로 갱신.
//
// MemoryInvitationStore / PostgresInvitationStore 모두에서 동작. 주기적 polling
// 으로 모든 조직의 초대를 스캔하므로 대규모 운영 환경에서는 별도 인덱스/배치
// 쿼리 최적화가 필요할 수 있음 (현재는 단순 + 정확).
//
// race-safe goroutine 종료 (Phase X 의 FeedbackLoop 패턴 재사용).
type InvitationCleaner struct {
	store    InvitationStore
	tenants  TenantLister // 옵션 — nil 이면 처리 대상이 없음
	interval time.Duration

	mu     sync.Mutex
	stopCh chan struct{}
	doneCh chan struct{}

	// stats (atomic for race-safe read)
	cycleCount   int64
	expiredCount int64

	// onError 는 정리 중 에러 보고 콜백 (옵션).
	onError func(err error)

	// nowFn 은 시간 소스 — 테스트 시 주입 가능.
	nowFn func() time.Time
}

// TenantLister 는 InvitationCleaner 가 스캔할 조직 ID 목록 제공자.
//
// MembershipStore 가 ListAll 같은 메서드를 갖지 않으므로 별도 추상화. 운영에서는
// 멤버십 store 의 ListUserTenants/ListTenantMembers 를 활용하거나 별도 캐시
// 사용 가능.
type TenantLister interface {
	AllTenants() []TenantID
}

// TenantListerFunc 는 함수를 TenantLister 로 감쌈.
type TenantListerFunc func() []TenantID

// AllTenants 호출.
func (f TenantListerFunc) AllTenants() []TenantID { return f() }

// MembershipBackedTenantLister 는 MembershipStore 를 TenantLister 로 사용.
//
// 모든 멤버십을 스캔하여 고유한 tenantID 만 반환. MemoryMembershipStore /
// PostgresMembershipStore 의 ListUserTenants/ListTenantMembers 만으로는
// 전체 조직 목록을 알 수 없으므로 별도 메서드 필요.
type MembershipBackedTenantLister struct {
	store MembershipStore
	// 모든 사용자 ID 를 알기 위한 함수 (운영에서는 user-service 와 통합).
	// nil 이면 빈 목록 반환.
	allUserIDs func() []string
}

// NewMembershipBackedTenantLister 생성.
func NewMembershipBackedTenantLister(store MembershipStore,
	allUserIDs func() []string) *MembershipBackedTenantLister {
	return &MembershipBackedTenantLister{store: store, allUserIDs: allUserIDs}
}

// AllTenants 는 모든 사용자의 멤버십을 스캔하여 고유 tenant 추출.
func (l *MembershipBackedTenantLister) AllTenants() []TenantID {
	if l == nil || l.store == nil || l.allUserIDs == nil {
		return nil
	}
	seen := make(map[TenantID]struct{})
	for _, uid := range l.allUserIDs() {
		for _, m := range l.store.ListUserTenants(uid) {
			seen[m.TenantID] = struct{}{}
		}
	}
	out := make([]TenantID, 0, len(seen))
	for t := range seen {
		out = append(out, t)
	}
	return out
}

// InvitationCleanerConfig 는 cleaner 동작 설정.
type InvitationCleanerConfig struct {
	// Interval 은 정리 주기 (기본 1시간).
	Interval time.Duration
	// OnError 는 store 에러 콜백.
	OnError func(err error)
	// Now 는 테스트용 시간 소스.
	Now func() time.Time
}

// NewInvitationCleaner 생성. tenants=nil 이면 cleaner 가 모든 사이클을 무사히
// 통과하지만 아무 것도 하지 않음 (MarkExpired 직접 호출 시에만 동작).
func NewInvitationCleaner(store InvitationStore, tenants TenantLister,
	cfg InvitationCleanerConfig) (*InvitationCleaner, error) {
	if store == nil {
		return nil, errors.New("store 필수")
	}
	if cfg.Interval <= 0 {
		cfg.Interval = 1 * time.Hour
	}
	if cfg.Now == nil {
		cfg.Now = time.Now
	}
	return &InvitationCleaner{
		store:    store,
		tenants:  tenants,
		interval: cfg.Interval,
		onError:  cfg.OnError,
		nowFn:    cfg.Now,
	}, nil
}

// MarkExpired 는 단일 조직의 pending 초대 중 만료된 것을 expired 로 갱신.
//
// 멱등 동작 — 이미 expired/accepted/revoked 인 초대는 변경 없음.
// 반환값: 갱신된 초대 개수.
func (c *InvitationCleaner) MarkExpired(tenantID TenantID) int {
	if c.store == nil {
		return 0
	}
	now := c.nowFn()
	list := c.store.ListByTenant(tenantID)
	expired := 0
	for _, inv := range list {
		if inv.Status != InvitationPending {
			continue
		}
		if !inv.ExpiresAt.Before(now) {
			continue
		}
		inv.Status = InvitationExpired
		if err := c.store.Update(*inv); err != nil {
			c.reportErr(err)
			continue
		}
		expired++
	}
	atomic.AddInt64(&c.expiredCount, int64(expired))
	return expired
}

// MarkExpiredAll 는 등록된 TenantLister 의 모든 조직을 순회하며 정리.
func (c *InvitationCleaner) MarkExpiredAll() int {
	atomic.AddInt64(&c.cycleCount, 1)
	if c.tenants == nil {
		return 0
	}
	total := 0
	for _, tid := range c.tenants.AllTenants() {
		total += c.MarkExpired(tid)
	}
	return total
}

// Start 는 백그라운드 goroutine 으로 주기적 cleanup 시작.
func (c *InvitationCleaner) Start(ctx context.Context) {
	c.mu.Lock()
	if c.stopCh != nil {
		c.mu.Unlock()
		return
	}
	stopCh := make(chan struct{})
	doneCh := make(chan struct{})
	c.stopCh = stopCh
	c.doneCh = doneCh
	c.mu.Unlock()
	go c.run(ctx, stopCh, doneCh)
}

// Stop 은 cleaner 종료 + 현재 사이클 완료까지 대기.
func (c *InvitationCleaner) Stop() {
	c.mu.Lock()
	if c.stopCh == nil {
		c.mu.Unlock()
		return
	}
	close(c.stopCh)
	doneCh := c.doneCh
	c.stopCh = nil
	c.mu.Unlock()
	if doneCh != nil {
		<-doneCh
	}
}

// CycleCount 는 지금까지 실행된 정리 사이클 수.
func (c *InvitationCleaner) CycleCount() int64 {
	return atomic.LoadInt64(&c.cycleCount)
}

// ExpiredCount 는 누적 만료 처리 개수.
func (c *InvitationCleaner) ExpiredCount() int64 {
	return atomic.LoadInt64(&c.expiredCount)
}

// SetTenants 는 런타임 중 TenantLister 교체.
func (c *InvitationCleaner) SetTenants(tenants TenantLister) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.tenants = tenants
}

func (c *InvitationCleaner) run(parentCtx context.Context, stopCh, doneCh chan struct{}) {
	defer close(doneCh)
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	// 첫 사이클은 즉시 실행
	_ = c.MarkExpiredAll()
	for {
		select {
		case <-stopCh:
			return
		case <-parentCtx.Done():
			return
		case <-ticker.C:
			_ = c.MarkExpiredAll()
		}
	}
}

func (c *InvitationCleaner) reportErr(err error) {
	if c.onError != nil {
		c.onError(err)
	}
}
