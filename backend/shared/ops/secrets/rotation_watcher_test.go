package secrets_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/manpasik/backend/shared/ops/secrets"
)

// fakeProvider 는 테스트용 Provider 구현.
type fakeProvider struct {
	mu       sync.Mutex
	store    map[string]*secrets.Secret
	getErr   error
}

func newFakeProvider() *fakeProvider {
	return &fakeProvider{store: make(map[string]*secrets.Secret)}
}

func (f *fakeProvider) Get(_ context.Context, path string) (*secrets.Secret, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.getErr != nil {
		return nil, f.getErr
	}
	s, ok := f.store[path]
	if !ok {
		return nil, errors.New("not found")
	}
	cp := *s
	return &cp, nil
}

func (f *fakeProvider) Set(_ context.Context, s *secrets.Secret) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	cp := *s
	f.store[s.Path] = &cp
	return nil
}

func (f *fakeProvider) SetGetError(err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.getErr = err
}

func (f *fakeProvider) Delete(_ context.Context, path string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.store, path)
	return nil
}

func (f *fakeProvider) List(_ context.Context, _ string) ([]string, error) { return nil, nil }
func (f *fakeProvider) Rotate(_ context.Context, path string) (*secrets.Secret, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	s := f.store[path]
	s.Version++
	cp := *s
	return &cp, nil
}
func (f *fakeProvider) Provider() string                     { return "fake" }
func (f *fakeProvider) HealthCheck(_ context.Context) error { return nil }

func TestRotationWatcher_NoChangeOnFirstObservation(t *testing.T) {
	p := newFakeProvider()
	_ = p.Set(context.Background(), &secrets.Secret{Path: "k", Value: "v1", Version: 1})

	w, _ := secrets.NewRotationWatcher(p, secrets.RotationWatcherConfig{
		Interval: 100 * time.Millisecond,
	})
	w.Watch("k")

	var fires int32
	w.AddListener(func(_ context.Context, _ secrets.RotationEvent) error {
		atomic.AddInt32(&fires, 1)
		return nil
	})

	w.CheckOnce(context.Background())
	if atomic.LoadInt32(&fires) != 0 {
		t.Errorf("첫 관찰에 알림 발생: %d", fires)
	}
}

func TestRotationWatcher_VersionChange(t *testing.T) {
	p := newFakeProvider()
	_ = p.Set(context.Background(), &secrets.Secret{Path: "k", Value: "v1", Version: 1})

	w, _ := secrets.NewRotationWatcher(p, secrets.RotationWatcherConfig{
		Interval: 50 * time.Millisecond,
	})
	w.Watch("k")

	var (
		mu     sync.Mutex
		events []secrets.RotationEvent
	)
	w.AddListener(func(_ context.Context, ev secrets.RotationEvent) error {
		mu.Lock()
		defer mu.Unlock()
		events = append(events, ev)
		return nil
	})

	// baseline
	w.CheckOnce(context.Background())

	// 버전 증가
	_ = p.Set(context.Background(), &secrets.Secret{Path: "k", Value: "v2", Version: 2})
	w.CheckOnce(context.Background())

	mu.Lock()
	defer mu.Unlock()
	if len(events) != 1 {
		t.Fatalf("events = %d, want 1", len(events))
	}
	ev := events[0]
	if ev.Path != "k" {
		t.Errorf("Path = %q", ev.Path)
	}
	if ev.NewSecret.Version != 2 {
		t.Errorf("NewSecret.Version = %d", ev.NewSecret.Version)
	}
	if ev.OldSecret == nil || ev.OldSecret.Version != 1 {
		t.Errorf("OldSecret = %+v", ev.OldSecret)
	}
}

func TestRotationWatcher_CompareByValue(t *testing.T) {
	p := newFakeProvider()
	_ = p.Set(context.Background(), &secrets.Secret{Path: "k", Value: "v1", Version: 1})

	w, _ := secrets.NewRotationWatcher(p, secrets.RotationWatcherConfig{
		Interval:  50 * time.Millisecond,
		CompareBy: "value",
	})
	w.Watch("k")

	var fires int32
	w.AddListener(func(_ context.Context, _ secrets.RotationEvent) error {
		atomic.AddInt32(&fires, 1)
		return nil
	})

	w.CheckOnce(context.Background()) // baseline

	// Version 동일하지만 Value 변경
	_ = p.Set(context.Background(), &secrets.Secret{Path: "k", Value: "v2", Version: 1})
	w.CheckOnce(context.Background())

	if atomic.LoadInt32(&fires) != 1 {
		t.Errorf("value 변경 미감지: fires = %d", fires)
	}
}

func TestRotationWatcher_NoChangeNoFire(t *testing.T) {
	p := newFakeProvider()
	_ = p.Set(context.Background(), &secrets.Secret{Path: "k", Value: "v1", Version: 1})

	w, _ := secrets.NewRotationWatcher(p, secrets.RotationWatcherConfig{
		Interval: 50 * time.Millisecond,
	})
	w.Watch("k")
	var fires int32
	w.AddListener(func(_ context.Context, _ secrets.RotationEvent) error {
		atomic.AddInt32(&fires, 1)
		return nil
	})

	w.CheckOnce(context.Background())
	w.CheckOnce(context.Background())
	w.CheckOnce(context.Background())

	if atomic.LoadInt32(&fires) != 0 {
		t.Errorf("변경 없는데 알림 발생: %d", fires)
	}
}

func TestRotationWatcher_Unwatch(t *testing.T) {
	p := newFakeProvider()
	_ = p.Set(context.Background(), &secrets.Secret{Path: "k", Value: "v1", Version: 1})

	w, _ := secrets.NewRotationWatcher(p, secrets.RotationWatcherConfig{
		Interval: 50 * time.Millisecond,
	})
	w.Watch("k")
	w.CheckOnce(context.Background()) // baseline 저장

	w.Unwatch("k")
	// Unwatch 후 변경되어도 알림 없어야
	_ = p.Set(context.Background(), &secrets.Secret{Path: "k", Value: "v2", Version: 2})

	var fires int32
	w.AddListener(func(_ context.Context, _ secrets.RotationEvent) error {
		atomic.AddInt32(&fires, 1)
		return nil
	})
	w.CheckOnce(context.Background())
	if atomic.LoadInt32(&fires) != 0 {
		t.Errorf("Unwatch 후 알림: %d", fires)
	}
}

func TestRotationWatcher_OnError(t *testing.T) {
	p := newFakeProvider()
	p.SetGetError(errors.New("vault down"))

	var (
		mu   sync.Mutex
		errs []error
	)
	w, _ := secrets.NewRotationWatcher(p, secrets.RotationWatcherConfig{
		Interval: 50 * time.Millisecond,
		OnError: func(_ string, err error) {
			mu.Lock()
			defer mu.Unlock()
			errs = append(errs, err)
		},
	})
	w.Watch("k")
	w.CheckOnce(context.Background())

	mu.Lock()
	defer mu.Unlock()
	if len(errs) == 0 {
		t.Error("OnError 미호출")
	}
}

func TestRotationWatcher_StartStop(t *testing.T) {
	p := newFakeProvider()
	_ = p.Set(context.Background(), &secrets.Secret{Path: "k", Value: "v1", Version: 1})

	w, _ := secrets.NewRotationWatcher(p, secrets.RotationWatcherConfig{
		Interval: 30 * time.Millisecond,
	})
	w.Watch("k")
	w.Start(context.Background())
	time.Sleep(80 * time.Millisecond)
	w.Stop()

	if w.CycleCount() < 2 {
		t.Errorf("CycleCount = %d", w.CycleCount())
	}
}

func TestRotationWatcher_StartIdempotent(t *testing.T) {
	p := newFakeProvider()
	w, _ := secrets.NewRotationWatcher(p, secrets.RotationWatcherConfig{
		Interval: 50 * time.Millisecond,
	})
	w.Start(context.Background())
	w.Start(context.Background()) // 두 번째 호출 안전
	w.Stop()
}

func TestRotationWatcher_NilProvider(t *testing.T) {
	if _, err := secrets.NewRotationWatcher(nil, secrets.RotationWatcherConfig{}); err == nil {
		t.Error("nil provider 통과")
	}
}

func TestRotationWatcher_MultipleListeners(t *testing.T) {
	p := newFakeProvider()
	_ = p.Set(context.Background(), &secrets.Secret{Path: "k", Value: "v1", Version: 1})

	w, _ := secrets.NewRotationWatcher(p, secrets.RotationWatcherConfig{
		Interval: 50 * time.Millisecond,
	})
	w.Watch("k")

	var aFires, bFires int32
	w.AddListener(func(_ context.Context, _ secrets.RotationEvent) error {
		atomic.AddInt32(&aFires, 1)
		return nil
	})
	w.AddListener(func(_ context.Context, _ secrets.RotationEvent) error {
		atomic.AddInt32(&bFires, 1)
		return nil
	})

	w.CheckOnce(context.Background())
	_ = p.Set(context.Background(), &secrets.Secret{Path: "k", Value: "v2", Version: 2})
	w.CheckOnce(context.Background())

	if atomic.LoadInt32(&aFires) != 1 || atomic.LoadInt32(&bFires) != 1 {
		t.Errorf("aFires=%d bFires=%d", aFires, bFires)
	}
}

func TestRotationWatcher_ListenerError_PropagatesToOnError(t *testing.T) {
	p := newFakeProvider()
	_ = p.Set(context.Background(), &secrets.Secret{Path: "k", Value: "v1", Version: 1})

	var errs []error
	var mu sync.Mutex
	w, _ := secrets.NewRotationWatcher(p, secrets.RotationWatcherConfig{
		Interval: 50 * time.Millisecond,
		OnError: func(_ string, err error) {
			mu.Lock()
			defer mu.Unlock()
			errs = append(errs, err)
		},
	})
	w.Watch("k")
	w.AddListener(func(_ context.Context, _ secrets.RotationEvent) error {
		return errors.New("reload failed")
	})

	w.CheckOnce(context.Background())
	_ = p.Set(context.Background(), &secrets.Secret{Path: "k", Value: "v2", Version: 2})
	w.CheckOnce(context.Background())

	mu.Lock()
	defer mu.Unlock()
	if len(errs) == 0 {
		t.Error("리스너 에러가 OnError 로 전달되지 않음")
	}
}
