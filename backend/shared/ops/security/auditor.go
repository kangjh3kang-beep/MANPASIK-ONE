// Package security는 보안 감사 + 취약점 스캔 + Dependabot 통합을 추상화합니다.
//
// 지원 도구: CodeQL (GitHub Action), OWASP ZAP DAST, Trivy 컨테이너 스캔, Dependabot.
// 본 패키지는 결과 집계와 정책 평가에 집중합니다.
package security

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

// Severity는 취약점 심각도입니다.
type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
	SeverityLow      Severity = "low"
	SeverityInfo     Severity = "info"
)

// SeverityScore는 심각도를 0~10 점수로 반환합니다 (CVSS 기준).
func SeverityScore(s Severity) float64 {
	switch s {
	case SeverityCritical:
		return 9.5
	case SeverityHigh:
		return 7.5
	case SeverityMedium:
		return 5.0
	case SeverityLow:
		return 3.0
	case SeverityInfo:
		return 1.0
	default:
		return 0
	}
}

// ScanType은 스캔 종류입니다.
type ScanType string

const (
	ScanCodeQL     ScanType = "codeql"      // SAST
	ScanZAP        ScanType = "zap"         // DAST
	ScanTrivy      ScanType = "trivy"       // 컨테이너
	ScanDependabot ScanType = "dependabot"  // 의존성
	ScanSecret     ScanType = "secret_scan" // 시크릿 누출
)

// Finding은 단일 취약점/이슈입니다.
type Finding struct {
	ID          string
	Type        ScanType
	Severity    Severity
	Title       string
	Description string
	CVE         string // CVE-2024-XXXXX
	CWE         string // CWE-79
	Component   string // 패키지명, 파일경로
	Version     string
	FixedIn     string // 수정 가능한 버전
	Reference   string // URL
	DiscoveredAt time.Time
	Acknowledged bool
}

// ScanReport는 단일 스캔 보고서입니다.
type ScanReport struct {
	ID          string
	Type        ScanType
	Target      string
	Findings    []*Finding
	StartedAt   time.Time
	CompletedAt time.Time
	Duration    time.Duration
	Summary     map[Severity]int
}

// CountBySeverity는 보고서 내 심각도별 건수를 계산합니다.
func (r *ScanReport) CountBySeverity() map[Severity]int {
	counts := make(map[Severity]int)
	for _, f := range r.Findings {
		counts[f.Severity]++
	}
	return counts
}

// ============================================================================
// 보안 정책
// ============================================================================

// Policy는 보안 정책입니다.
type Policy struct {
	Name                string
	BlockOnCritical     bool          // critical 발견 시 배포 차단
	BlockOnHigh         bool
	MaxMediumPerService int           // service당 medium 허용 수
	MaxAge              time.Duration // 발견 후 N일 이내 수정
	ExemptCVEs          []string      // 예외 CVE 목록
}

// DefaultPolicy는 운영 기본 정책입니다.
var DefaultPolicy = Policy{
	Name:                "default",
	BlockOnCritical:     true,
	BlockOnHigh:         true,
	MaxMediumPerService: 5,
	MaxAge:              30 * 24 * time.Hour,
}

// PolicyDecision은 정책 평가 결과입니다.
type PolicyDecision struct {
	Allowed       bool
	Reason        string
	BlockedBy     []*Finding
	WarningsCount int
}

// PolicyEvaluator는 보고서를 정책에 따라 평가합니다.
type PolicyEvaluator struct {
	policy Policy
}

// NewPolicyEvaluator는 새 평가자를 생성합니다.
func NewPolicyEvaluator(policy Policy) *PolicyEvaluator {
	return &PolicyEvaluator{policy: policy}
}

// Evaluate는 보고서를 평가하여 배포 가능 여부를 반환합니다.
func (e *PolicyEvaluator) Evaluate(report *ScanReport) *PolicyDecision {
	if report == nil {
		return &PolicyDecision{Allowed: false, Reason: "report is nil"}
	}

	exemptSet := make(map[string]bool)
	for _, cve := range e.policy.ExemptCVEs {
		exemptSet[cve] = true
	}

	decision := &PolicyDecision{Allowed: true}
	mediumCount := 0
	now := time.Now().UTC()

	for _, f := range report.Findings {
		if f.Acknowledged {
			continue
		}
		if exemptSet[f.CVE] {
			continue
		}

		// 심각도별 차단
		if f.Severity == SeverityCritical && e.policy.BlockOnCritical {
			decision.Allowed = false
			decision.BlockedBy = append(decision.BlockedBy, f)
		}
		if f.Severity == SeverityHigh && e.policy.BlockOnHigh {
			decision.Allowed = false
			decision.BlockedBy = append(decision.BlockedBy, f)
		}
		if f.Severity == SeverityMedium {
			mediumCount++
		}

		// 나이 정책
		if e.policy.MaxAge > 0 && now.Sub(f.DiscoveredAt) > e.policy.MaxAge {
			decision.Allowed = false
			decision.BlockedBy = append(decision.BlockedBy, f)
		}
	}

	if mediumCount > e.policy.MaxMediumPerService {
		decision.Allowed = false
		decision.Reason = fmt.Sprintf("medium count %d exceeds max %d", mediumCount, e.policy.MaxMediumPerService)
	}
	decision.WarningsCount = mediumCount

	if decision.Allowed && len(decision.BlockedBy) > 0 {
		decision.Allowed = false
	}
	if !decision.Allowed && decision.Reason == "" {
		decision.Reason = fmt.Sprintf("blocked by %d findings", len(decision.BlockedBy))
	}
	return decision
}

// ============================================================================
// 보고서 집계
// ============================================================================

// Aggregator는 여러 보고서를 통합합니다.
type Aggregator struct {
	mu      sync.RWMutex
	reports map[string]*ScanReport
}

// NewAggregator는 새 집계자를 생성합니다.
func NewAggregator() *Aggregator {
	return &Aggregator{reports: make(map[string]*ScanReport)}
}

// Add는 새 보고서를 추가합니다.
func (a *Aggregator) Add(report *ScanReport) error {
	if report == nil || report.ID == "" {
		return errors.New("report id required")
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.reports[report.ID] = report
	return nil
}

// Get은 보고서를 조회합니다.
func (a *Aggregator) Get(id string) (*ScanReport, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	r, ok := a.reports[id]
	if !ok {
		return nil, fmt.Errorf("report %s not found", id)
	}
	return r, nil
}

// AllFindings는 모든 보고서의 finding을 통합 반환합니다 (심각도 내림차순).
func (a *Aggregator) AllFindings() []*Finding {
	a.mu.RLock()
	defer a.mu.RUnlock()
	var all []*Finding
	for _, r := range a.reports {
		all = append(all, r.Findings...)
	}
	sort.Slice(all, func(i, j int) bool {
		return SeverityScore(all[i].Severity) > SeverityScore(all[j].Severity)
	})
	return all
}

// FindingsByType는 스캔 종류별 finding을 반환합니다.
func (a *Aggregator) FindingsByType(scanType ScanType) []*Finding {
	a.mu.RLock()
	defer a.mu.RUnlock()
	var result []*Finding
	for _, r := range a.reports {
		if r.Type != scanType {
			continue
		}
		result = append(result, r.Findings...)
	}
	return result
}

// CriticalCount는 critical 심각도 finding 수를 반환합니다.
func (a *Aggregator) CriticalCount() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	count := 0
	for _, r := range a.reports {
		for _, f := range r.Findings {
			if f.Severity == SeverityCritical && !f.Acknowledged {
				count++
			}
		}
	}
	return count
}

// AcknowledgeFinding은 finding을 ack 처리합니다.
func (a *Aggregator) AcknowledgeFinding(reportID, findingID string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	r, ok := a.reports[reportID]
	if !ok {
		return fmt.Errorf("report %s not found", reportID)
	}
	for _, f := range r.Findings {
		if f.ID == findingID {
			f.Acknowledged = true
			return nil
		}
	}
	return fmt.Errorf("finding %s not found", findingID)
}

// Count는 저장된 보고서 수를 반환합니다.
func (a *Aggregator) Count() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.reports)
}

// SummaryByScanType은 스캔 종류별 통계를 반환합니다.
func (a *Aggregator) SummaryByScanType() map[ScanType]int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	counts := make(map[ScanType]int)
	for _, r := range a.reports {
		counts[r.Type] += len(r.Findings)
	}
	return counts
}
