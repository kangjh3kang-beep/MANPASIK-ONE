// Package sla는 STAT 결과 표출 30초 SLA (P10-09b D-08) 모니터링입니다.
//
// 측정 시작 → 결과 표출까지의 latency를 추적하고 임계값 위반 시 알람을 발생시킵니다.
// Phase N tracing과 결합하여 분산 추적과 연동 가능합니다.
package sla

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/manpasik/backend/shared/medical/stats"
)

// ============================================================================
// SLA 정의
// ============================================================================

// Tier는 SLA 등급입니다.
type Tier string

const (
	TierSTAT       Tier = "stat"        // STAT 30초 (P10-09b D-08)
	TierUrgent     Tier = "urgent"      // 긴급 60초
	TierRoutine    Tier = "routine"     // 일반 5분
	TierBatch      Tier = "batch"       // 배치 1시간
)

// Target은 SLA 목표값입니다.
type Target struct {
	Tier            Tier
	BudgetMs        int64   // 응답 시간 예산 (ms)
	ToleranceMs     int64   // 허용 초과 (warning 임계, 100% = critical)
	BurnRateAlert   float64 // 1시간 윈도우 에러 예산 소진율 (예: 0.1 = 10%)
}

// DefaultTargets는 만파식 SLA 베이스라인입니다.
var DefaultTargets = map[Tier]Target{
	TierSTAT: {
		Tier: TierSTAT, BudgetMs: 30_000, ToleranceMs: 30_000, BurnRateAlert: 0.1,
	},
	TierUrgent: {
		Tier: TierUrgent, BudgetMs: 60_000, ToleranceMs: 60_000, BurnRateAlert: 0.1,
	},
	TierRoutine: {
		Tier: TierRoutine, BudgetMs: 300_000, ToleranceMs: 300_000, BurnRateAlert: 0.2,
	},
	TierBatch: {
		Tier: TierBatch, BudgetMs: 3_600_000, ToleranceMs: 3_600_000, BurnRateAlert: 0.3,
	},
}

// ============================================================================
// 측정 이벤트
// ============================================================================

// Measurement은 단일 STAT 측정 이벤트입니다.
type Measurement struct {
	ID              string
	SessionID       string
	Tier            Tier
	StartedAt       time.Time
	ResultAt        time.Time
	DurationMs      int64
	WithinSLA       bool
	Breach          bool
	BreachSeverity  string // "warning" | "critical"
	TestCode        string // LOINC 등
	Service         string // "ai-inference", "telemedicine" 등
}

// ComputeOutcome은 측정 결과를 계산합니다.
func ComputeOutcome(m *Measurement, target Target) {
	if m.ResultAt.IsZero() || m.StartedAt.IsZero() {
		m.WithinSLA = false
		return
	}
	m.DurationMs = m.ResultAt.Sub(m.StartedAt).Milliseconds()
	m.WithinSLA = m.DurationMs <= target.BudgetMs

	if !m.WithinSLA {
		m.Breach = true
		// budget 기준 50% 이상 초과면 critical, 그 미만은 warning
		if m.DurationMs >= target.BudgetMs+target.ToleranceMs/2 {
			m.BreachSeverity = "critical"
		} else {
			m.BreachSeverity = "warning"
		}
	}
}

// ============================================================================
// SLA 모니터
// ============================================================================

// Monitor는 SLA 측정을 수집·집계·알람합니다.
type Monitor struct {
	mu       sync.RWMutex
	targets  map[Tier]Target
	measurements []*Measurement
	maxHistory int
	alerts   []*Alert
}

// Alert는 SLA 위반 알람입니다.
type Alert struct {
	ID            string
	Tier          Tier
	Service       string
	Type          string // "breach", "burn_rate", "p99_spike"
	Severity      string
	Message       string
	BudgetMs      int64
	ActualMs      int64
	BurnRate      float64
	WindowStart   time.Time
	WindowEnd     time.Time
	TriggeredAt   time.Time
	Acknowledged  bool
}

// NewMonitor는 새 SLA 모니터를 생성합니다.
func NewMonitor(targets map[Tier]Target) *Monitor {
	if targets == nil {
		targets = DefaultTargets
	}
	return &Monitor{
		targets:    targets,
		maxHistory: 100_000,
	}
}

// Record는 측정 1건을 기록합니다.
func (m *Monitor) Record(measurement *Measurement) error {
	if measurement == nil {
		return errors.New("measurement is nil")
	}
	if measurement.SessionID == "" {
		return errors.New("session_id required")
	}
	if measurement.Tier == "" {
		measurement.Tier = TierSTAT
	}
	if measurement.ID == "" {
		measurement.ID = fmt.Sprintf("sla-%d", time.Now().UnixNano())
	}

	target, ok := m.targets[measurement.Tier]
	if !ok {
		return fmt.Errorf("unknown tier: %s", measurement.Tier)
	}

	ComputeOutcome(measurement, target)

	m.mu.Lock()
	defer m.mu.Unlock()
	m.measurements = append(m.measurements, measurement)
	if len(m.measurements) > m.maxHistory {
		m.measurements = m.measurements[len(m.measurements)-m.maxHistory:]
	}

	// breach 시 알람 자동 생성
	if measurement.Breach {
		m.alerts = append(m.alerts, &Alert{
			ID:           fmt.Sprintf("alert-%d", time.Now().UnixNano()),
			Tier:         measurement.Tier,
			Service:      measurement.Service,
			Type:         "breach",
			Severity:     measurement.BreachSeverity,
			Message:      fmt.Sprintf("SLA breach: %dms > %dms (tier %s)", measurement.DurationMs, target.BudgetMs, measurement.Tier),
			BudgetMs:     target.BudgetMs,
			ActualMs:     measurement.DurationMs,
			TriggeredAt:  time.Now().UTC(),
		})
	}
	return nil
}

// StartTimer는 측정 시작을 기록하고 종료 함수를 반환합니다.
//
// 사용 패턴:
//   stop := monitor.StartTimer(sessionID, sla.TierSTAT, "ai-inference", "2345-7")
//   ... 추론 수행 ...
//   stop() // 자동으로 ResultAt 기록 + Record 호출
func (m *Monitor) StartTimer(sessionID string, tier Tier, service, testCode string) func() *Measurement {
	measurement := &Measurement{
		SessionID: sessionID,
		Tier:      tier,
		Service:   service,
		TestCode:  testCode,
		StartedAt: time.Now().UTC(),
	}
	return func() *Measurement {
		measurement.ResultAt = time.Now().UTC()
		_ = m.Record(measurement)
		return measurement
	}
}

// ============================================================================
// 집계
// ============================================================================

// TierStats는 등급별 통계입니다.
type TierStats struct {
	Tier            Tier
	TotalSamples    int
	BreachCount     int
	BreachRate      float64
	MeanMs          float64
	P50Ms           float64
	P95Ms           float64
	P99Ms           float64
	BudgetMs        int64
	BudgetUtilization float64 // P95 / Budget
	WindowStart     time.Time
	WindowEnd       time.Time
}

// Stats는 등급의 윈도우 통계를 계산합니다.
func (m *Monitor) Stats(tier Tier, windowStart time.Time) *TierStats {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.statsLocked(tier, windowStart)
}

// statsLocked는 락이 이미 잡혀있는 상태에서 호출되는 내부 헬퍼입니다.
// EvaluateBurnRate처럼 단일 트랜잭션이 필요한 경우 사용.
func (m *Monitor) statsLocked(tier Tier, windowStart time.Time) *TierStats {
	target := m.targets[tier]
	end := time.Now().UTC()

	durations := []float64{}
	breachCount := 0
	for _, ms := range m.measurements {
		if ms.Tier != tier {
			continue
		}
		if !windowStart.IsZero() && ms.StartedAt.Before(windowStart) {
			continue
		}
		durations = append(durations, float64(ms.DurationMs))
		if ms.Breach {
			breachCount++
		}
	}

	tierStats := &TierStats{
		Tier:        tier,
		BudgetMs:    target.BudgetMs,
		WindowStart: windowStart,
		WindowEnd:   end,
		TotalSamples: len(durations),
		BreachCount:  breachCount,
	}
	if len(durations) == 0 {
		return tierStats
	}

	sort.Float64s(durations)
	tierStats.MeanMs = stats.Mean(durations)
	tierStats.P50Ms = stats.MedianSorted(durations)
	tierStats.P95Ms = stats.Percentile(durations, 95)
	tierStats.P99Ms = stats.Percentile(durations, 99)
	tierStats.BreachRate = float64(breachCount) / float64(len(durations))
	if target.BudgetMs > 0 {
		tierStats.BudgetUtilization = tierStats.P95Ms / float64(target.BudgetMs)
	}
	return tierStats
}

// EvaluateBurnRate는 에러 예산 소진율을 평가하고 임계값 초과 시 알람을 추가합니다.
//
// 호출 빈도: 5분 또는 1시간 단위 cron.
// Stats 평가와 알람 추가를 단일 락 트랜잭션으로 처리하여 race condition 방지.
func (m *Monitor) EvaluateBurnRate(tier Tier, windowStart time.Time) *Alert {
	m.mu.Lock()
	defer m.mu.Unlock()

	tierStats := m.statsLocked(tier, windowStart)
	target := m.targets[tier]

	if tierStats.TotalSamples == 0 {
		return nil
	}
	if tierStats.BreachRate < target.BurnRateAlert {
		return nil
	}

	alert := &Alert{
		ID:          fmt.Sprintf("burn-%d", time.Now().UnixNano()),
		Tier:        tier,
		Type:        "burn_rate",
		Severity:    "critical",
		Message:     fmt.Sprintf("burn rate %.2f%% > alert %.2f%% (n=%d)", tierStats.BreachRate*100, target.BurnRateAlert*100, tierStats.TotalSamples),
		BurnRate:    tierStats.BreachRate,
		WindowStart: tierStats.WindowStart,
		WindowEnd:   tierStats.WindowEnd,
		TriggeredAt: time.Now().UTC(),
	}
	m.alerts = append(m.alerts, alert)
	return alert
}

// AlertCount는 누적 알람 수를 반환합니다.
func (m *Monitor) AlertCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.alerts)
}

// AlertsBySeverity는 심각도별 알람을 반환합니다.
func (m *Monitor) AlertsBySeverity(severity string) []*Alert {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []*Alert
	for _, a := range m.alerts {
		if a.Severity == severity {
			out = append(out, a)
		}
	}
	return out
}

// AcknowledgeAlert는 알람을 확인 처리합니다.
func (m *Monitor) AcknowledgeAlert(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, a := range m.alerts {
		if a.ID == id {
			a.Acknowledged = true
			return nil
		}
	}
	return fmt.Errorf("alert %s not found", id)
}

// MeasurementCount는 누적 측정 수를 반환합니다.
func (m *Monitor) MeasurementCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.measurements)
}

// 통계 헬퍼는 backend/shared/medical/stats 패키지로 이동되었습니다.
