package dashboard

import (
	"context"
	"time"
)

// 모듈별 상태 인터페이스 (각 ops/medical 모듈이 직접 구현하지 않아도 어댑터로 통합).
//
// 이 파일은 외부 모듈에 의존하지 않고 인터페이스만 정의한다. 실제 어댑터는
// 각 모듈 사용처(서비스 main.go) 에서 anonymous Provider 로 주입.

// PendingCounter 는 OTLP HTTP / 기타 비동기 모듈에서 사용.
type PendingCounter interface {
	PendingCount() int
}

// HealthChecker 는 단순 HealthCheck() 메서드를 가진 모듈.
type HealthChecker interface {
	HealthCheck(ctx context.Context) error
}

// AdapterFromHealthChecker 는 HealthChecker 를 HealthProvider 로 변환.
func AdapterFromHealthChecker(name string, hc HealthChecker) HealthProvider {
	return NewSimpleProvider(name, func(ctx context.Context) ModuleSnapshot {
		err := hc.HealthCheck(ctx)
		m := ModuleSnapshot{Name: name, UpdatedAt: time.Now()}
		if err != nil {
			m.Status = StatusFailed
			m.Message = err.Error()
			return m
		}
		m.Status = StatusHealthy
		return m
	})
}

// AdapterFromPendingCounter 는 PendingCount() 를 헬스 체크에 사용.
//
//   - PendingCount = 0  → healthy
//   - 0 < count <= warningThreshold → healthy + 메트릭 노출
//   - warningThreshold < count <= failedThreshold → degraded
//   - count > failedThreshold → failed
func AdapterFromPendingCounter(name string, pc PendingCounter, warningThreshold, failedThreshold int) HealthProvider {
	return NewSimpleProvider(name, func(ctx context.Context) ModuleSnapshot {
		c := pc.PendingCount()
		m := ModuleSnapshot{
			Name:      name,
			UpdatedAt: time.Now(),
			Metrics:   map[string]interface{}{"pending": c},
		}
		switch {
		case c <= warningThreshold:
			m.Status = StatusHealthy
		case c <= failedThreshold:
			m.Status = StatusDegraded
			m.Message = "pending 누적"
		default:
			m.Status = StatusFailed
			m.Message = "pending 임계 초과"
		}
		return m
	})
}

// MetricsCollector 는 임의의 메트릭 키-값을 노출.
type MetricsCollector interface {
	CollectMetrics() map[string]interface{}
}

// AdapterFromMetricsCollector 는 메트릭만 수집하고 상태를 healthy 로 보고.
//
// 메트릭에 "status" 키가 있으면 그 값을 Status 로 사용.
func AdapterFromMetricsCollector(name string, mc MetricsCollector) HealthProvider {
	return NewSimpleProvider(name, func(ctx context.Context) ModuleSnapshot {
		metrics := mc.CollectMetrics()
		status := StatusHealthy
		message := ""
		if s, ok := metrics["status"].(string); ok {
			status = Status(s)
		}
		if msg, ok := metrics["message"].(string); ok {
			message = msg
		}
		return ModuleSnapshot{
			Name:      name,
			Status:    status,
			Message:   message,
			Metrics:   metrics,
			UpdatedAt: time.Now(),
		}
	})
}
