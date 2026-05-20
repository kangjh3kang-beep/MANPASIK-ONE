package canary

import (
	"context"
	"errors"
	"sync"
	"time"
)

// MetricsCollector 는 서비스의 운영 메트릭을 카나리 매니저가 사용할 수 있는
// MetricSnapshot 으로 추상화.
//
// 서비스마다 Prometheus / OpenTelemetry / 자체 카운터 등 다양하므로,
// 어댑터를 주입하는 패턴을 사용.
type MetricsCollector interface {
	// Collect 는 마지막 윈도우(예: 60초) 동안의 메트릭 스냅샷 반환.
	Collect(ctx context.Context) (*MetricSnapshot, error)
}

// MetricsCollectorFunc 는 함수를 MetricsCollector 로 감쌉니다.
type MetricsCollectorFunc func(ctx context.Context) (*MetricSnapshot, error)

// Collect 는 fn 호출.
func (f MetricsCollectorFunc) Collect(ctx context.Context) (*MetricSnapshot, error) {
	return f(ctx)
}

// FeedbackLoop 는 주기적으로 collector 에서 메트릭을 받아 manager.RecordMetrics
// 를 호출하고, 단계 진행/롤백 결정을 내림.
//
// 사용 예 (서비스 main.go):
//
//	mgr := canary.NewManager()
//	_ = mgr.Start(rollout)
//
//	collector := canary.MetricsCollectorFunc(func(ctx context.Context) (*canary.MetricSnapshot, error) {
//	    return &canary.MetricSnapshot{
//	        ErrorRate:    metrics.ErrorRate(),
//	        LatencyP99Ms: metrics.LatencyP99Ms(),
//	        RequestCount: metrics.RequestCount(),
//	    }, nil
//	})
//
//	loop := canary.NewFeedbackLoop(mgr, rollout.ID, collector, canary.FeedbackLoopConfig{
//	    Interval: 60 * time.Second,
//	    AutoAdvance: true,
//	    AutoRollbackThresholds: &canary.AutoRollbackThresholds{ErrorRate: 0.05, LatencyP99Ms: 2000},
//	})
//	loop.Start(ctx)
//	defer loop.Stop()
type FeedbackLoop struct {
	mgr        *Manager
	rolloutID  string
	collector  MetricsCollector
	cfg        FeedbackLoopConfig
	mu         sync.Mutex
	stopCh     chan struct{}
	doneCh     chan struct{}
	cycleCount int
}

// FeedbackLoopConfig 는 루프 동작 설정.
type FeedbackLoopConfig struct {
	// Interval 는 메트릭 수집 주기 (기본 60초).
	Interval time.Duration
	// AutoAdvance=true 면 매 사이클마다 AdvanceIfReady 호출.
	AutoAdvance bool
	// AutoRollbackThresholds 가 설정되면 매 사이클에 임계값 초과 시 자동 롤백.
	// nil 이면 자동 롤백 안 함.
	AutoRollbackThresholds *AutoRollbackThresholds
	// OnError 는 루프 내 오류 발생 시 콜백 (로그/알림용). nil 가능.
	OnError func(err error)
}

// AutoRollbackThresholds 는 자동 롤백 임계값.
//
// 메트릭이 이 임계값을 초과하면 즉시 manager.Rollback 호출.
type AutoRollbackThresholds struct {
	ErrorRate    float64
	LatencyP99Ms int64
}

// NewFeedbackLoop 생성. 필수 필드(mgr, rolloutID, collector) 누락 시 nil + nil 반환 회피.
func NewFeedbackLoop(mgr *Manager, rolloutID string, collector MetricsCollector, cfg FeedbackLoopConfig) (*FeedbackLoop, error) {
	if mgr == nil {
		return nil, errors.New("manager 필수")
	}
	if rolloutID == "" {
		return nil, errors.New("rolloutID 필수")
	}
	if collector == nil {
		return nil, errors.New("collector 필수")
	}
	if cfg.Interval <= 0 {
		cfg.Interval = 60 * time.Second
	}
	return &FeedbackLoop{
		mgr:       mgr,
		rolloutID: rolloutID,
		collector: collector,
		cfg:       cfg,
	}, nil
}

// Start 는 백그라운드 goroutine 으로 루프 시작.
//
// 이미 실행 중이면 no-op. ctx 또는 Stop() 호출 시 종료.
func (l *FeedbackLoop) Start(ctx context.Context) {
	l.mu.Lock()
	if l.stopCh != nil {
		l.mu.Unlock()
		return
	}
	stopCh := make(chan struct{})
	doneCh := make(chan struct{})
	l.stopCh = stopCh
	l.doneCh = doneCh
	l.mu.Unlock()

	// stopCh/doneCh 는 goroutine 실행 동안 안정적으로 유지하기 위해 클로저로 캡처
	go l.run(ctx, stopCh, doneCh)
}

// Stop 은 루프 종료. 현재 진행 중인 사이클 완료까지 대기.
func (l *FeedbackLoop) Stop() {
	l.mu.Lock()
	if l.stopCh == nil {
		l.mu.Unlock()
		return
	}
	close(l.stopCh)
	doneCh := l.doneCh
	l.stopCh = nil
	l.mu.Unlock()
	if doneCh != nil {
		<-doneCh
	}
}

// CycleCount 는 지금까지 실행된 사이클 수 반환.
func (l *FeedbackLoop) CycleCount() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.cycleCount
}

func (l *FeedbackLoop) run(parentCtx context.Context, stopCh, doneCh chan struct{}) {
	defer close(doneCh)
	ticker := time.NewTicker(l.cfg.Interval)
	defer ticker.Stop()

	// 첫 사이클은 즉시 실행 (운영 신호 빠르게 확보)
	l.cycle(parentCtx)
	for {
		select {
		case <-stopCh:
			return
		case <-parentCtx.Done():
			return
		case <-ticker.C:
			l.cycle(parentCtx)
		}
	}
}

func (l *FeedbackLoop) cycle(ctx context.Context) {
	l.mu.Lock()
	l.cycleCount++
	l.mu.Unlock()

	cycleCtx, cancel := context.WithTimeout(ctx, l.cfg.Interval)
	defer cancel()

	snap, err := l.collector.Collect(cycleCtx)
	if err != nil {
		l.reportErr(err)
		return
	}
	if snap == nil {
		return
	}

	if err := l.mgr.RecordMetrics(l.rolloutID, snap); err != nil {
		l.reportErr(err)
		return
	}

	// 자동 롤백 검사
	if th := l.cfg.AutoRollbackThresholds; th != nil {
		if shouldAutoRollback(snap, th) {
			_ = l.mgr.Rollback(l.rolloutID,
				"auto rollback: error_rate or latency_p99 exceeded thresholds")
			return
		}
	}

	// 자동 진행
	if l.cfg.AutoAdvance {
		if _, err := l.mgr.AdvanceIfReady(l.rolloutID); err != nil {
			l.reportErr(err)
		}
	}
}

func shouldAutoRollback(snap *MetricSnapshot, th *AutoRollbackThresholds) bool {
	if snap == nil || th == nil {
		return false
	}
	if th.ErrorRate > 0 && snap.ErrorRate > th.ErrorRate {
		return true
	}
	if th.LatencyP99Ms > 0 && snap.LatencyP99Ms > th.LatencyP99Ms {
		return true
	}
	return false
}

func (l *FeedbackLoop) reportErr(err error) {
	if l.cfg.OnError != nil {
		l.cfg.OnError(err)
	}
}
