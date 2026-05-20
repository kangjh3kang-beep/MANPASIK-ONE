// Package service: HIPAA 7년 보존 정책 + LOINC 기반 검색 + immutable 보장
//
// audit-service CRUD → CRUD+ 승격 모듈.
// HIPAA §164.530(j): 의료 기록은 작성일로부터 6년 보존, 만파식은 7년 안전 마진.
package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

// ============================================================================
// 보존 정책 상수 (HIPAA + 한국 의료법)
// ============================================================================

const (
	HIPAARetentionYears   = 6 // §164.530(j) 최소 6년
	ManpasikRetentionYears = 7 // 안전 마진 1년
	KoreaMedicalRecordYears = 10 // 한국 의료법 시행령 §15 (10년)
	DefaultRetentionYears   = 10 // 가장 보수적인 정책 채택
)

// RetentionPolicy는 데이터 보존 정책입니다.
type RetentionPolicy struct {
	Name             string
	Years            int
	Region           string // "us", "kr", "global"
	ImmutableDays    int    // 작성 후 N일 동안 수정/삭제 절대 불가
	LegalHoldEnabled bool   // 법적 분쟁 시 영구 보존
}

// 베이스라인 정책
var (
	PolicyHIPAA = RetentionPolicy{
		Name: "hipaa", Years: HIPAARetentionYears, Region: "us",
		ImmutableDays: 30, LegalHoldEnabled: true,
	}
	PolicyKoreaMedical = RetentionPolicy{
		Name: "korea_medical", Years: KoreaMedicalRecordYears, Region: "kr",
		ImmutableDays: 30, LegalHoldEnabled: true,
	}
	PolicyManpasikDefault = RetentionPolicy{
		Name: "manpasik_default", Years: DefaultRetentionYears, Region: "global",
		ImmutableDays: 30, LegalHoldEnabled: true,
	}
)

// ============================================================================
// Immutable 보장
// ============================================================================

// ImmutableGuard는 작성 후 N일 동안 수정/삭제를 차단합니다.
type ImmutableGuard struct {
	policy RetentionPolicy
}

// NewImmutableGuard는 새 가드를 생성합니다.
func NewImmutableGuard(policy RetentionPolicy) *ImmutableGuard {
	return &ImmutableGuard{policy: policy}
}

// CanModify는 entry가 immutable 기간을 벗어났는지 확인합니다.
func (g *ImmutableGuard) CanModify(createdAt time.Time) bool {
	if g.policy.ImmutableDays <= 0 {
		return true
	}
	cutoff := createdAt.Add(time.Duration(g.policy.ImmutableDays) * 24 * time.Hour)
	return time.Now().UTC().After(cutoff)
}

// CanDelete는 entry가 보존 기간을 초과했는지 확인합니다.
//
// LegalHold가 활성화된 경우 삭제 절대 불가.
func (g *ImmutableGuard) CanDelete(createdAt time.Time, legalHold bool) bool {
	if legalHold && g.policy.LegalHoldEnabled {
		return false
	}
	cutoff := createdAt.AddDate(g.policy.Years, 0, 0)
	return time.Now().UTC().After(cutoff)
}

// RetentionExpiresAt은 보존 만료 시각을 반환합니다.
func (g *ImmutableGuard) RetentionExpiresAt(createdAt time.Time) time.Time {
	return createdAt.AddDate(g.policy.Years, 0, 0)
}

// ============================================================================
// LOINC 기반 검색
// ============================================================================

// SearchableEntry는 LOINC 코드 검색이 가능한 감사 엔트리 인터페이스입니다.
type SearchableEntry interface {
	GetID() string
	GetCreatedAt() time.Time
	GetLOINCCodes() []string
	GetResourceType() string
}

// SearchQuery는 감사 로그 검색 질의입니다.
type SearchQuery struct {
	UserID        string
	ResourceType  string   // "Observation", "DiagnosticReport" 등
	LOINCCodes    []string // 검색할 LOINC 코드
	StartDate     time.Time
	EndDate       time.Time
	Action        string   // "create" | "read" | "update" | "delete"
	IncludeExpired bool    // 보존 만료된 항목 포함 여부
	Limit         int
	Offset        int
}

// SearchResult는 검색 결과입니다.
type SearchResult struct {
	Entries      []SearchableEntry
	TotalCount   int
	Query        *SearchQuery
	ExecutedAt   time.Time
	DurationMs   int64
}

// SearchEngine은 LOINC 기반 검색 엔진입니다.
type SearchEngine struct {
	mu      sync.RWMutex
	indexed map[string][]SearchableEntry // LOINC code → entries
}

// NewSearchEngine은 새 검색 엔진을 생성합니다.
func NewSearchEngine() *SearchEngine {
	return &SearchEngine{indexed: make(map[string][]SearchableEntry)}
}

// Index는 엔트리를 인덱싱합니다.
func (e *SearchEngine) Index(entry SearchableEntry) error {
	if entry == nil {
		return errors.New("entry is nil")
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, code := range entry.GetLOINCCodes() {
		e.indexed[code] = append(e.indexed[code], entry)
	}
	return nil
}

// Search는 질의에 맞는 엔트리를 반환합니다.
func (e *SearchEngine) Search(query *SearchQuery) *SearchResult {
	if query == nil {
		query = &SearchQuery{}
	}
	startTime := time.Now()

	e.mu.RLock()
	defer e.mu.RUnlock()

	candidates := e.gatherCandidates(query)
	filtered := e.filterCandidates(candidates, query)
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].GetCreatedAt().After(filtered[j].GetCreatedAt())
	})

	totalCount := len(filtered)

	// 페이지네이션
	if query.Offset > 0 && query.Offset < len(filtered) {
		filtered = filtered[query.Offset:]
	}
	if query.Limit > 0 && query.Limit < len(filtered) {
		filtered = filtered[:query.Limit]
	}

	return &SearchResult{
		Entries:    filtered,
		TotalCount: totalCount,
		Query:      query,
		ExecutedAt: time.Now().UTC(),
		DurationMs: time.Since(startTime).Milliseconds(),
	}
}

func (e *SearchEngine) gatherCandidates(query *SearchQuery) []SearchableEntry {
	if len(query.LOINCCodes) == 0 {
		// LOINC 미지정 → 전체 인덱스에서 dedup
		seen := make(map[string]bool)
		out := []SearchableEntry{}
		for _, entries := range e.indexed {
			for _, entry := range entries {
				if !seen[entry.GetID()] {
					seen[entry.GetID()] = true
					out = append(out, entry)
				}
			}
		}
		return out
	}

	// LOINC 코드별 결과를 union (OR)
	seen := make(map[string]bool)
	out := []SearchableEntry{}
	for _, code := range query.LOINCCodes {
		for _, entry := range e.indexed[code] {
			if !seen[entry.GetID()] {
				seen[entry.GetID()] = true
				out = append(out, entry)
			}
		}
	}
	return out
}

func (e *SearchEngine) filterCandidates(candidates []SearchableEntry, query *SearchQuery) []SearchableEntry {
	out := make([]SearchableEntry, 0, len(candidates))
	for _, entry := range candidates {
		if query.ResourceType != "" && entry.GetResourceType() != query.ResourceType {
			continue
		}
		createdAt := entry.GetCreatedAt()
		if !query.StartDate.IsZero() && createdAt.Before(query.StartDate) {
			continue
		}
		if !query.EndDate.IsZero() && createdAt.After(query.EndDate) {
			continue
		}
		out = append(out, entry)
	}
	return out
}

// IndexedCount는 인덱싱된 엔트리 총 수를 반환합니다 (코드별 dedup).
func (e *SearchEngine) IndexedCount() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	seen := make(map[string]bool)
	for _, entries := range e.indexed {
		for _, entry := range entries {
			seen[entry.GetID()] = true
		}
	}
	return len(seen)
}

// ============================================================================
// 보존 자동 정리 (cron)
// ============================================================================

// ExpirationReport는 자동 정리 결과입니다.
type ExpirationReport struct {
	ScannedCount   int
	ExpiredCount   int
	DeletedCount   int
	LegalHoldCount int  // 법적 보존으로 삭제 보류
	ErrorCount     int
	StartedAt      time.Time
	CompletedAt    time.Time
}

// AuditExpirationService는 만료 엔트리를 자동 정리합니다.
type AuditExpirationService struct {
	guard *ImmutableGuard
}

// NewAuditExpirationService는 새 만료 정리 서비스를 생성합니다.
func NewAuditExpirationService(policy RetentionPolicy) *AuditExpirationService {
	return &AuditExpirationService{guard: NewImmutableGuard(policy)}
}

// CheckRetention은 단일 엔트리의 보존 상태를 평가합니다.
//
// 반환:
//   - canDelete: 보존 기간 초과 + legal hold 미설정
//   - daysUntilExpiry: 만료까지 남은 일수 (음수면 이미 만료)
func (s *AuditExpirationService) CheckRetention(createdAt time.Time, legalHold bool) (canDelete bool, daysUntilExpiry int) {
	canDelete = s.guard.CanDelete(createdAt, legalHold)
	expiresAt := s.guard.RetentionExpiresAt(createdAt)
	daysUntilExpiry = int(time.Until(expiresAt).Hours() / 24)
	return
}

// EvaluateBatch는 배치의 만료 상태를 평가합니다 (실 삭제는 호출자 책임).
type RetentionCandidate struct {
	EntryID         string
	CreatedAt       time.Time
	LegalHold       bool
	CanDelete       bool
	DaysUntilExpiry int
	Reason          string
}

func (s *AuditExpirationService) EvaluateBatch(entries []SearchableEntry, legalHoldByID map[string]bool) []*RetentionCandidate {
	candidates := make([]*RetentionCandidate, 0, len(entries))
	for _, entry := range entries {
		legalHold := legalHoldByID[entry.GetID()]
		canDelete, days := s.CheckRetention(entry.GetCreatedAt(), legalHold)
		c := &RetentionCandidate{
			EntryID:         entry.GetID(),
			CreatedAt:       entry.GetCreatedAt(),
			LegalHold:       legalHold,
			CanDelete:       canDelete,
			DaysUntilExpiry: days,
		}
		switch {
		case legalHold && canDelete:
			c.CanDelete = false
			c.Reason = "legal hold prevents deletion"
		case canDelete:
			c.Reason = fmt.Sprintf("retention expired (%d days ago)", -days)
		case days <= 30:
			c.Reason = fmt.Sprintf("expiring soon (%d days)", days)
		default:
			c.Reason = "within retention"
		}
		candidates = append(candidates, c)
	}
	return candidates
}

// ============================================================================
// 보존 정책 평가 (다중 정책 지원)
// ============================================================================

// PolicySelector는 데이터 종류/지역에 맞는 정책을 선택합니다.
type PolicySelector struct {
	policies map[string]RetentionPolicy
	defaultPolicy RetentionPolicy
}

// NewPolicySelector는 새 정책 선택기를 생성합니다.
func NewPolicySelector(defaultPolicy RetentionPolicy) *PolicySelector {
	return &PolicySelector{
		policies:      make(map[string]RetentionPolicy),
		defaultPolicy: defaultPolicy,
	}
}

// SetPolicy는 특정 영역/데이터 종류의 정책을 등록합니다.
//
// key 예: "kr:medical" / "us:hipaa" / "global:audit"
func (s *PolicySelector) SetPolicy(key string, policy RetentionPolicy) {
	s.policies[key] = policy
}

// Select는 환자 지역/데이터 종류로 정책을 선택합니다.
func (s *PolicySelector) Select(region, dataType string) RetentionPolicy {
	key := fmt.Sprintf("%s:%s", strings.ToLower(region), strings.ToLower(dataType))
	if p, ok := s.policies[key]; ok {
		return p
	}
	regionKey := strings.ToLower(region) + ":default"
	if p, ok := s.policies[regionKey]; ok {
		return p
	}
	return s.defaultPolicy
}

// SelectMostRestrictive는 적용 가능한 정책 중 가장 보수적인 것을 선택합니다.
//
// 다국적 데이터(예: 한국 환자가 미국 병원 이용)에 사용.
func (s *PolicySelector) SelectMostRestrictive(_ context.Context, regions []string, dataType string) RetentionPolicy {
	policies := make([]RetentionPolicy, 0, len(regions))
	for _, r := range regions {
		policies = append(policies, s.Select(r, dataType))
	}
	if len(policies) == 0 {
		return s.defaultPolicy
	}

	mostStrict := policies[0]
	for _, p := range policies[1:] {
		if p.Years > mostStrict.Years {
			mostStrict = p
		}
		// LegalHold나 ImmutableDays도 더 큰 쪽 선택
		if p.ImmutableDays > mostStrict.ImmutableDays {
			mostStrict.ImmutableDays = p.ImmutableDays
		}
		if p.LegalHoldEnabled && !mostStrict.LegalHoldEnabled {
			mostStrict.LegalHoldEnabled = true
		}
	}
	return mostStrict
}
