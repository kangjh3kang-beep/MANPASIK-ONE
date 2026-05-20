// Package poise는 PoISE Loop (Plan-Observe-Improve-Standardize-Evaluate) 자가 개선 루프입니다.
//
// P10-09b D-10 차별화의 핵심.
// 화상진료 세션의 만족도/재내원/정확도 메트릭을 수집하여 모델 개선 후보를
// 인간 승인 게이트를 거쳐 적용합니다.
package poise

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/manpasik/backend/shared/medical/stats"
)

// ============================================================================
// 도메인 모델
// ============================================================================

// FeedbackKind는 피드백 종류입니다.
type FeedbackKind string

const (
	FeedbackSatisfaction FeedbackKind = "satisfaction"  // 만족도 (1~10)
	FeedbackReturn       FeedbackKind = "return"        // 재내원 여부 (true/false)
	FeedbackAccuracy     FeedbackKind = "accuracy"      // 정확도 (의사 검토 후)
	FeedbackSafetyAlert  FeedbackKind = "safety_alert"  // 안전 위반
	FeedbackResponseTime FeedbackKind = "response_time" // 응답 시간 (ms)
)

// Feedback은 단일 피드백 입력입니다.
type Feedback struct {
	ID          string
	SessionID   string
	UserID      string
	Kind        FeedbackKind
	Value       float64 // 만족도 1-10, 응답시간 ms, 정확도 0-1, 재내원 0/1
	Comment     string
	CollectedAt time.Time
}

// Metric은 집계된 KPI 메트릭입니다.
type Metric struct {
	Name        string
	Kind        FeedbackKind
	Mean        float64
	Median      float64
	StdDev      float64
	P95         float64
	P99         float64
	SampleSize  int
	WindowStart time.Time
	WindowEnd   time.Time
}

// ImprovementProposal는 모델/시스템 개선 제안입니다.
type ImprovementProposal struct {
	ID            string
	Trigger       string  // "satisfaction_drop", "accuracy_drop", "safety_violation_spike", "latency_spike"
	TargetModule  string  // "digital_physician", "scribe", "polaris", "reasoner"
	Severity      int     // 1-10
	Reason        string
	Recommendation string
	BaselineMetric *Metric
	CurrentMetric  *Metric
	ProposedAt    time.Time
	Status        string // "pending", "approved", "rejected", "applied"
	ApprovedBy    string
	ApprovedAt    *time.Time
	AppliedAt     *time.Time
	RejectReason  string
}

// ============================================================================
// 임계값 정책 (조정 가능)
// ============================================================================

// Thresholds는 개선 제안 트리거 임계값입니다.
type Thresholds struct {
	SatisfactionMinAvg     float64 // 평균 만족도 미달 시 트리거 (기본 8.0/10)
	AccuracyMinAvg         float64 // 정확도 미달 (기본 0.92)
	SafetyViolationsPerK   float64 // 1000세션당 안전 위반 (기본 1.0)
	ResponseTimeP95Ms      float64 // P95 응답시간 (기본 1500ms)
	ReturnRateMin          float64 // 재내원율 최소 (기본 0.30)
	MinSampleSize          int     // 분석 최소 샘플 (기본 30)
}

// DefaultThresholds는 기본 정책입니다.
var DefaultThresholds = Thresholds{
	SatisfactionMinAvg:   8.0,
	AccuracyMinAvg:       0.92,
	SafetyViolationsPerK: 1.0,
	ResponseTimeP95Ms:    1500,
	ReturnRateMin:        0.30,
	MinSampleSize:        30,
}

// ============================================================================
// PoISE Loop 엔진
// ============================================================================

// Loop는 PoISE 자가 개선 엔진입니다.
type Loop struct {
	mu          sync.RWMutex
	feedback    []*Feedback
	proposals   map[string]*ImprovementProposal
	thresholds  Thresholds
	maxFeedbackHistory int
}

// NewLoop는 새 PoISE Loop를 생성합니다.
func NewLoop(thresholds Thresholds) *Loop {
	return &Loop{
		proposals:          make(map[string]*ImprovementProposal),
		thresholds:         thresholds,
		maxFeedbackHistory: 100_000,
	}
}

// CollectFeedback은 피드백을 추가합니다.
func (l *Loop) CollectFeedback(fb *Feedback) error {
	if fb == nil {
		return errors.New("feedback is nil")
	}
	if fb.SessionID == "" {
		return errors.New("session_id required")
	}
	if fb.Kind == "" {
		return errors.New("kind required")
	}
	if fb.ID == "" {
		fb.ID = fmt.Sprintf("fb-%d", time.Now().UnixNano())
	}
	if fb.CollectedAt.IsZero() {
		fb.CollectedAt = time.Now().UTC()
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	l.feedback = append(l.feedback, fb)

	// 히스토리 윈도우
	if len(l.feedback) > l.maxFeedbackHistory {
		l.feedback = l.feedback[len(l.feedback)-l.maxFeedbackHistory:]
	}
	return nil
}

// AggregateMetric은 특정 종류의 피드백을 집계합니다.
func (l *Loop) AggregateMetric(kind FeedbackKind, since time.Time) *Metric {
	l.mu.RLock()
	defer l.mu.RUnlock()

	values := make([]float64, 0)
	end := time.Now().UTC()
	for _, fb := range l.feedback {
		if fb.Kind != kind {
			continue
		}
		if !since.IsZero() && fb.CollectedAt.Before(since) {
			continue
		}
		values = append(values, fb.Value)
	}

	if len(values) == 0 {
		return &Metric{Name: string(kind), Kind: kind, WindowStart: since, WindowEnd: end}
	}

	sort.Float64s(values)
	mean := stats.Mean(values)
	std := stats.StdDev(values, mean)

	return &Metric{
		Name:        string(kind),
		Kind:        kind,
		Mean:        mean,
		Median:      stats.MedianSorted(values),
		StdDev:      std,
		P95:         stats.Percentile(values, 95),
		P99:         stats.Percentile(values, 99),
		SampleSize:  len(values),
		WindowStart: since,
		WindowEnd:   end,
	}
}

// ============================================================================
// Plan: 개선 제안 트리거 평가
// ============================================================================

// EvaluateAndPropose는 임계값 위반 시 개선 제안을 생성합니다.
//
// 호출 빈도: 주 1회 (또는 N=1000 세션마다).
func (l *Loop) EvaluateAndPropose(since time.Time) []*ImprovementProposal {
	var proposals []*ImprovementProposal

	// 1. 만족도 검사
	if metric := l.AggregateMetric(FeedbackSatisfaction, since); metric.SampleSize >= l.thresholds.MinSampleSize {
		if metric.Mean < l.thresholds.SatisfactionMinAvg {
			proposals = append(proposals, l.proposeSatisfactionImprovement(metric))
		}
	}

	// 2. 정확도 검사
	if metric := l.AggregateMetric(FeedbackAccuracy, since); metric.SampleSize >= l.thresholds.MinSampleSize {
		if metric.Mean < l.thresholds.AccuracyMinAvg {
			proposals = append(proposals, l.proposeAccuracyImprovement(metric))
		}
	}

	// 3. 응답 시간 검사
	if metric := l.AggregateMetric(FeedbackResponseTime, since); metric.SampleSize >= l.thresholds.MinSampleSize {
		if metric.P95 > l.thresholds.ResponseTimeP95Ms {
			proposals = append(proposals, l.proposeLatencyImprovement(metric))
		}
	}

	// 4. 안전 위반 검사
	//
	// FeedbackSafetyAlert.Value는 다음 의미를 가집니다:
	//   - 0: 안전한 세션 (위반 없음)
	//   - 1: 경미한 위반 (rewrite/warn 처리됨)
	//   - 2: 심각한 위반 (block/escalate 처리됨)
	// 트리거 조건: (위반 합계 / 총 세션) × 1000 ≥ 임계값
	if metric := l.AggregateMetric(FeedbackSafetyAlert, since); metric.SampleSize > 0 {
		// 합계 = mean × sample_size
		violationCount := metric.Mean * float64(metric.SampleSize)
		violationsPerK := (violationCount / float64(metric.SampleSize)) * 1000
		// 총 세션 카운트도 같이 비교 (저샘플 false positive 방지)
		if metric.SampleSize >= l.thresholds.MinSampleSize && violationsPerK > l.thresholds.SafetyViolationsPerK {
			proposals = append(proposals, l.proposeSafetyImprovement(metric))
		}
	}

	// 등록
	l.mu.Lock()
	for _, p := range proposals {
		l.proposals[p.ID] = p
	}
	l.mu.Unlock()

	return proposals
}

func (l *Loop) proposeSatisfactionImprovement(metric *Metric) *ImprovementProposal {
	return &ImprovementProposal{
		ID:           fmt.Sprintf("imp-sat-%d", time.Now().UnixNano()),
		Trigger:      "satisfaction_drop",
		TargetModule: "digital_physician",
		Severity:     6,
		Reason: fmt.Sprintf("평균 만족도 %.2f가 임계값 %.2f 미달 (n=%d)",
			metric.Mean, l.thresholds.SatisfactionMinAvg, metric.SampleSize),
		Recommendation: "디지털 주치의 발화 톤·공감 표현 강화 (Polaris 안전망 + 추가 따뜻함 emotion 비율 ↑)",
		CurrentMetric:  metric,
		ProposedAt:     time.Now().UTC(),
		Status:         "pending",
	}
}

func (l *Loop) proposeAccuracyImprovement(metric *Metric) *ImprovementProposal {
	return &ImprovementProposal{
		ID:           fmt.Sprintf("imp-acc-%d", time.Now().UnixNano()),
		Trigger:      "accuracy_drop",
		TargetModule: "reasoner",
		Severity:     8,
		Reason: fmt.Sprintf("진단 정확도 %.3f가 임계값 %.3f 미달 (n=%d)",
			metric.Mean, l.thresholds.AccuracyMinAvg, metric.SampleSize),
		Recommendation: "감별진단 후보 가중치 재학습 + SHAP 검증 강화 (개별 모델 재학습 후보)",
		CurrentMetric:  metric,
		ProposedAt:     time.Now().UTC(),
		Status:         "pending",
	}
}

func (l *Loop) proposeLatencyImprovement(metric *Metric) *ImprovementProposal {
	return &ImprovementProposal{
		ID:           fmt.Sprintf("imp-lat-%d", time.Now().UnixNano()),
		Trigger:      "latency_spike",
		TargetModule: "all",
		Severity:     5,
		Reason: fmt.Sprintf("P95 응답 %.0fms가 임계값 %.0fms 초과 (n=%d)",
			metric.P95, l.thresholds.ResponseTimeP95Ms, metric.SampleSize),
		Recommendation: "추론 캐시 워밍 + DB 인덱스 검토 + canary 트래픽 50% 감속",
		CurrentMetric:  metric,
		ProposedAt:     time.Now().UTC(),
		Status:         "pending",
	}
}

func (l *Loop) proposeSafetyImprovement(metric *Metric) *ImprovementProposal {
	return &ImprovementProposal{
		ID:           fmt.Sprintf("imp-safe-%d", time.Now().UnixNano()),
		Trigger:      "safety_violation_spike",
		TargetModule: "polaris",
		Severity:     9,
		Reason: fmt.Sprintf("안전 위반 비율 %.2f/1k 세션 (임계 %.2f, n=%d)",
			metric.Mean, l.thresholds.SafetyViolationsPerK, metric.SampleSize),
		Recommendation: "Polaris 룰 추가 + 의사 인계 임계 강화 + 감사 로그 검토",
		CurrentMetric:  metric,
		ProposedAt:     time.Now().UTC(),
		Status:         "pending",
	}
}

// ============================================================================
// Improve: 인간 승인 게이트
// ============================================================================

// Approve는 제안을 승인합니다 (인간 게이트 — 자동 적용 금지).
func (l *Loop) Approve(proposalID, approverID string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	p, ok := l.proposals[proposalID]
	if !ok {
		return fmt.Errorf("proposal %s not found", proposalID)
	}
	if p.Status != "pending" {
		return fmt.Errorf("proposal status is %s", p.Status)
	}
	now := time.Now().UTC()
	p.Status = "approved"
	p.ApprovedBy = approverID
	p.ApprovedAt = &now
	return nil
}

// Reject는 제안을 거부합니다.
func (l *Loop) Reject(proposalID, approverID, reason string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	p, ok := l.proposals[proposalID]
	if !ok {
		return fmt.Errorf("proposal %s not found", proposalID)
	}
	if p.Status != "pending" {
		return fmt.Errorf("proposal status is %s", p.Status)
	}
	now := time.Now().UTC()
	p.Status = "rejected"
	p.ApprovedBy = approverID
	p.ApprovedAt = &now
	p.RejectReason = reason
	return nil
}

// MarkApplied는 승인된 제안의 적용 완료를 기록합니다 (배포 후).
func (l *Loop) MarkApplied(proposalID string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	p, ok := l.proposals[proposalID]
	if !ok {
		return fmt.Errorf("proposal %s not found", proposalID)
	}
	if p.Status != "approved" {
		return fmt.Errorf("only approved proposals can be marked applied (current: %s)", p.Status)
	}
	now := time.Now().UTC()
	p.Status = "applied"
	p.AppliedAt = &now
	return nil
}

// GetProposal은 제안을 조회합니다.
func (l *Loop) GetProposal(id string) (*ImprovementProposal, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	p, ok := l.proposals[id]
	if !ok {
		return nil, fmt.Errorf("proposal %s not found", id)
	}
	return p, nil
}

// ListProposals는 상태로 필터링한 제안 목록을 반환합니다.
func (l *Loop) ListProposals(status string) []*ImprovementProposal {
	l.mu.RLock()
	defer l.mu.RUnlock()
	var out []*ImprovementProposal
	for _, p := range l.proposals {
		if status == "" || p.Status == status {
			out = append(out, p)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].ProposedAt.After(out[j].ProposedAt)
	})
	return out
}

// CountByStatus는 상태별 제안 수를 반환합니다.
func (l *Loop) CountByStatus() map[string]int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	counts := make(map[string]int)
	for _, p := range l.proposals {
		counts[p.Status]++
	}
	return counts
}

// FeedbackCount는 누적 피드백 수를 반환합니다.
func (l *Loop) FeedbackCount() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.feedback)
}

// 통계 헬퍼는 backend/shared/medical/stats 패키지로 이동되었습니다.
// 본 파일에서는 stats.Mean / stats.StdDev / stats.Percentile / stats.MedianSorted 사용.
