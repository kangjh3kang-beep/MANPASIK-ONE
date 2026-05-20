package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/manpasik/backend/services/audit-service/internal/service"
)

// ============================================================================
// 테스트 헬퍼: SearchableEntry 구현
// ============================================================================

type mockEntry struct {
	id           string
	createdAt    time.Time
	loincCodes   []string
	resourceType string
}

func (m *mockEntry) GetID() string             { return m.id }
func (m *mockEntry) GetCreatedAt() time.Time   { return m.createdAt }
func (m *mockEntry) GetLOINCCodes() []string   { return m.loincCodes }
func (m *mockEntry) GetResourceType() string   { return m.resourceType }

// ============================================================================
// 보존 정책 테스트
// ============================================================================

func TestRetentionPolicy_HIPAA(t *testing.T) {
	if service.PolicyHIPAA.Years != 6 {
		t.Errorf("HIPAA Years = %d, want 6", service.PolicyHIPAA.Years)
	}
	if !service.PolicyHIPAA.LegalHoldEnabled {
		t.Error("HIPAA LegalHold 미활성")
	}
}

func TestRetentionPolicy_KoreaMedical(t *testing.T) {
	if service.PolicyKoreaMedical.Years != 10 {
		t.Errorf("Korea Years = %d, want 10", service.PolicyKoreaMedical.Years)
	}
}

// ============================================================================
// ImmutableGuard 테스트
// ============================================================================

func TestImmutableGuard_CannotModifyFresh(t *testing.T) {
	g := service.NewImmutableGuard(service.PolicyHIPAA)
	createdAt := time.Now().UTC().Add(-10 * 24 * time.Hour) // 10일 전 (immutable 30일 미만)

	if g.CanModify(createdAt) {
		t.Error("immutable 기간(30일) 내 수정 허용됨")
	}
}

func TestImmutableGuard_CanModifyOld(t *testing.T) {
	g := service.NewImmutableGuard(service.PolicyHIPAA)
	createdAt := time.Now().UTC().Add(-60 * 24 * time.Hour) // 60일 전

	if !g.CanModify(createdAt) {
		t.Error("60일 전 항목 수정 불가")
	}
}

func TestImmutableGuard_CanDelete_Expired(t *testing.T) {
	g := service.NewImmutableGuard(service.PolicyHIPAA)
	createdAt := time.Now().UTC().AddDate(-7, 0, 0) // 7년 전 (보존 6년 초과)

	if !g.CanDelete(createdAt, false) {
		t.Error("보존 만료된 항목 삭제 불가")
	}
}

func TestImmutableGuard_CannotDelete_NotYetExpired(t *testing.T) {
	g := service.NewImmutableGuard(service.PolicyHIPAA)
	createdAt := time.Now().UTC().AddDate(-2, 0, 0) // 2년 전

	if g.CanDelete(createdAt, false) {
		t.Error("보존 기간 내인데 삭제 허용됨")
	}
}

func TestImmutableGuard_LegalHold_BlocksDelete(t *testing.T) {
	g := service.NewImmutableGuard(service.PolicyHIPAA)
	createdAt := time.Now().UTC().AddDate(-10, 0, 0) // 만료됐어도

	if g.CanDelete(createdAt, true) {
		t.Error("legal hold 활성인데 삭제 허용됨")
	}
}

func TestImmutableGuard_RetentionExpiresAt(t *testing.T) {
	g := service.NewImmutableGuard(service.PolicyHIPAA)
	createdAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	expires := g.RetentionExpiresAt(createdAt)

	if expires.Year() != 2032 {
		t.Errorf("ExpiresAt Year = %d, want 2032 (HIPAA 6년)", expires.Year())
	}
}

// ============================================================================
// SearchEngine 테스트
// ============================================================================

func TestSearchEngine_IndexAndSearch(t *testing.T) {
	e := service.NewSearchEngine()
	now := time.Now().UTC()

	entries := []*mockEntry{
		{id: "e1", createdAt: now, loincCodes: []string{"2345-7"}, resourceType: "Observation"},
		{id: "e2", createdAt: now, loincCodes: []string{"2345-7", "4548-4"}, resourceType: "Observation"},
		{id: "e3", createdAt: now, loincCodes: []string{"4548-4"}, resourceType: "DiagnosticReport"},
	}
	for _, e := range entries {
		_ = e // mockEntry 사용
	}

	for _, entry := range entries {
		if err := e.Index(entry); err != nil {
			t.Fatalf("Index 실패: %v", err)
		}
	}

	// Glucose (2345-7) 검색
	result := e.Search(&service.SearchQuery{LOINCCodes: []string{"2345-7"}})
	if result.TotalCount != 2 {
		t.Errorf("glucose 검색 = %d, want 2", result.TotalCount)
	}

	// HbA1c (4548-4) 검색
	result = e.Search(&service.SearchQuery{LOINCCodes: []string{"4548-4"}})
	if result.TotalCount != 2 {
		t.Errorf("hba1c 검색 = %d, want 2", result.TotalCount)
	}

	// ResourceType 필터
	result = e.Search(&service.SearchQuery{ResourceType: "DiagnosticReport"})
	if result.TotalCount != 1 {
		t.Errorf("DiagnosticReport = %d, want 1", result.TotalCount)
	}
}

func TestSearchEngine_DateRange(t *testing.T) {
	e := service.NewSearchEngine()
	old := time.Now().UTC().AddDate(0, -3, 0)
	recent := time.Now().UTC()

	_ = e.Index(&mockEntry{id: "old", createdAt: old, loincCodes: []string{"X"}, resourceType: "Observation"})
	_ = e.Index(&mockEntry{id: "recent", createdAt: recent, loincCodes: []string{"X"}, resourceType: "Observation"})

	// 최근 1개월만
	since := time.Now().UTC().AddDate(0, -1, 0)
	result := e.Search(&service.SearchQuery{StartDate: since})
	if result.TotalCount != 1 {
		t.Errorf("recent only = %d, want 1", result.TotalCount)
	}
}

func TestSearchEngine_Pagination(t *testing.T) {
	e := service.NewSearchEngine()
	now := time.Now().UTC()

	for i := 0; i < 10; i++ {
		_ = e.Index(&mockEntry{
			id: string(rune('a' + i)), createdAt: now,
			loincCodes: []string{"P"}, resourceType: "Observation",
		})
	}

	result := e.Search(&service.SearchQuery{LOINCCodes: []string{"P"}, Limit: 3, Offset: 5})
	if len(result.Entries) != 3 {
		t.Errorf("Page Entries = %d, want 3", len(result.Entries))
	}
	if result.TotalCount != 10 {
		t.Errorf("TotalCount = %d, want 10", result.TotalCount)
	}
}

func TestSearchEngine_NilEntry(t *testing.T) {
	e := service.NewSearchEngine()
	if err := e.Index(nil); err == nil {
		t.Error("nil entry 통과")
	}
}

func TestSearchEngine_AllInOneRequest(t *testing.T) {
	e := service.NewSearchEngine()
	now := time.Now().UTC()

	_ = e.Index(&mockEntry{id: "a", createdAt: now, loincCodes: []string{"X", "Y"}, resourceType: "Observation"})

	// LOINC 미지정 + 전체 인덱스에서 dedup
	result := e.Search(&service.SearchQuery{})
	if result.TotalCount != 1 {
		t.Errorf("AllInOne TotalCount = %d, want 1 (a 하나만)", result.TotalCount)
	}
}

// ============================================================================
// AuditExpirationService 테스트
// ============================================================================

func TestExpirationService_CheckRetention(t *testing.T) {
	s := service.NewAuditExpirationService(service.PolicyHIPAA)

	// 7년 전 (만료)
	old := time.Now().UTC().AddDate(-7, 0, 0)
	canDelete, days := s.CheckRetention(old, false)
	if !canDelete {
		t.Error("7년 전 항목 삭제 불가")
	}
	if days >= 0 {
		t.Errorf("DaysUntilExpiry = %d, want < 0 (이미 만료)", days)
	}

	// 1년 전 (보존 중)
	recent := time.Now().UTC().AddDate(-1, 0, 0)
	canDelete2, days2 := s.CheckRetention(recent, false)
	if canDelete2 {
		t.Error("1년 전 항목 삭제 허용")
	}
	if days2 < 0 {
		t.Errorf("DaysUntilExpiry = %d, want > 0", days2)
	}
}

func TestExpirationService_LegalHold(t *testing.T) {
	s := service.NewAuditExpirationService(service.PolicyHIPAA)
	old := time.Now().UTC().AddDate(-10, 0, 0)

	canDelete, _ := s.CheckRetention(old, true)
	if canDelete {
		t.Error("legal hold 활성인데 삭제 허용")
	}
}

func TestExpirationService_EvaluateBatch(t *testing.T) {
	s := service.NewAuditExpirationService(service.PolicyHIPAA)
	now := time.Now().UTC()

	entries := []service.SearchableEntry{
		&mockEntry{id: "expired", createdAt: now.AddDate(-7, 0, 0)},
		&mockEntry{id: "fresh", createdAt: now.AddDate(-1, 0, 0)},
		&mockEntry{id: "expiring", createdAt: now.AddDate(-6, 0, -10)}, // 6년 + 10일 전 (만료 직전)
	}
	legalHolds := map[string]bool{
		"expired": true, // legal hold로 보호
	}

	candidates := s.EvaluateBatch(entries, legalHolds)
	if len(candidates) != 3 {
		t.Errorf("Candidates = %d, want 3", len(candidates))
	}

	for _, c := range candidates {
		switch c.EntryID {
		case "expired":
			if c.CanDelete {
				t.Error("expired에 legal hold인데 삭제 허용")
			}
		case "fresh":
			if c.CanDelete {
				t.Error("fresh 삭제 허용")
			}
		case "expiring":
			// 이미 만료 (6년 10일 전 → 2주 만료) 또는 만료 직전
			if c.DaysUntilExpiry > 30 {
				t.Errorf("expiring days = %d, want <= 30", c.DaysUntilExpiry)
			}
		}
	}
}

// ============================================================================
// PolicySelector 테스트
// ============================================================================

func TestPolicySelector_DefaultFallback(t *testing.T) {
	s := service.NewPolicySelector(service.PolicyManpasikDefault)
	policy := s.Select("unknown", "unknown")
	if policy.Name != "manpasik_default" {
		t.Errorf("Default policy = %q", policy.Name)
	}
}

func TestPolicySelector_RegionMatch(t *testing.T) {
	s := service.NewPolicySelector(service.PolicyManpasikDefault)
	s.SetPolicy("kr:medical", service.PolicyKoreaMedical)
	s.SetPolicy("us:hipaa", service.PolicyHIPAA)

	kr := s.Select("kr", "medical")
	if kr.Years != 10 {
		t.Errorf("KR medical = %d, want 10", kr.Years)
	}

	us := s.Select("us", "hipaa")
	if us.Years != 6 {
		t.Errorf("US HIPAA = %d, want 6", us.Years)
	}
}

func TestPolicySelector_MostRestrictive(t *testing.T) {
	s := service.NewPolicySelector(service.PolicyManpasikDefault)
	s.SetPolicy("kr:medical", service.PolicyKoreaMedical)
	s.SetPolicy("us:hipaa", service.PolicyHIPAA)

	// 한국 + 미국 동시 적용 → 가장 보수적인 한국 (10년) 선택
	policy := s.SelectMostRestrictive(context.Background(), []string{"kr", "us"}, "medical")
	if policy.Years < 10 {
		t.Errorf("MostRestrictive Years = %d, want >= 10", policy.Years)
	}
}

func TestSearchEngine_IndexedCount(t *testing.T) {
	e := service.NewSearchEngine()
	now := time.Now().UTC()

	_ = e.Index(&mockEntry{id: "a", createdAt: now, loincCodes: []string{"X", "Y"}})
	_ = e.Index(&mockEntry{id: "b", createdAt: now, loincCodes: []string{"Z"}})

	if e.IndexedCount() != 2 {
		t.Errorf("IndexedCount = %d, want 2", e.IndexedCount())
	}
}

func TestImmutableGuard_ZeroImmutableDays(t *testing.T) {
	policy := service.RetentionPolicy{Years: 6, ImmutableDays: 0}
	g := service.NewImmutableGuard(policy)

	// ImmutableDays=0이면 항상 수정 가능
	now := time.Now().UTC()
	if !g.CanModify(now) {
		t.Error("ImmutableDays=0인데 수정 불가")
	}
}
