// Package canary는 카나리 배포 롤아웃 로직을 제공합니다.
//
// 단계별 트래픽 비율 (5% → 25% → 50% → 100%) + 메트릭 기반 자동 게이팅 + 롤백.
package canary

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

// ============================================================================
// 도메인 모델
// ============================================================================

// Phase는 롤아웃 단계입니다.
type Phase struct {
	Name             string
	TrafficPercent   int           // 0~100
	MinDuration      time.Duration // 최소 단계 유지 시간
	MaxErrorRate     float64       // 진행 가능 최대 에러율 (0~1)
	MaxLatencyP99Ms  int64         // 진행 가능 최대 P99 latency
	MinSuccessChecks int           // 진행 위해 필요한 성공 체크 수
}

// 표준 롤아웃 단계 (5% → 25% → 50% → 100%)
var DefaultPhases = []Phase{
	{Name: "phase-1", TrafficPercent: 5, MinDuration: 10 * time.Minute, MaxErrorRate: 0.01, MaxLatencyP99Ms: 1000, MinSuccessChecks: 3},
	{Name: "phase-2", TrafficPercent: 25, MinDuration: 30 * time.Minute, MaxErrorRate: 0.01, MaxLatencyP99Ms: 1000, MinSuccessChecks: 3},
	{Name: "phase-3", TrafficPercent: 50, MinDuration: 60 * time.Minute, MaxErrorRate: 0.005, MaxLatencyP99Ms: 800, MinSuccessChecks: 3},
	{Name: "phase-4", TrafficPercent: 100, MinDuration: 0, MaxErrorRate: 0.005, MaxLatencyP99Ms: 800, MinSuccessChecks: 0},
}

// MetricSnapshot은 모니터링 메트릭 스냅샷입니다.
type MetricSnapshot struct {
	ErrorRate      float64
	LatencyP99Ms   int64
	RequestCount   int64
	Timestamp      time.Time
}

// State는 롤아웃 상태입니다.
const (
	StateNotStarted = "not_started"
	StateInProgress = "in_progress"
	StateCompleted  = "completed"
	StatePaused     = "paused"
	StateRolledBack = "rolled_back"
	StateFailed     = "failed"
)

// Rollout은 카나리 롤아웃 인스턴스입니다.
type Rollout struct {
	ID              string
	ServiceName     string
	NewVersion      string
	OldVersion      string
	Phases          []Phase
	CurrentPhase    int
	State           string
	SuccessChecks   int
	StartedAt       time.Time
	PromotedAt      time.Time
	CompletedAt     *time.Time
	PauseReason     string
	RollbackReason  string
	MetricsHistory  []*MetricSnapshot
}

// CurrentTrafficPercent는 현재 트래픽 비율을 반환합니다.
func (r *Rollout) CurrentTrafficPercent() int {
	if r.CurrentPhase < 0 || r.CurrentPhase >= len(r.Phases) {
		return 0
	}
	return r.Phases[r.CurrentPhase].TrafficPercent
}

// ============================================================================
// 롤아웃 매니저
// ============================================================================

// Manager는 카나리 롤아웃을 관리합니다.
type Manager struct {
	mu       sync.RWMutex
	rollouts map[string]*Rollout
}

// NewManager는 새 매니저를 생성합니다.
func NewManager() *Manager {
	return &Manager{rollouts: make(map[string]*Rollout)}
}

// Start는 새 카나리 롤아웃을 시작합니다.
func (m *Manager) Start(rollout *Rollout) error {
	if rollout == nil || rollout.ID == "" {
		return errors.New("rollout id required")
	}
	if rollout.ServiceName == "" || rollout.NewVersion == "" {
		return errors.New("service_name and new_version required")
	}
	if len(rollout.Phases) == 0 {
		rollout.Phases = DefaultPhases
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.rollouts[rollout.ID]; exists {
		return fmt.Errorf("rollout %s already exists", rollout.ID)
	}

	rollout.State = StateInProgress
	rollout.CurrentPhase = 0
	rollout.SuccessChecks = 0
	rollout.StartedAt = time.Now().UTC()
	rollout.PromotedAt = rollout.StartedAt
	m.rollouts[rollout.ID] = rollout
	return nil
}

// RecordMetrics는 현재 단계의 메트릭을 기록합니다.
func (m *Manager) RecordMetrics(rolloutID string, metrics *MetricSnapshot) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	r, ok := m.rollouts[rolloutID]
	if !ok {
		return fmt.Errorf("rollout %s not found", rolloutID)
	}
	if r.State != StateInProgress {
		return fmt.Errorf("rollout state is %s, cannot record metrics", r.State)
	}

	if metrics.Timestamp.IsZero() {
		metrics.Timestamp = time.Now().UTC()
	}
	r.MetricsHistory = append(r.MetricsHistory, metrics)

	phase := r.Phases[r.CurrentPhase]
	// 메트릭이 임계값 이내면 success check 증가
	if metrics.ErrorRate <= phase.MaxErrorRate && metrics.LatencyP99Ms <= phase.MaxLatencyP99Ms {
		r.SuccessChecks++
	} else {
		// 임계값 초과 시 success check 리셋
		r.SuccessChecks = 0
	}
	return nil
}

// AdvanceIfReady는 다음 단계로 진행할지 결정합니다.
//
// 조건:
//   - 현재 단계 MinDuration 경과
//   - SuccessChecks >= MinSuccessChecks
//   - 마지막 메트릭이 임계값 이내
func (m *Manager) AdvanceIfReady(rolloutID string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	r, ok := m.rollouts[rolloutID]
	if !ok {
		return false, fmt.Errorf("rollout %s not found", rolloutID)
	}
	if r.State != StateInProgress {
		return false, nil
	}
	if r.CurrentPhase >= len(r.Phases)-1 {
		// 이미 마지막 단계 → 완료 처리
		now := time.Now().UTC()
		r.State = StateCompleted
		r.CompletedAt = &now
		return true, nil
	}

	phase := r.Phases[r.CurrentPhase]
	if time.Since(r.PromotedAt) < phase.MinDuration {
		return false, nil
	}
	if r.SuccessChecks < phase.MinSuccessChecks {
		return false, nil
	}

	// 다음 단계로 진행
	r.CurrentPhase++
	r.SuccessChecks = 0
	r.PromotedAt = time.Now().UTC()

	// 마지막 단계로 진입했고 MinDuration이 0이면 즉시 완료
	if r.CurrentPhase == len(r.Phases)-1 && r.Phases[r.CurrentPhase].MinDuration == 0 {
		now := time.Now().UTC()
		r.State = StateCompleted
		r.CompletedAt = &now
	}
	return true, nil
}

// Rollback은 롤아웃을 이전 버전으로 되돌립니다.
func (m *Manager) Rollback(rolloutID, reason string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	r, ok := m.rollouts[rolloutID]
	if !ok {
		return fmt.Errorf("rollout %s not found", rolloutID)
	}
	if r.State == StateCompleted {
		return fmt.Errorf("cannot rollback completed rollout")
	}
	now := time.Now().UTC()
	r.State = StateRolledBack
	r.RollbackReason = reason
	r.CompletedAt = &now
	return nil
}

// AutoRollback은 임계값 초과 시 자동 롤백합니다.
//
// 호출 빈도: 메트릭 수집 직후 (예: 1분 간격).
func (m *Manager) AutoRollback(rolloutID string, threshold *MetricSnapshot) error {
	m.mu.RLock()
	r, ok := m.rollouts[rolloutID]
	m.mu.RUnlock()
	if !ok {
		return fmt.Errorf("rollout %s not found", rolloutID)
	}
	if r.State != StateInProgress {
		return nil
	}
	if len(r.MetricsHistory) == 0 {
		return nil
	}

	last := r.MetricsHistory[len(r.MetricsHistory)-1]
	if last.ErrorRate > threshold.ErrorRate {
		return m.Rollback(rolloutID, fmt.Sprintf("error_rate %f > threshold %f", last.ErrorRate, threshold.ErrorRate))
	}
	if last.LatencyP99Ms > threshold.LatencyP99Ms {
		return m.Rollback(rolloutID, fmt.Sprintf("latency_p99 %dms > threshold %dms", last.LatencyP99Ms, threshold.LatencyP99Ms))
	}
	return nil
}

// Pause는 롤아웃을 일시 정지합니다.
func (m *Manager) Pause(rolloutID, reason string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	r, ok := m.rollouts[rolloutID]
	if !ok {
		return fmt.Errorf("rollout %s not found", rolloutID)
	}
	if r.State != StateInProgress {
		return fmt.Errorf("cannot pause state %s", r.State)
	}
	r.State = StatePaused
	r.PauseReason = reason
	return nil
}

// Resume은 일시 정지된 롤아웃을 재개합니다.
func (m *Manager) Resume(rolloutID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	r, ok := m.rollouts[rolloutID]
	if !ok {
		return fmt.Errorf("rollout %s not found", rolloutID)
	}
	if r.State != StatePaused {
		return fmt.Errorf("cannot resume state %s", r.State)
	}
	r.State = StateInProgress
	r.PauseReason = ""
	r.PromotedAt = time.Now().UTC() // 재개 시각으로 갱신
	return nil
}

// Get은 롤아웃을 조회합니다.
func (m *Manager) Get(rolloutID string) (*Rollout, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	r, ok := m.rollouts[rolloutID]
	if !ok {
		return nil, fmt.Errorf("rollout %s not found", rolloutID)
	}
	return r, nil
}

// ListByService는 서비스의 모든 롤아웃을 반환합니다 (최신순).
func (m *Manager) ListByService(serviceName string) []*Rollout {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []*Rollout
	for _, r := range m.rollouts {
		if r.ServiceName == serviceName {
			result = append(result, r)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].StartedAt.After(result[j].StartedAt)
	})
	return result
}

// CountByState는 상태별 롤아웃 수를 반환합니다.
func (m *Manager) CountByState() map[string]int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	counts := make(map[string]int)
	for _, r := range m.rollouts {
		counts[r.State]++
	}
	return counts
}
