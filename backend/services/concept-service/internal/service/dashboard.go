// Package service: 컨셉/조직 통계 대시보드 + 멤버 초대 워크플로우
//
// 활동 트렌드, KPI 집계, 권한 검증, 초대 발송/수락/거절을 제공합니다.
package service

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

// ============================================================================
// 컨셉 통계
// ============================================================================

// ConceptStats는 컨셉별 통계입니다.
type ConceptStats struct {
	ConceptID         string
	Name              string
	Category          string
	UserCount         int
	DeviceCount       int
	MeasurementCount  int64
	AbnormalCount     int64
	AvgUsagePerDay    float64
	LastActivityAt    *time.Time
	GeneratedAt       time.Time
}

// OrganizationKPI는 조직 KPI입니다.
type OrganizationKPI struct {
	OrganizationID    string
	Name              string
	ActiveUsers       int
	TotalConcepts     int
	MeasurementsToday int64
	MeasurementsWeek  int64
	AbnormalRate      float64 // 0~1
	GeneratedAt       time.Time
}

// ActivityTrend는 시계열 활동 트렌드입니다.
type ActivityTrend struct {
	Date             time.Time
	ActiveUsers      int
	MeasurementCount int64
	AbnormalCount    int64
}

// ============================================================================
// 통계 엔진
// ============================================================================

// StatsEngine은 컨셉/조직 통계를 집계합니다.
type StatsEngine struct {
	mu                sync.RWMutex
	measurementCounts map[string]int64 // conceptID -> count
	abnormalCounts    map[string]int64
	dailyActivity     map[string]map[string]*ActivityTrend // conceptID -> date(YYYY-MM-DD) -> trend
	lastActivity      map[string]time.Time
}

// NewStatsEngine은 새 통계 엔진을 생성합니다.
func NewStatsEngine() *StatsEngine {
	return &StatsEngine{
		measurementCounts: make(map[string]int64),
		abnormalCounts:    make(map[string]int64),
		dailyActivity:     make(map[string]map[string]*ActivityTrend),
		lastActivity:      make(map[string]time.Time),
	}
}

// RecordMeasurement는 측정 활동을 기록합니다.
func (e *StatsEngine) RecordMeasurement(conceptID string, isAbnormal bool, ts time.Time) {
	if conceptID == "" {
		return
	}
	if ts.IsZero() {
		ts = time.Now().UTC()
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	e.measurementCounts[conceptID]++
	if isAbnormal {
		e.abnormalCounts[conceptID]++
	}

	dateKey := ts.Format("2006-01-02")
	if _, ok := e.dailyActivity[conceptID]; !ok {
		e.dailyActivity[conceptID] = make(map[string]*ActivityTrend)
	}
	if _, ok := e.dailyActivity[conceptID][dateKey]; !ok {
		e.dailyActivity[conceptID][dateKey] = &ActivityTrend{Date: ts.Truncate(24 * time.Hour)}
	}
	trend := e.dailyActivity[conceptID][dateKey]
	trend.MeasurementCount++
	if isAbnormal {
		trend.AbnormalCount++
	}

	e.lastActivity[conceptID] = ts
}

// GetConceptStats는 컨셉의 종합 통계를 반환합니다.
func (e *StatsEngine) GetConceptStats(c *Concept) *ConceptStats {
	if c == nil {
		return nil
	}

	e.mu.RLock()
	defer e.mu.RUnlock()

	stats := &ConceptStats{
		ConceptID:        c.ID,
		Name:             c.Name,
		Category:         c.Category,
		DeviceCount:      len(c.DeviceIDs),
		UserCount:        1, // owner
		MeasurementCount: e.measurementCounts[c.ID],
		AbnormalCount:    e.abnormalCounts[c.ID],
		GeneratedAt:      time.Now().UTC(),
	}

	if t, ok := e.lastActivity[c.ID]; ok {
		stats.LastActivityAt = &t
	}

	// 일평균 사용량 (활동 기간 기준)
	if dailyMap, ok := e.dailyActivity[c.ID]; ok && len(dailyMap) > 0 {
		stats.AvgUsagePerDay = float64(stats.MeasurementCount) / float64(len(dailyMap))
	}

	return stats
}

// GetTrends는 컨셉의 일별 활동 트렌드를 반환합니다 (최근 N일).
func (e *StatsEngine) GetTrends(conceptID string, days int) []*ActivityTrend {
	e.mu.RLock()
	defer e.mu.RUnlock()

	dailyMap, ok := e.dailyActivity[conceptID]
	if !ok {
		return nil
	}

	trends := make([]*ActivityTrend, 0, len(dailyMap))
	for _, t := range dailyMap {
		trends = append(trends, t)
	}

	sort.Slice(trends, func(i, j int) bool {
		return trends[i].Date.After(trends[j].Date)
	})

	if days > 0 && days < len(trends) {
		return trends[:days]
	}
	return trends
}

// GetOrgKPI는 조직 KPI를 계산합니다.
func (e *StatsEngine) GetOrgKPI(org *Organization, conceptCount int) *OrganizationKPI {
	if org == nil {
		return nil
	}

	e.mu.RLock()
	defer e.mu.RUnlock()

	kpi := &OrganizationKPI{
		OrganizationID: org.ID,
		Name:           org.Name,
		ActiveUsers:    len(org.MemberIDs),
		TotalConcepts:  conceptCount,
		GeneratedAt:    time.Now().UTC(),
	}

	// 오늘/이번주 측정 수 합산 (조직 멤버 컨셉)
	todayKey := time.Now().UTC().Format("2006-01-02")
	weekStart := time.Now().UTC().AddDate(0, 0, -7)

	var totalMeasure, totalAbnormal int64
	for _, dailyMap := range e.dailyActivity {
		for dateKey, trend := range dailyMap {
			if dateKey == todayKey {
				kpi.MeasurementsToday += trend.MeasurementCount
			}
			if !trend.Date.Before(weekStart) {
				kpi.MeasurementsWeek += trend.MeasurementCount
				totalMeasure += trend.MeasurementCount
				totalAbnormal += trend.AbnormalCount
			}
		}
	}

	if totalMeasure > 0 {
		kpi.AbnormalRate = float64(totalAbnormal) / float64(totalMeasure)
	}

	return kpi
}

// ============================================================================
// 역할 기반 권한 검증
// ============================================================================

// RoleType은 컨셉 멤버 역할입니다.
const (
	RoleOwner     = "owner"
	RoleAdmin     = "admin"
	RoleModerator = "moderator"
	RoleMember    = "member"
	RoleGuest     = "guest"
)

// rolePermissions는 역할별 허용 권한입니다.
var rolePermissions = map[string][]string{
	RoleOwner:     {"read", "write", "delete", "manage", "invite"},
	RoleAdmin:     {"read", "write", "delete", "invite"},
	RoleModerator: {"read", "write"},
	RoleMember:    {"read", "write"},
	RoleGuest:     {"read"},
}

// HasPermission은 역할이 특정 권한을 가지는지 확인합니다.
func HasPermission(role, action string) bool {
	perms, ok := rolePermissions[role]
	if !ok {
		return false
	}
	for _, p := range perms {
		if p == action {
			return true
		}
	}
	return false
}

// IsHigherOrEqualRole은 role1이 role2보다 권한이 높거나 같은지 확인합니다.
func IsHigherOrEqualRole(role1, role2 string) bool {
	rank := map[string]int{
		RoleOwner: 5, RoleAdmin: 4, RoleModerator: 3, RoleMember: 2, RoleGuest: 1,
	}
	return rank[role1] >= rank[role2]
}

// ============================================================================
// 초대 워크플로우
// ============================================================================

// InvitationStatus는 초대 상태입니다.
const (
	InviteStatusPending  = "pending"
	InviteStatusAccepted = "accepted"
	InviteStatusDeclined = "declined"
	InviteStatusExpired  = "expired"
	InviteStatusRevoked  = "revoked"
)

// Invitation은 컨셉/조직 초대입니다.
type Invitation struct {
	ID            string
	OrganizationID string
	ConceptID     string
	InviterID     string
	InviteeEmail  string
	InviteeUserID string
	Role          string
	Status        string
	Message       string
	CreatedAt     time.Time
	ExpiresAt     time.Time
	RespondedAt   *time.Time
}

// IsExpired는 만료 여부를 반환합니다.
func (i *Invitation) IsExpired() bool {
	return !i.ExpiresAt.IsZero() && time.Now().UTC().After(i.ExpiresAt)
}

// InvitationManager는 초대 워크플로우를 관리합니다.
type InvitationManager struct {
	mu          sync.RWMutex
	invitations map[string]*Invitation
}

// NewInvitationManager는 새 초대 매니저를 생성합니다.
func NewInvitationManager() *InvitationManager {
	return &InvitationManager{invitations: make(map[string]*Invitation)}
}

// Send는 새 초대를 생성합니다.
func (m *InvitationManager) Send(inv *Invitation) error {
	if inv == nil {
		return errors.New("invitation is nil")
	}
	if inv.ID == "" {
		inv.ID = fmt.Sprintf("inv-%d", time.Now().UnixNano())
	}
	if inv.InviteeEmail == "" && inv.InviteeUserID == "" {
		return errors.New("invitee_email or invitee_user_id required")
	}
	if inv.InviterID == "" {
		return errors.New("inviter_id required")
	}
	if inv.Role == "" {
		inv.Role = RoleMember
	}
	if _, ok := rolePermissions[inv.Role]; !ok {
		return fmt.Errorf("invalid role: %s", inv.Role)
	}

	now := time.Now().UTC()
	inv.Status = InviteStatusPending
	inv.CreatedAt = now
	if inv.ExpiresAt.IsZero() {
		inv.ExpiresAt = now.AddDate(0, 0, 7) // 7일 만료
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.invitations[inv.ID] = inv
	return nil
}

// Accept는 초대를 수락합니다.
func (m *InvitationManager) Accept(invID, userID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	inv, ok := m.invitations[invID]
	if !ok {
		return fmt.Errorf("invitation %s not found", invID)
	}
	if inv.Status != InviteStatusPending {
		return fmt.Errorf("invitation status is %s, cannot accept", inv.Status)
	}
	if inv.IsExpired() {
		inv.Status = InviteStatusExpired
		return errors.New("invitation expired")
	}

	now := time.Now().UTC()
	inv.Status = InviteStatusAccepted
	inv.RespondedAt = &now
	inv.InviteeUserID = userID
	return nil
}

// Decline은 초대를 거절합니다.
func (m *InvitationManager) Decline(invID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	inv, ok := m.invitations[invID]
	if !ok {
		return fmt.Errorf("invitation %s not found", invID)
	}
	if inv.Status != InviteStatusPending {
		return fmt.Errorf("invitation status is %s, cannot decline", inv.Status)
	}

	now := time.Now().UTC()
	inv.Status = InviteStatusDeclined
	inv.RespondedAt = &now
	return nil
}

// Revoke는 초대를 회수합니다 (발신자 측).
func (m *InvitationManager) Revoke(invID, requesterID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	inv, ok := m.invitations[invID]
	if !ok {
		return fmt.Errorf("invitation %s not found", invID)
	}
	if inv.InviterID != requesterID {
		return errors.New("only inviter can revoke")
	}
	if inv.Status != InviteStatusPending {
		return fmt.Errorf("cannot revoke in status %s", inv.Status)
	}
	inv.Status = InviteStatusRevoked
	return nil
}

// Get은 초대를 조회합니다.
func (m *InvitationManager) Get(invID string) (*Invitation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	inv, ok := m.invitations[invID]
	if !ok {
		return nil, fmt.Errorf("invitation %s not found", invID)
	}
	return inv, nil
}

// ListByOrganization은 조직의 모든 초대를 반환합니다.
func (m *InvitationManager) ListByOrganization(orgID string) []*Invitation {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []*Invitation
	for _, inv := range m.invitations {
		if inv.OrganizationID == orgID {
			result = append(result, inv)
		}
	}
	return result
}

// ExpireOld는 만료된 pending 초대를 expired 상태로 전환합니다.
func (m *InvitationManager) ExpireOld() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	expired := 0
	for _, inv := range m.invitations {
		if inv.Status == InviteStatusPending && inv.IsExpired() {
			inv.Status = InviteStatusExpired
			expired++
		}
	}
	return expired
}

// CountByStatus는 상태별 초대 건수를 반환합니다.
func (m *InvitationManager) CountByStatus() map[string]int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	counts := make(map[string]int)
	for _, inv := range m.invitations {
		counts[inv.Status]++
	}
	return counts
}
