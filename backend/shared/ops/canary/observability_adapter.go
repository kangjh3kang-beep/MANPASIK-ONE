package canary

import (
	"context"
	"errors"

	"github.com/manpasik/backend/shared/observability"
)

// observabilitySnapshotter 는 observability.Metrics 의 Snapshot 메서드를 추상화.
//
// 이 인터페이스로 캡슐화하여 테스트 시 메모리 모킹이 쉽도록 함.
type observabilitySnapshotter interface {
	Snapshot() observability.MetricsSnapshot
}

// NewObservabilityCollector 는 observability.Metrics 를 MetricsCollector 로 어댑팅.
//
// canary FeedbackLoop 에 주입하면 매 사이클마다 observability 의 누적
// 통계 (ErrorRate / LatencyP99 / TotalRequests) 를 자동으로 카나리 게이팅에 사용.
//
// 사용 예 (서비스 main.go):
//
//	metrics := observability.NewMetrics()  // 기존 메트릭 인스턴스
//	collector := canary.NewObservabilityCollector(metrics)
//	loop, _ := canary.NewFeedbackLoop(mgr, rolloutID, collector, canary.FeedbackLoopConfig{
//	    Interval: 60*time.Second,
//	    AutoAdvance: true,
//	    AutoRollbackThresholds: &canary.AutoRollbackThresholds{ErrorRate: 0.05},
//	})
//	loop.Start(ctx)
func NewObservabilityCollector(m *observability.Metrics) MetricsCollector {
	return &observabilityCollector{snapshotter: m}
}

// newObservabilityCollectorFromIface 는 테스트용 — 인터페이스 직접 주입.
func newObservabilityCollectorFromIface(s observabilitySnapshotter) MetricsCollector {
	return &observabilityCollector{snapshotter: s}
}

type observabilityCollector struct {
	snapshotter observabilitySnapshotter
}

// Collect 는 observability.MetricsSnapshot → canary.MetricSnapshot 변환.
//
// LatencyP99 는 milliseconds 로 변환. RequestCount 는 누적값 (FeedbackLoop 가
// 사이클 간 차이를 보고 싶다면 별도 윈도우 어댑터 필요).
func (c *observabilityCollector) Collect(_ context.Context) (*MetricSnapshot, error) {
	if c.snapshotter == nil {
		return nil, errors.New("observability metrics 미설정")
	}
	snap := c.snapshotter.Snapshot()
	return &MetricSnapshot{
		ErrorRate:    snap.ErrorRate,
		LatencyP99Ms: snap.LatencyP99.Milliseconds(),
		RequestCount: snap.TotalRequests,
	}, nil
}

// WindowedObservabilityCollector 는 누적값에서 사이클 간 차이를 계산하여
// "최근 윈도우" 의 메트릭만 보고하는 어댑터.
//
// 누적 메트릭의 한계 (시간이 지날수록 ErrorRate 가 평탄화) 를 보완.
//
// 동작:
//   - 첫 호출: 현재 누적값을 baseline 으로 저장, RequestCount=0 반환
//   - 이후 호출: (현재 누적 - 직전 누적) 으로 윈도우 통계 계산
//   - LatencyP99 는 누적값 사용 (윈도우 단위 latency 기록은 observability 측 변경 필요)
type WindowedObservabilityCollector struct {
	snapshotter observabilitySnapshotter
	prevReqs    int64
	prevErrs    int64
	hasBaseline bool
}

// NewWindowedObservabilityCollector 생성.
func NewWindowedObservabilityCollector(m *observability.Metrics) *WindowedObservabilityCollector {
	return &WindowedObservabilityCollector{snapshotter: m}
}

// Collect 는 직전 호출 이후 발생한 요청/에러만 집계.
func (w *WindowedObservabilityCollector) Collect(_ context.Context) (*MetricSnapshot, error) {
	if w.snapshotter == nil {
		return nil, errors.New("observability metrics 미설정")
	}
	snap := w.snapshotter.Snapshot()

	if !w.hasBaseline {
		w.prevReqs = snap.TotalRequests
		w.prevErrs = snap.TotalErrors
		w.hasBaseline = true
		return &MetricSnapshot{
			LatencyP99Ms: snap.LatencyP99.Milliseconds(),
		}, nil
	}

	deltaReqs := snap.TotalRequests - w.prevReqs
	deltaErrs := snap.TotalErrors - w.prevErrs
	w.prevReqs = snap.TotalRequests
	w.prevErrs = snap.TotalErrors

	if deltaReqs < 0 {
		deltaReqs = 0
	}
	if deltaErrs < 0 {
		deltaErrs = 0
	}

	rate := float64(0)
	if deltaReqs > 0 {
		rate = float64(deltaErrs) / float64(deltaReqs)
	}
	return &MetricSnapshot{
		ErrorRate:    rate,
		LatencyP99Ms: snap.LatencyP99.Milliseconds(),
		RequestCount: deltaReqs,
	}, nil
}
