package secrets

import (
	"context"
	"errors"
	"sync"
	"time"
)

// RotationEvent 는 시크릿이 회전(변경)되었을 때 리스너에게 전달되는 이벤트.
type RotationEvent struct {
	Path      string
	NewSecret *Secret
	OldSecret *Secret // 첫 관찰 시 nil
	Timestamp time.Time
}

// RotationListener 는 시크릿 변경 알림 콜백.
//
// 콜백 안에서 길게 블로킹하면 Watcher 의 다른 path 폴링이 지연되므로
// 비동기 처리 권장. 에러 반환 시 Watcher 가 OnError 콜백으로 보고.
type RotationListener func(ctx context.Context, ev RotationEvent) error

// RotationWatcherConfig 는 Watcher 동작 설정.
type RotationWatcherConfig struct {
	// Interval 은 폴링 주기 (기본 60초). 짧을수록 변경 감지가 빠르지만 부하 증가.
	Interval time.Duration
	// CompareBy 는 변경 판정 기준.
	//   "version" (기본): Secret.Version 비교
	//   "value":   Secret.Value 비교 (Provider가 Version 미지원 시)
	CompareBy string
	// OnError 는 폴링/리스너 호출 중 에러 발생 시 콜백 (로깅용).
	OnError func(path string, err error)
}

// RotationWatcher 는 등록된 path 들을 주기적으로 폴링하여 변경을 감지하고
// 모든 리스너에게 알림.
//
// 사용 예:
//
//	w := secrets.NewRotationWatcher(provider, secrets.RotationWatcherConfig{
//	    Interval: 30*time.Second, CompareBy: "version",
//	})
//	w.Watch("database/password")
//	w.Watch("auth/jwt_key")
//	w.AddListener(func(ctx context.Context, ev secrets.RotationEvent) error {
//	    log.Printf("secret rotated: %s v%d → v%d",
//	        ev.Path, ev.OldSecret.Version, ev.NewSecret.Version)
//	    return reloadAuthKey(ev.NewSecret.Value)
//	})
//	w.Start(ctx)
//	defer w.Stop()
type RotationWatcher struct {
	provider Provider
	cfg      RotationWatcherConfig

	mu        sync.Mutex
	paths     map[string]bool
	cache     map[string]*Secret
	listeners []RotationListener
	stopCh    chan struct{}
	doneCh    chan struct{}
	cycles    int
}

// NewRotationWatcher 생성. provider=nil 이면 panic 회피용 nil 반환.
func NewRotationWatcher(provider Provider, cfg RotationWatcherConfig) (*RotationWatcher, error) {
	if provider == nil {
		return nil, errors.New("provider 필수")
	}
	if cfg.Interval <= 0 {
		cfg.Interval = 60 * time.Second
	}
	if cfg.CompareBy == "" {
		cfg.CompareBy = "version"
	}
	return &RotationWatcher{
		provider: provider,
		cfg:      cfg,
		paths:    make(map[string]bool),
		cache:    make(map[string]*Secret),
	}, nil
}

// Watch 는 감시할 path 추가 (이미 등록되어 있으면 no-op).
func (w *RotationWatcher) Watch(path string) {
	if path == "" {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	w.paths[path] = true
}

// Unwatch 는 path 제거 (캐시도 삭제).
func (w *RotationWatcher) Unwatch(path string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.paths, path)
	delete(w.cache, path)
}

// AddListener 는 변경 알림 콜백 등록 (멱등 보장 안 됨 — 동일 함수 두 번 등록 시 두 번 호출).
func (w *RotationWatcher) AddListener(fn RotationListener) {
	if fn == nil {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	w.listeners = append(w.listeners, fn)
}

// Start 는 백그라운드 polling goroutine 실행. 이미 실행 중이면 no-op.
func (w *RotationWatcher) Start(ctx context.Context) {
	w.mu.Lock()
	if w.stopCh != nil {
		w.mu.Unlock()
		return
	}
	stopCh := make(chan struct{})
	doneCh := make(chan struct{})
	w.stopCh = stopCh
	w.doneCh = doneCh
	w.mu.Unlock()

	go w.run(ctx, stopCh, doneCh)
}

// Stop 은 polling 종료 + 현재 사이클 완료까지 대기.
func (w *RotationWatcher) Stop() {
	w.mu.Lock()
	if w.stopCh == nil {
		w.mu.Unlock()
		return
	}
	close(w.stopCh)
	doneCh := w.doneCh
	w.stopCh = nil
	w.mu.Unlock()
	if doneCh != nil {
		<-doneCh
	}
}

// CycleCount 는 지금까지 실행된 폴링 사이클 수.
func (w *RotationWatcher) CycleCount() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.cycles
}

// CheckOnce 는 현재 등록된 모든 path 를 한 번만 폴링하고 변경 감지 + 알림.
//
// Start 없이 수동 호출 가능 (테스트/일회성 갱신).
func (w *RotationWatcher) CheckOnce(ctx context.Context) {
	w.mu.Lock()
	w.cycles++
	pathList := make([]string, 0, len(w.paths))
	for p := range w.paths {
		pathList = append(pathList, p)
	}
	listeners := append([]RotationListener(nil), w.listeners...)
	w.mu.Unlock()

	for _, p := range pathList {
		w.checkPath(ctx, p, listeners)
	}
}

func (w *RotationWatcher) run(parentCtx context.Context, stopCh, doneCh chan struct{}) {
	defer close(doneCh)
	ticker := time.NewTicker(w.cfg.Interval)
	defer ticker.Stop()

	w.CheckOnce(parentCtx)
	for {
		select {
		case <-stopCh:
			return
		case <-parentCtx.Done():
			return
		case <-ticker.C:
			w.CheckOnce(parentCtx)
		}
	}
}

func (w *RotationWatcher) checkPath(ctx context.Context, path string, listeners []RotationListener) {
	cycleCtx, cancel := context.WithTimeout(ctx, w.cfg.Interval)
	defer cancel()

	cur, err := w.provider.Get(cycleCtx, path)
	if err != nil {
		w.reportErr(path, err)
		return
	}
	if cur == nil {
		return
	}

	w.mu.Lock()
	prev := w.cache[path]
	changed := w.isChanged(prev, cur)
	w.cache[path] = cur
	w.mu.Unlock()

	if !changed {
		return
	}

	ev := RotationEvent{
		Path:      path,
		NewSecret: cur,
		OldSecret: prev, // 첫 관찰 시 nil
		Timestamp: time.Now().UTC(),
	}
	for _, fn := range listeners {
		if err := fn(ctx, ev); err != nil {
			w.reportErr(path, err)
		}
	}
}

func (w *RotationWatcher) isChanged(prev, cur *Secret) bool {
	if prev == nil {
		// 첫 관찰: 변경 이벤트로 처리하지 않음 (baseline 만 설정)
		// 이 동작은 운영에서 디바운스 효과 — 첫 호출 시 모든 시크릿이 "변경"으로
		// 보이는 것을 방지.
		return false
	}
	if cur == nil {
		return false
	}
	switch w.cfg.CompareBy {
	case "value":
		return prev.Value != cur.Value
	default: // "version"
		return prev.Version != cur.Version
	}
}

func (w *RotationWatcher) reportErr(path string, err error) {
	if w.cfg.OnError != nil {
		w.cfg.OnError(path, err)
	}
}
