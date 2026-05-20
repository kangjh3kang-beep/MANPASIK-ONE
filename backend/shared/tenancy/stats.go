package tenancy

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// TenancyStats 는 tenancy 모듈의 운영 메트릭 스냅샷.
//
// 운영 대시보드 (/ops/dashboard) + Prometheus exposition 에 노출.
type TenancyStats struct {
	// 멤버십 메트릭
	TenantCount         int            // 활성 조직 수 (멤버가 1명 이상)
	ActiveMemberCount   int            // 전체 활성 멤버
	InactiveMemberCount int            // 비활성 멤버
	RoleDistribution    map[string]int // 역할별 분포 (admin: 5, member: 30, ...)

	// 초대 메트릭
	PendingInvitationCount  int
	AcceptedInvitationCount int
	ExpiredInvitationCount  int
	RevokedInvitationCount  int

	// 메타
	GeneratedAt time.Time
}

// StatsCollector 는 tenancy 메트릭 수집기.
//
// MembershipStore + InvitationStore 를 스캔하여 통계 산출. 캐시 미사용
// (호출 시점 정확값). 캐시가 필요하면 호출자가 별도 wrapper 적용.
type StatsCollector struct {
	memStore     MembershipStore
	invStore     InvitationStore
	tenantLister TenantLister

	mu       sync.RWMutex
	lastSnap *TenancyStats // Collect 결과 캐시 (옵션)
}

// NewStatsCollector 생성.
//
// memStore 는 필수. invStore/tenantLister 는 nil 허용 — 해당 메트릭만 0 반환.
func NewStatsCollector(memStore MembershipStore, invStore InvitationStore,
	tenantLister TenantLister) *StatsCollector {
	return &StatsCollector{
		memStore:     memStore,
		invStore:     invStore,
		tenantLister: tenantLister,
	}
}

// Collect 는 현재 시점 메트릭을 산출.
func (s *StatsCollector) Collect() *TenancyStats {
	stats := &TenancyStats{
		RoleDistribution: make(map[string]int),
		GeneratedAt:      time.Now().UTC(),
	}

	if s.tenantLister != nil && s.memStore != nil {
		tenants := s.tenantLister.AllTenants()
		s.aggregateMembership(stats, tenants)
		s.aggregateInvitations(stats, tenants)
	}

	s.mu.Lock()
	s.lastSnap = stats
	s.mu.Unlock()
	return stats
}

// LastSnapshot 은 마지막 Collect 결과 (없으면 nil).
func (s *StatsCollector) LastSnapshot() *TenancyStats {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.lastSnap == nil {
		return nil
	}
	cp := *s.lastSnap
	if s.lastSnap.RoleDistribution != nil {
		cp.RoleDistribution = make(map[string]int, len(s.lastSnap.RoleDistribution))
		for k, v := range s.lastSnap.RoleDistribution {
			cp.RoleDistribution[k] = v
		}
	}
	return &cp
}

// aggregateMembership 은 모든 조직을 순회하여 멤버 통계 산출.
func (s *StatsCollector) aggregateMembership(stats *TenancyStats, tenants []TenantID) {
	tenantsWithMembers := 0
	for _, tid := range tenants {
		members := s.memStore.ListTenantMembers(tid)
		if len(members) > 0 {
			tenantsWithMembers++
		}
		for _, m := range members {
			if m.Active {
				stats.ActiveMemberCount++
				stats.RoleDistribution[string(m.Role)]++
			} else {
				stats.InactiveMemberCount++
			}
		}
	}
	stats.TenantCount = tenantsWithMembers
}

// aggregateInvitations 는 모든 조직을 순회하여 초대 상태별 카운트.
func (s *StatsCollector) aggregateInvitations(stats *TenancyStats, tenants []TenantID) {
	if s.invStore == nil {
		return
	}
	now := time.Now()
	for _, tid := range tenants {
		invs := s.invStore.ListByTenant(tid)
		for _, inv := range invs {
			switch inv.Status {
			case InvitationPending:
				// 만료 시각 지났으면 expired 로 카운트 (cleaner 가 아직 정리 안 한 상태)
				if inv.ExpiresAt.Before(now) {
					stats.ExpiredInvitationCount++
				} else {
					stats.PendingInvitationCount++
				}
			case InvitationAccepted:
				stats.AcceptedInvitationCount++
			case InvitationExpired:
				stats.ExpiredInvitationCount++
			case InvitationRevoked:
				stats.RevokedInvitationCount++
			}
		}
	}
}

// WritePrometheusMetrics 는 stats 를 Prometheus text exposition format 으로 출력.
//
// 메트릭 명세:
//
//	manpasik_tenancy_tenants_total                  gauge — 활성 조직 수
//	manpasik_tenancy_members_active_total           gauge — 활성 멤버
//	manpasik_tenancy_members_inactive_total         gauge — 비활성 멤버
//	manpasik_tenancy_role_distribution{role="..."}  gauge — 역할별
//	manpasik_tenancy_invitations_total{status="..."} gauge — 상태별
func WritePrometheusMetrics(w http.ResponseWriter, stats *TenancyStats) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	type metric struct {
		name  string
		help  string
		value int
	}

	scalars := []metric{
		{"manpasik_tenancy_tenants_total", "활성 조직 수 (멤버가 1명 이상)", stats.TenantCount},
		{"manpasik_tenancy_members_active_total", "활성 멤버 수", stats.ActiveMemberCount},
		{"manpasik_tenancy_members_inactive_total", "비활성 멤버 수", stats.InactiveMemberCount},
	}
	for _, m := range scalars {
		_, _ = w.Write([]byte("# HELP " + m.name + " " + m.help + "\n"))
		_, _ = w.Write([]byte("# TYPE " + m.name + " gauge\n"))
		_, _ = w.Write([]byte(m.name + " " + intStr(m.value) + "\n"))
	}

	if len(stats.RoleDistribution) > 0 {
		_, _ = w.Write([]byte("# HELP manpasik_tenancy_role_distribution 역할별 활성 멤버 분포\n"))
		_, _ = w.Write([]byte("# TYPE manpasik_tenancy_role_distribution gauge\n"))
		for role, count := range stats.RoleDistribution {
			_, _ = w.Write([]byte("manpasik_tenancy_role_distribution{role=\"" +
				escapePromLabel(role) + "\"} " + intStr(count) + "\n"))
		}
	}

	_, _ = w.Write([]byte("# HELP manpasik_tenancy_invitations_total 상태별 초대 카운트\n"))
	_, _ = w.Write([]byte("# TYPE manpasik_tenancy_invitations_total gauge\n"))
	for _, e := range []struct {
		status string
		count  int
	}{
		{"pending", stats.PendingInvitationCount},
		{"accepted", stats.AcceptedInvitationCount},
		{"expired", stats.ExpiredInvitationCount},
		{"revoked", stats.RevokedInvitationCount},
	} {
		_, _ = w.Write([]byte("manpasik_tenancy_invitations_total{status=\"" +
			e.status + "\"} " + intStr(e.count) + "\n"))
	}
}

// intStr 는 int → 10진수 문자열.
func intStr(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	digits := []byte{}
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	if neg {
		digits = append([]byte{'-'}, digits...)
	}
	return string(digits)
}

// escapePromLabel 은 Prometheus 라벨 escape (백슬래시/따옴표/줄바꿈).
func escapePromLabel(s string) string {
	out := make([]byte, 0, len(s))
	for _, c := range []byte(s) {
		switch c {
		case '\\':
			out = append(out, '\\', '\\')
		case '"':
			out = append(out, '\\', '"')
		case '\n':
			out = append(out, '\\', 'n')
		default:
			out = append(out, c)
		}
	}
	return string(out)
}

// WriteStatsJSON 은 stats 를 JSON 응답으로 직렬화 (HTTP 핸들러용 헬퍼).
func WriteStatsJSON(w http.ResponseWriter, stats *TenancyStats) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	body := map[string]interface{}{
		"tenant_count":          stats.TenantCount,
		"active_member_count":   stats.ActiveMemberCount,
		"inactive_member_count": stats.InactiveMemberCount,
		"role_distribution":     stats.RoleDistribution,
		"invitations": map[string]int{
			"pending":  stats.PendingInvitationCount,
			"accepted": stats.AcceptedInvitationCount,
			"expired":  stats.ExpiredInvitationCount,
			"revoked":  stats.RevokedInvitationCount,
		},
		"generated_at": stats.GeneratedAt.Format(time.RFC3339),
	}
	_ = json.NewEncoder(w).Encode(body)
}

// CollectMetrics 는 dashboard.MetricsCollector 인터페이스 구현.
//
// 인터페이스 호환을 위한 어댑터 — 모든 수치 메트릭을 단일 map 으로 노출.
func (s *StatsCollector) CollectMetrics() map[string]interface{} {
	stats := s.Collect()
	out := map[string]interface{}{
		"status":                    "healthy",
		"tenant_count":              stats.TenantCount,
		"active_member_count":       stats.ActiveMemberCount,
		"inactive_member_count":     stats.InactiveMemberCount,
		"pending_invitations":       stats.PendingInvitationCount,
		"accepted_invitations":      stats.AcceptedInvitationCount,
		"expired_invitations":       stats.ExpiredInvitationCount,
		"revoked_invitations":       stats.RevokedInvitationCount,
	}
	for role, count := range stats.RoleDistribution {
		out["role_"+role] = count
	}
	return out
}
